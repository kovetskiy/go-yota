package yota

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Client struct {
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
		`<dd id="balance-holder"><span>(\d*(,\d*)?)</span>\s([а-я]+\.)`)
	reRemains = regexp.MustCompile(
		`<div class="tarriff-info">\n\s+<div class="time">\n\s*\n\s*<strong>(\d+)</strong>\s*<span>([а-я]+)\&nbsp;([а-я]+)</span>`)
)

const (
	urlLoginSuccess = "https://my.yota.ru/devices"
	urlLoginFail    = "https://my.yota.ru/selfcare/loginError"
	urlDevices      = "https://my.yota.ru/devices"
	urlUidByMail    = "https://my.yota.ru/selfcare/login/getUidByMail"
	urlChangeTariff = "https://my.yota.ru/selfcare/devices/changeOffer"
)

func NewClient(login, password string, httpClient *http.Client) *Client {
	cookies, _ := cookiejar.New(nil)

	if httpClient == nil {
		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Jar: cookies,
		}
	}

	cli := &Client{
		login:    login,
		password: password,
		http:     httpClient,
	}

	return cli
}

func (cli *Client) Login() error {
	if cli.uid == "" {
		uid, err := cli.getUid()
		if err != nil {
			return err
		}

		cli.uid = uid
	}

	payload := url.Values{
		"goto":       {urlLoginSuccess},
		"gotoOnFail": {urlLoginFail},
		"org":        {"customer"},
		"old-token":  {cli.login},
		"IDToken2":   {cli.password},
		"IDToken1":   {cli.uid},
	}

	resp, err := cli.http.PostForm(
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

func (cli *Client) GetTariffs() ([]Tariff, error) {
	resp, err := cli.http.Get(urlDevices)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	matches := reTariffs.FindSubmatch(body)
	if len(matches) == 0 {
		return nil, errors.New("could not find raw data")
	}
	rawData := matches[1]

	decoded := map[string]map[string]interface{}{}
	err = json.Unmarshal(rawData, &decoded)
	if err != nil {
		return nil, err
	}

	tariffs := []Tariff{}
	for product, data := range decoded {
		current := parseTariff(data["currentProduct"].(map[string]interface{}))
		steps := data["steps"].([]interface{})

		for _, step := range steps {
			tariff := parseTariff(step.(map[string]interface{}))
			if tariff.Code == current.Code {
				tariff.Active = true
			}

			tariff.Product = product

			tariffs = append(tariffs, tariff)
		}

		//yota sliders data is structured as {product: {tariffsData}}
		//but all clients have only one product
		break
	}

	return tariffs, nil
}

func parseTariff(rawData map[string]interface{}) Tariff {
	tariff := Tariff{
		Name: rawData["name"].(string),
		Code: rawData["code"].(string),
	}

	amount, _ := strconv.ParseFloat(rawData["amountNumber"].(string), 64)
	tariff.Amount = amount

	speed := rawData["speedNumber"].(string)
	if strings.Contains(speed, "max-value") {
		speed = "max"
	}
	tariff.Speed = speed

	return tariff
}

func (cli *Client) getUid() (string, error) {
	payload := url.Values{
		"value": {cli.login},
	}

	resp, err := cli.http.PostForm(
		urlUidByMail,
		payload,
	)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	status := string(body[0:2])

	if status != "ok" {
		return "", errors.New(
			fmt.Sprintf("response is not ok: %s", string(body)))
	}

	//stupid parsing...
	return string(body[3:]), nil
}

func (cli *Client) GetBalance() (balance float64, currency string, err error) {
	resp, err := cli.http.Get(urlDevices)
	defer resp.Body.Close()

	if err != nil {
		return 0, "", err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	matches := reBalance.FindStringSubmatch(string(body))
	if len(matches) == 0 {
		return 0, "", errors.New("could not find balance data")
	}

	balance, err = strconv.ParseFloat(strings.Replace(matches[1], ",", ".", -1), 64)
	currency = matches[3]
	if err != nil {
		return 0, "", err
	}

	return balance, currency, nil
}

func (cli *Client) GetRemains() (string, error) {
	resp, err := cli.http.Get(urlDevices)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	matches := reRemains.FindStringSubmatch(string(body))
	if len(matches) == 0 {
		return "", errors.New("could not find remains data")
	}

	remains := fmt.Sprintf("%s %s %s", matches[1], matches[2], matches[3])

	return remains, nil
}

func (cli *Client) ChangeTariff(tariff Tariff) error {
	payload := url.Values{
		"product":       {tariff.Product},
		"offerCode":     {tariff.Code},
		"currentDevice": {"1"},
	}

	resp, err := cli.http.PostForm(
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
