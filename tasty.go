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
	cli.AppHelpTemplate = `NAME:
	{{.Name}} - {{.Usage}}
 USAGE:
	{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
	{{if len .Authors}}
 AUTHOR:
	{{range .Authors}}{{ . }}{{end}}
	{{end}}{{if .Commands}}
 COMMANDS:
 {{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
 GLOBAL OPTIONS:
	{{range .VisibleFlags}}{{.}}
	{{end}}{{end}}{{if .Copyright }}
 COPYRIGHT:
	{{.Copyright}}
	{{end}}{{if .Version}}
 VERSION:
	{{.Version}}
	{{end}}
`
	cli.CommandHelpTemplate += "\nYMMV\n"
	cli.SubcommandHelpTemplate += "\nor something\n"

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
		// ArgsUsage: "[debug]",
		Commands: []*cli.Command{
			{
				Name:        "login",
				Aliases:     []string{"l"},
				Category:    "login",
				Usage:       "login to get session token",
				UsageText:   "login (defaults to sbx)",
				Description: "login to get session token that is good for 24 hours or until you logout",
				ArgsUsage:   "[env [sbx | prod]]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "debug",
						Aliases:  []string{"d"},
						Usage:    "displays additional messaging from HTTP requests and API calls",
						Category: "Debug",
					},
				},
				Action: initialLogin,
			},
			{
				Name:        "accounts",
				Aliases:     []string{"a"},
				Category:    "accounts",
				Usage:       "returns a list of your customer accounts",
				UsageText:   "accounts [--debug | -d]",
				Description: "returns a list of your customer accounts in your sbx or prod account",
				ArgsUsage:   "[]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "debug",
						Aliases:  []string{"d"},
						Usage:    "displays additional messaging from HTTP requests and API calls",
						Category: "Debug",
					},
				},
				Action: getAccounts,
			},
			{
				Name:        "me",
				Aliases:     []string{"info"},
				Category:    "customer",
				Usage:       "returns your customer information",
				UsageText:   "me [--debug | -d]",
				Description: "returns your customer information in your sbx or prod account",
				ArgsUsage:   "[]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "debug",
						Aliases:  []string{"d"},
						Usage:    "displays additional messaging from HTTP requests and API calls",
						Category: "Debug",
					},
				},
				Before: func(cCtx *cli.Context) error {
					fmt.Fprintf(cCtx.App.Writer, "You are logged in to your sbx account\n")
					return nil
				},
				Action: customerInfo,
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

	if cCtx.Bool("debug") {
		fmt.Printf("Debug flag set\n")
	}

	login.GetSessionToken()

	return nil
}

func Get(cmd ApiMsg) (response string) {
	_, token := login.GetCreds(sbxVaultUser, sbxVaultToken)

	url := httpclient.CreateURL("sbx", cmd.msg)
	respString := httpclient.ApiCall(token, url, cmd.method)

	// TODO: Debug
	// log.Println(respString)

	return respString
}

func customerInfo(cCtx *cli.Context) error {
	customerMe := ApiMsg{method: "GET", msg: "customers/me", model: "account"}
	cmd := customerMe

	if cCtx.Bool("debug") {
		fmt.Printf("Debug flag set\n")
	}

	respString := Get(cmd)
	jsondecode.PrintMe(respString)

	return nil
}

func getAccounts(cCtx *cli.Context) error {
	accountList := ApiMsg{method: "GET", msg: "customers/me/accounts", model: "account"}
	cmd := accountList

	if cCtx.Bool("debug") {
		fmt.Printf("Debug flag set\n")
	}

	respString := Get(cmd)
	jsondecode.PrintDataAccounts(respString)

	return nil

}
