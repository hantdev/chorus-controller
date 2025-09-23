#!/bin/bash

# Setup environment configuration for Chorus Controller
# This script helps you create a .env file from the example template

set -e

echo "🔧 Setting up Chorus Controller environment configuration..."

# Check if .env already exists
if [ -f ".env" ]; then
    echo "⚠️  .env file already exists!"
    read -p "Do you want to overwrite it? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ Setup cancelled."
        exit 1
    fi
fi

# Copy example to .env
if [ -f "env.example" ]; then
    cp env.example .env
    echo "✅ Created .env file from env.example"
else
    echo "❌ env.example not found!"
    exit 1
fi

# Generate random secrets
echo "🔐 Generating random secrets..."
JWT_SECRET=$(openssl rand -base64 32)
ENCRYPTION_KEY=$(openssl rand -base64 32)

sed -i.bak "s|your-secret-key-change-this-in-production-make-it-long-and-random|$JWT_SECRET|" .env
sed -i.bak "s|your-encryption-key-change-this-in-production-make-it-long-and-random|$ENCRYPTION_KEY|" .env
rm .env.bak 2>/dev/null || true

echo "✅ Generated random JWT secret and encryption key"

# Ask for database configuration
echo ""
echo "📊 Database Configuration:"
read -p "PostgreSQL DSN (default: postgres://postgres:postgres@localhost:5432/chorus?sslmode=disable): " postgres_dsn
if [ -z "$postgres_dsn" ]; then
    postgres_dsn="postgres://postgres:postgres@localhost:5432/chorus?sslmode=disable"
fi
sed -i.bak "s|postgres://postgres:postgres@localhost:5432/chorus?sslmode=disable|$postgres_dsn|" .env
rm .env.bak 2>/dev/null || true

# Ask for HTTP port
echo ""
echo "🌐 HTTP Server Configuration:"
read -p "HTTP Port (default: 8081): " http_port
if [ -z "$http_port" ]; then
    http_port="8081"
fi
sed -i.bak "s|HTTP_PORT=8081|HTTP_PORT=$http_port|" .env
rm .env.bak 2>/dev/null || true

# Ask for worker gRPC address
echo ""
echo "🔗 Worker gRPC Configuration:"
read -p "Worker gRPC Address (default: localhost:9670): " worker_addr
if [ -z "$worker_addr" ]; then
    worker_addr="localhost:9670"
fi
sed -i.bak "s|WORKER_GRPC_ADDR=localhost:9670|WORKER_GRPC_ADDR=$worker_addr|" .env
rm .env.bak 2>/dev/null || true

# Ask for JWT expiry
echo ""
echo "⏰ JWT Token Configuration:"
read -p "JWT Token Expiry (default: 24h): " jwt_expiry
if [ -z "$jwt_expiry" ]; then
    jwt_expiry="24h"
fi
sed -i.bak "s|JWT_EXPIRY=24h|JWT_EXPIRY=$jwt_expiry|" .env
rm .env.bak 2>/dev/null || true

# Ask for environment
echo ""
echo "🏗️  Environment:"
read -p "Environment (development/production, default: development): " env
if [ -z "$env" ]; then
    env="development"
fi
sed -i.bak "s|ENV=development|ENV=$env|" .env
rm .env.bak 2>/dev/null || true

echo ""
echo "🎉 Environment setup completed!"
echo ""
echo "📋 Configuration Summary:"
echo "   Database: $postgres_dsn"
echo "   HTTP Port: $http_port"
echo "   Worker gRPC: $worker_addr"
echo "   JWT Expiry: $jwt_expiry"
echo "   Environment: $env"
echo ""
echo "🔒 Security Notes:"
echo "   - Your .env file contains sensitive information"
echo "   - Never commit .env to version control"
echo "   - Keep your JWT secret secure"
echo ""
echo "🚀 Next steps:"
echo "   1. Review your .env file: cat .env"
echo "   2. Run database migrations: make migrate-apply"
echo "   3. Start the application: make run"
echo ""
echo "📖 For more information, see docs/authentication.md"
