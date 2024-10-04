package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func main() {
	names := []string{
		// "canonicalTests",
		"GoBench",
		"connect",
		"argo-cd",
		"kubernetes",
		"moby",
		"go-ethereum",
		"prometheus",
		// "etcd",
	}

	mainPath := "~/Uni/HiWi/ADVOCATE/examples/"

	if err := os.Chdir("../../toolchain"); err != nil {
		fmt.Printf("Failed to change directory: %v\n", err)
		return
	}

	wd, _ := os.Getwd()
	fmt.Printf("Changed working directory to: %s\n\n", wd)

	analysisTimeout := "1800" // 1800s = 0.5h
	maxProgAtSameTime := 2

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxProgAtSameTime)

	for _, name := range names {
		wg.Add(1)
		go runProg(name, mainPath, analysisTimeout, &wg, sem)
	}

	wg.Wait()

}

func runProg(name, mainPath, analysisTimeout string, wg *sync.WaitGroup, sem chan struct{}) {

	defer wg.Done()
	defer func() { <-sem }()
	sem <- struct{}{}

	fmt.Println("\n\nRun prog ", name)
	path := filepath.Join(mainPath, name)

	cmd := exec.Command("./tool", "test", "-a", "~/Uni/HiWi/ADVOCATE", "-f", path, "-m", "-s", "-t", "-N", name, "-T", analysisTimeout)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println(cmd.String())
	err := cmd.Run()

	if err != nil {
		fmt.Println("Could not run ", name, err)
		file, err2 := os.OpenFile("failed.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err2 != nil {
			fmt.Println("Error opening file:", err2)
			return
		}

		_, err2 = file.WriteString("# " + name)
		if err2 != nil {
			fmt.Println("Error writing string:", err2)
			file.Close()
			return
		}

		_, err2 = file.WriteString(err.Error())
		if err2 != nil {
			fmt.Println("Error writing string:", err2)
			file.Close()
			return
		}
		file.Close()
	}

	fmt.Println("Finished running", name)
}
