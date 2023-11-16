package utils

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/baoyxing/go-tools/net/fasthttp"
	"github.com/bytedance/gopkg/lang/mcache"
	"github.com/bytedance/sonic"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"sync"
	"sync/atomic"
	"testing"
)

type M3u8VerifyData struct {
	Filesize   int       `json:"filesize"`
	Tssize     []int     `json:"tssize"`
	Tsduration []float64 `json:"tsduration,"`
	Hash       string    `json:"hash"`
}
type VodM3u8Data struct {
	start int
	end   int
	url   string
}

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

	// 原始数据
	data := []byte("Hello, World!")
	length := len(data)
	fmt.Println("length:", length)
	// 截取从索引 7 到索引 12 的部分数据（索引从 0 开始）
	slice := data[7:length]

	// 输出截取后的数据
	fmt.Println(string(slice)) // 输出 "World!"
}

func TestParseM3u8Ts(t *testing.T) {
	//url := "http://4dcloud.8866.org:7090//home/cdn/t_disk1/416_14714489.m3u8"
	url := "http://4dcloud.8866.org:7090/home/cdn/t_disk1/416_avc14714489.m3u8"
	tsUrls, _, tsDuration, err := ParseM3u8Ts(context.Background(),
		url)

	if err != nil {
		fmt.Println("err:%", err)
	}
	content, err := DownloadVodTs(tsUrls, tsDuration)
	if err != nil {
		fmt.Println("err:%", err)
	} else {
		err := ioutil.WriteFile("test-sub.p2p", content, 0644)
		if err != nil {
			fmt.Println("写入文件时发生错误:", err)
			return
		}
	}

	//_, _, pieceSize, _, err := getFullVodM3u8DataSize(tsUrls)
	//fmt.Println("pieceSize:", pieceSize)
	//for _, value := range arr {
	//	length := 0
	//	for _, v1 := range value {
	//		length += v1.end - v1.start + 1
	//	}
	//	if length != pieceSize {
	//		fmt.Println("length:", length)
	//		for _, v1 := range value {
	//			fmt.Printf("0--url:%s,start:%d,end:%d\n", v1.url, v1.start, v1.end)
	//		}
	//		fmt.Println("++++++++++")
	//	}
	//
	//}

}

func DownloadVodAllTs(urls []string, tsDuration []float64) ([]byte, error) {
	total := len(urls)
	size, pieceTotal, pieceSize, ts, err := getFullVodM3u8DataSize(urls)
	if err != nil {
		return nil, err
	}
	fmt.Println("size:", size)
	verifyData := mcache.Malloc(0, int(pieceTotal)*40)
	poolSize := 30
	var vodM3u8DataMap sync.Map
	var errMap sync.Map
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(poolSize, func(i interface{}) {
		index := i.(int)
		url := urls[index]
		tsLength := ts[index]
		body := mcache.Malloc(0, tsLength)
		body, err := fasthttp.DownloadFullTsDataWithRetryTimes(url, "h", true, 2, 40)
		if len(body) != tsLength {
			fmt.Println("m3u8 长度不一致s")
		}
		if err != nil {
			errMap.Store(index, err)
		} else {
			if index == 1 {
				h := md5.New()
				io.WriteString(h, string(body))
				tsMD5 := fmt.Sprintf("%x", h.Sum(nil))
				fmt.Printf("1--------index:%d,tsMD5:%s\n", index, tsMD5)
			}
			vodM3u8DataMap.Store(index, body)

		}
		mcache.Free(body)
		wg.Done()
	})
	defer p.Release()
	for j := 0; j < total; j++ {
		wg.Add(1)
		p.Invoke(j)
	}
	wg.Wait()
	errCount := 0
	errMap.Range(func(key, value any) bool {
		index := key.(int)
		errMsg := value.(string)
		errCount++
		url := urls[index]
		fmt.Printf("获取TS内容失败,url:%s,err:%s\n", url, errMsg)
		return true
	})
	if errCount > 0 {
		return nil, errors.New("获取TS内容失败")
	}
	bodyList := make(map[int][]byte)
	vodM3u8DataMap.Range(func(key, value any) bool {
		index := key.(int)
		body := value.([]byte)
		if index == 1 {
			h := md5.New()
			io.WriteString(h, string(body))
			tsMD5 := fmt.Sprintf("%x", h.Sum(nil))
			fmt.Printf("2--------index:%d,tsMD5:%s\n", index, tsMD5)
		}
		if bodyList[index] == nil {
			bodyList[index] = make([]byte, 0)
		} else {
			fmt.Println("索引错位------------")
		}
		bodyList[index] = body
		return true
	})
	tmpBody := mcache.Malloc(0, int(size))
	for j := 0; j < total; j++ {
		body := bodyList[j]

		tmpBody = append(tmpBody, body...)
	}

	count := int(size) / pieceSize
	for j := 0; j < count; j++ {
		start := j * pieceSize
		end := (j + 1) * pieceSize
		body := tmpBody[start:end]
		//if start == 0 {
		//	h := md5.New()
		//	io.WriteString(h, string(body))
		//	tsMD5 := fmt.Sprintf("%x", h.Sum(nil))
		//	fmt.Println("2----tsMD5:", tsMD5)
		//}
		h := sha1.New()
		h.Write(body)
		hexStr := hex.EncodeToString(h.Sum(nil))
		verifyData = append(verifyData, []byte(hexStr)...)
	}
	if int(size)%pieceSize > 0 {
		start := count * pieceSize
		end := size
		body := tmpBody[start:end]
		h := sha1.New()
		h.Write(body)
		hexStr := hex.EncodeToString(h.Sum(nil))
		verifyData = append(verifyData, []byte(hexStr)...)
	}
	mcache.Free(tmpBody)
	m3u8VerifyData := &M3u8VerifyData{
		Filesize:   int(size),
		Tssize:     ts,
		Tsduration: tsDuration,
		Hash:       string(verifyData),
	}
	verifyData, err = sonic.Marshal(m3u8VerifyData)
	mcache.Free(verifyData)
	return verifyData, nil
}

