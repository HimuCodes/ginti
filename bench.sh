#!/usr/bin/env bash

# Ginti Performance Analysis Suite: ginti vs wc
# Focus: Latency, Throughput, Instruction Efficiency, and Cache Sensitivity.

BIN="./counter"
STD_BIN="wc"
TARGET_FILE="lots_of_words.txt"
CHUNK_DIR="test_chunks"
SMALL_FILE="small_test.txt"

# Colors for professional output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=====================================================${NC}"
echo -e "${BLUE}Ginti Performance Analysis Suite${NC}"
echo -e "${BLUE}=====================================================${NC}"

# 1. Verification & Environment Setup
if [[ ! -f "$BIN" ]]; then
    echo -e "${RED}Error: Binary '$BIN' not found. Rebuilding...${NC}"
    go build -o counter || exit 1
fi

if [[ ! -f "$TARGET_FILE" ]]; then
    echo -e "${RED}Error: Reference file '$TARGET_FILE' not found.${NC}"
    exit 1
fi

# Create a small file for boundary testing (1KB)
head -c 1024 /dev/urandom > "$SMALL_FILE"

# 2. Dependency Validation
for cmd in hyperfine perf /usr/bin/time; do
    if ! command -v "$cmd" &> /dev/null; then
        echo -e "${RED}Required tool '$cmd' is missing.${NC}"
        exit 1
    fi
done

# 3. Cache Strategy Setup
# Note: Clearing cache requires sudo. If not available, we rely on warm cache benchmarking.
CLEAR_CACHE="sync; echo 3 | sudo tee /proc/sys/vm/drop_caches > /dev/null"
if sudo -n true 2>/dev/null; then
    CACHE_MODE="Cold & Warm"
else
    CACHE_MODE="Warm Only (Sudo required for cold cache)"
fi

echo -e "Mode: ${GREEN}$CACHE_MODE${NC}"
echo -e "Target: ${GREEN}$TARGET_FILE (~1.1GB)${NC}"

# --- TEST SUITE 1: COLD CACHE (The 'Honest' Disk Test) ---
if [[ "$CACHE_MODE" == "Cold & Warm" ]]; then
    echo -e "\n${BLUE}Phase 1: Cold Cache Latency (Disk I/O Bound)${NC}"
    echo "-----------------------------------------------------"
    hyperfine --prepare "$CLEAR_CACHE" --runs 3 \
        "$BIN $TARGET_FILE" \
        "$STD_BIN $TARGET_FILE"
fi

# --- TEST SUITE 2: WARM CACHE (The 'Raw Throughput' Test) ---
echo -e "\n${BLUE}Phase 2: Warm Cache Throughput (CPU & RAM Bound)${NC}"
echo "-----------------------------------------------------"
hyperfine --warmup 5 --runs 10 \
    --export-markdown bench_results.md \
    "$BIN $TARGET_FILE" \
    "$STD_BIN $TARGET_FILE"

# --- TEST SUITE 3: CONCURRENCY SCALING ---
echo -e "\n${BLUE}Phase 3: Massively Parallel Workload (100MB Chunks)${NC}"
echo "-----------------------------------------------------"
mkdir -p "$CHUNK_DIR"
split -b 100M "$TARGET_FILE" "$CHUNK_DIR/chunk_"
hyperfine --warmup 3 --runs 5 \
    "$BIN $CHUNK_DIR/chunk_*" \
    "$STD_BIN $CHUNK_DIR/chunk_*"
rm -rf "$CHUNK_DIR"

# --- TEST SUITE 4: SMALL FILE OVERHEAD ---
echo -e "\n${BLUE}Phase 4: Small File Overhead (Threshold Check)${NC}"
echo "-----------------------------------------------------"
hyperfine --warmup 10 --runs 50 \
    "$BIN $SMALL_FILE" \
    "$STD_BIN $SMALL_FILE"
rm "$SMALL_FILE"

# --- TEST SUITE 5: INSTRUCTION ANALYSIS (perf stat) ---
echo -e "\n${BLUE}Phase 5: Hardware Metric Extraction (perf stat)${NC}"
echo "-----------------------------------------------------"
function run_perf() {
    echo -e "\n${GREEN}Metric: $1${NC}"
    perf stat -e cpu-cycles,instructions,branches,branch-misses,L1-dcache-load-misses "$2" "$TARGET_FILE" 2>&1 | \
    grep -E 'cycles|instructions|branches|branch-misses|dcache-load-misses|seconds'
}

run_perf "Standard wc" "$STD_BIN"
run_perf "Optimized ginti" "$BIN"

echo -e "\n${BLUE}=====================================================${NC}"
echo -e "${GREEN}Benchmarking Suite Complete${NC}"
echo -e "Detailed summary exported to bench_results.md"
echo -e "${BLUE}=====================================================${NC}"
