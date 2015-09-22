#### Count-Min-Log

A count–min-log sketch is a probabilistic data structure that serves as a frequency table of events in a stream of data.

**Creating** a new empty sketch of type Count-Min-Log (cml) with the id "sketch_2" and a capacity of 1000000:
```
curl -XPOST http://localhost:3596/cml/sketch_2 -d '{
  "capacity": 1000000
}'
```

**Adding** values to the sketch with id "sketch_2":
```
curl -XPUT http://localhost:3596/cml/sketch_2 -d '{
  "values": ["marvel", "hulk", "marvel"]
}'
```

**Getting** the frequency for the values "marvel" and "hulk" in "sketch_2":
```
curl -XGET http://localhost:3596/cml/sketch_2 -d '{
  "values": ["marvel", "hulk"]
}'
```
returns the current count for each of these values:
```
{  
  "result":{  
    "hulk":1,
    "marvel":2
  },
  "error":null
}
```

**Deleting** the sketch of type "cml" with id "sketch_2":
```
curl -XDELETE http://localhost:3596/cml/sketch_2
```
