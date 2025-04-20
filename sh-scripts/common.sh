#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# Test users
TEST_USER="testuser"
TEST_PASS="testpass123"
TEST_USER2="testuser2"
TEST_PASS2="testpass123"

# Auth service port
AUTH_PORT=8081

# Function to get auth token
get_auth_token() {
    local username=${1:-$TEST_USER}
    local password=${2:-$TEST_PASS}
    
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"${username}\",\"password\":\"${password}\"}" \
        http://localhost:${AUTH_PORT}/login)
    
    echo "$response" | jq -r '.token'
}

# Function to get auth token for a specific user
get_auth_token_for_user() {
    local username=$1
    local password=$2
    get_auth_token "$username" "$password"
}

# Function to create test users if they don't exist
create_test_users() {
    # Create first test user
    curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"${TEST_USER}\",\"password\":\"${TEST_PASS}\",\"email\":\"test@example.com\"}" \
        http://localhost:${AUTH_PORT}/register > /dev/null

    # Create second test user
    curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"${TEST_USER2}\",\"password\":\"${TEST_PASS2}\",\"email\":\"test2@example.com\"}" \
        http://localhost:${AUTH_PORT}/register > /dev/null
}

# Create test users at script start
create_test_users 