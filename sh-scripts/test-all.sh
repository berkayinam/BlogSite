#!/bin/bash

# Renk kodları
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "🚀 Full API Test Suite"
echo "--------------------"

# Base URLs
AUTH_URL="http://localhost:8085"
POST_URL="http://localhost:8081"

# 1. Test Auth Service
echo -e "\n${GREEN}1. Testing Auth Service...${NC}"
echo -e "\n${GREEN}1.1. Registering a new user...${NC}"
REGISTER_RESPONSE=$(curl -s -X POST "$AUTH_URL/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass123"
  }')
echo "Register Response: $REGISTER_RESPONSE"

echo -e "\n${GREEN}1.2. Logging in...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$AUTH_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass123"
  }')
echo "Login Response: $LOGIN_RESPONSE"

# Extract token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | grep -o '[^"]*$')

if [ -z "$TOKEN" ]; then
    echo -e "\n${RED}Failed to get token! Exiting...${NC}"
    exit 1
fi

echo -e "\n${GREEN}Successfully got token!${NC}"
echo "Token: $TOKEN"

# 2. Test Post Service
echo -e "\n${GREEN}2. Testing Post Service...${NC}"

echo -e "\n${GREEN}2.1. Creating a new post...${NC}"
CREATE_RESPONSE=$(curl -s -X POST "$POST_URL/posts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Test Başlık",
    "content": "Test İçerik",
    "author": "testuser"
  }')
echo "Create Response: $CREATE_RESPONSE"

sleep 1

echo -e "\n${GREEN}2.2. Listing all posts...${NC}"
LIST_RESPONSE=$(curl -s -X GET "$POST_URL/posts")
echo "List Response: $LIST_RESPONSE"

echo -e "\n${GREEN}2.3. Getting posts by author...${NC}"
AUTHOR_RESPONSE=$(curl -s -X GET "$POST_URL/posts/author?author=testuser")
echo "Author Posts Response: $AUTHOR_RESPONSE"

echo -e "\n${GREEN}2.4. Updating the post...${NC}"
UPDATE_RESPONSE=$(curl -s -X PUT "$POST_URL/posts/manage?title=Test%20Başlık" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Test Başlık",
    "content": "Güncellenmiş İçerik",
    "author": "testuser"
  }')
echo "Update Response: $UPDATE_RESPONSE"

echo -e "\n${GREEN}2.5. Verifying the update...${NC}"
VERIFY_RESPONSE=$(curl -s -X GET "$POST_URL/posts/author?author=testuser")
echo "Verify Response: $VERIFY_RESPONSE"

echo -e "\n${GREEN}2.6. Deleting the post...${NC}"
DELETE_RESPONSE=$(curl -s -X DELETE "$POST_URL/posts/manage?title=Test%20Başlık" \
  -H "Authorization: Bearer $TOKEN")
echo "Delete Response: $DELETE_RESPONSE"

echo -e "\n${GREEN}2.7. Verifying deletion...${NC}"
FINAL_RESPONSE=$(curl -s -X GET "$POST_URL/posts/author?author=testuser")
echo "Final Response: $FINAL_RESPONSE"

echo -e "\n${GREEN}Test suite completed! 🎉${NC}" 