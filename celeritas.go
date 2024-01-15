package celeritas

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const Version = "1.0.0"

type Celeritas struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string
}

func (c *Celeritas) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath: rootPath,
		folderNames: []string{
			"handlers",
			"data",
			"migrations",
			"views",
			"public",
			"tmp",
			"logs",
			"middleware",
		},
	}

	err := c.Init(pathConfig)

	if err != nil {
		return err
	}

	// check if dotenv file exists
	err = c.checkDotEnv(rootPath)

	if err != nil {
		return err
	}

	// read .env file
	err = godotenv.Load(rootPath + "/.env")

	if err != nil {
		return err
	}

	// create loggers
	infoLog, errorLog := c.StartLoggers()

	c.InfoLog = infoLog
	c.ErrorLog = errorLog
	c.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	c.Version = Version

	return nil
}

func (c *Celeritas) Init(p initPaths) error {
	root := p.rootPath

	for _, folder := range p.folderNames {
		// Create folders if it doesn't exist
		err := c.CreateDirIfNotExist(root + "/" + folder)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Celeritas) checkDotEnv(rootPath string) error {
	// check if dotenv file exists
	err := c.CreateFileIfNotExist(fmt.Sprintf("%s/.env", rootPath))

	if err != nil {
		return err
	}

	return nil
}

func (c *Celeritas) StartLoggers() (*log.Logger, *log.Logger) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}
