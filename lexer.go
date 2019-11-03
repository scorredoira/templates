package templates

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type kind byte

const (
	text kind = iota
	expression
	unescapedExp
	code
	backtickText
)

const eof = 0

type position struct {
	line   int
	column int
}

type token struct {
	kind kind
	str  string
	pos  position
}

type lexer struct {
	pos     position
	lastPos position
	kind    kind
	reader  *bufio.Reader
	tokens  []token
	buf     bytes.Buffer
}

func newLexer(reader io.Reader) *lexer {
	return &lexer{
		pos:     position{0, 0},
		lastPos: position{0, 0},
		reader:  bufio.NewReaderSize(reader, 4096),
	}
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d", p.line, p.column)
}

func (l *lexer) Run() error {
	l.kind = text

	for {
		c := l.next()
		if c == byte(eof) {
			l.emitToken(0)
			return nil
		}

		switch c {
		case '`':
			// go back and emit the rest of the text
			l.pos.column--
			l.buf.Truncate(l.buf.Len() - 1)
			l.reader.UnreadByte()
			l.emitToken(0)

			// now emit only the backtick
			l.next()
			l.emitToken(0)

		case '<':
			if l.peek() == '%' {
				l.next()
				if l.peek() == '=' {
					l.next()
					if l.peek() == '=' {
						l.next()
						l.emitToken(4)
						l.kind = unescapedExp
					} else {
						l.emitToken(3)
						l.kind = expression
					}
				} else {
					l.emitToken(2)
					l.kind = code
				}
				if err := l.readCode(); err != nil {
					return err
				}
			}
		}
	}
}

func (l *lexer) emitToken(rightTrim int) {
	ln := l.buf.Len()

	if rightTrim > 0 {
		ln -= rightTrim
		l.buf.Truncate(ln)
	}

	if ln == 0 {
		return
	}

	l.tokens = append(l.tokens, token{l.kind, l.buf.String(), l.lastPos})
	l.lastPos = l.pos
	l.buf.Reset()
}

func (l *lexer) readCode() error {
	for {
		c := l.next()
		if c == byte(eof) {
			return l.error("", "unterminated code block")
		}

		switch c {
		case '%':
			if l.peek() == '>' {
				l.next()
				l.emitToken(2)
				l.kind = text
				return nil
			}
		}
	}
}

func (l *lexer) readNext() byte {
	ch, err := l.reader.ReadByte()
	if err == io.EOF {
		return byte(eof)
	}
	return ch
}

func (l *lexer) peek() byte {
	ch := l.readNext()
	if ch != byte(eof) {
		l.reader.UnreadByte()
	}
	return ch
}

func (l *lexer) next() byte {
	ch := l.readNext()
	switch ch {
	case '\n', '\r':
		l.newline(ch)
		ch = '\n'
		l.buf.WriteByte(ch)
	case byte(eof):
		l.pos.line = eof
		l.pos.column = 0
	default:
		l.pos.column++
		l.buf.WriteByte(ch)
	}
	return ch
}

func (l *lexer) newline(ch byte) {
	l.pos.line += 1
	l.pos.column = 0
	next := l.peek()
	if ch == '\n' && next == '\r' || ch == '\r' && next == '\n' {
		l.reader.ReadByte()
	}
}

type Error struct {
	Pos     position
	Message string
	Token   string
}

func (e *Error) Position() position {
	return e.Pos
}

func (e *Error) Error() string {
	return e.Message
}

func (l *lexer) error(tok string, msg string) *Error {
	return &Error{l.pos, msg, tok}
}
