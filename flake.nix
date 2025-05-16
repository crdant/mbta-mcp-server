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
        # Define an overlay that adds our packages to nixpkgs
        overlay = final: prev: {
          semver-cli = final.buildGoModule {
            pname = "semver-cli";
            version = "1.1.1";
            src = final.fetchFromGitHub {
              owner = "maykonlsf";
              repo = "semver-cli";
              rev = "v1.1.1";
              sha256 = "sha256-Qj9RV2wW0i0hL5CDL4WCa7yKUIvmL2kkId1K8qNIfsw=";
            };
            vendorHash = "sha256-o87+Y0m2pmjijlEXQDFFekshCgWR+lYt83nU4w5faV0=";
            subPackages = [ "cmd/semver" ];
          };
        };

        # Import nixpkgs with our overlays
        pkgs = import nixpkgs {
          inherit system;
          overlays = [
            devshell.overlays.default
            overlay
          ];
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
            semver-cli  # Add the semver-cli package we defined
            # Container runtimes (users can choose their preference)
            docker
            nerdctl
            podman
            openssl
            coreutils   # Provides mkdir, rm and other basic tools
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