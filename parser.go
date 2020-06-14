package fastparse

import (
	"errors"
	"io"
)

// Parser is used to define an argument parser from the ParserManager. Note this (and functions on provided arguments) are not thread safe.
type Parser struct {
	initProperly bool
	pad *scratchpad
	r io.ReadSeeker
	argId int
}

// Done should be called when you are done with a scratchpad.
func (p *Parser) Done() {
	if !p.initProperly {
		panic("Parser objects should be created from a ParserManager")
	}
	p.pad.readd()
}

// Remainder is used to get the remainder of the reader as a string, ignoring arguments.
func (p *Parser) Remainder() (string, error) {
	if !p.initProperly {
		return "", errors.New("parser objects should be created from a ParserManager")
	}
	defer p.pad.reset()
	n, err := p.r.Read(p.pad.slice)
	if err != nil {
		return "", err
	}
	p.pad.len = n
	return p.pad.String(), nil
}

// Gets the raw argument information which we can use to create the Argument struct.
func (p *Parser) getRawArgInfo() (string, int) {
	defer p.pad.reset()
	raw := 0
	first := true
	quote := false
	ob := make([]byte, 1)
	for {
		// Read a char.
		_, err := p.r.Read(ob)
		if err != nil {
			// Return the current argument and raw length.
			return p.pad.String(), raw
		}
		raw++
		if ob[0] == '"' {
			if first {
				// Handle the start of a quote.
				quote = true
			} else if quote {
				// If this is within the quote, return the arg.
				return p.pad.String(), raw
			}
		} else if ob[0] == ' ' {
			// If this is the beginning, continue. If this isn't a quote, return. If it is, add to it.
			if first {
				continue
			} else if quote {
				_ = p.pad.AddByte(' ')
			} else {
				return p.pad.String(), raw
			}
		} else {
			// Just add to the argument.
			_ = p.pad.AddByte(ob[0])
		}

		// Set first to false.
		first = false
	}
}

// GetNextArg is used to get the next argument. If there are no additional arguments, the pointer will be nil.
func (p *Parser) GetNextArg() *Argument {
	if !p.initProperly {
		return nil
	}
	s, rawLen := p.getRawArgInfo()
	if rawLen == 0 {
		return nil
	}
	p.argId++
	return &Argument{Text: s, rawLen: rawLen, argId: p.argId, p: p}
}
