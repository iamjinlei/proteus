package markdown

import (
	"bytes"

	"github.com/gomarkdown/markdown/ast"
)

type htmlTag struct {
	buf        bytes.Buffer
	closingTag []byte
	closeFn    func(*htmlTag) ast.WalkStatus
}

func (b *htmlTag) close() ast.WalkStatus {
	return b.closeFn(b)
}

type htmlTagStack struct {
	stack []*htmlTag
}

func newHtmlTagStack() *htmlTagStack {
	return &htmlTagStack{}
}

func (s *htmlTagStack) len() int {
	return len(s.stack)
}

func (s *htmlTagStack) empty() bool {
	return s.len() == 0
}

func (s *htmlTagStack) top() *htmlTag {
	return s.stack[s.len()-1]
}

func (s *htmlTagStack) push(
	closingTag []byte,
	fn func(*htmlTag) ast.WalkStatus,
) {
	s.stack = append(s.stack, &htmlTag{
		closingTag: closingTag,
		closeFn:    fn,
	})
}

func (s *htmlTagStack) pop() {
	s.stack = s.stack[:len(s.stack)-1]
}
