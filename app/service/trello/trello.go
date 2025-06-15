package trello

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/khaledhikmat/tr-extractor/service/config"
)

type trelloService struct {
	CfgSvc config.IService
}

func New(cfgsvc config.IService) IService {
	return &trelloService{
		CfgSvc: cfgsvc,
	}
}

func (svc *trelloService) RetrieveProperties(_ int) ([]TRProperty, error) {
	var results []TRProperty

	boardID := svc.CfgSvc.GetTrelloPropertiesBoardID()

	customFieldDefs, err := fetchCustomFieldDefs(svc.CfgSvc.GetTrelloBaseURL(), svc.CfgSvc.GetTrelloAPIKey(), svc.CfgSvc.GetTrelloToken(), boardID)
	if err != nil {
		return results, err
	}

	searchURL := fmt.Sprintf("%s/boards/%s/cards?key=%s&token=%s",
		svc.CfgSvc.GetTrelloBaseURL(), boardID, svc.CfgSvc.GetTrelloAPIKey(), svc.CfgSvc.GetTrelloToken())
	props, err := fetchTrelloEntities[TRProperty](searchURL)
	if err != nil {
		return results, err
	}

	for _, prop := range props {
		customFields, err := fetchCustomFields(svc.CfgSvc.GetTrelloBaseURL(), svc.CfgSvc.GetTrelloAPIKey(), svc.CfgSvc.GetTrelloToken(), prop.ID)
		if err != nil {
			return results, err
		}

		// Exclude properties without custom fields
		if len(customFields) == 0 {
			continue
		}

		for _, cf := range customFields {
			field := TRField{}

			def := customFieldDefs[cf.IDCustomField]
			if def.Type == "list" {
				for _, opt := range def.Options {
					if opt.ID == cf.IDValue {
						field.Name = def.Name
						field.Type = "list"
						field.Value = opt.Value.Text
					}
				}
			} else {
				for k, v := range cf.Value {
					field.Name = def.Name
					field.Type = k
					field.Value = fmt.Sprintf("%v", v)
				}
			}

			prop.Fields = append(prop.Fields, field)
			if field.Name == "Location AR" {
				prop.LocationAR = field.Value
			} else if field.Name == "Location EN" {
				prop.LocationEN = field.Value
			} else if field.Name == "Lot" {
				prop.Lot = field.Value
			} else if field.Name == "Type" {
				prop.Type = field.Value
			} else if field.Name == "Status" {
				prop.Status = field.Value
			} else if field.Name == "Owner" {
				prop.Owner = field.Value
			} else if field.Name == "Area" {
				result, err := strconv.ParseFloat(field.Value, 64)
				if err != nil {
					result = 0
				}
				prop.Area = result
			} else if field.Name == "Shares" {
				result, err := strconv.ParseFloat(field.Value, 64)
				if err != nil {
					result = 0
				}
				prop.Shares = result
			} else if field.Name == "Organized" {
				result, err := strconv.ParseBool(field.Value)
				if err != nil {
					result = false
				}
				prop.Organized = result
			} else if field.Name == "Effects" {
				result, err := strconv.ParseBool(field.Value)
				if err != nil {
					result = false
				}
				prop.Effects = result
			}
		}

		attachments, err := fetchAttachments(svc.CfgSvc.GetTrelloBaseURL(), svc.CfgSvc.GetTrelloAPIKey(), svc.CfgSvc.GetTrelloToken(), prop.ID)
		if err != nil {
			return results, err
		}
		prop.Attachments = append(prop.Attachments, attachments...)

		comments, err := fetchComments(svc.CfgSvc.GetTrelloBaseURL(), svc.CfgSvc.GetTrelloAPIKey(), svc.CfgSvc.GetTrelloToken(), prop.ID)
		if err != nil {
			return results, err
		}
		prop.Comments = append(prop.Comments, comments...)

		results = append(results, prop)
	}

	return results, nil
}

