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
 * When inside of a listener, we may need access to one or more enclosing views. Therefore, each
 * view should save the current view, and each listener must have the ability to restore the
 * appropriate view. We eagerly generate all save view variables; they will be optimized away later.
 */
export function saveAndRestoreView(job) {
    for (const unit of job.units) {
        unit.create.prepend([
            ir.createVariableOp(unit.job.allocateXrefId(), {
                kind: ir.SemanticVariableKind.SavedView,
                name: null,
                view: unit.xref,
            }, new ir.GetCurrentViewExpr(), ir.VariableFlags.None),
        ]);
        for (const op of unit.create) {
            if (op.kind !== ir.OpKind.Listener && op.kind !== ir.OpKind.TwoWayListener) {
                continue;
            }
            // Embedded views always need the save/restore view operation.
            let needsRestoreView = unit !== job.root;
            if (!needsRestoreView) {
                for (const handlerOp of op.handlerOps) {
                    ir.visitExpressionsInOp(handlerOp, (expr) => {
                        if (expr instanceof ir.ReferenceExpr || expr instanceof ir.ContextLetReferenceExpr) {
                            // Listeners that reference() a local ref need the save/restore view operation.
                            needsRestoreView = true;
                        }
                    });
                }
            }
            if (needsRestoreView) {
                addSaveRestoreViewOperationToListener(unit, op);
            }
        }
    }
}
function addSaveRestoreViewOperationToListener(unit, op) {
    op.handlerOps.prepend([
        ir.createVariableOp(unit.job.allocateXrefId(), {
            kind: ir.SemanticVariableKind.Context,
            name: null,
            view: unit.xref,
        }, new ir.RestoreViewExpr(unit.xref), ir.VariableFlags.None),
    ]);
    // The "restore view" operation in listeners requires a call to `resetView` to reset the
    // context prior to returning from the listener operation. Find any `return` statements in
    // the listener body and wrap them in a call to reset the view.
    for (const handlerOp of op.handlerOps) {
        if (handlerOp.kind === ir.OpKind.Statement &&
            handlerOp.statement instanceof o.ReturnStatement) {
            handlerOp.statement.value = new ir.ResetViewExpr(handlerOp.statement.value);
        }
    }
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic2F2ZV9yZXN0b3JlX3ZpZXcuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi9wYWNrYWdlcy9jb21waWxlci9zcmMvdGVtcGxhdGUvcGlwZWxpbmUvc3JjL3BoYXNlcy9zYXZlX3Jlc3RvcmVfdmlldy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiQUFBQTs7Ozs7O0dBTUc7QUFFSCxPQUFPLEtBQUssQ0FBQyxNQUFNLCtCQUErQixDQUFDO0FBQ25ELE9BQU8sS0FBSyxFQUFFLE1BQU0sVUFBVSxDQUFDO0FBRy9COzs7O0dBSUc7QUFDSCxNQUFNLFVBQVUsa0JBQWtCLENBQUMsR0FBNEI7SUFDN0QsS0FBSyxNQUFNLElBQUksSUFBSSxHQUFHLENBQUMsS0FBSyxFQUFFLENBQUM7UUFDN0IsSUFBSSxDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUM7WUFDbEIsRUFBRSxDQUFDLGdCQUFnQixDQUNqQixJQUFJLENBQUMsR0FBRyxDQUFDLGNBQWMsRUFBRSxFQUN6QjtnQkFDRSxJQUFJLEVBQUUsRUFBRSxDQUFDLG9CQUFvQixDQUFDLFNBQVM7Z0JBQ3ZDLElBQUksRUFBRSxJQUFJO2dCQUNWLElBQUksRUFBRSxJQUFJLENBQUMsSUFBSTthQUNoQixFQUNELElBQUksRUFBRSxDQUFDLGtCQUFrQixFQUFFLEVBQzNCLEVBQUUsQ0FBQyxhQUFhLENBQUMsSUFBSSxDQUN0QjtTQUNGLENBQUMsQ0FBQztRQUVILEtBQUssTUFBTSxFQUFFLElBQUksSUFBSSxDQUFDLE1BQU0sRUFBRSxDQUFDO1lBQzdCLElBQUksRUFBRSxDQUFDLElBQUksS0FBSyxFQUFFLENBQUMsTUFBTSxDQUFDLFFBQVEsSUFBSSxFQUFFLENBQUMsSUFBSSxLQUFLLEVBQUUsQ0FBQyxNQUFNLENBQUMsY0FBYyxFQUFFLENBQUM7Z0JBQzNFLFNBQVM7WUFDWCxDQUFDO1lBRUQsOERBQThEO1lBQzlELElBQUksZ0JBQWdCLEdBQUcsSUFBSSxLQUFLLEdBQUcsQ0FBQyxJQUFJLENBQUM7WUFFekMsSUFBSSxDQUFDLGdCQUFnQixFQUFFLENBQUM7Z0JBQ3RCLEtBQUssTUFBTSxTQUFTLElBQUksRUFBRSxDQUFDLFVBQVUsRUFBRSxDQUFDO29CQUN0QyxFQUFFLENBQUMsb0JBQW9CLENBQUMsU0FBUyxFQUFFLENBQUMsSUFBSSxFQUFFLEVBQUU7d0JBQzFDLElBQUksSUFBSSxZQUFZLEVBQUUsQ0FBQyxhQUFhLElBQUksSUFBSSxZQUFZLEVBQUUsQ0FBQyx1QkFBdUIsRUFBRSxDQUFDOzRCQUNuRiwrRUFBK0U7NEJBQy9FLGdCQUFnQixHQUFHLElBQUksQ0FBQzt3QkFDMUIsQ0FBQztvQkFDSCxDQUFDLENBQUMsQ0FBQztnQkFDTCxDQUFDO1lBQ0gsQ0FBQztZQUVELElBQUksZ0JBQWdCLEVBQUUsQ0FBQztnQkFDckIscUNBQXFDLENBQUMsSUFBSSxFQUFFLEVBQUUsQ0FBQyxDQUFDO1lBQ2xELENBQUM7UUFDSCxDQUFDO0lBQ0gsQ0FBQztBQUNILENBQUM7QUFFRCxTQUFTLHFDQUFxQyxDQUM1QyxJQUF5QixFQUN6QixFQUF1QztJQUV2QyxFQUFFLENBQUMsVUFBVSxDQUFDLE9BQU8sQ0FBQztRQUNwQixFQUFFLENBQUMsZ0JBQWdCLENBQ2pCLElBQUksQ0FBQyxHQUFHLENBQUMsY0FBYyxFQUFFLEVBQ3pCO1lBQ0UsSUFBSSxFQUFFLEVBQUUsQ0FBQyxvQkFBb0IsQ0FBQyxPQUFPO1lBQ3JDLElBQUksRUFBRSxJQUFJO1lBQ1YsSUFBSSxFQUFFLElBQUksQ0FBQyxJQUFJO1NBQ2hCLEVBQ0QsSUFBSSxFQUFFLENBQUMsZUFBZSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsRUFDakMsRUFBRSxDQUFDLGFBQWEsQ0FBQyxJQUFJLENBQ3RCO0tBQ0YsQ0FBQyxDQUFDO0lBRUgsd0ZBQXdGO0lBQ3hGLDBGQUEwRjtJQUMxRiwrREFBK0Q7SUFDL0QsS0FBSyxNQUFNLFNBQVMsSUFBSSxFQUFFLENBQUMsVUFBVSxFQUFFLENBQUM7UUFDdEMsSUFDRSxTQUFTLENBQUMsSUFBSSxLQUFLLEVBQUUsQ0FBQyxNQUFNLENBQUMsU0FBUztZQUN0QyxTQUFTLENBQUMsU0FBUyxZQUFZLENBQUMsQ0FBQyxlQUFlLEVBQ2hELENBQUM7WUFDRCxTQUFTLENBQUMsU0FBUyxDQUFDLEtBQUssR0FBRyxJQUFJLEVBQUUsQ0FBQyxhQUFhLENBQUMsU0FBUyxDQUFDLFNBQVMsQ0FBQyxLQUFLLENBQUMsQ0FBQztRQUM5RSxDQUFDO0lBQ0gsQ0FBQztBQUNILENBQUMiLCJzb3VyY2VzQ29udGVudCI6WyIvKipcbiAqIEBsaWNlbnNlXG4gKiBDb3B5cmlnaHQgR29vZ2xlIExMQyBBbGwgUmlnaHRzIFJlc2VydmVkLlxuICpcbiAqIFVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlIGlzIGdvdmVybmVkIGJ5IGFuIE1JVC1zdHlsZSBsaWNlbnNlIHRoYXQgY2FuIGJlXG4gKiBmb3VuZCBpbiB0aGUgTElDRU5TRSBmaWxlIGF0IGh0dHBzOi8vYW5ndWxhci5pby9saWNlbnNlXG4gKi9cblxuaW1wb3J0ICogYXMgbyBmcm9tICcuLi8uLi8uLi8uLi9vdXRwdXQvb3V0cHV0X2FzdCc7XG5pbXBvcnQgKiBhcyBpciBmcm9tICcuLi8uLi9pcic7XG5pbXBvcnQgdHlwZSB7Q29tcG9uZW50Q29tcGlsYXRpb25Kb2IsIFZpZXdDb21waWxhdGlvblVuaXR9IGZyb20gJy4uL2NvbXBpbGF0aW9uJztcblxuLyoqXG4gKiBXaGVuIGluc2lkZSBvZiBhIGxpc3RlbmVyLCB3ZSBtYXkgbmVlZCBhY2Nlc3MgdG8gb25lIG9yIG1vcmUgZW5jbG9zaW5nIHZpZXdzLiBUaGVyZWZvcmUsIGVhY2hcbiAqIHZpZXcgc2hvdWxkIHNhdmUgdGhlIGN1cnJlbnQgdmlldywgYW5kIGVhY2ggbGlzdGVuZXIgbXVzdCBoYXZlIHRoZSBhYmlsaXR5IHRvIHJlc3RvcmUgdGhlXG4gKiBhcHByb3ByaWF0ZSB2aWV3LiBXZSBlYWdlcmx5IGdlbmVyYXRlIGFsbCBzYXZlIHZpZXcgdmFyaWFibGVzOyB0aGV5IHdpbGwgYmUgb3B0aW1pemVkIGF3YXkgbGF0ZXIuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBzYXZlQW5kUmVzdG9yZVZpZXcoam9iOiBDb21wb25lbnRDb21waWxhdGlvbkpvYik6IHZvaWQge1xuICBmb3IgKGNvbnN0IHVuaXQgb2Ygam9iLnVuaXRzKSB7XG4gICAgdW5pdC5jcmVhdGUucHJlcGVuZChbXG4gICAgICBpci5jcmVhdGVWYXJpYWJsZU9wPGlyLkNyZWF0ZU9wPihcbiAgICAgICAgdW5pdC5qb2IuYWxsb2NhdGVYcmVmSWQoKSxcbiAgICAgICAge1xuICAgICAgICAgIGtpbmQ6IGlyLlNlbWFudGljVmFyaWFibGVLaW5kLlNhdmVkVmlldyxcbiAgICAgICAgICBuYW1lOiBudWxsLFxuICAgICAgICAgIHZpZXc6IHVuaXQueHJlZixcbiAgICAgICAgfSxcbiAgICAgICAgbmV3IGlyLkdldEN1cnJlbnRWaWV3RXhwcigpLFxuICAgICAgICBpci5WYXJpYWJsZUZsYWdzLk5vbmUsXG4gICAgICApLFxuICAgIF0pO1xuXG4gICAgZm9yIChjb25zdCBvcCBvZiB1bml0LmNyZWF0ZSkge1xuICAgICAgaWYgKG9wLmtpbmQgIT09IGlyLk9wS2luZC5MaXN0ZW5lciAmJiBvcC5raW5kICE9PSBpci5PcEtpbmQuVHdvV2F5TGlzdGVuZXIpIHtcbiAgICAgICAgY29udGludWU7XG4gICAgICB9XG5cbiAgICAgIC8vIEVtYmVkZGVkIHZpZXdzIGFsd2F5cyBuZWVkIHRoZSBzYXZlL3Jlc3RvcmUgdmlldyBvcGVyYXRpb24uXG4gICAgICBsZXQgbmVlZHNSZXN0b3JlVmlldyA9IHVuaXQgIT09IGpvYi5yb290O1xuXG4gICAgICBpZiAoIW5lZWRzUmVzdG9yZVZpZXcpIHtcbiAgICAgICAgZm9yIChjb25zdCBoYW5kbGVyT3Agb2Ygb3AuaGFuZGxlck9wcykge1xuICAgICAgICAgIGlyLnZpc2l0RXhwcmVzc2lvbnNJbk9wKGhhbmRsZXJPcCwgKGV4cHIpID0+IHtcbiAgICAgICAgICAgIGlmIChleHByIGluc3RhbmNlb2YgaXIuUmVmZXJlbmNlRXhwciB8fCBleHByIGluc3RhbmNlb2YgaXIuQ29udGV4dExldFJlZmVyZW5jZUV4cHIpIHtcbiAgICAgICAgICAgICAgLy8gTGlzdGVuZXJzIHRoYXQgcmVmZXJlbmNlKCkgYSBsb2NhbCByZWYgbmVlZCB0aGUgc2F2ZS9yZXN0b3JlIHZpZXcgb3BlcmF0aW9uLlxuICAgICAgICAgICAgICBuZWVkc1Jlc3RvcmVWaWV3ID0gdHJ1ZTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICB9KTtcbiAgICAgICAgfVxuICAgICAgfVxuXG4gICAgICBpZiAobmVlZHNSZXN0b3JlVmlldykge1xuICAgICAgICBhZGRTYXZlUmVzdG9yZVZpZXdPcGVyYXRpb25Ub0xpc3RlbmVyKHVuaXQsIG9wKTtcbiAgICAgIH1cbiAgICB9XG4gIH1cbn1cblxuZnVuY3Rpb24gYWRkU2F2ZVJlc3RvcmVWaWV3T3BlcmF0aW9uVG9MaXN0ZW5lcihcbiAgdW5pdDogVmlld0NvbXBpbGF0aW9uVW5pdCxcbiAgb3A6IGlyLkxpc3RlbmVyT3AgfCBpci5Ud29XYXlMaXN0ZW5lck9wLFxuKSB7XG4gIG9wLmhhbmRsZXJPcHMucHJlcGVuZChbXG4gICAgaXIuY3JlYXRlVmFyaWFibGVPcDxpci5VcGRhdGVPcD4oXG4gICAgICB1bml0LmpvYi5hbGxvY2F0ZVhyZWZJZCgpLFxuICAgICAge1xuICAgICAgICBraW5kOiBpci5TZW1hbnRpY1ZhcmlhYmxlS2luZC5Db250ZXh0LFxuICAgICAgICBuYW1lOiBudWxsLFxuICAgICAgICB2aWV3OiB1bml0LnhyZWYsXG4gICAgICB9LFxuICAgICAgbmV3IGlyLlJlc3RvcmVWaWV3RXhwcih1bml0LnhyZWYpLFxuICAgICAgaXIuVmFyaWFibGVGbGFncy5Ob25lLFxuICAgICksXG4gIF0pO1xuXG4gIC8vIFRoZSBcInJlc3RvcmUgdmlld1wiIG9wZXJhdGlvbiBpbiBsaXN0ZW5lcnMgcmVxdWlyZXMgYSBjYWxsIHRvIGByZXNldFZpZXdgIHRvIHJlc2V0IHRoZVxuICAvLyBjb250ZXh0IHByaW9yIHRvIHJldHVybmluZyBmcm9tIHRoZSBsaXN0ZW5lciBvcGVyYXRpb24uIEZpbmQgYW55IGByZXR1cm5gIHN0YXRlbWVudHMgaW5cbiAgLy8gdGhlIGxpc3RlbmVyIGJvZHkgYW5kIHdyYXAgdGhlbSBpbiBhIGNhbGwgdG8gcmVzZXQgdGhlIHZpZXcuXG4gIGZvciAoY29uc3QgaGFuZGxlck9wIG9mIG9wLmhhbmRsZXJPcHMpIHtcbiAgICBpZiAoXG4gICAgICBoYW5kbGVyT3Aua2luZCA9PT0gaXIuT3BLaW5kLlN0YXRlbWVudCAmJlxuICAgICAgaGFuZGxlck9wLnN0YXRlbWVudCBpbnN0YW5jZW9mIG8uUmV0dXJuU3RhdGVtZW50XG4gICAgKSB7XG4gICAgICBoYW5kbGVyT3Auc3RhdGVtZW50LnZhbHVlID0gbmV3IGlyLlJlc2V0Vmlld0V4cHIoaGFuZGxlck9wLnN0YXRlbWVudC52YWx1ZSk7XG4gICAgfVxuICB9XG59XG4iXX0=