package provider

import (
	"regexp"
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
				Config:      testAccLogsDataSourceCloudWatchConflictConfig,
				ExpectError: regexp.MustCompile("Attribute \"queries\\[0]\\.cloudwatch\\[0]\\.metrics\" cannot be specified when"),
			},
			{
				Config: testAccLogsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_logs.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_logs.test", "json", testAccLogsDataSourceConfigExpectedJson),
				),
			},
			{
				Config: testAccLogsDataSourceEmptyConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_logs.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_logs.test", "json", testAccLogsDataSourceEmptyConfigExpectedJson),
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
    cloudwatch {
      logs {
        uid        = "cloudwatch"
        expression = "fields @timestamp, @message |\n sort @timestamp desc |\n limit 20"
        region     = "eu-west-2"

        log_group {
          arn  = "arn:aws:logs:eu-west-2:123456789012:log-group:/ecs/production/service-1:*"
          name = "/ecs/production/service-1"
        }

        log_group {
          arn  = "arn:aws:logs:eu-west-2:123456789012:log-group:/ecs/production/service-2:*"
          name = "/ecs/production/service-2"
        }
      }
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
        "uid": "cloudwatch",
        "name": "",
        "type": "cloudwatch",
        "typeLogoUrl": "",
        "access": "",
        "url": "",
        "isDefault": false,
        "jsonData": null,
        "secureJsonData": null
      },
      "queryMode": "Logs",
      "metricQueryType": 0,
      "metricEditorMode": 0,
      "expression": "fields @timestamp, @message |\n sort @timestamp desc |\n limit 20",
      "logGroups": [
        {
          "arn": "arn:aws:logs:eu-west-2:123456789012:log-group:/ecs/production/service-1:*",
          "name": "/ecs/production/service-1"
        },
        {
          "arn": "arn:aws:logs:eu-west-2:123456789012:log-group:/ecs/production/service-2:*",
          "name": "/ecs/production/service-2"
        }
      ],
      "region": "eu-west-2"
    }
  ]
}`

const testAccLogsDataSourceEmptyConfig = `
data "gdashboard_logs" "test" {
  title = "Test"
}
`

const testAccLogsDataSourceEmptyConfigExpectedJson = `{
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

const testAccLogsDataSourceCloudWatchConflictConfig = `
data "gdashboard_logs" "test" {
  title = "Test"

  queries {
    cloudwatch {
      logs {
        uid        = "cloudwatch"
        expression = ""
      }

      metrics {
        uid = "cloudwatch"
        statistic = "Sum"
        metric_name = "Name"
        namespace = "Namespace"
      }
    }
  }
}
`
