package nineone

import (
	"github.com/sirupsen/logrus"
	"os"
)

func NewLogrus(s string) *logrus.Logger{
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(os.Stdout)
	l.WithField("service", s)
	return l
}
