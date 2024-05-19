#!/bin/bash
if [ -z "$1" ]; then
    echo "Usage: $0 <directory>"
    exit 1
fi
dir_path=$1
directories=$(find "$dir_path" -type d -name "advocateTrace")
num_dirs=$(echo "$directories" | wc -l)
current_dir=1
for dir in $directories; do
    # print progress
    echo "Processing directory $current_dir of $num_dirs"
    echo "Processing directory: $dir"
    current_dir=$((current_dir+1))
done