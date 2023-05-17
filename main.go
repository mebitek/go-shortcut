package main

import (
	"encoding/json"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"os"
	"shortcuts/help"
	"shortcuts/utils"
	"strconv"
)

type Application struct {
	Name     string     `json:"name"`
	Category string     `json:"category"`
	Bindings []Bindings `json:"bindings"`
}

type Bindings struct {
	Index    int    `json:"index"`
	Function string `json:"function"`
	Keybind  string `json:"keybind"`
}

var applications = make([]Application, 0)

var pages = tview.NewPages()
var app = tview.NewApplication()
var form = tview.NewForm()

var bindingForm = tview.NewForm()
var applicationList = tview.NewList().ShowSecondaryText(false).SetSelectedBackgroundColor(tcell.Color133)
var flex = tview.NewFlex()
var detailFlex = tview.NewFlex()
var table = tview.NewTable().SetBorders(true)
var menu = tview.NewTextView().
	SetTextColor(tcell.Color133).
	SetText("(a) to add a new application - (d) to delete application - [ENTER] edit bindings - (?) help - (q) to quit")

var primitives = make(map[tview.Primitive]int)
var primitivesIndexMap = make(map[int]tview.Primitive)

var editMode = false

var dbFile = ""

func main() {

	initFocus()
	dbFile = utils.InitConfingDirectory()

	content, err := os.ReadFile(dbFile)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(content, &applications)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	var helpFlex = help.GetHelp(pages)

	addApplicationList()

	applicationList.SetSelectedFunc(func(index int, name string, secondName string, shortcut rune) {
		setBindingsTable(&applications[index])
		setDetailFlex(table, applications[index].Name+" Bindings")
	})
	applicationList.SetBorder(true).SetTitle("Applications").SetTitleColor(tcell.Color133)

	menu.SetBorder(true).SetTitle("Menu").SetTitleColor(tcell.Color133)
	setDetailFlex(table, "Bindings")

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(applicationList, 0, 1, true).
			AddItem(detailFlex, 0, 4, true), 0, 6, true).
		AddItem(menu, 0, 1, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 && !form.HasFocus() && !bindingForm.HasFocus() {
			// q
			app.Stop()
		} else if event.Rune() == 97 && (!form.HasFocus() && !bindingForm.HasFocus() && !editMode && !helpFlex.HasFocus()) {
			// a
			form.Clear(true)
			addApplicationForm()
			setDetailFlex(form, "Add Application")
			app.SetFocus(form)
		} else if event.Rune() == 100 && (!form.HasFocus() && !bindingForm.HasFocus() && !editMode && !helpFlex.HasFocus()) {
			// d
			deleteApplication()
		} else if event.Rune() == 9 && (!form.HasFocus() && !bindingForm.HasFocus() && !helpFlex.HasFocus()) {
			// tab
			primitive := app.GetFocus()
			actualPrimitiveIndex := primitives[primitive]
			app.SetFocus(getNextFocus(actualPrimitiveIndex + 1))
		} else if event.Rune() == 63 && !helpFlex.HasFocus() {
			pages.SwitchToPage("Help")
		}

		return event
	})

	pages.AddPage("Menu", flex, true, true)
	pages.AddPage("Help", helpFlex, true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func initFocus() {
	//primitives[menu] = 0
	primitives[applicationList] = 0
	primitives[table] = 1

	//primitivesIndexMap[0] = menu
	primitivesIndexMap[0] = applicationList
	primitivesIndexMap[1] = table

	applicationList.SetFocusFunc(func() {
		applicationList.SetBorderColor(tcell.Color133)
		applicationList.SetTitleColor(tcell.ColorWhite)
		menu.SetBorderColor(tcell.ColorWhite)
		menu.SetTitleColor(tcell.Color133)
		detailFlex.SetBorderColor(tcell.ColorWhite)
		detailFlex.SetTitleColor(tcell.Color133)

	})
	menu.SetFocusFunc(func() {
		menu.SetBorderColor(tcell.Color133)
		menu.SetTitleColor(tcell.ColorWhite)
		applicationList.SetBorderColor(tcell.ColorWhite)
		applicationList.SetTitleColor(tcell.Color133)
		detailFlex.SetBorderColor(tcell.ColorWhite)
		detailFlex.SetTitleColor(tcell.Color133)
	})

	table.SetFocusFunc(func() {
		detailFlex.SetBorderColor(tcell.Color133)
		detailFlex.SetTitleColor(tcell.ColorWhite)
		applicationList.SetBorderColor(tcell.ColorWhite)
		applicationList.SetTitleColor(tcell.Color133)
		menu.SetBorderColor(tcell.ColorWhite)
		menu.SetTitleColor(tcell.Color133)
	})
}

func getNextFocus(index int) tview.Primitive {

	if index == len(primitivesIndexMap) {
		index = 0
	}
	return primitivesIndexMap[index]
}

func setDetailFlex(element tview.Primitive, title string) {
	app.SetFocus(element)
	detailFlex.Clear()
	detailFlex.SetBorder(true)
	detailFlex.AddItem(tview.NewFlex(), 0, 2, false)
	detailFlex.AddItem(element, 0, 1, true)
	detailFlex.AddItem(tview.NewFlex(), 0, 2, false)
	detailFlex.SetTitle(title).SetTitleColor(tcell.Color133)
}

func addApplicationList() {
	applicationList.Clear()
	for index, application := range applications {
		applicationList.AddItem(application.Name+" - "+application.Category, " ", rune(49+index), nil)
	}
	applicationList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			return nil
		}
		return event
	})
}

