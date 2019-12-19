package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// CommandDebuggingDecorator logs every command and its response.
type CommandDebuggingDecorator struct{}

// Decorate implements the Decorator interface.
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
