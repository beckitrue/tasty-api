package main

import (
	"errors"
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

// TODO: move to a config file
const (
	sbxVaultUser  = "op://Private/tastytrade-sbx-api/username"
	sbxVaultToken = "op://Private/tastytrade-sbx-api/credential"
)

// set the prod and debug variables to default values
var prod bool = false
var debug bool = false

type ApiMsg struct {
	method string
	msg    string
	model  string
}

func init() {
	cli.AppHelpTemplate = `NAME:
	{{.Name}} - {{.Usage}}
 USAGE:
	{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
	{{if len .Authors}}
 AUTHOR:
	{{range .Authors}}{{ . }}{{end}}
	{{end}}{{if .Commands}}
 COMMANDS:
 {{range .Commands}}{{if not .HideHelp}} {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n " }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
 GLOBAL OPTIONS:
	{{range .VisibleFlags}}{{.}}
	{{end}}{{end}}{{if .Copyright }}
 COPYRIGHT:
	{{.Copyright}}
	{{end}}{{if .Version}}
 VERSION:
	{{.Version}}
	{{end}}
 WEBSITE: https://github.com/beckitrue/tasty-api
 THANK YOU: https://cli.urfave.org/
`
	cli.CommandHelpTemplate += "\nWEBSITE: https://github.com/beckitrue/tasty-api\nTHANK YOU: https://cli.urfave.org/\n"
	cli.SubcommandHelpTemplate += "\nWEBSITE: https://github.com/beckitrue/tasty-api\nTHANK YOU: https://cli.urfave.org/\n"

	cli.HelpFlag = &cli.BoolFlag{Name: "help", Aliases: []string{"h"}}
	cli.VersionFlag = &cli.BoolFlag{Name: "print-version", Aliases: []string{"V"}}

	// cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
	// 	fmt.Fprintf(w, "best of luck to you\n")
	// }
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
			{
				Name:  "Becki True",
				Email: "becki@beckitrue.com",
			},
		},
		Copyright: "(c) 2023 Me",
		HelpName:  "tasty",
		Usage:     "cli for securely calling the Tastytrade API",
		UsageText: "tasty- demonstrating the functionality of the API",
		// ArgsUsage: "[]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:               "debug",
				Category:           "",
				DefaultText:        "",
				FilePath:           "",
				Usage:              "displays additional messaging from HTTP requests and API calls",
				Required:           false,
				Hidden:             false,
				HasBeenSet:         false,
				Value:              false,
				Destination:        new(bool),
				Aliases:            []string{"d"},
				EnvVars:            []string{},
				Count:              new(int),
				DisableDefaultText: false,
				Action: func(*cli.Context, bool) error {
					debug = true
					return nil
				},
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "login",
				Aliases:     []string{"l"},
				Category:    "login",
				Usage:       "login to get session token",
				UsageText:   "login --prod for live account (defaults to sbx if flag is unset)",
				Description: "login to get session token that is good for 24 hours or until you logout",
				ArgsUsage:   "[]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:               "prod",
						Category:           "",
						DefaultText:        "",
						FilePath:           "",
						Usage:              "set this flag if you want to connect to your live account",
						Required:           false,
						Hidden:             false,
						HasBeenSet:         false,
						Value:              false,
						Destination:        new(bool),
						Aliases:            []string{"p"},
						EnvVars:            []string{},
						Count:              new(int),
						DisableDefaultText: false,
						Action: func(*cli.Context, bool) error {
							prod = true
							return nil
						},
					},
				},
				Action: initialLogin,
			},
			{
				Name:        "me",
				Aliases:     []string{"info"},
				Category:    "customer",
				Usage:       "returns your customer information",
				UsageText:   "me [options]",
				Description: "returns your customer information in your sbx or prod account",
				ArgsUsage:   "[--prod, --debug ]",
				Action:      customerInfo,
			},
			{
				Name:        "accounts",
				Aliases:     []string{"a"},
				Category:    "accounts",
				Usage:       "returns a list of your customer accounts",
				UsageText:   "accounts [--debug | -d]",
				Description: "returns a list of your customer accounts in your sbx or prod account",
				ArgsUsage:   "[]",
				Action:      getAccounts,
			},
		},
		Action: func(cCtx *cli.Context) error {
			cli.DefaultAppComplete(cCtx)
			cli.HandleExitCoder(errors.New("not an exit coder, though"))
			cli.ShowAppHelp(cCtx)
			cli.ShowCommandCompletions(cCtx, "nope")
			cli.ShowCommandHelp(cCtx, "also-nope")
			cli.ShowCompletions(cCtx)
			cli.ShowSubcommandHelp(cCtx)
			cli.ShowVersion(cCtx)

			cCtx.App.Setup()
			fmt.Printf("%#v\n", cCtx.App.VisibleCategories())
			fmt.Printf("%#v\n", cCtx.App.VisibleCommands())
			fmt.Printf("%#v\n", cCtx.App.VisibleFlags())

			ec := cli.Exit("ohwell", 86)
			fmt.Fprintf(cCtx.App.Writer, "%d", ec.ExitCode())
			fmt.Printf("made it!\n")
			return ec
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initialLogin(cCtx *cli.Context) error {
	// gets the session token and writes
	// it to 1Password

	login.GetSessionToken(debug)

	return nil
}

func checkDebugFlag(cCtx *cli.Context) {
	if cCtx.Bool("debug") {
		fmt.Printf("Debug flag set\n")
	}
}

func checkProdFlag(cCtx *cli.Context) {
	if cCtx.Bool("prod") {
		fmt.Printf("You are in your live accountt\n")
	} else {
		fmt.Printf("You are in your sbx account\n")
	}
}

func Get(cmd ApiMsg) (response string) {
	_, token := login.GetCreds(sbxVaultUser, sbxVaultToken)

	url := httpclient.CreateURL("sbx", cmd.msg)
	respString := httpclient.ApiCall(token, url, cmd.method, debug)

	// Debug by printing the whole response from the API call
	if debug {
		log.Println(respString)
	}

	return respString
}

func customerInfo(cCtx *cli.Context) error {
	customerMe := ApiMsg{method: "GET", msg: "customers/me", model: "account"}
	cmd := customerMe

	checkProdFlag(cCtx)

	respString := Get(cmd)
	jsondecode.PrintMe(respString)

	return nil
}

func getAccounts(cCtx *cli.Context) error {
	accountList := ApiMsg{method: "GET", msg: "customers/me/accounts", model: "account"}
	cmd := accountList

	checkProdFlag(cCtx)

	respString := Get(cmd)
	jsondecode.PrintDataAccounts(respString)

	return nil

}
