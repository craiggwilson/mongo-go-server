package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

// MessageHandler represents a way to handle an incoming message.
type MessageHandler interface {
	HandleMessage(context.Context, MessageResponseWriter, *MessageRequest) error
}

// MessageHandlerFunc is a functional implementation of MessageHandler.
type MessageHandlerFunc func(context.Context, MessageResponseWriter, *MessageRequest) error

// HandleMessage implements the MessageHandler interface.
func (f MessageHandlerFunc) HandleMessage(ctx context.Context, resp MessageResponseWriter, req *MessageRequest) error {
	return f(ctx, resp, req)
}

// MessageRequest represents the request that came in as well as any pertinent information.
type MessageRequest struct {
	// ConnectionID is the identifier of the connection.
	ConnectionID uint64
	// LocalAddr is thue address of the receiver.
	LocalAddr string
	// RemoteAddr is the address of the sender.
	RemoteAddr string
	// Message is the wire message representing the request.
	Message wiremessage.WireMessage
}

// MessageResponseWriter is used to write a response.
type MessageResponseWriter interface {
	WriteMessage(msg wiremessage.WireMessage) error
}
