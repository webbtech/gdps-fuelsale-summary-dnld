package graphql

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/machinebox/graphql"
	"github.com/pulpfree/gdps-fs-sum-dwnld/config"
	"github.com/pulpfree/gdps-fs-sum-dwnld/model"
)

const (
	timeLongFrmt  = "2006-01-02"
	timeShortFrmt = "20060102"
)

// Client struct
type Client struct {
	hdrs    http.Header
	client  *graphql.Client
	request *model.Request
}

// New graphql client
func New(req *model.Request, cfg *config.Config, authToken string) (c *Client) {

	hdrs := http.Header{}
	if len(authToken) > 0 {
		hdrs.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	}

	c = &Client{
		client:  graphql.NewClient(cfg.GraphqlURI),
		hdrs:    hdrs,
		request: req,
	}

	return c
}

// FuelSalesSummary method
func (c *Client) FuelSalesSummary() (rpt *model.Report, err error) {

	req := graphql.NewRequest(`
    query ($dateFrom: String!, $dateTo: String!) {
      fuelSaleRangeSummary(dateFrom: $dateFrom, dateTo: $dateTo) {
        dateStart
        dateEnd
        summary {
          stationID
          stationName
          hasDSL
          fuels {
            NL
            DSL
          }
        }
      }
    }
  `)

	req.Var("dateFrom", formattedDate(c.request.DateFrom))
	req.Var("dateTo", formattedDate(c.request.DateTo))
	req.Header = c.hdrs

	ctx := context.Background()
	err = c.client.Run(ctx, req, &rpt)
	if err != nil {
		log.Errorf("error running graphql client: %s", err.Error())
		return nil, err
	}

	rpt.FuelSales.DateFrom = intToDate(rpt.FuelSales.DateStart)
	rpt.FuelSales.DateTo = intToDate(rpt.FuelSales.DateEnd)

	return rpt, err
}

//
// ======================== Helper Functions =============================== //
//

// formattedDate function
func formattedDate(date time.Time) string {
	return date.Format(timeLongFrmt)
}

// intToDate function
// takes a integer date and convert to time struct
func intToDate(dte int) time.Time {
	date, _ := time.Parse(timeShortFrmt, strconv.Itoa(dte))
	return date
}
