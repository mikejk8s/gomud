package postgrespkg

import (
    "os"
)

var (
    POSTGRESS_HOST = os.Getenv("POSTGRES_HOST")
    POSTGRES_USER = os.Getenv("POSTGRES_USER")
    POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
    POSTGRES_SSLMODE = os.Getenv("POSTGRES_SSLMODE")
    POSTGRES_DB1 = os.Getenv("POSTGRES_DB1")
    POSTGRES_DB2 = os.Getenv("POSTGRES_DB2")
    RUNNING_ON_DOCKER = os.Getenv("RUNNING_ON_DOCKER")
)