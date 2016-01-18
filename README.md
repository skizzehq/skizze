![Skizze](http://i.imgur.com/9z47NdA.png)

[![Build Status](https://travis-ci.org/seiflotfy/skizze.svg?branch=master)](https://travis-ci.org/seiflotfy/skizze) [![license](http://img.shields.io/badge/license-Apache-blue.svg)](https://raw.githubusercontent.com/seiflotfy/skizze/master/LICENSE) [![Join the chat at https://gitter.im/seiflotfy/skizze](https://img.shields.io/badge/GITTER-join%20chat-green.svg)](https://gitter.im/seiflotfy/skizze?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Skizze ([ˈskɪt͡sə]: german for sketch) is a sketch data store to deal with all problems around counting and sketching using probabilistic data-structures.

Unlike a Key-Value store, Skizze does not store values, but rather appends values to defined sketches, allowing one to solve frequency and cardinality queries in near O(1) time, with minimal memory footprint.

<b> Current status ==> pre-Alpha </b>

## Motivation

Statistical analysis and mining of huge multi-terabyte data sets is a common task nowadays, especially in areas like web analytics and Internet advertising. Analysis of such large data sets often requires powerful distributed data stores like Hadoop and heavy data processing with techniques like MapReduce. This approach often leads to heavyweight high-latency analytical processes and poor applicability to realtime use cases. On the other hand, when one is interested only in simple additive metrics like total page views or average price of conversion, it is obvious that raw data can be efficiently summarized, for example, on a daily basis or using simple in-stream counters.  Computation of more advanced metrics like a number of unique visitor or most frequent items is more challenging and requires a lot of resources if implemented straightforwardly.

Skizze is a (fire and forget) service that provides a probabilistic data structures (sketches) storage that allows estimation of these and many other metrics, with a trade off in precision of the estimations for the memory consumption. These data structures can be used both as temporary data accumulators in query processing procedures and, perhaps more important, as a compact – sometimes astonishingly compact – replacement of raw data in stream-based computing.

## Example use cases (queries)?
* How many distinct elements are in the data set (i.e. what is the cardinality of the data set)?
* What are the most frequent elements (the terms “heavy hitters” and “top-k elements” are also used)?
* What are the frequencies of the most frequent elements?
* How many elements belong to the specified range (range query, in SQL it looks like `SELECT count(v) WHERE v >= c1 AND v < c2)?`
* Does the data set contain a particular element (membership query)?

## How to build and install

```
make dist
./bin/skizze
```

## Example usage:

```
./bin/skizze-cli
```

**Create** a new Domain (Collection of Sketches):
```{r, engine='bash', count_lines}
#CREATE DOM $name $estCardinality $topk
CREATE DOM demostream 10000000 100
```

**Add** values to the domain:
```{r, engine='bash', count_lines}
#ADD DOM $name $value1, $value2 ....
ADD DOM demostream zod joker grod zod zod grod
```

**Get** the *cardinality* of the domain:
```{r, engine='bash', count_lines}
# GET CARD $name
GET CARD demostream

# returns:
# Cardinality: 9
```

**Get** the *rankings* of the domain:
```{r, engine='bash', count_lines}
# GET RANK $name
GET RANK demostream

# returns:
# Rank: 1	  Value: zod	  Hits: 3
# Rank: 2	  Value: grod	  Hits: 2
# Rank: 3	  Value: joker	  Hits: 1
```

**Get** the *frequencies* of values in the domain:
```{r, engine='bash', count_lines}
# GET FREQ $name $value1 $value2 ...
GET FREQ demostream zod joker batman grod

# returns
# Value: zod	  Hits: 3
# Value: joker	  Hits: 1
# Value: batman	  Hits: 0
# Value: grod	  Hits: 2
```

**Get** the *membership* of values in the domain:
```{r, engine='bash', count_lines}
# GET MEMB $name $value1 $value2 ...
GET MEMB demostream zod joker batman grod

# returns
# Value: zod	  Member: true
# Value: joker	  Member: true
# Value: batman	  Member: false
# Value: grod	  Member: true
```

**List** all available sketches (created by domains):
```{r, engine='bash', count_lines}
LIST

# returns
# Name: demostream  Type: CARD
# Name: demostream  Type: FREQ
# Name: demostream  Type: MEMB
# Name: demostream  Type: RANK
```

**Create** a new sketch of type $type (CARD, MEMB, FREQ or RANK):
```{r, engine='bash', count_lines}
# CREATE CARD $name
CREATE CARD demosketch
```

**Add** values to the sketch of type $type (CARD, MEMB, FREQ or RANK):
```{r, engine='bash', count_lines}
#ADD $type $name $value1, $value2 ....
ADD CARD demostream zod joker grod zod zod grod
```

## In Progress:
### More REPL
* [x] Redesign data-structures main interface
* [x] Add new domains model
* [x] Add snapshotting
* [ ] Add AOF
* [ ] Add gRPC API
* [ ] Add REPL
  * [ ] DELETE DOM $name 							# delete domain and all its sketches
  * [ ] DELETE $type $name 						# delete a sketch of $type CARD, MEMB, FREQ, RANK and $name
  * [ ] ATTACH DOM $domName $sketchName $type 	# attach sketch $sketchName $type to $domName
  * [ ] LIST DOM 									# list all domains
  * [ ] SAVE 										# Explicityly save state of all domains and sketches
* [ ] New Docs
* [x] Clean up