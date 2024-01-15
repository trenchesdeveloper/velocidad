package velocidad

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const Version = "1.0.0"

type Velocidad struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string
	config   config
}

type config struct {
	port     string
	renderer string // template engine

}

func (v *Velocidad) New(rootPath string) error {
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

	err := v.Init(pathConfig)

	if err != nil {
		return err
	}

	// check if dotenv file exists
	err = v.checkDotEnv(rootPath)

	if err != nil {
		return err
	}

	// read .env file
	err = godotenv.Load(rootPath + "/.env")

	if err != nil {
		return err
	}

	// create loggers
	infoLog, errorLog := v.StartLoggers()

	v.InfoLog = infoLog
	v.ErrorLog = errorLog
	v.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	v.Version = Version
	v.RootPath = rootPath

	v.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	return nil
}

func (v *Velocidad) Init(p initPaths) error {
	root := p.rootPath

	for _, folder := range p.folderNames {
		// Create folders if it doesn't exist
		err := v.CreateDirIfNotExist(root + "/" + folder)

		if err != nil {
			return err
		}
	}

	return nil
}

func (v *Velocidad) checkDotEnv(rootPath string) error {
	// check if dotenv file exists
	err := v.CreateFileIfNotExist(fmt.Sprintf("%s/.env", rootPath))

	if err != nil {
		return err
	}

	return nil
}

func (v *Velocidad) StartLoggers() (*log.Logger, *log.Logger) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}
