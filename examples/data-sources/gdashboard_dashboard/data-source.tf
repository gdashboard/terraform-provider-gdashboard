data "gdashboard_timeseries" "jvm_memory" {
  title = "JVM Memory"

  targets {
    prometheus {
      uid  = "prometheus"
      expr = "sum(increase(jvm_memory_total{container_name='container'}[$__rate_interval]))"
    }
  }
}

data "gdashboard_dashboard" "jvm_dashboard" {
  title         = "JVM Dashboard"
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

    custom { # dropdown variable
      name = "custom"

      option {
        text  = "entry-1"
        value = "value"
      }

      option {
        text     = "entry-2"
        value    = "value"
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
        source = data.gdashboard_timeseries.jvm_memory.json
      }
    }
  }
}
