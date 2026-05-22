//go:build embed_osquery && linux && amd64

package osqueryplugin

import _ "embed"

//go:embed binaries/linux-amd64/osqueryd
var embeddedOsqueryd []byte
