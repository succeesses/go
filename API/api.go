package API

import (
	"awesomeProject/Model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// APInput1 接口1
func APInput1(c *gin.Context) {
	name, _ := strconv.Atoi(c.Param("MemberId"))
	Type := c.Param("Type")
	limit, _ := strconv.Atoi(c.Param("Limit"))
	offset, _ := strconv.Atoi(c.Param("Offset"))
	switch {
	case Type == "deposit": // 对充币进操作
		{
			result, err := Model.HandleAPI1Deposit(name, limit, offset)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    -1,
					"message": "没有相关记录与之匹配",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"data": result,
			})
		}
	case Type == "withdraw": // 对提币进行操作
		{
			result, err := Model.HandleAPI1Withdraw(name, limit, offset)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    -1,
					"message": "没有相关记录与之匹配",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"data": result,
			})
		}
	}
}

// APInput2 接口2
func APInput2(c *gin.Context) {
	Type := c.Param("Type")
	if Type == "add" { // 增加
		var member Model.IwalaMember
		member.Email = c.Request.FormValue("email")
		member.UserName = c.Request.FormValue("user_name")
		member.Name = c.Request.FormValue("name")
		err := Model.HandleAPI2Insert(member)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "添加失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "添加成功",
		})
	}
	if Type == "delete" {
		id := c.Request.FormValue("member_id")
		id1, err := strconv.ParseInt(id, 10, 64)
		err = Model.HandleAPI2Delete(id1)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "删除失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "删除成功",
		})
	}

	if Type == "select" {
		id := c.Request.FormValue("member_id")
		id1, err := strconv.ParseInt(id, 10, 64)
		result, err := Model.HandleAPI2Select(id1)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "抱歉未找到相关信息",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": result,
		})
	}
}

// APInput3 接口3
func APInput3(c *gin.Context) {

	time1 := c.Param("time1")
	time2 := c.Param("time2")
	result, err := Model.HanldeAPI3(time1, time2)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "抱歉未找到相关信息",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": result,
	})
}
