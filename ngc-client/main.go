package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"ngc-client/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	clients *client.Clients
)

var mainCmd = &cobra.Command{
	Use:   "ngc",
	Short: "Példa CLI auth + user API hívásokra",
}

var callCmd = &cobra.Command{
	Use:   "call [auth|user] [endpoint]",
	Short: "Hív egy végpontot (pl. user /users)",
	Args:  cobra.ExactArgs(2),	
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]
		path := args[1]

		resp, err := clients.DoUserRequest("GET", path)
		if err != nil {
			log.Fatalf("Hiba: %v", err)
		}
		defer resp.Body.Close()

		switch service {
		case "auth":
			resp, err = clients.DoAuthRequest("GET", path)
		case "user":
			resp, err = clients.DoUserRequest("GET", path)
		default:
			log.Fatalf("Ismeretlen szolgáltatás: %s (auth vagy user várható)", service)
		}

		if err != nil {
			log.Fatalf("Hiba a kérésnél: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Status: %d\n", resp.StatusCode)
		fmt.Println("Response:")
		// Szép JSON kimenet
		var prettyJSON map[string]interface{}
		if json.Unmarshal(body, &prettyJSON) == nil {
			pretty, _ := json.MarshalIndent(prettyJSON, "", "  ")
			fmt.Println(string(pretty))
		} else {
			fmt.Println(string(body))
		}
	},
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Nem lehet beolvasni a config fájlt: %v", err)
	}

	var cfg client.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Config unmarshal hiba: %v", err)
	}

	var err error
	clients, err = client.NewClients(cfg)
	if err != nil {
		log.Fatalf("Client inicializálási hiba: %v", err)
	}

	fmt.Printf("Betöltve környezet: %s\n", cfg.Environment)
}

func main() {
	mainCmd.AddCommand(callCmd)

	mainCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.yaml)")

	mainCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		initConfig()
	}

	if err := mainCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
