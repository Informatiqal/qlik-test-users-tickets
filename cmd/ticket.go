package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/informatiqal/qlik-test-users-tickets/Config"
	"github.com/informatiqal/qlik-test-users-tickets/qlik"
	"github.com/spf13/cobra"
)

var ticketCmd = &cobra.Command{
	Use:   "ticket",
	Short: "Web tickets operations",
	Long:  "Web tickets operations",
}

var ticketCreateCmd = (&cobra.Command{
	Use:     "create",
	Short:   "Create web ticket for specific user",
	Example: ".\\qs_test_users ticket create --suffix something --userId myUserId --vp smt --proxyId 6666666-7777-8888-9999-000000000000 --attributes [{\"group\": \"something\"},...]",
	Run: func(cmd *cobra.Command, args []string) {
		userDirectorySuffix, err := cmd.Flags().GetString("suffix")
		if err != nil {
			log.Fatal(err.Error())
		}

		userId, err := cmd.Flags().GetString("userId")
		if err != nil {
			log.Fatal(err.Error())
		}

		proxyId, err := cmd.Flags().GetString("proxyId")
		if err != nil {
			log.Fatal(err.Error())
		}

		virtualProxyPrefix, err := cmd.Flags().GetString("vp")
		if err != nil {
			log.Fatal(err.Error())
		}

		attributes, err := cmd.Flags().GetString("attributes")
		if err != nil {
			log.Fatal(err.Error())
		}
		// set to empty array string if no attributes were provided
		if attributes == "" {
			attributes = "[]"
		}

		jsonOutput, err := cmd.Flags().GetBool("json")
		if err != nil {
			log.Fatal(err.Error())
		}

		jsonFormatOutput, err := cmd.Flags().GetBool("jsonf")
		if err != nil {
			log.Fatal(err.Error())
		}

		ticket, err := qlik.CreateTicketForUser(
			userId,
			"TESTING_"+userDirectorySuffix,
			virtualProxyPrefix,
			attributes,
			proxyId,
		)
		if err != nil {
			log.Fatal(err.Error())
		}

		if !jsonOutput && !jsonFormatOutput {
			fmt.Println("Ticket  : " + ticket.Ticket)
			fmt.Println("QMC Link: " + ticket.Links.Qmc)
			fmt.Println("Hub Link: " + ticket.Links.Hub)
		} else {
			var jsonString []byte

			if jsonOutput {
				jsonString, err = json.Marshal(ticket)
				if err != nil {
					log.Fatal(err.Error())
				}
				fmt.Println(string(jsonString))
			} else {
				jsonString, err = json.MarshalIndent(ticket, "", "  ")
				if err != nil {
					log.Fatal(err.Error())
				}
				fmt.Println(string(jsonString))
			}

		}

		os.Exit(0)
	},
})

func init() {
	config.NewConfig()
	ticketCreateCmd.PersistentFlags().
		String("suffix", "", "Whats the user directory suffix for the requested user?")
	ticketCreateCmd.MarkPersistentFlagRequired("suffix")

	ticketCreateCmd.PersistentFlags().
		String("userId", "", "Whats the userId?")
	ticketCreateCmd.MarkPersistentFlagRequired("userId")

	ticketCreateCmd.PersistentFlags().
		String("proxyId", "", "(optional) On which proxy service the ticket should be created.")
	ticketCreateCmd.MarkPersistentFlagRequired("proxyId")

	ticketCreateCmd.PersistentFlags().
		String("vp", "", "(optional) Against which virtual proxy the ticket should be created. If not provided then \"/\" default will be used")

	ticketCreateCmd.PersistentFlags().
		String("attributes", "", "(optional) Additional attributes (JSON format) to be associated with the ticket")

	ticketCreateCmd.PersistentFlags().
		BoolP("json", "", false, "(optional) Print the ticket data in JSON format")

	ticketCreateCmd.PersistentFlags().
		BoolP("jsonf", "", false, "(optional) Print the ticket data in formatted JSON format")

	ticketCmd.AddCommand(ticketCreateCmd)
}
