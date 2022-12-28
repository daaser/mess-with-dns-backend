package main

import (
	"encoding/json"
	"fmt"

	"github.com/miekg/dns"
)

type UnknownRequest struct {
	Hdr dns.RR_Header
}

func ParseRecord(jsonString []byte) (dns.RR, error) {
	rr, err := parseRecord(jsonString)
	if err != nil {
		return nil, err
	}
	// make sure we have a valid RR
	// this prevents problems like invalid FQDNs in a record's fields
	msg := make([]byte, dns.Len(rr))
	_, err = dns.PackRR(rr, msg, 0, nil, false)
	if err != nil {
		return nil, fmt.Errorf("invalid RR: %s, %#v", err, rr)
	}
	return rr, nil
}

func parseRecord(jsonString []byte) (dns.RR, error) {
	var unknown UnknownRequest
	err := json.Unmarshal([]byte(jsonString), &unknown)
	if err != nil {
		return nil, err
	}

	switch unknown.Hdr.Rrtype {
	case dns.TypeA:
		var a dns.A
		err = json.Unmarshal([]byte(jsonString), &a)
		if err != nil {
			return nil, err
		}
		return &a, nil

	case dns.TypeAAAA:
		var aaaa dns.AAAA
		err = json.Unmarshal([]byte(jsonString), &aaaa)
		if err != nil {
			return nil, err
		}
		return &aaaa, nil

	case dns.TypeAFSDB:
		var afsdb dns.AFSDB
		err = json.Unmarshal([]byte(jsonString), &afsdb)
		if err != nil {
			return nil, err
		}
		return &afsdb, nil

	case dns.TypeANY:
		var any dns.ANY
		err = json.Unmarshal([]byte(jsonString), &any)
		if err != nil {
			return nil, err
		}
		return &any, nil

	case dns.TypeAPL:
		var apl dns.APL
		err = json.Unmarshal([]byte(jsonString), &apl)
		if err != nil {
			return nil, err
		}
		return &apl, nil

	case dns.TypeCAA:
		var caa dns.CAA
		err = json.Unmarshal([]byte(jsonString), &caa)
		if err != nil {
			return nil, err
		}
		return &caa, nil

	case dns.TypeCDNSKEY:
		var cdnskey dns.CDNSKEY
		err = json.Unmarshal([]byte(jsonString), &cdnskey)
		if err != nil {
			return nil, err
		}
		return &cdnskey, nil

	case dns.TypeCDS:
		var cds dns.CDS
		err = json.Unmarshal([]byte(jsonString), &cds)
		if err != nil {
			return nil, err
		}
		return &cds, nil

	case dns.TypeCERT:
		var cert dns.CERT
		err = json.Unmarshal([]byte(jsonString), &cert)
		if err != nil {
			return nil, err
		}
		return &cert, nil

	case dns.TypeCNAME:
		var cname dns.CNAME
		err = json.Unmarshal([]byte(jsonString), &cname)
		if err != nil {
			return nil, err
		}
		return &cname, nil

	case dns.TypeCSYNC:
		var csync dns.CSYNC
		err = json.Unmarshal([]byte(jsonString), &csync)
		if err != nil {
			return nil, err
		}
		return &csync, nil

	case dns.TypeDHCID:
		var dhcid dns.DHCID
		err = json.Unmarshal([]byte(jsonString), &dhcid)
		if err != nil {
			return nil, err
		}
		return &dhcid, nil

	case dns.TypeDLV:
		var dlv dns.DLV
		err = json.Unmarshal([]byte(jsonString), &dlv)
		if err != nil {
			return nil, err
		}
		return &dlv, nil

	case dns.TypeDNAME:
		var dname dns.DNAME
		err = json.Unmarshal([]byte(jsonString), &dname)
		if err != nil {
			return nil, err
		}
		return &dname, nil

	case dns.TypeDNSKEY:
		var dnskey dns.DNSKEY
		err = json.Unmarshal([]byte(jsonString), &dnskey)
		if err != nil {
			return nil, err
		}
		return &dnskey, nil

	case dns.TypeDS:
		var ds dns.DS
		err = json.Unmarshal([]byte(jsonString), &ds)
		if err != nil {
			return nil, err
		}
		return &ds, nil

	case dns.TypeEID:
		var eid dns.EID
		err = json.Unmarshal([]byte(jsonString), &eid)
		if err != nil {
			return nil, err
		}
		return &eid, nil

	case dns.TypeEUI48:
		var eui48 dns.EUI48
		err = json.Unmarshal([]byte(jsonString), &eui48)
		if err != nil {
			return nil, err
		}
		return &eui48, nil

	case dns.TypeEUI64:
		var eui64 dns.EUI64
		err = json.Unmarshal([]byte(jsonString), &eui64)
		if err != nil {
			return nil, err
		}
		return &eui64, nil

	case dns.TypeGID:
		var gid dns.GID
		err = json.Unmarshal([]byte(jsonString), &gid)
		if err != nil {
			return nil, err
		}
		return &gid, nil

	case dns.TypeGPOS:
		var gpos dns.GPOS
		err = json.Unmarshal([]byte(jsonString), &gpos)
		if err != nil {
			return nil, err
		}
		return &gpos, nil

	case dns.TypeHINFO:
		var hinfo dns.HINFO
		err = json.Unmarshal([]byte(jsonString), &hinfo)
		if err != nil {
			return nil, err
		}
		return &hinfo, nil

	case dns.TypeHIP:
		var hip dns.HIP
		err = json.Unmarshal([]byte(jsonString), &hip)
		if err != nil {
			return nil, err
		}
		return &hip, nil

	case dns.TypeHTTPS:
		var https dns.HTTPS
		err = json.Unmarshal([]byte(jsonString), &https)
		if err != nil {
			return nil, err
		}
		return &https, nil

	case dns.TypeKEY:
		var key dns.KEY
		err = json.Unmarshal([]byte(jsonString), &key)
		if err != nil {
			return nil, err
		}
		return &key, nil

	case dns.TypeKX:
		var kx dns.KX
		err = json.Unmarshal([]byte(jsonString), &kx)
		if err != nil {
			return nil, err
		}
		return &kx, nil

	case dns.TypeL32:
		var l32 dns.L32
		err = json.Unmarshal([]byte(jsonString), &l32)
		if err != nil {
			return nil, err
		}
		return &l32, nil

	case dns.TypeL64:
		var l64 dns.L64
		err = json.Unmarshal([]byte(jsonString), &l64)
		if err != nil {
			return nil, err
		}
		return &l64, nil

	case dns.TypeLOC:
		var loc dns.LOC
		err = json.Unmarshal([]byte(jsonString), &loc)
		if err != nil {
			return nil, err
		}
		return &loc, nil

	case dns.TypeLP:
		var lp dns.LP
		err = json.Unmarshal([]byte(jsonString), &lp)
		if err != nil {
			return nil, err
		}
		return &lp, nil

	case dns.TypeMB:
		var mb dns.MB
		err = json.Unmarshal([]byte(jsonString), &mb)
		if err != nil {
			return nil, err
		}
		return &mb, nil

	case dns.TypeMD:
		var md dns.MD
		err = json.Unmarshal([]byte(jsonString), &md)
		if err != nil {
			return nil, err
		}
		return &md, nil

	case dns.TypeMF:
		var mf dns.MF
		err = json.Unmarshal([]byte(jsonString), &mf)
		if err != nil {
			return nil, err
		}
		return &mf, nil

	case dns.TypeMG:
		var mg dns.MG
		err = json.Unmarshal([]byte(jsonString), &mg)
		if err != nil {
			return nil, err
		}
		return &mg, nil

	case dns.TypeMINFO:
		var minfo dns.MINFO
		err = json.Unmarshal([]byte(jsonString), &minfo)
		if err != nil {
			return nil, err
		}
		return &minfo, nil

	case dns.TypeMR:
		var mr dns.MR
		err = json.Unmarshal([]byte(jsonString), &mr)
		if err != nil {
			return nil, err
		}
		return &mr, nil

	case dns.TypeMX:
		var mx dns.MX
		err = json.Unmarshal([]byte(jsonString), &mx)
		if err != nil {
			return nil, err
		}
		return &mx, nil

	case dns.TypeNAPTR:
		var naptr dns.NAPTR
		err = json.Unmarshal([]byte(jsonString), &naptr)
		if err != nil {
			return nil, err
		}
		return &naptr, nil

	case dns.TypeNID:
		var nid dns.NID
		err = json.Unmarshal([]byte(jsonString), &nid)
		if err != nil {
			return nil, err
		}
		return &nid, nil

	case dns.TypeNIMLOC:
		var nimloc dns.NIMLOC
		err = json.Unmarshal([]byte(jsonString), &nimloc)
		if err != nil {
			return nil, err
		}
		return &nimloc, nil

	case dns.TypeNINFO:
		var ninfo dns.NINFO
		err = json.Unmarshal([]byte(jsonString), &ninfo)
		if err != nil {
			return nil, err
		}
		return &ninfo, nil

	case dns.TypeNS:
		var ns dns.NS
		err = json.Unmarshal([]byte(jsonString), &ns)
		if err != nil {
			return nil, err
		}
		return &ns, nil

	case dns.TypeNSEC:
		var nsec dns.NSEC
		err = json.Unmarshal([]byte(jsonString), &nsec)
		if err != nil {
			return nil, err
		}
		return &nsec, nil

	case dns.TypeNSEC3:
		var nsec3 dns.NSEC3
		err = json.Unmarshal([]byte(jsonString), &nsec3)
		if err != nil {
			return nil, err
		}
		return &nsec3, nil

	case dns.TypeNSEC3PARAM:
		var nsec3param dns.NSEC3PARAM
		err = json.Unmarshal([]byte(jsonString), &nsec3param)
		if err != nil {
			return nil, err
		}
		return &nsec3param, nil

	case dns.TypeNULL:
		var null dns.NULL
		err = json.Unmarshal([]byte(jsonString), &null)
		if err != nil {
			return nil, err
		}
		return &null, nil

	case dns.TypeOPENPGPKEY:
		var openpgpkey dns.OPENPGPKEY
		err = json.Unmarshal([]byte(jsonString), &openpgpkey)
		if err != nil {
			return nil, err
		}
		return &openpgpkey, nil

	case dns.TypeOPT:
		var opt dns.OPT
		err = json.Unmarshal([]byte(jsonString), &opt)
		if err != nil {
			return nil, err
		}
		return &opt, nil

	case dns.TypePTR:
		var ptr dns.PTR
		err = json.Unmarshal([]byte(jsonString), &ptr)
		if err != nil {
			return nil, err
		}
		return &ptr, nil

	case dns.TypePX:
		var px dns.PX
		err = json.Unmarshal([]byte(jsonString), &px)
		if err != nil {
			return nil, err
		}
		return &px, nil

	case dns.TypeRKEY:
		var rkey dns.RKEY
		err = json.Unmarshal([]byte(jsonString), &rkey)
		if err != nil {
			return nil, err
		}
		return &rkey, nil

	case dns.TypeRP:
		var rp dns.RP
		err = json.Unmarshal([]byte(jsonString), &rp)
		if err != nil {
			return nil, err
		}
		return &rp, nil

	case dns.TypeRRSIG:
		var rrsig dns.RRSIG
		err = json.Unmarshal([]byte(jsonString), &rrsig)
		if err != nil {
			return nil, err
		}
		return &rrsig, nil

	case dns.TypeRT:
		var rt dns.RT
		err = json.Unmarshal([]byte(jsonString), &rt)
		if err != nil {
			return nil, err
		}
		return &rt, nil

	case dns.TypeSIG:
		var sig dns.SIG
		err = json.Unmarshal([]byte(jsonString), &sig)
		if err != nil {
			return nil, err
		}
		return &sig, nil

	case dns.TypeSMIMEA:
		var smimea dns.SMIMEA
		err = json.Unmarshal([]byte(jsonString), &smimea)
		if err != nil {
			return nil, err
		}
		return &smimea, nil

	case dns.TypeSOA:
		var soa dns.SOA
		err = json.Unmarshal([]byte(jsonString), &soa)
		if err != nil {
			return nil, err
		}
		return &soa, nil

	case dns.TypeSPF:
		var spf dns.SPF
		err = json.Unmarshal([]byte(jsonString), &spf)
		if err != nil {
			return nil, err
		}
		return &spf, nil

	case dns.TypeSRV:
		var srv dns.SRV
		err = json.Unmarshal([]byte(jsonString), &srv)
		if err != nil {
			return nil, err
		}
		return &srv, nil

	case dns.TypeSSHFP:
		var sshfp dns.SSHFP
		err = json.Unmarshal([]byte(jsonString), &sshfp)
		if err != nil {
			return nil, err
		}
		return &sshfp, nil

	case dns.TypeSVCB:
		var svcb dns.SVCB
		err = json.Unmarshal([]byte(jsonString), &svcb)
		if err != nil {
			return nil, err
		}
		return &svcb, nil

	case dns.TypeTA:
		var ta dns.TA
		err = json.Unmarshal([]byte(jsonString), &ta)
		if err != nil {
			return nil, err
		}
		return &ta, nil

	case dns.TypeTALINK:
		var talink dns.TALINK
		err = json.Unmarshal([]byte(jsonString), &talink)
		if err != nil {
			return nil, err
		}
		return &talink, nil

	case dns.TypeTKEY:
		var tkey dns.TKEY
		err = json.Unmarshal([]byte(jsonString), &tkey)
		if err != nil {
			return nil, err
		}
		return &tkey, nil

	case dns.TypeTLSA:
		var tlsa dns.TLSA
		err = json.Unmarshal([]byte(jsonString), &tlsa)
		if err != nil {
			return nil, err
		}
		return &tlsa, nil

	case dns.TypeTSIG:
		var tsig dns.TSIG
		err = json.Unmarshal([]byte(jsonString), &tsig)
		if err != nil {
			return nil, err
		}
		return &tsig, nil

	case dns.TypeTXT:
		var txt dns.TXT
		err = json.Unmarshal([]byte(jsonString), &txt)
		if err != nil {
			return nil, err
		}
		return &txt, nil

	case dns.TypeUID:
		var uid dns.UID
		err = json.Unmarshal([]byte(jsonString), &uid)
		if err != nil {
			return nil, err
		}
		return &uid, nil

	case dns.TypeUINFO:
		var uinfo dns.UINFO
		err = json.Unmarshal([]byte(jsonString), &uinfo)
		if err != nil {
			return nil, err
		}
		return &uinfo, nil

	case dns.TypeURI:
		var uri dns.URI
		err = json.Unmarshal([]byte(jsonString), &uri)
		if err != nil {
			return nil, err
		}
		return &uri, nil

	case dns.TypeX25:
		var x25 dns.X25
		err = json.Unmarshal([]byte(jsonString), &x25)
		if err != nil {
			return nil, err
		}
		return &x25, nil

	case dns.TypeZONEMD:
		var zonemd dns.ZONEMD
		err = json.Unmarshal([]byte(jsonString), &zonemd)
		if err != nil {
			return nil, err
		}
		return &zonemd, nil

	case dns.TypeNSAPPTR:
		var nsapptr dns.NSAPPTR
		err = json.Unmarshal([]byte(jsonString), &nsapptr)
		if err != nil {
			return nil, err
		}
		return &nsapptr, nil
	}

	return nil, fmt.Errorf("unhandled RR type: %d", unknown.Hdr.Rrtype)
}
