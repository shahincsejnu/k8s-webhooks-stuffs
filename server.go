package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//Server contains the functions handling server requests
type Server struct {
	ServerTLSConf *tls.Config
	ClientTLSConf *tls.Config
	CaPEM         []byte
}

func (s Server) getCA(w http.ResponseWriter, req *http.Request) {
	if len(s.CaPEM) == 0 {
		fmt.Fprintf(w, "No certificate found\n")
		return
	}

	// if base64 parameter is set, return in base64 format
	req.ParseForm()
	if _, hasParam := req.Form["base64"]; hasParam {
		fmt.Fprintf(w, string(base64.StdEncoding.EncodeToString(s.CaPEM)))
		return
	}

	fmt.Fprintf(w, string(s.CaPEM))
}

func (s Server) postMutator(w http.ResponseWriter, r *http.Request) {
	fmt.Println("mutator webhook got Called")
	var request AdmissionReviewRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON body in invalid format: %s\n", err.Error()), http.StatusBadRequest)
		return
	}
	if request.APIVersion != "admission.k8s.io/v1" || request.Kind != "AdmissionReview" {
		http.Error(w, fmt.Sprintf("wrong APIVersion or kind: %s - %s", request.APIVersion, request.Kind), http.StatusBadRequest)
		return

	}
	fmt.Printf("debug: %+v\n", request.Request)
	response := AdmissionReviewResponse{
		APIVersion: "admission.k8s.io/v1",
		Kind:       "AdmissionReview",
		Response: Response{
			UID:     request.Request.UID,
			Allowed: true,
		},
	}

	// add label if we're creating a pod
	if request.Request.Kind.Group == "shahin.oka.com" && request.Request.Kind.Version == "v1alpha1" && request.Request.Kind.Kind == "Teployment" && request.Request.Operation == "CREATE" {
		// if there is a path till /metadata/labels this one will work
		patch := `[{"op": "add", "path": "/metadata/labels/myExtraLabel", "value": "webhook-was-here"}]`
		// if need also create the labels then add myExtraLabel then this one
		//patch := `[{"op": "add", "path": "/metadata/labels", "value": {"myExtraLabel": "webhook-was-here"}}]`
		patchEnc := base64.StdEncoding.EncodeToString([]byte(patch))
		response.Response.PatchType = "JSONPatch"
		response.Response.Patch = patchEnc
	}

	out, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON output marshal error: %s\n", err.Error()), http.StatusBadRequest)
		return
	}
	fmt.Printf("Got request, response: %s\n", string(out))

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		log.Fatal(err)
	}
}

// Validate handler accepts or rejects based on request contents
func (s Server) postValidator(w http.ResponseWriter, r *http.Request) {
	fmt.Println("validator webhook got called")
	var request AdmissionReviewRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON body in invalid format: %s\n", err.Error()), http.StatusBadRequest)
		return
	}
	if request.APIVersion != "admission.k8s.io/v1" || request.Kind != "AdmissionReview" {
		http.Error(w, fmt.Sprintf("wrong APIVersion or kind: %s - %s", request.APIVersion, request.Kind), http.StatusBadRequest)
		return

	}

	fmt.Printf("debug: %+v\n", request.Request)
	response := AdmissionReviewResponse{
		APIVersion: "admission.k8s.io/v1",
		Kind:       "AdmissionReview",
		Response: Response{
			UID:     request.Request.UID,
			Allowed: true,
		},
	}

	fmt.Println(request)

	if len(request.Request.Metadata.Labels) == 0 {
		response.Response.Allowed = false
	}

	out, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON output marshal error: %s\n", err.Error()), http.StatusBadRequest)
		return
	}
	fmt.Printf("Got request, response: %s\n", string(out))

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		log.Fatal(err)
	}
}