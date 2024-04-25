package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	rootCmd = &cobra.Command{
		Use:   "ephemyral",
		Short: "A CLI for managing ephemyral tasks",
		Long:  `Ephemyral is a CLI application for managing ephemyral tasks, including initialization, building, and testing.`,
	}
)

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ephemyral.yaml)")
}

func initConfig() {
	viper.SetConfigType("yaml")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		setDefaultConfig()
	}

	viper.AutomaticEnv()
	readConfig()
}

func setDefaultConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigName(".ephemyral")
}

func readConfig() {
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
