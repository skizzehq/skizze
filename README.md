# Counts
A domain-counter data store

### Problem
* Is this URI in my spam list? (spam list over a million entries)
* How many users like my post? (a like being subject to change)
* How may times did oliver watch this video? (counting frequencies)
* How many unique users visited my website in the last 3 hours? (sliding hyperloglog)


### TODO
- [x] Design and implement REST API 
- [x] Create counter manager
- [x] Integrate UniqueIncremental Counter (Hyperloglog++)
- [x] Integrate Unique (CuckooFilter and possibly play with the idea of CuckooLogLog)
- [ ] Integrate UniqueFrequency Counter (minCount)
- [ ] Integrate UniqueExpiring (Sliding Hyperloglog)
- [ ] Integrate Free (Just a plain +1 and -1 Counter)
- [ ] Store to Disk
- [ ] Replication on multiple servers
