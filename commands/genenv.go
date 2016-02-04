package commands

import (
	"bufio"
	"fmt"
	"github.com/danjac/podbaby/commands/Godeps/_workspace/src/golang.org/x/crypto/ssh/terminal"
	"github.com/danjac/podbaby/config"
	"os"
	"strings"
	"text/template"
)

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func readLine(prompt, defaultValue string, required bool) string {
	for {
		fmt.Printf("\r\n%s ", prompt)
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "" {
			text = defaultValue
		}
		if text != "" || !required {
			return text
		}
	}
}

func readSecret(prompt, defaultValue string, required bool) string {
	for {
		fmt.Printf("\r\n%s ", prompt)
		bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
		text := string(bytePassword)
		text = strings.TrimSpace(text)
		if text == "" {
			text = defaultValue
		}
		if text != "" || !required {
			return text
		}
		return text
	}
}

// Genenv interactively generates an environment settings file
// Usage: podbaby genenv [dest-file.env]

func Genenv(dst string) {
	if err := genenv(dst); err != nil {
		panic(err)
	}
}

func genenv(dst string) error {

	// does destination file exist?

	if fileExists(dst) {
		text := readLine("This file already exists. Overwrite (Y/N)?", "N", false)
		fmt.Println("TEXT:", text, len(text))
		if strings.ToLower(text) != "y" {
			return nil
		}
	}

	pgHost := readLine("PostgreSQL Host (localhost:5432)?", "localhost:5432", true)

	pgUser := readLine("PostgreSQL User?", "", false)
	pgName := readLine("PostgreSQL Database (podbaby)", "podbaby", true)
	pgPass := readSecret("PostgreSQL Password?", "", false)

	smtpHost := readLine("SMTP Host? (localhost)", "localhost", true)
	smtpUser := readLine("SMTP User?", "", false)
	smtpPass := readSecret("SMPT Password?", "", false)

	secretKey := readSecret("Secret key?", config.RandomKey(), true)
	secureCookieKey := readSecret("Secure cookie key?", config.RandomKey(), true)

	gaKey := readLine("Google analytics key?", "", false)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		pgUser, pgPass, pgHost, pgName)

	cfg := config.Default()
	cfg.DatabaseURL = dbURL
	cfg.SecretKey = secretKey
	cfg.SecureCookieKey = secureCookieKey
	cfg.GoogleAnalyticsID = gaKey
	cfg.Mail.Host = smtpHost
	cfg.Mail.User = smtpUser
	cfg.Mail.Password = smtpPass

	t, err := template.New("env").Parse(`
DB_URL="{{.DatabaseURL}}"
GOOGLE_ANALYTICS_ID="{{.GoogleAnalyticsID}}"
MAX_DB_CONNECTIONS="{{.MaxDBConnections}}"
SECRET_KEY="{{.SecretKey}}"
SECURE_COOKIE_KEY="{{.SecureCookieKey}}"
`)
	if err != nil {
		return err
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, cfg)
}
