# Explanation
# Input
# Output
# Usage
# Common Problems
This tool requires a go.mod at the project root otherwise the tests won't run.
This is the case for some repositories (eg Moby).
In this case you need to manually add a go.mod via `go mod init` in the project root and call the program with the flag `-m true` like so
```sh
./runFullWorkflowOnAllUnitTests -a <path-to-advocate> -f <path-to-folder> -m true
```