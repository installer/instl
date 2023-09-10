test-win:
	go run . -test -verbose > ./run.ps1
	powershell -executionpolicy bypass -File .\run.ps1
	rm .\run.ps1

test-linux:
	go run . -test -verbose > ./run.sh
	cat ./run.sh | bash
	rm ./run.sh
