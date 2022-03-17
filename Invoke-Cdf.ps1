param(
    [Parameter(Mandatory=$true, ValueFromPipeline=$true)] [string] $path
)

$f = New-TemporaryFile
& (Join-Path $PSScriptRoot "cdf.exe") "$path" -f "$($f.FullName)"
$target = Get-Content "$($f.FullName)"
Remove-Item $f.FullName | Out-Null
Set-Location $target
