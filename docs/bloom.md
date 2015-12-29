#### Bloom Filter

A Bloom filter is a representation of a set of n items, where the main requirement is to make membership queries; i.e., whether an item is a member of a set.

**Creating** a new empty sketch of type Bloom Filter (bloom) with the id "sketch_1" (body of request will be ignored):
```{r, engine='bash', count_lines}
curl -XPOST http://localhost:3596/bloom/sketch_1
```

* optional arguments:
  * capacity: the max capacity of values (does not apply to hllpp), default is 1000000.


**Adding** values to the sketch with id "sketch_1":
```{r, engine='bash', count_lines}
curl -XPUT http://localhost:3596/bloom/sketch_1 -d '{
  "values":[
    "image",
    "rick grimes"
  ]
}'
```

**Getting** the membership for the values "rick grimes" and "hulk" in "sketch_2":
```{r, engine='bash', count_lines}
curl -XGET http://localhost:3596/cml/sketch_2 -d '{
  "values": ["rick grimes", "hulk"]
}'
```
returns the current count for each of these values:
```json
{
  "result":{
    "hulk": false,
    "rick grimes": true
  },
  "error":null
}
```


**Deleting** the sketch of type "bloom" with id "sketch_1":
```{r, engine='bash', count_lines}
curl -XDELETE http://localhost:3596/bloom/sketch_1
```
