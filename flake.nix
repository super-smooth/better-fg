{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.systems.url = "github:nix-systems/default";
  inputs.flake-utils = {
    url = "github:numtide/flake-utils";
    inputs.systems.follows = "systems";
  };
  inputs.gomod2nix = {
    url = "github:nix-community/gomod2nix";
    inputs.nixpkgs.follows = "nixpkgs";
    inputs.flake-utils.follows = "flake-utils";
  };

  outputs =
    {
      nixpkgs,
      flake-utils,
      gomod2nix,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          system = system;
          config.allowUnfree = true;
        };
        gomod2nixPkgs = gomod2nix.legacyPackages.${system};
      in
      rec {
        packages.default = gomod2nixPkgs.buildGoApplication {
          pname = "better-fg";
          version = "0.1.0";

          src = ./.;

          # Use gomod2nix generated lock file
          modules = ./gomod2nix.toml;

          subPackages = [ "cmd/better-fg" ];

          meta = with pkgs.lib; {
            description = "A CLI tool for fuzzy searching background jobs";
            homepage = "https://github.com/super-smooth/better-fg";
            license = licenses.mit;
            maintainers = [ ];
            mainProgram = "better-fg";
          };
        };

        devShells.default = pkgs.mkShell {
          shellHook = ''
            export CGO_ENABLED="1"

            # Auto-install lefthook if in a git repo and lefthook.yml exists (skip in CI)
            if [ -z "$CI" ] && [ -d .git ] && [ -f lefthook.yml ]; then
              if ! lefthook version &> /dev/null; then
                echo "⚠️  lefthook not found in PATH"
              elif [ ! -f .git/hooks/lefthook ] && [ ! -f .git/hooks/pre-commit ]; then
                echo "🔧 Installing lefthook hooks..."
                lefthook install
              fi
            fi

            # Check if gomod2nix.toml needs to be generated
            if [ -f go.mod ] && [ ! -f gomod2nix.toml ]; then
              echo "📦 Generating gomod2nix.toml (first time setup)..."
              gomod2nix generate
              echo "✅ gomod2nix.toml generated!"
            fi
          '';

          packages = [
            # Go
            pkgs.go
            pkgs.gopls
            pkgs.golangci-lint
            pkgs.gofumpt

            # Build tools
            pkgs.git

            # Nix tools
            pkgs.nix

            # System libraries
            pkgs.gcc

            # Additional tools
            pkgs.lld

            # gomod2nix for generating lock file
            gomod2nixPkgs.gomod2nix

            # lefthook for git hooks
            pkgs.lefthook
            # commitlint
            pkgs.commitlint
          ];
        };
      }
    );
}
