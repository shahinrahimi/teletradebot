package bot

import "fmt"

const (
	START   string = "start"
	HELP    string = "help"
	ADD     string = "add"
	CANCEL  string = "cancel"
	EXECUTE string = "execute"
	CHECK   string = "check"
	REMOVE  string = "remove"
)

func GetCommandHelp() string {
	return fmt.Sprintf("Usage\n/%s\n/%s\n/%s\n/%s\n/%s\n/%s\n", HELP, ADD, CANCEL, EXECUTE, CHECK, REMOVE)
}
