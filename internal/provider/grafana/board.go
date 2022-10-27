package grafana

import (
	"bytes"
	"encoding/json"
)

type (
	// Board represents Grafana dashboard.
	Board struct {
		ID              uint     `json:"id,omitempty"`
		UID             string   `json:"uid,omitempty"`
		Slug            string   `json:"slug"`
		Title           string   `json:"title"`
		OriginalTitle   string   `json:"originalTitle"`
		Tags            []string `json:"tags"`
		Style           string   `json:"style"`
		Timezone        string   `json:"timezone"`
		Editable        bool     `json:"editable"`
		HideControls    bool     `json:"hideControls" graf:"hide-controls"`
		SharedCrosshair bool     `json:"sharedCrosshair" graf:"shared-crosshair"`
		Panels          []*Panel `json:"panels"`
		//		Rows            []*Row     `json:"rows"`
		Templating  Templating `json:"templating"`
		Annotations struct {
			List []Annotation `json:"list"`
		} `json:"annotations"`
		Refresh       *BoolString `json:"refresh,omitempty"`
		SchemaVersion uint        `json:"schemaVersion"`
		Version       uint        `json:"version"`
		Links         []Link      `json:"links"`
		Time          Time        `json:"time"`
		Timepicker    Timepicker  `json:"timepicker"`
		GraphTooltip  int         `json:"graphTooltip,omitempty"`
	}
	Time struct {
		From string `json:"from"`
		To   string `json:"to"`
	}
	Timepicker struct {
		Now              *bool    `json:"now,omitempty"`
		RefreshIntervals []string `json:"refresh_intervals"`
		TimeOptions      []string `json:"time_options"`
	}
	Templating struct {
		List []TemplateVar `json:"list"`
	}
	TemplateVar struct {
		Name        string      `json:"name"`
		Type        string      `json:"type"`
		Auto        bool        `json:"auto,omitempty"`
		AutoCount   *int        `json:"auto_count,omitempty"`
		Datasource  interface{} `json:"datasource"`
		Refresh     BoolInt     `json:"refresh"`
		Options     []Option    `json:"options"`
		IncludeAll  bool        `json:"includeAll"`
		AllFormat   string      `json:"allFormat"`
		AllValue    string      `json:"allValue"`
		Multi       bool        `json:"multi"`
		MultiFormat string      `json:"multiFormat"`
		Query       interface{} `json:"query"`
		Regex       string      `json:"regex"`
		Current     Current     `json:"current"`
		Label       string      `json:"label"`
		Hide        uint8       `json:"hide"`
		Sort        int         `json:"sort"`
	}
	// Option for templateVar
	Option struct {
		Text     string `json:"text"`
		Value    string `json:"value"`
		Selected bool   `json:"selected"`
	}
	// Current for templateVar
	Current struct {
		Tags  []*string          `json:"tags,omitempty"`
		Text  *StringSliceString `json:"text"`
		Value interface{}        `json:"value"` // TODO select more precise type
	}
	Annotation struct {
		Name        string      `json:"name"`
		Datasource  interface{} `json:"datasource"`
		ShowLine    bool        `json:"showLine"`
		IconColor   string      `json:"iconColor"`
		LineColor   string      `json:"lineColor"`
		IconSize    uint        `json:"iconSize"`
		Enable      bool        `json:"enable"`
		Query       string      `json:"query"`
		Expr        string      `json:"expr"`
		Step        string      `json:"step"`
		TextField   string      `json:"textField"`
		TextFormat  string      `json:"textFormat"`
		TitleFormat string      `json:"titleFormat"`
		TagsField   string      `json:"tagsField"`
		Tags        []string    `json:"tags"`
		TagKeys     string      `json:"tagKeys"`
		Type        string      `json:"type"`
	}
	// Link represents link to another dashboard or external weblink
	Link struct {
		Title       string   `json:"title"`
		Type        string   `json:"type"`
		AsDropdown  *bool    `json:"asDropdown,omitempty"`
		DashURI     *string  `json:"dashUri,omitempty"`
		Dashboard   *string  `json:"dashboard,omitempty"`
		Icon        *string  `json:"icon,omitempty"`
		IncludeVars bool     `json:"includeVars"`
		KeepTime    *bool    `json:"keepTime,omitempty"`
		Params      *string  `json:"params,omitempty"`
		Tags        []string `json:"tags,omitempty"`
		TargetBlank *bool    `json:"targetBlank,omitempty"`
		Tooltip     *string  `json:"tooltip,omitempty"`
		URL         *string  `json:"url,omitempty"`
	}
)

// Height of rows maybe passed as number (ex 200) or
// as string (ex "200px") or empty string
type Height string

func (h *Height) UnmarshalJSON(raw []byte) error {
	if raw == nil || bytes.Equal(raw, []byte(`"null"`)) {
		return nil
	}
	if raw[0] != '"' {
		tmp := []byte{'"'}
		raw = append(tmp, raw...)
		raw = append(raw, byte('"'))
	}
	var tmp string
	err := json.Unmarshal(raw, &tmp)
	*h = Height(tmp)
	return err
}
