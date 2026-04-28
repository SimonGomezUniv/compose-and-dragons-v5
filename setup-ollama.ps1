#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Setup and run Ollama in Docker for the game engine
.DESCRIPTION
    Starts Ollama in a Docker container on port 11434
    Then pulls recommended small models for game use
.NOTES
    All models run locally, no internet required after download
#>

$ErrorActionPreference = "Stop"

Write-Host "`n" -NoNewline
Write-Host "═══════════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "🦙 Ollama in Docker Setup (Windows)" -ForegroundColor Cyan
Write-Host "═══════════════════════════════════════════════════════════════" -ForegroundColor Cyan

# ============================================================================
# STEP 1: Check Docker
# ============================================================================
Write-Host "`n[1/5] Checking Docker Desktop..." -ForegroundColor Yellow
try {
    $dockerVersion = docker --version
    Write-Host "✅ Docker: $dockerVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker is not installed or not in PATH" -ForegroundColor Red
    Write-Host "   Download: https://www.docker.com/products/docker-desktop" -ForegroundColor Yellow
    exit 1
}

# ============================================================================
# STEP 2: Stop any previous Ollama container
# ============================================================================
Write-Host "`n[2/5] Cleaning up previous containers..." -ForegroundColor Yellow
docker compose -f ollama-compose.yml down 2>&1 | Out-Null
Start-Sleep -Seconds 2
Write-Host "✅ Cleaned" -ForegroundColor Green

# ============================================================================
# STEP 3: Start Ollama
# ============================================================================
Write-Host "`n[3/5] Starting Ollama container..." -ForegroundColor Yellow
try {
    docker compose -f ollama-compose.yml up -d
    Write-Host "✅ Container started" -ForegroundColor Green
} catch {
    Write-Host "❌ Failed to start container" -ForegroundColor Red
    Write-Host $_ -ForegroundColor Red
    exit 1
}

# ============================================================================
# STEP 4: Wait for Ollama to be ready
# ============================================================================
Write-Host "`n[4/5] Waiting for Ollama to initialize (60s max)..." -ForegroundColor Yellow
$ready = $false
$maxAttempts = 60

for ($i = 1; $i -le $maxAttempts; $i++) {
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:11434/api/tags" -ErrorAction SilentlyContinue
        $ready = $true
        break
    } catch {
        Write-Host "." -NoNewline -ForegroundColor Gray
        Start-Sleep -Seconds 1
    }
}

Write-Host ""
if ($ready) {
    Write-Host "✅ Ollama is ready" -ForegroundColor Green
} else {
    Write-Host "⚠️  Ollama is starting... (may take 1-2 minutes)" -ForegroundColor Yellow
}

# ============================================================================
# STEP 5: Pull recommended models
# ============================================================================
Write-Host "`n[5/5] Model selection..." -ForegroundColor Yellow

$models = @(
    @{ name = "mistral"; size = "4.1GB"; desc = "Mistral 7B - Fast & Smart" },
    @{ name = "neural-chat"; size = "3.8GB"; desc = "Neural Chat 7B - Good for dialogue" },
    @{ name = "orca-mini"; size = "1.3GB"; desc = "Orca Mini - Very small & fast" }
)

Write-Host "`nAvailable models for game use:" -ForegroundColor Cyan
for ($i = 0; $i -lt $models.Count; $i++) {
    Write-Host "  $($i+1). $($models[$i].desc) ($($models[$i].size))" -ForegroundColor White
}

$choice = Read-Host "`nSelect model (1-3, or 0 to skip)"

if ($choice -ge 1 -and $choice -le 3) {
    $selectedModel = $models[$choice - 1]
    Write-Host "`nPulling $($selectedModel.name) (this may take a few minutes)..." -ForegroundColor Yellow
    Write-Host "Size: $($selectedModel.size)" -ForegroundColor Gray
    
    docker exec ollama-server ollama pull $selectedModel.name
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Model pulled successfully" -ForegroundColor Green
    } else {
        Write-Host "⚠️  Model pull failed (but Ollama is still running)" -ForegroundColor Yellow
    }
} else {
    Write-Host "⏭️  Skipping model download" -ForegroundColor Yellow
}

# ============================================================================
# Summary
# ============================================================================
Write-Host "`n" -NoNewline
Write-Host "═══════════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "✅ Ollama is READY!" -ForegroundColor Green
Write-Host "═══════════════════════════════════════════════════════════════" -ForegroundColor Cyan

Write-Host "`n📍 Configuration for your Go code:" -ForegroundColor Cyan
Write-Host "`nChange your engine URL to:" -ForegroundColor White
Write-Host "`n   baseURL := ""http://localhost:11434""" -ForegroundColor Yellow
Write-Host "`n   or" -ForegroundColor Gray
Write-Host "`n   baseURL := ""http://localhost:11434/api/generate""" -ForegroundColor Yellow

Write-Host "`n🌐 Test endpoints:" -ForegroundColor Cyan
Write-Host "`n   List models:" -ForegroundColor Gray
Write-Host "   curl http://localhost:11434/api/tags" -ForegroundColor White
Write-Host "`n   Test completion:" -ForegroundColor Gray
Write-Host "   curl -X POST http://localhost:11434/api/generate \" -ForegroundColor White
Write-Host "     -H 'Content-Type: application/json' \" -ForegroundColor White
Write-Host "     -d '{\"model\": \"mistral\", \"prompt\": \"Hello\"}'" -ForegroundColor White

Write-Host "`n📋 Useful commands:" -ForegroundColor Cyan
Write-Host "`n   See logs:" -ForegroundColor Gray
Write-Host "   docker compose -f ollama-compose.yml logs -f" -ForegroundColor White
Write-Host "`n   Check models:" -ForegroundColor Gray
Write-Host "   docker exec ollama-server ollama list" -ForegroundColor White
Write-Host "`n   Pull another model:" -ForegroundColor Gray
Write-Host "   docker exec ollama-server ollama pull <model-name>" -ForegroundColor White
Write-Host "`n   Stop Ollama:" -ForegroundColor Gray
Write-Host "   docker compose -f ollama-compose.yml down" -ForegroundColor White

Write-Host "`n💾 Models storage:" -ForegroundColor Cyan
Write-Host "   Location: Docker volume 'ollama-models'" -ForegroundColor Gray
Write-Host "   (Persists even if container is deleted)" -ForegroundColor Gray

Write-Host "`n🔗 For Docker-to-Docker communication:" -ForegroundColor Cyan
Write-Host "   base_url: http://ollama:11434" -ForegroundColor Yellow
Write-Host "   (Use service name instead of localhost)" -ForegroundColor Gray

Write-Host "`n"
