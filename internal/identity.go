package internal

import (
	"os"

	"golang.org/x/crypto/ssh"
)

// Identity represents an SSH identity as they are stored in the config file
type Identity struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

// Connect connects to an SSH identity or returns an error if it failed
func (i *Identity) Connect(password string) error {
	sshconfig := &ssh.ClientConfig{
		User:            i.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	// TODO: Manage auth with SSH key

	connection, err := ssh.Dial("tcp", i.Address, sshconfig)
	if err != nil {
		return err
	}
	defer connection.Close()

	session, err := connection.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return err
	}

	// Redirect all std to user ones
	session.Stdout = os.Stdout
	session.Stdin = os.Stdin
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		return err
	}

	// FIXME: Need to allow calling back a command
	// FIXME: Need to allow left/right arrow
	// FIXME: Need to allow CTRL+L to clear
	// https://stackoverflow.com/questions/39018039/interactive-secure-shell-in-golang-not-capturing-all-keyboard
	// Perhaps you need to capture the signals in your Go program and send them over to the remote host using ssh.Session's method func (*Session) Signal
	session.Wait()

	return nil
}
