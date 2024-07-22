package jsondecode

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Me struct {
	Data struct {
		ID        string `json:"id"`
		FirstName string `json:"first-name"`
		LastName  string `json:"last-name"`
		Address   struct {
			StreetOne   string `json:"street-one"`
			City        string `json:"city"`
			StateRegion string `json:"state-region"`
			PostalCode  string `json:"postal-code"`
			Country     string `json:"country"`
			IsForeign   bool   `json:"is-foreign"`
			IsDomestic  bool   `json:"is-domestic"`
		} `json:"address"`
		MailingAddress struct {
			StreetOne   string `json:"street-one"`
			City        string `json:"city"`
			StateRegion string `json:"state-region"`
			PostalCode  string `json:"postal-code"`
			Country     string `json:"country"`
			IsForeign   bool   `json:"is-foreign"`
			IsDomestic  bool   `json:"is-domestic"`
		} `json:"mailing-address"`
		CustomerSuitability struct {
			ID                                int    `json:"id"`
			MaritalStatus                     string `json:"marital-status"`
			NumberOfDependents                int    `json:"number-of-dependents"`
			EmploymentStatus                  string `json:"employment-status"`
			Occupation                        string `json:"occupation"`
			EmployerName                      string `json:"employer-name"`
			JobTitle                          string `json:"job-title"`
			AnnualNetIncome                   int    `json:"annual-net-income"`
			NetWorth                          int    `json:"net-worth"`
			LiquidNetWorth                    int    `json:"liquid-net-worth"`
			StockTradingExperience            string `json:"stock-trading-experience"`
			CoveredOptionsTradingExperience   string `json:"covered-options-trading-experience"`
			UncoveredOptionsTradingExperience string `json:"uncovered-options-trading-experience"`
			FuturesTradingExperience          string `json:"futures-trading-experience"`
		} `json:"customer-suitability"`
		UsaCitizenshipType              string `json:"usa-citizenship-type"`
		IsForeign                       bool   `json:"is-foreign"`
		MobilePhoneNumber               string `json:"mobile-phone-number"`
		Email                           string `json:"email"`
		TaxNumberType                   string `json:"tax-number-type"`
		TaxNumber                       string `json:"tax-number"`
		BirthDate                       string `json:"birth-date"`
		ExternalID                      string `json:"external-id"`
		CitizenshipCountry              string `json:"citizenship-country"`
		SubjectToTaxWithholding         bool   `json:"subject-to-tax-withholding"`
		AgreedToMargining               bool   `json:"agreed-to-margining"`
		AgreedToTerms                   bool   `json:"agreed-to-terms"`
		HasIndustryAffiliation          bool   `json:"has-industry-affiliation"`
		HasPoliticalAffiliation         bool   `json:"has-political-affiliation"`
		HasListedAffiliation            bool   `json:"has-listed-affiliation"`
		IsProfessional                  bool   `json:"is-professional"`
		HasDelayedQuotes                bool   `json:"has-delayed-quotes"`
		HasPendingOrApprovedApplication bool   `json:"has-pending-or-approved-application"`
		IdentifiableType                string `json:"identifiable-type"`
		Person                          struct {
			ExternalID         string `json:"external-id"`
			FirstName          string `json:"first-name"`
			LastName           string `json:"last-name"`
			BirthDate          string `json:"birth-date"`
			CitizenshipCountry string `json:"citizenship-country"`
			UsaCitizenshipType string `json:"usa-citizenship-type"`
			MaritalStatus      string `json:"marital-status"`
			NumberOfDependents int    `json:"number-of-dependents"`
			EmploymentStatus   string `json:"employment-status"`
			Occupation         string `json:"occupation"`
			EmployerName       string `json:"employer-name"`
			JobTitle           string `json:"job-title"`
		} `json:"person"`
	} `json:"data"`
	Context string `json:"context"`
}

func PrintMe(data string) {
	dec := json.NewDecoder(strings.NewReader(data))
	var me Me

	dec.Decode(&me)

	// TODO: print more information
	fmt.Printf("Customer ID: %s\nFirst name: %s\nLast name:%s\nEmail:%s\n",
		me.Data.ID, me.Data.FirstName, me.Data.LastName, me.Data.Email)

}
