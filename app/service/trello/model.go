package trello

import "time"

type trPropCustomFieldItem struct {
	IDCustomField string                 `json:"idCustomField"`
	Value         map[string]interface{} `json:"value"`
	IDValue       string                 `json:"idValue"` // for list type
}

type trPropCustomFieldDef struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Options []struct {
		ID    string `json:"id"`
		Value struct {
			Text string `json:"text"`
		} `json:"value"`
	} `json:"options"`
}

type TRPropLabel struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type TRPropField struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type TRPropAttachment struct {
	ID   string    `json:"id"`
	Name string    `json:"name"`
	URL  string    `json:"url"`
	Date time.Time `json:"date"`
}

type TRPropComment struct {
	Data struct {
		Text string `json:"text"`
	} `json:"data"`
}

type TRProperty struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	LocationAR       string             `json:"locationAR"`
	LocationEN       string             `json:"locationEN"`
	Lot              string             `json:"lot"`
	Type             string             `json:"type"`
	Status           string             `json:"status"`
	Owner            string             `json:"owner"`
	Area             float64            `json:"area"`
	Shares           float64            `json:"shares"`
	Organized        bool               `json:"organized"`
	Effects          bool               `json:"effects"`
	Labels           []TRPropLabel      `json:"labels"`
	Fields           []TRPropField      `json:"fields"`
	Attachments      []TRPropAttachment `json:"attachments"`
	Comments         []TRPropComment    `json:"comments"`
	DateLastActivity time.Time          `json:"dateLastActivity"`
	UpdatedAt        time.Time          `json:"updatedAt"`
}
