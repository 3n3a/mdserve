package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/3n3a/mdserve/internal/assets"
	"github.com/3n3a/mdserve/internal/server"
)

var version = "dev"

func defaultPort() int {
	if v := os.Getenv("PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return 8000
}

func envBool(key string) bool {
	v := os.Getenv(key)
	return v == "1" || v == "true" || v == "yes"
}

func main() {
	host := flag.String("host", "127.0.0.1", "host to bind to")
	port := flag.Int("port", defaultPort(), "port to listen on (or $PORT)")
	user := flag.String("user", os.Getenv("MDSERVE_USER"), "basic-auth username (or $MDSERVE_USER); enables auth when set")
	pass := flag.String("pass", os.Getenv("MDSERVE_PASS"), "basic-auth password (or $MDSERVE_PASS)")
	offline := flag.Bool("offline", envBool("MDSERVE_OFFLINE"), "serve Pico CSS and Mermaid from the binary instead of the CDN")
	saveCheckboxes := flag.Bool("save-checkboxes", envBool("MDSERVE_SAVE_CHECKBOXES"), "persist task-list checkbox state to the browser's localStorage")
	showVersion := flag.Bool("version", false, "print version and exit")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: mdserve [flags] [content_dir]\n\nServe .md files as HTML.\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	contentDir := "."
	if flag.NArg() > 0 {
		contentDir = flag.Arg(0)
	}
	abs, err := filepath.Abs(contentDir)
	if err != nil {
		log.Fatalf("resolving content dir: %v", err)
	}

	if (*user == "") != (*pass == "") {
		log.Fatal("--user and --pass must be set together")
	}

	if *offline && assets.Empty() {
		log.Println("WARNING: --offline set but bundled assets are empty; pages will render unstyled. Did you run scripts/fetch-assets.sh?")
	}

	handler, err := server.New(server.Options{
		ContentDir:     abs,
		User:           *user,
		Pass:           *pass,
		Offline:        *offline,
		SaveCheckboxes: *saveCheckboxes,
	})
	if err != nil {
		log.Fatalf("building server: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", *host, *port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("mdserve %s — serving .md files from: %s", version, abs)
	if *user != "" {
		log.Print("basic auth: enabled")
	}
	log.Printf("Open http://%s", addr)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
