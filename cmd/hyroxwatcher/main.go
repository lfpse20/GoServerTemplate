package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// Config holds all runtime configuration from environment variables.
type Config struct {
	EventURL         string
	NotifyEmails     []string
	GmailUser        string
	GmailAppPassword string
	StatusFile       string
}

func loadConfig() Config {
	emails := strings.Split(os.Getenv("NOTIFY_EMAILS"), ",")
	var filtered []string
	for _, e := range emails {
		if e = strings.TrimSpace(e); e != "" {
			filtered = append(filtered, e)
		}
	}
	statusFile := os.Getenv("STATUS_FILE")
	if statusFile == "" {
		statusFile = "hyrox_status.json"
	}
	eventURL := os.Getenv("HYROX_EVENT_URL")
	if eventURL == "" {
		eventURL = "https://hyrox.com/event/hyrox-tampa/"
	}
	return Config{
		EventURL:         eventURL,
		NotifyEmails:     filtered,
		GmailUser:        os.Getenv("GMAIL_USER"),
		GmailAppPassword: os.Getenv("GMAIL_APP_PASSWORD"),
		StatusFile:       statusFile,
	}
}

// savedState is what we persist between runs to detect transitions.
type savedState struct {
	// map of category name -> available bool
	Categories map[string]bool `json:"categories"`
}

func loadState(path string) savedState {
	data, err := os.ReadFile(path)
	if err != nil {
		return savedState{Categories: map[string]bool{}}
	}
	var s savedState
	if err := json.Unmarshal(data, &s); err != nil {
		return savedState{Categories: map[string]bool{}}
	}
	if s.Categories == nil {
		s.Categories = map[string]bool{}
	}
	return s
}

func saveState(path string, s savedState) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func main() {
	cfg := loadConfig()

	log.Printf("checking availability for %s", cfg.EventURL)

	ctx := context.Background()
	statuses, err := CheckAvailability(ctx, cfg.EventURL)
	if err != nil {
		log.Fatalf("scrape failed: %v", err)
	}

	log.Printf("found %d categories", len(statuses))
	for _, s := range statuses {
		log.Printf("  %s: available=%v", s.Name, s.Available)
	}

	prev := loadState(cfg.StatusFile)
	current := savedState{Categories: map[string]bool{}}

	var newlyAvailable []string
	for _, s := range statuses {
		current.Categories[s.Name] = s.Available
		wasAvailable, seen := prev.Categories[s.Name]
		if s.Available && (!seen || !wasAvailable) {
			newlyAvailable = append(newlyAvailable, s.Name)
		}
	}

	if err := saveState(cfg.StatusFile, current); err != nil {
		log.Printf("warning: could not save state: %v", err)
	}

	if len(newlyAvailable) == 0 {
		log.Println("no new availability — no notification sent")
		return
	}

	subject := fmt.Sprintf("Hyrox Tampa tickets available: %s", strings.Join(newlyAvailable, ", "))
	body := fmt.Sprintf(
		"Tickets just became available for the following categories:\n\n%s\n\nBuy now: %s",
		strings.Join(newlyAvailable, "\n"),
		cfg.EventURL,
	)
	log.Printf("sending alert: %s", subject)

	if err := sendAlert(cfg, subject, body); err != nil {
		log.Fatalf("failed to send alert: %v", err)
	}
	log.Println("alert sent successfully")
}
