/*!
 * @license
 * Copyright Google LLC All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
import * as ir from '../../ir';
/**
 * Removes any `storeLet` calls that aren't referenced outside of the current view.
 */
export function optimizeStoreLet(job) {
    const letUsedExternally = new Set();
    // Since `@let` declarations can be referenced in child views, both in
    // the creation block (via listeners) and in the update block, we have
    // to look through all the ops to find the references.
    for (const unit of job.units) {
        for (const op of unit.ops()) {
            ir.visitExpressionsInOp(op, (expr) => {
                if (expr instanceof ir.ContextLetReferenceExpr) {
                    letUsedExternally.add(expr.target);
                }
            });
        }
    }
    // TODO(crisbeto): potentially remove the unused calls completely, pending discussion.
    for (const unit of job.units) {
        for (const op of unit.update) {
            ir.transformExpressionsInOp(op, (expression) => expression instanceof ir.StoreLetExpr && !letUsedExternally.has(expression.target)
                ? expression.value
                : expression, ir.VisitorContextFlag.None);
        }
    }
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic3RvcmVfbGV0X29wdGltaXphdGlvbi5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uL3BhY2thZ2VzL2NvbXBpbGVyL3NyYy90ZW1wbGF0ZS9waXBlbGluZS9zcmMvcGhhc2VzL3N0b3JlX2xldF9vcHRpbWl6YXRpb24udHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IkFBQUE7Ozs7OztHQU1HO0FBR0gsT0FBTyxLQUFLLEVBQUUsTUFBTSxVQUFVLENBQUM7QUFHL0I7O0dBRUc7QUFDSCxNQUFNLFVBQVUsZ0JBQWdCLENBQUMsR0FBbUI7SUFDbEQsTUFBTSxpQkFBaUIsR0FBRyxJQUFJLEdBQUcsRUFBYSxDQUFDO0lBRS9DLHNFQUFzRTtJQUN0RSxzRUFBc0U7SUFDdEUsc0RBQXNEO0lBQ3RELEtBQUssTUFBTSxJQUFJLElBQUksR0FBRyxDQUFDLEtBQUssRUFBRSxDQUFDO1FBQzdCLEtBQUssTUFBTSxFQUFFLElBQUksSUFBSSxDQUFDLEdBQUcsRUFBRSxFQUFFLENBQUM7WUFDNUIsRUFBRSxDQUFDLG9CQUFvQixDQUFDLEVBQUUsRUFBRSxDQUFDLElBQUksRUFBRSxFQUFFO2dCQUNuQyxJQUFJLElBQUksWUFBWSxFQUFFLENBQUMsdUJBQXVCLEVBQUUsQ0FBQztvQkFDL0MsaUJBQWlCLENBQUMsR0FBRyxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsQ0FBQztnQkFDckMsQ0FBQztZQUNILENBQUMsQ0FBQyxDQUFDO1FBQ0wsQ0FBQztJQUNILENBQUM7SUFFRCxzRkFBc0Y7SUFDdEYsS0FBSyxNQUFNLElBQUksSUFBSSxHQUFHLENBQUMsS0FBSyxFQUFFLENBQUM7UUFDN0IsS0FBSyxNQUFNLEVBQUUsSUFBSSxJQUFJLENBQUMsTUFBTSxFQUFFLENBQUM7WUFDN0IsRUFBRSxDQUFDLHdCQUF3QixDQUN6QixFQUFFLEVBQ0YsQ0FBQyxVQUFVLEVBQUUsRUFBRSxDQUNiLFVBQVUsWUFBWSxFQUFFLENBQUMsWUFBWSxJQUFJLENBQUMsaUJBQWlCLENBQUMsR0FBRyxDQUFDLFVBQVUsQ0FBQyxNQUFNLENBQUM7Z0JBQ2hGLENBQUMsQ0FBQyxVQUFVLENBQUMsS0FBSztnQkFDbEIsQ0FBQyxDQUFDLFVBQVUsRUFDaEIsRUFBRSxDQUFDLGtCQUFrQixDQUFDLElBQUksQ0FDM0IsQ0FBQztRQUNKLENBQUM7SUFDSCxDQUFDO0FBQ0gsQ0FBQyIsInNvdXJjZXNDb250ZW50IjpbIi8qIVxuICogQGxpY2Vuc2VcbiAqIENvcHlyaWdodCBHb29nbGUgTExDIEFsbCBSaWdodHMgUmVzZXJ2ZWQuXG4gKlxuICogVXNlIG9mIHRoaXMgc291cmNlIGNvZGUgaXMgZ292ZXJuZWQgYnkgYW4gTUlULXN0eWxlIGxpY2Vuc2UgdGhhdCBjYW4gYmVcbiAqIGZvdW5kIGluIHRoZSBMSUNFTlNFIGZpbGUgYXQgaHR0cHM6Ly9hbmd1bGFyLmlvL2xpY2Vuc2VcbiAqL1xuXG5pbXBvcnQgKiBhcyBvIGZyb20gJy4uLy4uLy4uLy4uL291dHB1dC9vdXRwdXRfYXN0JztcbmltcG9ydCAqIGFzIGlyIGZyb20gJy4uLy4uL2lyJztcbmltcG9ydCB7Q29tcGlsYXRpb25Kb2J9IGZyb20gJy4uL2NvbXBpbGF0aW9uJztcblxuLyoqXG4gKiBSZW1vdmVzIGFueSBgc3RvcmVMZXRgIGNhbGxzIHRoYXQgYXJlbid0IHJlZmVyZW5jZWQgb3V0c2lkZSBvZiB0aGUgY3VycmVudCB2aWV3LlxuICovXG5leHBvcnQgZnVuY3Rpb24gb3B0aW1pemVTdG9yZUxldChqb2I6IENvbXBpbGF0aW9uSm9iKTogdm9pZCB7XG4gIGNvbnN0IGxldFVzZWRFeHRlcm5hbGx5ID0gbmV3IFNldDxpci5YcmVmSWQ+KCk7XG5cbiAgLy8gU2luY2UgYEBsZXRgIGRlY2xhcmF0aW9ucyBjYW4gYmUgcmVmZXJlbmNlZCBpbiBjaGlsZCB2aWV3cywgYm90aCBpblxuICAvLyB0aGUgY3JlYXRpb24gYmxvY2sgKHZpYSBsaXN0ZW5lcnMpIGFuZCBpbiB0aGUgdXBkYXRlIGJsb2NrLCB3ZSBoYXZlXG4gIC8vIHRvIGxvb2sgdGhyb3VnaCBhbGwgdGhlIG9wcyB0byBmaW5kIHRoZSByZWZlcmVuY2VzLlxuICBmb3IgKGNvbnN0IHVuaXQgb2Ygam9iLnVuaXRzKSB7XG4gICAgZm9yIChjb25zdCBvcCBvZiB1bml0Lm9wcygpKSB7XG4gICAgICBpci52aXNpdEV4cHJlc3Npb25zSW5PcChvcCwgKGV4cHIpID0+IHtcbiAgICAgICAgaWYgKGV4cHIgaW5zdGFuY2VvZiBpci5Db250ZXh0TGV0UmVmZXJlbmNlRXhwcikge1xuICAgICAgICAgIGxldFVzZWRFeHRlcm5hbGx5LmFkZChleHByLnRhcmdldCk7XG4gICAgICAgIH1cbiAgICAgIH0pO1xuICAgIH1cbiAgfVxuXG4gIC8vIFRPRE8oY3Jpc2JldG8pOiBwb3RlbnRpYWxseSByZW1vdmUgdGhlIHVudXNlZCBjYWxscyBjb21wbGV0ZWx5LCBwZW5kaW5nIGRpc2N1c3Npb24uXG4gIGZvciAoY29uc3QgdW5pdCBvZiBqb2IudW5pdHMpIHtcbiAgICBmb3IgKGNvbnN0IG9wIG9mIHVuaXQudXBkYXRlKSB7XG4gICAgICBpci50cmFuc2Zvcm1FeHByZXNzaW9uc0luT3AoXG4gICAgICAgIG9wLFxuICAgICAgICAoZXhwcmVzc2lvbikgPT5cbiAgICAgICAgICBleHByZXNzaW9uIGluc3RhbmNlb2YgaXIuU3RvcmVMZXRFeHByICYmICFsZXRVc2VkRXh0ZXJuYWxseS5oYXMoZXhwcmVzc2lvbi50YXJnZXQpXG4gICAgICAgICAgICA/IGV4cHJlc3Npb24udmFsdWVcbiAgICAgICAgICAgIDogZXhwcmVzc2lvbixcbiAgICAgICAgaXIuVmlzaXRvckNvbnRleHRGbGFnLk5vbmUsXG4gICAgICApO1xuICAgIH1cbiAgfVxufVxuIl19