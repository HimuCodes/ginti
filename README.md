# ginti (गिनती) - High-Performance Word Counter

Ginti is a Go-based implementation of the standard `wc` utility, optimized for modern hardware architectures and high-throughput I/O. It utilizes concurrency, SIMD-accelerated instructions, and zero-copy memory mapping to outperform the standard GNU C implementation.

## Performance Evolution

Target: 1.1GB text file (`lots_of_words.txt`)
Environment: Arch Linux (x86_64), Multi-core CPU, NVMe SSD

| Implementation | Execution Time | Speed vs wc | CPU Instructions | Key Technique |
| :--- | :--- | :--- | :--- | :--- |
| Baseline | 6.43s | 0.4x | 110B | `bufio.ReadRune` (UTF-8) |
| Phase 1 | 1.49s | 1.7x | 21.4B | 32KB Buffer + ASCII State Machine |
| Level 1 | 0.31s | 8.4x | 24.4B | `ReadAt` + Multi-core Concurrency |
| **Level 3** | **0.24s** | **11.2x** | **11.7B** | **Zero-Copy Memory Mapping (mmap)** |
| Standard `wc` | 2.70s | 1.0x | 30.3B | Sequential C Implementation |

## Implementation Methodology

### Baseline: UTF-8 Safety (Discarded)
The original implementation utilized `bufio.ReadRune()` to handle multi-byte Unicode characters. While safe, the overhead of UTF-8 decoding and per-byte function calls resulted in 110 billion instructions for a 1.1GB file.

### Phase 1: Fast Loop & ASCII Optimization
Switched to raw byte processing with a 32KB buffer to minimize syscall frequency. Optimized whitespace detection using an ASCII range check (`c <= 32`). This approach remains UTF-8 safe as all whitespace characters in UTF-8 are represented by single bytes below 128, while multi-byte sequences always utilize bytes above 127.

### Level 1: Multi-Worker Concurrency
Parallelized single-file processing by dividing the file into `N` chunks (where `N` = `runtime.NumCPU()`). Used `os.File.ReadAt` for thread-safe concurrent access to the same file descriptor. Implemented a boundary-peeking algorithm (`ReadAt(pos-1)`) to ensure word fragments across chunks were counted exactly once.

### Level 2: SIMD & Lookup Tables
1. **Line Counting:** Integrated `bytes.Count` for newline detection. The Go standard library implements this in raw assembly using SIMD (AVX2/AVX-512) instructions, processing up to 64 bytes per clock cycle.
2. **Word Counting:** Replaced branchy `if/else` logic with a 256-byte pre-computed lookup table. This reduced branch mispredictions and ensured $O(1)$ whitespace detection per byte.

### Level 3: Zero-Copy Memory Mapping (mmap)
The final optimization utilized `syscall.Mmap` to map the file directly into the process's virtual address space. 
* **Impact:** Eliminated the expensive kernel-to-user space memory copy overhead.
* **Instruction Efficiency:** Dropped total instruction count from 24.4B to 11.7B (a 52% reduction from Level 1).
* **Memory Management:** Utilized the OS page cache for "lazy loading," ensuring physical RAM usage is managed efficiently by the kernel regardless of file size.

## Technical Notes
* **Binary:** `go build -o counter`
* **Benchmarks:** Managed via `bench.sh` using `hyperfine` and `perf stat`.
* **Safety:** Implements a fallback to `ReadAt` logic if `mmap` is unsupported by the filesystem.
