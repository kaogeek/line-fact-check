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
          
          buildInputs = with pkgs; [
            wire
            git
          ];
          
          # Check that wire-generated code is pristine
          preBuild = ''
            echo "Checking wire code generation..."
            cd cmd/api/di
            wire
            cd ../../../
            
            # Check if any files were modified
            if [ -n "$(git status --porcelain)" ]; then
              echo "Error: Wire-generated code differs from repository state"
              git diff
              exit 1
            else
              echo "Wire-generated code is pristine"
            fi
          '';
          
          meta = {
            inherit homepage;
            description = "${description} - factcheck";
          };
        };
      });

      # To build and load the image:
      # nix build .#dockerImages.factcheck && docker load < result
      dockerImages = forAllSystems ({ pkgs }: {
        factcheck = pkgs.dockerTools.buildImage {
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

            # Database for integration tests
            postgresql_16
          ];
        };

        # Integration test environment with Postgres setup
        test-integration = pkgs.mkShell {
          packages = with pkgs; [
            # Go tools
            go
            gopls
            gotools
            go-tools
            golangci-lint
            sqlc
            wire

            # Database
            postgresql_16

            # Utilities
            coreutils
            bash
            lsof
          ];

          # Set up Postgres environment and startup
          shellHook = ''
            # Set up Postgres environment
            export PGDATA="$PWD/.postgres"
            export PGHOST="localhost"
            export PGPORT=5432
            export PGUSER="postgres"
            export PGPASSWORD="postgres"
            export PGDATABASE="factcheck"
            
            echo "Setting up Postgres test environment..."
            
            # Initialize Postgres if not already done
            if [ ! -d "$PGDATA" ]; then
              echo "Initializing Postgres database..."
              initdb -D "$PGDATA" --auth=trust
              echo "host all all 127.0.0.1/32 trust" >> "$PGDATA/pg_hba.conf"
              echo "host all all ::1/128 trust" >> "$PGDATA/pg_hba.conf"
            fi
            
            # Check if port 5432 is already in use
            echo "Checking if port 5432 is available..."
            if lsof -i :5432 > /dev/null 2>&1; then
              echo "Warning: Port 5432 is already in use:"
              lsof -i :5432
              echo "Trying to use a different port..."
              export PGPORT=5433
              echo "Using port $PGPORT instead"
            fi
            
            # Start Postgres
            echo "Starting Postgres..."
            pg_ctl -D "$PGDATA" -l postgres.log start
            
            # Wait for Postgres to be ready with timeout (1:30 minutes)
            echo "Waiting for Postgres to be ready on port $PGPORT..."
            timeout 90 bash -c "
              until pg_isready -h localhost -p $PGPORT; do
                echo \"Waiting for Postgres on port $PGPORT...\"
                sleep 2
              done
              echo \"Postgres is ready on port $PGPORT!\"
            "
            
            if [ $? -ne 0 ]; then
              echo "Error: Postgres failed to start within 90 seconds"
              echo "Postgres log:"
              cat postgres.log
              echo "Process status:"
              ps aux | grep postgres
              echo "Port status:"
              lsof -i :$PGPORT || echo "No process using port $PGPORT"
              exit 1
            fi
            
            # Set up database schema if it exists
            if [ -f "factcheck/data/postgres/schema.sql" ]; then
              echo "Setting up database schema..."
              psql -d factcheck -f factcheck/data/postgres/schema.sql
            fi
            
            echo "Test environment ready! Run your integration tests now."
            echo "Postgres will be stopped when you exit this shell."
          '';
          
          # Clean up Postgres when shell exits
          shellExitHook = ''
            echo "Stopping Postgres..."
            pg_ctl -D "$PGDATA" stop
            echo "Postgres stopped."
          '';
        };
      });
    };
}
