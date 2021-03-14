package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

var loggingLevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "sensoplas",
	// RunE: func(cmd *cobra.Command, args []string) error {
	// 	fmt.Println("we okay in here")
	// 	return nil
	// },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "conf.yaml", "config file (default is ./conf.yaml)")
	rootCmd.PersistentFlags().StringVarP(&loggingLevel, "log", "l", "info", "debug logging (default is false)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
