package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"fmt"

	log "github.com/Sirupsen/logrus"
	app "github.com/jmcfarlane/notable/app"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new note",
	Long:  `Add a new note to notable server`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		log.Debugf("Creating a note on %s", viper.GetString("server"))
		log.WithFields(log.Fields{
			"Subject": viper.GetString("create.subject"),
			"Tags":    viper.GetString("create.tags"),
			"Content": viper.GetString("create.content"),
		}).Debug("Provided values")

		var note app.Note
		note.Content = viper.GetString("create.content")
		note.Tags = viper.GetString("create.tags")
		note.Subject = viper.GetString("create.subject")

		data, err := json.MarshalIndent(note, "", "  ")
		if err != nil {
			log.Errorf("Error while creating payload: %#v", err)
		}
		reader := bytes.NewReader(data)
		_, err = http.Post(fmt.Sprintf("%s/api/note/create", viper.GetString("server")), "application/json", reader)
		if err != nil {
			log.Errorf("Error while creating payload: %#v", err)
		}
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
