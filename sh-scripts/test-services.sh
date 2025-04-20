#!/bin/bash

# Import common variables and functions
source "$(dirname "$0")/common.sh"

# Test function
test_endpoint() {
    local method=$1
    local endpoint=$2
    local payload=$3
    local expected_status=$4
    local service_name=$5
    local port=$6

    echo -e "\nTesting ${service_name} - ${method} ${endpoint}"
    
    if [ -z "$payload" ]; then
        response=$(curl -s -w "\n%{http_code}" -X ${method} \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer ${AUTH_TOKEN}" \
            http://localhost:${port}${endpoint})
    else
        response=$(curl -s -w "\n%{http_code}" -X ${method} \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer ${AUTH_TOKEN}" \
            -d "${payload}" \
            http://localhost:${port}${endpoint})
    fi

    status_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | sed '$d')

    if [ "$status_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✓ Success${NC} (Status: ${status_code})"
        echo "Response: ${response_body}"
        echo "$response_body" | jq -r '.id // empty'
    else
        echo -e "${RED}✗ Failed${NC} (Expected: ${expected_status}, Got: ${status_code})"
        echo "Response: ${response_body}"
        return 1
    fi
}

# Get auth token
echo "Getting authentication token..."
AUTH_TOKEN=$(get_auth_token)
if [ -z "$AUTH_TOKEN" ]; then
    echo "Failed to get authentication token"
    exit 1
fi

echo "Starting API Tests..."

# Test Team Service
TEAM_PORT=8084
echo -e "\n=== Testing Team Service ==="

# Create team
TEAM_ID=$(test_endpoint "POST" "/teams" '{"name":"Test Team","description":"Test Description"}' 200 "Team Service" $TEAM_PORT)
if [ $? -ne 0 ]; then
    echo "Failed to create team"
    exit 1
fi

# Get team
test_endpoint "GET" "/teams/${TEAM_ID}" "" 200 "Team Service" $TEAM_PORT

# Get user teams
test_endpoint "GET" "/teams/user" "" 200 "Team Service" $TEAM_PORT

# Update team
test_endpoint "PUT" "/teams/${TEAM_ID}" '{"name":"Updated Team","description":"Updated Description"}' 200 "Team Service" $TEAM_PORT

# Invite member
INVITE_RESPONSE=$(test_endpoint "POST" "/teams/${TEAM_ID}/invite" "{\"username\":\"${TEST_USER2}\"}" 200 "Team Service" $TEAM_PORT)
INVITE_ID=$(echo "$INVITE_RESPONSE" | jq -r '.id')

# Switch to test user 2
AUTH_TOKEN=$(get_auth_token_for_user "${TEST_USER2}" "${TEST_PASS2}")

# Respond to invite
test_endpoint "POST" "/teams/invites/${INVITE_ID}/respond" '{"status":"accepted"}' 204 "Team Service" $TEAM_PORT

# Switch back to test user 1
AUTH_TOKEN=$(get_auth_token)

# Remove member
test_endpoint "DELETE" "/teams/${TEAM_ID}/members/${TEST_USER2}" "" 204 "Team Service" $TEAM_PORT

# Delete team
test_endpoint "DELETE" "/teams/${TEAM_ID}" "" 204 "Team Service" $TEAM_PORT

# Test Comment Service
COMMENT_PORT=8083
echo -e "\n=== Testing Comment Service ==="

# Create post ID for testing
POST_ID="6432a7d8f5c88e1234567892"

# Create comment
COMMENT_ID=$(test_endpoint "POST" "/posts/${POST_ID}/comments" '{"content":"Test comment"}' 200 "Comment Service" $COMMENT_PORT)

# Get comments for post
test_endpoint "GET" "/posts/${POST_ID}/comments" "" 200 "Comment Service" $COMMENT_PORT

# Update comment
test_endpoint "PUT" "/comments/${COMMENT_ID}" '{"content":"Updated comment"}' 200 "Comment Service" $COMMENT_PORT

# Like comment
test_endpoint "POST" "/comments/${COMMENT_ID}/like" "" 200 "Comment Service" $COMMENT_PORT

# Unlike comment
test_endpoint "POST" "/comments/${COMMENT_ID}/unlike" "" 200 "Comment Service" $COMMENT_PORT

# Like post
test_endpoint "POST" "/posts/${POST_ID}/like" "" 200 "Comment Service" $COMMENT_PORT

# Get post likes
test_endpoint "GET" "/posts/${POST_ID}/likes" "" 200 "Comment Service" $COMMENT_PORT

# Unlike post
test_endpoint "POST" "/posts/${POST_ID}/unlike" "" 200 "Comment Service" $COMMENT_PORT

# Send friend request to test user 2
test_endpoint "POST" "/friends/request" "{\"user2\":\"${TEST_USER2}\"}" 200 "Comment Service" $COMMENT_PORT

# Switch to test user 2
AUTH_TOKEN=$(get_auth_token_for_user "${TEST_USER2}" "${TEST_PASS2}")

# Get friend requests
REQUEST_RESPONSE=$(test_endpoint "GET" "/friends/requests" "" 200 "Comment Service" $COMMENT_PORT)
REQUEST_ID=$(echo "$REQUEST_RESPONSE" | jq -r '.[0].id')

# Accept friend request
test_endpoint "POST" "/friends/requests/${REQUEST_ID}/accept" "" 200 "Comment Service" $COMMENT_PORT

# Get friends list
test_endpoint "GET" "/friends" "" 200 "Comment Service" $COMMENT_PORT

# Switch back to test user 1
AUTH_TOKEN=$(get_auth_token)

# Delete comment
test_endpoint "DELETE" "/comments/${COMMENT_ID}" "" 200 "Comment Service" $COMMENT_PORT

echo -e "\n${GREEN}API Tests Completed Successfully!${NC}" 