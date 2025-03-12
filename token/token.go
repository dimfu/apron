package token

type Kind string

type Token struct {
	Kind    Kind
	Literal string
}

const (
	// METADATA
	NAME     = "META_NAME"
	TAGS     = "META_TAGS"
	SERVINGS = "META_SERVINGS"

	INSTRUCTION = "INSTRUCTION"
)

var Keywords = map[string]Kind{
	"name":        NAME,
	"tags":        TAGS,
	"servings":    SERVINGS,
	"instruction": INSTRUCTION,
}
