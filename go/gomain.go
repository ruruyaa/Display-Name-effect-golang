package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// --- Constants & Pre-defined Data ---

type Style struct {
	FontID   int   `json:"font_id"`
	EffectID int   `json:"effect_id"`
	Colors   []int `json:"colors"`
}

type Preset struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Style Style  `json:"style"`
}

var FONTS = []map[string]interface{}{
	{"label": "Bangers", "id": 1},
	{"label": "BioRhyme", "id": 2},
	{"label": "Cherry Bomb", "id": 3},
	{"label": "Chicle", "id": 4},
	{"label": "Compagnon", "id": 5},
	{"label": "Museo Moderno", "id": 6},
	{"label": "Neo-Castel", "id": 7},
	{"label": "Pixelify Sans", "id": 8},
	{"label": "Ribes", "id": 9},
	{"label": "Sinistre", "id": 10},
	{"label": "Default (GG Sans)", "id": 11},
	{"label": "Zilla Slab", "id": 12},
}

var EFFECTS = []map[string]interface{}{
	{"label": "Solid", "id": 1},
	{"label": "Gradient", "id": 2},
	{"label": "Neon", "id": 3},
	{"label": "Toon", "id": 4},
	{"label": "Pop", "id": 5},
	{"label": "Glow", "id": 6},
}

var COLOR_TESTS = []map[string]interface{}{
	{"label": "White", "colors": []int{16777215}},
	{"label": "Blue", "colors": []int{5865}},
	{"label": "Pink", "colors": []int{16711935}},
	{"label": "Purple", "colors": []int{8388736}},
	{"label": "White to Blue Gradient", "colors": []int{16777215, 5865}},
	{"label": "Pink to Purple Gradient", "colors": []int{16711935, 8388736}},
}

var TARGET_STYLE = Style{FontID: 10, EffectID: 3, Colors: []int{16777215}}

var STYLE_PRESETS = []Preset{
	{"sinistre-neon-white", "Sinistre Neon White", Style{FontID: 10, EffectID: 3, Colors: []int{16777215}}},
	{"ribes-neon-pink", "Ribes Neon Pink", Style{FontID: 9, EffectID: 3, Colors: []int{16711935}}},
	{"neo-castel-gradient-blue-white", "Neo-Castel Blue/White Gradient", Style{FontID: 7, EffectID: 2, Colors: []int{5865, 16777215}}},
	{"pixelify-pop-purple", "Pixelify Sans Pop Purple", Style{FontID: 8, EffectID: 5, Colors: []int{8388736}}},
	{"bangers-glow-pink-purple", "Bangers Pink/Purple Glow", Style{FontID: 1, EffectID: 6, Colors: []int{16711935, 8388736}}},
	{"cherry-toon-white", "Cherry Bomb Toon White", Style{FontID: 3, EffectID: 4, Colors: []int{16777215}}},
	{"zilla-solid-blue", "Zilla Slab Solid Blue", Style{FontID: 12, EffectID: 1, Colors: []int{5865}}},
}

var SUPPORTED_ERROR_STATUSES = map[int]bool{
	400: true, 401: true, 403: true, 404: true, 405: true, 409: true, 429: true, 500: true, 502: true, 503: true, 504: true,
}

// --- Discord Profile API Client ---

type APIResult struct {
	Ok                bool                   `json:"ok"`
	Status            int                    `json:"status"`
	StatusText        string                 `json:"statusText"`
	Method            string                 `json:"method"`
	Endpoint          string                 `json:"endpoint"`
	URL               string                 `json:"url"`
	Attempt           int                    `json:"attempt"`
	DurationMs        int64                  `json:"durationMs"`
	Headers           map[string]string      `json:"headers"`
	RetryAfterSeconds *float64               `json:"retryAfterSeconds"`
	RateLimited       bool                   `json:"rateLimited"`
	Raw               string                 `json:"raw"`
	Parsed            map[string]interface{} `json:"parsed"`
	ParseError        *string                `json:"parseError"`
	RequestBody       interface{}            `json:"requestBody"`
	Error             map[string]interface{} `json:"error"`
	VerifiedStyle     bool                   `json:"verifiedStyle,omitempty"`
	Verification      []interface{}          `json:"verification,omitempty"`
	GuildID           *string                `json:"guildId,omitempty"`
}

type DiscordProfileAPI struct {
	Token             string
	BaseURL           string
	Logger            func(map[string]interface{})
	MaxRetries        int
	RetryBaseDelayMs  int
	Client            *http.Client
}

func NewDiscordProfileAPI(token string, logger func(map[string]interface{}), maxRetries int) *DiscordProfileAPI {
	if logger == nil {
		logger = func(m map[string]interface{}) {}
	}
	return &DiscordProfileAPI{
		Token:            token,
		BaseURL:          "https://discord.com/api/v10",
		Logger:           logger,
		MaxRetries:       maxRetries,
		RetryBaseDelayMs: 1000,
		Client:           &http.Client{Timeout: 10 * time.Second},
	}
}

