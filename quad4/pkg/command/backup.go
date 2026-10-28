package command

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Quad4-Software/cloud/quad4/pkg/backup"
	"github.com/Quad4-Software/cloud/quad4/pkg/register"
	"github.com/Quad4-Software/cloud/pkg/config"
	"github.com/Quad4-Software/cloud/pkg/config/configlog"
	"github.com/Quad4-Software/cloud/pkg/config/parser"
	decomposedbs "github.com/opencloud-eu/reva/v2/pkg/storage/fs/decomposed/blobstore"
	decomposeds3bs "github.com/opencloud-eu/reva/v2/pkg/storage/fs/decomposeds3/blobstore"
)

// BackupCommand is the entrypoint for the backup command
func BackupCommand(cfg *config.Config) *cobra.Command {
	bckCmd := &cobra.Command{
		Use:   "backup",
		Short: "Quad4 Cloud backup functionality",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return configlog.ReturnError(parser.ParseConfig(cfg, true))
		},
	}
	bckCmd.AddCommand(ConsistencyCommand(cfg))
	return bckCmd
}

// ConsistencyCommand is the entrypoint for the consistency Command
func ConsistencyCommand(cfg *config.Config) *cobra.Command {
	consCmd := &cobra.Command{
		Use:   "consistency",
		Short: "check backup consistency",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				bs  backup.ListBlobstore
				err error
			)
			basePath, _ := cmd.Flags().GetString("basepath")
			blobstoreFlag, _ := cmd.Flags().GetString("blobstore")
			switch blobstoreFlag {
			case "decomposeds3":
				bs, err = decomposeds3bs.New(
					cfg.StorageUsers.Drivers.DecomposedS3.Endpoint,
					cfg.StorageUsers.Drivers.DecomposedS3.Region,
					cfg.StorageUsers.Drivers.DecomposedS3.Bucket,
					cfg.StorageUsers.Drivers.DecomposedS3.AccessKey,
					cfg.StorageUsers.Drivers.DecomposedS3.SecretKey,
					decomposeds3bs.Options{},
				)
			case "decomposed":
				bs, err = decomposedbs.New(basePath)
			case "none":
				bs = nil
			default:
				err = errors.New("blobstore type not supported")
			}
			if err != nil {
				fmt.Println(err)
				return err
			}
			fail, _ := cmd.Flags().GetBool("fail")
			if err := backup.CheckProviderConsistency(basePath, bs, fail); err != nil {
				fmt.Println(err)
				return err
			}

			return nil
		},
	}
	consCmd.Flags().StringP("basepath", "p", "", "the basepath of the decomposedfs (e.g. /var/tmp/quad4/storage/users)")
	_ = consCmd.MarkFlagRequired("basepath")
	consCmd.Flags().StringP("blobstore", "b", "decomposed", "the blobstore type. Can be (none, decomposed, decomposeds3). Default decomposed")
	consCmd.Flags().Bool("fail", false, "exit with non-zero status if consistency check fails")
	return consCmd
}

func init() {
	register.AddCommand(BackupCommand)
}
