package agiledoc

import (
	bf "github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	"testing"
)

func parse(in []byte, t *Tokenizer) []byte {
	return bf.Markdown(in, t, 0)
}
func callback(t *Tokenizer, c chan Token, cb func(...interface{})) chan struct{} {
	var quit chan struct{} = make(chan struct{})
	go func() {
		for quit == nil {
			select {
			case <-c:
				tok := <-c
				cb(tok.rawTxt)
			}
		}
	}()
	return quit
}
func read(path string) []byte {
	f, _ := os.Open(path)
	b, _ := ioutil.ReadAll(f)
	return b
}
func test(p string, cb func(...interface{})) {
	t, i := NewTokenizer()
	q := callback(t, i, cb)
	_ = parse(read(p), t)
	q <- struct{}{}
}
func TestTokens(t *testing.T) {
	test("home/j/src/go/src/github.com/russross/blackfriday/testdata/Tabs.text", (*t).Log)
}
