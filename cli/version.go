package cli

import (
	"fmt"
	"runtime"

	"github.com/hxia043/nest/types"

	"github.com/spf13/cobra"
)

func initVersionCommand() {
	/*
		version = &cobra.Command{
			Use:     "version",
			Short:   "Show the swic version information",
			Long:    "Show the swic version information. The version information is expected to follow semantic versionsing (https://semver.org/)",
			Example: "swic version",
			Run:     versionExecute,
		}
	*/
}

func versionExecute(cmd *cobra.Command, args []string) {
	fmt.Printf("version.BuildInfo{Version: \"%s\", GoVersion: \"%s\"}\n", types.Version, runtime.Version())
}
