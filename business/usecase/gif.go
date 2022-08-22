package usecase

import (
	"time"

	"github.com/bluele/gcache"
)

type gifGetter interface {
	GetGIF() (string, error)
}

type ProcessGIF struct {
	getter gifGetter
	cache  gcache.Cache
}

func NewProcessGIF(gif gifGetter, cache gcache.Cache) ProcessGIF {
	return ProcessGIF{getter: gif, cache: cache}
}

func (p *ProcessGIF) GetGif() (string, error) {
	tf := time.Now().Format("2006-01-02")
	urlc, err := p.cache.Get(tf)
	switch err {
	case gcache.KeyNotFoundError:
		url, err := p.getter.GetGIF()
		if err != nil {
			return "", err
		}
		p.cache.SetWithExpire(tf, url, 24*time.Hour)
		return url, nil
	case nil:
		return urlc.(string), nil
	}

	return "", err
}
