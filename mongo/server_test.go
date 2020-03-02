package mongo_test

import (
	"context"
	"testing"
	"time"

	"github.com/craiggwilson/mongo-go-server/mongo/mongotest"
)

func TestServer_ConnectAndClose(t *testing.T) {
	t.Run("new connection", func(t *testing.T) {
		s := mongotest.NewServer(nil)
		_ = s.Dial()

		done := make(chan struct{})

		go func() {
			time.Sleep(1 * time.Second)
			s.Close()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatalf("expected all connections to be closed, but some were not.")
		}
	})
	t.Run("active connection", func(t *testing.T) {
		s := mongotest.NewServer(nil)
		c := s.Dial()

		if _, err := c.Write([]byte{1, 0, 0, 0, 0}); err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}

		done := make(chan struct{})

		go func() {
			time.Sleep(1 * time.Second)
			s.Close()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatalf("expected all connections to be closed, but some were not.")
		}
	})
}

func TestServer_ConnectAndShutdown(t *testing.T) {
	t.Run("new connection", func(t *testing.T) {
		s := mongotest.NewServer(nil)
		c := s.Dial()

		done := make(chan struct{})

		go func() {
			time.Sleep(1 * time.Second)
			s.Shutdown(context.Background())
			close(done)
		}()

		select {
		case <-done:
			t.Fatalf("expected no connections to be closed, but all were")
		case <-time.After(1 * time.Second):
			_ = c.Close()
			select {
			case <-done:
			case <-time.After(1 * time.Second):
				t.Fatalf("expected all connections to be closed, but some were not.")
			}
		}
	})
	t.Run("active connection", func(t *testing.T) {
		s := mongotest.NewServer(nil)
		c := s.Dial()

		if _, err := c.Write([]byte{1, 0, 0, 0, 0}); err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}

		done := make(chan struct{})

		go func() {
			time.Sleep(1 * time.Second)
			s.Shutdown(context.Background())
			close(done)
		}()

		select {
		case <-done:
			t.Fatalf("expected no connections to be closed, but all were")
		case <-time.After(1 * time.Second):
			_ = c.Close()
			select {
			case <-done:
			case <-time.After(1 * time.Second):
				t.Fatalf("expected all connections to be closed, but some were not.")
			}
		}
	})
}
