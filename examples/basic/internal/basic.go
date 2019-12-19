package internal

import (
	"context"

	"github.com/craiggwilson/mongo-go-server/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type BuildInfoCommandHandler interface {
	BuildInfo(context.Context, *mongo.CommandRequest, *BuildInfoRequest) (*BuildInfoResponse, error)
}

type BuildInfoCommandHandlerFunc func(context.Context, *mongo.CommandRequest, *BuildInfoRequest) (*BuildInfoResponse, error)

func (f BuildInfoCommandHandlerFunc) BuildInfo(ctx context.Context, orig *mongo.CommandRequest, req *BuildInfoRequest) (*BuildInfoResponse, error) {
	return f(ctx, orig, req)
}

type BuildInfoRequest struct {
}

type BuildInfoResponse struct {
	OK           int32   `json:"ok" bson:"ok"`
	Version      string  `json:"version" bson:"version"`
	VersionArray []int32 `json:"versionArray" bson:"versionArray"`
}

type buildInfoCommandHandlerImpl struct {
	impl BuildInfoCommandHandler
}

func (h *buildInfoCommandHandlerImpl) HandleCommand(ctx context.Context, resp mongo.CommandResponseWriter, req *mongo.CommandRequest) error {
	var typedReq BuildInfoRequest
	if err := bson.Unmarshal(req.Document, &typedReq); err != nil {
		return &mongo.Error{
			Code:     mongo.CodeFailedToParse,
			CodeName: mongo.CodeToName(mongo.CodeFailedToParse),
			Message:  "invalid buildInfo command",
			Cause:    err,
		}
	}

	typedResp, err := h.impl.BuildInfo(ctx, req, &typedReq)
	if err != nil {
		return err
	}

	respDoc, err := bson.Marshal(typedResp)
	if err != nil {
		return &mongo.Error{
			Code:     mongo.CodeInternalError,
			CodeName: mongo.CodeToName(mongo.CodeInternalError),
			Message:  "failed marshaling buildInfo command output",
			Cause:    err,
		}
	}

	return resp.WriteSingleDocument(respDoc)
}

func RegisterBuildInfoCommandHandler(mux *mongo.CommandMux, h BuildInfoCommandHandler) {
	mux.Handlers["buildInfo"] = &buildInfoCommandHandlerImpl{impl: h}
	mux.Handlers["buildinfo"] = mux.Handlers["buildInfo"]
}

type IsMasterCommandHandler interface {
	IsMaster(context.Context, *mongo.CommandRequest, *IsMasterRequest) (*IsMasterResponse, error)
}

type IsMasterCommandHandlerFunc func(context.Context, *mongo.CommandRequest, *IsMasterRequest) (*IsMasterResponse, error)

func (f IsMasterCommandHandlerFunc) IsMaster(ctx context.Context, orig *mongo.CommandRequest, req *IsMasterRequest) (*IsMasterResponse, error) {
	return f(ctx, orig, req)
}

type IsMasterRequest struct {
	Compression []string `json:"compression" bson:"compression"`
}

type IsMasterResponse struct {
	OK                           int32    `json:"ok" bson:"ok"`
	Ismaster                     bool     `json:"ismaster" bson:"ismaster"`
	MaxBsonObjectSize            int32    `json:"maxBsonObjectSize" bson:"maxBsonObjectSize"`
	MaxMessageSizeBytes          int32    `json:"maxMessageSizeBytes" bson:"maxMessageSizeBytes"`
	MaxWriteBatchSize            int32    `json:"maxWriteBatchSize" bson:"maxWriteBatchSize"`
	LogicalSessionTimeoutMinutes int32    `json:"logicalSessionTimeoutMinutes" bson:"logicalSessionTimeoutMinutes"`
	MinWireVersion               int32    `json:"minWireVersion" bson:"minWireVersion"`
	MaxWireVersion               int32    `json:"maxWireVersion" bson:"maxWireVersion"`
	ReadOnly                     bool     `json:"readOnly" bson:"readOnly"`
	Compression                  []string `json:"compression" bson:"compression"`
}

type isMasterCommandHandlerImpl struct {
	impl IsMasterCommandHandler
}

func (h *isMasterCommandHandlerImpl) HandleCommand(ctx context.Context, resp mongo.CommandResponseWriter, req *mongo.CommandRequest) error {
	var typedReq IsMasterRequest
	if err := bson.Unmarshal(req.Document, &typedReq); err != nil {
		return &mongo.Error{
			Code:     mongo.CodeFailedToParse,
			CodeName: mongo.CodeToName(mongo.CodeFailedToParse),
			Message:  "invalid isMaster command",
			Cause:    err,
		}
	}

	typedResp, err := h.impl.IsMaster(ctx, req, &typedReq)
	if err != nil {
		return err
	}

	respDoc, err := bson.Marshal(typedResp)
	if err != nil {
		return &mongo.Error{
			Code:     mongo.CodeInternalError,
			CodeName: mongo.CodeToName(mongo.CodeInternalError),
			Message:  "failed marshaling isMaster command output",
			Cause:    err,
		}
	}

	return resp.WriteSingleDocument(respDoc)
}

func RegisterIsMasterCommandHandler(mux *mongo.CommandMux, h IsMasterCommandHandler) {
	mux.Handlers["isMaster"] = &isMasterCommandHandlerImpl{impl: h}
	mux.Handlers["ismaster"] = mux.Handlers["isMaster"]
}

type CustomCommandHandler interface {
	Custom(context.Context, *mongo.CommandRequest, *CustomRequest) (*CustomResponse, error)
}

type CustomCommandHandlerFunc func(context.Context, *mongo.CommandRequest, *CustomRequest) (*CustomResponse, error)

func (f CustomCommandHandlerFunc) Custom(ctx context.Context, orig *mongo.CommandRequest, req *CustomRequest) (*CustomResponse, error) {
	return f(ctx, orig, req)
}

type CustomRequest struct {
}

type CustomResponse struct {
	OK int32 `json:"ok" bson:"ok"`
}

type customCommandHandlerImpl struct {
	impl CustomCommandHandler
}

func (h *customCommandHandlerImpl) HandleCommand(ctx context.Context, resp mongo.CommandResponseWriter, req *mongo.CommandRequest) error {
	var typedReq CustomRequest
	if err := bson.Unmarshal(req.Document, &typedReq); err != nil {
		return &mongo.Error{
			Code:     mongo.CodeFailedToParse,
			CodeName: mongo.CodeToName(mongo.CodeFailedToParse),
			Message:  "invalid custom command",
			Cause:    err,
		}
	}

	typedResp, err := h.impl.Custom(ctx, req, &typedReq)
	if err != nil {
		return err
	}

	respDoc, err := bson.Marshal(typedResp)
	if err != nil {
		return &mongo.Error{
			Code:     mongo.CodeInternalError,
			CodeName: mongo.CodeToName(mongo.CodeInternalError),
			Message:  "failed marshaling custom command output",
			Cause:    err,
		}
	}

	return resp.WriteSingleDocument(respDoc)
}

func RegisterCustomCommandHandler(mux *mongo.CommandMux, h CustomCommandHandler) {
	mux.Handlers["custom"] = &customCommandHandlerImpl{impl: h}
}

type BasicService interface {
	BuildInfoCommandHandler
	IsMasterCommandHandler
	CustomCommandHandler
}

func RegisterBasicService(mux *mongo.CommandMux, svc BasicService) {
	RegisterBuildInfoCommandHandler(mux, svc)
	RegisterIsMasterCommandHandler(mux, svc)
	RegisterCustomCommandHandler(mux, svc)
}
