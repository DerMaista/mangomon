package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MonitorRule represents a single monitorrule line in the config
type MonitorRule struct {
	ID                  string
	Transform           int
	Scale               float64
	X, Y                int
	Width, Height       int
	RefreshRate         float64
	VariableRefreshRate int // 0 or 1
}

// ToString converts the rule back to the config string format
func (r MonitorRule) ToString() string {
	// monitorrule=name:eDP-1,width:1920,height:1080,refresh:60,x:0,y:0,scale:1.0,vrr:0,rr:0
	return fmt.Sprintf("monitorrule=name:%s,width:%d,height:%d,refresh:%.0f,x:%d,y:%d,scale:%.2f,vrr:%d,rr:%d",
		r.ID, r.Width, r.Height, r.RefreshRate, r.X, r.Y, r.Scale, r.VariableRefreshRate, r.Transform)
}

// ConfigParser handles reading and writing the MangoWC config
type ConfigParser struct {
	FilePath string
	Lines    []string // Store all lines to preserve comments/other configs
}

func NewParser(path string) (*ConfigParser, error) {
	// Default fallback if path is empty
	if path == "" {
		home, _ := os.UserHomeDir()
		path = home + "/.config/mango/config.conf"
	}
	return &ConfigParser{FilePath: path}, nil
}

func (p *ConfigParser) Parse() (map[string]MonitorRule, error) {
	file, err := os.Open(p.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]MonitorRule), nil
		}
		return nil, err
	}
	defer file.Close()

	var lines []string
	rules := make(map[string]MonitorRule)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "monitorrule=") {
			// Parse the rule
			val := strings.TrimPrefix(trimmed, "monitorrule=")
			parts := strings.Split(val, ",")

			// Expected format: key:value pairs
			// Example: name:eDP-1,width:1920,height:1080,refresh:60,x:0,y:0,scale:1.0,vrr:0,rr:0
			rule := MonitorRule{}
			for _, part := range parts {
				kv := strings.Split(strings.TrimSpace(part), ":")
				if len(kv) != 2 {
					continue
				}
				key := strings.TrimSpace(kv[0])
				val := strings.TrimSpace(kv[1])

				switch key {
				case "name":
					rule.ID = val
				case "width":
					rule.Width, _ = strconv.Atoi(val)
				case "height":
					rule.Height, _ = strconv.Atoi(val)
				case "refresh":
					rule.RefreshRate, _ = strconv.ParseFloat(val, 64)
				case "x":
					rule.X, _ = strconv.Atoi(val)
				case "y":
					rule.Y, _ = strconv.Atoi(val)
				case "scale":
					rule.Scale, _ = strconv.ParseFloat(val, 64)
				case "vrr":
					rule.VariableRefreshRate, _ = strconv.Atoi(val)
				case "rr":
					rule.Transform, _ = strconv.Atoi(val)
				}
			}

			if rule.ID != "" {
				rules[rule.ID] = rule
			}
		}
	}
	p.Lines = lines
	return rules, scanner.Err()
}

func (p *ConfigParser) Save(newRules []MonitorRule) error {
	// Reconstruct the file content
	// This is a naive implementation: it replaces existing monitorrule lines
	// and appends new ones if they didn't exist.
	// A better approach for preserving order:
	// 1. Iterate over p.Lines.
	// 2. If a line matches a rule we are updating, replace it. Mark rule as "written".
	// 3. If line matches specific rule that is NOT in newRules (deleted?), maybe comment it out or delete.
	// 4. Append remaining unwritten rules at the end (or specific section).

	// Let's go with: Replace existing, Append new.

	updatedLines := make([]string, 0, len(p.Lines))
	writtenIDs := make(map[string]bool)

	for _, line := range p.Lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "monitorrule=") {
			// Check if this line corresponds to one of our new rules
			currentID := strings.Split(strings.TrimPrefix(trimmed, "monitorrule="), ",")[0]

			foundRequest := false
			for _, nr := range newRules {
				if nr.ID == currentID {
					updatedLines = append(updatedLines, nr.ToString())
					writtenIDs[nr.ID] = true
					foundRequest = true
					break
				}
			}
			if !foundRequest {
				// Keep existing rules that we aren't modifying?
				// Or provided list is authoritative?
				// Usually for TUI, what isn't configured might be kept.
				// Let's keep it unless we want to delete.
				updatedLines = append(updatedLines, line)
			}
		} else {
			updatedLines = append(updatedLines, line)
		}
	}

	// Append new rules that weren't in the file
	for _, nr := range newRules {
		if !writtenIDs[nr.ID] {
			updatedLines = append(updatedLines, nr.ToString())
		}
	}

	file, err := os.Create(p.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range updatedLines {
		fmt.Fprintln(w, line)
	}
	// ... previous code ...
	return w.Flush()
}

// Profile management

func (p *ConfigParser) ListProfiles() ([]string, error) {
	profilesDir := strings.ReplaceAll(p.FilePath, "config.conf", "profiles")
	entries, err := os.ReadDir(profilesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var profiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".conf") {
			profiles = append(profiles, strings.TrimSuffix(e.Name(), ".conf"))
		}
	}
	return profiles, nil
}

func (p *ConfigParser) SaveProfile(name string, rules []MonitorRule) error {
	profilesDir := strings.ReplaceAll(p.FilePath, "config.conf", "profiles")
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s.conf", profilesDir, name)
	// For a profile, we just save the MonitorRule lines for now?
	// Or do we save the whole config?
	// Hyprmon saves the "monitors.conf" part usually.
	// We'll write just the rules.

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, r := range rules {
		fmt.Fprintln(w, r.ToString())
	}
	return w.Flush()
}

func (p *ConfigParser) LoadProfile(name string) (map[string]MonitorRule, error) {
	profilesDir := strings.ReplaceAll(p.FilePath, "config.conf", "profiles")
	path := fmt.Sprintf("%s/%s.conf", profilesDir, name)

	// Create a temporary parser for the profile file
	profParser := &ConfigParser{FilePath: path}
	return profParser.Parse()
}