func (api *DiscordProfileAPI) Patch(endpoint string, body interface{}) APIResult {
	return api.Request("PATCH", endpoint, body, api.MaxRetries)
}

func (api *DiscordProfileAPI) Get(endpoint string, maxRetries int) APIResult {
	return api.Request("GET", endpoint, nil, maxRetries)
}

func (api *DiscordProfileAPI) Request(method, endpoint string, body interface{}, maxRetries int) APIResult {
	url := api.BaseURL + endpoint
	attempt := 0
	maxAttempts := maxRetries + 1

	for attempt < maxAttempts {
		attempt++
		startTime := time.Now()
		var parseError *string
		var parsed map[string]interface{}
		var rawBody []byte
		var err error

		if body != nil {
			rawBody, _ = json.Marshal(body)
		}

		req, err := http.NewRequest(method, url, bytes.NewBuffer(rawBody))
		if err != nil {
			return api.buildErrorResult(err, method, endpoint, url, attempt, startTime, body)
		}

		req.Header.Set("Authorization", "Bot "+api.Token)
		req.Header.Set("User-Agent", "DiscordBot (https://discord.com, 1.0)")
		if method != "GET" {
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")
		}

		resp, err := api.Client.Do(req)
		if err != nil {
			result := api.buildErrorResult(err, method, endpoint, url, attempt, startTime, body)
			api.Logger(structToMap(result))
			if attempt < maxAttempts {
				time.Sleep(api.calculateBackoff(attempt))
				continue
			}
			return result
		}

		rawBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		raw := string(rawBytes)

		if len(rawBytes) > 0 {
			if jsonErr := json.Unmarshal(rawBytes, &parsed); jsonErr != nil {
				errMsg := jsonErr.Error()
				parseError = &errMsg
			}
		}

		duration := time.Since(startTime).Milliseconds()
		headers := api.pickHeaders(resp)
		retryAfterSeconds := api.getRetryAfterSeconds(resp, parsed)
		rateLimited := resp.StatusCode == 429

		result := APIResult{
			Ok:                resp.StatusCode >= 200 && resp.StatusCode < 300,
			Status:            resp.StatusCode,
			StatusText:        resp.Status,
			Method:            method,
			Endpoint:          endpoint,
			URL:               url,
			Attempt:           attempt,
			DurationMs:        duration,
			Headers:           headers,
			RetryAfterSeconds: retryAfterSeconds,
			RateLimited:       rateLimited,
			Raw:               raw,
			Parsed:            parsed,
			ParseError:        parseError,
			RequestBody:       body,
		}

		api.Logger(structToMap(result))

		if rateLimited && attempt < maxAttempts {
			delay := api.calculateBackoff(attempt)
			if retryAfterSeconds != nil {
				delay = time.Duration(*retryAfterSeconds * float64(time.Second))
			}
			time.Sleep(delay)
			continue
		}

		if (resp.StatusCode == 500 || resp.StatusCode == 502 || resp.StatusCode == 503 || resp.StatusCode == 504) && attempt < maxAttempts {
			time.Sleep(api.calculateBackoff(attempt))
			continue
		}

		return result
	}

	return APIResult{} // Should not reach here
}

func (api *DiscordProfileAPI) calculateBackoff(attempt int) time.Duration {
	delayMs := float64(api.RetryBaseDelayMs) * math.Pow(2, float64(attempt-1))
	return time.Duration(delayMs) * time.Millisecond
}

func (api *DiscordProfileAPI) buildErrorResult(err error, method, endpoint, url string, attempt int, startTime time.Time, body interface{}) APIResult {
	return APIResult{
		Ok:         false,
		Status:     0,
		StatusText: "NETWORK_ERROR",
		Method:     method,
		Endpoint:   endpoint,
		URL:        url,
		Attempt:    attempt,
		DurationMs: time.Since(startTime).Milliseconds(),
		Headers:    map[string]string{},
		Raw:        err.Error(),
		RequestBody: body,
		Error: map[string]interface{}{
			"name":    fmt.Sprintf("%T", err),
			"message": err.Error(),
			"stack":   nil,
		},
	}
}

func (api *DiscordProfileAPI) getRetryAfterSeconds(resp *http.Response, parsed map[string]interface{}) *float64 {
	if parsed != nil {
		if val, ok := parsed["retry_after"].(float64); ok {
			return &val
		}
	}
	if headerVal := resp.Header.Get("retry-after"); headerVal != "" {
		if val, err := strconv.ParseFloat(headerVal, 64); err == nil {
			return &val
		}
	}
	return nil
}

