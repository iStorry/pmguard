package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var installHooksCmd = &cobra.Command{
	Use:   "install-hooks",
	Short: "Print shell hook code to activate pmguard",
	Long: `Outputs shell functions to wrap npm, pnpm, yarn, and bun.
Add this to your ~/.zshrc or ~/.bashrc:

  eval "$(pmguard install-hooks)"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print(shellHooks)
		return nil
	},
}

// shellHooks is the shell code printed by install-hooks.
// Users add: eval "$(pmguard install-hooks)" to their shell config.
const shellHooks = `
# pmguard hooks — auto-generated, do not edit manually
# To remove: delete the eval line from your shell config

_pmguard_wrap() {
  local pm="$1"
  shift
  pmguard guard "$pm" "$@"
}

npm()  { _pmguard_wrap npm  "$@"; }
pnpm() { _pmguard_wrap pnpm "$@"; }
yarn() { _pmguard_wrap yarn "$@"; }
bun()  { _pmguard_wrap bun  "$@"; }
`
