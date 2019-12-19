package mongo

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

// ErrorLogger is a callback for logging errors that occur during serving.
type ErrorLogger interface {
	Logf(format string, args ...interface{})
}

// ErrAbortHandler is used to close a connection from within a handler.
var ErrAbortHandler = errors.New("abort")

// ErrServerClosed is returned when an operation on a closed server occurs.
var ErrServerClosed = errors.New("mongo: server closed")

// ListenAndServe starts up a server at the specified address with messages
// handled by the handler.
func ListenAndServe(ctx context.Context, addr string, handler MessageHandler) error {
	svr := &Server{
		Handler: handler,
	}

	return svr.ListenAndServe(ctx, addr)
}

func Serve(ctx context.Context, l net.Listener, handler MessageHandler) error {
	svr := &Server{
		Handler: handler,
	}

	return svr.Serve(ctx, l)
}

// Server serves the MongoDB wire protocol.
type Server struct {
	Handler     MessageHandler
	ErrorLogger ErrorLogger

	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	Compressors     []Compressor
	MaxDocumentSize int32
	MaxMessageSize  int32

	ConnectionDecorator ConnectionDecorator

	mu            sync.Mutex
	doneChan      chan struct{}
	listeners     map[*net.Listener]struct{}
	conns         map[*conn]struct{}
	inShutdown    int32
	currentConnID uint64
}

// Close immediately closes all active net.Listeners and any connections.
// For a graceful shutdown, use Shutdown.
func (s *Server) Close() error {
	atomic.StoreInt32(&s.inShutdown, 1)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.closeDoneChanLocked()
	err := s.closeListenersLocked()
	for c := range s.conns {
		c.rwc.Close()
		delete(s.conns, c)
	}

	return err
}

func (s *Server) ListenAndServe(ctx context.Context, addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return s.Serve(ctx, l)
}

func (s *Server) Serve(ctx context.Context, l net.Listener) error {
	if ctx == nil {
		ctx = context.Background()
	}

	l = &onceCloseListener{Listener: l}
	defer l.Close()

	if !s.trackListener(&l, true) {
		return ErrServerClosed
	}
	defer s.trackListener(&l, false)

	var tempDelay time.Duration // how long to sleep on accept failure

	ctx = context.WithValue(ctx, ServerContextKey, s)
	for {
		rw, err := l.Accept()
		if err != nil {
			select {
			case <-s.getDoneChan():
				return ErrServerClosed
			default:
			}

			var nerr net.Error
			if errors.As(err, &nerr) && nerr.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				// LOG
				time.Sleep(tempDelay)
				continue
			}

			return err
		}

		connCtx := ctx
		if s.ConnectionDecorator != nil {
			connCtx, rw = s.ConnectionDecorator.DecorateConnection(ctx, rw)
			if connCtx == nil || rw == nil {
				panic("ConnDecorator cannot return a nil context.Context or net.Conn")
			}
		}
		tempDelay = 0
		c := s.newConn(rw)
		s.trackConn(c, true)
		go c.serve(connCtx)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	atomic.StoreInt32(&s.inShutdown, 1)

	s.mu.Lock()
	err := s.closeListenersLocked()
	s.closeDoneChanLocked()
	s.mu.Unlock()

	// Poll all the opened connections until none
	// are left. On the way, close idle connections.
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		if s.closeIdleConns() {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (s *Server) closeDoneChanLocked() {
	ch := s.getDoneChanLocked()
	select {
	case <-ch:
		// Already closed. Don't close again.
	default:
		// Safe to close here. We're the only closer, guarded
		// by s.mu.
		close(ch)
	}
}

func (s *Server) closeIdleConns() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	allIdle := true
	for c := range s.conns {
		if !c.isIdle() {
			allIdle = false
			continue
		}

		_ = c.rwc.Close()
		delete(s.conns, c)
	}

	return allIdle
}

func (s *Server) closeListenersLocked() error {
	var err error
	for l := range s.listeners {
		if cerr := (*l).Close(); cerr != nil && err == nil {
			err = cerr
		}
		delete(s.listeners, l)
	}
	return err
}

func (s *Server) getDoneChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getDoneChanLocked()
}

func (s *Server) getDoneChanLocked() chan struct{} {
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	return s.doneChan
}

func (s *Server) handler() MessageHandler {
	h := s.Handler
	if h == nil {
		h = MessageHandlerFunc(func(ctx context.Context, resp MessageResponseWriter, req *MessageRequest) error {
			return newError(nil, 235, "no handlers configured")
		})
	}

	if len(s.Compressors) > 0 {
		h = &CompressionHandler{
			Compressors: s.Compressors,
			Handler:     h,
		}
	}

	return h
}

