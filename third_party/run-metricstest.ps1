function Get-Unused-Port
{
    $usedPorts = (Get-NetTCPConnection | select -ExpandProperty LocalPort) + (Get-NetUDPEndpoint | select -ExpandProperty LocalPort)
    5000..60000 | where { $usedPorts -notcontains $_ } | select -first 1
}

cd $PSScriptRoot\..

cd cmd\agent
echo "Building agent.exe"
go build -buildvcs=false -o agent.exe
echo "Done"

cd ../server
echo "Building server.exe"
go build -buildvcs=false -o server.exe
echo "Done"

cd ../..
echo "Running tests"
$b = git rev-parse --abbrev-ref HEAD
$iter = ($b -replace 'iter').Trim()

if ($iter -ge 1)
{
    metricstest '-test.v' '-test.run=^TestIteration1$' '-binary-path=cmd\server\server.exe'
}
if ($iter -ge 2)
{
    metricstest '-test.v' '-test.run=^TestIteration2[AB]*$' '-agent-binary-path=cmd\agent\agent.exe' '-source-path=.'
}
if ($iter -ge 3)
{
    metricstest '-test.v' '-test.run=^TestIteration3[AB]*$' '-binary-path=cmd\server\server.exe' '-agent-binary-path=cmd\agent\agent.exe' '-source-path=.'
}
if ($iter -ge 4)
{
    $port = Get-Unused-Port
    $address = "localhost:$port"
    metricstest '-test.v' '-test.run=^TestIteration4$' '-binary-path=cmd\server\server.exe' '-agent-binary-path=cmd\agent\agent.exe' "-server-port=$port" '-source-path=.'
}
if ($iter -ge 5)
{
    powershell -Command {
        $port = '15243'
        $address = "localhost:$port"
        $env:ADDRESS = $address
        metricstest '-test.v' '-test.run=^TestIteration5$' '-binary-path=cmd\server\server.exe' '-agent-binary-path=cmd\agent\agent.exe' "-server-port=$port" '-source-path=.'
    }
}

cd $PSScriptRoot
