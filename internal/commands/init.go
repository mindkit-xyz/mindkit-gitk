package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/mindkit-xyz/mindkit-gitk/internal/storage"
)

func NewInitCmd(store *storage.ObjectStorage, refStore *storage.ReferenceStorage) *cobra.Command {
	var bucketName string

	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Initialize a new Gitk repository",
		Long: `Initialize a new Gitk repository that uses BNB Greenfield for storage. 
This command creates a new .gitk directory with the repository configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			// Create .gitk directory
			gitkDir := filepath.Join(path, ".gitk")
			if err := os.MkdirAll(gitkDir, 0755); err != nil {
				return fmt.Errorf("failed to create .gitk directory: %w", err)
			}

			// Create config file
			config := fmt.Sprintf(`[storage]
bucket = "%s"
`, bucketName)

			if err := os.WriteFile(filepath.Join(gitkDir, "config"), []byte(config), 0644); err != nil {
				return fmt.Errorf("failed to write config file: %w", err)
			}

			// Initialize empty HEAD reference
			if err := refStore.SetReference(context.Background(), "HEAD", "refs/heads/main"); err != nil {
				return fmt.Errorf("failed to initialize HEAD reference: %w", err)
			}

			fmt.Printf("Initialized empty Gitk repository in %s\n", gitkDir)
			return nil
		},
	}

	cmd.Flags().StringVarP(&bucketName, "bucket", "b", "", "BNB Greenfield bucket name (required)")
	cmd.MarkFlagRequired("bucket")

	return cmd
}
