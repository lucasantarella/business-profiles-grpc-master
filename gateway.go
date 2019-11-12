package main

import (
	"flag"
	"github.com/joho/godotenv"
	"net/http"
	"os"

	"context"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "github.com/lucasantarella/business-profiles-grpc-golib"
	"google.golang.org/grpc"
)

func run() error {
	_ = godotenv.Load() // Ignore error

	grpcPort, present := os.LookupEnv("GRPC_PORT")
	if !present {
		grpcPort = ":50051"
	} else {
		grpcPort = ":" + grpcPort
	}

	httpListenPort, present := os.LookupEnv("HTTP_PORT")
	if !present {
		httpListenPort = ":9001"
	} else {
		httpListenPort = ":" + httpListenPort
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterProfilesServiceHandlerFromEndpoint(ctx, mux, "127.0.0.1"+grpcPort, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(httpListenPort, mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
