package mongo_test

import (
	"testing"

	"github.com/craiggwilson/mongo-go-server/mongo/mongotest"
)

func TestServer_ConnectAndClose(t *testing.T) {
	s := mongotest.NewServer(nil)

	c := s.Dial()

	if _, err := c.Write([]byte{1, 0, 0, 0, 0}); err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	s.Close()
}

// func TestServer_ConnectAndShutdown(t *testing.T) {
// 	ctx, shutdown := context.WithCancel(context.Background())
// 	hostPort, s := startTestServer(context.Background(), nil)

// 	c, err := net.Dial("tcp", hostPort)
// 	if err != nil {
// 		t.Fatalf("expected no error, but got %v", err)
// 	}

// 	watchedC := &closeWatchConn{Conn: c}

// 	if _, err = c.Write([]byte{1, 0, 0, 0, 0}); err != nil {
// 		t.Fatalf("expected no error, but got %v", err)
// 	}

// 	now := time.Now()
// 	done := make(chan struct{})

// 	go func() {
// 		shutdown()
// 		time.Sleep(3 * time.Second)
// 		close(done)
// 	}()

// 	<-done

// 	if watchedC.IsClosed() {
// 		t.Fatalf("expected shutdown to wait until connection was closed")
// 	}

// 	_ = s.Close()
// }
