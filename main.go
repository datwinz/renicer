package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Layout: lists in VBox, on left Border, with a Centered HBox on right screen.
// Grid is probably better than VBox and HBox, because it reserves a minimum space.
// Right screen can also be a Form (https://youtu.be/LWn1403gY9E?t=717)
// All items are automatically rendered at the minimum size.
// Combining layouts is explained here: https://youtu.be/LWn1403gY9E?t=1061
func mainLayout(
    wholeProcesses *widget.List,
    searchBar *widget.Entry,
    mainWindow *widget.Form,
    ) (*fyne.Container) {
    processes := container.New(layout.NewGridLayout(2), wholeProcesses, mainWindow)
    totalLayout := container.NewBorder(searchBar, nil, nil, nil, processes)
    return totalLayout
}

func main() {
    psPath := processPaths("ps")
    renicePath := processPaths("renice")
    manPath := processPaths("man")

    psOutput := findProcesses(psPath)

    a := app.New()
    w := a.NewWindow("Renicer")

    processListContent := widget.NewList(
        func() (int) {
            return len(psOutput)
        },
        func() (fyne.CanvasObject) {
            // This is the standard name for the items in the list
            return widget.NewLabel("Process")
        },
        func(i widget.ListItemID, o fyne.CanvasObject) {
            o.(*widget.Label).SetText(
                formatWholeLines(psOutput)[i],
            )
        },
    )
    psNameLabel := widget.NewLabel("process")
    psNiLabel := widget.NewLabel("0")
    var psPidValue string
    psNiEntry := widget.NewEntry()
    messageLabel := widget.NewLabel("")
    psSaveButtonFunction := func() {
        // Do some validation (value between -20 and 19
        value := psNiEntry.Text
        valueInt, err := strconv.Atoi(value)
        if err != nil {
            msg := "Can't convert to integer"
            fmt.Println(msg)
            messageLabel.SetText(msg)
        }
        if valueInt >= -20 && valueInt <= 19 {} else {
            msg := "Value should be between -20 and 19"
            fmt.Println(msg)
            messageLabel.SetText(msg)
            return
        }
        fmt.Printf("%q %q %q", psPidValue, value, psNameLabel.Text)
        // This always has exit status 1 for some reason
        if valueInt >=-20 && valueInt < 0 {
            fmt.Println("Spawn polkitd or mac-like window because value")
        }
        reniceCmd := exec.Command(renicePath, value, psPidValue)
        // I have to close the pipe somehow otherwise the program hangs
        ///stderr, err := reniceCmd.StderrPipe()
        ///slurp, _ := io.ReadAll(stderr)
        ///fmt.Printf("%s\n", slurp)
        // This doesnt work, I have to convert to slurp to string
        ///if strings.Contains(slurp, "exit status 1") {
        ///    fmt.Println("Spawn polkitd or mac-like window")
        ///}
        err = reniceCmd.Run()
        if err != nil {
            log.Println(err)
        }
        messageLabel.SetText("")
    }
    psManpageButtonFunction := func () {
        manCmd := exec.Command(manPath, psNameLabel.Text)
        manCmd.Run()
        // Somehow open terminal and show the process, this prints the pid
        // I made a small C script but it doesn't work. So maybe I can spawn a terminal with
        // the pid.
        fmt.Println(manCmd.Process)
    }

    processListContent.OnSelected = func(i widget.ListItemID) {
        j := psOutput[i]
        k := strings.Fields(j)[2]
        if strings.Contains(k, "/") {
            s := strings.Split(k, "/")
            k = s[len(s)-1]
        }
        l := strings.Fields(j)[1]
        m := strings.Fields(j)[0]

        psNameLabel.SetText(k)
        psNiLabel.SetText(l)
        psPidValue = m
    }

    search := &widget.Entry{PlaceHolder: "Search"}
    mainwindow := &widget.Form{
        Items: []*widget.FormItem{ // we can specify items in the constructor
            {Text: "Process:", Widget: psNameLabel},
            {Text: "Current nice value:", Widget: psNiLabel},
            {Text: "New nice value:", Widget: psNiEntry},
            {Widget: widget.NewButton("Save", psSaveButtonFunction)},
            {Widget: messageLabel},
            {Widget: widget.NewButton("man page", psManpageButtonFunction)},
        },
    }

    content := mainLayout(processListContent, search, mainwindow)

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

func formatWholeLines(processes []string) (formatted []string) {
    var allLines []string
    for i :=0; i < len(processes)-1; i++ {
        f := strings.Fields(processes[i])
        pid := f[0]
        ni := f[1]
        if strings.Contains(f[2], "/") {
            s := strings.Split(f[2], "/")
            comm := s[len(s)-1]
            allLines = append(allLines, pid + " " + ni + " " + comm)
        } else {
            comm := f[2]
            allLines = append(allLines, pid + " " + ni + " " + comm)
        }
    }
    return allLines
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
        case "comm":
            if strings.Contains(f[2], "/") {
                s := strings.Split(f[2], "/")
                comm := s[len(s)-1]
                allLines = append(allLines, []string{comm}...)
            } else {
                comm := f[2]
                allLines = append(allLines, []string{comm}...)
            }
        default:
            allLines = nil
        }
    }
    return allLines
}
