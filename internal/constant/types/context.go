package types

type key int

const (
	ContextKeyUser       key = 0
	ContextKeyVisibility key = 1
)

type VisibilityLevel int

const (
	VisibilityLevelPublic  VisibilityLevel = 1
	VisibilityLevelPrivate VisibilityLevel = 2
	VisibilityLevelServer  VisibilityLevel = 10
)
