package config

import "os"

var Debug = os.Getenv("DEBUG") == "true"
