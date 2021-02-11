package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	//third party deps
	"github.com/gookit/color"
)

//command line args
func commandLineUsage() {
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = commandLineUsage
	emailPtr := flag.String("email", "", "Email Address to Target")
	//todo add verbose flag to output debug stuff
	flag.Parse()
	//validate it's an email
	emailRe := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !(emailRe.MatchString(*emailPtr)) {
		color.Error.Printf("Invalid/No Email Provided")
		flag.Usage()
		os.Exit(1)
	}
	//define struct for getuserrealm json response
	type target struct {
		State                   int
		UserState               int
		Login                   string
		NameSpaceType           string
		DomainName              string
		FederationGlobalVersion int
		AuthURL                 string
		FederationBrandName     string
		CloudInstanceName       string
		CloudInstanceIssuerUri  string
	}

	splitArray := strings.Split(*emailPtr, "@")
	givenDomain := splitArray[1]
	returnedNames, err := net.LookupMX(givenDomain)
	if err != nil {
		color.Error.Printf("Could not get MX record: %v\n", err)
		os.Exit(1)
	}
	for _, mxRecord := range returnedNames {
		color.Info.Printf("MX Records for %v: %s\n", givenDomain, mxRecord.Host)
		if strings.Contains(mxRecord.Host, "pphosted.com") {
			color.Warn.Printf("MX Records indicate ProofPoint email filtering\n")
			break
		}
		if strings.Contains(mxRecord.Host, "mimecast.com") {
			color.Warn.Printf("MX Records indicate Mimecast email filtering\n")
			break
		}
	}

	//build URL for GetUserRealm
	url := "https://login.microsoftonline.com/getuserrealm.srf?login=" + *emailPtr

	// define http client options
	webClient := http.Client{
		Timeout: time.Second * 15, // Maximum of 2 secs
	}
	//make the request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, getErr := webClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	//	fmt.Printf("Full Response: %v\n", res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	//create struct for response
	target1 := target{}
	//read json
	jsonErr := json.Unmarshal(body, &target1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	//understanding the results
	color.Info.Printf("Results for %v:\n", givenDomain)
	if target1.NameSpaceType == "Unknown" {
		color.Warn.Printf("NameSpaceType was %v, meaning this domain does not exist in Office365 or Azure\n", target1.NameSpaceType)
	} else if target1.NameSpaceType == "Managed" {
		color.Info.Printf("NameSpaceType was %v, meaning this domain is not federated (all auth goes through Office365)\n", target1.NameSpaceType)
		color.Info.Printf("Recommendation: use Go365 - https://github.com/optiv/Go365\n")
	} else if target1.NameSpaceType == "Federated" {
		color.Info.Printf("NameSpaceType was %v, meaning this domain has been verified in Office365 but they use a third-party auth provider)\n", target1.NameSpaceType)
		color.Info.Printf("AuthURL for %v: %v\n", givenDomain, target1.AuthURL)
		if strings.Contains(target1.AuthURL, "idp/prp.wsf") {
			color.Warn.Printf("AuthURL looks like Ping Identity based on 'idp/prp.wsf' in URL (manually verify)\n")
		}
		if strings.Contains(target1.AuthURL, "okta.com") {
			color.Warn.Printf("AuthURL looks like Okta based on 'okta.com' in URL (manually verify)\n")
		}
		if strings.Contains(target1.AuthURL, "adfs/ls") {
			color.Warn.Printf("AuthURL looks like ADFS based on 'adfs/ls' in URL (manually verify)\n")
		}
		if strings.Contains(target1.AuthURL, "nidp/app") {
			color.Warn.Printf("AuthURL looks like NetIQ Access Manager based on 'nidp/app' in URL (manually verify)\n")
		}
	} else {
		color.Error.Printf("NameSpaceType looked weird, go investigate and submit an issue :)\n")
	}
	if target1.FederationBrandName != "" {
		color.Warn.Printf("FederationBrandName was set to %v\n", target1.FederationBrandName)
	}
}
