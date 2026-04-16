package cmd

import (
	"fmt"

	"github.com/istorry/pmguard/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or update pmguard configuration",
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		fmt.Println(cfg)
		return nil
	},
}

var configSetModeCmd = &cobra.Command{
	Use:   "set-mode [warn|redirect]",
	Short: "Set the behavior mode (warn or redirect)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mode := config.Mode(args[0])
		if mode != config.ModeWarn && mode != config.ModeRedirect {
			return fmt.Errorf("invalid mode %q — must be 'warn' or 'redirect'", mode)
		}

		cfg, _ := config.Load()
		cfg.Mode = mode
		if err := config.Save(cfg); err != nil {
			return err
		}

		fmt.Printf("✅ pmguard mode set to: %s\n", mode)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetModeCmd)
}
