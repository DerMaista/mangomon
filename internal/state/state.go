package state

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppState struct {
	GridSize int `json:"grid_size"`
}

func GetStatePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "state.json"
	}
	return filepath.Join(home, ".config", "mangomon", "state.json")
}

func Load() (AppState, error) {
	path := GetStatePath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return AppState{GridSize: 1}, nil // Default
		}
		return AppState{}, err
	}

	var state AppState
	if err := json.Unmarshal(data, &state); err != nil {
		return AppState{GridSize: 1}, nil
	}
	return state, nil
}

func Save(state AppState) error {
	path := GetStatePath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
