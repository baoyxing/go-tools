package utils

import (
	"context"
	"fmt"
	"testing"
)

func TestCheckM3u8MediaType(t *testing.T) {
	url := "http://4dcloud.8866.org:7090//home/cdn/t_disk1/416_14714489.m3u8"
	mediaType, err := CheckM3u8MediaType(url)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		switch mediaType {
		case Vod:
			fmt.Println("点播")
		case Live:
			fmt.Println("直播")
		}
	}
}

func TestParseM3u8Ts(t *testing.T) {
	tsUrls, duration, _, err := ParseM3u8Ts(context.Background(),
		"http://4dcloud.8866.org:7090/home/cdn/t_disk1/416_avc14714489.m3u8")
	if err != nil {
		fmt.Println("err:", err.Error())
	}
	for _, value := range tsUrls {
		fmt.Println("tsUrl:", value)
	}
	fmt.Println("duration:", duration)
}
