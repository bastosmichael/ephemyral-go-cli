package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	rootCmd = &cobra.Command{
		Use:   "ephemyral",
		Short: "Ephemyral is an AI-powered CLI application for managing tasks that leverage machine learning models.",
		Long:  `Ephemyral is an AI-powered CLI application designed to streamline and optimize various tasks associated with machine learning projects. By leveraging machine learning models, Ephemyral provides a set of robust commands that simplify building, testing, and managing ML workflows. This tool is tailored for software engineers, data scientists, and anyone managing AI-driven projects.`,
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

func GenerateAndMergeDocs() {
	// Generate individual markdown files for each command
	err := doc.GenMarkdownTree(rootCmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}

	// Read and concatenate the generated docs
	var documentation strings.Builder
	files, _ := ioutil.ReadDir("./docs")

	for _, file := range files {
		if !file.IsDir() {
			content, err := ioutil.ReadFile("./docs/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			documentation.WriteString(string(content) + "\n")
		}
	}

	// Combine the static content with the generated documentation
	newReadme := documentation.String()

	// Write the new README to the file system, replacing the existing one
	err = ioutil.WriteFile("docs/README.md", []byte(newReadme), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
