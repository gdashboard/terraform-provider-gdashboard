package grafana

type Datasource struct {
	ID                uint        `json:"id"`
	OrgID             uint        `json:"orgId"`
	UID               string      `json:"uid"`
	Name              string      `json:"name"`
	Type              string      `json:"type"`
	TypeLogoURL       string      `json:"typeLogoUrl"`
	Access            string      `json:"access"` // direct or proxy
	URL               string      `json:"url"`
	Password          *string     `json:"password,omitempty"`
	User              *string     `json:"user,omitempty"`
	Database          *string     `json:"database,omitempty"`
	BasicAuth         *bool       `json:"basicAuth,omitempty"`
	ReadOnly          *bool       `json:"readOnly,omitempty"`
	BasicAuthUser     *string     `json:"basicAuthUser,omitempty"`
	BasicAuthPassword *string     `json:"basicAuthPassword,omitempty"`
	IsDefault         bool        `json:"isDefault"`
	JSONData          interface{} `json:"jsonData"`
	SecureJSONData    interface{} `json:"secureJsonData"`
}

// DatasourceType Datasource type as described in
// http://docs.grafana.org/reference/http_api/#available-data-source-types
type DatasourceType struct {
	Metrics  bool   `json:"metrics"`
	Module   string `json:"module"`
	Name     string `json:"name"`
	Partials struct {
		Query string `json:"query"`
	} `json:"datasource"`
	PluginType  string `json:"pluginType"`
	ServiceName string `json:"serviceName"`
	Type        string `json:"type"`
}
