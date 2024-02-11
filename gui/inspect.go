package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/yieldray/middleman/impl"
)

var isInspectRunning = false

type launchInspectOptions struct {
	dbPath, addr, key, cert string
}

var launchInspectChan = make(chan launchInspectOptions, 1)
var shutdownInspectChan = make(chan any, 1)

func init() {
	var shutdown func() = nil

	go func() {
		for {
			launchInspectOptions := <-launchInspectChan
			var fatalErrorChan chan error
			fatalErrorChan, shutdown = impl.Inspect(
				launchInspectOptions.dbPath, launchInspectOptions.addr,
				launchInspectOptions.key, launchInspectOptions.cert)

			go func() {
				if err := <-fatalErrorChan; err != nil {
					dialog.ShowError(err, topWindow)
					return
				}
				shutdown()
			}()
		}
	}()

	go func() {
		for {
			<-shutdownInspectChan
			if shutdown != nil {
				shutdown()
			}
		}
	}()
}

func makeInspectTab(w fyne.Window) fyne.CanvasObject {
	a := fyne.CurrentApp()
	a.Lifecycle().SetOnStopped(func() {
		shutdownInspectChan <- true
	})

	label := widget.NewLabelWithStyle("Middleman", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})

	entryDbPath := widget.NewEntry()
	entryDbPath.SetText(a.Preferences().StringWithFallback("DbPath", "./middleman.db"))

	entryHttpAddr := widget.NewEntry()
	entryHttpAddr.SetText(a.Preferences().StringWithFallback("HttpAddr", "127.0.0.1:9981"))

	entryKeyPath := widget.NewEntry()
	entryKeyPath.SetText(a.Preferences().StringWithFallback("KeyPath", "./ca.key"))

	entryCertPath := widget.NewEntry()
	entryCertPath.SetText(a.Preferences().StringWithFallback("CertPath", "./ca.crt"))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "DbPath", Widget: entryDbPath},
			{Text: "HttpAddr", Widget: entryHttpAddr},
			{Text: "keyPath", Widget: entryKeyPath},
			{Text: "certPath", Widget: entryCertPath},
		},
	}

	var onSubmit func()
	var onCancel func()
	var updateBtn = func() {
		if isInspectRunning {
			form.OnSubmit = nil
			form.OnCancel = onCancel
			label.Text = "Middleman now running..."
		} else {
			form.OnSubmit = onSubmit
			form.OnCancel = nil
			label.Text = "Middleman now stopped."
		}
		label.Refresh()
		form.Refresh()
	}

	onSubmit = func() {
		dp := entryDbPath.Text
		addr := entryHttpAddr.Text
		key := entryKeyPath.Text
		cert := entryCertPath.Text

		a.Preferences().SetString("DbPath", dp)
		a.Preferences().SetString("HttpAddr", addr)
		a.Preferences().SetString("KeyPath", key)
		a.Preferences().SetString("CertPath", cert)

		launchInspectChan <- launchInspectOptions{dbPath: dp, addr: addr, key: key, cert: cert}

		isInspectRunning = !isInspectRunning
		updateBtn()
	}

	onCancel = func() {
		isInspectRunning = !isInspectRunning
		updateBtn()

		shutdownInspectChan <- true
	}

	updateBtn()

	grid := container.New(layout.NewGridLayout(1), label, form)
	return grid
}
