#!/bin/bash
set -e  # Exit on any error

echo "🚀 Starting build process..."

# Install migrate CLI
echo "📦 Installing golang-migrate..."
MIGRATE_VERSION="v4.18.3"
curl -L "https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz" | tar xvz
chmod +x migrate
sudo mv migrate /usr/local/bin/migrate || mv migrate /usr/bin/migrate 2>/dev/null || export PATH="$PWD:$PATH"

# Verify migrate installation
echo "✅ Verifying migrate installation..."
migrate -version

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "❌ ERROR: DATABASE_URL environment variable is not set"
    exit 1
fi

echo "🗄️  Running database migrations..."
migrate -path database/migrations -database "$DATABASE_URL" up

echo "🏗️  Building Go application..."
# Vercel handles the Go build automatically, but you can add custom build steps here if needed
# go build -o main .

echo "✅ Build process completed successfully!"