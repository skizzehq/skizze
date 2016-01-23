package server

import (
	"strings"

	pb "datamodel"

	"github.com/gogo/protobuf/proto"
	"golang.org/x/net/context"
)

func (s *serverStruct) CreateSketch(ctx context.Context, in *pb.Sketch) (*pb.Sketch, error) {
	info := pb.NewEmptyInfo()
	info.Name = in.GetName()
	info.Type = strings.ToLower(in.GetType().String())
	info.Properties.Size = uint(in.Properties.GetSize())
	info.Properties.MaxUniqueItems = uint(in.Properties.GetMaxUniqueItems())
	if err := s.manager.CreateSketch(info); err != nil {
		return nil, err
	}
	return in, nil
}

func (s *serverStruct) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddReply, error) {
	info := pb.NewEmptyInfo()
	if dom := in.GetDomain(); dom != nil {
		info.Name = dom.GetName()
		info.Type = pb.DOM
		err := s.manager.AddToDomain(info.Name, in.GetValues())
		if err != nil {
			return nil, err
		}
	} else if sketch := in.GetSketch(); sketch != nil {
		info.Name = sketch.GetName()
		info.Type = strings.ToLower(sketch.GetType().String())
		err := s.manager.AddToSketch(info.ID(), in.GetValues())
		if err != nil {
			return nil, err
		}
	}
	return &pb.AddReply{}, nil
}

func (s *serverStruct) GetMembership(ctx context.Context, in *pb.GetRequest) (*pb.GetMembershipReply, error) {
	reply := &pb.GetMembershipReply{}

	for _, sketch := range in.GetSketches() {
		info := pb.NewEmptyInfo()
		info.Name = sketch.GetName()
		info.Type = strings.ToLower(pb.SketchType_MEMB.String())
		res, err := s.manager.GetFromSketch(info.ID(), in.GetValues())
		if err != nil {
			return nil, err
		}
		result := &pb.MembershipResult{}
		values := res.([]*pb.Member)
		// FIXME: return in same order
		for _, v := range values {
			result.Memberships = append(result.Memberships, &pb.Membership{
				Value:    proto.String(v.Key),
				IsMember: proto.Bool(v.Member),
			})
		}
		reply.Results = append(reply.Results, result)
	}
	return reply, nil
}

func (s *serverStruct) GetFrequency(ctx context.Context, in *pb.GetRequest) (*pb.GetFrequencyReply, error) {
	reply := &pb.GetFrequencyReply{}

	for _, sketch := range in.GetSketches() {
		info := pb.NewEmptyInfo()
		info.Name = sketch.GetName()
		info.Type = strings.ToLower(pb.SketchType_FREQ.String())
		res, err := s.manager.GetFromSketch(info.ID(), in.GetValues())
		if err != nil {
			return nil, err
		}
		result := &pb.FrequencyResult{}
		// FIXME: return in same order
		for k, v := range res.(map[string]uint) {
			result.Frequencies = append(result.Frequencies, &pb.Frequency{
				Value: proto.String(k),
				Count: proto.Int64(int64(v)),
			})
		}
		reply.Results = append(reply.Results, result)
	}
	return reply, nil
}

func (s *serverStruct) GetCardinality(ctx context.Context, in *pb.GetRequest) (*pb.GetCardinalityReply, error) {
	reply := &pb.GetCardinalityReply{}

	for _, sketch := range in.GetSketches() {
		info := pb.NewEmptyInfo()
		info.Name = sketch.GetName()
		info.Type = strings.ToLower(pb.SketchType_CARD.String())
		res, err := s.manager.GetFromSketch(info.ID(), in.GetValues())
		if err != nil {
			return nil, err
		}
		result := &pb.CardinalityResult{
			Cardinality: proto.Int64(int64(res.(uint))),
		}
		reply.Results = append(reply.Results, result)
	}
	return reply, nil
}

func (s *serverStruct) GetRankings(ctx context.Context, in *pb.GetRequest) (*pb.GetRankingsReply, error) {
	reply := &pb.GetRankingsReply{}

	for _, sketch := range in.GetSketches() {
		info := pb.NewEmptyInfo()
		info.Name = sketch.GetName()
		info.Type = strings.ToLower(pb.SketchType_RANK.String())
		res, err := s.manager.GetFromSketch(info.ID(), in.GetValues())
		if err != nil {
			return nil, err
		}
		result := &pb.RankingsResult{}
		for _, v := range res.([]*pb.Element) {
			result.Rankings = append(result.Rankings, &pb.Rank{
				Value: proto.String(v.Key),
				Count: proto.Int64(int64(v.Count)),
			})
		}
		reply.Results = append(reply.Results, result)
	}
	return reply, nil
}

func (s *serverStruct) DeleteSketch(ctx context.Context, in *pb.Sketch) (*pb.Empty, error) {
	info := pb.NewEmptyInfo()
	info.Name = in.GetName()
	info.Type = strings.ToLower(in.GetType().String())
	return &pb.Empty{}, s.manager.DeleteSketch(info.ID())
}

func (s *serverStruct) ListAll(ctx context.Context, in *pb.Empty) (*pb.ListReply, error) {
	sketches := s.manager.GetSketches()
	filtered := &pb.ListReply{}
	for _, v := range sketches {
		var typ pb.SketchType
		switch v[1] {
		case pb.CML:
			typ = pb.SketchType_FREQ
		case pb.TopK:
			typ = pb.SketchType_RANK
		case pb.HLLPP:
			typ = pb.SketchType_CARD
		case pb.Bloom:
			typ = pb.SketchType_MEMB
		default:
			continue
		}
		filtered.Sketches = append(filtered.Sketches, &pb.Sketch{Name: proto.String(v[0]), Type: &typ})
	}
	return filtered, nil
}

func (s *serverStruct) GetSketch(ctx context.Context, in *pb.Sketch) (*pb.Sketch, error) {
	var err error
	info := pb.NewEmptyInfo()
	info.Name = in.GetName()
	info.Type = strings.ToLower(in.GetType().String())
	if info, err = s.manager.GetSketch(info.ID()); err != nil {
		return in, err
	}
	in.Properties = &pb.SketchProperties{
		MaxUniqueItems: proto.Int64(int64(info.Properties.MaxUniqueItems)),
		Size:           proto.Int64(int64(info.Properties.Size)),
	}
	return in, nil
}

func (s *serverStruct) List(ctx context.Context, in *pb.ListRequest) (*pb.ListReply, error) {
	sketches := s.manager.GetSketches()
	filtered := &pb.ListReply{}
	for _, v := range sketches {
		var typ pb.SketchType
		switch v[1] {
		case pb.CML:
			typ = pb.SketchType_FREQ
		case pb.TopK:
			typ = pb.SketchType_RANK
		case pb.HLLPP:
			typ = pb.SketchType_CARD
		case pb.Bloom:
			typ = pb.SketchType_MEMB
		default:
			continue
		}
		if in.GetType() == typ {
			filtered.Sketches = append(filtered.Sketches, &pb.Sketch{Name: proto.String(v[0]), Type: &typ})
		}
	}
	return filtered, nil
}
