data "gdashboard_stat" "status" {
  title       = "Status"
  description = "Stat description"

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

  queries {
    prometheus {
      uid     = "prometheus"
      expr    = "up{container_name='container'}"
      ref_id  = "Prometheus_Query"
      instant = true
    }
  }
}
