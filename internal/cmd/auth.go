package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/strscp/cli/internal/api"
	"golang.org/x/term"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with the Starscope API",
	Long:  "Log in by providing your Starscope API token. Tokens can be created at Settings > API Tokens in the Starscope dashboard.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		token := flagToken
		if token == "" {
			fmt.Println("Authenticate with the Starscope API.")
			fmt.Println()
			fmt.Println("Don't have an account yet? Sign up at https://starscope.app")
			fmt.Println("Generate a token at https://starscope.app/settings/api-tokens")
			fmt.Println()
			fmt.Print("Enter your Starscope API token: ")
			bytePw, err := term.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			if err != nil {
				// Fallback for non-terminal input (e.g., piping)
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					token = strings.TrimSpace(scanner.Text())
				}
			} else {
				token = strings.TrimSpace(string(bytePw))
			}
		}

		if token == "" {
			return fmt.Errorf("no token provided")
		}

		// Validate token
		apiURL := cfg.ResolveAPIURL(flagAPIURL)
		client := api.NewClient(apiURL, token)

		fmt.Print("Validating token... ")
		ws, err := client.Workspace.Show(context.Background())
		if err != nil {
			fmt.Println("failed")
			return fmt.Errorf("token validation failed: %w", err)
		}
		fmt.Println("ok")

		// Store token and workspace info
		profile := cfg.CurrentProfile()
		if err := cfg.StoreToken(profile, token); err != nil {
			return fmt.Errorf("storing token: %w", err)
		}

		cfg.SetProfileField(profile, "api_url", apiURL)
		cfg.SetProfileField(profile, "workspace_id", ws.ID)

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Printf("Authenticated as workspace '%s' (%s plan). Token stored securely.\n", ws.Name, ws.PlanTier)
		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored authentication",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		profile := cfg.CurrentProfile()
		if err := cfg.RemoveToken(profile); err != nil {
			return fmt.Errorf("removing token: %w", err)
		}

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Println("Logged out. Token removed.")
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		profile := cfg.CurrentProfile()
		token := cfg.ResolveToken(flagToken)
		apiURL := cfg.ResolveAPIURL(flagAPIURL)

		fmt.Printf("Profile:  %s\n", profile)
		fmt.Printf("API URL:  %s\n", apiURL)

		if token == "" {
			fmt.Println("Status:   Not authenticated")
			return nil
		}

		// Show masked token
		masked := token[:4] + strings.Repeat("*", len(token)-8) + token[len(token)-4:]
		fmt.Printf("Token:    %s\n", masked)

		// Validate with API
		client := api.NewClient(apiURL, token)
		ws, err := client.Workspace.Show(context.Background())
		if err != nil {
			fmt.Println("Status:   Invalid token")
			return nil
		}

		fmt.Printf("Workspace: %s (ID: %d)\n", ws.Name, ws.ID)
		fmt.Printf("Plan:     %s\n", ws.PlanTier)
		fmt.Println("Status:   Authenticated")
		return nil
	},
}

func init() {
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authStatusCmd)
}
