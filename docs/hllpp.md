#### HyperLogLog++

HyperLogLog is an algorithm for approximating the cardinality of elements in a sketch, and does not require a body for the POST and GET requests.

**Creating** a new empty sketch of type HyperLogLog++ (hllpp) with the id "sketch_1" (body of request will be ignored):
```
curl -XPOST http://localhost:3596/hllpp/sketch_1
```


**Adding** values to the sketch with id "sketch_1":
```
curl -XPUT http://localhost:3596/hllpp/sketch_1 -d '{  
  "values":[  
    "image",
    "rick grimes"
  ]
}'
```


**Retrieving** the cardinality of "sketch_1":
```
curl -XGET http://localhost:3596/hllpp/sketch_1
```
returns 
```
{  
  "result":2,
  "error":null
}
```


**Deleting** the sketch of type "hllpp" with id "sketch_1":
```
curl -XDELETE http://localhost:3596/hllpp/sketch_1
```
