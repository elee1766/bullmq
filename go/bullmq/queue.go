package bullmq

import (
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ObliterateOpts struct {
	// Use force = true to force obliteration even with active jobs in the queue
	Force bool `json:"force"`
	// Default count is 1000
	Count *int `json:"count,omitempty"`
}

type Queue struct {
	Token uuid.UUID

	JobOpts *BaseJobOptions
}

type QueueBase struct {
	ToKey   func(string) string
	Keys    map[string]string
	Closing func()

	closed    bool
	scripts   *Scripts
	redisConn redis.Cmder
}
