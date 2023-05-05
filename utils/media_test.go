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

	url := "http://120.198.234.199:8088/wh7f454c46tw2677549006_-863282200/rrs03.hw.gmcc.net/81/16/" +
		"20230501/277774463/index.m3u8?rrsip=rrs04.hw.gmcc.net%3A8088%2Crrsip%3Drrs01.hw.gmcc.net%3" +
		"A8088%2Crrsip%3Drrs02.hw.gmcc.net%3A8088&zoneoffset=0&servicetype=0&icpid=81&limitfl" +
		"ux=-1&limitdur=-1&tenantId=8601&accountinfo=dEiFMTXlfbO08hgQfV9MjV6w2atVr%2BswLl04e4Z" +
		"eAgC7NUIbnwRdErUP9KJyDIkbwiSfYdMnvv0bfRinK039g%2FMzc4RHlNj%2F%2BvTPUfzLeJnNxz%2BjWcG" +
		"7HsK6BGjg7yyMb865f9bedc6eaca798ff38411a11d24a%3A20230505093350%2C10500064262245%2C120" +
		".231.214.2%2C20230505093350%2Cad748ef296f08e7c1f5889ef8171387d%2C10001073820194%2C-1%2C1%2C0" +
		"%2C-1%2C2%2C1%2C100000203%2C%2C390856565%2C1%2C%2C390856645%2CEND&GuardEncType=2&RTS=16832" +
		"51900&from=202&hms_devid=2552&vqe=3&it=H4sIAAAAAAAAAy2MvQrCMBRG36ZjyE8bkyGTIrgEoeoqN8ltCMQW" +
		"k1rw7W2lwzccOOebC3i8nIwSUnbghFStboEH7To10MEzFjTIAzYV33YyovGQcxqjncKWPfrjk3FKmJaEC76ON" +
		"rft8pwhGvq37eflsOywpj2WJXk0oQ5kgUogxoIR5jSN5Jrhey95V370cv4QnQAAAA"

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
