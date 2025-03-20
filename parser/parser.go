package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/dimfu/apron/token"
	"github.com/fatih/color"
)

type Parser struct {
	Recipe Recipe
	tokens []token.Token
}

type Recipe struct {
	Name            string
	Metadata        map[string]string
	Ingredients     map[string]ingredient
	Materials       map[string]bool
	Timer           []string
	rawInstructions []string
	Instructions    []string
}

type ingredient struct {
	amount string
	rest   string
}

func New(tokens []token.Token) (*Parser, error) {
	if len(tokens) == 0 {
		return nil, errors.New("there is nothing to parse")
	}
	return &Parser{
		Recipe: Recipe{
			Name:            "",
			Metadata:        map[string]string{},
			Ingredients:     map[string]ingredient{},
			Materials:       map[string]bool{},
			Timer:           []string{},
			rawInstructions: []string{},
		},
		tokens: tokens,
	}, nil
}

func (p *Parser) Display() {
	fmt.Println(p.Recipe.Metadata[token.NAME])
	color.New(color.Underline).Println("\nIngredients:")
	maxWidth := 20

	for _, ingredient := range p.Recipe.Ingredients {
		paddedText := fmt.Sprintf("%-*s", maxWidth, ingredient.rest)
		coloredText := color.New(color.FgBlue).Sprint(paddedText)

		fmt.Printf("  %s %s\n", coloredText, ingredient.amount)
	}

	color.New(color.Underline).Println("\nInstructions:")
	for _, step := range p.Recipe.Instructions {
		fmt.Printf("  %-*s\n", maxWidth, step)
	}

}

func (p *Parser) extractFromTokens() {
	for _, t := range p.tokens {
		if t.Kind == token.INSTRUCTION {
			p.Recipe.rawInstructions = append(p.Recipe.rawInstructions, t.Literal)
		} else if strings.HasPrefix(string(t.Kind), "META_") {
			p.Recipe.Metadata[string(t.Kind)] = strings.TrimLeft(t.Literal, " ")
		}
	}
}

func parseAmount(elements []string) (*ingredient, error) {
	if len(elements) < 2 {
		return nil, errors.New("invalid input: expected at least two elements")
	}
	var (
		found  bool
		digits []rune
		unit   []rune
	)

	for _, c := range elements[0] {
		if unicode.IsDigit(c) || c == '.' || c == '/' {
			digits = append(digits, c)
			found = true
		} else {
			if !found {
				return nil, errors.New("expecting quantity value before unit")
			}
			// only add empty space at the first index that is not an empty space
			if len(unit) == 0 && c != ' ' {
				unit = append(unit, ' ')
			}
			unit = append(unit, c)
		}
	}

	if len(digits) == 0 {
		return nil, errors.New("no quantity found")
	}

	digitStr := string(digits)
	if _, err := strconv.ParseFloat(strings.Replace(digitStr, "/", ".", 1), 64); err != nil {
		return nil, errors.New("invalid quantity format")
	}
	digitStr += fmt.Sprintf("%s", string(unit)) // add the unit after the amount digit

	return &ingredient{
		amount: digitStr,
		rest:   elements[1],
	}, nil
}

// look for possible ingredient, material, timer property inside instruction string
func (p *Parser) processInstructions(input string) (string, error) {
	var sanitizedInstruction string
	for i := 0; i < len(input); i++ {
		if input[i] != '{' && !(i > 0 && (input[i-1] == '&' || input[i-1] == 't') && input[i] == '{') {
			sanitizedInstruction += string(input[i])
			continue
		}

		var prefix byte
		if i > 0 && (input[i-1] == '&' || input[i-1] == 't') {
			prefix = input[i-1]
		}

		// get string after prefix and before postfix
		element := p.getEnclosedString(input, &i, "}")

		// collect unit amount if the syntax provide parameter after the {} syntax
		if !(prefix == '&' || prefix == 't') && i < len(input) && input[i] == '(' {
			element = color.New(color.FgBlue).Sprintf(element)
			unitAmount := p.getEnclosedString(input, &i, ")")
			ingredient, err := parseAmount([]string{unitAmount, element})
			if err != nil {
				return "", err
			}

			// TODO: should it ensure there is no duplicate for the same ingredient element?
			if _, exists := p.Recipe.Ingredients[element]; !exists {
				p.Recipe.Ingredients[element] = *ingredient
			}
		} else if prefix == '&' || prefix == 't' {
			if prefix == '&' {
				element = color.New(color.FgYellow).Sprintf(element)
			}

			if prefix == 't' {
				element = color.New(color.FgMagenta).Sprintf(element)
			}

			sanitizedText := sanitizedInstruction
			if len(sanitizedText) > 0 {
				sanitizedInstruction = sanitizedInstruction[0 : len(sanitizedText)-1]
			}
		}

		i-- // move i to the closing parentheses
		sanitizedInstruction += element
		p.classifyElement(prefix, element)
	}
	return sanitizedInstruction, nil
}

// get string inside brackets from instruction
func (p *Parser) getEnclosedString(input string, idx *int, postfix string) string {
	start := *idx + 1
	end := strings.Index(string(input[start:]), postfix) + start
	element := input[start:end]
	*idx = end + 1
	return element
}

// classifies an element as an ingredient, material, or timer
func (p *Parser) classifyElement(prefix byte, element string) {
	switch prefix {
	case '&':
		p.Recipe.Materials[element] = true
	case 't':
		p.Recipe.Timer = append(p.Recipe.Timer, element)
	}
}

func (p *Parser) Parse() (*Recipe, error) {
	p.extractFromTokens()

	for _, input := range p.Recipe.rawInstructions {
		sanitizedInput, err := p.processInstructions(input)
		if err != nil {
			return nil, err
		}
		p.Recipe.Instructions = append(p.Recipe.Instructions, strings.TrimLeft(sanitizedInput, " "))
	}

	return &p.Recipe, nil
}
