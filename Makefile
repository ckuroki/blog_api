all: blog_api

blog_api: blog_api.go 
	go build blog_api.go

