package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/jmcfarlane/notable/app"
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
			"Subject":  viper.GetString("create.subject"),
			"Tags":     viper.GetString("create.tags"),
			"Password": "Really?",
		}).Debug("Provided values")
		if !validateCreateParams() {
			log.Errorf("Subject has to be set")
		}
		var note app.Note
		note.Tags = viper.GetString("create.tags")
		note.Subject = viper.GetString("create.subject")
		note.Password = viper.GetString("create.password")

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

func validateCreateParams() bool {
	return "" != viper.GetString("create.subject")
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("subject", "s", "", "Subject for the note")
	createCmd.Flags().StringP("tags", "t", "", "Tags to associate")
	createCmd.Flags().StringP("password", "p", "", "Password to encrypt")

	viper.BindPFlag("create.subject", createCmd.Flags().Lookup("subject"))
	viper.BindPFlag("create.tags", createCmd.Flags().Lookup("tags"))
	viper.BindPFlag("create.password", createCmd.Flags().Lookup("password"))
}
