package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLogsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccLogsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_logs.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_logs.test", "json", testAccLogsDataSourceConfigExpectedJson),
				),
			},
			{
				Config: testAccLogsDataSourceProviderDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_logs.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_logs.test", "json", testAccLogsDataSourceProviderDefaultsConfigExpectedJson),
				),
			},
		},
	})
}

const testAccLogsDataSourceConfig = `
data "gdashboard_logs" "test" {
  title       = "Test"
  description = "Logs description"

  graph {
    show_time          = true
	show_unique_labels = true
    show_common_labels = true
    wrap_lines         = true
    prettify_json      = true
    enable_log_details = false
    deduplication      = "exact"
    order              = "oldest_first"
  }

  queries {
    prometheus {
      uid          = "prometheus"
      expr         = "sum (jvm_memory_bytes_used{container_name='container', area='heap'}) / sum (jvm_memory_bytes_max{container_name='container', area='heap'}) * 100"
      min_interval = "30"
      instant      = true
    }
  }
}
`

const testAccLogsDataSourceConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "description": "Logs description",
  "transparent": false,
  "type": "logs",
  "options": {
    "showTime": true,
    "showLabels": true,
    "showCommonLabels": true,
    "wrapLogMessage": true,
    "prettifyLogMessage": true,
    "enableLogDetails": false,
    "dedupStrategy": "exact",
    "sortOrder": "Ascending"
  },
  "targets": [
    {
      "refId": "",
      "datasource": {
        "id": 0,
        "orgId": 0,
        "uid": "prometheus",
        "name": "",
        "type": "prometheus",
        "typeLogoUrl": "",
        "access": "",
        "url": "",
        "isDefault": false,
        "jsonData": null,
        "secureJsonData": null
      },
      "expr": "sum (jvm_memory_bytes_used{container_name='container', area='heap'}) / sum (jvm_memory_bytes_max{container_name='container', area='heap'}) * 100",
      "interval": "30",
      "instant": true
    }
  ]
}`

const testAccLogsDataSourceProviderDefaultsConfig = `
data "gdashboard_logs" "test" {
  title = "Test"
}
`

const testAccLogsDataSourceProviderDefaultsConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "transparent": false,
  "type": "logs",
  "options": {
    "showTime": false,
    "showLabels": false,
    "showCommonLabels": false,
    "wrapLogMessage": false,
    "prettifyLogMessage": false,
    "enableLogDetails": true,
    "dedupStrategy": "none",
    "sortOrder": "Descending"
  }
}`
