package utils

import (
	"fmt"
	"testing"
)

func TestCheckM3u8MediaType(t *testing.T) {
	//url := "http://rrs03.hw.gmcc.net:8088/81/16/20230501/277774463/index.m3u8?" +
	//	"rrsip=rrs04.hw.gmcc.net%3A8088%2Crrsip%3Drrs01.hw.gmcc.net%3A8088%2Crrsip%" +
	//	"3Drrs02.hw.gmcc.net%3A8088&zoneoffset=0&servicetype=0&icpid=&limitflux=-1&l" +
	//	"imitdur=-1&tenantId=8601&accountinfo=dEiFMTXlfbO08hgQfV9MjV6w2atVr%2BswLl04" +
	//	"e4ZeAgC7NUIbnwRdErUP9KJyDIkbwiSfYdMnvv0bfRinK039g%2FMzc4RHlNj%2F%2BvTPUfzLeJn" +
	//	"Nxz%2BjWcG7HsK6BGjg7yyMb865f9bedc6eaca798ff38411a11d24a%3A20230505093350%2C10" +
	//	"500064262245%2C120.231.214.2%2C20230505093350%2Cad748ef296f08e7c1f5889ef817138" +
	//	"7d%2C10001073820194%2C-1%2C1%2C0%2C-1%2C2%2C1%2C100000203%2C%2C390856565%2C1%2" +
	//	"C%2C390856645%2CEND&GuardEncType=2&it=H4sIAAAAAAAAAE2OwQ6CMBBE_6bHBopWOPSkMTE" +
	//	"xaCJ6NUO7VGKh2qKJfy8YDh72MJv3JjMEaNptVJ5JuUSdyXxRLCBMUS_zJml0mpoCckUs0rP0SjAN" +
	//	"59relt5M2uW0vqYi4WkhucjEeAmrpsqtg1XJjy5fXU1hDqN6ovBuNSkTG_5G5LA2kMXQ-p4fHT7" +
	//	"n4GaEUTXP61_OsWEKFeJ9fLAb4tp3DwQye29_nGrgIrEH9B2WSnT05x2CGUd8ATBTIEvvAAAA"

	url := "http://120.196.232.58:8088/wh7f454c46tw2764159674_813681295/rrs03.hw.gmcc.net/81/16/20230410/277670393/index.m3u8?rrsip=rrs04.hw.gmcc.net%3A8088%2Crrsip%3Drrs01.hw.gmcc.net%3A8088%2Crrsip%3Drrs02.hw.gmcc.net%3A8088&zoneoffset=0&servicetype=0&icpid=81&limitflux=-1&limitdur=-1&tenantId=8601&accountinfo=X22A2Pbdf%2BOXUTfbzan0hqDKDPSREU3wBme31JY%2BF2a6KqHJl86azYp4vyHpV9FcP07NQAVzqu92gKFP9Fjp%2B9wxvgkLFAqAZjI%2F7fl%2FksJ4E2JTOERcjrapLcCpOe2J63e99c12a6dd7571d87cd40d355d0c02%3A20230505111121%2C10500064262245%2C120.231.214.2%2C20230505111121%2C9d493c0a484e67cb7f279bacb7c1b9fa%2C10001073820194%2C-1%2C1%2C0%2C-1%2C2%2C1%2C100000203%2C%2C390540056%2C1%2C%2C390540134%2CEND&GuardEncType=2&RTS=1683256282&from=38&hms_devid=2250&vqe=3&it=H4sIAAAAAAAAAy2MvQrDIBhF38ZR1Kgxg1NLoEsopO1aPn8RbKQmDfTtm5QMdzhwzl0qWH85a9EGxoVRXCrDJVXKhK41xIWmY0JQg2b_HopukIWc0xSH4vbsMZ6elBFMO4lZw7YRdNsv-wxRk789fF7G1wO2dPR1TdZrNwe8wowhxuojLKlM-Jrhe6_5UH7Dt9zznQAAAA"

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
	//url, err := fasthttp.GetRedirectUrlWithRetryTimes(url, 3)
	//if err != nil {
	//	fmt.Println("err:", err)
	//} else {
	//	fmt.Println("url:", url)
	//}
}
