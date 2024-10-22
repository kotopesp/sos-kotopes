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
 * Generate a preamble sequence for each view creation block and listener function which declares
 * any variables that be referenced in other operations in the block.
 *
 * Variables generated include:
 *   * a saved view context to be used to restore the current view in event listeners.
 *   * the context of the restored view within event listener handlers.
 *   * context variables from the current view as well as all parent views (including the root
 *     context if needed).
 *   * local references from elements within the current view and any lexical parents.
 *
 * Variables are generated here unconditionally, and may optimized away in future operations if it
 * turns out their values (and any side effects) are unused.
 */
export function generateVariables(job) {
    recursivelyProcessView(job.root, /* there is no parent scope for the root view */ null);
}
/**
 * Process the given `ViewCompilation` and generate preambles for it and any listeners that it
 * declares.
 *
 * @param `parentScope` a scope extracted from the parent view which captures any variables which
 *     should be inherited by this view. `null` if the current view is the root view.
 */
function recursivelyProcessView(view, parentScope) {
    // Extract a `Scope` from this view.
    const scope = getScopeForView(view, parentScope);
    for (const op of view.create) {
        switch (op.kind) {
            case ir.OpKind.Template:
                // Descend into child embedded views.
                recursivelyProcessView(view.job.views.get(op.xref), scope);
                break;
            case ir.OpKind.Projection:
                if (op.fallbackView !== null) {
                    recursivelyProcessView(view.job.views.get(op.fallbackView), scope);
                }
                break;
            case ir.OpKind.RepeaterCreate:
                // Descend into child embedded views.
                recursivelyProcessView(view.job.views.get(op.xref), scope);
                if (op.emptyView) {
                    recursivelyProcessView(view.job.views.get(op.emptyView), scope);
                }
                break;
            case ir.OpKind.Listener:
            case ir.OpKind.TwoWayListener:
                // Prepend variables to listener handler functions.
                op.handlerOps.prepend(generateVariablesInScopeForView(view, scope, true));
                break;
        }
    }
    view.update.prepend(generateVariablesInScopeForView(view, scope, false));
}
/**
 * Process a view and generate a `Scope` representing the variables available for reference within
 * that view.
 */
function getScopeForView(view, parent) {
    const scope = {
        view: view.xref,
        viewContextVariable: {
            kind: ir.SemanticVariableKind.Context,
            name: null,
            view: view.xref,
        },
        contextVariables: new Map(),
        aliases: view.aliases,
        references: [],
        letDeclarations: [],
        parent,
    };
    for (const identifier of view.contextVariables.keys()) {
        scope.contextVariables.set(identifier, {
            kind: ir.SemanticVariableKind.Identifier,
            name: null,
            identifier,
            local: false,
        });
    }
    for (const op of view.create) {
        switch (op.kind) {
            case ir.OpKind.ElementStart:
            case ir.OpKind.Template:
                if (!Array.isArray(op.localRefs)) {
                    throw new Error(`AssertionError: expected localRefs to be an array`);
                }
                // Record available local references from this element.
                for (let offset = 0; offset < op.localRefs.length; offset++) {
                    scope.references.push({
                        name: op.localRefs[offset].name,
                        targetId: op.xref,
                        targetSlot: op.handle,
                        offset,
                        variable: {
                            kind: ir.SemanticVariableKind.Identifier,
                            name: null,
                            identifier: op.localRefs[offset].name,
                            local: false,
                        },
                    });
                }
                break;
            case ir.OpKind.DeclareLet:
                scope.letDeclarations.push({
                    targetId: op.xref,
                    targetSlot: op.handle,
                    variable: {
                        kind: ir.SemanticVariableKind.Identifier,
                        name: null,
                        identifier: op.declaredName,
                        local: false,
                    },
                });
                break;
        }
    }
    return scope;
}
/**
 * Generate declarations for all variables that are in scope for a given view.
 *
 * This is a recursive process, as views inherit variables available from their parent view, which
 * itself may have inherited variables, etc.
 */
