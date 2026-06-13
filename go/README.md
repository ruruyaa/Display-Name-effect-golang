# Go Integration Guide — Discord Display Name Styles
This directory contains the full concurrent **Golang** implementation of the Display Name Styles service. It integrates beautifully inside systems driven by the standard bwmarrin/discordgo package without blocking your main bot thread.
## Credits & Partners
 * **Main Server Partner**: KYRONIX
 * **High Partner**: Ruru Aka 2f9r
## 🏗️ Architecture Design
The Golang implementation mirrors the layered architecture of the TypeScript/JavaScript systems, but leverages Go's powerful concurrency model:
 1. **Transport Layer (DiscordProfileAPI)**: Safe, concurrent-ready network caller wrapped around Go's native net/http package. Supports exponential backoffs, rate limit evaluation, retry offsets, and raw logger callbacks.
 2. **Business Logic Layer (ProfileStyleService)**: Governs style resolution coordinates (JSON paths, custom fields), manages rotation state databases via local files, validates payloads, and handles discovery.
 3. **Startup Integration Layer (Ready Event Hook)**: Hooks into discordgo's standard event handlers on bot login, executing safely inside a lightweight goroutine (go func()) so it never blocks your bot's command listeners.
## 🛠️ Complete Installation and Setup
### 1. Requirements Installation
The module relies on the standard discordgo package for the bot session. Ensure it is installed in your Go module:
```bash
go get [github.com/bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)

```
### 2. Integration inside a typical discordgo Bot
Import your core structs and trigger the service as a background goroutine inside your client's Ready handler:
```go
// File: main.go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"[github.com/bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)"
	// Import your fuzzer package here if you placed it in a sub-directory
	// "yourproject/styles"
)

var styleTaskStarted bool = false

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		fmt.Println("Please set your DISCORD_TOKEN environment variable.")
		return
	}

	// Initialize the Discord bot session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	// Hook into the Ready event
	dg.AddHandler(onReady)

	dg.Identify.Intents = discordgo.IntentsDefault

	// Open the websocket connection
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	fmt.Println("[Bot] Running perfectly. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

// onReady triggers exactly once when the bot successfully connects to Discord's gateway
func onReady(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Printf("[Bot] Connected as %s#%s (ID: %s)\n", event.User.Username, event.User.Discriminator, event.User.ID)
	
	// Guard to ensure we only start the style discovery once per process
	if !styleTaskStarted {
		styleTaskStarted = true
		
		// Dispatch as a concurrent background goroutine immediately
		go runPresetsTask(s.Token)
	}
}

func runPresetsTask(token string) {
	// Panic-recovery block to ensure a fuzzer crash never kills the main bot
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[DisplayNameStyles] Safe background fallback caught: %v\n", r)
		}
	}()

	fmt.Println("[DisplayNameStyles] Dispatched background style service...")
	
	// Configure the style engine options
	opts := Options{
		Token:                 token,
		StyleMode:             "rotate", // Automatically rotate presets on each startup
		RunCompatibilityTests: false,
	}

	// Initialize and run the service we ported
	service := NewProfileStyleService(opts)
	report := service.Run()

	fmt.Printf("[DisplayNameStyles] Task complete. Endpoint Supported: %v\n", report["endpointSupported"])
}

```
## 📊 Environment Configurations
You can drive the Go service parameters by declaring direct environment variables inside your container deployment environment (e.g., in your Pterodactyl panel startup arguments or a .env file):
```env
# General switches
DISCORD_PROFILE_STYLE_ENABLED="true"

# Behavior controls
DISCORD_PROFILE_STYLE_MODE="rotate"            # 'rotate', 'random', or 'fixed'
DISCORD_PROFILE_STYLE_PRESET="ribes-neon-pink" # Override directly targeting this preset

# Network speeds
DISCORD_PROFILE_STYLE_REQUEST_DELAY_MS="1500"  # Delay spacing between compatibility checks

```

