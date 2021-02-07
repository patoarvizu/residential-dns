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
	"github.com/cloudflare/cloudflare-go"
)

type residentialDNSConfig struct {
	provider            string
	route53HostedZoneId string
	cloudflareZoneName  string
	cloudflareProxied   bool
	recordName          string
	ttl                 int64
	syncPeriodMinutes   int
}

func main() {
	var cfg = &residentialDNSConfig{}
	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fl.StringVar(&cfg.provider, "provider", "", "DNS provider where the record will be created. Valid values are: 'route53' and 'cloudflare'")
	fl.StringVar(&cfg.route53HostedZoneId, "route53-hosted-zone-id", "", "Route 53 hosted zone id")
	fl.StringVar(&cfg.cloudflareZoneName, "cloudflare-zone-name", "", "Cloudflare zone name")
	fl.BoolVar(&cfg.cloudflareProxied, "cloudflare-proxied", false, "Set to true if Cloudflare requests to this hostname will be proxied through Cloudflare")
	fl.StringVar(&cfg.recordName, "record-name", "", "DNS record name")
	fl.Int64Var(&cfg.ttl, "ttl", 60, "DNS record TTL")
	fl.IntVar(&cfg.syncPeriodMinutes, "sync-period-minutes", 15, "The amount of time, in minutes, to wait between syncs")
	fl.Parse(os.Args[1:])

	providerErrors := 0
	for {
		if providerErrors == 3 {
			panic("Controller has failed three times, restarting")
		}
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
		if cfg.provider == "route53" {
			awsSession := session.Must(session.NewSession())
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
				HostedZoneId: aws.String(cfg.route53HostedZoneId),
			}
			_, err = r53.ChangeResourceRecordSets(input)
			if err != nil {
				fmt.Printf("Route 53 error: %v", err)
				providerErrors++
				continue
			}
		} else if cfg.provider == "cloudflare" {
			api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
			if err != nil {
				fmt.Printf("Cloudflare error: %v", err)
				providerErrors++
				continue
			}
			id, err := api.ZoneIDByName(cfg.cloudflareZoneName)
			if err != nil {
				fmt.Printf("Cloudflare error: %v", err)
				providerErrors++
				continue
			}
			record := cloudflare.DNSRecord{
				Type:    "A",
				Name:    cfg.recordName,
				Content: string(ip),
				TTL:     int(cfg.ttl),
				Proxied: cfg.cloudflareProxied,
			}
			existing, err := api.DNSRecords(id, cloudflare.DNSRecord{
				Type: "A",
				Name: cfg.recordName,
			})
			if len(existing) == 1 {
				err = api.UpdateDNSRecord(id, existing[0].ID, record)
				if err != nil {
					fmt.Printf("Cloudflare error: %v", err)
					providerErrors++
					continue
				}
			} else {
				_, err = api.CreateDNSRecord(id, record)
				if err != nil {
					fmt.Printf("Cloudflare error: %v", err)
					providerErrors++
					continue
				}
			}
		}
		time.Sleep(time.Duration(cfg.syncPeriodMinutes) * time.Minute)
	}
}
