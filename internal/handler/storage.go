package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hantdev/chorus-controller/internal/domain"
)

// StorageHandler handles storage-related endpoints
type StorageHandler struct {
	storageService domain.StorageService
}

// NewStorageHandler creates a new storage handler
func NewStorageHandler(storageService domain.StorageService) *StorageHandler {
	return &StorageHandler{
		storageService: storageService,
	}
}

// ListStorages
// @Summary		List all storages
// @Description	Returns a list of all configured storage backends
// @Tags			storages
// @Accept			json
// @Produce		json
// @Success		200	{object}	map[string]interface{}
// @Failure		502	{object}	map[string]interface{}
// @Router			/storages [get]
func (h *StorageHandler) ListStorages(c *gin.Context) {
	resp, err := h.storageService.ListStorages(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListBuckets
// @Summary		List buckets available for replication
// @Description	Returns a list of buckets that can be used for replication jobs between specified storages
// @Tags			buckets
// @Accept			json
// @Produce		json
// @Param			user				query		string	false	"User identifier"
// @Param			from				query		string	false	"Source storage name"
// @Param			to					query		string	false	"Destination storage name"
// @Param			show_replicated		query		bool	false	"Whether to show already replicated buckets"
// @Success		200					{object}	map[string]interface{}
// @Failure		502					{object}	map[string]interface{}
// @Router			/buckets [get]
func (h *StorageHandler) ListBuckets(c *gin.Context) {
	var req domain.ListBucketsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	resp, err := h.storageService.ListBuckets(c.Request.Context(), &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"available_buckets":  resp.Buckets,
		"replicated_buckets": resp.ReplicatedBuckets,
	})
}

// CreateStorage
// @Summary		Create a storage configuration
// @Description	Creates a storage with simplified parameters. All other settings use default values.
// @Tags			storages
// @Accept			json
// @Produce		json
// @Param			request			body		domain.CreateStorageRequest	true	"Storage creation request"
// @Success		201				{string}	string	"Storage created successfully"
// @Failure		400				{object}	map[string]interface{}
// @Router			/storages [post]
func (h *StorageHandler) CreateStorage(c *gin.Context) {
	var req domain.CreateStorageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
	if err := h.storageService.CreateStorageFromRequest(c.Request.Context(), &req); err != nil {
		HandleError(c, err)
		return
	}
	c.Status(http.StatusCreated)
}

// ListStoragesDB
// @Summary		List storages from DB
// @Tags			storages
// @Produce		json
// @Success		200 {array} domain.Storage
// @Router			/storages/db [get]
func (h *StorageHandler) ListStoragesDB(c *gin.Context) {
	items, err := h.storageService.ListStorageFromDB(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetStorage
// @Summary		Get storage by ID
// @Tags			storages
// @Produce		json
// @Param			id			path		string	true	"Storage ID"
// @Success		200				{object}	domain.Storage
// @Failure		404				{object}	map[string]interface{}
// @Router			/storages/{id} [get]
func (h *StorageHandler) GetStorage(c *gin.Context) {
	id := c.Param("id")
	storage, err := h.storageService.GetStorageByID(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, storage)
}

// UpdateStorage
// @Summary		Update a storage configuration
// @Description	Updates an existing storage with new parameters
// @Tags			storages
// @Accept			json
// @Produce		json
// @Param			id			path		string	true	"Storage ID"
// @Param			request			body		domain.CreateStorageRequest	true	"Storage update request"
// @Success		200				{string}	string	"Storage updated successfully"
// @Failure		400				{object}	map[string]interface{}
// @Failure		404				{object}	map[string]interface{}
// @Router			/storages/{id} [put]
func (h *StorageHandler) UpdateStorage(c *gin.Context) {
	id := c.Param("id")
	var req domain.CreateStorageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
	if err := h.storageService.UpdateStorageByID(c.Request.Context(), id, &req); err != nil {
		HandleError(c, err)
		return
	}
	c.Status(http.StatusOK)
}

// DeleteStorage
// @Summary		Delete a storage configuration
// @Description	Deletes a storage and all its credentials
// @Tags			storages
// @Produce		json
// @Param			id			path		string	true	"Storage ID"
// @Success		200				{string}	string	"Storage deleted successfully"
// @Failure		404				{object}	map[string]interface{}
// @Router			/storages/{id} [delete]
func (h *StorageHandler) DeleteStorage(c *gin.Context) {
	id := c.Param("id")
	if err := h.storageService.DeleteStorageByID(c.Request.Context(), id); err != nil {
		HandleError(c, err)
		return
	}
	c.Status(http.StatusOK)
}
