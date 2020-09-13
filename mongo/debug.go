package mongo

import (
	"context"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// CommandDebuggingDecorator logs every command and its response.
type CommandDebuggingDecorator struct{}

// DecorateCommandHandler implements the CommandHandlerDecorator interface.
func (CommandDebuggingDecorator) DecorateCommandHandler(ctx context.Context, h CommandHandler) (context.Context, CommandHandler) {
	return ctx, CommandHandlerFunc(func(ctx context.Context, resp CommandResponseWriter, req *CommandRequest) error {
		log.Printf("received command: %s\n", req.Document)
		w := &debuggingCommandResponseWriter{CommandResponseWriter: resp}
		return h.HandleCommand(ctx, w, req)
	})
}

type debuggingCommandResponseWriter struct {
	CommandResponseWriter
}

func (w *debuggingCommandResponseWriter) WriteSingleDocument(doc bsoncore.Document) error {
	log.Printf("writing response document: %s\n", doc)
	return w.CommandResponseWriter.WriteSingleDocument(doc)
}

// ConnectionDebuggingDecorator logs when a connection is received and closed.
type ConnectionDebuggingDecorator struct{}

// DecorateConnection implements the ConnectionDecorator interface.
func (ConnectionDebuggingDecorator) DecorateConnection(ctx context.Context, c net.Conn) (context.Context, net.Conn) {
	log.Printf("received connection on %s from %s", c.LocalAddr(), c.RemoteAddr())
	return ctx, c
}

type debuggingConnection struct {
	net.Conn
}

func (c *debuggingConnection) Close() error {
	err := c.Conn.Close()
	if err != nil {
		log.Print("failed closing connection")
		return err
	}

	log.Print("closed connection")
	return nil
}
