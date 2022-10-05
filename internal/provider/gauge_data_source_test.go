package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGaugeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			/*{
				Config: testAccGaugeDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_gauge.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_gauge.test", "json", gaugeProviderCustomSeriesExpectedJson),
				),
			},*/
			{
				Config: testAccGaugeDataSourceProviderCustomDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_gauge.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_gauge.test", "json", gaugeProviderCustomDefaultsExpectedJson),
				),
			},
			{
				Config: testAccGaugeDataSourceProviderDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_gauge.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_gauge.test", "json", gaugeProviderDefaultsExpectedJson),
				),
			},
		},
	})
}

const testAccGaugeDataSourceConfig = `
data "gdashboard_gauge" "test" {
  title = "Test"

  legend {
	calculations = ["min", "max", "mean"]
	display_mode = "table"
    placement    = "bottom"
  }

  tooltip {
	mode = "multi"
  }

  field {
    unit     = "bytes"
    decimals = 1
    min      = 0
    max      = 10000
    
    color { 
      mode        = "palette-classic"
      fixed_color = "green"
      series_by   = "last"
    }
  }
}
`

const gaugeProviderCustomSeriesExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0
}`

const testAccGaugeDataSourceProviderCustomDefaultsConfig = `
provider "gdashboard" {
  defaults {
    gauge {
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

		thresholds {
		  mode = "percentage"
			
 		  step {
			color = "green"
		  }

		  step {
			color = "orange"
			value = 65
		  }

		  step {
			color = "red"
			value = 90
		  }
		}
      }

	  graph {
        orientation            = "vertical"
        show_threshold_labels  = true
        show_threshold_markers = false

        options {
          values      = true
          fields 	  = "/.*/"
		  calculation = "first"
        }
      }
	}
  }
}

data "gdashboard_gauge" "test" {
  title = "Test"
}
`

const gaugeProviderCustomDefaultsExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "transparent": false,
  "type": "gauge",
  "options": {
    "orientation": "vertical",
    "textMode": "",
    "colorMode": "",
    "graphMode": "",
    "justifyMode": "",
    "displayMode": "",
    "content": "",
    "mode": "",
    "showThresholdLabels": true,
    "ShowThresholdMarkers": false,
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
        "mode": "percentage",
        "steps": [
          {
            "color": "green",
            "value": null
          },
          {
            "color": "orange",
            "value": 65
          },
          {
            "color": "red",
            "value": 90
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

const testAccGaugeDataSourceProviderDefaultsConfig = `
data "gdashboard_gauge" "test" {
  title = "Test"
}
`

const gaugeProviderDefaultsExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "transparent": false,
  "type": "gauge",
  "options": {
    "orientation": "auto",
    "textMode": "",
    "colorMode": "",
    "graphMode": "",
    "justifyMode": "",
    "displayMode": "",
    "content": "",
    "mode": "",
    "showThresholdLabels": false,
    "ShowThresholdMarkers": true,
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
