package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/onetwopunch/google-groups/fetcher"
	"log"
)

func main() {
	keyFile := flag.String("key-file", "", "Service Account Key JSON file path")
	impersonate := flag.String("impersonate", "", "GSuite admin email to impersonate")
	subject := flag.String("subject", "", "GSuite user for which to fetch groups")
	recurseDepth := flag.Int("depth", 0, "(optional) Depth of recursion if desired. i.e user belongs to group, which belongs to group, etc")
	flag.Parse()

	if len(*impersonate) == 0 {
		log.Fatalf("A gsuite admin email must be provided for service account impersonation")
	}

	if len(*subject) == 0 {
		log.Fatalf("A gsuite email must be provided for the subject")
	}

	f, err := fetcher.NewDefaultGroupFetcher(*keyFile, *impersonate)
	if err != nil {
		log.Fatal(err)
	}

	groups, err := f.ListGoogleGroups(*subject, *recurseDepth)
	if err != nil {
		log.Fatal(err)
	}

	data, err := json.Marshal(groups)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
