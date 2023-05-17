package help

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GetHelp(pages *tview.Pages) *tview.Flex {
	help := tview.NewFlex()
	help.SetBorder(true)
	help.SetTitle("Help")
	help.SetDirection(tview.FlexRow)
	help.SetBorderColor(tcell.Color133)

	helpText :=
		"q:    \t\tquit\n" +
			"TAB:  \t\tswitch focus between application list and binding view\n" +
			"a:    \t\tadd application\n" +
			"d:    \t\tdelete application\n" +
			"ENTER:\t\ton binding view start table selection mode\n" +
			"a:    \t\tin table selection mode add binging to application\n" +
			"e:    \t\tin table selection mode edit selected binding\n" +
			"d:    \t\tin table selection mode delete selected binding\n" +
			"ESC:  \t\texit table selection mode"

	help.AddItem(tview.NewTextView().SetText("ESC: back to application"), 0, 1, true)
	help.AddItem(tview.NewTextView().SetText(helpText), 0, 8, true)
	help.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages.SwitchToPage("Menu")
			return event
		} else {
			return nil
		}
	})
	return help
}
