package cache

import "gitee.com/quant1x/gox/util/homedir"

var (
	default_cache_path = "~/.quant1x" // 数据根路径
)

func init() {
	// 初始化缓存路径
	rootPath, err := homedir.Expand(default_cache_path)
	if err != nil {
		panic(err)
	}
	default_cache_path = rootPath
}

// DefaultCachePath 数据缓存的根路径
func DefaultCachePath() string {
	return default_cache_path
}

// GetBlockPath 板块路径
func GetBlockPath() string {
	return DefaultCachePath() + "/block"
}
