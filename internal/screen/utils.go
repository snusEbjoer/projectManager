package screen

import (
	"bytes"
	"os/exec"

	"github.com/eiannone/keyboard"
)

func RuneKey(char string) Event {
	return Event{Key: 0, Char: char}
}

func OtherKey(key keyboard.Key) Event {
	return Event{Key: key, Char: "\x00"}
}

func Shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("zsh", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
