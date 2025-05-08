package trello

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	props, err := fetchProperties(searchURL)
	if err != nil {
		return results, err
	}

	// lgr.Logger.Info("\n📌 Retrieved properties:\n", slog.Int("count", len(props)))

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
			field := TRPropField{}

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

func fetchProperties(url string) ([]TRProperty, error) {
	var props []TRProperty
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

func fetchCustomFieldDefs(baseURL, apiKey, token, boardID string) (map[string]trPropCustomFieldDef, error) {
	url := fmt.Sprintf("%s/boards/%s/customFields?key=%s&token=%s", baseURL, boardID, apiKey, token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var defs []trPropCustomFieldDef
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &defs)

	defMap := make(map[string]trPropCustomFieldDef)
	for _, def := range defs {
		defMap[def.ID] = def
	}

	return defMap, nil
}

func fetchCustomFields(baseURL, apiKey, token, cardID string) ([]trPropCustomFieldItem, error) {
	var fields []trPropCustomFieldItem
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

func fetchAttachments(baseURL, apiKey, token, cardID string) ([]TRPropAttachment, error) {
	var attachments []TRPropAttachment
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

func fetchComments(baseURL, apiKey, token, cardID string) ([]TRPropComment, error) {
	var comments []TRPropComment
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
