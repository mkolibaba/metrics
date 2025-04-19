cd $PSScriptRoot\..

cd cmd\agent
echo "Building agent.exe"
go build -o agent.exe .\main.go
echo "Done"

cd ../server
echo "Building server.exe"
go build -o server.exe .\main.go
echo "Done"

cd ../..
echo "Running tests"
metricstest '-test.v' '-test.run=^TestIteration5$' '-binary-path=cmd\server\server.exe' '-agent-binary-path=cmd\agent\agent.exe' '-server-port=15243' '-source-path=.'

cd $PSScriptRoot