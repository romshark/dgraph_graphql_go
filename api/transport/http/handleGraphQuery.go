package http

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/api/graph"
)

var graphRespHead = []byte(`{"data":`)
var graphRespTail = []byte(`}`)

type graphResponseError struct {
	Code    string `json:"c"`
	Message string `json:"m"`
}

type graphResponse struct {
	Data  []byte              `json:"data"`
	Error *graphResponseError `json:"errors"`
}

// graphQuery represents the JSON graph query structure
type graphQuery struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

// handleGraphQuery handles a graph query request
func (t *Server) handleGraphQuery(
	resp http.ResponseWriter,
	req *http.Request,
) {
	handleUnexpectedErr := func(err error, logErr bool) {
		http.Error(
			resp,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		if logErr {
			t.errorLog.Print(err)
		}
	}

	// Decode graph query
	requestDecoderJSON := json.NewDecoder(req.Body)
	var graphQuery graphQuery
	if err := requestDecoderJSON.Decode(&graphQuery); err != nil {
		handleUnexpectedErr(errors.Wrap(err, "graph query JSON decode"), true)
		return
	}

	response, err := t.onGraphQuery(
		req.Context(),
		graph.Query{
			Query:         graphQuery.Query,
			OperationName: graphQuery.OperationName,
			Variables:     graphQuery.Variables,
		},
	)
	if err != nil {
		handleUnexpectedErr(err, false)
	}

	jsonEncoder := json.NewEncoder(resp)

	if response.Error != nil {
		// User error
		resp.WriteHeader(http.StatusBadRequest)

		if err := jsonEncoder.Encode(graphResponse{
			Error: &graphResponseError{
				Code:    response.Error.Code,
				Message: response.Error.Message,
			},
		}); err != nil {
			handleUnexpectedErr(
				errors.Wrap(err, "graph response JSON encode"),
				true,
			)
		}
		return
	}

	// Reply successfully
	resp.Header().Set("Content-Type", "application/json")

	if _, err := resp.Write(graphRespHead); err != nil {
		handleUnexpectedErr(errors.Wrap(err, "response head write"), true)
		return
	}
	if _, err := resp.Write(response.Data); err != nil {
		handleUnexpectedErr(errors.Wrap(err, "response data write"), true)
		return
	}
	if _, err := resp.Write(graphRespTail); err != nil {
		handleUnexpectedErr(errors.Wrap(err, "response tail write"), true)
		return
	}
}
