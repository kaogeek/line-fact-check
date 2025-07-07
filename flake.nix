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
          vendorHash = "sha256-1+dsaRdpIh9lNHkKQa7FflzeveXv10JaxXr4VRpPil8=";
          meta = {
            inherit homepage;
            description = "${description} - factcheck";
          };
        };

        # Docker images
        docker-factcheck = pkgs.dockerTools.buildImage {
          name = "factcheck";
          tag = version;
          copyToRoot = [ pkgs.bash pkgs.coreutils ];
          config = {
            Entrypoint = [ "${self.packages.${pkgs.system}.factcheck}/bin/api" ];
            ExposedPorts = {
              "8080/tcp" = {};
            };
          };
        };

        docker-postgres-integration-test = pkgs.dockerTools.buildImage {
          name = "postgres-integration-test";
          tag = version;
          copyToRoot = [ pkgs.bash pkgs.coreutils pkgs.postgresql_16 ];
          config = {
            Env = [
              "POSTGRES_USER=postgres"
              "POSTGRES_PASSWORD=postgres"
              "POSTGRES_DB=factcheck"
              "POSTGRES_HOST_AUTH_METHOD=trust"
              "PGDATA=/var/lib/postgresql/data"
            ];
            ExposedPorts = {
              "5432/tcp" = {};
            };
            Entrypoint = [ "${pkgs.bash}/bin/bash" ];
            Cmd = [ "-c" ''
              # Create postgres user (UID 999, same as official postgres image)
              groupadd -g 999 postgres
              useradd -u 999 -g postgres -s /bin/bash -m postgres
              
              # Create data directory with proper ownership
              mkdir -p /var/lib/postgresql/data
              chown -R postgres:postgres /var/lib/postgresql/data
              
              # Initialize database if not already done
              if [ ! -f /var/lib/postgresql/data/PG_VERSION ]; then
                echo "Initializing PostgreSQL database..."
                su postgres -c "${pkgs.postgresql_16}/bin/initdb -D /var/lib/postgresql/data -U postgres"
              fi
              
              # Start PostgreSQL as postgres user
              echo "Starting PostgreSQL..."
              exec su postgres -c "${pkgs.postgresql_16}/bin/postgres -D /var/lib/postgresql/data -c listen_addresses='*'"
            '' ];
            Volumes = {
              "/var/lib/postgresql/data" = {};
            };
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
            wire
          ];
        };
      });
    };
}
