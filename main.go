package main

import (
    "log"
    "os/exec"
    "strings"
    "fmt"
)    

func main() {
    psPath := processPaths("ps")
    //renicePath := processPaths("renice")
    proccessesSlice := findProcesses(psPath)
    fmt.Println(proccessesSlice)
}

func processPaths(processName string) (path string) {
    path, err := exec.LookPath(processName)
    if err != nil {
        log.Fatal(err)
    }
    return path
}

func findProcesses(psPath string) (processes []string) {
    psCmd := exec.Command(psPath, "ax", "-o pid,ni,comm")
    var outAll strings.Builder
    psCmd.Stdout = &outAll
    err := psCmd.Run()
    if err != nil {
        log.Fatal(err)
    }

    outSingle := strings.Split(outAll.String(), "\n")
    //fmt.Println(outSingle)
    return outSingle
}

// Do 'ps ax -o pid,ni,comm' and make sort by name, procces number nice value

// Put it in window something like this:
// __________________________________________________
// | pid | Process  | Ni |                          |
// __________________________________________________
// | 1   | init     | 0  |                          |
// | 2   | process2 | 0  |      process5            |
// | 3   | process3 | 0  |      Old value: 0        |
// | 4   | process4 | 0  |      New value : -20     |
// | 5   | process5 | 0  |                          |
// | 6   | process6 | 0  |      Save                |
// --------------------------------------------------

// Add search

// Show man pages or something for processes on double click