function generateVariablesInScopeForView(view, scope, isListener) {
    const newOps = [];
    if (scope.view !== view.xref) {
        // Before generating variables for a parent view, we need to switch to the context of the parent
        // view with a `nextContext` expression. This context switching operation itself declares a
        // variable, because the context of the view may be referenced directly.
        newOps.push(ir.createVariableOp(view.job.allocateXrefId(), scope.viewContextVariable, new ir.NextContextExpr(), ir.VariableFlags.None));
    }
    // Add variables for all context variables available in this scope's view.
    const scopeView = view.job.views.get(scope.view);
    for (const [name, value] of scopeView.contextVariables) {
        const context = new ir.ContextExpr(scope.view);
        // We either read the context, or, if the variable is CTX_REF, use the context directly.
        const variable = value === ir.CTX_REF ? context : new o.ReadPropExpr(context, value);
        // Add the variable declaration.
        newOps.push(ir.createVariableOp(view.job.allocateXrefId(), scope.contextVariables.get(name), variable, ir.VariableFlags.None));
    }
    for (const alias of scopeView.aliases) {
        newOps.push(ir.createVariableOp(view.job.allocateXrefId(), alias, alias.expression.clone(), ir.VariableFlags.AlwaysInline));
    }
    // Add variables for all local references declared for elements in this scope.
    for (const ref of scope.references) {
        newOps.push(ir.createVariableOp(view.job.allocateXrefId(), ref.variable, new ir.ReferenceExpr(ref.targetId, ref.targetSlot, ref.offset), ir.VariableFlags.None));
    }
    if (scope.view !== view.xref || isListener) {
        for (const decl of scope.letDeclarations) {
            newOps.push(ir.createVariableOp(view.job.allocateXrefId(), decl.variable, new ir.ContextLetReferenceExpr(decl.targetId, decl.targetSlot), ir.VariableFlags.None));
        }
    }
    if (scope.parent !== null) {
        // Recursively add variables from the parent scope.
        newOps.push(...generateVariablesInScopeForView(view, scope.parent, false));
    }
    return newOps;
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZ2VuZXJhdGVfdmFyaWFibGVzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vcGFja2FnZXMvY29tcGlsZXIvc3JjL3RlbXBsYXRlL3BpcGVsaW5lL3NyYy9waGFzZXMvZ2VuZXJhdGVfdmFyaWFibGVzLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBOzs7Ozs7R0FNRztBQUVILE9BQU8sS0FBSyxDQUFDLE1BQU0sK0JBQStCLENBQUM7QUFDbkQsT0FBTyxLQUFLLEVBQUUsTUFBTSxVQUFVLENBQUM7QUFJL0I7Ozs7Ozs7Ozs7Ozs7R0FhRztBQUNILE1BQU0sVUFBVSxpQkFBaUIsQ0FBQyxHQUE0QjtJQUM1RCxzQkFBc0IsQ0FBQyxHQUFHLENBQUMsSUFBSSxFQUFFLGdEQUFnRCxDQUFDLElBQUksQ0FBQyxDQUFDO0FBQzFGLENBQUM7QUFFRDs7Ozs7O0dBTUc7QUFDSCxTQUFTLHNCQUFzQixDQUFDLElBQXlCLEVBQUUsV0FBeUI7SUFDbEYsb0NBQW9DO0lBQ3BDLE1BQU0sS0FBSyxHQUFHLGVBQWUsQ0FBQyxJQUFJLEVBQUUsV0FBVyxDQUFDLENBQUM7SUFFakQsS0FBSyxNQUFNLEVBQUUsSUFBSSxJQUFJLENBQUMsTUFBTSxFQUFFLENBQUM7UUFDN0IsUUFBUSxFQUFFLENBQUMsSUFBSSxFQUFFLENBQUM7WUFDaEIsS0FBSyxFQUFFLENBQUMsTUFBTSxDQUFDLFFBQVE7Z0JBQ3JCLHFDQUFxQztnQkFDckMsc0JBQXNCLENBQUMsSUFBSSxDQUFDLEdBQUcsQ0FBQyxLQUFLLENBQUMsR0FBRyxDQUFDLEVBQUUsQ0FBQyxJQUFJLENBQUUsRUFBRSxLQUFLLENBQUMsQ0FBQztnQkFDNUQsTUFBTTtZQUNSLEtBQUssRUFBRSxDQUFDLE1BQU0sQ0FBQyxVQUFVO2dCQUN2QixJQUFJLEVBQUUsQ0FBQyxZQUFZLEtBQUssSUFBSSxFQUFFLENBQUM7b0JBQzdCLHNCQUFzQixDQUFDLElBQUksQ0FB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KlxuICogQGxpY2Vuc2VcbiAqIENvcHlyaWdodCBHb29nbGUgTExDIEFsbCBSaWdodHMgUmVzZXJ2ZWQuXG4gKlxuICogVXNlIG9mIHRoaXMgc291cmNlIGNvZGUgaXMgZ292ZXJuZWQgYnkgYW4gTUlULXN0eWxlIGxpY2Vuc2UgdGhhdCBjYW4gYmVcbiAqIGZvdW5kIGluIHRoZSBMSUNFTlNFIGZpbGUgYXQgaHR0cHM6Ly9hbmd1bGFyLmlvL2xpY2Vuc2VcbiAqL1xuXG5pbXBvcnQgKiBhcyBvIGZyb20gJy4uLy4uLy4uLy4uL291dHB1dC9vdXRwdXRfYXN0JztcbmltcG9ydCAqIGFzIGlyIGZyb20gJy4uLy4uL2lyJztcblxuaW1wb3J0IHR5cGUge0NvbXBvbmVudENvbXBpbGF0aW9uSm9iLCBWaWV3Q29tcGlsYXRpb25Vbml0fSBmcm9tICcuLi9jb21waWxhdGlvbic7XG5cbi8qKlxuICogR2VuZXJhdGUgYSBwcmVhbWJsZSBzZXF1ZW5jZSBmb3IgZWFjaCB2aWV3IGNyZWF0aW9uIGJsb2NrIGFuZCBsaXN0ZW5lciBmdW5jdGlvbiB3aGljaCBkZWNsYXJlc1xuICogYW55IHZhcmlhYmxlcyB0aGF0IGJlIHJlZmVyZW5jZWQgaW4gb3RoZXIgb3BlcmF0aW9ucyBpbiB0aGUgYmxvY2suXG4gKlxuICogVmFyaWFibGVzIGdlbmVyYXRlZCBpbmNsdWRlOlxuICogICAqIGEgc2F2ZWQgdmlldyBjb250ZXh0IHRvIGJlIHVzZWQgdG8gcmVzdG9yZSB0aGUgY3VycmVudCB2aWV3IGluIGV2ZW50IGxpc3RlbmVycy5cbiAqICAgKiB0aGUgY29udGV4dCBvZiB0aGUgcmVzdG9yZWQgdmlldyB3aXRoaW4gZXZlbnQgbGlzdGVuZXIgaGFuZGxlcnMuXG4gKiAgICogY29udGV4dCB2YXJpYWJsZXMgZnJvbSB0aGUgY3VycmVudCB2aWV3IGFzIHdlbGwgYXMgYWxsIHBhcmVudCB2aWV3cyAoaW5jbHVkaW5nIHRoZSByb290XG4gKiAgICAgY29udGV4dCBpZiBuZWVkZWQpLlxuICogICAqIGxvY2FsIHJlZmVyZW5jZXMgZnJvbSBlbGVtZW50cyB3aXRoaW4gdGhlIGN1cnJlbnQgdmlldyBhbmQgYW55IGxleGljYWwgcGFyZW50cy5cbiAqXG4gKiBWYXJpYWJsZXMgYXJlIGdlbmVyYXRlZCBoZXJlIHVuY29uZGl0aW9uYWxseSwgYW5kIG1heSBvcHRpbWl6ZWQgYXdheSBpbiBmdXR1cmUgb3BlcmF0aW9ucyBpZiBpdFxuICogdHVybnMgb3V0IHRoZWlyIHZhbHVlcyAoYW5kIGFueSBzaWRlIGVmZmVjdHMpIGFyZSB1bnVzZWQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBnZW5lcmF0ZVZhcmlhYmxlcyhqb2I6IENvbXBvbmVudENvbXBpbGF0aW9uSm9iKTogdm9pZCB7XG4gIHJlY3Vyc2l2ZWx5UHJvY2Vzc1ZpZXcoam9iLnJvb3QsIC8qIHRoZXJlIGlzIG5vIHBhcmVudCBzY29wZSBmb3IgdGhlIHJvb3QgdmlldyAqLyBudWxsKTtcbn1cblxuLyoqXG4gKiBQcm9jZXNzIHRoZSBnaXZlbiBgVmlld0NvbXBpbGF0aW9uYCBhbmQgZ2VuZXJhdGUgcHJlYW1ibGVzIGZvciBpdCBhbmQgYW55IGxpc3RlbmVycyB0aGF0IGl0XG4gKiBkZWNsYXJlcy5cbiAqXG4gKiBAcGFyYW0gYHBhcmVudFNjb3BlYCBhIHNjb3BlIGV4dHJhY3RlZCBmcm9tIHRoZSBwYXJlbnQgdmlldyB3aGljaCBjYXB0dXJlcyBhbnkgdmFyaWFibGVzIHdoaWNoXG4gKiAgICAgc2hvdWxkIGJlIGluaGVyaXRlZCBieSB0aGlzIHZpZXcuIGBudWxsYCBpZiB0aGUgY3VycmVudCB2aWV3IGlzIHRoZSByb290IHZpZXcuXG4gKi9cbmZ1bmN0aW9uIHJlY3Vyc2l2ZWx5UHJvY2Vzc1ZpZXcodmlldzogVmlld0NvbXBpbGF0aW9uVW5pdCwgcGFyZW50U2NvcGU6IFNjb3BlIHwgbnVsbCk6IHZvaWQge1xuICAvLyBFeHRyYWN0IGEgYFNjb3BlYCBmcm9tIHRoaXMgdmlldy5cbiAgY29uc3Qgc2NvcGUgPSBnZXRTY29wZUZvclZpZXcodmlldywgcGFyZW50U2NvcGUpO1xuXG4gIGZvciAoY29uc3Qgb3Agb2Ygdmlldy5jcmVhdGUpIHtcbiAgICBzd2l0Y2ggKG9wLmtpbmQpIHtcbiAgICAgIGNhc2UgaXIuT3BLaW5kLlRlbXBsYXRlOlxuICAgICAgICAvLyBEZXNjZW5kIGludG8gY2hpbGQgZW1iZWRkZWQgdmlld3MuXG4gICAgICAgIHJlY3Vyc2l2ZWx5UHJvY2Vzc1ZpZXcodmlldy5qb2Iudmlld3MuZ2V0KG9wLnhyZWYpISwgc2NvcGUpO1xuICAgICAgICBicmVhaztcbiAgICAgIGNhc2UgaXIuT3BLaW5kLlByb2plY3Rpb246XG4gICAgICAgIGlmIChvcC5mYWxsYmFja1ZpZXcgIT09IG51bGwpIHtcbiAgICAgICAgICByZWN1cnNpdmVseVByb2Nlc3NWaWV3KHZpZXcuam9iLnZpZXdzLmdldChvcC5mYWxsYmFja1ZpZXcpISwgc2NvcGUpO1xuICAgICAgICB9XG4gICAgICAgIGJyZWFrO1xuICAgICAgY2FzZSBpci5PcEtpbmQuUmVwZWF0ZXJDcmVhdGU6XG4gICAgICAgIC8vIERlc2NlbmQgaW50byBjaGlsZCBlbWJlZGRlZCB2aWV3cy5cbiAgICAgICAgcmVjdXJzaXZlbHlQcm9jZXNzVmlldyh2aWV3LmpvYi52aWV3cy5nZXQob3AueHJlZikhLCBzY29wZSk7XG4gICAgICAgIGlmIChvcC5lbXB0eVZpZXcpIHtcbiAgICAgICAgICByZWN1cnNpdmVseVByb2Nlc3NWaWV3KHZpZXcuam9iLnZpZXdzLmdldChvcC5lbXB0eVZpZXcpISwgc2NvcGUpO1xuICAgICAgICB9XG4gICAgICAgIGJyZWFrO1xuICAgICAgY2FzZSBpci5PcEtpbmQuTGlzdGVuZXI6XG4gICAgICBjYXNlIGlyLk9wS2luZC5Ud29XYXlMaXN0ZW5lcjpcbiAgICAgICAgLy8gUHJlcGVuZCB2YXJpYWJsZXMgdG8gbGlzdGVuZXIgaGFuZGxlciBmdW5jdGlvbnMuXG4gICAgICAgIG9wLmhhbmRsZXJPcHMucHJlcGVuZChnZW5lcmF0ZVZhcmlhYmxlc0luU2NvcGVGb3JWaWV3KHZpZXcsIHNjb3BlLCB0cnVlKSk7XG4gICAgICAgIGJyZWFrO1xuICAgIH1cbiAgfVxuXG4gIHZpZXcudXBkYXRlLnByZXBlbmQoZ2VuZXJhdGVWYXJpYWJsZXNJblNjb3BlRm9yVmlldyh2aWV3LCBzY29wZSwgZmFsc2UpKTtcbn1cblxuLyoqXG4gKiBMZXhpY2FsIHNjb3BlIG9mIGEgdmlldywgaW5jbHVkaW5nIGEgcmVmZXJlbmNlIHRvIGl0cyBwYXJlbnQgdmlldydzIHNjb3BlLCBpZiBhbnkuXG4gKi9cbmludGVyZmFjZSBTY29wZSB7XG4gIC8qKlxuICAgKiBgWHJlZklkYCBvZiB0aGUgdmlldyB0byB3aGljaCB0aGlzIHNjb3BlIGNvcnJlc3BvbmRzLlxuICAgKi9cbiAgdmlldzogaXIuWHJlZklkO1xuXG4gIHZpZXdDb250ZXh0VmFyaWFibGU6IGlyLlNlbWFudGljVmFyaWFibGU7XG5cbiAgY29udGV4dFZhcmlhYmxlczogTWFwPHN0cmluZywgaXIuU2VtYW50aWNWYXJpYWJsZT47XG5cbiAgYWxpYXNlczogU2V0PGlyLkFsaWFzVmFyaWFibGU+O1xuXG4gIC8qKlxuICAgKiBMb2NhbCByZWZlcmVuY2VzIGNvbGxlY3RlZCBmcm9tIGVsZW1lbnRzIHdpdGhpbiB0aGUgdmlldy5cbiAgICovXG4gIHJlZmVyZW5jZXM6IFJlZmVyZW5jZVtdO1xuXG4gIC8qKlxuICAgKiBgQGxldGAgZGVjbGFyYXRpb25zIGNvbGxlY3RlZCBmcm9tIHRoZSB2aWV3LlxuICAgKi9cbiAgbGV0RGVjbGFyYXRpb25zOiBMZXREZWNsYXJhdGlvbltdO1xuXG4gIC8qKlxuICAgKiBgU2NvcGVgIG9mIHRoZSBwYXJlbnQgdmlldywgaWYgYW55LlxuICAgKi9cbiAgcGFyZW50OiBTY29wZSB8IG51bGw7XG59XG5cbi8qKlxuICogSW5mb3JtYXRpb24gbmVlZGVkIGFib3V0IGEgbG9jYWwgcmVmZXJlbmNlIGNvbGxlY3RlZCBmcm9tIGFuIGVsZW1lbnQgd2l0aGluIGEgdmlldy5cbiAqL1xuaW50ZXJmYWNlIFJlZmVyZW5jZSB7XG4gIC8qKlxuICAgKiBOYW1lIGdpdmVuIHRvIHRoZSBsb2NhbCByZWZlcmVuY2UgdmFyaWFibGUgd2l0aGluIHRoZSB0ZW1wbGF0ZS5cbiAgICpcbiAgICogVGhpcyBpcyBub3QgdGhlIG5hbWUgd2hpY2ggd2lsbCBiZSB1c2VkIGZvciB0aGUgdmFyaWFibGUgZGVjbGFyYXRpb24gaW4gdGhlIGdlbmVyYXRlZFxuICAgKiB0ZW1wbGF0ZSBjb2RlLlxuICAgKi9cbiAgbmFtZTogc3RyaW5nO1xuXG4gIC8qKlxuICAgKiBgWHJlZklkYCBvZiB0aGUgZWxlbWVudC1saWtlIG5vZGUgd2hpY2ggdGhpcyByZWZlcmVuY2UgdGFyZ2V0cy5cbiAgICpcbiAgICogVGhlIHJlZmVyZW5jZSBtYXkgYmUgZWl0aGVyIHRvIHRoZSBlbGVtZW50IChvciB0ZW1wbGF0ZSkgaXRzZWxmLCBvciB0byBhIGRpcmVjdGl2ZSBvbiBpdC5cbiAgICovXG4gIHRhcmdldElkOiBpci5YcmVmSWQ7XG5cbiAgdGFyZ2V0U2xvdDogaXIuU2xvdEhhbmRsZTtcblxuICAvKipcbiAgICogQSBnZW5lcmF0ZWQgb2Zmc2V0IG9mIHRoaXMgcmVmZXJlbmNlIGFtb25nIGFsbCB0aGUgcmVmZXJlbmNlcyBvbiBhIHNwZWNpZmljIGVsZW1lbnQuXG4gICAqL1xuICBvZmZzZXQ6IG51bWJlcjtcblxuICB2YXJpYWJsZTogaXIuU2VtYW50aWNWYXJpYWJsZTtcbn1cblxuLyoqXG4gKiBJbmZvcm1hdGlvbiBhYm91dCBgQGxldGAgZGVjbGFyYXRpb24gY29sbGVjdGVkIGZyb20gYSB2aWV3LlxuICovXG5pbnRlcmZhY2UgTGV0RGVjbGFyYXRpb24ge1xuICAvKiogYFhyZWZJZGAgb2YgdGhlIGBAbGV0YCBkZWNsYXJhdGlvbiB0aGF0IHRoZSByZWZlcmVuY2UgaXMgcG9pbnRpbmcgdG8uICovXG4gIHRhcmdldElkOiBpci5YcmVmSWQ7XG5cbiAgLyoqIFNsb3QgaW4gd2hpY2ggdGhlIGRlY2xhcmF0aW9uIGlzIHN0b3JlZC4gKi9cbiAgdGFyZ2V0U2xvdDogaXIuU2xvdEhhbmRsZTtcblxuICAvKiogVmFyaWFibGUgcmVmZXJyaW5nIHRvIHRoZSBkZWNsYXJhdGlvbi4gKi9cbiAgdmFyaWFibGU6IGlyLklkZW50aWZpZXJWYXJpYWJsZTtcbn1cblxuLyoqXG4gKiBQcm9jZXNzIGEgdmlldyBhbmQgZ2VuZXJhdGUgYSBgU2NvcGVgIHJlcHJlc2VudGluZyB0aGUgdmFyaWFibGVzIGF2YWlsYWJsZSBmb3IgcmVmZXJlbmNlIHdpdGhpblxuICogdGhhdCB2aWV3LlxuICovXG5mdW5jdGlvbiBnZXRTY29wZUZvclZpZXcodmlldzogVmlld0NvbXBpbGF0aW9uVW5pdCwgcGFyZW50OiBTY29wZSB8IG51bGwpOiBTY29wZSB7XG4gIGNvbnN0IHNjb3BlOiBTY29wZSA9IHtcbiAgICB2aWV3OiB2aWV3LnhyZWYsXG4gICAgdmlld0NvbnRleHRWYXJpYWJsZToge1xuICAgICAga2luZDogaXIuU2VtYW50aWNWYXJpYWJsZUtpbmQuQ29udGV4dCxcbiAgICAgIG5hbWU6IG51bGwsXG4gICAgICB2aWV3OiB2aWV3LnhyZWYsXG4gICAgfSxcbiAgICBjb250ZXh0VmFyaWFibGVzOiBuZXcgTWFwPHN0cmluZywgaXIuU2VtYW50aWNWYXJpYWJsZT4oKSxcbiAgICBhbGlhc2VzOiB2aWV3LmFsaWFzZXMsXG4gICAgcmVmZXJlbmNlczogW10sXG4gICAgbGV0RGVjbGFyYXRpb25zOiBbXSxcbiAgICBwYXJlbnQsXG4gIH07XG5cbiAgZm9yIChjb25zdCBpZGVudGlmaWVyIG9mIHZpZXcuY29udGV4dFZhcmlhYmxlcy5rZXlzKCkpIHtcbiAgICBzY29wZS5jb250ZXh0VmFyaWFibGVzLnNldChpZGVudGlmaWVyLCB7XG4gICAgICBraW5kOiBpci5TZW1hbnRpY1ZhcmlhYmxlS2luZC5JZGVudGlmaWVyLFxuICAgICAgbmFtZTogbnVsbCxcbiAgICAgIGlkZW50aWZpZXIsXG4gICAgICBsb2NhbDogZmFsc2UsXG4gICAgfSk7XG4gIH1cblxuICBmb3IgKGNvbnN0IG9wIG9mIHZpZXcuY3JlYXRlKSB7XG4gICAgc3dpdGNoIChvcC5raW5kKSB7XG4gICAgICBjYXNlIGlyLk9wS2luZC5FbGVtZW50U3RhcnQ6XG4gICAgICBjYXNlIGlyLk9wS2luZC5UZW1wbGF0ZTpcbiAgICAgICAgaWYgKCFBcnJheS5pc0FycmF5KG9wLmxvY2FsUmVmcykpIHtcbiAgICAgICAgICB0aHJvdyBuZXcgRXJyb3IoYEFzc2VydGlvbkVycm9yOiBleHBlY3RlZCBsb2NhbFJlZnMgdG8gYmUgYW4gYXJyYXlgKTtcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIFJlY29yZCBhdmFpbGFibGUgbG9jYWwgcmVmZXJlbmNlcyBmcm9tIHRoaXMgZWxlbWVudC5cbiAgICAgICAgZm9yIChsZXQgb2Zmc2V0ID0gMDsgb2Zmc2V0IDwgb3AubG9jYWxSZWZzLmxlbmd0aDsgb2Zmc2V0KyspIHtcbiAgICAgICAgICBzY29wZS5yZWZlcmVuY2VzLnB1c2goe1xuICAgICAgICAgICAgbmFtZTogb3AubG9jYWxSZWZzW29mZnNldF0ubmFtZSxcbiAgICAgICAgICAgIHRhcmdldElkOiBvcC54cmVmLFxuICAgICAgICAgICAgdGFyZ2V0U2xvdDogb3AuaGFuZGxlLFxuICAgICAgICAgICAgb2Zmc2V0LFxuICAgICAgICAgICAgdmFyaWFibGU6IHtcbiAgICAgICAgICAgICAga2luZDogaXIuU2VtYW50aWNWYXJpYWJsZUtpbmQuSWRlbnRpZmllcixcbiAgICAgICAgICAgICAgbmFtZTogbnVsbCxcbiAgICAgICAgICAgICAgaWRlbnRpZmllcjogb3AubG9jYWxSZWZzW29mZnNldF0ubmFtZSxcbiAgICAgICAgICAgICAgbG9jYWw6IGZhbHNlLFxuICAgICAgICAgICAgfSxcbiAgICAgICAgICB9KTtcbiAgICAgICAgfVxuICAgICAgICBicmVhaztcblxuICAgICAgY2FzZSBpci5PcEtpbmQuRGVjbGFyZUxldDpcbiAgICAgICAgc2NvcGUubGV0RGVjbGFyYXRpb25zLnB1c2goe1xuICAgICAgICAgIHRhcmdldElkOiBvcC54cmVmLFxuICAgICAgICAgIHRhcmdldFNsb3Q6IG9wLmhhbmRsZSxcbiAgICAgICAgICB2YXJpYWJsZToge1xuICAgICAgICAgICAga2luZDogaXIuU2VtYW50aWNWYXJpYWJsZUtpbmQuSWRlbnRpZmllcixcbiAgICAgICAgICAgIG5hbWU6IG51bGwsXG4gICAgICAgICAgICBpZGVudGlmaWVyOiBvcC5kZWNsYXJlZE5hbWUsXG4gICAgICAgICAgICBsb2NhbDogZmFsc2UsXG4gICAgICAgICAgfSxcbiAgICAgICAgfSk7XG4gICAgICAgIGJyZWFrO1xuICAgIH1cbiAgfVxuXG4gIHJldHVybiBzY29wZTtcbn1cblxuLyoqXG4gKiBHZW5lcmF0ZSBkZWNsYXJhdGlvbnMgZm9yIGFsbCB2YXJpYWJsZXMgdGhhdCBhcmUgaW4gc2NvcGUgZm9yIGEgZ2l2ZW4gdmlldy5cbiAqXG4gKiBUaGlzIGlzIGEgcmVjdXJzaXZlIHByb2Nlc3MsIGFzIHZpZXdzIGluaGVyaXQgdmFyaWFibGVzIGF2YWlsYWJsZSBmcm9tIHRoZWlyIHBhcmVudCB2aWV3LCB3aGljaFxuICogaXRzZWxmIG1heSBoYXZlIGluaGVyaXRlZCB2YXJpYWJsZXMsIGV0Yy5cbiAqL1xuZnVuY3Rpb24gZ2VuZXJhdGVWYXJpYWJsZXNJblNjb3BlRm9yVmlldyhcbiAgdmlldzogVmlld0NvbXBpbGF0aW9uVW5pdCxcbiAgc2NvcGU6IFNjb3BlLFxuICBpc0xpc3RlbmVyOiBib29sZWFuLFxuKTogaXIuVmFyaWFibGVPcDxpci5VcGRhdGVPcD5bXSB7XG4gIGNvbnN0IG5ld09wczogaXIuVmFyaWFibGVPcDxpci5VcGRhdGVPcD5bXSA9IFtdO1xuXG4gIGlmIChzY29wZS52aWV3ICE9PSB2aWV3LnhyZWYpIHtcbiAgICAvLyBCZWZvcmUgZ2VuZXJhdGluZyB2YXJpYWJsZXMgZm9yIGEgcGFyZW50IHZpZXcsIHdlIG5lZWQgdG8gc3dpdGNoIHRvIHRoZSBjb250ZXh0IG9mIHRoZSBwYXJlbnRcbiAgICAvLyB2aWV3IHdpdGggYSBgbmV4dENvbnRleHRgIGV4cHJlc3Npb24uIFRoaXMgY29udGV4dCBzd2l0Y2hpbmcgb3BlcmF0aW9uIGl0c2VsZiBkZWNsYXJlcyBhXG4gICAgLy8gdmFyaWFibGUsIGJlY2F1c2UgdGhlIGNvbnRleHQgb2YgdGhlIHZpZXcgbWF5IGJlIHJlZmVyZW5jZWQgZGlyZWN0bHkuXG4gICAgbmV3T3BzLnB1c2goXG4gICAgICBpci5jcmVhdGVWYXJpYWJsZU9wKFxuICAgICAgICB2aWV3LmpvYi5hbGxvY2F0ZVhyZWZJZCgpLFxuICAgICAgICBzY29wZS52aWV3Q29udGV4dFZhcmlhYmxlLFxuICAgICAgICBuZXcgaXIuTmV4dENvbnRleHRFeHByKCksXG4gICAgICAgIGlyLlZhcmlhYmxlRmxhZ3MuTm9uZSxcbiAgICAgICksXG4gICAgKTtcbiAgfVxuXG4gIC8vIEFkZCB2YXJpYWJsZXMgZm9yIGFsbCBjb250ZXh0IHZhcmlhYmxlcyBhdmFpbGFibGUgaW4gdGhpcyBzY29wZSdzIHZpZXcuXG4gIGNvbnN0IHNjb3BlVmlldyA9IHZpZXcuam9iLnZpZXdzLmdldChzY29wZS52aWV3KSE7XG4gIGZvciAoY29uc3QgW25hbWUsIHZhbHVlXSBvZiBzY29wZVZpZXcuY29udGV4dFZhcmlhYmxlcykge1xuICAgIGNvbnN0IGNvbnRleHQgPSBuZXcgaXIuQ29udGV4dEV4cHIoc2NvcGUudmlldyk7XG4gICAgLy8gV2UgZWl0aGVyIHJlYWQgdGhlIGNvbnRleHQsIG9yLCBpZiB0aGUgdmFyaWFibGUgaXMgQ1RYX1JFRiwgdXNlIHRoZSBjb250ZXh0IGRpcmVjdGx5LlxuICAgIGNvbnN0IHZhcmlhYmxlID0gdmFsdWUgPT09IGlyLkNUWF9SRUYgPyBjb250ZXh0IDogbmV3IG8uUmVhZFByb3BFeHByKGNvbnRleHQsIHZhbHVlKTtcbiAgICAvLyBBZGQgdGhlIHZhcmlhYmxlIGRlY2xhcmF0aW9uLlxuICAgIG5ld09wcy5wdXNoKFxuICAgICAgaXIuY3JlYXRlVmFyaWFibGVPcChcbiAgICAgICAgdmlldy5qb2IuYWxsb2NhdGVYcmVmSWQoKSxcbiAgICAgICAgc2NvcGUuY29udGV4dFZhcmlhYmxlcy5nZXQobmFtZSkhLFxuICAgICAgICB2YXJpYWJsZSxcbiAgICAgICAgaXIuVmFyaWFibGVGbGFncy5Ob25lLFxuICAgICAgKSxcbiAgICApO1xuICB9XG5cbiAgZm9yIChjb25zdCBhbGlhcyBvZiBzY29wZVZpZXcuYWxpYXNlcykge1xuICAgIG5ld09wcy5wdXNoKFxuICAgICAgaXIuY3JlYXRlVmFyaWFibGVPcChcbiAgICAgICAgdmlldy5qb2IuYWxsb2NhdGVYcmVmSWQoKSxcbiAgICAgICAgYWxpYXMsXG4gICAgICAgIGFsaWFzLmV4cHJlc3Npb24uY2xvbmUoKSxcbiAgICAgICAgaXIuVmFyaWFibGVGbGFncy5BbHdheXNJbmxpbmUsXG4gICAgICApLFxuICAgICk7XG4gIH1cblxuICAvLyBBZGQgdmFyaWFibGVzIGZvciBhbGwgbG9jYWwgcmVmZXJlbmNlcyBkZWNsYXJlZCBmb3IgZWxlbWVudHMgaW4gdGhpcyBzY29wZS5cbiAgZm9yIChjb25zdCByZWYgb2Ygc2NvcGUucmVmZXJlbmNlcykge1xuICAgIG5ld09wcy5wdXNoKFxuICAgICAgaXIuY3JlYXRlVmFyaWFibGVPcChcbiAgICAgICAgdmlldy5qb2IuYWxsb2NhdGVYcmVmSWQoKSxcbiAgICAgICAgcmVmLnZhcmlhYmxlLFxuICAgICAgICBuZXcgaXIuUmVmZXJlbmNlRXhwcihyZWYudGFyZ2V0SWQsIHJlZi50YXJnZXRTbG90LCByZWYub2Zmc2V0KSxcbiAgICAgICAgaXIuVmFyaWFibGVGbGFncy5Ob25lLFxuICAgICAgKSxcbiAgICApO1xuICB9XG5cbiAgaWYgKHNjb3BlLnZpZXcgIT09IHZpZXcueHJlZiB8fCBpc0xpc3RlbmVyKSB7XG4gICAgZm9yIChjb25zdCBkZWNsIG9mIHNjb3BlLmxldERlY2xhcmF0aW9ucykge1xuICAgICAgbmV3T3BzLnB1c2goXG4gICAgICAgIGlyLmNyZWF0ZVZhcmlhYmxlT3A8aXIuVXBkYXRlT3A+KFxuICAgICAgICAgIHZpZXcuam9iLmFsbG9jYXRlWHJlZklkKCksXG4gICAgICAgICAgZGVjbC52YXJpYWJsZSxcbiAgICAgICAgICBuZXcgaXIuQ29udGV4dExldFJlZmVyZW5jZUV4cHIoZGVjbC50YXJnZXRJZCwgZGVjbC50YXJnZXRTbG90KSxcbiAgICAgICAgICBpci5WYXJpYWJsZUZsYWdzLk5vbmUsXG4gICAgICAgICksXG4gICAgICApO1xuICAgIH1cbiAgfVxuXG4gIGlmIChzY29wZS5wYXJlbnQgIT09IG51bGwpIHtcbiAgICAvLyBSZWN1cnNpdmVseSBhZGQgdmFyaWFibGVzIGZyb20gdGhlIHBhcmVudCBzY29wZS5cbiAgICBuZXdPcHMucHVzaCguLi5nZW5lcmF0ZVZhcmlhYmxlc0luU2NvcGVGb3JWaWV3KHZpZXcsIHNjb3BlLnBhcmVudCwgZmFsc2UpKTtcbiAgfVxuICByZXR1cm4gbmV3T3BzO1xufVxuIl19