package server

import (
	pb "datamodel/protobuf"
	"errors"

	"golang.org/x/net/context"
)

func (s *serverStruct) CreateSnapshot(ctx context.Context, in *pb.CreateSnapshotRequest) (*pb.CreateSnapshotReply, error) {
	//err := s.manager.Save()
	// FIXME: return a snapshot ID
	status := pb.SnapshotStatus_FAILED
	return &pb.CreateSnapshotReply{Status: &status}, errors.New("Snapshots not supported (yet)!")
}
