package types

import "context"

type TriggerRequest struct {
	JobId     string `json:"job_id"`
	DriveLink string `json:"drive_link"`
}

type Drive interface {
	GetFileList(string) ([]string, error)
	GetFileContent(string) ([]byte, error)
	GetUsername([]byte) ([]string, error)
}

type Queue interface {
	Publish(string, any) error
}

type StatusQueue struct {
	JobId     string
	FileId    string
	GithubIDs []string
}

type KVStore interface {
	Set(ctx context.Context, key any, value any) error
	Get(ctx context.Context, key any) (value any, err error)
}
