genTrace="../genTrace/genTrace.bash"
analyzer="../../analyzer/analyzer"
githubUrl=$1
git clone $githubUrl
cd $(basename "$githubUrl" .git)
fileToExecute=$(find . -name "*.go" -exec grep -q "func main()" {} \; -print -quit)
fileToExecute=$(echo "$fileToExecute" | sed 's|^\./||')
$genTrace $fileToExecute
echo "Run Analysis"
$analyzer -t adadvocateTrace

