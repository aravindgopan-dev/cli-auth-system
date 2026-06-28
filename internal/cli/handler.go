package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aravindgopan-dev/cli-auth-system/internal/repository"

	prompt "github.com/c-bata/go-prompt"
)

const (
	colorReset = "\033[0m"
	colorCyan  = "\033[1;36m"
	colorRed   = "\033[1;31m"
	colorGreen = "\033[1;32m"
)

type Authenticator interface {
	Register(ctx context.Context, username, password string) error
	PreLoginValidate(ctx context.Context, username string) (*repository.User, error)
	PasswordLogin(ctx context.Context, user *repository.User, password string) error
	VerifyTOTP(user *repository.User, code string) bool
	CreateSession(ctx context.Context, username string) (string, time.Time, error)
	Generate2FASecret(username string) (string, string, error)
	Enable2FA(ctx context.Context, user *repository.User, secret string) error
	Disable2FA(ctx context.Context, user *repository.User) error
}

type SessionStore interface {
	GetSession(ctx context.Context, token string) (*repository.Session, error)
	DeleteSession(ctx context.Context, token string) error
}

type CLIHandler struct {
	AuthService  Authenticator
	UserRepo     SessionStore
	CurrentToken string
	CurrentUser  *repository.User
}

func NewCLIHandler(authService Authenticator, userRepo SessionStore) *CLIHandler {
	return &CLIHandler{AuthService: authService, UserRepo: userRepo}
}

func (h *CLIHandler) Run() {
	banner := `
███████  ███████  ██████  ██    ██ ██████  ███████      ██████  ██      ██
██       ██      ██       ██    ██ ██   ██ ██          ██       ██      ██
███████  █████   ██       ██    ██ ██████  █████       ██       ██      ██
     ██  ██      ██       ██    ██ ██   ██ ██          ██       ██      ██
███████  ███████  ██████   ██████  ██   ██ ███████      ██████  ███████ ██
`
	fmt.Printf("%s%s%s\n", colorCyan, banner, colorReset)
	h.printHelpMenu()

	p := prompt.New(
		h.ExecuteCommand,
		h.CompleteCommand,
		prompt.OptionLivePrefix(h.ChangeLivePromptPrefix),
		prompt.OptionTitle("Secure-CLI Console"),
		prompt.OptionPrefixTextColor(prompt.Cyan),
		prompt.OptionInputTextColor(prompt.White),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSuggestionTextColor(prompt.White),
		prompt.OptionSelectedSuggestionBGColor(prompt.Cyan),
		prompt.OptionSelectedSuggestionTextColor(prompt.Black),
		prompt.OptionDescriptionBGColor(prompt.LightGray),
		prompt.OptionDescriptionTextColor(prompt.Black),
		prompt.OptionSelectedDescriptionBGColor(prompt.Turquoise),
		prompt.OptionSelectedDescriptionTextColor(prompt.Black),
	)
	p.Run()
}

func (h *CLIHandler) ChangeLivePromptPrefix() (string, bool) {
	if h.CurrentUser != nil {
		return fmt.Sprintf("%s> ", h.CurrentUser.Username), true
	}
	return "guest> ", true
}

func (h *CLIHandler) CompleteCommand(d prompt.Document) []prompt.Suggest {
	var suggestions []prompt.Suggest

	if h.CurrentUser == nil {
		suggestions = []prompt.Suggest{
			{Text: "register", Description: "Build a brand new user profile"},
			{Text: "login", Description: "Authenticate against security tokens"},
			{Text: "help", Description: "Show available guest routes"},
			{Text: "exit", Description: "Close system console link completely"},
		}
	} else {
		suggestions = []prompt.Suggest{
			{Text: "whoami", Description: "View session metrics and account rules"},
			{Text: "logout", Description: "Flush workspace authorization cache"},
			{Text: "help", Description: "Show active account workspace routes"},
		}
		if h.CurrentUser.TwoFAEnabled {
			suggestions = append(suggestions, prompt.Suggest{Text: "disable-2fa", Description: "Strip account of 2FA boundaries"})
		} else {
			suggestions = append(suggestions, prompt.Suggest{Text: "enable-2fa", Description: "Turn on secondary Google Authenticator TOTP factor"})
		}
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func (h *CLIHandler) ExecuteCommand(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	parts := strings.Fields(input)
	command := parts[0]
	args := parts[1:]

	ctx := context.Background()
	isAuthenticated := h.CurrentUser != nil

	if isAuthenticated {
		sess, err := h.UserRepo.GetSession(ctx, h.CurrentToken)
		if err != nil || time.Now().After(sess.ExpiresAt) {
			fmt.Println("\n[Session Expired]. Access tokens flushed.")
			h.performLogout(ctx)
			return
		}
	}

	guestMap := map[string]func(context.Context, []string) error{
		"register": h.Register,
		"login":    h.Login,
		"help":     h.Help,
		"exit":     h.Exit,
	}

	authMap := map[string]func(context.Context, []string) error{
		"whoami":      h.WhoAmI,
		"enable-2fa":  h.Enable2FA,
		"disable-2fa": h.Disable2FA,
		"logout":      h.Logout,
		"help":        h.Help,
	}

	var targetFunc func(context.Context, []string) error
	var exists bool

	if !isAuthenticated {
		targetFunc, exists = guestMap[command]
	} else {
		targetFunc, exists = authMap[command]
	}

	if !exists {
		fmt.Println("Unknown command. Hit [TAB] to parse operational suggestions.")
		return
	}

	if err := targetFunc(ctx, args); err != nil {
		fmt.Printf("%sError:%s %v\n", colorRed, colorReset, err)
	}
}