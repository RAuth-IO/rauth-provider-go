# Publishing Guide for RauthProvider Go Library

This guide will help you publish the Go rauth-provider library to make it available for others to use.

## 📋 Prerequisites

1. **GitHub Account**: You need a GitHub account to host the repository
2. **Go 1.16+**: Ensure you have Go 1.16 or later installed
3. **Git**: Make sure Git is installed and configured

## 🚀 Step-by-Step Publishing Process

### Step 1: Create GitHub Repository

1. Go to [GitHub](https://github.com) and create a new repository
2. Repository name: `rauth-provider-go`
3. Description: "A lightweight, plug-and-play Go library for phone number authentication using Rauth.io"
4. Make it **Public** (required for Go modules)
5. **Don't** initialize with README, .gitignore, or license (we already have these)

### Step 2: Initialize Git Repository

```bash
cd "Backend-Library/Go Lang"
git init
git add .
git commit -m "Initial commit: RauthProvider Go library"
```

### Step 3: Connect to GitHub

```bash
git remote add origin https://github.com/YOUR_USERNAME/rauth-provider-go.git
git branch -M main
git push -u origin main
```

### Step 4: Create Release Tags

Go modules use semantic versioning. Create tags for releases:

```bash
# For first release (v1.0.0)
git tag v1.0.0
git push origin v1.0.0

# For future releases
git tag v1.1.0
git push origin v1.1.0
```

### Step 5: Update Module Name (if needed)

If you want to use a different module name, update `go.mod`:

```go
module github.com/YOUR_USERNAME/rauth-provider-go
```

Then update all import paths in the code accordingly.

### Step 6: Test the Module

Create a test project to verify the module works:

```bash
mkdir test-module
cd test-module
go mod init test
```

Create `main.go`:
```go
package main

import (
    "log"
    "github.com/YOUR_USERNAME/rauth-provider-go/internal/domain"
    "github.com/YOUR_USERNAME/rauth-provider-go/pkg/rauthprovider"
)

func main() {
    config := &domain.Config{
        RauthAPIKey:   "test-key",
        AppID:         "test-app",
        WebhookSecret: "test-secret",
    }
    
    if err := rauthprovider.Init(config); err != nil {
        log.Fatal(err)
    }
    
    log.Println("Module imported successfully!")
}
```

Run the test:
```bash
go mod tidy
go run main.go
```

## 📦 Publishing to Go Module Proxy

Once your repository is on GitHub with proper tags, the Go module proxy will automatically pick it up:

1. **Go Module Proxy**: Automatically indexes public repositories
2. **No manual registration required** for public repositories
3. **Users can install** using `go get github.com/YOUR_USERNAME/rauth-provider-go`

## 🔧 Additional Publishing Options

### Option 1: Publish to pkg.go.dev

1. Your module will automatically appear on [pkg.go.dev](https://pkg.go.dev)
2. Users can browse documentation, examples, and download statistics
3. No additional steps required for public repositories

### Option 2: Create GitHub Release

1. Go to your GitHub repository
2. Click "Releases" → "Create a new release"
3. Tag version: `v1.0.0`
4. Title: `v1.0.0 - Initial Release`
5. Description: Include features, breaking changes, etc.
6. Upload compiled binaries if needed

### Option 3: Add to Awesome Go List

1. Fork the [Awesome Go](https://github.com/avelino/awesome-go) repository
2. Add your library to the appropriate section
3. Submit a pull request

## 📝 Documentation Updates

### Update README.md

Make sure your README includes:

1. **Installation instructions** with correct module name
2. **Quick start examples**
3. **API documentation**
4. **Usage examples**
5. **Contributing guidelines**

### Add Go Documentation Comments

Ensure all exported functions have proper documentation:

```go
// VerifySession verifies if a session is valid and matches the phone number.
// It returns true if the session is valid, false otherwise.
func VerifySession(ctx context.Context, sessionToken, userPhone string) (bool, error) {
    // implementation
}
```

## 🧪 Testing Before Publishing

### Run All Tests

```bash
go test ./...
```

### Run Examples

```bash
cd examples/basic
go run main.go
```

### Check for Issues

```bash
# Run linter
go vet ./...

# Check for security issues
go list -json -deps ./... | nancy sleuth
```

## 📊 Post-Publishing

### Monitor Usage

1. **pkg.go.dev**: Check download statistics
2. **GitHub**: Monitor stars, forks, and issues
3. **Go Module Proxy**: Check download counts

### Maintain the Library

1. **Respond to issues** promptly
2. **Review pull requests** regularly
3. **Update dependencies** when needed
4. **Release new versions** for features and bug fixes

## 🔗 Useful Links

- [Go Modules Documentation](https://golang.org/doc/modules)
- [pkg.go.dev](https://pkg.go.dev) - Go package discovery
- [Go Module Proxy](https://proxy.golang.org) - Module proxy
- [Semantic Versioning](https://semver.org/) - Version numbering

## 🎯 Example Publishing Commands

```bash
# Complete publishing workflow
cd "Backend-Library/Go Lang"

# Initialize git
git init
git add .
git commit -m "Initial commit: RauthProvider Go library"

# Connect to GitHub (replace with your username)
git remote add origin https://github.com/YOUR_USERNAME/rauth-provider-go.git
git branch -M main
git push -u origin main

# Create first release
git tag v1.0.0
git push origin v1.0.0

# Test the module
mkdir ../test-module
cd ../test-module
go mod init test
echo 'package main

import (
    "log"
    "github.com/YOUR_USERNAME/rauth-provider-go/internal/domain"
    "github.com/YOUR_USERNAME/rauth-provider-go/pkg/rauthprovider"
)

func main() {
    config := &domain.Config{
        RauthAPIKey:   "test-key",
        AppID:         "test-app",
        WebhookSecret: "test-secret",
    }
    
    if err := rauthprovider.Init(config); err != nil {
        log.Fatal(err)
    }
    
    log.Println("Module imported successfully!")
}' > main.go

go mod tidy
go run main.go
```

## 🎉 Success!

Once you've completed these steps, your Go library will be:

- ✅ Available via `go get github.com/YOUR_USERNAME/rauth-provider-go`
- ✅ Listed on pkg.go.dev
- ✅ Accessible through the Go module proxy
- ✅ Ready for others to use in their projects

Users can now install your library with:
```bash
go get github.com/YOUR_USERNAME/rauth-provider-go
``` 