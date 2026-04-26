package upgrade

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/go-nunu/nunu/config"
	"github.com/spf13/cobra"
)

var CmdUpgrade = &cobra.Command{
	Use:     "upgrade",
	Short:   "Upgrade the nunu command.",
	Long:    "Upgrade the nunu command.",
	Example: "nunu upgrade",
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Printf("go install %s\n", config.NunuCmd)
		cmd := exec.Command("go", "install", config.NunuCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("go install %s: %w", config.NunuCmd, err)
		}
		fmt.Printf("\n🎉 Nunu upgrade successfully!\n\n")
		return nil
	},
}
