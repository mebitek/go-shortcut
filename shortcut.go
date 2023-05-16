package main

import "github.com/pterm/pterm"

func main() {

	panel1 := pterm.DefaultBox.Sprint("1")
	panel2 := pterm.DefaultBox.WithTitle("title").Sprint("2")
	panel3 := pterm.DefaultBox.WithTitle("bottom center title").WithTitleBottomCenter().Sprint("3")
	menu := pterm.DefaultHeader.Sprint("q to quit")

	panels, _ := pterm.DefaultPanel.WithPanels(pterm.Panels{
		{{Data: panel1}, {Data: panel2}},
		{{Data: panel3}},
		{{Data: menu}},
	}).Srender()

	pterm.DefaultBox.WithTitle("Lorem Ipsum").WithTitleTopCenter().WithRightPadding(0).WithBottomPadding(0).Println(panels)

	for true {
		result, _ := pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show()

		if result == "q" {
			break
		}
	}

}
