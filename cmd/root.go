/*
 */
package cmd

import (
	"context"
	"os"

	"github.com/docktermj/demo-entity-search/httpserver"
	"github.com/senzing/go-cmdhelping/cmdhelper"
	"github.com/senzing/go-cmdhelping/option"
	"github.com/senzing/go-cmdhelping/option/optiontype"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Short string = "demo-entity-search short description"
	Use   string = "demo-entity-search"
	Long  string = `
demo-entity-search long description.
    `
)

// ----------------------------------------------------------------------------
// Context variables
// ----------------------------------------------------------------------------

var SomethingToSay = option.ContextVariable{
	Arg:     "something-to-say",
	Default: option.OsLookupEnvString("SENZING_TOOLS_SOMETHING_TO_SAY", "Main says 'Hi!'"),
	Envar:   "SENZING_TOOLS_SOMETHING_TO_SAY",
	Help:    "Just a phrase to say [%s]",
	Type:    optiontype.String,
}

var ContextVariablesForMultiPlatform = []option.ContextVariable{
	option.HttpPort,
	option.ServerAddress,
}

var ContextVariables = append(ContextVariablesForMultiPlatform, ContextVariablesForOsArch...)

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	cmdhelper.Init(RootCmd, ContextVariables)
}

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Used in construction of cobra.Command
func PreRun(cobraCommand *cobra.Command, args []string) {
	cmdhelper.PreRun(cobraCommand, args, Use, ContextVariables)
}

// Used in construction of cobra.Command
func RunE(_ *cobra.Command, _ []string) error {
	ctx := context.Background()
	httpServer := &httpserver.HttpServerImpl{
		ServerAddress: viper.GetString(option.ServerAddress.Arg),
		ServerPort:    viper.GetInt(option.HttpPort.Arg),
	}
	return httpServer.Serve(ctx)
}

// Used in construction of cobra.Command
func Version() string {
	return cmdhelper.Version(githubVersion, githubIteration)
}

// ----------------------------------------------------------------------------
// Command
// ----------------------------------------------------------------------------

// RootCmd represents the command.
var RootCmd = &cobra.Command{
	Use:     Use,
	Short:   Short,
	Long:    Long,
	PreRun:  PreRun,
	RunE:    RunE,
	Version: Version(),
}
