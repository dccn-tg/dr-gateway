package dr

import (
	"os"
	"testing"
)

func TestGetAllCollections(t *testing.T) {

	initLogger()

	pathCfg := os.Getenv("DR_GATEWAY_CONFIG")

	cfg, err := loadDRConfig(pathCfg)
	if err != nil {
		t.Fatal(err)
	}

	colls, err := GetAllCollections(cfg)
	if err != nil {
		t.Fatal(err)
	}

	for c := range colls {
		t.Logf("user: %+v\n", c)
	}
}
