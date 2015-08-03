# Counts
A domain-counter data store

### Problem
* How many unique users visited my website last week? (once visited can't take it back)
* Is this URI in my spam list? (spam list over a million entries)
* How many users like my post? (a like being subject to change)
* How may times did oliver watch this video? (counting frequencies)

### TODO
- [ ] Design and implement REST API 
- [ ] Create counter manager
- [ ] Integrate Immutable Counter (Hyperloglog++)
- [ ] Integrate Mutable Counter (CuckooFilter and possibly play with the idea of Cubic HyperLogLog)
- [ ] Integrate Frequency Counter (minCount)
- [ ] Store to Disk
- [ ] Replication on multiple servers
