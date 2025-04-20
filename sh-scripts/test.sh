#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Test function
test_endpoint() {
    local method=$1
    local endpoint=$2
    local payload=$3
    local auth_token=$4
    local expected_status=$5
    local service_name=$6

    echo -e "\nTesting ${service_name} - ${method} ${endpoint}"
    
    if [ -z "$payload" ]; then
        response=$(curl -s -w "\n%{http_code}" -X ${method} \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer ${auth_token}" \
            http://localhost:${port}${endpoint})
    else
        response=$(curl -s -w "\n%{http_code}" -X ${method} \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer ${auth_token}" \
            -d "${payload}" \
            http://localhost:${port}${endpoint})
    fi

    status_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | sed '$d')

    if [ "$status_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✓ Success${NC} (Status: ${status_code})"
        echo "Response: ${response_body}"
    else
        echo -e "${RED}✗ Failed${NC} (Expected: ${expected_status}, Got: ${status_code})"
        echo "Response: ${response_body}"
    fi
}

# JWT token for testing (replace with a valid token)
AUTH_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIn0.J2mRHPH6YwHIjx7kvsKrHXWkn9LBQGCg1RqMEZ4i-Ok"

echo "Starting API Tests..."

# Test Team Service
port=8084
echo -e "\n=== Testing Team Service ==="

# Create team
test_endpoint "POST" "/teams" '{"name":"Test Team","description":"Test Description"}' "$AUTH_TOKEN" 200 "Team Service"

# Get the created team ID from the response
TEAM_ID="6432a7d8f5c88e1234567890"  # Replace with actual ID from create response

# Get team
test_endpoint "GET" "/teams/${TEAM_ID}" "" "$AUTH_TOKEN" 200 "Team Service"

# Get user teams
test_endpoint "GET" "/teams/user" "" "$AUTH_TOKEN" 200 "Team Service"

# Update team
test_endpoint "PUT" "/teams/${TEAM_ID}" '{"name":"Updated Team","description":"Updated Description"}' "$AUTH_TOKEN" 200 "Team Service"

# Invite member
test_endpoint "POST" "/teams/${TEAM_ID}/invite" '{"username":"newuser"}' "$AUTH_TOKEN" 200 "Team Service"

# Get the invite ID from the response
INVITE_ID="6432a7d8f5c88e1234567891"  # Replace with actual ID from invite response

# Respond to invite
test_endpoint "POST" "/teams/invites/${INVITE_ID}/respond" '{"status":"accepted"}' "$AUTH_TOKEN" 204 "Team Service"

# Remove member
test_endpoint "DELETE" "/teams/${TEAM_ID}/members/newuser" "" "$AUTH_TOKEN" 204 "Team Service"

# Delete team
test_endpoint "DELETE" "/teams/${TEAM_ID}" "" "$AUTH_TOKEN" 204 "Team Service"

# Test Comment Service
port=8083
echo -e "\n=== Testing Comment Service ==="

# Create comment
POST_ID="6432a7d8f5c88e1234567892"  # Replace with actual post ID
test_endpoint "POST" "/posts/${POST_ID}/comments" '{"content":"Test comment"}' "$AUTH_TOKEN" 200 "Comment Service"

# Get comments for post
test_endpoint "GET" "/posts/${POST_ID}/comments" "" "$AUTH_TOKEN" 200 "Comment Service"

# Get the comment ID from the response
COMMENT_ID="6432a7d8f5c88e1234567893"  # Replace with actual ID from create response

# Update comment
test_endpoint "PUT" "/comments/${COMMENT_ID}" '{"content":"Updated comment"}' "$AUTH_TOKEN" 200 "Comment Service"

# Like comment
test_endpoint "POST" "/comments/${COMMENT_ID}/like" "" "$AUTH_TOKEN" 200 "Comment Service"

# Unlike comment
test_endpoint "POST" "/comments/${COMMENT_ID}/unlike" "" "$AUTH_TOKEN" 200 "Comment Service"

# Like post
test_endpoint "POST" "/posts/${POST_ID}/like" "" "$AUTH_TOKEN" 200 "Comment Service"

# Get post likes
test_endpoint "GET" "/posts/${POST_ID}/likes" "" "$AUTH_TOKEN" 200 "Comment Service"

# Unlike post
test_endpoint "POST" "/posts/${POST_ID}/unlike" "" "$AUTH_TOKEN" 200 "Comment Service"

# Send friend request
test_endpoint "POST" "/friends/request" '{"user2":"frienduser"}' "$AUTH_TOKEN" 200 "Comment Service"

# Get friend requests
test_endpoint "GET" "/friends/requests" "" "$AUTH_TOKEN" 200 "Comment Service"

# Get the friend request ID from the response
REQUEST_ID="6432a7d8f5c88e1234567894"  # Replace with actual ID from request response

# Accept friend request
test_endpoint "POST" "/friends/requests/${REQUEST_ID}/accept" "" "$AUTH_TOKEN" 200 "Comment Service"

# Get friends list
test_endpoint "GET" "/friends" "" "$AUTH_TOKEN" 200 "Comment Service"

# Delete comment
test_endpoint "DELETE" "/comments/${COMMENT_ID}" "" "$AUTH_TOKEN" 200 "Comment Service"

echo -e "\nAPI Tests Completed!" 