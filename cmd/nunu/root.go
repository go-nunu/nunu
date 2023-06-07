package nunu

import (
	"fmt"
	"github.com/go-nunu/nunu/internal/command/create"
	"github.com/go-nunu/nunu/internal/command/project"
	"github.com/go-nunu/nunu/internal/command/run"
	"github.com/go-nunu/nunu/internal/command/upgrade"
	"github.com/go-nunu/nunu/internal/command/wire"
	"github.com/go-nunu/nunu/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "nunu",
	Example: "nunu new demo-api",
	Short:   "\n _   _                   \n| \\ | |_   _ _ __  _   _ \n|  \\| | | | | '_ \\| | | |\n| |\\  | |_| | | | | |_| |\n|_| \\_|\\__,_|_| |_|\\__,_| \n \n" + "\x1B[38;2;66;211;146mA\x1B[39m \x1B[38;2;67;209;149mC\x1B[39m\x1B[38;2;68;206;152mL\x1B[39m\x1B[38;2;69;204;155mI\x1B[39m \x1B[38;2;70;201;158mt\x1B[39m\x1B[38;2;71;199;162mo\x1B[39m\x1B[38;2;72;196;165mo\x1B[39m\x1B[38;2;73;194;168ml\x1B[39m \x1B[38;2;74;192;171mf\x1B[39m\x1B[38;2;75;189;174mo\x1B[39m\x1B[38;2;76;187;177mr\x1B[39m \x1B[38;2;77;184;180mb\x1B[39m\x1B[38;2;78;182;183mu\x1B[39m\x1B[38;2;79;179;186mi\x1B[39m\x1B[38;2;80;177;190ml\x1B[39m\x1B[38;2;81;175;193md\x1B[39m\x1B[38;2;82;172;196mi\x1B[39m\x1B[38;2;83;170;199mn\x1B[39m\x1B[38;2;83;167;202mg\x1B[39m \x1B[38;2;84;165;205mg\x1B[39m\x1B[38;2;85;162;208mo\x1B[39m \x1B[38;2;86;160;211ma\x1B[39m\x1B[38;2;87;158;215mp\x1B[39m\x1B[38;2;88;155;218ml\x1B[39m\x1B[38;2;89;153;221mi\x1B[39m\x1B[38;2;90;150;224mc\x1B[39m\x1B[38;2;91;148;227ma\x1B[39m\x1B[38;2;92;145;230mt\x1B[39m\x1B[38;2;93;143;233mi\x1B[39m\x1B[38;2;94;141;236mo\x1B[39m\x1B[38;2;95;138;239mn\x1B[39m\x1B[38;2;96;136;243m.\x1B[39m",
	Version: fmt.Sprintf("\n _   _                   \n| \\ | |_   _ _ __  _   _ \n|  \\| | | | | '_ \\| | | |\n| |\\  | |_| | | | | |_| |\n|_| \\_|\\__,_|_| |_|\\__,_| \n \nNunu %s - Copyright (c) 2023 Nunu\nReleased under the MIT License.\n\n", version.Version),
}

func init() {
	rootCmd.AddCommand(project.NewCmd)
	rootCmd.AddCommand(create.CreateCmd)
	rootCmd.AddCommand(wire.WireCmd)
	rootCmd.AddCommand(run.RunCmd)
	rootCmd.AddCommand(upgrade.UpgradeCmd)
	create.CreateCmd.AddCommand(create.CreateHandlerCmd)
	create.CreateCmd.AddCommand(create.CreateServiceCmd)
	create.CreateCmd.AddCommand(create.CreateDaoCmd)
	create.CreateCmd.AddCommand(create.CreateModelCmd)
	create.CreateCmd.AddCommand(create.CreateAllCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
