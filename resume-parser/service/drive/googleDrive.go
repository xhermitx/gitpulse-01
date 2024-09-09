package drive

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/xhermitx/gitpulse-01/resume-parser/config"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	pattern = `github\.com/[a-zA-Z0-9]+(\-[a-zA-Z0-9]*)*`
	offset  = 11 // To match id after "github.com/"
)

type GoogleDrive struct {
	DriveService *drive.Service
}

func NewGoogleDrive(ds *drive.Service) *GoogleDrive {
	return &GoogleDrive{
		DriveService: ds,
	}
}

func NewGoogleService() (*drive.Service, error) {
	ctx := context.Background()

	// Service Account File
	jsonKey, err := os.ReadFile(config.Envs.ServiceAccount)
	if err != nil {
		log.Println("Error reading Credentials")
		return nil, err
	}

	config, err := google.JWTConfigFromJSON(jsonKey, drive.DriveReadonlyScope)
	if err != nil {
		log.Println("Error getting JWT Configs")
		return nil, err
	}
	client := config.Client(ctx)

	DriveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Println("Error creating Service")
		return nil, err
	}
	return DriveService, nil
}

func (g *GoogleDrive) GetFileList(folderId string) ([]string, error) {
	// Read files from the Folder
	query := fmt.Sprintf("'%s' in parents", folderId)
	res, err := g.DriveService.Files.List().
		Q(query).
		Fields("nextPageToken, files(id, name)").
		Do()
	if err != nil {
		log.Println("Error fetching the file list")
		return nil, err
	}

	var fileList []string
	for _, file := range res.Files {
		// TODO: Store file.Id:file.Name along with its
		//       status(parsed/unparsed) in Redis
		// 		 To be referred when there is an error
		fileList = append(fileList, file.Id)
	}
	return fileList, nil
}

func (g *GoogleDrive) GetFileContent(fileId string) ([]byte, error) {
	file, err := g.DriveService.Files.Get(fileId).Fields("mimeType").Do()
	if err != nil {
		return nil, err
	}
	if file.MimeType != "application/pdf" {
		// FIXME: Add support for other file types
		return nil, err
	}

	resp, err := g.DriveService.Files.Get(fileId).Download()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (g *GoogleDrive) GetUsername(fileContent []byte) ([]string, error) {
	var (
		fileText = string(fileContent)
		pattern  = regexp.MustCompile(pattern)
		uniqIDs  = make(map[string]bool)
	)

	matches := pattern.FindAllString(fileText, -1)
	for _, match := range matches {
		uniqIDs[match[offset:]] = true
	}
	if len(uniqIDs) == 0 {
		return nil, fmt.Errorf("no username found in file")
	}

	userIDs := make([]string, 0, len(uniqIDs))
	for key := range uniqIDs {
		userIDs = append(userIDs, key)
	}

	return userIDs, nil
}
