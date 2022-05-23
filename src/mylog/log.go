package mylog

import (
	"log"
	"os"
	"path"
)

var (
	logDirectory     = "/var/log/govwa/"
	warnLogFileName  = "govwa.warning.log"
	infoLogFileName  = "govwa.info.log"
	buildLogFileName = "govwa.build.log"
	debugLogFileName = "govwa.debug.log"
	errorLogFileName = "govwa.error.log"
	Warn             log.Logger
	Info             log.Logger
	Build            log.Logger
	Debug            log.Logger
	Error            log.Logger
)

// LogInit funtion initializes logs to file
func Init() {
	// WARN
	warnLogFile, err := os.OpenFile(
		path.Join(logDirectory, warnLogFileName),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}
	//defer warnLogFile.Close()
	Warn = *log.New(warnLogFile, "[ warn   ]", log.LstdFlags)

	// INFO
	infoLogFile, err := os.OpenFile(
		path.Join(logDirectory, infoLogFileName),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}
	//defer infoLogFile.Close()
	Info = *log.New(infoLogFile, "[ info   ]", log.LstdFlags)

	// BUILD
	buildLogFile, err := os.OpenFile(
		path.Join(logDirectory, buildLogFileName),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}
	//defer buildLogFile.Close()
	Build = *log.New(buildLogFile, "[ build  ]", log.LstdFlags)

	// DEBUG
	debugLogFile, err := os.OpenFile(
		path.Join(logDirectory, debugLogFileName),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}
	//defer debugLogFile.Close()
	Debug = *log.New(debugLogFile, "[ debug  ]", log.LstdFlags)

	// ERROR
	errorLogFile, err := os.OpenFile(
		path.Join(logDirectory, errorLogFileName),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}
	//defer errorLogFile.Close()
	Error = *log.New(errorLogFile, "[ error  ]", log.LstdFlags)
}
