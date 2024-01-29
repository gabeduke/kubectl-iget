/*
Copyright © 2024 Gabriel Duke <gabeduke@gmail.com>
*/
package cmd

import (
	"fmt"
	"github.com/gabeduke/kubectl-iget/pkg/kube"
	"github.com/gabeduke/kubectl-iget/pkg/ui"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "iget [resource] [flags]",
	Short: "Interactive Get for Kubernetes objects.",
	Long: `kubectl-iget is an interactive CLI plugin for Kubernetes that extends 'kubectl get'. 
           It allows users to dynamically construct object selectors and filters for querying Kubernetes resources.
           Use it just like 'kubectl get' and benefit from an enhanced, interactive experience.`,
	Args: cobra.MinimumNArgs(1), // Ensure at least one argument for the resource type
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		objectType := ""
		if len(args) > 0 {
			objectType = args[0] // Assuming the first argument is the object type
		}

		config := ui.UIConfig{
			ObjectType: objectType,
			// Set other config fields based on flags
		}

		kubeClient, _ := kube.NewKubeClient() // Add error handling
		uiManager := ui.NewUIManager(kubeClient, config)
		uiManager.Start()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubectl-iget.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".kubectl-iget" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kubectl-iget")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
