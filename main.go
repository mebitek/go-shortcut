package main

import (
	"encoding/json"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"os"
	"shortcuts/utils"
	"strconv"
)

type Application struct {
	Name     string     `json:"name"`
	Category string     `json:"category"`
	Bindings []Bindings `json:"bindings"`
}

type Bindings struct {
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
	SetText("(a) to add a new application - (A) to add binding to selected application (d) to delete application - (q) to quit")

var primitives = make(map[tview.Primitive]int)
var primitivesIndexMap = make(map[int]tview.Primitive)

var dbFile = ""

func main() {

	initFocusMap()
	dbFile = utils.InitConfingDirectory()

	content, err := os.ReadFile(dbFile)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(content, &applications)

	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	addContactList()

	applicationList.SetSelectedFunc(func(index int, name string, secondName string, shortcut rune) {
		setBindingsTable(&applications[index])
		setDetailFlex(table, "Bindings")
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
			app.Stop()
		} else if event.Rune() == 97 && (applicationList.HasFocus() || menu.HasFocus()) {
			form.Clear(true)
			addApplicationForm()
			setDetailFlex(form, "Add Application")
			app.SetFocus(form)
		} else if event.Rune() == 65 && (applicationList.HasFocus() || menu.HasFocus()) {
			bindingForm.Clear(true)
			addBindingForm()
			index := applicationList.GetCurrentItem()
			application := applications[index]
			setDetailFlex(bindingForm, "Add Binding to "+application.Name)
			app.SetFocus(bindingForm)
		} else if event.Rune() == 100 && (applicationList.HasFocus() || menu.HasFocus()) {
			deleteApplication()
		} else if event.Rune() == 9 && (!form.HasFocus() && !bindingForm.HasFocus()) {
			primitive := app.GetFocus()
			actualPrimitiveIndex := primitives[primitive]
			app.SetFocus(getNextFocus(actualPrimitiveIndex + 1))
		}

		return event
	})

	pages.AddPage("Menu", flex, true, true)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func initFocusMap() {
	primitives[menu] = 0
	primitives[applicationList] = 1
	primitives[table] = 2

	primitivesIndexMap[0] = menu
	primitivesIndexMap[1] = applicationList
	primitivesIndexMap[2] = table
}

func getNextFocus(index int) tview.Primitive {

	if index == len(primitivesIndexMap) {
		index = 0
	}
	return primitivesIndexMap[index]
}

func setDetailFlex(element tview.Primitive, title string) {
	app.SetFocus(applicationList)
	detailFlex.Clear()
	detailFlex.SetBorder(true)
	detailFlex.AddItem(tview.NewFlex(), 0, 2, false)
	detailFlex.AddItem(element, 0, 1, true)
	detailFlex.AddItem(tview.NewFlex(), 0, 2, false)
	detailFlex.SetTitle(title).SetTitleColor(tcell.Color133)
}

func addContactList() {
	applicationList.Clear()
	for index, application := range applications {
		applicationList.AddItem(application.Name+" - "+application.Category, " ", rune(49+index), nil)
	}
}

func deleteApplication() {
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
		addContactList()
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

func addBindingForm() *tview.Form {
	bindingForm.Clear(true)
	bindingForm.SetFieldTextColor(tcell.Color133)
	bindingForm.SetLabelColor(tcell.Color133)
	bindingForm.SetFieldBackgroundColor(tcell.ColorWhite)
	bindingForm.SetFieldTextColor(tcell.ColorBlack)
	bindingForm.SetButtonBackgroundColor(tcell.Color133)

	binding := Bindings{}
	bindingForm.AddInputField("Function", "", 20, nil, func(function string) {
		binding.Function = function
	})

	bindingForm.AddInputField("Key Binding", "", 20, nil, func(keybind string) {
		binding.Keybind = keybind
	})

	bindingForm.AddButton("Save", func() {

		var index = applicationList.GetCurrentItem()
		application := applications[index]
		application.Bindings = append(application.Bindings, binding)
		applications[index] = application
		saveList(dbFile)
		setDetailFlex(table, "Bindings")
		setBindingsTable(&application)
		index = applicationList.GetItemCount() - 1
		applicationList.SetCurrentItem(index)
		bindingForm.Clear(true)
		app.SetFocus(applicationList)

	})

	bindingForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			setDetailFlex(table, "Bindings")
		}
		return event
	})

	return bindingForm
}

func setBindingsTable(application *Application) {
	var cols = 3
	table.Clear()
	detailFlex.SetTitle(application.Name + " - Bindings")
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
		}
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
		}
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		rows, cols := table.GetSelectable()
		if event.Rune() == 101 && (rows || cols) {
			// edit
		} else if event.Rune() == 100 && (rows || cols) {
			// delete
			row, col := table.GetSelection()
			if row == 0 {
				return event
			}
			_ = col
			application.Bindings = removeBinding(application.Bindings, row-1)
			saveList(dbFile)
			table.RemoveRow(row)
			table.SetSelectable(false, false)
		}
		return event
	})
}
