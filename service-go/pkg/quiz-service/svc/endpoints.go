// Code generated by ncraft. DO NOT EDIT.
// Rerunning ncraft will overwrite this file.
// Version: 0.1.0
// Version Date:

package svc

// This file contains methods to make individual endpoints from services,
// request and response types to serve those endpoints, as well as encoders and
// decoders for those types, for all of our supported transport serialization
// formats.

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	"github.com/liankui/prenatal/go/pkg/prenatal"
	"github.com/mojo-lang/core/go/pkg/mojo/core"

	// this service api
	pb "github.com/liankui/prenatal/go/pkg/prenatal/v1"
)

var (
	_ = prenatal.Question{}
	_ = core.Null{}
	_ = prenatal.Answer{}
)

// Endpoints collects all of the endpoints that compose an add service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	CreateQuestionEndpoint endpoint.Endpoint
	GetQuestionEndpoint    endpoint.Endpoint
	UpdateQuestionEndpoint endpoint.Endpoint
	DeleteQuestionEndpoint endpoint.Endpoint
	CreateAnswerEndpoint   endpoint.Endpoint
}

// Endpoints

func (e Endpoints) CreateQuestion(ctx context.Context, in *pb.CreateQuestionRequest) (*prenatal.Question, error) {
	response, err := e.CreateQuestionEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*prenatal.Question), nil
}

func (e Endpoints) GetQuestion(ctx context.Context, in *pb.GetQuestionRequest) (*prenatal.Question, error) {
	response, err := e.GetQuestionEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*prenatal.Question), nil
}

func (e Endpoints) UpdateQuestion(ctx context.Context, in *pb.UpdateQuestionRequest) (*prenatal.Question, error) {
	response, err := e.UpdateQuestionEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*prenatal.Question), nil
}

func (e Endpoints) DeleteQuestion(ctx context.Context, in *pb.DeleteQuestionRequest) (*core.Null, error) {
	response, err := e.DeleteQuestionEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*core.Null), nil
}

func (e Endpoints) CreateAnswer(ctx context.Context, in *pb.CreateAnswerRequest) (*prenatal.Answer, error) {
	response, err := e.CreateAnswerEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*prenatal.Answer), nil
}

// Make Endpoints

func MakeCreateQuestionEndpoint(s pb.QuizServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.CreateQuestionRequest)
		v, err := s.CreateQuestion(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeGetQuestionEndpoint(s pb.QuizServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.GetQuestionRequest)
		v, err := s.GetQuestion(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeUpdateQuestionEndpoint(s pb.QuizServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UpdateQuestionRequest)
		v, err := s.UpdateQuestion(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeDeleteQuestionEndpoint(s pb.QuizServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.DeleteQuestionRequest)
		v, err := s.DeleteQuestion(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeCreateAnswerEndpoint(s pb.QuizServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.CreateAnswerRequest)
		v, err := s.CreateAnswer(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// WrapAllExcept wraps each Endpoint field of struct Endpoints with a
// go-kit/kit/endpoint.Middleware.
// Use this for applying a set of middlewares to every endpoint in the service.
// Optionally, endpoints can be passed in by name to be excluded from being wrapped.
// WrapAllExcept(middleware, "Status", "Ping")
func (e *Endpoints) WrapAllExcept(middleware endpoint.Middleware, excluded ...string) {
	included := map[string]struct{}{
		"create_question": struct{}{},
		"get_question":    struct{}{},
		"update_question": struct{}{},
		"delete_question": struct{}{},
		"create_answer":   struct{}{},
	}

	for _, ex := range excluded {
		if _, ok := included[ex]; !ok {
			panic(fmt.Sprintf("Excluded endpoint '%s' does not exist; see middlewares/endpoints.go", ex))
		}
		delete(included, ex)
	}

	for inc, _ := range included {
		if inc == "create_question" {
			e.CreateQuestionEndpoint = middleware(e.CreateQuestionEndpoint)
		}
		if inc == "get_question" {
			e.GetQuestionEndpoint = middleware(e.GetQuestionEndpoint)
		}
		if inc == "update_question" {
			e.UpdateQuestionEndpoint = middleware(e.UpdateQuestionEndpoint)
		}
		if inc == "delete_question" {
			e.DeleteQuestionEndpoint = middleware(e.DeleteQuestionEndpoint)
		}
		if inc == "create_answer" {
			e.CreateAnswerEndpoint = middleware(e.CreateAnswerEndpoint)
		}
	}
}

// LabeledMiddleware will get passed the endpoint name when passed to
// WrapAllLabeledExcept, this can be used to write a generic metrics
// middleware which can send the endpoint name to the metrics collector.
type LabeledMiddleware func(string, endpoint.Endpoint) endpoint.Endpoint

// WrapAllLabeledExcept wraps each Endpoint field of struct Endpoints with a
// LabeledMiddleware, which will receive the name of the endpoint. See
// LabeldMiddleware. See method WrapAllExept for details on excluded
// functionality.
func (e *Endpoints) WrapAllLabeledExcept(middleware func(string, endpoint.Endpoint) endpoint.Endpoint, excluded ...string) {
	included := map[string]struct{}{
		"create_question": struct{}{},
		"get_question":    struct{}{},
		"update_question": struct{}{},
		"delete_question": struct{}{},
		"create_answer":   struct{}{},
	}

	for _, ex := range excluded {
		if _, ok := included[ex]; !ok {
			panic(fmt.Sprintf("Excluded endpoint '%s' does not exist; see middlewares/endpoints.go", ex))
		}
		delete(included, ex)
	}

	for inc, _ := range included {
		if inc == "create_question" {
			e.CreateQuestionEndpoint = middleware("create_question", e.CreateQuestionEndpoint)
		}
		if inc == "get_question" {
			e.GetQuestionEndpoint = middleware("get_question", e.GetQuestionEndpoint)
		}
		if inc == "update_question" {
			e.UpdateQuestionEndpoint = middleware("update_question", e.UpdateQuestionEndpoint)
		}
		if inc == "delete_question" {
			e.DeleteQuestionEndpoint = middleware("delete_question", e.DeleteQuestionEndpoint)
		}
		if inc == "create_answer" {
			e.CreateAnswerEndpoint = middleware("create_answer", e.CreateAnswerEndpoint)
		}
	}
}
