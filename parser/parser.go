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
			Instructions: []string{},
		},
		tokens: tokens,
	}, nil
}

func (p *Parser) Parse() (*Recipe, error) {
	p.instructions()
	if err := p.ingredients(); err != nil {
		return nil, err
	}
	return &p.Recipe, nil
}

func (p *Parser) instructions() {
	for _, t := range p.tokens {
		if t.Kind == token.INSTRUCTION {
			p.Recipe.Instructions = append(p.Recipe.Instructions, t.Literal)
		}
	}
}

func (p *Parser) ingredients() error {
	for _, input := range p.Recipe.Instructions {
		i := 0
		for i < len(input) {
			if input[i] == '{' && (input[i-1] != '&' && input[i-1] != 't') {
				start := i + 1
				end := strings.Index(string(input[start:]), "}") + start
				if end == -1 {
					return errors.New("unclosed {")
				}
				element := input[start:end]
				i = end + 1

				var unitAmount ingredient
				if i < len(input) && input[i] == '(' {
					start = i + 1
					end = strings.Index(string(input[start:]), ")") + start
					if end == -1 {
						return errors.New("unclosed )")
					}
					unitAmount.unit = input[start:end]
					i = end + 1
				}
				p.Recipe.Ingredients[element] = append(p.Recipe.Ingredients[element], ingredient{
					unit:   unitAmount.unit,
					amount: 0,
				})
			}
			i++
		}
	}
	return nil
}
