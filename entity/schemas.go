package entity

type UserCreateSchema struct {
	ID uint `form:"user_id" json:"user_id" validate:"required"`
}

type SegmentCreateSchema struct {
	Name    string `form:"segment_name" json:"segment_name" validate:"required"`
	Percent int    `json:"auto_add_percent" validate:"gte=0"`
}

type SegmentUpdateSchema struct {
	Name    string `form:"segment_name" json:"segment_name" validate:"required"`
	NewName string `form:"new_segment_name" json:"new_segment_name" validate:"required"`
}

type AddSegmentsSchema struct {
	ID           uint   `form:"user_id" json:"user_id" validate: "required"`
	SegmentNames string `form:"segment_names" json:"segment_names" validate:"required"`
}

type SegmentWithUsersSchema struct {
	Name  string
	Users []uint
}

type SegmentSchema struct {
	Name string
}

type UserWithSegmentsSchema struct {
	ID       uint
	Segments []string
}

type RemoveUserFromSegmentSchema struct {
	ID      uint   `form:"user_id" json:"user_id" validate:"required"`
	Segment string `form:"segment_name" json:"segment_name" validate:"required"`
}

type AddAndRemoveSegmentsSchema struct {
	ID             uint     `form:"user_id" json:"user_id" validate: "required"`
	AddSegments    []string `form:"add_segment_names" json:"add_segment_names" validate:"required"`
	RemoveSegments []string `form:"remove_segment_names" json:"remove_segment_names" validate:"required"`
}

type SegmentLogRequestSchema struct {
	UserID     uint   `form:"user_id" json:"user_id" validate: "required"`
	DateAfter  string `json:"date_after" validate:"required"`
	DateBefore string `json:"date_before" validate:"required"`
}
