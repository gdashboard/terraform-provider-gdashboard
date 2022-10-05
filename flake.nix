{
  description = "Grafana Dashboard Provider";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/22.05";
    nixpkgs-unstable.url = "github:nixos/nixpkgs";
    hestia.url = "github:iRevive/hestia-nix";
    flake-utils.url = "github:numtide/flake-utils";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, nixpkgs-unstable, flake-utils, hestia, ... }:
    let
      pkgs-unstable = import nixpkgs-unstable {
        localSystem = "aarch64-darwin";
      };

      terraform-overrides = final: prev: {
        go_1_19 = pkgs-unstable.go_1_19;
      };
    in
    flake-utils.lib.simpleFlake { # todo use forEachSystem to keep flake crossplatform
      inherit self nixpkgs;
      name = "gdashboard";
      overlay = hestia.overlays.default;
      preOverlays = [terraform-overrides];
      shell = ./shell-impl.nix;
    };
}
