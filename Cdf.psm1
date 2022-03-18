function Invoke-Cdf {
    param(
        [Parameter(ValueFromPipeline=$true)] [string] $path
    )

    if (-not $path) {
        return
    }
    $f = New-TemporaryFile
    & (Join-Path $PSScriptRoot "cdf.exe") "$path" -f "$($f.FullName)"
    $target = Get-Content "$($f.FullName)"
    Remove-Item $f.FullName | Out-Null
    Set-Location $target
}

Export-ModuleMember "Invoke-Cdf"
