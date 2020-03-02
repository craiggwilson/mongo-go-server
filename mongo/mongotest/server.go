package mongotest

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/craiggwilson/mongo-go-server/mongo"
)

// NewServer makes a Server that is already started.
func NewServer(handler mongo.MessageHandler) *Server {
	s := NewUnstartedServer(handler)
	s.Start()
	return s
}

// NewUnstartedServer makes a Server that is not started.
func NewUnstartedServer(handler mongo.MessageHandler) *Server {
	return &Server{
		Listener: newLocalListener(),
		Config: &mongo.Server{
			Handler: handler,
		},
	}
}

type Server struct {
	Listener net.Listener
	Config   *mongo.Server
	HostPort string

	conns  map[net.Conn]mongo.ConnState
	mu     sync.Mutex
	wg     sync.WaitGroup
	closed bool
}

func (s *Server) Close() {
	s.mu.Lock()
	if !s.closed {
		s.closed = true
		_ = s.Listener.Close()
		for c, st := range s.conns {
			if st == mongo.StateNew || st == mongo.StateInactive {
				_ = c.Close()
			}
		}
	}
	s.mu.Unlock()

	s.wg.Wait()
}

func (s *Server) Dial() net.Conn {
	c, err := net.Dial("tcp", s.Listener.Addr().String())
	if err != nil {
		panic(err)
	}

	return c
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Config.Shutdown(ctx)
}

func (s *Server) Start() {
	if s.HostPort != "" {
		panic("Server already started")
	}

	s.wrap()
	s.goServe()
}

func (s *Server) goServe() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		_ = s.Config.Serve(context.Background(), s.Listener)
	}()
}

func (s *Server) wrap() {
	oldHook := s.Config.ConnStateHook
	s.Config.ConnStateHook = mongo.ConnStateHookFunc(func(ctx context.Context, c net.Conn, st mongo.ConnState, d time.Duration) {
		s.mu.Lock()
		defer s.mu.Unlock()
		switch st {
		case mongo.StateNew:
			s.wg.Add(1)
			if _, exists := s.conns[c]; exists {
				panic("invalid state transition")
			}
			if s.conns == nil {
				s.conns = make(map[net.Conn]mongo.ConnState)
			}

			s.conns[c] = st
			if s.closed {
				_ = c.Close()
			}
		case mongo.StateActive:
			if oldState, ok := s.conns[c]; ok {
				if oldState != mongo.StateNew && oldState != mongo.StateInactive {
					panic("invalid state transition")
				}

				s.conns[c] = st
			}
		case mongo.StateInactive:
			if oldState, ok := s.conns[c]; ok {
				if oldState != mongo.StateActive {
					panic("invalid state transition")
				}
				s.conns[c] = st
			}
			if s.closed {
				_ = c.Close()
			}
		case mongo.StateClosed:
			if _, ok := s.conns[c]; ok {
				delete(s.conns, c)
				s.wg.Done()
			}
		}

		if oldHook != nil {
			oldHook.OnConnStateChange(ctx, c, st, d)
		}
	})
}

type closeWatchConn struct {
	net.Conn

	closed bool
}

func (c *closeWatchConn) Close() error {
	c.closed = true
	return c.Conn.Close()
}

func (c *closeWatchConn) IsClosed() bool {
	return c.IsClosed()
}

func newLocalListener() net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if l, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			panic(fmt.Sprintf("mongotest: failed to listen on a port: %v", err))
		}
	}
	return l
}
