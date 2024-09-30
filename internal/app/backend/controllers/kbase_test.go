package controllers

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend/mclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend/trsclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend/vdatabase"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestController_KnowledgeBase(t *testing.T) {
	type args struct {
		statusCode    int
		response      string
		requestString string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "mock",
			args: args{
				statusCode:    http.StatusOK,
				response:      `Welcome!`,
				requestString: `{"question": "Welcome!"}`,
			},
		},
	}

	controller := Controller{
		TRSParserClient: &trsclient.Mock{},
		ModelClient:     &mclient.Mock{},
		VectorDatabase:  &vdatabase.Mock{},
	}

	router := httprouter.New()
	router.GET("/knowledge_base", controller.KnowledgeBase)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			requestBody := strings.NewReader(tt.args.requestString)
			request, err := http.NewRequest("GET", "/knowledge_base", requestBody)

			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			require.Equal(t, tt.args.statusCode, recorder.Code)

			require.Equal(t, tt.args.response, recorder.Body.String())
		})
	}
}
