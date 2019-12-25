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

type JobsHandler struct {
	AuthRepo  repos.AuthRepo
	JobRepo   repos.JobRepo
	Scheduler scheduler.Scheduler
}

func NewJobsHandler(auth repos.AuthRepo, r repos.JobRepo, s scheduler.Scheduler) *JobsHandler {
	return &JobsHandler{auth, r, s}
}

func (h *JobsHandler) BasicMiddleware(c *gin.Context) {
	traceId := utils.NewRandomHex(10)
	c.Set(CtxKeyRequestId, traceId)
	log.Trace(traceId).Infof("Method:%s URI:%s", c.Request.Method, c.Request.URL)
}

func (h *JobsHandler) JobsIndex(c *gin.Context) {
	traceId := c.GetString(CtxKeyRequestId)
	status := c.Request.URL.Query().Get("status")
	expression := c.Request.URL.Query().Get("expression")
	query := models.NewQuery(status, expression)

	ents, err := h.JobRepo.Search(query)
	if err != nil {
		log.Trace(traceId).Errorf("查询定时任务失败:%v", err)
		render.Response(c.Writer, http.StatusPreconditionFailed, err)
		return
	}

	render.JSON(c.Writer, http.StatusOK, ents)
}

func (h *JobsHandler) JobsCreate(c *gin.Context) {
	traceId := c.GetString(CtxKeyRequestId)
	rawBody, _ := httputil.DumpRequest(c.Request, true)
	log.Trace(traceId).Infof("用户请求:%s", rawBody)

	ent := models.NewCronJob()
	if err := json.NewDecoder(c.Request.Body).Decode(ent); err != nil {
		log.Trace(traceId).Warnf("解析请求体失败:%v", err)
		render.Response(c.Writer, http.StatusBadRequest, "Failed to decode request body:"+err.Error())
		return
	}

	if errors, valid := ent.Validate(); !valid {
		log.Trace(traceId).Warnf("请求体不合法:%v", errors)
		render.Response(c.Writer, http.StatusBadRequest, errors)
		return
	}

	if err := h.JobRepo.Create(ent); err != nil {
		log.Trace(traceId).Errorf("存储定时任务到数据库失败:%v", err)
		render.Response(c.Writer, http.StatusUnprocessableEntity, "An error occurred during creating ent:"+err.Error())
		return
	}

	if err := h.Scheduler.Create(ent); err != nil {
		log.Trace(traceId).Errorf("启动定时任务失败:%v", err)
		render.Response(c.Writer, http.StatusInternalServerError, "An error occurred during scheduling ent:"+err.Error())
		return
	}

	log.Trace(traceId).Infof("定时任务创建成功:%+v", ent)
	render.JSON(c.Writer, http.StatusCreated, ent)
}

func (h *JobsHandler) JobsDetail(c *gin.Context) {
	traceId := c.GetString(CtxKeyRequestId)
	params := c.Request.Context().Value(violetear.ParamsKey).(violetear.Params)
	id, err := strconv.Atoi(params[":id"].(string))
	if err != nil {
		log.Trace(traceId).Warnf("解析请求参数错误:%v", err)
		render.Response(c.Writer, http.StatusBadRequest, "Missing param :id err:"+err.Error())
		return
	}

	ent, err := h.JobRepo.FindById(id)
	if err != nil {
		log.Trace(traceId).Errorf("未找到任务:%v", err)
		render.Response(c.Writer, http.StatusNotFound, "CronJob not found:"+err.Error())
		return
	}

	render.JSON(c.Writer, http.StatusOK, ent)
}

func (h *JobsHandler) JobsUpdate(c *gin.Context) {
	traceId := c.GetString(CtxKeyRequestId)
	params := c.Request.Context().Value(violetear.ParamsKey).(violetear.Params)
	id, err := strconv.Atoi(params[":id"].(string))
	if err != nil {
		log.Trace(traceId).Warnf("解析请求参数错误:%v", err)
		render.Response(c.Writer, http.StatusBadRequest, "Missing param :id err:"+err.Error())
		return
	}

	ent, err := h.JobRepo.FindById(id)
	if err != nil {
		log.Trace(traceId).Errorf("未找到任务:%v", err)
		render.Response(c.Writer, http.StatusNotFound, "CronJob not found:"+err.Error())
		return
	}

	newJob := models.NewCronJob()
	if err := json.NewDecoder(c.Request.Body).Decode(newJob); err != nil {
		log.Trace(traceId).Warnf("解析请求体失败:%v", err)
		render.Response(c.Writer, http.StatusBadRequest, "Failed to decode request body:"+err.Error())
		return
	}

	if errors, valid := newJob.Validate(); !valid {
		log.Trace(traceId).Warnf("请求体不合法:%v", err)
		render.Response(c.Writer, http.StatusBadRequest, errors)
		return
	}

	ent.SetAttributes(newJob)
	if err := h.JobRepo.Update(ent); err != nil {
		log.Trace(traceId).Errorf("更新数据库中的定时任务失败:%v", err)
		render.Response(c.Writer, http.StatusUnprocessableEntity, "An error occurred during updating ent:"+err.Error())
		return
	}

	if err := h.Scheduler.Update(ent); err != nil {
		log.Trace(traceId).Errorf("更新内存中的定时任务失败:%v", err)
		render.Response(c.Writer, http.StatusInternalServerError, "An error occurred during scheduling ent:"+err.Error())
		return
	}

	render.JSON(c.Writer, http.StatusOK, ent)
}

func (h *JobsHandler) JobsDelete(c *gin.Context) {
	traceId := c.GetString(CtxKeyRequestId)
	params := c.Request.Context().Value(violetear.ParamsKey).(violetear.Params)
	id, err := strconv.Atoi(params[":id"].(string))
	if err != nil {
		log.Trace(traceId).Warnf("解析请求参数错误:%v", err)
		render.Response(c.Writer, http.StatusBadRequest, "Missing param :id")
		return
	}

	ent, err := h.JobRepo.FindById(id)
	if err != nil {
		log.Trace(traceId).Errorf("未找到任务:%v", err)
		render.Response(c.Writer, http.StatusNotFound, "CronJob not found")
		return
	}

	if err := h.JobRepo.Delete(ent); err != nil {
		log.Trace(traceId).Errorf("删除数据库中的定时任务失败:%v", err)
		render.Response(c.Writer, http.StatusUnprocessableEntity, "An error occurred during deleting ent")
		return
	}

	if err := h.Scheduler.Delete(ent.Id); err != nil {
		log.Trace(traceId).Errorf("删除数据库中的定时任务失败:%v", err)
		render.Response(c.Writer, http.StatusInternalServerError, "An error occurred during deleting scheduled ent")
		return
	}

	render.JSON(c.Writer, http.StatusNoContent, nil)
}
