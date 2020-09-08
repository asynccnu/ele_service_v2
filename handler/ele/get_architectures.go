package ele

// type GetArchitectureResponse struct {
// 	Count int
// 	List  []*service.ArchitectureInfo
// }

// // 获取楼栋信息
// func GetArchitectures(c *gin.Context) {
// 	area := c.DefaultQuery("area", "")
// 	if area == "" {
// 		handler.SendBadRequest(c, errno.ErrQuery, nil, "area is required")
// 		return
// 	}

// 	architectures, err := service.GetArchitectures(area)
// 	if err != nil {
// 		log.Error("GetArchitectures function error")
// 		handler.SendError(c, err, nil, err.Error())
// 		return
// 	}

// 	handler.SendResponse(c, nil, &GetArchitectureResponse{
// 		Count: len(architectures.List),
// 		List:  architectures.List,
// 	})
// }
