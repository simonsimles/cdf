package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		return
	}
	var arg = os.Args[1]

    pwd, _ := os.Getwd()
	var state = InitDirectoryWalk(arg)
	result := Run(state)
    os.Chdir(pwd)

	if len(os.Args) == 4 && os.Args[2] == "-f" {
        filename := os.Args[3]
        os.WriteFile(filename, []byte(result), 0666)
	} else {
		fmt.Print(result)
	}
}
