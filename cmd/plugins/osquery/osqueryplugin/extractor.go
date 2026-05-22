package osqueryplugin

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrNoEmbeddedBinary is returned when the build did not include an embedded
// osqueryd (the default — built without `-tags embed_osquery`, or on an
// unsupported platform).
var ErrNoEmbeddedBinary = errors.New("no osqueryd binary embedded in this build")

// extractEmbedded materializes the bundled osqueryd to a content-addressed
// cache directory and returns its path. Subsequent calls for the same embed
// bytes are no-ops — the cached file is reused.
//
// Cache location: /var/cache/dtac/osquery/<sha256>/osqueryd, with fallback
// to $TMPDIR/dtac-osquery/<sha256>/osqueryd when /var/cache isn't writable.
func extractEmbedded() (string, error) {
	if len(embeddedOsqueryd) == 0 {
		return "", ErrNoEmbeddedBinary
	}

	sum := sha256.Sum256(embeddedOsqueryd)
	hash := hex.EncodeToString(sum[:])

	for _, root := range cacheRoots() {
		path, err := tryExtractInto(root, hash, embeddedOsqueryd)
		if err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not extract embedded osqueryd: no writable cache directory")
}

// cacheRoots returns the candidate cache base directories, in preference
// order. The first one we can write to wins.
func cacheRoots() []string {
	return []string{
		"/var/cache/dtac/osquery",
		filepath.Join(os.TempDir(), "dtac-osquery"),
	}
}

// tryExtractInto writes the binary into <root>/<hash>/osqueryd atomically,
// or reuses an existing valid copy. Returns the binary path on success.
func tryExtractInto(root, hash string, payload []byte) (string, error) {
	dir := filepath.Join(root, hash)
	target := filepath.Join(dir, "osqueryd")

	// Fast path: cached binary already present and the right size. We trust
	// the hash because it's part of the directory name; a partial write would
	// have been atomically renamed only if it succeeded.
	if info, err := os.Stat(target); err == nil && info.Size() == int64(len(payload)) {
		return target, nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	// Atomic write: temp file in the same directory, then rename.
	tmp, err := os.CreateTemp(dir, "osqueryd.*.tmp")
	if err != nil {
		return "", err
	}
	tmpName := tmp.Name()
	cleanup := func() { _ = os.Remove(tmpName) }

	if _, err := tmp.Write(payload); err != nil {
		_ = tmp.Close()
		cleanup()
		return "", err
	}
	if err := tmp.Chmod(0o755); err != nil {
		_ = tmp.Close()
		cleanup()
		return "", err
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return "", err
	}
	if err := os.Rename(tmpName, target); err != nil {
		cleanup()
		return "", err
	}
	return target, nil
}
