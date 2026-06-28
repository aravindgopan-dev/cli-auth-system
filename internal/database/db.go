package database

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

func InitDB(connString string) *pgx.Conn {
	ctx := context.Background()
	var conn *pgx.Conn
	var err error

	for i := 0; i < 10; i++ {
		conn, err = pgx.Connect(ctx, connString)
		if err == nil {
			break
		}
		log.Printf("Database connection failed (attempt %d/10): %v. Retrying in 2 seconds...", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Unable to connect to database after 10 attempts: %v", err)
	}

	runMigrations(ctx, conn)

	return conn
}

func runMigrations(ctx context.Context, conn *pgx.Conn) {
	paths := []string{"./migration", "../migration", "/app/migration", "../../migration"}
	var dir string
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			dir = p
			break
		}
	}
	if dir == "" {
		log.Println("Warning: migration directory not found, skipping database migrations")
		return
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Warning: failed to read migration directory: %v", err)
		return
	}

	var upFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			upFiles = append(upFiles, entry.Name())
		}
	}
	sort.Strings(upFiles)

	for _, file := range upFiles {
		content, err := os.ReadFile(filepath.Join(dir, file))
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", file, err)
		}
		_, err = conn.Exec(ctx, string(content))
		if err != nil {
			log.Fatalf("Failed to execute migration %s: %v", file, err)
		}
	}
}