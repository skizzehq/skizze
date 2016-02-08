package storage

import (
	"os"
	"testing"

	"config"
	pb "datamodel/protobuf"

	"github.com/golang/protobuf/proto"

	"utils"
)

func createDom(id string) *pb.Domain {
	dom := new(pb.Domain)
	dom.Name = utils.Stringp(id)
	types := []pb.SketchType{pb.SketchType_MEMB, pb.SketchType_FREQ, pb.SketchType_RANK, pb.SketchType_CARD}
	for _, ty := range types {
		sketch := &pb.Sketch{}
		sketch.Name = dom.Name
		sketch.Type = &ty
		sketch.Properties = &pb.SketchProperties{
			Size:           utils.Int64p(100),
			MaxUniqueItems: utils.Int64p(10),
		}
		dom.Sketches = append(dom.Sketches, sketch)
	}
	return dom
}

func createSketch(id string, typ pb.SketchType) *pb.Sketch {
	sketch := &pb.Sketch{}
	sketch.Name = utils.Stringp(id)
	sketch.Type = &typ
	sketch.Properties = &pb.SketchProperties{
		Size:           utils.Int64p(100),
		MaxUniqueItems: utils.Int64p(10),
	}
	return sketch
}

func TestCreateDeleteDom(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	_ = os.Remove("test.log")

	aof := NewAOF("test.log")
	go aof.Run()

	dom := createDom("test1")
	aof.AppendDomOp(CreateDom, dom)

	dom = createDom("test2")
	aof.AppendDomOp(CreateDom, dom)

	aof.Stop()

	// Create new AOF
	aof = NewAOF("test.log")
	go aof.Run()
	for {
		e, err2 := aof.Read()
		if err2 != nil {
			if err2.Error() != "EOF" {
				t.Error("Expected no error, got", err2)
			}
			break
		}
		dom := &pb.Domain{}
		err := proto.Unmarshal(e.args, dom)
		if err != nil {
			t.Error("Expected no error, got", err)
		}
	}

	dom = createDom("test3")

	aof.AppendDomOp(CreateDom, dom)

	dom = new(pb.Domain)
	dom.Name = utils.Stringp("test1")

	aof.AppendDomOp(DeleteDom, dom)

	aof.Stop()

	aof = NewAOF("test.log")
	go aof.Run()
	for {
		e, err := aof.Read()
		if err != nil {
			if err.Error() != "EOF" {
				t.Error("Expected no error, got", err)
			}
			break
		}
		dom := &pb.Domain{}
		err = proto.Unmarshal(e.args, dom)
		if err != nil {
			t.Error("Expected no error, got", err)
		}
	}
}

func TestCreateDeleteSketch(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	_ = os.Remove("test.log")

	aof := NewAOF("test.log")
	go aof.Run()

	sketch := createSketch("skz1", pb.SketchType_CARD)
	aof.AppendDomOp(CreateDom, sketch)

	sketch = createSketch("skz2", pb.SketchType_FREQ)
	aof.AppendDomOp(CreateDom, sketch)

	aof.Stop()

	// Create new AOF
	aof = NewAOF("test.log")
	go aof.Run()

	for {
		e, err2 := aof.Read()
		if err2 != nil {
			if err2.Error() != "EOF" {
				t.Error("Expected no error, got", err2)
			}
			break
		}
		sketch := &pb.Sketch{}
		err := proto.Unmarshal(e.args, sketch)
		if err != nil {
			t.Error("Expected no error, got", err)
		}
	}

	sketch = createSketch("skz3", pb.SketchType_RANK)
	aof.AppendDomOp(CreateDom, sketch)

	sketch = createSketch("skz1", pb.SketchType_RANK)
	aof.AppendDomOp(DeleteDom, sketch)

	addReq := &pb.AddRequest{
		Sketch: sketch,
		Values: []string{"foo", "bar", "hello", "world"},
	}
	aof.AppendAddOp(addReq)

	aof.Stop()

	aof = NewAOF("test.log")
	go aof.Run()
	for {
		e, err := aof.Read()
		if err != nil {
			if err.Error() != "EOF" {
				t.Error("Expected no error, got", err)
			}
			break
		}
		if e.op == Add {
			req := &pb.AddRequest{}
			err = proto.Unmarshal(e.args, req)
		} else {
			sketch := &pb.Sketch{}
			err = proto.Unmarshal(e.args, sketch)
		}
		if err != nil {
			t.Error("Expected no error, got", err)
		}
	}
}

func TestStop(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	_ = os.Remove("test.log")

	aof := NewAOF("test.log")
	go aof.Run()

	sketch := createSketch("skz1", pb.SketchType_CARD)
	aof.AppendDomOp(CreateDom, sketch)

	aof.Stop()

	sketch = createSketch("skz2", pb.SketchType_FREQ)
	aof.AppendDomOp(CreateDom, sketch)
}
