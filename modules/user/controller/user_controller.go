package controller

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/Mobilizes/materi-be-alpro/modules/user/service"
    "github.com/Mobilizes/materi-be-alpro/modules/user/validation"
    "github.com/Mobilizes/materi-be-alpro/pkg/utils"
)

type UserController struct {
    service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
    return &UserController{service: service}
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
    req, err := validation.ValidateCreateUser(c)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    user, err := ctrl.service.CreateUser(req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat user")
        return
    }

    utils.SuccessResponse(c, http.StatusCreated, "User berhasil dibuat", user)
}
