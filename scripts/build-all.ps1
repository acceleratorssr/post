$outputDir = "../output"

if (-Not (Test-Path -Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir
}

$services = @("..\\article", "interactive", "search", "sso", "user", "bff", "recommend")

foreach ($service in $services) {
    Write-Host "Building $service..."
    Set-Location $service
    $result = go build -o $outputDir .

    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error building $service"
        exit 1
    }
    Set-Location ..
}

Set-Location "scripts"

Write-Host "All services built successfully!"
