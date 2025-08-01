# line-fact-check

line-fact-check is a free and open source fact checker platform.

It is structured as a multi-language monorepo, with focus on cross-platform development.

# Running this project with Docker Compose (WIP 90%)

> Note: this [Docker Compose setup](/docker-compose.yaml) uses **your host's Docker daemon** socket mount
> for it to be able to load images to host Docker daemon.
> 
> This has security implications, but since everything is controlled by Nix, the risk is minimized.
> We still need a better way to pin Docker Compose's services to images built by the build service.

Because Nix use is encouraged in the build process, the best way to bring up the project's environment
is to use Docker Compose to fire up a `nixos/nix` container to build and run images locally without
having to install Nix on your local machine.

To pre-build, build, or rebuild the images, run

> Note that this has a side-effect of loading whatever images we're building with `compose` tag

```sh
# Start build-images service
# This will build the Docker images and load it to your host Docker daemon

docker-compose up build-images
```

After `build-images` has finished, we can bring up other components:

```sh
# Bring up all components
docker-compose up
```

Or you can refer to individual Docker Compose service:

```sh
# Start PostgreSQL listening on local port 5432
docker-compose up -d factcheck-postgres

# Start our API
docker-compose up -d factcheck-api

# Start our frontend NGINX
docker-compose up -d backoffice-webapp
```

To simplify and minimize Nix build time, we use stateful multi-stage Docker Compose definition for each compose run.
Each run is also attached to a Docker volume, which is shared across runs, so that we can cache Nix store across builds.

You must manually manage your Nix store cache volum for the compose.

## Docker Compose build service `build-images`

`build-images` is our build service for our compose.

The service *is stateful* in that it copies your working directory into the container
when it starts. That is used to build images for the entire compose.

The built images is tagged with `latest`, so that downstream app services can find them.

## Docker Compose app component services

The rest of services are our app components: the HTTP backend, frontend, and PostgreSQL database.
To rebuild images for the compose after code changes, bring current compose down and bring up
a new `build-images`.

If rebuilding does not work (i.e. app components still point to old versions),
try removing the old images before bringing up `build-images` again,
or use [our cheat sheet](#cheat-sheet-1-using-nixos-container-to-build-docker-images-on-macos).
to manually build new images.

# Nix flake

Nix flake is heavily used in our GitHub Actions CICD pipelines,
and it does most of the heavy lifting of setting stuff up and
have reproducible builds across all hosts.

Thanks to Nix flake, we can even access this repository without having to clone it first,
via Nix URL `github:kaogeek/line-fact-check`:

```sh
# Enter default test environment from GitHub's default branch,
# with PostgreSQL running on port 5432 and all the NodeJS tooling from Flake
nix develop \
    --extra-experimental-features nix-command \
    --extra-experimental-features flakes \
    github:kaogeek/line-fact-check

# Build output docker-factcheck from branch main
nix build \
    --extra-experimental-features nix-command \
    --extra-experimental-features flakes \
    github:kaogeek/line-fact-check/main#docker-factcheck
```

This ensures that every one of us and the pipelines always have the same build and test environment.

# Nix sheet cheat

The project's outputs or artifacts are all defined as flake outputs.

For example, our backend in [/factcheck](./factcheck/) is built into Go binary as flake output `#factcheck`. Another output `docker-factcheck` is just `#factcheck` inside a Docker image.

# Contributing without Nix

We can run the NixOS container and use it to try to interact with Nix:

```sh
# Try Nix command
docker run -it nixos/nix:latest
```

## Cheat sheet 1: using NixOS container to build Docker images on macOS

The simplest way to run this project via Nix is by referencing the remote GitHub repository.
We do this by running a NixOS container (with mounted volume).

We then use Nix inside that to build our image within the container.
We finish off by copying the result to host machine filesystem,
mounted at /workspace within the container

```sh
# Cheat sheet 1
docker run --rm -v $(pwd):/workspace nixos/nix:latest sh -c \
    "nix build github:kaogeek/line-fact-check/main#docker-factcheck --extra-experimental-features nix-command --extra-experimental-features flakes && cp result /workspace/"
```

However, this command is "stateless" in that every build will start from scratch with long ass build time.

We can fix that by mounting a Docker volume to build containers, and we can mount that volume to
our new builders to benefit from the cache. Let's say our new Docker volume is `line-fact-check-nix-store`
then we just give docker run this volume option `-v line-fact-check-nix-store:/nix/store`:

```sh
# Cheat sheet 1, revised with caching Docker volume
docker run --rm -v $(pwd):/workspace -v line-fact-check-nix-store:/nix/store nixos/nix:latest sh -c \
    "nix build github:kaogeek/line-fact-check/main#docker-factcheck --extra-experimental-features nix-command --extra-experimental-features flakes && cp result /workspace/"
```

## Cheat sheet 2: Cached build of Docker images from local files
The command above, while simple, will write result to your `./result`. This is stateful, and will break if run back-to-back. This script will fix that by appending our flake "version" to the name of file.

Differences from cheat sheet 1:

- Read from local directory instead of remote GitHub repository with `.#docker-factcheck`

- Unique output, where FLAKE_VERSION is appended to the result image filename

```sh
# Cheat sheet 2
docker run --rm -v $(pwd):/workspace -v line-fact-check-nix-store:/nix/store nixos/nix:latest sh -c '
    cp -r /workspace /source
    cd source

    export FLAKE_VERSION="$(nix eval --extra-experimental-features nix-command --extra-experimental-features flakes  .#version.text --raw)"

    # Build docker-factcheck with unique result name using the actual flake version
    nix build .#docker-factcheck --extra-experimental-features nix-command --extra-experimental-features flakes

    # Copy the result to workspace
    cp result /workspace/result-factcheck-$FLAKE_VERSION
    echo "copied to /workspace/result-factcheck-$FLAKE_VERSION"
'
```
