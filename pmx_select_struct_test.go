package pmx_test

import (
	"context"
	"testing"

	"github.com/iimrudy/pmx"
	"github.com/iimrudy/pmx/test"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
)

type SelectStructSuite struct {
	suite.Suite
	conn *pgx.Conn
}

func (s *SelectStructSuite) SetupSuite() {
	s.conn = test.Connect(context.Background())
}

func (s *SelectStructSuite) TearDownSuite() {
	s.NoError(s.conn.Close(context.Background()))
}

func TestSelectStruct(t *testing.T) {
	suite.Run(t, new(SelectStructSuite))
}

func (s *SelectStructSuite) TestStructPointer() {
	var projection test.Projection
	err := pmx.Select(context.Background(), s.conn, &projection,
		"select $1 as id, $2 as name, $3::jsonb as metadata, $4::jsonb as slice",
		"projection-id",
		"projection-name",
		map[string]int{"index": 1},
		[]string{"value"},
	)
	s.Equal(test.Projection{
		ID:       "projection-id",
		Name:     "projection-name",
		Metadata: map[string]int{"index": 1},
		Slice:    []string{"value"},
	}, projection)
	s.NoError(err)
}

func (s *SelectStructSuite) TestNull() {
	var projection test.Projection
	err := pmx.Select(context.Background(), s.conn, &projection, "select $1 as id, $2 as name", "projection-id", nil)
	s.Equal(test.Projection{ID: "projection-id"}, projection)
	s.NoError(err)
}

func (s *SelectStructSuite) TestUnmapped() {
	var projection test.Projection
	err := pmx.Select(context.Background(), s.conn, &projection, "select $1 as id, $2 as unmapped", "projection-id", "x")
	s.Equal(test.Projection{ID: "projection-id"}, projection)
	s.NoError(err)
}

func (s *SelectStructSuite) TestNoRows() {
	var projection test.Projection
	err := pmx.Select(context.Background(), s.conn, &projection, "select 1 limit 0")
	s.ErrorIs(err, pmx.ErrNoRows)
}

func (s *SelectStructSuite) TestStructValue() {
	var projection test.Projection
	err := pmx.Select(context.Background(), s.conn, projection, "select 1")
	s.ErrorIs(err, pmx.ErrInvalidRef)
}

func (s *SelectStructSuite) TestMapPointer() {
	projection := map[string]string{}
	err := pmx.Select(context.Background(), s.conn, &projection, "select 1")
	s.ErrorIs(err, pmx.ErrInvalidRef)
}

func (s *SelectStructSuite) TestMapValue() {
	projection := map[string]string{}
	err := pmx.Select(context.Background(), s.conn, projection, "select 1")
	s.ErrorIs(err, pmx.ErrInvalidRef)
}

func (s *SelectStructSuite) TestEmbeddedStruct() {
	liked := true
	var event test.EnrichedEvent
	err := pmx.Select(context.Background(), s.conn, &event,
		"select $1::bigint as position, $2::text as recorded_by, $3::boolean as liked, $4::bigint as likes, $5::bigint as views",
		int64(1), "user-1", true, uint64(10), uint64(100),
	)
	s.NoError(err)
	s.Equal(int64(1), event.Position)
	s.Equal("user-1", event.RecordedBy)
	s.Equal(&liked, event.Liked)
	s.Equal(uint64(10), event.Likes)
	s.Equal(uint64(100), event.Views)
}

func (s *SelectStructSuite) TestEmbeddedStructNull() {
	var event test.EnrichedEvent
	err := pmx.Select(context.Background(), s.conn, &event,
		"select $1::bigint as position, $2::text as recorded_by",
		int64(5), nil,
	)
	s.NoError(err)
	s.Equal(int64(5), event.Position)
	s.Nil(event.Liked)
	s.Equal(uint64(0), event.Likes)
}

func (s *SelectStructSuite) TestEmbeddedStructUnmapped() {
	var event test.EnrichedEvent
	err := pmx.Select(context.Background(), s.conn, &event,
		"select $1::bigint as position, $2::text as unmapped",
		int64(1), "x",
	)
	s.NoError(err)
	s.Equal(int64(1), event.Position)
	s.Equal("", event.RecordedBy)
}
