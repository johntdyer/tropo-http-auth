package main

import (
	"bufio"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	user     string
	password string
)

// PapiResponse - Provisioning API response
type PapiResponse struct {
	Address string `json:"address"`
	Roles   []struct {
		Href     string `json:"href"`
		Role     string `json:"role"`
		RoleName string `json:"roleName"`
	}
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	if os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
		log.Debug("Starting in debug mode")

	} else {
		log.SetLevel(log.InfoLevel)
	}

}

// fetchProvisioningAddressConfig - Return struct from PAPI of user role config data
func fetchUserRoles(user string, password string) (*PapiResponse, error) {

	client := &http.Client{}
	respJSON := &PapiResponse{}

	req, err := http.NewRequest("GET", "https://api.tropo.com/users/"+user+"/roles", nil)
	req.Close = true
	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)

	// defer close
	defer resp.Body.Close()

	if err != nil {
		log.Error("PAPI Error : %s", err)
		return respJSON, err
	}

	fields := log.Fields{
		"responseCode": resp.StatusCode,
		"user":         user,
	}

	if resp.StatusCode == 200 {

		body, _ := ioutil.ReadAll(resp.Body)

		log.WithFields(fields).Debug("PAPI Request successful")

		err := json.Unmarshal(body, &respJSON.Roles)
		if err != nil {
			log.Error("PAPI Error : %s", err)
			return respJSON, err
		}

	} else {
		log.Warnf("Non-2xx response from papi - %v", resp.StatusCode)
	}
	return respJSON, nil
}

func main() {
	// user := ""
	// password := ""
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// check if err == io.EOF
			break
		}

		// Has username been set?
		if user == "" {
			user = strings.Trim(string(line), "\n")
		} else {
			// Second line is password
			password = strings.Trim(string(line), "\n")
		}

	}

	// Get environment values passed into mod_authnz_external and add to log fields
	fieldData := log.Fields{
		"user":     user,
		"ip":       os.Getenv("IP"),
		"uri":      os.Getenv("URI"),
		"httpHost": os.Getenv("HOST"),
	}

	// Get users roles from Provisioning
	roles, err := fetchUserRoles(user, password)
	if err != nil {
		log.WithFields(fieldData).Error("Failed to fetch roles: %s", err)
	}

	// Search all roles returned from PAPI and search for EMPLOYEE role
	for _, role := range roles.Roles {
		if role.RoleName == "EMPLOYEE" {
			log.WithFields(fieldData).Debugf("%s is an EMPLOYEE", user)
			os.Exit(0)
		}
	}

	// If we find no empoyee role we log and then exit 1
	log.WithFields(fieldData).Warnf("%s is not an EMPLOYEE", user)
	os.Exit(1)

}
