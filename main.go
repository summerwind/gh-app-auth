package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v45/github"
	flag "github.com/spf13/pflag"
)

func fail(msg string) {
	fmt.Fprintf(os.Stderr, msg)
}

func main() {
	var (
		appID         int
		privateKey    string
		targetAccount string
		help          bool
	)

	flag.IntVarP(&appID, "app-id", "i", 0, "GitHub App ID")
	flag.StringVarP(&privateKey, "private-key", "f", "", "Path to the private key file of GitHub App")
	flag.StringVarP(&targetAccount, "account", "a", "", "Target account name")
	flag.BoolVarP(&help, "help", "h", false, "Print usage and exit")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if appID == 0 {
		fail("GitHub App ID must be specified")
		os.Exit(1)
	}
	if privateKey == "" {
		fail("Private key file must be specified")
		os.Exit(1)
	}
	if targetAccount == "" {
		fail("Target account must be specified")
		os.Exit(1)
	}

	ctx := context.Background()

	app, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, int64(appID), privateKey)
	if err != nil {
		fail(fmt.Sprintf("failed to authenticate as GitHub App: %v", err))
		os.Exit(1)
	}

	appClient := github.NewClient(&http.Client{Transport: app})

	installationID := int64(0)
	opts := &github.ListOptions{PerPage: 10}

	for {
		installations, resp, err := appClient.Apps.ListInstallations(ctx, opts)
		if err != nil {
			fail(fmt.Sprintf("failed to get installations: %v", err))
			os.Exit(1)
		}

		for _, inst := range installations {
			account := inst.GetAccount()
			if account.GetLogin() == targetAccount {
				installationID = inst.GetID()
				break
			}
		}

		if installationID != 0 || resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	if installationID == 0 {
		fail("the GitHub App is not installed on the target account")
		os.Exit(1)
	}

	instToken, _, err := appClient.Apps.CreateInstallationToken(ctx, installationID, &github.InstallationTokenOptions{})
	if err != nil {
		fail(fmt.Sprintf("failed to create installation token: %v", err))
		os.Exit(1)
	}
	fmt.Println(instToken.GetToken())
}
