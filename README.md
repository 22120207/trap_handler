# trap_handler

```bash
curl -X POST "$DISCORD_WEBHOOK_URL" \
    -H "Content-Type: application/json" \
    -d @- <<EOF
        {
        "username": "CI/CD Bot",
        "embeds": [
            {
            "title": "✅ Build Success",
            "color": 3066993,
            "fields": [
                { "name": "Author", "value": "${GITHUB_ACTOR}", "inline": true },
                { "name": "Repository", "value": "**${GITHUB_REPOSITORY}**", "inline": true },
                { "name": "Branch", "value": "${GITHUB_REF}", "inline": true },
                { "name": "Date", "value": "$(TZ='GMT+7' date)", "inline": false },
                { "name": "Commit Message", "value": "\"${COMMIT_MESSAGE}\"", "inline": false }
            ]
            }
        ]
        }
EOF
```