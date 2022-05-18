package trash

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hsblhsn/trash/cli"
)

type Trash struct {
	cfg *cli.Config
}

func New(cfg *cli.Config) *Trash {
	return &Trash{
		cfg: cfg,
	}
}

// Run the trash command.
// It checks if the trash directory exists and if not, it creates it.
// Then it moves the files to the trash directory one by one.
//
func (t *Trash) Run() error {
	if err := t.createTrashDir(); err != nil {
		return err
	}
	if t.cfg.Interactivity == cli.OnceInteractive {
		if !t.askForConfirmation(strings.Join(t.cfg.Files, ", ")) {
			return errors.New("aborted")
		}
	}

	var errCount int
	for _, v := range t.cfg.Files {
		if err := t.moveToTrash(v); err != nil {
			fmt.Println("error: ", err)
			errCount++
		}
	}

	if errCount > 0 {
		return fmt.Errorf("%d errors occurred", errCount)
	}
	return nil
}

// moveToTrash moves a file to the trash directory.
func (t *Trash) moveToTrash(path string) error {
	// return error if the file does not exist
	if !fileExists(path) {
		if t.cfg.Interactivity == cli.NeverInteractive {
			return nil
		}
		return fmt.Errorf("file %q does not exist", path)
	}

	// get the absolute destination trash path of the file.
	// if the file already exists in the trash,
	// try to find a unique name for the file by appending date-time to the filename.
	dst := t.getTrashPath(path)
	for fileExists(dst) {
		dst = t.getSuffixedPath(path)
		time.Sleep(time.Millisecond * 100)
	}

	// ask for confirmation before moving the file.
	if t.cfg.Interactivity == cli.AlwaysInteractive {
		if !t.askForConfirmation(path) {
			return fmt.Errorf("aborted: %q", path)
		}
	}

	// move the file to the trash.
	t.log(path)
	return os.Rename(path, dst)
}

// createTrashDir creates the trash directory.
// If the directory already exists, it is assumed to be a no op.
func (t *Trash) createTrashDir() error {
	if !fileExists(t.cfg.TrashDir) {
		t.log("trash: creating %q", t.cfg.TrashDir)
		return os.MkdirAll(t.cfg.TrashDir, 0o755)
	}
	return nil
}

// getTrashPath returns the path to the trash directory for the given file.
func (t *Trash) getTrashPath(path string) string {
	return filepath.Join(t.cfg.TrashDir, filepath.Base(path))
}

// getSuffixedPath returns a new path with the given path and the current date and time.
func (t *Trash) getSuffixedPath(path string) string {
	now := time.Now()
	date := now.Format("2006.01.02.15.04.05")
	ext := filepath.Ext(path)
	name := strings.TrimSuffix(path, ext)
	return filepath.Join(t.cfg.TrashDir, name+"_"+date+ext)
}

// askForConfirmation asks the user for confirmation. A user must type in "y" to give consent.
func (t *Trash) askForConfirmation(path string) bool {
	fmt.Printf("prompt: move %q to trash? [y/N] ", path)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}
	return strings.ToLower(response) == "y"
}

// log prints a message to stderr if the verbose flag is set.
func (t *Trash) log(msg string, args ...interface{}) {
	if t.cfg.Verbose {
		fmt.Fprintf(os.Stderr, msg+"\n", args...)
	}
}

// fileExists returns true if the file exists.
func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
