# Skizze

[![Build Status](https://travis-ci.org/seiflotfy/skizze.svg?branch=master)](https://travis-ci.org/seiflotfy/skizze)
[![license](http://img.shields.io/badge/license-Apache-blue.svg)](https://raw.githubusercontent.com/seiflotfy/counts/master/LICENSE)
[![Join the chat at https://gitter.im/seiflotfy/counts](https://img.shields.io/badge/GITTER-join%20chat-green.svg)](https://gitter.im/seiflotfy/counts?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)


A domain-counter data store to deal with all problems around counting using probabilistic data-structures.

Unlike a Key-Value store, Skizze does not store values, but rather keeps track of keys added to a domain, allowing you to solve frequency and cardinality queries in near O(1) time, with minimal memory footprint.

<b>Note:</b> Data structures that can grow too big to reside in memory are read and written from/to disk directly via open stream to make sure we can maintain a high number of counters.

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
* Is this URI in my spam list? (spam list over a million entries)
* How many users like my post? (a like being subject to change)
* How may times did oliver watch this video? (counting frequencies)
* How many unique users visited my website in the last 3 hours? (sliding hyperloglog)


## API-Documentation

	GET	/
	Lists all available counters.

	MERGE	/
	Merges multiple HyperLogLog counters.

	POST	/<key>
	Creates a new Counter.

	GET	/<key>
	Returns the count/cardinality of a counter.

	PUT	/<key>
	Updates a counter.
	Adds values to a cardinality/counter or increments a counter.

	PURGE	/<key>
	Purges values from a counter.

	DELETE	/<key>
	Deletes a counter.


## Milestones
- [x] Design and implement REST API
- [x] Create counter manager
- [x] Integrate Cardinality Counter (Hyperloglog++)
- [x] Integrate CardinalityPurgable Counter (CuckooFilter and possibly play with the idea of CuckooLogLog)
- [ ] Integrate Frequency Counter (Count-Min-Log sketch)
- [ ] Integrate FrequencyPurgable Counter (Count-Min sketch)
- [ ] Integrate Expiring Counter (Sliding Hyperloglog)
- [x] Integrate Top (TopK)
- [x] Store to Disk
- [ ] Replication on multiple servers