func (api *DiscordProfileAPI) pickHeaders(resp *http.Response) map[string]string {
	headers := make(map[string]string)
	keys := []string{"content-type", "date", "x-ratelimit-bucket", "x-ratelimit-limit", "x-ratelimit-remaining", "x-ratelimit-reset", "x-ratelimit-reset-after", "retry-after", "x-discord-trace-id"}
	for _, key := range keys {
		if val := resp.Header.Get(key); val != "" {
			headers[key] = val
		}
	}
	return headers
}

// --- Service Core ---

type Options struct {
	LogDir                string
	ForceDiscovery        bool
	RunCompatibilityTests bool
	RequestDelayMs        int
	MaxRetries            int
	TargetStyle           *Style
	StylePreset           string
	StyleMode             string
	GuildID               string
	Token                 string
}

type ProfileStyleService struct {
	Options             Options
	Report              map[string]interface{}
	RunID               string
	LogFile             string
	ReportFile          string
	CacheFile           string
	RotationFile        string
	API                 *DiscordProfileAPI
	SelectedStylePreset *Preset
}

func NewProfileStyleService(opts Options) *ProfileStyleService {
	if opts.LogDir == "" {
		cwd, _ := os.Getwd()
		opts.LogDir = filepath.Join(cwd, "logs", "display-name-styles")
	}
	if opts.RequestDelayMs == 0 {
		opts.RequestDelayMs = 1500
	}
	if opts.MaxRetries == 0 {
		opts.MaxRetries = 2
	}
	if opts.StyleMode == "" {
		opts.StyleMode = "rotate"
	}

	runID := strings.ReplaceAll(strings.ReplaceAll(time.Now().UTC().Format(time.RFC3339Nano), ":", "-"), ".", "-")

	return &ProfileStyleService{
		Options:      opts,
		Report:       createEmptyReport(),
		RunID:        runID,
		LogFile:      filepath.Join(opts.LogDir, fmt.Sprintf("%s.jsonl", runID)),
		ReportFile:   filepath.Join(opts.LogDir, "latest-report.md"),
		CacheFile:    filepath.Join(opts.LogDir, "working-config.json"),
		RotationFile: filepath.Join(opts.LogDir, "preset-rotation.json"),
	}
}

func createEmptyReport() map[string]interface{} {
	return map[string]interface{}{
		"title":                     "Display Name Styles Report",
		"generatedAt":               time.Now().UTC().Format(time.RFC3339),
		"endpointSupported":         "UNKNOWN",
		"botTokenSupported":         "UNKNOWN",
		"payloadFormat":             "UNKNOWN",
		"acceptedFontIds":           []int{},
		"acceptedEffectIds":         []int{},
		"acceptedColors":            []string{},
		"finalWorkingConfiguration": nil,
		"finalTargetConfiguration":  nil,
		"selectedStylePreset":       nil,
		"availableStylePresets":     STYLE_PRESETS,
		"unsupportedFields":         []string{},
		"endpoints":                 make(map[string]interface{}),
		"notes":                     []string{},
	}
}

func (s *ProfileStyleService) Run() map[string]interface{} {
	os.MkdirAll(s.Options.LogDir, os.ModePerm)

	if os.Getenv("DISCORD_PROFILE_STYLE_ENABLED") == "false" {
		s.writeSummary("Display Name Styles disabled by DISCORD_PROFILE_STYLE_ENABLED=false.")
		return s.Report
	}

	token := s.Options.Token
	if token == "" {
		token = os.Getenv("DISCORD_TOKEN")
	}
	if token == "" {
		s.writeSummary("No Discord bot token available for Display Name Styles startup task.")
		return s.Report
	}

	s.API = NewDiscordProfileAPI(token, func(entry map[string]interface{}) {
		s.logEvent("api-response", s.sanitizeLogEntry(entry))
	}, s.Options.MaxRetries)

	s.resolveTargetStyle()
	s.Report["botTokenSupported"] = "UNKNOWN"
	s.Report["selectedStylePreset"] = s.SelectedStylePreset
	s.Report["finalTargetConfiguration"] = s.Options.TargetStyle

	cached := s.loadWorkingConfig()
	if cached != nil && !s.Options.ForceDiscovery {
		s.writeSummary("Using saved Display Name Styles configuration from a previous successful run.")
		responses := s.applyStyleToConfiguredScope(cached, *s.Options.TargetStyle, "cached-apply", "Saved Working Configuration")
		
		applied := false
		for _, r := range responses {
			if s.isStyleConfirmed(r, *s.Options.TargetStyle) {
				applied = true
				break
			}
		}

		if applied {
			s.saveWorkingConfig(cached) // Simplified fallback
			s.finalizeReport("Applied saved working configuration.")
			return s.Report
		}
		s.writeSummary("Saved Display Name Styles configuration failed verification; running fresh discovery.")
	}

	endpoints := s.getCandidateEndpoints()
	s.writeSummary(fmt.Sprintf("Starting Display Name Styles discovery with %d endpoint(s).", len(endpoints)))

	working := s.detectWorkingConfiguration(endpoints)
	if working != nil {
		s.applyFinalStyle(working)
		s.saveWorkingConfig(working)
		if s.Options.RunCompatibilityTests {
			s.runCompatibilityMatrix(working)
			s.applyFinalStyle(working)
		}
	} else {
		s.Report["endpointSupported"] = "NO"
		if s.Report["botTokenSupported"] != "YES" {
			s.Report["botTokenSupported"] = "NO"
		}
	}

	msg := "No working Display Name Styles configuration was found."
	if working != nil {
		msg = "Display Name Styles startup task completed."
	}
	s.finalizeReport(msg)

	return s.Report
}

