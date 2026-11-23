package requests

type BookRequest struct {
	RoomID    int    `json:"room_id"`
	Title     string `json:"title"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
