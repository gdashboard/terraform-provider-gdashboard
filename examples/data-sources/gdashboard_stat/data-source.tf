data "gdashboard_stat" "status" {
  title = "Status"

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

  targets {
    prometheus {
      uid     = "prometheus"
      expr    = "up{container_name='container'}"
      instant = true
    }
  }
}