package main

import (
	"flag"
	"fmt"

	"github.com/tyspice/zhuzh/internal/ui"
)

// set during build
var version string = "dev"

func main() {
	// flags
	versionFlag := flag.Bool("v", false, "Print version")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		return
	}

	ui.Run()
}
