package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
)

var seq = []byte{'m', 'u', 'l', '(', ',', ')'}

type MulReader struct {
	State  int // index of slice seq
	reader *bufio.Reader

	buf bytes.Buffer
}

func NewMulReader(b *bufio.Reader) *MulReader {
	return &MulReader{
		State:  0,
		reader: b,
	}
}

func (m *MulReader) Reset() {
	m.buf.Reset()
	m.State = 0
}

// return false if EOF
func (m *MulReader) Read() (int, int, error) {
	reader := m.reader
	var a, b int
	for {
		char, err := reader.ReadByte()
		if err != nil {
			return 0, 0, err
		}
		fmt.Println(m.State, string(char), a, b, m.buf.String())
		switch m.State {
		case 0:
			if char != 'm' {
				m.Reset()
				continue
			}
			m.State = 1

		case 1:
			if char != 'u' {
				m.Reset()
				continue
			}
			m.State = 2

		case 2:
			if char != 'l' {
				m.Reset()
				continue
			}
			m.State = 3

		case 3:
			if char != '(' {
				m.Reset()
				continue
			}
			m.State = 4

		case 4:
			if unicode.IsDigit(rune(char)) {
				err = m.buf.WriteByte(char)
				if err != nil {
					return 0, 0, err
				}
				m.State = 5
			} else {
				m.Reset()
			}

		case 5:
			if unicode.IsDigit(rune(char)) {
				// digit, same state
				err = m.buf.WriteByte(char)
				if err != nil {
					return 0, 0, err
				}
				m.State = 5
			} else if char == ',' {
				// exit number mode, assign first number, reset buffer
				a, err = strconv.Atoi(m.buf.String())
				if err != nil {
					return 0, 0, err
				}
				m.buf.Reset()
				m.State = 6
			} else {
				m.Reset()
				continue
			}

		case 6:
			if unicode.IsDigit(rune(char)) {
				// digit, same state
				err = m.buf.WriteByte(char)
				if err != nil {
					return 0, 0, err
				}
				m.State = 6
			} else if char == ')' {
				// assign second number, reset buffer
				// yield result
				b, err = strconv.Atoi(m.buf.String())
				if err != nil {
					return 0, 0, err
				}
				m.Reset()
				return a, b, nil
			} else {
				m.Reset()
			}

		}
	}
}

func main() {
	reader := NewMulReader(bufio.NewReader(os.Stdin))
	var score int
	for {
		a, b, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		score += a * b
	}
	fmt.Println(score)
}