func deleteApplication() {
	if applicationList.GetItemCount() == 0 {
		return
	}

	index := applicationList.GetCurrentItem()
	applicationList.RemoveItem(index)
	applications = removeApplication(applications, index)
	saveList(dbFile)
	applicationList.SetCurrentItem(index - 1)
}

func removeApplication(slice []Application, s int) []Application {
	return append(slice[:s], slice[s+1:]...)
}
func removeBinding(slice []Bindings, s int) []Bindings {
	return append(slice[:s], slice[s+1:]...)
}

func saveList(dbFile string) {
	output, err := json.Marshal(applications)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(dbFile, output, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
}

func addApplicationForm() *tview.Form {
	form.Clear(true)
	form.SetFieldTextColor(tcell.Color133)
	form.SetLabelColor(tcell.Color133)
	form.SetFieldBackgroundColor(tcell.ColorWhite)
	form.SetFieldTextColor(tcell.ColorBlack)
	form.SetButtonBackgroundColor(tcell.Color133)

	application := Application{}

	form.AddInputField("Name", "", 20, nil, func(name string) {
		application.Name = name
	})

	form.AddInputField("Category", "", 20, nil, func(category string) {
		application.Category = category
	})

	form.AddButton("Save", func() {
		applications = append(applications, application)
		addApplicationList()
		saveList(dbFile)
		setDetailFlex(table, "Bindings")
		index := applicationList.GetItemCount() - 1
		applicationList.SetCurrentItem(index)
		form.Clear(true)
		app.SetFocus(applicationList)

	})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			setDetailFlex(table, "Bindings")
		}
		return event
	})

	return form
}

func addBindingForm(binding Bindings, edit bool) *tview.Form {
	bindingForm.Clear(true)
	bindingForm.SetFieldTextColor(tcell.Color133)
	bindingForm.SetLabelColor(tcell.Color133)
	bindingForm.SetFieldBackgroundColor(tcell.ColorWhite)
	bindingForm.SetFieldTextColor(tcell.ColorBlack)
	bindingForm.SetButtonBackgroundColor(tcell.Color133)

	bindingForm.AddInputField("Function", binding.Function, 20, nil, func(function string) {
		binding.Function = function
	})

	bindingForm.AddInputField("Key Binding", binding.Keybind, 20, nil, func(keybind string) {
		binding.Keybind = keybind
	})

	bindingForm.AddButton("Save", func() {
		var index = applicationList.GetCurrentItem()
		application := applications[index]
		if !edit {
			binding.Index = len(application.Bindings)
			application.Bindings = append(application.Bindings, binding)
		} else {
			application.Bindings[binding.Index] = binding
		}
		applications[index] = application
		saveList(dbFile)
		setDetailFlex(table, "Bindings - EDIT")
		setBindingsTable(&application)
		index = applicationList.GetItemCount() - 2
		applicationList.SetCurrentItem(index)
		bindingForm.Clear(true)
		app.SetFocus(table)

	})

	bindingForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			setDetailFlex(table, "Bindings - EDIT")
			app.SetFocus(table)
		}
		return event
	})

	return bindingForm
}

func setBindingsTable(application *Application) {
	var cols = 3
	table.Clear()
	table.SetCell(0, 0, tview.NewTableCell("Index").SetTextColor(tcell.Color133))
	table.SetCell(0, 1, tview.NewTableCell("Function").SetTextColor(tcell.Color133))
	table.SetCell(0, 2, tview.NewTableCell("Binding").SetTextColor(tcell.Color133))
	var i = 0
	for r := 1; r <= len(application.Bindings); r++ {
		for c := 0; c < cols; c++ {
			var value = ""
			if c == 0 {
				value = strconv.Itoa(i)
			} else if c == 1 {
				value = application.Bindings[i].Function
			} else if c == 2 {
				value = application.Bindings[i].Keybind
			}
			table.SetCell(r, c,
				tview.NewTableCell(value).
					SetAlign(tview.AlignCenter))
		}
		i++
	}
	table.Select(1, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			table.SetSelectable(false, false)
			detailFlex.SetTitle(application.Name + " Bindings")
			editMode = false
		}
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
			detailFlex.SetTitle(detailFlex.GetTitle() + " - EDIT")
			editMode = true
		}
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		rows, cols := table.GetSelectable()
		if event.Rune() == 97 && editMode {
			bindingForm.Clear(true)
			addBindingForm(Bindings{}, false)
			index := applicationList.GetCurrentItem()
			application := applications[index]
			setDetailFlex(bindingForm, "Add Binding to "+application.Name)
			app.SetFocus(bindingForm)
		} else if event.Rune() == 101 && (rows || cols) {
			// edit
			row, col := table.GetSelection()
			if row == 0 {
				return event
			}
			_ = col
			if table.GetRowCount() == 1 {
				return event
			}
			addBindingForm(application.Bindings[row-1], true)
			setDetailFlex(bindingForm, "Edit Binding to "+application.Name)
			app.SetFocus(bindingForm)
		} else if event.Rune() == 100 && (rows || cols) {
			// delete
			row, col := table.GetSelection()
			if row == 0 {
				return event
			}
			_ = col
			if table.GetRowCount() == 1 {
				return event
			}
			application.Bindings = removeBinding(application.Bindings, row-1)
			saveList(dbFile)
			table.RemoveRow(row)
		}
		return event
	})
}
