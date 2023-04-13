package fasthttp

import (
	"github.com/baoyxing/go-tools/net/consts"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"math/rand"
	"time"
)

type Func func() error

type BodyFunc func() ([]byte, error)

type HTTPFunc func() (*fasthttp.Response, error)

type BodyHTTPFunc func() ([]byte, error)

func RetryDoWithErr(fn Func, retries int, sleep time.Duration) error {
	if sleep == 0 {
		sleep = consts.DefaultSleep
	}
	if err := fn(); err != nil {
		retries--
		if retries <= 0 {
			return err
		}
		sleep += (time.Duration(rand.Int63n(int64(sleep)))) / 2
		time.Sleep(sleep)
		return RetryDoWithErr(fn, retries, 2*sleep)
	}

	return nil
}

func RetryDoHTTPWithErr(fn HTTPFunc, retries int, sleep time.Duration) (*fasthttp.Response, error) {
	var res *fasthttp.Response
	err := RetryDoWithErr(func() error {
		var err error
		res, err = fn()
		return err
	}, retries, sleep)

	return res, err
}

func RetryDoWithEmptyBody(fn BodyFunc, retries int, sleep time.Duration) ([]byte, error) {
	if sleep == 0 {
		sleep = consts.DefaultSleep
	}
	body, err := fn()
	if len(body) == 0 || err != nil {
		retries--
		if retries <= 0 {
			errMsg := ""
			if err != nil {
				errMsg += err.Error()
			}
			if len(body) == 0 {
				errMsg += "OR body empty"
			}

			return nil, errors.New(errMsg)
		}
		sleep += (time.Duration(rand.Int63n(int64(sleep)))) / 2
		time.Sleep(sleep)
		return RetryDoWithEmptyBody(fn, retries, 2*sleep)
	}
	return nil, err
}

func RetryDoHTTPWithEmptyBody(fn BodyHTTPFunc, retries int, sleep time.Duration) ([]byte, error) {
	var res []byte
	_, err := RetryDoWithEmptyBody(func() ([]byte, error) {
		var err error
		res, err = fn()
		return res, err
	}, retries, sleep)
	return res, err
}
