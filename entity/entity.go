package entity

type User struct {
	ID       uint       `gorm:"primarykey"`
	Segments []*Segment `gorm:"many2many:user_segments;constraint:OnDelete:CASCADE;"`
}

type Segment struct {
	ID    uint    `gorm:"primarykey"`
	Name  string  `gorm:"unique"`
	Users []*User `gorm:"many2many:user_segments;constraint:OnDelete:CASCADE;"`
}
