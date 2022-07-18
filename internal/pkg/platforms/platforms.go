package platforms

type Platform uint8

func (p Platform) String() string {
	switch p {
	case Windows:
		return "windows"
	case MacOS:
		return "macos"
	case Linux:
		return "linux"
	case Unknown:
		return "unknown"
	default:
		return "unknown"
	}
}

const (
	Unknown Platform = iota
	Linux
	MacOS
	Windows
)
