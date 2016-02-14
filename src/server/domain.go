package server

import (
	"datamodel"
	pb "datamodel/protobuf"

	"storage"

	"golang.org/x/net/context"
)

func (s *serverStruct) createDomain(ctx context.Context, in *pb.Domain) (*pb.Domain, error) {
	info := datamodel.NewEmptyInfo()
	info.Name = in.Name
	// FIXME: A Domain's info should have an array of properties for each Sketch (or just an array
	// of Sketches, like what the proto has). This is just a hack to choose the first Sketch and
	// use it's info for now
	info.Properties.MaxUniqueItems = in.GetSketches()[0].GetProperties().MaxUniqueItems
	info.Properties.Size = in.GetSketches()[0].GetProperties().Size
	if info.Properties.Size == nil || *info.Properties.Size == 0 {
		var defaultSize int64 = 100
		info.Properties.Size = &defaultSize
	}
	// FIXME: We should be passing a pb.Domain and not a datamodel.Info to manager.CreateDomain
	err := s.manager.CreateDomain(info)
	if err != nil {
		return nil, err
	}
	return in, nil
}

func (s *serverStruct) CreateDomain(ctx context.Context, in *pb.Domain) (*pb.Domain, error) {
	if err := s.storage.Append(storage.CreateDom, in); err != nil {
		return nil, err
	}
	return s.createDomain(ctx, in)
}

func (s *serverStruct) ListDomains(ctx context.Context, in *pb.Empty) (*pb.ListDomainsReply, error) {
	res := s.manager.GetDomains()
	names := make([]string, len(res), len(res))
	for i, n := range res {
		names[i] = n[0]
	}
	doms := &pb.ListDomainsReply{
		Names: names,
	}
	return doms, nil
}

func (s *serverStruct) deleteDomain(ctx context.Context, in *pb.Domain) (*pb.Empty, error) {
	return &pb.Empty{}, s.manager.DeleteDomain(in.GetName())
}

func (s *serverStruct) DeleteDomain(ctx context.Context, in *pb.Domain) (*pb.Empty, error) {
	if err := s.storage.Append(storage.DeleteDom, in); err != nil {
		return nil, err
	}
	return s.deleteDomain(ctx, in)
}

func (s *serverStruct) GetDomain(ctx context.Context, in *pb.Domain) (*pb.Domain, error) {
	return s.manager.GetDomain(in.GetName())
}
