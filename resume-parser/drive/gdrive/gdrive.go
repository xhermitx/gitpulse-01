package gdrive

import "bytes"

type GoogleDrive struct {
	FolderId string
}

func NewGoogleDrive(folderId string) *GoogleDrive {
	return &GoogleDrive{
		FolderId: folderId,
	}
}

func (g *GoogleDrive) GetFileList() ([]string, error) {
	return nil, nil
}

func (g *GoogleDrive) GetFileContent(fileName string) ([]byte, error) {
	return nil, nil
}

func (g *GoogleDrive) GetUsername(fileContent *bytes.Buffer) (string, error) {
	return "", nil
}
