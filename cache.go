package httpx

import (
	"fmt"
	"time"
)

type Cache struct {
	ctx *Ctx
}

func (ctx *Ctx) Cache() Cache {
	return Cache{ctx}
}

func (cache Cache) Prohibit() {
	cache.ctx.SetHeader("Cache-Control", "no-store")
}

func (cache Cache) Public(maxAge time.Duration) {
	cache.ctx.SetHeader("Cache-Control", fmt.Sprintf("public, max-age=%.0f", maxAge.Seconds()))
}

func (cache Cache) Private(maxAge time.Duration) {
	cache.ctx.SetHeader("Cache-Control", fmt.Sprintf("private, max-age=%.0f", maxAge.Seconds()))
}
