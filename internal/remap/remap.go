package remap

import "github.com/istorry/pmguard/internal/detect"

// flagMap translates flags from one PM to another.
// Key: [fromPM][toPM][originalFlag] = translatedFlag
var flagMap = map[detect.PackageManager]map[detect.PackageManager]map[string]string{
	detect.Pnpm: {
		detect.Bun: {
			"-D":         "-d",
			"--save-dev": "-d",
			"-g":         "-g",
		},
		detect.Yarn: {
			"-D": "-D",
			"-g": "global add", // handled specially
		},
		detect.Npm: {
			"-D": "--save-dev",
			"-g": "-g",
		},
	},
	detect.Npm: {
		detect.Bun: {
			"--save-dev": "-d",
			"-D":         "-d",
			"-g":         "-g",
		},
		detect.Pnpm: {
			"--save-dev": "-D",
			"-D":         "-D",
		},
	},
	detect.Yarn: {
		detect.Bun: {
			"-D":    "-d",
			"--dev": "-d",
		},
		detect.Pnpm: {
			"-D":    "-D",
			"--dev": "-D",
		},
	},
	detect.Bun: {
		detect.Pnpm: {
			"-d": "-D",
		},
		detect.Npm: {
			"-d": "--save-dev",
		},
	},
}

// subcommandMap translates subcommands between PMs.
// e.g. "pnpm add" -> "bun add" (same), "npm install" -> "bun install" (same)
var subcommandMap = map[detect.PackageManager]map[detect.PackageManager]map[string]string{
	detect.Npm: {
		detect.Bun:  {"install": "install", "i": "install", "ci": "install --frozen-lockfile", "run": "run"},
		detect.Pnpm: {"install": "install", "i": "install", "run": "run"},
		detect.Yarn: {"install": "install", "i": "install", "run": "run"},
	},
	detect.Pnpm: {
		detect.Bun:  {"install": "install", "add": "add", "run": "run", "dlx": "x"},
		detect.Npm:  {"install": "install", "add": "install", "run": "run", "dlx": "npx"},
		detect.Yarn: {"install": "install", "add": "add", "run": "run"},
	},
	detect.Yarn: {
		detect.Bun:  {"install": "install", "add": "add", "run": "run"},
		detect.Pnpm: {"install": "install", "add": "add", "run": "run"},
		detect.Npm:  {"install": "install", "add": "install", "run": "run"},
	},
	detect.Bun: {
		detect.Pnpm: {"install": "install", "add": "add", "run": "run", "x": "dlx"},
		detect.Npm:  {"install": "install", "add": "install", "run": "run", "x": "npx"},
		detect.Yarn: {"install": "install", "add": "add", "run": "run"},
	},
}

// RemapArgs translates args from one PM's convention to another's.
// e.g. ["add", "-D", "react"] from pnpm -> bun becomes ["add", "-d", "react"]
func RemapArgs(from, to detect.PackageManager, args []string) []string {
	if len(args) == 0 {
		return args
	}

	result := make([]string, 0, len(args))

	// Translate subcommand (first arg)
	subcommand := args[0]
	if subs, ok := subcommandMap[from][to]; ok {
		if mapped, ok := subs[subcommand]; ok {
			subcommand = mapped
		}
	}
	result = append(result, subcommand)

	// Translate remaining flags
	flags, ok := flagMap[from][to]
	for _, arg := range args[1:] {
		if ok {
			if mapped, exists := flags[arg]; exists {
				result = append(result, mapped)
				continue
			}
		}
		result = append(result, arg)
	}

	return result
}