func (s *ProfileStyleService) resolveTargetStyle() {
	if s.Options.TargetStyle != nil {
		s.SelectedStylePreset = &Preset{Key: "custom-options", Label: "Custom Options", Style: *s.Options.TargetStyle}
		return
	}

	preset := s.selectPreset()
	// Deep copy to prevent accidental mutation of constants
	colorsCopy := make([]int, len(preset.Style.Colors))
	copy(colorsCopy, preset.Style.Colors)
	s.Options.TargetStyle = &Style{FontID: preset.Style.FontID, EffectID: preset.Style.EffectID, Colors: colorsCopy}
	s.SelectedStylePreset = preset
	s.writeSummary(fmt.Sprintf("Selected Display Name Style preset: %s (%s).", preset.Label, preset.Key))
}

func (s *ProfileStyleService) selectPreset() *Preset {
	presetKey := strings.ToLower(strings.TrimSpace(s.Options.StylePreset))
	if presetKey != "" {
		for _, p := range STYLE_PRESETS {
			if p.Key == presetKey || strings.ToLower(p.Label) == presetKey {
				return &p
			}
		}
		s.Report["notes"] = append(s.Report["notes"].([]string), fmt.Sprintf("Unknown Display Name Style preset '%s', falling back to rotation mode.", s.Options.StylePreset))
	}

	mode := strings.ToLower(s.Options.StyleMode)
	if mode == "random" {
		return &STYLE_PRESETS[rand.Intn(len(STYLE_PRESETS))]
	}
	if mode == "fixed" {
		return &STYLE_PRESETS[0]
	}

	// Rotation state
	state := s.loadRotationState()
	index := 1
	if val, ok := state["nextIndex"].(float64); ok {
		index = int(val)
	}

	preset := STYLE_PRESETS[index%len(STYLE_PRESETS)]
	s.saveRotationState(map[string]interface{}{
		"nextIndex":     (index + 1) % len(STYLE_PRESETS),
		"lastPresetKey": preset.Key,
		"updatedAt":     time.Now().UTC().Format(time.RFC3339),
	})
	return &preset
}

func (s *ProfileStyleService) getCandidateEndpoints() []map[string]interface{} {
	guildID := s.Options.GuildID
	endpoints := []map[string]interface{}{}

	if guildID != "" {
		endpoints = append(endpoints, map[string]interface{}{
			"key":              "guild-members-me",
			"label":            fmt.Sprintf("PATCH /guilds/%s/members/@me", guildID),
			"endpoint":         fmt.Sprintf("/guilds/%s/members/@me", guildID),
			"endpointTemplate": "/guilds/{guild_id}/members/@me",
			"guildId":          guildID,
			"guildScoped":      true,
			"payloadFormats":   []string{"B", "A"},
		})
		endpoints = append(endpoints, map[string]interface{}{
			"key":              "guild-profile-me",
			"label":            fmt.Sprintf("PATCH /guilds/%s/profile/@me", guildID),
			"endpoint":         fmt.Sprintf("/guilds/%s/profile/@me", guildID),
			"endpointTemplate": "/guilds/{guild_id}/profile/@me",
			"guildId":          guildID,
			"guildScoped":      true,
			"payloadFormats":   []string{"B", "A"},
		})
	} else {
		s.Report["notes"] = append(s.Report["notes"].([]string), "No guild id was available, so guild-specific profile endpoints were skipped.")
	}

	endpoints = append(endpoints, map[string]interface{}{
		"key":              "users-me",
		"label":            "PATCH /users/@me",
		"endpoint":         "/users/@me",
		"endpointTemplate": "/users/@me",
		"guildId":          nil,
		"guildScoped":      false,
		"payloadFormats":   []string{"B", "A"},
	})

	return endpoints
}

