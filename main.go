package main

import (
	"encoding/json"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"os"
	"shortcuts/utils"
)

type Shortcut struct {
	Name     string     `json:"name"`
	Category string     `json:"category"`
	Bindings []Bindings `json:"bindings"`
}

type Bindings struct {
	Function string `json:"function"`
	Keybind  string `json:"keybind"`
}

var shortcuts = make([]Shortcut, 0)

var pages = tview.NewPages()
var app = tview.NewApplication()
var form = tview.NewForm()

var bindingForm = tview.NewForm()
var shortcutList = tview.NewList().ShowSecondaryText(false).SetSelectedBackgroundColor(tcell.Color133)
var flex = tview.NewFlex()

var detailFlex = tview.NewFlex()

var table = tview.NewTable().SetBorders(true)
var text = tview.NewTextView().
	SetTextColor(tcell.Color133).
	SetText("(a) to add a new shortcut - (A) to add binding to selected item (d) to delete shortcut - (q) to quit")

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

	err = json.Unmarshal(content, &shortcuts)

	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	addContactList()

	shortcutList.SetSelectedFunc(func(index int, name string, secondName string, shortcut rune) {
		setBindingsTable(&shortcuts[index])
		setDetailFlex(table, "Bindings")
	})

	shortcutList.SetBorder(true).SetTitle("Shortcuts").SetTitleColor(tcell.Color133)
	text.SetBorder(true).SetTitle("Menu").SetTitleColor(tcell.Color133)
	setDetailFlex(table, "Bindings")

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(shortcutList, 0, 1, true).
			AddItem(detailFlex, 0, 4, true), 0, 6, true).
		AddItem(text, 0, 1, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 && !form.HasFocus() && !bindingForm.HasFocus() {
			app.Stop()
		} else if event.Rune() == 97 && !form.HasFocus() && !bindingForm.HasFocus() {
			form.Clear(true)
			addShortcutForm()
			setDetailFlex(form, "Add Shortcut")
			app.SetFocus(form)
		} else if event.Rune() == 65 && !form.HasFocus() && !bindingForm.HasFocus() {
			bindingForm.Clear(true)
			addBinding()
			index := shortcutList.GetCurrentItem()
			shortcut := shortcuts[index]
			setDetailFlex(bindingForm, "Add Binding to "+shortcut.Name)
			app.SetFocus(bindingForm)
		} else if event.Rune() == 100 && !form.HasFocus() && !bindingForm.HasFocus() {
			deleteShortcut()
		} else if event.Rune() == 9 && !form.HasFocus() && !bindingForm.HasFocus() {
			primitive := app.GetFocus()
			actualPrimitiveIndex := primitives[primitive]
			app.SetFocus(getNextFocus(actualPrimitiveIndex + 1))
		}

		return event
	})

	pages.AddPage("Menu", flex, true, true)
	pages.AddPage("Add Contact", form, true, false)
	pages.AddPage("Add Binding", bindingForm, true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func initFocusMap() {
	primitives[text] = 0
	primitives[shortcutList] = 1
	primitives[table] = 2

	primitivesIndexMap[0] = text
	primitivesIndexMap[1] = shortcutList
	primitivesIndexMap[2] = table
}

func getNextFocus(index int) tview.Primitive {

	if index == len(primitivesIndexMap) {
		index = 0
	}
	return primitivesIndexMap[index]
}

func setDetailFlex(element tview.Primitive, title string) {
	detailFlex.Clear()
	detailFlex.SetBorder(true)
	detailFlex.AddItem(tview.NewFlex(), 0, 4, true)
	detailFlex.AddItem(element, 0, 1, true)
	detailFlex.AddItem(tview.NewFlex(), 0, 4, true)
	detailFlex.SetTitle(title).SetTitleColor(tcell.Color133)
}

func addContactList() {
	shortcutList.Clear()
	for index, contact := range shortcuts {
		shortcutList.AddItem(contact.Name+" - "+contact.Category, " ", rune(49+index), nil)
	}
}

func deleteShortcut() {
	index := shortcutList.GetCurrentItem()
	shortcutList.RemoveItem(index)
	shortcuts = remove(shortcuts, index)
	saveList(dbFile)
	shortcutList.SetCurrentItem(index - 1)
}

func remove(slice []Shortcut, s int) []Shortcut {
	return append(slice[:s], slice[s+1:]...)
}

func saveList(dbFile string) {
	output, err := json.Marshal(shortcuts)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(dbFile, output, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
}

func addShortcutForm() *tview.Form {
	form.Clear(true)
	form.SetFieldTextColor(tcell.Color133)
	form.SetLabelColor(tcell.Color133)
	form.SetFieldBackgroundColor(tcell.ColorWhite)
	form.SetFieldTextColor(tcell.ColorBlack)
	form.SetButtonBackgroundColor(tcell.Color133)
	form.SetBorder(true)

	shortcut := Shortcut{}

	form.AddInputField("Name", "", 20, nil, func(name string) {
		shortcut.Name = name
	})

	form.AddInputField("Category", "", 20, nil, func(category string) {
		shortcut.Category = category
	})

	form.AddButton("Save", func() {
		shortcuts = append(shortcuts, shortcut)
		addContactList()
		saveList(dbFile)
		setDetailFlex(table, "Bindings")
		index := shortcutList.GetItemCount() - 1
		shortcutList.SetCurrentItem(index)
		form.Clear(true)
		app.SetFocus(shortcutList)

	})
	return form
}

func addBinding() *tview.Form {
	bindingForm.Clear(true)
	bindingForm.SetFieldTextColor(tcell.Color133)
	bindingForm.SetLabelColor(tcell.Color133)
	bindingForm.SetFieldBackgroundColor(tcell.ColorWhite)
	bindingForm.SetFieldTextColor(tcell.ColorBlack)
	bindingForm.SetButtonBackgroundColor(tcell.Color133)
	bindingForm.SetBorder(true)

	binding := Bindings{}
	bindingForm.AddInputField("Function", "", 20, nil, func(function string) {
		binding.Function = function
	})

	bindingForm.AddInputField("Key Binding", "", 20, nil, func(keybind string) {
		binding.Keybind = keybind
	})

	bindingForm.AddButton("Save", func() {

		var index = shortcutList.GetCurrentItem()
		shortcut := shortcuts[index]
		shortcut.Bindings = append(shortcut.Bindings, binding)
		shortcuts[index] = shortcut
		saveList(dbFile)
		setDetailFlex(table, "Bindings")
		setBindingsTable(&shortcut)
		index = shortcutList.GetItemCount() - 1
		shortcutList.SetCurrentItem(index)
		bindingForm.Clear(true)
		app.SetFocus(shortcutList)

	})
	return bindingForm
}

func setBindingsTable(shortcut *Shortcut) {
	var cols = 2
	table.Clear()
	detailFlex.SetTitle(shortcut.Name + " - Bindings")
	table.SetCell(0, 0, tview.NewTableCell("Function").SetTextColor(tcell.Color133))
	table.SetCell(0, 1, tview.NewTableCell("Binding").SetTextColor(tcell.Color133))
	var i = 0
	for r := 1; r <= len(shortcut.Bindings); r++ {
		for c := 0; c < cols; c++ {
			var value = ""
			if c == 0 {
				value = shortcut.Bindings[i].Function
			} else if c == 1 {
				value = shortcut.Bindings[i].Keybind
			}
			table.SetCell(r, c,
				tview.NewTableCell(value).
					SetAlign(tview.AlignCenter))
		}
		i++
	}
}
