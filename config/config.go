package config

import "os"

// Debug указывает, в каком режиме работает приложение (true - debug, false - release)
var Debug = os.Getenv("DEBUG") == "true"
