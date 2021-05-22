package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	apiKey    string
	channelId string

	requestURL string

	lastCache     string
	lastCacheTime time.Time
	lastCacheLock sync.Mutex
)

type parameter struct {
	key   string
	value string
}

type response struct {
	Etag  string `json:"etag"`
	Items []struct {
		Etag string `json:"etag"`
		ID   struct {
			Kind    string `json:"kind"`
			VideoID string `json:"videoId"`
		} `json:"id"`
		Kind    string `json:"kind"`
		Snippet struct {
			ChannelID            string `json:"channelId"`
			ChannelTitle         string `json:"channelTitle"`
			Description          string `json:"description"`
			LiveBroadcastContent string `json:"liveBroadcastContent"`
			PublishTime          string `json:"publishTime"`
			PublishedAt          string `json:"publishedAt"`
			Thumbnails           struct {
				Default struct {
					Height int64  `json:"height"`
					URL    string `json:"url"`
					Width  int64  `json:"width"`
				} `json:"default"`
				High struct {
					Height int64  `json:"height"`
					URL    string `json:"url"`
					Width  int64  `json:"width"`
				} `json:"high"`
				Medium struct {
					Height int64  `json:"height"`
					URL    string `json:"url"`
					Width  int64  `json:"width"`
				} `json:"medium"`
			} `json:"thumbnails"`
			Title string `json:"title"`
		} `json:"snippet"`
	} `json:"items"`
	Kind          string `json:"kind"`
	NextPageToken string `json:"nextPageToken"`
	PageInfo      struct {
		ResultsPerPage int64 `json:"resultsPerPage"`
		TotalResults   int64 `json:"totalResults"`
	} `json:"pageInfo"`
	RegionCode string `json:"regionCode"`
}

const (
	emptyRequestURL = "request url is empty"

	cacheTime  = 2
	maxResults = "3"
)

func init() {
	apiKey = os.Getenv("apiKey")
	channelId = os.Getenv("channelId")

	ctx := context.Background()

	var err error
	requestURL, err = makeRequestURL(ctx)
	if err != nil {
		log.Println(err)
	}

	lastCacheLock.Lock()
	defer lastCacheLock.Unlock()

	lastCache, err = getData(ctx)
	if err != nil {
		log.Println(err)
	}

	lastCacheTime = time.Now()

}

func getParameters() []parameter {
	return []parameter{
		{
			key:   "part",
			value: "snippet",
		},
		{
			key:   "channelId",
			value: channelId,
		},
		{
			key:   "maxResults",
			value: maxResults,
		},
		{
			key:   "order",
			value: "date",
		},
		{
			key:   "type",
			value: "video",
		},
		{
			key:   "key",
			value: apiKey,
		},
	}

}

func makeRequestURL(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/youtube/v3/search", nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	for _, param := range getParameters() {
		q.Add(param.key, param.value)
	}

	req.URL.RawQuery = q.Encode()
	return req.URL.String(), nil
}

func getData(ctx context.Context) (r string, err error) {
	if requestURL == "" {
		return r, errors.New(emptyRequestURL)
	}

	resp, err := http.Get(requestURL)
	if err != nil {
		return r, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	// unmarshall to response type to confirm it's valid
	var respStruct response
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return r, err
	}

	// then marshal back to string
	rByte, err := json.Marshal(respStruct)
	if err != nil {
		return r, err
	}

	return string(rByte), nil

}

func cacheExpired() bool {
	return time.Now().Unix() > lastCacheTime.Add(time.Minute*cacheTime).Unix()
}

func Handler(ctx context.Context) (string, error) {
	lastCacheLock.Lock()
	defer lastCacheLock.Unlock()

	var err error
	if lastCache == "" {
		log.Println("no cache, getting data")
		lastCache, err = getData(ctx)
		if err != nil {
			return lastCache, err
		}

		lastCacheTime = time.Now()

	} else if cacheExpired() {
		log.Println("cache expired, getting data")
		lastCache, err = getData(ctx)
		if err != nil {
			return lastCache, err
		}

	} else {
		log.Println("hitting cache")
	}

	lastCacheTime = time.Now()

	return lastCache, nil
}

func main() {
	lambda.Start(Handler)
}
