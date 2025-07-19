# RauthProvider Go Library Publishing Script (PowerShell)
# This script helps you publish the Go library to GitHub

param(
    [Parameter(Mandatory=$true)]
    [string]$GitHubUsername
)

Write-Host "🚀 RauthProvider Go Library Publishing Script" -ForegroundColor Green
Write-Host "==============================================" -ForegroundColor Green

# Check if git is installed
try {
    git --version | Out-Null
} catch {
    Write-Host "❌ Git is not installed. Please install Git first." -ForegroundColor Red
    exit 1
}

# Check if we're in the right directory
if (-not (Test-Path "go.mod")) {
    Write-Host "❌ Please run this script from the Go library directory" -ForegroundColor Red
    Write-Host "   cd 'Backend-Library/Go Lang'" -ForegroundColor Yellow
    exit 1
}

Write-Host ""
Write-Host "📋 Publishing Steps:" -ForegroundColor Cyan
Write-Host "1. Initialize Git repository"
Write-Host "2. Add all files"
Write-Host "3. Create initial commit"
Write-Host "4. Connect to GitHub"
Write-Host "5. Push to GitHub"
Write-Host "6. Create release tag"
Write-Host ""

$continue = Read-Host "Do you want to continue? (y/N)"
if ($continue -ne "y" -and $continue -ne "Y") {
    Write-Host "❌ Publishing cancelled" -ForegroundColor Red
    exit 1
}

# Step 1: Initialize Git repository
Write-Host "📁 Initializing Git repository..." -ForegroundColor Yellow
if (-not (Test-Path ".git")) {
    git init
} else {
    Write-Host "⚠️  Git repository already exists" -ForegroundColor Yellow
}

# Step 2: Add all files
Write-Host "📦 Adding files to Git..." -ForegroundColor Yellow
git add .

# Step 3: Create initial commit
Write-Host "💾 Creating initial commit..." -ForegroundColor Yellow
$commitMessage = @"
Initial commit: RauthProvider Go library

Features:
- Session management with TTL
- Webhook support for real-time updates
- HTTP middleware integration
- Signature-based verification
- In-memory session tracking with API fallback
- Clean architecture design
"@

git commit -m $commitMessage

# Step 4: Connect to GitHub
Write-Host "🔗 Connecting to GitHub..." -ForegroundColor Yellow
try {
    git remote add origin "https://github.com/$GitHubUsername/rauth-provider-go.git"
} catch {
    Write-Host "⚠️  Remote origin already exists" -ForegroundColor Yellow
}

# Step 5: Push to GitHub
Write-Host "⬆️  Pushing to GitHub..." -ForegroundColor Yellow
git branch -M main
git push -u origin main

# Step 6: Create release tag
Write-Host "🏷️  Creating release tag v1.0.0..." -ForegroundColor Yellow
git tag v1.0.0
git push origin v1.0.0

Write-Host ""
Write-Host "✅ Publishing completed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "📊 Next steps:" -ForegroundColor Cyan
Write-Host "1. Visit: https://github.com/$GitHubUsername/rauth-provider-go"
Write-Host "2. Create a GitHub Release for v1.0.0"
Write-Host "3. Add description and release notes"
Write-Host ""
Write-Host "🔗 Your library is now available at:" -ForegroundColor Green
Write-Host "   go get github.com/$GitHubUsername/rauth-provider-go"
Write-Host ""
Write-Host "📖 Documentation will be available at:" -ForegroundColor Green
Write-Host "   https://pkg.go.dev/github.com/$GitHubUsername/rauth-provider-go"
Write-Host ""
Write-Host "🎉 Congratulations! Your Go library is now published!" -ForegroundColor Green 