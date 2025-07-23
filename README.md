# line-fact-check

line-fact-check is a free and open source fact checker platform.

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

# Nix outputs

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

We can easily build Docker images from any Git commits in our repository.

We do this by running a NixOS container (with mounted volume). We then use Nix inside that to build our image within the container. We finish off by copying the result

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