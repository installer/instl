. "../header.txt"

# Setup variables
$owner="{{ .Owner }}"
$repo="{{ .Repo }}"
$github_url="https://github.com/$owner/$repo"
$install_location=""

# Import libraries
. "./lib/colors.ps1"

# Print "INSTL" header
. "../shared/intro.ps1"

# Installation

# GitHub public API
$latestReleaseURL = "https://api.github.com/repos/$owner/$repo/releases/latest"
$latestRelease = Invoke-WebRequest $latestReleaseURL | ConvertFrom-Json
$tagName = $latestRelease.tag_name
info "Found latest release of $repo (version: $tagName)"

. "../footer.txt"
