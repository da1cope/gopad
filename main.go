package main

import (
	"io"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Editor struct {
	win         fyne.Window
	text        *widget.Entry
	currentFile string
	darkMode    bool
}

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme()) // Start in light mode

	w := a.NewWindow("GoPad - Simple Text Editor")
	w.SetMaster()

	e := &Editor{
		win:      w,
		text:     widget.NewMultiLineEntry(),
		darkMode: false,
	}
	e.text.SetPlaceHolder("Start typing...")
	e.text.Wrapping = fyne.TextWrapWord

	// Toolbar with Dark Mode toggle
	darkModeBtn := widget.NewButton("Dark Mode", e.toggleTheme)
	toolbar := container.NewHBox(
		widget.NewButton("New", e.newFile),
				     widget.NewButton("Open", e.openFile),
				     widget.NewButton("Save", e.saveFile),
				     widget.NewButton("│", nil), // visual separator
				     darkModeBtn,
	)

	// Menu (added "Toggle Theme" + shortcut Ctrl+D)
	fileMenu := fyne.NewMenu("File",
				 fyne.NewMenuItem("New", e.newFile),
				 fyne.NewMenuItem("Open...", e.openFile),
				 fyne.NewMenuItemSeparator(),
				 fyne.NewMenuItem("Save", e.saveFile),
				 fyne.NewMenuItem("Save As...", e.saveFileAs),
				 fyne.NewMenuItemSeparator(),
				 fyne.NewMenuItem("Toggle Dark Mode\tCtrl+D", e.toggleTheme),
				 fyne.NewMenuItemSeparator(),
				 fyne.NewMenuItem("Quit", a.Quit),
	)

	// Keyboard shortcuts
	fileMenu.Items[0].Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierControl}
	fileMenu.Items[1].Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierControl}
	fileMenu.Items[3].Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
	fileMenu.Items[6].Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyD, Modifier: fyne.KeyModifierControl} // Ctrl+D

	w.SetMainMenu(fyne.NewMainMenu(fileMenu))
	w.SetContent(container.NewBorder(toolbar, nil, nil, nil, e.text))

	w.Resize(fyne.NewSize(800, 600))
	w.CenterOnScreen()
	w.ShowAndRun()
}

func (e *Editor) toggleTheme() {
	e.darkMode = !e.darkMode
	if e.darkMode {
		fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	} else {
		fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	}
}

func (e *Editor) newFile() {
	e.text.SetText("")
	e.currentFile = ""
	e.win.SetTitle("GoPad - Simple Text Editor")
}

func (e *Editor) openFile() {
	dialog.ShowFileOpen(func(r fyne.URIReadCloser, err error) {
		if r == nil || err != nil {
			return
		}
		defer r.Close()
		data, _ := io.ReadAll(r)
		e.text.SetText(string(data))
		e.currentFile = r.URI().Path()
		e.win.SetTitle("GoPad - " + r.URI().Name())
	}, e.win)
}

func (e *Editor) saveFile() {
	if e.currentFile == "" {
		e.saveFileAs()
		return
	}
	e.writeFile(e.currentFile)
}

func (e *Editor) saveFileAs() {
	dialog.ShowFileSave(func(w fyne.URIWriteCloser, err error) {
		if w == nil || err != nil {
			return
		}
		defer w.Close()
		path := w.URI().Path()
		e.writeFile(path)
		e.currentFile = path
		e.win.SetTitle("GoPad - " + w.URI().Name())
	}, e.win)
}

func (e *Editor) writeFile(path string) {
	if err := os.WriteFile(path, []byte(e.text.Text), 0644); err != nil {
		dialog.ShowError(err, e.win)
	}
}
