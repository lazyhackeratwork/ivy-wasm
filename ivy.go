package main

import (
	"syscall/js"

	"lazyhackergo.com/browser"
	"robpike.io/ivy/config"
	"robpike.io/ivy/exec"
	"robpike.io/ivy/mobile"
	"robpike.io/ivy/value"
)

var signal = make(chan int)

var (
	conf    config.Config
	context value.Context
)

func main() {

	ivy := js.NewEventCallback(false, false, false, cbRunIvy)
	defer ivy.Close()

	window := browser.GetWindow()

	window.Document.GetElementById("expression").AddEventListener(browser.EventKeyUp, ivy)
	f := window.Document.GetElementById("expression")
	window.Document.GetElementById("loadspinner").SetAttribute("class", "")
	f.Focus()

	// TODO: Pull the values from the UI
	conf.SetFormat("")
	conf.SetMaxBits(1e9)
	conf.SetMaxDigits(1e4)
	conf.SetOrigin(1)
	conf.SetPrompt("")

	context = exec.NewContext(&conf)

	keepalive()
}

func cbRunIvy(e js.Value) {

	println("running Ivy")
	window := browser.GetWindow()
	if e.Get("keyCode").Int() == 13 {

		expr := window.Document.GetElementById("expression").Value()
		res, err := mobile.Eval(expr)

		if err != nil {
			window.Console.Warn(err.Error())
			return
		}
		element := window.Document.GetElementById("ivy-out")
		a := element.InnerHTML()
		element.SetInnerHTML(a + "> " + expr + "<br/>" + res + "<br/>")
		f := window.Document.GetElementById("expression")
		f.SetValue("")
		window.ScrollTo(0, window.InnerHeight())
	}
}

func keepalive() {
	for {
		m := <-signal
		if m == 0 {
			println("quit signal received")
			break
		}
	}
}
