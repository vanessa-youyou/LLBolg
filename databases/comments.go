package databases

import (
	"LlBlog/models"
	"strconv"
)

// CommentDelete 删除评论
func CommentDelete(cm *models.CommentInfo, userId uint) (bool, error) {
	// 检验操作人是不是评论人员：
	if cm.UserID != userId{
		return false, nil
	}

	// 删除这条评论 在 redis中的记录
	SetName := strconv.Itoa(int(cm.ID))
	SetName += "LikeComment:"
	Redis.Del(SetName)

	// 删除此条评论
	err := DB.Model(cm).Where("id = ?",cm.ID).Delete(&models.CommentInfo{}).Error
	if err != nil{
		return false, err
	}
	return true, nil
}

// LikeComment 评论点赞 点赞/取消点赞操作	-redis
func LikeComment(cm *models.CommentInfo, userId uint) (bool, error){
	// set的名字为评论的id 内容为 用户id
	CommentId := strconv.Itoa(int(cm.ID))
	CommentId += "LikeComment:"
	UserId := strconv.Itoa(int(userId))
	if !Redis.SIsMember(CommentId, UserId).Val(){
		// 如果没有 就把这个数据存入
		Redis.SAdd(CommentId, UserId)
	}else{
		// 如果存在 就删除
		Redis.SRem(CommentId, UserId)
	}
	return true, nil
}