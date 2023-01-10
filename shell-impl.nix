{ pkgs }:

let
  hestia = pkgs.hestia;
  colored = hestia.ansi.colored;
in
hestia.shell.mkShell {
  name = "gdashboard";

  shellScripts = [];

  packages = [
    pkgs.go
    pkgs.go-tools
    pkgs.terraform
  ];
}
