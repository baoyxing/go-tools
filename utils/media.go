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
		return ParseM3u8Ts(ctx, url)
	}
	return nil, nil
}
