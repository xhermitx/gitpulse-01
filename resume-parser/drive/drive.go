package drive

import "bytes"

type Drive interface {
	GetFileList() ([]string, error)
	GetFileContent(string) ([]byte, error)
	GetUsername(bytes.Buffer) (string, error)
}
