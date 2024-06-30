package tastymenu

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func init() {

	cli.AppHelpTemplate = `NAME:
	{{.Name}} - {{.Usage}}
 USAGE:
	{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
	{{if len .Authors}}
 AUTHOR:
	{{range .Authors}}{{ . }}{{end}}
	{{end}}{{if .Commands}}
 GLOBAL OPTIONS:
	{{range .VisibleFlags}}{{.}}
	{{end}}{{end}}{{if .Copyright }}
 COMMANDS:
 {{range .Commands}}{{if not .HideHelp}} {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n " }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
COPYRIGHT:
	{{.Copyright}}
	{{end}}{{if .Version}}
 VERSION:
	{{.Version}}
	{{end}}
 WEBSITE: https://github.com/beckitrue/tasty-api/wiki
 THANK YOU: https://cli.urfave.org/
`
	cli.CommandHelpTemplate += "\nWEBSITE: https://github.com/beckitrue/tasty-api/wiki\nTHANK YOU: https://cli.urfave.org/\n"
	cli.SubcommandHelpTemplate += "\nWEBSITE: https://github.com/beckitrue/tasty-api/wiki\nTHANK YOU: https://cli.urfave.org/\n"

	cli.HelpFlag = &cli.BoolFlag{Name: "help", Aliases: []string{"h"}}
	cli.VersionFlag = &cli.BoolFlag{Name: "version", Aliases: []string{"v"}}

	//	cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
	//	 	fmt.Fprintf(w, "run tasty-api --help to see the help menu\n")
	//	}
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

func Menu() {

	var debug bool

	app := &cli.App{
		Name:     "tasty-menu",
		Version:  "v1.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Becki True",
				Email: "becki@beckitrue.com",
			},
		},
		Copyright: "(c) 2023 Me",
		HelpName:  "tasty-menu",
		Usage:     "cli for securely calling the Tastytrade API",
		UsageText: "tasty-menu [option] <cmd> [args]",
		// ArgsUsage: "[]",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Value:       false,
				Usage:       "enable debugging to see more output",
				Destination: &debug,
			},
		},
		Commands: []*cli.Command{
			{Name: "set-env",
				Category:    "config",
				Usage:       "set the environment you want to interact with: sbx or money",
				UsageText:   "set-env [sbx | money]",
				Description: "use this command to switch between your sbx and money accounts",
				ArgsUsage:   "[sbx | money]",
				Action: func(cCtx *cli.Context) error {
					if debug {
						fmt.Printf("debugging enabled\n")
					}
					env := (cCtx.Args().Get(0))
					if cCtx.NArg() == 0 {
						cli.ShowCommandHelp(cCtx, "You need to enter an environment: [sbx or money]")
					} else {
						if env != "sbx" && env != "money" {
							cli.ShowCommandHelp(cCtx, "You need to enter an environment: [sbx or money]")
						} else {
							fmt.Printf("You are working in the %s environment\n", env)
						}
					}

					return nil
				},
			},
			{
				Name:        "login",
				Aliases:     []string{"l"},
				Category:    "login",
				Usage:       "login to get session and remember tokens",
				UsageText:   "login",
				Description: "login to environment set using set-env command to get session token that is good for 24 hours or until you logout",
				ArgsUsage:   "[]",
				Action: func(cCtx *cli.Context) error {
					if debug {
						fmt.Printf("debugging enabled\n")
					}
					fmt.Fprintf(cCtx.App.Writer, "logging in...\n")
					return nil
				},
			},
			{
				Name:        "logout",
				Category:    "login",
				Usage:       "disables your session token",
				UsageText:   "logout",
				Description: "disables your session token, logging you out",
				Action: func(cCtx *cli.Context) error {
					if debug {
						fmt.Printf("debugging enabled\n")
					}
					fmt.Printf("logging out\n")
					return nil
				},
			},
			{
				Name:        "me",
				Aliases:     []string{"info"},
				Category:    "customer",
				Usage:       "returns your customer information",
				UsageText:   "me",
				Description: "returns your customer information in your sbx or money account",
				ArgsUsage:   "[]",
				Action: func(cCtx *cli.Context) error {
					if debug {
						fmt.Printf("debugging enabled\n")
					}
					fmt.Printf("getting your customer information...\n")
					return nil
				},
			},
			{
				Name:        "accounts",
				Aliases:     []string{"a"},
				Category:    "accounts",
				Usage:       "returns a list of your customer accounts",
				UsageText:   "accounts",
				Description: "returns a list of your customer accounts in your sbx or money account",
				ArgsUsage:   "[]",
				Action: func(cCtx *cli.Context) error {
					if debug {
						fmt.Printf("debugging enabled\n")
					}
					fmt.Printf("getting your accounts...\n")
					return nil
				},
			},
			{
				Name:        "set-account",
				Aliases:     []string{"sa"},
				Category:    "accounts",
				Usage:       "sets the account you want to interact with",
				UsageText:   "set-account [account id]",
				Description: "sets the account you want to interact with",
				ArgsUsage:   "[enter your account id]",
				Action: func(cCtx *cli.Context) error {
					if debug {
						fmt.Printf("debugging enabled\n")
					}
					if cCtx.NArg() > 0 {
						fmt.Printf("Setting account to: %s...\n", cCtx.Args().Get(0))
					} else {
						fmt.Printf("You didn't enter an acount number\n")
					}
					return nil
					// TODO: error checking on the input
				},
			},
			{
				Name:        "get-account",
				Aliases:     []string{"ga"},
				Category:    "accounts",
				Usage:       "gets the account you set to interact with",
				UsageText:   "get-account",
				Description: "gets the account you set to interact with",
				Action: func(cCtx *cli.Context) error {
					if debug {
						fmt.Printf("debugging enabled\n")
					}
					fmt.Printf("We're working with account id: %s\n", cCtx.Args().Get(0))
					return nil
					// TODO: error checking on the input
				},
			},
			{
				Name:        "get-positions",
				Aliases:     []string{"positions"},
				Category:    "accounts",
				Usage:       "lists the account positions",
				UsageText:   "get-positions [account id if you haven't set one]",
				Description: "lists the account positions for the account you set to interact with using the set-account command",
				ArgsUsage:   "[enter your account id]",
				Action: func(cCtx *cli.Context) error {
					if debug {
						fmt.Printf("debugging enabled\n")
					}
					fmt.Printf("Getting your positions in account id: %s...\n", cCtx.Args().Get(0))
					return nil
					// TODO: error checking on the input
				},
			},
		},
		CommandNotFound: func(cCtx *cli.Context, command string) {
			fmt.Fprintf(cCtx.App.Writer, "Command not found: %q\n", command)
		},
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			if isSubcommand {
				return err
			}

			fmt.Fprintf(cCtx.App.Writer, "WRONG: %#v\n", err)
			return err
		},
		Action: func(cCtx *cli.Context) error {
			cli.HandleExitCoder(errors.New("not an exit coder, though"))
			cli.ShowAppHelp(cCtx)
			cli.ShowVersion(cCtx)

			ec := cli.Exit("You didn't enter a command. Exiting", 86)
			fmt.Fprintf(cCtx.App.Writer, "%d", ec.ExitCode())
			fmt.Printf(" logged exit code\n")
			return ec
		},
	}

	app.Run(os.Args)
}
