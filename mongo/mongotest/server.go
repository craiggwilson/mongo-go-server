package mongotest

import (
	"context"
	"fmt"
	"net"

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
}

func (s *Server) Close() error {
	return s.Config.Close()
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

	go func() {
		_ = s.Config.Serve(context.Background(), s.Listener)
	}()
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
