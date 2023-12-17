package main

import (
	"example/user/tasty/httpclient"
	"example/user/tasty/jsondecode"
	"example/user/tasty/login"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

const (
	sbxVaultUser  = "op://Private/tastytrade-sbx-api/username"
	sbxVaultToken = "op://Private/tastytrade-sbx-api/credential"
)

type ApiMsg struct {
	method string
	msg    string
	model  string
}

func init() {
	cli.AppHelpTemplate += "\nCUSTOMIZED: you bet ur muffins\n"
	cli.CommandHelpTemplate += "\nYMMV\n"
	cli.SubcommandHelpTemplate += "\nor something\n"

	cli.HelpFlag = &cli.BoolFlag{Name: "halp"}
	cli.VersionFlag = &cli.BoolFlag{Name: "print-version", Aliases: []string{"V"}}

	cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
		fmt.Fprintf(w, "best of luck to you\n")
	}
	cli.VersionPrinter = func(cCtx *cli.Context) {
		fmt.Fprintf(cCtx.App.Writer, "version=%s\n", cCtx.App.Version)
	}
	cli.OsExiter = func(cCtx int) {
		fmt.Fprintf(cli.ErrWriter, "refusing to exit %d\n", cCtx)
	}
	cli.ErrWriter = io.Discard
	cli.FlagStringer = func(fl cli.Flag) string {
		return fmt.Sprintf("\t\t%s", fl.Names()[0])
	}
}

type hexWriter struct{}

func (w *hexWriter) Write(p []byte) (int, error) {
	for _, b := range p {
		fmt.Printf("%x", b)
	}
	fmt.Printf("\n")

	return len(p), nil
}

type genericType struct {
	s string
}

func (g *genericType) Set(value string) error {
	g.s = value
	return nil
}

func (g *genericType) String() string {
	return g.s
}

func main() {
	app := &cli.App{
		Name:     "tasty",
		Version:  "v1.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Becki True",
				Email: "becki@beckitrue.com",
			},
		},
		Copyright: "(c) 2023 Me",
		HelpName:  "tasty",
		Usage:     "cli for Tastytrade API",
		UsageText: "tasty- demonstrating the functionality of the API",
		ArgsUsage: "[debug]",
		Commands: []*cli.Command{
			&cli.Command{
				Name:        "login",
				Aliases:     []string{"l"},
				Category:    "login",
				Usage:       "login to get session token that is good for 24 hours or until you logout",
				UsageText:   "login [global commands] env [sbx | prod]",
				Description: "main command",
				ArgsUsage:   "[sbx | prod]",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "forever", Aliases: []string{"forevvarr"}},
				},
				Action: initialLogin,
			},
			&cli.Command{
				Name:        "accounts",
				Aliases:     []string{"acct"},
				Category:    "accounts",
				Usage:       "returns a list of your customer accounts",
				UsageText:   "accounts [global commands]",
				Description: "returns a list of your customer accounts in your sbx or prod account",
				ArgsUsage:   "[]",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "forever", Aliases: []string{"forevvarr"}},
				},
				Action: getAccounts,
			},
			&cli.Command{
				Name:        "me",
				Aliases:     []string{"info"},
				Category:    "accounts",
				Usage:       "returns your customer accounts",
				UsageText:   "me [global commands]",
				Description: "returns your customer information in your sbx or prod account",
				ArgsUsage:   "[]",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "forever", Aliases: []string{"forevvarr"}},
				},
				Before: func(cCtx *cli.Context) error {
					fmt.Fprintf(cCtx.App.Writer, "You are logged in to your sbx account\n")
					return nil
				},
				Action: customerInfo,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initialLogin(cCtx *cli.Context) error {
	login.GetSessionToken()

	return nil
}

func customerInfo(cCtx *cli.Context) error {
	customerMe := ApiMsg{method: "GET", msg: "customers/me", model: "account"}
	cmd := customerMe
	_, token := login.GetCreds(sbxVaultUser, sbxVaultToken)

	url := httpclient.CreateURL("sbx", cmd.msg)
	respString := httpclient.ApiCall(token, url, cmd.method)
	log.Println(respString)
	jsondecode.PrintMe(respString)

	return nil
}

func getAccounts(cCtx *cli.Context) error {
	accountList := ApiMsg{method: "GET", msg: "customers/me/accounts", model: "account"}
	cmd := accountList
	username, token := login.GetCreds(sbxVaultUser, sbxVaultToken)
	fmt.Printf("Hello %s\ttoken:%s\n", username, token)

	url := httpclient.CreateURL("sbx", cmd.msg)
	respString := httpclient.ApiCall(token, url, cmd.method)
	log.Println(respString)
	jsondecode.PrintDataAccounts(respString)

	return nil

}

// 	// // setting env to sbx for safety during testing
// 	env := "sbx"

// 	// login to TastyTrade to get the session token
// 	// and save it to 1Password for future API calls
// 	// https://developer.tastytrade.com/api-guides/sessions/
// 	// login.GetSessionToken()

// 	// get session token from 1Password
// 	username, token := login.GetCreds(sbxVaultUser, sbxVaultToken)
// 	fmt.Printf("Hello %s\ttoken:%s\n", username, token)

// 	// command we'll use for testing
// accountList := ApiMsg{method: "GET", msg: "customers/me/accounts", model: "account"}

// 	cmd := accountList
// 	log.Printf("running command: %s", cmd)

// 	// configure the API URL
// 	url := httpclient.CreateURL(env, cmd.msg)

// 	// make the API call
// 	respString := httpclient.ApiCall(token, url, cmd.method)
// 	log.Println(respString)

// 	// perform the appropriate JSON decoding and
// 	// print output
// 	jsondecode.PrintDataAccounts(respString)

// }
