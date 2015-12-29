# API reference

Skizze is communicated with via a RESTful API. All methods apply on all different types of sketches (with optional parameters)

## Quick Overview
<b>Note:</b> Data structures that can grow too big to reside in memory are read and written from/to disk directly via open stream to make sure we can maintain a high number of sketches.

### Sketch Types

| type  | purpose         | Sketch               | Description                                | Notes |
| ---   | ---             | ---                  | ---                                        | ---   |
| card  | cardinality     | HyperLogLog++        | query unique items from all added values   | capacity up to billions, does not support purging added values |
| freq  | frequency       | Count-Min-Log Sketch | query frequency of unique values added     | N/A |
| rank  | rank + frequncy | Top-k Sketch         | query the top k values added to the sketch | N/A |
| memb  | membership      | Bloom Filter         | query sketch membership of a value         | N/A |
| dict  | frequency       | Dictionary           | query frequency of unique values added     | infinte capacity (lots of memory), 100% accurate |

### RESTful API

| Method | Route      | Parameters                           | Task |
| ---    | ---        | ---                                  | --- |
| GET    | /          | N/A                                  | Lists all available sketches (sketches) |
| MERGE  | /          | not implemented yet                  | Merges multiple sketches of the same <type> if they support merging |
| POST   | /$type/$id | {"capacity": uint64}                 | Creates a new <type> sketch with id: <id> |
| GET    | /$type/$id | (optional) {"values": [string, ...]} | Get cardinality/frequency/rank of a sketch (for given values if supported by the sketch type) |
| PUT    | /$type/$id | {"values": [string, ...]}            | Updates a sketch by adding values to it |
| PURGE  | /$type/$id | {"values": [string, ...]}            | Updates a sketch by purging values from it |
| DELETE | /$type/$id | N/A                                  | Deletes a sketch. |

### Example requests:


**Creating** a new empty sketch of type card with the id "sketch_1":
```{r, engine='bash', count_lines}
curl -XPOST http://localhost:3596/card/sketch_1
```

* optional arguments:
	* capacity: the max capacity of values (does not apply to card), default is 1000000.


**Adding** values to the sketch with id "sketch_1":
```{r, engine='bash', count_lines}
curl -XPUT http://localhost:3596/card/sketch_1 -d '{
  "values": ["image", "rick grimes"]
}'
```

* required arugments:
	* values: an array of values to be inserted into the sketch


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

**Listing** all available sketches:
```{r, engine='bash', count_lines}
curl -XGET http://localhost:3596
```
returns
```json
{
  "result":[
    "card/sketch_1"
  ],
  "error":null
}
```

**Deleting** the sketch of type "card" with id "sketch_1":
```{r, engine='bash', count_lines}
curl -XDELETE http://localhost:3596/card/sketch_1
```
---
For the API of each sketch type (implementation) look at the following type specific examples:
* [card (cardinality - HyperLogLog++)](card.md) (cardinality)
* [freq (frequeny - Count-Min-Log)](freq.md) (frequency)
* [rank (ranked - Top-K)](rank.md) (ranking)
* [memb (membership - Bloom Filter)](memb.md) (membership)
* [dict (dictionary - Dictionary)](dict.md) (dictionary)
