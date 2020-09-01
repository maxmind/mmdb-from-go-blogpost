package main

import (
	"log"
	"net"
	"os"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

func main() {
	// Load the database we wish to enrich.
	writer, err := mmdbwriter.Load("GeoLite2-Country.mmdb", mmdbwriter.Options{})
	if err != nil {
		log.Fatal(err)
	}

	// Define and insert the new data.
	_, sreNet, err := net.ParseCIDR("56.0.0.0/16")
	if err != nil {
		log.Fatal(err)
	}
	sreData := mmdbtype.Map{
		"AcmeCorp.DeptName": mmdbtype.String("SRE"),
		"AcmeCorp.Environments": mmdbtype.Slice{
			mmdbtype.String("development"),
			mmdbtype.String("staging"),
			mmdbtype.String("production"),
		},
	}
	if err := writer.InsertFunc(sreNet, inserter.TopLevelMergeWith(sreData)); err != nil {
		log.Fatal(err)
	}

	_, devNet, err := net.ParseCIDR("56.1.0.0/16")
	if err != nil {
		log.Fatal(err)
	}
	devData := mmdbtype.Map{
		"AcmeCorp.DeptName": mmdbtype.String("Development"),
		"AcmeCorp.Environments": mmdbtype.Slice{
			mmdbtype.String("development"),
			mmdbtype.String("staging"),
		},
	}
	if err := writer.InsertFunc(devNet, inserter.TopLevelMergeWith(devData)); err != nil {
		log.Fatal(err)
	}

	_, mgmtNet, err := net.ParseCIDR("56.2.0.0/16")
	if err != nil {
		log.Fatal(err)
	}
	mgmtData := mmdbtype.Map{
		"AcmeCorp.DeptName": mmdbtype.String("Management"),
		"AcmeCorp.Environments": mmdbtype.Slice{
			mmdbtype.String("development"),
			mmdbtype.String("staging"),
		},
	}
	if err := writer.InsertFunc(mgmtNet, inserter.TopLevelMergeWith(mgmtData)); err != nil {
		log.Fatal(err)
	}

	// Write the newly enriched DB to the filesystem.
	fh, err := os.Create("GeoLite2-Country-with-Department-Data.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}
