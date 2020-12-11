package internal

// Identity represents an SSH identity as they are stored in the config file
type Identity struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

// Connect connects to an SSH identity or returns an error if it failed
func (i *Identity) Connect(password string) error {
	return nil
}
