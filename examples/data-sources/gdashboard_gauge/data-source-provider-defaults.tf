provider "gdashboard" {
  defaults {
    gauge {
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

      graph {
        orientation            = "horizontal"
        show_threshold_labels  = true
        show_threshold_markers = true

        options {
          calculation = "lastNotNull"
        }
      }
    }
  }
}

data "gdashboard_gauge" "jvm_memory" {
  title = "JVM Memory"

  queries {
    prometheus {
      uid           = "prometheus"
      expr          = "sum(increase(jvm_memory_total{container_name='container'}[$__rate_interval]))"
      min_step      = "30"
      legend_format = "{{job_type}}"
      instant       = true
    }
  }
}

data "gdashboard_gauge" "native_memory" {
  title = "Native Memory"

  queries {
    prometheus {
      uid           = "prometheus"
      expr          = "sum(increase(native_total{container_name='container'}[$__rate_interval]))"
      min_step      = "30"
      legend_format = "{{job_type}}"
      instant       = true
    }
  }
}
