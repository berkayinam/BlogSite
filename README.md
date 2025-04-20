Medium.com dan sıkıldıgım için kendi blog sitemi yapıyorum

minikube ortamına gidiş: eval $(minikube docker-env)

:dönüş eval $(minikube docker-env -u)

kubectl apply -f auth-deployment.yaml   kubectl delete deployment auth-deployment 

çıkılan konteynerleri temizleme docker rm $(docker ps -a -q -f "status=exited")



Public Endpoints:
GET /teams - List all teams
GET /teams/:id - Get team details

Protected Endpoints (requires JWT token):
POST /teams - Create a new team
POST /teams/invite - Invite a user to team
POST /teams/invite/respond - Accept/reject team invitation
POST /teams/join/request - Request to join a team

✅ Authentication is working (tokens are being validated)
✅ Team creation is working
✅ Team listing is working
✅ Team membership is working (you were automatically added as admin)
✅ Authorization is working (the invite endpoint correctly checks for admin role)


GET    /posts           - List all posts
POST   /posts           - Create a new post (requires auth)
GET    /posts/author    - Get posts by author (query param: author)
GET    /posts/search    - Search posts (query param: q)
PUT    /posts/manage    - Update a post (query param: title, requires auth)
DELETE /posts/manage    - Delete a post (query param: title, requires auth)




myblog-project/
├── auth-service/
│   ├── internal/
│   │   ├── db.go
│   │   ├── handler.go
│   │   ├── models.go
│   │   └── utils.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── post-service/
│   ├── internal/
│   │   ├── auth.go
│   │   ├── handlers.go
│   │   ├── models.go
│   │   ├── mongo.go
│   │   └── post_repository.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── sh-scripts/
│   ├── clean-db.sh      # MongoDB'yi temizlemek için
│   └── test-all.sh      # Tüm API testleri için
├── kubernetes/          # Kubernetes deployment dosyaları
└── README.md


binam@localhost:~/myblog-project/sh-scripts$ cd /home/binam/myblog-project/sh-scripts && ./test-all.sh
🚀 Full API Test Suite
--------------------

1. Testing Auth Service...

1.1. Registering a new user...
Register Response: Username already exists

1.2. Logging in...
Login Response: {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDU0MTc4MDAsInVzZXJuYW1lIjoidGVzdHVzZXIifQ.4gQu2fEZbxg-NJODNxKAjbiuk_sa26dgUMwumh2X1kg"}

Successfully got token!
Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDU0MTc4MDAsInVzZXJuYW1lIjoidGVzdHVzZXIifQ.4gQu2fEZbxg-NJODNxKAjbiuk_sa26dgUMwumh2X1kg

2. Testing Post Service...

2.1. Creating a new post...
Create Response: {"message":"Post created"}

2.2. Listing all posts...
List Response: [{"title":"Test Başlık","content":"Test İçerik","author":"testuser","createdAt":"2025-04-20T14:16:40.095Z"}]

2.3. Getting posts by author...
Author Posts Response: [{"title":"Test Başlık","content":"Test İçerik","author":"testuser","createdAt":"2025-04-20T14:16:40.095Z"}]

2.4. Updating the post...
Update Response: {"message":"Post updated"}

2.5. Verifying the update...
Verify Response: [{"title":"Test Başlık","content":"Güncellenmiş İçerik","author":"testuser","createdAt":"0001-01-01T00:00:00Z"}]

2.6. Deleting the post...
Delete Response: {"message":"Post deleted"}

2.7. Verifying deletion...
Final Response: null

Test suite completed! 🎉


# Auth service
cd auth-service && JWT_SECRET=supppppersecretkeyimxD go run main.go

# Post service (yeni terminal)
cd post-service && JWT_SECRET=supppppersecretkeyimxD go run main.go

# Veritabanını temizle (opsiyonel)
cd sh-scripts && ./clean-db.sh

# Tüm API testlerini çalıştır
./test-all.sh