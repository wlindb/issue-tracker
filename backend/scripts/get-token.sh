#!/bin/sh
curl -s -X POST "http://localhost:8180/realms/issue-tracker/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=issue-tracker-app" \
  -d "client_secret=test-secret" \
  -d "username=testuser" \
  -d "password=password" \
  | grep -o '"access_token":"[^"]*"' \
  | cut -d'"' -f4
