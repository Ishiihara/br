package cmd

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/overvenus/br/backup"
	"github.com/spf13/cobra"
)

var defaultBacker *backup.Backer
var defaultBackerMu = sync.Mutex{}

// SetDefaultBacker sets the default backer for command line usage.
func SetDefaultBacker(ctx context.Context, pdAddrs string) {
	defaultBackerMu.Lock()
	defer defaultBackerMu.Unlock()
	print(pdAddrs)
	var err error
	defaultBacker, err = backup.NewBacker(ctx, pdAddrs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// GetDefaultBacker returns the default backer for command line usage.
func GetDefaultBacker() *backup.Backer {
	defaultBackerMu.Lock()
	defer defaultBackerMu.Unlock()
	return defaultBacker
}

// NewMetaCommand return a meta subcommand.
func NewMetaCommand() *cobra.Command {
	meta := &cobra.Command{
		Use:   "meta <subcommand>",
		Short: "show meta data of a cluster",
	}
	meta.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "show cluster version",
		Run: func(cmd *cobra.Command, _ []string) {
			backer := GetDefaultBacker()
			v, err := backer.GetClusterVersion()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cmd.Println(v)
		},
	})
	meta.AddCommand(&cobra.Command{
		Use:   "safepoint",
		Short: "show the current GC safepoint of cluster",
		Run: func(cmd *cobra.Command, _ []string) {
			backer := GetDefaultBacker()
			sp, err := backer.GetGCSaftPoint()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cmd.Printf("Timestamp { Physical: %d, Logical: %d }\n",
				sp.Physical, sp.Logical)
		},
	})
	return meta
}
