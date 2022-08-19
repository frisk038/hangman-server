package giphy

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
)

const apiKey = "uNVGhSvMfigaxr9VytrIm5ZNA3FDJITI"
const apiURL = "https://api.giphy.com/v1/gifs/search"
const apiQWord = "clapping"
const apiLimit = 50
const apiOffset = 10

type Conf struct {
	key    string
	url    string
	qWord  string
	limit  int
	offset int
}

type Client struct {
	client http.Client
	conf   Conf
}

func NewGiphy() *Client {
	return &Client{
		client: *http.DefaultClient,
		conf: Conf{
			key:    apiKey,
			url:    apiURL,
			qWord:  apiQWord,
			limit:  apiLimit,
			offset: apiOffset,
		},
	}
}

func (g *Client) GetGIF() (string, error) {
	req, err := http.NewRequest("GET", g.conf.url, nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Add("api_key", g.conf.key)
	q.Add("q", g.conf.qWord)
	q.Add("limit", strconv.Itoa(g.conf.limit))
	q.Add("offset", strconv.Itoa(rand.Intn(g.conf.offset)))
	req.URL.RawQuery = q.Encode()

	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	r := struct {
		Data []struct {
			ID     string `json:"id"`
			Images struct {
				Downsized struct {
					URL string `json:"url"`
				} `json:"downsized"`
			} `json:"images"`
			URL string `json:"url"`
		} `json:"data"`
		Meta struct {
			Status int    `json:"status"`
			MSG    string `json:"msg"`
		} `json:"meta"`
	}{}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", err
	}

	return r.Data[rand.Intn(g.conf.limit)].Images.Downsized.URL, nil
}
