package yota

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	username     string
	password     string
	accessToken  string
	refreshToken string
	jwt          string
	httpClient   *http.Client
}

type Tariff struct {
	Name      string
	Amount    int
	Code      string
	Speed     string
	SpeedType string
	Active    bool
}

const (
	hUserAgent       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36"
	hAccept          = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3"
	hAcceptLang      = "ru-RU,en-US,en;q=0.9"
	hContentTypeJson = "application/json"
	hContentTypeForm = "application/x-www-form-urlencoded"

	urlAccessToken  = "https://id.yota.ru/sso/oauth2/access_token?skipAutoLogin=true"
	urlTokenInfo    = "https://my.yota.ru/wa/v1/auth/tokenInfo"
	urlLoginSuccess = "https://my.yota.ru/wa/v1/auth/loginSuccess"
	urlDevices      = "https://my.yota.ru/wa/v1/devices/devices"
	urlStatusLegal  = "https://my.yota.ru/wa/v1/profile/statusLegal"
	urlGetBalance   = "https://my.yota.ru/wa/v1/finance/getBalance"
	urlInfo         = "https://my.yota.ru/wa/v1/profile/info"
	urlChangeTariff = "https://my.yota.ru/wa/v1/devices/changeOffer/change"
	urlPayments     = "https://my.yota.ru/wa/v1/finance/future/payments"
	urlOpHistory    = "https://my.yota.ru/wa/v1/finance/getOperationHistory"

	basicAuthStr = "Basic bmV3X2xrX3Jlc3Q6cGFzc3dvcmQ="
	dtLayout     = "2006-01-02T15:04:05.000Z"
)

func NewClient(login, password string, httpClient *http.Client) *Client {
	cookies, _ := cookiejar.New(nil)

	if httpClient == nil {
		httpClient = &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 20,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			},
			Jar:     cookies,
			Timeout: 1 * time.Minute,
		}
	}

	cli := &Client{
		username:   login,
		password:   password,
		httpClient: httpClient,
	}

	return cli
}

func (cli *Client) getDevices() DevicesInfo {
	req, err := http.NewRequest(http.MethodGet, urlDevices, nil)
	req.Header.Add("User-Agent", hUserAgent)
	req.Header.Set("Accept", hContentTypeJson)
	req.Header.Set("Authorization", basicAuthStr)
	res, err := cli.httpClient.Do(req)
	if err != nil {
		log.Printf("req do err: %s", err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(body), &jsonData)

	var devicesInfo DevicesInfo
	err = json.Unmarshal(body, &devicesInfo)
	if err != nil {
		log.Println(err)
	}

	return devicesInfo
}

func (cli *Client) getTokenInfo() (err error) {
	req, err := http.NewRequest(http.MethodPost, urlTokenInfo, nil)
	req.Header.Add("User-Agent", hUserAgent)
	req.Header.Set("Accept", hContentTypeJson)
	res, err := cli.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	return
}

func (cli *Client) getExecution() (execution string, err error) {
	data := url.Values{
		"client_id":     {"yota_mya"},
		"client_secret": {"password"},
		"realm":         {"/customer"},
		"service":       {"dispatcher"},
		"grant_type":    {"urn:roox:params:oauth:grant-type:m2m"},
		"response_type": {"token cookie"},
	}

	req, err := http.NewRequest(http.MethodPost, urlAccessToken, strings.NewReader(data.Encode()))
	req.Header.Add("User-Agent", hUserAgent)
	req.Header.Set("Accept", hContentTypeJson)
	req.Header.Set("Content-Type", hContentTypeForm)
	res, err := cli.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(body), &jsonData)
	execution = jsonData["execution"].(string)
	return
}

func (cli *Client) getAccessToken(execution string) (err error) {
	data := url.Values{
		"execution":     {execution},
		"username":      {cli.username},
		"password":      {cli.password},
		"_eventId":      {"next"},
		"response_type": {"token cookie"},
		"client_id":     {"yota_mya"},
		"client_secret": {"password"},
		"service":       {"dispatcher"},
		"grant_type":    {"urn:roox:params:oauth:grant-type:m2m"},
		"realm":         {"/customer"},
	}

	req, err := http.NewRequest(http.MethodPost, urlAccessToken, strings.NewReader(data.Encode()))
	req.Header.Add("User-Agent", hUserAgent)
	req.Header.Set("Accept", hContentTypeJson)
	req.Header.Set("Content-Type", hContentTypeForm)
	res, err := cli.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(body), &jsonData)
	cli.accessToken = jsonData["access_token"].(string)
	cli.refreshToken = jsonData["refresh_token"].(string)
	cli.jwt = jsonData["JWTToken"].(string)

	return
}

func (cli *Client) getLoginSuccess() (err error) {
	req, err := http.NewRequest(http.MethodGet, urlLoginSuccess, nil)
	req.Header.Add("User-Agent", hUserAgent)
	req.Header.Set("Accept", hContentTypeJson)
	req.Header.Set("Authorization", basicAuthStr)
	res, err := cli.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(body), &jsonData)
	redirect := jsonData["redirect"].(string)
	if redirect == "/devices" {
		return
	}
	return errors.New("ErrLoginFail")
}

