# Tüm testleri çalıştırmak için
go test ./internal/... -v

# Belirli bir testi çalıştırmak için
go test ./internal/... -v -run TestGetAllPosts

# Test coverage'ı görmek için
go test ./internal/... -cover


Post Service Test Documentation
Prerequisites
Go 1.x installed
MongoDB ning locally or accessible via network
Access to the project repository
Test Environment Setup
First, make sure MongoDB is ning. You can  it locally or use a containerized version:


Navigate to the post-service directory:


Install dependencies:


ning the Tests
 all tests


 specific test files


 individual test functions


Test Configuration
The tests use environment variables for configuration:
TEST_MONGO_URI: MongoDB connection string (default: "mongodb://localhost:27017")
Example with custom MongoDB URI:


Test Coverage
To  tests with coverage:


Test Suite Overview
Repository Tests (post_repository_test.go)
TestCreatePost
Tests post creation functionality
Verifies proper storage in MongoDB
Checks error handling
TestGetAllPosts
Tests retrieval of all posts
Verifies correct number of posts returned
Checks sorting and filtering
TestGetPostsByAuthor
Tests filtering posts by author
Verifies author-specific post retrieval
Checks error handling for invalid authors
TestUpdatePost
Tests post update functionality
Verifies content modification
Checks error handling for non-existent posts
TestDeletePost
Tests post deletion
Verifies complete removal from database
Checks error handling
Handler Tests (handlers_test.go)
TestCreatePostHandler
Tests HTTP POST endpoint
Verifies authentication middleware
Checks response codes and error handling
TestListPostsHandler
Tests HTTP GET endpoint for all posts
Verifies response format
Checks pagination and sorting
TestGetPostsByAuthorHandler
Tests HTTP GET endpoint for author-specific posts
Verifies query parameter handling
Checks response format and error cases
TestUpdatePostHandler
Tests HTTP PUT endpoint
Verifies authentication requirements
Checks response codes and error handling
TestDeletePostHandler
Tests HTTP DELETE endpoint
Verifies authentication requirements
Checks response codes and success cases
Common Test Scenarios
Authentication Tests
All protected endpoints require valid authentication
Tests include both authenticated and unauthenticated requests
Verifies proper handling of invalid authentication
Error Handling Tests
Invalid input data
Missing required fields
Database connection issues
Duplicate entries
Non-existent resources
Data Validation Tests
Input sanitization
Required field validation
Data type validation
Length and format constraints
Troubleshooting
Common issues and solutions:
MongoDB Connection Issues

Solution: Verify MongoDB is running and accessible at the specified URI
Test Timeouts

Solution: Increase test timeout duration:


Clean Up Issues
If tests fail due to leftover data, manually clean the test database:


Best Practices
Always  tests in a clean environment
Use proper setup and teardown functions
Don't rely on test execution order
Clean up resources after tests
Use meaningful test names and descriptions
Include both positive and negative test cases
Contributing New Tests
When adding new tests:
Follow existing test patterns
Include setup and teardown
Add proper documentation
Test both success and failure cases
Verify edge cases
Include performance considerations
This documentation should help you get started with ning and understanding the tests in the post service. Let me know if you need any clarification or have questions about specific test cases!