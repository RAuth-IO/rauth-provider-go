#!/bin/bash

# RauthProvider Go Library Publishing Script
# This script helps you publish the Go library to GitHub

set -e

echo "🚀 RauthProvider Go Library Publishing Script"
echo "=============================================="

# Check if git is installed
if ! command -v git &> /dev/null; then
    echo "❌ Git is not installed. Please install Git first."
    exit 1
fi

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Please run this script from the Go library directory"
    echo "   cd 'Backend-Library/Go Lang'"
    exit 1
fi

# Get GitHub username
read -p "Enter your GitHub username: " GITHUB_USERNAME

if [ -z "$GITHUB_USERNAME" ]; then
    echo "❌ GitHub username is required"
    exit 1
fi

echo ""
echo "📋 Publishing Steps:"
echo "1. Initialize Git repository"
echo "2. Add all files"
echo "3. Create initial commit"
echo "4. Connect to GitHub"
echo "5. Push to GitHub"
echo "6. Create release tag"
echo ""

read -p "Do you want to continue? (y/N): " -n 1 -r
echo

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Publishing cancelled"
    exit 1
fi

# Step 1: Initialize Git repository
echo "📁 Initializing Git repository..."
if [ ! -d ".git" ]; then
    git init
else
    echo "⚠️  Git repository already exists"
fi

# Step 2: Add all files
echo "📦 Adding files to Git..."
git add .

# Step 3: Create initial commit
echo "💾 Creating initial commit..."
git commit -m "Initial commit: RauthProvider Go library

Features:
- Session management with TTL
- Webhook support for real-time updates
- HTTP middleware integration
- Signature-based verification
- In-memory session tracking with API fallback
- Clean architecture design"

# Step 4: Connect to GitHub
echo "🔗 Connecting to GitHub..."
git remote add origin "https://github.com/$GITHUB_USERNAME/rauth-provider-go.git" 2>/dev/null || {
    echo "⚠️  Remote origin already exists"
}

# Step 5: Push to GitHub
echo "⬆️  Pushing to GitHub..."
git branch -M main
git push -u origin main

# Step 6: Create release tag
echo "🏷️  Creating release tag v1.0.0..."
git tag v1.0.0
git push origin v1.0.0

echo ""
echo "✅ Publishing completed successfully!"
echo ""
echo "📊 Next steps:"
echo "1. Visit: https://github.com/$GITHUB_USERNAME/rauth-provider-go"
echo "2. Create a GitHub Release for v1.0.0"
echo "3. Add description and release notes"
echo ""
echo "🔗 Your library is now available at:"
echo "   go get github.com/$GITHUB_USERNAME/rauth-provider-go"
echo ""
echo "📖 Documentation will be available at:"
echo "   https://pkg.go.dev/github.com/$GITHUB_USERNAME/rauth-provider-go"
echo ""
echo "🎉 Congratulations! Your Go library is now published!" 