package entity

type StickySession struct {
	Enabled     bool
	MaxLifeTime int
}

func NewStickySession(enabled bool, maxLifeTime int) StickySession {
	return StickySession{Enabled: enabled, MaxLifeTime: maxLifeTime}
}
