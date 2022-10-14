data "gdashboard_stat" "test" {
  title       = "Test"
  description = "Stat description"

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

    mappings {
      value {
        value        = "1"
        display_text = "UP"
        color        = "green"
      }

      value {
        value        = "0"
        display_text = "DOWN"
        color        = "red"
      }

      special {
        match        = "null+nan"
        display_text = "DOWN"
        color        = "red"
      }
    }
  }

  targets {
    prometheus {
      uid     = "prometheus"
      expr    = "up{container_name='container'}"
      instant = true
    }
  }
}