package ddd

type Entity interface {
	SameIdentityAs(e interface{}) bool
}
