package user

type Role int

const (
	Owner Role = iota
	Guest
)

var RoleNames = map[Role]string{
	Owner: "owner",
	Guest: "guest",
}

var RoleValues = map[string]Role{
	"owner": Owner,
	"guest": Guest,
}

func (r Role) String() string {
	if name, ok := RoleNames[r]; ok {
		return name
	}
	return "unknown"
}
