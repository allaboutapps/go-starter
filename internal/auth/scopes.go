package auth

type Scope string

const (
	ScopeApp Scope = "app"
)

func (s Scope) String() string {
	return string(s)
}
