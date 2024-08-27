/**
 * @license
 * Copyright Google LLC All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
import * as ir from '../../ir';
/**
 * Generate `ir.AdvanceOp`s in between `ir.UpdateOp`s that ensure the runtime's implicit slot
 * context will be advanced correctly.
 */
export function generateAdvance(job) {
    for (const unit of job.units) {
        // First build a map of all of the declarations in the view that have assigned slots.
        const slotMap = new Map();
        for (const op of unit.create) {
            if (!ir.hasConsumesSlotTrait(op)) {
                continue;
            }
            else if (op.handle.slot === null) {
                throw new Error(`AssertionError: expected slots to have been allocated before generating advance() calls`);
            }
            slotMap.set(op.xref, op.handle.slot);
        }
        // Next, step through the update operations and generate `ir.AdvanceOp`s as required to ensure
        // the runtime's implicit slot counter will be set to the correct slot before executing each
        // update operation which depends on it.
        //
        // To do that, we track what the runtime's slot counter will be through the update operations.
        let slotContext = 0;
        for (const op of unit.update) {
            let consumer = null;
            if (ir.hasDependsOnSlotContextTrait(op)) {
                consumer = op;
            }
            else {
                ir.visitExpressionsInOp(op, (expr) => {
                    if (consumer === null && ir.hasDependsOnSlotContextTrait(expr)) {
                        consumer = expr;
                    }
                });
            }
            if (consumer === null) {
                continue;
            }
            if (!slotMap.has(consumer.target)) {
                // We expect ops that _do_ depend on the slot counter to point at declarations that exist in
                // the `slotMap`.
                throw new Error(`AssertionError: reference to unknown slot for target ${consumer.target}`);
            }
            const slot = slotMap.get(consumer.target);
            // Does the slot counter need to be adjusted?
            if (slotContext !== slot) {
                // If so, generate an `ir.AdvanceOp` to advance the counter.
                const delta = slot - slotContext;
                if (delta < 0) {
                    throw new Error(`AssertionError: slot counter should never need to move backwards`);
                }
                ir.OpList.insertBefore(ir.createAdvanceOp(delta, consumer.sourceSpan), op);
                slotContext = slot;
            }
        }
    }
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZ2VuZXJhdGVfYWR2YW5jZS5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uL3BhY2thZ2VzL2NvbXBpbGVyL3NyYy90ZW1wbGF0ZS9waXBlbGluZS9zcmMvcGhhc2VzL2dlbmVyYXRlX2FkdmFuY2UudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IkFBQUE7Ozs7OztHQU1HO0FBRUgsT0FBTyxLQUFLLEVBQUUsTUFBTSxVQUFVLENBQUM7QUFHL0I7OztHQUdHO0FBQ0gsTUFBTSxVQUFVLGVBQWUsQ0FBQyxHQUFtQjtJQUNqRCxLQUFLLE1BQU0sSUFBSSxJQUFJLEdBQUcsQ0FBQyxLQUFLLEVBQUUsQ0FBQztRQUM3QixxRkFBcUY7UUFDckYsTUFBTSxPQUFPLEdBQUcsSUFBSSxHQUFHLEVBQXFCLENBQUM7UUFDN0MsS0FBSyxNQUFNLEVBQUUsSUFBSSxJQUFJLENBQUMsTUFBTSxFQUFFLENBQUM7WUFDN0IsSUFBSSxDQUFDLEVBQUUsQ0FBQyxvQkFBb0IsQ0FBQyxFQUFFLENBQUMsRUFBRSxDQUFDO2dCQUNqQyxTQUFTO1lBQ1gsQ0FBQztpQkFBTSxJQUFJLEVBQUUsQ0FBQyxNQUFNLENBQUMsSUFBSSxLQUFLLElBQUksRUFBRSxDQUFDO2dCQUNuQyxNQUFNLElBQUksS0FBSyxDQUNiLHlGQUF5RixDQUMxRixDQUFDO1lBQ0osQ0FBQztZQUVELE9BQU8sQ0FBQyxHQUFHLENBQUMsRUFBRSxDQUFDLElBQUksRUFBRSxFQUFFLENBQUMsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ3ZDLENBQUM7UUFFRCw4RkFBOEY7UUFDOUYsNEZBQTRGO1FBQzVGLHdDQUF3QztRQUN4QyxFQUFFO1FBQ0YsOEZBQThGO1FBQzlGLElBQUksV0FBVyxHQUFHLENBQUMsQ0FBQztRQUNwQixLQUFLLE1BQU0sRUFBRSxJQUFJLElBQUksQ0FBQyxNQUFNLEVBQUUsQ0FBQztZQUM3QixJQUFJLFFBQVEsR0FBMEMsSUFBSSxDQUFDO1lBRTNELElBQUksRUFBRSxDQUFDLDRCQUE0QixDQUFDLEVBQUUsQ0FBQyxFQUFFLENBQUM7Z0JBQ3hDLFFBQVEsR0FBRyxFQUFFLENBQUM7WUFDaEIsQ0FBQztpQkFBTSxDQUFDO2dCQUNOLEVBQUUsQ0FBQyxvQkFBb0IsQ0FBQyxFQUFFLEVBQUUsQ0FBQyxJQUFJLEVBQUUsRUFBRTtvQkFDbkMsSUFBSSxRQUFRLEtBQUssSUFBSSxJQUFJLEVBQUUsQ0FBQyw0QkFBNEIsQ0FBQyxJQUFJLENBQUMsRUFBRSxDQUFDO3dCQUMvRCxRQUFRLEdBQUcsSUFBSSxDQUFDO29CQUNsQixDQUFDO2dCQUNILENBQUMsQ0FBQyxDQUFDO1lBQ0wsQ0FBQztZQUVELElBQUksUUFBUSxLQUFLLElBQUksRUFBRSxDQUFDO2dCQUN0QixTQUFTO1lBQ1gsQ0FBQztZQUVELElBQUksQ0FBQyxPQUFPLENBQUMsR0FBRyxDQUFDLFFBQVEsQ0FBQyxNQUFNLENBQUMsRUFBRSxDQUFDO2dCQUNsQyw0RkFBNEY7Z0JBQzVGLGlCQUFpQjtnQkFDakIsTUFBTSxJQUFJLEtBQUssQ0FBQyx3REFBd0QsUUFBUSxDQUFDLE1BQU0sRUFBRSxDQUFDLENBQUM7WUFDN0YsQ0FBQztZQUVELE1BQU0sSUFBSSxHQUFHLE9BQU8sQ0FBQyxHQUFHLENBQUMsUUFBUSxDQUFDLE1BQU0sQ0FBRSxDQUFDO1lBRTNDLDZDQUE2QztZQUM3QyxJQUFJLFdBQVcsS0FBSyxJQUFJLEVBQUUsQ0FBQztnQkFDekIsNERBQTREO2dCQUM1RCxNQUFNLEtBQUssR0FBRyxJQUFJLEdBQUcsV0FBVyxDQUFDO2dCQUNqQyxJQUFJLEtBQUssR0FBRyxDQUFDLEVBQUUsQ0FBQztvQkFDZCxNQUFNLElBQUksS0FBSyxDQUFDLGtFQUFrRSxDQUFDLENBQUM7Z0JBQ3RGLENBQUM7Z0JBRUQsRUFBRSxDQUFDLE1BQU0sQ0FBQyxZQUFZLENBQWMsRUFBRSxDQUFDLGVBQWUsQ0FBQyxLQUFLLEVBQUUsUUFBUSxDQUFDLFVBQVUsQ0FBQyxFQUFFLEVBQUUsQ0FBQyxDQUFDO2dCQUN4RixXQUFXLEdBQUcsSUFBSSxDQUFDO1lBQ3JCLENBQUM7UUFDSCxDQUFDO0lBQ0gsQ0FBQztBQUNILENBQUMiLCJzb3VyY2VzQ29udGVudCI6WyIvKipcbiAqIEBsaWNlbnNlXG4gKiBDb3B5cmlnaHQgR29vZ2xlIExMQyBBbGwgUmlnaHRzIFJlc2VydmVkLlxuICpcbiAqIFVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlIGlzIGdvdmVybmVkIGJ5IGFuIE1JVC1zdHlsZSBsaWNlbnNlIHRoYXQgY2FuIGJlXG4gKiBmb3VuZCBpbiB0aGUgTElDRU5TRSBmaWxlIGF0IGh0dHBzOi8vYW5ndWxhci5pby9saWNlbnNlXG4gKi9cblxuaW1wb3J0ICogYXMgaXIgZnJvbSAnLi4vLi4vaXInO1xuaW1wb3J0IHR5cGUge0NvbXBpbGF0aW9uSm9ifSBmcm9tICcuLi9jb21waWxhdGlvbic7XG5cbi8qKlxuICogR2VuZXJhdGUgYGlyLkFkdmFuY2VPcGBzIGluIGJldHdlZW4gYGlyLlVwZGF0ZU9wYHMgdGhhdCBlbnN1cmUgdGhlIHJ1bnRpbWUncyBpbXBsaWNpdCBzbG90XG4gKiBjb250ZXh0IHdpbGwgYmUgYWR2YW5jZWQgY29ycmVjdGx5LlxuICovXG5leHBvcnQgZnVuY3Rpb24gZ2VuZXJhdGVBZHZhbmNlKGpvYjogQ29tcGlsYXRpb25Kb2IpOiB2b2lkIHtcbiAgZm9yIChjb25zdCB1bml0IG9mIGpvYi51bml0cykge1xuICAgIC8vIEZpcnN0IGJ1aWxkIGEgbWFwIG9mIGFsbCBvZiB0aGUgZGVjbGFyYXRpb25zIGluIHRoZSB2aWV3IHRoYXQgaGF2ZSBhc3NpZ25lZCBzbG90cy5cbiAgICBjb25zdCBzbG90TWFwID0gbmV3IE1hcDxpci5YcmVmSWQsIG51bWJlcj4oKTtcbiAgICBmb3IgKGNvbnN0IG9wIG9mIHVuaXQuY3JlYXRlKSB7XG4gICAgICBpZiAoIWlyLmhhc0NvbnN1bWVzU2xvdFRyYWl0KG9wKSkge1xuICAgICAgICBjb250aW51ZTtcbiAgICAgIH0gZWxzZSBpZiAob3AuaGFuZGxlLnNsb3QgPT09IG51bGwpIHtcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKFxuICAgICAgICAgIGBBc3NlcnRpb25FcnJvcjogZXhwZWN0ZWQgc2xvdHMgdG8gaGF2ZSBiZWVuIGFsbG9jYXRlZCBiZWZvcmUgZ2VuZXJhdGluZyBhZHZhbmNlKCkgY2FsbHNgLFxuICAgICAgICApO1xuICAgICAgfVxuXG4gICAgICBzbG90TWFwLnNldChvcC54cmVmLCBvcC5oYW5kbGUuc2xvdCk7XG4gICAgfVxuXG4gICAgLy8gTmV4dCwgc3RlcCB0aHJvdWdoIHRoZSB1cGRhdGUgb3BlcmF0aW9ucyBhbmQgZ2VuZXJhdGUgYGlyLkFkdmFuY2VPcGBzIGFzIHJlcXVpcmVkIHRvIGVuc3VyZVxuICAgIC8vIHRoZSBydW50aW1lJ3MgaW1wbGljaXQgc2xvdCBjb3VudGVyIHdpbGwgYmUgc2V0IHRvIHRoZSBjb3JyZWN0IHNsb3QgYmVmb3JlIGV4ZWN1dGluZyBlYWNoXG4gICAgLy8gdXBkYXRlIG9wZXJhdGlvbiB3aGljaCBkZXBlbmRzIG9uIGl0LlxuICAgIC8vXG4gICAgLy8gVG8gZG8gdGhhdCwgd2UgdHJhY2sgd2hhdCB0aGUgcnVudGltZSdzIHNsb3QgY291bnRlciB3aWxsIGJlIHRocm91Z2ggdGhlIHVwZGF0ZSBvcGVyYXRpb25zLlxuICAgIGxldCBzbG90Q29udGV4dCA9IDA7XG4gICAgZm9yIChjb25zdCBvcCBvZiB1bml0LnVwZGF0ZSkge1xuICAgICAgbGV0IGNvbnN1bWVyOiBpci5EZXBlbmRzT25TbG90Q29udGV4dE9wVHJhaXQgfCBudWxsID0gbnVsbDtcblxuICAgICAgaWYgKGlyLmhhc0RlcGVuZHNPblNsb3RDb250ZXh0VHJhaXQob3ApKSB7XG4gICAgICAgIGNvbnN1bWVyID0gb3A7XG4gICAgICB9IGVsc2Uge1xuICAgICAgICBpci52aXNpdEV4cHJlc3Npb25zSW5PcChvcCwgKGV4cHIpID0+IHtcbiAgICAgICAgICBpZiAoY29uc3VtZXIgPT09IG51bGwgJiYgaXIuaGFzRGVwZW5kc09uU2xvdENvbnRleHRUcmFpdChleHByKSkge1xuICAgICAgICAgICAgY29uc3VtZXIgPSBleHByO1xuICAgICAgICAgIH1cbiAgICAgICAgfSk7XG4gICAgICB9XG5cbiAgICAgIGlmIChjb25zdW1lciA9PT0gbnVsbCkge1xuICAgICAgICBjb250aW51ZTtcbiAgICAgIH1cblxuICAgICAgaWYgKCFzbG90TWFwLmhhcyhjb25zdW1lci50YXJnZXQpKSB7XG4gICAgICAgIC8vIFdlIGV4cGVjdCBvcHMgdGhhdCBfZG9fIGRlcGVuZCBvbiB0aGUgc2xvdCBjb3VudGVyIHRvIHBvaW50IGF0IGRlY2xhcmF0aW9ucyB0aGF0IGV4aXN0IGluXG4gICAgICAgIC8vIHRoZSBgc2xvdE1hcGAuXG4gICAgICAgIHRocm93IG5ldyBFcnJvcihgQXNzZXJ0aW9uRXJyb3I6IHJlZmVyZW5jZSB0byB1bmtub3duIHNsb3QgZm9yIHRhcmdldCAke2NvbnN1bWVyLnRhcmdldH1gKTtcbiAgICAgIH1cblxuICAgICAgY29uc3Qgc2xvdCA9IHNsb3RNYXAuZ2V0KGNvbnN1bWVyLnRhcmdldCkhO1xuXG4gICAgICAvLyBEb2VzIHRoZSBzbG90IGNvdW50ZXIgbmVlZCB0byBiZSBhZGp1c3RlZD9cbiAgICAgIGlmIChzbG90Q29udGV4dCAhPT0gc2xvdCkge1xuICAgICAgICAvLyBJZiBzbywgZ2VuZXJhdGUgYW4gYGlyLkFkdmFuY2VPcGAgdG8gYWR2YW5jZSB0aGUgY291bnRlci5cbiAgICAgICAgY29uc3QgZGVsdGEgPSBzbG90IC0gc2xvdENvbnRleHQ7XG4gICAgICAgIGlmIChkZWx0YSA8IDApIHtcbiAgICAgICAgICB0aHJvdyBuZXcgRXJyb3IoYEFzc2VydGlvbkVycm9yOiBzbG90IGNvdW50ZXIgc2hvdWxkIG5ldmVyIG5lZWQgdG8gbW92ZSBiYWNrd2FyZHNgKTtcbiAgICAgICAgfVxuXG4gICAgICAgIGlyLk9wTGlzdC5pbnNlcnRCZWZvcmU8aXIuVXBkYXRlT3A+KGlyLmNyZWF0ZUFkdmFuY2VPcChkZWx0YSwgY29uc3VtZXIuc291cmNlU3BhbiksIG9wKTtcbiAgICAgICAgc2xvdENvbnRleHQgPSBzbG90O1xuICAgICAgfVxuICAgIH1cbiAgfVxufVxuIl19