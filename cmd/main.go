package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"sort"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/ssh-manager/internal"
	"github.com/urfave/cli/v2"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	manager, err := internal.Manager{}.New(fmt.Sprintf("%s/.ssh-manager.json", usr.HomeDir))
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:                 "ssh-manager",
		Version:              "0.0.1",
		Usage:                "",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "ls",
				Usage: "List user identities",
				Action: func(c *cli.Context) error {
					return listIdentities(manager)
				},
			},
			{
				Name:  "add",
				Usage: "Add a new identity",
				Action: func(c *cli.Context) error {
					prompt := promptui.Prompt{
						Label: "Identity Name",
					}

					name, err := prompt.Run()
					if err != nil {
						return err
					}

					prompt = promptui.Prompt{
						Label: "Username",
					}

					username, err := prompt.Run()
					if err != nil {
						return err
					}

					prompt = promptui.Prompt{
						Label: "Address",
					}

					address, err := prompt.Run()
					if err != nil {
						return err
					}

					prompt = promptui.Prompt{
						Label: "Port",
					}

					port, err := prompt.Run()
					if err != nil {
						return err
					}

					intPort, err := strconv.Atoi(port)
					if err != nil {
						return err
					}

					prompt = promptui.Prompt{
						Label: "Description",
					}

					description, err := prompt.Run()
					if err != nil {
						return err
					}

					manager.AddIdentity(name, username, address, description, intPort)

					return nil
				},
			},
			{
				Name:  "rm",
				Usage: "Remove an existing identity",
				Action: func(c *cli.Context) error {
					if c.Args().Get(0) == "" {
						return listIdentities(manager)
					}

					err := manager.RemoveIdentity(c.Args().Get(0))
					if err != nil {
						return err
					}

					fmt.Println("Removed", c.Args().Get(0))

					return nil
				},
			},
			{
				Name:  "connect",
				Usage: "Connect to an SSH identity",
				Action: func(c *cli.Context) error {
					if c.Args().Get(0) == "" {
						return listIdentities(manager)
					}

					identity, err := manager.GetIdentity(c.Args().Get(0))
					if err != nil {
						return err
					}

					err = identity.Connect()
					if err != nil {
						log.Fatal(err)
					}

					return nil
				},
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func listIdentities(manager internal.Manager) error {
	identities, err := manager.GetIdentities()
	if err != nil {
		return err
	}

	if len(identities) == 0 {
		fmt.Println("No identities configured")
		return nil
	}

	fmt.Println("SSH Identities:")
	fmt.Println("==============================")

	for _, identity := range identities {
		fmt.Println()
		fmt.Println("Name:", identity.Name)
		fmt.Println("Username:", identity.Username)
		fmt.Println("Address:", identity.Address)
		fmt.Println("Port:", identity.Port)
		fmt.Println("Description:", identity.Description)
	}

	return nil
}
