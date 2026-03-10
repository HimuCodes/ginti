#!/usr/bin/env bash

# Define our targets
BIN="./counter"
STD_BIN="wc"
FILE="lots_of_words.txt"
CHUNK_DIR="test_chunks"

echo "====================================================="
echo "🚀 Starting Benchmark: $BIN vs $STD_BIN"
echo "📂 Target File: $FILE"
echo "====================================================="

# 1. Sanity Checks
if [ ! -f "$BIN" ]; then
    echo "❌ Error: Binary '$BIN' not found. Please run 'go build -o counter' first."
    exit 1
fi

if [ ! -f "$FILE" ]; then
    echo "❌ Error: File '$FILE' not found."
    exit 1
fi

# 2. Dependency Checks
for cmd in hyperfine perf /usr/bin/time split; do
    if ! command -v $cmd &> /dev/null; then
        echo "❌ Error: '$cmd' is not installed."
        echo "   Install missing tools on Arch with: sudo pacman -S hyperfine perf time"
        exit 1
    fi
done

echo -e "\n⏳ Phase 1: Speed Test (Single Large File - Sequential)"
echo "-----------------------------------------------------"
hyperfine --warmup 2 "$BIN $FILE" "$STD_BIN $FILE"

echo -e "\n⏳ Phase 2: Speed Test (Multiple Files - Concurrency)"
echo "-----------------------------------------------------"
echo "Splitting $FILE into 100MB chunks to test goroutine scaling..."
mkdir -p "$CHUNK_DIR"
split -b 100M "$FILE" "$CHUNK_DIR/chunk_"

hyperfine --warmup 2 "$BIN $CHUNK_DIR/chunk_"* "$STD_BIN $CHUNK_DIR/chunk_"*

echo "Cleaning up chunks..."
rm -rf "$CHUNK_DIR"

echo -e "\n🧠 Phase 3: Hardware & CPU Cycles (perf stat)"
echo "-----------------------------------------------------"
echo "Running standard wc..."
perf stat -e cpu-cycles,instructions,branches,branch-misses $STD_BIN $FILE > /dev/null

echo -e "\nRunning your Go counter..."
perf stat -e cpu-cycles,instructions,branches,branch-misses $BIN $FILE > /dev/null

echo -e "\n💾 Phase 4: Peak Memory Usage & System Time (/usr/bin/time)"
echo "-----------------------------------------------------"
echo "Running standard wc..."
/usr/bin/time -v $STD_BIN $FILE > /dev/null

echo -e "\nRunning your Go counter..."
/usr/bin/time -v $BIN $FILE > /dev/null

echo -e "\n✅ Benchmarking Complete!"
