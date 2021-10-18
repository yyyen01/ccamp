export tag=v1.0
root:
	export ROOT=ccamp/lesson2

build:
	echo "building lesson2 binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release: build
	echo "building lesson2 container"
	docker build -t ccamp/lesson2:${tag} .
	echo "built done"

push: release
	echo "pushing ccamp/lesson2"
	docker tag ccamp/lesson2:v1.0 yyyen01/lesson2:v1.0
	echo "finish tagging"
	#docker login -u yyyen01 -p S4s0t9000!
	docker push yyyen01/lesson2:v1.0