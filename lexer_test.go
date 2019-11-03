package templates

import (
	"fmt"
	"strings"
	"testing"
)

func TestLex(t *testing.T) {
	data := []struct {
		s     string
		t     []token
		isErr bool
	}{
		{
			`Foo"<%= a %>":`,
			[]token{
				{kind: text, str: `Foo"`},
				{kind: expression, str: " a "},
				{kind: text, str: `":`},
			},
			false,
		},
		{
			"<% print(\"foo\") %>",
			[]token{
				{kind: code, str: " print(\"foo\") "},
			},
			false,
		},
		{
			"test <%= foo  %>bar<%== expr  %><% qw3  %>END",
			[]token{
				{kind: text, str: "test "},
				{kind: expression, str: " foo  "},
				{kind: text, str: "bar"},
				{kind: unescapedExp, str: " expr  "},
				{kind: code, str: " qw3  "},
				{kind: text, str: "END"},
			},
			false,
		},
		{
			"test <%= foo",
			[]token{},
			true,
		},
	}

	for i, d := range data {
		if err := test(d.s, d.t, d.isErr); err != nil {
			t.Fatalf("%d) %v", i, err)
		}
	}
}

func test(s string, tokens []token, isErr bool) error {
	l := newLexer(strings.NewReader(s))
	if err := l.Run(); err != nil {
		if isErr {
			return nil // the error was expected
		}
		return err
	}

	if len(tokens) != len(l.tokens) {
		return fmt.Errorf("found %d tokens, expected %d", len(l.tokens), len(tokens))
	}

	for i, t := range tokens {
		lt := l.tokens[i]
		if lt.kind != t.kind {
			return fmt.Errorf("%d. found %v, expected %v", i, lt.kind, t.kind)
		}
		if lt.str != t.str {
			return fmt.Errorf("%d. found %s, expected %s", i, lt.str, t.str)
		}
	}

	return nil
}
