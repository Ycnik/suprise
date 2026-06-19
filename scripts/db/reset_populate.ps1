param(
    [string]$Container = "postgres",
    [string]$Database = "soldat",
    [string]$User = "soldat",
    [string]$Password = "p"
)

$ErrorActionPreference = "Stop"

$scriptPath = Join-Path $PSScriptRoot "reset_populate.sql"
Get-Content -Raw -Path $scriptPath | docker exec -i -e PGPASSWORD=$Password $Container psql -U $User -d $Database
