package console

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
)

type Console struct {
	clear map[string]func() error
}

func NewConsole() *Console {
	cleaners := map[string]func() error{
		"darwin": func() error {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			return cmd.Run()
		},
		"linux": func() error {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			return cmd.Run()
		},
		"windows": func() error {
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			return cmd.Run()
		},
	}

	return &Console{
		clear: cleaners,
	}
}

func (c Console) Clear() error {
	value, ok := c.clear[runtime.GOOS]
	if ok {
		return value()
	}

	return errors.New("os not found for cleaning")
}
