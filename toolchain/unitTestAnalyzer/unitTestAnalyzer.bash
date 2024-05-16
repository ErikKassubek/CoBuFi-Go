#!/bin/bash
if [ -z "$1" ]; then
    echo "Usage: $0 <directory>"
    exit 1
fi

dir="$1"
pathToUnitTestOverheadInserter="../unitTestOverheadInserter/unitTestOverheadInserter.go"
pathToUnitTestOverheadRemover="../unitTestOverheadRemover/unitTestOverheadRemover.go"
pathToPatchedGoRuntime="../../go-patch/bin/go"
pathToGoRoot="../../go-patch/"
cd "$dir"
#make folder to store the result
rm -r "advocateResults"
mkdir "advocateResults"

export $pathToGoRoot
test_files=$(find "$dir" -name "*_test.go")
for file in $test_files; do
    package_path=$(dirname "$file")
    test_functions=$(grep -oE "[a-zA-Z0-9_]+ *Test[a-zA-Z0-9_]*" $file | sed 's/ *\(t *\*testing\.T\)//' | sed 's/func //')
    for test_func in $test_functions; do
        echo "Running test: $test_func" in "$file of package $package_path"
        $pathToPatchedGoRuntime run "$pathToUnitTestOverheadInserter" -f "$file" -t "$test_func"
        echo "go test -count=1 -run=$test_func $package_path"
        $pathToPatchedGoRuntime test -count=1 -run="$test_func" "$package_path"
        $pathToPatchedGoRuntime run "$pathToUnitTestOverheadRemover" -f "$file" -t "$test_func"
        packageName=$(basename "$package_path")
        fileName=$(basename "$file")
        mkdir -p "advocateResults/$packageName/$fileName/$test_func"
        mv "$package_path/advocateTrace" "advocateResults/$packageName/$fileName/$test_func/advocateTrace"
    done
done