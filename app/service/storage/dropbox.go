package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/khaledhikmat/tr-extractor/service/config"
)

type dropboxService struct {
	CfgSvc config.IService
}

func NewDropbox(cfgsvc config.IService) IService {
	return &dropboxService{
		CfgSvc: cfgsvc,
	}
}

func (svc *dropboxService) Upload(filePath, folder, identifier string) (string, error) {
	dropboxPath := fmt.Sprintf("%s%s#%s", svc.CfgSvc.GetDropboxUploadPath(), folder, identifier)
	fmt.Printf("Uploading to Dropbox: %s\n", dropboxPath)
	// if 0 == 0 {
	// 	return filePath, nil
	// }

	defer func() {
		// Delete the local file
		_ = os.Remove(filePath)
	}()

	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	body, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/upload", bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+svc.CfgSvc.GetDropboxAccessToken())
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Dropbox-API-Arg", fmt.Sprintf(`{"path":"/%s","mode":"add","autorename":true}`, dropboxPath))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Uploading to Dropbox failed: %s\n", responseBody)
		return "", fmt.Errorf("failed to upload to Dropbox: %s", responseBody)
	}

	// Now create a shared link
	// sharedURL, err := svc.createDropboxShareLink(dropboxPath)
	// if err != nil {
	// 	return "", err
	// }

	// return sharedURL, nil
	return dropboxPath, nil
}

func (svc *dropboxService) createDropboxShareLink(path string) (string, error) {
	reqBody := map[string]interface{}{
		"path": path,
		"settings": map[string]string{
			"requested_visibility": "public",
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/create_shared_link_with_settings", bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+svc.CfgSvc.GetDropboxAccessToken())
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to create share link: %s\n%s", resp.Status, body)
	}

	var result struct {
		URL string `json:"url"`
	}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// Optional: convert ?dl=0 to ?raw=1 or ?dl=1 if you want direct download
	sharedURL := result.URL
	if len(sharedURL) > 0 {
		sharedURL = sharedURL[:len(sharedURL)-1] + "1" // change ?dl=0 to ?dl=1
	}

	return sharedURL, nil
}
