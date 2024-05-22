
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    -patch|--patched-go-runtime)
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
    -a|--analyzer)
      pathToAnalyzer="$2"
      shift
      shift
      ;;
    -f|--folder)
      dir="$2"
      shift
      shift
      ;;
    -t|--testName)
      testName="$2"
      shift
      shift
      ;;
    -package|--packageName)
      packageName="$2"
      shift
      shift
      ;;
    -h|--help)
      echo "Usage: $0 -p <patched-go-runtime> -g <go-root> -i <overhead-inserter> -r <overhead-remover> -a <analyzer> -f <folder> -t <testName> -package <packageName>"
      exit 0
      ;;
    *)
      shift
      ;;
  esac
done

if [ -z "$pathToPatchedGoRuntime" ] || [ -z "$pathToAnalyzer" ] || [ -z "$pathToGoRoot" ] || [ -z "$pathToOverheadInserter" ] || [ -z "$pathToOverheadRemover" ] || [ -z "$dir" ] || [ -z $testName ] || [ -z $packageName];then
  echo "Usage: $0 -p <patched-go-runtime> -g <go-root> -i <overhead-inserter> -r <overhead-remover> -a <analyzer> -f <folder>"
  exit 1
fi




cd "$dir"
echo  "In directory: $dir"
export GOROOT=$pathToGoRoot
echo "Goroot exported"
#Remove Overhead just in case
#Add Overhead
#Run test
#Remove Overhead
#Run Analyzer
#Loop through every rewritten traces
## Remove Overhead just in case
## Apply reorder overhead
## Run test
## Remove reorder overhead
echo "$pathToOverheadInserter -f $file -t $test_func"
$pathToOverheadInserter -f $file -t $test_func
echo "$pathToPatchedGoRuntime test -count=1 -run=$test_func $package_path"
$pathToPatchedGoRuntime test -count=1 -run="$test_func" "$package_path"
echo "$pathToOverheadRemover -f $file -t $test_func"
$pathToOverheadRemover -f "$file" -t "$test_func"
packageName=$(basename "$package_path")
fileName=$(basename "$file")
unset GOROOT