package config_test

import (
	"os"
	"testing"

	"github.com/kaogeek/line-fact-check/factcheck/internal/config"
)

func TestParseEnvConfig(t *testing.T) {
	setRequired := func(addr string, db string) func() {
		if addr != "" {
			err := os.Setenv("FACTCHECKAPI_LISTEN_ADDRESS", addr)
			if err != nil {
				panic(err)
			}
		}
		if db != "" {
			err := os.Setenv("POSTGRES_DB", db)
			if err != nil {
				panic(err)
			}
		}
		return func() {
			err := os.Unsetenv("FACTCHECKAPI_LISTEN_ADDRESS")
			if err != nil {
				panic(err)
			}
			err = os.Unsetenv("POSTGRES_DB")
			if err != nil {
				panic(err)
			}
		}
	}

	t.Run("error - missing required", func(t *testing.T) {
		conf, err := config.New()
		if err == nil {
			t.Fatal("unexpected nil error", conf)
		}
	})

	t.Run("error - missing required ListenAddress", func(t *testing.T) {
		defer setRequired("", "some_db")()
		conf, err := config.New()
		if err == nil {
			t.Fatal("unexpected nil error", conf)
		}
	})

	t.Run("error - missing required DBName", func(t *testing.T) {
		defer setRequired(":8080", "")()
		conf, err := config.New()
		if err == nil {
			t.Fatal("unexpected nil error", conf)
		}
	})

	t.Run("normal - set required only", func(t *testing.T) {
		addr := ":8888"
		os.Setenv("FACTCHECKAPI_LISTEN_ADDRESS", addr)
		os.Setenv("POSTGRES_DB", "some_db")
		defer func() {
			os.Setenv("FACTCHECKAPI_LISTEN_ADDRESS", addr)
			os.Unsetenv("POSTGRES_DB")
		}()
		conf, err := config.New()
		if err != nil {
			t.Fatal(err)
		}
		if conf.HTTP.ListenAddr != addr {
			t.Fatalf("unexpected listen address: %+v", conf.HTTP)
		}
		if conf.HTTP.TimeoutMsRead == 0 {
			t.Fatalf("unexpected 0 timeout read: %+v", conf.HTTP)
		}
		if conf.HTTP.TimeoutMsWrite == 0 {
			t.Fatalf("unexpected 0 timeout write: %+v", conf.HTTP)
		}
	})

	t.Run("normal", func(t *testing.T) {
		addr := ":8888"
		os.Setenv("FACTCHECKAPI_LISTEN_ADDRESS", addr)
		os.Setenv("FACTCHECKAPI_TIMEOUTMS_READ", "2000")
		os.Setenv("FACTCHECKAPI_TIMEOUTMS_WRITE", "1000")
		os.Setenv("POSTGRES_DB", "some_db")
		defer func() {
			os.Unsetenv("FACTCHECKAPI_LISTEN_ADDRESS")
			os.Unsetenv("FACTCHECKAPI_TIMEOUTMS_READ")
			os.Unsetenv("FACTCHECKAPI_TIMEOUTMS_WRITE")
			os.Unsetenv("POSTGRES_DB")
		}()
		conf, err := config.New()
		if err != nil {
			t.Fatal(err)
		}
		if conf.HTTP.ListenAddr != addr {
			t.Fatalf("unexpected listen address: %+v", conf.HTTP)
		}
		if conf.HTTP.TimeoutMsRead != 2000 {
			t.Fatalf("unexpected timeout read: %+v", conf.HTTP)
		}
		if conf.HTTP.TimeoutMsWrite != 1000 {
			t.Fatalf("unexpected timeout write: %+v", conf.HTTP)
		}
	})
}
