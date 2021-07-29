package yota

import (
	"log"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
	"github.com/zazab/zhash"
)

var (
	client *Client
)

func TestClient_Login(t *testing.T) {
	t.Log("Testing Login")
	err := client.Login()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestClient_GetTariffs(t *testing.T) {
	tariffs, err := client.GetTariffs()
	if err != nil {
		t.Errorf("Can't get tariffs: %s", err.Error())
	}
	for _, tariff := range tariffs {
		t.Logf("%+v", tariff)
	}
}

func TestClient_GetBalance(t *testing.T) {
	balance, err := client.GetBalance()
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("balance: %.2f, currency: %s", balance.Amount, balance.CurrencyCode)
}

func TestClient_GetRemains(t *testing.T) {
	remains, err := client.GetCurrentInfo()
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%+v", remains)
}

func TestClient_GetPayments(t *testing.T) {
	payments, err := client.GetPayments()
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%+v", payments)
}

func TestClient_GetOperationHistory(t *testing.T) {
	oh, err := client.GetOperationHistory()
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%+v", oh)
}

func getConfig() (zhash.Hash, error) {
	var configData map[string]interface{}
	home, err := homedir.Dir()
	var path = home + "/.config/yotarc"

	_, err = toml.DecodeFile(path, &configData)
	if err != nil {
		return zhash.Hash{}, err
	}

	return zhash.HashFromMap(configData), nil
}

func setup() {
	log.Println("setup")
	config, err := getConfig()
	if err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	username, err := config.GetString("username")
	if err != nil {
		log.Fatal(err)
	}

	password, err := config.GetString("password")
	if err != nil {
		log.Fatal(err)
	}
	client = NewClient(username, password, nil)
	err = client.Login()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func shutdown() {
	log.Println("shutdown")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}
