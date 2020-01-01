// +build !windows

package maxmind

import (
	"strings"
)

func formatPath(path string) string {
	return strings.Replace(path, `\`, `/`, -1)
}
