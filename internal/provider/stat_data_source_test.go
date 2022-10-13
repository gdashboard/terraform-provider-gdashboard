package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStatDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccStatDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_stat.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_stat.test", "json", testAccStatDataSourceConfigExpectedJson),
				),
			},
			{
				Config: testAccStatDataSourceProviderCustomDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_stat.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_stat.test", "json", testAccStatDataSourceProviderCustomDefaultsConfigExpectedJson),
				),
			},
			{
				Config: testAccStatDataSourceProviderDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_stat.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_stat.test", "json", testAccStatDataSourceProviderDefaultsConfigExpectedJson),
				),
			},
		},
	})
}

const testAccStatDataSourceConfig = `
data "gdashboard_stat" "test" {
  title       = "Test"
  description = "Stat description"

  graph {
    orientation = "vertical"
    text_mode   = "value"
    color_mode  = "background"
    graph_mode  = "none"

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
	unit = "p"

    mappings {
      value {
        value        = "1"
        display_text = "UP"
        color        = "green"
      }

      value {
        value        = "0"
        display_text = "DOWN"
        color        = "red"
      }

      special {
        match        = "null+nan"
        display_text = "DOWN"
        color        = "red"
      }
    }
  }

  targets {
    prometheus {
      uid     = "prometheus"
      expr    = "up{container_name='container'}"
      instant = true
    }
  }
	
}
`

const testAccStatDataSourceConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "description": "Stat description",
  "transparent": false,
  "type": "stat",
  "colors": null,
  "colorValue": false,
  "colorBackground": false,
  "decimals": 0,
  "format": "",
  "gauge": {
    "maxValue": 0,
    "minValue": 0,
    "show": false,
    "thresholdLabels": false,
    "thresholdMarkers": false
  },
  "nullPointMode": "",
  "sparkline": {},
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
      "expr": "up{container_name='container'}",
      "instant": true
    }
  ],
  "thresholds": "",
  "valueFontSize": "",
  "valueMaps": null,
  "valueName": "",
  "options": {
    "orientation": "vertical",
    "textMode": "value",
    "colorMode": "background",
    "graphMode": "none",
    "justifyMode": "",
    "displayMode": "",
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
  "fieldConfig": {
    "defaults": {
      "unit": "p",
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
      },
      "mappings": [
        {
          "type": "value",
          "options": {
            "0": {
              "color": "red",
              "text": "DOWN",
              "index": 1
            },
            "1": {
              "color": "green",
              "text": "UP",
              "index": 0
            }
          }
        },
        {
          "type": "special",
          "options": {
            "match": "null+nan",
            "result": {
              "color": "red",
              "text": "DOWN",
              "index": 2
            }
          }
        }
      ]
    }
  }
}`

const testAccStatDataSourceProviderCustomDefaultsConfig = `
provider "gdashboard" {
  defaults {
    stat {
	  graph {
        orientation = "vertical"
        text_mode   = "value"
        color_mode  = "background"
        graph_mode  = "none"

		text_size {
	  	  title = 10
	      value = 15
	    }

        options {
          values      = true
          fields 	  = "/.*/"
		  calculation = "first"
        }
      }
	}
  }
}

data "gdashboard_stat" "test" {
  title = "Test"
}
`

const testAccStatDataSourceProviderCustomDefaultsConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "transparent": false,
  "type": "stat",
  "colors": null,
  "colorValue": false,
  "colorBackground": false,
  "decimals": 0,
  "format": "",
  "gauge": {
    "maxValue": 0,
    "minValue": 0,
    "show": false,
    "thresholdLabels": false,
    "thresholdMarkers": false
  },
  "nullPointMode": "",
  "sparkline": {},
  "thresholds": "",
  "valueFontSize": "",
  "valueMaps": null,
  "valueName": "",
  "options": {
    "orientation": "vertical",
    "textMode": "value",
    "colorMode": "background",
    "graphMode": "none",
    "justifyMode": "",
    "displayMode": "",
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

const testAccStatDataSourceProviderDefaultsConfig = `
data "gdashboard_stat" "test" {
  title = "Test"
}
`

const testAccStatDataSourceProviderDefaultsConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "transparent": false,
  "type": "stat",
  "colors": null,
  "colorValue": false,
  "colorBackground": false,
  "decimals": 0,
  "format": "",
  "gauge": {
    "maxValue": 0,
    "minValue": 0,
    "show": false,
    "thresholdLabels": false,
    "thresholdMarkers": false
  },
  "nullPointMode": "",
  "sparkline": {},
  "thresholds": "",
  "valueFontSize": "",
  "valueMaps": null,
  "valueName": "",
  "options": {
    "orientation": "auto",
    "textMode": "auto",
    "colorMode": "value",
    "graphMode": "area",
    "justifyMode": "",
    "displayMode": "",
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
