package accounts

import (
	"fmt"
	"tastyapi/login"
)

// get account information
func GetAccountInfo() {
	info := login.GetStoredToken("sbx")
	fmt.Println(info)
}
