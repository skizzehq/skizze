package server

import (
	pb "datamodel/protobuf"

	"golang.org/x/net/context"
)

func (s *serverStruct) CreateSnapshot(ctx context.Context, in *pb.CreateSnapshotRequest) (*pb.CreateSnapshotReply, error) {
	err := s.manager.Save()
	// FIXME: return a snapshot ID
	status := pb.SnapshotStatus_PENDING
	return &pb.CreateSnapshotReply{Status: &status}, err
}
