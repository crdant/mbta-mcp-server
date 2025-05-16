{
  description = "MBTA MCP Server";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            golangci-lint
            goreleaser
            docker
            docker-compose
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