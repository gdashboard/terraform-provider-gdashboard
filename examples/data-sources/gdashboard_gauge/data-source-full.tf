data "gdashboard_gauge" "jvm_memory" {
  title = "JVM Memory"

  field {
    unit = "percent"

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
    orientation            = "horizontal"
    show_threshold_labels  = true
    show_threshold_markers = true

    options {
      calculation = "lastNotNull"
    }
  }

  queries {
    prometheus {
      uid           = "prometheus"
      expr          = "sum(increase(jvm_memory_total{container_name='container'}[$__rate_interval]))"
      min_step      = "30"
      legend_format = "{{job_type}}"
      ref_id        = "Prometheus_Query"
      instant       = true
    }
  }
}
