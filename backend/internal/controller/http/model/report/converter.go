package report

import "github.com/kotopesp/sos-kotopes/internal/core"

func (r *CreateRequestBodyReport) ToCoreReport() core.Report {
	if r == nil {
		return core.Report{}
	}
	return core.Report{
		PostID: r.PostID,
		UserID: r.UserID,
		Reason: r.Reason,
	}
}