func (s *Server) logf(format string, args ...interface{}) {
	if s.ErrorLogger != nil {
		s.ErrorLogger.Logf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

func (s *Server) newConn(rwc net.Conn) *conn {
	return &conn{
		connID: atomic.AddUint64(&s.currentConnID, 1),
		server: s,
		rwc:    rwc,
	}
}

func (s *Server) trackConn(c *conn, add bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conns == nil {
		s.conns = make(map[*conn]struct{})
	}

	if add {
		s.conns[c] = struct{}{}
	} else {
		delete(s.conns, c)
	}
}

func (s *Server) trackListener(l *net.Listener, add bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listeners == nil {
		s.listeners = make(map[*net.Listener]struct{})
	}
	if add {
		if s.shuttingDown() {
			return false
		}
		s.listeners[l] = struct{}{}
	} else {
		delete(s.listeners, l)
	}

	return true
}

func (s *Server) shuttingDown() bool {
	return atomic.LoadInt32(&s.inShutdown) != 0
}

// onceCloseListener wraps a net.Listener, protecting it from
// multiple Close calls.
type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceCloseListener) close() { oc.closeErr = oc.Listener.Close() }

// ChainConnectionDecorators chains together multiple connection decorators.
func ChainConnectionDecorators(decorators ...ConnectionDecorator) ConnectionDecorator {
	return ConnectionDecoratorFunc(func(ctx context.Context, c net.Conn) (context.Context, net.Conn) {
		for _, dec := range decorators {
			ctx, c = dec.DecorateConnection(ctx, c)
		}

		return ctx, c
	})
}

// ConnectionDecorator decorates a connection and/or its context.
type ConnectionDecorator interface {
	DecorateConnection(context.Context, net.Conn) (context.Context, net.Conn)
}

// ConnectionDecoratorFunc is a function implementation of ConnectionDecorator.
type ConnectionDecoratorFunc func(context.Context, net.Conn) (context.Context, net.Conn)

// DecorateConnection implements the ConnectionDecorator interface.
func (f ConnectionDecoratorFunc) DecorateConnection(ctx context.Context, c net.Conn) (context.Context, net.Conn) {
	return f(ctx, c)
}

type conn struct {
	connID uint64
	server *Server
	rwc    net.Conn

	localAddr  string
	remoteAddr string

	r    *connReader
	bufr *bufio.Reader
	bufw *bufio.Writer
}

func (c *conn) finalFlush() {
	if c.bufr != nil {
		putBufioReader(c.bufr)
		c.bufr = nil
	}

	if c.bufw != nil {
		c.bufw.Flush()
		putBufioWriter(c.bufw)
		c.bufw = nil
	}
}

func (c *conn) isIdle() bool {
	// TODO: implement this...
	return false
}

func (c *conn) readMessage(ctx context.Context) (wiremessage.WireMessage, error) {
	var readDeadline time.Time
	if d := c.server.ReadTimeout; d != 0 {
		readDeadline = time.Now().Add(d)
	}
	_ = c.rwc.SetReadDeadline(readDeadline)

	return c.r.read(c.bufr)
}

func (c *conn) serve(ctx context.Context) {
	c.localAddr = c.rwc.LocalAddr().String()
	c.remoteAddr = c.rwc.RemoteAddr().String()
	defer func() {
		if err := recover(); err != nil && err != ErrAbortHandler {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			c.server.logf("mongo: panic serving %v: %v\n%s", c.remoteAddr, err, buf)
		}
		_ = c.rwc.Close()
		c.finalFlush()
		c.server.trackConn(c, false)
	}()

	maxMessageSize := c.server.MaxMessageSize
	if maxMessageSize == 0 {
		maxMessageSize = DefaultMaxMessageSize
	}

	c.r = &connReader{maxMessageSize: maxMessageSize}
	c.bufr = newBufioReader(c.rwc)
	c.bufw = newBufioWriterSize(c.rwc, 2<<10) // start with 2K

	handler := c.server.handler()

	for {
		msg, err := c.readMessage(ctx)
		if err != nil {
			switch {
			case err == io.EOF:
				return
			default:
				c.server.logf("mongo: error reading request: %v", err)
				return
			}
		}

		req := &MessageRequest{
			ConnectionID: c.connID,
			LocalAddr:    c.localAddr,
			RemoteAddr:   c.remoteAddr,
			Message:      msg,
		}
		resp := &connResponseWriter{bufw: c.bufw}

		if err = handler.HandleMessage(ctx, resp, req); err != nil {
			var cerr *Error
			if errors.As(err, &cerr) {
				err = WriteSingleDocument(resp, req, cerr.MarshalBSON())
			}

			if err != nil {
				c.server.logf("mongo: error handling request: %v", err)
			}
		}

		if err = resp.bufw.Flush(); err != nil {
			c.server.logf("mongo: error writing response: %v", err)
		}
	}
}

type connReader struct {
	tmp            [4]byte
	maxMessageSize int32
}

func (r *connReader) read(bufr *bufio.Reader) (wiremessage.WireMessage, error) {
	_, err := io.ReadFull(bufr, r.tmp[:])
	if err != nil {
		return nil, err
	}

	size := int32(binary.LittleEndian.Uint32(r.tmp[:]))

	if size > r.maxMessageSize || size < 16 {
		return nil, fmt.Errorf("message size is out of bounds: %v (max size: %v)", size, r.maxMessageSize)
	}

	msg := make([]byte, size)
	msg[0] = r.tmp[0]
	msg[1] = r.tmp[1]
	msg[2] = r.tmp[2]
	msg[3] = r.tmp[3]
	if _, err = io.ReadFull(bufr, msg[4:]); err != nil {
		return nil, err
	}

	return wiremessage.WireMessage(msg), nil
}

type connResponseWriter struct {
	bufw *bufio.Writer
	err  error
}

func (r *connResponseWriter) WriteMessage(msg wiremessage.WireMessage) error {
	if r.err != nil {
		return r.err
	}

	_, r.err = r.bufw.Write(msg)
	if r.err == nil {
		r.err = r.bufw.Flush()
	}
	return r.err
}
