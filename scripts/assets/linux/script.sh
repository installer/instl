source ../header.txt

# Setup variables
owner="{{ .Owner }}"
repo="{{ .Repo }}"
github_url="https://github.com/$owner/$repo"
install_location=""

# Import libraries
source ./lib/colors.sh

# Print "INSTL" header
source ./lib/intro.sh

# Installation

# GitHub public API
latestReleaseURL="https://api.github.com/repos/{{ .Owner }}/{{ .Repo }}/releases/latest"
tagName="$(curl -sS $latestReleaseURL | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')"
info "Found latest release of $repo (version: $tagName)"


source ../footer.txt
