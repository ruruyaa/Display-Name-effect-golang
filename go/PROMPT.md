# AI Prompts for Golang Display Name Styles System

This document contains specialized prompt templates designed for AI coding assistants. When utilizing these prompts, ensure that they are focused **strictly** on generating clean, idiomatic (`gofmt`-compliant), and strongly typed Golang code.

---

## 🎯 Master Generation Prompt

Use this prompt when instructing an AI to reconstruct or refactor the Display Name Styles codebase in Golang:

```text
Role: High-Performance Discord Golang Engineer
Partner Credits: okay & Ruru

Instruction:
Generate a completely self-contained, strongly typed, and professional Golang package for styling a Discord Bot account's profile name using custom experimental display name styles.

Strict Rules:
1. ONLY produce valid Golang code matching standard `gofmt` styling conventions. DO NOT output Python, JavaScript, or TypeScript syntax.
2. Utilize Go's native "net/http" package for concurrent-safe HTTP transmission.
3. Import core dependencies: os, encoding/json, time, path/filepath, bytes, net/http.
4. Ensure all struct signatures, methods, and configurations mirror the core model parameters provided using proper JSON struct tags.
5. All asynchronous background tasks must utilize native Goroutines (`go func()`) to ensure the main bot thread is never blocked.
6. Support dynamic logging by appending JSONL format lines to files, and write diagnostic report files cleanly in local directories using "os.OpenFile".
7. Implement rigorous exponential backoff retry algorithms when encountering Discord 429 Rate Limits (parsing the "retry-after" header) or 5xx temporary server failures.
8. Safely catch runtime panics using "defer recover()" to prevent standard event loop terminations or total bot crashes.

Style Specs:
- Support 12 discrete fonts (including IDs 1-12 such as Bangers, BioRhyme, Cherry Bomb, Neo-Castel, etc.)
- Support 6 custom effects (Solid, Gradient, Neon, Toon, Pop, Glow)
- Support integer slices ([]int) for display_name_colors.
- Support both nested Payload Format A ("display_name_styles": { ... }) and flat Payload Format B ("display_name_font_id": ...)

Code Structure:
The resulting package must export two core structs:
1. DiscordProfileAPI: Handles raw HTTP requests, retries, headers inspection, and is logged recursively.
2. ProfileStyleService: Resolves target presets, reads and writes working cached configurations from local filesystem JSON files, performs live capabilities/endpoint discovery, runs compatibility validation matrices, and summarizes accomplishments in a Markdown format.

