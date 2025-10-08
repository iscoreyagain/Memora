# Memora
a mini-Redis re-implementation with some inspired features by [DragonflyDB](https://www.dragonflydb.io/)

# Supported Features:
- _Multi-threaded_ and [**shared-nothing**](https://en.wikipedia.org/wiki/Shared-nothing_architecture) architect
- Multiplexing I/O using **epoll** for Linux and **kqueue** for macOS/BSD
- [RESP Protocol](https://redis.io/docs/latest/develop/reference/protocol-spec/)
- **Graceful shutdown** - Allowed to terminate the program elegantly
- **Approximated LRU** eviction 
- Commands:

  ○ PING

  ○ SET, GET, DEL, TTL, INFO

  ○ RPUSH, LPUSH, LLEN

  ○ SADD, SREM, SMEMBERS, SCARD, SINTER, SINTERSTORE, SDIFF, SDIFFSTORE, SUNION, SUNIONSTORE

  
- Data structures:

  ○ Hashtable
  
  ○ [B+ Tree](https://www.dragonflydb.io/blog/dragonfly-new-sorted-set)
  
  ○ [Bloom Filter](https://en.wikipedia.org/wiki/Bloom_filter)
  
  ○ [Skiplist](https://en.wikipedia.org/wiki/Skip_list)

  ○ [Count-min Sketch](https://en.wikipedia.org/wiki/Count%E2%80%93min_sketch)
  
  ○ [Sorted Set](https://redis.io/docs/latest/develop/data-types/sorted-sets/)

  ○ [Listpack](https://deepwiki.com/antirez/listpack/2.1-memory-layout-and-encoding#backlen-field-and-backward-traversal)

  ○ Quicklist
  
