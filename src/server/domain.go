package server

import (
	pb "datamodel"

	"golang.org/x/net/context"
)

func (s *serverStruct) CreateDomain(ctx context.Context, in *pb.Domain) (*pb.Domain, error) {
	info := pb.NewEmptyInfo()
	info.Name = in.GetName()
	info.Type = in.GetType().String()
	info.Properties.Capacity = uint(in.GetDefaults().GetCapacity())
	info.Properties.Rank = uint(in.GetDefaults().GetRank())
	err := s.manager.CreateDomain(info)
	if err != nil {
		return nil, err
	}
	return in, nil
}
