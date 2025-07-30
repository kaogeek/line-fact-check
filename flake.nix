rec {
  description = "Fact checker LINE Bot";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      homepage = "https://github.com/kaogeek/line-fact-check";
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
      version = "${builtins.substring 0 8 lastModifiedDate}-${builtins.toString self.lastModified}";

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
        version = pkgs.writeTextFile {
          name = "factcheck-version";
          text = version;
          destination = "/version.txt";
        };

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
          vendorHash = "sha256-dpggYNs50dyXcfZYH/TCw/uLB2Nl/Kf0nGqQfrNpQm8=";
          meta = {
            inherit homepage;
            description = "${description} - factcheck";
          };
        };

        backoffice-webapp = pkgs.buildNpmPackage {
          inherit version;
          pname = "backoffice-webapp";
          src = ./backoffice-webapp;
          
          npmDepsHash = "sha256-oWGQ/jfqkm3emN96Wx9i3hhryS8BOQhoIZ3TdRqU/yI="; # This will be updated after first build
          
          installPhase = ''
            mkdir -p $out
            cp -r dist/* $out/
          '';
          
          meta = {
            inherit homepage;
            description = "${description} - backoffice webapp";
          };
        };

        # To build and load the image:
        # nix build .#docker-factcheck && docker load < result
        docker-factcheck = pkgs.dockerTools.buildImage {
          name = "factcheck";
          tag = version;
          config = {
            Entrypoint = [ "${self.packages.${pkgs.system}.factcheck}/bin/api" ];
            ExposedPorts = {
              "8080/tcp" = {};
            };
            Healthcheck = {
              Test = [
                "CMD"
                "${self.packages.${pkgs.system}.factcheck}/bin/healthcheck"
              ];
              StartPeriod = 5000000000; # 5 seconds
              Interval = 15000000000; # 15 seconds
              Timeout = 5000000000;  # 5 seconds
              Retries = 3;
            };
          };
        };

        docker-foo = pkgs.dockerTools.buildImage {
          name = "factcheck/foo";
          tag = version;
          config = {
            Entrypoint = [ "${self.packages.${pkgs.system}.foo}/bin/foo" ];
            ExposedPorts = {
              "8080/tcp" = {};
            };
          };
        };

        # PostgreSQL Docker image for integration tests
        # nix build .#docker-postgres-factcheck && docker load < result
        docker-postgres-factcheck = pkgs.dockerTools.buildImage {
          name = "factcheck/postgres";
          tag = "16";
          fromImage = if pkgs.lib.strings.hasPrefix "aarch64" pkgs.system then pkgs.dockerTools.pullImage {
            imageName = "postgres";
            imageDigest = "sha256:f756bec6e37ffdae0f94b499d49a48cd35dd4eea1edb06711b7fbd9c5708d6d5";
            finalImageName = "postgres";
            finalImageTag = "16";
            sha256 = "sha256-d9a0i0VhnCEqQrRdkFCSkZC2OTOtjVoxRZWocK0zCUg=";
          } else pkgs.dockerTools.pullImage {
            imageName = "postgres";
            imageDigest = "sha256:85df95e724f37d83e3cdc771bc7928519a0febcaeda7f3b497986a0d04959c0d";
            finalImageName = "postgres";
            finalImageTag = "16";
            sha256 = "sha256-Gj6yGJdD0/eyMkXDrbEPPuPd89018+zyQ6EG7faC83g=";
          };
          copyToRoot = pkgs.runCommand "postgres-init-schema-factcheck" {} ''
            mkdir -p $out/docker-entrypoint-initdb.d
            cp ${./factcheck/internal/data/postgres/schema.sql} $out/docker-entrypoint-initdb.d/01-schema.sql
          '';
          config = {
            Entrypoint = [ "docker-entrypoint.sh" ];
            Cmd = [ "postgres" ];
            Env = [
              "POSTGRES_DB=factcheck"
            ];
          };
        };

        # NGINX Docker image for serving frontend static files
        # nix build .#docker-backoffice-webapp && docker load < result
        docker-backoffice-webapp = pkgs.dockerTools.buildImage {
          name = "factcheck/backoffice-webapp";
          tag = version;
          fromImage = if pkgs.lib.strings.hasPrefix "aarch64" pkgs.system then pkgs.dockerTools.pullImage {
            imageName = "arm64v8/nginx";
            imageDigest = "sha256:fb634803c8e82bf44e6260504c84c1420bcd965ac32d002273d489cb7c6057d9";
            finalImageName = "arm64v8/nginx";
            finalImageTag = "stable-alpine";
            sha256 = "sha256-srzTrhkuvXpOvS42Keir0rZkCsDaymlXEO3fggLu8vE=";
          } else pkgs.dockerTools.pullImage {
            imageName = "nginx";
            imageDigest = "sha256:64a376d12f051d5b97f8825514a7621bfd613ebad1cb1876354b9a42c9b17891";
            finalImageName = "nginx";
            finalImageTag = "alpine";
            sha256 = "sha256-YZcEgn7GaL0LD8EbijdyNKR29XvV9YbCrAA3VlwbXG0=";
          };
          copyToRoot = pkgs.runCommand "nginx-config-backoffice-webapp" {} ''
            mkdir -p $out/etc/nginx/conf.d
            mkdir -p $out/usr/share/nginx/html
            # Copy static files
            cp -r ${self.packages.${pkgs.system}.backoffice-webapp}/* $out/usr/share/nginx/html/
          '';
          config = {
            Entrypoint = [ "/docker-entrypoint.sh" ];
            Cmd = [ "nginx" "-g" "daemon off;" ];
            ExposedPorts = {
              "80/tcp" = {};
            };
          };
        };
      });

      devShells = forAllSystems ({ pkgs }: let
        packagesBackend = with pkgs; [
          # Development - server
          go
          gopls
          gotools
          go-tools
          golangci-lint
          sqlc
          wire
        ];

        packagesFrontend = with pkgs; [
          # Development - webapp
          nodejs_22
          nodePackages.npm
        ];

        packagesItTest = with pkgs; [
          docker
          docker-compose
          coreutils
          bash
        ];

        in {
        # Default shell has everything, but nothing running
        default = pkgs.mkShell {
          packages = packagesBackend ++ packagesFrontend ++ packagesItTest;
          shellHook = ''
            echo "Entering Nix default devShell"
            export FACTCHECK_VERSION=${version}
            echo "FACTCHECK_VERSION=$FACTCHECK_VERSION"
          '';
        };

        # Shell with Go and code-gen tools
        go-develop = pkgs.mkShell {
          packages = packagesBackend;
          shellHook = ''
            echo "Entering Nix shell go-develop"
          '';
        };

        # Shell for running integration tests with PostgreSQL
        go-it-test = pkgs.mkShell {
          packages = packagesBackend ++ packagesItTest;

          FACTCHECKAPI_LISTEN_ADDRESS = ":8080";
          FACTCHECKAPI_TIMEOUTMS_READ = "3000";
          FACTCHECKAPI_TIMEOUTMS_WRITE = "3000";
          POSTGRES_USER = "postgres";
          POSTGRES_PASSWORD = "postgres";
          POSTGRES_DB = "factcheck";
          POSTGRES_PORT = "5432";

          shellHook = ''
            echo "Entering Nix shell go-it-test"
            echo "Loading PostgreSQL image from Nix..."
            docker load < ${self.packages.${pkgs.system}.docker-postgres-factcheck}
            
            echo "Starting PostgreSQL container for integration tests..."
            docker run -d \
              --name postgres-it-test \
              -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
              -e POSTGRES_USER=$POSTGRES_USER \
              -e POSTGRES_DB=$POSTGRES_DB \
              -p 5432:$POSTGRES_PORT \
              factcheck/postgres:16

            echo "PostgreSQL container started on $POSTGRES_HOST:$POSTGRES_PORT"
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
