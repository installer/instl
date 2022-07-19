. "../header.txt"

# Import libraries
. "./lib/colors.ps1"

# Setup variables
$verboseString = "{{ .Verbose }}"
$verbose = $false
if ($verboseString -eq "true")
{
    $verbose = $true
}

verbose "Setting up variables"
$owner = "{{ .Owner }}"
$repo = "{{ .Repo }}"

$tmpDir = "$env:TEMP\instl-cache"
verbose "Temporary directory: $tmpDir"
# Remove previous cache, if it exists
if (test-path $tmpDir)
{
    rm -r -fo $tmpDir
}
New-Item -Path $tmpDir -ItemType Directory > $null

$installLocation = "$HOME\instl\$repo"
verbose "Install location: $installLocation"
# Remove previous installation, if it exists
if (test-path $installLocation)
{
    rm -r -fo $installLocation
}
New-Item -Path $installLocation -ItemType Directory > $null

# Print "INSTL" header
. "../shared/intro.ps1"

# Installation

# GitHub public API
$latestReleaseURL = "https://api.github.com/repos/$owner/$repo/releases/latest"
$latestRelease = Invoke-WebRequest $latestReleaseURL | ConvertFrom-Json
$tagName = $latestRelease.tag_name
info "Found latest release of $repo (version: $tagName)"

# Get list of assets
verbose "Getting list of assets"
$assetsRaw = $latestRelease.assets
# Map to array
$assets = $assetsRaw | %{ $_.browser_download_url }
$assetCount = $assets.Count
info "Found $assetCount assets in '$tagName' release - searching for one that fits your system..."

# Get host architecture
$arch = $env:PROCESSOR_ARCHITECTURE
# arch to lowercase
$arch = $arch.ToLower()
verbose "Host architecture: $arch"

# Set aliases for architecture
$amd64 = @("amd64", "x86_64", "x86-64", "x64")
$amd32 = @("386", "i386", "i686", "x86")
$windows = @("windows", "win")

$currentArchAliases = @()
if ($arch -eq "amd64")
{
    $currentArchAliases = $amd64
}
elseif ($arch -eq "x86")
{
    $currentArchAliases = $amd32
}
else
{
    error "Unsupported architecture: $arch"
}
verbose "Current architecture aliases: $currentArchAliases"

# Create hastable of assets and a score
$assetMap = @{ }
# Loop through assets
foreach ($asset in $assets)
{
    if ( $asset.ToLower().Contains("darwin"))
    {
        continue
    }
    $assetMap[$asset] = 0
    # Loop through windows aliases
    $windows | %{
        $windowsAlias = $_
        # If asset contains architecture alias, increase score
        if ( $asset.ToLower().Contains($windowsAlias))
        {
            verbose "Asset $asset contains windows alias $windowsAlias"
            $assetMap[$asset] = $assetMap[$asset] + 1
        }
    }

    # Loop through architecture aliases
    $currentArchAliases | %{
        $archAlias = $_
        # If asset contains architecture alias, increase score
        if ( $asset.ToLower().Contains($archAlias))
        {
            verbose "Asset $asset contains architecture alias $archAlias"
            $assetMap[$asset] = $assetMap[$asset] + 1
        }
    }
}

# Get highest score
$highestScore = 0
$highestScoreAsset = ""
foreach ($Key in $assetMap.Keys)
{
    $asset = $Key
    $score = $assetMap[$Key]
    verbose "Asset: $asset, score: $score"
    if ($score -gt $highestScore)
    {
        $highestScore = $score
        $highestScoreAsset = $asset
    }
}


$assetURL = $highestScoreAsset
$assetName = $assetURL.Split('/')[-1]
info "Found asset with highest match score: $assetName"

# Downoad asset
info "Downloading asset..."
$assetPath = "$tmpDir\$assetName"
Invoke-Expression "Invoke-WebRequest -Uri $assetURL -OutFile $assetPath"
verbose "Asset downloaded to $assetPath"

# Extract asset if it is a zip file
if ( $assetName.EndsWith(".zip"))
{
    verbose "Extracting asset..."
    Expand-Archive -Path $assetPath -Destination $installLocation\
    verbose "Asset extracted to $extractDir"
}

# Find binary file in install path
$binaryFile = (Get-ChildItem -Path $installLocation -Filter "*.exe")[0]
$binaryFile = $installLocation + "\" + $binaryFile
$binaryName = $binaryFile.Split('\')[-1]
verbose "Binary file: $binaryFile"

# Change PATH to include install location
verbose "Changing PATH to include install location"
$oldPath = (Get-ItemProperty -Path 'Registry::HKEY_CURRENT_USER\Environment' -Name PATH).path
# if oldPath does not contain install location, add it
if (!$oldPath.Contains($installLocation))
{
    verbose "PATH does not contain install location, adding it"
    $newPath = $oldPath + ";" + $installLocation
    Set-ItemProperty -Path 'Registry::HKEY_CURRENT_USER\Environment' -Name PATH -Value $newPath
}

info "Running clean up..."
if (test-path $tmpDir)
{
    verbose "Removing temporary directory"
    rm -r -fo $tmpDir
}

echo ""
success "You can now run '$binaryName' in your terminal."
info "You might have to restart your terminal session for the changes to take effect."

. "../footer.txt"