func DownloadVodTs(urls []string, tsDuration []float64) ([]byte, error) {
	total := len(urls)
	size, pieceTotal, pieceSize, ts, err := getFullVodM3u8DataSize(urls)
	if err != nil {
		return nil, err
	}
	poolSize := 30
	poolTotal := total / poolSize
	verifyData := mcache.Malloc(0, int(pieceTotal)*40)
	surplusBody := make([]byte, 0)
	for i := 0; i < poolTotal; i++ {
		var vodM3u8DataMap sync.Map
		var errMap sync.Map
		var wg sync.WaitGroup
		p, _ := ants.NewPoolWithFunc(poolSize, func(i interface{}) {
			index := i.(int)
			url := urls[index]
			tsLength := ts[index]
			body := mcache.Malloc(0, tsLength)
			body, err := fasthttp.DownloadFullTsDataWithRetryTimes(url, "h", true, 2, 40)
			if err != nil {
				errMap.Store(index, err)
			} else {
				vodM3u8DataMap.Store(index, body)

			}
			mcache.Free(body)
			wg.Done()
		})
		defer p.Release()
		for j := 0; j < poolSize; j++ {
			wg.Add(1)
			p.Invoke(j + i*poolSize)
		}
		wg.Wait()
		errCount := 0
		errMap.Range(func(key, value any) bool {
			index := key.(int)
			errMsg := value.(string)
			errCount++
			url := urls[index]
			fmt.Printf("获取TS内容失败,url:%s,err:%s\n", url, errMsg)
			return true
		})
		if errCount > 0 {
			return nil, errors.New("获取TS内容失败")
		}
		bodyList := make(map[int][]byte)
		tsLength := len(surplusBody)
		vodM3u8DataMap.Range(func(key, value any) bool {
			index := key.(int)
			body := value.([]byte)
			tsLength += ts[index]
			if bodyList[index] == nil {
				bodyList[index] = make([]byte, 0)
			}

			bodyList[index] = body
			return true
		})
		tmpBody := mcache.Malloc(0, tsLength)
		if len(surplusBody) > 0 {
			tmpBody = append(tmpBody, surplusBody...)
		}
		surplusBody = make([]byte, 0)
		for j := 0; j < poolSize; j++ {
			index := j + i*poolSize
			body := bodyList[index]
			tmpBody = append(tmpBody, body...)
		}
		count := tsLength / pieceSize
		for j := 0; j < count; j++ {
			start := j * pieceSize
			end := (j + 1) * pieceSize
			body := tmpBody[start:end]
			h := sha1.New()
			h.Write(body)
			hexStr := hex.EncodeToString(h.Sum(nil))
			verifyData = append(verifyData, []byte(hexStr)...)
		}
		surplus := tsLength % pieceSize
		if surplus > 0 {
			start := count * pieceSize
			end := tsLength
			body := tmpBody[start:end]
			surplusBody = append(surplusBody, body...)
		}
		mcache.Free(tmpBody)
	}
	if len(surplusBody) > 0 {
		fmt.Println("------------------有剩余")
	}
	surplus := total % poolSize
	if surplus > 0 {
		fmt.Printf("--------total:%d,poolSize:%d,surplus:%d\n", total, poolSize, surplus)
		var vodM3u8DataMap sync.Map
		var errMap sync.Map
		var wg sync.WaitGroup
		p, _ := ants.NewPoolWithFunc(poolSize, func(i interface{}) {
			index := i.(int)
			url := urls[index]
			tsLength := ts[index]
			body := mcache.Malloc(0, tsLength)
			body, err := fasthttp.DownloadFullTsDataWithRetryTimes(url, "h", true, 2, 40)
			if err != nil {
				errMap.Store(index, err)
			} else {
				vodM3u8DataMap.Store(index, body)

			}
			mcache.Free(body)
			wg.Done()
		})
		defer p.Release()
		for j := poolTotal * poolSize; j < total; j++ {
			wg.Add(1)
			p.Invoke(j)
		}
		wg.Wait()
		errCount := 0
		errMap.Range(func(key, value any) bool {
			index := key.(int)
			errMsg := value.(string)
			errCount++
			url := urls[index]
			fmt.Printf("获取TS内容失败,url:%s,err:%s\n", url, errMsg)
			return true
		})
		if errCount > 0 {
			return nil, errors.New("获取TS内容失败")
		}
		bodyList := make(map[int][]byte)
		tsLength := len(surplusBody)
		vodM3u8DataMap.Range(func(key, value any) bool {
			index := key.(int)
			body := value.([]byte)
			tsLength += ts[index]
			if bodyList[index] == nil {
				bodyList[index] = make([]byte, 0)
			}

			bodyList[index] = body
			return true
		})
		tmpBody := mcache.Malloc(0, tsLength)
		if len(surplusBody) > 0 {
			tmpBody = append(tmpBody, surplusBody...)
		}
		surplusBody = make([]byte, 0)
		//fmt.Println("poolTotal * pieceSize:", poolTotal*pieceSize)
		//fmt.Println("total:", total)
		for j := poolTotal * poolSize; j < total; j++ {
			index := j
			body := bodyList[index]
			h := md5.New()
			io.WriteString(h, string(body))
			tsMD5 := fmt.Sprintf("%x", h.Sum(nil))
			fmt.Printf("--------url:%s,index:%d,tsMD5:%s\n", urls[index], index, tsMD5)
			tmpBody = append(tmpBody, body...)
		}
		count := tsLength / pieceSize
		for j := 0; j < count; j++ {
			start := j * pieceSize
			end := (j + 1) * pieceSize
			body := tmpBody[start:end]
			h := sha1.New()
			h.Write(body)
			hexStr := hex.EncodeToString(h.Sum(nil))
			verifyData = append(verifyData, []byte(hexStr)...)
		}
		surplus := tsLength % pieceSize
		if surplus > 0 {
			start := count * pieceSize
			end := tsLength
			body := tmpBody[start:end]
			surplusBody = append(surplusBody, body...)
		}
		mcache.Free(tmpBody)

	}
	if len(surplusBody) > 0 {
		h := sha1.New()
		h.Write(surplusBody)
		hexStr := hex.EncodeToString(h.Sum(nil))
		verifyData = append(verifyData, []byte(hexStr)...)
	}
	m3u8VerifyData := &M3u8VerifyData{
		Filesize:   int(size),
		Tssize:     ts,
		Tsduration: tsDuration,
		Hash:       string(verifyData),
	}
	verifyData, err = sonic.Marshal(m3u8VerifyData)
	mcache.Free(verifyData)
	return verifyData, nil

}

