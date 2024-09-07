package keyword

type Type string

const (
	Name Type = "name"
)

func ValidType(t string) bool {
	return Type(t) == Name
}
