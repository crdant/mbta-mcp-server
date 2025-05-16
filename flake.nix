{
  description = "MBTA MCP Server";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    devshell.url = "github:numtide/devshell";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, devshell, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ devshell.overlays.default ];
        };

        # Build semver-cli directly
        semver-cli = pkgs.buildGoModule {
          pname = "semver-cli";
          version = "1.0.0";
          src = pkgs.fetchFromGitHub {
            owner = "maykonlf";
            repo = "semver-cli";
            rev = "v1.0.0";
            sha256 = "0gwgwc1m432ri3mzap7bx4ddgcml1g3a4xaqbjn2xnxcdyiz8555";
          };
          vendorHash = "sha256-y3BX6bGOLvxTSaQcLyo8PJ4L/U/+GsbONSv29xO9Mj8=";
          subPackages = [ "cmd/semver-cli" ];
        };
      in
      {
        devShells.default = pkgs.devshell.mkShell {
          imports = [ (pkgs.devshell.importTOML ./devshell.toml) ];
          packages = with pkgs; [
            go
            golangci-lint
            goreleaser
            apko
            melange
            crane
            cosign
            gnumake
            semver-cli
          ];
          env = [
            {
              name = "GOFLAGS";
              value = "-mod=vendor";
            }
            {
              name = "GOPRIVATE";
              value = "github.com/crdant/*";
            }
          ];
        };

        packages.default = pkgs.buildGoModule {
          pname = "mbta-mcp-server";
          version = "dev";
          src = ./.;
          vendorHash = null; # Will be computed on first build
          subPackages = [ "cmd/server" ];
          ldflags = [ "-s" "-w" "-X main.Version=dev" ];
          meta = with pkgs.lib; {
            description = "An MCP server that communicates with the MBTA API to provide Boston-area transit information";
            homepage = "https://github.com/crdant/mbta-mcp-server";
            license = licenses.mit;
            maintainers = [ "crdant" ];
          };
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/mbta-mcp-server";
        };
      }
    );
}