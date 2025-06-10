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
    in

    {
      packages = forAllSystems ({ pkgs }: {
        foo = pkgs.buildGoModule {
          inherit version;
          pname = "foo";
          env = {
            GOWORK = "off";
          };
          src = ./.;
          modRoot = "./foo";
          vendorHash = ""; # Bad hash
          meta = {
            inherit homepage;
            description = "${description} - foo";
          };
        };

        bar = pkgs.buildGoModule {
          inherit version;
          pname = "bar";
          env = {
            GOWORK = "off";
          };
          src = ./.;
          modRoot = "./bar";
          vendorHash = ""; # Bad hash
          meta = {
            inherit homepage;
            description = "${description} - bar";
          };
        };
      });

      devShells = forAllSystems ({ pkgs }: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            nixd
            nixpkgs-fmt

            bash-language-server
            shellcheck
            shfmt

            coreutils
            lowdown

            go
            gopls
            gotools
            go-tools
          ];
        };
      });
    };
}
