package config

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	var c AppConfig
	c.Load()
	fmt.Println(c.Db.DataSource)
}
