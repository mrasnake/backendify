package main

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

// The Server object represents the transport layer of the design architecture building
// most functionality as Methods off of the object. The ideology of the transport layer
// is to abstract away all communication protocol from the service layer, parsing the incoming
// request, packageing it for the service layer, formatting and writing the outgoing response.
type Server struct {
	Service *Service
	Router  *mux.Router
}

// NewServer is the constructor for the Server object. After accepting the Service
// object the Server is then initialized to define the router before being returned.
func NewServer(srvc *Service) *Server {
	out := &Server{
		Service: srvc,
	}
	out.InitServer()
	return out
}

// GetStatus is a health check endpoint, the purpose is simply to
// respond with status code 200 to confirm the service is up and running.
func (s *Server) GetStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// GetCompany retrieves company data from other backend services.
// On the Transport layer it parses the query data from the URL,
// makes a call to the service layer, marshals the response into JSON,
// then finally writes the response to the ResponseWriter.
func (s *Server) GetCompany(w http.ResponseWriter, r *http.Request) {

	// Parse the id and country code from the URL Query String.
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	id := params.Get("id")
	iso := params.Get("country_iso")

	// Format data into a service layer request.
	req := &GetCompanyRequest{
		ID:   id,
		Code: iso,
	}

	// Call the service layer.
	resp, err := s.Service.GetCompany(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Marshal the service response into JSON.
	out, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Write the JSON object to the ResponseWriter.
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

	return
}