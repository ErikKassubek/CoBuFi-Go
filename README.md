# AdvocateGo
> [!NOTE]
> This program is a work in progress and may result in incorrect or incomplete results.

## What is AdvocateGo
AdvocateGo is an analysis tool for Go programs.
It detects concurrency bugs and gives  diagnostic insight.
This is achieved through `happens-before-relation` and `vector-clocks`

Furthermore it is also able to produce traces which can be fed back into the program in order to experience the predicted bug.

AdvocateGo tries to detect the following situations:
- A01: Send on closed channel
- A02: Receive on closed channel
- A03: Close on closed channel
- A04: Concurrent recv
- A05: Select case without partner
- P01: Possible send on closed channel
- P02: Possible receive on closed channel
- P03: Possible negative waitgroup counter
- P04: Possible unlock of not locked mutex
- L00: Leak on routine without blocking element
- L01: Leak on unbuffered channel with possible partner
- L02: Leak on unbuffered channel without possible partner
- L03: Leak on buffered channel with possible partner
- L04: Leak on buffered channel without possible partner
- L05: Leak on nil channel
- L06: Leak on select with possible partner
- L07: Leak on select without possible partner
- L08: Leak on mutex
- L09: Leak on waitgroup
- L10: Leak on cond

A more in detail explanation of how it works can be found [here](./doc/Analysis.md).
## Usage
![Flowchart of AdvocateGoProcess](doc/img/architecture.png "Architecture")


### Preparation
Before Advocate can be used, it must first be build.

First you need to build the [analyzer](https://github.com/ErikKassubek/ADVOCATE/tree/main/analyzer)
and if you want to use it the [toolchain](https://github.com/ErikKassubek/ADVOCATE/tree/main/toolchain).
If you do not wish to use the script but to run the step manually, you do not need to build the toolchain.
The two programs are go programs and can just be build using `go build`.

Additionally, the modified go runtime must be build. The runtime can be found in [go-patch](https://github.com/ErikKassubek/ADVOCATE/tree/main/go-patch).
To build it run the
```shell
./src/make.bash
```
or
```shell
./src/make.bat
```
script. This will create a go executable in the `bin` directory.


If you do not wish to use the toolchain script (see below), you must set your
`GOROOT` environment variable to this new runtime:
```shell
export GOROOT=$HOME/ADVOCATE/go-patch/
```
If you use the toolchain script, this will be done automatically.


