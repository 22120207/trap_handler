package models

// DiscordEmbed represents a Discord embed message
type DiscordEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
}

// DiscordButton represents an individual button in Discord
type DiscordButton struct {
	Type     int    `json:"type"`  // 2 for button
	Style    int    `json:"style"` // 1=Primary, 2=Secondary, 3=Success, 4=Danger, 5=Link
	Label    string `json:"label"` // Button text
	CustomID string `json:"custom_id,omitempty"`
	URL      string `json:"url,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

// A struct of Discord Payload Data
type DiscordPayload struct {
	Content string         `json:"content"`
	Embeds  []DiscordEmbed `json:"embeds"`
}

// Embed colors
const (
	ColorDefault = 0
	ColorRed     = 15158332
	ColorGreen   = 3066993
	ColorYellow  = 15844367
	ColorBlue    = 3447003
)
