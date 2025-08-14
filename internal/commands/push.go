package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mindkit-xyz/mindkit-gitk/internal/storage"
)

func NewPushCommand(store *storage.ObjectStorage, refStore *storage.ReferenceStorage) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push [<remote>] [<branch>]",
		Short: "Update remote refs along with associated objects",
		Long: `Updates remote refs using local refs, while sending objects
necessary to complete the given refs.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get current HEAD reference
			headHash, err := refStore.GetReference(cmd.Context(), "HEAD")
			if err != nil {
				return fmt.Errorf("failed to get HEAD: %w", err)
			}

			// Push objects to BNB Greenfield
			if err := pushObjects(cmd.Context(), store, headHash); err != nil {
				return fmt.Errorf("failed to push objects: %w", err)
			}

			// Update remote reference
			remote := "origin"
			branch := "main"
			if len(args) > 0 {
				remote = args[0]
			}
			if len(args) > 1 {
				branch = args[1]
			}

			remoteRef := fmt.Sprintf("refs/remotes/%s/%s", remote, branch)
			if err := refStore.SetReference(cmd.Context(), remoteRef, headHash); err != nil {
				return fmt.Errorf("failed to update remote ref: %w", err)
			}

			fmt.Printf("Successfully pushed to %s/%s\n", remote, branch)
			return nil
		},
	}

	return cmd
}

func pushObjects(ctx context.Context, store *storage.ObjectStorage, hash string) error {
	// TODO: Implement pushing all objects reachable from the given hash
	// This would involve traversing the commit graph and ensuring all
	// objects are stored in BNB Greenfield
	return nil
}
