package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntMarkdownList_Markdown(t *testing.T) {
	t.Run("returns valid markdown", func(t *testing.T) {
		got := IntMarkdownList{1, 2, 3}.Markdown()
		want := "\n  - *1*\n  - *2*\n  - *3*\n"
		assert.Equal(t, want, got)
	})
}

func TestNewIntMarkdownList(t *testing.T) {
	got := NewIntMarkdownList([]int{1, 2, 3})
	want := IntMarkdownList{1, 2, 3}

	assert.Equal(t, want, got)
}

func TestIntMarkdownList_ToInt64(t *testing.T) {
	got := IntMarkdownList{1, 2, 3}.ToInt64()
	want := []int64{1, 2, 3}

	assert.Equal(t, want, got)
}
