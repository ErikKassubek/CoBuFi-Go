package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	names := []string{
		// "GoBench",
		// "canonicalTests",
		"prometheus",
		"etcd",
		"kubernetes",
		"go-ethereum",
		"connect",
		"argo-cd",
		"moby",
	}

	mainPath := "~/Uni/HiWi/ADVOCATE/examples/"

	if err := os.Chdir("../../toolchain"); err != nil {
		fmt.Printf("Failed to change directory: %v\n", err)
		return
	} else {
		wd, _ := os.Getwd()
		fmt.Printf("Changed working directory to: %s\n\n", wd)
	}

	analysisTimeout := "3600" // 3600s = 1h

	for _, name := range names {
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

}
