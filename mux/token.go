package mux

import "fmt"

type Token struct {
	Kind  TokenKind
	Lit   rune
	Param string
}

type TokenKind uint8

const (
	Undefined TokenKind = iota
	Lit
	Sep
	Param
	End
)

func (kind TokenKind) String() string {
	switch kind {
	case Undefined:
		return "undefined"
	case Lit:
		return "lit"
	case Sep:
		return "sep"
	case Param:
		return "param"
	case End:
		return "end"
	default:
		return fmt.Sprintf("unk-%d", kind)
	}
}

func (tok Token) String() string {
	if tok.Param != "" {
		return fmt.Sprintf("%s='%s'", tok.Kind, tok.Param)
	}
	if tok.Lit != 0 {
		return fmt.Sprintf("%s='%c'", tok.Kind, tok.Lit)
	}
	return tok.Kind.String()
}
