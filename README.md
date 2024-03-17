# renicer

App to show nice values of processes and renice them. It is basically a wrapper around the ```ps``` and ```renice``` commands. 

I has some extra functionality, namely that you can show the manpages of the listed processes. That uses the ```man``` and the ```mandoc``` command. It authorises on MacOS with ```osascript```. In the future it wil authorise on Linux with ```polkit```.

## Usage

Nice values dictate the scheduling priority of processes on *nixes. I.e. it changes how much process power the command or application gets from you CPU. The value ranges from -20 to 20, the *lower* the value, the *higher* the priority.

In this app you can change those values. Search for the name of the command you want in the search bar on top. You can also search for the process ID or the nice value if you want. Put in a new nice value and click on safe. In a few cases you need superuser privileges, in those cases the app asks for authorisation.

## Dependencies

### Commands/binaries in $PATH

* ```ps```
* ```renice```
* ```man```
* ```mandoc```

#### MacOS

* ```osascript```

## Build

```bash
go install fyne.io/fyne/v2/cmd/fyne@latest
FYNE_THEME=light; fyne package --icon resources/renicer_light.png --name renicer
# FYNE_THEME=dark; fyne package --icon resources/renicer_dark.png --name renicer
```

## Basic design outline

Do 'ps ax -o pid,ni,comm' and make sort by name, process number nice value
I.i.r.c. Linux has different words for the options, but if I look it up the only difference
is that in Linux you can also use cmd instead of comm

Put it in window something like this:

```
__________________________________________________
| pid | Process  | Ni |                          |
__________________________________________________
| 1   | init     | 0  |                          |
| 2   | process2 | 0  |      process5            |
| 3   | process3 | 0  |      Old value: 0        |
| 4   | process4 | 0  |      New value : -20     |
| 5   | process5 | 0  |                          |
| 6   | process6 | 0  |      Save                |
--------------------------------------------------
```

Use polkit for linux and some privilege escalation for mac

Add search

Show man pages or something for processes on double click
