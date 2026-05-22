//go:build !embed_osquery

package osqueryplugin

// embeddedOsqueryd is empty in default builds. Build with `-tags embed_osquery`
// (and a per-platform binary placed under binaries/<os-arch>/osqueryd) to ship
// osqueryd inside the plugin. The mage `BundleOsquery` target does both steps.
var embeddedOsqueryd []byte
