package linkding

import (
	"os"
	"path/filepath"
	"strings"
)

func parseAllowWrite(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "true", "yes", "write", "on":
		return true
	default:
		return false
	}
}

// EffectiveAllowWrite resolves write policy: LINKDING_ALLOW_WRITE env, then
// ~/.config/linkding/allow_write (or $LINKDING_CONFIG_DIR/allow_write), then
// allow_write in config.json (first matching file in search order).
func EffectiveAllowWrite(getenv func(string) string, fromConfig bool, configSource string) (allowed bool, source string) {
	if v := getenv("LINKDING_ALLOW_WRITE"); v != "" {
		return parseAllowWrite(v), "LINKDING_ALLOW_WRITE"
	}
	dir, err := configDirFrom(getenv)
	if err == nil {
		b, err := os.ReadFile(filepath.Join(dir, "allow_write"))
		if err == nil && parseAllowWrite(string(bytesTrimSpace(b))) {
			return true, "allow_write file"
		}
	}
	if fromConfig {
		src := configSource
		if src == "" {
			src = "config.json"
		}
		return true, src
	}
	return false, "default"
}
