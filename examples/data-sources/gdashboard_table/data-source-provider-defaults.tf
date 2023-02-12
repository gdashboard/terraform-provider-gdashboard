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
  title = "Example"

  graph {
    graph {
      show_header = false

      column {
        align      = "right"
        filterable = true
        width      = 30
        min_width  = 50
      }

      cell {
        inspectable = true
      }
    }
  }

  queries {
    prometheus {
      uid     = "prometheus"
      expr    = "up{container_name='container_1'}"
      instant = true
    }
  }
}
