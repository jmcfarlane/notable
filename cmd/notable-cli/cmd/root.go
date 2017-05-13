package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var (
	cfgFile       string
	notableServer string
	printDebug    bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "notable-cli",
	Short: "A notable command line client",
	Long: `A command line client for notable. Notable is a
simple note taking application`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.notable-cli.yaml)")
	RootCmd.PersistentFlags().StringVar(&notableServer, "server", "http://localhost:8080", "Base url for notable server")
	RootCmd.PersistentFlags().BoolVar(&printDebug, "debug", false, "Enable to print debug")

	viper.BindPFlag("server", RootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".notable-cli")   // name of config file (without extension)
	viper.AddConfigPath("$HOME/.notable") // adding notable config directory as first search path
	viper.AddConfigPath(".")              // adding current directory as second search path
	viper.AddConfigPath("$HOME")          // adding home directory as third search path
	viper.AutomaticEnv()                  // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
}
