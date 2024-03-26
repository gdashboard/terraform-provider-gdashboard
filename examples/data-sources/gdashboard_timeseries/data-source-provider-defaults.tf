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
    }
  }
}

data "gdashboard_timeseries" "jvm_memory" {
  title = "JVM Memory"

  field {
    unit     = "bytes"
    decimals = 1
    min      = 0
    max      = 10000
  }

  queries {
    prometheus {
      uid           = "prometheus"
      expr          = "sum(increase(jvm_memory_total{container_name='container'}[$__rate_interval]))"
      instant       = false
      ref_id        = "Prometheus_Query"
      min_step      = "30"
      legend_format = "Memory total"
    }
  }
}

data "gdashboard_timeseries" "native_memory" {
  title = "Native Memory"

  field {
    unit     = "bytes"
    decimals = 1
    min      = 0
    max      = 30000
  }

  queries {
    prometheus {
      uid           = "prometheus"
      expr          = "sum(increase(native_memory_total{container_name='container'}[$__rate_interval]))"
      instant       = false
      ref_id        = "Prometheus_Query"
      min_step      = "30"
      legend_format = "Memory total"
    }
  }
}
