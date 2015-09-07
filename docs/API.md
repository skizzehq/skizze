# Skizze API

Skizze is communicated with via a RESTful API. All methods apply on all different types of sketches (with optional parameters)

## Quick Overview
<b>Note:</b> Data structures that can grow too big to reside in memory are read and written from/to disk directly via open stream to make sure we can maintain a high number of sketches.

### Sketch Types

| type  | purpose     | Sketch               | Description                              | Notes |
| ---   | ---         | ---                  | ---                                      | ---   |
| hllpp | cardinality | HyperLogLog++        | query unique items from all added values | capacity up to billions, does not support purging added values |
| cml   | frequency   | Count-Min-Log Sketch | query frequency of unique values added   | N/A |
| topk  | rank + frequncy | Top-k Sketch | query the top k values added to the sketch | N/A |

### RESTful API

| Method | Route      | Parameters                   | Task |
| ---    | ---        | ---                          | --- |
| GET    | /          | N/A                          | Lists all available sketches (sketches) |
| MERGE  | /          | not implemented yet          | Merges multiple sketches of the same <type> if they support merging |
| POST   | /$type/$id | {"capacity": uint64}         | Creates a new <type> sketch with id: <id> |
| GET    | /$type/$id | (optional) {"values": [string, string]} | Get cardinality/frequency/rank of a sketch (for given values if supported by the sketch type) |
| PUT    | /$type/$id | {"values": [string, string]} | Updates a sketch by adding values to it |
| PURGE  | /$type/$id | {"values": [string, string]} | Updates a sketch by purging values from it |
| DELETE | /$type/$id | N/A                          | Deletes a sketch. |

