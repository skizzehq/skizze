#### Dictionary

A dictionray is 100% accurate counting data structure that serves as a frequency table of events in a stream of data.

**Creating** a new empty sketch of type Dictionary (dict) with the id "sketch_2" and a capacity of 1000000:
```{r, engine='bash', count_lines}
curl -XPOST http://localhost:3596/dict/sketch_2
```

**Adding** values to the sketch with id "sketch_2":
```{r, engine='bash', count_lines}
curl -XPUT http://localhost:3596/dict/sketch_2 -d '{
  "values": ["marvel", "hulk", "marvel"]
}'
```

**Getting** the frequency for the values "marvel" and "hulk" in "sketch_2":
```{r, engine='bash', count_lines}
curl -XGET http://localhost:3596/dict/sketch_2 -d '{
  "values": ["marvel", "hulk"]
}'
```
returns the current count for each of these values:
```json
{
  "result":{
    "hulk":1,
    "marvel":2
  },
  "error":null
}
```

**Deleting** the sketch of type "dict" with id "sketch_2":
```{r, engine='bash', count_lines}
curl -XDELETE http://localhost:3596/dict/sketch_2
```
