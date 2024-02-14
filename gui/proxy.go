package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/yieldray/middleman/cli/impl"
	"github.com/yieldray/middleman/cli/interceptor"
)

var isProxyRunning = false

type launchProxyOptions struct {
	server          interceptor.ProxyServer
	addr, key, cert string
}

var launchProxyChan = make(chan launchProxyOptions, 1)
var shutdownProxyChan = make(chan any, 1)

func init() {
	var shutdown func() = nil

	go func() {
		for {
			launchProxyOptions := <-launchProxyChan
			var fatalErrorChan chan error
			fatalErrorChan, shutdown = impl.Proxy(
				launchProxyOptions.server, launchProxyOptions.addr,
				launchProxyOptions.key, launchProxyOptions.cert)

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
			<-shutdownProxyChan
			if shutdown != nil {
				shutdown()
			}
		}
	}()
}

func makeProxyTab(w fyne.Window) fyne.CanvasObject {
	a := fyne.CurrentApp()
	a.Lifecycle().SetOnStopped(func() {
		shutdownProxyChan <- true
	})

	label := widget.NewLabelWithStyle("Middleman", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})

	entryProxyServer := widget.NewEntry()
	entryProxyServer.SetText(a.Preferences().StringWithFallback("ProxyServer", "https://cros.deno.dev/"))

	entryHttpAddr := widget.NewEntry()
	entryHttpAddr.SetText(a.Preferences().StringWithFallback("HttpProxyAddr", "127.0.0.1:9980"))

	entryKeyPath := widget.NewEntry()
	entryKeyPath.SetText(a.Preferences().StringWithFallback("KeyPath", "./ca.key"))

	entryCertPath := widget.NewEntry()
	entryCertPath.SetText(a.Preferences().StringWithFallback("CertPath", "./ca.crt"))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ProxyServer", Widget: entryProxyServer},
			{Text: "HttpAddr", Widget: entryHttpAddr},
			{Text: "keyPath", Widget: entryKeyPath},
			{Text: "certPath", Widget: entryCertPath},
		},
	}

	var onSubmit func()
	var onCancel func()
	var updateBtn = func() {
		if isProxyRunning {
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
		ps := entryProxyServer.Text
		addr := entryHttpAddr.Text
		key := entryKeyPath.Text
		cert := entryCertPath.Text

		a.Preferences().SetString("ProxyServer", ps)
		a.Preferences().SetString("HttpProxyAddr", addr)
		a.Preferences().SetString("KeyPath", key)
		a.Preferences().SetString("CertPath", cert)

		proxyServer, err := interceptor.NewProxyServer(ps)
		if err != nil {
			dialog.ShowError(err, w)
			return
		} else {
			launchProxyChan <- launchProxyOptions{server: proxyServer, addr: addr, key: key, cert: cert}
		}

		isProxyRunning = !isProxyRunning
		updateBtn()
	}

	onCancel = func() {
		isProxyRunning = !isProxyRunning
		updateBtn()

		shutdownProxyChan <- true
	}

	updateBtn()

	grid := container.New(layout.NewGridLayout(1), label, form)
	return grid
}
