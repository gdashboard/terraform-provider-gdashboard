provider "gdashboard" {
  defaults {

    dashboard {
      editable      = false
      graph_tooltip = "shared-tooltip"
      style         = "light"
      default_time_range {
        from = "now-12h"
        to   = "now-3h"
      }
    }

    stat {
      graph {
        orientation = "vertical"
        text_mode   = "value"
        color_mode  = "background"
        graph_mode  = "none"

        options {
          values      = true
          fields      = "/.*/"
          calculation = "first"
        }
      }
    }

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
          fields      = "/.*/"
          calculation = "first"
        }
      }
    }

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
        draw_style         = "bars"
        line_interpolation = "smooth"
        line_width         = 10
        fill_opacity       = 30
        gradient_mode      = "hue"
        line_style         = "dash"
        span_nulls         = true
        show_points        = "never"
        point_size         = 22
        stack_series       = "percent"
      }
    }
  }
}
