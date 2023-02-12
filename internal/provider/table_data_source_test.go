package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTableDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccTableDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_table.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_table.test", "json", testAccTableDataSourceConfigExpectedJson),
				),
			},
			{
				Config: testAccTableDataSourceProviderCustomDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_table.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_table.test", "json", testAccTableDataSourceProviderCustomDefaultsConfigExpectedJson),
				),
			},
			{
				Config: testAccTableDataSourceProviderDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_table.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_table.test", "json", testAccTableDataSourceProviderDefaultsConfigExpectedJson),
				),
			},
		},
	})
}

const testAccTableDataSourceConfig = `
data "gdashboard_table" "test" {
  title       = "Test"
  description = "Table description"

  graph {
	show_header = false

	column {
	  align 	 = "right"
	  filterable = true
	  width 	 = 30
      min_width  = 50
	}

	cell {
	  inspectable  = true
	  display_mode = "basic" 
	}

	footer {
	  pagination   = true
	  fields 	   = ["a", "b"]
      calculations = ["min", "max"]
	}
  }

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

  queries {
    prometheus {
      uid     = "prometheus"
      expr    = "up{container_name='container'}"
      instant = true
    }
  }
	
}
`

const testAccTableDataSourceConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "description": "Table description",
  "transparent": false,
  "type": "table",
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
  "options": {
    "showHeader": false,
    "footer": {
      "show": true,
      "enablePagination": true,
      "fields": [
        "a",
        "b"
      ],
      "reducer": [
        "a",
        "b"
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
        },
        "align": "right",
        "displayMode": "basic",
        "inspect": true,
        "filterable": true,
        "width": 30,
        "minWidth": 50
      }
    }
  }
}`

const testAccTableDataSourceProviderCustomDefaultsConfig = `
provider "gdashboard" {
  defaults {
    table {
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
	}
  }
}

data "gdashboard_table" "test" {
  title = "Test"
}
`

const testAccTableDataSourceProviderCustomDefaultsConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "transparent": false,
  "type": "table",
  "options": {
    "showHeader": true,
    "footer": {
      "show": false,
      "enablePagination": false
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

const testAccTableDataSourceProviderDefaultsConfig = `
data "gdashboard_table" "test" {
  title = "Test"
}
`

const testAccTableDataSourceProviderDefaultsConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "transparent": false,
  "type": "table",
  "options": {
    "showHeader": true,
    "footer": {
      "show": false,
      "enablePagination": false
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
