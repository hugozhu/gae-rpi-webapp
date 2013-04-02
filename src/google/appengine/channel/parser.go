package channel

import (
	"strings"
)

type Element struct {
	Key   string
	Val   *[]Element
	Child *Element
}

func (e *Element) HasChild() bool {
	return e.Child != nil
}

func (e *Element) ToString() string {
	s := ""
	if e.Key != "" {
		if e.Child == nil && e.Val == nil {
			s = "[" + e.Key
		} else {
			s = "[\"" + e.Key + "\", "
		}

	}
	if e.Child != nil {
		s = s + e.Child.ToString()
	}
	if e.Val != nil {
		for _, ee := range *e.Val {
			s = s + ee.ToString()
		}
	}
	s = s + "]"
	return s
}

type Parser struct {
	parent *Element
	arr    *[]Element
}

func (p *Parser) Root() *Element {
	return p.parent
}

func (p *Parser) Parse(bytes []byte) *Element {
	p.parse(bytes, 0)
	return p.Root()
}

func (p *Parser) parse(bytes []byte, pos1 int) int {
	if p.arr == nil {
		p.arr = &[]Element{}
	}
	key := ""
	for i := pos1; i < len(bytes); i++ {
		b := bytes[i]
		switch b {
		case '[':
			if bytes[i+1] == ']' {
				i = i + 1
			} else {
				key = strings.Trim(string(bytes[pos1:i]), "\"',\n")
				pos1 = i
				i = p.parse(bytes, i+1)
			}
		case ']':
			e := Element{
				Key: string(bytes[pos1:i]),
			}

			if key != "" {
				if key[0:1] == "[" { //ignore
					return i
				}
				p.parent = &Element{
					Key:   key,
					Val:   p.arr,
					Child: p.parent,
				}
				p.arr = &[]Element{}
			} else {
				tmp := append(*p.arr, e)
				p.arr = &tmp
			}
			return i
		}
	}
	return len(bytes)
}
