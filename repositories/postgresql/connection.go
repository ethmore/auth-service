package postgresql

import (
	"database/sql"
	"fmt"

	"strconv"
	"sync"

	"auth-and-db-service/dotEnv"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *sql.DB
var db2 *gorm.DB
var err error
var err2 error

var lock = &sync.Mutex{}
var singleInstance *single

type single struct {
}

func Connect() *single {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			host := dotEnv.GoDotEnvVariable("HOST")
			port := dotEnv.GoDotEnvVariable("PORT")
			user := dotEnv.GoDotEnvVariable("USER")
			password := dotEnv.GoDotEnvVariable("PASSWORD")
			dbname := dotEnv.GoDotEnvVariable("DBNAME")
			intPort, _ := strconv.Atoi(port)
			psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, intPort, user, password, dbname)
			db2, err2 = gorm.Open(postgres.Open(psqlconn), &gorm.Config{})

			db, err = sql.Open("postgres", psqlconn)
			if err != nil {
				panic(err)
			}

			singleInstance = &single{}
			fmt.Println("PostgreSQL Connected!")
		} else {
			fmt.Println("PostgreSQL connection already created.")
		}
	} else {
		fmt.Println("PostgreSQL connection already created.")
	}
	return singleInstance
}
