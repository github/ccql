package text

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestParseHostsSingle(t *testing.T) {
	hostsList := `localhost`
	hosts, err := ParseHosts(hostsList, "")

	if err != nil {
		t.Error(err.Error())
	}
	if len(hosts) != 1 {
		t.Errorf("expected 1 host; got %+v", len(hosts))
	}
	if hosts[0] != `localhost:3306` {
		t.Errorf("got unexpected host `%+v`", hosts[0])
	}
}

func TestParseHostsMulti(t *testing.T) {
	useCases := []string{
		`host1 host2:4408 host3`,
		` host1 host2:4408 host3 `,
		`host1, host2:4408, host3,`,
		`host1
            host2:4408
            host3`,
		`host1  ,
            host2:4408 host3`,
		`host1, ,, host2:4408, host3, ,,`,
	}
	expected := `host1:3306,host2:4408,host3:3306`
	for _, hostsList := range useCases {
		hosts, err := ParseHosts(hostsList, "")

		if err != nil {
			t.Error(err.Error())
		}
		if len(hosts) != 3 {
			t.Errorf("expected 3 hosts; got %+v", len(hosts))
		}
		result := strings.Join(hosts, ",")
		if result != expected {
			t.Errorf("got unexpected results: `%+v`", result)
		}
	}
}

func TestParseHostFiles(t *testing.T) {
	s := `host1 host2:4408 host3`

	tmpFile, err := ioutil.TempFile(os.TempDir(), "querytest-")

	if err != nil {
		t.Errorf("error creating temporary file: %s", err.Error())
	}

	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.Write([]byte(s)); err != nil {
		t.Errorf("error while trying to write to temporary file %s: %s", tmpFile.Name(), err.Error())
	}

	hosts, err := ParseHosts("", tmpFile.Name())

	if err != nil {
		t.Error(err.Error())
	}
	if len(hosts) != 3 {
		t.Errorf("expected 1 host; got %+v", len(hosts))
	}
	if hosts[0] != `host1:3306` {
		t.Errorf("got unexpected host `%+v`", hosts[0])
	}

	if hosts[1] != `host2:4408` {
		t.Errorf("got unexpected host `%+v`", hosts[1])
	}

	if hosts[2] != `host3:3306` {
		t.Errorf("got unexpected host `%+v`", hosts[2])
	}

}

func TestSplitNonEmpty(t *testing.T) {
	s := "the, quick,, brown,fox ,,"
	splits := SplitNonEmpty(s, ",")

	if len(splits) != 4 {
		t.Errorf("expected 4 tokens; got %+v", len(splits))
	}
	join := strings.Join(splits, ";")
	expected := "the;quick;brown;fox"
	if join != expected {
		t.Errorf("expected tokens: `%+v`. Got: `%+v`", expected, join)
	}
}
