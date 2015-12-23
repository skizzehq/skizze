# API reference

Skizze is communicated with via a RESTful API. All methods apply on all different types of sketches (with optional parameters)

## Quick Overview
<b>Note:</b> Data structures that can grow too big to reside in memory are read and written from/to disk directly via open stream to make sure we can maintain a high number of sketches.

### Sketch Types

| type  | purpose     | Sketch               | Description                              | Notes |
| ---   | ---         | ---                  | ---                                      | ---   |
| hllpp | cardinality | HyperLogLog++        | query unique items from all added values | capacity up to billions, does not support purging added values |
| cml   | frequency   | Count-Min-Log Sketch | query frequency of unique values added   | N/A |
| topk  | rank + frequncy | Top-k Sketch | query the top k values added to the sketch | N/A |
| bloom | membership | Bloom Filter | query unique items from all added values | N/A |
| dictionary | frequency | Dictionary | query frequency of unique values added | infinte capacity (lots of memory), 100% accurate |

### RESTful API

| Method | Route      | Parameters                   | Task |
| ---    | ---        | ---                          | --- |
| GET    | /          | N/A                          | Lists all available sketches (sketches) |
| MERGE  | /          | not implemented yet          | Merges multiple sketches of the same <type> if they support merging |
| POST   | /$type/$id | {"capacity": uint64}         | Creates a new <type> sketch with id: <id> |
| GET    | /$type/$id | (optional) {"values": [string, ...]} | Get cardinality/frequency/rank of a sketch (for given values if supported by the sketch type) |
| PUT    | /$type/$id | {"values": [string, ...]} | Updates a sketch by adding values to it |
| PURGE  | /$type/$id | {"values": [string, ...]} | Updates a sketch by purging values from it |
| DELETE | /$type/$id | N/A                          | Deletes a sketch. |

### Example requests:


**Creating** a new empty sketch of type HyperLogLog++ (hllpp) with the id "sketch_1":
```{r, engine='bash', count_lines}
curl -XPOST http://localhost:3596/hllpp/sketch_1
```


**Adding** values to the sketch with id "sketch_1":
```{r, engine='bash', count_lines}
curl -XPUT http://localhost:3596/hllpp/sketch_1 -d '{
  "values": ["image", "rick grimes"]
}'
```


**Retrieving** the cardinality of "sketch_1":
```{r, engine='bash', count_lines}
curl -XGET http://localhost:3596/hllpp/sketch_1
```
returns
```json
{
  "result":2,
  "error":null
}
```

**Listing** all available sketches:
```{r, engine='bash', count_lines}
curl -XGET http://localhost:3596
```
returns
```json
{
  "result":[
    "hllpp/sketch_1"
  ],
  "error":null
}
```

**Deleting** the sketch of type "hllpp" with id "sketch_1":
```{r, engine='bash', count_lines}
curl -XDELETE http://localhost:3596/hllpp/sketch_1
```
---
For the API of each sketch type (implementation) look at the following type specific examples:
* [HyperLogLog++ (hllpp)](hllpp.md) (cardinality)
* [Count-Min-Log (cml)](cml.md) (frequency)
* [Top-K (topk)](topk.md) (ranking)
* [Bloom Filter (bloom)](bloom.md) (membership)
* [Dictionary (dict)](dict.md) (frequency)
