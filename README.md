# mangomon

A TUI for configuring MangoWC monitor rules. Similar to hyprmon for Hyprland.

## Features

- Spatial monitor arrangement with arrow keys
- Resolution and refresh rate selection (reads from `/sys/class/drm`)
- Scale adjustment
- Transform/rotation editing
- Mirror configuration
- Variable Refresh Rate

## Installation

### Pre-built binary

A pre-built Linux binary is included in the repository. Just run:

```
./mangomon
```

### Build from source

Requires Go 1.21+

```
go build -o mangomon main.go
```

## Usage

Run `mangomon` and use the following keys:

| Key | Action |
|-----|--------|
| Tab | Cycle between monitors |
| Arrow keys | Move selected monitor |
| Shift+Arrow | Move faster |
| G | Cycle grid size |
| R | Open scale picker |
| F | Open resolution/mode picker |
| T | Open transform/rotation picker |
| V | Open VRR picker |
| M | Open mirror picker |
| S | Save config |
| Q | Quit |

## Dependencies

Requires `mmsg` from MangoWC to be available in PATH for querying connected outputs.

## Notes

- Restart MangoWC to apply layout changes
- Config is saved to `~/.config/mango/config.conf`

## Author

Created by [thatsjor](https://github.com/thatsjor) (Jordan F.)

## License

MIT
