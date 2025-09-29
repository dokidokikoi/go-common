package meta

type DoUpdates struct {
	Columns []string

	Updates   []string
	Values    []any
	DoNothing bool
	UpdateAll bool
}

type CreateOption struct {
	Omit      []string
	DoUpdates *DoUpdates
}

type CreateCollectionOption struct {
	Omit      []string
	DoUpdates *DoUpdates
}
