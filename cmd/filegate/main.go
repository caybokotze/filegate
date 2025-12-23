package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/anacrolix/dms/dlna/dms"
	"github.com/anacrolix/ffprobe"
	alog "github.com/anacrolix/log"
	"github.com/filegate/filegate/internal/tunnel"
	"github.com/filegate/filegate/internal/webdav"
)

const (
	passwordLength = 12
	passwordChars  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// WebDAVCmd handles the webdav subcommand
type WebDAVCmd struct {
	Local bool   `help:"Run in local mode (LAN only, no relay)" short:"l"`
	Port  int    `help:"Port to listen on (local mode only)" default:"8080" short:"p"`
	User  string `help:"Username for Basic Auth" default:"admin" short:"u"`
	Pass  string `help:"Password for Basic Auth (auto-generated if not provided)"`
	Relay string `help:"Relay server WebSocket URL" default:"wss://filegate.app/tunnel" hidden:""`
}

func (cmd *WebDAVCmd) Run() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Generate password if not provided
	password := cmd.Pass
	if password == "" {
		password, err = generatePassword(passwordLength)
		if err != nil {
			return fmt.Errorf("failed to generate password: %w", err)
		}
	}

	// Create WebDAV server
	srv, err := webdav.New(webdav.Config{
		Root:     cwd,
		Username: cmd.User,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("failed to create WebDAV server: %w", err)
	}

	if cmd.Local {
		runLocalMode(srv, cwd, cmd.Port, cmd.User, password)
	} else {
		runRemoteMode(srv, cwd, cmd.Relay, cmd.User, password)
	}
	return nil
}

// DLNACmd handles the dlna subcommand
type DLNACmd struct {
	Port int    `help:"Port to listen on" default:"8080" short:"p"`
	Name string `help:"Server name (defaults to hostname)" short:"n"`
}

func (cmd *DLNACmd) Run() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	runDLNAMode(cwd, cmd.Port, cmd.Name)
	return nil
}

var CLI struct {
	Webdav  WebDAVCmd `cmd:"" default:"withargs" help:"Expose directory via WebDAV (default: public URL via relay)"`
	Dlna    DLNACmd   `cmd:"" help:"Expose directory via DLNA for smart TVs"`
	Version kong.VersionFlag `help:"Show version" short:"v"`
}

var version = "dev"

func main() {
	ctx := kong.Parse(&CLI,
		kong.Name("filegate"),
		kong.Description("Expose the current directory via WebDAV or DLNA"),
		kong.UsageOnError(),
		kong.Vars{"version": version},
	)

	err := ctx.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func runDLNAMode(cwd string, port int, name string) {
	// Get hostname for friendly name
	hostname := name
	if hostname == "" {
		hostname, _ = os.Hostname()
		if hostname == "" {
			hostname = "filegate"
		}
	}

	// Create listener first to verify port is available
	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}

	// Create a logger for the DLNA server
	logger := alog.NewLogger("dms")
	logger.SetHandlers(alog.DiscardHandler)

	// Allow all clients (0.0.0.0/0)
	_, allowAll, _ := net.ParseCIDR("0.0.0.0/0")

	// Configure DLNA server
	server := &dms.Server{
		HTTPConn:       ln,
		FriendlyName:   hostname,
		RootObjectPath: cwd,
		NoTranscode:    true,  // Don't transcode - serve files directly
		NoProbe:        true,  // Disable probing to avoid dms library bugs
		NotifyInterval: 30 * time.Second,
		IgnoreHidden:   true,
		AllowedIpNets:  []*net.IPNet{allowAll},
		Logger:         logger,
		Icons:          defaultIcons(),
	}

	// Initialize the server
	if err := server.Init(); err != nil {
		ln.Close()
		log.Fatalf("Failed to initialize DLNA server: %v", err)
	}

	// Now print startup message after successful init
	ips := getLocalIPs()

	fmt.Println("Starting filegate in DLNA mode...")
	if ffprobe.Available() {
		fmt.Println("Media probing: \033[36mffprobe\033[0m")
	} else {
		fmt.Println("Media probing: \033[33minternal (ffprobe not found)\033[0m")
	}
	if isCommandAvailable("ffmpegthumbnailer") {
		fmt.Println("Thumbnails: \033[36menabled\033[0m")
	} else {
		fmt.Println("Thumbnails: \033[33mdisabled (ffmpegthumbnailer not found)\033[0m")
	}
	fmt.Println()
	fmt.Printf("Serving: %s\n", cwd)
	fmt.Println()
	fmt.Printf("Server name: %s\n", hostname)
	fmt.Println()
	fmt.Println("Your smart TV should discover this server automatically.")
	fmt.Println("Look for it in your TV's media/DLNA sources.")
	fmt.Println()
	fmt.Println("Access URLs:")
	for _, ip := range ips {
		fmt.Printf("  http://%s:%d\n", ip, port)
	}
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop")

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		fmt.Println("\nShutting down...")
		server.Close()
		ln.Close()
	}()

	if err := server.Run(); err != nil {
		log.Fatalf("DLNA server error: %v", err)
	}
}

