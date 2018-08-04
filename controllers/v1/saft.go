package controllers

import (
	"saft_parser/util"

	"saft_parser/services"

	"github.com/gin-gonic/gin"
)

type (
	// SAFTController represents the controller for operating on SAFT files
	SAFTController struct {
		saftServ *services.SAFTService
	}
)

// NewRoleController is the constructor of RoleController
func NewSAFTController(saftServ *services.SAFTService) *SAFTController {
	return &SAFTController{
		saftServ: saftServ,
	}
}

// CreateAction creates a new role
func (this SAFTController) UploadAction(c *gin.Context) {

	fileHeader, err := c.FormFile("SAFT-XML")
	if err != nil {
		e := util.HandleErrorResponse(util.INVALID_REQUEST, nil, err.Error())
		c.JSON(e.HttpCode, e)
		return
	}

	resp, e := this.saftServ.UploadAction(fileHeader)

	if e != nil {
		c.JSON(e.HttpCode, e)
		return
	}

	c.JSON(200, resp)
}
