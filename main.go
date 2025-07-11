package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
	"trap_handler/services"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var redisService *services.RedisService
var dcimService *services.DCIMService

func init() {
	// Load the .env file
	err := godotenv.Load("/root/scripts/trap_handler/.env")
	if err != nil {
		panic("Error loading .env file")
	}

	redisService, err = services.NewRedisService()

	if err != nil {
		panic(err)
	}

	dcimService = services.NewDCIMService()
}

func main() {
	// Create or open the log file for writing (use appropriate file path)
	logFile, err := os.OpenFile("/root/scripts/trap_handler/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer func(logFile *os.File) {
		_ = logFile.Close()
	}(logFile)

	// Create a custom logger that writes to the log file
	logger := log.New(logFile, "CUSTOM: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Notify log
	logger.Print("Start processing the trap")

	// Scan

	scanner := bufio.NewScanner(os.Stdin)

	var rawContent string

	// Read the input

	for scanner.Scan() {
		line := scanner.Text()
		rawContent += line + "\n"
	}

	// User trap content sample
	//<UNKNOWN>
	//UDP: [103.90.226.6]:51961->[14.225.207.112]:162
	//iso.3.6.1.2.1.1.3.0 36:18:51:30.34
	//iso.3.6.1.6.3.1.1.4.1.0 iso.3.6.1.4.1.6876.4.3.0.203
	//iso.3.6.1.4.1.6876.4.3.308.0 2
	//iso.3.6.1.4.1.6876.4.3.304.0 "Gray"
	//iso.3.6.1.4.1.6876.4.3.305.0 "Yellow"
	//iso.3.6.1.4.1.6876.4.3.306.0 "Mem > 80%25 - Tests - Metric Memory Host consumed % = 89%"
	//iso.3.6.1.4.1.6876.4.3.307.0 "firewall9.vietnix.vn"
	//iso.3.6.1.6.3.18.1.3.0 103.90.226.6
	//iso.3.6.1.6.3.18.1.4.0 "public"
	//iso.3.6.1.6.3.1.1.4.3.0 iso.3.6.1.4.1.6876.4.3
	//
	// Get line 8 and 9

	// System trap content sample
	//	static.vnpt.vn					\\ line 0
	//	UDP: [14.225.208.13]:47074->[14.225.207.112]:162
	//	iso.3.6.1.2.1.1.3.0 213:0:09:08.55
	//	iso.3.6.1.6.3.1.1.4.1.0 iso.3.6.1.4.1.6876.4.3.0.203
	//	iso.3.6.1.4.1.6876.4.3.308.0 4
	//	iso.3.6.1.4.1.6876.4.3.304.0 "Normal"
	//	iso.3.6.1.4.1.6876.4.3.305.0 "Warning"		\\ line 6
	//	iso.3.6.1.4.1.6876.4.3.306.0 "Memory Exhaustion on cloudvcenter - Event: Stats monitor detected resource utilization status change. (1717615)
	//	Summary: vCenter Memory Resource status changed from Green to Yellow on cloudvcenter.vietnix.vn for continuous Memory utilization 85% in 0 mins
	//	Date: 02/22/2024 06:53:18 AM
	//	User name: VCENTER1.VIETNIX.VN\\machine-ac2dc892-0c79-4233-8658-c44d2f628075
	//	Arguments:
	//	eventTypeId = vim.event.ResourceExhaustionStatusChangedEvent
	//	severity = info
	//	resourceName = mem_usage
	//	oldStatus = green
	//	newStatus = yellow \\ line 16
	//	reason = for continuous Memory utilization 85% in 0 mins
	//	nodeType = vcenter
	//	_sourcehost_ = cloudvcenter.vietnix.vn
	//	"
	//	iso.3.6.1.4.1.6876.4.3.307.0 "Datacenters"
	//	iso.3.6.1.6.3.18.1.3.0 10.23.99.3
	//	iso.3.6.1.6.3.18.1.4.0 "public"
	//	iso.3.6.1.6.3.1.1.4.3.0 iso.3.6.1.4.1.6876.4.3

	// OR
	//<UNKNOWN>
	//	UDP: [103.90.226.6]:44828->[14.225.207.112]:162
	//iso.3.6.1.2.1.1.3.0 22:22:47:32.73
	//iso.3.6.1.6.3.1.1.4.1.0 iso.3.6.1.4.1.6876.4.3.0.203
	//iso.3.6.1.4.1.6876.4.3.308.0 2
	//iso.3.6.1.4.1.6876.4.3.304.0 "Gray"
	//iso.3.6.1.4.1.6876.4.3.305.0 "Gray"		// line 6
	//iso.3.6.1.4.1.6876.4.3.306.0 "alarm.HostConnectivityAlarm - Event: Host connection lost (26813982)
	//Summary: Host cpuv4-s1.vietnix.vn in VPS-CPUv4 is not responding
	//Date: 02/22/2024 07:14:23 AM
	//Host: cpuv4-s1.vietnix.vn
	//Resource pool: cpuv4-s1.vietnix.vn
	//Data center: VPS-CPUv4
	//"
	//iso.3.6.1.4.1.6876.4.3.307.0 "cpuv4-s1.vietnix.vn"
	//iso.3.6.1.6.3.18.1.3.0 103.90.226.6
	//iso.3.6.1.6.3.18.1.4.0 "public"
	//iso.3.6.1.6.3.1.1.4.3.0 iso.3.6.1.4.1.6876.4.3

	//logger.Print(rawContent)

	lines := strings.Split(rawContent, "\n")

	alarmName := ""
	currentMetric := ""

	if !strings.Contains(lines[8], "iso.") {
		// This is system trap
		alarmName = strings.Split(lines[7], `"`)[1]
		// alarmName: alarm.HostConnectivityAlarm - Event: Host connection lost (26813982)
		if !strings.Contains(alarmName, "HostConnectivityAlarm") {
			// Do nothing
			return
		}
		// Get Event: Host connection lost (26813982)
		alarmName = strings.Split(alarmName, " - Event: ")[1]
		alarmName = strings.Split(alarmName, " (")[0]
		currentMetric = lines[8]
	} else {
		// this is user trap
		metricsRaw := strings.Split(lines[7], `"`)[1]
		alarmNameRaw := strings.Split(metricsRaw, `-`)
		currentMetric = alarmNameRaw[len(alarmNameRaw)-1]
		alarmName = strings.Replace(metricsRaw, " -"+currentMetric, "", -1)
		currentMetric = strings.TrimSpace(currentMetric)
	}

	vcenterIP := strings.Split(lines[len(lines)-4], ` `)[1]
	vcenter := "cloudvcenter.vietnix.vn"
	if vcenterIP == "103.90.226.6" {
		vcenter = "vcenter.vietnix.vn"
	}
	host := strings.Split(lines[len(lines)-5], `"`)[1]

	//// Print the content to log file
	//logger.Print(alarmName + "\n" + currentMetric + "\n" + host)

	if err = scanner.Err(); err != nil {
		logger.Fatal("Error reading standard input:", err)
	}

	// Notify to telegram
	botToken := os.Getenv("BOT_TOKEN")
	chatId := os.Getenv("CHAT_ID")
	threadId := os.Getenv("THREAD_ID")
	isThread := os.Getenv("IS_THREAD")

	// env variables for discord
	criticalThread := "1383266646499790888"
	nonCriticalThread := "1383018751154327582"

	var DISCORD_MENTIONED_IDS = []string{
		"618423557021761547",
		"758578700081037312",
		"1375668363258494996",
	}

	criticalBotToken := os.Getenv("CRITICAL_BOT_TOKEN")
	//criticalChatId := os.Getenv("CRITICAL_CHAT_ID")
	criticalThreadId := os.Getenv("CRITICAL_THREAD_ID")
	//criticalIsThread := os.Getenv("CRITICAL_IS_THREAD")

	// Replace special characters
	alarmName = strings.Replace(alarmName, "[CUSTOM] ", "", -1)

	template := "[SIGN] *CẢNH BÁO TÀI NGUYÊN [TARGET]* [SIGN]%0A%0A[CONTENT]%0A[CRITICAL_MESSAGE]%0A_* Send from: 14.225.207.112 _"

	// Using redis to check double alarm
	key := base64.StdEncoding.EncodeToString([]byte("alarm-vm-vcenter-:" + host + ":" + alarmName))
	// REdis here

	// Check if the key exists
	ok, err := redisService.CheckKeyRedis(key)
	if ok && !strings.Contains(host, "vietnix.vn") {
		logger.Print("The key exists: " + key)
		return
	}
	if err != nil {
		if !strings.Contains(err.Error(), "redis: nil") {
			logger.Print(err)
		}
	}

	//logger.Println(content)
	currentMetricRaw := strings.Split(currentMetric, "=")
	if len(currentMetricRaw) >= 2 {
		currentMetric = strings.TrimSpace(currentMetricRaw[0]) + ": *" + currentMetricRaw[1] + "*"
	} else {
		log.Println("Current metric:" + currentMetric)
	}

	logger.Print("Telegram Notify: " + alarmName)

	duration := 12 * 60 * time.Minute

	// Send notify to non critical thread
	if strings.Contains(alarmName, "VM") {
		if strings.Contains(host, "firewall") {
			// Do nothing
			return
		}
		content := " - Alert Type: " + alarmName + "%0A - " + currentMetric + "%0A - VM: *" + host + "*%0A - vCenter: " + vcenter
		message := strings.Replace(template, "[TARGET]", "VM", -1)
		message = strings.Replace(message, "[CONTENT]", content, -1)
		message = strings.Replace(message, "[SIGN]", "⚠️", -1)
		message = strings.Replace(message, "[CRITICAL_MESSAGE]", "", -1)

		_, err := NotifyTelegram(botToken, chatId, threadId, isThread, message)
		if err != nil {

			logger.Print(err)
			return
		}

		err = services.NotifyDiscord(message, "warning", nonCriticalThread, false, DISCORD_MENTIONED_IDS)

		if err != nil {

			logger.Print(err)
			return
		}
	} else {
		if strings.Contains(alarmName, "CPU") {
			// Do nothing
			return
		}

		content := " - Alert Type: " + alarmName + "%0A - " + currentMetric + "%0A - Node: *" + host + "*%0A - vCenter: " + vcenter
		if strings.Contains(alarmName, "Datastore") || strings.Contains(alarmName, "Mem") {
			duration = 3 * 24 * time.Hour
		}

		// Check if host is "Server for vietnix" role in DCIM
		res, err := dcimService.IsServerForVietnix(host)
		if err != nil {
			logger.Print(err.Error())
			content += "%0A - Err inCheck Vietnix Server Role: " + err.Error()
		}
		if !res {
			return
		}
		// END Check if host is "Server for vietnix" role in DCIM

		message := strings.Replace(template, "[TARGET]", "NODE", -1)
		message = strings.Replace(message, "[CONTENT]", content, -1)
		message = strings.Replace(message, "[SIGN]", "❗️", -1)
		message = strings.Replace(message, "[CRITICAL_MESSAGE]", "", -1)
		//message = strings.ReplaceAll(message, "[CRITICAL_MESSAGE]", "️%0A[@nguyenhoang91](tg://user?id=331113301) [@thucdduy](tg://user?id=159728680) [@Ox54616E5461](tg://user?id=362157387) [@nhatkini](tg://user?id=334166509) [@KhanhTruong](tg://user?id=702361047) [@imlowkey](tg://user?id=482047370) [@nightbarron](tg://user?id=1753149166)%0A")
		if strings.Contains(alarmName, "Datastore > 90 Percent") || strings.Contains(alarmName, "Mem > 90 Percent") {
			err = createTaskInWorkplace(6, "high", message)
			if err != nil {
				logger.Print(err)
			}
		}

		_, err = NotifyTelegram(criticalBotToken, "-1001682572909", criticalThreadId, "false", message)
		if err != nil {

			logger.Print(err)
			return
		}

		err = services.NotifyDiscord(message, "", criticalThread, true, DISCORD_MENTIONED_IDS)

		if err != nil {
			logger.Print(err)
			return
		}
	}

	// Add key to redis
	err = redisService.AddKeyRedis(key, "-", duration)
	if err != nil {
		logger.Print(err)
		return
	}

}

func NotifyTelegram(DefaultBotToken, DefaultChatID, ThreadID, isThread, messages string) (string, error) {
	baseCommand := `curl https://api.telegram.org/bot` + DefaultBotToken + `/sendMessage
	-X POST
	-s --connect-timeout 10
	-d chat_id=` + DefaultChatID + `
	-d parse_mode=Markdown
	-d text="` + messages + `"`

	if isThread == "true" {
		baseCommand = baseCommand + ` -d message_thread_id=` + ThreadID
	}

	// TODO: execute curl command
	out, err := RunCmd(baseCommand)
	if err != nil {
		return string(out), err
	}
	//os.R(base_command, nil, nil)

	return string(out), nil
}

func RunCmd(s ...string) ([]byte, error) {

	// Gen command and Remove duplicate whitespace from a command
	space := regexp.MustCompile(`\s+`)
	command := space.ReplaceAllString(fmt.Sprint(strings.Join(s[:], " ")), " ")

	//log.Infof("Run cmd> %s", command)

	// Chuyển CombinedOutput() => Output()
	str, err := exec.Command("bash", "-c", command).Output()
	return str, err
	//return exec.Command("bash", "-c", command).CombinedOutput()
}

func createTaskInWorkplace(whmcsClientId int, priority, description string) error {
	// Nguyen Hung
	whmcsClientId = 6
	url := "https://api.vietnix.vn/vtasks/vietnixbot"
	method := "POST"

	// Replace all %0A to \n
	description = strings.Replace(description, "%0A", "\n", -1)
	description = strings.Replace(description, "%2A", "*", -1)
	description = strings.Replace(description, "%25", "%", -1)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-telegram-id", "999999999999")
	req.Header.Add("api-key", "kAxEIi5d52TEMCxlKVY6tMyPyRqY3Vu9")

	data := map[string]interface{}{
		"vTaskCase":     "lrwAygJkO58W6X3",
		"priority":      priority,
		"whmcsClientId": whmcsClientId,
		"description":   description,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(strings.NewReader(string(payload)))

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	//log.Info("CreateTaskInWorkplace: Workplace API response: ", result)
	return err
}

func checkKeyRedis(client redis.Client, key string) (bool, error) {
	val, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	if val == "" {
		return false, nil
	}
	return true, nil
}

func addKeyRedis(client redis.Client, key string, duration time.Duration) error {
	err := client.Set(context.Background(), key, "1", duration).Err()
	if err != nil {
		return err
	}
	return nil
}
