#!/bin/bash
CODE_FILE="b.c"
TEST_CASES_DIR="tests"
LOG_DIR="logs"

mkdir -p "$LOG_DIR"

gcc "$CODE_FILE" -o ali

for test_case in "$TEST_CASES_DIR"/*.txt; do
    TEST_NAME=$(basename "$test_case" .txt)
    ./ali < "$test_case" > "$LOG_DIR/$TEST_NAME.log"
done



