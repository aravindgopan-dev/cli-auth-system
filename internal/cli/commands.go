package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aravindgopan-dev/cli-auth-system/internal/repository"
	prompt "github.com/c-bata/go-prompt"
	"github.com/mdp/qrterminal/v3"
)

func (h *CLIHandler) Register(ctx context.Context, args []string) error {
	username, password := h.promptCredentials()
	if err := h.AuthService.Register(ctx, username, password); err != nil {
		return err
	}
	fmt.Println("Registration successful!")
	return nil
}

func (h *CLIHandler) Login(ctx context.Context, args []string) error {
	username, password := h.promptCredentials()
	user, err := h.AuthService.PreLoginValidate(ctx, username)
	if err != nil {
		return err
	}

	if err := h.AuthService.PasswordLogin(ctx, user, password); err != nil {
		return err
	}

	if user.TwoFAEnabled {
		code := prompt.Input("Enter 2FA TOTP Code: ", func(d prompt.Document) []prompt.Suggest { return nil })
		code = strings.TrimSpace(code)
		if !h.AuthService.VerifyTOTP(user, code) {
			return fmt.Errorf("invalid MFA validation code string sequence")
		}
	}

	token, _, err := h.AuthService.CreateSession(ctx, user.Username)
	if err != nil {
		return fmt.Errorf("failed to process active connection parameters")
	}

	h.CurrentToken = token
	h.CurrentUser = user
	fmt.Println("\nLogin Successful!")
	h.displayWhoAmI(ctx)
	return nil
}

func (h *CLIHandler) Exit(ctx context.Context, args []string) error {
	fmt.Println("Goodbye!")
	os.Exit(0)
	return nil
}

func (h *CLIHandler) WhoAmI(ctx context.Context, args []string) error {
	h.displayWhoAmI(ctx)
	return nil
}

func (h *CLIHandler) Enable2FA(ctx context.Context, args []string) error {
	if h.CurrentUser.TwoFAEnabled {
		fmt.Println("2FA is already enabled.")
		return nil
	}
	secret, url, _ := h.AuthService.Generate2FASecret(h.CurrentUser.Username)
	fmt.Printf("Secret seed string: %s\nURI target parameters: %s\n", secret, url)
	
	fmt.Println("\nScan this QR code with your authenticator app:")
	config := qrterminal.Config{
		Level:      qrterminal.L,
		Writer:     os.Stdout,
		HalfBlocks: true,
	}
	qrterminal.GenerateWithConfig(url, config)
	fmt.Println()
	
	code := prompt.Input("Verify app TOTP token number sequence: ", func(d prompt.Document) []prompt.Suggest { return nil })
	code = strings.TrimSpace(code)

	if h.AuthService.VerifyTOTP(&repository.User{TwoFASecret: secret}, code) {
		_ = h.AuthService.Enable2FA(ctx, h.CurrentUser, secret)
		h.CurrentUser.TwoFAEnabled = true
		fmt.Println("2FA enabled successfully!")
	} else {
		fmt.Println("Invalid verification code. Canceled configuration steps.")
	}
	return nil
}

func (h *CLIHandler) Disable2FA(ctx context.Context, args []string) error {
	_ = h.AuthService.Disable2FA(ctx, h.CurrentUser)
	h.CurrentUser.TwoFAEnabled = false
	fmt.Println("2FA configuration deactivated.")
	return nil
}

func (h *CLIHandler) Logout(ctx context.Context, args []string) error {
	h.performLogout(ctx)
	fmt.Println("Logged out successfully.")
	return nil
}

func (h *CLIHandler) Help(ctx context.Context, args []string) error {
	h.printHelpMenu()
	return nil
}

// Internal Local Shared Utilities
func (h *CLIHandler) promptCredentials() (string, string) {
	u := prompt.Input("Username: ", func(d prompt.Document) []prompt.Suggest { return nil })
	p := prompt.Input("Password: ", func(d prompt.Document) []prompt.Suggest { return nil })
	return strings.TrimSpace(u), strings.TrimSpace(p)
}

func (h *CLIHandler) performLogout(ctx context.Context) {
	_ = h.UserRepo.DeleteSession(ctx, h.CurrentToken)
	h.CurrentUser = nil
	h.CurrentToken = ""
}

func (h *CLIHandler) printHelpMenu() {
	if h.CurrentUser == nil {
		fmt.Println("\nAvailable commands: register, login, help, exit")
	} else {
		fmt.Println("\nAvailable commands: whoami, enable-2fa, disable-2fa, logout, help")
	}
}

func (h *CLIHandler) displayWhoAmI(ctx context.Context) {
	// CLI Handler can read session mappings straight because Repo satisfies SessionStore
	sess, _ := h.UserRepo.GetSession(ctx, h.CurrentToken)
	fmt.Println("\n--- User Context Workspace ---")
	fmt.Printf("Username: %s\n", h.CurrentUser.Username)
	fmt.Printf("Created On: %s\n", h.CurrentUser.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("2FA Enabled: %t\n", h.CurrentUser.TwoFAEnabled)
	fmt.Printf("Session Expiration Timestamp: %s\n", sess.ExpiresAt.Format("15:04:05"))
	
	if h.CurrentUser.LastLogin.Valid {
		fmt.Printf("Last Success System Entry: %s\n", h.CurrentUser.LastLogin.Time.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Println("Last Success System Entry: Never")
	}
	fmt.Println("------------------------------")
}