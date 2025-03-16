package report

import "github.com/kotopesp/sos-kotopes/internal/core"

// ToCoreReport - converts CreateRequestBodyReport to core.Report structure.
func (r *CreateRequestBodyReport) ToCoreReport(userID, postID int) core.Report {
	if r == nil {
		return core.Report{}
	}
	return core.Report{
		UserID: userID,
		PostID: postID,
		Reason: core.ReportReason(r.Reason),
	}
}
