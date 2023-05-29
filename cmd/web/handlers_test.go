package main

import (
	"net/http"
	"testing"

	"github.com/sxc/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplicaiton(t)
	// app := &application{
	// 	errorLog: log.New(io.Discard, "", 0),
	// 	infoLog:  log.New(io.Discard, "", 0),
	// }

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

// 	ts := httptest.NewTLSServer(app.routes())
// 	defer ts.Close()

// 	rs, err := ts.Client().Get(ts.URL + "/ping")

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	assert.Equal(t, rs.StatusCode, http.StatusOK)

// 	defer rs.Body.Close()

// 	body, err := io.ReadAll(rs.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	bytes.TrimSpace(body)

// 	assert.Equal(t, string(body), "OK")
// }
