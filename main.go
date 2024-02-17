package main

import (
	//"fmt"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)    

func main() {
    psPath := processPaths("ps")
    //renicePath := processPaths("renice")
    proccessesSlice := findProcesses(psPath)

    // This prints in cli how i want it too look in program
    for i := 0; i < len(proccessesSlice)-1; i++ {
        f := strings.Fields(proccessesSlice[i])
        if strings.Contains(f[2], "/") {
            s := strings.Split(f[2], "/")
            pid := f[0]
            ni := f[1]
            comm := s[len(s)-1]
            line := []string{pid, ni, comm}
            fmt.Println(strings.Join(line, " "))
        } else {
            pid := f[0]
            ni := f[1]
            s := f[2]
            line := []string{pid, ni, s}
            fmt.Println(strings.Join(line, " "))
        }
    }

	a := app.New()
	w := a.NewWindow("Renicer")

    // This shows the program window, it uses the function formatProcesses with pretty much
    // the same code as above. But I something is going wrong, it only prints the first line.
    content := widget.NewList(
        func() int {
            return len(proccessesSlice)
        },
        func() fyne.CanvasObject {
            // This is the standard name for the items in the list
            return widget.NewLabel("Process")
        },
        func(i widget.ListItemID, o fyne.CanvasObject) {
            o.(*widget.Label).SetText(formatProcesses(proccessesSlice))
        })
    w.SetContent(content)
	w.ShowAndRun()
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

func formatProcesses(processes []string) (formatted string) {
    for i := 0; i < len(processes)-1; i++ {
        f := strings.Fields(processes[i])
        if strings.Contains(f[2], "/") {
            s := strings.Split(f[2], "/")
            pid := f[0]
            ni := f[1]
            comm := s[len(s)-1]
            line := []string{pid, ni, comm}
            return strings.Join(line, " ")
        } else {
            pid := f[0]
            ni := f[1]
            s := f[2]
            line := []string{pid, ni, s}
            return strings.Join(line, " ")
        }
    }
    return "error"
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
