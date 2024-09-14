package bot

import (
	"fmt"
)

const (
	START    string = "start"
	HELP     string = "help"
	INFO     string = "info"
	ALIAS    string = "alias"
	ADD      string = "add"
	LIST     string = "list"
	RESET    string = "reset"
	EXECUTE  string = "execute"
	CLOSE    string = "close"
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
		ALIAS,
	}
	return commands
}

func GetSpecifiedCommands() []string {
	commands := []string{
		ADD,
		RESET,
		EXECUTE,
		CLOSE,
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
		ALIAS,
		RESET,
		EXECUTE,
		CLOSE,
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
