package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tyspice/zhuzh/internal/chatgpt"
	"github.com/tyspice/zhuzh/internal/ui"
)

// set during build
var version string = "dev"

func main() {

	// flags
	versionFlag := flag.Bool("v", false, "Print version")
	flag.Parse()

	prompt := flag.Arg(0)

	if *versionFlag {
		fmt.Println(version)
		return
	}

	// If the user hasn't supplied a prompt run the ui
	// else assume they want to run the prompt inline
	if prompt == "" {
		ui.Run()
	} else {
		gptClient := chatgpt.NewClient()
		gptClient.SetInstructions("")
		res, err := gptClient.Subscribe()
		gptClient.Ask(prompt)
		for {
			select {
			case next := <-res:
				if next.Done {
					fmt.Print("\n")
					return
				}
				fmt.Print(next.Delta)
			case e := <-err:
				fmt.Fprintf(os.Stderr, "error: %v\n", e)
				return
			}
		}
	}
}
