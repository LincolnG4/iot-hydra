package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LincolnG4/iot-hydra/internal/runtimer"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func newTestRouter(m *runtimer.MockPodmanManager) *gin.Engine {
	gin.SetMode(gin.TestMode)
	logger := zerolog.Nop()
	app := &application{
		PodmanRuntime: m,
		logger:        &logger,
		config:        &config{Addr: ":0"},
	}
	return app.mount()
}

func TestCreateContainer_Success(t *testing.T) {
	m := new(runtimer.MockPodmanManager)
	r := newTestRouter(m)

	payload := newContainerPayload{Name: "foo", Image: "docker/bar"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/v1/containers/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	m.On("CreateContainer", runtimer.Container{Name: "foo", Image: "docker/bar"}).Return(nil).Once()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateContainer_BadJSON(t *testing.T) {
	m := new(runtimer.MockPodmanManager)
	r := newTestRouter(m)

	req := httptest.NewRequest(http.MethodPost, "/v1/containers/", bytes.NewBufferString("{"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCheckContainer_Success(t *testing.T) {
	m := new(runtimer.MockPodmanManager)
	r := newTestRouter(m)

	m.On("CheckContainer", "alpha").Return(runtimer.Container{Name: "alpha", Image: "img", State: "running"}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/v1/containers/alpha", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestStartStopDelete_Success(t *testing.T) {
	m := new(runtimer.MockPodmanManager)
	r := newTestRouter(m)

	m.On("StartContainer", "alpha").Return(nil).Once()
	reqStart := httptest.NewRequest(http.MethodPost, "/v1/containers/alpha/start", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqStart)
	assert.Equal(t, http.StatusCreated, w.Code)

	m.On("StopContainer", "alpha").Return(nil).Once()
	reqStop := httptest.NewRequest(http.MethodPost, "/v1/containers/alpha/stop", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, reqStop)
	assert.Equal(t, http.StatusCreated, w.Code)

	m.On("DeleteContainer", "alpha").Return(nil).Once()
	reqDel := httptest.NewRequest(http.MethodDelete, "/v1/containers/alpha", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, reqDel)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestListContainers_Success(t *testing.T) {
	m := new(runtimer.MockPodmanManager)
	r := newTestRouter(m)

	list := []runtimer.Container{{Name: "a", Image: "i1", State: "running"}, {Name: "b", Image: "i2", State: "exited"}}
	m.On("ListContainers").Return(list, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/v1/containers/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}
