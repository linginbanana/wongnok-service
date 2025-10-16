package profile

type User struct {
    ID       string `gorm:"primaryKey" json:"id"`
    Name     string `json:"name"`
    Nickname string `json:"nickname"`
    Avatar   string `json:"avatar"`
}