func (svc *trelloService) RetrieveInheritanceConfinments(_ int) ([]TRInheritanceConfinement, error) {
	var results []TRInheritanceConfinement

	boardID := svc.CfgSvc.GetTrelloInheritanceConfinmentsBoardID()

	customFieldDefs, err := fetchCustomFieldDefs(svc.CfgSvc.GetTrelloBaseURL(), svc.CfgSvc.GetTrelloAPIKey(), svc.CfgSvc.GetTrelloToken(), boardID)
	if err != nil {
		return results, err
	}

	searchURL := fmt.Sprintf("%s/boards/%s/cards?key=%s&token=%s",
		svc.CfgSvc.GetTrelloBaseURL(), boardID, svc.CfgSvc.GetTrelloAPIKey(), svc.CfgSvc.GetTrelloToken())
	entities, err := fetchTrelloEntities[TRInheritanceConfinement](searchURL)
	if err != nil {
		return results, err
	}

	for _, entity := range entities {
		customFields, err := fetchCustomFields(svc.CfgSvc.GetTrelloBaseURL(), svc.CfgSvc.GetTrelloAPIKey(), svc.CfgSvc.GetTrelloToken(), entity.ID)
		if err != nil {
			return results, err
		}

		// Exclude properties without custom fields
		if len(customFields) == 0 {
			continue
		}

		for _, cf := range customFields {
			field := TRField{}

			def := customFieldDefs[cf.IDCustomField]
			if def.Type == "list" {
				for _, opt := range def.Options {
					if opt.ID == cf.IDValue {
						field.Name = def.Name
						field.Type = "list"
						field.Value = opt.Value.Text
					}
				}
			} else {
				for k, v := range cf.Value {
					field.Name = def.Name
					field.Type = k
					field.Value = fmt.Sprintf("%v", v)
				}
			}

			entity.Fields = append(entity.Fields, field)
			if field.Name == "Generation" {
				result, err := strconv.ParseInt(field.Value, 10, 64)
				if err != nil {
					result = 0
				}
				entity.Generation = result
			} else if field.Name == "Title" {
				entity.Title = field.Value
			}
		}

		if entity.Title == "" {
			entity.Title = entity.Name
		}

		attachments, err := fetchAttachments(svc.CfgSvc.GetTrelloBaseURL(), svc.CfgSvc.GetTrelloAPIKey(), svc.CfgSvc.GetTrelloToken(), entity.ID)
		if err != nil {
			return results, err
		}
		entity.Attachments = append(entity.Attachments, attachments...)

		comments, err := fetchComments(svc.CfgSvc.GetTrelloBaseURL(), svc.CfgSvc.GetTrelloAPIKey(), svc.CfgSvc.GetTrelloToken(), entity.ID)
		if err != nil {
			return results, err
		}
		entity.Comments = append(entity.Comments, comments...)

		results = append(results, entity)
	}

	return results, nil
}

func (svc *trelloService) RetrieveSupportiveDocs(_ int) ([]TRSupportiveDoc, error) {
	var results []TRSupportiveDoc

	boardID := svc.CfgSvc.GetTrelloSupportiveDocsBoardID()

	customFieldDefs, err := fetchCustomFieldDefs(svc.CfgSvc.GetTrelloBaseURL(), svc.CfgSvc.GetTrelloAPIKey(), svc.CfgSvc.GetTrelloToken(), boardID)
	if err != nil {
		return results, err
	}

	searchURL := fmt.Sprintf("%s/boards/%s/cards?key=%s&token=%s",
		svc.CfgSvc.GetTrelloBaseURL(), boardID, svc.CfgSvc.GetTrelloAPIKey(), svc.CfgSvc.GetTrelloToken())
	entities, err := fetchTrelloEntities[TRSupportiveDoc](searchURL)
	if err != nil {
		return results, err
	}

	for _, entity := range entities {
		customFields, err := fetchCustomFields(svc.CfgSvc.GetTrelloBaseURL(), svc.CfgSvc.GetTrelloAPIKey(), svc.CfgSvc.GetTrelloToken(), entity.ID)
		if err != nil {
			return results, err
		}

		// Exclude properties without custom fields
		if len(customFields) == 0 {
			continue
		}

		for _, cf := range customFields {
			field := TRField{}

			def := customFieldDefs[cf.IDCustomField]
			if def.Type == "list" {
				for _, opt := range def.Options {
					if opt.ID == cf.IDValue {
						field.Name = def.Name
						field.Type = "list"
						field.Value = opt.Value.Text
					}
				}
			} else {
				for k, v := range cf.Value {
					field.Name = def.Name
					field.Type = k
					field.Value = fmt.Sprintf("%v", v)
				}
			}

			entity.Fields = append(entity.Fields, field)
			if field.Name == "Category" {
				entity.Category = field.Value
			} else if field.Name == "Title" {
				entity.Title = field.Value
			}
		}

		if entity.Title == "" {
			entity.Title = entity.Name
		}

		attachments, err := fetchAttachments(svc.CfgSvc.GetTrelloBaseURL(), svc.CfgSvc.GetTrelloAPIKey(), svc.CfgSvc.GetTrelloToken(), entity.ID)
		if err != nil {
			return results, err
		}
		entity.Attachments = append(entity.Attachments, attachments...)

		comments, err := fetchComments(svc.CfgSvc.GetTrelloBaseURL(), svc.CfgSvc.GetTrelloAPIKey(), svc.CfgSvc.GetTrelloToken(), entity.ID)
		if err != nil {
			return results, err
		}
		entity.Comments = append(entity.Comments, comments...)

		results = append(results, entity)
	}

	return results, nil
}

