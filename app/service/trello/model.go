package trello

import "time"

type trCustomFieldItem struct {
	IDCustomField string                 `json:"idCustomField"`
	Value         map[string]interface{} `json:"value"`
	IDValue       string                 `json:"idValue"` // for list type
}

type trCustomFieldDef struct {
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

type TRLabel struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type TRField struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type TRAttachment struct {
	ID   string    `json:"id"`
	Name string    `json:"name"`
	URL  string    `json:"url"`
	Date time.Time `json:"date"`
}

type TRComment struct {
	Data struct {
		Text string `json:"text"`
	} `json:"data"`
}

type TRProperty struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	LocationAR       string         `json:"locationAR"`
	LocationEN       string         `json:"locationEN"`
	Lot              string         `json:"lot"`
	Type             string         `json:"type"`
	Status           string         `json:"status"`
	Owner            string         `json:"owner"`
	Area             float64        `json:"area"`
	Shares           float64        `json:"shares"`
	Organized        bool           `json:"organized"`
	Effects          bool           `json:"effects"`
	Labels           []TRLabel      `json:"labels"`
	Fields           []TRField      `json:"fields"`
	Attachments      []TRAttachment `json:"attachments"`
	Comments         []TRComment    `json:"comments"`
	DateLastActivity time.Time      `json:"dateLastActivity"`
	UpdatedAt        time.Time      `json:"updatedAt"`
}

type TRInheritanceConfinement struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Title            string         `json:"title"`
	Generation       int64          `json:"generation"`
	Labels           []TRLabel      `json:"labels"`
	Fields           []TRField      `json:"fields"`
	Attachments      []TRAttachment `json:"attachments"`
	Comments         []TRComment    `json:"comments"`
	DateLastActivity time.Time      `json:"dateLastActivity"`
	UpdatedAt        time.Time      `json:"updatedAt"`
}

type TRSupportiveDoc struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Title            string         `json:"title"`
	Category         string         `json:"category"`
	Labels           []TRLabel      `json:"labels"`
	Fields           []TRField      `json:"fields"`
	Attachments      []TRAttachment `json:"attachments"`
	Comments         []TRComment    `json:"comments"`
	DateLastActivity time.Time      `json:"dateLastActivity"`
	UpdatedAt        time.Time      `json:"updatedAt"`
}
