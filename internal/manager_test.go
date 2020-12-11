package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Return an instance of manager", func(t *testing.T) {
		path := "/tmp/.ssh-manager.json"
		defer os.Remove(path)

		want := Manager{path, nil}
		got, _ := Manager{}.New(path)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Create configuration file to provided path if it doesn't exist", func(t *testing.T) {
		path := "/tmp/.ssh-manager.json"
		defer os.Remove(path)

		Manager{}.New(path)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist", path)
		}
	})

	t.Run("Return an error when creating the configuration file failed", func(t *testing.T) {
		_, err := Manager{}.New("/etc/.ssh-manager.json")
		if err == nil {
			t.Error("Expected error to be returned, got nil")
		}
	})

	t.Run("Return an error when a file exists with the same name, but doesn't contain expected content", func(t *testing.T) {
		path := "/tmp/.ssh-manager.json"
		defer os.Remove(path)

		ioutil.WriteFile(path, []byte("{}"), 0755)

		_, err := Manager{}.New(path)
		if err == nil {
			t.Error("Expected error to be returned, got nil")
		}
	})
}

func TestGetIdentities(t *testing.T) {
	t.Run("Returns a list of identities", func(t *testing.T) {
		path := "/tmp/.ssh-manager.json"
		defer os.Remove(path)

		want := []Identity{
			{"a", "1.1.1.1", "description a"},
			{"b", "2.2.2.2", "description b"},
		}

		data, _ := json.Marshal(configuration{want})

		ioutil.WriteFile(path, data, 0755)

		manager, _ := Manager{}.New(path)

		identities, _ := manager.GetIdentities()
		if len(identities) == 0 {
			t.Error("Expected identities to exist")
			return
		}

		if !reflect.DeepEqual(want, identities) {
			t.Errorf("want %v, got %v", want, identities)
		}
	})

	t.Run("Returns an error if the configuration file does not exists", func(t *testing.T) {
		path := "/tmp/.ssh-manager.json"
		defer os.Remove(path)

		manager := Manager{path, nil}

		if _, err := manager.GetIdentities(); err == nil {
			t.Error("Expected error to be returned, got nil")
		}
	})

	t.Run("Return an error when configuration parsing failed", func(t *testing.T) {
		path := "/tmp/.ssh-manager.json"
		defer os.Remove(path)

		manager := Manager{path, nil}

		ioutil.WriteFile(path, []byte("{'identities': ['error!']}"), 0755)

		if _, err := manager.GetIdentities(); err == nil {
			t.Error("Expected error to be returned, got nil")
		}
	})
}

func TestGetIdentity(t *testing.T) {
	t.Run("Return Identity", func(t *testing.T) {
		path := "/tmp/.ssh-manager.json"
		defer os.Remove(path)

		manager, _ := Manager{}.New(path)

		want := Identity{"test", "1.1.1.1", "description"}
		identities := []Identity{want}

		data, _ := json.Marshal(configuration{identities})
		ioutil.WriteFile(path, data, 0755)

		got, _ := manager.GetIdentity("test")

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Return an error when identity does not exists", func(t *testing.T) {
		path := "/tmp/.ssh-manager.json"
		defer os.Remove(path)

		manager, _ := Manager{}.New(path)
		_, err := manager.GetIdentity("test")
		if err == nil {
			t.Error("Expected error to be returned, got nil")
		}
	})

	t.Run("Returns an error if the configuration file does not exists", func(t *testing.T) {
		path := "/tmp/.ssh-manager.json"
		defer os.Remove(path)

		manager := Manager{path, nil}

		if _, err := manager.GetIdentity("test"); err == nil {
			t.Error("Expected error to be returned, got nil")
		}
	})

	t.Run("Return an error when configuration parsing failed", func(t *testing.T) {
		path := "/tmp/.ssh-manager.json"
		defer os.Remove(path)

		manager := Manager{path, nil}

		ioutil.WriteFile(path, []byte("{'identities': ['error!']}"), 0755)

		if _, err := manager.GetIdentity("test"); err == nil {
			t.Error("Expected error to be returned, got nil")
		}
	})
}

func TestAddIdentity(t *testing.T) {
	t.Run("Add Identity to file", func(t *testing.T) {
		path := "/tmp/.ssh-manager.json"
		defer os.Remove(path)

		manager, _ := Manager{}.New(path)
		manager.AddIdentity("test", "1.1.1.1", "description")

		want := []Identity{{"test", "1.1.1.1", "description"}}
		if err := manager.read(path); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(want, manager.identities) {
			t.Errorf("got %v, want %v", manager.identities, want)
		}
	})

	t.Run("Returns an error when an identity exists with the same name", func(t *testing.T) {
		path := "/tmp/.ssh-manager.json"
		defer os.Remove(path)

		manager, _ := Manager{}.New(path)
		manager.AddIdentity("test", "1.1.1.1", "description")

		if err := manager.AddIdentity("test", "1.1.1.1", "description"); err == nil {
			t.Error("Expected error to be returned, got nil")
		}
	})

	t.Run("Returns an error when adding identity failed", func(t *testing.T) {
		// FIXME: How to test that? We'll probably need to mock the FS
	})

	t.Run("Returns an error if the configuration file does not exists", func(t *testing.T) {
		path := "/tmp/.ssh-manager.json"
		defer os.Remove(path)

		manager := Manager{path, nil}

		if err := manager.AddIdentity("test", "1.1.1.1", "description"); err == nil {
			t.Error("Expected error to be returned, got nil")
		}
	})
}

func TestRemoveIdentity(t *testing.T) {
	t.Run("Remove identity from file", func(t *testing.T) {
		path := "/tmp/.ssh-manager.json"
		defer os.Remove(path)

		manager, _ := Manager{}.New(path)
		manager.AddIdentity("test", "1.1.1.1", "description")

		err := manager.RemoveIdentity("test")
		if err != nil {
			t.Error(err)
		}

		if identities, _ := manager.GetIdentities(); len(identities) > 0 {
			t.Errorf("Expected to not have identities, found %v", identities)
		}
	})

	t.Run("Returns an error when identity does not exists in file", func(t *testing.T) {
		path := "/tmp/.ssh-manager.json"
		defer os.Remove(path)

		manager, _ := Manager{}.New(path)

		if err := manager.RemoveIdentity("test"); err == nil {
			t.Error("Expected error to be returned, got nil")
		}
	})

	t.Run("Returns an error when removing identity failed", func(t *testing.T) {
		// FIXME: How to test that? We'll probably need to mock the FS
	})

	t.Run("Returns an error if the configuration file does not exists", func(t *testing.T) {
		path := "/tmp/.ssh-manager.json"
		defer os.Remove(path)

		manager := Manager{path, nil}

		if err := manager.RemoveIdentity("test"); err == nil {
			t.Error("Expected error to be returned, got nil")
		}
	})
}
