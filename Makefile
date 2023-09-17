test-win:
	go run . -test -verbose > ./run.ps1
	powershell -executionpolicy bypass -File .\run.ps1

test-linux:
	go run . -test -verbose > ./run.sh
	cat ./run.sh | bash

build:
	docker build -t marvinjwendt/instl .

run: build
	docker run -it --rm -p "80:80" marvinjwendt/instl

publish: build
	docker push marvinjwendt/instl
