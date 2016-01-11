package text

import (
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
