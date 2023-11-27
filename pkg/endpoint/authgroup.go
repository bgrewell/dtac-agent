package endpoint

// AuthGroup is a type for the different types of authentication groups supported.
type AuthGroup string

// String returns the string representation of the AuthGroup.
func (a AuthGroup) String() string {
	return string(a)
}

const (
	// AuthGroupAdmin is the admin authentication group.
	AuthGroupAdmin AuthGroup = "admin"
	// AuthGroupOperator is the operator authentication group.
	AuthGroupOperator AuthGroup = "operator"
	// AuthGroupUser is the user authentication group.
	AuthGroupUser AuthGroup = "user"
	// AuthGroupGuest is the guest authentication group.
	AuthGroupGuest AuthGroup = "guest"
)
