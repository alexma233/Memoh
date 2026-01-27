package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/containerd/errdefs"
	"github.com/containerd/containerd/v2/pkg/namespaces"
	"github.com/labstack/echo/v4"

	"github.com/memohai/memoh/internal/config"
	ctr "github.com/memohai/memoh/internal/containerd"
)

type ContainerdHandler struct {
	service ctr.Service
	cfg     config.MCPConfig
	namespace string
}

type CreateContainerRequest struct {
	ContainerID string `json:"container_id"`
	Image       string `json:"image,omitempty"`
	Snapshotter string `json:"snapshotter,omitempty"`
}

type CreateContainerResponse struct {
	ContainerID string `json:"container_id"`
	Image       string `json:"image"`
	Snapshotter string `json:"snapshotter"`
	Started     bool   `json:"started"`
}

type CreateSnapshotRequest struct {
	ContainerID  string `json:"container_id"`
	SnapshotName string `json:"snapshot_name"`
}

type CreateSnapshotResponse struct {
	ContainerID  string `json:"container_id"`
	SnapshotName string `json:"snapshot_name"`
	Snapshotter  string `json:"snapshotter"`
}

func NewContainerdHandler(service ctr.Service, cfg config.MCPConfig, namespace string) *ContainerdHandler {
	return &ContainerdHandler{
		service:   service,
		cfg:       cfg,
		namespace: namespace,
	}
}

func (h *ContainerdHandler) Register(e *echo.Echo) {
	group := e.Group("/mcp")
	group.POST("/containers", h.CreateContainer)
	group.DELETE("/containers/:id", h.DeleteContainer)
	group.POST("/snapshots", h.CreateSnapshot)
}

// CreateContainer godoc
// @Summary Create and start MCP container
// @Tags containerd
// @Param payload body CreateContainerRequest true "Create container payload"
// @Success 200 {object} CreateContainerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /mcp/containers [post]
func (h *ContainerdHandler) CreateContainer(c echo.Context) error {
	var req CreateContainerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if strings.TrimSpace(req.ContainerID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "container_id is required")
	}

	image := strings.TrimSpace(req.Image)
	if image == "" {
		image = h.cfg.BusyboxImage
	}
	if image == "" {
		image = config.DefaultBusyboxImg
	}
	snapshotter := strings.TrimSpace(req.Snapshotter)
	if snapshotter == "" {
		snapshotter = h.cfg.Snapshotter
	}
	if snapshotter == "" {
		snapshotter = "overlayfs"
	}

	_, err := h.service.CreateContainer(c.Request().Context(), ctr.CreateContainerRequest{
		ID:          req.ContainerID,
		ImageRef:    image,
		Snapshotter: snapshotter,
	})
	if err != nil && !errdefs.IsAlreadyExists(err) {
		return echo.NewHTTPError(http.StatusInternalServerError, "snapshotter="+snapshotter+" image="+image+" err="+err.Error())
	}

	started := false
	if _, err := h.service.StartTask(c.Request().Context(), req.ContainerID, &ctr.StartTaskOptions{
		UseStdio: false,
	}); err == nil {
		started = true
	}

	return c.JSON(http.StatusOK, CreateContainerResponse{
		ContainerID: req.ContainerID,
		Image:       image,
		Snapshotter: snapshotter,
		Started:     started,
	})
}

// DeleteContainer godoc
// @Summary Delete MCP container
// @Tags containerd
// @Param id path string true "Container ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /mcp/containers/{id} [delete]
func (h *ContainerdHandler) DeleteContainer(c echo.Context) error {
	containerID := strings.TrimSpace(c.Param("id"))
	if containerID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "container id is required")
	}
	_ = h.service.DeleteTask(c.Request().Context(), containerID, &ctr.DeleteTaskOptions{Force: true})
	if err := h.service.DeleteContainer(c.Request().Context(), containerID, &ctr.DeleteContainerOptions{
		CleanupSnapshot: true,
	}); err != nil {
		if errdefs.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, "container not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// CreateSnapshot godoc
// @Summary Create container snapshot
// @Tags containerd
// @Param payload body CreateSnapshotRequest true "Create snapshot payload"
// @Success 200 {object} CreateSnapshotResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /mcp/snapshots [post]
func (h *ContainerdHandler) CreateSnapshot(c echo.Context) error {
	var req CreateSnapshotRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if strings.TrimSpace(req.ContainerID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "container_id is required")
	}
	container, err := h.service.GetContainer(c.Request().Context(), req.ContainerID)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, "container not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	ctx := c.Request().Context()
	if strings.TrimSpace(h.namespace) != "" {
		ctx = namespaces.WithNamespace(ctx, h.namespace)
	}
	info, err := container.Info(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	snapshotName := strings.TrimSpace(req.SnapshotName)
	if snapshotName == "" {
		snapshotName = req.ContainerID + "-" + time.Now().Format("20060102150405")
	}
	if err := h.service.CommitSnapshot(c.Request().Context(), info.Snapshotter, snapshotName, info.SnapshotKey); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, CreateSnapshotResponse{
		ContainerID:  req.ContainerID,
		SnapshotName: snapshotName,
		Snapshotter:  info.Snapshotter,
	})
}
