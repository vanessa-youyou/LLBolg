package databases

import "LlBlog/models"

// databases应该关注于单纯的通用逻辑
// 例如通过id获取记录这种
func GetUserById(id uint) (*models.UserInfo, error) {
	var u models.UserInfo
	if err := DB.Where("id = ?", id).Find(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
