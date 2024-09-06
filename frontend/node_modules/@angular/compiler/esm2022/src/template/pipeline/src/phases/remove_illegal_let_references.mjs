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
 * It's not allowed to access a `@let` declaration before it has been defined. This is enforced
 * already via template type checking, however it can trip some of the assertions in the pipeline.
 * E.g. the naming phase can fail because we resolved the variable here, but the variable doesn't
 * exist anymore because the optimization phase removed it since it's invalid. To avoid surfacing
 * confusing errors to users in the case where template type checking isn't running (e.g. in JIT
 * mode) this phase detects illegal forward references and replaces them with `undefined`.
 * Eventually users will see the proper error from the template type checker.
 */
export function removeIllegalLetReferences(job) {
    for (const unit of job.units) {
        for (const op of unit.update) {
            if (op.kind !== ir.OpKind.Variable ||
                op.variable.kind !== ir.SemanticVariableKind.Identifier ||
                !(op.initializer instanceof ir.StoreLetExpr)) {
                continue;
            }
            const name = op.variable.identifier;
            let current = op;
            while (current && current.kind !== ir.OpKind.ListEnd) {
                ir.transformExpressionsInOp(current, (expr) => expr instanceof ir.LexicalReadExpr && expr.name === name ? o.literal(undefined) : expr, ir.VisitorContextFlag.None);
                current = current.prev;
            }
        }
    }
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicmVtb3ZlX2lsbGVnYWxfbGV0X3JlZmVyZW5jZXMuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi9wYWNrYWdlcy9jb21waWxlci9zcmMvdGVtcGxhdGUvcGlwZWxpbmUvc3JjL3BoYXNlcy9yZW1vdmVfaWxsZWdhbF9sZXRfcmVmZXJlbmNlcy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiQUFBQTs7Ozs7O0dBTUc7QUFFSCxPQUFPLEtBQUssQ0FBQyxNQUFNLCtCQUErQixDQUFDO0FBQ25ELE9BQU8sS0FBSyxFQUFFLE1BQU0sVUFBVSxDQUFDO0FBRy9COzs7Ozs7OztHQVFHO0FBQ0gsTUFBTSxVQUFVLDBCQUEwQixDQUFDLEdBQW1CO0lBQzVELEtBQUssTUFBTSxJQUFJLElBQUksR0FBRyxDQUFDLEtBQUssRUFBRSxDQUFDO1FBQzdCLEtBQUssTUFBTSxFQUFFLElBQUksSUFBSSxDQUFDLE1BQU0sRUFBRSxDQUFDO1lBQzdCLElBQ0UsRUFBRSxDQUFDLElBQUksS0FBSyxFQUFFLENBQUMsTUFBTSxDQUFDLFFBQVE7Z0JBQzlCLEVBQUUsQ0FBQyxRQUFRLENBQUMsSUFBSSxLQUFLLEVBQUUsQ0FBQyxvQkFBb0IsQ0FBQyxVQUFVO2dCQUN2RCxDQUFDLENBQUMsRUFBRSxDQUFDLFdBQVcsWUFBWSxFQUFFLENBQUMsWUFBWSxDQUFDLEVBQzVDLENBQUM7Z0JBQ0QsU0FBUztZQUNYLENBQUM7WUFFRCxNQUFNLElBQUksR0FBRyxFQUFFLENBQUMsUUFBUSxDQUFDLFVBQVUsQ0FBQztZQUNwQyxJQUFJLE9BQU8sR0FBdUIsRUFBRSxDQUFDO1lBQ3JDLE9BQU8sT0FBTyxJQUFJLE9BQU8sQ0FBQyxJQUFJLEtBQUssRUFBRSxDQUFDLE1BQU0sQ0FBQyxPQUFPLEVBQUUsQ0FBQztnQkFDckQsRUFBRSxDQUFDLHdCQUF3QixDQUN6QixPQUFPLEVBQ1AsQ0FBQyxJQUFJLEVBQUUsRUFBRSxDQUNQLElBQUksWUFBWSxFQUFFLENBQUMsZUFBZSxJQUFJLElBQUksQ0FBQyxJQUFJLEtBQUssSUFBSSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsT0FBTyxDQUFDLFNBQVMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxJQUFJLEVBQ3hGLEVBQUUsQ0FBQyxrQkFBa0IsQ0FBQyxJQUFJLENBQzNCLENBQUM7Z0JBQ0YsT0FBTyxHQUFHLE9BQU8sQ0FBQyxJQUFJLENBQUM7WUFDekIsQ0FBQztRQUNILENBQUM7SUFDSCxDQUFDO0FBQ0gsQ0FBQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQGxpY2Vuc2VcbiAqIENvcHlyaWdodCBHb29nbGUgTExDIEFsbCBSaWdodHMgUmVzZXJ2ZWQuXG4gKlxuICogVXNlIG9mIHRoaXMgc291cmNlIGNvZGUgaXMgZ292ZXJuZWQgYnkgYW4gTUlULXN0eWxlIGxpY2Vuc2UgdGhhdCBjYW4gYmVcbiAqIGZvdW5kIGluIHRoZSBMSUNFTlNFIGZpbGUgYXQgaHR0cHM6Ly9hbmd1bGFyLmlvL2xpY2Vuc2VcbiAqL1xuXG5pbXBvcnQgKiBhcyBvIGZyb20gJy4uLy4uLy4uLy4uL291dHB1dC9vdXRwdXRfYXN0JztcbmltcG9ydCAqIGFzIGlyIGZyb20gJy4uLy4uL2lyJztcbmltcG9ydCB7Q29tcGlsYXRpb25Kb2J9IGZyb20gJy4uL2NvbXBpbGF0aW9uJztcblxuLyoqXG4gKiBJdCdzIG5vdCBhbGxvd2VkIHRvIGFjY2VzcyBhIGBAbGV0YCBkZWNsYXJhdGlvbiBiZWZvcmUgaXQgaGFzIGJlZW4gZGVmaW5lZC4gVGhpcyBpcyBlbmZvcmNlZFxuICogYWxyZWFkeSB2aWEgdGVtcGxhdGUgdHlwZSBjaGVja2luZywgaG93ZXZlciBpdCBjYW4gdHJpcCBzb21lIG9mIHRoZSBhc3NlcnRpb25zIGluIHRoZSBwaXBlbGluZS5cbiAqIEUuZy4gdGhlIG5hbWluZyBwaGFzZSBjYW4gZmFpbCBiZWNhdXNlIHdlIHJlc29sdmVkIHRoZSB2YXJpYWJsZSBoZXJlLCBidXQgdGhlIHZhcmlhYmxlIGRvZXNuJ3RcbiAqIGV4aXN0IGFueW1vcmUgYmVjYXVzZSB0aGUgb3B0aW1pemF0aW9uIHBoYXNlIHJlbW92ZWQgaXQgc2luY2UgaXQncyBpbnZhbGlkLiBUbyBhdm9pZCBzdXJmYWNpbmdcbiAqIGNvbmZ1c2luZyBlcnJvcnMgdG8gdXNlcnMgaW4gdGhlIGNhc2Ugd2hlcmUgdGVtcGxhdGUgdHlwZSBjaGVja2luZyBpc24ndCBydW5uaW5nIChlLmcuIGluIEpJVFxuICogbW9kZSkgdGhpcyBwaGFzZSBkZXRlY3RzIGlsbGVnYWwgZm9yd2FyZCByZWZlcmVuY2VzIGFuZCByZXBsYWNlcyB0aGVtIHdpdGggYHVuZGVmaW5lZGAuXG4gKiBFdmVudHVhbGx5IHVzZXJzIHdpbGwgc2VlIHRoZSBwcm9wZXIgZXJyb3IgZnJvbSB0aGUgdGVtcGxhdGUgdHlwZSBjaGVja2VyLlxuICovXG5leHBvcnQgZnVuY3Rpb24gcmVtb3ZlSWxsZWdhbExldFJlZmVyZW5jZXMoam9iOiBDb21waWxhdGlvbkpvYik6IHZvaWQge1xuICBmb3IgKGNvbnN0IHVuaXQgb2Ygam9iLnVuaXRzKSB7XG4gICAgZm9yIChjb25zdCBvcCBvZiB1bml0LnVwZGF0ZSkge1xuICAgICAgaWYgKFxuICAgICAgICBvcC5raW5kICE9PSBpci5PcEtpbmQuVmFyaWFibGUgfHxcbiAgICAgICAgb3AudmFyaWFibGUua2luZCAhPT0gaXIuU2VtYW50aWNWYXJpYWJsZUtpbmQuSWRlbnRpZmllciB8fFxuICAgICAgICAhKG9wLmluaXRpYWxpemVyIGluc3RhbmNlb2YgaXIuU3RvcmVMZXRFeHByKVxuICAgICAgKSB7XG4gICAgICAgIGNvbnRpbnVlO1xuICAgICAgfVxuXG4gICAgICBjb25zdCBuYW1lID0gb3AudmFyaWFibGUuaWRlbnRpZmllcjtcbiAgICAgIGxldCBjdXJyZW50OiBpci5VcGRhdGVPcCB8IG51bGwgPSBvcDtcbiAgICAgIHdoaWxlIChjdXJyZW50ICYmIGN1cnJlbnQua2luZCAhPT0gaXIuT3BLaW5kLkxpc3RFbmQpIHtcbiAgICAgICAgaXIudHJhbnNmb3JtRXhwcmVzc2lvbnNJbk9wKFxuICAgICAgICAgIGN1cnJlbnQsXG4gICAgICAgICAgKGV4cHIpID0+XG4gICAgICAgICAgICBleHByIGluc3RhbmNlb2YgaXIuTGV4aWNhbFJlYWRFeHByICYmIGV4cHIubmFtZSA9PT0gbmFtZSA/IG8ubGl0ZXJhbCh1bmRlZmluZWQpIDogZXhwcixcbiAgICAgICAgICBpci5WaXNpdG9yQ29udGV4dEZsYWcuTm9uZSxcbiAgICAgICAgKTtcbiAgICAgICAgY3VycmVudCA9IGN1cnJlbnQucHJldjtcbiAgICAgIH1cbiAgICB9XG4gIH1cbn1cbiJdfQ==