func getFullVodM3u8DataSize(tsUrls []string) (uint64, int, int, []int, error) {
	poolSize := 50
	count := len(tsUrls)
	if count < poolSize {
		poolSize = count
	}
	var tmpTsLengthMap sync.Map
	var errMap sync.Map
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(poolSize, func(i interface{}) {
		index := i.(int)
		url := tsUrls[index]
		size, err := fasthttp.GetTsSizeWithRetryTimes(url, 2)
		if err != nil {
			errMap.LoadOrStore(index, err.Error())
		} else {
			tmpTsLengthMap.LoadOrStore(index, size)
		}
		wg.Done()
	})
	defer p.Release()
	for index, _ := range tsUrls {
		wg.Add(1)
		p.Invoke(index)
	}
	wg.Wait()
	errCount := 0
	errMap.Range(func(key, value any) bool {
		index := key.(int)
		errMsg := value.(string)
		errCount++
		url := tsUrls[index]
		fmt.Printf("获取TS长度失败,url:%s,err:%s\n", url, errMsg)
		return true
	})
	if errCount > 0 {
		return 0, 0, 0, nil, errors.New("获取ts长度失败")
	}
	totalSize := uint64(0)
	tsLengthMap := make(map[int]uint64, 0)
	tmpTsLengthMap.Range(func(key, value any) bool {
		index := key.(int)
		size := value.(uint64)
		totalSize += size
		tsLengthMap[index] = size
		return true
	})
	tsLength := make([]int, 0, count)
	for index, _ := range tsUrls {
		size := tsLengthMap[index]
		if size == 0 {
			return 0, 0, 0, nil, errors.New("获取ts长度失败")
		}
		tsLength = append(tsLength, int(size))
	}
	pieceSize := CalcPieceSize(totalSize)
	pieceTotal := int(totalSize) / pieceSize
	if int(totalSize)%pieceSize > 0 {
		pieceTotal += 1
	}

	return totalSize, pieceTotal, pieceSize, tsLength, nil
}

