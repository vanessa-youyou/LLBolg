package databases

import "LlBlog/models"

func UserLoginCheck(id uint, psw string) bool {
	var u models.UserInfo
	if err := DB.Where("id = ? AND password = ?", id, psw).Find(&u).Error; err != nil {
		return false
	}
	return true
}
