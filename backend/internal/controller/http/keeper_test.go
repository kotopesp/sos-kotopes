package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/keeper"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateKeeper(t *testing.T) {
	app, _ := newTestApp(t)

	route := "/api/v1/keepers/%d"

	tests := []struct {
		name          string
		KeeperID      int
		UserID        int
		token         string
		keeper        keeper.UpdateKeeper
		mockBehaviour func(keeper.UpdateKeeper) core.UpdateKeeper
		wantCode      int
	}{
		{
			name:     "sucess",
			KeeperID: 1,
			UserID:   1,
			token:    token,
			keeper: keeper.UpdateKeeper{
				Description: &[]string{gofakeit.Sentence(10)}[0],
			},
			// mockBehaviour: func(uk keeper.UpdateKeeper) core.UpdateKeeper {
			// 	coreUpdateKeeper := uk.ToCoreUpdateKeeper()
			// 	dependencies.keeperService.EXPECT().
			// 		UpdateKeeper(mock.Anything, 1, 1, coreUpdateKeeper).
			// 		Return(coreUpdateKeeper, nil).Once()

			// },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// coreKeeper := tt.mockBehaviour(tt.keeper)

			body, err := json.Marshal(tt.keeper)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf(route, tt.KeeperID), bytes.NewBuffer(body))

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))

			res, err := app.Test(req, -1)
			require.NoError(t, err)

			body, err = io.ReadAll(res.Body)
			require.NoError(t, err)

			err = res.Body.Close()
			require.NoError(t, err)

			if tt.wantCode == http.StatusOK {
				var data struct {
					Data keeper.ResponseKeeper `json:"data"`
				}

				err := json.Unmarshal(body, &data)
				require.NoError(t, err)

				//assert.Equal(t, keeper.ToModelResponseKeeper(coreKeeper), data.Data)
			}

			assert.Equal(t, tt.wantCode, res.StatusCode)
		})
	}
}
