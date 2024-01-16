package main

import (
	"log"
	"net"
	"strings"
)

type domainChecks struct {
	HasMX       bool   `json:"hasMX"`
	HasSPF      bool   `json:"hasSPF"`
	SPFRecord   string `json:"spfRecord"`
	HasDMARC    bool   `json:"hasDMARC"`
	DMARCRecord string `json:"dmarcRecord"`
}

func checkDomain(domain string) (domainChecks, error) {
	hasMX := false
	hasSPF := false
	hasDMARC := false
	var spfRecord, dmarcRecord string

	// Function to check if a specific record is present in the TXT records
	checkRecord := func(records []string, prefix string) (bool, string) {
		for _, record := range records {
			if strings.HasPrefix(record, prefix) {
				return true, record
			}
		}
		return false, ""
	}

	// Lookup MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error looking up MX records: %v\n", err)
	}

	// Check if MX records are present
	hasMX = len(mxRecords) > 0

	// Lookup TXT records
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error looking up TXT records: %v\n", err)
	}

	// Check for SPF record
	hasSPF, spfRecord = checkRecord(txtRecords, "v=spf1")

	// Lookup DMARC records
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error looking up DMARC records: %v\n", err)
	}

	// Check for DMARC record
	hasDMARC, dmarcRecord = checkRecord(dmarcRecords, "v=DMARC1")

	return domainChecks{
		HasMX:       hasMX,
		HasSPF:      hasSPF,
		SPFRecord:   spfRecord,
		HasDMARC:    hasDMARC,
		DMARCRecord: dmarcRecord,
	}, err
}
