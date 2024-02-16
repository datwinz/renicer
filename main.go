package main

import (
	"fmt"
    "log"
	"os/exec"
)    

func main() {
    path, err := exec.LookPath("ps")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("ps is available at %s\n", path)
}
