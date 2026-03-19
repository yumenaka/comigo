package main

import (
	"fmt"
	"os"

	"github.com/yumenaka/comigo/cmd/tui"
)

func main() {
	if err := tui.Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
