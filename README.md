# line-fact-check

line-fact-check is a free and open source fact checker platform.

It is structured as a multi-language monorepo, with focus on cross-platform development.

# Running this project

Because Nix use is encouraged in the build process, the best way to bring up the project's environment
is to use Docker Compose to fire up a `nixos/nix` container to build and run images locally.

To simplify and reduce Nix builds, we use stateful multi-stage Docker Compose definition for each compose run.

To pre-build the images, run

```sh
# Start build-images service
# This will build the Docker images and load it to your host Docker daemon

docker-compose up build-images
```

After `build-images` has finished, we can bring up other components:

```sh
docker-compose up
```

## Docker Compose build service `build-images`

`build-images` is our build service for our compose.

The service *is stateful* in that it copies your working directory into the container
when it starts. That is used to build images for the entire compose.

The built images is tagged with `latest`, so that downstream app services can find them.

To rebuild images for the compose after code changes,
bring current compose down and bring up a new `build-images`.

## Docker Compose app services

The rest of services are our app components: the HTTP backend, frontend, and PostgreSQL database.

# Nix flake

We use Nix flake to pin everything and have reproducible builds.

Nix flake is heavily used in our GitHub Actions CICD pipelines,
and it does most of the heavy lifting of setting stuff up.

Thanks to Nix flake, we can build this repository without having to clone it first:

```sh
# Build from branch main
nix build\ 
    --extra-experimental-features nix-command \ 
    --extra-experimental-features flakes \ 
    github:kaogeek/line-fact-check/main#docker-factcheck
```

This ensures that every one of us and the pipelines always have the same build and test environment.

# Nix sheet cheat

The project's outputs or artifacts are all defined as flake outputs.

For example, our backend in [/factcheck](./factcheck/) is built into Go binary as flake output `#factcheck`. Another output `docker-factcheck` is just `#factcheck` inside a Docker image.

# Contributing without Nix

However, most of our devs don't have Nix installed locally. But most have Docker.

To hack around this, the official NixOS image `nixos/nix` is actually very handy.
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
docker run --rm -v $(pwd):/workspace nixos/nix:latest sh -c \
    "nix build github:kaogeek/line-fact-check/main#docker-factcheck --extra-experimental-features nix-command --extra-experimental-features flakes && cp result /workspace/"
```

## Cheat sheet 2: Cached build of Docker images from local files
The command above, while simple, will write result to your `./result`. This is stateful, and will break if run back-to-back. This script will fix that by appending our flake "version" to the name of file.

Differences to example above

- Read from local directory instead of remote GitHub repository with `.#docker-factcheck`

- Cache Nix store with simple Docker volume with `-v line-fact-check-nix-store:/nix/store`.

- Unique output, where FLAKE_VERSION is appended to the result image filename

```sh
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
