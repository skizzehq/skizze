package server

import (
	"strings"

	"github.com/gogo/protobuf/proto"
	pb "github.com/seiflotfy/skizze/datamodel"

	"golang.org/x/net/context"
)

func (s *serverStruct) CreateSketch(ctx context.Context, in *pb.Sketch) (*pb.Sketch, error) {
	info := pb.NewEmptyInfo()
	info.Name = in.GetName()
	info.Type = strings.ToLower(in.GetType().String())
	info.Properties.Rank = uint(in.Defaults.GetRank())
	info.Properties.Capacity = uint(in.Defaults.GetCapacity())
	if err := s.manager.CreateSketch(info); err != nil {
		return nil, err
	}
	return in, nil
}

func (s *serverStruct) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddReply, error) {
	sketch := in.GetSketch()
	info := pb.NewEmptyInfo()
	info.Name = sketch.GetName()
	info.Type = strings.ToLower(sketch.GetType().String())
	err := s.manager.AddToSketch(info.ID(), in.GetValues())
	if err != nil {
		return nil, err
	}
	return &pb.AddReply{}, nil
}

func (s *serverStruct) GetMembership(ctx context.Context, in *pb.GetRequest) (*pb.GetMembershipReply, error) {
	sketch := in.GetSketch()
	info := pb.NewEmptyInfo()
	info.Name = sketch.GetName()
	info.Type = strings.ToLower(in.GetSketch().GetType().String())
	res, err := s.manager.GetFromSketch(info.ID(), in.GetValues())
	if err != nil {
		return nil, err
	}
	reply := &pb.GetMembershipReply{}
	// FIXME: return in same order
	for k, v := range res.(map[string]bool) {
		reply.Memberships = append(reply.Memberships, &pb.Membership{
			Value:    &k,
			IsMember: &v,
		})
	}
	return reply, nil
}

func (s *serverStruct) GetFrequency(ctx context.Context, in *pb.GetRequest) (*pb.GetFrequencyReply, error) {
	sketch := in.GetSketch()
	info := pb.NewEmptyInfo()
	info.Name = sketch.GetName()
	info.Type = strings.ToLower(in.GetSketch().GetType().String())
	res, err := s.manager.GetFromSketch(info.ID(), in.GetValues())
	if err != nil {
		return nil, err
	}
	reply := &pb.GetFrequencyReply{}
	// FIXME: return in same order
	for k, v := range res.(map[string]uint) {
		reply.Frequencies = append(reply.Frequencies, &pb.Frequency{
			Value:     &k,
			Frequency: proto.Int64(int64(v)),
		})
	}
	return reply, nil
}

func (s *serverStruct) GetCardinality(ctx context.Context, in *pb.GetRequest) (*pb.GetCardinalityReply, error) {
	sketch := in.GetSketch()
	info := pb.NewEmptyInfo()
	info.Name = sketch.GetName()
	info.Type = strings.ToLower(in.GetSketch().GetType().String())
	res, err := s.manager.GetFromSketch(info.ID(), in.GetValues())
	if err != nil {
		return nil, err
	}
	reply := &pb.GetCardinalityReply{
		Cardinality: proto.Int64(int64(res.(uint))),
	}
	return reply, nil
}

func (s *serverStruct) GetRank(ctx context.Context, in *pb.GetRequest) (*pb.GetRankReply, error) {
	sketch := in.GetSketch()
	info := pb.NewEmptyInfo()
	info.Name = sketch.GetName()
	info.Type = strings.ToLower(in.GetSketch().GetType().String())
	res, err := s.manager.GetFromSketch(info.ID(), in.GetValues())
	if err != nil {
		return nil, err
	}
	reply := &pb.GetRankReply{}
	for _, v := range res.([]*pb.Element) {
		reply.Ranks = append(reply.Ranks, &pb.Rank{
			Value: &v.Key,
			Rank:  proto.Int64(int64(v.Count)),
		})
	}
	return reply, nil
}
