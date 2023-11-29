package utils

import (
	"fmt"
	"testing"
)

func TestHash1(t *testing.T) {
	content := Hash1([]byte("http://4dcloud.8866.org:7090/home/cdn/t_disk1/416_14714489.m3u8"))
	fmt.Println(content)
}
