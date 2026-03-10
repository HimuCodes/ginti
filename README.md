# ginti (गिनती) - High-Performance Word Counter

Ginti is a Go-based implementation of the standard `wc` utility, optimized for modern hardware architectures and high-throughput I/O. It utilizes concurrency, SIMD-accelerated instructions, and zero-copy memory mapping to outperform the standard GNU C implementation.

## Performance Evolution

Target: 1.1GB text file (`lots_of_words.txt`)
Environment: Arch Linux (x86_64), Zen 3 Architecture, NVMe SSD

| Level | Execution Time | Speed vs wc | CPU Instructions | Key Technique |
| :--- | :--- | :--- | :--- | :--- |
| Baseline | 6.43s | 0.4x | 110B | `bufio.ReadRune` (UTF-8) |
| Level 1 | 1.49s | 1.7x | 21.4B | 32KB Buffer + ASCII State Machine |
| Level 2 | 0.31s | 8.4x | 24.4B | `ReadAt` + Multi-core Concurrency |
| **Level 3** | **0.24s** | **11.2x** | **11.7B** | **Zero-Copy Memory Mapping (mmap)** |
| Standard `wc` | 2.70s | 1.0x | 30.3B | Sequential C Implementation |

## System Specifications

The benchmarks above were conducted on the following hardware:

*   **Machine:** Lenovo IdeaPad Gaming 3 15ACH6
*   **CPU:** AMD Ryzen 5 5600H (6 Cores, 12 Threads, Zen 3 Architecture)
*   **RAM:** 16GB DDR4-3200 MT/s Dual Channel (Hynix + Crucial)
*   **Storage:** KIOXIA EXCERIA PLUS G3 NVMe SSD (PCIe Gen4 x4)
*   **OS:** Arch Linux (Kernel: 6.19.6-arch1-1)
*   **Compiler:** Go 1.24+

## Reproducing Results

To verify these performance claims, download the 1.1GB reference file used in the benchmarks:
[Download lots_of_words.txt (1.1GB)](https://www.dropbox.com/scl/fi/bknr56grghsbpr7a2c5bs/lots_of_words.txt?rlkey=w9znbedr3rhzrck38yb2de46r&st=wk3h2uyb&dl=1)

## Implementation Methodology

### Level 1: Fast Loop & ASCII Optimization
Switched from `ReadRune` to raw byte processing with a 32KB buffer to minimize syscall frequency. Optimized whitespace detection using an ASCII range check (`c <= 32`). This approach remains UTF-8 safe as all whitespace characters in UTF-8 are represented by single bytes below 128, while multi-byte sequences always utilize bytes above 127.

### Level 2: Multi-Worker Concurrency & SIMD
1. **Parallelization:** Partitioned single-file processing across all available CPU cores using `os.File.ReadAt`. Implemented a boundary-peeking algorithm (`ReadAt(pos-1)`) to ensure word fragments across chunks were counted exactly once.
2. **SIMD Line Counting:** Integrated `bytes.Count` for newline detection, leveraging assembly-optimized instructions (AVX2/AVX-512).
3. **Lookup Tables:** Replaced branchy logic with a 256-byte pre-computed lookup table for O(1) whitespace detection.

### Level 3: Zero-Copy Memory Mapping (mmap)
Mapped the file directly into the process's virtual address space via `syscall.Mmap`. This eliminated kernel-to-user space memory copying, reducing total CPU instruction count by over 50% compared to Level 2.

## Notice on Generative AI Usage
Some parts of this codebase were developed with the assistance of Generative AI tools. While AI has been valuable for refining documentation and implementing specific, complex optimizations (such as memory mapping and SIMD integration), it does not replace original research or architectural design. Every AI-assisted implementation has been strictly reviewed and verified for accuracy.

Note: All development up to and including the concurrency branch was performed entirely without AI assistance.

## Technical Notes
*   **Binary:** `go build -o counter`
*   **Benchmarks:** Managed via `bench.sh` using `hyperfine` and `perf stat`.
*   **Safety:** Implements a fallback to `ReadAt` logic if `mmap` is unsupported by the filesystem.
