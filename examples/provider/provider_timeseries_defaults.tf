provider "gdashboard" {
  defaults {
    timeseries {
      legend {
        calculations = ["min", "max", "mean"]
        display_mode = "table"
        placement    = "bottom"
      }
    }
  }
}

# The panel has a configured legend: from defaults - 3 calculations, table display mode, and placement on the bottom
data "gdashboard_timeseries" "defaults" {
  title = "Test"
}

# The panel has a configured legend: from defaults - 3 calculations, table display mode; override - a placement on the right
data "gdashboard_timeseries" "override_defaults" {
  title = "Test"

  legend {
    placement = "right" # override one field
  }
}
