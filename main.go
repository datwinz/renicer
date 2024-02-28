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
    mandocPath := processPaths("mandoc")

    psOutput := findProcesses(psPath)

    a := app.New()
    w := a.NewWindow("Renicer")

    processListContent := makeListContent(
        len(psOutput), 
        "Process",
        formatWholeLines(psOutput),
    )

    formNameLabel := widget.NewLabel("process")
    formNiLabel := widget.NewLabel("0")
    var formPidValue string
    formNiEntry := widget.NewEntry()
    formMessageLabel := widget.NewLabel("")
    formSaveButtonFunction := func() {
        newValue := formNiEntry.Text
        newValueInt, err := strconv.Atoi(newValue)
        if err != nil {
            msg := "New nice value should be a number"
            fmt.Println(msg)
            formMessageLabel.SetText(msg)
            return
        }
        if newValueInt < -20 && newValueInt > 20 {
            msg := "New nice value should be between -20 and 20"
            fmt.Println(msg)
            formMessageLabel.SetText(msg)
            return
        }
        oldValueInt, err := strconv.Atoi(formPidValue)
        if err != nil {
            msg := "Existing nice value isn't a number"
            fmt.Println(msg)
        }
        fmt.Printf("%q %q %q", formPidValue, newValue, formNameLabel.Text)
        // Users other than the super-user may only alter the priority of processes they own,
        // and can only monotonically increase their ``nice value'' within the range 0 to
        // PRIO_MAX (20).
        if newValueInt >=-20 && newValueInt < 0 {
            fmt.Println("Spawn polkitd or mac-like window because value")
        } else if newValueInt < oldValueInt {
            fmt.Println("Spawn polkitd or mac-like window because value")
        }
        reniceCmd := exec.Command(renicePath, newValue, formPidValue)
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
        manPagePath := exec.Command(manPath, "-w", formNameLabel.Text)
        var b strings.Builder
        manPagePath.Stdout = &b
        err := manPagePath.Run()
        if err != nil {
            fmt.Println("Couldn't find path of manpage")
        }
        manFilePath := strings.TrimSpace(b.String())

        mandocCmd := exec.Command(mandocPath, "-Tmarkdown", manFilePath)
        var c strings.Builder
        mandocCmd.Stdout = &c
        err = mandocCmd.Run()
        if err != nil {
            fmt.Println(err)
        }

        w2 := a.NewWindow("manpage")
        text := widget.NewRichTextFromMarkdown(c.String())
        w2.SetContent(text)
        w2.Resize(w.Content().Size())
        text.Wrapping = 2
        // add scrolling somehow
        w2.Show()
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
        searchedListContent := makeListContent(len(searchResult), "Process", searchResult)
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

func makeListContent(lenOfList int, templateString string, labelText []string) (*widget.List) {
    List := widget.NewList(
        func() (int) {
            return lenOfList
        },
        func () (fyne.CanvasObject) {
            return widget.NewLabel(templateString)
        },
        func(i widget.ListItemID, o fyne.CanvasObject) {
            o.(*widget.Label).SetText(
                labelText[i],
            )
        },
    )
    return List
}
