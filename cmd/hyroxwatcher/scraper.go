package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// CategoryStatus holds the availability of one ticket category.
type CategoryStatus struct {
	Name      string
	Available bool
}

// CheckAvailability navigates the Hyrox event page, clicks through the
// "Buy Tickets" -> "Athlete Ticket" flow, and returns per-category status.
func CheckAvailability(ctx context.Context, eventURL string) ([]CategoryStatus, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		chromedp.WindowSize(1280, 900),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()

	timeoutCtx, cancelTimeout := context.WithTimeout(browserCtx, 90*time.Second)
	defer cancelTimeout()

	// Navigate and wait for the page to load.
	if _, err := chromedp.Run(timeoutCtx, chromedp.Steps(
		chromedp.Navigate(eventURL),
		chromedp.Sleep(3*time.Second),
	)); err != nil {
		return nil, fmt.Errorf("navigating to event page: %w", err)
	}

	// Try to click "Buy Tickets" — attempt multiple selectors.
	clickBuy(timeoutCtx)

	if _, err := chromedp.Run(timeoutCtx, chromedp.Sleep(2*time.Second)); err != nil {
		return nil, fmt.Errorf("waiting after buy click: %w", err)
	}

	// Try to click "Athlete Ticket" in the modal — best effort.
	clickAthlete(timeoutCtx)

	if _, err := chromedp.Run(timeoutCtx, chromedp.Sleep(3*time.Second)); err != nil {
		return nil, fmt.Errorf("waiting after athlete click: %w", err)
	}

	// Read the full body text to parse availability.
	bodyText, err := chromedp.Run(timeoutCtx, chromedp.Text("body"))
	if err != nil {
		return nil, fmt.Errorf("reading page body: %w", err)
	}

	return parseAvailability(bodyText), nil
}

func clickBuy(ctx context.Context) {
	selectors := []string{
		`.buy-tickets`,
		`[data-action="buy-tickets"]`,
		`a[href*="ticket"]`,
		`button[class*="ticket"]`,
	}
	tryClick(ctx, selectors)
}

func clickAthlete(ctx context.Context) {
	selectors := []string{
		`[class*="athlete"]`,
		`[data-type="athlete"]`,
		`.ticket-type-athlete`,
	}
	tryClick(ctx, selectors)
}

// tryClick attempts each CSS selector and clicks the first one found.
func tryClick(ctx context.Context, selectors []string) {
	for _, sel := range selectors {
		count, err := chromedp.Run(ctx, chromedp.Evaluate[int](
			fmt.Sprintf(`document.querySelectorAll(%q).length`, sel),
		))
		if err != nil || count == 0 {
			continue
		}
		if _, err := chromedp.Run(ctx, chromedp.Click(sel)); err == nil {
			return
		}
	}
}

// parseAvailability looks for sold-out or available signals in the page text.
// The category keywords will need tuning after a discovery run —
// see the README section "Selector Discovery".
func parseAvailability(body string) []CategoryStatus {
	lines := strings.Split(body, "\n")
	var statuses []CategoryStatus
	soldOutKeywords := []string{"sold out", "soldout", "unavailable", "not available"}

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !looksLikeCategory(line) {
			continue
		}
		// Check this line and the next few for sold-out language.
		window := strings.Join(lines[i:minInt(i+5, len(lines))], " ")
		lower := strings.ToLower(window)
		available := true
		for _, kw := range soldOutKeywords {
			if strings.Contains(lower, kw) {
				available = false
				break
			}
		}
		statuses = append(statuses, CategoryStatus{Name: line, Available: available})
	}

	// If nothing parsed, return a sentinel so callers know the page structure
	// needs selector updates (see README).
	if len(statuses) == 0 {
		statuses = append(statuses, CategoryStatus{Name: "unknown", Available: false})
	}
	return statuses
}

func looksLikeCategory(s string) bool {
	if len(s) < 3 || len(s) > 80 {
		return false
	}
	lower := strings.ToLower(s)
	hints := []string{"open", "pro", "doubles", "relay", "wave", "athlete", "elite"}
	for _, h := range hints {
		if strings.Contains(lower, h) {
			return true
		}
	}
	return false
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
