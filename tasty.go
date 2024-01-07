package main

import (
	"errors"
	"github.com/beckitrue/tasty-api/httpclient"
	"github.com/beckitrue/tasty-api/jsondecode"
	"github.com/beckitrue/tasty-api/login"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"encoding/json"
	"bytes"

	"github.com/urfave/cli/v2"
)

// TODO: move to a config file
const (
	sbxVaultUser  = "op://SBX/tastytrade-sbx-api/username"
	sbxVaultToken = "op://SBX/tastytrade-sbx-api/credential"
)

// set the debug variables to default value
var debug bool

type ApiMsg struct {
	method string
	msg    string
	model  string
}

type WorkingEnv struct {
    Environment string `json:"environment"`
	Account string `json:account`
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
 WEBSITE: https://github.com/beckitrue/tasty-api/wiki
 THANK YOU: https://cli.urfave.org/
`
	cli.CommandHelpTemplate += "\nWEBSITE: https://github.com/beckitrue/tasty-api/wiki\nTHANK YOU: https://cli.urfave.org/\n"
	cli.SubcommandHelpTemplate += "\nWEBSITE: https://github.com/beckitrue/tasty-api/wiki\nTHANK YOU: https://cli.urfave.org/\n"

	cli.HelpFlag = &cli.BoolFlag{Name: "help", Aliases: []string{"h"}}
	cli.VersionFlag = &cli.BoolFlag{Name: "version", Aliases: []string{"V"}}

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
		UsageText: "tasty-api [option] <cmd> [flag]" ,
		// ArgsUsage: "[]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:               "debug",
				Category:           "",
				DefaultText:        "",
				FilePath:           "",
				Usage:              "displays additional messaging from HTTP requests and API calls that can be used to help identify issues",
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
			{   Name:        "set-env",
				Category:    "config",
				Usage:       "set the environment you want to interact with: sbx or money",
				UsageText:   "set-env [sbx | money]",
				Description: "use this command to switch between your sbx and money accounts",
				ArgsUsage:   "[sbx | money]",
				Action: func(cCtx *cli.Context) error {
					env := ((cCtx.Args().Get(0)))
					setEnv(env)
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
				Action: initialLogin,
			},
			{
				Name:        "logout",
				Category:    "login",
				Usage:       "disables your session token",
				UsageText:   "logout",
				Description: "disables your session token, logging you out",
				Action:      customerLogout,
			},
			{
				Name:        "me",
				Aliases:     []string{"info"},
				Category:    "customer",
				Usage:       "returns your customer information",
				UsageText:   "me [options]",
				Description: "returns your customer information in your sbx or money account",
				ArgsUsage:   "[]",
				Action:      customerInfo,
			},
			{
				Name:        "accounts",
				Aliases:     []string{"a"},
				Category:    "accounts",
				Usage:       "returns a list of your customer accounts",
				UsageText:   "accounts [--debug | -d]",
				Description: "returns a list of your customer accounts in your sbx or money account",
				ArgsUsage:   "[]",
				Action:      getAccounts,
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
					if cCtx.NArg() > 0 {
						fmt.Printf("OK, we'll be working with account id: %s\n", cCtx.Args().Get(0))
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
					fmt.Printf("We're working with account id: %s\n", "pull this from a file")
					return nil
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
					fmt.Printf("We're working with account id: %s\n", "pull this from a file")
					// TODO: write the function to call api
					return nil
				},
			},
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

func setEnv(env string) {

	// check for valid input
	if (env != "sbx") && (env != "money") {
		fmt.Printf("You didn't enter sbx or money. Follow the set-env command with either sbx or money\n")
	} else {
	    fmt.Printf("You are setting your working environment to: %s\n", env)

		settings := WorkingEnv{
			Environment: env,
			Account: "1234",
		}

		// write data to file
		if err := writeToJSON("/home/becki/environment.json", settings); err != nil {
			fmt.Println(err)
		}
    }
}

func writeToJSON(filename string, settings WorkingEnv) error {
	jsonFile, err := os.Create(filename)
	if err != nil {
       return fmt.Errorf("error creating JSON file: %v", err)
    }
	defer jsonFile.Close()

	var Marshal = func(v interface{}) (io.Reader, error) {
		b, err := json.Marshal(settings)

	    if err != nil {
		    return nil,err
	    }
		return bytes.NewReader(b), nil 
	}

	r, err := Marshal(settings)
	if err != nil {
		return err
	}

	_, err = io.Copy(jsonFile, r)

	
	// Close file and return any errors
	if err := jsonFile.Close(); err != nil {
        return fmt.Errorf("error closing JSON file: %v", err)
    }

	return err
}

func initialLogin(cCtx *cli.Context) error {
	// gets the session token and writes
	// it to 1Password

	login.GetSessionToken(debug)

	return nil
}

func customerLogout(cCtx *cli.Context) error {
	// disables the user's session token

	login.DisableToken(debug)

	return nil
}

func checkDebugFlag(cCtx *cli.Context) {
	if cCtx.Bool("debug") {
		fmt.Printf("Debug flag set\n")
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

	respString := Get(cmd)
	jsondecode.PrintMe(respString)

	return nil
}

func getAccounts(cCtx *cli.Context) error {
	accountList := ApiMsg{method: "GET", msg: "customers/me/accounts", model: "account"}
	cmd := accountList

	respString := Get(cmd)
	jsondecode.PrintDataAccounts(respString)

	return nil

}
