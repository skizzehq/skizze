# Skizze

[![Build Status](https://travis-ci.org/seiflotfy/skizze.svg?branch=master)](https://travis-ci.org/seiflotfy/skizze)
[![license](http://img.shields.io/badge/license-Apache-blue.svg)](https://raw.githubusercontent.com/seiflotfy/counts/master/LICENSE)
[![Join the chat at https://gitter.im/seiflotfy/counts](https://img.shields.io/badge/GITTER-join%20chat-green.svg)](https://gitter.im/seiflotfy/counts?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)


A domain-sketch data store to deal with all problems around counting and sketching using probabilistic data-structures.

Unlike a Key-Value store, Skizze does not store values, but rather appends values to a sketch for a specified domain, allowing you to solve frequency and cardinality queries in near O(1) time, with minimal memory footprint.

<b>Note:</b> Data structures that can grow too big to reside in memory are read and written from/to disk directly via open stream to make sure we can maintain a high number of sketches.

#### Current status ==> pre-Alpha

## Motivation

From [Synopses for Massive Data: Samples, Histograms, Wavelets, Sketches](http://db.cs.berkeley.edu/cs286/papers/synopses-fntdb2012.pdf)
By Graham Cormode, Minos Garofalakis, Peter J. Haas and Chris Jermaine

#####The Need for Synopses
The use of synopses is essential for managing the massive data that arises in modern information management scenarios. When handling large datasets, from gigabytes to petabytes in size, it is often impractical to operate on them in full. Instead, it is much more convenient to build a synopsis, and then use this synopsis to analyze the data. This approach captures a variety of use-cases:

* A search engine collects logs of every search made, amounting to billions of queries every day. It would be too slow, and energy-intensive, to look for trends and patterns on the full data. Instead, it is preferable to use a synopsis that is guaranteed to preserve most of the as-yet undiscovered patterns in the data.
* A team of analysts for a retail chain would like to study the impact of different promotions and pricing strategies on sales of different items. It is not cost-effective to give each analyst the resources needed to study the national sales data in full, but by working with synopses of the data, each analyst can perform their explorations on their own laptops.
* A large cellphone provider wants to track the health of its network by studying statistics of calls made in different regions, on hardware from different manufacturers, under different levels of contention, and so on. The volume of information is too large to retain in a database, but instead the provider can build a synopsis of the data as it is observed live, and then use the synopsis off-line for further analysis.

These examples expose a variety of settings. The full data may reside in a traditional data warehouse, where it is indexed and accessible, but is too costly to work on in full. In other cases, the data is stored as flat files in a distributed file system; or it may never be stored in full, but be accessible only as it is observed in a streaming fashion. Sometimes synopsis construction is a one-time process, and sometimes we need to update the synopsis as the base data is modified or as accuracy requirements change. In all cases though, being able to construct a high quality synopsis enables much faster and more scalable data analysis.


## Other example problems?
* I want to know if a uri is in my spam list (spam list over a million entries)
* I want to know the number of users like my post (a like being subject to change)
* I want to know how may times oliver watched a video (counting frequencies)
* I want to know how many unique users visited my website in the last 3 hours? (sliding hyperloglog)
* I want to know the top 10 countries of users experiencing a crash (sliding hyperloglog)

## API
### RESTful API

| Method | Route | Parameters | Task |
| --- | --- | --- | --- |
| GET | / | N/A |Lists all available domains (sketches). |
| MERGE | / | not implmented yet | Merges multiple HyperLogLog counters. |
| POST | /<key> | {"domainName": string, "domainType": string, "capacity": uint64} | Creates a new Counter. DomainType is mandatory. DomainTypes can be found below. |
| GET | /<key> | N/A | Updates a domain. Adds values to a cardinality/counter to a domain. |
| PUT | /<key> | {"values": [string, string]} | Updates a domain. Adds values to a cardinality/counter to a domain. |
| PURGE | /<key> | {"values": [string, string]} | Purges values from a domain. |
| DELETE | /<key> | N/A | Deletes a domain. |

### DomainType
 - [x] <b>"cardinality"</b>: query unique items of all added values
  	* HyperLogLog
  	* does not support purging added values
  	* merge available soon
  	* capacity up to billions
 - [x] <b>"pcardinality"</b>: query unique items of all added values
 	* CuckooFilter
 	* allows puring added values
 	* requires way more space than Cardinality (1 byte per unique value)
 	* recommended capacity < 10.000.000
 	* disk usage is very intensive for now, caching coming soon
 - [x] <b>"frequency"</b>: query occurance frequenct of values
  	* Count-Min-Log Sketch
  	* integration under development
  	* recommended capacity < 1.000.000)
 - [ ] <b>"expiring"</b>: query cardinality withing the last n time units
 	* Sliding Hyper-Log-Log
 	* like HyperLogLog but with expiring entries
 - [x] <b>"topk"</b>: query the top k values added to the sketch
 	* Top-K Sketch


## Milestones
- [x] Design and implement REST API
- [x] Create domain manager
- [x] Integrate Cardinality Sketch (Hyperloglog++)
- [x] Integrate CardinalityPurgable Sketch (CuckooFilte)
- [x] Integrate Frequency Sketch (Count-Min-Log sketch)
- [ ] Integrate Expiring Sketch (Sliding Hyperloglog)
- [x] Integrate Top (TopK)
- [x] Store to Disk
- [ ] Replication on multiple servers
