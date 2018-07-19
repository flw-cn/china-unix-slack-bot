package plugin

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type BaseConfig struct {
	ID      string `json:"ID" yaml:"ID" flag:"||service ID"`
	Debug   bool   `json:"Debug" yaml:"Debug"`
	LogFile string `json:"LogFile" yaml:"LogFile"`
}

type Base struct {
	Service
	Logger      Logger
	config      BaseConfig
	debug       bool
	initialized bool
}

func NewBase(config BaseConfig) *Base {
	return &Base{
		config: config,
		Logger: log.New(ioutil.Discard, "", log.LstdFlags),
	}
}

func (b *Base) ID() string {
	return b.config.ID
}

func (b *Base) SetLogger(logger Logger) {
	b.Logger = logger
}

func (b *Base) SetLogOutput(w io.Writer) {
	b.Logger.SetOutput(w)
}

func (b *Base) SetDebug(debug bool) {
	b.debug = debug
}

func (b *Base) GetDebug() bool {
	return b.debug
}

func (b *Base) Debug(args ...interface{}) {
	if !b.debug {
		return
	}

	b.Logger.Print(args...)
}

func (b *Base) Debugf(format string, args ...interface{}) {
	if !b.debug {
		return
	}

	b.Logger.Printf(format, args...)
}

func (b *Base) Init() error {
	if b.config.LogFile == "" {
	} else if b.config.LogFile == "-" {
		b.Logger.SetOutput(os.Stderr)
	} else {
		file, err := os.Create(b.config.LogFile)
		if err == nil {
			b.Logger.SetOutput(file)
		} else {
			b.Logger.SetOutput(os.Stderr)
			b.Logger.Printf("Can't open log file. Error: %v", err)
			b.Logger.SetOutput(ioutil.Discard)
		}
	}

	b.initialized = true
	return nil
}

func (b *Base) Start() error {
	if !b.initialized {
		return errors.New("not initialized yet")
	}

	return nil
}

func (b *Base) Stop() {
}
