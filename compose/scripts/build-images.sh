#!/bin/bash
set -e

echo "Building Docker images from Nix flake"

cp -r /source /workspace
cd /workspace

# Print the system architecture that Nix is using
echo "Nix system architecture:"
nix eval --impure --expr 'builtins.currentSystem' --extra-experimental-features nix-command
echo ""

# Clean up any existing result symlinks
rm -f result

echo 'Building factcheck image'
nix build .#docker-factcheck --extra-experimental-features nix-command --extra-experimental-features flakes
docker load < result

# Add a latest tag for easier reference
docker images factcheck --format "table {{.Repository}}:{{.Tag}}" | grep -v "TAG" | head -1 | xargs -I {} docker tag {} factcheck:latest

# Clean up for next build
rm -f result

echo 'Building PostgreSQL image...'
nix build .#docker-postgres-factcheck --extra-experimental-features nix-command --extra-experimental-features flakes
docker load < result

echo 'All builds completed successfully' 