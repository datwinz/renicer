package main

import (
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

    fmt.Println(
        formatLines(proccessesSlice, "pid"),
        formatLines(proccessesSlice, "ni"),
        formatLines(proccessesSlice, "comm"))

    // Layout: lists in VBox, on left Border, with a Centered HBox on right screen.
    // Grid is probably better than VBox and HBox, because it reserves a minimum space.
    // Right screen can also be a Form (https://youtu.be/LWn1403gY9E?t=717)
    // All items are automatically rendered at the minimum size.
    // Combining layouts is explained here: https://youtu.be/LWn1403gY9E?t=1061
	a := app.New()
	w := a.NewWindow("Renicer")

    content := widget.NewList(
        func() (int) {
            return len(proccessesSlice)
        },
        func() (fyne.CanvasObject) {
            // This is the standard name for the items in the list
            return widget.NewLabel("Process")
        },
        func(i widget.ListItemID, o fyne.CanvasObject) {
            o.(*widget.Label).SetText(
                formatLines(proccessesSlice, "pid")[i],
            )
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

func formatLines(processes []string, outputfield string) (formatted []string) {
    var allLines []string
    for i := 0; i < len(processes)-1; i++ {
        f := strings.Fields(processes[i])
        switch outputfield {
        case "pid":
            pid := f[0]
            allLines = append(allLines, []string{pid}...)
        case "ni":
            ni := f[1]
            allLines = append(allLines, []string{ni}...)
        default:
            if strings.Contains(f[2], "/") {
                s := strings.Split(f[2], "/")
                comm := s[len(s)-1]
                allLines = append(allLines, []string{comm}...)
            } else {
                comm := f[2]
                allLines = append(allLines, []string{comm}...)
            }
        }
    }
    return allLines
}

// Do 'ps ax -o pid,ni,comm' and make sort by name, procces number nice value
// I.i.r.c. Linux has different words for the options, but if I look it up the only difference
// is that in Linux you can also use cmd instead of comm

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
