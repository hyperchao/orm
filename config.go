package orm

const (
	TagPrimaryKey    = "primary"
	TagAutoIncrement = "autoincrement"
	TagVersion       = "version"
)

type config struct {
	tagName              string
	enableOptimisticLock bool
	rewriteQuery         bool
}

var (
	defaultConfig = config{
		tagName:              "orm",
		enableOptimisticLock: false,
		rewriteQuery:         true,
	}
)

func SetTagName(tag string) {
	defaultConfig.tagName = tag
}

func SetEnableOptimisticLock(enabled bool) {
	defaultConfig.enableOptimisticLock = enabled
}

func SetRewriteQuery(enabled bool) {
	defaultConfig.rewriteQuery = enabled
}

func WithTagName(tag string) func(c *config) {
	return func(c *config) {
		c.tagName = tag
	}
}

func WithEnableOptimisticLock(enabled bool) func(c *config) {
	return func(c *config) {
		c.enableOptimisticLock = enabled
	}
}

func WithRewriteQuery(enabled bool) func(c *config) {
	return func(c *config) {
		c.rewriteQuery = enabled
	}
}
