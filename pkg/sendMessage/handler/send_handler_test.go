package send_handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	instance_model "github.com/evolution-foundation/evolution-go/pkg/instance/model"
	send_service "github.com/evolution-foundation/evolution-go/pkg/sendMessage/service"
	"github.com/gin-gonic/gin"
)

type fakeSendService struct{}

func (f fakeSendService) SendText(data *send_service.TextStruct, instance *instance_model.Instance) (*send_service.MessageSendStruct, error) {
	return &send_service.MessageSendStruct{}, nil
}
func (f fakeSendService) SendLink(data *send_service.LinkStruct, instance *instance_model.Instance) (*send_service.MessageSendStruct, error) {
	return &send_service.MessageSendStruct{}, nil
}
func (f fakeSendService) SendMediaUrl(data *send_service.MediaStruct, instance *instance_model.Instance) (*send_service.MessageSendStruct, error) {
	return &send_service.MessageSendStruct{}, nil
}
func (f fakeSendService) SendMediaFile(data *send_service.MediaStruct, fileData []byte, instance *instance_model.Instance) (*send_service.MessageSendStruct, error) {
	return &send_service.MessageSendStruct{}, nil
}
func (f fakeSendService) SendPoll(data *send_service.PollStruct, instance *instance_model.Instance) (*send_service.MessageSendStruct, error) {
	return &send_service.MessageSendStruct{}, nil
}
func (f fakeSendService) SendSticker(data *send_service.StickerStruct, instance *instance_model.Instance) (*send_service.MessageSendStruct, error) {
	return &send_service.MessageSendStruct{}, nil
}
func (f fakeSendService) SendLocation(data *send_service.LocationStruct, instance *instance_model.Instance) (*send_service.MessageSendStruct, error) {
	return &send_service.MessageSendStruct{}, nil
}
func (f fakeSendService) SendContact(data *send_service.ContactStruct, instance *instance_model.Instance) (*send_service.MessageSendStruct, error) {
	return &send_service.MessageSendStruct{}, nil
}
func (f fakeSendService) SendButton(data *send_service.ButtonStruct, instance *instance_model.Instance) (*send_service.MessageSendStruct, error) {
	return &send_service.MessageSendStruct{}, nil
}
func (f fakeSendService) SendList(data *send_service.ListStruct, instance *instance_model.Instance) (*send_service.MessageSendStruct, error) {
	return &send_service.MessageSendStruct{}, nil
}
func (f fakeSendService) SendCarousel(data *send_service.CarouselStruct, instance *instance_model.Instance) (*send_service.MessageSendStruct, error) {
	return &send_service.MessageSendStruct{}, nil
}
func (f fakeSendService) SendStatusText(data *send_service.StatusTextStruct, instance *instance_model.Instance) (*send_service.MessageSendStruct, error) {
	return &send_service.MessageSendStruct{}, nil
}
func (f fakeSendService) SendStatusMediaUrl(data *send_service.StatusMediaStruct, instance *instance_model.Instance) (*send_service.MessageSendStruct, error) {
	return &send_service.MessageSendStruct{}, nil
}
func (f fakeSendService) SendStatusMediaFile(data *send_service.StatusMediaStruct, fileData []byte, instance *instance_model.Instance) (*send_service.MessageSendStruct, error) {
	return &send_service.MessageSendStruct{}, nil
}

func testSendRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := &sendHandler{sendMessageService: fakeSendService{}}
	router.Use(func(ctx *gin.Context) {
		ctx.Set("instance", &instance_model.Instance{Id: "test-instance", Name: "test-instance"})
		ctx.Next()
	})
	router.POST("/send/text", handler.SendText)
	router.POST("/send/button", handler.SendButton)
	router.POST("/send/list", handler.SendList)
	router.POST("/send/carousel", handler.SendCarousel)
	return router
}

func postJSON(t *testing.T, router http.Handler, path string, payload map[string]interface{}) *httptest.ResponseRecorder {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func TestSendEndpointsAcceptMinimalPayloads(t *testing.T) {
	router := testSendRouter()

	tests := []struct {
		name    string
		path    string
		payload map[string]interface{}
	}{
		{
			name: "send button",
			path: "/send/button",
			payload: map[string]interface{}{
				"number":      "5511917210173",
				"title":       "Teste",
				"description": "Teste",
				"footer":      "Teste",
				"buttons": []map[string]interface{}{
					{"type": "reply", "id": "1", "displayText": "OK"},
				},
			},
		},
		{
			name: "send text",
			path: "/send/text",
			payload: map[string]interface{}{
				"number": "5511917210173",
				"text":   "Teste",
			},
		},
		{
			name: "send list",
			path: "/send/list",
			payload: map[string]interface{}{
				"number":      "5511917210173",
				"title":       "Teste",
				"description": "Teste",
				"footerText":  "Teste",
				"buttonText":  "Abrir",
				"sections": []map[string]interface{}{
					{
						"title": "Opcoes",
						"rows": []map[string]interface{}{
							{"title": "OK", "rowId": "1"},
						},
					},
				},
			},
		},
		{
			name: "send carousel",
			path: "/send/carousel",
			payload: map[string]interface{}{
				"number": "5511917210173",
				"cards": []map[string]interface{}{
					{
						"header": map[string]interface{}{"title": "Teste"},
						"body":   map[string]interface{}{"text": "Teste"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := postJSON(t, router, tt.path, tt.payload)
			if recorder.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
			}
		})
	}
}
