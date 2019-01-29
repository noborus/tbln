package tbln

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
)

type Attribute struct {
	Name string
	Type string
}

type Decoder struct {
	s         *Scanner
	comment   []string
	ext       map[string]string
	Attribute []Attribute
	Ncol      int
	Nrow      int
}

func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{
		s:   NewScanner(reader),
		ext: make(map[string]string),
	}
}

func (dec *Decoder) Decode() ([]string, error) {
	s := dec.s
	for {
		t, err := s.Scan()
		if err != nil {
			if err == io.EOF {
				return nil, err
			}
			log.Printf("%s", err)
			continue
		}
		switch t {
		case Record:
			if dec.Ncol <= 0 {
				dec.Ncol = len(s.Record)
			}
			if dec.Ncol == len(s.Record) {
				dec.Nrow++
				return s.Record, nil
			} else {
				return nil, fmt.Errorf("number of column is invalid")
			}
		case Comment:
			dec.comment = append(dec.comment, s.Comment)
		case Extra:
			dec.analyzeExt(s.Extra)
		}
	}
	return nil, fmt.Errorf("row invalid")
}

func (dec *Decoder) analyzeExt(ext []string) error {
	escrep := regexp.MustCompile(`\|(\|+)`)
	switch ext[0] {
	case "name":
		body := strings.TrimRight(ext[1][2:], " |")
		rec := strings.Split(body, " | ")
		// Unescape vertical bars || -> |
		dec.Attribute = make([]Attribute, len(rec))
		for i, column := range rec {
			if strings.Contains(column, "|") {
				rec[i] = escrep.ReplaceAllString(column, "$1")
			}
			dec.Attribute[i].Name = column
		}
	case "type":
		body := strings.TrimRight(ext[1][2:], " |")
		rec := strings.Split(body, " | ")
		// Unescape vertical bars || -> |
		for i, column := range rec {
			if strings.Contains(column, "|") {
				rec[i] = escrep.ReplaceAllString(column, "$1")
			}
			dec.Attribute[i].Type = column
		}
	}
	dec.ext[ext[0]] = ext[1]
	return nil
}
