package grafana

import (
	"bytes"
	"encoding/json"
)

type (
	// Board represents Grafana dashboard.
	Board struct {
		ID          uint     `json:"id,omitempty"`
		UID         string   `json:"uid,omitempty"`
		Title       string   `json:"title"`
		Description string   `json:"description,omitempty"`
		Tags        []string `json:"tags,omitempty"`
		Style       string   `json:"style"`
		Timezone    string   `json:"timezone"`
		WeekStart   string   `json:"weekStart,omitempty"`
		LiveNow     bool     `json:"liveNow"`
		Editable    bool     `json:"editable"`
		Panels      []Panel  `json:"panels"`
		//		Rows            []*Row     `json:"rows"`
		Templating    Templating  `json:"templating"`
		Annotations   Annotations `json:"annotations"`
		Refresh       *BoolString `json:"refresh,omitempty"`
		SchemaVersion uint        `json:"schemaVersion"`
		Version       int64       `json:"version"`
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
		Hidden           *bool    `json:"hidden,omitempty"`
		NowDelay         string   `json:"now_delay,omitempty"`
		RefreshIntervals []string `json:"refresh_intervals"`
		TimeOptions      []string `json:"time_options"`
	}
	Annotations struct {
		List []Annotation `json:"list"`
	}
	Templating struct {
		List []TemplateVar `json:"list"`
	}
	TemplateVarDataSource struct {
		UID  string `json:"uid"`
		Type string `json:"type"`
	}
	TemplateVarQueryPrometheus struct {
		Query string `json:"query"`
		RefID string `json:"refId"`
	}
	TemplateVarAdHocFilter struct {
		Condition string `json:"condition"`
		Key       string `json:"key"`
		Operator  string `json:"operator"`
		Value     string `json:"value"`
	}
	TemplateVar struct {
		Type        string                   `json:"type"`                  // constant, custom, textbox, adhoc, datasource, query, interval
		Name        string                   `json:"name"`                  // constant, custom, textbox, adhoc, datasource, query, interval
		Description string                   `json:"description,omitempty"` // constant, custom, textbox, adhoc, datasource, query, interval
		Label       string                   `json:"label"`                 // constant, custom, textbox, adhoc, datasource, query, interval
		Hide        uint8                    `json:"hide"`                  // ________, custom, textbox, adhoc, datasource, query, interval
		Auto        bool                     `json:"auto,omitempty"`        // ________, ______, _______, _____, __________, _____, interval
		AutoCount   *int64                   `json:"auto_count,omitempty"`  // ________, ______, _______, _____, __________, _____, interval
		AutoMin     *string                  `json:"auto_min,omitempty"`    // ________, ______, _______, _____, __________, _____, interval
		Datasource  *TemplateVarDataSource   `json:"datasource,omitempty"`  // ________, ______, _______, adhoc, __________, _____, ________
		Filters     []TemplateVarAdHocFilter `json:"filters,omitempty"`     // ________, ______, _______, adhoc, __________, _____, ________
		Refresh     BoolInt                  `json:"refresh"`               // ________, ______, _______, _____, datasource, query, ________
		Options     []Option                 `json:"options"`               // ________, custom, _______, _____, __________, _____, interval
		IncludeAll  bool                     `json:"includeAll"`            // ________, custom, _______, _____, datasource, query, ________
		AllValue    string                   `json:"allValue"`              // ________, custom, _______, _____, datasource, query, ________
		Multi       bool                     `json:"multi"`                 // ________, custom, _______, _____, datasource, query, ________
		Query       interface{}              `json:"query"`                 // constant, custom, textbox, _____, datasource, query, ________
		Regex       string                   `json:"regex"`                 // ________, ______, _______, _____, datasource, query, interval
		Current     Option                   `json:"current"`               // ________, custom, _______, _____, __________, _____, interval
		Sort        int                      `json:"sort"`                  // ________, ______, _______, _____, __________, query, ________
		Definition  string                   `json:"definition,omitempty"`  // ________, ______, _______, _____, __________, query, ________
	}
	// Option for templateVar
	Option struct {
		Text     *string `json:"text"`
		Value    string  `json:"value"`
		Selected bool    `json:"selected"`
	}
	Annotation struct {
		Name       string               `json:"name"`
		Datasource AnnotationDataSource `json:"datasource"`
		IconColor  string               `json:"iconColor"`
		Enable     bool                 `json:"enable,omitempty"`
		Hide       *bool                `json:"hide,omitempty"`
		Type       *string              `json:"type,omitempty"`
		// Grafana
		Target *AnnotationGrafanaTarget `json:"target,omitempty"`
		// Prometheus
		Expr            *string `json:"expr,omitempty"`
		Step            *string `json:"step,omitempty"`
		UseValueForTime *bool   `json:"useValueForTime,omitempty"`
		TitleFormat     *string `json:"titleFormat,omitempty"`
		TextFormat      *string `json:"textFormat,omitempty"`
		TagKeys         *string `json:"tagKeys,omitempty"`
	}
	AnnotationDataSource struct {
		UID  string `json:"uid"`
		Type string `json:"type"`
	}
	AnnotationGrafanaTarget struct {
		Limit    int64    `json:"limit"`
		MatchAny bool     `json:"matchAny"`
		Tags     []string `json:"tags"`
		Type     string   `json:"type"`
	}
	// Link represents link to another dashboard or external weblink
	Link struct {
		Title       string   `json:"title"`
		Type        string   `json:"type"`
		AsDropdown  *bool    `json:"asDropdown,omitempty"`
		Icon        *string  `json:"icon,omitempty"`
		IncludeVars bool     `json:"includeVars"`
		KeepTime    *bool    `json:"keepTime,omitempty"`
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
