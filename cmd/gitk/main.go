package main

import (
	"fmt"
	"os"

	gsdk "github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/mindkit-xyz/mindkit-gitk/internal/commands"
	"github.com/mindkit-xyz/mindkit-gitk/internal/storage"
	"github.com/mindkit-xyz/mindkit-gitk/internal/mindkit"
)

func main() {
	// Initialize configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.gitk")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize BNB Greenfield client
	greenfieldClient, err := gsdk.New(
		viper.GetString("greenfield.endpoint"),
		viper.GetString("greenfield.chainId"),
		viper.GetString("greenfield.privateKey"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing Greenfield client: %v\n", err)
		os.Exit(1)
	}

	// Initialize storage
	objStorage := storage.NewObjectStorage(greenfieldClient, 
		viper.GetString("storage.bucket"), 
		viper.GetString("storage.prefix"))
	refStorage := storage.NewReferenceStorage(greenfieldClient, 
		viper.GetString("storage.bucket"), 
		viper.GetString("storage.prefix"))

	// Initialize MindKit client
	mindkitClient := mindkit.NewClient(mindkit.Config{
		BaseURL: viper.GetString("mindkit.baseURL"),
		APIKey:  viper.GetString("mindkit.apiKey"),
	})
	ai := mindkit.NewAI(mindkitClient)

	// Create root command
	rootCmd := &cobra.Command{
		Use:   "gitk",
		Short: "Gitk is a Git implementation using BNB Greenfield storage",
		Long: `Gitk is a decentralized Git implementation that uses BNB Greenfield for storage
and integrates with MindKit's AI capabilities for enhanced project management.`,
	}

	// Add subcommands
	rootCmd.AddCommand(
		commands.NewInitCmd(objStorage, refStorage),
		commands.NewAddCommand(objStorage),
		commands.NewCommitCommand(objStorage, refStorage, ai),
		commands.NewPushCommand(objStorage, refStorage),
	)

	// Execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
