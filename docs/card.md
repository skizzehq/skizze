#### Cardinality - card

HyperLogLog++ is an algorithm for approximating the cardinality of elements in a sketch, and does not require a body for the POST and GET requests.

**Creating** a new empty sketch of type HyperLogLog++ (card) with the id "sketch_1" (body of request will be ignored):
```{r, engine='bash', count_lines}
curl -XPOST http://localhost:3596/card/sketch_1
```


**Adding** values to the sketch with id "sketch_1":
```{r, engine='bash', count_lines}
curl -XPUT http://localhost:3596/card/sketch_1 -d '{
  "values":[
    "image",
    "rick grimes"
  ]
}'
```


**Retrieving** the cardinality of "sketch_1":
```{r, engine='bash', count_lines}
curl -XGET http://localhost:3596/card/sketch_1
```
returns
```json
{
  "result":2,
  "error":null
}
```


**Deleting** the sketch of type "card" with id "sketch_1":
```{r, engine='bash', count_lines}
curl -XDELETE http://localhost:3596/card/sketch_1
```
