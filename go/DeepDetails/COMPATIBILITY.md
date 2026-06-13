# 📦 Library Compatibility Details — Golang

This document outlines the framework compatibility and library-specific integrations for the **Discord Display Name Styles** service in Golang.

---

## 🟢 discordgo Integration

The subsystem is fully compatible with the standard **`bwmarrin/discordgo`** package. It operates cleanly as a concurrent Goroutine dispatched from standard startup event handlers (such as the `Ready` event).

### 🔌 Safe Task Dispatching

Always coordinate style updating inside a separate Goroutine to avoid blocking the main websocket heartbeat or standard command handling events:

```go
package main

import (
	"fmt"
	"[github.com/bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)"
	// Import your custom styles package
	// "yourproject/styles" 
)

var styleServiceDispatched bool = false

func main() {
	dg, _ := discordgo.New("Bot YOUR_TOKEN_HERE")
	
	// Register the Ready handler
	dg.AddHandler(onReady)
	
	dg.Open()
	
	// Keep connection alive
	select {}
}

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Printf("Logged in as %s\n", s.State.User.Username)
	
	if !styleServiceDispatched {
		styleServiceDispatched = true
		
		// Schedule as a background goroutine so it doesn't block the event loop
		go runCustomStyleSetup(s.Token)
	}
}

func runCustomStyleSetup(token string) {
	opts := Options{
		Token: token,
		StyleMode: "rotate",
	}
	
	service := NewProfileStyleService(opts)
	report := service.Run()
	
	fmt.Printf("[Styles] Active format: %v\n", report["payloadFormat"])
}
```

---

## 🔑 Permissions & Server Scopes

To modify server nickname attributes, the bot requires specific scopes:

### 1. Change Nickname Permissions
When executing `PATCH /guilds/{guild_id}/members/@me`, the bot requires the `Change Nickname` permission in that specific guild.
- **Golang Verification Check**:
  ```go
  // Assuming you have the channelID and the session (s)
  perms, err := s.State.UserChannelPermissions(s.State.User.ID, channelID)
  
  // Use a bitwise AND to check for the specific permission flag
  if err == nil && (perms & discordgo.PermissionChangeNickname == discordgo.PermissionChangeNickname) {
      fmt.Println("Bot is permitted to style its nickname!")
  }
  ```

### 2. Arikawa and Disgo (Alternative Libraries)
- **Arikawa / Disgo**: The Display Name Styles service is completely decoupled from the WebSocket client. It is natively compatible with alternative Go Discord libraries like `arikawa` or `disgo` without any code changes to the core service. 
- You only need to extract the raw Bot Token string from your chosen library's configuration and pass it directly into the `Options{ Token: "..." }` struct when initializing the `ProfileStyleService`. If token extraction is difficult, simply configure `DISCORD_TOKEN` as an OS environment variable and the service will fall back to it automatically.

---

## 🛠️ Golang Version Requirements

- Minimum Go Version: **Go 1.18+** (Recommended for optimal native JSON parsing and module management).
- Core HTTP Dependency: The native `net/http` package (used inside `DiscordProfileAPI` for executing stable, retrying REST PATCH actions). **No external third-party HTTP modules are required.**
- 
