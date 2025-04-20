#!/bin/bash

# Renk kodları
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "🧹 Cleaning MongoDB Collections"
echo "----------------------------"

# MongoDB connection string
MONGO_URI="mongodb://localhost:27017"

# Clean auth database
echo -e "\n${GREEN}1. Cleaning auth database...${NC}"
mongosh "$MONGO_URI/authdb" --eval "db.users.drop()"

# Clean posts database
echo -e "\n${GREEN}2. Cleaning posts database...${NC}"
mongosh "$MONGO_URI/blogdb" --eval "db.posts.drop()"

echo -e "\n${GREEN}Database cleanup completed! 🎉${NC}" 