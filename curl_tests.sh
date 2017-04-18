echo "post article"
curl -X POST \
	-H 'Content-Type: application/json' \
	-d '{
    "text": "Hello world",
    "created_by": "ckuroki@gmail.com"
    }' localhost:8080/api/articles
echo
echo "post comment"
curl -X POST \
	-H 'Content-Type: application/json' \
	-d '{
    "article_id": 1,
    "text": "Nice article!",
    "created_by": "niceguy@mail.com"
    }' localhost:8080/api/comments
echo
echo "post comment"
curl -X POST \
	-H 'Content-Type: application/json' \
	-d '{
    "article_id": 1,
    "text": "Very interesting...",
    "created_by": "interested@gmail.com"
    }' localhost:8080/api/comments
echo
echo "post comment"
curl -X POST \
	-H 'Content-Type: application/json' \
	-d '{
    "article_id": 1,
    "text": "Worst article ever!",
    "created_by": "pesismist@worst.com"
    }' localhost:8080/api/comments
echo
echo "put comment"
curl -X PUT \
	-H 'Content-Type: application/json' \
	-d '{
    "text": "Best article ever!"
    }' localhost:8080/api/comments/3
echo
echo "get comment 1"
curl localhost:8080/api/comments/1
echo
echo "get comment 2"
curl localhost:8080/api/comments/4
echo
echo "get article comments"
curl localhost:8080/api/articles/1/comments
echo
echo "get article list"
curl localhost:8080/api/articles
echo

