package chat_api

import (
	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/gin-gonic/gin"
)

func (ChatApi) UserChatDeleteView(c *gin.Context) {
	cr := middleware.GetBind[models.RemoveRequest](c)

	var chatList []models.ChatModel
	global.DB.Find(&chatList, "id in ?", cr.IDList)

	claims := jwts.GetClaims(c)
	chatMap := common.ScanMapV2(models.UserChatActionModel{}, common.ScanOption{
		Where: global.DB.Where("user_id = ? and chat_id in ?", claims.UserID, cr.IDList),
		Key:   "ChatID",
	})

	var addChatAc []models.UserChatActionModel
	var updateChatAcIdList []uint
	for _, model := range chatList {
		// 判断这个消息是不是删过了
		chat, ok := chatMap[model.ID]
		if !ok {
			// 找不到的情况
			addChatAc = append(addChatAc, models.UserChatActionModel{
				UserID:   claims.UserID,
				ChatID:   model.ID,
				IsDelete: true,
			})
			continue
		}
		if chat.IsDelete {
			continue
		}
		updateChatAcIdList = append(updateChatAcIdList, chat.ID)
	}

	if len(addChatAc) > 0 {
		err := global.DB.Create(&addChatAc).Error
		if err != nil {
			res.FailWithMsg("删除消息失败", c)
			return
		}
	}
	if len(updateChatAcIdList) > 0 {
		global.DB.Model(&models.UserChatActionModel{}).Where("id in ?", updateChatAcIdList).Update("is_delete", true)
	}
	res.OkWithMsg("消息删除成功", c)
}
