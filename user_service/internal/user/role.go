package user

type Role string

const (
	Owner Role = "owner"
	Guest Role = "guest"
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
