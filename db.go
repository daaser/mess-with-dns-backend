package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/miekg/dns"
)

// connect to planetscale
func connect() (*sql.DB, error) {
	// get connection string from environment
	connStr := os.Getenv("DSN")
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTables(db *sql.DB) error {
	if os.Getenv("DEV") == "true" {
		fmt.Println("creating tables...")
		err := loadSQLFile(db, "create.sql")
		if err != nil {
			return err
		}
		// initialize the serials table
		// check if serials table has anything in it
		rows, err := db.Query("SELECT * FROM dns_serials")
		if err != nil {
			return err
		}
		if rows.Next() {
			// if it has something in it, we don't need to do anything
			return nil
		}
		_, err = db.Exec("INSERT INTO dns_serials (serial) VALUES (10)")
		if err != nil {
			return err
		}
	}
	return nil
}

func loadSQLFile(db *sql.DB, sqlFile string) error {
	file, err := os.ReadFile(sqlFile)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	for _, q := range strings.Split(string(file), ";") {
		q := strings.TrimSpace(q)
		if q == "" {
			continue
		}
		if _, err := tx.Exec(q); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func GetSerial(db *sql.DB) (uint32, error) {
	var serial uint32
	err := db.QueryRow("SELECT serial FROM dns_serials").Scan(&serial)
	if err != nil {
		return 0, err
	}
	return serial, nil
}

func IncrementSerial(tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE dns_serials SET serial = serial + 1")
	if err != nil {
		return err
	}
	// get new serial
	var serial uint32
	err = tx.QueryRow("SELECT serial FROM dns_serials").Scan(&serial)
	if err != nil {
		return err
	}
	// commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	soaSerial = serial
	return nil
}

func DeleteRecord(db *sql.DB, id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM dns_records WHERE id = ?", id)
	if err != nil {
		return err
	}
	return IncrementSerial(tx)
}

func DeleteOldRecords(db *sql.DB) {
	// delete records where created_at timestamp is more than a week old
	_, err := db.Exec("DELETE FROM dns_records WHERE created_at < NOW() - INTERVAL 1 DAY")
	if err != nil {
		panic(err)
	}
}

func DeleteOldRequests(db *sql.DB) {
	// delete requests where created_at timestamp is more than a day
	// if we don't put the limit I get a "resources exhausted" error
	// 1 day ago, postgres
	_, err := db.Exec("DELETE FROM dns_requests WHERE created_at < NOW() - INTERVAL 1 DAY")
	if err != nil {
		panic(err)
	}
}

func UpdateRecord(db *sql.DB, id int, record dns.RR) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	jsonString, err := json.Marshal(record)
	if err != nil {
		return err
	}
	name := record.Header().Name
	_, err = tx.Exec(
		"UPDATE dns_records SET name = ?, subdomain = ?, rrtype = ?, content = ? WHERE id = ?",
		name,
		ExtractSubdomain(name),
		record.Header().Rrtype,
		jsonString,
		id,
	)
	if err != nil {
		return err
	}
	return IncrementSerial(tx)
}

func InsertRecord(db *sql.DB, record dns.RR) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	jsonString, err := json.Marshal(record)
	if err != nil {
		return err
	}
	name := record.Header().Name
	_, err = tx.Exec(
		"INSERT INTO dns_records (name, subdomain, rrtype, content) VALUES (?, ?, ?, ?)",
		name,
		ExtractSubdomain(name),
		record.Header().Rrtype,
		jsonString,
	)
	if err != nil {
		return err
	}
	return IncrementSerial(tx)
}

func uncommittedTransaction(db *sql.DB) (*sql.Tx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED")
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func GetRecordsForName(db *sql.DB, subdomain string) (map[int]dns.RR, error) {
	// we're stricter about the isolation level here because it's weird if you delete
	// a record, but it still exists after
	rows, err := db.Query("SELECT id, content FROM dns_records WHERE subdomain = ?", subdomain)
	if err != nil {
		return nil, err
	}
	records := make(map[int]dns.RR)
	for rows.Next() {
		var content []byte
		var id int
		err = rows.Scan(&id, &content)
		if err != nil {
			return nil, err
		}
		record, err := ParseRecord(content)
		if err != nil {
			return nil, err
		}
		records[id] = record
	}
	return records, nil
}

func LogRequest(
	db *sql.DB,
	request *dns.Msg,
	response *dns.Msg,
	src_ip net.IP,
	src_host string,
) error {
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return err
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}
	name := request.Question[0].Name
	subdomain := ExtractSubdomain(name)
	err = StreamRequest(subdomain, jsonRequest, jsonResponse, src_ip.String(), src_host)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO dns_requests (name, subdomain, request, response, src_ip, src_host) VALUES (?, ?, ?, ?, ?, ?)",
		name,
		subdomain,
		jsonRequest,
		jsonResponse,
		src_ip.String(),
		src_host,
	)
	if err != nil {
		return err
	}
	return nil
}

