{ pkgs }:

let
  hestia = pkgs.hestia;
  colored = hestia.ansi.colored;
in
hestia.shell.mkShell {
  name = "grafana-dashboard-builder";

  shellScripts = [];

  packages = [
    pkgs.go_1_19
    pkgs.terraform
  ];
}
