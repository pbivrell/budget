package influx

import (
	"fmt"
	"net/url"
	"time"

	client "github.com/influxdata/influxdb1-client"
)

type Client struct {
	measurement string
	tags        map[string]string
	fields      map[string]interface{}
	precision   string
	series      string
	points      []point
	*client.Client
}

type point struct {
	fields map[string]interface{}
	tags   map[string]string
	time   time.Time
}

func New(dbHost string) (*Client, error) {
	host, err := url.Parse(dbHost)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse series host as a url: %v", err)
	}
	con, err := client.NewClient(client.Config{URL: *host})
	if err != nil {
		return nil, fmt.Errorf("Failed to create series client: %v", err)
	}
	return &Client{
		"",
		make(map[string]string),
		make(map[string]interface{}),
		"",
		"",
		make([]point, 0),
		con,
	}, nil
}

func (con *Client) Measurement(name string) *Client {
	con.measurement = name
	return con
}

func (con *Client) Precision(precision string) *Client {
	con.precision = precision
	return con
}

func (con *Client) Series(series string) *Client {
	con.series = series
	return con
}

func (con *Client) Point(tags map[string]string, fields map[string]interface{}, inTime time.Time) *Client {
	t := make(map[string]string)
	for k, v := range tags {
		t[k] = v
	}
	f := make(map[string]interface{})
	for k, v := range fields {
		f[k] = v
	}
	con.points = append(con.points, point{tags: t, fields: f, time: inTime})
	return con
}

func (con *Client) PointWrite(tags map[string]string, fields map[string]interface{}, inTime time.Time) error {
	t := make(map[string]string)
	for k, v := range tags {
		t[k] = v
	}
	f := make(map[string]interface{})
	for k, v := range fields {
		f[k] = v
	}
	pts := make([]client.Point, 1)
	pts[0] = client.Point{
		Measurement: con.measurement,
		Tags:        t,
		Fields:      f,
		Time:        inTime,
		Precision:   con.precision,
	}

	bps := client.BatchPoints{
		Points:   pts,
		Database: con.series,
	}
	_, err := con.Write(bps)
	if err != nil {
		return fmt.Errorf("Failed to write data to result series: %v", err)
	}

	return nil

}

func (con *Client) BatchWrite() error {
	pts := make([]client.Point, len(con.points))
	for i, v := range con.points {
		pts[i] = client.Point{
			Measurement: con.measurement,
			Tags:        v.tags,
			Fields:      v.fields,
			Time:        v.time,
			Precision:   con.precision,
		}
	}
	bps := client.BatchPoints{
		Points:   pts,
		Database: con.series,
	}
	_, err := con.Write(bps)
	if err != nil {
		return fmt.Errorf("Failed to write %d data points to result series: %v", len(pts), err)
	}
	return nil
}

func (con *Client) Reset() {
	con.measurement = ""
	con.tags = make(map[string]string)
	con.fields = make(map[string]interface{})
	con.precision = ""
	con.series = ""
	con.points = make([]point, 0)
}
