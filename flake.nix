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

        # To build and load the image:
        # nix build .#docker-factcheck && docker load < result
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

        # PostgreSQL Docker image for integration tests
        # nix build .#docker-postgres-it-test && docker load < result
        docker-postgres-it-test = pkgs.dockerTools.buildImage {
          name = "postgres-factcheck";
          tag = "16";
          fromImage = pkgs.dockerTools.pullImage {
            imageName = "postgres";
            imageDigest = "sha256:c0aab7962b283cf24a0defa5d0d59777f5045a7be59905f21ba81a20b1a110c9";
            sha256 = if pkgs.system == "x86_64-darwin" then
              "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
            else if pkgs.system == "aarch64-darwin" then
              "sha256-EdHeqBwnd84kFi2QEFbDT+eE/F1r09OFDVvp56MS+RQ="
            else
              "sha256-TWrE5ZILio0f+WKvyWjOvCIc6+diPhPeVQoPR32JSdw=";
            finalImageName = "postgres";
            finalImageTag = "16";
          };
          runAsRoot = ''
            cp ${./factcheck/data/postgres/schema.sql} /docker-entrypoint-initdb.d/01-schema.sql
          '';
          config = {
            Entrypoint = [ "docker-entrypoint.sh" ];
            Cmd = [ "postgres" ];
            Env = [
              "POSTGRES_PASSWORD=postgres"
              "POSTGRES_USER=postgres"
              "POSTGRES_DB=factcheck"
            ];
          };
        };
      });

      devShells = forAllSystems ({ pkgs }: let
        packagesDevelop = with pkgs; [
          # Development - server
          go
          gopls
          gotools
          go-tools
          golangci-lint
          sqlc
          wire
        ];
        packagesItTest = with pkgs; [
          docker
          docker-compose
          coreutils
          bash
        ];
        packagesExtra = with pkgs; [
          # Basic LSPs
          nixd
          nixpkgs-fmt
          bash-language-server
          shellcheck
          shfmt
          lowdown
        ];

        in {
        default = pkgs.mkShell {
          packages = packagesDevelop ++ packagesExtra ++ packagesItTest;
        };

        # Shell for running integration tests with PostgreSQL
        shell-it-test = pkgs.mkShell {
          packages = packagesDevelop ++ packagesItTest;
          shellHook = ''
            echo "Loading PostgreSQL image from Nix..."
            docker load < ${self.packages.${pkgs.system}.docker-postgres-it-test}
            
            echo "Starting PostgreSQL container for integration tests..."
            docker run -d \
              --name postgres-it-test \
              -e POSTGRES_PASSWORD=postgres \
              -e POSTGRES_USER=postgres \
              -e POSTGRES_DB=factcheck \
              -p 5432:5432 \
              postgres-factcheck:16
            echo "PostgreSQL container started on localhost:5432"
            
            echo "Waiting for PostgreSQL to be ready..."
            timeout=90
            counter=0
            while [ $counter -lt $timeout ]; do
              if docker exec postgres-it-test pg_isready -U postgres -d factcheck > /dev/null 2>&1; then
                echo "PostgreSQL is ready!"
                break
              fi
              echo "Waiting for PostgreSQL... ($counter/$timeout seconds)"
              sleep 2
              counter=$((counter + 2))
            done
            
            if [ $counter -ge $timeout ]; then
              echo "Error: PostgreSQL did not become ready within $timeout seconds"
              docker logs postgres-it-test
              exit 1
            fi
            
            echo "Use 'docker stop postgres-it-test && docker rm postgres-it-test' to clean up"
          '';
        };
      });
    };
}
