package conf_test

import (
	"os"
	"source-score/pkg/conf"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	SamplePort   = "8099"
	SamplePwd    = "sample-pwd"
	SampleServer = "sample-server"
	SampleSUPwd  = "super-pwd"
)

var _ = Describe("Conf Tests", func() {
	When("dotenv path is not set", func() {
		os.Setenv("APP_USER_PASSWORD", SamplePwd)
		os.Setenv("PG_SERVER", SampleServer)
		os.Setenv("PORT", SamplePort)
		os.Setenv("SUPER_USER_PASSWORD", SampleSUPwd)

		It("should load the environment variables into the config", func() {
			os.Unsetenv("DOTENV_PATH")
			conf.LoadConfig()

			Expect(conf.Cfg.AppUserPassword).To(BeEquivalentTo(SamplePwd))
			Expect(conf.Cfg.PgServer).To(BeEquivalentTo(SampleServer))
			Expect(conf.Cfg.Port).To(BeEquivalentTo("8999"))
			Expect(conf.Cfg.SuperUserPassword).To(BeEquivalentTo("user-pwd"))
		})
	})

	When("dotenv path is set", func() {
		It("should load the environment variables into the config", func() {
			os.Setenv("DOTENV_PATH", "./conf.yaml")
			conf.LoadConfig()

			Expect(conf.Cfg.AppUserPassword).To(BeEquivalentTo("env-pwd"))
			Expect(conf.Cfg.PgServer).To(BeEquivalentTo("env-server"))
			Expect(conf.Cfg.Port).To(BeEquivalentTo("8999"))
			Expect(conf.Cfg.SuperUserPassword).To(BeEquivalentTo("user-pwd"))
		})
	})
})
