# RabbitMQ Go Concurrency Benchmark

This project explores how concurrent worker pools in Go interact with RabbitMQ under load, with the goal of understanding performance characteristics, tuning parallelism, and using profiling tools like `pprof` to interpret CPU and memory behavior.

---

## 🧠 Motivation

- Learn how Go handles concurrency under realistic workloads
- Understand RabbitMQ's behavior with multiple consumers
- Benchmark performance with different worker pool sizes
- Profile CPU and memory usage to find bottlenecks
- Investigate oversubscription, goroutine scheduling, and backpressure

---

## 📦 Project Structure

```
.
├── main.go              # Switch between different run modes
├── benchmark.go         # Orchestrates test runs with different worker counts
├── consumer.go          # Worker pool and RabbitMQ consumption logic
├── producer.go          # Message publisher to RabbitMQ
├── profiles/            # CPU & heap profiles saved per test run
├── go.mod / go.sum
└── README.md            # This file
```

---

## 🚀 How It Works

### 1. Worker Pool Model

- Fixed-size goroutine pool (`N` workers) consuming messages from a shared channel
- Messages published to RabbitMQ before each benchmark run
- Workers consume messages with simulated delays (`time.Sleep(delayMs)`)

### 2. Benchmark Loop

- Runs tests with varying worker counts: `10, 20, 40, 80, 160, 320, 640`
- Measures:
  - Duration
  - Memory Allocation (`Alloc`, `Sys`)
  - Garbage Collections (`NumGC`)
  - CPU Usage (optional)
  - pprof CPU/heap profiles

### 3. Profiling (`pprof`)

- CPU profile saved per test (e.g., `cpu-20workers.pprof`)
- Heap profile saved using `pprof.WriteHeapProfile(...)`
- All saved in `profiles/` directory

➜  go-rabbitmq-demo go run . --mode=benchmark    

### Start RabbitMQ using Docker

If RabbitMQ is not running you will get:
                                                       
========================================
Benchmarking with 5 workers
2025/07/22 10:01:44 Publishing messages...
2025/07/22 10:01:44 Failed to connect to RabbitMQ: dial tcp [::1]:5672: connect: connection refused

---

Start RabbitMQ container:

```
docker compose up -d
```

Verify it’s running:

```
docker ps
```

Access RabbitMQ management UI:

```
http://localhost:15672
```

Login using:

- Username: guest
- Password: guest

Stop RabbitMQ later:

```
docker-compose down
```

- 5672 is used for AMQP protocol (your Go app will talk to RabbitMQ through this).
- 15672 is the HTTP port for the RabbitMQ Management UI.


## 📊 Sample Output

```
========== Benchmark Results ==========
Workers    Messages   Delay      Duration       
----------------------------------------------------------------------------------------------------
5          1000       100        20.218244958s  
10         1000       100        10.127634125s  
20         1000       100        5.048434125s   
40         1000       100        2.521909833s   
80         1000       100        1.319024916s   
160        1000       100        711.081875ms   
320        1000       100        410.482583ms   
640        1000       100        211.548542ms   
```

---

## 🔬 Profiling Results

- View profiles using `go tool pprof`:

```bash
go tool pprof -http=:8080 ./go-rabbitmq-demo profiles/cpu-320workers.pprof
```

- Visual tools:
  - **Top**: Shows time spent in each function
  - **Graph**: Call graph with nodes sized by CPU cost
  - **Flamegraph**: Stack-based heat visualization

- Observed issues:
  - Logging and I/O (e.g., `syscall.write`) dominate low-worker runs
  - Scheduling contention (e.g., `semasleep`, `pthread_cond_wait`) in high-worker runs
  - CPU usage drops at high concurrency → indication of oversubscription?

---

## 🧵 Lessons Learned

- **Too few workers** = underutilization, long durations
- **Too many workers** = context-switching, contention, diminishing returns
- **Optimal worker count** depends on workload, delay, and system limits
- Logging can **skew performance**—keep it minimal or async
- `pprof` is invaluable for **seeing what the CPU is *actually* doing**

---

## 🛠️ Possible Enhancements

- Add async logging or disable logs during benchmarking
- Use `GOMAXPROCS` tuning per test
- Automatically generate HTML flamegraphs (`pprof -http`)
- Export metrics to Prometheus / Grafana
- Dockerize the setup for consistent benchmarking

---

## 📚 References

- Go `pprof` Guide: https://golang.org/pkg/net/http/pprof/
- RabbitMQ Go Client: https://pkg.go.dev/github.com/streadway/amqp
- Go Memory Management: https://blog.golang.org/memory

---

## 📝 Final Notes

This project was built as a learning exercise and internal benchmark suite to better understand concurrency, messaging systems, and performance profiling in Go. It uses real-world tools and patterns you’d find in production systems.
