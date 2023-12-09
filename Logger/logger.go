package logger

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/justinas/alice"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

var Zero zerolog.Logger
var Chain alice.Chain

func init() {
	pwd, _ := os.Executable()
	dir := filepath.Dir(pwd)

	appLogConsole := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	appLogFile := &lumberjack.Logger{
		Filename:  dir + "/logs/app.log",
		LocalTime: true,
		Compress:  true,
	}
	appLogFile.Rotate()
	appLog := zerolog.MultiLevelWriter(appLogFile, appLogConsole)
	Zero = zerolog.New(appLog).With().Timestamp().Logger()

	httpLogFile := &lumberjack.Logger{
		Filename:  dir + "/logs/http.log",
		LocalTime: true,
		Compress:  true,
	}
	httpLogFile.Rotate()
	httpLog := zerolog.MultiLevelWriter(httpLogFile)
	httpLogger := zerolog.New(httpLog).With().Timestamp().Logger()

	Chain = alice.New()
	Chain = Chain.Append(hlog.NewHandler(httpLogger))

	Chain = Chain.Append(
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Stringer("url", r.URL).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("")
		}),
	)
	Chain = Chain.Append(hlog.RemoteAddrHandler("ip"))
	Chain = Chain.Append(hlog.UserAgentHandler("user_agent"))
	Chain = Chain.Append(hlog.RefererHandler("referer"))
	Chain = Chain.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

}
