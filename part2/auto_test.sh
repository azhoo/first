#!/bin/bash

# Directories
CODES_DIR="./c_codes"       # Directory with C code files
TESTS_DIR="./test_cases"    # Directory with test case inputs
REFS_DIR="./ref_outputs"    # Directory with expected outputs
LOG_FILE="./results.log"    # File to log results

# Clear log file
> $LOG_FILE

for c_file in "$CODES_DIR"/*.c; do
    base_name=$(basename "$c_file" .c)
    echo "Testing $base_name.c..." >> $LOG_FILE

    # Compile the C file
    gcc "$c_file" -o "$CODES_DIR/$base_name.out" 
    if [ $? -ne 0 ]; then
        echo "  Compilation failed." >> $LOG_FILE
        continue
    fi
    total=0
    passed=0
    for test_file in "$TESTS_DIR/$base_name"_input*.txt; do
       # Step 1: Create the reference file path
        test_base_name=$(basename "$test_file" _input*.txt)
       
        ref_file="${REFS_DIR}/${test_base_name}_output.txt"

        # Step 2: Run the compiled program with input and save output
        "$CODES_DIR/$base_name.out" < "$test_file" > temp_output.txt


        if diff -q temp_output.txt "$ref_file";  then
            ((passed++))
        fi
            ((total++))
    done
    
    percentage=$((100 * passed / total))
    echo "  Passed $passed/$total tests ($percentage%)." >> $LOG_FILE

    # Clean up
    rm -f "$CODES_DIR/$base_name.out" temp_output.txt
done

echo "Testing complete. Results saved in $LOG_FILE."
