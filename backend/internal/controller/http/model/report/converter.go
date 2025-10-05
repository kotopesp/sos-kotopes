package report

import "github.com/kotopesp/sos-kotopes/internal/core"

func (r *CreateRequestBodyReport) ToCoreReport(userID int) core.Report {
	if r == nil {
		return core.Report{}
	}

	return core.Report{
		UserID:         userID,
		ReportableID:   r.TargetID,
		ReportableType: r.TargetType,
		Reason:         core.ReportReason(r.Reason),
	}
}