func (cli *Client) getStatusLegal() (err error) {
	req, err := http.NewRequest(http.MethodGet, urlStatusLegal, nil)
	req.Header.Add("User-Agent", hUserAgent)
	req.Header.Set("Accept", hContentTypeJson)
	req.Header.Set("Authorization", basicAuthStr)
	res, err := cli.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(body), &jsonData)
	legalize := jsonData["legalize"].(bool)
	if legalize {
		return
	}
	return errors.New("ErrNotLegalized")
}

func (cli *Client) Login() (err error) {
	err = cli.getTokenInfo()
	if err != nil {
		return
	}
	execution, err := cli.getExecution()
	if err != nil {
		return
	}
	err = cli.getAccessToken(execution)
	if err != nil {
		return
	}
	err = cli.getLoginSuccess()
	if err != nil {
		return
	}
	err = cli.getStatusLegal()
	if err != nil {
		return
	}

	err = cli.getTokenInfo()
	if err != nil {
		return
	}

	return nil
}

func (cli *Client) GetTariffs() (tariffs []Tariff, err error) {
	devicesInfo := cli.getDevices()
	current := devicesInfo.Devices[0].Slider.CurrentProduct
	for _, s := range devicesInfo.Devices[0].Slider.Steps {
		var t Tariff
		t.Code = s.Code
		if s.Code == current.Code {
			t.Active = true
		}
		if strings.Contains(s.Speed, "maxvalue") {
			t.Speed = "max"
		} else {
			t.Speed = s.Speed
		}
		t.SpeedType = s.SpeedType
		t.Amount = s.Amount
		t.Name = s.OfferDescription

		tariffs = append(tariffs, t)
	}
	return
}

func (cli *Client) GetBalance() (balance Balance, err error) {
	req, err := http.NewRequest(http.MethodGet, urlGetBalance, nil)
	req.Header.Add("User-Agent", hUserAgent)
	req.Header.Set("Accept", hContentTypeJson)
	req.Header.Set("Authorization", basicAuthStr)
	res, err := cli.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &balance)
	if err != nil {
		return
	}
	return
}

func (cli *Client) GetUserInfo() (ui UserInfo, err error) {
	req, err := http.NewRequest(http.MethodGet, urlInfo, nil)
	req.Header.Add("User-Agent", hUserAgent)
	req.Header.Set("Accept", hContentTypeJson)
	req.Header.Set("Authorization", basicAuthStr)
	res, err := cli.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &ui)
	if err != nil {
		return
	}
	return
}

func (cli *Client) GetCurrentInfo() (ci CurrentInfo, err error) {
	devicesInfo := cli.getDevices()
	d := devicesInfo.Devices[0]
	ci.BeginDate = d.Product.BeginDate
	ci.EndDate = d.Product.EndDate
	ci.Price.Amount = d.Product.Price.Amount
	ci.Price.CurrencyCode = d.Product.Price.CurrencyCode
	ci.Speed.SpeedValue = d.OfferingSpeed.SpeedValue
	ci.Speed.UnitOfMeasure = d.OfferingSpeed.UnitOfMeasure
	return
}

func (cli *Client) ChangeOfferTo(code string) (err error) {
	devicesInfo := cli.getDevices()
	payload := new(ChangeTariff)
	payload.CurrentProductID = devicesInfo.Devices[0].Product.ProductID
	payload.DisablingAutoprolong = false
	payload.OfferCode = code
	payload.ResourceID.Key = devicesInfo.Devices[0].PhysicalResource.ResourceID.Key
	payload.ResourceID.Type = devicesInfo.Devices[0].PhysicalResource.ResourceID.Type
	b, err := json.Marshal(payload)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, urlChangeTariff, bytes.NewReader(b))
	req.Header.Add("User-Agent", hUserAgent)
	req.Header.Set("Accept", hContentTypeJson)
	req.Header.Set("Authorization", basicAuthStr)
	req.Header.Set("Content-Type", hContentTypeJson)
	res, err := cli.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(body), &jsonData)
	message := jsonData["message"].(string)
	if message == "OK" {
		return
	} else {
		return errors.New("ErrChangeOffer code: " + jsonData["code"].(string))
	}
}

func (cli *Client) GetPayments() (pi PaymentsInfo, err error) {
	req, err := http.NewRequest(http.MethodGet, urlPayments, nil)
	req.Header.Add("User-Agent", hUserAgent)
	req.Header.Set("Accept", hContentTypeJson)
	req.Header.Set("Authorization", basicAuthStr)

	now := time.Now()
	endTs := now.AddDate(0, 1, 0)
	q := req.URL.Query()
	q.Add("startTs", now.Format(dtLayout))
	q.Add("endTs", endTs.Format(dtLayout))
	req.URL.RawQuery = q.Encode()

	res, err := cli.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	if err = json.Unmarshal(body, &pi); err != nil {
		return
	}
	return
}

func (cli *Client) GetOperationHistory() (oh []OperationHistory, err error) {
	req, err := http.NewRequest(http.MethodGet, urlOpHistory, nil)
	req.Header.Add("User-Agent", hUserAgent)
	req.Header.Set("Accept", hContentTypeJson)
	req.Header.Set("Authorization", basicAuthStr)

	now := time.Now()
	startTs := now.AddDate(0, -6, 0)
	q := req.URL.Query()
	q.Add("fromDate", startTs.Format(dtLayout))
	q.Add("toDate", now.Format(dtLayout))
	req.URL.RawQuery = q.Encode()

	res, err := cli.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	if err = json.Unmarshal(body, &oh); err != nil {
		return
	}
	return
}
