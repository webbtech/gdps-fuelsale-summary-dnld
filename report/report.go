package report

import (
	"fmt"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/pulpfree/gdps-fs-sum-dwnld/awsservices"
	"github.com/pulpfree/gdps-fs-sum-dwnld/config"
	"github.com/pulpfree/gdps-fs-sum-dwnld/graphql"
	"github.com/pulpfree/gdps-fs-sum-dwnld/model"
	"github.com/pulpfree/gdps-fs-sum-dwnld/xlsx"
)

// ReportName constants
const (
	reportFileName = "FuelSalesSummaryReport"
	timeFrmt       = "2006-01-02"
)

// Report struct
type Report struct {
	authToken string
	cfg       *config.Config
	request   *model.Request
	file      *xlsx.XLSX
	filenm    string
}

// New function
func New(req *model.Request, cfg *config.Config, authToken string) (r *Report, err error) {
	r = &Report{
		authToken: authToken,
		cfg:       cfg,
		request:   req,
	}
	r.setFileName()
	return r, err
}

// Create method
func (r *Report) Create() (err error) {
	// Init graphql and xlsx packages
	client := graphql.New(r.request, r.cfg, r.authToken)
	r.file, err = xlsx.NewFile()
	if err != nil {
		return err
	}

	// Fetch and create Fuel Sales Summary
	fs, err := client.FuelSalesSummary()
	if err != nil {
		log.Errorf("Error fetching FuelSalesSummary: %s", err)
		return err
	}

	err = r.file.NLSales(fs)
	if err != nil {
		return err
	}

	err = r.file.DSLSales(fs)
	if err != nil {
		return err
	}

	return err
}

// SaveToDisk method
func (r *Report) SaveToDisk(dir string) (fp string, err error) {

	filePath := path.Join(dir, r.getFileName())
	fp, err = r.file.OutputToDisk(filePath)
	if err != nil {
		return "", err
	}
	return fp, err
}

// CreateSignedURL method
func (r *Report) CreateSignedURL() (url string, err error) {

	output, err := r.file.OutputFile()
	if err != nil {
		return "", err
	}

	s3Serv, err := awsservices.NewS3(r.cfg)
	filePrefix := path.Join(r.cfg.S3FilePrefix, r.getFileName())

	return s3Serv.GetSignedURL(filePrefix, &output)
}

//
// ======================== Helper Functions =============================== //
//

func (r *Report) setFileName() {
	r.filenm = fmt.Sprintf("%s_%s-%s.xlsx", reportFileName, r.request.DateFrom.Format(timeFrmt), r.request.DateTo.Format(timeFrmt))
}

func (r *Report) getFileName() string {
	return r.filenm
}
