package mockdatacontainer

import (
	"reflect"
	"testing"

	"common"
)

func TestKNearestNeigbours(t *testing.T) {
	c := New()
	c.Insert(&common.Location{1, 5, 5})
	c.Insert(&common.Location{2, 4, 3})
	c.Insert(&common.Location{3, 2, 1})
	c.Insert(&common.Location{4, 2, 2})

	want := []*common.Location{
		{3, 2, 1},
		{4, 2, 2},
		{2, 4, 3},
		{1, 5, 5},
	}
	got := c.KNearestNeighbour(&common.Location{0, 0, 0}, 4)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateFoodTruck #1: Got %v, want %v", got, want)
	}
}
