test-windows:
	@echo "# Testing windows installation..."
	@echo "## Deleting old binary..."
	@rm -f "$env:USERPROFILE\instl\instl-demo"
	@echo "## Running instl..."
	@go run . -test -verbose > ./run.ps1
	@powershell -executionpolicy bypass -File ./run.ps1
	@echo
	@echo
	@echo "## Testing binary..."
	@echo
	@start "$env:USERPROFILE\instl\instl-demo\instl-demo.exe"

test:
	@echo "# Testing linux installation..."
	@echo "## Deleting old binary..."
	@rm -f "$HOME/.local/bin/instl-demo"
	@rm -rf "$HOME/.local/bin/.instl/instl-demo/instl-demo"
	@echo "## Running instl..."
	@go run . -test -verbose > ./run.sh
	@cat ./run.sh | bash
	@echo
	@echo
	@echo "## Testing binary..."
	@echo
	@instl-demo

build:
	docker build -t marvinjwendt/instl .

run: build
	docker run -it --rm -p "80:80" marvinjwendt/instl

publish: build
	docker push marvinjwendt/instl
