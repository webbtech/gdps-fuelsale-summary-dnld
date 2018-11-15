package xlsx

import (
	"os"
	"testing"
	"time"

	"github.com/pulpfree/gdps-fs-sum-dwnld/config"
	"github.com/pulpfree/gdps-fs-sum-dwnld/graphql"
	"github.com/pulpfree/gdps-fs-sum-dwnld/model"
	"github.com/stretchr/testify/suite"
)

const (
	dateFrom         = "2018-08-01"
	dateTo           = "2018-08-31"
	defaultsFilePath = "../config/defaults.yml"
	filePath         = "../tmp/testfile.xlsx"
	stationID        = "d03224a7-f1df-4863-bcaa-5c6e61af11fc"
	timeFormat       = "2006-01-02"
)

// Suite struct
type Suite struct {
	suite.Suite
	cfg     *config.Config
	request *model.Request
	file    *XLSX
	graphql *graphql.Client
}

// SetupTest method
func (suite *Suite) SetupTest() {

	os.Setenv("Stage", "test")
	dteFrom, err := time.Parse(timeFormat, dateFrom)
	dteTo, err := time.Parse(timeFormat, dateTo)
	suite.request = &model.Request{
		DateFrom: dteFrom,
		DateTo:   dteTo,
	}
	suite.cfg = &config.Config{DefaultsFilePath: defaultsFilePath}
	err = suite.cfg.Load()
	suite.NoError(err)
	suite.IsType(new(config.Config), suite.cfg)

	suite.file, err = NewFile()
	suite.NoError(err)
	suite.IsType(new(XLSX), suite.file)

	suite.graphql = graphql.New(suite.request, suite.cfg, "")
	suite.IsType(new(graphql.Client), suite.graphql)
}

// TestOutput method
func (suite *Suite) TestOutput() {

	// Fetch all report data
	fs, err := suite.graphql.FuelSalesSummary()
	suite.NoError(err)
	suite.IsType(new(model.Report), fs)

	// Create sheet tabs
	err = suite.file.NLSales(fs)
	suite.NoError(err)

	err = suite.file.DSLSales(fs)
	suite.NoError(err)

	_, err = suite.file.OutputToDisk(filePath)
	suite.NoError(err)
}

// TestXLSXSuite function
func TestXLSXSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
