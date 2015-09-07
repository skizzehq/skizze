# Skizze

[![Build Status](https://travis-ci.org/seiflotfy/skizze.svg?branch=master)](https://travis-ci.org/seiflotfy/skizze) [![license](http://img.shields.io/badge/license-Apache-blue.svg)](https://raw.githubusercontent.com/seiflotfy/skizze/master/LICENSE) [![Join the chat at https://gitter.im/seiflotfy/skizze](https://img.shields.io/badge/GITTER-join%20chat-green.svg)](https://gitter.im/seiflotfy/skizze?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Skizze ([ˈskɪt͡sə]: german for sketch) is a sketch data store to deal with all problems around counting and sketching using probabilistic data-structures.

Unlike a Key-Value store, Skizze does not store values, but rather appends values to a defined sketches, allowing one to solve frequency and cardinality queries in near O(1) time, with minimal memory footprint.

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

## API

See [API](file:///Users/seif/Projects/go/src/github.com/seiflotfy/skizze/api.html)
