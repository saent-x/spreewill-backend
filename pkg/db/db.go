package db

import (
	"fmt"
	"log"
	"os"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Init() {
	err := mgm.SetDefaultConfig(nil, "spreewill_lab", options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@localhost:27018", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))))

	if err != nil {
		log.Fatalln(err)
	}
}