func runLocalMode(srv *webdav.Server, cwd string, port int, username, password string) {
	addr := fmt.Sprintf(":%d", port)

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      srv,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	// Get local IPs
	ips := getLocalIPs()

	fmt.Println("Starting filegate in local mode...")
	fmt.Println()
	fmt.Printf("Serving: %s\n", cwd)
	fmt.Println()
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password: %s\n", password)
	fmt.Println()
	fmt.Println("Access URLs:")
	for _, ip := range ips {
		fmt.Printf("  http://%s:%d\n", ip, port)
	}
	fmt.Printf("  http://localhost:%d\n", port)
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop")

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		fmt.Println("\nShutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		httpServer.Shutdown(ctx)
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func runRemoteMode(srv *webdav.Server, cwd, relayURL, username, password string) {
	fmt.Println("Starting filegate...")
	fmt.Println()
	fmt.Printf("Serving: %s\n", cwd)
	fmt.Println()
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password: %s\n", password)
	fmt.Println()
	fmt.Println("Connecting to relay server...")

	ctx, cancel := context.WithCancel(context.Background())

	// Handle shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Println("\nShutting down...")
		cancel()
	}()

	client := tunnel.New(tunnel.Config{
		RelayURL: relayURL,
		Handler:  srv,
		OnConnected: func(subdomain, fullURL string) {
			fmt.Println()
			fmt.Printf("Connected! Your WebDAV is available at:\n")
			fmt.Printf("  %s\n", fullURL)
			fmt.Println()
			fmt.Println("Press Ctrl+C to stop")
		},
		OnDisconnected: func(err error) {
			fmt.Printf("\nDisconnected: %v\n", err)
		},
		OnReconnecting: func(attempt int) {
			fmt.Printf("Reconnecting (attempt %d)...\n", attempt)
		},
	})

	if err := client.Connect(ctx); err != nil && err != context.Canceled {
		log.Fatalf("Connection error: %v", err)
	}
}

// generatePassword creates a random password
func generatePassword(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(passwordChars))))
		if err != nil {
			return "", err
		}
		result[i] = passwordChars[num.Int64()]
	}
	return string(result), nil
}

// isCommandAvailable checks if a command exists in PATH
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// getLocalIPs returns all non-loopback IPv4 addresses
func getLocalIPs() []string {
	var ips []string

	interfaces, err := net.Interfaces()
	if err != nil {
		return ips
	}

	for _, iface := range interfaces {
		// Skip loopback and down interfaces
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Only include IPv4, non-loopback addresses
			if ip != nil && ip.To4() != nil && !ip.IsLoopback() {
				ips = append(ips, ip.String())
			}
		}
	}

	return ips
}

// defaultIcons returns a simple embedded icon for DLNA clients
func defaultIcons() []dms.Icon {
	// Simple 48x48 PNG icon (folder icon)
	// This is a minimal valid PNG to prevent serveIcon from panicking
	pngData := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x30,
		0x08, 0x02, 0x00, 0x00, 0x00, 0xd8, 0x60, 0x6e, 0xd5, 0x00, 0x00, 0x00,
		0x3c, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0xed, 0xc1, 0x01, 0x01, 0x00,
		0x00, 0x00, 0x82, 0x20, 0xff, 0xaf, 0x6e, 0x48, 0x40, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x70, 0x0b, 0x0b, 0xf0, 0x00, 0x01, 0xfe, 0xe9, 0x37,
		0x85, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60,
		0x82,
	}

	return []dms.Icon{
		{
			Width:    48,
			Height:   48,
			Depth:    24,
			Mimetype: "image/png",
			Bytes:    pngData,
		},
	}
}
