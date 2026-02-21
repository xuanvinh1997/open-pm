package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	dbURL := flag.String("database", "", "Database URL")
	path := flag.String("path", "./migrations", "Migrations path")
	flag.Parse()

	if *dbURL == "" {
		*dbURL = os.Getenv("OPENPM_DATABASE_URL")
	}
	if *dbURL == "" {
		fmt.Fprintln(os.Stderr, "error: database URL required (use -database flag or OPENPM_DATABASE_URL env var)")
		os.Exit(1)
	}

	m, err := migrate.New("file://"+*path, *dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: migrate [up|down|force|version]")
		os.Exit(1)
	}

	switch args[0] {
	case "up":
		if len(args) > 1 {
			n, _ := strconv.Atoi(args[1])
			err = m.Steps(n)
		} else {
			err = m.Up()
		}
	case "down":
		n := 1
		if len(args) > 1 {
			n, _ = strconv.Atoi(args[1])
		}
		err = m.Steps(-n)
	case "force":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "usage: migrate force <version>")
			os.Exit(1)
		}
		v, _ := strconv.Atoi(args[1])
		err = m.Force(v)
	case "version":
		v, dirty, verr := m.Version()
		if verr != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", verr)
			os.Exit(1)
		}
		fmt.Printf("version: %d, dirty: %v\n", v, dirty)
		return
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		os.Exit(1)
	}

	if err != nil && err != migrate.ErrNoChange {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if err == migrate.ErrNoChange {
		fmt.Println("no change")
	} else {
		fmt.Println("done")
	}
}
