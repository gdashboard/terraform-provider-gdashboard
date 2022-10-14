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
			{
				Config: testAccGaugeDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_gauge.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_gauge.test", "json", testAccGaugeDataSourceConfigExpectedJson),
				),
			},
			{
				Config: testAccGaugeDataSourceProviderCustomDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_gauge.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_gauge.test", "json", testAccGaugeDataSourceProviderCustomDefaultsConfigExpectedJson),
				),
			},
			{
				Config: testAccGaugeDataSourceProviderDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_gauge.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_gauge.test", "json", testAccGaugeDataSourceProviderDefaultsConfigExpectedJson),
				),
			},
		},
	})
}

const testAccGaugeDataSourceConfig = `
data "gdashboard_gauge" "test" {
  title       = "Test"
  description = "Gauge description"

  graph {
    orientation 		   = "vertical"
	show_threshold_labels  = true
    show_threshold_markers = false

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

  field {
    unit = "percent"
	min  = 0
	max  = 100

    color {
      mode = "thresholds"
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

const testAccGaugeDataSourceConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "description": "Gauge description",
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
    "showThresholdMarkers": false,
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
  ],
  "fieldConfig": {
    "defaults": {
      "unit": "percent",
      "min": 0,
      "max": 100,
      "color": {
        "mode": "thresholds",
        "fixedColor": "green",
        "seriesBy": "last"
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

		text_size {
		  title = 10
		  value = 15
		}
      }
	}
  }
}

data "gdashboard_gauge" "test" {
  title = "Test"
}
`

const testAccGaugeDataSourceProviderCustomDefaultsConfigExpectedJson = `{
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
    "showThresholdMarkers": false,
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

const testAccGaugeDataSourceProviderDefaultsConfigExpectedJson = `{
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
    "showThresholdMarkers": true,
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
