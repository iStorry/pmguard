package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/istorry/pmguard/internal/config"
	"github.com/istorry/pmguard/internal/detect"
	"github.com/istorry/pmguard/internal/remap"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pmguard",
	Short: "pmguard — use the right package manager, automatically",
	Long: `pmguard intercepts package manager commands and warns or redirects
if you're using the wrong one for the current project.

Add this to your shell config to activate:
  eval "$(pmguard install-hooks)"`,
}

// guardCmd is the core: called by shell hooks as e.g. "pmguard guard pnpm install"
var guardCmd = &cobra.Command{
	Use:                "guard <invoked-pm> [args...]",
	Short:              "Check and optionally redirect a package manager command",
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		invokedPM := detect.PackageManager(args[0])
		pmArgs := args[1:]

		cfg, err := config.Load()
		if err != nil {
			cfg = config.Default
		}

		detectedPM := detect.Detect()

		// No lockfile found — just pass through transparently
		if detectedPM == detect.Unknown {
			return passthrough(string(invokedPM), pmArgs)
		}

		// Correct PM — pass through
		if invokedPM == detectedPM {
			return passthrough(string(invokedPM), pmArgs)
		}

		// Wrong PM detected
		switch cfg.Mode {
		case config.ModeWarn:
			fmt.Fprintf(os.Stderr, "\n⚠️  pmguard: This project uses %s, but you ran %s.\n", detectedPM, invokedPM)
			fmt.Fprintf(os.Stderr, "   Run: %s %s\n\n", detectedPM, joinArgs(pmArgs))
			os.Exit(1)

		case config.ModeRedirect:
			remapped := remap.RemapArgs(invokedPM, detectedPM, pmArgs)
			fmt.Fprintf(os.Stderr, "⚡ pmguard: redirecting %s → %s\n", invokedPM, detectedPM)
			return passthrough(string(detectedPM), remapped)
		}

		return nil
	},
}

func passthrough(pm string, args []string) error {
	bin, err := exec.LookPath(pm)
	if err != nil {
		return fmt.Errorf("could not find %s in PATH: %w", pm, err)
	}
	// Replace current process with the PM (no extra process spawned)
	return syscall.Exec(bin, append([]string{pm}, args...), os.Environ())
}

func joinArgs(args []string) string {
	result := ""
	for i, a := range args {
		if i > 0 {
			result += " "
		}
		result += a
	}
	return result
}

func Execute() {
	rootCmd.AddCommand(guardCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(installHooksCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
