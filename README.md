# Tasty Trade Open API CLI Written in Golang

I'm experimenting with the [Tasty Trade Open API](https://support.tastyworks.com/support/solutions/articles/43000700385-tastytrade-open-api) and writing a CLI in golang. I'm doing this with a couple of goals in mind:

* Use the Tasty Trade API to look for trading opportunities by automating some of the evaluation criteria I do manually
* Learn to write go code

Go is compiled, so it should be portable and able to run on most computers. The 1Password implentation is a bit trickier, and less universal. I'll get that to work in my environment and leave the portability for later.

## Security

I'm using 1Password to store my credentials and using the [1Password CLI](https://developer.1password.com/docs/cli/get-started/) to [update](https://developer.1password.com/docs/cli/item-edit) and 
[retrieve](https://developer.1password.com/docs/cli/reference/commands/read) the credentials in my code. There are no secrets hardcoded in the code or in environment variables. 

## Session flow

1. Login to your 1Password desktop application
1. Run `tasty login env [sbx | prod] <username>` to login. Your credentials will be saved to your 1Password vault
1. Run `tasty me` to verify that you've logged in sucessfully, and to see a list of your accounts
1. You're set to interact with your accounts using any of the commands below
1. Logout when you want to close your session by running `tasty logout`

## Commands

| Command | Description |
| ------- | ----------- |
|`tasty login [--prod]`| Gets session and remember tokens for environment. Defaults to `sbx` environment unless the `prod` option is set |
| `tasty logout` | Deactivates your session and remember tokens |
| `tasty me` | Your customer information |
| `tasty accounts` | List of your customer accounts |
| `tasty set-account` | Sets the account id that you want to interact with in subsequent commands
| `tasty get-account` | Returns the account number you set previously
| `tasty positions <account_number>` | Gets a list of your account postions. You must reference an account number if you haven't set one using `tasty set-account` |
| `tasty balances <account_number>` | Returns the monetary value of your account. You must reference an account number if you haven't set one using `tasty set-account` |
| `tasty vol-data <watchlist_name>` | Returns volatility data for each equity in the list included in the command |
| `tasty watchlist create <watchlist_name>` | Creates a user watchlist
| `tasty watchlist <watchlist_name> <JSON list of instruments to include in the list>` | Updates the watchlist with the list of entities passed in the command |
| `tasty watchlist get [all \| <watchlist_name>]` | Returns a list of all watchlists or the specific list by name. Defaults to `all` |
| `tasty watchlist delete <watchlist_name>` | Deletes the watchlist 
