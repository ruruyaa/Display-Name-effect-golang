# Discord Display Name Styles: Golang Document Workspace

Welcome to the **Discord Display Name Styles** Golang document workspace. This workspace holds the exhaustive documentation, prompt lists, specifications, and code blueprints for implementing custom profile and displayName-styling systems inside modern Discord bots using Go.

## Credits & Verification
- **Created & maintained by**: [KyronixStudio](https://github.com/kyronixstudio)
- **Dev**: `dray.me`, `6fck`, `2f9r`
- **GlowForNodejs**: [GlowForNodejs](https://github.com/kyronixstudio/GlowForNodejs)

---

# Join Our discord
- [KyronixStudio](https://discord.gg/FBUEj8daSk)

---

This project documents the underlying experimental capabilities of modern Discord Profile APIs, including custom fonts, border/glow effects, selective coloring, name layouts, rate limiting, and startup discovery.

---

## Repository Directory Map & Navigation

The workspace is organized exclusively for our high-performance Golang architecture. Below you will find interactive navigation structures to easily browse through self-contained code implementations, language-specific README files, structured prompts, and secondary deep-dive details.

### Quick Navigation Dashboard

Use the dashboard below to jump to specific Golang modules:

| Component / File Type | Golang |
| :--- | :--- |
| **Main Implementation** | [gomain.go](./go/gomain.go) |
| **Language Guide** | [README.md](./go/README.md) |
| **AI Prompt Script** | [PROMPT.md](./go/PROMPT.md) |
| **Typography (Fonts)** | [FONT.md](./go/DeepDetails/FONT.md) |
| **Visual Effects** | [EFFECTS.md](./go/DeepDetails/EFFECTS.md) |
| **Colors Index** | [COLORS.md](./go/DeepDetails/COLORS.md) |
| **API Endpoints Spec** | [ENDPOINTS.md](./go/DeepDetails/ENDPOINTS.md) |
| **Experimental Flags** | [EXPERIMENT.md](./go/DeepDetails/EXPERIMENT.md) |
| **Compatibility Guide** | [COMPATIBILITY.md](./go/DeepDetails/COMPATIBILITY.md) |

---

### Clickable File Tree Map

Explore the file hierarchy interactively:

[project-root](./)<br>
├── [README.md](./README.md) *(Root entry point - this file)*<br>
└── [go/](./go/) *(Golang Project Directory)*<br>
&nbsp;&nbsp;&nbsp;&nbsp;├── [gomain.go](./go/gomain.go) *(Standard Golang concurrent service)*<br>
&nbsp;&nbsp;&nbsp;&nbsp;├── [README.md](./go/README.md) *(Golang integration guide & discordgo setup)*<br>
&nbsp;&nbsp;&nbsp;&nbsp;├── [PROMPT.md](./go/PROMPT.md) *(AI code prompt for custom Golang generation)*<br>
&nbsp;&nbsp;&nbsp;&nbsp;└── [DeepDetails/](./go/DeepDetails/) *(Golang deep-dive details)*<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├── [FONT.md](./go/DeepDetails/FONT.md) *(Typography database: IDs 1-12)*<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├── [EFFECTS.md](./go/DeepDetails/EFFECTS.md) *(Visual effects: Solid, Glow, Neon)*<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├── [COLORS.md](./go/DeepDetails/COLORS.md) *(RGB decimal/hex color index)*<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├── [ENDPOINTS.md](./go/DeepDetails/ENDPOINTS.md) *(Golang REST API endpoints)*<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├── [EXPERIMENT.md](./go/DeepDetails/EXPERIMENT.md) *(Profile experiment gates)*<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;└── [COMPATIBILITY.md](./go/DeepDetails/COMPATIBILITY.md) *(discordgo, Arikawa, Disgo support)*


## Key System Design Architecture

In our Golang engine, the display name styles system operates on a highly concurrent three-layer model:

1. **Transport Layer (`DiscordProfileAPI`)**: Focuses on rate-limiting, exponential backoff, raw response capture, parsed JSON structures, diagnostic tracking, and authenticated requests using Go's native `net/http`.
2. **Business Logic Layer (`ProfileStyleService`)**: Resolves styles, loads/saves working configurations, runs startup capability discovery, and manages preset rotations.
3. **Startup Integration Layer (`Ready` Event Hook)**: Hooks into the bot's standard websocket startup flow as a safe, concurrent background Goroutine. It is completely non-blocking, so API failures or rate limits do not block standard bot features (music commands, moderation, dashboards).

---

## Core Documentation Sections

The `go/DeepDetails` folder contains sub-files detailing every aspect of the system:
- **FONT**: Enumerates all 12 fonts (Bangers, BioRhyme, Cherry Bomb, Chicle, Compagnon, Museo Moderno, Neo-Castel, Pixelify Sans, Ribes, Sinistre, GG Sans, and Zilla Slab).
- **EFFECTS**: Describes visual styles (Solid, Gradient, Neon, Toon, Pop, Glow).
- 
