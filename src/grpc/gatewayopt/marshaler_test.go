package gatewayopt

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	bookpb "safeweb.app/grpc/gatewayopt/internal/example"
)

func TestJSONPb_Marshal(t *testing.T) {
	marshaler := JSONPbMarshaler{
		JSONPb:                  &runtime.JSONPb{EmitDefaults: true},
		removeOriRespDataFields: true,
	}

	books := []*bookpb.Book{
		{
			Id:       1,
			Name:     "Book 1",
			AuthorId: 1,
		},
		{
			Id:       2,
			Name:     "Book 2",
			AuthorId: 2,
		},
	}

	user := &bookpb.User{
		Id:   1,
		Name: "Writer 1",
	}

	cases := []struct {
		Name   string
		Proto  proto.Message
		Golden string
	}{
		{
			Name: "Data contains 2 fields",
			Proto: &bookpb.ListBooksWithAuthorResponse{
				Books:   books,
				Author:  user,
				Message: "A best-seller book.",
			},
			Golden: "ListBooksWithAuthorResponse",
		},
		{
			Name: "Data contain 1 field which is an array",
			Proto: &bookpb.ListBooksResponse{
				Books:   books,
				Message: "A list of books",
			},
			Golden: "ListBooksResponse",
		},
		{
			Name: "Data contain 1 field which is an object",
			Proto: &bookpb.GetBookResponse{
				Book:    books[0],
				Message: "A book",
			},
			Golden: "GetBookResponse",
		},
		{
			Name: "Data contains 2 fields, one is empty object",
			Proto: &bookpb.ListBooksWithAuthorResponse{
				Books:   books,
				Author:  nil,
				Message: "A best-seller book.",
			},
			Golden: "ListBooksWithAuthorResponse_Empty",
		},
		{
			Name: "Data is an empty object",
			Proto: &bookpb.GetBookResponse{
				Book:    nil,
				Message: "A book",
			},
			Golden: "GetBookResponse_Empty",
		},
		{
			Name: "Data is an empty array",
			Proto: &bookpb.ListBooksResponse{
				Books:   nil,
				Message: "A book",
			},
			Golden: "ListBooksResponse_Empty",
		},
	}

	g := goldie.New(t)
	for _, tc := range cases {
		bs, err := marshaler.Marshal(tc.Proto)
		assert.Nil(t, err, "Cannot marshal data (%s)", tc.Name)
		g.Assert(t, tc.Golden, bs)
	}
}

func TestJSONPb_Marshal_NotRemoveOriRespDataFields(t *testing.T) {
	marshaler := JSONPbMarshaler{
		JSONPb:                  &runtime.JSONPb{EmitDefaults: true},
		removeOriRespDataFields: false,
	}

	books := []*bookpb.Book{
		{
			Id:       1,
			Name:     "Book 1",
			AuthorId: 1,
		},
		{
			Id:       2,
			Name:     "Book 2",
			AuthorId: 2,
		},
	}

	user := &bookpb.User{
		Id:   1,
		Name: "Writer 1",
	}

	cases := []struct {
		Name   string
		Proto  proto.Message
		Golden string
	}{
		{
			Name: "Data contains 2 fields",
			Proto: &bookpb.ListBooksWithAuthorResponse{
				Books:   books,
				Author:  user,
				Message: "A best-seller book.",
			},
			Golden: "ListBooksWithAuthorResponse_WithOrig",
		},
		{
			Name: "Data contain 1 field which is an array",
			Proto: &bookpb.ListBooksResponse{
				Books:   books,
				Message: "A list of books",
			},
			Golden: "ListBooksResponse_WithOrig",
		},
		{
			Name: "Data contain 1 field which is an object",
			Proto: &bookpb.GetBookResponse{
				Book:    books[0],
				Message: "A book",
			},
			Golden: "GetBookResponse_WithOrig",
		},
		{
			Name: "Data contains 2 fields, one is empty object",
			Proto: &bookpb.ListBooksWithAuthorResponse{
				Books:   books,
				Author:  nil,
				Message: "A best-seller book.",
			},
			Golden: "ListBooksWithAuthorResponse_Empty_WithOrig",
		},
		{
			Name: "Data is an empty object",
			Proto: &bookpb.GetBookResponse{
				Book:    nil,
				Message: "A book",
			},
			Golden: "GetBookResponse_Empty_WithOrig",
		},
		{
			Name: "Data is an empty array",
			Proto: &bookpb.ListBooksResponse{
				Books:   nil,
				Message: "A book",
			},
			Golden: "ListBooksResponse_Empty_WithOrig",
		},
	}

	g := goldie.New(t)
	for _, tc := range cases {
		bs, err := marshaler.Marshal(tc.Proto)
		assert.Nil(t, err, "Cannot marshal data (%s)", tc.Name)
		g.Assert(t, tc.Golden, bs)
	}
}
