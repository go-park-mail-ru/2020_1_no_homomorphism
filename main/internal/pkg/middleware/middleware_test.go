package middleware

import (
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_main/logger"
	"github.com/bmizerany/assert"
	"github.com/steinfletcher/apitest"
	"net/http"
	"os"
	"testing"
)

type TestStruct struct {
	id    string
	start uint64
	end   uint64
	t     *testing.T
}

func (ts TestStruct) MockHandler(w http.ResponseWriter, r *http.Request) {
	ctxID, ok := r.Context().Value("id").(string)
	assert.T(ts.t, ok, true)
	assert.Equal(ts.t, ts.id, ctxID)

	ctxStart, ok := r.Context().Value("start").(uint64)
	assert.T(ts.t, ok, true)
	assert.Equal(ts.t, ts.start, ctxStart)

	ctxEnd, ok := r.Context().Value("end").(uint64)
	assert.T(ts.t, ok, true)
	assert.Equal(ts.t, ts.end, ctxEnd)
}

func TestBoundedVars(t *testing.T) {
	t.Run("BoundedVars-OK", func(t *testing.T) {
		id := "1"

		var start uint64 = 2
		var end uint64 = 3

		ts := TestStruct{id, start, end, t}

		boundedVars := BoundedVars(ts.MockHandler, logger.NewLogger(os.Stdout))

		handler := SetTripleVars(
			boundedVars,
			ts.id,
			fmt.Sprint(ts.start),
			fmt.Sprint(ts.end),
		)

		apitest.New("BoundedVars-OK").
			Handler(handler).
			Method("Get").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("BoundedVars-NoVars", func(t *testing.T) {

		ts := TestStruct{}

		handler := BoundedVars(ts.MockHandler, logger.NewLogger(os.Stdout))

		apitest.New("BoundedVars-NoVars").
			Handler(handler).
			Method("Get").
			URL("/users/albums").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("BoundedVars-FailedToParse", func(t *testing.T) {

		ts := TestStruct{}

		boundedVars := BoundedVars(ts.MockHandler, logger.NewLogger(os.Stdout))
		handler := SetTripleVars(boundedVars, "no int", "2", "no int too")

		apitest.New("BoundedVars-FailedToParse").
			Handler(handler).
			Method("Get").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

