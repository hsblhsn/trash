package main

import (
	"fmt"
	"os"

	"github.com/hsblhsn/trash"
	"github.com/hsblhsn/trash/cli"
)

func main() {
	cfg := cli.ParseFlags()
	t := trash.New(cfg)
	if err := t.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "trash: %v\n", err)
		os.Exit(1)
	}
}
