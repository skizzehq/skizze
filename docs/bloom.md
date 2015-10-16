#### Bloom Filter

A Bloom filter is a representation of a set of n items, where the main requirement is to make membership queries; i.e., whether an item is a member of a set.

**Creating** a new empty sketch of type Bloom Filter (bloom) with the id "sketch_1" (body of request will be ignored):
```{r, engine='bash', count_lines}
curl -XPOST http://localhost:3596/bloom/sketch_1
```


**Adding** values to the sketch with id "sketch_1":
```{r, engine='bash', count_lines}
curl -XPUT http://localhost:3596/bloom/sketch_1 -d '{
  "values":[
    "image",
    "rick grimes"
  ]
}'
```



**Deleting** the sketch of type "bloom" with id "sketch_1":
```{r, engine='bash', count_lines}
curl -XDELETE http://localhost:3596/bloom/sketch_1
```
