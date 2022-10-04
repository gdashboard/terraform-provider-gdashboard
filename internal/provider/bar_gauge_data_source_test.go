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
			/*{
				Config: testAccBarGaugeDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_bar_gauge.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_bar_gauge.test", "json", barGaugeProviderCustomSeriesExpectedJson),
				),
			},*/
			{
				Config: testAccBarGaugeDataSourceProviderCustomDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_bar_gauge.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_bar_gauge.test", "json", barGaugeProviderCustomDefaultsExpectedJson),
				),
			},
			{
				Config: testAccBarGaugeDataSourceProviderDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_bar_gauge.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_bar_gauge.test", "json", barGaugeProviderDefaultsExpectedJson),
				),
			},
		},
	})
}

const testAccBarGaugeDataSourceConfig = `
data "gdashboard_bar_gauge" "test" {
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

  graph {
    fill_opacity = 10
    show_points  = "always"
    span_nulls   = true
  }

  targets {
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

	  ref_id        = "CW_Query"
	  period        = "30"
      legend_format = "Request Count"
	}
  }
	
}
`

const barGaugeProviderCustomSeriesExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
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
      "legendFormat": "Request Count",
      "namespace": "AWS/ApplicationELB",
      "metricName": "HTTPCode_Target_2XX_Count",
      "statistics": [
        "Sum"
      ],
      "dimensions": {
        "LoadBalancer": "lb_arn_suffix",
        "TargetGroup": "target_group"
      },
      "period": "30",
      "region": "af-south-1"
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
        "mode": "",
        "steps": null
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
          "type": ""
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

const barGaugeProviderCustomDefaultsExpectedJson = `{
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
        "mode": "",
        "steps": null
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

const barGaugeProviderDefaultsExpectedJson = `{
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
        "mode": "",
        "steps": null
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
