package mongo

import (
	"context"
	"hash/crc32"
	"strings"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

// CommandHandler executes a given command.
type CommandHandler interface {
	HandleCommand(context.Context, CommandResponseWriter, *CommandRequest) error
}

// CommandHandlerFunc is a function implementation of the CommandHandler interface.
type CommandHandlerFunc func(context.Context, CommandResponseWriter, *CommandRequest) error

// HandleCommand implements the CommandHandler interface.
func (f CommandHandlerFunc) HandleCommand(ctx context.Context, resp CommandResponseWriter, req *CommandRequest) error {
	return f(ctx, resp, req)
}

// CommandRequest holds information about the command to execute.
type CommandRequest struct {
	*MessageRequest

	DatabaseName string
	Document     bsoncore.Document
}

// CommandResponseWriter writes command responses.
type CommandResponseWriter interface {
	WriteSingleDocument(doc bsoncore.Document) error
}

type commandResponseWriter struct {
	MessageResponseWriter

	req *MessageRequest
}

func (w *commandResponseWriter) WriteSingleDocument(doc bsoncore.Document) error {
	return WriteSingleDocument(w.MessageResponseWriter, w.req, doc)
}

// ChainCommandHandlerDecorators chains together multiple CommandHandlerDecorators.
func ChainCommandHandlerDecorators(decorators ...CommandHandlerDecorator) CommandHandlerDecorator {
	return CommandHandlerDecoratorFunc(func(ctx context.Context, h CommandHandler) (context.Context, CommandHandler) {
		for _, dec := range decorators {
			ctx, h = dec.DecorateCommandHandler(ctx, h)
		}
		return ctx, h
	})
}

// CommandHandlerDecorator decorates a CommandHandler.
type CommandHandlerDecorator interface {
	DecorateCommandHandler(context.Context, CommandHandler) (context.Context, CommandHandler)
}

// CommandHandlerDecoratorFunc is a functional implementation of a CommandHandlerDecorator.
type CommandHandlerDecoratorFunc func(context.Context, CommandHandler) (context.Context, CommandHandler)

// Decorate implements the CommandDecorator interface.
func (f CommandHandlerDecoratorFunc) DecorateCommandHandler(ctx context.Context, h CommandHandler) (context.Context, CommandHandler) {
	return f(ctx, h)
}

// New makes an empty Mux.
func NewCommandMux() *CommandMux {
	return &CommandMux{
		Handlers: make(map[string]CommandHandler),
	}
}

// CommandMux is a MessageHandler implementation that works with commands instead of messages.
type CommandMux struct {
	Handlers map[string]CommandHandler
	Fallback CommandHandler

	CommandHandlerDecorator CommandHandlerDecorator
}

// HandleMessage implements the MessageHandler interface.
func (m *CommandMux) HandleMessage(ctx context.Context, resp MessageResponseWriter, req *MessageRequest) error {
	cmdReq, err := readMainDocument(resp, req)
	if err != nil {
		return err
	}

	// if we had no error, but also no request, then there is nothing further to do.
	if cmdReq == nil {
		return nil
	}

	if err = cmdReq.Document.Validate(); err != nil {
		return newError(err, CodeInvalidBSON, "invalid command document")
	}

	nameElement, err := cmdReq.Document.IndexErr(0)
	if err != nil {
		return newError(err, CodeProtocolError, "failed reading command name")
	}

	name := nameElement.Key()

	var handler CommandHandler
	if m.Handlers != nil {
		var ok bool
		if handler, ok = m.Handlers[name]; !ok {
			handler = m.Fallback
		}
	}

	if handler == nil {
		return newErrorf(nil, CodeCommandNotFound, "no such command: '%s'", name)
	}

	if m.CommandHandlerDecorator != nil {
		ctx, handler = m.CommandHandlerDecorator.DecorateCommandHandler(ctx, handler)
		if ctx == nil || handler == nil {
			panic("CommandMux.Decorator cannot return a nil context or handler")
		}
	}

	cmdResp := &commandResponseWriter{
		MessageResponseWriter: resp,
		req:                   req,
	}

	return handler.HandleCommand(ctx, cmdResp, cmdReq)
}

func readMainDocument(resp MessageResponseWriter, req *MessageRequest) (*CommandRequest, error) {
	length, _, _, opCode, rem, ok := wiremessage.ReadHeader(req.Message)
	if !ok {
		return nil, newError(nil, CodeProtocolError, "could not read msg header")
	}

	if length != int32(len(req.Message)) {
		return nil, newError(nil, CodeProtocolError, "message length is invalid")
	}

	switch opCode {
	case wiremessage.OpQuery:
		return readOpQueryMainDocument(resp, req, rem)
	case wiremessage.OpMsg:
		return readOpMsgMainDocument(resp, req, rem)
	}

	return nil, newErrorf(nil, CodeProtocolError, "unsupported wiremessage opcode for command: %s", opCode)
}

func readOpQueryMainDocument(resp MessageResponseWriter, req *MessageRequest, rem []byte) (*CommandRequest, error) {
	_, rem, ok := wiremessage.ReadQueryFlags(rem)
	if !ok {
		return nil, newError(nil, CodeProtocolError, "could not read query flags")
	}
	fullCollectionName, rem, ok := wiremessage.ReadQueryFullCollectionName(rem)
	if !ok {
		return nil, newError(nil, CodeProtocolError, "could not read full collection name")
	}
	idx := strings.Index(fullCollectionName, ".")
	if idx < 0 {
		return nil, newErrorf(nil, CodeProtocolError, "could not read collection name, got %q", fullCollectionName)
	}
	databaseName := fullCollectionName[:idx]

	_, rem, ok = wiremessage.ReadQueryNumberToSkip(rem)
	if !ok {
		return nil, newError(nil, CodeProtocolError, "could not read number to skip")
	}

	_, rem, ok = wiremessage.ReadQueryNumberToReturn(rem)
	if !ok {
		return nil, newError(nil, CodeProtocolError, "could not read number to return")
	}

	var doc bsoncore.Document
	doc, _, ok = wiremessage.ReadQueryQuery(rem)
	if !ok {
		return nil, newError(nil, CodeProtocolError, "could not read query")
	}

	return &CommandRequest{
		MessageRequest: req,
		DatabaseName:   databaseName,
		Document:       doc,
	}, nil
}

func readOpMsgMainDocument(resp MessageResponseWriter, req *MessageRequest, rem []byte) (*CommandRequest, error) {
	flags, rem, ok := wiremessage.ReadMsgFlags(rem)
	if !ok {
		return nil, newError(nil, CodeProtocolError, "could not read msg flags")
	}
	if flags&wiremessage.ChecksumPresent != 0 {
		if err := validateChecksum(req); err != nil {
			return nil, err
		}

		// remove the checksum from the back so it isn't considered when reading sections
		rem = rem[:len(rem)-4]
	}

	var doc bsoncore.Document
	for len(rem) > 0 {
		var stype wiremessage.SectionType
		stype, rem, ok = wiremessage.ReadMsgSectionType(rem)
		if !ok {
			return nil, newError(nil, CodeProtocolError, "could not read section type")
		}
		switch stype {
		case wiremessage.SingleDocument:
			doc, rem, ok = wiremessage.ReadMsgSectionSingleDocument(rem)
			if !ok {
				return nil, newError(nil, CodeProtocolError, "could not read type 0 section")
			}
		case wiremessage.DocumentSequence:
			// we don't care about this type yet, so let's just skip it...
			seqLength, _, ok := bsoncore.ReadLength(rem)
			if !ok {
				return nil, newError(nil, CodeProtocolError, "could not read type 1 section: not enough bytes")
			}
			if seqLength < 0 || seqLength >= int32(len(rem)) {
				return nil, newError(nil, CodeProtocolError, "could not read type 1 section: invalid length")
			}

			rem = rem[seqLength:]
		default:
			return nil, newErrorf(nil, CodeProtocolError, "unknown message section type %d", stype)
		}
	}

	if len(doc) == 0 {
		return nil, newError(nil, CodeProtocolError, "invalid OP_MSG: no type 0 section found")
	}

	dbValue, err := doc.LookupErr("$db")
	if err != nil {
		return nil, newError(nil, CodeProtocolError, "OP_MSG does not contain required field \"$db\"")
	}

	databaseName, ok := dbValue.StringValueOK()
	if !ok {
		return nil, newError(nil, CodeProtocolError, "\"$db\" must be a string")
	}

	return &CommandRequest{
		MessageRequest: req,
		DatabaseName:   databaseName,
		Document:       doc,
	}, nil
}

func validateChecksum(req *MessageRequest) error {
	if len(req.Message) < 4 {
		return newError(nil, CodeProtocolError, "message is invalid")
	}
	providedChecksum, _, ok := wiremessage.ReadMsgChecksum(req.Message[len(req.Message)-4:])
	if !ok {
		return newError(nil, CodeProtocolError, "could not read checksum")
	}
	actualCheckSum := crc32.Checksum(req.Message[0:len(req.Message)-4], crc32.MakeTable(crc32.Castagnoli))
	if providedChecksum != actualCheckSum {
		return newError(nil, CodeProtocolError, "checksum comparison failed")
	}

	return nil
}
