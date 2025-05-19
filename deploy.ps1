# PowerShell script for deploying User-Backend to GCP

# Variables
$SERVICE_NAME = "default"

# Check if required tools are installed
function Test-CommandExists {
    param ($command)
    $exists = Get-Command $command -ErrorAction SilentlyContinue
    return $exists -ne $null
}

Write-Host "Checking for required tools..." -ForegroundColor Yellow

if (-not (Test-CommandExists "gcloud")) {
    Write-Host "Error: gcloud not found. Please install Google Cloud SDK." -ForegroundColor Red
    exit 1
}

# Get Project ID
$PROJECT_ID = & gcloud config get-value project 2>$null
if ([string]::IsNullOrWhiteSpace($PROJECT_ID)) {
    Write-Host "Error: Could not determine Google Cloud project ID." -ForegroundColor Red
    Write-Host "Please set your project ID with: gcloud config set project YOUR_PROJECT_ID" -ForegroundColor Yellow
    exit 1
}

# Load environment variables from .env file
$ENV_VARS = @{}
if (Test-Path -Path ".env") {
    Write-Host "Loading environment variables from .env file..." -ForegroundColor Green
    
    Get-Content ".env" | ForEach-Object {
        $line = $_
        
        # Skip comments and empty lines
        if (-not [string]::IsNullOrWhiteSpace($line) -and -not $line.StartsWith("#")) {
            $parts = $line.Split("=", 2)
            $var_name = $parts[0].Trim()
            $var_value = $parts[1].Trim()
            
            # Add to environment variables hashtable
            $ENV_VARS[$var_name] = $var_value
        }
    }
    
    Write-Host "Environment variables loaded successfully." -ForegroundColor Green
} else {
    Write-Host "Error: .env file not found in current directory." -ForegroundColor Red
    exit 1
}

# Confirm deployment
Write-Host "You are about to deploy $SERVICE_NAME to Google App Engine." -ForegroundColor Yellow
Write-Host "Project: $PROJECT_ID" -ForegroundColor Yellow
Write-Host "Press ENTER to continue or CTRL+C to cancel..." -ForegroundColor Yellow
Read-Host

Write-Host "Creating temporary app.yaml with environment variables..." -ForegroundColor Green

# Create a temporary app.yaml with environment variables
$appYamlContent = Get-Content "app.yaml" -Raw

# Create the environment variables section
$envVarsSection = "env_variables:`n"
foreach ($envVar in $ENV_VARS.GetEnumerator()) {
    $envVarsSection += "  $($envVar.Key): `"$($envVar.Value)`"`n"
}

# Replace the env_variables section in app.yaml
$appYamlContent = $appYamlContent -replace "env_variables:.*?(?=network:|$)", $envVarsSection

# Write to temporary file
$appYamlContent | Out-File -FilePath "app_temp.yaml" -Encoding utf8

Write-Host "Deploying to Google App Engine with environment variables..." -ForegroundColor Green
& gcloud app deploy app_temp.yaml

# Clean up
Remove-Item "app_temp.yaml" -Force

Write-Host "Deployment completed successfully!" -ForegroundColor Green
Write-Host "To view your deployed application, run:" -ForegroundColor Green
Write-Host "gcloud app browse" -ForegroundColor Cyan 