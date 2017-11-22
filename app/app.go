package main

import (
	"diviner/app/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	rootCmd := cobra.Command{Use: "app"}
	rootCmd.PersistentFlags().StringP("ski", "s", "", "import private key from SKI")
	viper.BindPFlag("ski", rootCmd.PersistentFlags().Lookup("ski"))
	rootCmd.PersistentFlags().String("host", "localhost:50051", "diviner service host")
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))

	rootCmd.AddCommand(cmd.NewMemberCmd())
	rootCmd.AddCommand(cmd.NewEventCmd())
	rootCmd.AddCommand(cmd.NewMarketCmd())
	rootCmd.AddCommand(cmd.NewTxCmd())

	rootCmd.PersistentPreRun = func(xcmd *cobra.Command, args []string) {
		cmd.Init()
	}
	defer cmd.CloseConnection()
	rootCmd.Execute()

}
