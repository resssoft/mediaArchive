package models

type CoubItem struct {
	Id                   int                `json:"id,omitempty"`
	Title                string             `json:"title,omitempty"`
	Permalink            string             `json:"permalink,omitempty"`
	Type                 string             `json:"type,omitempty"`
	ChannelId            int                `json:"channel_id,omitempty"`
	Duration             float64            `json:"duration,omitempty"`
	CreatedAt            string             `json:"created_at,omitempty"`
	UpdatedAt            string             `json:"updated_at,omitempty"`
	LikeCount            int                `json:"like_count,omitempty"`
	OriginalSound        bool               `json:"original_sound,omitempty"`
	HasSound             bool               `json:"has_sound,omitempty"`
	FileVersions         CoubFileVersions   `json:"file_versions,omitempty"`
	AudioVersions        CoubAudioVersions  `json:"audio_versions,omitempty"`
	ImageVersions        CoubImageVersions  `json:"image_versions,omitempty"`
	FirstFrameVersions   CoubImageVersions  `json:"first_frame_versions,omitempty"`
	Size                 []int              `json:"size,omitempty"`
	AgeRestricted        bool               `json:"age_restricted,omitempty"`
	AgeRestrictedByAdmin bool               `json:"age_restricted_by_admin,omitempty"`
	Banned               bool               `json:"banned,omitempty"`
	ExternalDownload     CoubExternalSource `json:"external_download,omitempty"`
	Channel              CoubSmallChannel   `json:"channel,omitempty"`
	PercentCone          int                `json:"percent_done,omitempty"`
	Tags                 []CoubTag          `json:"tags,omitempty"`
	Picture              string             `json:"picture,omitempty"`
	AudioFileUrl         string             `json:"audio_file_url,omitempty"`
	Categories           []CoubCategories   `json:"categories,omitempty"`
	LikesCount           int                `json:"likes_count,omitempty"`
}

type CoubTag struct {
	Id    int    `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
}

type CoubCategories struct {
	Id                 int    `json:"id,omitempty"`
	Title              string `json:"title,omitempty"`
	Permalink          string `json:"permalink,omitempty"`
	SubscriptionsCount int    `json:"subscriptions_count,omitempty"`
	BigImageUrl        string `json:"big_image_url,omitempty"`
	SmallImageUrl      string `json:"small_image_url,omitempty"`
	MedImageUrl        string `json:"med_image_url,omitempty"`
}

type CoubExternalSource struct {
	Type        string `json:"type,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
	Url         string `json:"url,omitempty"`
	HasEmbed    bool   `json:"has_embed,omitempty"`
}

type CoubSmallChannel struct {
	Id              int               `json:"id,omitempty"`
	Permalink       string            `json:"permalink,omitempty"`
	Title           string            `json:"title,omitempty"`
	Description     *string           `json:"description,omitempty"`
	FollowersCount  int               `json:"followers_count,omitempty"`
	FollowingCount  int               `json:"following_count,omitempty"`
	AvatarVersions  CoubImageVersions `json:"avatar_versions,omitempty"`
	BackgroundImage string            `json:"background_image,omitempty"`
	CoubsCount      int               `json:"coubs_count,omitempty"`
	RecoubsCount    int               `json:"recoubs_count,omitempty"`
}

type CoubFileVersions struct {
	Html5 CoubMediaTypes `json:"html5,omitempty"`
	Share CoubShareType  `json:"share,omitempty"`
}

type CoubAudioVersions struct {
	Template string          `json:"template,omitempty"`
	Versions []string        `json:"versions,omitempty"`
	Chunks   CoubAudioChunks `json:"chunks,omitempty"`
}

type CoubAudioChunks struct {
	Template string   `json:"template,omitempty"`
	Versions []string `json:"versions,omitempty"`
	Chunks   []int    `json:"chunks,omitempty"`
}

type CoubImageVersions struct {
	Template string   `json:"template,omitempty"`
	Versions []string `json:"versions,omitempty"`
}

type CoubShareType struct {
	Default string `json:"default,omitempty"`
}

type CoubMediaTypes struct {
	Video CoubMediaQuality `json:"video,omitempty"`
	Audio CoubMediaQuality `json:"audio,omitempty"`
}

type CoubMediaQuality struct {
	Higher CoubMediaData `json:"higher,omitempty"`
	High   CoubMediaData `json:"high,omitempty"`
	Med    CoubMediaData `json:"med,omitempty"`
}

type CoubMediaData struct {
	Url  string `json:"url,omitempty"`
	Size int    `json:"size,omitempty"`
}

func (ci *CoubItem) GetCategoriesNames() []string {
	categories := make([]string, 0)
	for _, category := range ci.Categories {
		categories = append(categories, category.Title)
	}
	return categories
}

func (ci *CoubItem) GetTagNames() []string {
	tags := make([]string, 0)
	for _, tag := range ci.Tags {
		tags = append(tags, tag.Title)
	}
	return tags
}
