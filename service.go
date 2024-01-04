package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// The Service object represents the service layer of the design architecture building
// most functionality as Methods off of the object. The ideology of the service layer
// is to perform data validation, data processing and business logic while leaving
// communication and data persistence/storage to other layers to be abstracted away.
type Service struct {
	Backends map[string]string
}

// NewService is the constructor for the Service object. First calling readArgs()
// to fill out the 'Backends' map, the function then returns the Service object.
func NewService() (*Service, error) {
	m, err := readArgs()
	if err != nil {
		return nil, err
	}
	return &Service{
		Backends: m,
	}, nil
}

// readArgs is used by the NewService constructor to parse the commandline arguments.
// Reading in addition backend addresses as arguments, validating their proper format
// and storing them as a map as a member variable for the Service object.
func readArgs() (map[string]string, error) {

	// Check for any arguments other than the run program command.
	params := os.Args[1:]
	m := make(map[string]string)
	if len(params) < 1 {
		return nil, errors.New("no valid parameters")
	}

	// Loop through the args validating their format and adding them to the map.
	for _, arg := range params {
		if string(arg[2]) != "=" {
			return nil, errors.New("invalid parameter format")
		}
		_, err := url.ParseRequestURI(arg[3:])
		if err != nil {
			return nil, errors.New("invalid parameter, must contain valid URL")
		}
		m[arg[:2]] = arg[3:]
	}

	return m, nil
}

// GetCompanyRequest is the data abstraction protocol for receiving the request from the transport layer.
type GetCompanyRequest struct {
	ID   string
	Code string
}

// GetCompanyResponse is the data abstraction protocol for sending the response to the transport layer.
type GetCompanyResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	ActiveUntil string `json:"active_until"`
}

// V1Response is the object into which the data is marshalled from the V1 backend response.
type V1Response struct {
	CompanyName string `json:"cn"`
	CreatedOn   string `json:"created_on"`
	ClosedOn    string `json:"closed_on"`
}

// V2Response is the object into which the data is marshalled from the V2 backend response.
type V2Response struct {
	CompanyName string `json:"company_name"`
	TaxID       string `json:"tin"`
	DissolvedOn string `json:"dissolved_on"`
}

// Validate is called on the GetCompanyRequest object to insure
// the ID and Code fields are present and Code has an exact length of 2.
func (c GetCompanyRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ID, validation.Required),
		validation.Field(&c.Code, validation.Required, validation.Length(2, 2)),
	)
}

// GetCompany retrieves company data from other backend services.
// On the Service layer it validates the incoming request, forms the backend
// request, processes the response and return it to the transport layer.
func (s *Service) GetCompany(req *GetCompanyRequest) (*GetCompanyResponse, error) {

	// validate incoming request object.
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// format outgoing request url.
	url := s.formRequest(req)

	// perform GET request to appropriate backend service.
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	// process the backend response.
	out, err := processResponse(resp)
	if err != nil {
		return nil, err
	}
	out.Id = req.ID

	return out, nil
}

// formRequest uses the incoming request from the transport layer and
// stored backend addresses to format the outgoing request url string.
// example: http://localhost:9001/companies/1234
func (s *Service) formRequest(req *GetCompanyRequest) string {
	return s.Backends[req.Code] + "/companies/" + req.ID
}

// processResponse parses the http response from the GET call to the backend. The
// function uses the response headers to determine the backend version and subsequent
// data type to marshal the response body and begins to format the Service layer response.
func processResponse(resp *http.Response) (*GetCompanyResponse, error) {

	// Initialize service layer response.
	out := &GetCompanyResponse{
		Active: false,
	}

	// Pull the content-type from the header and us it to determine V1 vs V2.
	contentType := resp.Header.Get("Content-Type")
	if contentType == "application/x-company-v1" {

		// Decode response into V1Response object.
		var obj V1Response
		err := json.NewDecoder(resp.Body).Decode(&obj)
		if err != nil {
			return nil, fmt.Errorf("unable to decode response: %w", err)
		}

		// Parse critical data and determine if the company is still "active"
		out.ActiveUntil = obj.ClosedOn
		active, err := isActive(out.ActiveUntil)
		if err != nil {
			return nil, err
		}
		out.Active = active
		out.Name = obj.CompanyName

		return out, nil
	} else if contentType == "application/x-company-v2" {

		// Decode response into V2Response object.
		var obj V2Response
		err := json.NewDecoder(resp.Body).Decode(&obj)
		if err != nil {
			return nil, fmt.Errorf("unable to decode response: %w", err)
		}

		// Parse critical data and determine if the company is still "active"
		out.ActiveUntil = obj.DissolvedOn
		active, err := isActive(out.ActiveUntil)
		if err != nil {
			return nil, err
		}
		out.Active = active
		out.Name = obj.CompanyName

		return out, nil

	} else {
		return nil, errors.New("invalid response")
	}
}

// isActive uses the closed/dissolved-on date compared against time.Now() to
// determine if the requested company is still active. Since these are optional
// fields it is assumed that if no date is provided the company is still considered.
func isActive(close string) (bool, error) {
	// If no date is provided company is "active".
	if close == "" {
		return true, nil
	}

	// Convert date string into time object.
	formattedDate, err := time.Parse(time.RFC3339, close)
	if err != nil {
		return false, fmt.Errorf("Error while parsing the date time: %w", err)
	}

	// If not past closed/dissolved date, company is still "active".
	if formattedDate.After(time.Now()) {
		return true, nil
	}
	return false, nil
}