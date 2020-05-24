package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	plainFormatter := new(PlainFormatter)
	plainFormatter.TimestampFormat = "2006/01/02 15:04:05"
	plainFormatter.LevelDesc = []string{"[panic]", "[fetal]", "[error]", "[warn]", "[info]", "[debug]", "[trace]"}
	log.SetFormatter(plainFormatter)
	log.SetLevel(logrus.DebugLevel)
}

type PlainFormatter struct {
	TimestampFormat string
	LevelDesc       []string
}

func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := fmt.Sprintf(entry.Time.Format(f.TimestampFormat))
	return []byte(fmt.Sprintf("%s %s %s\n", timestamp, f.LevelDesc[entry.Level], entry.Message)), nil
}

func Debug(format string, a ...interface{}) {
	log.Debugf(format, a...)
}

func Info(format string, a ...interface{}) {
	log.Infof(format, a...)
}

func Error(format string, a ...interface{}) {
	log.Errorf(format, a...)
}

func Fatal(format string, a ...interface{}) {
	log.Fatalf(format, a...)
}
