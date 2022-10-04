data "gdashboard_timeseries" "jvm_memory" {
  title = "Test"

  targets {
    prometheus {
      uid  = "prometheus"
      expr = "sum(increase(jvm_memory_total{container_name='container'}[$__rate_interval]))"
    }
  }
}

data "gdashboard_dashboard" "test" {
  title         = "Test"
  uid           = "test-uid"
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
        source = data.gdashboard_timeseries.jvm_memory
      }
    }
  }
}