func (s *ProfileStyleService) detectWorkingConfiguration(endpoints []map[string]interface{}) map[string]interface{} {
	for _, endpoint := range endpoints {
		formats := endpoint["payloadFormats"].([]string)
		for _, format := range formats {
			payload := s.buildPayload(format, *s.Options.TargetStyle)

			response := s.testPayload(map[string]interface{}{
				"phase":         "endpoint-discovery",
				"endpoint":      endpoint,
				"payloadFormat": format,
				"payload":       payload,
				"style":         *s.Options.TargetStyle,
			})

			s.recordEndpointResult(endpoint, format, response)
			if s.isStyleConfirmed(response, *s.Options.TargetStyle) {
				working := map[string]interface{}{
					"endpoint":         endpoint["endpoint"],
					"endpointTemplate": endpoint["endpointTemplate"],
					"endpointLabel":    endpoint["label"],
					"endpointKey":      endpoint["key"],
					"guildId":          endpoint["guildId"],
					"guildScoped":      endpoint["guildScoped"],
					"payloadFormat":    format,
					"style":            *s.Options.TargetStyle,
					"discoveredAt":     time.Now().UTC().Format(time.RFC3339),
				}

				s.Report["endpointSupported"] = "YES"
				s.Report["botTokenSupported"] = "YES"
				s.Report["payloadFormat"] = format
				s.Report["finalWorkingConfiguration"] = working
				return working
			}

			if response.Ok {
				s.Report["notes"] = append(s.Report["notes"].([]string), fmt.Sprintf("%s returned %d, but the response did not confirm Display Name Styles were applied.", endpoint["label"], response.Status))
			}
			s.captureUnsupportedFields(response)
			time.Sleep(time.Duration(s.Options.RequestDelayMs) * time.Millisecond)
		}
	}
	return nil
}

func (s *ProfileStyleService) runCompatibilityMatrix(working map[string]interface{}) {
	s.writeSummary("Running Display Name Styles compatibility matrix for known fonts, effects, and colors.")

	for _, font := range FONTS {
		style := *s.Options.TargetStyle
		style.FontID = font["id"].(int)
		response := s.testStyleVariant(working, "font", font["label"].(string), style)
		if response.Ok {
			s.addUniqueInt(&s.Report, "acceptedFontIds", font["id"].(int))
		}
		time.Sleep(time.Duration(s.Options.RequestDelayMs) * time.Millisecond)
	}

	for _, effect := range EFFECTS {
		style := *s.Options.TargetStyle
		style.EffectID = effect["id"].(int)
		response := s.testStyleVariant(working, "effect", effect["label"].(string), style)
		if response.Ok {
			s.addUniqueInt(&s.Report, "acceptedEffectIds", effect["id"].(int))
		}
		time.Sleep(time.Duration(s.Options.RequestDelayMs) * time.Millisecond)
	}

	for _, colorTest := range COLOR_TESTS {
		style := *s.Options.TargetStyle
		colors := colorTest["colors"].([]int)
		if len(colors) > 1 {
			style.EffectID = 2
		}
		style.Colors = colors
		response := s.testStyleVariant(working, "color", colorTest["label"].(string), style)
		if response.Ok {
			colorStrs := []string{}
			for _, c := range colors {
				colorStrs = append(colorStrs, strconv.Itoa(c))
			}
			s.addUniqueStr(&s.Report, "acceptedColors", strings.Join(colorStrs, ","))
		}
		time.Sleep(time.Duration(s.Options.RequestDelayMs) * time.Millisecond)
	}
}

func (s *ProfileStyleService) testStyleVariant(working map[string]interface{}, category, label string, style Style) APIResult {
	payload := s.buildPayload(working["payloadFormat"].(string), style)
	return s.testPayload(map[string]interface{}{
		"phase": fmt.Sprintf("compatibility-%s", category),
		"endpoint": map[string]interface{}{
			"key":      working["endpointKey"],
			"label":    working["endpointLabel"],
			"endpoint": working["endpoint"],
			"guildId":  working["guildId"],
		},
		"payloadFormat": working["payloadFormat"],
		"payload":       payload,
		"style":         style,
		"label":         label,
	})
}

func (s *ProfileStyleService) applyFinalStyle(working map[string]interface{}) []APIResult {
	presetLabel := "Custom Display Name Style"
	if s.SelectedStylePreset != nil {
		presetLabel = s.SelectedStylePreset.Label
	}
	s.writeSummary(fmt.Sprintf("Applying final target Display Name Style: %s.", presetLabel))
	responses := s.applyStyleToConfiguredScope(working, *s.Options.TargetStyle, "final-apply", presetLabel)
	
	appliedGuildIds := []string{}
	confirmedCount := 0
	for _, r := range responses {
		if s.isStyleConfirmed(r, *s.Options.TargetStyle) {
			confirmedCount++
			if r.GuildID != nil && *r.GuildID != "" {
				appliedGuildIds = append(appliedGuildIds, *r.GuildID)
			}
		}
	}

	if confirmedCount > 0 {
		working["style"] = *s.Options.TargetStyle
		working["appliedAt"] = time.Now().UTC().Format(time.RFC3339)
		if working["guildScoped"].(bool) {
			working["appliedGuildIds"] = appliedGuildIds
		}
		s.Report["finalWorkingConfiguration"] = working
		s.saveWorkingConfig(working)
	}
	return responses
}

