// Copyright 2014 The erzha Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package kernel

import (
	"context"
	"testing"
)

type handler struct {
}

func (h *handler) Serve(ctx context.Context, p *Server) {

}

func testBoot(t *testing.T) {
	h := &handler{}
	Boot(h)
}
