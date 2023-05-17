# go-shortcut
simple TUI app to store shortcuts reminders

![Screenshot 2023-05-17-11-39-17](https://github.com/mebitek/go-shortcut/assets/1067967/632c7d53-6d58-4a91-9596-d4526cc8fcb9)

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


