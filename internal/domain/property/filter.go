package property

type PropertyFilter struct {
	UserID    *int64
	Type      []PropertyType
	ManagerID *int64
	Active    *bool
	Limit     int
	Offset    int
}
