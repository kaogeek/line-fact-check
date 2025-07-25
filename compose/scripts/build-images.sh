#!/bin/bash
set -e

echo "Building Docker images from Nix flake"

cp -r /source /workspace
cd /workspace

echo "Nix system architecture:"
nix eval --impure --expr 'builtins.currentSystem' --extra-experimental-features nix-command

rm -f result
echo 'Building factcheck image'
nix build .#docker-factcheck --extra-experimental-features nix-command --extra-experimental-features flakes
docker load < result
docker images factcheck --format "table {{.Repository}}:{{.Tag}}" | grep -v "TAG" | head -1 | xargs -I {} docker tag {} factcheck:compose

rm -f result
echo 'Building PostgreSQL image...'
nix build .#docker-postgres-factcheck --extra-experimental-features nix-command --extra-experimental-features flakes
docker load < result
docker images factcheck/postgres --format "table {{.Repository}}:{{.Tag}}" | grep -v "TAG" | head -1 | xargs -I {} docker tag {} factcheck/postgres:compose

rm -f result
echo 'Building backoffice webapp image...'
nix build .#docker-backoffice-webapp --extra-experimental-features nix-command --extra-experimental-features flakes
docker load < result
docker images factcheck/backoffice-webapp --format "table {{.Repository}}:{{.Tag}}" | grep -v "TAG" | head -1 | xargs -I {} docker tag {} factcheck/backoffice-webapp:compose

echo 'All builds completed successfully'
