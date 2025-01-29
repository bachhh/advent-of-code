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

// read until the char matched the seq, return true if fully matched
// return false + index of the mismatched byte + last byte read, if input is not fully matched
func (m *MulReader) readSeq(seq []byte) (bool, error) {
	reader := m.reader
	for i := range seq {
		char, err := reader.ReadByte()
		if err != nil {
			return false, err
		}
		if char != seq[i] {
			return false, nil
		}
	}
	return true, nil
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

		case 0: // default state
			switch char {
			case 'm':
				m.Reset()
				m.State = 1
			case 'd':
				m.State = 7
			default:
				m.Reset()
			}

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
		case 7:
			switch char {
			case 'o':
				m.State = 8
			default:
				m.Reset()
			}
		case 8:
			switch char {
			case 'n':
				m.State = 9
			default:
				m.Reset()
			}
		case 9:
			switch char {
			case '\'':
				m.State = 10
			default:
				m.Reset()
			}
		case 10:
			switch char {
			case 't':
				m.State = 11
			default:
				m.Reset()
			}
		case 11:
			switch char {
			case '(':
				m.State = 12
			default:
				m.Reset()
			}

		case 12:
			switch char {
			case ')':
				m.State = 13
			default:
				m.Reset()
			}

		case 13: // trap state
			switch char {
			case 'd':
				m.State = 14
			default:
				m.State = 13
			}

		case 14:
			switch char {
			case 'o':
				m.State = 15
			default:
				m.State = 13
			}

		case 15:
			switch char {
			case '(':
				m.State = 16
			default:
				m.State = 13
			}

		case 16:
			switch char {
			case ')':
				m.Reset() // reset back to default state

			default:
				m.State = 13
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

// don't()d
// don't()d
// do()d
// do()d
