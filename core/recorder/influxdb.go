package recorder

import (
	fcom "github.com/hyperbench/hyperbench-common/common"
	"github.com/influxdata/influxdb/client"
	"net/url"
	"time"
)

type influxdb struct {
	url       *url.URL
	database  string
	username  string
	password  string
	benchmark string

	client *client.Client
}

func (i *influxdb) process(report fcom.Report) {
	go i.send(report)
}

func (i *influxdb) send(report fcom.Report) {
	pts := make([]client.Point, 0, len(report.Cur.Results))

	for _, r := range report.Cur.Results {
		pts = append(pts, client.Point{
			Measurement: "current",
			Tags: map[string]string{
				"label": r.Label,
			},
			Time: time.Unix(0, r.Time),
			Fields: map[string]interface{}{
				"send":         r.Num,
				"duration":     r.Duration,
				"send_rate":    float64(r.Num) * float64(time.Second) / float64(r.Duration),
				"succeeded":    r.Statuses[fcom.Success],
				"failed":       r.Statuses[fcom.Failure],
				"confirmed":    r.Statuses[fcom.Confirm],
				"unknown":      r.Statuses[fcom.Unknown],
				"send_avg":     r.Send.Avg,
				"send_p0":      r.Send.P0,
				"send_p50":     r.Send.P50,
				"send_p90":     r.Send.P90,
				"send_p95":     r.Send.P95,
				"send_p99":     r.Send.P99,
				"send_p100":    r.Send.P100,
				"confirm_avg":  r.Confirm.Avg,
				"confirm_p0":   r.Confirm.P0,
				"confirm_p50":  r.Confirm.P50,
				"confirm_p90":  r.Confirm.P90,
				"confirm_p95":  r.Confirm.P95,
				"confirm_p99":  r.Confirm.P99,
				"confirm_p100": r.Confirm.P100,
				"write_avg":    r.Write.Avg,
				"write_p0":     r.Write.P0,
				"write_p50":    r.Write.P50,
				"write_p90":    r.Write.P90,
				"write_p95":    r.Write.P95,
				"write_p99":    r.Write.P99,
				"write_p100":   r.Write.P100,
			},
		})
	}

	bps := client.BatchPoints{
		Points:   pts,
		Database: i.database,
	}

	_, err := i.client.Write(bps)
	if err != nil {
		fcom.GetLogger("influx").Error(err)
	}
}

func (i *influxdb) release() {
}

func newInfluxdb(benchmark string, URL string, database string, username string, password string) (*influxdb, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	i := &influxdb{
		url:       u,
		database:  database,
		username:  username,
		password:  password,
		benchmark: benchmark,
	}
	err = i.makeClient()
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (i *influxdb) makeClient() (err error) {
	i.client, err = client.NewClient(client.Config{
		URL:      *i.url,
		Username: i.username,
		Password: i.password,
		Timeout:  30 * time.Second,
	})
	return
}
