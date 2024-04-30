package cmd

import (
	"fmt"
	"os"
	"log"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	rootCmd = &cobra.Command{
		Use:   "ephemyral",
		Short: "A CLI for managing ephemyral tasks",
		Long:  `Ephemyral is an AI-powered CLI application for managing tasks that leverage machine learning models, including initialization, building, and testing of ML-driven workflows.`,
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

func GenerateMarkdownDocs() {
	err := doc.GenMarkdownTree(rootCmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
