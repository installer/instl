package platforms

import "fmt"

var platformMap = map[string]Platform{
	"windows": Windows,
	"win":     Windows,

	"macos":  MacOS,
	"darwin": MacOS,
	"mac":    MacOS,

	"linux": Linux,
}

func Parse(platform string) (Platform, error) {
	if p, ok := platformMap[platform]; ok {
		return p, nil
	}

	return 0, fmt.Errorf("unknown platform: %s", platform)
}
