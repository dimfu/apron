package scanner

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/dimfu/apron/token"
)

type Scanner struct {
	tokens   []token.Token
	keywords map[string]token.Kind
	source   []byte
	current  int
	line     int
}

func New(source []byte) *Scanner {
	return &Scanner{
		tokens:   []token.Token{},
		keywords: token.Keywords,
		source:   source,
		current:  0,
		line:     1,
	}
}

func (s *Scanner) Scan() {
	for !s.isAtEnd() {
		s.next()
	}
}

func (s *Scanner) next() {
	c := s.advance()
	switch c {
	case '>':
		if s.match('>') && !s.isAtEnd() {
			err := s.metadata()
			if err != nil {
				log.Fatal(err)
			}
		}
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		s.line++
		break
	case '"':
		break
	}
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

	s.tokens = append(s.tokens, token.Token{
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
