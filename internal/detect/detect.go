package detect

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type PackageManager string

const (
	Bun     PackageManager = "bun"
	Pnpm    PackageManager = "pnpm"
	Yarn    PackageManager = "yarn"
	Npm     PackageManager = "npm"
	Unknown PackageManager = ""
)

// lockfiles maps lockfile names to their package managers
var lockfiles = map[string]PackageManager{
	"bun.lockb":         Bun,
	"bun.lock":          Bun,
	"pnpm-lock.yaml":    Pnpm,
	"yarn.lock":         Yarn,
	"package-lock.json": Npm,
}

type packageJSON struct {
	PackageManager string `json:"packageManager"`
}

// Detect walks up from cwd looking for a lockfile or packageManager field.
// Returns Unknown if nothing is found.
func Detect() PackageManager {
	dir, err := os.Getwd()
	if err != nil {
		return Unknown
	}

	for {
		// 1. Check lockfiles first (most reliable)
		for lockfile, pm := range lockfiles {
			if fileExists(filepath.Join(dir, lockfile)) {
				return pm
			}
		}

		// 2. Check packageManager field in package.json
		if pm := fromPackageJSON(dir); pm != Unknown {
			return pm
		}

		// Walk up to parent
		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached root
		}
		dir = parent
	}

	return Unknown
}

func fromPackageJSON(dir string) PackageManager {
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return Unknown
	}

	var pkg packageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return Unknown
	}

	// packageManager field looks like "pnpm@9.0.0" or just "pnpm"
	name := strings.SplitN(pkg.PackageManager, "@", 2)[0]
	for _, pm := range []PackageManager{Bun, Pnpm, Yarn, Npm} {
		if name == string(pm) {
			return pm
		}
	}

	return Unknown
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
