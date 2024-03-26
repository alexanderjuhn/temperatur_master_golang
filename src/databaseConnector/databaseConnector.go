package databaseConnector

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	viper "github.com/spf13/viper"
)

var (
	host             string
	port             int
	user             string
	password         string
	dbname           string
	RabbitConnection rabbitConnectionParam
)

type rabbitConnectionParam struct {
	HostRabbit     string
	PortRabbit     string
	UserRabbit     string
	PasswordRabbit string
	Queue          string
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/appname/")
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath("/config/")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	log.Println("Reading config at " + viper.GetViper().ConfigFileUsed())
	// Get PostgreSQL vars
	host = viper.GetString("database.hostname")
	port = viper.GetInt("database.port")
	user = viper.GetString("database.user")
	password = viper.GetString("database.password")
	dbname = viper.GetString("database.dbname")

	// Get RabbitMQ vars
	var rabbitConnection rabbitConnectionParam
	rabbitConnection.HostRabbit = viper.GetString("rabbitmq.hostname")
	rabbitConnection.PortRabbit = viper.GetString("rabbitmq.port")
	rabbitConnection.UserRabbit = viper.GetString("rabbitmq.user")
	rabbitConnection.PasswordRabbit = viper.GetString("rabbitmq.password")
	rabbitConnection.Queue = viper.GetString("rabbitmq.queue")
	RabbitConnection = rabbitConnection

	// Enable debug
	debug = viper.GetBool("debug")
	log.Println("Done reading config")
}

func GetDatabaseConnection() *sql.DB {
	if db != nil {
		return db
	}
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	return db
}
