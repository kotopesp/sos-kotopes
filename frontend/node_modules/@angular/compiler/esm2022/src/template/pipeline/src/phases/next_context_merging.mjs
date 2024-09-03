/**
 * @license
 * Copyright Google LLC All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
import * as o from '../../../../output/output_ast';
import * as ir from '../../ir';
/**
 * Merges logically sequential `NextContextExpr` operations.
 *
 * `NextContextExpr` can be referenced repeatedly, "popping" the runtime's context stack each time.
 * When two such expressions appear back-to-back, it's possible to merge them together into a single
 * `NextContextExpr` that steps multiple contexts. This merging is possible if all conditions are
 * met:
 *
 *   * The result of the `NextContextExpr` that's folded into the subsequent one is not stored (that
 *     is, the call is purely side-effectful).
 *   * No operations in between them uses the implicit context.
 */
export function mergeNextContextExpressions(job) {
    for (const unit of job.units) {
        for (const op of unit.create) {
            if (op.kind === ir.OpKind.Listener || op.kind === ir.OpKind.TwoWayListener) {
                mergeNextContextsInOps(op.handlerOps);
            }
        }
        mergeNextContextsInOps(unit.update);
    }
}
function mergeNextContextsInOps(ops) {
    for (const op of ops) {
        // Look for a candidate operation to maybe merge.
        if (op.kind !== ir.OpKind.Statement ||
            !(op.statement instanceof o.ExpressionStatement) ||
            !(op.statement.expr instanceof ir.NextContextExpr)) {
            continue;
        }
        const mergeSteps = op.statement.expr.steps;
        // Try to merge this `ir.NextContextExpr`.
        let tryToMerge = true;
        for (let candidate = op.next; candidate.kind !== ir.OpKind.ListEnd && tryToMerge; candidate = candidate.next) {
            ir.visitExpressionsInOp(candidate, (expr, flags) => {
                if (!ir.isIrExpression(expr)) {
                    return expr;
                }
                if (!tryToMerge) {
                    // Either we've already merged, or failed to merge.
                    return;
                }
                if (flags & ir.VisitorContextFlag.InChildOperation) {
                    // We cannot merge into child operations.
                    return;
                }
                switch (expr.kind) {
                    case ir.ExpressionKind.NextContext:
                        // Merge the previous `ir.NextContextExpr` into this one.
                        expr.steps += mergeSteps;
                        ir.OpList.remove(op);
                        tryToMerge = false;
                        break;
                    case ir.ExpressionKind.GetCurrentView:
                    case ir.ExpressionKind.Reference:
                    case ir.ExpressionKind.ContextLetReference:
                        // Can't merge past a dependency on the context.
                        tryToMerge = false;
                        break;
                }
                return;
            });
        }
    }
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibmV4dF9jb250ZXh0X21lcmdpbmcuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi9wYWNrYWdlcy9jb21waWxlci9zcmMvdGVtcGxhdGUvcGlwZWxpbmUvc3JjL3BoYXNlcy9uZXh0X2NvbnRleHRfbWVyZ2luZy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiQUFBQTs7Ozs7O0dBTUc7QUFFSCxPQUFPLEtBQUssQ0FBQyxNQUFNLCtCQUErQixDQUFDO0FBQ25ELE9BQU8sS0FBSyxFQUFFLE1BQU0sVUFBVSxDQUFDO0FBSS9COzs7Ozs7Ozs7OztHQVdHO0FBQ0gsTUFBTSxVQUFVLDJCQUEyQixDQUFDLEdBQW1CO0lBQzdELEtBQUssTUFBTSxJQUFJLElBQUksR0FBRyxDQUFDLEtBQUssRUFBRSxDQUFDO1FBQzdCLEtBQUssTUFBTSxFQUFFLElBQUksSUFBSSxDQUFDLE1BQU0sRUFBRSxDQUFDO1lBQzdCLElBQUksRUFBRSxDQUFDLElBQUksS0FBSyxFQUFFLENBQUMsTUFBTSxDQUFDLFFBQVEsSUFBSSxFQUFFLENBQUMsSUFBSSxLQUFLLEVBQUUsQ0FBQyxNQUFNLENBQUMsY0FBYyxFQUFFLENBQUM7Z0JBQzNFLHNCQUFzQixDQUFDLEVBQUUsQ0FBQyxVQUFVLENBQUMsQ0FBQztZQUN4QyxDQUFDO1FBQ0gsQ0FBQztRQUNELHNCQUFzQixDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsQ0FBQztJQUN0QyxDQUFDO0FBQ0gsQ0FBQztBQUVELFNBQVMsc0JBQXNCLENBQUMsR0FBMkI7SUFDekQsS0FBSyxNQUFNLEVBQUUsSUFBSSxHQUFHLEVBQUUsQ0FBQztRQUNyQixpREFBaUQ7UUFDakQsSUFDRSxFQUFFLENBQUMsSUFBSSxLQUFLLEVBQUUsQ0FBQyxNQUFNLENBQUMsU0FBUztZQUMvQixDQUFDLENBQUMsRUFBRSxDQUFDLFNBQVMsWUFBWSxDQUFDLENBQUMsbUJBQW1CLENBQUM7WUFDaEQsQ0FBQyxDQUFDLEVBQUUsQ0FBQyxTQUFTLENBQUMsSUFBSSxZQUFZLEVBQUUsQ0FBQyxlQUFlLENBQUMsRUFDbEQsQ0FBQztZQUNELFNBQVM7UUFDWCxDQUFDO1FBRUQsTUFBTSxVQUFVLEdBQUcsRUFBRSxDQUFDLFNBQVMsQ0FBQyxJQUFJLENBQUMsS0FBSyxDQUFDO1FBRTNDLDBDQUEwQztRQUMxQyxJQUFJLFVBQVUsR0FBRyxJQUFJLENBQUM7UUFDdEIsS0FDRSxJQUFJLFNBQVMsR0FBRyxFQUFFLENBQUMsSUFBSyxFQUN4QixTQUFTLENBQUMsSUFBSSxLQUFLLEVBQUUsQ0FBQyxNQUFNLENBQUMsT0FBTyxJQUFJLFVBQVUsRUFDbEQsU0FBUyxHQUFHLFNBQVMsQ0FBQyxJQUFLLEVBQzNCLENBQUM7WUFDRCxFQUFFLENBQUMsb0JBQW9CLENBQUMsU0FBUyxFQUFFLENBQUMsSUFBSSxFQUFFLEtBQUssRUFBRSxFQUFFO2dCQUNqRCxJQUFJLENBQUMsRUFBRSxDQUFDLGNBQWMsQ0FBQyxJQUFJLENBQUMsRUFBRSxDQUFDO29CQUM3QixPQUFPLElBQUksQ0FBQztnQkFDZCxDQUFDO2dCQUVELElBQUksQ0FBQyxVQUFVLEVBQUUsQ0FBQztvQkFDaEIsbURBQW1EO29CQUNuRCxPQUFPO2dCQUNULENBQUM7Z0JBRUQsSUFBSSxLQUFLLEdBQUcsRUFBRSxDQUFDLGtCQUFrQixDQUFDLGdCQUFnQixFQUFFLENBQUM7b0JBQ25ELHlDQUF5QztvQkFDekMsT0FBTztnQkFDVCxDQUFDO2dCQUVELFFBQVEsSUFBSSxDQUFDLElBQUksRUFBRSxDQUFDO29CQUNsQixLQUFLLEVBQUUsQ0FBQyxjQUFjLENBQUMsV0FBVzt3QkFDaEMseURBQXlEO3dCQUN6RCxJQUFJLENBQUMsS0FBSyxJQUFJLFVBQVUsQ0FBQzt3QkFDekIsRUFBRSxDQUFDLE1BQU0sQ0FBQyxNQUFNLENBQUMsRUFBaUIsQ0FBQyxDQUFDO3dCQUNwQyxVQUFVLEdBQUcsS0FBSyxDQUFDO3dCQUNuQixNQUFNO29CQUNSLEtBQUssRUFBRSxDQUFDLGNBQWMsQ0FBQyxjQUFjLENBQUM7b0JBQ3RDLEtBQUssRUFBRSxDQUFDLGNBQWMsQ0FBQyxTQUFTLENBQUM7b0JBQ2pDLEtBQUssRUFBRSxDQUFDLGNBQWMsQ0FBQyxtQkFBbUI7d0JBQ3hDLGdEQUFnRDt3QkFDaEQsVUFBVSxHQUFHLEtBQUssQ0FBQzt3QkFDbkIsTUFBTTtnQkFDVixDQUFDO2dCQUNELE9BQU87WUFDVCxDQUFDLENBQUMsQ0FBQztRQUNMLENBQUM7SUFDSCxDQUFDO0FBQ0gsQ0FBQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQGxpY2Vuc2VcbiAqIENvcHlyaWdodCBHb29nbGUgTExDIEFsbCBSaWdodHMgUmVzZXJ2ZWQuXG4gKlxuICogVXNlIG9mIHRoaXMgc291cmNlIGNvZGUgaXMgZ292ZXJuZWQgYnkgYW4gTUlULXN0eWxlIGxpY2Vuc2UgdGhhdCBjYW4gYmVcbiAqIGZvdW5kIGluIHRoZSBMSUNFTlNFIGZpbGUgYXQgaHR0cHM6Ly9hbmd1bGFyLmlvL2xpY2Vuc2VcbiAqL1xuXG5pbXBvcnQgKiBhcyBvIGZyb20gJy4uLy4uLy4uLy4uL291dHB1dC9vdXRwdXRfYXN0JztcbmltcG9ydCAqIGFzIGlyIGZyb20gJy4uLy4uL2lyJztcblxuaW1wb3J0IHR5cGUge0NvbXBpbGF0aW9uSm9ifSBmcm9tICcuLi9jb21waWxhdGlvbic7XG5cbi8qKlxuICogTWVyZ2VzIGxvZ2ljYWxseSBzZXF1ZW50aWFsIGBOZXh0Q29udGV4dEV4cHJgIG9wZXJhdGlvbnMuXG4gKlxuICogYE5leHRDb250ZXh0RXhwcmAgY2FuIGJlIHJlZmVyZW5jZWQgcmVwZWF0ZWRseSwgXCJwb3BwaW5nXCIgdGhlIHJ1bnRpbWUncyBjb250ZXh0IHN0YWNrIGVhY2ggdGltZS5cbiAqIFdoZW4gdHdvIHN1Y2ggZXhwcmVzc2lvbnMgYXBwZWFyIGJhY2stdG8tYmFjaywgaXQncyBwb3NzaWJsZSB0byBtZXJnZSB0aGVtIHRvZ2V0aGVyIGludG8gYSBzaW5nbGVcbiAqIGBOZXh0Q29udGV4dEV4cHJgIHRoYXQgc3RlcHMgbXVsdGlwbGUgY29udGV4dHMuIFRoaXMgbWVyZ2luZyBpcyBwb3NzaWJsZSBpZiBhbGwgY29uZGl0aW9ucyBhcmVcbiAqIG1ldDpcbiAqXG4gKiAgICogVGhlIHJlc3VsdCBvZiB0aGUgYE5leHRDb250ZXh0RXhwcmAgdGhhdCdzIGZvbGRlZCBpbnRvIHRoZSBzdWJzZXF1ZW50IG9uZSBpcyBub3Qgc3RvcmVkICh0aGF0XG4gKiAgICAgaXMsIHRoZSBjYWxsIGlzIHB1cmVseSBzaWRlLWVmZmVjdGZ1bCkuXG4gKiAgICogTm8gb3BlcmF0aW9ucyBpbiBiZXR3ZWVuIHRoZW0gdXNlcyB0aGUgaW1wbGljaXQgY29udGV4dC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIG1lcmdlTmV4dENvbnRleHRFeHByZXNzaW9ucyhqb2I6IENvbXBpbGF0aW9uSm9iKTogdm9pZCB7XG4gIGZvciAoY29uc3QgdW5pdCBvZiBqb2IudW5pdHMpIHtcbiAgICBmb3IgKGNvbnN0IG9wIG9mIHVuaXQuY3JlYXRlKSB7XG4gICAgICBpZiAob3Aua2luZCA9PT0gaXIuT3BLaW5kLkxpc3RlbmVyIHx8IG9wLmtpbmQgPT09IGlyLk9wS2luZC5Ud29XYXlMaXN0ZW5lcikge1xuICAgICAgICBtZXJnZU5leHRDb250ZXh0c0luT3BzKG9wLmhhbmRsZXJPcHMpO1xuICAgICAgfVxuICAgIH1cbiAgICBtZXJnZU5leHRDb250ZXh0c0luT3BzKHVuaXQudXBkYXRlKTtcbiAgfVxufVxuXG5mdW5jdGlvbiBtZXJnZU5leHRDb250ZXh0c0luT3BzKG9wczogaXIuT3BMaXN0PGlyLlVwZGF0ZU9wPik6IHZvaWQge1xuICBmb3IgKGNvbnN0IG9wIG9mIG9wcykge1xuICAgIC8vIExvb2sgZm9yIGEgY2FuZGlkYXRlIG9wZXJhdGlvbiB0byBtYXliZSBtZXJnZS5cbiAgICBpZiAoXG4gICAgICBvcC5raW5kICE9PSBpci5PcEtpbmQuU3RhdGVtZW50IHx8XG4gICAgICAhKG9wLnN0YXRlbWVudCBpbnN0YW5jZW9mIG8uRXhwcmVzc2lvblN0YXRlbWVudCkgfHxcbiAgICAgICEob3Auc3RhdGVtZW50LmV4cHIgaW5zdGFuY2VvZiBpci5OZXh0Q29udGV4dEV4cHIpXG4gICAgKSB7XG4gICAgICBjb250aW51ZTtcbiAgICB9XG5cbiAgICBjb25zdCBtZXJnZVN0ZXBzID0gb3Auc3RhdGVtZW50LmV4cHIuc3RlcHM7XG5cbiAgICAvLyBUcnkgdG8gbWVyZ2UgdGhpcyBgaXIuTmV4dENvbnRleHRFeHByYC5cbiAgICBsZXQgdHJ5VG9NZXJnZSA9IHRydWU7XG4gICAgZm9yIChcbiAgICAgIGxldCBjYW5kaWRhdGUgPSBvcC5uZXh0ITtcbiAgICAgIGNhbmRpZGF0ZS5raW5kICE9PSBpci5PcEtpbmQuTGlzdEVuZCAmJiB0cnlUb01lcmdlO1xuICAgICAgY2FuZGlkYXRlID0gY2FuZGlkYXRlLm5leHQhXG4gICAgKSB7XG4gICAgICBpci52aXNpdEV4cHJlc3Npb25zSW5PcChjYW5kaWRhdGUsIChleHByLCBmbGFncykgPT4ge1xuICAgICAgICBpZiAoIWlyLmlzSXJFeHByZXNzaW9uKGV4cHIpKSB7XG4gICAgICAgICAgcmV0dXJuIGV4cHI7XG4gICAgICAgIH1cblxuICAgICAgICBpZiAoIXRyeVRvTWVyZ2UpIHtcbiAgICAgICAgICAvLyBFaXRoZXIgd2UndmUgYWxyZWFkeSBtZXJnZWQsIG9yIGZhaWxlZCB0byBtZXJnZS5cbiAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cblxuICAgICAgICBpZiAoZmxhZ3MgJiBpci5WaXNpdG9yQ29udGV4dEZsYWcuSW5DaGlsZE9wZXJhdGlvbikge1xuICAgICAgICAgIC8vIFdlIGNhbm5vdCBtZXJnZSBpbnRvIGNoaWxkIG9wZXJhdGlvbnMuXG4gICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG5cbiAgICAgICAgc3dpdGNoIChleHByLmtpbmQpIHtcbiAgICAgICAgICBjYXNlIGlyLkV4cHJlc3Npb25LaW5kLk5leHRDb250ZXh0OlxuICAgICAgICAgICAgLy8gTWVyZ2UgdGhlIHByZXZpb3VzIGBpci5OZXh0Q29udGV4dEV4cHJgIGludG8gdGhpcyBvbmUuXG4gICAgICAgICAgICBleHByLnN0ZXBzICs9IG1lcmdlU3RlcHM7XG4gICAgICAgICAgICBpci5PcExpc3QucmVtb3ZlKG9wIGFzIGlyLlVwZGF0ZU9wKTtcbiAgICAgICAgICAgIHRyeVRvTWVyZ2UgPSBmYWxzZTtcbiAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICAgIGNhc2UgaXIuRXhwcmVzc2lvbktpbmQuR2V0Q3VycmVudFZpZXc6XG4gICAgICAgICAgY2FzZSBpci5FeHByZXNzaW9uS2luZC5SZWZlcmVuY2U6XG4gICAgICAgICAgY2FzZSBpci5FeHByZXNzaW9uS2luZC5Db250ZXh0TGV0UmVmZXJlbmNlOlxuICAgICAgICAgICAgLy8gQ2FuJ3QgbWVyZ2UgcGFzdCBhIGRlcGVuZGVuY3kgb24gdGhlIGNvbnRleHQuXG4gICAgICAgICAgICB0cnlUb01lcmdlID0gZmFsc2U7XG4gICAgICAgICAgICBicmVhaztcbiAgICAgICAgfVxuICAgICAgICByZXR1cm47XG4gICAgICB9KTtcbiAgICB9XG4gIH1cbn1cbiJdfQ==