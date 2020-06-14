package fastparse

import (
	"bytes"
	"io"
	"sync"
	"testing"
)

func TestParser(t *testing.T) {
	m := NewParserManager(6, 1)
	r := bytes.NewReader([]byte("a \"test\""))
	p := m.Parser(r)
	defer p.Done()
	a := p.GetNextArg()
	if a == nil {
		t.Fatal("a shouldn't be nil")
		return
	}
	if a.Text != "a" {
		t.Fatal("should be a")
		return
	}
	a = p.GetNextArg()
	if a.Text != "test" {
		t.Fatal("should be test")
		return
	}
	err := a.Rewind()
	if err != nil {
		t.Fatal(err)
		return
	}
	s, err := p.Remainder()
	if err != nil {
		t.Fatal(err)
		return
	}
	if s != "\"test\"" {
		t.Fatal("should be \"test\"")
		return
	}
}

func BenchmarkParser_FullyPreAllocated(b *testing.B) {
	n := 10000
	m := NewParserManager(2000, n)
	readers := make([]io.ReadSeeker, n)
	for i := range readers {
		readers[i] = bytes.NewReader([]byte("\"test\""))
	}
	wg := sync.WaitGroup{}
	wg.Add(n)
	b.ResetTimer()
	for i := 0; i < n; i++ {
		r := readers[i]
		go func() {
			defer wg.Done()
			p := m.Parser(r)
			p.GetNextArg()
			p.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkParser_HalfPreAllocated(b *testing.B) {
	n := 10000
	m := NewParserManager(2000, n)
	readers := make([]io.ReadSeeker, n*2)
	for i := range readers {
		readers[i] = bytes.NewReader([]byte("\"test\""))
	}
	wg := sync.WaitGroup{}
	wg.Add(n*2)
	b.ResetTimer()
	for i := 0; i < n*2; i++ {
		r := readers[i]
		go func() {
			defer wg.Done()
			p := m.Parser(r)
			p.GetNextArg()
			p.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkParser_NotPreAllocated(b *testing.B) {
	n := 10000
	m := NewParserManager(2000, 0)
	readers := make([]io.ReadSeeker, n)
	for i := range readers {
		readers[i] = bytes.NewReader([]byte("\"test\""))
	}
	wg := sync.WaitGroup{}
	wg.Add(n)
	b.ResetTimer()
	for i := 0; i < n; i++ {
		r := readers[i]
		go func() {
			defer wg.Done()
			p := m.Parser(r)
			p.GetNextArg()
			p.Done()
		}()
	}
	wg.Wait()
}
