package day

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	sec := time.Now().Unix()
	tz := "America/Cordoba"
	actual := Format(sec, "YYYY-MM-DD HH:mm:ss", tz)
	fmt.Println("actual:", actual)
	l, _ := time.LoadLocation(tz)
	expected := time.Unix(sec, 0).In(l).Format("2006-01-02 15:04:05")
	assert.Equal(t, expected, actual)
}
