package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/informatiqal/qlik-test-users-tickets/API/qlik"
	"github.com/informatiqal/qlik-test-users-tickets/Config"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Create test users inside Qlik Sense",
	Long:  "Create test users inside Qlik Sense",
}

var usersCreateCmd = (&cobra.Command{
	Use:     "create",
	Short:   "Create test users in Qlik",
	Example: ".\\qs_test_users users create --suffix something --users \"user1;user2;user3\" ",

	Run: func(cmd *cobra.Command, args []string) {
		userDirectorySuffix, err := cmd.Flags().GetString("suffix")
		if err != nil {
			log.Fatal(err.Error())
		}
		if userDirectorySuffix == "" {
			log.Fatal("--suffix must not be empty")
		}

		userNamesRaw, err := cmd.Flags().GetString("users")
		if err != nil {
			log.Fatal(err.Error())
		}
		if userNamesRaw == "" {
			log.Fatal("--users must not be empty")
		}

		userNames := strings.Split(userNamesRaw, ";")

		qlik.CreateTestUsersCmd(
			userNames,
			userDirectorySuffix,
		)

		os.Exit(0)
	},
})

func init() {
	config.NewConfig()
	usersCreateCmd.PersistentFlags().
		String("suffix", "", "Whats the user directory suffix under which the user(s) will be created? The final user directory will be TESTING_<suffix>.")
	usersCreateCmd.PersistentFlags().
		String("users", "", "List of semi-colon separated user names. Ideally wrap this value in double quotes.")

	usersCmd.AddCommand(usersCreateCmd)
}
