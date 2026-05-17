package conf_test

import (
	"os"
	"source-score/pkg/conf"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	SamplePort  = "8099"
	SamplePwd   = "sample-pwd"
	SampleHost  = "sample-host"
	SampleSUPwd = "super-pwd"
	SampleDbName = "test-db"
)

var _ = Describe("Conf Tests", func() {
	When("dotenv path is not set", func() {
		os.Setenv("APP_USER_PASSWORD", SamplePwd)
		os.Setenv("PG_HOST", SampleHost)
		os.Setenv("PORT", SamplePort)
		os.Setenv("SUPER_USER_PASSWORD", SampleSUPwd)

		It("should load the environment variables into the config", func() {
			os.Unsetenv("DOTENV_PATH")
			conf.LoadConfig()

			Expect(conf.Cfg.AppUserPassword).To(BeEquivalentTo(SamplePwd))
			Expect(conf.Cfg.PgHost).To(BeEquivalentTo(SampleHost))
			Expect(conf.Cfg.Port).To(BeEquivalentTo(SamplePort))
			Expect(conf.Cfg.SuperUserPassword).To(BeEquivalentTo(SampleSUPwd))
		})
	})

	When("dotenv path is set", func() {
		It("should load the environment variables into the config", func() {
			os.Setenv("DOTENV_PATH", "./conf.yaml")
			conf.LoadConfig()

			Expect(conf.Cfg.AppUserPassword).To(BeEquivalentTo("env-pwd"))
			Expect(conf.Cfg.PgHost).To(BeEquivalentTo("env-host"))
			Expect(conf.Cfg.Port).To(BeEquivalentTo("8999"))
			Expect(conf.Cfg.SuperUserPassword).To(BeEquivalentTo("user-pwd"))
		})
	})

	When("DbName field is read", func() {
		It("should use default value when DB_NAME environment variable is not set", func() {
			os.Unsetenv("DOTENV_PATH")
			os.Unsetenv("DB_NAME")
			os.Setenv("APP_USER_PASSWORD", SamplePwd)
			os.Setenv("PG_HOST", SampleHost)
			os.Setenv("PORT", SamplePort)
			os.Setenv("SUPER_USER_PASSWORD", SampleSUPwd)

			conf.LoadConfig()

			Expect(conf.Cfg.DbName).To(BeEquivalentTo("sourcescore"))
		})

		It("should be overwritten when DB_NAME environment variable is set", func() {
			os.Unsetenv("DOTENV_PATH")
			os.Setenv("DB_NAME", SampleDbName)
			os.Setenv("APP_USER_PASSWORD", SamplePwd)
			os.Setenv("PG_HOST", SampleHost)
			os.Setenv("PORT", SamplePort)
			os.Setenv("SUPER_USER_PASSWORD", SampleSUPwd)

			conf.LoadConfig()

			Expect(conf.Cfg.DbName).To(BeEquivalentTo(SampleDbName))
		})
	})
})
