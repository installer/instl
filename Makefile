test-windows:
	@echo "Testing if binary exists in Windows..."
	@go run . -test -verbose > ./run.ps1
	@powershell -executionpolicy bypass -File ./run.ps1
	@powershell -executionpolicy bypass -Command "if (Test-Path $$env:USERPROFILE\instl\instl-demo\instl-demo.exe) { echo 'Binary found in Windows'; exit 0 } else { echo 'Binary not found in Windows'; exit 1 }"

test-linux:
	@echo "Testing if binary exists in Linux..."
	@go run . -test -verbose > ./run.sh
	@cat ./run.sh | bash
	@(test -f $$HOME/.local/bin/instl-demo && echo "Binary found in Linux" && exit 0) || (echo "Binary not found in Linux" && exit 1)

test-macos:
	@echo "Testing if binary exists in macOS..."
	@go run . -test -verbose > ./run.sh
	@cat ./run.sh | bash
	@(test -f $$HOME/.local/bin/instl-demo && echo "Binary found in macOS" && exit 0) || (echo "Binary not found in macOS" && exit 1)

build:
	docker build -t marvinjwendt/instl .

run: build
	docker run -it --rm -p "80:80" marvinjwendt/instl

publish: build
	docker push marvinjwendt/instl
