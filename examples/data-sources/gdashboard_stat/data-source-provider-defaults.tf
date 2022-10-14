provider "gdashboard" {
  defaults {
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

        text_size {
          title = 10
          value = 15
        }
      }

      field {
        unit = "p"
      }
    }
  }
}

data "gdashboard_stat" "status_1" {
  title = "Container 1 Status"

  targets {
    prometheus {
      uid     = "prometheus"
      expr    = "up{container_name='container_1'}"
      instant = true
    }
  }
}

data "gdashboard_stat" "status_2" {
  title = "Container 2 Status"

  targets {
    prometheus {
      uid     = "prometheus"
      expr    = "up{container_name='container_2'}"
      instant = true
    }
  }
}
