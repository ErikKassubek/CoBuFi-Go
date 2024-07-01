while [[ $# -gt 0 ]]; do
	key="$1"
	case $key in
	-a | --advocate)
		pathToAdvocate="$2"
		shift
		shift
		;;
	-f | --folder)
		dir="$2"
		shift
		shift
		;;
	-m | --modulemode)
		modulemode="$2"
		shift
		shift
		;;
	-t | --test-name)
		testName="$2"
		shift
		shift
		;;
	-p | --package)
		package="$2"
		shift
		shift
		;;
	-tf | --test-file)
		file="$2"
		shift
		shift
		;;
	*)
		shift
		;;
	esac
done

pathToPatchedGoRuntime="$pathToAdvocate/go-patch/bin/go"
pathToGoRoot="$pathToAdvocate/go-patch"
pathToOverheadInserter="$pathToAdvocate/toolchain/unitTestOverheadInserter/unitTestOverheadInserter"
pathToOverheadRemover="$pathToAdvocate/toolchain/unitTestOverheadRemover/unitTestOverheadRemover"
pathToAnalyzer="$pathToAdvocate/analyzer/analyzer"

if [ -z "$pathToAdvocate" ]; then
	echo "Path to advocate is empty"
	exit 1
fi
if [ -z "$dir" ]; then
	echo "Directory is empty"
	exit 1
fi
if [ -z "$testName" ]; then
	echo "Test name is empty"
	exit 1
fi
if [ -z "$package" ]; then
	echo "Package is empty"
	exit 1
fi
if [ -z "$file" ]; then
	echo "Test file is empty"
	exit 1
fi

cd "$dir"
echo "In directory: $dir"
export GOROOT=$pathToGoRoot
echo "Goroot exported"
touch advocateCommand.log
echo $file >>advocateCommand.log
echo $testName >>advocateCommand.log
echo "Remove Overhead just in case" 
echo "$pathToOverheadRemover -f $file -t $testName" >>advocateCommand.log
$pathToOverheadRemover -f $file -t $testName 
echo "Add Overhead"
echo "$pathToOverheadInserter -f $file -t $testName" >>advocateCommand.log
$pathToOverheadInserter -f $file -t $testName >>advocateCommand.log
if [ $? -ne 0 ]; then
	echo "Error in adding overhead"
	exit 1
fi
echo "Run test"
if [ "$modulemode" == "true" ]; then
	echo "$pathToPatchedGoRuntime test -count=1 -run=$testName -mod=mod ./$package" >>advocateCommand.log
	$pathToPatchedGoRuntime test -count=1 -run=$testName -mod=mod "./$package"
else
	echo "$pathToPatchedGoRuntime test -count=1 -run=$testName ./$package" >>advocateCommand.log
	$pathToPatchedGoRuntime test -count=1 -run=$testName "./$package"
fi
if [ $? -ne 0 ]; then
	echo "Remove Overhead"
	$pathToOverheadRemover -f $file -t $testName
	echo "Error in running test, therefor overhead removed and full workflow stopped."
	exit 1
fi
echo "Remove Overhead"
echo "$pathToOverheadRemover -f $file -t $testName" >>advocateCommand.log
$pathToOverheadRemover -f $file -t $testName
echo "$pathToAnalyzer -t $dir/$package/advocateTrace" >>advocateCommand.log
$pathToAnalyzer -t "$dir/$package/advocateTrace"
rewritten_traces=$(find "./$package" -type d -name "rewritten_trace*")
for trace in $rewritten_traces; do
	rtracenum=$(echo $trace | grep -o '[0-9]*$')
	echo "Apply reorder overhead"
	echo $pathToOverheadInserter -f $file -t $testName -r true -n "$rtracenum" >>advocateCommand.log
	$pathToOverheadInserter -f $file -t $testName -r true -n "$rtracenum"  >>advocateCommand.log
	if [ "$modulemode" == "true" ]; then
		echo "$pathToPatchedGoRuntime test -count=1 -run=$testName -mod=mod ./$package" >>advocateCommand.log
		$pathToPatchedGoRuntime test -count=1 -run=$testName -mod=mod "./$package" 2>&1 | tee -a "$trace/reorder_output.txt"
	else
		echo "$pathToPatchedGoRuntime test -count=1 -run=$testName ./$package" >>advocateCommand.log
		$pathToPatchedGoRuntime test -count=1 -run=$testName "./$package" 2>&1 | tee -a "$trace/reorder_output.txt"
	fi
	echo "Remove reorder overhead"
	echo "$pathToOverheadRemover -f $file -t $testName" >>advocateCommand.log
	$pathToOverheadRemover -f $file -t $testName
done
lineCount=$(wc -l < "$dir/$package/results_machine.log")
echo "Results machine log line count: $lineCount"
if [ $lineCount -ge 1 ]; then
	echo "Generate bug report"
	#usage ./analyzer -e -t advocateTrace -i 1
	# for every line in results_machine.log, get the line number and run the analyzer
	for i in $(seq 1 $lineCount); do
		echo "$pathToAnalyzer -e -t $dir/$package/advocateTrace -i $i" >>advocateCommand.log
		$pathToAnalyzer -e -t "$dir/$package/advocateTrace" -i $i
	done
	
else
	echo "Results machine log is empty => no need to generate bug report"
fi

unset GOROOT
