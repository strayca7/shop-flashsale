// used to define the user model

package userdomain

type User struct {
    Uid int `json:"uid"`
    Name string `json:"name"`
    Address []string `json:"address"`
	IsAdmin bool `json:"is_admin"`
}