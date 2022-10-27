package goutils

import (
	"fmt"
	"strings"
)

var clientSecrets map[string]string

// Load client secrets from environment variables
func EnableAPISecretKeys() {
	apiClients := Env("API_CLIENTS", []string{})
	if len(apiClients) != 0 {
		fmt.Printf("\r\n┌─────── CLIENT_SECRET: ─────────\r\n")

		clientSecrets = make(map[string]string)
		for _, client := range apiClients {
			clientSecrets[client] = Env(fmt.Sprintf("API_%s_SECRET", strings.TrimSpace(strings.ToUpper(client))), "")
			fmt.Printf("│ %s: %s\r\n", client, clientSecrets[client])
		}

		fmt.Println("└──────────────────────────────────────")
	}
}
