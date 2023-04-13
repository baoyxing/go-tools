package utils

import (
	"math"
	"net"
	"net/url"
	"path"
	"strings"
)

func ResolveURL(u *url.URL, p string) string {
	if strings.HasPrefix(p, "https://") || strings.HasPrefix(p, "http://") {
		return p
	}
	var baseURL string
	if strings.Index(p, "/") == 0 {
		baseURL = u.Scheme + "://" + u.Host
	} else {
		tU := u.String()
		baseURL = tU[0:strings.LastIndex(tU, "/")]
	}
	return baseURL + path.Join("/", p)
}

func GetNoFileUrl(url string) string {
	last := reFind(url, '/')
	noSuffixUrl := []byte(url)[0 : last+1]

	return string(noSuffixUrl)
}

func GetNoSuffixUrl(url string) string {

	last, isExist := psFind(url, '?')
	noSuffixUrl := []byte(url)
	if isExist {
		noSuffixUrl = []byte(url)[0:last]
	}
	return string(noSuffixUrl)
}

func GetLoalIPv4() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "未知内网IP"
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "未知内网IP"
		}
		for _, addr := range addrs {
			ip := getIpv4FromAddr(addr)
			if ip == nil {
				continue
			}
			return ip.String()
		}
	}
	return "未知内网IP"
}

func CalcPieceSize(size uint64) int {
	if size <= 2097152 {
		return 16384
	}
	n := int(math.Sqrt(float64(size) / 4096))
	i := 3
	for {
		tmp := n & (32768 >> i)
		if tmp > 0 {
			return tmp * 1024
		}
		if i < 16 {
			i++
		} else {
			break
		}
	}
	return 16384
}

func getIpv4FromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

//倒序查找
func reFind(url string, symbol byte) int {
	startIndex := len(url) - 1
	var i int
	for i = startIndex; []byte(url)[i] != symbol; i-- {
	}
	return i
}

//正序查找
func psFind(url string, symbol byte) (int, bool) {
	urlByte := []byte(url)
	urlByteLength := len(urlByte)
	index := 0
	isExist := false
	for i := 0; i < urlByteLength; i++ {
		if urlByte[i] == symbol {
			index = i
			isExist = true
			break
		}
	}
	return index, isExist
}
