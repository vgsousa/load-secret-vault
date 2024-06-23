package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/vault/api"
)

type VaultConfig struct {
	Address string
	Token   string
}

func getVaultSecret(client *api.Client, secretPath, key string) (string, error) {

	secret, err := client.KVv2("secret").Get(context.Background(), secretPath)
	if err != nil {
		return "", err
	}
	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("no data found at path: %s", secretPath)
	}

	value, ok := secret.Data[key].(string)

	if !ok {
		return "", fmt.Errorf("key %s not found in secret data", key)
	}

	return value, nil
}

// loads input file
func loadEnvFile(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// saves environment variables to a output file
func saveEnvFile(filePath, content string) error {
	return ioutil.WriteFile(filePath, []byte(content), 0644)
}

func main() {

	if len(os.Args) < 4 {
		fmt.Println("Use: go run main.go $secretPath $inputFile $outputFile")
		return
	}

	secretPath := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]

	vaultConfig := VaultConfig{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
	}

	client, err := api.NewClient(&api.Config{
		Address: vaultConfig.Address,
	})
	if err != nil {
		log.Fatalf("Error initializing Vault client: %v", err)
	}
	client.SetToken(vaultConfig.Token)

	// Load input file
	envContent, err := loadEnvFile(inputFile)

	if err != nil {
		log.Fatalf("Error loading %s file: %v", inputFile, err)
	}

	// Use regex to find placeholders in the form ${PLACEHOLDER}
	re := regexp.MustCompile(`\${([^}]+)}`)
	matches := re.FindAllStringSubmatch(envContent, -1)

	replacements := make(map[string]string)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		placeholder := match[0]
		secretItem := match[1]

		// Fetch the secret value from the Vault
		secretValue, err := getVaultSecret(client, secretPath, secretItem)
		if err != nil {
			log.Printf("Error fetching secret for %s: %v", secretItem, err)
			continue
		}

		replacements[placeholder] = secretValue
	}

	// Replace placeholders with actual secret values
	for placeholder, value := range replacements {
		envContent = strings.ReplaceAll(envContent, placeholder, value)
	}

	if err := saveEnvFile(outputFile, envContent); err != nil {
		log.Fatalf("Error saving %s file: %v", outputFile, err)
	}

	log.Println("Processed successfully.")
}
