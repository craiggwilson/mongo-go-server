package main

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/craiggwilson/mongo-go-server/examples/basic/internal"
	"github.com/craiggwilson/mongo-go-server/mongo"
)

func main() {
	mux := mongo.NewCommandMux()
	mux.CommandHandlerDecorator = mongo.CommandDebuggingDecorator{}
	internal.RegisterBasicService(mux, &basicService{
		MaxWireVersion: 6,
		VersionArray:   []int32{4, 2, 0},
	})

	svr := &mongo.Server{
		Handler:             mux,
		ConnectionDecorator: mongo.TLSConnectionDecorator(&tls.Config{}),
	}

	log.Println("serving MongoDB...")
	if err := svr.ListenAndServe(context.Background(), mongo.DefaultAddr); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

type basicService struct {
	LogicalSessionTimeoutMinutes int32
	MaxBatchSize                 int32
	MaxWireVersion               int32
	MinWireVersion               int32
	ReadOnly                     bool
	VersionArray                 []int32
}

func (svc *basicService) BuildInfo(_ context.Context, _ *mongo.CommandRequest, _ *internal.BuildInfoRequest) (*internal.BuildInfoResponse, error) {
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

func (svc *basicService) Custom(ctx context.Context, _ *mongo.CommandRequest, _ *internal.CustomRequest) (*internal.CustomResponse, error) {
	return nil, &mongo.Error{
		Code:    10,
		Message: "AHAHAHA",
	}
}

func (svc *basicService) IsMaster(ctx context.Context, _ *mongo.CommandRequest, req *internal.IsMasterRequest) (*internal.IsMasterResponse, error) {
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
