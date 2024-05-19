#!/bin/bash
# Check if directory path is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <directory>"
    exit 1
fi
dir_path=$1
directories=$(find "$dir_path" -type d -name "advocateTrace")
for dir in $directories; do
    echo "Directory: $dir"
    for file in "$dir"/*; do
        echo "File: $file"
        cat "$file"
    done
done