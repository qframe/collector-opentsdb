package models

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

var (
	t1 = time.Unix(0,0)
	t2 = time.Unix(10,0)
	a1 = Annotation{
		"tsuid1",
		"short discription",
		"Some lenghty description",
		map[string]string{"key1": "val1"},
		t1,
		t2,
	}
)

func TestNewAnnotation(t *testing.T) {
	got := NewAnnotation("tsuid1", "short discription", "Some lenghty description", t1)
	assert.Equal(t, a1.tsuid, got.tsuid)
	assert.Equal(t, a1.description, got.description)
	assert.Equal(t, a1.startTime, got.startTime)
	got.AddTag("key1", "val1")
	assert.Equal(t, "val1", got.custom["key1"])
}
