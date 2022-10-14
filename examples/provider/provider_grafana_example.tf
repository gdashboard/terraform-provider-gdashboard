# Define Grafana provider
provider "grafana" {
  url  = "https://grafana.example.com/"
  auth = var.grafana_auth
}

# Create your dashboard
resource "grafana_dashboard" "my_dashboard" {
  config_json = data.gdashboard_dashboard.dashboard.json
}

data "gdashboard_dashboard" "dashboard" {
  title = "My dashboard"

  layout {
    row {
      panel {
        size = {
          height = 8
          width  = 10
        }
        source = data.gdashboard_stat.status.json
      }
    }
  }
}

data "gdashboard_stat" "status" {
  title       = "Status"
  description = "Shows the status of the container"

  field {
    mappings {
      value {
        value        = "1"
        display_text = "UP"
        color        = "green"
      }

      special {
        match        = "null+nan"
        display_text = "DOWN"
        color        = "red"
      }
    }
  }

  queries {
    prometheus {
      uid     = "prometheus"
      expr    = "up{container_name='container'}"
      instant = true
    }
  }
}
