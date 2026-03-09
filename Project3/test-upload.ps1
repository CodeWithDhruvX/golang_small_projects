# Test document upload API
$token = (Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -ContentType "application/json" -Body '{"username":"test","password":"test"}').token
$headers = @{Authorization = "Bearer $token"}

# Read file content
$filePath = "sample-document.txt"
$fileContent = Get-Content $filePath -Raw
$fileName = Split-Path $filePath -Leaf

# Create multipart form data
$boundary = "----WebKitFormBoundary$(Get-Random)"
$body = @"

--$boundary
Content-Disposition: form-data; name="file"; filename="$fileName"
Content-Type: text/plain

$fileContent

--$boundary--
"@

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/documents/upload" -Method POST -ContentType "multipart/form-data; boundary=$boundary" -Headers $headers -Body $body
    Write-Host "Upload successful:" -ForegroundColor Green
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "Upload failed:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $errorBody = $reader.ReadToEnd()
        Write-Host "Error response: $errorBody" -ForegroundColor Red
    }
}
