package api

import (
	"context"
	"net/http"
	"time"

	pb "github.com/clyso/chorus/proto/gen/go/chorus"
	"github.com/gin-gonic/gin"
	"github.com/hantdev/chorus-controller/internal/worker"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type Handlers struct {
	worker *worker.Client
}

func New(workerClient *worker.Client) *Handlers { return &Handlers{worker: workerClient} }

// Health
// @Summary		Health check
// @Description	Returns the health status of the controller
// @Tags			health
// @Accept			json
// @Produce		json
// @Success		200	{object}	map[string]interface{}
// @Router			/health [get]
func (h *Handlers) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
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
func (h *Handlers) ListStorages(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	client, conn, err := h.worker.Dial(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
		return
	}
	defer func() {
		_ = conn.Close()
	}()
	resp, err := client.GetStorages(ctx, &emptypb.Empty{})
	if err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
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
func (h *Handlers) ListBuckets(c *gin.Context) {
	user := c.Query("user")
	from := c.Query("from")
	to := c.Query("to")
	showReplicated := c.Query("show_replicated") == "true"

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	client, conn, err := h.worker.Dial(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
		return
	}
	defer func() { _ = conn.Close() }()

	req := &pb.ListBucketsForReplicationRequest{
		User:           user,
		From:           from,
		To:             to,
		ShowReplicated: showReplicated,
	}
	resp, err := client.ListBucketsForReplication(ctx, req)
	if err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"available_buckets": resp.Buckets, "replicated_buckets": resp.ReplicatedBuckets})
}

// CreateStorage
// @Summary		Create a new storage
// @Description	Adds a new storage configuration at runtime
// @Tags			storages
// @Accept			json
// @Produce		json
// @Param			storage	body		UpsertStorage	true	"Storage configuration"
// @Success		201		{string}	string			"Storage created successfully"
// @Failure		400		{object}	map[string]interface{}
// @Failure		502		{object}	map[string]interface{}
// @Router			/storages [post]
func (h *Handlers) CreateStorage(c *gin.Context) {
	var in UpsertStorage
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, errJSON(err))
		return
	}
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	client, conn, err := h.worker.Dial(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
		return
	}
	defer func() { _ = conn.Close() }()
	if _, err := client.AddStorage(ctx, mapUpsertToPB(in)); err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
		return
	}
	c.Status(http.StatusCreated)
}

// UpdateStorage
// @Summary		Update an existing storage
// @Description	Updates an existing storage configuration at runtime
// @Tags			storages
// @Accept			json
// @Produce		json
// @Param			name		path		string			true	"Storage name"
// @Param			storage		body		UpsertStorage	true	"Updated storage configuration"
// @Success		200			{string}	string			"Storage updated successfully"
// @Failure		400			{object}	map[string]interface{}
// @Failure		502			{object}	map[string]interface{}
// @Router			/storages/{name} [patch]
func (h *Handlers) UpdateStorage(c *gin.Context) {
	var in UpsertStorage
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, errJSON(err))
		return
	}
	in.Name = c.Param("name")
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	client, conn, err := h.worker.Dial(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
		return
	}
	defer func() { _ = conn.Close() }()
	if _, err := client.UpdateStorage(ctx, mapUpsertToPB(in)); err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
		return
	}
	c.Status(http.StatusOK)
}

// DeleteStorage
// @Summary		Delete a storage
// @Description	Deletes a storage configuration by name
// @Tags			storages
// @Accept			json
// @Produce		json
// @Param			name	path		string	true	"Storage name"
// @Success		200		{string}	string	"Storage deleted successfully"
// @Failure		502		{object}	map[string]interface{}
// @Router			/storages/{name} [delete]
func (h *Handlers) DeleteStorage(c *gin.Context) {
	name := c.Param("name")
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	client, conn, err := h.worker.Dial(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
		return
	}
	defer func() { _ = conn.Close() }()
	if _, err := client.DeleteStorage(ctx, &pb.DeleteStorageRequest{Name: name}); err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
		return
	}
	c.Status(http.StatusOK)
}

