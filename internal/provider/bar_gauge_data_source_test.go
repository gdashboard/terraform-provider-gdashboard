package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBarGaugeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccBarGaugeDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_bar_gauge.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_bar_gauge.test", "json", testAccBarGaugeDataSourceConfigExpectedJson),
				),
			},
			{
				Config: testAccBarGaugeDataSourceProviderCustomDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_bar_gauge.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_bar_gauge.test", "json", testAccBarGaugeDataSourceProviderCustomDefaultsConfigExpectedJson),
				),
			},
			{
				Config: testAccBarGaugeDataSourceProviderDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_bar_gauge.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_bar_gauge.test", "json", testAccBarGaugeDataSourceProviderDefaultsConfigExpectedJson),
				),
			},
		},
	})
}

const testAccBarGaugeDataSourceConfig = `
data "gdashboard_bar_gauge" "test" {
  title       = "Test"
  description = "Bar gauge description"

  field {
    decimals = 0
  }

  graph {
    orientation  = "horizontal"
    display_mode = "basic"

    options {
      values      = true
      fields 	  = "/.*/"
	  calculation = "first"
    }

	text_size {
	  title = 10
	  value = 15
	}
  }

  targets {
    prometheus {
      uid           = "prometheus"
      expr          = "sort_desc(sum(increase(data[$__range])) by (job_type))"
      min_interval  = "30"
      legend_format = "{{job_type}}"
      instant       = true
	  ref_id		= "Prometheus_Query"
    }
  }
	
}
`

const testAccBarGaugeDataSourceConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "description": "Bar gauge description",
  "transparent": false,
  "type": "bargauge",
  "options": {
    "orientation": "horizontal",
    "textMode": "",
    "colorMode": "",
    "graphMode": "",
    "justifyMode": "",
    "displayMode": "basic",
    "content": "",
    "mode": "",
    "text": {
      "titleSize": 10,
      "valueSize": 15
    },
    "reduceOptions": {
      "values": true,
      "fields": "/.*/",
      "calcs": [
        "first"
      ]
    }
  },
  "targets": [
    {
      "refId": "Prometheus_Query",
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
      "expr": "sort_desc(sum(increase(data[$__range])) by (job_type))",
      "interval": "30",
      "legendFormat": "{{job_type}}",
      "instant": true
    }
  ],
  "fieldConfig": {
    "defaults": {
      "unit": "",
      "decimals": 0,
      "color": {
        "mode": "palette-classic",
        "fixedColor": "green",
        "seriesBy": "last"
      },
      "thresholds": {
        "mode": "absolute",
        "steps": [
          {
            "color": "green",
            "value": null
          }
        ]
      },
      "custom": {
        "axisPlacement": "",
        "barAlignment": 0,
        "drawStyle": "",
        "fillOpacity": 0,
        "gradientMode": "",
        "lineInterpolation": "",
        "lineWidth": 0,
        "pointSize": 0,
        "showPoints": "",
        "spanNulls": false,
        "hideFrom": {
          "legend": false,
          "tooltip": false,
          "viz": false
        },
        "lineStyle": {
          "fill": ""
        },
        "scaleDistribution": {
          "type": ""
        },
        "stacking": {
          "group": "",
          "mode": ""
        },
        "thresholdsStyle": {
          "mode": ""
        }
      }
    }
  }
}`

const testAccBarGaugeDataSourceProviderCustomDefaultsConfig = `
provider "gdashboard" {
  defaults {
    bar_gauge {
      field {
        unit     = "bytes"
        decimals = 1
        min      = 0
        max      = 10000
        
        color { 
          mode        = "palette-classic"
          fixed_color = "red"
          series_by   = "first"
        }
      }

	  graph {
        orientation  = "vertical"
        display_mode = "lcd"

        options {
          values      = true
          fields 	  = "/.*/"
		  calculation = "first"
        }
      }
	}
  }
}

data "gdashboard_bar_gauge" "test" {
  title = "Test"
}
`

const testAccBarGaugeDataSourceProviderCustomDefaultsConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "transparent": false,
  "type": "bargauge",
  "options": {
    "orientation": "vertical",
    "textMode": "",
    "colorMode": "",
    "graphMode": "",
    "justifyMode": "",
    "displayMode": "lcd",
    "content": "",
    "mode": "",
    "text": {},
    "reduceOptions": {
      "values": true,
      "fields": "/.*/",
      "calcs": [
        "first"
      ]
    }
  },
  "fieldConfig": {
    "defaults": {
      "unit": "bytes",
      "decimals": 1,
      "min": 0,
      "max": 10000,
      "color": {
        "mode": "palette-classic",
        "fixedColor": "red",
        "seriesBy": "first"
      },
      "thresholds": {
        "mode": "absolute",
        "steps": [
          {
            "color": "green",
            "value": null
          }
        ]
      },
      "custom": {
        "axisPlacement": "",
        "barAlignment": 0,
        "drawStyle": "",
        "fillOpacity": 0,
        "gradientMode": "",
        "lineInterpolation": "",
        "lineWidth": 0,
        "pointSize": 0,
        "showPoints": "",
        "spanNulls": false,
        "hideFrom": {
          "legend": false,
          "tooltip": false,
          "viz": false
        },
        "lineStyle": {
          "fill": ""
        },
        "scaleDistribution": {
          "type": ""
        },
        "stacking": {
          "group": "",
          "mode": ""
        },
        "thresholdsStyle": {
          "mode": ""
        }
      }
    }
  }
}`

const testAccBarGaugeDataSourceProviderDefaultsConfig = `
data "gdashboard_bar_gauge" "test" {
  title = "Test"
}
`

const testAccBarGaugeDataSourceProviderDefaultsConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "transparent": false,
  "type": "bargauge",
  "options": {
    "orientation": "auto",
    "textMode": "",
    "colorMode": "",
    "graphMode": "",
    "justifyMode": "",
    "displayMode": "gradient",
    "content": "",
    "mode": "",
    "text": {},
    "reduceOptions": {
      "values": false,
      "fields": "",
      "calcs": [
        "lastNotNull"
      ]
    }
  },
  "fieldConfig": {
    "defaults": {
      "unit": "",
      "color": {
        "mode": "palette-classic",
        "fixedColor": "green",
        "seriesBy": "last"
      },
      "thresholds": {
        "mode": "absolute",
        "steps": [
          {
            "color": "green",
            "value": null
          }
        ]
      },
      "custom": {
        "axisPlacement": "",
        "barAlignment": 0,
        "drawStyle": "",
        "fillOpacity": 0,
        "gradientMode": "",
        "lineInterpolation": "",
        "lineWidth": 0,
        "pointSize": 0,
        "showPoints": "",
        "spanNulls": false,
        "hideFrom": {
          "legend": false,
          "tooltip": false,
          "viz": false
        },
        "lineStyle": {
          "fill": ""
        },
        "scaleDistribution": {
          "type": ""
        },
        "stacking": {
          "group": "",
          "mode": ""
        },
        "thresholdsStyle": {
          "mode": ""
        }
      }
    }
  }
}`