func (s *ProfileStyleService) applyStyleToConfiguredScope(working map[string]interface{}, style Style, phase, label string) []APIResult {
	guildIDs := []string{""}
	if working["guildScoped"].(bool) {
		gID := ""
		if val, ok := working["guildId"].(string); ok {
			gID = val
		}
		guildIDs = []string{gID}
	}

	var responses []APIResult
	for _, guildID := range guildIDs {
		payload := s.buildPayload(working["payloadFormat"].(string), style)
		
		endpointStr := working["endpoint"].(string)
		if working["guildScoped"].(bool) && guildID != "" {
			template, ok := working["endpointTemplate"].(string)
			if ok {
				endpointStr = strings.ReplaceAll(template, "{guild_id}", guildID)
			}
		}

		endpoint := map[string]interface{}{
			"key":         working["endpointKey"],
			"label":       working["endpointLabel"],
			"endpoint":    endpointStr,
			"guildId":     guildID,
			"guildScoped": working["guildScoped"],
		}

		response := s.testPayload(map[string]interface{}{
			"phase":         phase,
			"endpoint":      endpoint,
			"payloadFormat": working["payloadFormat"],
			"payload":       payload,
			"style":         style,
			"label":         label,
		})
		
		if guildID != "" {
			response.GuildID = &guildID
		}
		responses = append(responses, response)
		time.Sleep(time.Duration(s.Options.RequestDelayMs) * time.Millisecond)
	}

	return responses
}

func (s *ProfileStyleService) isStyleConfirmed(response APIResult, expectedStyle Style) bool {
	if !response.Ok {
		return false
	}
	if response.VerifiedStyle {
		return true
	}
	return s.responseContainsStyle(response, expectedStyle)
}

func (s *ProfileStyleService) responseContainsStyle(response APIResult, expectedStyle Style) bool {
	return s.objectContainsStyle(response.Parsed, expectedStyle)
}

func (s *ProfileStyleService) objectContainsStyle(value map[string]interface{}, expectedStyle Style) bool {
	if value == nil {
		return false
	}

	if stylesObj, ok := value["display_name_styles"].(map[string]interface{}); ok {
		if s.styleMatches(stylesObj, expectedStyle) {
			return true
		}
	}
	
	if s.styleMatches(value, expectedStyle) {
		return true
	}

	// Recursive check (simplified for map structures in Go)
	for _, v := range value {
		if childMap, ok := v.(map[string]interface{}); ok {
			if s.objectContainsStyle(childMap, expectedStyle) {
				return true
			}
		}
	}

	return false
}

func (s *ProfileStyleService) styleMatches(value map[string]interface{}, expectedStyle Style) bool {
	if value == nil {
		return false
	}
	fID, fOK := value["font_id"].(float64)
	eID, eOK := value["effect_id"].(float64)
	if !fOK && !eOK {
		// Try flat structure
		fID, fOK = value["display_name_font_id"].(float64)
		eID, eOK = value["display_name_effect_id"].(float64)
	}

	if !fOK || !eOK || int(fID) != expectedStyle.FontID || int(eID) != expectedStyle.EffectID {
		return false
	}

	var colorInterface []interface{}
	if cols, ok := value["colors"].([]interface{}); ok {
		colorInterface = cols
	} else if cols, ok := value["display_name_colors"].([]interface{}); ok {
		colorInterface = cols
	} else {
		return false
	}

	if len(colorInterface) != len(expectedStyle.Colors) {
		return false
	}

	for i, c := range colorInterface {
		if val, ok := c.(float64); !ok || int(val) != expectedStyle.Colors[i] {
			return false
		}
	}

	return true
}

func (s *ProfileStyleService) testPayload(params map[string]interface{}) APIResult {
	phase := params["phase"].(string)
	endpoint := params["endpoint"].(map[string]interface{})
	payloadFormat := params["payloadFormat"].(string)
	payload := params["payload"]
	style := params["style"].(Style)
	label, _ := params["label"].(string)

	s.logDisplayNameStyleTest(map[string]interface{}{
		"phase":         phase,
		"endpoint":      endpoint,
		"payloadFormat": payloadFormat,
		"payload":       payload,
		"style":         style,
		"label":         label,
		"status":        "STARTED",
	})

	response := s.API.Patch(endpoint["endpoint"].(string), payload)
	response.VerifiedStyle = s.responseContainsStyle(response, style)

	resultStr := "FAILURE"
	if s.isStyleConfirmed(response, style) {
		resultStr = "SUCCESS_CONFIRMED"
	} else if response.Ok {
		resultStr = "SUCCESS_UNCONFIRMED"
	}

	s.logDisplayNameStyleTest(map[string]interface{}{
		"phase":         phase,
		"endpoint":      endpoint,
		"payloadFormat": payloadFormat,
		"payload":       payload,
		"style":         style,
		"label":         label,
		"status":        response.Status,
		"response":      response,
		"result":        resultStr,
	})

	if !response.Ok && SUPPORTED_ERROR_STATUSES[response.Status] {
		s.Report["notes"] = append(s.Report["notes"].([]string), fmt.Sprintf("%s returned %d during %phase.", endpoint["label"], response.Status, phase))
	}

	return response
}

