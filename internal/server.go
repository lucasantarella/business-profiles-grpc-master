package internal

import (
	"context"
	"database/sql"
	pb "github.com/lucasantarella/business-profiles-grpc-golib"
	"log"
	"lucasantarella.com/businesscards/models"
	"lucasantarella.com/businesscards/utils"
	"time"
)

type Server struct {
	Db *sql.DB
}

func (s *Server) CreateProfile(context.Context, *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {
	defer utils.TimeTrack(time.Now(), "CreateProfile")
	log.Printf("onCreateProfile")
	panic("implement me")
}

func (s *Server) GetProfile(ctx context.Context, in *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	defer utils.TimeTrack(time.Now(), "GetProfile")
	log.Printf("onGetProfile")

	// create model
	profile := models.Profiles{}
	err := profile.Find(s.Db, int64(in.Id))

	if err != nil {
		resp := pb.GetProfileResponse{Ok: false, Error: err.Error()}
		return &resp, err
	}

	return &pb.GetProfileResponse{Ok: true, Profile: profile.ToPbProfile()}, nil
}

func (s *Server) UpdateProfile(context.Context, *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	defer utils.TimeTrack(time.Now(), "UpdateProfile")
	log.Printf("onUpdateProfile")
	panic("implement me")
}
