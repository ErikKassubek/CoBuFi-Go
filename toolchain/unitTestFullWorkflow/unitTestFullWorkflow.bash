
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
    -t|--test-name)
      testName="$2"
      shift
      shift
      ;;
    -package)
      package="$2"
      shift
      shift
      ;;
    -tf|--test-file)
      file="$2"
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
fi
if [ -z "$pathToGoRoot" ]; then
  echo "Path to go root is empty"
fi
if [ -z "$pathToOverheadInserter" ]; then
  echo "Path to overhead inserter is empty"
fi
if [ -z "$pathToOverheadRemover" ]; then
  echo "Path to overhead remover is empty"
fi
if [ -z "$pathToAnalyzer" ]; then
  echo "Path to analyzer is empty"
fi
if [ -z "$dir" ]; then
  echo "Directory is empty"
fi
if [ -z "$testName" ]; then
  echo "Test name is empty"
fi
if [ -z "$package" ]; then
  echo "Package is empty"
fi
if [ -z "$file" ]; then
  echo "Test file is empty"
fi

if [ -z "$pathToPatchedGoRuntime" ] || [ -z "$pathToGoRoot" ] || [ -z "$pathToOverheadInserter" ] || [ -z "$pathToOverheadRemover" ] || [ -z "$pathToAnalyzer" ] || [ -z "$dir" ] || [ -z "$testName" ] || [ -z "$package" ] || [ -z "$file" ]; then
  echo "Usage: $0 -patch|--patched-go-runtime <pathToPatchedGoRuntime> -g|--go-root <pathToGoRoot> -i|--overhead-inserter <pathToOverheadInserter> -r|--overhead-remover <pathToOverheadRemover> -a|--analyzer <pathToAnalyzer> -f|--folder <directory> -t|--test-name <testName> -package <package> -tf|--test-file <testFile>"
  exit 1
fi




cd "$dir"
echo  "In directory: $dir"
export GOROOT=$pathToGoRoot
echo "Goroot exported"
#Remove Overhead just in case
echo "Remove Overhead just in case"
#echo "$pathToOverheadRemover -f $file -t $testName"
echo "Remover: $pathToOverheadRemover"
echo "File: $file"
echo "Test Name: $testName"
$pathToOverheadRemover -f $file -t $testName
#Add Overhead
#echo "Add Overhead"
#echo "$pathToOverheadInserter -f $file -t $testName"
#$pathToOverheadInserter -f $package/$file -t $testName
##Run test
#echo "$pathToPatchedGoRuntime test -count=1 -run=$testName $package"
#"$pathToPatchedGoRuntime test -count=1 -run=$testName $package"
##Remove Overhead
#echo "$pathToOverheadRemover -f $file -t $testName"
#"$pathToOverheadRemover -f $package/$file -t $testName"
#Run Analyzer
#Loop through every rewritten traces
## Remove Overhead just in case
## Apply reorder overhead
## Run test
## Remove reorder overhead
unset GOROOT