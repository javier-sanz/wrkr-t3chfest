package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aodin/date"
	yaml "gopkg.in/yaml.v2"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/kshvakov/clickhouse"
)

type Config struct {

	// FlushFrequency DB commit frequency
	FlushFrequency time.Duration `yaml:"flushFrequency"`

	// DatabaseURL database url where
	DatabaseURL string `yaml:"databaseURL"`

	// DataFolder folder where data can be found
	DataFolder string `yaml:"dataFolder"`
}

func loadConfig(configFile string) (Config, error) {
	res := Config{}
	source, err := ioutil.ReadFile(configFile)
	if err != nil {
		return res, err
	}
	err = yaml.Unmarshal(source, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func createQuestionMarks(len int) string {
	marks := make([]string, len)
	for i := 0; i < len; i++ {
		marks[i] = "?"
	}
	return strings.Join(marks, ", ")
}

func insertDateData(curDate date.Date, folder string, db *sql.DB, tableName string) error {
	log.Println("Loading day: " + curDate.String())
	fileName := folder + "/" + curDate.String() + ".csv"
	csvFile, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	line, err := reader.Read()
	if err == io.EOF {
		log.Println("Empty file")
		return nil
	} else if err != nil {
		panic(err)
	}
	fields := strings.Join(line, ", ")
	marks := createQuestionMarks(len(line))
	var transaction *sql.Tx
	var statement *sql.Stmt
	for i := 0; ; i++ {
		if i%5000 == 0 {
			if transaction != nil {
				log.Println("Commiting insertion transaction")
				transaction.Commit()
			}
			transaction, err = db.Begin()
			if err != nil {
				panic(err)
				return err
			}
			statement, err = transaction.Prepare("INSERT INTO " + tableName +
				" (" + fields + ") VALUES (" + marks + ")")
			if err != nil {
				panic(err)
			}
		}
		line, err = reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
			continue
		}

		_, err = statement.Exec(createSQLParams(line)...)

		if err != nil {
			panic(err)
			return err

		}

	}
	return nil
}

func newConnectionClickHouseConnection(databaseURL string) *sql.DB {
	dbConn, err := sql.Open("clickhouse", databaseURL)
	if err != nil {
		panic("Could not connect with ClickHouse db: " + err.Error())
	}
	return dbConn
}

func createSQLParams(values []string) []interface{} {
	interfaceValues := make([]interface{}, len(values))
	if myDate, err := time.Parse("2006-01-02", values[0]); err == nil {
		interfaceValues[0] = myDate
	} else {
		panic("Error parsin date: " + err.Error())
	}
	interfaceValues[1] = values[1]
	interfaceValues[2] = values[2]
	if capacityBytes, err := (strconv.ParseUint(values[3], 10, 64)); err == nil {
		interfaceValues[3] = capacityBytes
	} else {
		if values[3] == "-1" {
			interfaceValues[3] = 0
		} else {
			panic("Error parsing capacity_bytes: " + err.Error())
		}
	}
	if failure, err := (strconv.ParseUint(values[4], 10, 8)); err == nil {
		interfaceValues[4] = uint8(failure)
	} else {
		panic("Error parsing failure: " + err.Error())
	}

	for i := 5; i < len(values); i++ {
		if capacityBytes, err := (strconv.ParseUint(values[i], 10, 64)); err == nil {
			interfaceValues[i] = capacityBytes
		} else {
			interfaceValues[i] = 0
		}
	}
	return interfaceValues

}

func createClickhouseDB(db *sql.DB) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS default.drive_stats
		(
			date Date,
			serial_number String,
			model String,
			capacity_bytes UInt64,
			failure UInt8,
			smart_1_normalized UInt64,
			smart_1_raw UInt64,
			smart_2_normalized UInt64,
			smart_2_raw UInt64,
			smart_3_normalized UInt64,
			smart_3_raw UInt64,
			smart_4_normalized UInt64,
			smart_4_raw UInt64,
			smart_5_normalized UInt64,
			smart_5_raw UInt64,
			smart_7_normalized UInt64,
			smart_7_raw UInt64,
			smart_8_normalized UInt64,
			smart_8_raw UInt64,
			smart_9_normalized UInt64,
			smart_9_raw UInt64,
			smart_10_normalized UInt64,
			smart_10_raw UInt64,
			smart_11_normalized UInt64,
			smart_11_raw UInt64,
			smart_12_normalized UInt64,
			smart_12_raw UInt64,
			smart_13_normalized UInt64,
			smart_13_raw UInt64,
			smart_15_normalized UInt64,
			smart_15_raw UInt64,
			smart_22_normalized UInt64,
			smart_22_raw UInt64,   
			smart_183_normalized UInt64,
			smart_183_raw UInt64,
			smart_184_normalized UInt64,
			smart_184_raw UInt64,
			smart_187_normalized UInt64,
			smart_187_raw UInt64,
			smart_188_normalized UInt64,
			smart_188_raw UInt64,
			smart_189_normalized UInt64,
			smart_189_raw UInt64,
			smart_190_normalized UInt64,
			smart_190_raw UInt64,
			smart_191_normalized UInt64,
			smart_191_raw UInt64,
			smart_192_normalized UInt64,
			smart_192_raw UInt64,
			smart_193_normalized UInt64,
			smart_193_raw UInt64,
			smart_194_normalized UInt64,
			smart_194_raw UInt64,
			smart_195_normalized UInt64,
			smart_195_raw UInt64,
			smart_196_normalized UInt64,
			smart_196_raw UInt64,
			smart_197_normalized UInt64,
			smart_197_raw UInt64,
			smart_198_normalized UInt64,
			smart_198_raw UInt64,
			smart_199_normalized UInt64,
			smart_199_raw UInt64,
			smart_200_normalized UInt64,
			smart_200_raw UInt64,
			smart_201_normalized UInt64,
			smart_201_raw UInt64,
			smart_220_normalized UInt64,
			smart_220_raw UInt64,
			smart_222_normalized UInt64,
			smart_222_raw UInt64,
			smart_223_normalized UInt64,
			smart_223_raw UInt64,
			smart_224_normalized UInt64,
			smart_224_raw UInt64,   
			smart_225_normalized UInt64,
			smart_225_raw UInt64,
			smart_226_normalized UInt64,
			smart_226_raw UInt64,  
			smart_240_normalized UInt64,
			smart_240_raw UInt64,
			smart_241_normalized UInt64,
			smart_241_raw UInt64,
			smart_242_normalized UInt64,
			smart_242_raw UInt64,
			smart_250_normalized UInt64,
			smart_250_raw UInt64,
			smart_251_normalized UInt64,
			smart_251_raw UInt64,
			smart_252_normalized UInt64,
			smart_252_raw UInt64,
			smart_254_normalized UInt64,
			smart_254_raw UInt64,
			smart_255_normalized UInt64,
			smart_255_raw UInt64
		)
		ENGINE = MergeTree(date, (serial_number, date), 8192);`)
	if err != nil {
		panic("Could not create clickhouse drives table: " + err.Error())
	}
}

func newConnectionMariaDB(databaseURL string) *sql.DB {
	dbConn, err := sql.Open("mysql", databaseURL)
	if err != nil {
		panic("Could not connect with MariaDB db: " + err.Error())
	}
	return dbConn
}

func createMariaDB(db *sql.DB) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS drive_stats (
			date TEXT NOT NULL,
			serial_number TEXT NOT NULL,
			model TEXT NOT NULL,
			capacity_bytes BIGINT NOT NULL,
			failure BIGINT NOT NULL,
			smart_1_normalized BIGINT,
			smart_1_raw BIGINT,
			smart_2_normalized BIGINT,
			smart_2_raw BIGINT,
			smart_3_normalized BIGINT,
			smart_3_raw BIGINT,
			smart_4_normalized BIGINT,
			smart_4_raw BIGINT,
			smart_5_normalized BIGINT,
			smart_5_raw BIGINT,
			smart_7_normalized BIGINT,
			smart_7_raw BIGINT,
			smart_8_normalized BIGINT,
			smart_8_raw BIGINT,
			smart_9_normalized BIGINT,
			smart_9_raw BIGINT,
			smart_10_normalized BIGINT,
			smart_10_raw BIGINT,
			smart_11_normalized BIGINT,
			smart_11_raw BIGINT,
			smart_12_normalized BIGINT,
			smart_12_raw BIGINT,
			smart_13_normalized BIGINT,
			smart_13_raw BIGINT,
			smart_15_normalized BIGINT,
			smart_15_raw BIGINT,
			smart_22_normalized BIGINT,
			smart_22_raw BIGINT,   
			smart_183_normalized BIGINT,
			smart_183_raw BIGINT,
			smart_184_normalized BIGINT,
			smart_184_raw BIGINT,
			smart_187_normalized BIGINT,
			smart_187_raw BIGINT,
			smart_188_normalized BIGINT,
			smart_188_raw BIGINT,
			smart_189_normalized BIGINT,
			smart_189_raw BIGINT,
			smart_190_normalized BIGINT,
			smart_190_raw BIGINT,
			smart_191_normalized BIGINT,
			smart_191_raw BIGINT,
			smart_192_normalized BIGINT,
			smart_192_raw BIGINT,
			smart_193_normalized BIGINT,
			smart_193_raw BIGINT,
			smart_194_normalized BIGINT,
			smart_194_raw BIGINT,
			smart_195_normalized BIGINT,
			smart_195_raw BIGINT,
			smart_196_normalized BIGINT,
			smart_196_raw BIGINT,
			smart_197_normalized BIGINT,
			smart_197_raw BIGINT,
			smart_198_normalized BIGINT,
			smart_198_raw BIGINT,
			smart_199_normalized BIGINT,
			smart_199_raw BIGINT,
			smart_200_normalized BIGINT,
			smart_200_raw BIGINT,
			smart_201_normalized BIGINT,
			smart_201_raw BIGINT,
			smart_220_normalized BIGINT,
			smart_220_raw BIGINT,
			smart_222_normalized BIGINT,
			smart_222_raw BIGINT,
			smart_223_normalized BIGINT,
			smart_223_raw BIGINT,
			smart_224_normalized BIGINT,
			smart_224_raw BIGINT,   
			smart_225_normalized BIGINT,
			smart_225_raw BIGINT,
			smart_226_normalized BIGINT,
			smart_226_raw BIGINT,  
			smart_240_normalized BIGINT,
			smart_240_raw BIGINT,
			smart_241_normalized BIGINT,
			smart_241_raw BIGINT,
			smart_242_normalized BIGINT,
			smart_242_raw BIGINT,
			smart_250_normalized BIGINT,
			smart_250_raw BIGINT,
			smart_251_normalized BIGINT,
			smart_251_raw BIGINT,
			smart_252_normalized BIGINT,
			smart_252_raw BIGINT,
			smart_254_normalized BIGINT,
			smart_254_raw BIGINT,
			smart_255_normalized BIGINT,
			smart_255_raw BIGINT
			);`)
	if err != nil {
		panic("Could not create MariaDB drives table: " + err.Error())
	}
}

func main() {

	var clickhouse bool
	if len(os.Args) > 1 && os.Args[1] == "clickhouse" {
		clickhouse = true
	}

	configFile := "configMariaDB.yaml"
	if clickhouse {
		configFile = "configClickhouse.yaml"
	}
	config, err := loadConfig(configFile)

	if err != nil {
		panic("Config file not found")
	}

	var db *sql.DB
	var tableName string
	if clickhouse {
		db = newConnectionClickHouseConnection(config.DatabaseURL)
		createClickhouseDB(db)
		tableName = "default.drive_stats"
	} else {
		db = newConnectionMariaDB(config.DatabaseURL)
		createMariaDB(db)
		tableName = "drive_stats"

	}
	defer db.Close()

	for curDate := date.New(2017, 1, 1); date.New(2017, 9, 30).After(curDate); curDate = curDate.AddDays(1) {
		err := insertDateData(curDate, config.DataFolder, db, tableName)
		if err != nil {
			fmt.Println("Error loading date: " + curDate.String() + ". Error: " + err.Error())
		}
	}

}