func (svc *trelloService) DownloadAttachment(url string) (string, string, string, error) {
	// Extract the card ID and attachment ID from the URL
	cardID, attachmentID, extension, err := extractTrelloIDsAndExt(url)
	if err != nil {
		return "", "", "", err
	}

	baseURL := svc.CfgSvc.GetTrelloBaseURL()
	apiKey := svc.CfgSvc.GetTrelloAPIKey()
	token := svc.CfgSvc.GetTrelloReadToken()

	filename := fmt.Sprintf("%s%s", attachmentID, extension)
	downloadURL := fmt.Sprintf("%s/cards/%s/attachments/%s/download", baseURL, cardID, attachmentID)
	client := &http.Client{}

	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return "", "", "", err
	}

	authHeader := fmt.Sprintf(`OAuth oauth_consumer_key="%s", oauth_token="%s"`, apiKey, token)
	req.Header.Set("Authorization", authHeader)

	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	localPath := filepath.Join(svc.CfgSvc.GetTrelloDownloadPath(), filename)
	out, err := os.Create(localPath)
	if err != nil {
		return "", "", "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", "", "", err
	}

	return localPath, attachmentID, extension, nil
}

func extractTrelloIDsAndExt(url string) (cardID, attachmentID, extension string, err error) {
	// Regex: capture cardID, attachmentID, and file extension
	re := regexp.MustCompile(`cards/([a-f0-9]+)/attachments/([a-f0-9]+)/download/[^/]+(\.[a-zA-Z0-9]+)$`)
	matches := re.FindStringSubmatch(url)

	if len(matches) != 4 {
		return "", "", "", fmt.Errorf("no match found")
	}

	return matches[1], matches[2], matches[3], nil
}

func fetchTrelloEntities[T any](url string) ([]T, error) {
	var props []T
	resp, err := http.Get(url)
	if err != nil {
		return props, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return props, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &props); err != nil {
		return props, err
	}

	return props, nil
}

func fetchCustomFieldDefs(baseURL, apiKey, token, boardID string) (map[string]trCustomFieldDef, error) {
	url := fmt.Sprintf("%s/boards/%s/customFields?key=%s&token=%s", baseURL, boardID, apiKey, token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var defs []trCustomFieldDef
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &defs)

	defMap := make(map[string]trCustomFieldDef)
	for _, def := range defs {
		defMap[def.ID] = def
	}

	return defMap, nil
}

func fetchCustomFields(baseURL, apiKey, token, cardID string) ([]trCustomFieldItem, error) {
	var fields []trCustomFieldItem
	url := fmt.Sprintf("%s/cards/%s/customFieldItems?key=%s&token=%s", baseURL, cardID, apiKey, token)
	resp, err := http.Get(url)
	if err != nil {
		return fields, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &fields)
	return fields, nil
}

func fetchAttachments(baseURL, apiKey, token, cardID string) ([]TRAttachment, error) {
	var attachments []TRAttachment
	url := fmt.Sprintf("%s/cards/%s/attachments?key=%s&token=%s", baseURL, cardID, apiKey, token)
	resp, err := http.Get(url)
	if err != nil {
		return attachments, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &attachments)
	return attachments, nil
}

func fetchComments(baseURL, apiKey, token, cardID string) ([]TRComment, error) {
	var comments []TRComment
	url := fmt.Sprintf("%s/cards/%s/actions?filter=commentCard&key=%s&token=%s", baseURL, cardID, apiKey, token)
	resp, err := http.Get(url)
	if err != nil {
		return comments, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &comments)
	return comments, nil
}
