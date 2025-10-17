{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.systems.url = "github:nix-systems/default";
  inputs.flake-utils = {
    url = "github:numtide/flake-utils";
    inputs.systems.follows = "systems";
  };

  outputs =
    {
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          system = system;
          config.allowUnfree = true;
        };
      in
      rec {
        packages.default = pkgs.buildGoModule {
          pname = "better-fg";
          version = "0.1.0";

          src = ./.;

          vendorHash = "sha256-7BIdujoAlK10/9sBuBL1UNvrSkFos5VvxHVGbOiNFfo=";

          subPackages = [ "cmd/better-fg" ];

          meta = with pkgs.lib; {
            description = "A CLI tool for fuzzy searching background jobs";
            homepage = "https://github.com/super-smooth/better-fg";
            license = licenses.mit; # Update with actual license
            maintainers = [ ];
            mainProgram = "better-fg";
          };
        };

        devShells.default = pkgs.mkShell {
          shellHook = ''
            export CGO_ENABLED="1"
          '';

          packages = [
            # go
            pkgs.go
            pkgs.gopls
            pkgs.golangci-lint
            pkgs.gofumpt

            # Build tools
            pkgs.git

            # Node.js for commitlint and git hooks
            pkgs.nodejs
            pkgs.nodePackages.npm
            pkgs.commitlint

            # System libraries
            pkgs.gcc

            # Additional tools that might be useful
            pkgs.lld

          ];

        };
      }
    );
}
