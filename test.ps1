go run . -test -verbose > ./run.ps1
"./run.ps1" | iex
rm .\run.ps1
