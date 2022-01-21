package model

type Error struct {
	ErrorCode int      `json:"errorCode"`
	Timestamp int64    `json:"timestamp"`
	Errors    []string `json:"errors"`
}
