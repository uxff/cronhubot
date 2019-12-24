package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nbari/violetear"
	"github.com/uxff/cronhubot/pkg/models"
	"github.com/uxff/cronhubot/pkg/render"
	"github.com/uxff/cronhubot/pkg/repos"
	"github.com/uxff/cronhubot/pkg/scheduler"
	"github.com/uxff/cronhubot/pkg/utils"
	"github.com/uxff/cronhubot/pkg/log"
)

const (
	CtxKeyRequestId = "__requestId"
)

type EventsHandler struct {
	AuthRepo  repos.AuthRepo
	EventRepo repos.EventRepo
	Scheduler scheduler.Scheduler
}

func NewEventsHandler(auth repos.AuthRepo, r repos.EventRepo, s scheduler.Scheduler) *EventsHandler {
	return &EventsHandler{auth, r, s}
}

func (h *EventsHandler) BasicMiddleware(c *gin.Context) {
	traceId := utils.NewRandomHex(10)
	c.Set(CtxKeyRequestId, traceId)
	log.Trace(traceId).Infof("Method:%s URI:%s", c.Request.Method, c.Request.URL)
}

func (h *EventsHandler) EventsIndex(c *gin.Context) {
	traceId := c.GetString(CtxKeyRequestId)
	status := c.Request.URL.Query().Get("status")
	expression := c.Request.URL.Query().Get("expression")
	query := models.NewQuery(status, expression)

	events, err := h.EventRepo.Search(query)
	if err != nil {
		log.Trace(traceId).Errorf("查询定时任务失败:%v", err)
		render.Response(c.Writer, http.StatusPreconditionFailed, err)
		return
	}

	render.JSON(c.Writer, http.StatusOK, events)
}

func (h *EventsHandler) EventsCreate(c *gin.Context) {
	traceId := c.GetString(CtxKeyRequestId)
	rawBody, _ := httputil.DumpRequest(c.Request, true)
	log.Trace(traceId).Infof("用户请求:%s", rawBody)

	event := models.NewEvent()
	if err := json.NewDecoder(c.Request.Body).Decode(event); err != nil {
		log.Trace(traceId).Warnf("解析请求体失败:%v", err)
		render.Response(c.Writer, http.StatusBadRequest, "Failed to decode request body:"+err.Error())
		return
	}

	if errors, valid := event.Validate(); !valid {
		log.Trace(traceId).Warnf("请求体不合法:%v", errors)
		render.Response(c.Writer, http.StatusBadRequest, errors)
		return
	}

	if err := h.EventRepo.Create(event); err != nil {
		log.Trace(traceId).Errorf("存储定时任务到数据库失败:%v", err)
		render.Response(c.Writer, http.StatusUnprocessableEntity, "An error occurred during creating event:"+err.Error())
		return
	}

	if err := h.Scheduler.Create(event); err != nil {
		log.Trace(traceId).Errorf("启动定时任务失败:%v", err)
		render.Response(c.Writer, http.StatusInternalServerError, "An error occurred during scheduling event:"+err.Error())
		return
	}

	log.Trace(traceId).Infof("定时任务创建成功:%+v", event)
	render.JSON(c.Writer, http.StatusCreated, event)
}

func (h *EventsHandler) EventsShow(c *gin.Context) {
	traceId := c.GetString(CtxKeyRequestId)
	params := c.Request.Context().Value(violetear.ParamsKey).(violetear.Params)
	id, err := strconv.Atoi(params[":id"].(string))
	if err != nil {
		log.Trace(traceId).Warnf("解析请求参数错误:%v", err)
		render.Response(c.Writer, http.StatusBadRequest, "Missing param :id err:"+err.Error())
		return
	}

	event, err := h.EventRepo.FindById(id)
	if err != nil {
		log.Trace(traceId).Errorf("未找到任务:%v", err)
		render.Response(c.Writer, http.StatusNotFound, "Event not found:"+err.Error())
		return
	}

	render.JSON(c.Writer, http.StatusOK, event)
}

func (h *EventsHandler) EventsUpdate(c *gin.Context) {
	traceId := c.GetString(CtxKeyRequestId)
	params := c.Request.Context().Value(violetear.ParamsKey).(violetear.Params)
	id, err := strconv.Atoi(params[":id"].(string))
	if err != nil {
		log.Trace(traceId).Warnf("解析请求参数错误:%v", err)
		render.Response(c.Writer, http.StatusBadRequest, "Missing param :id err:"+err.Error())
		return
	}

	event, err := h.EventRepo.FindById(id)
	if err != nil {
		log.Trace(traceId).Errorf("未找到任务:%v", err)
		render.Response(c.Writer, http.StatusNotFound, "Event not found:"+err.Error())
		return
	}

	newEvent := models.NewEvent()
	if err := json.NewDecoder(c.Request.Body).Decode(newEvent); err != nil {
		log.Trace(traceId).Warnf("解析请求体失败:%v", err)
		render.Response(c.Writer, http.StatusBadRequest, "Failed to decode request body:"+err.Error())
		return
	}

	if errors, valid := newEvent.Validate(); !valid {
		log.Trace(traceId).Warnf("请求体不合法:%v", err)
		render.Response(c.Writer, http.StatusBadRequest, errors)
		return
	}

	event.SetAttributes(newEvent)
	if err := h.EventRepo.Update(event); err != nil {
		log.Trace(traceId).Errorf("更新数据库中的定时任务失败:%v", err)
		render.Response(c.Writer, http.StatusUnprocessableEntity, "An error occurred during updating event:"+err.Error())
		return
	}

	if err := h.Scheduler.Update(event); err != nil {
		log.Trace(traceId).Errorf("更新内存中的定时任务失败:%v", err)
		render.Response(c.Writer, http.StatusInternalServerError, "An error occurred during scheduling event:"+err.Error())
		return
	}

	render.JSON(c.Writer, http.StatusOK, event)
}

func (h *EventsHandler) EventsDelete(c *gin.Context) {
	traceId := c.GetString(CtxKeyRequestId)
	params := c.Request.Context().Value(violetear.ParamsKey).(violetear.Params)
	id, err := strconv.Atoi(params[":id"].(string))
	if err != nil {
		log.Trace(traceId).Warnf("解析请求参数错误:%v", err)
		render.Response(c.Writer, http.StatusBadRequest, "Missing param :id")
		return
	}

	event, err := h.EventRepo.FindById(id)
	if err != nil {
		log.Trace(traceId).Errorf("未找到任务:%v", err)
		render.Response(c.Writer, http.StatusNotFound, "Event not found")
		return
	}

	if err := h.EventRepo.Delete(event); err != nil {
		log.Trace(traceId).Errorf("删除数据库中的定时任务失败:%v", err)
		render.Response(c.Writer, http.StatusUnprocessableEntity, "An error occurred during deleting event")
		return
	}

	if err := h.Scheduler.Delete(event.Id); err != nil {
		log.Trace(traceId).Errorf("删除数据库中的定时任务失败:%v", err)
		render.Response(c.Writer, http.StatusInternalServerError, "An error occurred during deleting scheduled event")
		return
	}

	render.JSON(c.Writer, http.StatusNoContent, nil)
}
