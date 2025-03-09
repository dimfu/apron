package token

type Kind string

type Token struct {
	Kind    Kind
	Literal string
}

const (
	// METADATA
	NAME     = "Name"
	TAGS     = "Tags"
	SERVINGS = "Servings"

	INSTRUCTION = "Instruction"
)

var Keywords = map[string]Kind{
	"name":        NAME,
	"tags":        TAGS,
	"servings":    SERVINGS,
	"instruction": INSTRUCTION,
}
