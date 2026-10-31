//go:build unix

package linkding

import (
	"fmt"
	"os"
)

var configPermWarn func(string)

func init() {
	configPermWarn = func(path string) {
		st, err := os.Stat(path)
		if err != nil {
			return
		}
		if st.Mode().Perm()&0o077 != 0 {
			fmt.Fprintf(os.Stderr, "linkding: warning: %s is group/world-readable (use chmod 600)\n", path)
		}
	}
}

func auditConfigFilePerm(path string) {
	if configPermWarn != nil {
		configPermWarn(path)
	}
}
