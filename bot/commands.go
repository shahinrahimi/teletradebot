package bot

import (
	"fmt"
)

const (
	START    string = "start"
	HELP     string = "help"
	INFO     string = "info"
	ADD      string = "add"
	LIST     string = "list"
	CANCEL   string = "cancel"
	EXECUTE  string = "execute"
	CHECK    string = "check"
	REMOVE   string = "remove"
	VIEW     string = "view"
	DESCRIBE string = "describe"
)

func GetGeneralCommands() []string {
	commands := []string{
		START,
		HELP,
		INFO,
		LIST,
	}
	return commands
}

func GetSpecifiedCommands() []string {
	commands := []string{
		ADD,
		CANCEL,
		EXECUTE,
		CHECK,
		REMOVE,
		VIEW,
		DESCRIBE,
	}
	return commands
}

func GetAllCommands() []string {
	commands := []string{
		START,
		HELP,
		INFO,
		ADD,
		LIST,
		CANCEL,
		EXECUTE,
		CHECK,
		REMOVE,
		VIEW,
		DESCRIBE,
	}
	return commands
}

func GetCommandHelp() string {
	var commandsStr string = ""
	gCommands := GetGeneralCommands()
	for _, c := range gCommands {
		commandsStr = commandsStr + "/" + c + "\n"
	}
	sCommands := GetSpecifiedCommands()
	for _, c := range sCommands {
		commandsStr = commandsStr + "/" + c + " [id]\n"
	}
	return fmt.Sprintf("Usage\n\n%s", commandsStr)
}
