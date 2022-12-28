package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"net"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/suite"
)

// some integration-style tests
var connString = "root:mysecretpassword@tcp(localhost:3306)/mysql?tls=false"

func connectTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	return db, mock
}

func makeA(name string, ip string) *dns.A {
	return &dns.A{
		Hdr: dns.RR_Header{
			Name:   name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    0,
		},
		A: net.ParseIP(ip),
	}
}

func makeCNAME(name string, target string) *dns.CNAME {
	return &dns.CNAME{
		Hdr: dns.RR_Header{
			Name:   name,
			Rrtype: dns.TypeCNAME,
			Class:  dns.ClassINET,
			Ttl:    0,
		},
		Target: target,
	}
}

func makeQuestion(name string, qtype uint16) *dns.Msg {
	return &dns.Msg{
		Question: []dns.Question{
			{
				Name:   name,
				Qtype:  qtype,
				Qclass: dns.ClassINET,
			},
		},
	}
}

type RecordSuite struct {
	suite.Suite
	db     *sql.DB
	mock   sqlmock.Sqlmock
	prefix string
	name   string
}

func (rs *RecordSuite) SetupTest() {
	rs.db, rs.mock = connectTestDB(rs.T())
	rs.prefix = randString(10)
	rs.name = rs.prefix + ".flatbo.at."
}

func (rs *RecordSuite) scaffoldMocks(rr dns.RR, dnsType uint16) {
	content, _ := json.Marshal(rr)
	rs.mock.ExpectBegin()
	rs.mock.ExpectExec("INSERT INTO dns_records").
		WithArgs(rs.name, rs.prefix, dnsType, content).
		WillReturnResult(driver.ResultNoRows)
	rs.mock.ExpectExec("UPDATE dns_serials").WillReturnResult(driver.ResultNoRows)
	rs.mock.ExpectQuery("SELECT serial").
		WillReturnRows(sqlmock.NewRows([]string{"serial"}).AddRow(11))
	rs.mock.ExpectCommit()

	rows := sqlmock.NewRows([]string{"content"}).AddRow(content)
	rs.mock.ExpectBegin()
	rs.mock.ExpectExec("SET TRANSACTION").WillReturnResult(driver.ResultNoRows)
	rs.mock.ExpectQuery("SELECT content FROM dns_records").
		WithArgs(rs.name).
		WillReturnRows(rows)
	rs.mock.ExpectCommit()
}

func TestRecordSuite(t *testing.T) {
	suite.Run(t, new(RecordSuite))
}

func (rs *RecordSuite) TestARecord() {
	record := makeA(rs.name, "1.2.3.4")
	rs.scaffoldMocks(record, dns.TypeA)

	err := InsertRecord(rs.db, makeA(rs.name, "1.2.3.4"))
	rs.NoError(err)

	response := dnsResponse(rs.db, makeQuestion(rs.name, dns.TypeA))
	// check that we got NOERROR and 1 answer
	rs.Equal(dns.RcodeSuccess, response.Rcode)
	rs.Equal(1, len(response.Answer))
}

func (rs *RecordSuite) TestCNAMERecord() {
	record := makeCNAME(rs.name, "example.com.")
	rs.scaffoldMocks(record, dns.TypeCNAME)

	err := InsertRecord(rs.db, record)
	rs.NoError(err)

	response := dnsResponse(rs.db, makeQuestion(rs.name, dns.TypeA))
	// check that we got NOERROR and 1 answer
	rs.Equal(dns.RcodeSuccess, response.Rcode)
	rs.Equal(1, len(response.Answer))
}

func (rs *RecordSuite) TestHTTPSRecord() {
	record := makeA(rs.name, "1.2.3.4")
	rs.scaffoldMocks(record, dns.TypeA)

	err := InsertRecord(rs.db, record)
	rs.NoError(err)

	response := dnsResponse(rs.db, makeQuestion(rs.name, dns.TypeHTTPS))
	// check that we got NOERROR and 1 answer
	rs.Equal(dns.RcodeSuccess, response.Rcode)
	rs.Equal(1, len(response.Answer))
}

func (rs *RecordSuite) TestNoError() {
	record := makeA(rs.name, "1.2.3.4")
	rs.scaffoldMocks(record, dns.TypeA)

	err := InsertRecord(rs.db, record)
	rs.NoError(err)

	response := dnsResponse(rs.db, makeQuestion(rs.name, dns.TypeAAAA))
	// check that we got NOERROR and 0 answers
	rs.Equal(dns.RcodeSuccess, response.Rcode)
	rs.Equal(0, len(response.Answer))
}

func (rs *RecordSuite) TestNXDOMAIN() {
	rows := sqlmock.NewRows([]string{"content"})
	rs.mock.ExpectBegin()
	rs.mock.ExpectExec("SET TRANSACTION").WillReturnResult(driver.ResultNoRows)
	rs.mock.ExpectQuery("SELECT content FROM dns_records").
		WithArgs(rs.name).
		WillReturnRows(rows)
	rs.mock.ExpectCommit()

	response := dnsResponse(rs.db, makeQuestion(rs.name, dns.TypeA))

	// check that we got NXDOMAIN
	rs.Equal(dns.RcodeNameError, response.Rcode)
}
