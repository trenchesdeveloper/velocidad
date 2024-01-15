package celeritas

import (
	"fmt"

	"github.com/joho/godotenv"
)

const Version = "1.0.0"

type Celeritas struct {
	AppName string
	Debug   bool
	Version string
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