package sketches

import (
	"errors"

	"github.com/seiflotfy/skizze/sketches/abstract"
	"github.com/seiflotfy/skizze/sketches/wrappers/count-min-log"
	"github.com/seiflotfy/skizze/sketches/wrappers/hllpp"
	"github.com/seiflotfy/skizze/sketches/wrappers/topk"
)

/*
SketchProxy ...
*/
type SketchProxy struct {
	*abstract.Info
	sketch abstract.Sketch
}

/*
Add ...
*/
func (sp *SketchProxy) Add(values [][]byte) (bool, error) {
	return sp.sketch.AddMultiple(values)
}

/*
Remove ...
*/
func (sp *SketchProxy) Remove(values [][]byte) (bool, error) {
	return sp.sketch.RemoveMultiple(values)
}

/*
Count ...
*/
func (sp *SketchProxy) Count(values []string) interface{} {
	if sp.Type == abstract.CML {
		bvalues := make([][]byte, len(values), len(values))
		for i, value := range values {
			bvalues[i] = []byte(value)
		}
		return sp.sketch.GetFrequency(bvalues)
	} else if sp.Type == abstract.TopK {
		return sp.sketch.GetFrequency(nil)
	}
	return sp.sketch.GetCount()
}

func createSketch(info *abstract.Info) (*SketchProxy, error) {
	var sketch abstract.Sketch
	var err error

	switch info.Type {
	case abstract.HLLPP:
		sketch, err = hllpp.NewSketch(info)
	case abstract.TopK:
		sketch, err = topk.NewSketch(info)
	case abstract.CML:
		sketch, err = cml.NewSketch(info)
	default:
		return nil, errors.New("Invalid sketch type: " + info.Type)
	}
	if err != nil {
		return nil, errors.New("Error creating new sketch")
	}
	sp := SketchProxy{info, sketch}

	return &sp, nil
}

func loadSketch(info *abstract.Info) (*SketchProxy, error) {
	var sketch abstract.Sketch
	var err error
	switch info.Type {
	case abstract.HLLPP:
		sketch, err = hllpp.NewSketchFromData(info)
	case abstract.TopK:
		sketch, err = topk.NewSketchFromData(info)
	case abstract.CML:
		sketch, err = cml.NewSketchFromData(info)
	default:
		logger.Info.Println("Invalid sketch type", info.Type)
	}
	sp := SketchProxy{info, sketch}

	if err != nil {
		return nil, errors.New("Error creating new sketch")
	}
	return &sp, nil
}
