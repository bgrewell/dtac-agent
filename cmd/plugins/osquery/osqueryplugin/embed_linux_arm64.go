//go:build embed_osquery && linux && arm64

package osqueryplugin

import _ "embed"

//go:embed binaries/linux-arm64/osqueryd
var embeddedOsqueryd []byte