func DownloadSubVodTs(urls []string, tsDuration []float64) ([]byte, error) {
	size, pieceTotal, pieceSize, ts, arr, err := getVodM3u8DataSize(urls)
	if err != nil {
		return nil, err
	}
	//for _, value := range arr {
	//	length := 0
	//	for _, v1 := range value {
	//		length += v1.end - v1.start + 1
	//	}
	//	if length != pieceSize {
	//		repo.log.CtxInfof(ctx, "length:%d", length)
	//		for _, v1 := range value {
	//			repo.log.CtxInfof(ctx, "0--url:%s,start:%d,end:%d", v1.url, v1.start, v1.end)
	//		}
	//	}
	//
	//}
	poolSize := 30
	poolTotal := pieceTotal / poolSize
	verifyData := mcache.Malloc(0, int(pieceTotal)*40)
	var downloadedCount uint32
	_, cancel := context.WithCancel(context.TODO())
	for i := 0; i < poolTotal; i++ {
		var vodM3u8DataMap sync.Map
		var wg sync.WaitGroup
		p, _ := ants.NewPoolWithFunc(poolSize, func(i interface{}) {
			index := i.(int)
			vodM3u8DataList := arr[index]
			body := make([]byte, 0)
			for _, value := range vodM3u8DataList {
				tempBody, _ := fasthttp.DownloadTsDataWithRetryTimes(value.url, value.start, value.end, 3, 40)
				body = append(body, tempBody...)
			}
			length := len(body)
			if length != pieceSize {

			}
			count := atomic.LoadUint32(&downloadedCount)
			count++
			atomic.StoreUint32(&downloadedCount, count)
			vodM3u8DataMap.LoadOrStore(index, body)

			mcache.Free(body)
			wg.Done()
		})
		defer p.Release()
		for j := 0; j < poolSize; j++ {
			wg.Add(1)
			p.Invoke(j + i*poolSize)
		}
		wg.Wait()
		bodyList := make(map[int][]byte, 0)
		vodM3u8DataMap.Range(func(key, value any) bool {
			index := key.(int)
			body := value.([]byte)
			if bodyList[index] == nil {
				bodyList[index] = make([]byte, 0)
			}
			bodyList[index] = body
			return true
		})
		for j := 0; j < poolSize; j++ {
			index := j + i*poolSize
			body := bodyList[index]
			h := sha1.New()
			h.Write(body)
			hexStr := hex.EncodeToString(h.Sum(nil))
			verifyData = append(verifyData, []byte(hexStr)...)
		}

	}
	index := poolTotal * poolSize
	surplus := pieceTotal % poolSize
	if surplus > 0 {
		var vodM3u8DataMap sync.Map
		var wg sync.WaitGroup
		p, _ := ants.NewPoolWithFunc(poolSize, func(i interface{}) {
			index := i.(int)
			vodM3u8DataList := arr[index]
			body := mcache.Malloc(0, pieceSize)
			for _, value := range vodM3u8DataList {
				tempBody, _ := fasthttp.DownloadTsDataWithRetryTimes(value.url, value.start, value.end, 3, 40)
				body = append(body, tempBody...)
			}
			count := atomic.LoadUint32(&downloadedCount)
			count++
			atomic.StoreUint32(&downloadedCount, count)
			vodM3u8DataMap.LoadOrStore(index, body)
			mcache.Free(body)
			wg.Done()
		})
		defer p.Release()
		for j := 0; j < surplus; j++ {
			wg.Add(1)
			p.Invoke(j + index)
		}
		wg.Wait()
		bodyList := make(map[int][]byte, 0)
		vodM3u8DataMap.Range(func(key, value any) bool {
			index := key.(int)
			body := value.([]byte)
			if bodyList[index] == nil {
				bodyList[index] = make([]byte, 0)
			}
			bodyList[index] = body
			return true
		})
		for j := 0; j < surplus; j++ {
			index := j + index
			body := bodyList[index]
			h := sha1.New()
			h.Write(body)
			hexStr := hex.EncodeToString(h.Sum(nil))
			verifyData = append(verifyData, []byte(hexStr)...)
		}
	}
	m3u8VerifyData := &M3u8VerifyData{
		Filesize:   int(size),
		Tssize:     ts,
		Tsduration: tsDuration,
		Hash:       string(verifyData),
	}
	cancel()
	verifyData, err = sonic.Marshal(m3u8VerifyData)
	mcache.Free(verifyData)
	return verifyData, nil
}

