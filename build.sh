#!/bin/bash
set -e

echo "🚀 Starting migration process..."

# Create local bin directory
mkdir -p ./bin

# Install migrate CLI locally
echo "📦 Installing golang-migrate..."
MIGRATE_VERSION="v4.18.3"
curl -L "https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz" | tar xvz
chmod +x migrate
mv migrate ./bin/
export PATH="$PWD/bin:$PATH"

echo "✅ Verifying migrate installation..."
migrate -version

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "❌ ERROR: DATABASE_URL environment variable is not set"
    exit 1
fi

echo "🗄️  Running database migrations..."
migrate -path database/migrations -database "$DATABASE_URL" up

echo "✅ Migrations completed!"