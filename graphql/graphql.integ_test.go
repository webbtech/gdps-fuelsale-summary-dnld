package graphql

import (
	"os"
	"testing"
	"time"

	"github.com/pulpfree/gdps-fs-sum-dwnld/config"
	"github.com/pulpfree/gdps-fs-sum-dwnld/model"
	"github.com/stretchr/testify/suite"
)

const (
	dateFrom         = "2018-08-01"
	dateTo           = "2018-08-31"
	defaultsFilePath = "../config/defaults.yml"
	timeFormat       = "2006-01-02"
)

// IntegSuite struct
type IntegSuite struct {
	suite.Suite
	client  *Client
	cfg     *config.Config
	request *model.Request
}

// SetupTest method
func (suite *IntegSuite) SetupTest() {
	os.Setenv("Stage", "test")
	dteFrom, err := time.Parse(timeFormat, dateFrom)
	dteTo, err := time.Parse(timeFormat, dateTo)
	req := &model.Request{
		DateFrom: dteFrom,
		DateTo:   dteTo,
	}
	suite.cfg = &config.Config{DefaultsFilePath: defaultsFilePath}
	err = suite.cfg.Load()
	suite.NoError(err)
	suite.IsType(new(config.Config), suite.cfg)

	suite.client = New(req, suite.cfg, "")
	suite.NoError(err)
	suite.IsType(new(Client), suite.client)
}

// TestFuelSalesSummary method
func (suite *IntegSuite) TestFuelSalesSummary() {
	res, err := suite.client.FuelSalesSummary()
	suite.NoError(err)
	suite.IsType(new(model.Report), res)
	suite.IsType(new(model.FuelSales), res.FuelSales)
}

// TestIntegSuite function
func TestIntegSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
