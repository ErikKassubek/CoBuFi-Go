# Explanation
# Example
# Common Problems
This tool requires a go.mod at the project root otherwise the tests won't run.
This is the case for some repositories (eg Moby).
In this case you need to manually add a go.mod via `go mod init` in the project root and call the program with the flag `-m true` like so
```sh
./unitTestFullWorkflow.bash -a <path-advocate> -f <path-kubernetes-root> -m <true> -tf <path-kuberbentes-root>/plugin/pkg/admission/deny/admission_test.go -p plugin/pkg/admission/deny -t TestAdmission 
```