package styling

import (
	"strings"
)

func TrimColourResetChar(input string) string {
	return strings.TrimSuffix(input, "\x1b[0m")
}
