package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/Sirupsen/logrus"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new note",
	Long:  `Add a new note to notable server`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		log.Debugf("Creating a note on %s", viper.GetString("server"))
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("subject", "s", "", "Subject for the note")
	createCmd.Flags().StringP("tags", "t", "", "Tags to associate")
	createCmd.Flags().StringP("content", "c", "", "Content for the note")

	viper.BindPFlag("create.subject", createCmd.Flags().Lookup("subject"))
	viper.BindPFlag("create.tags", createCmd.Flags().Lookup("tags"))
	viper.BindPFlag("create.content", createCmd.Flags().Lookup("content"))
}
