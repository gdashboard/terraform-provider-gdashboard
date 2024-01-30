{
  description = "Grafana Dashboard Provider";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/release-23.11";
    hestia.url = "github:iRevive/hestia-nix";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, hestia, ... }:
    flake-utils.lib.simpleFlake { 
      inherit self nixpkgs;
      name = "gdashboard";
      overlay = hestia.overlays.default;
      config = { allowUnfree = true; };
      shell = ./shell-impl.nix;
    };
}
