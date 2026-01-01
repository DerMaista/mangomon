package system

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Output struct {
	Name string
}

func GetOutputs() ([]Output, error) {
	cmd := exec.Command("mmsg", "-O")
	outputBytes, err := cmd.Output()
	if err != nil {
		return []Output{
			{Name: "eDP-1"},
			{Name: "HDMI-A-1"},
		}, nil
	}

	outputStr := string(outputBytes)
	lines := strings.Split(outputStr, "\n")
	var outputs []Output

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {

			parts := strings.Fields(trimmed)
			if len(parts) > 0 {
				outputs = append(outputs, Output{Name: parts[0]})
			}
		}
	}

	return outputs, nil
}

// Mode represents a display mode
type Mode struct {
	Width, Height int
	Rate          float64
}

func GetModes(output string) ([]Mode, error) {
	sysPath := "/sys/class/drm"
	files, err := os.ReadDir(sysPath)
	if err != nil {
		return getFallbackModes(), fmt.Errorf("failed to read %s: %w", sysPath, err)
	}

	var modeFile string
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, "-"+output) {
			modeFile = fmt.Sprintf("%s/%s/modes", sysPath, name)
			break
		}
	}

	if modeFile == "" {
		return getFallbackModes(), fmt.Errorf("could not find drm folder for output %s", output)
	}

	content, err := os.ReadFile(modeFile)
	if err != nil {
		return getFallbackModes(), fmt.Errorf("failed to read modes file %s: %w", modeFile, err)
	}

	lines := strings.Split(string(content), "\n")
	var modes []Mode

	seen := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if seen[line] {
			continue
		}
		seen[line] = true

		parts := strings.Split(line, "x")
		if len(parts) != 2 {
			continue
		}
		w, _ := strconv.Atoi(parts[0])
		h, _ := strconv.Atoi(parts[1])

		modes = append(modes, Mode{Width: w, Height: h, Rate: 60.0})

		if w >= 1920 {
			modes = append(modes, Mode{Width: w, Height: h, Rate: 120.0})
			modes = append(modes, Mode{Width: w, Height: h, Rate: 144.0})
			modes = append(modes, Mode{Width: w, Height: h, Rate: 165.0})
			modes = append(modes, Mode{Width: w, Height: h, Rate: 240.0})
		}
	}

	if len(modes) == 0 {
		return getFallbackModes(), nil
	}

	return modes, nil
}

func getFallbackModes() []Mode {
	return []Mode{
		{Width: 3840, Height: 2160, Rate: 144.0},
		{Width: 3840, Height: 2160, Rate: 60.0},
		{Width: 2560, Height: 1440, Rate: 165.0},
		{Width: 2560, Height: 1440, Rate: 144.0},
		{Width: 2560, Height: 1440, Rate: 60.0},
		{Width: 1920, Height: 1080, Rate: 144.0},
		{Width: 1920, Height: 1080, Rate: 60.0},
	}
}
