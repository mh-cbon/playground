// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

func main() {

	var opts = struct {
		windowsize string
		browser    string
		out        string
		action     string
		timeout    time.Duration
		sleep      time.Duration
	}{}

	flag.StringVar(&opts.windowsize, "windowsize", "800*600", "the window size, width*height")
	flag.StringVar(&opts.browser, "browser", "", "the browser target, one of chrome|edge (os dependant)")
	flag.StringVar(&opts.out, "out", "-", "out destination")
	flag.StringVar(&opts.action, "action", "initial", "action to do, one of initial|checkbox")
	flag.DurationVar(&opts.timeout, "timeout", time.Second*5*3, "the overall timeout of the process")
	flag.DurationVar(&opts.sleep, "sleep", time.Second, "duration of a peause inejected before waiting and screenshooting the selector (required by the engine, sometimes...)")

	flag.Parse()

	// create context
	ctxt, cancel := context.WithTimeout(context.Background(), opts.timeout)
	defer cancel()

	// create chrome instance
	var browser *chromedp.CDP
	{
		browserOpts := []chromedp.Option{
			chromedp.WithLog(log.Printf),
		}
		{
			if y := getBrowser(opts.browser); y != nil {
				r, err := runner.New(y)
				if err != nil {
					log.Fatal(err)
				}
				browserOpts = append(browserOpts, chromedp.WithRunner(r))
			}
		}
		c, err := chromedp.New(ctxt, browserOpts...)
		if err != nil {
			log.Fatal(err)
		}
		browser = c
	}

	// run task list
	var buf []byte

	if opts.action == "initial" {
		err := browser.Run(ctxt, screenshotInitialLoading(`http://localhost:8080/`, opts.windowsize, opts.sleep, &buf))
		if err != nil {
			log.Fatal(err)
		}
	} else if opts.action == "checkbox" {
		err := browser.Run(ctxt, screenshotCheckbox(`http://localhost:8080/`, opts.windowsize, opts.sleep, &buf))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		panic("nop")
	}

	// shutdown chrome

	if err := browser.Shutdown(ctxt); err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	if err := browser.Wait(); err != nil {
		log.Fatal(err)
	}

	var dst io.Writer
	if opts.out == "-" {
		dst = os.Stdout
	} else if f, err := os.Create(opts.out); err != nil {
		log.Fatal(err)
	} else {
		defer f.Close()
		dst = f
	}

	dst.Write(buf)
}

func screenshotInitialLoading(urlstr, windowsize string, sleep time.Duration, res *[]byte) chromedp.Tasks {

	var y *runtime.RemoteObject
	setWindowSize := ` console.log("window size not defined")`
	if windowsize != "" {
		f := strings.Split(windowsize, "*")
		if len(f) != 2 {
			log.Fatal("invalid window size, should be height*width, got", windowsize)
		}
		setWindowSize = fmt.Sprintf(`document.body.style.width="%vpx";document.body.style.height="%vpx";`, f[0], f[1])
	}
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(sleep),
		chromedp.Evaluate(setWindowSize, &y),
		chromedp.Sleep(sleep),
		chromedp.Screenshot("body>div", res, chromedp.NodeReady, chromedp.ByQuery),
	}
}

func screenshotCheckbox(urlstr, windowsize string, sleep time.Duration, res *[]byte) chromedp.Tasks {

	var y *runtime.RemoteObject
	setWindowSize := ` console.log("window size not defined")`
	if windowsize != "" {
		f := strings.Split(windowsize, "*")
		if len(f) != 2 {
			log.Fatal("invalid window size, should be height*width, got", windowsize)
		}
		setWindowSize = fmt.Sprintf(`document.body.style.width="%vpx";document.body.style.height="%vpx";`, f[0], f[1])
	}
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(sleep),
		chromedp.Evaluate(setWindowSize, &y),
		chromedp.Click("#importsBox", chromedp.ByQuery),
		chromedp.Sleep(sleep),
		chromedp.Screenshot("body>div", res, chromedp.NodeReady, chromedp.ByQuery),
	}
}
