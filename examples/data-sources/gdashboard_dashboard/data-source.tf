data "gdashboard_timeseries" "jvm_memory" {
  title = "JVM Memory"

  queries {
    prometheus {
      uid  = "prometheus"
      expr = "sum(increase(jvm_memory_total{container_name='container'}[$__rate_interval]))"
    }
  }
}

data "gdashboard_timeseries" "http_requests" {
  title = "HTTP Requests"

  queries {
    prometheus {
      uid  = "prometheus"
      expr = "sum(increase(http_request_total{container_name='container'}[$__rate_interval]))"
    }
  }
}

data "gdashboard_timeseries" "http_status" {
  title = "HTTP Status"

  queries {
    prometheus {
      uid  = "prometheus"
      expr = "sum(increase(http_status_total{container_name='container'}[$__rate_interval]))"
    }
  }
}

data "gdashboard_dashboard" "jvm_dashboard" {
  title         = "JVM Dashboard"
  description   = "JVM details"
  style         = "light"
  graph_tooltip = "shared-crosshair"

  time {
    default_range {
      from = "now-1h"
      to   = "now+1h"
    }
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
    section {
      title = "JVM"

      panel {
        size = {
          height = 8
          width  = 10
        }
        source = data.gdashboard_timeseries.jvm_memory.json
      }
    }

    section { // each panel is on a new row
      title = "HTTP"

      row { // force a new row/line
        panel {
          size = {
            height = 8
            width  = 10
          }
          source = data.gdashboard_timeseries.http_requests.json
        }
      }

      row { // force a new row/line
        panel {
          size = {
            height = 8
            width  = 10
          }
          source = data.gdashboard_timeseries.http_status.json
        }
      }
    }
  }
}
