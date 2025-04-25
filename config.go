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
	batchSize            int
}

var (
	defaultConfig = config{
		tagName:              "orm",
		enableOptimisticLock: false,
		rewriteQuery:         true,
		batchSize:            200,
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

func SetBatchSize(batchSize int) {
	defaultConfig.batchSize = batchSize
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

func WithBatchSize(batchSize int) func(c *config) {
	return func(c *config) {
		c.batchSize = batchSize
	}
}
