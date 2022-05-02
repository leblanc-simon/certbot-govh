package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ovh/go-ovh/ovh"
	"golang.org/x/net/publicsuffix"
	"gopkg.in/ini.v1"
)

type DnsEntry struct {
	FieldType string `json:"fieldType"`
	SubDomain string `json:"subDomain"`
	Target    string `json:"target"`
	Ttl       int    `json:"ttl"`
}

type DnsReturn struct {
	Id int `json:"id"`
}

type OvhCredential struct {
	Endpoint          string
	ApplicationKey    string
	ApplicationSecret string
	ConsumerKey       string
}

func checkOvhCredential(ovhCredential *OvhCredential, domain string, primaryDomain string) error {
	homeDirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	cfg, err := ini.Load(homeDirname + "/.ovh/" + primaryDomain + ".ini")
	if err != nil {
		return err
	}

	ovhCredential.Endpoint = cfg.Section("default").Key("endpoint").String()
	ovhCredential.ApplicationKey = cfg.Section(ovhCredential.Endpoint).Key("application_key").String()
	ovhCredential.ApplicationSecret = cfg.Section(ovhCredential.Endpoint).Key("application_secret").String()
	ovhCredential.ConsumerKey = cfg.Section(ovhCredential.Endpoint).Key("consumer_key").String()

	if ovhCredential.Endpoint == "" {
		return errors.New("endpoint is not defined")
	}

	if ovhCredential.ApplicationKey == "" {
		return errors.New("ApplicationKey is not defined")
	}

	if ovhCredential.ApplicationSecret == "" {
		return errors.New("ApplicationSecret is not defined")
	}

	if ovhCredential.ConsumerKey == "" {
		return errors.New("ConsumerKey is not defined")
	}

	return nil
}

// Create the validation DNS entry
func addTokenToDns(client *ovh.Client, domain string, primaryDomain string) {
	validation := os.Getenv("CERTBOT_VALIDATION")

	dnsEntry := &DnsEntry{FieldType: "TXT", SubDomain: "_acme-challenge." + domain + ".", Target: validation, Ttl: 60}
	var dnsReturn DnsReturn

	// Create the entry
	err := client.Post("/domain/zone/"+primaryDomain+"/record", dnsEntry, &dnsReturn)
	if err != nil {
		log.Fatalf("Error: %q\n", err)
	}

	// Refresh DNS Zone
	err = client.Post("/domain/zone/"+primaryDomain+"/refresh", nil, nil)
	if err != nil {
		log.Fatalf("Error: %q\n", err)
	}

	// Print the zone ID : certbot set this value via CERTBOT_AUTH_OUTPUT environment variable
	fmt.Printf("%d", dnsReturn.Id)
}

// Remove the created entry from DNS Zone
func removeTokenToDns(client *ovh.Client, domain string, primaryDomain string) {
	zoneId := strings.TrimSpace(os.Getenv("CERTBOT_AUTH_OUTPUT"))
	if zoneId == "" {
		log.Fatal("No zone ID in CERTBOT_AUTH_OUTPUT")
	}

	// Delete the DNS entry
	err := client.Delete("/domain/zone/"+primaryDomain+"/record/"+zoneId, nil)
	if err != nil {
		log.Fatalf("Error: %q\n", err)
	}

	// Refresh DNS Zone
	err = client.Post("/domain/zone/"+primaryDomain+"/refresh", nil, nil)
	if err != nil {
		log.Fatalf("Error: %q\n", err)
	}
}

func help(optionalExitCode ...int) {
	fmt.Printf("%s\n", os.Args[0])

	exitCode := 0
	if len(optionalExitCode) > 0 {
		exitCode = optionalExitCode[0]
	}

	os.Exit(exitCode)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("This program require one parameter exactly")
	}

	action := os.Args[1]

	if action == "help" || action == "-h" || action == "--help" {
		help()
	}

	/*
	   CERTBOT_DOMAIN: The domain being authenticated
	   CERTBOT_VALIDATION: The validation string
	   CERTBOT_TOKEN: Resource name part of the HTTP-01 challenge (HTTP-01 only)
	   CERTBOT_REMAINING_CHALLENGES: Number of challenges remaining after the current challenge
	   CERTBOT_ALL_DOMAINS: A comma-separated list of all domains challenged for the current certificate
	*/

	// Get the domain and primary domain from certbot
	domain := os.Getenv("CERTBOT_DOMAIN")
	primaryDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		log.Fatalf("%s : no primary domain found", domain)
	}

	// Get OVH credential for primary domain
	var ovhCredential OvhCredential
	err = checkOvhCredential(&ovhCredential, domain, primaryDomain)
	if err != nil {
		log.Fatalf("Error: %q\n", err)
	}

	// Init the OVH HTTP Client
	client, err := ovh.NewClient(
		ovhCredential.Endpoint,
		ovhCredential.ApplicationKey,
		ovhCredential.ApplicationSecret,
		ovhCredential.ConsumerKey,
	)
	if err != nil {
		log.Fatalf("Error: %q\n", err)
	}

	// Process action
	switch action {
	case "create":
		addTokenToDns(client, domain, primaryDomain)
	case "delete":
		removeTokenToDns(client, domain, primaryDomain)
	default:
		fmt.Print("Action is not found\n")
		help(1)
	}
}
