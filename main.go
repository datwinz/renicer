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
    searchBarButton *widget.Button,
    mainWindow *widget.Form,
    ) (*fyne.Container) {
    processes := container.New(layout.NewGridLayout(2), wholeProcesses, mainWindow)
    searchBarLayout := container.New(layout.NewGridLayout(2),
        container.NewPadded(searchBar), container.NewPadded(searchBarButton))
    totalLayout := container.NewBorder(searchBarLayout, nil, nil, nil, processes)
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

    formNameLabel := widget.NewLabel("process")
    formNiLabel := widget.NewLabel("0")
    var formPidValue string
    formNiEntry := widget.NewEntry()
    formMessageLabel := widget.NewLabel("")
    formSaveButtonFunction := func() {
        value := formNiEntry.Text
        valueInt, err := strconv.Atoi(value)
        if err != nil {
            msg := "New nice value should be a number"
            fmt.Println(msg)
            formMessageLabel.SetText(msg)
            return
        }
        if valueInt < -20 && valueInt > 20 {
            msg := "New nice value should be between -20 and 20"
            fmt.Println(msg)
            formMessageLabel.SetText(msg)
            return
        }
        fmt.Printf("%q %q %q", formPidValue, value, formNameLabel.Text)
        //Users other than the super-user may only alter the priority of processes they own,
        //and can only monotonically increase their ``nice value'' within the range 0 to
        //PRIO_MAX (20)
        if valueInt >=-20 && valueInt < 0 {
            fmt.Println("Spawn polkitd or mac-like window because value")
        }
        reniceCmd := exec.Command(renicePath, value, formPidValue)
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
        formMessageLabel.SetText("")
    }
    formManpageButtonFunction := func () {
        manCmd := exec.Command(manPath, formNameLabel.Text)
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

        formNameLabel.SetText(k)
        formNiLabel.SetText(l)
        formPidValue = m
    }

    searchBar := widget.NewEntry()
    searchBar.SetPlaceHolder("Search...")
    searchBarButton := widget.NewButton("Search", func () {
        searchBar.OnSubmitted(searchBar.Text)
    })

    mainWindow := &widget.Form{
        Items: []*widget.FormItem{ // we can specify items in the constructor
            {Text: "Process:", Widget: formNameLabel},
            {Text: "Current nice value:", Widget: formNiLabel},
            {Text: "New nice value:", Widget: formNiEntry},
            {Widget: widget.NewButton("Save", formSaveButtonFunction)},
            {Widget: formMessageLabel},
            {Widget: widget.NewButton("man page", formManpageButtonFunction)},
        },
    }

    searchBar.OnSubmitted = func(searchTerm string) {
        var searchResult []string
        for i := 0; i < len(psOutput) - 1; i++ {
            allLines := formatWholeLines(psOutput)[i]
            if strings.Contains(allLines, searchTerm) {
                searchResult = append(searchResult, allLines)
            }
        }
        searchedListContent := widget.NewList(
            func() (int) {
                return len(searchResult)
            },
            func () (fyne.CanvasObject) {
                return widget.NewLabel("Process")
            },
            func(j widget.ListItemID, p fyne.CanvasObject) {
                p.(*widget.Label).SetText(
                    searchResult[j],
                )
            },
        )
        searchedListContent.OnSelected = func(i widget.ListItemID) {
            j := searchResult[i]
            k := strings.Fields(j)[2]
            l := strings.Fields(j)[1]
            m := strings.Fields(j)[0]

            formNameLabel.SetText(k)
            formNiLabel.SetText(l)
            formPidValue = m
        }
        fmt.Println(searchResult)
        fmt.Println(len(searchResult))
        content := mainLayout(searchedListContent,
            searchBar,
            searchBarButton,
            mainWindow)
        w.SetContent(content)
    }

    content := mainLayout(processListContent, searchBar, searchBarButton, mainWindow)

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
    for i := 0; i < len(processes)-1; i++ {
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
