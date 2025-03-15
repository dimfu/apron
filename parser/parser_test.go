package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dimfu/apron/token"
)

func TestParseMetadata(t *testing.T) {
	p, err := New([]token.Token{
		{
			Kind:    token.NAME,
			Literal: "Curry",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	recipe, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(recipe.Metadata)
}

func TestParseAmount(t *testing.T) {
	tests := []struct {
		amounts  []string
		expected *ingredient
		hasError bool
	}{
		{
			amounts:  []string{"20g", "chocolate"},
			expected: &ingredient{amount: "20 g", rest: "chocolate"},
			hasError: false,
		},
		{
			amounts:  []string{"chocolate", "20g"},
			expected: nil,
			hasError: true,
		},
		{
			amounts:  []string{"20 g", "chocolate"},
			expected: &ingredient{amount: "20 g", rest: "chocolate"},
			hasError: false,
		},
		{
			amounts:  []string{"1.2 g", "chocolate"},
			expected: &ingredient{amount: "1.2 g", rest: "chocolate"},
			hasError: false,
		},
		{
			amounts:  []string{".5g", "of water"},
			expected: &ingredient{amount: ".5 g", rest: "of water"},
			hasError: false,
		},
		{
			amounts:  []string{"..5g", "of water"},
			expected: nil,
			hasError: true,
		},
	}

	for _, test := range tests {
		res, err := parseAmount(test.amounts)
		if test.hasError {
			if err == nil {
				t.Errorf("expected error for input %s, but got none", test.amounts)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for input %s: %v", test.amounts, err)
			}
			if res.amount != test.expected.amount {
				t.Errorf("amount '%v' should be '%v'", res.amount, test.expected.amount)
			}
			if res.rest != test.expected.rest {
				t.Errorf("rest '%v' should be '%v'", res.rest, test.expected.rest)
			}
		}

	}
}

func TestInstructions(t *testing.T) {
	p, err := New([]token.Token{
		{
			Kind:    token.INSTRUCTION,
			Literal: "Next, add {curry powder}(2 tbsp) and {garam masala}(1 tbsp) - cook for t{1 minute}.",
		},
		{
			Kind:    token.INSTRUCTION,
			Literal: "To the &{wok} add {red pepper}, {chilli}, {garlic} and {sweet corn}(160 g) - cook for another t{4 minutes}.",
		},
		{
			Kind:    token.INSTRUCTION,
			Literal: "To the empty &{wok}, add {butter}(60 g) (heat until melted), once melted add the plain {flour}(30g) and mix until golden for t{3 minutes}.",
		},
		{
			Kind:    token.INSTRUCTION,
			Literal: "Roast for t{20 minutes}, then mix it and roast for another t{20 minutes}.",
		},
	})

	expected := []string{
		"Next, add curry powder and garam masala - cook for 1 minute.",
		"To the wok add red pepper, chilli, garlic and sweet corn - cook for another 4 minutes.",
		"To the empty wok, add butter (heat until melted), once melted add the plain flour and mix until golden for 3 minutes.",
		"Roast for 20 minutes, then mix it and roast for another 20 minutes.",
	}

	if err != nil {
		t.Fatal(err)
	}
	recipe, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(recipe.Instructions) != len(expected) {
		t.Fatalf("Expected %d instructions, got %d", len(expected), len(recipe.Instructions))
	}

	for i, instr := range recipe.Instructions {
		if instr != expected[i] {
			t.Errorf("Instruction %d mismatch:\nGot:      %q\nExpected: %q", i, instr, expected[i])
		}
	}

	if !reflect.DeepEqual(recipe.Instructions, expected) {
		t.Errorf("Instructions do not match:\nGot:      %+v\nExpected: %+v", recipe.Instructions, expected)
	}
}

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
		{
			Kind:    token.INSTRUCTION,
			Literal: "To the empty &{wok}, add {butter}(60 g) (heat until melted), once melted add the plain {flour}(30g) and mix until golden for t{3 minutes}.",
		},
	})

	tests := []struct {
		key      string
		expected *ingredient
	}{
		{
			key: "butter",
			expected: &ingredient{
				amount: "60",
				rest:   "g butter",
			},
		},
		{
			key: "curry powder",
			expected: &ingredient{
				amount: "2",
				rest:   "tbsp curry powder",
			},
		},
		{
			key: "flour",
			expected: &ingredient{
				amount: "30",
				rest:   "g flour",
			},
		},
		{
			key: "sweet corn",
			expected: &ingredient{
				amount: "160",
				rest:   "g sweet corn",
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	recipe, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		ingredient, exists := recipe.Ingredients[test.key]
		if !exists {
			t.Fatalf("ingredient key %s should be exist in map", test.key)
		}
		if !reflect.DeepEqual(ingredient, *test.expected) {
			t.Fatalf("ingredient %v should be %v", ingredient, test.expected)
		}
	}
}
