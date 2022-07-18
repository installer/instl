package scripts

import (
	"path"
	"regexp"
	"strings"
)

// CombineLinuxScripts combines all Linux scripts into one.
func CombineLinuxScripts() (string, error) {
	return CombineBash("assets/linux/script.template.sh")
}

// CombineWindowsScripts combines all Windows scripts into one.
func CombineWindowsScripts() (string, error) {
	return CombinePowerShell("assets/windows/script.template.ps1")
}

// CombineBash combines multiple Bash files into one.
func CombineBash(filePath string) (string, error) {
	return CombineFile(filePath, "source (.*)")
}

// CombinePowerShell combines multiple PowerShell files into one.
func CombinePowerShell(filePath string) (string, error) {
	return CombineFile(filePath, `\. "(.*)"`)
}

// CombineFile combines multiple files into one.
// A regex is used to find the path to the file that should be injected.
func CombineFile(filePath, pattern string) (string, error) {
	regex := regexp.MustCompile(pattern)

	contentRaw, err := scripts.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	content := string(contentRaw)

	// For each line in the file, check if it matches the pattern
	// If it does, replace the line with the content of the file
	// If it doesn't, keep the line as is
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if regex.MatchString(line) {
			foundPath := strings.TrimSpace(regex.FindStringSubmatch(line)[1])
			p := path.Join(path.Dir(filePath), foundPath)
			contentRaw, err := scripts.ReadFile(p)
			if err != nil {
				return "", err
			}
			if strings.HasSuffix(foundPath, ".txt") {
				lines[i] = string(contentRaw)
			} else {
				lines[i] = "# --- Sourced from file: " + foundPath + " ---\n\n" + string(contentRaw) + "\n# --- End of " + foundPath + " ---\n\n"
			}
		}
	}

	return strings.Join(lines, "\n"), nil
}
