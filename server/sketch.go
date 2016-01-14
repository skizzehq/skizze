package server

import (
	pb "github.com/seiflotfy/skizze/datamodel"

	"golang.org/x/net/context"
)

func (s *server) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddReply, error) {
	return nil, nil
}
func (s *server) GetMembership(ctx context.Context, in *pb.GetRequest) (*pb.GetMembershipReply, error) {
	return nil, nil
}
func (s *server) GetFrequency(ctx context.Context, in *pb.GetRequest) (*pb.GetFrequencyReply, error) {
	return nil, nil
}
func (s *server) GetCardinality(ctx context.Context, in *pb.GetRequest) (*pb.GetCardinalityReply, error) {
	return nil, nil
}
func (s *server) GetRank(ctx context.Context, in *pb.GetRequest) (*pb.GetRankReply, error) {
	return nil, nil
}
