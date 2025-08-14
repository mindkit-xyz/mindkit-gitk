package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/mindkit-xyz/mindkit-gitk/internal/storage"
)

func NewAddCommand(store *storage.ObjectStorage) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [<path>...]",
		Short: "Add file contents to the index",
		Long: `Updates the index using the current content found in the working tree,
preparing the content staged for the next commit.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("nothing specified, nothing added")
			}

			for _, path := range args {
				if err := addPath(cmd.Context(), store, path); err != nil {
					return fmt.Errorf("failed to add %s: %w", path, err)
				}
			}

			return nil
		},
	}

	return cmd
}

func addPath(ctx context.Context, store *storage.ObjectStorage, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return addDirectory(ctx, store, path)
	}

	return addFile(ctx, store, path)
}

func addDirectory(ctx context.Context, store *storage.ObjectStorage, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if err := addPath(ctx, store, path); err != nil {
			return err
		}
	}

	return nil
}

func addFile(ctx context.Context, store *storage.ObjectStorage, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Store the file content in BNB Greenfield
	hash := calculateHash(data)
	if err := store.Store(ctx, hash, data); err != nil {
		return err
	}

	fmt.Printf("added '%s'\n", path)
	return nil
}

func calculateHash(data []byte) string {
	// Implement SHA-1 hash calculation for Git compatibility
	// TODO: Implement proper Git hash calculation
	return fmt.Sprintf("%x", data[:20])
}
