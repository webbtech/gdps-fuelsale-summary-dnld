package validate

import (
	"os"
	"testing"
	"time"

	"github.com/pulpfree/gdps-fs-sum-dwnld/model"
	"github.com/stretchr/testify/suite"
)

const (
	dateFrom     = "2018-08-01"
	dateTo       = "2018-08-31"
	dateToFuture = "2118-08-31"
	timeFormat   = "2006-01-02"
)

// IntegSuite struct
type IntegSuite struct {
	suite.Suite
	request *model.Request
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
	suite.NoError(err)
	suite.IsType(new(model.Request), suite.request)
}

// TestDate method
func (suite *IntegSuite) TestDate() {
	dteFrom, err := Date(dateFrom)
	dteTo, err := Date(dateTo)
	suite.NoError(err)
	suite.IsType(time.Time{}, dteFrom)
	suite.IsType(time.Time{}, dteTo)

	dteToFuture, err := Date(dateToFuture)
	suite.IsType(time.Time{}, dteToFuture)
	suite.Error(err)
}

// TestRequestInput method
func (suite *IntegSuite) TestRequestInput() {
	req := &model.RequestInput{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	res, err := RequestInput(req)
	suite.NoError(err)
	suite.IsType(&model.Request{}, res)
}

// TestUnitSuite function
func TestUnitSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
