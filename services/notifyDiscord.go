package services

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
	"trap_handler/models"
)

/*
NotifyDiscordWithMessage sends a Discord notification using a raw message string.
Parameters:
  - message: string
  - alertType: "warning", "resolved", "spam"
*/
func NotifyDiscord(message, alertType, threadID string, isCritical bool, mentionIDs []string) error {

	// Parse ASCII character
	message = strings.ReplaceAll(message, "%0A", "\n")
	lines := strings.Split(message, "\n")

	title := "❗️❗️❗️ 🚨 CẢNH BÁO ❗️❗️❗️"
	description := message
	if len(lines) > 1 {
		title = strings.TrimSpace(lines[0])
		description = strings.Join(lines[1:], "\n")
	}
	// Determine embed color
	color := models.ColorDefault

	// If alertType is "warning", use yellow color
	if strings.ToLower(alertType) == "warning" {
		color = models.ColorYellow
	}

	// If alertType is "spam", use blue color
	if strings.ToLower(alertType) == "spam" {
		color = models.ColorBlue
	}

	// Build embed
	embed := models.DiscordEmbed{
		Title:       title,
		Description: description,
		Color:       color,
	}

	// Just mentions if it is critical
	mentionPayload := ""
	if isCritical {
		color = models.ColorRed
		//Build mentions
		if len(mentionIDs) > 0 {
			for _, id := range mentionIDs {
				mentionPayload += "<@" + id + "> "
			}
		}
	}

	// Build the payload
	payload := models.DiscordPayload{
		Content: mentionPayload,
		Embeds:  []models.DiscordEmbed{embed},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal Discord payload:", err)
		return err
	}

	jsonStr := string(jsonData)

	// Construct curl command using string concatenation only
	baseCommand := `curl \
	-X POST \
	-H "Content-Type: application/json" \
	-s --connect-timeout 10 \
	-d '` + jsonStr + `' \
	"https://discord.com/api/webhooks/1384379242753818724/_LNnAZAOL55chbhrLj6lKwYxHzEUuYll_aD8pKzSFCpHpeVUIf3ypTEoPkDzpJ1oYYtM?thread_id=` + threadID

	log.Println("DISCORD CURL CMD:", baseCommand)

	output, err := exec.Command("bash", "-c", baseCommand).CombinedOutput()
	if err != nil {
		log.Println("Error sending to Discord:", string(output))
		return err
	}

	log.Println("Discord Output:", string(output))
	return nil
}

func BuildDiscordFiringMessage(message string) string {
	description := ""

	return description
}

func GetMentionedIDs(discordUsers map[string]string) []string {
	var listOfIDs []string

	for _, userID := range discordUsers {
		listOfIDs = append(listOfIDs, userID)
	}

	return listOfIDs
}
