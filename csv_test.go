package struct2sql_test

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"testing"
	"time"

	s2s "fallen.world/mingchoi/struct2sql"
)

func TestSlice2Strings(t *testing.T) {
	expected := "ID,Name,Level,Health,Boss,Spawn,Birth,MID\n1,Monster 怪物,18,1029.8,true,2018-01-01 00:00:00,1960-04-12 00:00:00,1\n2,Monster2 怪物2,29,3049.8,false,2018-01-01 00:00:00,1960-04-12 00:00:00,2\n3,Monster3 怪物3,32,5529.8,true,2018-01-01 00:00:00,1960-04-12 00:00:00,2\n"

	nowTime, _ := time.Parse(time.RFC3339, "2018-01-01T00:00:00+08:00")
	oldTime, _ := time.Parse(time.RFC3339, "1960-04-12T00:00:00+08:00")
	monsters := []Monster{
		Monster{
			ID:     1,
			Name:   "Monster 怪物",
			Level:  18,
			Health: 1029.8,
			Boss:   true,
			Spawn:  nowTime,
			Birth:  oldTime,
			MID:    1,
		},
		Monster{
			ID:     2,
			Name:   "Monster2 怪物2",
			Level:  29,
			Health: 3049.8,
			Boss:   false,
			Spawn:  nowTime,
			Birth:  oldTime,
			MID:    2,
		},
		Monster{
			ID:     3,
			Name:   "Monster3 怪物3",
			Level:  32,
			Health: 5529.8,
			Boss:   true,
			Spawn:  nowTime,
			Birth:  oldTime,
			MID:    2,
		},
	}

	// Convert struct slice to string slice
	list := s2s.Slice2Strings(monsters)

	// Convert stringslice to CSV
	buf := bytes.NewBuffer([]byte{})
	w := csv.NewWriter(buf)
	w.WriteAll(list)

	// Test Case
	if buf.String() != expected {
		t.Error("CSV convert failed to match excepted")
		fmt.Println(buf.String())
	}

	// Write to file
	/*
		err := s2s.Slice2CSV(filepath.Join(os.TempDir(), "test.csv"), monsters, true)
		if err != nil {
			t.Error(err)
		}
	*/
}
