package refreshsession

import "github.com/kotopesp/sos-kotopes/internal/core"

func (rs *RefreshSession) ToCoreRefreshSession(refreshToken *string) core.RefreshSession {
	var refreshSession core.RefreshSession
	refreshSession.FingerprintHash = rs.Fingerprint

	if refreshToken == nil {
		return refreshSession
	}

	refreshSession.RefreshToken = *refreshToken

	return refreshSession
}
