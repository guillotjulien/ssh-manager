package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

type configuration struct {
	Identities []Identity `json:"identities"`
}

// Manager manage read and write of identity in configuration file
type Manager struct {
	path       string
	identities []Identity
}

// New create a new instance of manager and return it
func (m Manager) New(path string) (manager Manager, err error) {
	// When file does not exist, create it with an empty configuration structure
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := m.createEmptyConfiguration(path); err != nil {
			return Manager{}, err
		}
	}

	err = m.read(path)
	if err != nil {
		return Manager{}, err
	}

	return Manager{path, nil}, nil
}

// GetIdentities returns every identity found in configuration file or an error
// when the configration file was not read properly
func (m Manager) GetIdentities() (identities []Identity, err error) {
	// TODO: Speed up things by not reading from the file each time
	err = m.read(m.path)
	if err != nil {
		return nil, err
	}

	return m.identities, nil
}

// GetIdentity returns an identity from the configuration file by its name.
// When identity is not found, an error is returned.
func (m Manager) GetIdentity(name string) (identity Identity, err error) {
	// TODO: Speed up things by not reading from the file each time
	err = m.read(m.path)
	if err != nil {
		return Identity{}, err
	}

	// TODO: Use hash map instead of simple array
	for _, identity := range m.identities {
		if identity.Name == name {
			return identity, nil
		}
	}

	return Identity{}, fmt.Errorf("No identity found matching name %s", name)
}

// AddIdentity adds a new identity to the configuration file. In case an identity
// with the same name already exists, an error will be returned.
func (m *Manager) AddIdentity(name, username, address, description string, port int) error {
	var existingIdentity Identity

	// We first check if the file is here and the format is okay
	if err := m.read(m.path); err != nil {
		return err
	}

	// TODO: Use hash map instead of simple array
	for _, identity := range m.identities {
		if identity.Name == name {
			existingIdentity = identity
			break
		}
	}

	emptyIdentity := Identity{}
	if existingIdentity != emptyIdentity {
		return fmt.Errorf("An identity with name %s already exists", name)
	}

	m.identities = append(m.identities, Identity{
		name,
		username,
		address,
		port,
		description,
	})

	return m.write()
}

// RemoveIdentity removes an identity from the configuration file. In case the
// identity does not exists, or removal resulted in an error, we return the error.
func (m *Manager) RemoveIdentity(name string) error {
	// We first check if the file is here and the format is okay
	if err := m.read(m.path); err != nil {
		return err
	}

	filteredIdentities := []Identity{}

	// TODO: Use hash map instead of simple array
	for _, identity := range m.identities {
		if identity.Name != name {
			filteredIdentities = append(filteredIdentities, identity)
		}
	}

	if len(m.identities) == len(filteredIdentities) {
		return fmt.Errorf("No identity found matching name %s", name)
	}

	m.identities = filteredIdentities

	return m.write()
}

// read reads from the configuration file, creates an array of identity and pass
// it to the Manager.
func (m *Manager) read(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	config := configuration{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	// Test if the file content match what we expected
	if (reflect.DeepEqual(config, configuration{})) {
		return fmt.Errorf("Configuration file exists at path %s but doesn't have expected content", m.path)
	}

	// Assign all identities to the manager
	m.identities = config.Identities

	return nil
}

func (m *Manager) write() error {
	config := configuration{m.identities}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(m.path, data, 0755)
	if err != nil {
		return err
	}

	return nil
}

func (m Manager) createEmptyConfiguration(path string) error {
	data, err := json.Marshal(configuration{[]Identity{}})
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, data, 0755)
	if err != nil {
		return err
	}

	return nil
}
