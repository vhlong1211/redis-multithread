 # Redis-Multithread

A simplified Redis-like server implemented in Go with multithreading support.  
This project was built as a learning exercise to understand how in-memory databases like Redis work internally, focusing on data structures, command handling, and concurrency.

---

## 🚀 Features
- Basic Redis-like commands (`SET`, `GET`, `DEL`, …)
- Support for multithreaded request handling
- Custom thread pool implementation
- Lightweight in-memory storage with Go data structures

---

## 🛠 Tech Stack
- **Language**: Go
- **Paradigm**: Concurrent / Multithreaded programming
- **Concepts**:  
  - Goroutines & channels  
  - Thread-per-connection vs. thread-pool models  
  - In-memory key-value storage  