func getVodM3u8DataSize(tsUrls []string) (uint64, int, int, []int, [][]VodM3u8Data, error) {
	poolSize := 50
	count := len(tsUrls)
	if count < poolSize {
		poolSize = count
	}
	var tmpTsLengthMap sync.Map
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(poolSize, func(i interface{}) {
		index := i.(int)
		url := tsUrls[index]
		size, _ := fasthttp.GetTsSizeWithRetryTimes(url, 3)
		tmpTsLengthMap.LoadOrStore(index, size)
		wg.Done()
	})
	defer p.Release()
	for index, _ := range tsUrls {
		wg.Add(1)
		p.Invoke(index)
	}
	wg.Wait()
	totalSize := uint64(0)
	tsLengthMap := make(map[int]uint64, 0)
	tmpTsLengthMap.Range(func(key, value any) bool {
		index := key.(int)
		size := value.(uint64)
		totalSize += size
		tsLengthMap[index] = size
		return true
	})
	tsLength := make([]int, 0, count)
	for index, _ := range tsUrls {
		size := tsLengthMap[index]
		if size == 0 {
			fmt.Println("获取ts失败")
			return 0, 0, 0, nil, nil, errors.New("获取ts失败")
		}
		tsLength = append(tsLength, int(size))
	}
	pieceSize := CalcPieceSize(totalSize)
	fmt.Println("pieceSize:", pieceSize)
	pieceTotal := int(totalSize) / pieceSize
	if int(totalSize)%pieceSize > 0 {
		pieceTotal += 1
	}
	pieceData := make([][]VodM3u8Data, pieceTotal)
	index := uint64(0)
	surplus := int(0)
	for key, value := range tsLength {
		if surplus == 0 {
			if value < pieceSize {
				pieceData[index] = append(pieceData[index], VodM3u8Data{
					start: 0,
					end:   value - 1,
					url:   tsUrls[key],
				})
				surplus = pieceSize - value
				continue
			}
			count := value / pieceSize
			surplus = value % pieceSize
			for i := 0; i < count; i++ {
				start := i * pieceSize
				end := (i+1)*pieceSize - 1
				pieceData[index] = append(pieceData[index], VodM3u8Data{
					start: start,
					end:   end,
					url:   tsUrls[key],
				})
				index++
			}
			if surplus > 0 {
				start := count * pieceSize
				end := value - 1
				pieceData[index] = append(pieceData[index], VodM3u8Data{
					start: start,
					end:   end,
					url:   tsUrls[key],
				})
				surplus = pieceSize - surplus
			}

		} else if surplus > 0 {
			if value < surplus {
				pieceData[index] = append(pieceData[index], VodM3u8Data{
					start: 0,
					end:   value - 1,
					url:   tsUrls[key],
				})
				surplus = surplus - value
				continue
			}
			pieceData[index] = append(pieceData[index], VodM3u8Data{
				start: 0,
				end:   surplus - 1,
				url:   tsUrls[key],
			})
			index++
			count := (value - surplus) / pieceSize
			for i := 0; i < count; i++ {
				start := i*pieceSize + surplus
				end := (i+1)*pieceSize + surplus - 1
				pieceData[index] = append(pieceData[index], VodM3u8Data{
					start: start,
					end:   end,
					url:   tsUrls[key],
				})
				index++
			}
			tempSurplus := surplus
			surplus = (value - surplus) % pieceSize
			if surplus > 0 {
				start := count*pieceSize + tempSurplus
				end := value - 1
				pieceData[index] = append(pieceData[index], VodM3u8Data{
					start: start,
					end:   end,
					url:   tsUrls[key],
				})
				surplus = pieceSize - surplus
			}
		}
	}
	return totalSize, pieceTotal, pieceSize, tsLength, pieceData, nil
}
