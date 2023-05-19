package utils

import (
	"fmt"
	"testing"
)

func TestCheckM3u8MediaType(t *testing.T) {
	//ctx := app.RequestContext{
	//	Request:    protocol.Request{},
	//	Response:   protocol.Response{},
	//	Errors:     nil,
	//	Params:     nil,
	//	HTMLRender: nil,
	//	Keys:       nil,
	//}
	//hlog.CtxDebugf(context.Background(), "ctx:%v", ctx.Keys)
	//
	//err := rpc.NewBizStatusError(rpc.ErrorTypeDBHandle, errors.New("测试"))
	//fmt.Println("err:", err)

	url := "http://120.196.232.57:8088/wh7f454c46tw1038680586_-1109007784/" +
		"rrs03.hw.gmcc.net/PLTV/81/224/3221226737/index.m3u8?rrsip=rrs04.hw.g" +
		"mcc.net%3A8088%2Crrsip%3Drrs01.hw.gmcc.net%3A8088%2Crrsip%3Drrs02.hw.g" +
		"mcc.net%3A8088&zoneoffset=0&servicetype=1&icpid=81&limitflux=-1&limitd" +
		"ur=-1&tenantId=8601&accountinfo=~~V2.0~98CenL6PgXDZiEwY6CXl4w2a275dab923" +
		"3e528221dc87525216373~tlhyd06rjoii89cnYkMsHz0wskWEU03WBdHNKEoP0d8PrUG" +
		"kiGkjctOaZC1FcYA97e60da0037beb67f0a58a00047b53891~ExtInfoKMsv8R36b80IN" +
		"pqgoxAdpQ%3D%3D7831d729597ec84626dadb2edabd54bc%3A20230519095722%2C105" +
		"00064262245%2C223.74.107.134%2C20230519095722%2Curn%3ACims%3AliveTV%3AXTV" +
		"12306259%2C10001073820194%2C-1%2C0%2C1%2C%2C%2C2%2C%2C%2C%2C2%2C%2C37414" +
		"1397%2CEND&GuardEncType=2&RTS=1684461442&from=38&hms_devid=2250&online=" +
		"1684461442&vqe=3&it=H4sIAAAAAAAAAzXMQQvCIBiA4X-zo6hTmwdPRdBFglXX-NRvItgk" +
		"XYP-fSzq_j7vUsHj6WA0DqDDoKR3QXBHnfJTEEI4J7mXO9E1fNpi-s5DzmmOtoSN3cb9nXFKmF" +
		"aE95wwobrLtjxmiIZ-a_t6OKyG_-mIdU0eTWgTWaERiLFihCWVmZwzvK81_5IPQi-0Wp0AAAA"
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
