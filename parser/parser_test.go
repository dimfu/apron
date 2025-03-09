package parser

import (
	"fmt"
	"testing"

	"github.com/dimfu/apron/token"
)

func TestIngredients(t *testing.T) {
	p, err := New([]token.Token{
		{
			Kind:    token.INSTRUCTION,
			Literal: "Next, add {curry powder}(2 tbsp) and garam masala - cook for t{1 minute}.",
		},
		{
			Kind:    token.INSTRUCTION,
			Literal: "To the &{wok} add {red pepper}, {chilli}, {garlic} and {sweet corn}(160 g) - cook for another t{4 minutes}",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	recipe, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(recipe.Ingredients)
}
