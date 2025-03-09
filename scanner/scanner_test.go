package scanner

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dimfu/apron/token"
)

func TestTokens(t *testing.T) {
	source := `
/*
ts works?
*/
>> name: Curry udon noodle soup
>> tags: vegan
// maybe?
>> servings: 4

- Heat a &{wok} on medium heat and add {olive oil}(1 tbsp) and {onion} - cook for t{4 minutes}.
- To the &{wok} add {red pepper}, {chilli}, {garlic} and {sweet corn}(160 g) - cook for another t{4 minutes}.
- Remove the veggies from the {wok} and place in a {bowl} for use later.
- To the empty &{wok}, add {butter}(60 g) (heat until melted), once melted add the plain {flour}(30g) and mix until golden for t{3 minutes}.
/*
ya it does work
- but how bout dis?
*/
- Next, add {curry powder}(2 tbsp) and garam masala - cook for t{1 minute}.
- Slowly pour in the {vegetable stock}(1.2l) making sure mix so it doesnt clump.
// - more commented instruction
- Next, add in the {soy sauce}(4 tbsp) and {tomato ketchup}(60g) – cook for t{1 minute}.
- Once the broth is finished remove it from the heat and cook the {udon noodles}(500g) according to the instructions.
- Serve in a bowl with noodles, broth, cooked veggies, sliced {spring onion}(4), pickled red cabbage, {toasted sesame seeds}(1 tbsp) and {chilli flakes}(1 tsp).
	`

	tests := []struct {
		expectedKind    token.Kind
		expectedLiteral string
	}{
		{token.NAME, "Curry udon noodle soup"},
		{token.TAGS, "vegan"},
		{token.SERVINGS, "4"},
		{token.INSTRUCTION, "Heat a &{wok} on medium heat and add {olive oil}(1 tbsp) and {onion} - cook for t{4 minutes}."},
		{token.INSTRUCTION, "To the &{wok} add {red pepper}, {chilli}, {garlic} and {sweet corn}(160 g) - cook for another t{4 minutes}."},
		{token.INSTRUCTION, "Remove the veggies from the {wok} and place in a {bowl} for use later."},
		{token.INSTRUCTION, "To the empty &{wok}, add {butter}(60 g) (heat until melted), once melted add the plain {flour}(30g) and mix until golden for t{3 minutes}."},
		{token.INSTRUCTION, "Next, add {curry powder}(2 tbsp) and garam masala - cook for t{1 minute}."},
		{token.INSTRUCTION, "Slowly pour in the {vegetable stock}(1.2l) making sure mix so it doesnt clump."},
		{token.INSTRUCTION, "Next, add in the {soy sauce}(4 tbsp) and {tomato ketchup}(60g) – cook for t{1 minute}."},
		{token.INSTRUCTION, "Once the broth is finished remove it from the heat and cook the {udon noodles}(500g) according to the instructions."},
		{token.INSTRUCTION, "Serve in a bowl with noodles, broth, cooked veggies, sliced {spring onion}(4), pickled red cabbage, {toasted sesame seeds}(1 tbsp) and {chilli flakes}(1 tsp)."},
	}

	scanner, err := New([]byte(source))
	if err != nil {
		t.Fatalf("error while scanning tokens: %v\n", err)
	}

	for i := range tests {
		if !reflect.DeepEqual(tests[i].expectedKind, scanner.Tokens[i].Kind) {
			t.Fatalf("should expecting %v as token but got %v", tests[i].expectedKind, scanner.Tokens[i])
		}
	}

	fmt.Println(scanner.Tokens)
}
