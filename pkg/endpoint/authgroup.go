package endpoint

type AuthGroup string

func (a AuthGroup) String() string {
	return string(a)
}

const (
	AuthGroupAdmin    AuthGroup = "admin"
	AuthGroupOperator AuthGroup = "operator"
	AuthGroupUser     AuthGroup = "user"
	AuthGroupGuest    AuthGroup = "guest"
)
