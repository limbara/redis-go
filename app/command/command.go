package command

const (
	CommandArgTypeKey       CommandArgType = "key"
	CommandArgTypeString    CommandArgType = "string"
	CommandArgTypePureToken CommandArgType = "pure-token"
	CommandArgTypeInteger   CommandArgType = "integer"
	CommandArgTypeUnixTime  CommandArgType = "unix-time"
)

type CommandArgType string

type CommandArg struct {
	name    string
	argType CommandArgType
	value   any
}

type CommandArgTokenized struct {
	token string
	CommandArg
}

type CommandExecutor interface {
	Execute() (data any, e error)
}

type Command[Arg CommandArg | CommandArgTokenized] struct {
	name string
	args []Arg
}
