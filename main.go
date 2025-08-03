package main

import "github.com/tyspice/zhuzh/internal/ui"

func main() {

	ui.Run()

	// prompt := os.Args[1]
	// responseChan, errorChan := chatgpt.StreamResponse(prompt)

	// for {
	// 	select {
	// 	case resp, ok := <-responseChan:
	// 		if !ok {
	// 			fmt.Print("\n\nAll done!\n")
	// 			return // Channel closed, streaming finished
	// 		}
	// 		fmt.Print(resp)
	// 	case err := <-errorChan:
	// 		if err != nil {
	// 			fmt.Printf("Error: %v\n", err)
	// 			return
	// 		}
	// 	}
	// }
}
