package main

import (
	"context"
	pb "github.com/lucasantarella/business-profiles-grpc-golib"
	"google.golang.org/grpc"
	"log"
	"lucasantarella.com/businesscards/utils"
	"os"
	"strconv"
	"time"
)

const (
	address          = "localhost:9001"
	defaultProfileId = 1
)

func main() {
	defer utils.TimeTrack(time.Now(), "Client Request")
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProfilesServiceClient(conn)

	// Contact the server and print out its response.
	profileId := defaultProfileId
	if len(os.Args) > 1 {
		profileId, err = strconv.Atoi(os.Args[1])
		if err != nil {
			panic("invalid int in argument")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetProfile(ctx, &pb.GetProfileRequest{Id: uint64(profileId)})
	if err != nil {
		log.Fatalf("could not validate: %v", err)
	}
	log.Printf("%v", r)
}
