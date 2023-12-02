// Code generated by ncraft. DO NOT EDIT.
// Rerunning ncraft will overwrite this file.
// Version: 0.1.0
// Version Date:

package svc

// This file provides server-side bindings for the HTTP transport.
// It utilizes the transport/http.Server.

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/gorilla/mux"
	"github.com/json-iterator/go"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
	"github.com/pkg/errors"

	httptransport "github.com/go-kit/kit/transport/http"
	mjhttp "github.com/mojo-lang/http/go/pkg/mojo/http"
	pagination "github.com/ncraft-io/ncraft-gokit/pkg/pagination"
	nhttp "github.com/ncraft-io/ncraft-gokit/pkg/transport/http"
	stdopentracing "github.com/opentracing/opentracing-go"

	"github.com/mojo-lang/core/go/pkg/mojo/core"

	"github.com/liankui/prenatal/go/pkg/prenatal"

	// this service api
	pb "github.com/liankui/prenatal/go/pkg/prenatal/v1"
)

const contentType = "application/json; charset=utf-8"

var (
	_ = fmt.Sprint
	_ = bytes.Compare
	_ = strconv.Atoi
	_ = httptransport.NewServer
	_ = ioutil.NopCloser
	_ = pb.NewQuizClient
	_ = io.Copy
	_ = errors.Wrap
	_ = mjhttp.UnmarshalQueryParam
)

var (
	_ = prenatal.Question{}
	_ = core.Null{}
	_ = prenatal.Answer{}
)

var cfg *nhttp.Config

func init() {
	cfg = nhttp.NewConfig()
}

// RegisterHttpHandler register a set of endpoints available on predefined paths to the router.
func RegisterHttpHandler(router *mux.Router, endpoints Endpoints, tracer stdopentracing.Tracer, logger log.Logger) {
	serverOptions := []httptransport.ServerOption{
		httptransport.ServerBefore(headersToContext, queryToContext),
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerAfter(httptransport.SetContentType(contentType)),
	}

	addTracerOption := func(methodName string) []httptransport.ServerOption {
		if tracer != nil {
			return append(serverOptions, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, methodName, logger)))
		}
		return serverOptions
	}

	router.Methods("POST").Path("/v1/questions").Handler(
		httptransport.NewServer(
			endpoints.CreateQuestionEndpoint,
			DecodeHTTPCreateQuestionZeroRequest,
			EncodeHTTPGenericResponse,
			addTracerOption("create_question")...,
		//append(serverOptions, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "create_question", logger)))...,
		))

	router.Methods("GET").Path("/v1/questions/{id}").Handler(
		httptransport.NewServer(
			endpoints.GetQuestionEndpoint,
			DecodeHTTPGetQuestionZeroRequest,
			EncodeHTTPGenericResponse,
			addTracerOption("get_question")...,
		//append(serverOptions, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "get_question", logger)))...,
		))

	router.Methods("PUT").Path("/v1/questions/{id}").Handler(
		httptransport.NewServer(
			endpoints.UpdateQuestionEndpoint,
			DecodeHTTPUpdateQuestionZeroRequest,
			EncodeHTTPGenericResponse,
			addTracerOption("update_question")...,
		//append(serverOptions, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "update_question", logger)))...,
		))

	router.Methods("DELETE").Path("/v1/questions/{id}").Handler(
		httptransport.NewServer(
			endpoints.DeleteQuestionEndpoint,
			DecodeHTTPDeleteQuestionZeroRequest,
			EncodeHTTPGenericResponse,
			addTracerOption("delete_question")...,
		//append(serverOptions, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "delete_question", logger)))...,
		))

	router.Methods("POST").Path("/v1/Answer").Handler(
		httptransport.NewServer(
			endpoints.CreateAnswerEndpoint,
			DecodeHTTPCreateAnswerZeroRequest,
			EncodeHTTPGenericResponse,
			addTracerOption("create_answer")...,
		//append(serverOptions, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "create_answer", logger)))...,
		))
}

