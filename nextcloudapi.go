package nextcloudapi

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// Package structs
type NextcloudApi struct {
	NextcloudConfig *NextcloudConfig
}
type NextcloudConfig struct {
	Username     string
	Password     string
	NextcloudUrl string
}
type ApiRoutes struct {
	users  string
	groups string
}

var (
	apiUrlPath = "/ocs/v1.php"
	apiRoutes  = ApiRoutes{
		users:  apiUrlPath + "/cloud/users",
		groups: apiUrlPath + "/cloud/groups",
	}
)

func getCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
}

// Call Helper for executing HTTP requests to Nextcloud API with basic auth
//
// return: OCS, error
func (ncapi *NextcloudApi) Call(url, method string, data io.Reader) (OCS, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(method, ncapi.NextcloudConfig.NextcloudUrl+url, data)
	if err != nil {
		return OCS{}, fmt.Errorf("%s: %s", getCurrentFuncName(), err)
	}

	req.Header.Set("OCS-APIRequest", "true")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(ncapi.NextcloudConfig.Username, ncapi.NextcloudConfig.Password)

	response, err := client.Do(req)
	if err != nil {
		return OCS{}, fmt.Errorf("%s: %s", getCurrentFuncName(), err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return OCS{}, fmt.Errorf("%s: %s", getCurrentFuncName(), err)
	}
	var osc OCS
	err = xml.Unmarshal(body, &osc)
	if err != nil {
		return OCS{}, fmt.Errorf("%s: %s", getCurrentFuncName(), err)
	}

	return osc, nil
}

// ListUsers search for users or list all if no arguments passed.
//
// return: OCS, error
func (ncapi *NextcloudApi) ListUsers(searchString ...string) (OCS, error) {
	//log.Println(searchString)
	urlpath := apiRoutes.users
	if len(searchString) > 0 {
		urlpath = fmt.Sprintf("%s?search=%s", urlpath, url.QueryEscape(searchString[0]))
	}
	//log.Println(urlpath)
	resp, err := ncapi.Call(urlpath, "GET", nil)
	if err != nil {
		return OCS{}, fmt.Errorf("%s: %s", getCurrentFuncName(), err)
	}

	return resp, nil
}

// AddUser
//
// return: OCS, error
func (ncapi *NextcloudApi) AddUser(user User) (OCS, error) {
	data := url.Values{
		"userid":      {user.Userid},
		"password":    {user.Password},
		"displayName": {user.DisplayName},
		"email":       {user.Email},
		"quota":       {user.Quota},
		"language":    {user.Language},
	}
	if len(user.Groups) > 0 {
		for gr := 0; gr < len(user.Groups); gr++ {
			data.Add("groups[]", user.Groups[gr])
		}
	}
	if len(user.Subadmin) > 0 {
		for gr := 0; gr < len(user.Subadmin); gr++ {
			data.Add("subadmin[]", user.Subadmin[gr])
		}
	}
	resp, err := ncapi.Call(apiRoutes.users, "POST", strings.NewReader(data.Encode()))
	if err != nil {
		return OCS{}, fmt.Errorf("%s: %s", getCurrentFuncName(), err)
	}

	return resp, nil
}

// AddUserSimple just enter email and user will be added with userid as in email before @
//
// return: OCS, error
func (ncapi *NextcloudApi) AddUserSimple(email string) (OCS, error) {
	var username string
	match, err := regexp.MatchString("^[a-zA-Z0-9.!#$%&’*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-][a-zA-Z0-9-]+)+$", email)
	if err != nil {
		return OCS{}, fmt.Errorf("%s: %s", getCurrentFuncName(), err)
	}
	if !match {
		return OCS{}, fmt.Errorf("%s: email '%s' is not valid", getCurrentFuncName(), email)
	}
	username = strings.Split(email, "@")[0]
	user := User{
		Userid: username,
		Email:  email,
	}

	return ncapi.AddUser(user)
}

// DelUser
//
// return: OCS, error
func (ncapi *NextcloudApi) DelUser(username string) (OCS, error) {
	resp, err := ncapi.Call(fmt.Sprintf(apiRoutes.users+"/%s", username), "DELETE", nil)
	if err != nil {
		return OCS{}, fmt.Errorf("%s: %s", getCurrentFuncName(), err)
	}

	return resp, nil
}

// AddUserInGroup
//
// return: OCS, error
func (ncapi *NextcloudApi) AddUserInGroup(username, group string) (OCS, error) {
	data := url.Values{
		"groupid": {group},
	}
	resp, err := ncapi.Call(fmt.Sprintf(apiRoutes.users+"/%s/groups", username),
		"POST", strings.NewReader(data.Encode()))
	if err != nil {
		return OCS{}, fmt.Errorf("%s: %s", getCurrentFuncName(), err)
	}

	return resp, nil
}

// DeleteUserFromGroup
//
// return: OCS, error
func (ncapi *NextcloudApi) DeleteUserFromGroup(username, group string) (OCS, error) {
	data := url.Values{
		"groupid": {group},
	}
	resp, err := ncapi.Call(fmt.Sprintf(apiRoutes.users+"/%s/groups", username),
		"DELETE", strings.NewReader(data.Encode()))
	if err != nil {
		return OCS{}, fmt.Errorf("%s: %s", getCurrentFuncName(), err)
	}

	return resp, nil
}

// ResendEmail resend welcome email message for password setup
//
// return: OCS, error
func (ncapi *NextcloudApi) ResendEmail(username string) (OCS, error) {
	resp, err := ncapi.Call(fmt.Sprintf(apiRoutes.users+"/%s/welcome", username), "POST", nil)
	if err != nil {
		return OCS{}, fmt.Errorf("%s: %s", getCurrentFuncName(), err)
	}

	return resp, nil
}

// DisableUser
//
// return: OCS, error
func (ncapi *NextcloudApi) DisableUser(username string) (OCS, error) {
	resp, err := ncapi.Call(fmt.Sprintf(apiRoutes.users+"/%s/disable", username), "PUT", nil)
	if err != nil {
		return OCS{}, fmt.Errorf("%s: %s", getCurrentFuncName(), err)
	}

	return resp, nil
}

// EnableUser
//
// return: OCS, error
func (ncapi *NextcloudApi) EnableUser(username string) (OCS, error) {
	resp, err := ncapi.Call(fmt.Sprintf(apiRoutes.users+"/%s/enable", username), "PUT", nil)
	if err != nil {
		return OCS{}, fmt.Errorf("%s: %s", getCurrentFuncName(), err)
	}

	return resp, nil
}
