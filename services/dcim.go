package services

import (
	"encoding/json"
	"log"
	"strings"
	"time"
	"trap_handler/helpers"
)

const (
	DCIM_URL   = "https://dcim.vietnix.vn/api"
	DCIM_TOKEN = "01aa6b319e2db6007e6012d609f65a5bbb3897cf"
)

type DCIMService struct {
	token string
	url   string
}

func NewDCIMService() *DCIMService {
	return &DCIMService{token: DCIM_TOKEN, url: DCIM_URL}
}

func (d *DCIMService) IsServerForVietnix(serverName string) (bool, error) {
	result := false
	redisKey := "CACHED_SERVER_FOR_VIETNIX_LST"
	// Query from cache
	cachedData, err := redisService.GetKeyRedis(redisKey)
	if err == nil {
		deviceList := strings.Split(cachedData, ",")
		for _, device := range deviceList {
			if device == serverName {
				return true, nil
			}
		}

		// If not found in cache, return false
		return false, nil
	} // else, query from DCIM

	// Query from DCIM
	queryStr := "/dcim/devices?role_id=16&limit=0&offset=0"
	rawResponse, err := d.query("GET", queryStr)
	if err != nil {
		time.Sleep(1 * time.Second) // Retry after 1 second
		rawResponse, err = d.query("GET", queryStr)
		if err != nil {
			log.Println("error in get all device by role:", err.Error())
			return false, err
		}
	}

	var deviceList []string

	for _, device := range rawResponse["results"].([]interface{}) {
		deviceName := device.(map[string]interface{})["name"].(string)
		if helpers.IsDomainFormat(deviceName) {
			deviceList = append(deviceList, deviceName)
			if deviceName == serverName {
				result = true
			}
		}
	}

	// Add backup to redis to cache
	deviceListStr := strings.Join(deviceList, ",")
	_ = redisService.AddKeyRedis(redisKey, deviceListStr, 3*time.Minute)

	// ====

	return result, nil
}

func (d *DCIMService) query(reqMethod, queryStr string) (response map[string]interface{}, err error) {
	response = nil

	reqUrl := d.url + queryStr
	reqHeader := map[string]string{
		"Authorization": "Token " + d.token,
		"Content-Type":  "application/json",
	}

	responseRaw, err := helpers.RequestToAPI(reqUrl, reqMethod, reqHeader, nil, 30)
	if err != nil {
		return
	}

	// Convert json to map
	err = json.Unmarshal(responseRaw, &response)
	if err != nil {
		return
	}

	return
}
