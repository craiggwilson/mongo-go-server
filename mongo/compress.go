package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

// ErrUnknownCompresor represents an unknown compressor error.
var ErrUnknownCompressor = errors.New("mongo: unknown compressor")

// CompressorError exposes an error that occurs during compression or decompression.
type CompressorError struct {
	CompressorID wiremessage.CompressorID
	Err          error
}

// Error implements the errors.Error interface.
func (e *CompressorError) Error() string {
	return fmt.Sprintf("mongo-compressor(%d): %v", e.CompressorID, e.Err)
}

// Unwrap implements the errors.Unwrap interface.
func (e *CompressorError) Unwrap() error {
	return e.Err
}

// Compressor handles compression and decompression of a wire message.
type Compressor interface {
	// ID identifies the Compressor.
	ID() wiremessage.CompressorID
	// Name is the name of the Compressor.
	Name() string
	// Compress compresses the bytes.
	Compress(b []byte) ([]byte, error)
	// Decompress decompresses the bytes.
	Decompress(uncompressed []byte, compressed []byte) error
}

// CompressionHandler handles messages in a compressed format and then
// delegates the processing of the uncompressed message to an inner handler.
type CompressionHandler struct {
	Compressors []Compressor
	Handler     MessageHandler
}

// HandleMessage implements the Handler interface.
func (h *CompressionHandler) HandleMessage(ctx context.Context, resp MessageResponseWriter, req *MessageRequest) error {
	msg, compressor, err := h.decompress(ctx, req)
	if err != nil {
		return err
	}

	req.Message = msg
	rw := &compressedResponseWriter{wrapped: resp, compressor: compressor}

	return h.Handler.HandleMessage(ctx, rw, req)
}

func (h *CompressionHandler) decompress(ctx context.Context, req *MessageRequest) (wiremessage.WireMessage, Compressor, error) {
	length, reqid, respto, opcode, rem, ok := wiremessage.ReadHeader(req.Message)
	if !ok || length < 16 || len(req.Message) < int(length) {
		return nil, nil, errors.New("malformed wire message: insufficient bytes")
	}
	if opcode != wiremessage.OpCompressed {
		return req.Message, nil, nil
	}

	opcode, rem, ok = wiremessage.ReadCompressedOriginalOpCode(rem)
	if !ok {
		return nil, nil, errors.New("malformed OP_COMPRESSED: missing original opcode")
	}

	uncompressedSize, rem, ok := wiremessage.ReadCompressedUncompressedSize(rem)
	if !ok {
		return nil, nil, errors.New("malformed OP_COMPRESSED: missing uncompressed size")
	}
	if uncompressedSize < 0 {
		return nil, nil, errors.New("malformed OP_COMPRESSED: uncompressed size is out of bounds")
	}

	compressorID, rem, ok := wiremessage.ReadCompressedCompressorID(rem)
	if !ok {
		return nil, nil, errors.New("malformed OP_COMPRESSED: missing compressor ID")
	}

	compressor := h.compressorByID(ctx, req, compressorID)
	if compressor == nil {
		return nil, nil, ErrUnknownCompressor
	}

	compressedSize := length - 25 // header (16) + original opcode (4) + uncompressed size (4) + compressor ID (1)
	if compressedSize < 0 {
		return nil, compressor, errors.New("malformed OP_COMPRESSED: compressed size is out of bounds")
	}

	cmsg, _, ok := wiremessage.ReadCompressedCompressedMessage(rem, compressedSize)
	if !ok {
		return nil, compressor, errors.New("malformed OP_COMPRESSED: insufficient bytes for compressed wiremessage")
	}

	uncompressed := make([]byte, 0, uncompressedSize+16)
	uncompressed = wiremessage.AppendHeader(uncompressed, uncompressedSize+16, reqid, respto, opcode)
	uncompressed = uncompressed[0:cap(uncompressed)] // change the length of the uncompressed slice up to its capacity

	err := compressor.Decompress(uncompressed[16:], cmsg)
	if err != nil {
		return nil, compressor, &CompressorError{
			CompressorID: compressor.ID(),
			Err:          err,
		}
	}

	return wiremessage.WireMessage(uncompressed), compressor, nil
}

func (h *CompressionHandler) compressorByID(ctx context.Context, req *MessageRequest, id wiremessage.CompressorID) Compressor {
	compressors := h.Compressors
	if len(compressors) == 0 {
		if svr := ServerFromContext(ctx); svr != nil {
			compressors = svr.Compressors
		}
	}

	for _, c := range compressors {
		if c.ID() == id {
			return c
		}
	}

	return nil
}

type compressedResponseWriter struct {
	wrapped    MessageResponseWriter
	compressor Compressor
}

func (rw *compressedResponseWriter) WriteMessage(msg wiremessage.WireMessage) error {
	if rw.compressor == nil {
		return rw.wrapped.WriteMessage(msg)
	}

	_, reqid, respto, origcode, rem, ok := wiremessage.ReadHeader(msg)
	if !ok {
		return errors.New("wiremessage is too short to compress, less than 16 bytes")
	}

	compressed, err := rw.compressor.Compress(rem)
	if err != nil {
		return &CompressorError{
			CompressorID: rw.compressor.ID(),
			Err:          err,
		}
	}

	// TODO: can use a buffer here somehow?
	_, compressedMsg := wiremessage.AppendHeaderStart(nil, reqid, respto, wiremessage.OpCompressed)
	compressedMsg = wiremessage.AppendCompressedOriginalOpCode(compressedMsg, origcode)
	compressedMsg = wiremessage.AppendCompressedUncompressedSize(compressedMsg, int32(len(rem)))
	compressedMsg = wiremessage.AppendCompressedCompressorID(compressedMsg, rw.compressor.ID())
	compressedMsg = wiremessage.AppendCompressedCompressedMessage(compressedMsg, compressed)
	bsoncore.UpdateLength(compressedMsg, 0, int32(len(compressedMsg)))

	return rw.wrapped.WriteMessage(compressedMsg)
}
