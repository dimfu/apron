package token

type Kind string

type Token struct {
	Kind    Kind
	Literal string
}

const (
	// METADATA
	Name     = "Name"
	Tags     = "Tags"
	Servings = "Servings"
)

var Keywords = map[string]Kind{
	"name":     Name,
	"tags":     Tags,
	"servings": Servings,
}
