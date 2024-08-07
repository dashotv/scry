package app

import (
	"os"

	"github.com/streamingfast/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/term"
)

func init() {
	initializers = append(initializers, setupLogger)
}

func setupLogger(app *Application) (err error) {
	switch app.Config.Logger {
	case "dev":
		isTTY := term.IsTerminal(int(os.Stderr.Fd()))
		verbosity := 1
		logStdoutWriter := zapcore.Lock(os.Stderr)
		log := zap.New(zapcore.NewCore(logging.NewEncoder(verbosity, isTTY), logStdoutWriter, zapcore.DebugLevel))
		app.Log = log.Sugar().Named("app")
	case "release":
		log, err := zap.NewProduction()
		app.Log = log.Sugar().Named("app")
		if err != nil {
			return err
		}
	}

	return nil
}
