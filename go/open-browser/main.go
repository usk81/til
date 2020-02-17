package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func main() {
	uri := "https://www.google.com"
	openbrowser(uri)
}

func openbrowser(uri string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", uri).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", uri).Start()
	case "darwin":
		err = exec.Command("open", uri).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
