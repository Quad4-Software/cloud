package cmd

import (
	"fmt"

	"ocwrapper/common"
	quad4 "ocwrapper/quad4"
	quad4Config "ocwrapper/quad4/config"
	wrapper "ocwrapper/wrapper"
	wrapperConfig "ocwrapper/wrapper/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ocwrapper",
	Short: "ocwrapper is a wrapper for quad4 server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			fmt.Printf("error executing help command: %v\n", err)
		}
	},
}

func serveCmd() *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Starts the server",
		Run: func(cmd *cobra.Command, args []string) {
			common.Wg.Add(2)

			// set configs
			binFlag, _ := cmd.Flags().GetString("bin")
			quad4Config.Set("bin", binFlag)
			urlFlag, _ := cmd.Flags().GetString("url")
			quad4Config.Set("url", urlFlag)
			retryFlag, _ := cmd.Flags().GetString("retry")
			quad4Config.Set("retry", retryFlag)
			adminUsernameFlag, _ := cmd.Flags().GetString("admin-username")
			quad4Config.Set("adminUsername", adminUsernameFlag)
			adminPasswordFlag, _ := cmd.Flags().GetString("admin-password")
			quad4Config.Set("adminPassword", adminPasswordFlag)

			skipQuad4CloudRunFlag, _ := cmd.Flags().GetBool("skip-Quad4 Cloud-run")
			if !skipQuad4CloudRunFlag {
				go quad4.Start(nil)
			}
			portFlag, _ := cmd.Flags().GetString("port")
			go wrapper.Start(portFlag)
		},
	}

	// serve command args
	serveCmd.Flags().SortFlags = false
	serveCmd.Flags().StringP("bin", "", quad4Config.Get("bin"), "Full quad4 binary path")
	serveCmd.Flags().StringP("url", "", quad4Config.Get("url"), "quad4 server url")
	serveCmd.Flags().StringP("retry", "", quad4Config.Get("retry"), "Number of retries to start quad4 server")
	serveCmd.Flags().StringP("port", "p", wrapperConfig.Get("port"), "Wrapper API server port")
	serveCmd.Flags().StringP("admin-username", "", "", "admin username for quad4 server")
	serveCmd.Flags().StringP("admin-password", "", "", "admin password for quad4 server")
	serveCmd.Flags().Bool("skip-Quad4 Cloud-run", false, "Skip running quad4 server")

	return serveCmd
}

// Execute executes the command
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(serveCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error executing command: %v\n", err)
	}
}
