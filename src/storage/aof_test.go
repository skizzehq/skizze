package storage

import (
	"config"
	pb "datamodel/protobuf"
	"path/filepath"
	"testing"

	"github.com/gogo/protobuf/proto"

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

	path := filepath.Join(config.GetConfig().DataDir, "skizze.aof")
	aof := NewAOF(path)
	aof.Run()

	dom := createDom("test1")
	err := aof.Append(CreateDom, dom)
	if err != nil {
		t.Error("Expected no error, got", err)
	}

	dom = createDom("test2")
	err = aof.Append(CreateDom, dom)
	if err != nil {
		t.Error("Expected no error, got", err)
	}

	// Create new AOF
	aof = NewAOF(path)

	for {
		e, err2 := aof.Read()
		if err2 != nil {
			if err2.Error() != "EOF" {
				t.Error("Expected no error, got", err2)
			}
			break
		}
		dom := &pb.Domain{}
		err = proto.Unmarshal(e.raw, dom)
		if err != nil {
			t.Error("Expected no error, got", err)
			break
		}
	}
	aof.Run()

	dom = createDom("test3")

	if err = aof.Append(CreateDom, dom); err != nil {
		t.Error("Expected no error, got", err)
	}

	dom = new(pb.Domain)
	dom.Name = utils.Stringp("test1")

	if err = aof.Append(DeleteDom, dom); err != nil {
		t.Error("Expected no error, got", err)
	}

	aof = NewAOF(path)
	for {
		e, err := aof.Read()
		if err != nil {
			if err.Error() != "EOF" {
				t.Error("Expected no error, got", err)
			}
			break
		}
		dom := &pb.Domain{}
		err = proto.Unmarshal(e.raw, dom)
		if err != nil {
			t.Error("Expected no error, got", err)
			break
		}
	}
}

func TestCreateDeleteSketch(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	path := filepath.Join(config.GetConfig().DataDir, "skizze.aof")
	aof := NewAOF(path)
	aof.Run()

	sketch := createSketch("skz1", pb.SketchType_CARD)
	err := aof.Append(CreateDom, sketch)
	if err != nil {
		t.Error("Expected no error, got", err)
	}

	sketch = createSketch("skz2", pb.SketchType_FREQ)
	err = aof.Append(CreateDom, sketch)
	if err != nil {
		t.Error("Expected no error, got", err)
	}

	// Create new AOF
	aof = NewAOF(path)
	for {
		e, err2 := aof.Read()
		if err2 != nil {
			if err2.Error() != "EOF" {
				t.Error("Expected no error, got", err2)
			}
			break
		}
		sketch := &pb.Sketch{}
		err = proto.Unmarshal(e.raw, sketch)
		if err != nil {
			t.Error("Expected no error, got", err)
		}
	}
	aof.Run()

	sketch = createSketch("skz3", pb.SketchType_RANK)
	err = aof.Append(CreateDom, sketch)
	if err != nil {
		t.Error("Expected no error, got", err)
	}

	sketch = createSketch("skz1", pb.SketchType_RANK)
	if err = aof.Append(DeleteDom, sketch); err != nil {
		t.Error("Expected no error, got", err)
	}

	addReq := &pb.AddRequest{
		Sketch: sketch,
		Values: []string{"foo", "bar", "hello", "world"},
	}
	if err = aof.Append(4, addReq); err != nil {
		t.Error("Expected no error, got", err)
	}

	aof = NewAOF(path)
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
			err = proto.Unmarshal(e.raw, req)
		} else {
			sketch := &pb.Sketch{}
			err = proto.Unmarshal(e.raw, sketch)
		}
		if err != nil {
			t.Error("Expected no error, got", err)
		}
	}
}
