package cli

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

const (
	NeverInteractive  = -1
	OnceInteractive   = 0
	AlwaysInteractive = 1
)

type Config struct {
	TrashDir      string
	Interactivity int // -1 never, 0 once, 1 always
	Recursive     bool
	Verbose       bool
	Files         []string
}

func ParseFlags() *Config {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	dir := filepath.Join(u.HomeDir, ".Trash")

	cfg := &Config{
		TrashDir:      dir,
		Interactivity: OnceInteractive,
		Recursive:     false,
		Verbose:       false,
	}
	flag.BoolVar(&cfg.Recursive, "r", false, "this flag is ignored")
	flag.BoolVar(&cfg.Verbose, "v", false, "explain what is being done")
	force := flag.Bool("f", false, "ignore nonexistent files, never prompt")
	// because std flag package doesn't support multiple flags in one argument.
	force2 := flag.Bool("rf", false, "alias for -f")
	alwaysInteractive := flag.Bool("i", false, "prompt before every removal")
	onceInteractive := flag.Bool("I", false, "prompt once before removing any removals")

	flag.Usage = func() {
		fmt.Println("Usage: trash [options] [files...]")
		fmt.Printf("moves files to the trash (%s)\n\n", dir)
		flag.PrintDefaults()
	}

	flag.Parse()
	cfg.Files = flag.Args()
	if len(cfg.Files) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if *force || *force2 {
		cfg.Interactivity = NeverInteractive
	}

	if *alwaysInteractive {
		cfg.Interactivity = AlwaysInteractive
	} else if *onceInteractive {
		cfg.Interactivity = OnceInteractive
	}
	return cfg
}
