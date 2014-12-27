package yota

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
)

type Yota struct {
	login    string
	password string
	uid      string
	http     *http.Client
}

type Tariff struct {
	Product string
	Name    string
	Amount  float64
	Code    string
	Speed   string
	Active  bool
}

var (
	reTariffs = regexp.MustCompile(
		`var sliderData = (.*);\n`)
	reBalance = regexp.MustCompile(
		`<dd id="balance-holder"><span>(\d+)</span>`)
)

const (
	urlLoginSuccess = "https://my.yota.ru/devices"
	urlLoginFail    = "https://my.yota.ru/selfcare/loginError"
	urlDevices      = "https://my.yota.ru/devices"
	urlUidByMail    = "https://my.yota.ru/selfcare/login/getUidByMail"
	urlChangeTariff = "https://my.yota.ru/selfcare/devices/changeOffer"
)

func New(login, password string) *Yota {
	cookies, _ := cookiejar.New(nil)

	http := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Jar: cookies,
	}

	yota := &Yota{
		login:    login,
		password: password,
		http:     http,
	}

	return yota
}

func (yota *Yota) Login() error {
	if yota.uid == "" {
		uid, err := yota.getUid()
		if err != nil {
			return err
		}

		yota.uid = uid
	}

	payload := url.Values{
		"goto":       {urlLoginSuccess},
		"gotoOnFail": {urlLoginFail},
		"org":        {"customer"},
		"old-token":  {yota.login},
		"IDToken2":   {yota.password},
		"IDToken1":   {yota.uid},
	}

	resp, err := yota.http.PostForm(
		"https://login.yota.ru/UI/Login",
		payload,
	)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	redirected := resp.Request.URL.String()
	if redirected != urlLoginSuccess {
		return errors.New("redirected to not success url")
	}

	return nil
}

func (yota *Yota) GetTariffs() ([]Tariff, error) {
	resp, err := yota.http.Get(urlDevices)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	matches := reTariffs.FindStringSubmatch(string(body[:]))
	if len(matches) == 0 {
		return nil, errors.New("could not find raw data")
	}
	rawData := matches[1]

	decoded := map[string]interface{}{}
	err = json.Unmarshal([]byte(rawData), &decoded)
	if err != nil {
		return nil, err
	}

	tariffs := []Tariff{}
	for product, data := range decoded {
		current := parseTariff(
			data.(map[string]interface{})["currentProduct"].(map[string]interface{}))
		steps := data.(map[string]interface{})["steps"].([]interface{})

		for _, step := range steps {
			tariff := parseTariff(step.(map[string]interface{}))
			if tariff.Code == current.Code {
				tariff.Active = true
			}

			tariff.Product = product

			tariffs = append(tariffs, tariff)
		}

		break
	}

	return tariffs, nil
}

func parseTariff(step map[string]interface{}) Tariff {
	tariff := Tariff{
		Name:  step["name"].(string),
		Code:  step["code"].(string),
	}

	amount, _ := strconv.ParseFloat(step["amountNumber"].(string), 64)
	tariff.Amount = amount

	speed := step["speedNumber"].(string)
	if speed == "<div class=\"max-value\">Макс.</div>" {
		speed = "max"
	}
	tariff.Speed = speed

	return tariff
}

func (yota *Yota) getUid() (string, error) {
	payload := url.Values{
		"value": {yota.login},
	}

	resp, err := yota.http.PostForm(
		urlUidByMail,
		payload,
	)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	status := string(body[0:2])

	if status != "ok" {
		return "", errors.New("status is not ok")
	}

	//stupid parsing...
	return string(body[3:]), nil
}

func (yota *Yota) GetBalance() (int, error) {
	resp, err := yota.http.Get(urlDevices)
	defer resp.Body.Close()

	if err != nil {
		return 0, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	matches := reBalance.FindStringSubmatch(string(body[:]))
	if len(matches) == 0 {
		return 0, errors.New("could not find balance data")
	}

	balance, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (yota *Yota) ChangeTariff(tariff Tariff) error {
	payload := url.Values{
		"product":       {tariff.Product},
		"offerCode":     {tariff.Code},
		"currentDevice": {"1"},
	}

	resp, err := yota.http.PostForm(
		urlChangeTariff,
		payload,
	)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
