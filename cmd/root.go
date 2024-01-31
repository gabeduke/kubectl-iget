/*
Copyright © 2024 Gabriel Duke <gabeduke@gmail.com>
*/
package cmd

import (
	"fmt"
	"github.com/gabeduke/kubectl-iget/pkg/kube"
	"github.com/gabeduke/kubectl-iget/pkg/ui"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
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

		// Configure the genericclioptions
		streams := genericiooptions.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		}
		cmd.SetOut(streams.Out)
		cmd.SetErr(streams.ErrOut)

		objectType := ""
		if len(args) > 0 {
			objectType = args[0] // Assuming the first argument is the object type
		}

		config := ui.UIConfig{
			Object: objectType,
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
	flags := pflag.NewFlagSet("kubectl-iget", pflag.ExitOnError)
	pflag.CommandLine = flags

	// This set of flags is the one used for the kubectl configuration such as:
	// cache-dir, cluster-name, namespace, kube-config, insecure, timeout, impersonate,
	// ca-file and so on
	kubeConfigFlags := genericclioptions.NewConfigFlags(false)

	// It is a set of flags related to a specific resource such as: label selector
	//(-L), --all-namespaces, --schema and so on.
	kubeResouceBuilderFlags := genericclioptions.NewResourceBuilderFlags()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubectl-iget.yaml)")
	flags.AddFlagSet(rootCmd.PersistentFlags())
	kubeConfigFlags.AddFlags(flags)

	// TODO why aren't these flags added?
	kubeResouceBuilderFlags.AddFlags(flags)
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
