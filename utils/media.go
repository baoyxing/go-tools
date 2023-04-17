package utils

import (
	"bytes"
	"context"
	"github.com/baoyxing/go-tools/net/fasthttp"
	"github.com/grafov/m3u8"
	"github.com/grafov/m3u8/example/template"
	nUrl "net/url"
	"time"
)

type MediaType uint

const (
	Live MediaType = iota + 1 //直播
	Vod                       //点播
)

type MediaPlaylist struct {
	TargetDuration float64
	SeqNo          uint64 // EXT-X-MEDIA-SEQUENCE
	Segments       []*MediaSegment
	Ver            uint8
	Count          uint
	MediaLength    uint64
}

type MediaSegment struct {
	Size            uint64
	Url             string
	Duration        float64
	ProgramDateTime time.Time
	Start           uint64
	End             uint64
}

func CheckM3u8MediaType(url string) (MediaType, error) {
	customTags := []m3u8.CustomDecoder{
		&template.CustomPlaylistTag{},
		&template.CustomSegmentTag{},
	}
	body, err := fasthttp.Get(url)
	p, listType, err := m3u8.DecodeWith(body, false, customTags)

	if err != nil {
		return Vod, err
	}
	switch listType {
	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist)
		if mediapl.Closed {
			return Vod, nil
		} else {
			return Live, nil
		}
	case m3u8.MASTER:
		noSuffixUrl := GetNoSuffixUrl(url)
		masterpl := p.(*m3u8.MasterPlaylist)
		_url, err := nUrl.Parse(noSuffixUrl)
		if err != nil {
			return Vod, err
		}
		url = ResolveURL(_url, masterpl.Variants[0].URI)
		return CheckM3u8MediaType(url)
	}
	return Vod, err
}

func ParseM3u8(ctx context.Context, url string) (*MediaPlaylist, error) {
	body, err := fasthttp.Get(url)
	if err != nil {
		return nil, err
	}
	customTags := []m3u8.CustomDecoder{
		&template.CustomPlaylistTag{},
		&template.CustomSegmentTag{},
	}
	p, listType, err := m3u8.DecodeWith(*bytes.NewBuffer(body), false, customTags)
	if err != nil {
		return nil, err
	}
	switch listType {
	case m3u8.MEDIA:
		noSuffixUrl := GetNoFileUrl(url)
		_, err := nUrl.Parse(noSuffixUrl)
		if err != nil {
			return nil, err
		}
		mediapl := p.(*m3u8.MediaPlaylist)
		count := mediapl.Count()
		mediaPlaylist := &MediaPlaylist{
			TargetDuration: mediapl.TargetDuration,
			SeqNo:          mediapl.SeqNo,
			Segments:       make([]*MediaSegment, 0, mediapl.Count()),
			Ver:            mediapl.Version(),
			Count:          count,
			MediaLength:    0,
		}
		for _, segment := range mediapl.Segments {
			if segment != nil {
				mediaSegment := &MediaSegment{
					Url:             segment.URI,
					Duration:        segment.Duration,
					ProgramDateTime: segment.ProgramDateTime,
				}
				mediaPlaylist.Segments = append(mediaPlaylist.Segments, mediaSegment)
			}
		}
		return mediaPlaylist, nil
	case m3u8.MASTER:
		noSuffixUrl := GetNoSuffixUrl(url)
		masterpl := p.(*m3u8.MasterPlaylist)
		_url, err := nUrl.Parse(noSuffixUrl)
		if err != nil {
			return nil, err
		}
		url = ResolveURL(_url, masterpl.Variants[0].URI)
		return ParseM3u8(ctx, url)
	}
	return nil, nil
}

func ParseM3u8Ts(ctx context.Context, url string) (tsUrls []string, duration float64, tsDuration []float64, err error) {
	customTags := []m3u8.CustomDecoder{
		&template.CustomPlaylistTag{},
		&template.CustomSegmentTag{},
	}
	body, err := fasthttp.Get(url)
	p, listType, err := m3u8.DecodeWith(*bytes.NewBuffer(body), false, customTags)
	if err != nil {
		return nil, 0, nil, err
	}
	if tsUrls == nil {
		tsUrls = make([]string, 0)
	}
	if tsDuration == nil {
		tsDuration = make([]float64, 0)
	}
	switch listType {
	case m3u8.MEDIA:
		noSuffixUrl := GetNoFileUrl(url)
		_url, err := nUrl.Parse(noSuffixUrl)
		if err != nil {
			return nil, 0, nil, err
		}
		mediapl := p.(*m3u8.MediaPlaylist)
		for _, segment := range mediapl.Segments {
			if segment != nil {
				segmentUrl := ResolveURL(_url, segment.URI)
				duration += segment.Duration
				tsUrls = append(tsUrls, segmentUrl)
				tsDuration = append(tsDuration, segment.Duration)
			}
		}
	case m3u8.MASTER:
		noSuffixUrl := GetNoSuffixUrl(url)
		masterpl := p.(*m3u8.MasterPlaylist)
		_url, err := nUrl.Parse(noSuffixUrl)
		if err != nil {
			return nil, 0, nil, err
		}
		url = ResolveURL(_url, masterpl.Variants[0].URI)
		tsUrls, duration, tsDuration, err = ParseM3u8Ts(ctx, url)
	}
	return

}
