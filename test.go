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
	q := js.NewEventCallback(false, false, false, cbQuit)
	defer q.Close()

	ivy := js.NewEventCallback(false, false, false, cbRunIvy)
	defer ivy.Close()

	window := browser.GetWindow()
	window.Document.GetElementById("runButton").SetAttribute("disabled", true)
	window.Document.GetElementById("quit").AddEventListener(browser.EventClick, q)
	window.Document.GetElementById("quit").SetProperty("disabled", false)

	window.Document.GetElementById("ivy").AddEventListener(browser.EventClick, ivy)
	window.Document.GetElementById("ivy").SetProperty("disabled", false)

	// TODO: Pull the values from the UI
	conf.SetFormat("")
	conf.SetMaxBits(1e9)
	conf.SetMaxDigits(1e4)
	conf.SetOrigin(1)
	conf.SetPrompt("")

	context = exec.NewContext(&conf)

	keepalive()
}

func cb(args []js.Value) {
	println("callback")
}

func cbQuit(e js.Value) {
	println("got Quit event callback!")
	window := browser.GetWindow()
	window.Document.GetElementById("runButton").SetProperty("disabled", false)
	window.Document.GetElementById("quit").SetAttribute("disabled", true)
	signal <- 0
}

func cbRunIvy(e js.Value) {

	println("running Ivy")
	window := browser.GetWindow()
	expr := window.Document.GetElementById("expression").Value()
	res, err := mobile.Eval(expr)

	if err != nil {
		window.Console.Warn(err.Error())
		return
	}
	element := window.Document.GetElementById("ivy-out")
	element.InnerHTML(res)

}

func keepalive() {
	for {
		m := <-signal
		if m == 0 {
			println("quit signal received")
			break
		}
	}
	// select {} also seems to work but the following doesn't:
	// select {
	//    case m <-signal:
	//       // do something
	//    default:
	//       // wait
	// }
}
