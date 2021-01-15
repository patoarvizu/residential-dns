package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/sts"
)

type residentialDNSConfig struct {
	hostedZoneId      string
	recordName        string
	ttl               int64
	syncPeriodMinutes int
}

func main() {
	var cfg = &residentialDNSConfig{}
	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fl.StringVar(&cfg.hostedZoneId, "hosted-zone-id", "", "Route 53 hosted zone id")
	fl.StringVar(&cfg.recordName, "record-name", "", "DNS record name")
	fl.Int64Var(&cfg.ttl, "ttl", 60, "DNS record TTL")
	fl.IntVar(&cfg.syncPeriodMinutes, "sync-period-minutes", 15, "The amount of time, in minutes, to wait between syncs")
	fl.Parse(os.Args[1:])

	for {
		url := "https://api.ipify.org?format=text"
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		ip, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Printf("IP found: %s\n", ip)
		awsSession := session.Must(session.NewSession())
		stsSvc := sts.New(session.New())
		i := &sts.GetCallerIdentityInput{}
		r, err := stsSvc.GetCallerIdentity(i)
		if err != nil {
			fmt.Printf("Error: %v", err)
			time.Sleep(time.Duration(cfg.syncPeriodMinutes) * time.Minute)
			continue
		}
		fmt.Printf("Identity: %v", r)
		r53 := route53.New(awsSession)
		input := &route53.ChangeResourceRecordSetsInput{
			ChangeBatch: &route53.ChangeBatch{
				Changes: []*route53.Change{
					{
						Action: aws.String("UPSERT"),
						ResourceRecordSet: &route53.ResourceRecordSet{
							Name: aws.String(cfg.recordName),
							ResourceRecords: []*route53.ResourceRecord{
								{
									Value: aws.String(string(ip)),
								},
							},
							TTL:  aws.Int64(cfg.ttl),
							Type: aws.String("A"),
						},
					},
				},
			},
			HostedZoneId: aws.String(cfg.hostedZoneId),
		}
		_, err = r53.ChangeResourceRecordSets(input)
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
		time.Sleep(time.Duration(cfg.syncPeriodMinutes) * time.Minute)
	}
}
