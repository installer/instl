source ../header.txt

# Import libraries
source ./lib/colors.sh
source ./lib/map.sh

# Setup variables
verbose="{{ .Verbose }}"
if [ "$verbose" = "true" ]; then
  verbose=true
else
  verbose=false
fi

# Pass variables from the go server into the script
verbose "Setting up variables"
owner="{{ .Owner }}"
repo="{{ .Repo }}"

verbose "Creating temporary directory"
tmpDir="$(mktemp -d)"

binaryLocation="$HOME/.local/bin"
verbose "Binary location: $binaryLocation"
mkdir -p "$binaryLocation"

installLocation="$HOME/.local/bin/.instl/$repo"
verbose "Install location: $installLocation"
# Remove installLocation directory if it exists
if [ -d "$installLocation" ]; then
  verbose "Removing existing install location"
  rm -rf "$installLocation"
fi
mkdir -p "$installLocation"

# Print "INSTL" header
source ../shared/intro.ps1

# Installation
curlOpts=("-sS")
if [ -n "$GH_TOKEN" ]; then
  verbose "Using authentication with GH_TOKEN"
  curlOpts+=("--header \"Authorization: Bearer $GH_TOKEN\"")
elif [ -n "$GITHUB_TOKEN" ]; then
  verbose "Using authentication with GITHUB_TOKEN"
  curlOpts+=("--header \"Authorization: Bearer $GITHUB_TOKEN\"")
else
  verbose "No authentication"
fi

# GitHub public API
latestReleaseURL="https://api.github.com/repos/$owner/$repo/releases/latest"
verbose "Getting latest release from GitHub"
getReleaseArgs=$curlOpts
getReleaseArgs+=("$latestReleaseURL")
releaseJSON="$(curl ${getReleaseArgs[@]})"
tagName="$(echo "$releaseJSON" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')"
info "Found latest release of $repo (version: $tagName)"

# Get list of assets
verbose "Parsing release json"
assets=$(echo "$releaseJSON" | grep "browser_download_url" | sed -E 's/.*"([^"]+)".*/\1/')
verbose "Found assets:"
verbose "$assets"
assetCount="$(echo "$assets" | wc -l | sed -E 's/^[[:space:]]*//')"
info "Found $assetCount assets in '$tagName' - searching for one that fits your system..."

# Get architecture of host
arch="$(uname -m)"
# Convert arch to lowercase
arch="$(echo "$arch" | tr '[:upper:]' '[:lower:]')"
verbose "Host architecture: $arch"

# Set aliases for architectures
amd64=("amd64" "x86_64" "x86-64" "x64")
amd32=("386" "i386" "i686" "x86")
arm=("arm" "armv7" "armv6" "armv8" "armv8l" "armv7l" "armv6l" "armv8l" "armv7l" "armv6l")

currentArchAliases=()
# Set the right aliases for current host system
if [ $arch == "x86_64" ]; then
  currentArchAliases=("${amd64[@]}")
elif [ $arch == "i386" ] || [ $arch == "i686" ]; then
  currentArchAliases=("${amd32[@]}")
# Else if starts with "arm"
elif [[ $arch =~ ^arm ]]; then
  currentArchAliases=("${arm[@]}")
else
  error "Unsupported architecture: $arch"
  exit 1
fi
verbose "Current architecture aliases: ${currentArchAliases[*]}"

# Get operating system of host
os="$(uname -s)"
# Convert os to lowercase
os="$(echo "$os" | tr '[:upper:]' '[:lower:]')"
verbose "Host operating system: $os"

# Set aliases for operating systems
linux=("linux")
darwin=("darwin" "macos" "osx")

currentOsAliases=()
# If current os is linux, add linux aliases to the curentOsAliases array
if [ "${os}" == "linux" ]; then
  currentOsAliases+=(${linux[@]})
elif [ "${os}" == "darwin" ]; then
  currentOsAliases+=(${darwin[@]})
fi
verbose "Current operating system aliases: ${currentOsAliases[*]}"

# Create map of assets and a score
for asset in $assets; do
  score=0

  # Get file name from asset path
  fileName="$(echo "$asset" | awk -F'/' '{ print $NF; }')"
  # Set filename to lowercase
  fileName="$(echo "$fileName" | tr '[:upper:]' '[:lower:]')"

  # Set score to one, if the file name contains the current os
  for osAlias in "${currentOsAliases[@]}"; do
    if [[ "${fileName}" == *"$osAlias"* ]]; then
      score=10
      break
    fi
  done

  # Add two to the score for every alias that matches the current architecture
  for archAlias in "${currentArchAliases[@]}"; do
    if [[ "${fileName}" == *"$archAlias"* ]]; then
      verbose "Adding one to score for asset $fileName because it matches architecture $archAlias"
      score=$((score + 2))
    fi
  done

  # Add one to the score if the file name contains .tar or .tar.gz or .tar.bz2
  if [[ "$fileName" == *".tar" ]] || [[ "$fileName" == *".tar.gz" ]] || [[ "$fileName" == *".tar.bz2" ]]; then
    verbose "Adding one to score for asset $fileName because it is a .tar or .tar.gz or .tar.bz2 file"
    score=$((score + 1))
  fi

  # Add two to the score if the file name contains the repo name
  if [[ "$fileName" == *"$repo"* ]]; then
    verbose "Adding two to score for asset $fileName because it contains the repo name"
    score=$((score + 2))
  fi

  # Add one to the score if the file name is exactly the repo name
  if [[ "$fileName" == "$repo" ]]; then
    verbose "Adding one to score for asset $fileName because it is exactly the repo name"
    score=$((score + 1))
  fi

  # Initialize asset with score
  map_put assets "$fileName" "$score"
