package rid

import "github.com/onexstack/onexstack/pkg/id"

const defaultABC = "abcdefghijklmnopqrstuvwxyz1234567890"

type ResourceID string

const (
	// 定义用户资源标识符
	UserID ResourceID = "user"
	// 定义blog资源标识符
	PostID ResourceID = "post"
)

// 将资源标识符转换成字符串
func (rid ResourceID) String() string {
	return string(rid)
}

// 创建带有前缀的唯一标识符
func (rid ResourceID) New(counter uint64) string {
	// 使用自定义选项生成唯一标识符
	uniqueStr := id.NewCode(
		counter,
		id.WithCodeChars([]rune(defaultABC)),
		id.WithCodeL(6),
		id.WithCodeSalt(Salt()),
	)
	return rid.String() + "-" + uniqueStr
}
