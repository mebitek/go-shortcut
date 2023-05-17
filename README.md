# go-shortcut
simple TUI app to store shortcuts reminders

# dependecies 

[TVIEW](https://github.com/rivo/tview)

# run

```
go get github.com/rivo/tview
```

```
go run main.go
```

# shortcuts

| function | key                                                    |
|----------|--------------------------------------------------------|
| q        | quit                                                   |
| TAB      | switch focus between application list and binding view |
| a        | add application                                        |
| d        | delete application                                     |
| enter    | on binding view start table selection mode             |
| a        | in table selection mode add binging to application     |
| e        | in table selection mode edit selected binding          |
| d        | in table selection mode delete selected binding        |
| ESC      | exit table selection mode                              |

