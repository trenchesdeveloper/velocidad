package velocidad

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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
	Routes   *chi.Mux
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
	v.Routes = v.routes().(*chi.Mux)

	v.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	return nil
}

// ListenAndServe starts the HTTP server and listens for incoming requests.
// It creates an instance of http.Server with the specified address, handler,
// error log, idle timeout, read timeout, and write timeout.
// It then starts the server and logs the port it is listening on.
// If an error occurs while starting the server, it logs the error and exits the program.
func (v *Velocidad) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler:      v.routes(),
		ErrorLog:     v.ErrorLog,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	v.InfoLog.Printf("Starting server on port %s", os.Getenv("PORT"))

	err := srv.ListenAndServe()

	v.ErrorLog.Fatal(err)
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
