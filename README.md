# renicer

App to show nice values of processes and renice them. It is basically a wrapper around the ```ps``` and ```renice``` commands. 

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
