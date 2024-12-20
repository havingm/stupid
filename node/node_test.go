package node

import (
	leaf "github.com/stupid"
	"testing"
)

func TestNode(t *testing.T) {
	n := NewNode(":8000")
	leaf.Run(n)
}
