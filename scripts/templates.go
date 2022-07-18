package scripts

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/installer/installer/internal/pkg/config"
	"github.com/installer/installer/internal/pkg/platforms"
)

var (
	windowsTemplate *template.Template
	linuxTemplate   *template.Template
)

func ParseTemplateForPlatform(platform platforms.Platform, config config.Config) (string, error) {
	var tpl bytes.Buffer
	var t *template.Template
	var err error

	switch platform {
	case platforms.Windows:
		t, err = GetWindowsTemplate()
	case platforms.MacOS:
		t, err = GetLinuxTemplate()
	case platforms.Linux:
		t, err = GetLinuxTemplate()
	default:
		return "", errors.New("unknown platform")
	}
	if err != nil {
		return "", err
	}

	err = t.Execute(&tpl, config)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func GetWindowsTemplate() (*template.Template, error) {
	if windowsTemplate != nil {
		return windowsTemplate, nil
	}
	windowsScript, err := CombineWindowsScripts()
	if err != nil {
		return nil, err
	}

	tpl, err := template.New("windows").Parse(windowsScript)
	if err != nil {
		return nil, err
	}
	windowsTemplate = tpl

	return windowsTemplate, nil
}

func GetLinuxTemplate() (*template.Template, error) {
	if linuxTemplate != nil {
		return linuxTemplate, nil
	}
	linuxScript, err := CombineLinuxScripts()
	if err != nil {
		return nil, err
	}

	tpl, err := template.New("linux").Parse(linuxScript)
	if err != nil {
		return nil, err
	}
	linuxTemplate = tpl

	return linuxTemplate, nil
}
