package server

import (
	"golang.org/x/net/context"

	pb "datamodel/protobuf"
)

func (s *serverStruct) GetSnapshot(ctx context.Context, in *pb.GetSnapshotRequest) (*pb.GetSnapshotReply, error) {
	return nil, nil
}
