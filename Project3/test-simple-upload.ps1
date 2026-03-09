# Test simple document upload
$token = (Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -ContentType "application/json" -Body '{"username":"test","password":"test"}').token
$headers = @{Authorization = "Bearer $token"}

# Read file content
$filePath = "test-simple.txt"
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
    Write-Host "Document ID: $($response.document_id)"
    Write-Host "Status: $($response.status)"
    Write-Host "Message: $($response.message)"
    
    # Now test chat with this document
    $chatBody = @{
        session_id = "0878bbe8-365c-4c25-a0aa-289bad539757"
        message = "What is supervised learning?"
        model = "llama3.1:8b"
    } | ConvertTo-Json
    
    Write-Host "`nTesting chat with RAG..." -ForegroundColor Yellow
    $chatResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/chat/" -Method POST -ContentType "application/json" -Headers $headers -Body $chatBody
    Write-Host "Chat Response:" -ForegroundColor Cyan
    $chatResponse.message
} catch {
    Write-Host "Error:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $errorBody = $reader.ReadToEnd()
        Write-Host "Error response: $errorBody" -ForegroundColor Red
    }
}
