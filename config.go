package orm

type config struct {
	tagName              string
	enableOptimisticLock bool
	versionTag           string
	rewriteQuery         bool
}

var (
	defaultConfig = config{
		tagName:              "orm",
		enableOptimisticLock: false,
		versionTag:           "version",
		rewriteQuery:         true,
	}
)

func SetTagName(tag string) {
	defaultConfig.tagName = tag
}

func SetEnableOptimisticLock(enabled bool) {
	defaultConfig.enableOptimisticLock = enabled
}

func SetVersionTag(tag string) {
	defaultConfig.versionTag = tag
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

func WithVersionTag(tag string) func(c *config) {
	return func(c *config) {
		c.versionTag = tag
	}
}

func WithRewriteQuery(enabled bool) func(c *config) {
	return func(c *config) {
		c.rewriteQuery = enabled
	}
}
