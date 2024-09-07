package types

type TriggerRequest struct {
	JobId     string `json:"job_id"`
	DriveLink string `json:"drive_link"`
}

type Drive interface {
	GetFileList() ([]string, error)
	GetFileContent(string) ([]byte, error)
	GetUsername([]byte) ([]string, error)
}
