package handler

import (
	"database/sql"
	"errors"
	"github.com/GoncharovFyodor/hezzltest/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (s *Server) CreateGood(c *gin.Context) {
	var input domain.CreateGoodRequest

	projectID, err := strconv.Atoi(c.Query("projectId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid project id",
		})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}

	good, err := s.services.Goods.CreateGood(c.Request.Context(), projectID, input)

	if err != nil {
		s.log.Infof("ошибка при создании товара: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	c.JSON(http.StatusOK, good)
}

func (s *Server) UpdateGood(c *gin.Context) {
	var input domain.UpdateGoodRequest

	projectID, err := strconv.Atoi(c.Query("projectId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid project id",
		})
		return
	}

	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}

	good, err := s.services.Goods.UpdateGood(c.Request.Context(), projectID, id, input)

	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  3,
			"message": "errors.common.notFound",
			"details": gin.H{},
		})
		return
	}

	if err != nil {
		s.log.Infof("ошибка при обновлении товара: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	c.JSON(http.StatusOK, good)
}

func (s *Server) DeleteGood(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Query("projectId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid project id",
		})
		return
	}

	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	good, err := s.services.Goods.DeleteGood(c.Request.Context(), projectID, id)

	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  3,
			"message": "errors.common.notFound",
			"details": gin.H{},
		})
		return
	}

	if err != nil {
		s.log.Infof("ошибка при удалении товара: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	c.JSON(http.StatusOK, good)
}

func (s *Server) GetGoods(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))

	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(c.Query("offset"))

	if err != nil {
		offset = 0
	}

	goodsList, err := s.services.Goods.GetGoods(c.Request.Context(), limit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	c.JSON(http.StatusOK, goodsList)
}

func (s *Server) ReprioritizeGood(c *gin.Context) {
	var input domain.ReprioritizeRequest

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}

	projectID, err := strconv.Atoi(c.Query("projectId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid project id",
		})
		return
	}

	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	good, err := s.services.Goods.ReprioritizeGood(c.Request.Context(), projectID, id, input)

	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  3,
			"message": "errors.common.notFound",
			"details": gin.H{},
		})
		return
	}

	if err != nil {
		s.log.Infof("ошибка при изменении приоритета товара: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	c.JSON(http.StatusOK, good)
}
