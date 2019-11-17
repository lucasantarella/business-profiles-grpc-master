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

func (s *Server) GetProfileSocialLinks(ctx context.Context, in *pb.GetProfileSocialLinksRequest) (*pb.GetProfileSocialLinksResponse, error) {
	defer utils.TimeTrack(time.Now(), "GetProfileSocialLinks")
	log.Printf("onGetProfileSocialLinks")

	// create model
	socialModel := models.ProfilesSocial{}
	socials, err := socialModel.FindByProfileID(s.Db, int64(in.ProfileId))

	if err != nil {
		return &pb.GetProfileSocialLinksResponse{}, err
	}

	socialPbSet := make([]*pb.ProfileSocialLink, len(socials))
	for i, social := range socials {
		socialPbSet[i] = &pb.ProfileSocialLink{
			Provider: pb.ProfileSocialLink_SocialProvider(social.Type),
			Value:    social.Value.String,
		}
	}
	return &pb.GetProfileSocialLinksResponse{Links: socialPbSet}, nil
}

func (s *Server) GetProfileExperiences(context.Context, *pb.GetProfileExperiencesRequest) (*pb.GetProfileExperiencesResponse, error) {
	panic("implement me")
}

func (s *Server) CreateProfile(context.Context, *pb.CreateProfileRequest) (*pb.Profile, error) {
	defer utils.TimeTrack(time.Now(), "CreateProfile")
	log.Printf("onCreateProfile")
	panic("implement me")
}

func (s *Server) GetProfile(ctx context.Context, in *pb.GetProfileRequest) (*pb.Profile, error) {
	defer utils.TimeTrack(time.Now(), "GetProfile")
	log.Printf("onGetProfile")

	// create model
	profile := models.Profiles{}
	err := profile.Find(s.Db, int64(in.Id))

	if err != nil {
		return nil, err
	}

	return profile.ToPbProfile(), nil
}

func (s *Server) UpdateProfile(context.Context, *pb.UpdateProfileRequest) (*pb.Profile, error) {
	defer utils.TimeTrack(time.Now(), "UpdateProfile")
	log.Printf("onUpdateProfile")
	panic("implement me")
}