func (s *ProfileStyleService) buildPayload(format string, style Style) interface{} {
	if format == "A" {
		return map[string]interface{}{"display_name_styles": style}
	}
	return map[string]interface{}{
		"display_name_font_id":   style.FontID,
		"display_name_effect_id": style.EffectID,
		"display_name_colors":    style.Colors,
	}
}

func (s *ProfileStyleService) recordEndpointResult(endpoint map[string]interface{}, payloadFormat string, response APIResult) {
	label := endpoint["label"].(string)
	endpointsMap := s.Report["endpoints"].(map[string]interface{})
	
	if _, exists := endpointsMap[label]; !exists {
		endpointsMap[label] = make(map[string]interface{})
	}
	
	formatMap := endpointsMap[label].(map[string]interface{})
	
	supported := "NO"
	if s.isStyleConfirmed(response, s.Options.TargetStyle != nil && *s.Options.TargetStyle != Style{} ? *s.Options.TargetStyle : Style{}) {
		supported = "YES"
	} else if response.Ok {
		supported = "UNCONFIRMED"
	}

	errMsg := response.StatusText
	errCode := 0
	if response.Parsed != nil {
		if msg, ok := response.Parsed["message"].(string); ok {
			errMsg = msg
		}
		if code, ok := response.Parsed["code"].(float64); ok {
			errCode = int(code)
		}
	}

	formatMap[payloadFormat] = map[string]interface{}{
		"status":      response.Status,
		"supported":   supported,
		"rateLimited": response.RateLimited,
		"errorCode":   errCode,
		"message":     errMsg,
	}
}

func (s *ProfileStyleService) captureUnsupportedFields(response APIResult) {
	if response.Parsed == nil {
		return
	}
	if errors, ok := response.Parsed["errors"].(map[string]interface{}); ok {
		fields := s.extractErrorFields(errors, "")
		for _, f := range fields {
			s.addUniqueStr(&s.Report, "unsupportedFields", f)
		}
	}
}

func (s *ProfileStyleService) extractErrorFields(errors map[string]interface{}, prefix string) []string {
	var fields []string
	for key, value := range errors {
		if key == "_errors" {
			continue
		}
		next := key
		if prefix != "" {
			next = prefix + "." + key
		}
		
		if valMap, ok := value.(map[string]interface{}); ok {
			if _, hasErrors := valMap["_errors"]; hasErrors {
				fields = append(fields, next)
			}
			fields = append(fields, s.extractErrorFields(valMap, next)...)
		}
	}
	return fields
}

func (s *ProfileStyleService) writeSummary(message string) {
	s.Report["notes"] = append(s.Report["notes"].([]string), message)
	s.logEvent("summary", map[string]interface{}{"message": message})
	fmt.Printf("[DisplayNameStyles] %s | Main: KyronixStudio | High Partner: dray.me\n", message)
}

func (s *ProfileStyleService) logDisplayNameStyleTest(entry map[string]interface{}) {
	entry["separator"] = "━━━━━━━━━━━━━━━━━━"
	entry["title"] = "Display Name Style Test"
	s.logEvent("display-name-style-test", s.sanitizeLogEntry(entry))
}

func (s *ProfileStyleService) logEvent(eventType string, data interface{}) {
	entry := map[string]interface{}{
		"type":      eventType,
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		"data":      data,
	}
	b, _ := json.Marshal(entry)
	f, err := os.OpenFile(s.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		f.Write(append(b, '\n'))
		f.Close()
	}
}

func (s *ProfileStyleService) sanitizeLogEntry(entry interface{}) interface{} {
	if entry == nil {
		return nil
	}
	
	// Complex deep sanitization logic can go here. For simplicity in Go, 
	// we avoid full reflection-based recursion unless strictly needed, 
	// but we can scrub known keys if it's a map.
	if m, ok := entry.(map[string]interface{}); ok {
		clean := make(map[string]interface{})
		for k, v := range m {
			if strings.Contains(strings.ToLower(k), "authorization") {
				clean[k] = "[REDACTED]"
			} else {
				clean[k] = v // Shallow pass for this Go port for performance
			}
		}
		return clean
	}
	return entry
}

func (s *ProfileStyleService) saveWorkingConfig(config map[string]interface{}) {
	b, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(s.CacheFile, b, 0644)
}

