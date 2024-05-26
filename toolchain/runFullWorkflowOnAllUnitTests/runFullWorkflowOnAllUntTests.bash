#!/bin/bash
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    -p|--patched-go-runtime)
      pathToPatchedGoRuntime="$2"
      shift
      shift
      ;;
    -g|--go-root)
      pathToGoRoot="$2"
      shift
      shift
      ;;
    -i|--overhead-inserter)
      pathToOverheadInserter="$2"
      shift
      shift
      ;;
    -r|--overhead-remover)
      pathToOverheadRemover="$2"
      shift
      shift
      ;;
    -f|--folder)
      dir="$2"
      shift
      shift
      ;;
    -a|--analyzer)
      pathToAnalyzer="$2"
      shift
      shift
      ;;
    *)
      shift
      ;;
  esac
done


if [ -z "$pathToPatchedGoRuntime" ]; then
  echo "Path to patched go runtime is empty"
  exit 1
fi
if [ -z "$pathToGoRoot" ]; then
  echo "Path to go root is empty"
  exit 1
fi
if [ -z "$pathToOverheadInserter" ]; then
  echo "Path to overhead inserter is empty"
  exit 1
fi
if [ -z "$pathToOverheadRemover" ]; then
  echo "Path to overhead remover is empty"
  exit 1
fi
if [ -z "$dir" ]; then
  echo "Directory is empty"
  exit 1
fi
if [ -z "$pathToAnalyzer" ]; then
  echo "Path to analyzer is empty"
  exit 1
fi




cd "$dir"
echo  "In directory: $dir"

test_files=$(find "$dir" -name "*_test.go")
total_files=$(echo "$test_files" | wc -l)
current_file=1
#echo "Test files: $test_files"
for file in $test_files; do
    echo "Progress: $current_file/$total_files"
    echo "Processing file: $file"
    package_path=$(dirname "$file")
    test_functions=$(grep -oE "[a-zA-Z0-9_]+ *Test[a-zA-Z0-9_]*" $file | sed 's/ *\(t *\*testing\.T\)//' | sed 's/func //')
    for test_func in $test_functions; do
        packageName=$(basename "$package_path")
        fileName=$(basename "$file")
        #runfullworkflow for single test pass all 
        echo "Running full workflow for test: $test_func in package: $package_path in file: $file"
        adjustedPackagePath=$(echo "$package_path" | sed "s|$dir||g")
        /home/mario/Desktop/thesis/ADVOCATE/toolchain/unitTestFullWorkflow/unitTestFullWorkflow.bash -p $pathToPatchedGoRuntime -g $pathToGoRoot -i $pathToOverheadInserter -r $pathToOverheadRemover -package $adjustedPackagePath -f $dir -tf $file -a $pathToAnalyzer -t $test_func

    done
    current_file=$((current_file+1))
done
