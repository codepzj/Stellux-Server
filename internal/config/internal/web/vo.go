package web

import (
	"time"

	"github.com/codepzj/Stellux-Server/internal/config/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ConfigVO 网站配置VO
type ConfigVO struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Type      string    `json:"type"`
	Content   Content   `json:"content"`
}

// ConfigSummaryVO 网站配置摘要VO
type ConfigSummaryVO struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ConfigDtoToDomain DTO转域模型
func (h *ConfigHandler) ConfigDtoToDomain(dto ConfigDto) *domain.Config {
	return &domain.Config{
		Type:    dto.Type,
		Content: h.ContentDtoToDomain(dto.Content),
	}
}

// ConfigUpdateDtoToDomain 更新DTO转域模型
func (h *ConfigHandler) ConfigUpdateDtoToDomain(dto ConfigUpdateDto) *domain.Config {
	objId, _ := bson.ObjectIDFromHex(dto.ID)
	return &domain.Config{
		Id:      objId,
		Type:    dto.Type,
		Content: h.ContentDtoToDomain(dto.Content),
	}
}

// ContentDtoToDomain 网站内容DTO转域模型
func (h *ConfigHandler) ContentDtoToDomain(dto Content) domain.Content {
	var repos []domain.Repo
	for _, repo := range dto.Repositories {
		repos = append(repos, domain.Repo{
			Name: repo.Name,
			URL:  repo.URL,
			Desc: repo.Desc,
		})
	}

	var skills []domain.Skill
	for _, skill := range dto.Skills {
		skills = append(skills, domain.Skill{
			Category: skill.Category,
			Items:    skill.Items,
		})
	}

	var timeline []domain.Timeline
	for _, item := range dto.Timeline {
		timeline = append(timeline, domain.Timeline{
			Year:  item.Year,
			Title: item.Title,
			Desc:  item.Desc,
		})
	}

	return domain.Content{
		Title:             dto.Title,
		Description:       dto.Description,
		Avatar:            dto.Avatar,
		Name:              dto.Name,
		Bio:               dto.Bio,
		Github:            dto.Github,
		Blog:              dto.Blog,
		Location:          dto.Location,
		TechStacks:        dto.TechStacks,
		Repositories:      repos,
		Quote:             dto.Quote,
		Motto:             dto.Motto,
		ShowRecentPosts:   dto.ShowRecentPosts,
		RecentPostsCount:  dto.RecentPostsCount,
		Skills:            skills,
		Timeline:          timeline,
		Interests:         dto.Interests,
		FocusItems:        dto.FocusItems,
		SEOTitle:          dto.SEOTitle,
		SEOKeywords:       dto.SEOKeywords,
		SEODescription:    dto.SEODescription,
		SEOAuthor:         dto.SEOAuthor,
		RobotsMeta:        dto.RobotsMeta,
		CanonicalURL:      dto.CanonicalURL,
		OGTitle:           dto.OGTitle,
		OGDescription:     dto.OGDescription,
		OGImage:           dto.OGImage,
		TwitterCard:       dto.TwitterCard,
	}
}

// ConfigToVO 域模型转VO
func (h *ConfigHandler) ConfigToVO(config *domain.Config) *ConfigVO {
	return &ConfigVO{
		ID:        config.Id.Hex(),
		CreatedAt: config.CreatedAt,
		UpdatedAt: config.UpdatedAt,
		Type:      config.Type,
		Content:   h.ContentToVO(config.Content),
	}
}

// ContentToVO 网站内容域模型转VO
func (h *ConfigHandler) ContentToVO(content domain.Content) Content {
	var repos []Repo
	for _, repo := range content.Repositories {
		repos = append(repos, Repo{
			Name: repo.Name,
			URL:  repo.URL,
			Desc: repo.Desc,
		})
	}

	var skills []Skill
	for _, skill := range content.Skills {
		skills = append(skills, Skill{
			Category: skill.Category,
			Items:    skill.Items,
		})
	}

	var timeline []Timeline
	for _, item := range content.Timeline {
		timeline = append(timeline, Timeline{
			Year:  item.Year,
			Title: item.Title,
			Desc:  item.Desc,
		})
	}

	return Content{
		Title:             content.Title,
		Description:       content.Description,
		Avatar:            content.Avatar,
		Name:              content.Name,
		Bio:               content.Bio,
		Github:            content.Github,
		Blog:              content.Blog,
		Location:          content.Location,
		TechStacks:        content.TechStacks,
		Repositories:      repos,
		Quote:             content.Quote,
		Motto:             content.Motto,
		ShowRecentPosts:   content.ShowRecentPosts,
		RecentPostsCount:  content.RecentPostsCount,
		Skills:            skills,
		Timeline:          timeline,
		Interests:         content.Interests,
		FocusItems:        content.FocusItems,
		SEOTitle:          content.SEOTitle,
		SEOKeywords:       content.SEOKeywords,
		SEODescription:    content.SEODescription,
		SEOAuthor:         content.SEOAuthor,
		RobotsMeta:        content.RobotsMeta,
		CanonicalURL:      content.CanonicalURL,
		OGTitle:           content.OGTitle,
		OGDescription:     content.OGDescription,
		OGImage:           content.OGImage,
		TwitterCard:       content.TwitterCard,
	}
}

// ConfigListToVOList 配置列表转VO列表
func (h *ConfigHandler) ConfigListToVOList(configs []*domain.Config) []*ConfigVO {
	var vos []*ConfigVO
	for _, config := range configs {
		vos = append(vos, h.ConfigToVO(config))
	}
	return vos
}

// ConfigToSummaryVO 配置转摘要VO
func (h *ConfigHandler) ConfigToSummaryVO(config *domain.Config) *ConfigSummaryVO {
	return &ConfigSummaryVO{
		ID:        config.Id.Hex(),
		Type:      config.Type,
		Title:     config.Content.Title,
		UpdatedAt: config.UpdatedAt,
	}
}

// ConfigListToSummaryVOList 配置列表转摘要VO列表
func (h *ConfigHandler) ConfigListToSummaryVOList(configs []*domain.Config) []*ConfigSummaryVO {
	var vos []*ConfigSummaryVO
	for _, config := range configs {
		vos = append(vos, h.ConfigToSummaryVO(config))
	}
	return vos
}
