#### Top-K Sketch
Top-K Sketch (topk) stores the k most popular elements of a data stream.

**Creating** a new empty sketch of type "topk" with the id "sketch_3" and a capacity of 10:
```
curl -XPOST http://localhost:3596/topk/sketch_3 -d '{"capacity": 10}'
```
<br>
**Adding** values to the sketch with id "sketch_3": 
```
curl -XPUT http://localhost:3596/topk/sketch_3 -d '{"values": ["dc", "batman"]}'
```
<br>
**Getting** the top k items in "sketch_3" (body of request will be ignored):
```
curl -XGET http://localhost:3596/topk/sketch_3
```
returns the current top k values:
```
{  
  "result":[  
    {  
      "Key":"batman",
      "Count":1,
      "Error":0
    },
    {  
      "Key":"dc",
      "Count":1,
      "Error":0
    }
  ],
  "error":null
}
```
<br>
**Deleting** the sketch of type "topk" with id "sketch_3":
```
curl -XDELETE http://localhost:3596/topk/sketch_3
```