// ErrorEncoder writes the error to the ResponseWriter, by default a content
// type of application/json, a body of json with key "error" and the value
// error.Error(), and a status code of 500. If the error implements Headerer,
// the provided headers will be applied to the response. If the error
// implements json.Marshaler, and the marshaling succeeds, the JSON encoded
// form of the error will be used. If the error implements StatusCoder, the
// provided StatusCode will be used instead of 500.
func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	e, ok := err.(*core.Error)
	if !ok {
		e = core.NewErrorFrom(500, err.Error())
	}

	w.Header().Set("Content-Type", contentType)
	if headerer, ok := err.(httptransport.Headerer); ok {
		for k := range headerer.Headers() {
			w.Header().Set(k, headerer.Headers().Get(k))
		}
	}

	var body []byte
	code := http.StatusInternalServerError

	if enveloped := nhttp.IsEnvelopeStyle(ctx, cfg.GetStyle()); enveloped {
		envelope := &nhttp.EnvelopedResponse{}
		envelope.Error = e

		if cfg.GetEnvelop().MappingCode {
			if sc, ok := err.(httptransport.StatusCoder); ok {
				code = sc.StatusCode()
			}
		} else {
			code = http.StatusOK
		}

		var response interface{}
		if cfg.GetEnvelop().ErrorWrapped {
			response = envelope.ToErrorWrapped()
		} else {
			response = envelope
		}

		jsonBody, marshalErr := jsoniter.ConfigFastest.Marshal(response)
		if marshalErr != nil {
			logs.Warnw("failed to marshal the error response to json", "error", marshalErr)
		} else {
			body = jsonBody
		}
	} else {
		if sc, ok := err.(httptransport.StatusCoder); ok {
			code = sc.StatusCode()
		}
		if marshaler, ok := err.(json.Marshaler); ok {
			if jsonBody, marshalErr := marshaler.MarshalJSON(); marshalErr == nil {
				body = jsonBody
			}
		}

		if jsonBody, marshalErr := jsoniter.ConfigFastest.Marshal(e); marshalErr == nil {
			body = jsonBody
		}
	}

	w.WriteHeader(code)
	w.Write(body)
}

// Server Decode

// DecodeHTTPCreateQuestionZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded create_question request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPCreateQuestionZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.CreateQuestionRequest

	// to support gzip input
	var reader io.ReadCloser
	var err error
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		defer reader.Close()
		if err != nil {
			return nil, nhttp.WrapError(err, 400, "failed to read the gzip content")
		}
	default:
		reader = r.Body
	}

	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, nhttp.WrapError(err, 400, "cannot read body of http request")
	}
	if len(buf) > 0 {
		req.Question = &prenatal.Question{}
		if err = jsoniter.ConfigFastest.Unmarshal(buf, req.Question); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, nhttp.WrapError(err,
				http.StatusBadRequest,
				fmt.Sprintf("request body '%s': cannot parse non-json request body", buf),
			)
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := core.NewUrlQueryFrom(r.URL.Query())
	_ = queryParams

	parsedQueryParams := make(map[string]bool)
	_ = parsedQueryParams

	questionInitialized := false
	if req.Question == nil {
		questionInitialized = true
		req.Question = &prenatal.Question{}
	}
	err = mjhttp.UnmarshalQueryParam(queryParams, req.Question, "question")
	if err != nil {
		if core.IsNotFoundError(err) {
			if questionInitialized {
				req.Question = nil
			}
		} else {
			return nil, nhttp.WrapError(err, 400, "cannot unmarshal the question  query parameter")
		}
	}

	return &req, nil
}

// DecodeHTTPGetQuestionZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded get_question request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPGetQuestionZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.GetQuestionRequest

	// to support gzip input
	var reader io.ReadCloser
	var err error
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		defer reader.Close()
		if err != nil {
			return nil, nhttp.WrapError(err, 400, "failed to read the gzip content")
		}
	default:
		reader = r.Body
	}

	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, nhttp.WrapError(err, 400, "cannot read body of http request")
	}
	if len(buf) > 0 {
		if err = jsoniter.ConfigFastest.Unmarshal(buf, &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, nhttp.WrapError(err,
				http.StatusBadRequest,
				fmt.Sprintf("request body '%s': cannot parse non-json request body", buf),
			)
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := core.NewUrlQueryFrom(r.URL.Query())
	_ = queryParams

	parsedQueryParams := make(map[string]bool)
	_ = parsedQueryParams

	err = mjhttp.UnmarshalPathParam(pathParams, &req.Id, "id")
	if err != nil && !core.IsNotFoundError(err) {
		return nil, nhttp.WrapError(err, 400, "cannot unmarshal the id  query parameter")
	}

	return &req, nil
}

// DecodeHTTPUpdateQuestionZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded update_question request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPUpdateQuestionZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.UpdateQuestionRequest

	// to support gzip input
	var reader io.ReadCloser
	var err error
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		defer reader.Close()
		if err != nil {
			return nil, nhttp.WrapError(err, 400, "failed to read the gzip content")
		}
	default:
		reader = r.Body
	}

	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, nhttp.WrapError(err, 400, "cannot read body of http request")
	}
	if len(buf) > 0 {
		req.Question = &prenatal.Question{}
		if err = jsoniter.ConfigFastest.Unmarshal(buf, req.Question); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, nhttp.WrapError(err,
				http.StatusBadRequest,
				fmt.Sprintf("request body '%s': cannot parse non-json request body", buf),
			)
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := core.NewUrlQueryFrom(r.URL.Query())
	_ = queryParams

	parsedQueryParams := make(map[string]bool)
	_ = parsedQueryParams

	questionInitialized := false
	if req.Question == nil {
		questionInitialized = true
		req.Question = &prenatal.Question{}
	}
	err = mjhttp.UnmarshalQueryParam(queryParams, req.Question, "question")
	if err != nil {
		if core.IsNotFoundError(err) {
			if questionInitialized {
				req.Question = nil
			}
		} else {
			return nil, nhttp.WrapError(err, 400, "cannot unmarshal the question  query parameter")
		}
	}

	return &req, nil
}

// DecodeHTTPDeleteQuestionZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded delete_question request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPDeleteQuestionZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.DeleteQuestionRequest

	// to support gzip input
	var reader io.ReadCloser
	var err error
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		defer reader.Close()
		if err != nil {
			return nil, nhttp.WrapError(err, 400, "failed to read the gzip content")
		}
	default:
		reader = r.Body
	}

	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, nhttp.WrapError(err, 400, "cannot read body of http request")
	}
	if len(buf) > 0 {
		if err = jsoniter.ConfigFastest.Unmarshal(buf, &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, nhttp.WrapError(err,
				http.StatusBadRequest,
				fmt.Sprintf("request body '%s': cannot parse non-json request body", buf),
			)
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := core.NewUrlQueryFrom(r.URL.Query())
	_ = queryParams

	parsedQueryParams := make(map[string]bool)
	_ = parsedQueryParams

	err = mjhttp.UnmarshalPathParam(pathParams, &req.Id, "id")
	if err != nil && !core.IsNotFoundError(err) {
		return nil, nhttp.WrapError(err, 400, "cannot unmarshal the id  query parameter")
	}

	return &req, nil
}

// DecodeHTTPCreateAnswerZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded create_answer request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPCreateAnswerZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.CreateAnswerRequest

	// to support gzip input
	var reader io.ReadCloser
	var err error
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		defer reader.Close()
		if err != nil {
			return nil, nhttp.WrapError(err, 400, "failed to read the gzip content")
		}
	default:
		reader = r.Body
	}

	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, nhttp.WrapError(err, 400, "cannot read body of http request")
	}
	if len(buf) > 0 {
		req.Answer = &prenatal.Answer{}
		if err = jsoniter.ConfigFastest.Unmarshal(buf, req.Answer); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, nhttp.WrapError(err,
				http.StatusBadRequest,
				fmt.Sprintf("request body '%s': cannot parse non-json request body", buf),
			)
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := core.NewUrlQueryFrom(r.URL.Query())
	_ = queryParams

	parsedQueryParams := make(map[string]bool)
	_ = parsedQueryParams

	answerInitialized := false
	if req.Answer == nil {
		answerInitialized = true
		req.Answer = &prenatal.Answer{}
	}
	err = mjhttp.UnmarshalQueryParam(queryParams, req.Answer, "answer")
	if err != nil {
		if core.IsNotFoundError(err) {
			if answerInitialized {
				req.Answer = nil
			}
		} else {
			return nil, nhttp.WrapError(err, 400, "cannot unmarshal the answer  query parameter")
		}
	}

	return &req, nil
}

// EncodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if writer, ok := response.(nhttp.ResponseWriter); ok {
		return writer.WriteHttpResponse(ctx, w)
	}

	if reflect.ValueOf(response).IsNil() {
		response = nil
	}
	if _, ok := response.(*core.Null); ok {
		response = nil
	}

	enveloped := nhttp.IsEnvelopeStyle(ctx, cfg.GetStyle())
	if enveloped {
		code := core.NewErrorCode(200)
		message := "OK"
		if sc := cfg.GetEnvelop().SuccessCode; len(sc) > 0 {
			if c, err := core.ParseErrorCode(sc); err != nil {
				logs.Warnw("failed to parse the user setting success code, will use \"200\" indeed.", "code", sc, "error", err)
			} else {
				code = c
			}
		}
		if msg := cfg.GetEnvelop().SuccessMessage; len(msg) > 0 {
			message = msg
		}

		totalCount := int32(0)
		nextPageToken := ""
		if p, ok := response.(pagination.Paginater); ok {
			totalCount = p.GetTotalCount()
			nextPageToken = p.GetNextPageToken()
		}

		response = &nhttp.EnvelopedResponse{
			Error: &core.Error{
				Code:    code,
				Message: message,
			},
			TotalCount:    totalCount,
			NextPageToken: nextPageToken,
			Data:          response,
		}
	}

	if response == nil {
		return nil
	}

	if p, ok := response.(pagination.Paginater); ok && !enveloped {
		total := p.GetTotalCount()
		if total > 0 {
			w.Header().Set("X-Total-Count", strconv.Itoa(int(total)))
		}

		next := p.GetNextPageToken()
		if len(next) > 0 {
			path, _ := ctx.Value("http-request-path").(string)
			if len(path) == 0 {
				path = "/?next-page-token=" + next
			} else {
				url, _ := core.ParseUrl(path)
				url.Query.Add("next-page-token", next)
				path = url.Format()
			}
			w.Header().Set("Link", fmt.Sprintf("<%s>; rel=\"next\"", path))
		}
	}

	return nhttp.NewResponseJsonWriter(response).WriteHttpResponse(ctx, w)
}

// Helper functions

func headersToContext(ctx context.Context, r *http.Request) context.Context {
	for k, _ := range r.Header {
		// The key is added both in http format (k) which has had
		// http.CanonicalHeaderKey called on it in transport as well as the
		// strings.ToLower which is the grpc metadata format of the key so
		// that it can be accessed in either format
		ctx = context.WithValue(ctx, k, r.Header.Get(k))
		ctx = context.WithValue(ctx, strings.ToLower(k), r.Header.Get(k))
	}

	// add the access key to context
	accessKey := r.URL.Query().Get("access_key")
	if len(accessKey) > 0 {
		ctx = context.WithValue(ctx, "access_key", accessKey)
	}

	// Tune specific change.
	// also add the request url
	ctx = context.WithValue(ctx, "http-request-path", r.URL.Path)
	ctx = context.WithValue(ctx, "transport", "HTTPJSON")

	return ctx
}

func queryToContext(ctx context.Context, r *http.Request) context.Context {
	check := func(values []string) bool {
		for _, value := range values {
			if value == "true" {
				return true
			}
		}
		return false
	}
	for key, values := range r.URL.Query() {
		switch key {
		case "envelope":
			if check(values) {
				ctx = context.WithValue(ctx, "envelope", true)
			}
		}
	}
	return ctx
}
