package scanner

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dimfu/apron/token"
)

type Scanner struct {
	Tokens   []token.Token
	keywords map[string]token.Kind
	source   []byte
	start    int
	current  int
	line     int
}

func New(source []byte) (*Scanner, error) {
	scanner := &Scanner{
		Tokens:   []token.Token{},
		keywords: token.Keywords,
		source:   source,
		current:  0,
		start:    0,
		line:     1,
	}
	if err := scanner.scan(); err != nil {
		return nil, err
	}
	return scanner, nil
}

func (s *Scanner) scan() error {
	for !s.isAtEnd() {
		s.start = s.current
		err := s.next()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Scanner) next() error {
	c := s.advance()
	switch c {
	case '>':
		if s.match('>') && !s.isAtEnd() {
			err := s.metadata()
			if err != nil {
				return err
			}
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			for !s.isAtEnd() {
				if s.peek() == '*' && s.peekNext() == '/' {
					s.advance()
					s.advance()
					break
				}
				s.advance()
			}

			if s.isAtEnd() {
				return errors.New("comment not properly closed")
			}
		}
		break
	case '-':
		for s.peek() != '\n' && !s.isAtEnd() {
			s.advance()
		}

		input := string(s.source[s.start+1 : s.current-1])
		stack := []rune{}
		for _, c := range input {
			if c == '(' || c == '{' {
				stack = append(stack, c)
			} else if c == ')' || c == '}' {
				if len(stack) == 0 {
					return fmt.Errorf("unexpected closing parentheses at line %d\n", s.line)
				}
				top := stack[len(stack)-1]
				if (top == '(' && c == ')') || (top == '{' && c == '}') {
					stack = stack[:len(stack)-1]
				} else {
					return fmt.Errorf("mismatched parentheses at line %d\n", s.line)
				}
			}
		}

		if len(stack) != 0 {
			return fmt.Errorf("missing closing parentheses at line %d", s.line)
		}

		s.Tokens = append(s.Tokens, token.Token{
			Kind:    token.INSTRUCTION,
			Literal: input,
		})
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		s.line++
		break
	}

	return nil
}

func (s *Scanner) metadata() error {
	newlineIdx := s.newlineIdx()

	if newlineIdx == -1 {
		newlineIdx = len(s.source)
		return errors.New("could not find new line")
	}

	line := string(s.source[s.current:newlineIdx])
	parts := strings.Split(line, ":")

	if len(parts) < 2 {
		return fmt.Errorf("invalid metadata at line %d", s.line)
	}

	kind, err := s.lookupKeywords(strings.TrimSpace(parts[0]))
	if err != nil {
		return err
	}

	s.Tokens = append(s.Tokens, token.Token{
		Kind:    kind,
		Literal: parts[1],
	})

	s.current = newlineIdx + 1
	return nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) newlineIdx() int {
	idx := strings.IndexByte(string(s.source[s.current:]), '\n')
	if idx == -1 {
		return -1
	} else {
		idx += s.current
	}
	return idx
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) match(target byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != target {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) lookupKeywords(identifier string) (token.Kind, error) {
	tok, ok := s.keywords[identifier]
	if ok {
		return tok, nil
	} else {
		return "", fmt.Errorf("cannot find %s keyword\n", identifier)
	}
}
