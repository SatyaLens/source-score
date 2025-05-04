package conf_test

import (
	"os"
	"source-score/pkg/conf"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	SamplePwd    = "sample-pwd"
	SampleServer = "sample-server"
)

var _ = Describe("Conf Tests", func() {
	When("dotenv path is not set", func() {
		os.Setenv("PG_USER_PASSWORD", SamplePwd)
		os.Setenv("PG_SERVER", SampleServer)

		It("should load the environment variables into the config", func() {
			os.Unsetenv("DOTENV_PATH")
			conf.LoadConfig()

			Expect(conf.Cfg.AppUserPassword).To(BeEquivalentTo(SamplePwd))
			Expect(conf.Cfg.PgServer).To(BeEquivalentTo(SampleServer))
		})
	})

	When("dotenv path is set", func() {
		It("should load the environment variables into the config", func() {
			os.Setenv("DOTENV_PATH", "./conf.yaml")
			conf.LoadConfig()

			Expect(conf.Cfg.AppUserPassword).To(BeEquivalentTo("env-pwd"))
			Expect(conf.Cfg.PgServer).To(BeEquivalentTo("env-server"))
		})
	})
})
