package send_service

import (
	"encoding/json"
	"strings"
	"testing"
)

func validReplyButton(id string) Button {
	return Button{Type: "reply", Id: id, DisplayText: "OK " + id}
}

func validCopyButton() Button {
	return Button{Type: "copy", DisplayText: "Copiar", CopyCode: "PROMO2026"}
}

func validURLButton() Button {
	return Button{Type: "url", DisplayText: "Abrir", URL: "https://example.com"}
}

func validCallButton() Button {
	return Button{Type: "call", DisplayText: "Ligar", PhoneNumber: "+5511999999999"}
}

func validPixButton() Button {
	return Button{Type: "pix", Currency: "BRL", Name: "SmartPedido", KeyType: "email", Key: "pix@example.com"}
}

func testButtonPayload(buttons ...Button) *ButtonStruct {
	return &ButtonStruct{
		Number:      "5511917210173",
		Title:       "Teste",
		Description: "Teste",
		Footer:      "Teste",
		Buttons:     buttons,
	}
}

func TestValidateButtonDataCombinations(t *testing.T) {
	tests := []struct {
		name       string
		buttons    []Button
		wantErr    string
		wantNative string
	}{
		{name: "1 reply", buttons: []Button{validReplyButton("1")}, wantNative: "quick_reply"},
		{name: "2 replies", buttons: []Button{validReplyButton("1"), validReplyButton("2")}, wantNative: "quick_reply"},
		{name: "3 replies", buttons: []Button{validReplyButton("1"), validReplyButton("2"), validReplyButton("3")}, wantNative: "quick_reply"},
		{name: "reply + url", buttons: []Button{validReplyButton("1"), validURLButton()}, wantErr: "reply buttons cannot be mixed"},
		{name: "reply + call", buttons: []Button{validReplyButton("1"), validCallButton()}, wantErr: "reply buttons cannot be mixed"},
		{name: "reply + copy", buttons: []Button{validReplyButton("1"), validCopyButton()}, wantErr: "reply buttons cannot be mixed"},
		{name: "reply + pix", buttons: []Button{validReplyButton("1"), validPixButton()}, wantErr: "reply buttons cannot be mixed"},
		{name: "url", buttons: []Button{validURLButton()}, wantNative: "mixed"},
		{name: "call", buttons: []Button{validCallButton()}, wantNative: "mixed"},
		{name: "copy", buttons: []Button{validCopyButton()}, wantNative: "mixed"},
		{name: "copy + url", buttons: []Button{validCopyButton(), validURLButton()}, wantNative: "mixed"},
		{name: "copy + call", buttons: []Button{validCopyButton(), validCallButton()}, wantNative: "mixed"},
		{name: "url + call", buttons: []Button{validURLButton(), validCallButton()}, wantNative: "mixed"},
		{name: "copy + url + call", buttons: []Button{validCopyButton(), validURLButton(), validCallButton()}, wantNative: "mixed"},
		{name: "pix", buttons: []Button{validPixButton()}, wantNative: "payment_info"},
		{name: "pix + outro botao", buttons: []Button{validPixButton(), validURLButton()}, wantErr: "pix button cannot be combined"},
		{name: "4 replies", buttons: []Button{validReplyButton("1"), validReplyButton("2"), validReplyButton("3"), validReplyButton("4")}, wantErr: "maximum of 3 reply"},
		{name: "unknown type", buttons: []Button{{Type: "bad", DisplayText: "Bad"}}, wantErr: "type must be one of"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan, err := validateButtonData(testButtonPayload(tt.buttons...))
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if plan.NativeFlowName != tt.wantNative {
				t.Fatalf("expected native flow %q, got %q", tt.wantNative, plan.NativeFlowName)
			}
		})
	}
}

func TestBuildNativeFlowButtons(t *testing.T) {
	plan, err := validateButtonData(testButtonPayload(validCopyButton(), validURLButton(), validCallButton()))
	if err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}

	buttons, err := buildNativeFlowButtons(plan)
	if err != nil {
		t.Fatalf("unexpected build error: %v", err)
	}
	if len(buttons) != 3 {
		t.Fatalf("expected 3 native buttons, got %d", len(buttons))
	}

	wantNames := []string{"cta_copy", "cta_url", "cta_call"}
	for i, wantName := range wantNames {
		if buttons[i].GetName() != wantName {
			t.Fatalf("button %d expected name %q, got %q", i, wantName, buttons[i].GetName())
		}
		var params map[string]string
		if err := json.Unmarshal([]byte(buttons[i].GetButtonParamsJSON()), &params); err != nil {
			t.Fatalf("button %d has invalid params json: %v", i, err)
		}
		if params["display_text"] == "" {
			t.Fatalf("button %d missing display_text in params", i)
		}
	}
}
