package internal

import (
	"fmt"
	"os"
	"os/exec"
)

// Identity represents an SSH identity as they are stored in the config file
type Identity struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Address     string `json:"address"`
	Port        int    `json:"port"`
	Description string `json:"description"`
}

// Connect connects to an SSH identity or returns an error if it failed
func (i *Identity) Connect() error {
	var port string

	if i.Port != 0 {
		port = fmt.Sprintf("-p %d", i.Port)
	}

	cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", i.Username, i.Address), port)

	// Redirect all std to user ones
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
