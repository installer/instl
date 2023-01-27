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

verbose "Setting up variables"
owner="{{ .Owner }}"
repo="{{ .Repo }}"
verbose "Creating temporary directory"
tmpDir="$(mktemp -d)"

binaryLocation="$HOME/.local/bin"
verbose "Binary location: $binaryLocation"
mkdir -p $binaryLocation

installLocation="$HOME/.local/bin/.instl/$repo"
verbose "Install location: $installLocation"
# Remove installLocation directory if it exists
if [ -d "$installLocation" ]; then
  verbose "Removing existing install location"
  rm -rf "$installLocation"
fi
mkdir -p $installLocation

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
releaseJSON="$(curl "${getReleaseArgs[@]}")"
tagName="$(echo "$releaseJSON" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')"
info "Found latest release of $repo (version: $tagName)"

# Get list of assets
verbose "Parsing release json"
assets=$(echo "$releaseJSON" | grep "browser_download_url" | sed -E 's/.*"([^"]+)".*/\1/')
verbose "Found assets:"
verbose "$assets"
assetCount="$(echo "$assets" | wc -l | sed -E 's/^[[:space:]]*//')"
info "Found $assetCount assets in '$tagName' release - searching for one that fits your system..."

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
      score=1
      break
    fi
  done

  # Skip assets that do not match the operating system of the host
  if [ $score != 1 ]; then
    verbose "Skipping asset $fileName because it does not match the current operating system"
    continue
  fi

  verbose "$fileName matches current operating system"

  # Add one to the score for every alias that matches the current architecture
  for archAlias in "${currentArchAliases[@]}"; do
    if [[ "${fileName}" == *"$archAlias"* ]]; then
      verbose "Adding one to score for asset $fileName because it matches architecture $archAlias"
      score=$((score + 1))
    fi
  done

  # Initialize asset with score
  map_put assets "$fileName" "$score"
done

# Get map entry with highest score
verbose "Finding asset with highest score"
maxScore=0
maxKey=""
for key in $(map_keys assets); do
  score="$(map_get assets "$key")"
  if [ $score -gt $maxScore ]; then
    maxScore=$score
    maxKey=$key
  fi
done

assetName="$maxKey"
# Get asset URL from release assets
assetURL="$(echo "$assets" | grep -i "$assetName")"

info "Found asset with highest match score: $assetName"

info "Downloading asset..."
# Download asset
downloadAssetArgs=$curlOpts
downloadAssetArgs+=(-L "$assetURL" -o "$tmpDir/$assetName")
curl "${downloadAssetArgs[@]}"

# Unpack asset if it is a  tar, tar.gz or tar.bz2 file
if [[ "$assetName" == *".tar" ]]; then
  verbose "Unpacking .tar asset"
  tar -xf "$tmpDir/$assetName" -C "$tmpDir"
elif [[ "$assetName" == *".tar.gz" ]]; then
  verbose "Unpacking .tar.gz asset"
  tar -xzf "$tmpDir/$assetName" -C "$tmpDir"
elif [[ "$assetName" == *".tar.bz2" ]]; then
  verbose "Unpacking .tar.bz2 asset"
  tar -xjf "$tmpDir/$assetName" -C "$tmpDir"
fi

# Removing packed asset
verbose "Removing packed asset"
rm "$tmpDir/$assetName"

# Copy files to install location
info "Installing '$repo'"
verbose "Copying files to install location"
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

# Remove previous symlink if it exists
verbose "Removing previous symlink"
rm "$binaryLocation/$repo" > /dev/null 2>&1 || true
# Create symlink to binary
verbose "Creating symlink '$binaryLocation/$binaryName' -> '$binary'"
ln -s "$binary" "$binaryLocation/$binaryName"

# Add binaryLocation to PATH, if it is not already in PATH
if ! echo "$PATH" | grep -q "$binaryLocation"; then
  verbose "Adding $binaryLocation to PATH"
  # Append binaryLocation to .profile, if it is not already in .profile
  if ! grep -q -s "export PATH=$binaryLocation:\$PATH" "$HOME/.profile"; then
    verbose "Appending $binaryLocation to $HOME/.profile"
    echo "export PATH=$binaryLocation:\$PATH" >>"$HOME/.profile"
  fi
fi

info "Running clean up..."
verbose "Removing temporary directory"
rm -rf "$tmpDir"

echo
success "You can now run '$binaryName' in your terminal."
info "You might have to restart your terminal session for the changes to take effect."

source ../footer.txt