done

# Get map entry with highest score
verbose "Finding asset with highest score"
maxScore=0
maxKey=""
for asset in $(map_keys assets); do
  score="$(map_get assets "$asset")"
  if [ $score -gt $maxScore ]; then
    maxScore=$score
    maxKey=$asset
  fi
  verbose "Asset: $asset, score: $score"
done

assetName="$maxKey"

# Check if asset name is still empty
if [ -z "$assetName" ]; then
  error "Could not find any assets that fit your system"
  exit 1
fi

# Get asset URL from release assets
assetURL="$(echo "$assets" | grep -i "$assetName")"

info "Found asset with highest match score: $assetName"

info "Downloading asset..."
# Download asset
downloadAssetArgs=$curlOpts
downloadAssetArgs+=(-L "$assetURL" -o "$tmpDir/$assetName")
verbose "${downloadAssetArgs[@]}"
curl "${downloadAssetArgs[@]}"

# Unpack asset if it is compressed
if [[ "$assetName" == *".tar" ]]; then
  verbose "Unpacking .tar asset to $tmpDir"
  tar -xf "$tmpDir/$assetName" -C "$tmpDir"
  verbose "Removing packed asset ($tmpDir/$assetName)"
  rm "$tmpDir/$assetName"
elif [[ "$assetName" == *".tar.gz" ]]; then
  verbose "Unpacking .tar.gz asset to $tmpDir"
  tar -xzf "$tmpDir/$assetName" -C "$tmpDir"
  verbose "Removing packed asset ($tmpDir/$assetName)"
  rm "$tmpDir/$assetName"
elif [[ "$assetName" == *".gz" ]]; then
  verbose "Unpacking .gz asset to $tmpDir/$repo"
  gunzip -c "$tmpDir/$assetName" > "$tmpDir/$repo"
  verbose "Removing packed asset ($tmpDir/$assetName)"
  rm "$tmpDir/$assetName"
  verbose "Setting asset name to $repo, because it is a .gz file"
  assetName="$repo"
  verbose "Marking asset as executable"
  chmod +x "$tmpDir/$repo"
elif [[ "$assetName" == *".tar.bz2" ]]; then
  verbose "Unpacking .tar.bz2 asset to $tmpDir"
  tar -xjf "$tmpDir/$assetName" -C "$tmpDir"
  verbose "Removing packed asset"
  rm "$tmpDir/$assetName"
elif [[ "$assetName" == *".zip" ]]; then
  verbose "Unpacking .zip asset to $tmpDir/$repo"
  unzip "$tmpDir/$assetName" -d "$tmpDir"
  verbose "Removing packed asset ($tmpDir/$assetName)"
  rm "$tmpDir/$assetName"
else
  verbose "Asset is not a tar or zip file. Skipping unpacking."
fi

# Copy files to install location
info "Installing '$repo'"
verbose "Copying files to install location"
# verbose print tmpDir files
verbose "Files in $tmpDir:"
verbose "$(ls "$tmpDir")"
cp -r "$tmpDir"/* "$installLocation"

# Find binary in install location to symlink to it later
verbose "Finding binary in install location"
binary=""
for file in "$installLocation"/*; do
  # Check if the file is a binary file
  verbose "Checking if $file is a binary file"
  if [ -f "$file" ] && [ -x "$file" ]; then
    binary="$file"
    verbose "Found binary: $binary"
    break
  fi
done

# Get name of binary
binaryName="$(basename "$binary")"
verbose "Binary name: $binaryName"

# Check if binary is empty
if [ -z "$binary" ]; then
  error "Could not find binary in install location"
  exit 1
fi

# Remove previous symlink if it exists
verbose "Removing previous symlink"
rm "$binaryLocation/$repo" >/dev/null 2>&1 || true
# Create symlink to binary
verbose "Creating symlink '$binaryLocation/$binaryName' -> '$binary'"
ln -s "$binary" "$binaryLocation/$binaryName"

# Add binaryLocation to PATH, if it is not already in PATH
if ! echo "$PATH" | grep -q "$binaryLocation"; then
  verbose "Adding $binaryLocation to PATH"
  # Array of common shell configuration files
  config_files=("$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.zshrc" "$HOME/.profile")
  for config in "${config_files[@]}"; do
    # Check if the file exists
    if [ -f "$config" ]; then
      # Check if binaryLocation is already in the file
      if ! grep -q -s "export PATH=$binaryLocation:\$PATH" "$config"; then
        verbose "Appending $binaryLocation to $config"
        echo "" >>"$config"
        echo "# Instl.sh installer binary location" >>"$config"
        echo "export PATH=$binaryLocation:\$PATH" >>"$config"
      else
        verbose "$binaryLocation is already in $config"
      fi
    else
      verbose "$config does not exist. Skipping append."
    fi
  done
else
  verbose "$binaryLocation is already in PATH"
fi

info "Running clean up..."
verbose "Removing temporary directory"
rm -rf "$tmpDir"

echo
success "You can now run '$binaryName' in your terminal!"
info "You might have to restart your terminal session for the changes to take effect"

source ../footer.txt