func (s *ProfileStyleService) loadWorkingConfig() map[string]interface{} {
	b, err := os.ReadFile(s.CacheFile)
	if err != nil {
		return nil
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(b, &parsed); err == nil {
		if parsed["endpoint"] != nil && parsed["payloadFormat"] != nil && parsed["style"] != nil {
			return parsed
		}
	}
	return nil
}

func (s *ProfileStyleService) saveRotationState(state map[string]interface{}) {
	b, _ := json.MarshalIndent(state, "", "  ")
	os.WriteFile(s.RotationFile, b, 0644)
}

func (s *ProfileStyleService) loadRotationState() map[string]interface{} {
	b, err := os.ReadFile(s.RotationFile)
	if err != nil {
		return make(map[string]interface{})
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(b, &parsed); err == nil {
		return parsed
	}
	return make(map[string]interface{})
}

func (s *ProfileStyleService) finalizeReport(message string) {
	s.Report["generatedAt"] = time.Now().UTC().Format(time.RFC3339)
	s.Report["notes"] = append(s.Report["notes"].([]string), message)
	
	reportContent := s.renderReport()
	os.WriteFile(s.ReportFile, []byte(reportContent), 0644)
	
	s.logEvent("report", s.Report)
	fmt.Printf("[DisplayNameStyles] %s | Powered by KyronixStudio & dray.me\n", message)
}

func (s *ProfileStyleService) renderReport() string {
	b1, _ := json.MarshalIndent(s.Report["finalWorkingConfiguration"], "", "  ")
	
	targetConfig := map[string]interface{}{}
	if s.Options.TargetStyle != nil {
		targetConfig["display_name_styles"] = *s.Options.TargetStyle
	}
	b2, _ := json.MarshalIndent(targetConfig, "", "  ")
	b3, _ := json.MarshalIndent(s.Report["availableStylePresets"], "", "  ")
	b4, _ := json.MarshalIndent(s.Report["endpoints"], "", "  ")

	var acceptedFonts []string
	for _, v := range s.Report["acceptedFontIds"].([]int) { acceptedFonts = append(acceptedFonts, strconv.Itoa(v)) }
	
	var acceptedEffects []string
	for _, v := range s.Report["acceptedEffectIds"].([]int) { acceptedEffects = append(acceptedEffects, strconv.Itoa(v)) }

	lines := []string{
		"# Display Name Styles Report",
		"",
		"## Credits",
		"- **Main Server**: KyronixStudio",
		"- **High Partner**: dray.me",
		"",
		fmt.Sprintf("Generated At: %s", s.Report["generatedAt"]),
		"",
		fmt.Sprintf("Endpoint Supported: %s", s.Report["endpointSupported"]),
		fmt.Sprintf("Bot Token Supported: %s", s.Report["botTokenSupported"]),
		fmt.Sprintf("Payload Format: %s", s.Report["payloadFormat"]),
		fmt.Sprintf("Accepted Font IDs: %s", strings.Join(acceptedFonts, ", ")),
		fmt.Sprintf("Accepted Effect IDs: %s", strings.Join(acceptedEffects, ", ")),
		fmt.Sprintf("Accepted Colors: %s", strings.Join(s.Report["acceptedColors"].([]string), " | ")),
		fmt.Sprintf("Unsupported Fields: %s", strings.Join(s.Report["unsupportedFields"].([]string), ", ")),
		"",
		"## Final Working Configuration",
		"```json", string(b1), "```",
		"",
		"## Final Target Configuration",
		"```json", string(b2), "```",
		"",
		"## Available Style Presets",
		"```json", string(b3), "```",
		"",
		"## Endpoint Results",
		"```json", string(b4), "```",
		"",
		"## Notes",
	}

	for _, note := range s.Report["notes"].([]string) {
		lines = append(lines, "- "+note)
	}
	lines = append(lines, "")

	return strings.Join(lines, "\n")
}

func (s *ProfileStyleService) addUniqueStr(report *map[string]interface{}, key, value string) {
	list := (*report)[key].([]string)
	for _, v := range list {
		if v == value {
			return
		}
	}
	(*report)[key] = append(list, value)
}

func (s *ProfileStyleService) addUniqueInt(report *map[string]interface{}, key string, value int) {
	list := (*report)[key].([]int)
	for _, v := range list {
		if v == value {
			return
		}
	}
	(*report)[key] = append(list, value)
}

func structToMap(obj interface{}) map[string]interface{} {
	data, _ := json.Marshal(obj)
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result
}

// --- Main Execution (Equivalent to Python's async initialize) ---

func main() {
	// You can pass your Discord token here or via the DISCORD_TOKEN environment variable
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		fmt.Println("Warning: No DISCORD_TOKEN provided in environment variables.")
	}

	opts := Options{
		Token:                 token,
		ForceDiscovery:        os.Getenv("DISCORD_PROFILE_STYLE_FORCE_DISCOVERY") == "true",
		RunCompatibilityTests: os.Getenv("DISCORD_PROFILE_STYLE_RUN_COMPATIBILITY") == "true",
		GuildID:               os.Getenv("DISCORD_PROFILE_STYLE_GUILD_ID"),
	}

	service := NewProfileStyleService(opts)
	service.Run()
}