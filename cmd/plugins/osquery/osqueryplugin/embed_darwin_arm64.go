//go:build embed_osquery && darwin && arm64

package osqueryplugin

import _ "embed"

//go:embed binaries/darwin-arm64/osqueryd
var embeddedOsqueryd []byte
