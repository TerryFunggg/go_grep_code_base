package gui

import (
	"fmt"
	"grep_code_base/grep"
	"path/filepath"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
)

type App struct {
	*tview.Application
}

type CodeText struct {
	*tview.TextView
	Buffer *grep.Result
}

type FileList struct {
	*tview.List
	InfoView   *tview.TextView
	CodeBuffer *CodeText
	ListData   *[]grep.Result
}

type AppGrip struct {
	*tview.Grid
	CodeTextView *CodeText
	ListFileView *FileList
}

func NewApp() *App {
	return &App{
		Application: tview.NewApplication(),
	}
}

func NewCodeText() *CodeText {
	return &CodeText{
		TextView: tview.NewTextView(),
	}
}

func NewAppGrid() *AppGrip {
	return &AppGrip{
		Grid: tview.NewGrid(),
	}
}

func NewFileList() *FileList {
	return &FileList{
		List: tview.NewList(),
	}
}
func NewInfoView(text string) *tview.TextView {
	return tview.NewTextView().
		SetTextAlign(tview.AlignRight).
		SetText(text)
}

func NewHelpView(text string) *tview.TextView {
	return tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText(text)
}

// Overwrite the default List handler
// cus I don't want Left/Right key react
func (r *FileList) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return r.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {

		previousItem := r.GetCurrentItem()
		var current int

		switch key := event.Key(); key {
		case tcell.KeyDown:
			current = previousItem + 1
			if current > r.GetItemCount()-1 {
				current = 0
			}

			r.SetCurrentItem(current)
			r.CodeBuffer.Clear()

			data := *r.ListData
			r.CodeBuffer.Buffer = &data[current]
			r.InfoView.SetText(data[current].Path)
			fmt.Fprint(r.CodeBuffer, string(data[current].Body))

		case tcell.KeyUp:
			current = previousItem - 1
			if current < 0 {
				current = r.GetItemCount() - 1
			}

			r.SetCurrentItem(current)
			r.CodeBuffer.Clear()

			data := *r.ListData
			r.CodeBuffer.Buffer = &data[current]
			r.InfoView.SetText(data[current].Path)
			fmt.Fprint(r.CodeBuffer, string(data[current].Body))
		}
	})
}

func Show(result *[]grep.Result) {

	err := clipboard.Init()

	if err != nil {
		// clipboard error, but still can provide other function
		fmt.Println("Clipboard error", err)
	}

	app := NewApp()
	grid := NewAppGrid()

	helpStr := "Left/Right: List/Code Area, Ctrl-y: Yank Code"
	helpView := NewHelpView(helpStr)
	infoView := NewInfoView("")

	codeView := NewCodeText()
	codeView.
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	// init fist buffer
	firstItem := (*result)[0]
	fmt.Fprint(codeView, string(firstItem.Body))
	codeView.SetBorder(true)

	index := 0

	fileList := func(files *[]grep.Result) *FileList {
		listView := NewFileList()

		r := 'a'
		for _, result := range *files {
			dir, file := filepath.Split(result.Path)
			dir = strings.TrimRight(dir, "/")
			dirs := strings.Split(dir, "/")
			subDir := dirs[len(dirs)-1]
			listView.AddItem(file, subDir, '-', nil)
			r++
		}

		return listView
	}

	listView := fileList(result)
	listView.ListData = result
	listView.CodeBuffer = codeView
	listView.InfoView = infoView
	listView.SetBorder(true)

	codeView.SetBorder(true)

	grid.SetRows(0, 1).
		SetColumns(30, 0)

	grid.AddItem(listView, 0, 0, 1, 1, 0, 100, true).
		AddItem(codeView, 0, 1, 1, 2, 0, 100, false).
		AddItem(helpView, 1, 0, 1, 1, 0, 0, false).
		AddItem(infoView, 1, 1, 1, 2, 0, 0, false)

	grid.CodeTextView = codeView
	grid.ListFileView = listView

	pages := tview.NewPages().
		AddPage("main", grid, true, true)

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			app.SetFocus(listView)
		case tcell.KeyRight:
			app.SetFocus(codeView)
		}

		return event
	})

	yankedText := "Yank into clipboard!"
	codeView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlY:
			go func() {
				infoView.SetText(yankedText)
				time.Sleep(time.Second * 1)
				infoView.SetText(codeView.Buffer.Path)
			}()

			clipboard.Write(clipboard.FmtText, (*result)[index].Body)
		}

		return event
	})

	if err := app.SetRoot(pages, true).SetFocus(grid).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}
