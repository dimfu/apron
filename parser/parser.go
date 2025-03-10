package parser

import (
	"errors"
	"strings"

	"github.com/dimfu/apron/token"
)

type Parser struct {
	Recipe Recipe
	tokens []token.Token
}

type Recipe struct {
	Name         string
	Metadata     map[string]string
	Ingredients  map[string][]ingredient
	Materials    map[string]bool
	Timer        []string
	Instructions []string
}

type ingredient struct {
	unit   string
	amount int
}

func New(tokens []token.Token) (*Parser, error) {
	if len(tokens) == 0 {
		return nil, errors.New("there is nothing to parse")
	}
	return &Parser{
		Recipe: Recipe{
			Name:         "",
			Metadata:     map[string]string{},
			Ingredients:  map[string][]ingredient{},
			Materials:    map[string]bool{},
			Timer:        []string{},
			Instructions: []string{},
		},
		tokens: tokens,
	}, nil
}

func (p *Parser) Parse() (*Recipe, error) {
	if err := p.instructions(); err != nil {
		return nil, err
	}
	return &p.Recipe, nil
}

// parses every ingredients, materials and timer that is provided in the instruction token
func (p *Parser) instructions() error {
	for _, t := range p.tokens {
		if t.Kind == token.INSTRUCTION {
			p.Recipe.Instructions = append(p.Recipe.Instructions, t.Literal)
		}
	}

	for _, input := range p.Recipe.Instructions {
		i := 0
		for i < len(input) {
			// check for ingredients {}, materials &{} and timer t{}
			if input[i] == '{' || (i > 0 && (input[i-1] == '&' || input[i-1] == 't') && input[i] == '{') {
				var prefix byte
				if i > 0 && (input[i-1] == '&' || input[i-1] == 't') {
					prefix = input[i-1]
				}

				// get string after prefix and before postfix
				start := i + 1
				end := strings.Index(string(input[start:]), "}") + start
				element := input[start:end]
				i = end + 1

				// collect unit amount if the syntax provide parameter after the {} syntax
				var unitAmount ingredient
				if !(prefix == '&' || prefix == 't') && i < len(input) && input[i] == '(' {
					start = i + 1
					end = strings.Index(string(input[start:]), ")") + start
					unitAmount.unit = input[start:end]
					i = end + 1
				}

				switch prefix {
				case '&':
					p.Recipe.Materials[element] = true
					break
				case 't':
					p.Recipe.Timer = append(p.Recipe.Timer, element)
				default:
					p.Recipe.Ingredients[element] = append(p.Recipe.Ingredients[element], ingredient{
						unit:   unitAmount.unit,
						amount: 0,
					})
				}
			}
			i++
		}
	}

	return nil
}
