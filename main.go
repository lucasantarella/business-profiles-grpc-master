package main

import (
	"database/sql"
	_ "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	pb "github.com/lucasantarella/business-profiles-grpc-golib"
	"google.golang.org/grpc"
	"log"
	"lucasantarella.com/businesscards/internal"
	_ "lucasantarella.com/businesscards/utils"
	"net"
	"os"
)

func main() {
	_ = godotenv.Load() // Ignore error

	// Load listen port from env, and substitute if not present
	port, present := os.LookupEnv("GRPC_PORT")
	if !present {
		port = ":50051"
	} else {
		port = ":" + port
	}

	db, err := sql.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASS")+"@tcp("+os.Getenv("DB_HOST")+")/"+os.Getenv("DB_NAME")+"?charset=utf8&parseTime=True")
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(1)
	defer db.Close()

	log.Printf("listening on port: %s", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterProfilesServiceServer(s, &internal.Server{Db: db})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
