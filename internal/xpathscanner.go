package internal

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type XPathScanner struct {
	expr string
	pos  int

	curr          rune
	kind          LexKind
	name          string
	prefix        string
	numval        float64
	strval        string
	canBeFunction bool
}

func (s *XPathScanner) NextChar() bool {
	if s.pos < len(s.expr) {
		s.curr = rune(s.expr[s.pos])
		s.pos += 1
		return true
	}
	s.curr = rune(0)
	return false
}

func (s *XPathScanner) NextLex() bool {
	s.skipSpace()
	switch s.curr {
	case 0:
		s.kind = LexEof
		return false
	case ',', '@', '(', ')', '|', '*', '[', ']', '+', '-', '=', '#', '$':
		s.kind, _ = lexKinds[int(s.curr)]
		s.NextChar()
	case '<':
		s.kind = LexLt
		s.NextChar()
		if s.curr == '=' {
			s.kind = LexLe
			s.NextChar()
		}
	case '>':
		s.kind = LexGt
		s.NextChar()
		if s.curr == '=' {
			s.kind = LexGe
			s.NextChar()
		}
	case '!':
		s.kind = LexBang
		s.NextChar()
		if s.curr == '=' {
			s.kind = LexNe
			s.NextChar()
		}
	case '.':
		s.kind = LexDot
		s.NextChar()
		if s.curr == '.' {
			s.kind = LexDotDot
			s.NextChar()
		} else if unicode.IsDigit(s.curr) {
			s.kind = LexNumber
			s.numval = scanFraction(s)
		}
	case '/':
		s.kind = LexSlash
		s.NextChar()
		if s.curr == '/' {
			s.kind = LexSlashSlash
			s.NextChar()
		}
	case '"', '\'':
		s.kind = LexString
		s.strval = scanString(s)
	default:
		if unicode.IsDigit(s.curr) {
			s.kind = LexNumber
			s.numval = scanNumber(s)
		} else if isElemChar(s.curr) {
			s.kind = LexName
			s.name = scanName(s)
			s.prefix = ""
			// "foo:bar" is one lexem not three because it doesn't allow spaces in between
			// We should distinct it from "foo::" and need process "foo ::" as well
			if s.curr == ':' {
				s.NextChar()
				// can be "foo:bar" or "foo::"
				if s.curr == ':' {
					// "foo::"
					s.NextChar()
					s.kind = LexAxe
				} else { // "foo:*", "foo:bar" or "foo: "
					s.prefix = s.name
					if s.curr == '*' {
						s.NextChar()
						s.name = "*"
					} else if isElemChar(s.curr) {
						s.name = scanName(s)
					} else {
						panic(fmt.Sprintf("%s has an invalid qualified name.", s.expr))
					}
				}
			} else {
				s.skipSpace()
				if s.curr == ':' {
					s.NextChar()
					// it can be "foo ::" or just "foo :"
					if s.curr == ':' {
						s.NextChar()
						s.kind = LexAxe
					} else {
						panic(fmt.Sprintf("%s has an invalid qualified name.", s.expr))
					}
				}
			}
			s.skipSpace()
			s.canBeFunction = s.curr == '('
		} else {
			panic(fmt.Sprintf("%s has an invalid token.", s.expr))
		}
	}
	return true
}

func (s *XPathScanner) skipSpace() {
	for unicode.Is(whitespace, s.curr) && s.NextChar() {
		//
	}
}

func scanFraction(s *XPathScanner) float64 {
	start := s.pos - 2
	len := 1
	for unicode.IsDigit(s.curr) {
		s.NextChar()
		len++
	}
	v, err := strconv.ParseFloat(s.expr[start:start+len], 64)
	if err != nil {
		panic(err)
	}
	return v
}

func scanNumber(s *XPathScanner) float64 {
	start := s.pos - 1
	len := 0
	for unicode.IsDigit(s.curr) {
		s.NextChar()
		len++
	}
	if s.curr == '.' {
		s.NextChar()
		len++
		for unicode.IsDigit(s.curr) {
			s.NextChar()
			len++
		}
	}
	v, err := strconv.ParseFloat(s.expr[start:start+len], 64)
	if err != nil {
		panic(err)
	}
	return v
}

func scanString(s *XPathScanner) string {
	end := s.curr
	s.NextChar()
	start := s.pos - 1
	len := 0
	for s.curr != end {
		if !s.NextChar() {
			panic(errors.New("unclosed string."))
		}
		len++
	}
	s.NextChar()
	return s.expr[start : start+len]
}

func scanName(s *XPathScanner) string {
	start := s.pos - 1
	var len int
	for len = 0; isElemChar(s.curr); s.NextChar() {
		len++
	}
	return s.expr[start : start+len]
}

func isElemChar(r rune) bool {
	return string(r) != ":" && string(r) != "/" &&
		(unicode.Is(first, r) || unicode.Is(second, r) || string(r) == "*")
}

var lexKinds = map[int]LexKind{
	',':  LexComma,
	'/':  LexSlash,
	'@':  LexAt,
	'.':  LexDot,
	'(':  LexLParens,
	')':  LexRParens,
	'[':  LexLBracket,
	']':  LexRBracket,
	'*':  LexStar,
	'+':  LexPlus,
	'-':  LexMinus,
	'=':  LexEq,
	'<':  LexLt,
	'>':  LexGt,
	'!':  LexBang,
	'$':  LexDollar,
	'\'': LexApos,
	'"':  LexQuote,
	'|':  LexUnion,
	'N':  LexNe,
	'L':  LexLe,
	'G':  LexGe,
	'A':  LexAnd,
	'O':  LexOr,
	'D':  LexDotDot,
	'S':  LexSlashSlash,
	'n':  LexName,
	's':  LexString,
	'd':  LexNumber,
	'a':  LexAxe,
	'E':  LexEof,
}
