# ginti

**ginti (गिनती)**, meaning "to count" or "counter" in Hindi, is a fast, concurrent `wc` clone written in Go.

## Overview

`ginti` is a command-line utility for counting lines, words, and bytes. This branch specifically implements concurrent file processing to maximize throughput when analyzing multiple files.

## Features

*   **Concurrent Processing**: Spawns a goroutine per file, utilizing a wait group and channels to aggregate results safely and efficiently.
*   **Standard Flags**: Supports `-c` (bytes), `-l` (lines), and `-w` (words) options.
*   **Standard Input**: Falls back to `stdin` when no file arguments are provided.
*   **Clean Output**: Uses `text/tabwriter` for aligned, readable output. Optional `-header` flag available.

## Installation

Install using `go install`:

```sh
go install github.com/HimuCodes/ginti@latest
```

Alternatively, clone the repository and build from source:

```sh
git clone https://github.com/HimuCodes/ginti.git
cd ginti
go build -o ginti
```

## Usage

Execute with files:

```sh
./ginti -l -w file1.txt file2.txt
```

Execute via standard input:

```sh
cat file.txt | ./ginti -c
```

## Concurrency Model

When processing multiple files, `ginti` avoids blocking on sequential I/O. The workload distribution works as follows:

1.  A goroutine is initialized for every file argument.
2.  File parsing (`GetCounts`) occurs concurrently across all files.
3.  Results and potential errors are sent over a unified channel back to the main routine.
4.  The main routine collects the results, maintains the original argument order, and computes the aggregate totals.

This architecture significantly reduces wall-clock time for batch file processing.
