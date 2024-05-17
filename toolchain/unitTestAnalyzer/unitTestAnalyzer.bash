#!/bin/bash
if [ -z "$1" ]; then
    echo "Usage: $0 <directory>"
    exit 1
fi


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
      pathToOverHeaderInserter="$2"
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
    *)
      shift
      ;;
  esac
done

if [ -z "$pathToPatchedGoRuntime" ] || [ -z "$pathToGoRoot" ] || [ -z "$pathToOverHeaderInserter" ] || [ -z "$pathToOverheadRemover" ] || [-z "$dir" ];then
  echo "Usage: $0 -p <patched-go-runtime> -g <go-root> -i <overhead-inserter> -r <overhead-remover> -f <folder>"
  exit 1
fi




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
        $pathToPatchedGoRuntime run "$pathToOverheadInserter" -f "$file" -t "$test_func"
        echo "go test -count=1 -run=$test_func $package_path"
        $pathToPatchedGoRuntime test -count=1 -run="$test_func" "$package_path"
        $pathToPatchedGoRuntime run "$pathToOverheadRemover" -f "$file" -t "$test_func"
        packageName=$(basename "$package_path")
        fileName=$(basename "$file")
        mkdir -p "advocateResults/$packageName/$fileName/$test_func"
        mv "$package_path/advocateTrace" "advocateResults/$packageName/$fileName/$test_func/advocateTrace"
    done
done