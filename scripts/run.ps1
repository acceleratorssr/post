$servicesDirectory = "../output"
$errorLog = @()
$executionOrder = @("sso.exe", "user.exe", "article.exe", "interactive.exe", "search.exe", "bff.exe")

foreach ($service in $executionOrder) {
    $exeFile = Join-Path -Path $servicesDirectory -ChildPath $service

    if (Test-Path $exeFile) {
        Write-Host "Starting process for: $exeFile"

        $process = Start-Process -FilePath $exeFile -PassThru -ErrorAction Stop
        Start-Sleep -Seconds 1 # 等1s 没死就当它正常启动了

        if (-not $process.HasExited) {
            Write-Host "$exeFile is running successfully."
        } else {
            $errorMsg = "$exeFile exited unexpectedly."
            Write-Host $errorMsg
            $errorLog += $errorMsg
        }
    } else {
        Write-Host "Executable not found: $exeFile"
        $errorLog += "Executable not found: $exeFile"
    }
}

if ($errorLog.Count -gt 0) {
    Write-Host ">! !-> Errors encountered during execution:"
    $errorLog | ForEach-Object { Write-Host $_ }
} else {
    Write-Host "All processes started successfully."
}
