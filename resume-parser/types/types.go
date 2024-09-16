package types

import (
	"context"
	"time"
)

type TriggerRequest struct {
	JobId     string `json:"job_id"`
	DriveLink string `json:"drive_link"`
}

type Drive interface {
	GetFileList(string) (map[string]string, error)
	GetFileContent(string) ([]byte, error)
	GetUsername([]byte) ([]string, error)
}

type Queue interface {
	Publish(string, any) error
}

type JobQueue struct {
	JobId     string   `json:"job_id"`
	Filename  string   `json:"filename"`
	GithubIDs []string `json:"github_ids"`
	Status    bool     `json:"status"`
}

type KVStore interface {
	Get(ctx context.Context, key string) (value string, err error)
	Set(ctx context.Context, key string, value any, t time.Duration) error
	Append(ctx context.Context, key string, value string) error
}

type UnparsedFilesCache struct {
	JobId         string
	UnparsedFiles []string // Name of the files
}
