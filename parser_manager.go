package fastparse

import (
	"io"
	"sync"
)

// ParserManager is used to manage the parsing of arguments.
type ParserManager struct {
	initProperly bool
	pool *sync.Pool
}

// Used to allocate a fresh scratchpad on the fly.
func freshScratchpadAlloc(MessageLen int) func() interface{} {
	return func() interface{} {
		return &scratchpad{
			slice:    make([]byte, MessageLen),
			sliceLen: MessageLen,
		}
	}
}

// NewParserManager is used to create a new parser manager.
// MessageLen is the maximum length of each message, and PreAllocatedPads are the amount of scratchpads which will be pre-allocated.
// Each pre-allocated scratchpad uses the amount of bytes which a message would take, but they mean that the memory will not need to be allocated on the fly later.
func NewParserManager(MessageLen, PreAllocatedPads int) *ParserManager {
	p := &ParserManager{
		initProperly: true,
		pool:         &sync.Pool{New: freshScratchpadAlloc(MessageLen)},
	}
	for i := 0; i < PreAllocatedPads; i++ {
		p.pool.Put(&scratchpad{
			slice:    make([]byte, MessageLen),
			sliceLen: MessageLen,
			pool:     p.pool,
		})
	}
	return p
}

// Parser is used to get a parser from the manager.
func (m *ParserManager) Parser(r io.ReadSeeker) *Parser {
	if !m.initProperly {
		panic("The ParserManager needs to be created with NewParserManager")
	}
	pad := m.pool.Get().(*scratchpad)
	pad.reset()
	return &Parser{initProperly: true, pad: pad, r: r}
}
