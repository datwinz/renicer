package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
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
    logFile, err := os.OpenFile("/tmp/renicelog", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Println("[INFO] logfile:", err)
    }
    log.SetOutput(logFile)

    psPath := processPaths("ps")
    renicePath := processPaths("renice")
    manPath := processPaths("man")
    mandocPath := processPaths("mandoc")

    psOutput := findProcesses(psPath)

    app := app.New()
    window := app.NewWindow("Renicer")

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
            log.Println("[INFO] form:", msg)
            formMessageLabel.SetText(msg)
            return
        }
        if newValueInt < -20 && newValueInt > 20 {
            msg := "New nice value should be between -20 and 20"
            log.Println("[INFO] form:", msg)
            formMessageLabel.SetText(msg)
            return
        }
        oldValueInt, err := strconv.Atoi(formNiLabel.Text)
        if err != nil {
            msg := "Existing nice value isn't a number"
            log.Println("[INFO] form:", msg)
        }
        log.Printf("form new values: %q %q %q", formPidValue, newValue, formNameLabel.Text)
        // Users other than the super-user may only alter the priority of processes they own,
        // and can only monotonically increase their ``nice value'' within the range 0 to
        // PRIO_MAX (20).
        if newValueInt >=-20 && newValueInt < 0 {
            authorisationConditional(formMessageLabel, newValue, formPidValue)
            newPsOutput := findProcesses(psPath)
            processListContent = makeListContent(
                len(newPsOutput),
                "Process",
                formatWholeLines(newPsOutput),
            )
            // This doesn't work. As in processListContent gets updated, but the content of w
            // doesn't. I can't redraw the whole screen as in the searchBar.OnSubmitted
            // function because of recursion: the var content needs the formSaveButtonFunction
            // and the redrawing of the whole screen needs the var content.
            window.Content().Refresh()
            return
        } else if newValueInt < oldValueInt {
            authorisationConditional(formMessageLabel, newValue, formPidValue)
            return
        }
        reniceCmd := exec.Command(renicePath, newValue, formPidValue)
        err = reniceCmd.Run()
        if err != nil {
            log.Println("[INFO] renice:", err)
            authorisationConditional(formMessageLabel, newValue, formPidValue)
            return
        }
        formMessageLabel.SetText("")
    }
    formManpageButtonFunction := func () {
        manPagePath := exec.Command(manPath, "-w", formNameLabel.Text)
        var b strings.Builder
        manPagePath.Stdout = &b
        err := manPagePath.Run()
        if err != nil {
            log.Println("[INFO] man: Couldn't find path of manpage")
        }
        manFilePath := strings.TrimSpace(b.String())

        mandocCmd := exec.Command(mandocPath, "-Tmarkdown", manFilePath)
        var c strings.Builder
        mandocCmd.Stdout = &c
        err = mandocCmd.Run()
        if err != nil {
            log.Println("[INFO] man:", err)
        }

        windowChild := app.NewWindow("manpage")
        text := widget.NewRichTextFromMarkdown(c.String())
        textContainer := container.NewScroll(text)
        windowChild.SetContent(textContainer)
        windowChild.Resize(window.Content().Size())
        text.Wrapping = 2
        windowChild.Show()
    }
    mainWindowForm := &widget.Form{
        Items: []*widget.FormItem{ // we can specify items in the constructor
            {Text: "Process:", Widget: formNameLabel},
            {Text: "Current nice value:", Widget: formNiLabel},
            {Text: "New nice value:", Widget: formNiEntry},
            {Widget: widget.NewButton("Save", formSaveButtonFunction)},
            {Widget: formMessageLabel},
            {Widget: widget.NewButton("man page", formManpageButtonFunction)},
        },
    }

    processListContent.OnSelected = func(i widget.ListItemID) {
        selectedProcess := psOutput[i]
        selectedProcessName := strings.Fields(selectedProcess)[2]
        if strings.Contains(selectedProcessName, "/") {
            selectedProcessName = path.Base(selectedProcessName)
        }
        selectedProcessNice := strings.Fields(selectedProcess)[1]
        selectedProcessPid := strings.Fields(selectedProcess)[0]

        formNameLabel.SetText(selectedProcessName)
        formNiLabel.SetText(selectedProcessNice)
        formPidValue = selectedProcessPid
    }

    searchBar := widget.NewEntry()
    searchBar.SetPlaceHolder("Search...")
    searchBarButton := widget.NewButton("Search", func () {
        searchBar.OnSubmitted(searchBar.Text)
    })

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
            selectedProcess := searchResult[i]
            selectedProcessName := strings.Fields(selectedProcess)[2]
            selectedProcessNice := strings.Fields(selectedProcess)[1]
            selectedProcessPid := strings.Fields(selectedProcess)[0]

            formNameLabel.SetText(selectedProcessName)
            formNiLabel.SetText(selectedProcessNice)
            formPidValue = selectedProcessPid
        }
        log.Println("[INFO] search result:", searchResult)
        log.Println("[INFO] search length:", len(searchResult))
        content := mainLayout(searchedListContent,
            searchBar,
            searchBarButton,
            mainWindowForm)
        window.SetContent(content)
    }

    content := mainLayout(processListContent, searchBar, searchBarButton, mainWindowForm)


    window.SetContent(content)
    window.ShowAndRun()
}

func processPaths(processName string) (path string) {
    path, err := exec.LookPath(processName)
    if err != nil {
        log.Fatal("[ERROR] process path:", err)
    }
    return path
}

func findProcesses(psPath string) (processes []string) {
    psCmd := exec.Command(psPath, "ax", "-o", "pid,ni,comm")
    var outAll strings.Builder
    psCmd.Stdout = &outAll
    err := psCmd.Run()
    if err != nil {
        log.Fatal("[ERROR] processes:", err)
    }

    outSingle := strings.Split(outAll.String(), "\n")
    return outSingle
}

func formatWholeLines(processes []string) (formatted []string) {
    var allLines []string
    for i := 0; i < len(processes)-1; i++ {
        f := strings.Fields(processes[i])
        pid := f[0]
        ni := f[1]
        if strings.Contains(f[2], "/") {
            comm := path.Base(f[2])
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
                comm := path.Base(f[2])
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

func macAuthorisation(niceValue string, pidValue string) {
    osaPath := processPaths("osascript")
    osaInnerScript := "renice " + niceValue + " " + pidValue
    osaOuterScript := "do shell script \"" +
        osaInnerScript +
        "\" with administrator privileges"
    osaReniceCmd := exec.Command(osaPath, "-e", osaOuterScript)
    err := osaReniceCmd.Run()
    if err != nil {
        log.Println("[INFO] osascript renice:", err)
    }
}

func linuxPolkitAuthorisation(niceValue string, pidValue string) {
    pkexecPath := processPaths("pkexec")
    pkexecCmd := exec.Command(pkexecPath, "renice", niceValue, pidValue)
    err := pkexecCmd.Run()
    if err != nil {
        log.Println("[INFO] pkexec renice:", err)
    }
}

func authorisationConditional(pidLabel *widget.Label, niceValue string, pidValue string) {
    if runtime.GOOS == "darwin" {
        macAuthorisation(niceValue, pidValue)
        pidLabel.SetText("")
        //w.Content().Refresh()
        return
    } else if runtime.GOOS == "linux" {
        linuxPolkitAuthorisation(niceValue, pidValue)
        pidLabel.SetText("")
        //w.Content().Refresh()
        return
    }
}
