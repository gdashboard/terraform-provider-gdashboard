package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDashboardDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccDashboardDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "json", testAccDashboardDataSourceConfigExpectedJson),
				),
			},
			{
				Config: testAccDashboardDataSourceProviderDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "json", testAccDashboardDataSourceProviderDefaultsConfigExpectedJson),
				),
			},
			{
				Config: testAccDashboardDataSourceProviderCustomDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "json", testAccDashboardDataSourceProviderCustomDefaultsConfigExpectedJson),
				),
			},
		},
	})
}

const testAccDashboardDataSourceConfig = `
data "gdashboard_dashboard" "test" {
  title         = "Test"
  uid 	        = "test-uid"
  editable      = false
  style         = "light"
  graph_tooltip = "shared-crosshair"
  
  time {
    from = "now-1h"
    to   = "now+1h"
  }

  variables {
	const {
	  name  = "var"
	  value = "const-value"
	}

    custom {
	  name = "custom"
	  hide = "label" 

      option {
		text  = "entry-1"
		value = "value"
	  }

	  option {
		text  	 = "entry-2"
		value 	 = "value"
		selected = true
	  }
    }
  }

  layout {
	row {
	  panel {
	  	size = {
		  height = 8
		  width  = 10
		}
	  	source = "{\"title\": \"Panel 1\"}"
	  }

	  panel {
	  	size = {
		  height = 3
		  width  = 24
		}
	  	source = "{\"title\": \"Panel 2\"}"
	  } 
	}

	row {
	  panel {
	  	size = {
		  height = 4
		  width  = 24
		}
	  	source = "{\"title\": \"Panel 3\"}"
	  }

	  panel {
	  	size = {
		  height = 3
		  width  = 3
		}
	  	source = "{\"title\": \"Panel 4\"}"
	  } 
	}
  }
}
`

const testAccDashboardDataSourceConfigExpectedJson = `{
  "uid": "test-uid",
  "slug": "",
  "title": "Test",
  "originalTitle": "",
  "tags": null,
  "style": "light",
  "timezone": "",
  "editable": false,
  "hideControls": false,
  "sharedCrosshair": false,
  "panels": [
    {
      "editable": false,
      "error": false,
      "gridPos": {
        "h": 8,
        "w": 10,
        "x": 0,
        "y": 0
      },
      "id": 0,
      "isNew": false,
      "span": 0,
      "title": "Panel 1",
      "transparent": false,
      "type": "",
      "title": "Panel 1"
    },
    {
      "editable": false,
      "error": false,
      "gridPos": {
        "h": 3,
        "w": 24,
        "x": 10,
        "y": 0
      },
      "id": 0,
      "isNew": false,
      "span": 0,
      "title": "Panel 2",
      "transparent": false,
      "type": "",
      "title": "Panel 2"
    },
    {
      "editable": false,
      "error": false,
      "gridPos": {
        "h": 4,
        "w": 24,
        "x": 0,
        "y": 9
      },
      "id": 0,
      "isNew": false,
      "span": 0,
      "title": "Panel 3",
      "transparent": false,
      "type": "",
      "title": "Panel 3"
    },
    {
      "editable": false,
      "error": false,
      "gridPos": {
        "h": 3,
        "w": 3,
        "x": 24,
        "y": 9
      },
      "id": 0,
      "isNew": false,
      "span": 0,
      "title": "Panel 4",
      "transparent": false,
      "type": "",
      "title": "Panel 4"
    }
  ],
  "templating": {
    "list": [
      {
        "name": "custom",
        "type": "custom",
        "datasource": null,
        "refresh": false,
        "options": [
          {
            "text": "entry-1",
            "value": "value",
            "selected": false
          },
          {
            "text": "entry-2",
            "value": "value",
            "selected": true
          }
        ],
        "includeAll": false,
        "allFormat": "",
        "allValue": "",
        "multi": false,
        "multiFormat": "",
        "query": "entry-1 : value, entry-2 : value",
        "regex": "",
        "current": {
          "text": [
            "entry-2"
          ],
          "value": "value"
        },
        "label": "",
        "hide": 1,
        "sort": 0
      },
      {
        "name": "var",
        "type": "constant",
        "datasource": null,
        "refresh": false,
        "options": null,
        "includeAll": false,
        "allFormat": "",
        "allValue": "",
        "multi": false,
        "multiFormat": "",
        "query": "const-value",
        "regex": "",
        "current": {
          "text": null,
          "value": null
        },
        "label": "",
        "hide": 0,
        "sort": 0
      }
    ]
  },
  "annotations": {
    "list": null
  },
  "schemaVersion": 0,
  "version": 1,
  "links": null,
  "time": {
    "from": "now-1h",
    "to": "now+1h"
  },
  "timepicker": {
    "refresh_intervals": null,
    "time_options": null
  },
  "graphTooltip": 1
}`

const testAccDashboardDataSourceProviderCustomDefaultsConfig = `
provider "gdashboard" {
  defaults {
    dashboard {
	  editable		= false
	  graph_tooltip = "shared-tooltip"
 	  style 		= "light"
      time {
        from = "now-12h"
		to   = "now-3h"
      }
    }
  }
}

data "gdashboard_dashboard" "test" {
  title = "Test"

  layout {

  }
}
`

const testAccDashboardDataSourceProviderCustomDefaultsConfigExpectedJson = `{
  "slug": "",
  "title": "Test",
  "originalTitle": "",
  "tags": null,
  "style": "light",
  "timezone": "",
  "editable": false,
  "hideControls": false,
  "sharedCrosshair": false,
  "panels": [],
  "templating": {
    "list": []
  },
  "annotations": {
    "list": null
  },
  "schemaVersion": 0,
  "version": 1,
  "links": null,
  "time": {
    "from": "now-12h",
    "to": "now-3h"
  },
  "timepicker": {
    "refresh_intervals": null,
    "time_options": null
  },
  "graphTooltip": 2
}`

const testAccDashboardDataSourceProviderDefaultsConfig = `
data "gdashboard_dashboard" "test" {
  title = "Test"

  layout {

  }
}
`

const testAccDashboardDataSourceProviderDefaultsConfigExpectedJson = `{
  "slug": "",
  "title": "Test",
  "originalTitle": "",
  "tags": null,
  "style": "dark",
  "timezone": "",
  "editable": true,
  "hideControls": false,
  "sharedCrosshair": false,
  "panels": [],
  "templating": {
    "list": []
  },
  "annotations": {
    "list": null
  },
  "schemaVersion": 0,
  "version": 1,
  "links": null,
  "time": {
    "from": "now-6h",
    "to": "now"
  },
  "timepicker": {
    "refresh_intervals": null,
    "time_options": null
  }
}`
