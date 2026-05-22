//go:build embed_osquery && !((linux && amd64) || (linux && arm64) || (darwin && arm64))

package osqueryplugin

// embeddedOsqueryd stays empty on platforms we don't ship a bundled binary for
// (e.g. darwin/amd64, windows/*). Building with `-tags embed_osquery` is still
// allowed on those platforms — the runtime extractor falls through to
// config.binary_path or a system osqueryd lookup.
var embeddedOsqueryd []byte
