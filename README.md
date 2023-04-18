# Terraform Provider GDashboard

The provider offers a handy syntax to define Grafana dashboards: time series, gauge, bar gauge, stat, etc.

## Using the provider

Please, see [provider documentation](https://registry.terraform.io/providers/gdashboard/gdashboard/latest/docs).

The provider defines only data sources. Each data source computes a JSON that is compatible with Grafana API.
Therefore, this provider is not particularly useful on its own, but it can be
used to generate a JSON compatible with Grafana API, which can then be used
with [Grafana provider](https://registry.terraform.io/providers/grafana/grafana/latest/docs) to provision a dashboard.

## Examples

```terraform
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

data "gdashboard_dashboard" "dashboard" {
  title = "My dashboard"

  layout {
    section {
      title = "Basic Details"
              
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

# Define provider
provider "grafana" {
  url  = "https://my.grafana.com" # use your API endpoint
  auth = var.grafana_auth
}

# Create your dashboard
resource "grafana_dashboard" "my_dashboard" {
  config_json = data.gdashboard_dashboard.dashboard.json
}
```

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.18

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up-to-date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Local Testing

Build the provider locally:
```shell
$ go build -o terraform-provider-gdashboard
```

Add a path to the **folder** with the binary to the `dev.tfrc`:
```hocon
provider_installation {
  dev_overrides {
    "gdashboard/gdashboard" = "../terraform-provider-gdashboard"
  }
  direct {}
}
```

Configure terraform cli:
```shell
$ export TF_CLI_CONFIG_FILE=./dev.tfrc
```

Run changes:
```shell
$ terraform plan
```
