package kod

import (
	"context"
	"testing"
)

func TestGetNoKodInContext(t *testing.T) {
	_, err := Get[Main](context.Background())
	if err == nil {
		t.Fatal("unexpected Get success")
	}
}
