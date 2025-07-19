#!/bin/bash

# Test script for Fiber + RauthProvider integration
echo "🧪 Testing Fiber + RauthProvider Integration"
echo "============================================="

# Set environment variables (replace with your actual values)
export RAUTH_API_KEY="your_api_key_here"
export RAUTH_APP_ID="your_app_id_here"
export RAUTH_WEBHOOK_SECRET="your_webhook_secret_here"

# Start the server in background
echo "🚀 Starting Fiber server..."
cd simple
go run . &
SERVER_PID=$!
cd ..

# Wait for server to start
sleep 3

# Test health endpoint
echo ""
echo "📊 Testing health endpoint..."
curl -s http://localhost:3000/health | jq .

# Test login endpoint (this will fail without valid session token)
echo ""
echo "🔐 Testing login endpoint..."
curl -s -X POST http://localhost:3000/login \
  -H "Content-Type: application/json" \
  -d '{"session_token": "test_token", "user_phone": "+1234567890"}' | jq .

# Test protected endpoint without auth
echo ""
echo "🚫 Testing protected endpoint without auth..."
curl -s http://localhost:3000/protected | jq .

# Test protected endpoint with invalid auth
echo ""
echo "❌ Testing protected endpoint with invalid auth..."
curl -s -H "Authorization: Bearer invalid_token" \
  -H "X-User-Phone: +1234567890" \
  http://localhost:3000/protected | jq .

# Cleanup
echo ""
echo "🧹 Cleaning up..."
kill $SERVER_PID

echo ""
echo "✅ Test completed!"
echo ""
echo "📝 To test with real data:"
echo "1. Set your actual RAUTH_API_KEY, RAUTH_APP_ID, and RAUTH_WEBHOOK_SECRET"
echo "2. Get a valid session token from your Rauth integration"
echo "3. Run: curl -X POST http://localhost:3000/login \\"
echo "   -H 'Content-Type: application/json' \\"
echo "   -d '{\"session_token\": \"YOUR_SESSION_TOKEN\", \"user_phone\": \"YOUR_PHONE\"}'"
echo ""
echo "4. Use the session token in protected routes:"
echo "   curl -H 'Authorization: Bearer YOUR_SESSION_TOKEN' \\"
echo "   -H 'X-User-Phone: YOUR_PHONE' \\"
echo "   http://localhost:3000/protected" 