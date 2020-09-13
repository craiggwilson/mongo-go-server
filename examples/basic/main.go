package main

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/craiggwilson/mongo-go-server/examples/basic/internal"
	"github.com/craiggwilson/mongo-go-server/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

var nextConnectionID = int32(0)

func main() {
	mux := mongo.NewCommandMux()
	mux.CommandHandlerDecorator = mongo.CommandDebuggingDecorator{}
	internal.RegisterBasicService(mux, &basicService{
		MaxWireVersion: 6,
		VersionArray:   []int32{4, 2, 0},
	})

	svr := &mongo.Server{
		Handler: mux,
		ConnectionDecorator: mongo.ChainConnectionDecorators(
			mongo.ConnectionDebuggingDecorator{},
			mongo.ConnectionDecoratorFunc(func(ctx context.Context, c net.Conn) (context.Context, net.Conn) {
				return context.WithValue(ctx, "basic-connection-state", &CustomConnectionState{
					connectionID: atomic.AddInt32(&nextConnectionID, 1),
				}), c
			}),
		),
	}

	log.Println("serving MongoDB...")
	if err := svr.ListenAndServe(context.Background(), mongo.DefaultAddr); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

type CustomConnectionState struct {
	connectionID int32
}

type basicService struct {
	LogicalSessionTimeoutMinutes int32
	MaxBatchSize                 int32
	MaxWireVersion               int32
	MinWireVersion               int32
	ReadOnly                     bool
	VersionArray                 []int32
}

func (svc *basicService) HandleAggregate(_ context.Context, raw *mongo.CommandRequest, req *internal.AggregateRequest) (*internal.AggregateResponse, error) {
	log.Printf("aggregate: pipeline: %s", req.Pipeline)

	_, batch := bsoncore.AppendArrayStart(nil)
	batch, _ = bsoncore.AppendArrayEnd(batch, 0)

	return &internal.AggregateResponse{
		OK: 1,
		Cursor: internal.CursorFirst{
			FirstBatch: batch,
			ID:         0,
			NS:         "",
		},
	}, nil
}

func (svc *basicService) HandleBuildInfo(_ context.Context, _ *mongo.CommandRequest, _ *internal.BuildInfoRequest) (*internal.BuildInfoResponse, error) {
	versionStrArray := make([]string, 0, len(svc.VersionArray))
	for _, p := range svc.VersionArray {
		versionStrArray = append(versionStrArray, strconv.Itoa(int(p)))
	}

	return &internal.BuildInfoResponse{
		OK:           1,
		Version:      strings.Join(versionStrArray, "."),
		VersionArray: svc.VersionArray,
	}, nil
}

func (svc *basicService) HandleGetLastError(ctx context.Context, _ *mongo.CommandRequest, _ *internal.GetLastErrorRequest) (*internal.GetLastErrorResponse, error) {
	state := ctx.Value("basic-connection-state").(*CustomConnectionState)

	return &internal.GetLastErrorResponse{
		OK:           1,
		WrittenTo:    "null",
		Err:          "null",
		ConnectionID: state.connectionID,
	}, nil
}

func (svc *basicService) HandleIsMaster(ctx context.Context, _ *mongo.CommandRequest, req *internal.IsMasterRequest) (*internal.IsMasterResponse, error) {
	svr := mongo.ServerFromContext(ctx)
	maxDocumentSize := svr.MaxDocumentSize
	if maxDocumentSize == 0 {
		maxDocumentSize = mongo.DefaultMaxDocumentSize
	}
	maxMessageSize := svr.MaxMessageSize
	if maxMessageSize == 0 {
		maxMessageSize = mongo.DefaultMaxMessageSize
	}

	compression := make([]string, 0, len(svr.Compressors))
	for _, sc := range svr.Compressors {
		for _, uc := range req.Compression {
			if sc.Name() == uc {
				compression = append(compression, sc.Name())
				break
			}
		}
	}

	return &internal.IsMasterResponse{
		OK:                           1,
		Ismaster:                     true,
		LogicalSessionTimeoutMinutes: svc.LogicalSessionTimeoutMinutes,
		MaxBsonObjectSize:            maxDocumentSize,
		MaxMessageSizeBytes:          maxMessageSize,
		MaxWriteBatchSize:            svc.MaxBatchSize,
		MaxWireVersion:               svc.MaxWireVersion,
		MinWireVersion:               svc.MinWireVersion,
		ReadOnly:                     svc.ReadOnly,
		Compression:                  compression,
	}, nil
}
