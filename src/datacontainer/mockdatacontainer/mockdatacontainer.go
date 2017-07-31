package mockdatacontainer

import "common"

type Container struct {
	capacity int64
	locs []*common.Location
}

func New(cap int64) *Container {
	return &Container{
		capacity: cap,
	}
}

func (c *Container) Insert(loc *common.Location) {
	c.locs = append(c.locs, loc)
}

func (c *Container) Remove(loc *common.Location) {
	for i, l := range c.locs {
		if l.ID == loc.ID {
			c.locs = append(c.locs[:i], c.locs[i+1:]...)
		}
	}
}

func (c *Container) KNearestNeighbour(loc *common.Location, k int) []*common.Location {
	if k > len(c.locs) {
		k = len(c.locs)
	}
	used := map[int64]bool{}
	neighbours := []*common.Location{}
	for i := 0;i < k;i++ {
		var best_loc *common.Location
		for _, l := range c.locs {
			if _, fnd := used[l.ID]; fnd {
				continue
			}
			if best_loc == nil || c.Distance(best_loc, loc) > c.Distance(l, loc) {
				best_loc = l
			}
		}
		used[best_loc.ID] = true
		neighbours = append(neighbours, best_loc)
	}
	return neighbours
}

func (c *Container) Distance(loc1 *common.Location, loc2 *common.Location) float64 {
	diff_lat := loc1.Lat - loc2.Lat
	diff_lon := loc1.Lon - loc2.Lon
	return diff_lat * diff_lat + diff_lon * diff_lon
}

func (c *Container) Serialize() string {
	return ""
}

func (c *Container) Deserialize(s string) {
}

func (c *Container) Generate(locs []*common.Location) {
	c.locs = locs
}