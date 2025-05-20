package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	apiKey  = "<key>"
	token   = "<key>"
	cardID  = "6809bcb36d41988a3d6f2aac" // your card ID
	baseURL = "https://api.trello.com/1"
)

type Attachment struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func main() {
	attachments := fetchAttachments(cardID)

	fmt.Println("📎 Attachments on card:")
	for _, att := range attachments {
		fmt.Printf("- %s (%s)\n  → %s\n", att.Name, att.ID, att.URL)
	}

	if len(attachments) > 0 {
		fmt.Println("\n⬇️  Attempting to download first attachment...")
		err := downloadAttachment(cardID, attachments[0].ID, attachments[0].Name)
		if err != nil {
			fmt.Println("❌ Download failed:", err)
		} else {
			fmt.Println("✅ Download successful!")
		}
	}
}

func fetchAttachments(cardID string) []Attachment {
	url := fmt.Sprintf("%s/cards/%s/attachments?fields=name,url,mimeType&key=%s&token=%s", baseURL, cardID, apiKey, token)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		panic(fmt.Errorf("API error %d: %s", resp.StatusCode, string(body)))
	}

	var attachments []Attachment
	json.NewDecoder(resp.Body).Decode(&attachments)
	return attachments
}

func downloadAttachment(cardID, attachmentID, filename string) error {
	downloadURL := fmt.Sprintf("%s/cards/%s/attachments/%s/download", baseURL, cardID, attachmentID)
	client := &http.Client{}

	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return err
	}

	// Mimic a browser
	//req.Header.Set("Authorization", fmt.Sprintf("OAuth %s %s", apiKey, token))
	//req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", token))
	// req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", apiKey))
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s %s", apiKey, token))
	//req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	//req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Authorization", fmt.Sprintf("OAuth oauth_consumer_key=\"%s\", oauth_token=\"%s\"", apiKey, token))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("Downloaded %s\n", filename)

	localPath := filepath.Join("downloads", filename)
	//os.MkdirAll("downloads", os.ModePerm)
	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Downloaded:", localPath)
	return err
}

func downloadAttachment2(cardID, downloadURL, filename string) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return err
	}

	// Mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	out, err := os.Create("downloads/" + filename)
	if err != nil {
		return err
	}
	defer out.Close()

	os.MkdirAll("downloads", os.ModePerm)
	_, err = io.Copy(out, resp.Body)
	return err
}
