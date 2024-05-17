#!/bin/bash
pathToPatchedGoRuntime="../../go-patch/bin/go"
pathToGoRoot="../../go-patch/"
pathToOverHeaderInserter="../overHeadInserter/inserter"
pathToOverheadRemover="../overHeadRemover/remover"

echo "Running full workflow on $1"
echo "Step 0: Remove Overhead just in case"
$pathToOverheadRemover $1
echo "Step 1: Add Overhead to file"
$pathToOverHeaderInserter -f $1
echo "Step 2: Run with patched go runetime"
echo "Step 2.1: save current go root and set adjusted goroot"
export $pathToGoRoot
echo "Step 2.2: run program"
$pathToPatchedGoRuntime run $1
echo "Step 3: Analyze Trace"
echo "Step 3.1: Unset goroot"
unset GOROOT
echo "Step 3.3: Remove Overhead"
$pathToOverheadRemover $1
exit 0