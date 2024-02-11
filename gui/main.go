package gui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var navBar = map[string]func(w fyne.Window) fyne.CanvasObject{
	"proxy":   makeProxyTab,
	"inspect": makeInspectTab,
}

func makeNav(setTab func(uid string, t func(w fyne.Window) fyne.CanvasObject)) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			if uid == "" {
				return []string{"proxy", "inspect"}
			}
			return []string{}
		},
		IsBranch: func(uid string) bool {
			return true
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			_, ok := navBar[uid]
			if !ok {
				fyne.LogError("Missing tab panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(uid)
			obj.(*widget.Label).TextStyle = fyne.TextStyle{}
		},
		OnSelected: func(uid string) {
			if t, ok := navBar[uid]; ok {
				a.Preferences().SetString("navBar", uid)
				setTab(uid, t)
			}
		},
	}

	currentPref := a.Preferences().StringWithFallback("navBar", "proxy")
	tree.Select(currentPref)

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}

var topWindow fyne.Window

func Main() {
	a := app.NewWithID("io.github.yieldray.middleman")
	a.SetIcon(theme.ComputerIcon())

	w := a.NewWindow("Middleman")
	topWindow = w
	w.SetMaster()

	content := container.NewStack()
	title := widget.NewLabel("Component name")

	setTab := func(uid string, t func(w fyne.Window) fyne.CanvasObject) {
		if fyne.CurrentDevice().IsMobile() {
			child := a.NewWindow(uid)
			topWindow = child
			child.SetContent(t(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = w
			})
			return
		}

		title.SetText(uid)

		content.Objects = []fyne.CanvasObject{t(w)}
		content.Refresh()
	}

	tabSide := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()), nil, nil, nil, content)

	split := container.NewHSplit(makeNav(setTab), tabSide)
	split.Offset = 0.2
	w.SetContent(split)

	// menu

	menu := fyne.NewMainMenu(fyne.NewMenu("About", fyne.NewMenuItem("Github", func() {
		u, _ := url.Parse("https://github.com/YieldRay/middleman")
		_ = a.OpenURL(u)
	})))

	w.SetMainMenu(menu)

	// // system tray
	// if desk, ok := a.(desktop.App); ok {
	// 	h := fyne.NewMenuItem("Hi", func() {})
	// 	h.Icon = theme.HomeIcon()
	// 	menu := fyne.NewMenu("Menu", h)
	// 	h.Action = func() {
	// 		h.Label = "Welcome"
	// 		menu.Refresh()
	// 	}
	// 	desk.SetSystemTrayMenu(menu)
	// }

	// run
	w.Resize(fyne.NewSize(640, 460))
	w.ShowAndRun()
}
