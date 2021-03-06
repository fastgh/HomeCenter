package requests

import (
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/app/models"
	"net/http"
)

type CreateDownloadTaskRequest struct {
	Url       string `json:"url"`
	Path      string `json:"path"`
	ThreadNum int64  `json:"thread_num"`
}

func (self *CreateDownloadTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	if err := models.DownloadManager.NewTask(self.Url, self.ThreadNum, self.Path); err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.Success()
	}
	return
}

func NewCreateDownloadTaskRequest() *CreateDownloadTaskRequest {
	return &CreateDownloadTaskRequest{}
}

type QueryDownloadTaskRequest struct{}

func (self *QueryDownloadTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	tasks, err := models.DownloadManager.Tasks()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.DataTable(tasks, len(tasks))
	}
	return
}

func NewQueryDownloadTaskRequest() *QueryDownloadTaskRequest {
	return &QueryDownloadTaskRequest{}
}

type RemoveDownloadTaskRequest struct {
	TaskId string `json:"id"`
}

func (self *RemoveDownloadTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	models.DownloadManager.RemoveTask(self.TaskId)
	response.Success()
	return
}

func NewRemoveDownloadTaskRequest() *RemoveDownloadTaskRequest {
	return &RemoveDownloadTaskRequest{}
}

type CancelDownloadTaskRequest struct {
	TaskId string `json:"id"`
}

func (self *CancelDownloadTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	models.DownloadManager.CancelTask(self.TaskId)
	response.Success()
	return
}

func NewCancelDownloadTaskRequest() *CancelDownloadTaskRequest {
	return &CancelDownloadTaskRequest{}
}

type StartDownloadTaskRequest struct {
	TaskId string `json:"id"`
}

func (self *StartDownloadTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	models.DownloadManager.StartTask(self.TaskId)
	response.Success()
	return
}

func NewStartDownloadTaskRequest() *StartDownloadTaskRequest {
	return &StartDownloadTaskRequest{}
}

type GetDownloadSettingsRequest struct{}

func (self *GetDownloadSettingsRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	settings := models.DownloadManager.GetSettings()
	response.SuccessWithDetail(settings)
	return
}

func NewGetDownloadSettingsRequest() *GetDownloadSettingsRequest {
	return &GetDownloadSettingsRequest{}
}

type UpdateDownloadSettingsRequest struct {
	Path          string `json:"path"`           // 下载路径
	ThreadNum     int64  `json:"thread_num"`     // 默认的线程大小
	NotifyOpen    bool   `json:"notify_open"`    // 是否开启通知
	NotifyMessage string `json:"notify_message"` // 通知消息
}

func (self *UpdateDownloadSettingsRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	err := models.DownloadManager.UpdateSettings(self)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.Success()
	}
	return
}

func NewUpdateDownloadSettingsRequest() *UpdateDownloadSettingsRequest {
	return &UpdateDownloadSettingsRequest{}
}
