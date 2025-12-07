package dao

type Article struct {
	Id         int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title      string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	Content    string `gorm:"type=BLOB" bson:"content,omitempty"`
	CoverImage string `gorm:"type=varchar(1024)" bson:"cover_image,omitempty"` // 封面图片URL
	// 作者
	AuthorId int64 `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8 `bson:"status,omitempty"`
	Ctime    int64 `bson:"ctime,omitempty"`
	Utime    int64 `bson:"utime,omitempty" gorm:"index"`
}
type PublishedArticle Article

// PublishedArticleV1 s3 专属
type PublishedArticleV1 struct {
	Id         int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title      string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	CoverImage string `gorm:"type=varchar(1024)" bson:"cover_image,omitempty"` // 封面图片URL
	AuthorId   int64  `gorm:"index" bson:"author_id,omitempty"`
	Status     uint8  `bson:"status,omitempty"`
	Ctime      int64  `bson:"ctime,omitempty"`
	Utime      int64  `bson:"utime,omitempty"`
}
