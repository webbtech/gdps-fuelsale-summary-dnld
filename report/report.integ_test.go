package report

import (
	"fmt"
	"net/http"
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
	cfg     *config.Config
	request *model.Request
	report  *Report
}

// SetupTest method
func (suite *IntegSuite) SetupTest() {
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

	suite.report, err = New(suite.request, suite.cfg, "")
	suite.NoError(err)
	suite.IsType(new(Report), suite.report)

	err = suite.report.Create()
	suite.NoError(err)
}

// TestIntegSuite function
func TestIntegSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}

// TestConfig method
func (suite *IntegSuite) TestConfig() {
	suite.NotEqual("", suite.cfg.AWSRegion, "Expected AWSRegion to be populated")
}

// TestSaveToDisk method
func (suite *IntegSuite) TestSaveToDisk() {
	err := suite.report.Create()
	suite.NoError(err)

	fp, err := suite.report.SaveToDisk("../tmp")
	suite.NoError(err)
	suite.NotEqual("", fp, "Expected file path to be populated")
}

// TestCreateSignedURL method
func (suite *IntegSuite) TestCreateSignedURL() {
	err := suite.report.Create()
	suite.NoError(err)

	mockFN := fmt.Sprintf("%s_%s-%s.xlsx", reportFileName, dateFrom, dateTo)
	fn := suite.report.getFileName()
	suite.Equal(mockFN, fn)

	url, err := suite.report.CreateSignedURL()
	suite.NoError(err)
	suite.NotEqual("", url, "Expected url to be populated")

	response, err := http.Get(url)
	suite.NoError(err)
	defer response.Body.Close()
	suite.Equal(200, response.StatusCode, "Expect response code to be 200")
}
