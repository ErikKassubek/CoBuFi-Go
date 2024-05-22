
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
    -package|--package)
      package="$2"
      shift
      shift
      ;;
    -h|--help)
      echo "Usage: $0 -p <patched-go-runtime> -g <go-root> -i <overhead-inserter> -r <overhead-remover> -a <analyzer> -f <folder> -t <testName> -package <package>"
      exit 0
      ;;
    *)
      shift
      ;;
  esac
done

if [ -z "$pathToPatchedGoRuntime" ] || [ -z "$pathToAnalyzer" ] || [ -z "$pathToGoRoot" ] || [ -z "$pathToOverheadInserter" ] || [ -z "$pathToOverheadRemover" ] || [ -z "$dir" ] || [ -z $testName ] || [ -z $package];then
  echo "Usage: $0 -p <patched-go-runtime> -g <go-root> -i <overhead-inserter> -r <overhead-remover> -a <analyzer> -f <folder>"
  exit 1
fi




cd "$dir"
echo  "In directory: $dir"
export GOROOT=$pathToGoRoot
echo "Goroot exported"
#Remove Overhead just in case
echo "$pathToOverheadRemover -f $package/$file -t $testName"
"$pathToOverheadRemover -f $package/$file -t $testName"
#Add Overhead
echo "$pathToOverheadInserter -f $file -t $testName"
"$pathToOverheadInserter -f $package/$file -t $testName"
#Run test
echo "$pathToPatchedGoRuntime test -count=1 -run=$testName $package"
"$pathToPatchedGoRuntime test -count=1 -run=$testName $package"
#Remove Overhead
echo "$pathToOverheadRemover -f $file -t $testName"
"$pathToOverheadRemover -f $package/$file -t $testName"
#Run Analyzer
#Loop through every rewritten traces
## Remove Overhead just in case
## Apply reorder overhead
## Run test
## Remove reorder overhead
unset GOROOT