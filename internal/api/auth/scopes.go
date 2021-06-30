package auth

type Scope string

const (
	AuthScopeApp Scope = "app"
)

func (s Scope) String() string {
	return string(s)
}
