package fastparse

import (
	"errors"
	"io"
)

// Argument is used to define an argument.
type Argument struct {
	// The text from the argument.
	Text string

	// Internal data.
	rawLen int
	argId int
	p *Parser
}

// Rewind is used to rewind the reader to before an argument was read. This is useful for some argument verification situations.
// Note that if you want to rewind multiple arguments, you need to run it on every argument you parsed AFTER the one you wish to rewind in order of last to first.
func (a *Argument) Rewind() (err error) {
	if a.argId != a.p.argId {
		return errors.New("please rewind the arguments in order of last fetched first")
	}
	_, err = a.p.r.Seek(int64(a.rawLen*-1), io.SeekCurrent)
	return
}
