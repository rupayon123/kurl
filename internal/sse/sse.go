package sse

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kavix/kurl/color"
)

type Event struct {
	ID    string
	Event string
	Data  string
	Retry time.Duration
}

type Options struct {
	URL        string
	Headers    []string
	FilterType string
	OutputFile string
	NoColor    bool
}

// ParseStream reads SSE events line-by-line from a Reader according to W3C EventSource spec
func ParseStream(r io.Reader, handler func(Event)) error {
	scanner := bufio.NewScanner(r)
	// Allow large JSON payload lines while keeping scanner growth bounded.
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	scanner.Split(eventStreamLines())
	firstLine := true
	var current Event
	hasData := false

	for scanner.Scan() {
		line := scanner.Text()
		if firstLine {
			line = strings.TrimPrefix(line, "\ufeff")
			firstLine = false
		}
		if line == "" {
			// Dispatch event on empty line
			if hasData {
				handler(current)
			}
			current = Event{ID: current.ID, Retry: current.Retry}
			hasData = false
			continue
		}

		if strings.HasPrefix(line, ":") {
			// Comment line, ignore
			continue
		}

		field, value, found := strings.Cut(line, ":")
		if !found {
			field = line
			value = ""
		} else if strings.HasPrefix(value, " ") {
			value = value[1:]
		}

		switch field {
		case "event":
			current.Event = value
		case "data":
			if hasData {
				current.Data += "\n" + value
			} else {
				current.Data = value
			}
			hasData = true
		case "id":
			if !strings.ContainsRune(value, '\x00') {
				current.ID = value
			}
		case "retry":
			if value != "" && strings.IndexFunc(value, func(r rune) bool { return r < '0' || r > '9' }) == -1 {
				if ms, err := time.ParseDuration(value + "ms"); err == nil {
					current.Retry = ms
				}
			}
		}
	}

	return scanner.Err()
}

// eventStreamLines accepts CR, LF, and CRLF, including split CRLF pairs.
func eventStreamLines() bufio.SplitFunc {
	skipLF := false
	return func(data []byte, atEOF bool) (int, []byte, error) {
		if len(data) == 0 {
			return 0, nil, nil
		}
		if skipLF {
			skipLF = false
			if data[0] == '\n' {
				return 1, nil, nil
			}
		}
		if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
			skipLF = data[i] == '\r'
			return i + 1, data[:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	}
}

// RunSSE connects to an SSE endpoint and streams colorized events
func RunSSE(ctx context.Context, opts Options) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, opts.URL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	for _, h := range opts.Headers {
		name, value, ok := strings.Cut(h, ":")
		if ok {
			req.Header.Add(strings.TrimSpace(name), strings.TrimSpace(value))
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %s", resp.Status)
	}

	mediaType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil || mediaType != "text/event-stream" {
		return fmt.Errorf("expected Content-Type text/event-stream, got %q", resp.Header.Get("Content-Type"))
	}

	var logFile *os.File
	if opts.OutputFile != "" {
		f, err := os.Create(opts.OutputFile)
		if err != nil {
			return fmt.Errorf("failed to open output log file: %w", err)
		}
		defer f.Close()
		logFile = f
	}

	enabled := color.AutoEnabled(os.Stdout) && !opts.NoColor
	fmt.Printf("📡 Connected to SSE stream at %s\n\n", opts.URL)

	return ParseStream(resp.Body, func(ev Event) {
		eventType := ev.Event
		if eventType == "" {
			eventType = "message"
		}
		if opts.FilterType != "" && eventType != opts.FilterType {
			return
		}

		timestamp := time.Now().Format("15:04:05")
		header := fmt.Sprintf("[%s] EVENT: %s", timestamp, eventType)
		if ev.ID != "" {
			header += fmt.Sprintf(" (ID: %s)", ev.ID)
		}

		coloredHeader := color.Wrap(enabled, color.Bold+color.Cyan, header)
		fmt.Println(coloredHeader)

		if ev.Data != "" {
			fmt.Printf("  %s\n\n", ev.Data)
		}

		if logFile != nil {
			logLine := fmt.Sprintf("[%s] %s: %s\n", timestamp, eventType, ev.Data)
			logFile.WriteString(logLine)
		}
	})
}
