package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/mindkit-xyz/mindkit-gitk/internal/storage"
	"github.com/mindkit-xyz/mindkit-gitk/internal/mindkit"
)

func NewCommitCommand(store *storage.ObjectStorage, refStore *storage.ReferenceStorage, ai *mindkit.AI) *cobra.Command {
	var message string
	var useAI bool

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Record changes to the repository",
		Long: `Creates a new commit containing the current contents of the index and
the given log message describing the changes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if useAI {
				// Generate commit message using AI
				msg, err := ai.GenerateCommitMessage(cmd.Context(), "")
				if err != nil {
					return fmt.Errorf("failed to generate commit message: %w", err)
				}
				message = msg
			} else if message == "" {
				return fmt.Errorf("please provide a commit message")
			}

			// Create commit object
			commit := &Commit{
				Message:   message,
				Author:    "user@example.com", // TODO: Get from config
				Date:     time.Now(),
				Parent:   "", // TODO: Get current HEAD
			}

			// Store commit object
			commitData := commit.Serialize()
			hash := calculateHash(commitData)
			if err := store.Store(cmd.Context(), hash, commitData); err != nil {
				return fmt.Errorf("failed to store commit: %w", err)
			}

			// Update HEAD reference
			if err := refStore.SetReference(cmd.Context(), "HEAD", hash); err != nil {
				return fmt.Errorf("failed to update HEAD: %w", err)
			}

			fmt.Printf("[%s] %s\n", hash[:7], message)
			return nil
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "commit message")
	cmd.Flags().BoolVar(&useAI, "ai", false, "use AI to generate commit message")

	return cmd
}

type Commit struct {
	Message string
	Author  string
	Date    time.Time
	Parent  string
}

func (c *Commit) Serialize() []byte {
	// TODO: Implement proper Git commit object serialization
	return []byte(fmt.Sprintf("commit %s\n%s\n%s\n%s",
		c.Parent,
		c.Author,
		c.Date.Format(time.RFC3339),
		c.Message,
	))
}
