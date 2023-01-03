all:
	@go build -o server main.go 

export v= 11

dockerbuild:
	docker build --platform linux/amd64 . -t tayyab7891/sportsapiv1:v$v

dockerpush:
	docker push tayyab7891/sportsapiv1:v$v