// CreateReplication
// @Summary		Create a new replication job
// @Description	Configures a new replication job for specified buckets between storages
// @Tags			replications
// @Accept			json
// @Produce		json
// @Param			replication	body		CreateReplicationRequest	true	"Replication configuration"
// @Success		201			{string}	string				"Replication job created successfully"
// @Failure		400			{object}	map[string]interface{}
// @Failure		502			{object}	map[string]interface{}
// @Router			/replications [post]
func (h *Handlers) CreateReplication(c *gin.Context) {
	var req CreateReplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errJSON(err))
		return
	}
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	client, conn, err := h.worker.Dial(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
		return
	}
	defer func() { _ = conn.Close() }()
	addReq := &pb.AddReplicationRequest{
		User:            req.User,
		From:            req.From,
		To:              req.To,
		Buckets:         req.Buckets,
		IsForAllBuckets: len(req.Buckets) == 0,
		ToBucket:        req.ToBucket,
	}
	if req.AgentURL != "" {
		addReq.AgentUrl = &req.AgentURL
	}
	if _, err := client.AddReplication(ctx, addReq); err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
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
func (h *Handlers) ListReplications(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	client, conn, err := h.worker.Dial(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
		return
	}
	defer func() { _ = conn.Close() }()
	resp, err := client.ListReplications(ctx, &emptypb.Empty{})
	if err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
		return
	}
	c.JSON(http.StatusOK, resp.Replications)
}

// SwitchZeroDowntime
// @Summary		Switch main and follower buckets without downtime
// @Description	Switches main and follower buckets for a replication job without blocking writes
// @Tags			replications
// @Accept			json
// @Produce		json
// @Param			replication	body		object	true	"Replication identifier"
// @Success		202			{string}	string				"Switch initiated successfully"
// @Failure		400			{object}	map[string]interface{}
// @Failure		502			{object}	map[string]interface{}
// @Router			/replications/switch/zero-downtime [post]
func (h *Handlers) SwitchZeroDowntime(c *gin.Context) {
	var in pb.ReplicationRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, errJSON(err))
		return
	}
	if in.ToBucket == "" {
		in.ToBucket = in.Bucket
	}
	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()
	client, conn, err := h.worker.Dial(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
		return
	}
	defer func() { _ = conn.Close() }()
	if _, err = client.SwitchBucketZeroDowntime(ctx, &pb.SwitchBucketZeroDowntimeRequest{ReplicationId: &in}); err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
		return
	}
	c.Status(http.StatusAccepted)
}

// PauseReplication
// @Summary        Pause a replication job
// @Description    Pauses an active replication job
// @Tags           replications
// @Accept         json
// @Produce        json
// @Param          replication body replicationIdent true "Replication identifier"
// @Success        200 {string} string "Replication paused successfully"
// @Failure        400 {object} map[string]interface{}
// @Failure        502 {object} map[string]interface{}
// @Router         /replications/pause [post]
func (h *Handlers) PauseReplication(c *gin.Context) { h.replicationAction(c, "pause") }

// ResumeReplication
// @Summary        Resume a paused replication job
// @Description    Resumes a paused replication job
// @Tags           replications
// @Accept         json
// @Produce        json
// @Param          replication body replicationIdent true "Replication identifier"
// @Success        200 {string} string "Replication resumed successfully"
// @Failure        400 {object} map[string]interface{}
// @Failure        502 {object} map[string]interface{}
// @Router         /replications/resume [post]
func (h *Handlers) ResumeReplication(c *gin.Context) { h.replicationAction(c, "resume") }

// DeleteReplication
// @Summary        Delete a replication job
// @Description    Deletes a replication job
// @Tags           replications
// @Accept         json
// @Produce        json
// @Param          replication body replicationIdent true "Replication identifier"
// @Success        200 {string} string "Replication deleted successfully"
// @Failure        400 {object} map[string]interface{}
// @Failure        502 {object} map[string]interface{}
// @Router         /replications [delete]
func (h *Handlers) DeleteReplication(c *gin.Context) { h.replicationAction(c, "delete") }

func (h *Handlers) replicationAction(c *gin.Context, action string) {
	var id replicationIdent
	if err := c.ShouldBindJSON(&id); err != nil {
		c.JSON(http.StatusBadRequest, errJSON(err))
		return
	}
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	client, conn, err := h.worker.Dial(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, errJSON(err))
		return
	}
	defer func() { _ = conn.Close() }()
	toBucket := id.ToBucket
	if toBucket == "" {
		toBucket = id.Bucket
	}
	req := &pb.ReplicationRequest{User: id.User, Bucket: id.Bucket, From: id.From, To: id.To, ToBucket: toBucket}
	var rpcErr error
	switch action {
	case "pause":
		_, rpcErr = client.PauseReplication(ctx, req)
	case "resume":
		_, rpcErr = client.ResumeReplication(ctx, req)
	case "delete":
		_, rpcErr = client.DeleteReplication(ctx, req)
	}
	if rpcErr != nil {
		c.JSON(http.StatusBadGateway, errJSON(rpcErr))
		return
	}
	c.Status(http.StatusOK)
}
