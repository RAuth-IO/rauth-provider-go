# Test script for Fiber + RauthProvider integration (Windows PowerShell)
Write-Host "🧪 Testing Fiber + RauthProvider Integration" -ForegroundColor Green
Write-Host "=============================================" -ForegroundColor Green

# Set environment variables (replace with your actual values)
$env:RAUTH_API_KEY = "your_api_key_here"
$env:RAUTH_APP_ID = "your_app_id_here"
$env:RAUTH_WEBHOOK_SECRET = "your_webhook_secret_here"

# Start the server in background
Write-Host "🚀 Starting Fiber server..." -ForegroundColor Yellow
$job = Start-Job -ScriptBlock {
    Set-Location $using:PWD
    cd simple
    go run .
}

# Wait for server to start
Start-Sleep -Seconds 3

# Test health endpoint
Write-Host ""
Write-Host "📊 Testing health endpoint..." -ForegroundColor Cyan
try {
    $response = Invoke-RestMethod -Uri "http://localhost:3000/health" -Method Get
    $response | ConvertTo-Json -Depth 3
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

# Test login endpoint (this will fail without valid session token)
Write-Host ""
Write-Host "🔐 Testing login endpoint..." -ForegroundColor Cyan
try {
    $body = @{
        session_token = "test_token"
        user_phone = "+1234567890"
    } | ConvertTo-Json

    $response = Invoke-RestMethod -Uri "http://localhost:3000/login" -Method Post -Body $body -ContentType "application/json"
    $response | ConvertTo-Json -Depth 3
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

# Test protected endpoint without auth
Write-Host ""
Write-Host "🚫 Testing protected endpoint without auth..." -ForegroundColor Cyan
try {
    $response = Invoke-RestMethod -Uri "http://localhost:3000/protected" -Method Get
    $response | ConvertTo-Json -Depth 3
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

# Test protected endpoint with invalid auth
Write-Host ""
Write-Host "❌ Testing protected endpoint with invalid auth..." -ForegroundColor Cyan
try {
    $headers = @{
        "Authorization" = "Bearer invalid_token"
        "X-User-Phone" = "+1234567890"
    }
    
    $response = Invoke-RestMethod -Uri "http://localhost:3000/protected" -Method Get -Headers $headers
    $response | ConvertTo-Json -Depth 3
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

# Cleanup
Write-Host ""
Write-Host "🧹 Cleaning up..." -ForegroundColor Yellow
Stop-Job $job
Remove-Job $job

Write-Host ""
Write-Host "✅ Test completed!" -ForegroundColor Green
Write-Host ""
Write-Host "📝 To test with real data:" -ForegroundColor Yellow
Write-Host "1. Set your actual RAUTH_API_KEY, RAUTH_APP_ID, and RAUTH_WEBHOOK_SECRET"
Write-Host "2. Get a valid session token from your Rauth integration"
Write-Host "3. Run: Invoke-RestMethod -Uri 'http://localhost:3000/login' -Method Post -Body '{\"session_token\": \"YOUR_SESSION_TOKEN\", \"user_phone\": \"YOUR_PHONE\"}' -ContentType 'application/json'"
Write-Host ""
Write-Host "4. Use the session token in protected routes:"
Write-Host "   Invoke-RestMethod -Uri 'http://localhost:3000/protected' -Method Get -Headers @{\"Authorization\"=\"Bearer YOUR_SESSION_TOKEN\"; \"X-User-Phone\"=\"YOUR_PHONE\"}" 