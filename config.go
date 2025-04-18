package orm

type config struct {
	tagName              string
	enableOptimisticLock bool
	versionTag           string
}

var (
	defaultConfig = config{
		tagName:              "orm",
		enableOptimisticLock: false,
		versionTag:           "version",
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

type callOpt func(c *config)

func WithTagName(tag string) callOpt {
	return func(c *config) {
		c.tagName = tag
	}
}

func WithEnableOptimisticLock(enabled bool) callOpt {
	return func(c *config) {
		c.enableOptimisticLock = enabled
	}
}

func WithVersionTag(tag string) callOpt {
	return func(c *config) {
		c.versionTag = tag
	}
}
