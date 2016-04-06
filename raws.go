package main

import (
	"bufio"
	"io"
	"strings"
)

type RawsTokenizer struct {
	r *bufio.Reader
}

func NewRawsTokenizer(r io.Reader) *RawsTokenizer {
	return &RawsTokenizer{bufio.NewReader(r)}
}

func (t *RawsTokenizer) Next() ([]string, error) {
	_, err := t.r.ReadString('[')
	if err != nil {
		return nil, err
	}
	s, err := t.r.ReadString(']')
	if err != nil {
		return nil, err
	}

	s = s[:len(s)-1]

	return strings.Split(s, ":"), nil
}
