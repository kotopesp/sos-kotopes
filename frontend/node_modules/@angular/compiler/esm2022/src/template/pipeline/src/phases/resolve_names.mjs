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
 * Resolves lexical references in views (`ir.LexicalReadExpr`) to either a target variable or to
 * property reads on the top-level component context.
 *
 * Also matches `ir.RestoreViewExpr` expressions with the variables of their corresponding saved
 * views.
 */
export function resolveNames(job) {
    for (const unit of job.units) {
        processLexicalScope(unit, unit.create, null);
        processLexicalScope(unit, unit.update, null);
    }
}
function processLexicalScope(unit, ops, savedView) {
    // Maps names defined in the lexical scope of this template to the `ir.XrefId`s of the variable
    // declarations which represent those values.
    //
    // Since variables are generated in each view for the entire lexical scope (including any
    // identifiers from parent templates) only local variables need be considered here.
    const scope = new Map();
    // Symbols defined within the current scope. They take precedence over ones defined outside.
    const localDefinitions = new Map();
    // First, step through the operations list and:
    // 1) build up the `scope` mapping
    // 2) recurse into any listener functions
    for (const op of ops) {
        switch (op.kind) {
            case ir.OpKind.Variable:
                switch (op.variable.kind) {
                    case ir.SemanticVariableKind.Identifier:
                        if (op.variable.local) {
                            if (localDefinitions.has(op.variable.identifier)) {
                                continue;
                            }
                            localDefinitions.set(op.variable.identifier, op.xref);
                        }
                        else if (scope.has(op.variable.identifier)) {
                            continue;
                        }
                        scope.set(op.variable.identifier, op.xref);
                        break;
                    case ir.SemanticVariableKind.Alias:
                        // This variable represents some kind of identifier which can be used in the template.
                        if (scope.has(op.variable.identifier)) {
                            continue;
                        }
                        scope.set(op.variable.identifier, op.xref);
                        break;
                    case ir.SemanticVariableKind.SavedView:
                        // This variable represents a snapshot of the current view context, and can be used to
                        // restore that context within listener functions.
                        savedView = {
                            view: op.variable.view,
                            variable: op.xref,
                        };
                        break;
                }
                break;
            case ir.OpKind.Listener:
            case ir.OpKind.TwoWayListener:
                // Listener functions have separate variable declarations, so process them as a separate
                // lexical scope.
                processLexicalScope(unit, op.handlerOps, savedView);
                break;
        }
    }
    // Next, use the `scope` mapping to match `ir.LexicalReadExpr` with defined names in the lexical
    // scope. Also, look for `ir.RestoreViewExpr`s and match them with the snapshotted view context
    // variable.
    for (const op of ops) {
        if (op.kind == ir.OpKind.Listener || op.kind === ir.OpKind.TwoWayListener) {
            // Listeners were already processed above with their own scopes.
            continue;
        }
        ir.transformExpressionsInOp(op, (expr) => {
            if (expr instanceof ir.LexicalReadExpr) {
                // `expr` is a read of a name within the lexical scope of this view.
                // Either that name is defined within the current view, or it represents a property from the
                // main component context.
                if (localDefinitions.has(expr.name)) {
                    return new ir.ReadVariableExpr(localDefinitions.get(expr.name));
                }
                else if (scope.has(expr.name)) {
                    // This was a defined variable in the current scope.
                    return new ir.ReadVariableExpr(scope.get(expr.name));
                }
                else {
                    // Reading from the component context.
                    return new o.ReadPropExpr(new ir.ContextExpr(unit.job.root.xref), expr.name);
                }
            }
            else if (expr instanceof ir.RestoreViewExpr && typeof expr.view === 'number') {
                // `ir.RestoreViewExpr` happens in listener functions and restores a saved view from the
                // parent creation list. We expect to find that we captured the `savedView` previously, and
                // that it matches the expected view to be restored.
                if (savedView === null || savedView.view !== expr.view) {
                    throw new Error(`AssertionError: no saved view ${expr.view} from view ${unit.xref}`);
                }
                expr.view = new ir.ReadVariableExpr(savedView.variable);
                return expr;
            }
            else {
                return expr;
            }
        }, ir.VisitorContextFlag.None);
    }
    for (const op of ops) {
        ir.visitExpressionsInOp(op, (expr) => {
            if (expr instanceof ir.LexicalReadExpr) {
                throw new Error(`AssertionError: no lexical reads should remain, but found read of ${expr.name}`);
            }
        });
    }
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicmVzb2x2ZV9uYW1lcy5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uL3BhY2thZ2VzL2NvbXBpbGVyL3NyYy90ZW1wbGF0ZS9waXBlbGluZS9zcmMvcGhhc2VzL3Jlc29sdmVfbmFtZXMudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IkFBQUE7Ozs7OztHQU1HO0FBRUgsT0FBTyxLQUFLLENBQUMsTUFBTSwrQkFBK0IsQ0FBQztBQUNuRCxPQUFPLEtBQUssRUFBRSxNQUFNLFVBQVUsQ0FBQztBQUcvQjs7Ozs7O0dBTUc7QUFDSCxNQUFNLFVBQVUsWUFBWSxDQUFDLEdBQW1CO0lBQzlDLEtBQUssTUFBTSxJQUFJLElBQUksR0FBRyxDQUFDLEtBQUssRUFBRSxDQUFDO1FBQzdCLG1CQUFtQixDQUFDLElBQUksRUFBRSxJQUFJLENBQUMsTUFBTSxFQUFFLElBQUksQ0FBQyxDQUFDO1FBQzdDLG1CQUFtQixDQUFDLElBQUksRUFBRSxJQUFJLENBQUMsTUFBTSxFQUFFLElBQUksQ0FBQyxDQUFDO0lBQy9DLENBQUM7QUFDSCxDQUFDO0FBRUQsU0FBUyxtQkFBbUIsQ0FDMUIsSUFBcUIsRUFDckIsR0FBb0QsRUFDcEQsU0FBMkI7SUFFM0IsK0ZBQStGO0lBQy9GLDZDQUE2QztJQUM3QyxFQUFFO0lBQ0YseUZBQXlGO0lBQ3pGLG1GQUFtRjtJQUNuRixNQUFNLEtBQUssR0FBRyxJQUFJLEdBQUcsRUFBcUIsQ0FBQztJQUUzQyw0RkFBNEY7SUFDNUYsTUFBTSxnQkFBZ0IsR0FBRyxJQUFJLEdBQUcsRUFBcUIsQ0FBQztJQUV0RCwrQ0FBK0M7SUFDL0Msa0NBQWtDO0lBQ2xDLHlDQUF5QztJQUN6QyxLQUFLLE1BQU0sRUFBRSxJQUFJLEdBQUcsRUFBRSxDQUFDO1FBQ3JCLFFBQVEsRUFBRSxDQUFDLElBQUksRUFBRSxDQUFDO1lBQ2hCLEtBQUssRUFBRSxDQUFDLE1BQU0sQ0FBQyxRQUFRO2dCQUNyQixRQUFR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LElBQUksS0FBSyxDQUNiLHFFQUFxRSxJQUFJLENBQUMsSUFBSSxFQUFFLENBQ2pGLENBQUM7WUFDSixDQUFDO1FBQ0gsQ0FBQyxDQUFDLENBQUM7SUFDTCxDQUFDO0FBQ0gsQ0FBQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQGxpY2Vuc2VcbiAqIENvcHlyaWdodCBHb29nbGUgTExDIEFsbCBSaWdodHMgUmVzZXJ2ZWQuXG4gKlxuICogVXNlIG9mIHRoaXMgc291cmNlIGNvZGUgaXMgZ292ZXJuZWQgYnkgYW4gTUlULXN0eWxlIGxpY2Vuc2UgdGhhdCBjYW4gYmVcbiAqIGZvdW5kIGluIHRoZSBMSUNFTlNFIGZpbGUgYXQgaHR0cHM6Ly9hbmd1bGFyLmlvL2xpY2Vuc2VcbiAqL1xuXG5pbXBvcnQgKiBhcyBvIGZyb20gJy4uLy4uLy4uLy4uL291dHB1dC9vdXRwdXRfYXN0JztcbmltcG9ydCAqIGFzIGlyIGZyb20gJy4uLy4uL2lyJztcbmltcG9ydCB7Q29tcGlsYXRpb25Kb2IsIENvbXBpbGF0aW9uVW5pdH0gZnJvbSAnLi4vY29tcGlsYXRpb24nO1xuXG4vKipcbiAqIFJlc29sdmVzIGxleGljYWwgcmVmZXJlbmNlcyBpbiB2aWV3cyAoYGlyLkxleGljYWxSZWFkRXhwcmApIHRvIGVpdGhlciBhIHRhcmdldCB2YXJpYWJsZSBvciB0b1xuICogcHJvcGVydHkgcmVhZHMgb24gdGhlIHRvcC1sZXZlbCBjb21wb25lbnQgY29udGV4dC5cbiAqXG4gKiBBbHNvIG1hdGNoZXMgYGlyLlJlc3RvcmVWaWV3RXhwcmAgZXhwcmVzc2lvbnMgd2l0aCB0aGUgdmFyaWFibGVzIG9mIHRoZWlyIGNvcnJlc3BvbmRpbmcgc2F2ZWRcbiAqIHZpZXdzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gcmVzb2x2ZU5hbWVzKGpvYjogQ29tcGlsYXRpb25Kb2IpOiB2b2lkIHtcbiAgZm9yIChjb25zdCB1bml0IG9mIGpvYi51bml0cykge1xuICAgIHByb2Nlc3NMZXhpY2FsU2NvcGUodW5pdCwgdW5pdC5jcmVhdGUsIG51bGwpO1xuICAgIHByb2Nlc3NMZXhpY2FsU2NvcGUodW5pdCwgdW5pdC51cGRhdGUsIG51bGwpO1xuICB9XG59XG5cbmZ1bmN0aW9uIHByb2Nlc3NMZXhpY2FsU2NvcGUoXG4gIHVuaXQ6IENvbXBpbGF0aW9uVW5pdCxcbiAgb3BzOiBpci5PcExpc3Q8aXIuQ3JlYXRlT3A+IHwgaXIuT3BMaXN0PGlyLlVwZGF0ZU9wPixcbiAgc2F2ZWRWaWV3OiBTYXZlZFZpZXcgfCBudWxsLFxuKTogdm9pZCB7XG4gIC8vIE1hcHMgbmFtZXMgZGVmaW5lZCBpbiB0aGUgbGV4aWNhbCBzY29wZSBvZiB0aGlzIHRlbXBsYXRlIHRvIHRoZSBgaXIuWHJlZklkYHMgb2YgdGhlIHZhcmlhYmxlXG4gIC8vIGRlY2xhcmF0aW9ucyB3aGljaCByZXByZXNlbnQgdGhvc2UgdmFsdWVzLlxuICAvL1xuICAvLyBTaW5jZSB2YXJpYWJsZXMgYXJlIGdlbmVyYXRlZCBpbiBlYWNoIHZpZXcgZm9yIHRoZSBlbnRpcmUgbGV4aWNhbCBzY29wZSAoaW5jbHVkaW5nIGFueVxuICAvLyBpZGVudGlmaWVycyBmcm9tIHBhcmVudCB0ZW1wbGF0ZXMpIG9ubHkgbG9jYWwgdmFyaWFibGVzIG5lZWQgYmUgY29uc2lkZXJlZCBoZXJlLlxuICBjb25zdCBzY29wZSA9IG5ldyBNYXA8c3RyaW5nLCBpci5YcmVmSWQ+KCk7XG5cbiAgLy8gU3ltYm9scyBkZWZpbmVkIHdpdGhpbiB0aGUgY3VycmVudCBzY29wZS4gVGhleSB0YWtlIHByZWNlZGVuY2Ugb3ZlciBvbmVzIGRlZmluZWQgb3V0c2lkZS5cbiAgY29uc3QgbG9jYWxEZWZpbml0aW9ucyA9IG5ldyBNYXA8c3RyaW5nLCBpci5YcmVmSWQ+KCk7XG5cbiAgLy8gRmlyc3QsIHN0ZXAgdGhyb3VnaCB0aGUgb3BlcmF0aW9ucyBsaXN0IGFuZDpcbiAgLy8gMSkgYnVpbGQgdXAgdGhlIGBzY29wZWAgbWFwcGluZ1xuICAvLyAyKSByZWN1cnNlIGludG8gYW55IGxpc3RlbmVyIGZ1bmN0aW9uc1xuICBmb3IgKGNvbnN0IG9wIG9mIG9wcykge1xuICAgIHN3aXRjaCAob3Aua2luZCkge1xuICAgICAgY2FzZSBpci5PcEtpbmQuVmFyaWFibGU6XG4gICAgICAgIHN3aXRjaCAob3AudmFyaWFibGUua2luZCkge1xuICAgICAgICAgIGNhc2UgaXIuU2VtYW50aWNWYXJpYWJsZUtpbmQuSWRlbnRpZmllcjpcbiAgICAgICAgICAgIGlmIChvcC52YXJpYWJsZS5sb2NhbCkge1xuICAgICAgICAgICAgICBpZiAobG9jYWxEZWZpbml0aW9ucy5oYXMob3AudmFyaWFibGUuaWRlbnRpZmllcikpIHtcbiAgICAgICAgICAgICAgICBjb250aW51ZTtcbiAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICBsb2NhbERlZmluaXRpb25zLnNldChvcC52YXJpYWJsZS5pZGVudGlmaWVyLCBvcC54cmVmKTtcbiAgICAgICAgICAgIH0gZWxzZSBpZiAoc2NvcGUuaGFzKG9wLnZhcmlhYmxlLmlkZW50aWZpZXIpKSB7XG4gICAgICAgICAgICAgIGNvbnRpbnVlO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgc2NvcGUuc2V0KG9wLnZhcmlhYmxlLmlkZW50aWZpZXIsIG9wLnhyZWYpO1xuICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgICAgY2FzZSBpci5TZW1hbnRpY1ZhcmlhYmxlS2luZC5BbGlhczpcbiAgICAgICAgICAgIC8vIFRoaXMgdmFyaWFibGUgcmVwcmVzZW50cyBzb21lIGtpbmQgb2YgaWRlbnRpZmllciB3aGljaCBjYW4gYmUgdXNlZCBpbiB0aGUgdGVtcGxhdGUuXG4gICAgICAgICAgICBpZiAoc2NvcGUuaGFzKG9wLnZhcmlhYmxlLmlkZW50aWZpZXIpKSB7XG4gICAgICAgICAgICAgIGNvbnRpbnVlO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgc2NvcGUuc2V0KG9wLnZhcmlhYmxlLmlkZW50aWZpZXIsIG9wLnhyZWYpO1xuICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgICAgY2FzZSBpci5TZW1hbnRpY1ZhcmlhYmxlS2luZC5TYXZlZFZpZXc6XG4gICAgICAgICAgICAvLyBUaGlzIHZhcmlhYmxlIHJlcHJlc2VudHMgYSBzbmFwc2hvdCBvZiB0aGUgY3VycmVudCB2aWV3IGNvbnRleHQsIGFuZCBjYW4gYmUgdXNlZCB0b1xuICAgICAgICAgICAgLy8gcmVzdG9yZSB0aGF0IGNvbnRleHQgd2l0aGluIGxpc3RlbmVyIGZ1bmN0aW9ucy5cbiAgICAgICAgICAgIHNhdmVkVmlldyA9IHtcbiAgICAgICAgICAgICAgdmlldzogb3AudmFyaWFibGUudmlldyxcbiAgICAgICAgICAgICAgdmFyaWFibGU6IG9wLnhyZWYsXG4gICAgICAgICAgICB9O1xuICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgIH1cbiAgICAgICAgYnJlYWs7XG4gICAgICBjYXNlIGlyLk9wS2luZC5MaXN0ZW5lcjpcbiAgICAgIGNhc2UgaXIuT3BLaW5kLlR3b1dheUxpc3RlbmVyOlxuICAgICAgICAvLyBMaXN0ZW5lciBmdW5jdGlvbnMgaGF2ZSBzZXBhcmF0ZSB2YXJpYWJsZSBkZWNsYXJhdGlvbnMsIHNvIHByb2Nlc3MgdGhlbSBhcyBhIHNlcGFyYXRlXG4gICAgICAgIC8vIGxleGljYWwgc2NvcGUuXG4gICAgICAgIHByb2Nlc3NMZXhpY2FsU2NvcGUodW5pdCwgb3AuaGFuZGxlck9wcywgc2F2ZWRWaWV3KTtcbiAgICAgICAgYnJlYWs7XG4gICAgfVxuICB9XG5cbiAgLy8gTmV4dCwgdXNlIHRoZSBgc2NvcGVgIG1hcHBpbmcgdG8gbWF0Y2ggYGlyLkxleGljYWxSZWFkRXhwcmAgd2l0aCBkZWZpbmVkIG5hbWVzIGluIHRoZSBsZXhpY2FsXG4gIC8vIHNjb3BlLiBBbHNvLCBsb29rIGZvciBgaXIuUmVzdG9yZVZpZXdFeHByYHMgYW5kIG1hdGNoIHRoZW0gd2l0aCB0aGUgc25hcHNob3R0ZWQgdmlldyBjb250ZXh0XG4gIC8vIHZhcmlhYmxlLlxuICBmb3IgKGNvbnN0IG9wIG9mIG9wcykge1xuICAgIGlmIChvcC5raW5kID09IGlyLk9wS2luZC5MaXN0ZW5lciB8fCBvcC5raW5kID09PSBpci5PcEtpbmQuVHdvV2F5TGlzdGVuZXIpIHtcbiAgICAgIC8vIExpc3RlbmVycyB3ZXJlIGFscmVhZHkgcHJvY2Vzc2VkIGFib3ZlIHdpdGggdGhlaXIgb3duIHNjb3Blcy5cbiAgICAgIGNvbnRpbnVlO1xuICAgIH1cbiAgICBpci50cmFuc2Zvcm1FeHByZXNzaW9uc0luT3AoXG4gICAgICBvcCxcbiAgICAgIChleHByKSA9PiB7XG4gICAgICAgIGlmIChleHByIGluc3RhbmNlb2YgaXIuTGV4aWNhbFJlYWRFeHByKSB7XG4gICAgICAgICAgLy8gYGV4cHJgIGlzIGEgcmVhZCBvZiBhIG5hbWUgd2l0aGluIHRoZSBsZXhpY2FsIHNjb3BlIG9mIHRoaXMgdmlldy5cbiAgICAgICAgICAvLyBFaXRoZXIgdGhhdCBuYW1lIGlzIGRlZmluZWQgd2l0aGluIHRoZSBjdXJyZW50IHZpZXcsIG9yIGl0IHJlcHJlc2VudHMgYSBwcm9wZXJ0eSBmcm9tIHRoZVxuICAgICAgICAgIC8vIG1haW4gY29tcG9uZW50IGNvbnRleHQuXG4gICAgICAgICAgaWYgKGxvY2FsRGVmaW5pdGlvbnMuaGFzKGV4cHIubmFtZSkpIHtcbiAgICAgICAgICAgIHJldHVybiBuZXcgaXIuUmVhZFZhcmlhYmxlRXhwcihsb2NhbERlZmluaXRpb25zLmdldChleHByLm5hbWUpISk7XG4gICAgICAgICAgfSBlbHNlIGlmIChzY29wZS5oYXMoZXhwci5uYW1lKSkge1xuICAgICAgICAgICAgLy8gVGhpcyB3YXMgYSBkZWZpbmVkIHZhcmlhYmxlIGluIHRoZSBjdXJyZW50IHNjb3BlLlxuICAgICAgICAgICAgcmV0dXJuIG5ldyBpci5SZWFkVmFyaWFibGVFeHByKHNjb3BlLmdldChleHByLm5hbWUpISk7XG4gICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIC8vIFJlYWRpbmcgZnJvbSB0aGUgY29tcG9uZW50IGNvbnRleHQuXG4gICAgICAgICAgICByZXR1cm4gbmV3IG8uUmVhZFByb3BFeHByKG5ldyBpci5Db250ZXh0RXhwcih1bml0LmpvYi5yb290LnhyZWYpLCBleHByLm5hbWUpO1xuICAgICAgICAgIH1cbiAgICAgICAgfSBlbHNlIGlmIChleHByIGluc3RhbmNlb2YgaXIuUmVzdG9yZVZpZXdFeHByICYmIHR5cGVvZiBleHByLnZpZXcgPT09ICdudW1iZXInKSB7XG4gICAgICAgICAgLy8gYGlyLlJlc3RvcmVWaWV3RXhwcmAgaGFwcGVucyBpbiBsaXN0ZW5lciBmdW5jdGlvbnMgYW5kIHJlc3RvcmVzIGEgc2F2ZWQgdmlldyBmcm9tIHRoZVxuICAgICAgICAgIC8vIHBhcmVudCBjcmVhdGlvbiBsaXN0LiBXZSBleHBlY3QgdG8gZmluZCB0aGF0IHdlIGNhcHR1cmVkIHRoZSBgc2F2ZWRWaWV3YCBwcmV2aW91c2x5LCBhbmRcbiAgICAgICAgICAvLyB0aGF0IGl0IG1hdGNoZXMgdGhlIGV4cGVjdGVkIHZpZXcgdG8gYmUgcmVzdG9yZWQuXG4gICAgICAgICAgaWYgKHNhdmVkVmlldyA9PT0gbnVsbCB8fCBzYXZlZFZpZXcudmlldyAhPT0gZXhwci52aWV3KSB7XG4gICAgICAgICAgICB0aHJvdyBuZXcgRXJyb3IoYEFzc2VydGlvbkVycm9yOiBubyBzYXZlZCB2aWV3ICR7ZXhwci52aWV3fSBmcm9tIHZpZXcgJHt1bml0LnhyZWZ9YCk7XG4gICAgICAgICAgfVxuICAgICAgICAgIGV4cHIudmlldyA9IG5ldyBpci5SZWFkVmFyaWFibGVFeHByKHNhdmVkVmlldy52YXJpYWJsZSk7XG4gICAgICAgICAgcmV0dXJuIGV4cHI7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgcmV0dXJuIGV4cHI7XG4gICAgICAgIH1cbiAgICAgIH0sXG4gICAgICBpci5WaXNpdG9yQ29udGV4dEZsYWcuTm9uZSxcbiAgICApO1xuICB9XG5cbiAgZm9yIChjb25zdCBvcCBvZiBvcHMpIHtcbiAgICBpci52aXNpdEV4cHJlc3Npb25zSW5PcChvcCwgKGV4cHIpID0+IHtcbiAgICAgIGlmIChleHByIGluc3RhbmNlb2YgaXIuTGV4aWNhbFJlYWRFeHByKSB7XG4gICAgICAgIHRocm93IG5ldyBFcnJvcihcbiAgICAgICAgICBgQXNzZXJ0aW9uRXJyb3I6IG5vIGxleGljYWwgcmVhZHMgc2hvdWxkIHJlbWFpbiwgYnV0IGZvdW5kIHJlYWQgb2YgJHtleHByLm5hbWV9YCxcbiAgICAgICAgKTtcbiAgICAgIH1cbiAgICB9KTtcbiAgfVxufVxuXG4vKipcbiAqIEluZm9ybWF0aW9uIGFib3V0IGEgYFNhdmVkVmlld2AgdmFyaWFibGUuXG4gKi9cbmludGVyZmFjZSBTYXZlZFZpZXcge1xuICAvKipcbiAgICogVGhlIHZpZXcgYGlyLlhyZWZJZGAgd2hpY2ggd2FzIHNhdmVkIGludG8gdGhpcyB2YXJpYWJsZS5cbiAgICovXG4gIHZpZXc6IGlyLlhyZWZJZDtcblxuICAvKipcbiAgICogVGhlIGBpci5YcmVmSWRgIG9mIHRoZSB2YXJpYWJsZSBpbnRvIHdoaWNoIHRoZSB2aWV3IHdhcyBzYXZlZC5cbiAgICovXG4gIHZhcmlhYmxlOiBpci5YcmVmSWQ7XG59XG4iXX0=