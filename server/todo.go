package server

import (
	pb "github.com/seiflotfy/skizze/datamodel"

	"golang.org/x/net/context"
)

func (s *server) List(ctx context.Context, in *pb.ListRequest) (*pb.ListReply, error) {
	return nil, nil
}
func (s *server) ListAll(ctx context.Context, in *pb.Empty) (*pb.ListReply, error) {
	return nil, nil
}
func (s *server) ListDomains(ctx context.Context, in *pb.Empty) (*pb.ListDomainsReply, error) {
	return nil, nil
}
func (s *server) SetDefaults(ctx context.Context, in *pb.Defaults) (*pb.Defaults, error) {
	return nil, nil
}
func (s *server) GetDefaults(ctx context.Context, in *pb.Empty) (*pb.Defaults, error) {
	return nil, nil
}
func (s *server) CreateDomain(ctx context.Context, in *pb.Domain) (*pb.Domain, error) {
	return nil, nil
}
func (s *server) DeleteDomain(ctx context.Context, in *pb.Domain) (*pb.Empty, error) {
	return nil, nil
}
func (s *server) GetDomain(ctx context.Context, in *pb.Domain) (*pb.Domain, error) {
	return nil, nil
}
func (s *server) CreateSketch(ctx context.Context, in *pb.Sketch) (*pb.Sketch, error) {
	return nil, nil
}
func (s *server) DeleteSketch(ctx context.Context, in *pb.Sketch) (*pb.Empty, error) {
	return nil, nil
}
func (s *server) GetSketch(ctx context.Context, in *pb.Sketch) (*pb.Sketch, error) {
	return nil, nil
}

func (s *server) CreateSnapshot(ctx context.Context, in *pb.CreateSnapshotRequest) (*pb.CreateSnapshotReply, error) {
	return nil, nil
}
func (s *server) GetSnapshot(ctx context.Context, in *pb.GetSnapshotRequest) (*pb.GetSnapshotReply, error) {
	return nil, nil
}