func StreamRequest(
	subdomain string,
	request []byte,
	response []byte,
	src_ip string,
	src_host string,
) error {
	fmt.Println("writing", subdomain)
	// get base domain
	x := map[string]interface{}{
		"created_at": time.Now().Unix(),
		"request":    string(request),
		"response":   string(response),
		"src_ip":     src_ip,
		"src_host":   src_host,
	}
	jsonString, err := json.Marshal(x)
	if err != nil {
		return err
	}
	WriteToStreams(subdomain, jsonString)
	return nil
}

func DeleteRequestsForDomain(db *sql.DB, subdomain string) error {
	_, err := db.Exec("DELETE FROM dns_requests WHERE subdomain = ?", subdomain)
	if err != nil {
		return err
	}
	return nil
}

func GetRequests(db *sql.DB, subdomain string) ([]map[string]interface{}, error) {
	tx, err := uncommittedTransaction(db)
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(
		`SELECT id, UNIX_TIMESTAMP(created_at), request, response, src_ip, src_host
FROM dns_requests
WHERE subdomain = ?
ORDER BY created_at
DESC LIMIT 30`,
		subdomain,
	)
	if err != nil {
		return make([]map[string]interface{}, 0), err
	}
	requests := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id int
		var created_at float32
		var request []byte
		var response []byte
		var src_ip string
		var src_host string
		err = rows.Scan(&id, &created_at, &request, &response, &src_ip, &src_host)
		if err != nil {
			return make([]map[string]interface{}, 0), err
		}
		x := map[string]interface{}{
			"id":         id,
			"created_at": int64(created_at),
			"request":    string(request),
			"response":   string(response),
			"src_ip":     src_ip,
			"src_host":   src_host,
		}
		requests = append(requests, x)
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func GetRecords(db *sql.DB, name string, rrtype uint16) ([]dns.RR, int, error) {
	tx, err := uncommittedTransaction(db)
	if err != nil {
		return nil, 0, err
	}
	// first get all the records
	rows, err := tx.Query(
		"SELECT content FROM dns_records WHERE name = ? ORDER BY created_at DESC",
		name,
	)
	if err != nil {
		return nil, 0, err
	}
	// next parse them
	var records []dns.RR
	for rows.Next() {
		var content []byte
		err = rows.Scan(&content)
		if err != nil {
			return nil, 0, err
		}
		record, err := ParseRecord(content)
		if err != nil {
			return nil, 0, err
		}
		records = append(records, record)
	}
	// now filter them
	filtered := make([]dns.RR, 0)
	for _, record := range records {
		if shouldReturn(rrtype, record.Header().Rrtype) {
			filtered = append(filtered, record)
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, 0, err
	}
	return filtered, len(records), nil
}

func shouldReturn(queryType uint16, recordType uint16) bool {
	if queryType == recordType {
		return true
	}
	if recordType == dns.TypeCNAME {
		return true
	}
	if queryType == dns.TypeHTTPS && (recordType == dns.TypeA || recordType == dns.TypeAAAA) {
		return true
	}
	return false
}
