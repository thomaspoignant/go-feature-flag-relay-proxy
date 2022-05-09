package model

import "time"

type Info struct{}

// InfoResponse is the object returned by the info API
type InfoResponse struct {
	LatestCacheRefresh time.Time `json:"cacheRefresh"`
}
