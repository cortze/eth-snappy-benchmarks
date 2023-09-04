package csvs

const (
	BlockNumberPosition int64 = iota
)

type FetchableBLock struct {
	Slot string
}

func (c *CSV) ReadTargetBlocks() ([]*FetchableBLock, error) {
	lines, err := c.items()
	if err != nil {
		return nil, err
	}

	lines = lines[1:] // remove the header

	fetchableBlocks := make([]*FetchableBLock, 0, len(lines))

	for _, line := range lines {
		blockRow := new(FetchableBLock)
		blockRow.Slot = line[BlockNumberPosition]

		fetchableBlocks = append(fetchableBlocks, blockRow)
	}
	return fetchableBlocks, nil
}
