rec {
  description = "Fact checker LINE Bot";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      homepage = "https://github.com/kaogeek/line-fact-check";
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
      version = builtins.substring 0 8 lastModifiedDate;

      # The set of systems to provide outputs for
      allSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];

      # A function that provides a system-specific Nixpkgs for the desired systems
      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });

      goEnvs = {
        GOTOOLCHAIN = "auto";
        GOWORK = "off";
      };
    in

    {
      packages = forAllSystems ({ pkgs }: {
        foo = pkgs.buildGoModule {
          inherit version;
          pname = "foo";
          env = goEnvs;
          src = ./.;
          modRoot = "./foo";
          vendorHash = "sha256-7zDTgVJ2yu6lkf6hwNdpAnC+VLEmL6iJGTKBOzPtlYM=";
          meta = {
            inherit homepage;
            description = "${description} - foo";
          };
        };

        factcheck = pkgs.buildGoModule {
          inherit version;
          pname = "factcheck";
          env = goEnvs;
          src = ./.;
          modRoot = "./factcheck";
          vendorHash = "sha256-ptfGpM0aHeH+vOLdea1a7UI6LfOQY9NY+48Ce3L9Q0Y=";
          meta = {
            inherit homepage;
            description = "${description} - factcheck";
          };
        };
      });

      devShells = forAllSystems ({ pkgs }: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            coreutils

            # Basic LSPs
            nixd
            nixpkgs-fmt
            bash-language-server
            shellcheck
            shfmt
            lowdown

            # Development - server
            go
            gopls
            gotools
            go-tools
            golangci-lint
            sqlc
          ];
        };
      });
    };
}
