// +build !windows

package main

import (
	"log"
	"os/exec"

	"github.com/chromedp/chromedp/runner"
)

func getBrowser(browser string) runner.CommandLineOption {
	if browser == "edge" {
		log.Fatal("edge support is not implemented")
	} else if browser == "opera" {
		return runner.HeadlessPathPort(findOpera(), 9222)
	} else if browser == "safari" {
		log.Fatal("safari support is not implemented")
	} else if browser == "firefox" {
		return firefoxHeadlessPathPort(findFirefox(), 9222)
	} else if browser == "android" {
		log.Fatal("android support is not implemented")
	}
	return nil
}

func firefoxHeadlessPathPort(path string, port int) runner.CommandLineOption {
	if path == "" {
		path, _ = exec.LookPath("headless_shell")
	}

	return func(m map[string]interface{}) error {
		m["-P"] = "screenshots"
		m["exec-path"] = path
		m["remote-debugging-port"] = port
		m["headless"] = true
		return nil
	}
}

func findFirefox() string {
	path, err := exec.LookPath(`firefox`)
	if err == nil {
		return path
	}
	return "firefox"
}

func findOpera() string {
	path, err := exec.LookPath(`opera`)
	if err == nil {
		return path
	}
	return "opera"
}
