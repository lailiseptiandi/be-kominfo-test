package requests

type RoomRequest struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}
