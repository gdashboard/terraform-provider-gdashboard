data "gdashboard_timeseries" "jvm_memory" {
  title = "JVM memory"

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

  overrides {
    by_query_id {
      query_id = "Prometheus_Query"
      field {
        color {
          mode        = "fixed"
          fixed_color = "red"
        }
      }
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
      ref_id        = "Prometheus_Query"
      min_interval  = "30"
      legend_format = "Memory total"
    }
  }

}
