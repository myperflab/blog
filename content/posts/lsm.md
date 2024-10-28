---
title: "LSM Tree Explained"
thumbnailImagePosition: left
thumbnailImage: /images/lsm.jpeg
date: 2024-10-24
categories:
- DS Algorithms
tags:
- Data Structure
showPagination: true
draft: False
---
<!--more-->

LSM Tree -- **Log-Structured Merge Tree**

An LSM Tree is a data structure commonly used in databases to efficiently handle high-volume write operations. It’s particularly popular in systems where data is mostly written in large quantities (e.g., time-series databases, key-value stores, and log management systems). LSM Tree works by organizing data in layers and leveraging a combination of in-memory and on-disk data structures to ensure highly writes efficient without dropping read speed. However, read speed will be not as efficiet as in rdbms (e.g. postgress DB) as those are read optimized. In lsm tree implementation, read speed depends on the multiple parameters which will be discussed in the next sections.

**Why LSM Tree is write efficient?**

It writes into memory rather than disk (storage) which makes it superfast. Lets explore how it works and its terminology and architecture.

In LSM Tree, there two types of operations which are important.
1. **User Operations in LSM Tree (User facing)**

![](/images/lsmtree_flow.png)

2. **Backend Operations (System-Managed)**

![](/images/lsmtree_flow_system_managed_ops.png)

Lets go deeper inunderstanding the above operations along with understanding how LSM Tree works:

**Write (insert)/Update Operation:**

![](/images/lsm_write.png)
1. When a user initiates a Create operation, Data gets inserted a new record or key-value pair into the database. In an LSM (Log-Structured Merge) Tree, data is initially stored an in-memory structure called the Memtable (a write-friendly structure, often implemented as a sorted tree or hash table), so this write is not directly made to disk which makes it write efficient.
2. Data in the Memtable is kept sorted, facilitating fast searches and inserts without needing immediate access to disk storage.
3. For durability, each insertion operation is also written to the Write-Ahead Log (WAL).
4. The WAL is a sequential log file where each write is recorded before it’s added to the Memtable, ensuring data isn’t lost if there’s a system crash before it’s persisted to disk.
5. The WAL is typically stored on disk and append-only, making it fast for sequential writes.
6. Over time, as more data is inserted, the Memtable fills up. And when it reaches a predefined threshold (based on memory limits), a Flush Operation is triggered, which moves the contents of the Memtable to disk.
7. The flush process writes data to a Sorted String Table (SSTable) on disk, which is an immutable, sorted file format optimized for LSM Trees.
8. SSTables are write-once and append-only, meaning that each SSTable is created once and not modified after being written. In case of an update to an existing key, a duplicate entry is recorded in the latest sstable and during compaction step old entry gets deleted. So at some point of time duplicate entries are possible to exist. And if read happens for the same key, it always gets the latest data even if more duplicates are present becasue read always happens from the latest sstable to oldatest table.
9. SSTables store the flushed Memtable data in a sorted manner, enabling efficient retrieval and minimizing the need for random disk access during reads.
10. Over time, multiple SSTables accumulate on disk, which can cause inefficiency and redundancy.
11. A Compaction process periodically merges these SSTables, removing duplicates and consolidating data to reduce disk space usage.
12. This process ensures that reads are more efficient by minimizing the number of SSTables that need to be accessed to retrieve data.
13. To speed up reads, each SSTable includes a Bloom Filter and Index structures.
14. Bloom Filters provide a probabilistic way to test if a key exists in an SSTable without scanning the entire file.
15. An index helps locate specific keys within the SSTable quickly, so reads don’t have to search the entire structure. All the sstable are indexed for the efficient read operation.

**Read Operation:**

![](/images/lsm_read.png)
The read operation in an LSM Tree starts when a user requests to retrieve a specific key or data. Since data is spread across both memory (in the Memtable) and disk (in multiple SSTables), the system (DB) must check multiple locations to retrieve the latest value for a key. **The search process begins in the Memtable**, as it contains the most recent data, including updates not yet flushed to disk. If the key is found here, the value is returned immediately. However, if the key is not found in the Memtable, **the system then checks the SSTables on disk**. To speed up this process, each SSTable is accompanied by Bloom Filters and indexes. Bloom Filters, probabilistic data structures, help quickly rule out SSTables that do not contain the key, reducing unnecessary file accesses. If a Bloom Filter indicates the possible presence of a key, the system uses the SSTable’s internal index to locate the exact position of the key in the file.

In cases where there are multiple versions of a key (due to delayed compaction), the system relies on timestamps to ensure it returns the most recent version. Thus, even with multiple sources of data, the LSM Tree can effectively manage read requests, maintaining efficient access by balancing in-memory structures, on-disk files, and background maintenance tasks like compaction and indexing. This architecture allows LSM Trees to handle high write throughput while keeping reads optimized and manageable across both memory and disk.

**Notes:** Bloom Filters and Index are associated with the SSTable. So to get the required data, DB iterates through SSTables, Bloom Filters, and Index until data is retrieved or all SSTables are exhausted.
