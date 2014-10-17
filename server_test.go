package kernel

import (
	"testing"
	"code.google.com/p/go.net/context"
)

type handler struct {

}

func (h *handler) Serve(ctx context.Context, p *Server) {

}

func testBoot(t *testing.T) {
	h := &handler{}
	Boot(h)
}
