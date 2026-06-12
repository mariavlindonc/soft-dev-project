Get-Content ../backend/.env | ForEach-Object {
    if ($_ -match '(.+)=(.+)') {
        [Environment]::SetEnvironmentVariable($matches[1], $matches[2])
    }
}

Write-Host "Creating schema..."
cmd /c "mysql -h %DB_HOST% -P %DB_PORT% -u %DB_USER% -p%DB_PASSWORD% %DB_NAME% --ssl-ca=../backend/certs/ca.pem < schema.sql"

Write-Host "Seeding data..."
cmd /c "mysql -h %DB_HOST% -P %DB_PORT% -u %DB_USER% -p%DB_PASSWORD% %DB_NAME% --ssl-ca=../backend/certs/ca.pem < data.sql"

Write-Host "Database initialized!"