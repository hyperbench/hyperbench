package recorder

import (
	"context"
	"fmt"
	fcom "github.com/hyperbench/hyperbench-common/common"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/spf13/viper"
	"time"
)

type influxdb struct {
	url       string
	database  string
	username  string
	password  string
	benchmark string

	writer api.WriteAPI
	client influxdb2.Client
}

func (i *influxdb) process(report fcom.Report) {
	go i.send(report)
}

func (i *influxdb) processStatistic(rs *fcom.RemoteStatistic) {
	//go i.send(report)
	viper.ConfigFileUsed()
	point := influxdb2.NewPoint("statistic", map[string]string{}, map[string]interface{}{
		"Start":          time.Unix(0, rs.Start).Format("2006-01-02 15:04:05"),
		"End":            time.Unix(0, rs.End).Format("2006-01-02 15:04:05"),
		"Tps":            rs.CTps,
		"Bps":            rs.Bps,
		"SendTps":        rs.Tps,
		"SendTxNum":      rs.SentTx,
		"MissTxNum":      rs.MissedTx,
		"TxNum":          rs.TxNum,
		"BlockNum":       rs.BlockNum,
		"blockChainType": viper.GetString(fcom.ClientTypePath),
	}, time.Now())
	i.writer.WritePoint(point)
}

func (i *influxdb) send(report fcom.Report) {

	for _, r := range report.Cur.Results {
		point := influxdb2.NewPoint("current", map[string]string{
			"label": r.Label,
		}, map[string]interface{}{
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
		}, time.Unix(0, r.Time))
		i.writer.WritePoint(point)
	}
}

func (i *influxdb) release() {
	i.client.Close()
}

func newInfluxdb(benchmark string, URL string, database string, username string, password string) (*influxdb, error) {
	//u, err := url.Parse(URL)
	//if err != nil {
	//	return nil, err
	//}
	i := &influxdb{
		url:       URL,
		database:  database,
		username:  username,
		password:  password,
		benchmark: benchmark,
	}
	err := i.makeClient()
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (i *influxdb) makeClient() (err error) {
	client := influxdb2.NewClient(i.url, fmt.Sprintf("%s:%s", i.username, i.password))
	_, err = client.Ping(context.Background())
	if err != nil {
		return err
	}
	writeAPI := client.WriteAPI("", i.database)
	i.writer = writeAPI
	i.client = client
	//i.client, err = client.NewClient(client.Config{
	//	URL:      *i.url,
	//	Username: i.username,
	//	Password: i.password,
	//	Timeout:  30 * time.Second,
	//})
	return
}
