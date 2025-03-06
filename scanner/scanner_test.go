package scanner

import (
	"reflect"
	"testing"

	"github.com/dimfu/apron/token"
)

func TestTokens(t *testing.T) {
	source := `
>> name: Curry udon noodle soup
>> tags: vegan
>> servings: 4

- Heat a &{wok} on medium heat and add {olive oil}(1 tbsp) and {onion} - cook for t{4 minutes}.
- To the &{wok} add {red pepper}, {chilli}, {garlic} and {sweet corn}(160 g) - cook for another t{4 minutes}.
- Remove the veggies from the {wok} and place in a {bowl} for use later.
- To the empty &{wok}, add {butter}(60 g) (heat until melted), once melted add the plain {flour}(30g) and mix until golden for t{3 minutes}.
- Next, add {curry powder}(2 tbsp) and garam masala - cook for t{1 minute}.
- Slowly pour in the {vegetable stock}(1.2l) making sure mix so it doesnt clump.
- Next, add in the {soy sauce}(4 tbsp) and {tomato ketchup}(60g) – cook for t{1 minute}.
- Once the broth is finished remove it from the heat and cook the {udon noodles}(500g) according to the instructions.
- Serve in a bowl with noodles, broth, cooked veggies, sliced {spring onion}(4), pickled red cabbage, {toasted sesame seeds}(1 tbsp) and {chilli flakes}(1 tsp).
	`

	tests := []struct {
		expectedKind    token.Kind
		expectedLiteral string
	}{
		{token.Name, "Curry udon noodle soup"},
		{token.Tags, "vegan"},
		{token.Servings, "4"},
	}

	scanner := New([]byte(source))
	scanner.Scan()

	for i := range tests {
		if !reflect.DeepEqual(tests[i].expectedKind, scanner.tokens[i].Kind) {
			t.Fatalf("should expecting %v as token but got %v", tests[i].expectedKind, scanner.tokens[i])
		}
	}

}
