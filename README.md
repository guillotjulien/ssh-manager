# SSH Manager

A very simple and opinionated SSH identity manager.

- Only support Linux, not tested on other platforms, use at your own risks!
- CLI use only, I don't plan to create a GUI at all
- Do not save passwords, there is many password managers out there, pick your poison
- Save identities in a single JSON file
    - You can version the file using git
    - If you want to drop this tool, read the JSON, and you are done!

## Usage

`ssh-manager <COMMAND>`

### Command list

- `ssh-manager ls`: List every existing identities.
- `ssh-manager add`: Add a new identity, you will be prompted for a unique name, the address of the server and the description, nothing else.
- `ssh-manager rm <IDENTITY_NAME>`: Delete an identity from config (no rollback!)
- `ssh-manager connect <IDENTITY_NAME>`: connect to an ssh identity (will prompt for password). If `<IDENTITY_NAME>` is omitted, it will list the identities and ask you to choose one.

## TODO (possibles, nothing say I will do it):

- Manage SSH keys per identity (only the path is stored)
- Allow modification of existing identity (for now, delete and recreate)
