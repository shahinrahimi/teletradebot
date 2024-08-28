package bot

import "fmt"

const (
	START    string = "start"
	HELP     string = "help"
	ADD      string = "add"
	LIST     string = "list"
	CANCEL   string = "cancel"
	EXECUTE  string = "execute"
	CHECK    string = "check"
	REMOVE   string = "remove"
	VIEW     string = "view"
	DESCRIBE string = "describe"
)

func GetCommandHelp() string {
	return fmt.Sprintf("Usage\n/%s\n/%s\n/%s\n/%s\n/%s\n/%s\n/%s\n/%s\n/%s", HELP, ADD, LIST, CANCEL, EXECUTE, CHECK, REMOVE, VIEW, DESCRIBE)
}
