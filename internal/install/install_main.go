package install

import (
	"fmt"
	"os"

	"answer/internal/base/translator"
	"answer/internal/cli"
)

var (
	port     = os.Getenv("INSTALL_PORT")
	confPath = ""
)

func Run(configPath string) {
	confPath = configPath
	// initialize translator for return internationalization error when installing.
	_, err := translator.NewTranslator(&translator.I18n{BundleDir: cli.I18nPath})
	if err != nil {
		panic(err)
	}

	installServer := NewInstallHTTPServer()
	if len(port) == 0 {
		port = "80"
	}
	fmt.Printf("[SUCCESS] answer installation service will run at: http://localhost:%s/install/ \n", port)
	if err = installServer.Run(":" + port); err != nil {
		panic(err)
	}
}
