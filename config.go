package main

import (
	"time"
	"fmt"
)

type configValues struct {
	homePage                  string
	sessionExpirationTime     time.Duration
	sessionStoreCleanInterval time.Duration
}

var config = InitializeConfig()

func InitializeConfig() configValues {
	fmt.Print("Initializing configs\n")

	values := configValues{}
	values.homePage = "/home"
	values.sessionExpirationTime = time.Hour
	values.sessionStoreCleanInterval = time.Hour

	return values
}
