package dao

import (
	"context"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Setting struct {
	Key   string     `bson:"key"`
	Value SiteConfig `bson:"value"`
}

type SiteConfig struct {
	SiteTitle       string `bson:"site_title"`       // 网站标题
	SiteSubtitle    string `bson:"site_subtitle"`   // 网站副标题
	SiteUrl         string `bson:"site_url"`         // 网站地址
	SiteFavicon     string `bson:"site_favicon"`     // 网站图标
	SiteAuthor      string `bson:"site_author"`      // 网站作者
	SiteAnimateText string `bson:"site_animate_text"` // 网站打字机文本
	SiteAvatar      string `bson:"site_avatar"`      // 网站头像
	SiteKeywords    string `bson:"site_keywords"`    // 网站关键词
	SiteDescription string `bson:"site_description"` // 网站描述
	SiteCopyright   string `bson:"site_copyright"`   // 网站版权
	SiteIcp         string `bson:"site_icp"`         // 网站备案号
	SiteIcpLink     string `bson:"site_icplink"`    // 网站备案号链接
	GithubUsername  string `bson:"github_username"`  // Github用户名
}

type ISettingDao interface {
	Upsert(ctx context.Context, setting *Setting) error
	GetSetting(ctx context.Context, key string) (*Setting, error)
}

var _ ISettingDao = (*SettingDao)(nil)

func NewSettingDao(db *mongox.Database) *SettingDao {
	return &SettingDao{coll: mongox.NewCollection[Setting](db, "setting")}
}

type SettingDao struct {
	coll *mongox.Collection[Setting]
}

func (s *SettingDao) Upsert(ctx context.Context, setting *Setting) error {
	updateResult, err := s.coll.Updater().Filter(query.Eq("key", setting.Key)).Updates(update.NewBuilder().Set("key", setting.Key).Set("value", setting.Value).Build()).Upsert(ctx)
	if err != nil {
		return err
	}
	if updateResult.ModifiedCount == 0 && updateResult.UpsertedID == nil {
		return errors.New("设置新增或更新失败")
	}
	return nil
}

func (s *SettingDao) GetSetting(ctx context.Context, key string) (*Setting, error) {
	setting, err := s.coll.Finder().Filter(query.Eq("key", key)).FindOne(ctx)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &Setting{
			Key:   key,
			Value: SiteConfig{},
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return &Setting{
		Key:   setting.Key,
		Value: setting.Value,
	}, nil
}
