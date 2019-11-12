package main

import (
	"database/sql"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"lucasantarella.com/businesscards/internal"
	"net"
	"net/http"
	"os"
	"strings"

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

	db, err := sql.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASS")+"@tcp("+os.Getenv("DB_HOST")+")/"+os.Getenv("DB_NAME")+"?charset=utf8&parseTime=True")
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(1)
	defer db.Close()

	log.Printf("grpc listening on port: %s", grpcPort)
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterProfilesServiceServer(s, &internal.Server{Db: db})

	go s.Serve(lis)

	conn, err := grpc.DialContext(
		context.Background(),
		"127.0.0.1"+grpcPort,
		grpc.WithInsecure(),
	)

	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	router := runtime.NewServeMux()
	if err = pb.RegisterProfilesServiceHandler(context.Background(), router, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	log.Printf("http listening on port: %s", httpListenPort)
	return http.ListenAndServe(httpListenPort, httpGrpcRouter(s, router))
}

func httpGrpcRouter(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	})
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
