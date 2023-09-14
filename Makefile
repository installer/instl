test-win:
	go run . -test -verbose > ./run.ps1
	powershell -executionpolicy bypass -File .\run.ps1

test-linux:
	go run . -test -verbose > ./run.sh
	cat ./run.sh | bash

build-image:
	docker build -t marvinjwendt/instl .

publish: build-image
	docker push marvinjwendt/instl
