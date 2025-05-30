package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Global flags
var (
	StoragePath string
	Port        int
)

func Parse() (err error) {
	var endpoints bool
	flag.IntVar(&Port, "port", 4000, "port to serve on")
	flag.StringVar(&StoragePath, "dir", "data", "S3 path to save images from posts and comments")
	flag.BoolVar(&endpoints, "endpoints", false, "to show endpoints")

	flag.Usage = func() {
		PrintHelp()
	}

	flag.Parse()

	if Port < 1024 || Port > 65535 {
		return fmt.Errorf("incorrect range port, port must me between 1024 and 65535")
	}

	if endpoints {
		PrintEndPoints()
		os.Exit(0)
	}

	return nil
}

func PrintHelp() {
	fmt.Println(`1337bo4rd

Usage:
  1337bo4rd [--port <N>] [--dir <S>] 
  1337bo4rd --help

Options:
  --help       Show this screen.
  --port N     Port number.
  --dir S      Directory of S3 (images from post and comments).
  --endpoints  Show the api endpoints.
  `)
}

func PrintEndPoints() {
	title := `
╔═╗┬ ┬┌─┐┌─┐┌─┐┬─┐┌┬┐  ╔╦╗┌─┐┌┬┐┌─┐┬─┐
╠═╝├─┤├─┤│ ┬│ │├┬┘ │───║║║├┤  │ ├┤ ├┬┘
╩  ┴ ┴┴ ┴└─┘└─┘┴└─ ┴   ╩ ╩└─┘ ┴ └─┘┴└─
`

	endpoints := []struct {
		Method string
		Path   string
		Desc   string
		Auth   bool
	}{
		{"GET", "/", "Main page (post catalog)", false},
		{"GET", "/catalog", "View catalog (tiled layout)", false},
		{"GET", "/archive", "View all archived posts", false},
		{"POST", "/post", "Create new post (+image upload)", true},
		{"GET", "/post/{id}", "View post with comments", false},
		{"POST", "/post/{id}/comment", "Add comment to post", true},
		{"GET", "/archive/post/{id}", "View archived post", false},
		{"GET", "/create", "Post creation form (HTML)", true},
		{"GET", "/auth/init", "Initialize session (sets cookie)", false},
		{"POST", "/auth/name", "Change display name", true},
	}

	fmt.Println(title)
	fmt.Println("╔══════════╦════════════════════╦═════════════════════════════════════╗")
	fmt.Println("║ Method   ║ Path               ║ Description                         ║")
	fmt.Println("╠══════════╬════════════════════╬═════════════════════════════════════╣")

	for _, ep := range endpoints {
		auth := "  "
		if ep.Auth {
			auth = "🔒"
		}
		fmt.Printf("║ %-8s ║ %-18s ║ %-33s %s║\n",
			ep.Method,
			ep.Path,
			ep.Desc,
			auth)
	}

	fmt.Println("╚══════════╩════════════════════╩═════════════════════════════════════╝")
	fmt.Println("\nKey:")
	fmt.Println("  🔒 - Requires valid session cookie")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("Session Management:")
	fmt.Println("- Automatic cookie generation on first visit")
	fmt.Println("- Session expires after 1 week")
	fmt.Println("- Avatars sourced from Rick and Morty API")
	fmt.Println(strings.Repeat("═", 60))
}
