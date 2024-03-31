# renicer

App to show nice values of processes and renice them. It is basically a wrapper around the ```ps``` and ```renice``` commands.  It can also show manpages of commands that have them.

It has some extra functionality, namely that you can show the manpages of the listed processes. That uses the ```man``` and the ```mandoc``` command. It authorises on MacOS with ```osascript``` and on Linux with ```pkexec``` (using [polkit](https://github.com/polkit-org/polkit?tab=readme-ov-file)).

## Usage

<img width="483" alt="renicer_example" src="https://github.com/datwinz/renicer/assets/77931290/b0650064-445a-4cbf-a01e-0405dc1314e4">

Nice values dictate the scheduling priority of processes on *nixes. I.e. it changes how much process power the command or application gets from you CPU. The value ranges from -20 to 20, the *lower* the value, the *higher* the priority.

In this app you can change those values. Search for the name of the command you want in the search bar on top. You can also search for the process ID or the nice value if you want. Put in a new nice value and click on save. In a few cases you need superuser privileges, in those cases the app asks for authorisation.

## Installing

On MacOS you can use [brew](https://brew.sh):

```bash
brew install datwinz/formulae-and-casks/renicer
```

Otherwise you can download the zip file under "Releases". Unzip it and move it to your applications folder.

On Linux you can install with the Makefile:

```bash
# Install in ~/.local/
make user-install
# System install
sudo make install
```

Or move the renicer binary into your path.

## Dependencies

### Commands/binaries in $PATH

* ```ps```
* ```renice```
* ```man```
* ```mandoc```

#### MacOS

* ```osascript```

#### Linux

* ```pkexec```

## Build

You need to install [go](https://go.dev/dl/) on your platform of choice. You also need [git](https://git-scm.com/downloads). Then run:

```bash
git clone https://github.com/datwinz/renicer
cd renicer
go install fyne.io/fyne/v2/cmd/fyne@latest
fyne package --icon resources/renicer_light.png --name renicer
# FYNE_THEME=dark; fyne package --icon resources/renicer_dark.png --name renicer
```

## Basic design outline

- [x] Do 'ps ax -o pid,ni,comm' and make sort by name, process number nice value.
I.i.r.c. Linux has different words for the options, but if I look it up the only difference
is that in Linux you can also use cmd instead of comm.
- [x] Put it in window something like this:

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

- [x] Add search.
- [x] Show man pages or something for processes on double click.
- [x] Use polkit for linux and some privilege escalation for mac.
