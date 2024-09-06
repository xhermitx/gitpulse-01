package gdrive

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/xhermitx/gitpulse-01/resume-parser/config"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

type GoogleDrive struct {
	FolderId string
}

func NewGoogleDrive(folderId string) *GoogleDrive {
	return &GoogleDrive{
		FolderId: folderId,
	}
}

func (g *GoogleDrive) GetFileList() ([]string, error) {
	ctx := context.Background()

	// SERVICE ACCOUNT FILE
	jsonKey, err := os.ReadFile(config.Envs.ServiceAccount)
	if err != nil {
		log.Println("Error reading Credentials")
		return nil, err
	}

	// CREATE CONFIGS USING AUTHENTICATION
	config, err := google.JWTConfigFromJSON(data, drive.DriveReadonlyScope)
	if err != nil {
		log.Println("Error getting JWT Configs")
		return nil, err
	}
	client := config.Client(ctx)

	// CREATE A NEW DRIVE CLIENT
	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Println("Error creating Service")
		return nil, err
	}

	// QUERY TO READ FILES RESIDING IN FOLDERS
	query := fmt.Sprintf("'%s' in parents", g.FolderId)

	fileList, err := driveService.Files.List().
		Q(query).
		Fields("nextPageToken, files(id, name)").
		Do()
	if err != nil {
		log.Println("Error fetching the file list")
		return nil, err
	}

	return fileList, nil
}

func (g *GoogleDrive) GetFileContent(fileName string) ([]byte, error) {
	return nil, nil
}

func (g *GoogleDrive) GetUsername(fileContent *bytes.Buffer) (string, error) {
	return "", nil
}
