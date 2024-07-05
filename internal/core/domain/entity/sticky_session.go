package entity

type StickySession struct {
	Enabled     bool
	MaxLifeTime int64
}

func NewStickySession(enabled bool, maxLifeTime int64) StickySession {
	return StickySession{Enabled: enabled, MaxLifeTime: maxLifeTime}
}
