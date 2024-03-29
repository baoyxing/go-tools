package fasthttp

import (
	"crypto/md5"
	"fmt"
	"github.com/baoyxing/go-tools/net/consts"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func Get(httpUrl string) ([]byte, error) {
	rsp, err := RetryDoHTTPWithErr(func() (*fasthttp.Response, error) {
		return get(httpUrl, 3*time.Second)
	}, 3, 0)
	if err != nil {
		return nil, err
	}
	defer fasthttp.ReleaseResponse(rsp)
	return rsp.Body(), nil
}

func DownloadTsDataWithRetryTimes(httpUrl string, start, end, retryTimes int, timeOut time.Duration) ([]byte, error) {
	return RetryDoHTTPWithEmptyBody(func() ([]byte, error) {
		return downloadTsData(httpUrl, start, end, timeOut)
	}, retryTimes, 0)
}

func DownloadFullTsDataWithRetryTimes(httpUrl, key string, isCheckMD5 bool, retryTimes int, timeOut time.Duration) ([]byte, error) {
	return RetryDoHTTPWithEmptyBody(func() ([]byte, error) {
		return downFullTsData(httpUrl, key, isCheckMD5, timeOut)
	}, retryTimes, 0)
}

// downFullTsData 下载全量TS
func downFullTsData(url, key string, isCheckMD5 bool, timeOut time.Duration) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(consts.MethodGet)
	req.SetRequestURI(url)
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; " +
		"Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
	rsp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(rsp)
	if err := fasthttp.DoTimeout(req, rsp, timeOut*time.Second); err != nil {
		return nil, err
	}
	body := rsp.Body()
	if isCheckMD5 {
		md5Str, err := getM3u8Md5WithUrl(url, key)
		if err != nil {
			return nil, err
		}
		isOk, err := checkM3u8TsMD5(body, md5Str)
		if err != nil {
			return nil, err
		}
		if !isOk {
			return nil, errors.New("下载ts数据不完整")
		}
	}
	return body, nil
}

// range 方式下载TS
func downloadTsData(url string, start, end int, timeOut time.Duration) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(consts.MethodGet)
	rangeValue := "bytes=" + strconv.Itoa(start) + "-" + strconv.Itoa(end)
	req.Header.Set("Range", rangeValue)
	req.SetRequestURI(url)
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; " +
		"Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
	rsp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(rsp)
	if err := fasthttp.DoTimeout(req, rsp, timeOut*time.Second); err != nil {
		return nil, err
	}
	return rsp.Body(), nil
}

func GetTsSizeWithRetryTimes(httpUrL string, retryTimes int) (uint64, error) {
	rsp, err := RetryDoHTTPWithErr(func() (*fasthttp.Response, error) {
		return GetTsSize(httpUrL)
	}, retryTimes, 0)
	if err != nil {
		return 0, err
	}
	defer fasthttp.ReleaseResponse(rsp)
	contentRange := string(rsp.Header.Peek("Content-Range"))
	size := int64(0)
	if contentRange != "" {
		arr := strings.Split(contentRange, "/")
		if len(arr) >= 2 {
			size, _ = strconv.ParseInt(arr[1], 10, 64)
		}
	}
	return uint64(size), nil
}

func GetTsSize(httpUrL string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(consts.MethodGet)
	req.Header.Set("Range", "bytes=0-0")
	req.SetRequestURI(httpUrL)
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; " +
		"Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
	rsp := fasthttp.AcquireResponse()
	if err := fasthttp.DoTimeout(req, rsp, 5*time.Second); err != nil {
		return nil, err
	}
	return rsp, nil
}

func GetFileSizeWithRetryTimes(httpUrL string, retryTimes int) (uint64, error) {
	rsp, err := RetryDoHTTPWithErr(func() (*fasthttp.Response, error) {
		return GetFileSize(httpUrL)
	}, retryTimes, 0)
	if err != nil {
		return 0, err
	}
	defer fasthttp.ReleaseResponse(rsp)
	contentRange := string(rsp.Header.Peek("Content-Range"))
	size := int64(0)
	if contentRange != "" {
		arr := strings.Split(contentRange, "/")
		if len(arr) >= 2 {
			size, _ = strconv.ParseInt(arr[1], 10, 64)
		}
	}
	return uint64(size), nil
}

func GetFileSize(httpUrL string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(consts.MethodHead)
	req.SetRequestURI(httpUrL)
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; " +
		"Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
	rsp := fasthttp.AcquireResponse()
	if err := fasthttp.DoTimeout(req, rsp, 3*time.Second); err != nil {
		return nil, err
	}
	return rsp, nil
}

func GetRedirectUrlWithRetryTimes(httpUrl string, retryTimes int) (string, error) {
	rsp, err := RetryDoHTTPWithErr(func() (*fasthttp.Response, error) {
		return get(httpUrl, 3*time.Second)
	}, retryTimes, 0)
	if err != nil {
		return "", err
	}
	defer fasthttp.ReleaseResponse(rsp)
	return string(rsp.Header.Peek("Location")), nil
}

func get(httpUrl string, timeout time.Duration) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(consts.MethodGet)
	req.SetRequestURI(httpUrl)
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; " +
		"Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
	rsp := fasthttp.AcquireResponse()
	if err := fasthttp.DoTimeout(req, rsp, timeout); err != nil {
		return nil, err
	}
	return rsp, nil
}

func getM3u8Md5WithUrl(urlPath string, key string) (string, error) {
	u, err := url.Parse(urlPath)
	if err != nil {
		return "", err
	}
	md5Str := u.Query().Get(key)
	if md5Str == "" {
		return "", errors.New("未带校验值")
	}
	return u.Query().Get("h"), nil
}

func checkM3u8TsMD5(body []byte, md5Value string) (bool, error) {
	h := md5.New()
	_, err := io.WriteString(h, string(body))
	if err != nil {
		return false, err
	}
	tsMD5 := fmt.Sprintf("%x", h.Sum(nil))
	if tsMD5 != md5Value {
		return false, nil
	}

	return true, nil

}
