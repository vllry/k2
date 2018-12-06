package main

import "testing"

func TestGenerateContainerTargetList(t *testing.T) {
	tables := []struct {
		containerList []string
		prefix        string
		replicas      int
	}{
		{make([]string, 0), "test", 0},
		{make([]string, 0), "test", 1},
		{[]string{"test-1", "test-2"}, "test", 1},
		{[]string{"test-1", "test-2"}, "test", 0},
		{[]string{"test-1", "test-2", "test-3", "test-4", "test-5"}, "test", 2},
	}

	for _, table := range tables {
		containers, changed := generateContainerTargetList(table.containerList, table.prefix, table.replicas)
		if len(containers) != table.replicas {
			t.Errorf("Replicas did not match - expected %d, got %d", table.replicas, len(containers))
		}
		if (len(table.containerList) != table.replicas) != changed {
			t.Errorf("Expected replica change to be %t", !changed)
		}
	}
}
