package fasthttp

import (
	"github.com/baoyxing/go-tools/net/consts"
	"github.com/valyala/fasthttp"
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

func DownloadTsDataWithRetryTimes(httpUrl string, start, end, retryTimes int) ([]byte, error) {
	return RetryDoHTTPWithEmptyBody(func() ([]byte, error) {
		return downloadTsData(httpUrl, start, end)
	}, retryTimes, 0)
}

func downloadTsData(url string, start, end int) ([]byte, error) {
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
	if err := fasthttp.DoTimeout(req, rsp, 30*time.Second); err != nil {
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
