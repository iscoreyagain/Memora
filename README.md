# Memora
a mini-Redis re-implementation with some inspired features by [DragonflyDB](https://www.dragonflydb.io/)

# Supported Features:
- _Multi-threaded_ architect
- Multiplexing I/O using **epoll** for Linux and **kqueue** for macOS/BSD
- RESP Protocol
- **Graceful shutdown** - Allowed to terminate the program elegantly
- **Approximated LRU** eviction
- Commands: (To be Continued)
- Data structures:

  ○ Hashtable
  
  ○ [B+ Tree](https://www.dragonflydb.io/blog/dragonfly-new-sorted-set)
  
  ○ [Bloom Filter](https://en.wikipedia.org/wiki/Bloom_filter)
  
  ○ [Skiplist](https://en.wikipedia.org/wiki/Skip_list)

  ○ [Count-min Sketch](https://en.wikipedia.org/wiki/Count%E2%80%93min_sketch)
  
  ○ [Sorted Set](https://redis.io/docs/latest/develop/data-types/sorted-sets/)
  
