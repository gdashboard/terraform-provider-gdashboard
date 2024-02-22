package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTimeseriesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccTimeseriesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_timeseries.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_timeseries.test", "json", testAccTimeseriesDataSourceConfigExpectedJson),
				),
			},
			{
				Config: testAccTimeseriesDataSourceProviderCustomDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_timeseries.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_timeseries.test", "json", testAccTimeseriesDataSourceProviderCustomDefaultsConfigExpectedJson),
				),
			},
			{
				Config: testAccTimeseriesDataSourceProviderDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_timeseries.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_timeseries.test", "json", testAccTimeseriesDataSourceProviderDefaultsConfigExpectedJson),
				),
			},
		},
	})
}

const testAccTimeseriesDataSourceConfig = `
data "gdashboard_timeseries" "test" {
  title       = "Test"
  description = "Timeseries description"

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

	thresholds {
      show_as = "line"
    }
  }

  graph {
    fill_opacity = 10
    show_points  = "always"
    span_nulls   = true
  }

  queries {
    prometheus {
      uid           = "prometheus"
      expr          = "sum(increase(jvm_memory_total{container_name='container'}[$__rate_interval]))"
      instant       = false
	  ref_id		= "Prometheus_Query"
      min_interval  = "30"
      legend_format = "Memory total"
    }

	cloudwatch {
	  uid         = "cloudwatch"
	  namespace   = "AWS/ApplicationELB"
	  metric_name = "HTTPCode_Target_2XX_Count"
      statistic   = "Sum"
      match_exact = true
	  region      = "af-south-1"
	  
	  dimension {
	    name = "LoadBalancer"
	    value = "lb_arn_suffix"
	  }
	  
	  dimension {
	    name = "TargetGroup"
	    value = "target_group"
	  }	

	  ref_id = "CW_Query"
	  period = "30"
      label  = "Request Count"
	}
  }
	
}
`

const testAccTimeseriesDataSourceConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "description": "Timeseries description",
  "transparent": false,
  "type": "timeseries",
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
      "expr": "sum(increase(jvm_memory_total{container_name='container'}[$__rate_interval]))",
      "interval": "30",
      "legendFormat": "Memory total"
    },
    {
      "refId": "CW_Query",
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
      "queryMode": "Metrics",
      "metricQueryType": 0,
      "metricEditorMode": 0,
      "namespace": "AWS/ApplicationELB",
      "metricName": "HTTPCode_Target_2XX_Count",
      "statistic": "Sum",
      "dimensions": {
        "LoadBalancer": "lb_arn_suffix",
        "TargetGroup": "target_group"
      },
      "matchExact": true,
      "period": "30",
      "region": "af-south-1",
      "label": "Request Count"
    }
  ],
  "options": {
    "legend": {
      "calcs": [
        "min",
        "max",
        "mean"
      ],
      "displayMode": "table",
      "placement": "bottom"
    },
    "tooltip": {
      "mode": "multi"
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
        "axisPlacement": "auto",
        "barAlignment": 0,
        "drawStyle": "line",
        "fillOpacity": 10,
        "gradientMode": "none",
        "lineInterpolation": "linear",
        "lineWidth": 1,
        "pointSize": 5,
        "showPoints": "always",
        "spanNulls": true,
        "hideFrom": {
          "legend": false,
          "tooltip": false,
          "viz": false
        },
        "lineStyle": {
          "fill": "solid"
        },
        "scaleDistribution": {
          "type": "linear"
        },
        "stacking": {
          "group": "",
          "mode": "none"
        },
        "thresholdsStyle": {
          "mode": "line"
        }
      }
    }
  }
}`

const testAccTimeseriesDataSourceProviderCustomDefaultsConfig = `
provider "gdashboard" {
  defaults {
    timeseries {
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
          fixed_color = "red"
          series_by   = "first"
        }
      }

	  graph {
        draw_style		   = "bars"
	    line_interpolation = "smooth"
        line_width         = 10
        fill_opacity       = 30
        gradient_mode      = "hue"
        line_style         = "dash"
        span_nulls         = true
        show_points  	   = "never"
        point_size   	   = 22
        stack_series   	   = "percent"
      }
	}
  }
}

data "gdashboard_timeseries" "test" {
  title = "Test"
}
`

const testAccTimeseriesDataSourceProviderCustomDefaultsConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "transparent": false,
  "type": "timeseries",
  "options": {
    "legend": {
      "calcs": [
        "min",
        "max",
        "mean"
      ],
      "displayMode": "table",
      "placement": "bottom"
    },
    "tooltip": {
      "mode": "multi"
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
        "axisPlacement": "auto",
        "barAlignment": 0,
        "drawStyle": "bars",
        "fillOpacity": 30,
        "gradientMode": "hue",
        "lineInterpolation": "smooth",
        "lineWidth": 10,
        "pointSize": 22,
        "showPoints": "never",
        "spanNulls": true,
        "hideFrom": {
          "legend": false,
          "tooltip": false,
          "viz": false
        },
        "lineStyle": {
          "fill": "dash"
        },
        "scaleDistribution": {
          "type": "linear"
        },
        "stacking": {
          "group": "",
          "mode": "percent"
        },
        "thresholdsStyle": {
          "mode": ""
        }
      }
    }
  }
}`

const testAccTimeseriesDataSourceProviderDefaultsConfig = `
data "gdashboard_timeseries" "test" {
  title = "Test"
}
`

const testAccTimeseriesDataSourceProviderDefaultsConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "transparent": false,
  "type": "timeseries",
  "options": {
    "legend": {
      "calcs": null,
      "displayMode": "list",
      "placement": "bottom"
    },
    "tooltip": {
      "mode": "single"
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
        "axisPlacement": "auto",
        "barAlignment": 0,
        "drawStyle": "line",
        "fillOpacity": 0,
        "gradientMode": "none",
        "lineInterpolation": "linear",
        "lineWidth": 1,
        "pointSize": 5,
        "showPoints": "auto",
        "spanNulls": false,
        "hideFrom": {
          "legend": false,
          "tooltip": false,
          "viz": false
        },
        "lineStyle": {
          "fill": "solid"
        },
        "scaleDistribution": {
          "type": "linear"
        },
        "stacking": {
          "group": "",
          "mode": "none"
        },
        "thresholdsStyle": {
          "mode": ""
        }
      }
    }
  }
}`
