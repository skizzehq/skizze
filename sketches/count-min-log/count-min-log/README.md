# Count-Min-Log
[![GoDoc](https://godoc.org/github.com/seiflotfy/count-min-log?status.svg)](https://godoc.org/github.com/seiflotfy/count-min-log)

[Count-Min-Log sketch: Approximately counting with approximate counters - Guillaume Pitel & Geoffroy Fouquier](http://iswag-symposium.org/2015/pdfs/shortpaper1.pdf)

<b>TL;DR:</b> Count-Min-Log Sketch for improved Average Relative Error on low frequency events

Count-Min Sketch is a widely adopted algorithm for approximate event counting in large scale processing. However, the original version of the Count-Min-Sketch (CMS) suffers of some deficiences, especially if one is interested in the low-frequency items, such as in text- mining related tasks. Several variants of CMS have been proposed to compensate for the high relative error for low-frequency events, but the proposed solutions tend to correct the errors instead of preventing them. In this paper, we propose the Count-Min-Log sketch, which uses logarithm-based, approximate counters instead of linear counters to improve the average relative error of CMS at constant memory footprint.

## Example Usage

8-bit version

```go
import cml

...

sk, err := cml.NewDefaultSketch8()
sk.IncreaseCount([]byte("scott pilgrim"))
...

sk.GetCount("scott pilgrim")

// >> 1

```

16-bit version

```go
import cml

...

sk, err := cml.NewDefaultSketch16()
sk.IncreaseCount([]byte("scott pilgrim"))
...

sk.GetCount("scott pilgrim")

// >> 1

```
