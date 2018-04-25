package main

import (
	"bytes"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"database/sql"
	"net/http"

	"github.com/json-iterator/go"
	"github.com/kshvakov/clickhouse"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func logger(logs chan *http.Request) {
	// Make DB connection string
	var connectParams = struct {
		host,
		port,
		debug,
		user,
		password,
		DB string
	}{
		"127.0.0.1",
		":9000",
		"?debug=false",
		"",
		"",
		"",
	}

	if v, present := os.LookupEnv("CH_HOST"); present {
		connectParams.host = v
	}

	if v, present := os.LookupEnv("CH_PORT"); present {
		connectParams.port = ":" + v
	}

	if v, present := os.LookupEnv("CH_DEBUG"); present && v == "true" {
		connectParams.debug = "?debug=true"
	}

	if v, present := os.LookupEnv("CH_USER"); present {
		connectParams.user = "&username=" + v
	}

	if v, present := os.LookupEnv("CH_PASSWORD"); present {
		connectParams.password = "&password=" + v
	}

	if v, present := os.LookupEnv("CH_DB"); present {
		connectParams.DB = "&database=" + v
	}

	// Prepare DB connection
	connect, err := sql.Open("clickhouse", "tcp://"+
		connectParams.host+
		connectParams.port+
		connectParams.debug+
		connectParams.user+
		connectParams.password+
		connectParams.DB+
		"")
	if err != nil {
		log.Fatal(err)
	}

	connect.SetConnMaxLifetime(4 * time.Minute) // ClickHouse will drop connection after 5 minues of inactive
	connect.SetMaxIdleConns(1)
	connect.SetMaxOpenConns(4)

	// Check DB connection
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			log.Println(err)
		}
		return
	}

	_, err = connect.Exec(`
		CREATE TABLE IF NOT EXISTS version_requests (
			date			Date DEFAULT today(),
			dt				DateTime DEFAULT now(),
			method			String,
			uri				String,
			proto			String,
			header			String,
			body			String,
			host			String,
			remote_host		String,
			remote_port		UInt16,
			buffer_used		UInt32,
			buffer_size		UInt32,
			response		String,
			code			UInt16
		) ENGINE = MergeTree(date, (dt, remote_host, remote_port), 8192)
	`)

	if err != nil {
		log.Fatal(err)
	}

	// Insert data from channel
	for {
		var (
			tx   *sql.Tx
			stmt *sql.Stmt
			err  error
		)

		// Wait (blocking) for first element in channel and iterate them
		for i := range logs {
			if tx == nil {
		// Prepare query
				tx, _ = connect.Begin()
				stmt, err = tx.Prepare(`INSERT INTO access (
				method,
				uri,
				proto,
				header, 
				body, 
				host,
				remote_host,
				remote_port,
				buffer_used,
				buffer_size,
				response,
				code
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
		if err != nil {
			log.Fatal(err)
		}
			}

			// Prepare values
			// Header as JSON
			headerJSON, _ := json.Marshal(&i.Header)

			// Body as string
			bodyBuf := new(bytes.Buffer)
			bodyBuf.ReadFrom(i.Body)
			bodyStr := bodyBuf.String()

			// Remote host/port
			remoteHost, remotePortStr, _ := net.SplitHostPort(i.RemoteAddr)
			var remotePort uint16
			if v, err := strconv.ParseUint(remotePortStr, 10, 16); err == nil {
				remotePort = uint16(v)
			}

			// Make query
			if _, err := stmt.Exec(
				i.Method,
				i.RequestURI,
				i.Proto,
				headerJSON,
				bodyStr,
				i.Host,
				remoteHost,
				remotePort,
				len(logs),
				bufferSize,
				response,
				200,
			); err != nil {
				log.Println(err)
			}

			// If no items left - commit query and whait for next batch
			if len(logs) == 0 {
				break
			}
		}

		if err := tx.Commit(); err != nil {
			log.Println(err)
		}
		stmt.Close()
	}
}
