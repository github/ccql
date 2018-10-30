package sql

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestParseQueriesSingle(t *testing.T) {
	queriesText := `select 1`
	queries, err := ParseQueries(queriesText, "")

	if err != nil {
		t.Error(err.Error())
	}
	if len(queries) != 1 {
		t.Errorf("expected 1 query; got %+v", len(queries))
	}
	if queries[0] != `select 1` {
		t.Errorf("got unexpected query `%+v`", queries[0])
	}
}

func TestParseQueriesMulti(t *testing.T) {
	useCases := []string{
		`select 1; select 2; select 3`,
		`select 1; select 2; select 3 `,
		` select 1; select 2 ; select 3;`,
		`select 1; select 2; select 3; `,
		`select 1; ;;select 2; select 3;`,
		`select 1;
                    select 2; select 3;`,
	}
	expected := `select 1;select 2;select 3`
	for _, queriesText := range useCases {
		queries, err := ParseQueries(queriesText, "")

		if err != nil {
			t.Error(err.Error())
		}
		if len(queries) != 3 {
			t.Errorf("expected 3 queries; got %+v", len(queries))
		}
		result := strings.Join(queries, ";")
		if result != expected {
			t.Errorf("got unexpected results: `%+v`", result)
		}
	}
}

func TestParseQueriesEmpty(t *testing.T) {
	useCases := []string{
		``,
		`  `,
		`;`,
		`  ;;; ; ; `,
	}
	for _, queriesText := range useCases {
		queries, err := ParseQueries(queriesText, "")

		if err != nil {
			t.Error(err.Error())
		}
		if len(queries) != 0 {
			t.Errorf("expected 0 queries; got %+v", len(queries))
		}
	}
}

func TestParseQueriesQuotes(t *testing.T) {
	queriesText := `select 1; select '2;2' ; select 3`
	queries, err := ParseQueries(queriesText, "")

	if err != nil {
		t.Error(err.Error())
	}
	if len(queries) != 3 {
		t.Errorf("expected 3 queries; got %+v", len(queries))
		t.Errorf("%+v", strings.Join(queries, "::"))
	}
	result := strings.Join(queries, ";")
	if result != `select 1;select '2;2';select 3` {
		t.Errorf("got unexpected results: `%+v`", result)
	}
}

func TestParseQueriesFile(t *testing.T) {
	queriesText := `select 1; select '2;2' ; select 3`
	tmpFile, err := ioutil.TempFile(os.TempDir(), "querytest-")

	if err != nil {
		t.Errorf("error creating temporary file: %s", err.Error())
	}

	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.Write([]byte(queriesText)); err != nil {
		t.Errorf("error while trying to write to temporary file %s: %s", tmpFile.Name(), err.Error())
	}

	queries, err := ParseQueries("", tmpFile.Name())
	if err != nil {
		t.Error(err.Error())
	}

	if len(queries) != 3 {
		t.Errorf("expected 3 queries; got %+v", len(queries))
		t.Errorf("%+v", strings.Join(queries, "::"))
	}
	result := strings.Join(queries, ";")
	if result != `select 1;select '2;2';select 3` {
		t.Errorf("got unexpected results: `%+v`", result)
	}
}
