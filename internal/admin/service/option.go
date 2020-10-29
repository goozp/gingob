package service

import (
	"fmt"

	"github.com/puti-projects/puti/internal/admin/dao"
	"github.com/puti-projects/puti/internal/pkg/config"
	"github.com/puti-projects/puti/internal/pkg/errno"
	optionCache "github.com/puti-projects/puti/internal/pkg/option"
	"github.com/puti-projects/puti/internal/pkg/theme"
	"github.com/puti-projects/puti/internal/utils"
)

// OptionUpdateRequest update request include params which are dynamitly
type OptionUpdateRequest struct {
	Parms map[string]interface{}
}

// OptionListResponse option list and extra data if exist
type OptionListResponse struct {
	Options   map[string]string `json:"options"`
	ExtraData interface{}       `json:"extraData"`
}

// OptionListTheme the list of all theme
type ListTheme struct {
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
}

// ListOption get option list by setting type
func ListOption(settingType string) (*OptionListResponse, error) {
	// get options list
	options, err := GetOptionsByType(settingType)
	if err != nil {
		return nil, errno.New(errno.ErrDatabase, err)
	}

	rsp := &OptionListResponse{
		Options: options,
	}

	if "theme" == settingType {
		rsp.ExtraData = generalThemeExtra()
	}

	return rsp, nil
}

func generalThemeExtra() []*ListTheme {
	var themes []*ListTheme
	// theme pkg have already load all theme in dir
	for _, t := range theme.Themes {
		var themeThumbnail string
		// thumbnail
		if exist, _ := utils.PathExists(config.StaticPath("theme/"+t) + "/thumbnail.jpg"); exist == true {
			themeThumbnail = config.StaticPath("theme/"+t) + "/thumbnail.jpg"
		} else {
			themeThumbnail = config.StaticPath("assets/images/") + "/image_default.png"
		}

		themeInfo := &ListTheme{
			Name:      t,
			Thumbnail: themeThumbnail,
		}
		themes = append(themes, themeInfo)
	}
	return themes
}

// GetOptionsByType get default options by setting type
func GetOptionsByType(settingType string) (map[string]string, error) {
	optionNames := GetDefaultOptionNamesByType(settingType)
	options, err := dao.Engine.GetAllOptions(optionNames)
	if err != nil {
		return nil, errno.New(errno.ErrDatabase, err)
	}

	optionsMap := make(map[string]string)
	for _, option := range options {
		optionsMap[option.OptionName] = option.OptionValue
	}

	return optionsMap, nil
}

// UpdateOptions update options
func UpdateOptions(options map[string]interface{}) error {
	// update options
	if err := dao.Engine.UpdateOptions(options); err != nil {
		return errno.New(errno.ErrDatabase, err)
	}

	// update options cache
	for optionName, optionValue := range options {
		optionCache.Options.Put(optionName, fmt.Sprintf("%v", optionValue))
	}

	return nil
}

// GetDefaultOptionsByType get default setting's option name
func GetDefaultOptionNamesByType(settingType string) []string {
	var optionNames []string
	switch settingType {
	case "general":
		optionNames = []string{"blog_name", "blog_description", "site_url", "admin_email", "users_can_register", "timezone_string", "site_language"}
	case "property":
		optionNames = []string{"site_description", "site_keywords", "footer_copyright"}
	case "menu":
		optionNames = []string{""}
	case "theme":
		optionNames = []string{"current_theme"}
	case "pubic-account":
		optionNames = []string{""}
	case "github":
		optionNames = []string{""}
	case "discuss":
		optionNames = []string{"article_comment_status", "page_comment_status", "comment_need_register", "show_comment_page", "comment_per_page", "comment_page_first", "comment_page_top", "comment_before_show", "show_avatar"}
	case "media":
		optionNames = []string{"image_thumbnail_width", "image_thumbnail_height", "image_medium_width", "image_medium_height", "image_large_width", "image_large_height"}
	case "reading":
		optionNames = []string{"show_on_front", "show_on_front_page", "posts_per_page", "open_XML"}
	case "writing":
		optionNames = []string{"default_category", "default_link_category"}
	default:
		optionNames = []string{}
	}

	return optionNames
}