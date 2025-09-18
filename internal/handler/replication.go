package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hantdev/chorus-controller/internal/domain"
)

// ReplicationHandler handles replication-related endpoints
type ReplicationHandler struct {
	replicationService domain.ReplicationService
}

// NewReplicationHandler creates a new replication handler
func NewReplicationHandler(replicationService domain.ReplicationService) *ReplicationHandler {
	return &ReplicationHandler{
		replicationService: replicationService,
	}
}

// CreateReplication
// @Summary		Create a new replication job
// @Description	Configures a new replication job for specified buckets between storages
// @Tags			replications
// @Accept			json
// @Produce		json
// @Param			replication	body		domain.CreateReplicationRequest	true	"Replication configuration"
// @Success		201			{string}	string				"Replication job created successfully"
// @Failure		400			{object}	map[string]interface{}
// @Failure		502			{object}	map[string]interface{}
// @Router			/replications [post]
func (h *ReplicationHandler) CreateReplication(c *gin.Context) {
	var req domain.CreateReplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	err := h.replicationService.CreateReplication(c.Request.Context(), &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

// ListReplications
// @Summary		List all replication jobs
// @Description	Returns a list of all configured replication jobs with their statuses
// @Tags			replications
// @Accept			json
// @Produce		json
// @Success		200	{array}	object
// @Failure		502	{object}	map[string]interface{}
// @Router			/replications [get]
func (h *ReplicationHandler) ListReplications(c *gin.Context) {
	replications, err := h.replicationService.ListReplications(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, replications)
}

// PauseReplication
// @Summary        Pause a replication job
// @Description    Pauses an active replication job
// @Tags           replications
// @Accept         json
// @Produce        json
// @Param          replication body domain.ReplicationIdentifier true "Replication identifier"
// @Success        200 {string} string "Replication paused successfully"
// @Failure        400 {object} map[string]interface{}
// @Failure        502 {object} map[string]interface{}
// @Router         /replications/pause [post]
func (h *ReplicationHandler) PauseReplication(c *gin.Context) {
	h.replicationAction(c, "pause")
}

// ResumeReplication
// @Summary        Resume a paused replication job
// @Description    Resumes a paused replication job
// @Tags           replications
// @Accept         json
// @Produce        json
// @Param          replication body domain.ReplicationIdentifier true "Replication identifier"
// @Success        200 {string} string "Replication resumed successfully"
// @Failure        400 {object} map[string]interface{}
// @Failure        502 {object} map[string]interface{}
// @Router         /replications/resume [post]
func (h *ReplicationHandler) ResumeReplication(c *gin.Context) {
	h.replicationAction(c, "resume")
}

// DeleteReplication
// @Summary        Delete a replication job
// @Description    Deletes a replication job
// @Tags           replications
// @Accept         json
// @Produce        json
// @Param          replication body domain.ReplicationIdentifier true "Replication identifier"
// @Success        200 {string} string "Replication deleted successfully"
// @Failure        400 {object} map[string]interface{}
// @Failure        502 {object} map[string]interface{}
// @Router         /replications [delete]
func (h *ReplicationHandler) DeleteReplication(c *gin.Context) {
	h.replicationAction(c, "delete")
}

// SwitchZeroDowntime
// @Summary		Switch main and follower buckets without downtime
// @Description	Switches main and follower buckets for a replication job without blocking writes
// @Tags			replications
// @Accept			json
// @Produce		json
// @Param			replication	body		domain.ReplicationIdentifier	true	"Replication identifier"
// @Success		202			{string}	string				"Switch initiated successfully"
// @Failure		400			{object}	map[string]interface{}
// @Failure		502			{object}	map[string]interface{}
// @Router			/replications/switch/zero-downtime [post]
func (h *ReplicationHandler) SwitchZeroDowntime(c *gin.Context) {
	var id domain.ReplicationIdentifier
	if err := c.ShouldBindJSON(&id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	err := h.replicationService.SwitchZeroDowntime(c.Request.Context(), &id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusAccepted)
}

// replicationAction handles pause, resume, and delete actions
func (h *ReplicationHandler) replicationAction(c *gin.Context, action string) {
	var id domain.ReplicationIdentifier
	if err := c.ShouldBindJSON(&id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	var err error
	switch action {
	case "pause":
		err = h.replicationService.PauseReplication(c.Request.Context(), &id)
	case "resume":
		err = h.replicationService.ResumeReplication(c.Request.Context(), &id)
	case "delete":
		err = h.replicationService.DeleteReplication(c.Request.Context(), &id)
	}

	if err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
