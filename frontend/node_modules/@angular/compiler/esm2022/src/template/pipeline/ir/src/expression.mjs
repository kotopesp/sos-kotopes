/**
 * @license
 * Copyright Google LLC All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
var _a, _b, _c, _d, _e, _f, _g, _h;
import * as o from '../../../../output/output_ast';
import { ExpressionKind, OpKind } from './enums';
import { Interpolation } from './ops/update';
import { ConsumesVarsTrait, DependsOnSlotContext, UsesVarOffset, } from './traits';
/**
 * Check whether a given `o.Expression` is a logical IR expression type.
 */
export function isIrExpression(expr) {
    return expr instanceof ExpressionBase;
}
/**
 * Base type used for all logical IR expressions.
 */
export class ExpressionBase extends o.Expression {
    constructor(sourceSpan = null) {
        super(null, sourceSpan);
    }
}
/**
 * Logical expression representing a lexical read of a variable name.
 */
export class LexicalReadExpr extends ExpressionBase {
    constructor(name) {
        super();
        this.name = name;
        this.kind = ExpressionKind.LexicalRead;
    }
    visitExpression(visitor, context) { }
    isEquivalent(other) {
        // We assume that the lexical reads are in the same context, which must be true for parent
        // expressions to be equivalent.
        // TODO: is this generally safe?
        return this.name === other.name;
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions() { }
    clone() {
        return new LexicalReadExpr(this.name);
    }
}
/**
 * Runtime operation to retrieve the value of a local reference.
 */
export class ReferenceExpr extends ExpressionBase {
    constructor(target, targetSlot, offset) {
        super();
        this.target = target;
        this.targetSlot = targetSlot;
        this.offset = offset;
        this.kind = ExpressionKind.Reference;
    }
    visitExpression() { }
    isEquivalent(e) {
        return e instanceof ReferenceExpr && e.target === this.target;
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions() { }
    clone() {
        return new ReferenceExpr(this.target, this.targetSlot, this.offset);
    }
}
export class StoreLetExpr extends ExpressionBase {
    static { _a = ConsumesVarsTrait, _b = DependsOnSlotContext; }
    constructor(target, value, sourceSpan) {
        super();
        this.target = target;
        this.value = value;
        this.sourceSpan = sourceSpan;
        this.kind = ExpressionKind.StoreLet;
        this[_a] = true;
        this[_b] = true;
    }
    visitExpression() { }
    isEquivalent(e) {
        return (e instanceof StoreLetExpr && e.target === this.target && e.value.isEquivalent(this.value));
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions(transform, flags) {
        this.value = transformExpressionsInExpression(this.value, transform, flags);
    }
    clone() {
        return new StoreLetExpr(this.target, this.value, this.sourceSpan);
    }
}
export class ContextLetReferenceExpr extends ExpressionBase {
    constructor(target, targetSlot) {
        super();
        this.target = target;
        this.targetSlot = targetSlot;
        this.kind = ExpressionKind.ContextLetReference;
    }
    visitExpression() { }
    isEquivalent(e) {
        return e instanceof ContextLetReferenceExpr && e.target === this.target;
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions() { }
    clone() {
        return new ContextLetReferenceExpr(this.target, this.targetSlot);
    }
}
/**
 * A reference to the current view context (usually the `ctx` variable in a template function).
 */
export class ContextExpr extends ExpressionBase {
    constructor(view) {
        super();
        this.view = view;
        this.kind = ExpressionKind.Context;
    }
    visitExpression() { }
    isEquivalent(e) {
        return e instanceof ContextExpr && e.view === this.view;
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions() { }
    clone() {
        return new ContextExpr(this.view);
    }
}
/**
 * A reference to the current view context inside a track function.
 */
export class TrackContextExpr extends ExpressionBase {
    constructor(view) {
        super();
        this.view = view;
        this.kind = ExpressionKind.TrackContext;
    }
    visitExpression() { }
    isEquivalent(e) {
        return e instanceof TrackContextExpr && e.view === this.view;
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions() { }
    clone() {
        return new TrackContextExpr(this.view);
    }
}
/**
 * Runtime operation to navigate to the next view context in the view hierarchy.
 */
export class NextContextExpr extends ExpressionBase {
    constructor() {
        super();
        this.kind = ExpressionKind.NextContext;
        this.steps = 1;
    }
    visitExpression() { }
    isEquivalent(e) {
        return e instanceof NextContextExpr && e.steps === this.steps;
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions() { }
    clone() {
        const expr = new NextContextExpr();
        expr.steps = this.steps;
        return expr;
    }
}
/**
 * Runtime operation to snapshot the current view context.
 *
 * The result of this operation can be stored in a variable and later used with the `RestoreView`
 * operation.
 */
export class GetCurrentViewExpr extends ExpressionBase {
    constructor() {
        super();
        this.kind = ExpressionKind.GetCurrentView;
    }
    visitExpression() { }
    isEquivalent(e) {
        return e instanceof GetCurrentViewExpr;
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions() { }
    clone() {
        return new GetCurrentViewExpr();
    }
}
/**
 * Runtime operation to restore a snapshotted view.
 */
export class RestoreViewExpr extends ExpressionBase {
    constructor(view) {
        super();
        this.view = view;
        this.kind = ExpressionKind.RestoreView;
    }
    visitExpression(visitor, context) {
        if (typeof this.view !== 'number') {
            this.view.visitExpression(visitor, context);
        }
    }
    isEquivalent(e) {
        if (!(e instanceof RestoreViewExpr) || typeof e.view !== typeof this.view) {
            return false;
        }
        if (typeof this.view === 'number') {
            return this.view === e.view;
        }
        else {
            return this.view.isEquivalent(e.view);
        }
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions(transform, flags) {
        if (typeof this.view !== 'number') {
            this.view = transformExpressionsInExpression(this.view, transform, flags);
        }
    }
    clone() {
        return new RestoreViewExpr(this.view instanceof o.Expression ? this.view.clone() : this.view);
    }
}
/**
 * Runtime operation to reset the current view context after `RestoreView`.
 */
export class ResetViewExpr extends ExpressionBase {
    constructor(expr) {
        super();
        this.expr = expr;
        this.kind = ExpressionKind.ResetView;
    }
    visitExpression(visitor, context) {
        this.expr.visitExpression(visitor, context);
    }
    isEquivalent(e) {
        return e instanceof ResetViewExpr && this.expr.isEquivalent(e.expr);
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions(transform, flags) {
        this.expr = transformExpressionsInExpression(this.expr, transform, flags);
    }
    clone() {
        return new ResetViewExpr(this.expr.clone());
    }
}
export class TwoWayBindingSetExpr extends ExpressionBase {
    constructor(target, value) {
        super();
        this.target = target;
        this.value = value;
        this.kind = ExpressionKind.TwoWayBindingSet;
    }
    visitExpression(visitor, context) {
        this.target.visitExpression(visitor, context);
        this.value.visitExpression(visitor, context);
    }
    isEquivalent(other) {
        return this.target.isEquivalent(other.target) && this.value.isEquivalent(other.value);
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions(transform, flags) {
        this.target = transformExpressionsInExpression(this.target, transform, flags);
        this.value = transformExpressionsInExpression(this.value, transform, flags);
    }
    clone() {
        return new TwoWayBindingSetExpr(this.target, this.value);
    }
}
/**
 * Read of a variable declared as an `ir.VariableOp` and referenced through its `ir.XrefId`.
 */
export class ReadVariableExpr extends ExpressionBase {
    constructor(xref) {
        super();
        this.xref = xref;
        this.kind = ExpressionKind.ReadVariable;
        this.name = null;
    }
    visitExpression() { }
    isEquivalent(other) {
        return other instanceof ReadVariableExpr && other.xref === this.xref;
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions() { }
    clone() {
        const expr = new ReadVariableExpr(this.xref);
        expr.name = this.name;
        return expr;
    }
}
export class PureFunctionExpr extends ExpressionBase {
    static { _c = ConsumesVarsTrait, _d = UsesVarOffset; }
    constructor(expression, args) {
        super();
        this.kind = ExpressionKind.PureFunctionExpr;
        this[_c] = true;
        this[_d] = true;
        this.varOffset = null;
        /**
         * Once extracted to the `ConstantPool`, a reference to the function which defines the computation
         * of `body`.
         */
        this.fn = null;
        this.body = expression;
        this.args = args;
    }
    visitExpression(visitor, context) {
        this.body?.visitExpression(visitor, context);
        for (const arg of this.args) {
            arg.visitExpression(visitor, context);
        }
    }
    isEquivalent(other) {
        if (!(other instanceof PureFunctionExpr) || other.args.length !== this.args.length) {
            return false;
        }
        return (other.body !== null &&
            this.body !== null &&
            other.body.isEquivalent(this.body) &&
            other.args.every((arg, idx) => arg.isEquivalent(this.args[idx])));
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions(transform, flags) {
        if (this.body !== null) {
            // TODO: figure out if this is the right flag to pass here.
            this.body = transformExpressionsInExpression(this.body, transform, flags | VisitorContextFlag.InChildOperation);
        }
        else if (this.fn !== null) {
            this.fn = transformExpressionsInExpression(this.fn, transform, flags);
        }
        for (let i = 0; i < this.args.length; i++) {
            this.args[i] = transformExpressionsInExpression(this.args[i], transform, flags);
        }
    }
    clone() {
        const expr = new PureFunctionExpr(this.body?.clone() ?? null, this.args.map((arg) => arg.clone()));
        expr.fn = this.fn?.clone() ?? null;
        expr.varOffset = this.varOffset;
        return expr;
    }
}
export class PureFunctionParameterExpr extends ExpressionBase {
    constructor(index) {
        super();
        this.index = index;
        this.kind = ExpressionKind.PureFunctionParameterExpr;
    }
    visitExpression() { }
    isEquivalent(other) {
        return other instanceof PureFunctionParameterExpr && other.index === this.index;
    }
    isConstant() {
        return true;
    }
    transformInternalExpressions() { }
    clone() {
        return new PureFunctionParameterExpr(this.index);
    }
}
export class PipeBindingExpr extends ExpressionBase {
    static { _e = ConsumesVarsTrait, _f = UsesVarOffset; }
    constructor(target, targetSlot, name, args) {
        super();
        this.target = target;
        this.targetSlot = targetSlot;
        this.name = name;
        this.args = args;
        this.kind = ExpressionKind.PipeBinding;
        this[_e] = true;
        this[_f] = true;
        this.varOffset = null;
    }
    visitExpression(visitor, context) {
        for (const arg of this.args) {
            arg.visitExpression(visitor, context);
        }
    }
    isEquivalent() {
        return false;
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions(transform, flags) {
        for (let idx = 0; idx < this.args.length; idx++) {
            this.args[idx] = transformExpressionsInExpression(this.args[idx], transform, flags);
        }
    }
    clone() {
        const r = new PipeBindingExpr(this.target, this.targetSlot, this.name, this.args.map((a) => a.clone()));
        r.varOffset = this.varOffset;
        return r;
    }
}
export class PipeBindingVariadicExpr extends ExpressionBase {
    static { _g = ConsumesVarsTrait, _h = UsesVarOffset; }
    constructor(target, targetSlot, name, args, numArgs) {
        super();
        this.target = target;
        this.targetSlot = targetSlot;
        this.name = name;
        this.args = args;
        this.numArgs = numArgs;
        this.kind = ExpressionKind.PipeBindingVariadic;
        this[_g] = true;
        this[_h] = true;
        this.varOffset = null;
    }
    visitExpression(visitor, context) {
        this.args.visitExpression(visitor, context);
    }
    isEquivalent() {
        return false;
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions(transform, flags) {
        this.args = transformExpressionsInExpression(this.args, transform, flags);
    }
    clone() {
        const r = new PipeBindingVariadicExpr(this.target, this.targetSlot, this.name, this.args.clone(), this.numArgs);
        r.varOffset = this.varOffset;
        return r;
    }
}
export class SafePropertyReadExpr extends ExpressionBase {
    constructor(receiver, name) {
        super();
        this.receiver = receiver;
        this.name = name;
        this.kind = ExpressionKind.SafePropertyRead;
    }
    // An alias for name, which allows other logic to handle property reads and keyed reads together.
    get index() {
        return this.name;
    }
    visitExpression(visitor, context) {
        this.receiver.visitExpression(visitor, context);
    }
    isEquivalent() {
        return false;
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions(transform, flags) {
        this.receiver = transformExpressionsInExpression(this.receiver, transform, flags);
    }
    clone() {
        return new SafePropertyReadExpr(this.receiver.clone(), this.name);
    }
}
export class SafeKeyedReadExpr extends ExpressionBase {
    constructor(receiver, index, sourceSpan) {
        super(sourceSpan);
        this.receiver = receiver;
        this.index = index;
        this.kind = ExpressionKind.SafeKeyedRead;
    }
    visitExpression(visitor, context) {
        this.receiver.visitExpression(visitor, context);
        this.index.visitExpression(visitor, context);
    }
    isEquivalent() {
        return false;
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions(transform, flags) {
        this.receiver = transformExpressionsInExpression(this.receiver, transform, flags);
        this.index = transformExpressionsInExpression(this.index, transform, flags);
    }
    clone() {
        return new SafeKeyedReadExpr(this.receiver.clone(), this.index.clone(), this.sourceSpan);
    }
}
export class SafeInvokeFunctionExpr extends ExpressionBase {
    constructor(receiver, args) {
        super();
        this.receiver = receiver;
        this.args = args;
        this.kind = ExpressionKind.SafeInvokeFunction;
    }
    visitExpression(visitor, context) {
        this.receiver.visitExpression(visitor, context);
        for (const a of this.args) {
            a.visitExpression(visitor, context);
        }
    }
    isEquivalent() {
        return false;
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions(transform, flags) {
        this.receiver = transformExpressionsInExpression(this.receiver, transform, flags);
        for (let i = 0; i < this.args.length; i++) {
            this.args[i] = transformExpressionsInExpression(this.args[i], transform, flags);
        }
    }
    clone() {
        return new SafeInvokeFunctionExpr(this.receiver.clone(), this.args.map((a) => a.clone()));
    }
}
export class SafeTernaryExpr extends ExpressionBase {
    constructor(guard, expr) {
        super();
        this.guard = guard;
        this.expr = expr;
        this.kind = ExpressionKind.SafeTernaryExpr;
    }
    visitExpression(visitor, context) {
        this.guard.visitExpression(visitor, context);
        this.expr.visitExpression(visitor, context);
    }
    isEquivalent() {
        return false;
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions(transform, flags) {
        this.guard = transformExpressionsInExpression(this.guard, transform, flags);
        this.expr = transformExpressionsInExpression(this.expr, transform, flags);
    }
    clone() {
        return new SafeTernaryExpr(this.guard.clone(), this.expr.clone());
    }
}
export class EmptyExpr extends ExpressionBase {
    constructor() {
        super(...arguments);
        this.kind = ExpressionKind.EmptyExpr;
    }
    visitExpression(visitor, context) { }
    isEquivalent(e) {
        return e instanceof EmptyExpr;
    }
    isConstant() {
        return true;
    }
    clone() {
        return new EmptyExpr();
    }
    transformInternalExpressions() { }
}
export class AssignTemporaryExpr extends ExpressionBase {
    constructor(expr, xref) {
        super();
        this.expr = expr;
        this.xref = xref;
        this.kind = ExpressionKind.AssignTemporaryExpr;
        this.name = null;
    }
    visitExpression(visitor, context) {
        this.expr.visitExpression(visitor, context);
    }
    isEquivalent() {
        return false;
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions(transform, flags) {
        this.expr = transformExpressionsInExpression(this.expr, transform, flags);
    }
    clone() {
        const a = new AssignTemporaryExpr(this.expr.clone(), this.xref);
        a.name = this.name;
        return a;
    }
}
export class ReadTemporaryExpr extends ExpressionBase {
    constructor(xref) {
        super();
        this.xref = xref;
        this.kind = ExpressionKind.ReadTemporaryExpr;
        this.name = null;
    }
    visitExpression(visitor, context) { }
    isEquivalent() {
        return this.xref === this.xref;
    }
    isConstant() {
        return false;
    }
    transformInternalExpressions(transform, flags) { }
    clone() {
        const r = new ReadTemporaryExpr(this.xref);
        r.name = this.name;
        return r;
    }
}
export class SlotLiteralExpr extends ExpressionBase {
    constructor(slot) {
        super();
        this.slot = slot;
        this.kind = ExpressionKind.SlotLiteralExpr;
    }
    visitExpression(visitor, context) { }
    isEquivalent(e) {
        return e instanceof SlotLiteralExpr && e.slot === this.slot;
    }
    isConstant() {
        return true;
    }
    clone() {
        return new SlotLiteralExpr(this.slot);
    }
    transformInternalExpressions() { }
}
export class ConditionalCaseExpr extends ExpressionBase {
    /**
     * Create an expression for one branch of a conditional.
     * @param expr The expression to be tested for this case. Might be null, as in an `else` case.
     * @param target The Xref of the view to be displayed if this condition is true.
     */
    constructor(expr, target, targetSlot, alias = null) {
        super();
        this.expr = expr;
        this.target = target;
        this.targetSlot = targetSlot;
        this.alias = alias;
        this.kind = ExpressionKind.ConditionalCase;
    }
    visitExpression(visitor, context) {
        if (this.expr !== null) {
            this.expr.visitExpression(visitor, context);
        }
    }
    isEquivalent(e) {
        return e instanceof ConditionalCaseExpr && e.expr === this.expr;
    }
    isConstant() {
        return true;
    }
    clone() {
        return new ConditionalCaseExpr(this.expr, this.target, this.targetSlot);
    }
    transformInternalExpressions(transform, flags) {
        if (this.expr !== null) {
            this.expr = transformExpressionsInExpression(this.expr, transform, flags);
        }
    }
}
export class ConstCollectedExpr extends ExpressionBase {
    constructor(expr) {
        super();
        this.expr = expr;
        this.kind = ExpressionKind.ConstCollected;
    }
    transformInternalExpressions(transform, flags) {
        this.expr = transform(this.expr, flags);
    }
    visitExpression(visitor, context) {
        this.expr.visitExpression(visitor, context);
    }
    isEquivalent(e) {
        if (!(e instanceof ConstCollectedExpr)) {
            return false;
        }
        return this.expr.isEquivalent(e.expr);
    }
    isConstant() {
        return this.expr.isConstant();
    }
    clone() {
        return new ConstCollectedExpr(this.expr);
    }
}
/**
 * Visits all `Expression`s in the AST of `op` with the `visitor` function.
 */
export function visitExpressionsInOp(op, visitor) {
    transformExpressionsInOp(op, (expr, flags) => {
        visitor(expr, flags);
        return expr;
    }, VisitorContextFlag.None);
}
export var VisitorContextFlag;
(function (VisitorContextFlag) {
    VisitorContextFlag[VisitorContextFlag["None"] = 0] = "None";
    VisitorContextFlag[VisitorContextFlag["InChildOperation"] = 1] = "InChildOperation";
})(VisitorContextFlag || (VisitorContextFlag = {}));
function transformExpressionsInInterpolation(interpolation, transform, flags) {
    for (let i = 0; i < interpolation.expressions.length; i++) {
        interpolation.expressions[i] = transformExpressionsInExpression(interpolation.expressions[i], transform, flags);
    }
}
/**
 * Transform all `Expression`s in the AST of `op` with the `transform` function.
 *
 * All such operations will be replaced with the result of applying `transform`, which may be an
 * identity transformation.
 */
export function transformExpressionsInOp(op, transform, flags) {
    switch (op.kind) {
        case OpKind.StyleProp:
        case OpKind.StyleMap:
        case OpKind.ClassProp:
        case OpKind.ClassMap:
        case OpKind.Binding:
            if (op.expression instanceof Interpolation) {
                transformExpressionsInInterpolation(op.expression, transform, flags);
            }
            else {
                op.expression = transformExpressionsInExpression(op.expression, transform, flags);
            }
            break;
        case OpKind.Property:
        case OpKind.HostProperty:
        case OpKind.Attribute:
            if (op.expression instanceof Interpolation) {
                transformExpressionsInInterpolation(op.expression, transform, flags);
            }
            else {
                op.expression = transformExpressionsInExpression(op.expression, transform, flags);
            }
            op.sanitizer =
                op.sanitizer && transformExpressionsInExpression(op.sanitizer, transform, flags);
            break;
        case OpKind.TwoWayProperty:
            op.expression = transformExpressionsInExpression(op.expression, transform, flags);
            op.sanitizer =
                op.sanitizer && transformExpressionsInExpression(op.sanitizer, transform, flags);
            break;
        case OpKind.I18nExpression:
            op.expression = transformExpressionsInExpression(op.expression, transform, flags);
            break;
        case OpKind.InterpolateText:
            transformExpressionsInInterpolation(op.interpolation, transform, flags);
            break;
        case OpKind.Statement:
            transformExpressionsInStatement(op.statement, transform, flags);
            break;
        case OpKind.Variable:
            op.initializer = transformExpressionsInExpression(op.initializer, transform, flags);
            break;
        case OpKind.Conditional:
            for (const condition of op.conditions) {
                if (condition.expr === null) {
                    // This is a default case.
                    continue;
                }
                condition.expr = transformExpressionsInExpression(condition.expr, transform, flags);
            }
            if (op.processed !== null) {
                op.processed = transformExpressionsInExpression(op.processed, transform, flags);
            }
            if (op.contextValue !== null) {
                op.contextValue = transformExpressionsInExpression(op.contextValue, transform, flags);
            }
            break;
        case OpKind.Listener:
        case OpKind.TwoWayListener:
            for (const innerOp of op.handlerOps) {
                transformExpressionsInOp(innerOp, transform, flags | VisitorContextFlag.InChildOperation);
            }
            break;
        case OpKind.ExtractedAttribute:
            op.expression =
                op.expression && transformExpressionsInExpression(op.expression, transform, flags);
            op.trustedValueFn =
                op.trustedValueFn && transformExpressionsInExpression(op.trustedValueFn, transform, flags);
            break;
        case OpKind.RepeaterCreate:
            op.track = transformExpressionsInExpression(op.track, transform, flags);
            if (op.trackByFn !== null) {
                op.trackByFn = transformExpressionsInExpression(op.trackByFn, transform, flags);
            }
            break;
        case OpKind.Repeater:
            op.collection = transformExpressionsInExpression(op.collection, transform, flags);
            break;
        case OpKind.Defer:
            if (op.loadingConfig !== null) {
                op.loadingConfig = transformExpressionsInExpression(op.loadingConfig, transform, flags);
            }
            if (op.placeholderConfig !== null) {
                op.placeholderConfig = transformExpressionsInExpression(op.placeholderConfig, transform, flags);
            }
            if (op.resolverFn !== null) {
                op.resolverFn = transformExpressionsInExpression(op.resolverFn, transform, flags);
            }
            break;
        case OpKind.I18nMessage:
            for (const [placeholder, expr] of op.params) {
                op.params.set(placeholder, transformExpressionsInExpression(expr, transform, flags));
            }
            for (const [placeholder, expr] of op.postprocessingParams) {
                op.postprocessingParams.set(placeholder, transformExpressionsInExpression(expr, transform, flags));
            }
            break;
        case OpKind.DeferWhen:
            op.expr = transformExpressionsInExpression(op.expr, transform, flags);
            break;
        case OpKind.StoreLet:
            op.value = transformExpressionsInExpression(op.value, transform, flags);
            break;
        case OpKind.Advance:
        case OpKind.Container:
        case OpKind.ContainerEnd:
        case OpKind.ContainerStart:
        case OpKind.DeferOn:
        case OpKind.DisableBindings:
        case OpKind.Element:
        case OpKind.ElementEnd:
        case OpKind.ElementStart:
        case OpKind.EnableBindings:
        case OpKind.I18n:
        case OpKind.I18nApply:
        case OpKind.I18nContext:
        case OpKind.I18nEnd:
        case OpKind.I18nStart:
        case OpKind.IcuEnd:
        case OpKind.IcuStart:
        case OpKind.Namespace:
        case OpKind.Pipe:
        case OpKind.Projection:
        case OpKind.ProjectionDef:
        case OpKind.Template:
        case OpKind.Text:
        case OpKind.I18nAttributes:
        case OpKind.IcuPlaceholder:
        case OpKind.DeclareLet:
            // These operations contain no expressions.
            break;
        default:
            throw new Error(`AssertionError: transformExpressionsInOp doesn't handle ${OpKind[op.kind]}`);
    }
}
/**
 * Transform all `Expression`s in the AST of `expr` with the `transform` function.
 *
 * All such operations will be replaced with the result of applying `transform`, which may be an
 * identity transformation.
 */
export function transformExpressionsInExpression(expr, transform, flags) {
    if (expr instanceof ExpressionBase) {
        expr.transformInternalExpressions(transform, flags);
    }
    else if (expr instanceof o.BinaryOperatorExpr) {
        expr.lhs = transformExpressionsInExpression(expr.lhs, transform, flags);
        expr.rhs = transformExpressionsInExpression(expr.rhs, transform, flags);
    }
    else if (expr instanceof o.UnaryOperatorExpr) {
        expr.expr = transformExpressionsInExpression(expr.expr, transform, flags);
    }
    else if (expr instanceof o.ReadPropExpr) {
        expr.receiver = transformExpressionsInExpression(expr.receiver, transform, flags);
    }
    else if (expr instanceof o.ReadKeyExpr) {
        expr.receiver = transformExpressionsInExpression(expr.receiver, transform, flags);
        expr.index = transformExpressionsInExpression(expr.index, transform, flags);
    }
    else if (expr instanceof o.WritePropExpr) {
        expr.receiver = transformExpressionsInExpression(expr.receiver, transform, flags);
        expr.value = transformExpressionsInExpression(expr.value, transform, flags);
    }
    else if (expr instanceof o.WriteKeyExpr) {
        expr.receiver = transformExpressionsInExpression(expr.receiver, transform, flags);
        expr.index = transformExpressionsInExpression(expr.index, transform, flags);
        expr.value = transformExpressionsInExpression(expr.value, transform, flags);
    }
    else if (expr instanceof o.InvokeFunctionExpr) {
        expr.fn = transformExpressionsInExpression(expr.fn, transform, flags);
        for (let i = 0; i < expr.args.length; i++) {
            expr.args[i] = transformExpressionsInExpression(expr.args[i], transform, flags);
        }
    }
    else if (expr instanceof o.LiteralArrayExpr) {
        for (let i = 0; i < expr.entries.length; i++) {
            expr.entries[i] = transformExpressionsInExpression(expr.entries[i], transform, flags);
        }
    }
    else if (expr instanceof o.LiteralMapExpr) {
        for (let i = 0; i < expr.entries.length; i++) {
            expr.entries[i].value = transformExpressionsInExpression(expr.entries[i].value, transform, flags);
        }
    }
    else if (expr instanceof o.ConditionalExpr) {
        expr.condition = transformExpressionsInExpression(expr.condition, transform, flags);
        expr.trueCase = transformExpressionsInExpression(expr.trueCase, transform, flags);
        if (expr.falseCase !== null) {
            expr.falseCase = transformExpressionsInExpression(expr.falseCase, transform, flags);
        }
    }
    else if (expr instanceof o.TypeofExpr) {
        expr.expr = transformExpressionsInExpression(expr.expr, transform, flags);
    }
    else if (expr instanceof o.WriteVarExpr) {
        expr.value = transformExpressionsInExpression(expr.value, transform, flags);
    }
    else if (expr instanceof o.LocalizedString) {
        for (let i = 0; i < expr.expressions.length; i++) {
            expr.expressions[i] = transformExpressionsInExpression(expr.expressions[i], transform, flags);
        }
    }
    else if (expr instanceof o.NotExpr) {
        expr.condition = transformExpressionsInExpression(expr.condition, transform, flags);
    }
    else if (expr instanceof o.TaggedTemplateExpr) {
        expr.tag = transformExpressionsInExpression(expr.tag, transform, flags);
        expr.template.expressions = expr.template.expressions.map((e) => transformExpressionsInExpression(e, transform, flags));
    }
    else if (expr instanceof o.ArrowFunctionExpr) {
        if (Array.isArray(expr.body)) {
            for (let i = 0; i < expr.body.length; i++) {
                transformExpressionsInStatement(expr.body[i], transform, flags);
            }
        }
        else {
            expr.body = transformExpressionsInExpression(expr.body, transform, flags);
        }
    }
    else if (expr instanceof o.WrappedNodeExpr) {
        // TODO: Do we need to transform any TS nodes nested inside of this expression?
    }
    else if (expr instanceof o.ReadVarExpr ||
        expr instanceof o.ExternalExpr ||
        expr instanceof o.LiteralExpr) {
        // No action for these types.
    }
    else {
        throw new Error(`Unhandled expression kind: ${expr.constructor.name}`);
    }
    return transform(expr, flags);
}
/**
 * Transform all `Expression`s in the AST of `stmt` with the `transform` function.
 *
 * All such operations will be replaced with the result of applying `transform`, which may be an
 * identity transformation.
 */
export function transformExpressionsInStatement(stmt, transform, flags) {
    if (stmt instanceof o.ExpressionStatement) {
        stmt.expr = transformExpressionsInExpression(stmt.expr, transform, flags);
    }
    else if (stmt instanceof o.ReturnStatement) {
        stmt.value = transformExpressionsInExpression(stmt.value, transform, flags);
    }
    else if (stmt instanceof o.DeclareVarStmt) {
        if (stmt.value !== undefined) {
            stmt.value = transformExpressionsInExpression(stmt.value, transform, flags);
        }
    }
    else if (stmt instanceof o.IfStmt) {
        stmt.condition = transformExpressionsInExpression(stmt.condition, transform, flags);
        for (const caseStatement of stmt.trueCase) {
            transformExpressionsInStatement(caseStatement, transform, flags);
        }
        for (const caseStatement of stmt.falseCase) {
            transformExpressionsInStatement(caseStatement, transform, flags);
        }
    }
    else {
        throw new Error(`Unhandled statement kind: ${stmt.constructor.name}`);
    }
}
/**
 * Checks whether the given expression is a string literal.
 */
export function isStringLiteral(expr) {
    return expr instanceof o.LiteralExpr && typeof expr.value === 'string';
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZXhwcmVzc2lvbi5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uL3BhY2thZ2VzL2NvbXBpbGVyL3NyYy90ZW1wbGF0ZS9waXBlbGluZS9pci9zcmMvZXhwcmVzc2lvbi50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiQUFBQTs7Ozs7O0dBTUc7O0FBRUgsT0FBTyxLQUFLLENBQUMsTUFBTSwrQkFBK0IsQ0FBQztBQUluRCxPQUFPLEVBQUMsY0FBYyxFQUFFLE1BQU0sRUFBQyxNQUFNLFNBQVMsQ0FBQztBQUkvQyxPQUFPLEVBQUMsYUFBYSxFQUFnQixNQUFNLGNBQWMsQ0FBQztBQUMxRCxPQUFPLEVBQ0wsaUJBQWlCLEVBQ2pCLG9CQUFvQixFQUVwQixhQUFhLEdBRWQsTUFBTSxVQUFVLENBQUM7QUFxQ2xCOztHQUVHO0FBQ0gsTUFBTSxVQUFVLGNBQWMsQ0FBQyxJQUFrQjtJQUMvQyxPQUFPLElBQUksWUFBWSxjQUFjLENBQUM7QUFDeEMsQ0FBQztBQUVEOztHQUVHO0FBQ0gsTUFBTSxPQUFnQixjQUFlLFNBQVEsQ0FBQyxDQUFDLFVBQVU7SUFHdkQsWUFBWSxhQUFxQyxJQUFJO1FBQ25ELEtBQUssQ0FBQyxJQUFJLEVBQUUsVUFBVSxDQUFDLENBQUM7SUFDMUIsQ0FBQztDQVVGO0FBRUQ7O0dBRUc7QUFDSCxNQUFNLE9BQU8sZUFBZ0IsU0FBUSxjQUFjO0lBR2pELFlBQXFCLElBQVk7UUFDL0IsS0FBSyxFQUFFLENBQUM7UUFEVyxTQUFJLEdBQUosSUFBSSxDQUFRO1FBRmYsU0FBSSxHQUFHLGNBQWMsQ0FBQyxXQUFXLENBQUM7SUFJcEQsQ0FBQztJQUVRLGVBQWUsQ0FBQyxPQUE0QixFQUFFLE9BQVksSUFBUyxDQUFDO0lBRXBFLFlBQVksQ0FBQyxLQUFzQjtRQUMxQywwRkFBMEY7UUFDMUYsZ0NBQWdDO1FBQ2hDLGdDQUFnQztRQUNoQyxPQUFPLElBQUksQ0FBQyxJQUFJLEtBQUssS0FBSyxD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QUFDLENBQUM7UUFDbEYsSUFBSSxDQUFDLEtBQUssR0FBRyxnQ0FBZ0MsQ0FBQyxJQUFJLENBQUMsS0FBSyxFQUFFLFNBQVMsRUFBRSxLQUFLLENBQUMsQ0FBQztJQUM5RSxDQUFDO1NBQU0sSUFBSSxJQUFJLFlBQVksQ0FBQyxDQUFDLGFBQWEsRUFBRSxDQUFDO1FBQzNDLElBQUksQ0FBQyxRQUFRLEdBQUcsZ0NBQWdDLENBQUMsSUFBSSxDQUFDLFFBQVEsRUFBRSxTQUFTLEVBQUUsS0FBSyxDQUFDLENBQUM7UUFDbEYsSUFBSSxDQUFDLEtBQUssR0FBRyxnQ0FBZ0MsQ0FBQyxJQUFJLENBQUMsS0FBSyxFQUFFLFNBQVMsRUFBRSxLQUFLLENBQUMsQ0FBQztJQUM5RSxDQUFDO1NBQU0sSUFBSSxJQUFJLFlBQVksQ0FBQyxDQUFDLFlBQVksRUFBRSxDQUFDO1FBQzFDLElBQUksQ0FBQyxRQUFRLEdBQUcsZ0NBQWdDLENBQUMsSUFBSSxDQUFDLFFBQVEsRUFBRSxTQUFTLEVBQUUsS0FBSyxDQUFDLENBQUM7UUFDbEYsSUFBSSxDQUFDLEtBQUssR0FBRyxnQ0FBZ0MsQ0FBQyxJQUFJLENBQUMsS0FBSyxFQUFFLFNBQVMsRUFBRSxLQUFLLENBQUMsQ0FBQztRQUM1RSxJQUFJLENBQUMsS0FBSyxHQUFHLGdDQUFnQyxDQUFDLElBQUksQ0FBQyxLQUFLLEVBQUUsU0FBUyxFQUFFLEtBQUssQ0FBQyxDQUFDO0lBQzlFLENBQUM7U0FBTSxJQUFJLElBQUksWUFBWSxDQUFDLENBQUMsa0JBQWtCLEVBQUUsQ0FBQztRQUNoRCxJQUFJLENBQUMsRUFBRSxHQUFHLGdDQUFnQyxDQUFDLElBQUksQ0FBQyxFQUFFLEVBQUUsU0FBUyxFQUFFLEtBQUssQ0FBQyxDQUFDO1FBQ3RFLEtBQUssSUFBSSxDQUFDLEdBQUcsQ0FBQyxFQUFFLENBQUMsR0FBRyxJQUFJLENBQUMsSUFBSSxDQUFDLE1BQU0sRUFBRSxDQUFDLEVBQUUsRUFBRSxDQUFDO1lBQzFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDLEdBQUcsZ0NB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QUMsU0FBUyxFQUFFLENBQUM7WUFDM0MsK0JBQStCLENBQUMsYUFBYSxFQUFFLFNBQVMsRUFBRSxLQUFLLENBQUMsQ0FBQztRQUNuRSxDQUFDO0lBQ0gsQ0FBQztTQUFNLENBQUM7UUFDTixNQUFNLElBQUksS0FBSyxDQUFDLDZCQUE2QixJQUFJLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxDQUFDLENBQUM7SUFDeEUsQ0FBQztBQUNILENBQUM7QUFFRDs7R0FFRztBQUNILE1BQU0sVUFBVSxlQUFlLENBQUMsSUFBa0I7SUFDaEQsT0FBTyxJQUFJLFlBQVksQ0FBQyxDQUFDLFdBQVcsSUFBSSxPQUFPLElBQUksQ0FBQyxLQUFLLEtBQUssUUFBUSxDQUFDO0FBQ3pFLENBQUMiLCJzb3VyY2VzQ29udGVudCI6WyIvKipcbiAqIEBsaWNlbnNlXG4gKiBDb3B5cmlnaHQgR29vZ2xlIExMQyBBbGwgUmlnaHRzIFJlc2VydmVkLlxuICpcbiAqIFVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlIGlzIGdvdmVybmVkIGJ5IGFuIE1JVC1zdHlsZSBsaWNlbnNlIHRoYXQgY2FuIGJlXG4gKiBmb3VuZCBpbiB0aGUgTElDRU5TRSBmaWxlIGF0IGh0dHBzOi8vYW5ndWxhci5pby9saWNlbnNlXG4gKi9cblxuaW1wb3J0ICogYXMgbyBmcm9tICcuLi8uLi8uLi8uLi9vdXRwdXQvb3V0cHV0X2FzdCc7XG5pbXBvcnQgdHlwZSB7UGFyc2VTb3VyY2VTcGFufSBmcm9tICcuLi8uLi8uLi8uLi9wYXJzZV91dGlsJztcblxuaW1wb3J0ICogYXMgdCBmcm9tICcuLi8uLi8uLi8uLi9yZW5kZXIzL3IzX2FzdCc7XG5pbXBvcnQge0V4cHJlc3Npb25LaW5kLCBPcEtpbmR9IGZyb20gJy4vZW51bXMnO1xuaW1wb3J0IHtTbG90SGFuZGxlfSBmcm9tICcuL2hhbmRsZSc7XG5pbXBvcnQgdHlwZSB7WHJlZklkfSBmcm9tICcuL29wZXJhdGlvbnMnO1xuaW1wb3J0IHR5cGUge0NyZWF0ZU9wfSBmcm9tICcuL29wcy9jcmVhdGUnO1xuaW1wb3J0IHtJbnRlcnBvbGF0aW9uLCB0eXBlIFVwZGF0ZU9wfSBmcm9tICcuL29wcy91cGRhdGUnO1xuaW1wb3J0IHtcbiAgQ29uc3VtZXNWYXJzVHJhaXQsXG4gIERlcGVuZHNPblNsb3RDb250ZXh0LFxuICBEZXBlbmRzT25TbG90Q29udGV4dE9wVHJhaXQsXG4gIFVzZXNWYXJPZmZzZXQsXG4gIFVzZXNWYXJPZmZzZXRUcmFpdCxcbn0gZnJvbSAnLi90cmFpdHMnO1xuXG4vKipcbiAqIEFuIGBvLkV4cHJlc3Npb25gIHN1YnR5cGUgcmVwcmVzZW50aW5nIGEgbG9naWNhbCBleHByZXNzaW9uIGluIHRoZSBpbnRlcm1lZGlhdGUgcmVwcmVzZW50YXRpb24uXG4gKi9cbmV4cG9ydCB0eXBlIEV4cHJlc3Npb24gPVxuICB8IExleGljYWxSZWFkRXhwclxuICB8IFJlZmVyZW5jZUV4cHJcbiAgfCBDb250ZXh0RXhwclxuICB8IE5leHRDb250ZXh0RXhwclxuICB8IEdldEN1cnJlbnRWaWV3RXhwclxuICB8IFJlc3RvcmVWaWV3RXhwclxuICB8IFJlc2V0Vmlld0V4cHJcbiAgfCBSZWFkVmFyaWFibGVFeHByXG4gIHwgUHVyZUZ1bmN0aW9uRXhwclxuICB8IFB1cmVGdW5jdGlvblBhcmFtZXRlckV4cHJcbiAgfCBQaXBlQmluZGluZ0V4cHJcbiAgfCBQaXBlQmluZGluZ1ZhcmlhZGljRXhwclxuICB8IFNhZmVQcm9wZXJ0eVJlYWRFeHByXG4gIHwgU2FmZUtleWVkUmVhZEV4cHJcbiAgfCBTYWZlSW52b2tlRnVuY3Rpb25FeHByXG4gIHwgRW1wdHlFeHByXG4gIHwgQXNzaWduVGVtcG9yYXJ5RXhwclxuICB8IFJlYWRUZW1wb3JhcnlFeHByXG4gIHwgU2xvdExpdGVyYWxFeHByXG4gIHwgQ29uZGl0aW9uYWxDYXNlRXhwclxuICB8IENvbnN0Q29sbGVjdGVkRXhwclxuICB8IFR3b1dheUJpbmRpbmdTZXRFeHByXG4gIHwgQ29udGV4dExldFJlZmVyZW5jZUV4cHJcbiAgfCBTdG9yZUxldEV4cHI7XG5cbi8qKlxuICogVHJhbnNmb3JtZXIgdHlwZSB3aGljaCBjb252ZXJ0cyBleHByZXNzaW9ucyBpbnRvIGdlbmVyYWwgYG8uRXhwcmVzc2lvbmBzICh3aGljaCBtYXkgYmUgYW5cbiAqIGlkZW50aXR5IHRyYW5zZm9ybWF0aW9uKS5cbiAqL1xuZXhwb3J0IHR5cGUgRXhwcmVzc2lvblRyYW5zZm9ybSA9IChleHByOiBvLkV4cHJlc3Npb24sIGZsYWdzOiBWaXNpdG9yQ29udGV4dEZsYWcpID0+IG8uRXhwcmVzc2lvbjtcblxuLyoqXG4gKiBDaGVjayB3aGV0aGVyIGEgZ2l2ZW4gYG8uRXhwcmVzc2lvbmAgaXMgYSBsb2dpY2FsIElSIGV4cHJlc3Npb24gdHlwZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGlzSXJFeHByZXNzaW9uKGV4cHI6IG8uRXhwcmVzc2lvbik6IGV4cHIgaXMgRXhwcmVzc2lvbiB7XG4gIHJldHVybiBleHByIGluc3RhbmNlb2YgRXhwcmVzc2lvbkJhc2U7XG59XG5cbi8qKlxuICogQmFzZSB0eXBlIHVzZWQgZm9yIGFsbCBsb2dpY2FsIElSIGV4cHJlc3Npb25zLlxuICovXG5leHBvcnQgYWJzdHJhY3QgY2xhc3MgRXhwcmVzc2lvbkJhc2UgZXh0ZW5kcyBvLkV4cHJlc3Npb24ge1xuICBhYnN0cmFjdCByZWFkb25seSBraW5kOiBFeHByZXNzaW9uS2luZDtcblxuICBjb25zdHJ1Y3Rvcihzb3VyY2VTcGFuOiBQYXJzZVNvdXJjZVNwYW4gfCBudWxsID0gbnVsbCkge1xuICAgIHN1cGVyKG51bGwsIHNvdXJjZVNwYW4pO1xuICB9XG5cbiAgLyoqXG4gICAqIFJ1biB0aGUgdHJhbnNmb3JtZXIgYWdhaW5zdCBhbnkgbmVzdGVkIGV4cHJlc3Npb25zIHdoaWNoIG1heSBiZSBwcmVzZW50IGluIHRoaXMgSVIgZXhwcmVzc2lvblxuICAgKiBzdWJ0eXBlLlxuICAgKi9cbiAgYWJzdHJhY3QgdHJhbnNmb3JtSW50ZXJuYWxFeHByZXNzaW9ucyhcbiAgICB0cmFuc2Zvcm06IEV4cHJlc3Npb25UcmFuc2Zvcm0sXG4gICAgZmxhZ3M6IFZpc2l0b3JDb250ZXh0RmxhZyxcbiAgKTogdm9pZDtcbn1cblxuLyoqXG4gKiBMb2dpY2FsIGV4cHJlc3Npb24gcmVwcmVzZW50aW5nIGEgbGV4aWNhbCByZWFkIG9mIGEgdmFyaWFibGUgbmFtZS5cbiAqL1xuZXhwb3J0IGNsYXNzIExleGljYWxSZWFkRXhwciBleHRlbmRzIEV4cHJlc3Npb25CYXNlIHtcbiAgb3ZlcnJpZGUgcmVhZG9ubHkga2luZCA9IEV4cHJlc3Npb25LaW5kLkxleGljYWxSZWFkO1xuXG4gIGNvbnN0cnVjdG9yKHJlYWRvbmx5IG5hbWU6IHN0cmluZykge1xuICAgIHN1cGVyKCk7XG4gIH1cblxuICBvdmVycmlkZSB2aXNpdEV4cHJlc3Npb24odmlzaXRvcjogby5FeHByZXNzaW9uVmlzaXRvciwgY29udGV4dDogYW55KTogdm9pZCB7fVxuXG4gIG92ZXJyaWRlIGlzRXF1aXZhbGVudChvdGhlcjogTGV4aWNhbFJlYWRFeHByKTogYm9vbGVhbiB7XG4gICAgLy8gV2UgYXNzdW1lIHRoYXQgdGhlIGxleGljYWwgcmVhZHMgYXJlIGluIHRoZSBzYW1lIGNvbnRleHQsIHdoaWNoIG11c3QgYmUgdHJ1ZSBmb3IgcGFyZW50XG4gICAgLy8gZXhwcmVzc2lvbnMgdG8gYmUgZXF1aXZhbGVudC5cbiAgICAvLyBUT0RPOiBpcyB0aGlzIGdlbmVyYWxseSBzYWZlP1xuICAgIHJldHVybiB0aGlzLm5hbWUgPT09IG90aGVyLm5hbWU7XG4gIH1cblxuICBvdmVycmlkZSBpc0NvbnN0YW50KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiBmYWxzZTtcbiAgfVxuXG4gIG92ZXJyaWRlIHRyYW5zZm9ybUludGVybmFsRXhwcmVzc2lvbnMoKTogdm9pZCB7fVxuXG4gIG92ZXJyaWRlIGNsb25lKCk6IExleGljYWxSZWFkRXhwciB7XG4gICAgcmV0dXJuIG5ldyBMZXhpY2FsUmVhZEV4cHIodGhpcy5uYW1lKTtcbiAgfVxufVxuXG4vKipcbiAqIFJ1bnRpbWUgb3BlcmF0aW9uIHRvIHJldHJpZXZlIHRoZSB2YWx1ZSBvZiBhIGxvY2FsIHJlZmVyZW5jZS5cbiAqL1xuZXhwb3J0IGNsYXNzIFJlZmVyZW5jZUV4cHIgZXh0ZW5kcyBFeHByZXNzaW9uQmFzZSB7XG4gIG92ZXJyaWRlIHJlYWRvbmx5IGtpbmQgPSBFeHByZXNzaW9uS2luZC5SZWZlcmVuY2U7XG5cbiAgY29uc3RydWN0b3IoXG4gICAgcmVhZG9ubHkgdGFyZ2V0OiBYcmVmSWQsXG4gICAgcmVhZG9ubHkgdGFyZ2V0U2xvdDogU2xvdEhhbmRsZSxcbiAgICByZWFkb25seSBvZmZzZXQ6IG51bWJlcixcbiAgKSB7XG4gICAgc3VwZXIoKTtcbiAgfVxuXG4gIG92ZXJyaWRlIHZpc2l0RXhwcmVzc2lvbigpOiB2b2lkIHt9XG5cbiAgb3ZlcnJpZGUgaXNFcXVpdmFsZW50KGU6IG8uRXhwcmVzc2lvbik6IGJvb2xlYW4ge1xuICAgIHJldHVybiBlIGluc3RhbmNlb2YgUmVmZXJlbmNlRXhwciAmJiBlLnRhcmdldCA9PT0gdGhpcy50YXJnZXQ7XG4gIH1cblxuICBvdmVycmlkZSBpc0NvbnN0YW50KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiBmYWxzZTtcbiAgfVxuXG4gIG92ZXJyaWRlIHRyYW5zZm9ybUludGVybmFsRXhwcmVzc2lvbnMoKTogdm9pZCB7fVxuXG4gIG92ZXJyaWRlIGNsb25lKCk6IFJlZmVyZW5jZUV4cHIge1xuICAgIHJldHVybiBuZXcgUmVmZXJlbmNlRXhwcih0aGlzLnRhcmdldCwgdGhpcy50YXJnZXRTbG90LCB0aGlzLm9mZnNldCk7XG4gIH1cbn1cblxuZXhwb3J0IGNsYXNzIFN0b3JlTGV0RXhwclxuICBleHRlbmRzIEV4cHJlc3Npb25CYXNlXG4gIGltcGxlbWVudHMgQ29uc3VtZXNWYXJzVHJhaXQsIERlcGVuZHNPblNsb3RDb250ZXh0T3BUcmFpdFxue1xuICBvdmVycmlkZSByZWFkb25seSBraW5kID0gRXhwcmVzc2lvbktpbmQuU3RvcmVMZXQ7XG4gIHJlYWRvbmx5IFtDb25zdW1lc1ZhcnNUcmFpdF0gPSB0cnVlO1xuICByZWFkb25seSBbRGVwZW5kc09uU2xvdENvbnRleHRdID0gdHJ1ZTtcblxuICBjb25zdHJ1Y3RvcihcbiAgICByZWFkb25seSB0YXJnZXQ6IFhyZWZJZCxcbiAgICBwdWJsaWMgdmFsdWU6IG8uRXhwcmVzc2lvbixcbiAgICBvdmVycmlkZSBzb3VyY2VTcGFuOiBQYXJzZVNvdXJjZVNwYW4sXG4gICkge1xuICAgIHN1cGVyKCk7XG4gIH1cblxuICBvdmVycmlkZSB2aXNpdEV4cHJlc3Npb24oKTogdm9pZCB7fVxuXG4gIG92ZXJyaWRlIGlzRXF1aXZhbGVudChlOiBvLkV4cHJlc3Npb24pOiBib29sZWFuIHtcbiAgICByZXR1cm4gKFxuICAgICAgZSBpbnN0YW5jZW9mIFN0b3JlTGV0RXhwciAmJiBlLnRhcmdldCA9PT0gdGhpcy50YXJnZXQgJiYgZS52YWx1ZS5pc0VxdWl2YWxlbnQodGhpcy52YWx1ZSlcbiAgICApO1xuICB9XG5cbiAgb3ZlcnJpZGUgaXNDb25zdGFudCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gZmFsc2U7XG4gIH1cblxuICBvdmVycmlkZSB0cmFuc2Zvcm1JbnRlcm5hbEV4cHJlc3Npb25zKFxuICAgIHRyYW5zZm9ybTogRXhwcmVzc2lvblRyYW5zZm9ybSxcbiAgICBmbGFnczogVmlzaXRvckNvbnRleHRGbGFnLFxuICApOiB2b2lkIHtcbiAgICB0aGlzLnZhbHVlID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24odGhpcy52YWx1ZSwgdHJhbnNmb3JtLCBmbGFncyk7XG4gIH1cblxuICBvdmVycmlkZSBjbG9uZSgpOiBTdG9yZUxldEV4cHIge1xuICAgIHJldHVybiBuZXcgU3RvcmVMZXRFeHByKHRoaXMudGFyZ2V0LCB0aGlzLnZhbHVlLCB0aGlzLnNvdXJjZVNwYW4pO1xuICB9XG59XG5cbmV4cG9ydCBjbGFzcyBDb250ZXh0TGV0UmVmZXJlbmNlRXhwciBleHRlbmRzIEV4cHJlc3Npb25CYXNlIHtcbiAgb3ZlcnJpZGUgcmVhZG9ubHkga2luZCA9IEV4cHJlc3Npb25LaW5kLkNvbnRleHRMZXRSZWZlcmVuY2U7XG5cbiAgY29uc3RydWN0b3IoXG4gICAgcmVhZG9ubHkgdGFyZ2V0OiBYcmVmSWQsXG4gICAgcmVhZG9ubHkgdGFyZ2V0U2xvdDogU2xvdEhhbmRsZSxcbiAgKSB7XG4gICAgc3VwZXIoKTtcbiAgfVxuXG4gIG92ZXJyaWRlIHZpc2l0RXhwcmVzc2lvbigpOiB2b2lkIHt9XG5cbiAgb3ZlcnJpZGUgaXNFcXVpdmFsZW50KGU6IG8uRXhwcmVzc2lvbik6IGJvb2xlYW4ge1xuICAgIHJldHVybiBlIGluc3RhbmNlb2YgQ29udGV4dExldFJlZmVyZW5jZUV4cHIgJiYgZS50YXJnZXQgPT09IHRoaXMudGFyZ2V0O1xuICB9XG5cbiAgb3ZlcnJpZGUgaXNDb25zdGFudCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gZmFsc2U7XG4gIH1cblxuICBvdmVycmlkZSB0cmFuc2Zvcm1JbnRlcm5hbEV4cHJlc3Npb25zKCk6IHZvaWQge31cblxuICBvdmVycmlkZSBjbG9uZSgpOiBDb250ZXh0TGV0UmVmZXJlbmNlRXhwciB7XG4gICAgcmV0dXJuIG5ldyBDb250ZXh0TGV0UmVmZXJlbmNlRXhwcih0aGlzLnRhcmdldCwgdGhpcy50YXJnZXRTbG90KTtcbiAgfVxufVxuXG4vKipcbiAqIEEgcmVmZXJlbmNlIHRvIHRoZSBjdXJyZW50IHZpZXcgY29udGV4dCAodXN1YWxseSB0aGUgYGN0eGAgdmFyaWFibGUgaW4gYSB0ZW1wbGF0ZSBmdW5jdGlvbikuXG4gKi9cbmV4cG9ydCBjbGFzcyBDb250ZXh0RXhwciBleHRlbmRzIEV4cHJlc3Npb25CYXNlIHtcbiAgb3ZlcnJpZGUgcmVhZG9ubHkga2luZCA9IEV4cHJlc3Npb25LaW5kLkNvbnRleHQ7XG5cbiAgY29uc3RydWN0b3IocmVhZG9ubHkgdmlldzogWHJlZklkKSB7XG4gICAgc3VwZXIoKTtcbiAgfVxuXG4gIG92ZXJyaWRlIHZpc2l0RXhwcmVzc2lvbigpOiB2b2lkIHt9XG5cbiAgb3ZlcnJpZGUgaXNFcXVpdmFsZW50KGU6IG8uRXhwcmVzc2lvbik6IGJvb2xlYW4ge1xuICAgIHJldHVybiBlIGluc3RhbmNlb2YgQ29udGV4dEV4cHIgJiYgZS52aWV3ID09PSB0aGlzLnZpZXc7XG4gIH1cblxuICBvdmVycmlkZSBpc0NvbnN0YW50KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiBmYWxzZTtcbiAgfVxuXG4gIG92ZXJyaWRlIHRyYW5zZm9ybUludGVybmFsRXhwcmVzc2lvbnMoKTogdm9pZCB7fVxuXG4gIG92ZXJyaWRlIGNsb25lKCk6IENvbnRleHRFeHByIHtcbiAgICByZXR1cm4gbmV3IENvbnRleHRFeHByKHRoaXMudmlldyk7XG4gIH1cbn1cblxuLyoqXG4gKiBBIHJlZmVyZW5jZSB0byB0aGUgY3VycmVudCB2aWV3IGNvbnRleHQgaW5zaWRlIGEgdHJhY2sgZnVuY3Rpb24uXG4gKi9cbmV4cG9ydCBjbGFzcyBUcmFja0NvbnRleHRFeHByIGV4dGVuZHMgRXhwcmVzc2lvbkJhc2Uge1xuICBvdmVycmlkZSByZWFkb25seSBraW5kID0gRXhwcmVzc2lvbktpbmQuVHJhY2tDb250ZXh0O1xuXG4gIGNvbnN0cnVjdG9yKHJlYWRvbmx5IHZpZXc6IFhyZWZJZCkge1xuICAgIHN1cGVyKCk7XG4gIH1cblxuICBvdmVycmlkZSB2aXNpdEV4cHJlc3Npb24oKTogdm9pZCB7fVxuXG4gIG92ZXJyaWRlIGlzRXF1aXZhbGVudChlOiBvLkV4cHJlc3Npb24pOiBib29sZWFuIHtcbiAgICByZXR1cm4gZSBpbnN0YW5jZW9mIFRyYWNrQ29udGV4dEV4cHIgJiYgZS52aWV3ID09PSB0aGlzLnZpZXc7XG4gIH1cblxuICBvdmVycmlkZSBpc0NvbnN0YW50KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiBmYWxzZTtcbiAgfVxuXG4gIG92ZXJyaWRlIHRyYW5zZm9ybUludGVybmFsRXhwcmVzc2lvbnMoKTogdm9pZCB7fVxuXG4gIG92ZXJyaWRlIGNsb25lKCk6IFRyYWNrQ29udGV4dEV4cHIge1xuICAgIHJldHVybiBuZXcgVHJhY2tDb250ZXh0RXhwcih0aGlzLnZpZXcpO1xuICB9XG59XG5cbi8qKlxuICogUnVudGltZSBvcGVyYXRpb24gdG8gbmF2aWdhdGUgdG8gdGhlIG5leHQgdmlldyBjb250ZXh0IGluIHRoZSB2aWV3IGhpZXJhcmNoeS5cbiAqL1xuZXhwb3J0IGNsYXNzIE5leHRDb250ZXh0RXhwciBleHRlbmRzIEV4cHJlc3Npb25CYXNlIHtcbiAgb3ZlcnJpZGUgcmVhZG9ubHkga2luZCA9IEV4cHJlc3Npb25LaW5kLk5leHRDb250ZXh0O1xuXG4gIHN0ZXBzID0gMTtcblxuICBjb25zdHJ1Y3RvcigpIHtcbiAgICBzdXBlcigpO1xuICB9XG5cbiAgb3ZlcnJpZGUgdmlzaXRFeHByZXNzaW9uKCk6IHZvaWQge31cblxuICBvdmVycmlkZSBpc0VxdWl2YWxlbnQoZTogby5FeHByZXNzaW9uKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIGUgaW5zdGFuY2VvZiBOZXh0Q29udGV4dEV4cHIgJiYgZS5zdGVwcyA9PT0gdGhpcy5zdGVwcztcbiAgfVxuXG4gIG92ZXJyaWRlIGlzQ29uc3RhbnQoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIGZhbHNlO1xuICB9XG5cbiAgb3ZlcnJpZGUgdHJhbnNmb3JtSW50ZXJuYWxFeHByZXNzaW9ucygpOiB2b2lkIHt9XG5cbiAgb3ZlcnJpZGUgY2xvbmUoKTogTmV4dENvbnRleHRFeHByIHtcbiAgICBjb25zdCBleHByID0gbmV3IE5leHRDb250ZXh0RXhwcigpO1xuICAgIGV4cHIuc3RlcHMgPSB0aGlzLnN0ZXBzO1xuICAgIHJldHVybiBleHByO1xuICB9XG59XG5cbi8qKlxuICogUnVudGltZSBvcGVyYXRpb24gdG8gc25hcHNob3QgdGhlIGN1cnJlbnQgdmlldyBjb250ZXh0LlxuICpcbiAqIFRoZSByZXN1bHQgb2YgdGhpcyBvcGVyYXRpb24gY2FuIGJlIHN0b3JlZCBpbiBhIHZhcmlhYmxlIGFuZCBsYXRlciB1c2VkIHdpdGggdGhlIGBSZXN0b3JlVmlld2BcbiAqIG9wZXJhdGlvbi5cbiAqL1xuZXhwb3J0IGNsYXNzIEdldEN1cnJlbnRWaWV3RXhwciBleHRlbmRzIEV4cHJlc3Npb25CYXNlIHtcbiAgb3ZlcnJpZGUgcmVhZG9ubHkga2luZCA9IEV4cHJlc3Npb25LaW5kLkdldEN1cnJlbnRWaWV3O1xuXG4gIGNvbnN0cnVjdG9yKCkge1xuICAgIHN1cGVyKCk7XG4gIH1cblxuICBvdmVycmlkZSB2aXNpdEV4cHJlc3Npb24oKTogdm9pZCB7fVxuXG4gIG92ZXJyaWRlIGlzRXF1aXZhbGVudChlOiBvLkV4cHJlc3Npb24pOiBib29sZWFuIHtcbiAgICByZXR1cm4gZSBpbnN0YW5jZW9mIEdldEN1cnJlbnRWaWV3RXhwcjtcbiAgfVxuXG4gIG92ZXJyaWRlIGlzQ29uc3RhbnQoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIGZhbHNlO1xuICB9XG5cbiAgb3ZlcnJpZGUgdHJhbnNmb3JtSW50ZXJuYWxFeHByZXNzaW9ucygpOiB2b2lkIHt9XG5cbiAgb3ZlcnJpZGUgY2xvbmUoKTogR2V0Q3VycmVudFZpZXdFeHByIHtcbiAgICByZXR1cm4gbmV3IEdldEN1cnJlbnRWaWV3RXhwcigpO1xuICB9XG59XG5cbi8qKlxuICogUnVudGltZSBvcGVyYXRpb24gdG8gcmVzdG9yZSBhIHNuYXBzaG90dGVkIHZpZXcuXG4gKi9cbmV4cG9ydCBjbGFzcyBSZXN0b3JlVmlld0V4cHIgZXh0ZW5kcyBFeHByZXNzaW9uQmFzZSB7XG4gIG92ZXJyaWRlIHJlYWRvbmx5IGtpbmQgPSBFeHByZXNzaW9uS2luZC5SZXN0b3JlVmlldztcblxuICBjb25zdHJ1Y3RvcihwdWJsaWMgdmlldzogWHJlZklkIHwgby5FeHByZXNzaW9uKSB7XG4gICAgc3VwZXIoKTtcbiAgfVxuXG4gIG92ZXJyaWRlIHZpc2l0RXhwcmVzc2lvbih2aXNpdG9yOiBvLkV4cHJlc3Npb25WaXNpdG9yLCBjb250ZXh0OiBhbnkpOiB2b2lkIHtcbiAgICBpZiAodHlwZW9mIHRoaXMudmlldyAhPT0gJ251bWJlcicpIHtcbiAgICAgIHRoaXMudmlldy52aXNpdEV4cHJlc3Npb24odmlzaXRvciwgY29udGV4dCk7XG4gICAgfVxuICB9XG5cbiAgb3ZlcnJpZGUgaXNFcXVpdmFsZW50KGU6IG8uRXhwcmVzc2lvbik6IGJvb2xlYW4ge1xuICAgIGlmICghKGUgaW5zdGFuY2VvZiBSZXN0b3JlVmlld0V4cHIpIHx8IHR5cGVvZiBlLnZpZXcgIT09IHR5cGVvZiB0aGlzLnZpZXcpIHtcbiAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG5cbiAgICBpZiAodHlwZW9mIHRoaXMudmlldyA9PT0gJ251bWJlcicpIHtcbiAgICAgIHJldHVybiB0aGlzLnZpZXcgPT09IGUudmlldztcbiAgICB9IGVsc2Uge1xuICAgICAgcmV0dXJuIHRoaXMudmlldy5pc0VxdWl2YWxlbnQoZS52aWV3IGFzIG8uRXhwcmVzc2lvbik7XG4gICAgfVxuICB9XG5cbiAgb3ZlcnJpZGUgaXNDb25zdGFudCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gZmFsc2U7XG4gIH1cblxuICBvdmVycmlkZSB0cmFuc2Zvcm1JbnRlcm5hbEV4cHJlc3Npb25zKFxuICAgIHRyYW5zZm9ybTogRXhwcmVzc2lvblRyYW5zZm9ybSxcbiAgICBmbGFnczogVmlzaXRvckNvbnRleHRGbGFnLFxuICApOiB2b2lkIHtcbiAgICBpZiAodHlwZW9mIHRoaXMudmlldyAhPT0gJ251bWJlcicpIHtcbiAgICAgIHRoaXMudmlldyA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKHRoaXMudmlldywgdHJhbnNmb3JtLCBmbGFncyk7XG4gICAgfVxuICB9XG5cbiAgb3ZlcnJpZGUgY2xvbmUoKTogUmVzdG9yZVZpZXdFeHByIHtcbiAgICByZXR1cm4gbmV3IFJlc3RvcmVWaWV3RXhwcih0aGlzLnZpZXcgaW5zdGFuY2VvZiBvLkV4cHJlc3Npb24gPyB0aGlzLnZpZXcuY2xvbmUoKSA6IHRoaXMudmlldyk7XG4gIH1cbn1cblxuLyoqXG4gKiBSdW50aW1lIG9wZXJhdGlvbiB0byByZXNldCB0aGUgY3VycmVudCB2aWV3IGNvbnRleHQgYWZ0ZXIgYFJlc3RvcmVWaWV3YC5cbiAqL1xuZXhwb3J0IGNsYXNzIFJlc2V0Vmlld0V4cHIgZXh0ZW5kcyBFeHByZXNzaW9uQmFzZSB7XG4gIG92ZXJyaWRlIHJlYWRvbmx5IGtpbmQgPSBFeHByZXNzaW9uS2luZC5SZXNldFZpZXc7XG5cbiAgY29uc3RydWN0b3IocHVibGljIGV4cHI6IG8uRXhwcmVzc2lvbikge1xuICAgIHN1cGVyKCk7XG4gIH1cblxuICBvdmVycmlkZSB2aXNpdEV4cHJlc3Npb24odmlzaXRvcjogby5FeHByZXNzaW9uVmlzaXRvciwgY29udGV4dDogYW55KTogYW55IHtcbiAgICB0aGlzLmV4cHIudmlzaXRFeHByZXNzaW9uKHZpc2l0b3IsIGNvbnRleHQpO1xuICB9XG5cbiAgb3ZlcnJpZGUgaXNFcXVpdmFsZW50KGU6IG8uRXhwcmVzc2lvbik6IGJvb2xlYW4ge1xuICAgIHJldHVybiBlIGluc3RhbmNlb2YgUmVzZXRWaWV3RXhwciAmJiB0aGlzLmV4cHIuaXNFcXVpdmFsZW50KGUuZXhwcik7XG4gIH1cblxuICBvdmVycmlkZSBpc0NvbnN0YW50KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiBmYWxzZTtcbiAgfVxuXG4gIG92ZXJyaWRlIHRyYW5zZm9ybUludGVybmFsRXhwcmVzc2lvbnMoXG4gICAgdHJhbnNmb3JtOiBFeHByZXNzaW9uVHJhbnNmb3JtLFxuICAgIGZsYWdzOiBWaXNpdG9yQ29udGV4dEZsYWcsXG4gICk6IHZvaWQge1xuICAgIHRoaXMuZXhwciA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKHRoaXMuZXhwciwgdHJhbnNmb3JtLCBmbGFncyk7XG4gIH1cblxuICBvdmVycmlkZSBjbG9uZSgpOiBSZXNldFZpZXdFeHByIHtcbiAgICByZXR1cm4gbmV3IFJlc2V0Vmlld0V4cHIodGhpcy5leHByLmNsb25lKCkpO1xuICB9XG59XG5cbmV4cG9ydCBjbGFzcyBUd29XYXlCaW5kaW5nU2V0RXhwciBleHRlbmRzIEV4cHJlc3Npb25CYXNlIHtcbiAgb3ZlcnJpZGUgcmVhZG9ubHkga2luZCA9IEV4cHJlc3Npb25LaW5kLlR3b1dheUJpbmRpbmdTZXQ7XG5cbiAgY29uc3RydWN0b3IoXG4gICAgcHVibGljIHRhcmdldDogby5FeHByZXNzaW9uLFxuICAgIHB1YmxpYyB2YWx1ZTogby5FeHByZXNzaW9uLFxuICApIHtcbiAgICBzdXBlcigpO1xuICB9XG5cbiAgb3ZlcnJpZGUgdmlzaXRFeHByZXNzaW9uKHZpc2l0b3I6IG8uRXhwcmVzc2lvblZpc2l0b3IsIGNvbnRleHQ6IGFueSk6IHZvaWQge1xuICAgIHRoaXMudGFyZ2V0LnZpc2l0RXhwcmVzc2lvbih2aXNpdG9yLCBjb250ZXh0KTtcbiAgICB0aGlzLnZhbHVlLnZpc2l0RXhwcmVzc2lvbih2aXNpdG9yLCBjb250ZXh0KTtcbiAgfVxuXG4gIG92ZXJyaWRlIGlzRXF1aXZhbGVudChvdGhlcjogVHdvV2F5QmluZGluZ1NldEV4cHIpOiBib29sZWFuIHtcbiAgICByZXR1cm4gdGhpcy50YXJnZXQuaXNFcXVpdmFsZW50KG90aGVyLnRhcmdldCkgJiYgdGhpcy52YWx1ZS5pc0VxdWl2YWxlbnQob3RoZXIudmFsdWUpO1xuICB9XG5cbiAgb3ZlcnJpZGUgaXNDb25zdGFudCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gZmFsc2U7XG4gIH1cblxuICBvdmVycmlkZSB0cmFuc2Zvcm1JbnRlcm5hbEV4cHJlc3Npb25zKHRyYW5zZm9ybTogRXhwcmVzc2lvblRyYW5zZm9ybSwgZmxhZ3M6IFZpc2l0b3JDb250ZXh0RmxhZykge1xuICAgIHRoaXMudGFyZ2V0ID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24odGhpcy50YXJnZXQsIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgIHRoaXMudmFsdWUgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbih0aGlzLnZhbHVlLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgfVxuXG4gIG92ZXJyaWRlIGNsb25lKCk6IFR3b1dheUJpbmRpbmdTZXRFeHByIHtcbiAgICByZXR1cm4gbmV3IFR3b1dheUJpbmRpbmdTZXRFeHByKHRoaXMudGFyZ2V0LCB0aGlzLnZhbHVlKTtcbiAgfVxufVxuXG4vKipcbiAqIFJlYWQgb2YgYSB2YXJpYWJsZSBkZWNsYXJlZCBhcyBhbiBgaXIuVmFyaWFibGVPcGAgYW5kIHJlZmVyZW5jZWQgdGhyb3VnaCBpdHMgYGlyLlhyZWZJZGAuXG4gKi9cbmV4cG9ydCBjbGFzcyBSZWFkVmFyaWFibGVFeHByIGV4dGVuZHMgRXhwcmVzc2lvbkJhc2Uge1xuICBvdmVycmlkZSByZWFkb25seSBraW5kID0gRXhwcmVzc2lvbktpbmQuUmVhZFZhcmlhYmxlO1xuICBuYW1lOiBzdHJpbmcgfCBudWxsID0gbnVsbDtcbiAgY29uc3RydWN0b3IocmVhZG9ubHkgeHJlZjogWHJlZklkKSB7XG4gICAgc3VwZXIoKTtcbiAgfVxuXG4gIG92ZXJyaWRlIHZpc2l0RXhwcmVzc2lvbigpOiB2b2lkIHt9XG5cbiAgb3ZlcnJpZGUgaXNFcXVpdmFsZW50KG90aGVyOiBvLkV4cHJlc3Npb24pOiBib29sZWFuIHtcbiAgICByZXR1cm4gb3RoZXIgaW5zdGFuY2VvZiBSZWFkVmFyaWFibGVFeHByICYmIG90aGVyLnhyZWYgPT09IHRoaXMueHJlZjtcbiAgfVxuXG4gIG92ZXJyaWRlIGlzQ29uc3RhbnQoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIGZhbHNlO1xuICB9XG5cbiAgb3ZlcnJpZGUgdHJhbnNmb3JtSW50ZXJuYWxFeHByZXNzaW9ucygpOiB2b2lkIHt9XG5cbiAgb3ZlcnJpZGUgY2xvbmUoKTogUmVhZFZhcmlhYmxlRXhwciB7XG4gICAgY29uc3QgZXhwciA9IG5ldyBSZWFkVmFyaWFibGVFeHByKHRoaXMueHJlZik7XG4gICAgZXhwci5uYW1lID0gdGhpcy5uYW1lO1xuICAgIHJldHVybiBleHByO1xuICB9XG59XG5cbmV4cG9ydCBjbGFzcyBQdXJlRnVuY3Rpb25FeHByXG4gIGV4dGVuZHMgRXhwcmVzc2lvbkJhc2VcbiAgaW1wbGVtZW50cyBDb25zdW1lc1ZhcnNUcmFpdCwgVXNlc1Zhck9mZnNldFRyYWl0XG57XG4gIG92ZXJyaWRlIHJlYWRvbmx5IGtpbmQgPSBFeHByZXNzaW9uS2luZC5QdXJlRnVuY3Rpb25FeHByO1xuICByZWFkb25seSBbQ29uc3VtZXNWYXJzVHJhaXRdID0gdHJ1ZTtcbiAgcmVhZG9ubHkgW1VzZXNWYXJPZmZzZXRdID0gdHJ1ZTtcblxuICB2YXJPZmZzZXQ6IG51bWJlciB8IG51bGwgPSBudWxsO1xuXG4gIC8qKlxuICAgKiBUaGUgZXhwcmVzc2lvbiB3aGljaCBzaG91bGQgYmUgbWVtb2l6ZWQgYXMgYSBwdXJlIGNvbXB1dGF0aW9uLlxuICAgKlxuICAgKiBUaGlzIGV4cHJlc3Npb24gY29udGFpbnMgaW50ZXJuYWwgYFB1cmVGdW5jdGlvblBhcmFtZXRlckV4cHJgcywgd2hpY2ggYXJlIHBsYWNlaG9sZGVycyBmb3IgdGhlXG4gICAqIHBvc2l0aW9uYWwgYXJndW1lbnQgZXhwcmVzc2lvbnMgaW4gYGFyZ3MuXG4gICAqL1xuICBib2R5OiBvLkV4cHJlc3Npb24gfCBudWxsO1xuXG4gIC8qKlxuICAgKiBQb3NpdGlvbmFsIGFyZ3VtZW50cyB0byB0aGUgcHVyZSBmdW5jdGlvbiB3aGljaCB3aWxsIG1lbW9pemUgdGhlIGBib2R5YCBleHByZXNzaW9uLCB3aGljaCBhY3RcbiAgICogYXMgbWVtb2l6YXRpb24ga2V5cy5cbiAgICovXG4gIGFyZ3M6IG8uRXhwcmVzc2lvbltdO1xuXG4gIC8qKlxuICAgKiBPbmNlIGV4dHJhY3RlZCB0byB0aGUgYENvbnN0YW50UG9vbGAsIGEgcmVmZXJlbmNlIHRvIHRoZSBmdW5jdGlvbiB3aGljaCBkZWZpbmVzIHRoZSBjb21wdXRhdGlvblxuICAgKiBvZiBgYm9keWAuXG4gICAqL1xuICBmbjogby5FeHByZXNzaW9uIHwgbnVsbCA9IG51bGw7XG5cbiAgY29uc3RydWN0b3IoZXhwcmVzc2lvbjogby5FeHByZXNzaW9uIHwgbnVsbCwgYXJnczogby5FeHByZXNzaW9uW10pIHtcbiAgICBzdXBlcigpO1xuICAgIHRoaXMuYm9keSA9IGV4cHJlc3Npb247XG4gICAgdGhpcy5hcmdzID0gYXJncztcbiAgfVxuXG4gIG92ZXJyaWRlIHZpc2l0RXhwcmVzc2lvbih2aXNpdG9yOiBvLkV4cHJlc3Npb25WaXNpdG9yLCBjb250ZXh0OiBhbnkpIHtcbiAgICB0aGlzLmJvZHk/LnZpc2l0RXhwcmVzc2lvbih2aXNpdG9yLCBjb250ZXh0KTtcbiAgICBmb3IgKGNvbnN0IGFyZyBvZiB0aGlzLmFyZ3MpIHtcbiAgICAgIGFyZy52aXNpdEV4cHJlc3Npb24odmlzaXRvciwgY29udGV4dCk7XG4gICAgfVxuICB9XG5cbiAgb3ZlcnJpZGUgaXNFcXVpdmFsZW50KG90aGVyOiBvLkV4cHJlc3Npb24pOiBib29sZWFuIHtcbiAgICBpZiAoIShvdGhlciBpbnN0YW5jZW9mIFB1cmVGdW5jdGlvbkV4cHIpIHx8IG90aGVyLmFyZ3MubGVuZ3RoICE9PSB0aGlzLmFyZ3MubGVuZ3RoKSB7XG4gICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuXG4gICAgcmV0dXJuIChcbiAgICAgIG90aGVyLmJvZHkgIT09IG51bGwgJiZcbiAgICAgIHRoaXMuYm9keSAhPT0gbnVsbCAmJlxuICAgICAgb3RoZXIuYm9keS5pc0VxdWl2YWxlbnQodGhpcy5ib2R5KSAmJlxuICAgICAgb3RoZXIuYXJncy5ldmVyeSgoYXJnLCBpZHgpID0+IGFyZy5pc0VxdWl2YWxlbnQodGhpcy5hcmdzW2lkeF0pKVxuICAgICk7XG4gIH1cblxuICBvdmVycmlkZSBpc0NvbnN0YW50KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiBmYWxzZTtcbiAgfVxuXG4gIG92ZXJyaWRlIHRyYW5zZm9ybUludGVybmFsRXhwcmVzc2lvbnMoXG4gICAgdHJhbnNmb3JtOiBFeHByZXNzaW9uVHJhbnNmb3JtLFxuICAgIGZsYWdzOiBWaXNpdG9yQ29udGV4dEZsYWcsXG4gICk6IHZvaWQge1xuICAgIGlmICh0aGlzLmJvZHkgIT09IG51bGwpIHtcbiAgICAgIC8vIFRPRE86IGZpZ3VyZSBvdXQgaWYgdGhpcyBpcyB0aGUgcmlnaHQgZmxhZyB0byBwYXNzIGhlcmUuXG4gICAgICB0aGlzLmJvZHkgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihcbiAgICAgICAgdGhpcy5ib2R5LFxuICAgICAgICB0cmFuc2Zvcm0sXG4gICAgICAgIGZsYWdzIHwgVmlzaXRvckNvbnRleHRGbGFnLkluQ2hpbGRPcGVyYXRpb24sXG4gICAgICApO1xuICAgIH0gZWxzZSBpZiAodGhpcy5mbiAhPT0gbnVsbCkge1xuICAgICAgdGhpcy5mbiA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKHRoaXMuZm4sIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgIH1cblxuICAgIGZvciAobGV0IGkgPSAwOyBpIDwgdGhpcy5hcmdzLmxlbmd0aDsgaSsrKSB7XG4gICAgICB0aGlzLmFyZ3NbaV0gPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbih0aGlzLmFyZ3NbaV0sIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgIH1cbiAgfVxuXG4gIG92ZXJyaWRlIGNsb25lKCk6IFB1cmVGdW5jdGlvbkV4cHIge1xuICAgIGNvbnN0IGV4cHIgPSBuZXcgUHVyZUZ1bmN0aW9uRXhwcihcbiAgICAgIHRoaXMuYm9keT8uY2xvbmUoKSA/PyBudWxsLFxuICAgICAgdGhpcy5hcmdzLm1hcCgoYXJnKSA9PiBhcmcuY2xvbmUoKSksXG4gICAgKTtcbiAgICBleHByLmZuID0gdGhpcy5mbj8uY2xvbmUoKSA/PyBudWxsO1xuICAgIGV4cHIudmFyT2Zmc2V0ID0gdGhpcy52YXJPZmZzZXQ7XG4gICAgcmV0dXJuIGV4cHI7XG4gIH1cbn1cblxuZXhwb3J0IGNsYXNzIFB1cmVGdW5jdGlvblBhcmFtZXRlckV4cHIgZXh0ZW5kcyBFeHByZXNzaW9uQmFzZSB7XG4gIG92ZXJyaWRlIHJlYWRvbmx5IGtpbmQgPSBFeHByZXNzaW9uS2luZC5QdXJlRnVuY3Rpb25QYXJhbWV0ZXJFeHByO1xuXG4gIGNvbnN0cnVjdG9yKHB1YmxpYyBpbmRleDogbnVtYmVyKSB7XG4gICAgc3VwZXIoKTtcbiAgfVxuXG4gIG92ZXJyaWRlIHZpc2l0RXhwcmVzc2lvbigpOiB2b2lkIHt9XG5cbiAgb3ZlcnJpZGUgaXNFcXVpdmFsZW50KG90aGVyOiBvLkV4cHJlc3Npb24pOiBib29sZWFuIHtcbiAgICByZXR1cm4gb3RoZXIgaW5zdGFuY2VvZiBQdXJlRnVuY3Rpb25QYXJhbWV0ZXJFeHByICYmIG90aGVyLmluZGV4ID09PSB0aGlzLmluZGV4O1xuICB9XG5cbiAgb3ZlcnJpZGUgaXNDb25zdGFudCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gdHJ1ZTtcbiAgfVxuXG4gIG92ZXJyaWRlIHRyYW5zZm9ybUludGVybmFsRXhwcmVzc2lvbnMoKTogdm9pZCB7fVxuXG4gIG92ZXJyaWRlIGNsb25lKCk6IFB1cmVGdW5jdGlvblBhcmFtZXRlckV4cHIge1xuICAgIHJldHVybiBuZXcgUHVyZUZ1bmN0aW9uUGFyYW1ldGVyRXhwcih0aGlzLmluZGV4KTtcbiAgfVxufVxuXG5leHBvcnQgY2xhc3MgUGlwZUJpbmRpbmdFeHByXG4gIGV4dGVuZHMgRXhwcmVzc2lvbkJhc2VcbiAgaW1wbGVtZW50cyBDb25zdW1lc1ZhcnNUcmFpdCwgVXNlc1Zhck9mZnNldFRyYWl0XG57XG4gIG92ZXJyaWRlIHJlYWRvbmx5IGtpbmQgPSBFeHByZXNzaW9uS2luZC5QaXBlQmluZGluZztcbiAgcmVhZG9ubHkgW0NvbnN1bWVzVmFyc1RyYWl0XSA9IHRydWU7XG4gIHJlYWRvbmx5IFtVc2VzVmFyT2Zmc2V0XSA9IHRydWU7XG5cbiAgdmFyT2Zmc2V0OiBudW1iZXIgfCBudWxsID0gbnVsbDtcblxuICBjb25zdHJ1Y3RvcihcbiAgICByZWFkb25seSB0YXJnZXQ6IFhyZWZJZCxcbiAgICByZWFkb25seSB0YXJnZXRTbG90OiBTbG90SGFuZGxlLFxuICAgIHJlYWRvbmx5IG5hbWU6IHN0cmluZyxcbiAgICByZWFkb25seSBhcmdzOiBvLkV4cHJlc3Npb25bXSxcbiAgKSB7XG4gICAgc3VwZXIoKTtcbiAgfVxuXG4gIG92ZXJyaWRlIHZpc2l0RXhwcmVzc2lvbih2aXNpdG9yOiBvLkV4cHJlc3Npb25WaXNpdG9yLCBjb250ZXh0OiBhbnkpOiB2b2lkIHtcbiAgICBmb3IgKGNvbnN0IGFyZyBvZiB0aGlzLmFyZ3MpIHtcbiAgICAgIGFyZy52aXNpdEV4cHJlc3Npb24odmlzaXRvciwgY29udGV4dCk7XG4gICAgfVxuICB9XG5cbiAgb3ZlcnJpZGUgaXNFcXVpdmFsZW50KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiBmYWxzZTtcbiAgfVxuXG4gIG92ZXJyaWRlIGlzQ29uc3RhbnQoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIGZhbHNlO1xuICB9XG5cbiAgb3ZlcnJpZGUgdHJhbnNmb3JtSW50ZXJuYWxFeHByZXNzaW9ucyhcbiAgICB0cmFuc2Zvcm06IEV4cHJlc3Npb25UcmFuc2Zvcm0sXG4gICAgZmxhZ3M6IFZpc2l0b3JDb250ZXh0RmxhZyxcbiAgKTogdm9pZCB7XG4gICAgZm9yIChsZXQgaWR4ID0gMDsgaWR4IDwgdGhpcy5hcmdzLmxlbmd0aDsgaWR4KyspIHtcbiAgICAgIHRoaXMuYXJnc1tpZHhdID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24odGhpcy5hcmdzW2lkeF0sIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgIH1cbiAgfVxuXG4gIG92ZXJyaWRlIGNsb25lKCkge1xuICAgIGNvbnN0IHIgPSBuZXcgUGlwZUJpbmRpbmdFeHByKFxuICAgICAgdGhpcy50YXJnZXQsXG4gICAgICB0aGlzLnRhcmdldFNsb3QsXG4gICAgICB0aGlzLm5hbWUsXG4gICAgICB0aGlzLmFyZ3MubWFwKChhKSA9PiBhLmNsb25lKCkpLFxuICAgICk7XG4gICAgci52YXJPZmZzZXQgPSB0aGlzLnZhck9mZnNldDtcbiAgICByZXR1cm4gcjtcbiAgfVxufVxuXG5leHBvcnQgY2xhc3MgUGlwZUJpbmRpbmdWYXJpYWRpY0V4cHJcbiAgZXh0ZW5kcyBFeHByZXNzaW9uQmFzZVxuICBpbXBsZW1lbnRzIENvbnN1bWVzVmFyc1RyYWl0LCBVc2VzVmFyT2Zmc2V0VHJhaXRcbntcbiAgb3ZlcnJpZGUgcmVhZG9ubHkga2luZCA9IEV4cHJlc3Npb25LaW5kLlBpcGVCaW5kaW5nVmFyaWFkaWM7XG4gIHJlYWRvbmx5IFtDb25zdW1lc1ZhcnNUcmFpdF0gPSB0cnVlO1xuICByZWFkb25seSBbVXNlc1Zhck9mZnNldF0gPSB0cnVlO1xuXG4gIHZhck9mZnNldDogbnVtYmVyIHwgbnVsbCA9IG51bGw7XG5cbiAgY29uc3RydWN0b3IoXG4gICAgcmVhZG9ubHkgdGFyZ2V0OiBYcmVmSWQsXG4gICAgcmVhZG9ubHkgdGFyZ2V0U2xvdDogU2xvdEhhbmRsZSxcbiAgICByZWFkb25seSBuYW1lOiBzdHJpbmcsXG4gICAgcHVibGljIGFyZ3M6IG8uRXhwcmVzc2lvbixcbiAgICBwdWJsaWMgbnVtQXJnczogbnVtYmVyLFxuICApIHtcbiAgICBzdXBlcigpO1xuICB9XG5cbiAgb3ZlcnJpZGUgdmlzaXRFeHByZXNzaW9uKHZpc2l0b3I6IG8uRXhwcmVzc2lvblZpc2l0b3IsIGNvbnRleHQ6IGFueSk6IHZvaWQge1xuICAgIHRoaXMuYXJncy52aXNpdEV4cHJlc3Npb24odmlzaXRvciwgY29udGV4dCk7XG4gIH1cblxuICBvdmVycmlkZSBpc0VxdWl2YWxlbnQoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIGZhbHNlO1xuICB9XG5cbiAgb3ZlcnJpZGUgaXNDb25zdGFudCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gZmFsc2U7XG4gIH1cblxuICBvdmVycmlkZSB0cmFuc2Zvcm1JbnRlcm5hbEV4cHJlc3Npb25zKFxuICAgIHRyYW5zZm9ybTogRXhwcmVzc2lvblRyYW5zZm9ybSxcbiAgICBmbGFnczogVmlzaXRvckNvbnRleHRGbGFnLFxuICApOiB2b2lkIHtcbiAgICB0aGlzLmFyZ3MgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbih0aGlzLmFyZ3MsIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICB9XG5cbiAgb3ZlcnJpZGUgY2xvbmUoKTogUGlwZUJpbmRpbmdWYXJpYWRpY0V4cHIge1xuICAgIGNvbnN0IHIgPSBuZXcgUGlwZUJpbmRpbmdWYXJpYWRpY0V4cHIoXG4gICAgICB0aGlzLnRhcmdldCxcbiAgICAgIHRoaXMudGFyZ2V0U2xvdCxcbiAgICAgIHRoaXMubmFtZSxcbiAgICAgIHRoaXMuYXJncy5jbG9uZSgpLFxuICAgICAgdGhpcy5udW1BcmdzLFxuICAgICk7XG4gICAgci52YXJPZmZzZXQgPSB0aGlzLnZhck9mZnNldDtcbiAgICByZXR1cm4gcjtcbiAgfVxufVxuXG5leHBvcnQgY2xhc3MgU2FmZVByb3BlcnR5UmVhZEV4cHIgZXh0ZW5kcyBFeHByZXNzaW9uQmFzZSB7XG4gIG92ZXJyaWRlIHJlYWRvbmx5IGtpbmQgPSBFeHByZXNzaW9uS2luZC5TYWZlUHJvcGVydHlSZWFkO1xuXG4gIGNvbnN0cnVjdG9yKFxuICAgIHB1YmxpYyByZWNlaXZlcjogby5FeHByZXNzaW9uLFxuICAgIHB1YmxpYyBuYW1lOiBzdHJpbmcsXG4gICkge1xuICAgIHN1cGVyKCk7XG4gIH1cblxuICAvLyBBbiBhbGlhcyBmb3IgbmFtZSwgd2hpY2ggYWxsb3dzIG90aGVyIGxvZ2ljIHRvIGhhbmRsZSBwcm9wZXJ0eSByZWFkcyBhbmQga2V5ZWQgcmVhZHMgdG9nZXRoZXIuXG4gIGdldCBpbmRleCgpIHtcbiAgICByZXR1cm4gdGhpcy5uYW1lO1xuICB9XG5cbiAgb3ZlcnJpZGUgdmlzaXRFeHByZXNzaW9uKHZpc2l0b3I6IG8uRXhwcmVzc2lvblZpc2l0b3IsIGNvbnRleHQ6IGFueSk6IGFueSB7XG4gICAgdGhpcy5yZWNlaXZlci52aXNpdEV4cHJlc3Npb24odmlzaXRvciwgY29udGV4dCk7XG4gIH1cblxuICBvdmVycmlkZSBpc0VxdWl2YWxlbnQoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIGZhbHNlO1xuICB9XG5cbiAgb3ZlcnJpZGUgaXNDb25zdGFudCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gZmFsc2U7XG4gIH1cblxuICBvdmVycmlkZSB0cmFuc2Zvcm1JbnRlcm5hbEV4cHJlc3Npb25zKFxuICAgIHRyYW5zZm9ybTogRXhwcmVzc2lvblRyYW5zZm9ybSxcbiAgICBmbGFnczogVmlzaXRvckNvbnRleHRGbGFnLFxuICApOiB2b2lkIHtcbiAgICB0aGlzLnJlY2VpdmVyID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24odGhpcy5yZWNlaXZlciwgdHJhbnNmb3JtLCBmbGFncyk7XG4gIH1cblxuICBvdmVycmlkZSBjbG9uZSgpOiBTYWZlUHJvcGVydHlSZWFkRXhwciB7XG4gICAgcmV0dXJuIG5ldyBTYWZlUHJvcGVydHlSZWFkRXhwcih0aGlzLnJlY2VpdmVyLmNsb25lKCksIHRoaXMubmFtZSk7XG4gIH1cbn1cblxuZXhwb3J0IGNsYXNzIFNhZmVLZXllZFJlYWRFeHByIGV4dGVuZHMgRXhwcmVzc2lvbkJhc2Uge1xuICBvdmVycmlkZSByZWFkb25seSBraW5kID0gRXhwcmVzc2lvbktpbmQuU2FmZUtleWVkUmVhZDtcblxuICBjb25zdHJ1Y3RvcihcbiAgICBwdWJsaWMgcmVjZWl2ZXI6IG8uRXhwcmVzc2lvbixcbiAgICBwdWJsaWMgaW5kZXg6IG8uRXhwcmVzc2lvbixcbiAgICBzb3VyY2VTcGFuOiBQYXJzZVNvdXJjZVNwYW4gfCBudWxsLFxuICApIHtcbiAgICBzdXBlcihzb3VyY2VTcGFuKTtcbiAgfVxuXG4gIG92ZXJyaWRlIHZpc2l0RXhwcmVzc2lvbih2aXNpdG9yOiBvLkV4cHJlc3Npb25WaXNpdG9yLCBjb250ZXh0OiBhbnkpOiBhbnkge1xuICAgIHRoaXMucmVjZWl2ZXIudmlzaXRFeHByZXNzaW9uKHZpc2l0b3IsIGNvbnRleHQpO1xuICAgIHRoaXMuaW5kZXgudmlzaXRFeHByZXNzaW9uKHZpc2l0b3IsIGNvbnRleHQpO1xuICB9XG5cbiAgb3ZlcnJpZGUgaXNFcXVpdmFsZW50KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiBmYWxzZTtcbiAgfVxuXG4gIG92ZXJyaWRlIGlzQ29uc3RhbnQoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIGZhbHNlO1xuICB9XG5cbiAgb3ZlcnJpZGUgdHJhbnNmb3JtSW50ZXJuYWxFeHByZXNzaW9ucyhcbiAgICB0cmFuc2Zvcm06IEV4cHJlc3Npb25UcmFuc2Zvcm0sXG4gICAgZmxhZ3M6IFZpc2l0b3JDb250ZXh0RmxhZyxcbiAgKTogdm9pZCB7XG4gICAgdGhpcy5yZWNlaXZlciA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKHRoaXMucmVjZWl2ZXIsIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgIHRoaXMuaW5kZXggPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbih0aGlzLmluZGV4LCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgfVxuXG4gIG92ZXJyaWRlIGNsb25lKCk6IFNhZmVLZXllZFJlYWRFeHByIHtcbiAgICByZXR1cm4gbmV3IFNhZmVLZXllZFJlYWRFeHByKHRoaXMucmVjZWl2ZXIuY2xvbmUoKSwgdGhpcy5pbmRleC5jbG9uZSgpLCB0aGlzLnNvdXJjZVNwYW4pO1xuICB9XG59XG5cbmV4cG9ydCBjbGFzcyBTYWZlSW52b2tlRnVuY3Rpb25FeHByIGV4dGVuZHMgRXhwcmVzc2lvbkJhc2Uge1xuICBvdmVycmlkZSByZWFkb25seSBraW5kID0gRXhwcmVzc2lvbktpbmQuU2FmZUludm9rZUZ1bmN0aW9uO1xuXG4gIGNvbnN0cnVjdG9yKFxuICAgIHB1YmxpYyByZWNlaXZlcjogby5FeHByZXNzaW9uLFxuICAgIHB1YmxpYyBhcmdzOiBvLkV4cHJlc3Npb25bXSxcbiAgKSB7XG4gICAgc3VwZXIoKTtcbiAgfVxuXG4gIG92ZXJyaWRlIHZpc2l0RXhwcmVzc2lvbih2aXNpdG9yOiBvLkV4cHJlc3Npb25WaXNpdG9yLCBjb250ZXh0OiBhbnkpOiBhbnkge1xuICAgIHRoaXMucmVjZWl2ZXIudmlzaXRFeHByZXNzaW9uKHZpc2l0b3IsIGNvbnRleHQpO1xuICAgIGZvciAoY29uc3QgYSBvZiB0aGlzLmFyZ3MpIHtcbiAgICAgIGEudmlzaXRFeHByZXNzaW9uKHZpc2l0b3IsIGNvbnRleHQpO1xuICAgIH1cbiAgfVxuXG4gIG92ZXJyaWRlIGlzRXF1aXZhbGVudCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gZmFsc2U7XG4gIH1cblxuICBvdmVycmlkZSBpc0NvbnN0YW50KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiBmYWxzZTtcbiAgfVxuXG4gIG92ZXJyaWRlIHRyYW5zZm9ybUludGVybmFsRXhwcmVzc2lvbnMoXG4gICAgdHJhbnNmb3JtOiBFeHByZXNzaW9uVHJhbnNmb3JtLFxuICAgIGZsYWdzOiBWaXNpdG9yQ29udGV4dEZsYWcsXG4gICk6IHZvaWQge1xuICAgIHRoaXMucmVjZWl2ZXIgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbih0aGlzLnJlY2VpdmVyLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICBmb3IgKGxldCBpID0gMDsgaSA8IHRoaXMuYXJncy5sZW5ndGg7IGkrKykge1xuICAgICAgdGhpcy5hcmdzW2ldID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24odGhpcy5hcmdzW2ldLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICB9XG4gIH1cblxuICBvdmVycmlkZSBjbG9uZSgpOiBTYWZlSW52b2tlRnVuY3Rpb25FeHByIHtcbiAgICByZXR1cm4gbmV3IFNhZmVJbnZva2VGdW5jdGlvbkV4cHIoXG4gICAgICB0aGlzLnJlY2VpdmVyLmNsb25lKCksXG4gICAgICB0aGlzLmFyZ3MubWFwKChhKSA9PiBhLmNsb25lKCkpLFxuICAgICk7XG4gIH1cbn1cblxuZXhwb3J0IGNsYXNzIFNhZmVUZXJuYXJ5RXhwciBleHRlbmRzIEV4cHJlc3Npb25CYXNlIHtcbiAgb3ZlcnJpZGUgcmVhZG9ubHkga2luZCA9IEV4cHJlc3Npb25LaW5kLlNhZmVUZXJuYXJ5RXhwcjtcblxuICBjb25zdHJ1Y3RvcihcbiAgICBwdWJsaWMgZ3VhcmQ6IG8uRXhwcmVzc2lvbixcbiAgICBwdWJsaWMgZXhwcjogby5FeHByZXNzaW9uLFxuICApIHtcbiAgICBzdXBlcigpO1xuICB9XG5cbiAgb3ZlcnJpZGUgdmlzaXRFeHByZXNzaW9uKHZpc2l0b3I6IG8uRXhwcmVzc2lvblZpc2l0b3IsIGNvbnRleHQ6IGFueSk6IGFueSB7XG4gICAgdGhpcy5ndWFyZC52aXNpdEV4cHJlc3Npb24odmlzaXRvciwgY29udGV4dCk7XG4gICAgdGhpcy5leHByLnZpc2l0RXhwcmVzc2lvbih2aXNpdG9yLCBjb250ZXh0KTtcbiAgfVxuXG4gIG92ZXJyaWRlIGlzRXF1aXZhbGVudCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gZmFsc2U7XG4gIH1cblxuICBvdmVycmlkZSBpc0NvbnN0YW50KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiBmYWxzZTtcbiAgfVxuXG4gIG92ZXJyaWRlIHRyYW5zZm9ybUludGVybmFsRXhwcmVzc2lvbnMoXG4gICAgdHJhbnNmb3JtOiBFeHByZXNzaW9uVHJhbnNmb3JtLFxuICAgIGZsYWdzOiBWaXNpdG9yQ29udGV4dEZsYWcsXG4gICk6IHZvaWQge1xuICAgIHRoaXMuZ3VhcmQgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbih0aGlzLmd1YXJkLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICB0aGlzLmV4cHIgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbih0aGlzLmV4cHIsIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICB9XG5cbiAgb3ZlcnJpZGUgY2xvbmUoKTogU2FmZVRlcm5hcnlFeHByIHtcbiAgICByZXR1cm4gbmV3IFNhZmVUZXJuYXJ5RXhwcih0aGlzLmd1YXJkLmNsb25lKCksIHRoaXMuZXhwci5jbG9uZSgpKTtcbiAgfVxufVxuXG5leHBvcnQgY2xhc3MgRW1wdHlFeHByIGV4dGVuZHMgRXhwcmVzc2lvbkJhc2Uge1xuICBvdmVycmlkZSByZWFkb25seSBraW5kID0gRXhwcmVzc2lvbktpbmQuRW1wdHlFeHByO1xuXG4gIG92ZXJyaWRlIHZpc2l0RXhwcmVzc2lvbih2aXNpdG9yOiBvLkV4cHJlc3Npb25WaXNpdG9yLCBjb250ZXh0OiBhbnkpOiBhbnkge31cblxuICBvdmVycmlkZSBpc0VxdWl2YWxlbnQoZTogRXhwcmVzc2lvbik6IGJvb2xlYW4ge1xuICAgIHJldHVybiBlIGluc3RhbmNlb2YgRW1wdHlFeHByO1xuICB9XG5cbiAgb3ZlcnJpZGUgaXNDb25zdGFudCgpIHtcbiAgICByZXR1cm4gdHJ1ZTtcbiAgfVxuXG4gIG92ZXJyaWRlIGNsb25lKCk6IEVtcHR5RXhwciB7XG4gICAgcmV0dXJuIG5ldyBFbXB0eUV4cHIoKTtcbiAgfVxuXG4gIG92ZXJyaWRlIHRyYW5zZm9ybUludGVybmFsRXhwcmVzc2lvbnMoKTogdm9pZCB7fVxufVxuXG5leHBvcnQgY2xhc3MgQXNzaWduVGVtcG9yYXJ5RXhwciBleHRlbmRzIEV4cHJlc3Npb25CYXNlIHtcbiAgb3ZlcnJpZGUgcmVhZG9ubHkga2luZCA9IEV4cHJlc3Npb25LaW5kLkFzc2lnblRlbXBvcmFyeUV4cHI7XG5cbiAgcHVibGljIG5hbWU6IHN0cmluZyB8IG51bGwgPSBudWxsO1xuXG4gIGNvbnN0cnVjdG9yKFxuICAgIHB1YmxpYyBleHByOiBvLkV4cHJlc3Npb24sXG4gICAgcHVibGljIHhyZWY6IFhyZWZJZCxcbiAgKSB7XG4gICAgc3VwZXIoKTtcbiAgfVxuXG4gIG92ZXJyaWRlIHZpc2l0RXhwcmVzc2lvbih2aXNpdG9yOiBvLkV4cHJlc3Npb25WaXNpdG9yLCBjb250ZXh0OiBhbnkpOiBhbnkge1xuICAgIHRoaXMuZXhwci52aXNpdEV4cHJlc3Npb24odmlzaXRvciwgY29udGV4dCk7XG4gIH1cblxuICBvdmVycmlkZSBpc0VxdWl2YWxlbnQoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIGZhbHNlO1xuICB9XG5cbiAgb3ZlcnJpZGUgaXNDb25zdGFudCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gZmFsc2U7XG4gIH1cblxuICBvdmVycmlkZSB0cmFuc2Zvcm1JbnRlcm5hbEV4cHJlc3Npb25zKFxuICAgIHRyYW5zZm9ybTogRXhwcmVzc2lvblRyYW5zZm9ybSxcbiAgICBmbGFnczogVmlzaXRvckNvbnRleHRGbGFnLFxuICApOiB2b2lkIHtcbiAgICB0aGlzLmV4cHIgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbih0aGlzLmV4cHIsIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICB9XG5cbiAgb3ZlcnJpZGUgY2xvbmUoKTogQXNzaWduVGVtcG9yYXJ5RXhwciB7XG4gICAgY29uc3QgYSA9IG5ldyBBc3NpZ25UZW1wb3JhcnlFeHByKHRoaXMuZXhwci5jbG9uZSgpLCB0aGlzLnhyZWYpO1xuICAgIGEubmFtZSA9IHRoaXMubmFtZTtcbiAgICByZXR1cm4gYTtcbiAgfVxufVxuXG5leHBvcnQgY2xhc3MgUmVhZFRlbXBvcmFyeUV4cHIgZXh0ZW5kcyBFeHByZXNzaW9uQmFzZSB7XG4gIG92ZXJyaWRlIHJlYWRvbmx5IGtpbmQgPSBFeHByZXNzaW9uS2luZC5SZWFkVGVtcG9yYXJ5RXhwcjtcblxuICBwdWJsaWMgbmFtZTogc3RyaW5nIHwgbnVsbCA9IG51bGw7XG5cbiAgY29uc3RydWN0b3IocHVibGljIHhyZWY6IFhyZWZJZCkge1xuICAgIHN1cGVyKCk7XG4gIH1cblxuICBvdmVycmlkZSB2aXNpdEV4cHJlc3Npb24odmlzaXRvcjogby5FeHByZXNzaW9uVmlzaXRvciwgY29udGV4dDogYW55KTogYW55IHt9XG5cbiAgb3ZlcnJpZGUgaXNFcXVpdmFsZW50KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiB0aGlzLnhyZWYgPT09IHRoaXMueHJlZjtcbiAgfVxuXG4gIG92ZXJyaWRlIGlzQ29uc3RhbnQoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIGZhbHNlO1xuICB9XG5cbiAgb3ZlcnJpZGUgdHJhbnNmb3JtSW50ZXJuYWxFeHByZXNzaW9ucyhcbiAgICB0cmFuc2Zvcm06IEV4cHJlc3Npb25UcmFuc2Zvcm0sXG4gICAgZmxhZ3M6IFZpc2l0b3JDb250ZXh0RmxhZyxcbiAgKTogdm9pZCB7fVxuXG4gIG92ZXJyaWRlIGNsb25lKCk6IFJlYWRUZW1wb3JhcnlFeHByIHtcbiAgICBjb25zdCByID0gbmV3IFJlYWRUZW1wb3JhcnlFeHByKHRoaXMueHJlZik7XG4gICAgci5uYW1lID0gdGhpcy5uYW1lO1xuICAgIHJldHVybiByO1xuICB9XG59XG5cbmV4cG9ydCBjbGFzcyBTbG90TGl0ZXJhbEV4cHIgZXh0ZW5kcyBFeHByZXNzaW9uQmFzZSB7XG4gIG92ZXJyaWRlIHJlYWRvbmx5IGtpbmQgPSBFeHByZXNzaW9uS2luZC5TbG90TGl0ZXJhbEV4cHI7XG5cbiAgY29uc3RydWN0b3IocmVhZG9ubHkgc2xvdDogU2xvdEhhbmRsZSkge1xuICAgIHN1cGVyKCk7XG4gIH1cblxuICBvdmVycmlkZSB2aXNpdEV4cHJlc3Npb24odmlzaXRvcjogby5FeHByZXNzaW9uVmlzaXRvciwgY29udGV4dDogYW55KTogYW55IHt9XG5cbiAgb3ZlcnJpZGUgaXNFcXVpdmFsZW50KGU6IEV4cHJlc3Npb24pOiBib29sZWFuIHtcbiAgICByZXR1cm4gZSBpbnN0YW5jZW9mIFNsb3RMaXRlcmFsRXhwciAmJiBlLnNsb3QgPT09IHRoaXMuc2xvdDtcbiAgfVxuXG4gIG92ZXJyaWRlIGlzQ29uc3RhbnQoKSB7XG4gICAgcmV0dXJuIHRydWU7XG4gIH1cblxuICBvdmVycmlkZSBjbG9uZSgpOiBTbG90TGl0ZXJhbEV4cHIge1xuICAgIHJldHVybiBuZXcgU2xvdExpdGVyYWxFeHByKHRoaXMuc2xvdCk7XG4gIH1cblxuICBvdmVycmlkZSB0cmFuc2Zvcm1JbnRlcm5hbEV4cHJlc3Npb25zKCk6IHZvaWQge31cbn1cblxuZXhwb3J0IGNsYXNzIENvbmRpdGlvbmFsQ2FzZUV4cHIgZXh0ZW5kcyBFeHByZXNzaW9uQmFzZSB7XG4gIG92ZXJyaWRlIHJlYWRvbmx5IGtpbmQgPSBFeHByZXNzaW9uS2luZC5Db25kaXRpb25hbENhc2U7XG5cbiAgLyoqXG4gICAqIENyZWF0ZSBhbiBleHByZXNzaW9uIGZvciBvbmUgYnJhbmNoIG9mIGEgY29uZGl0aW9uYWwuXG4gICAqIEBwYXJhbSBleHByIFRoZSBleHByZXNzaW9uIHRvIGJlIHRlc3RlZCBmb3IgdGhpcyBjYXNlLiBNaWdodCBiZSBudWxsLCBhcyBpbiBhbiBgZWxzZWAgY2FzZS5cbiAgICogQHBhcmFtIHRhcmdldCBUaGUgWHJlZiBvZiB0aGUgdmlldyB0byBiZSBkaXNwbGF5ZWQgaWYgdGhpcyBjb25kaXRpb24gaXMgdHJ1ZS5cbiAgICovXG4gIGNvbnN0cnVjdG9yKFxuICAgIHB1YmxpYyBleHByOiBvLkV4cHJlc3Npb24gfCBudWxsLFxuICAgIHJlYWRvbmx5IHRhcmdldDogWHJlZklkLFxuICAgIHJlYWRvbmx5IHRhcmdldFNsb3Q6IFNsb3RIYW5kbGUsXG4gICAgcmVhZG9ubHkgYWxpYXM6IHQuVmFyaWFibGUgfCBudWxsID0gbnVsbCxcbiAgKSB7XG4gICAgc3VwZXIoKTtcbiAgfVxuXG4gIG92ZXJyaWRlIHZpc2l0RXhwcmVzc2lvbih2aXNpdG9yOiBvLkV4cHJlc3Npb25WaXNpdG9yLCBjb250ZXh0OiBhbnkpOiBhbnkge1xuICAgIGlmICh0aGlzLmV4cHIgIT09IG51bGwpIHtcbiAgICAgIHRoaXMuZXhwci52aXNpdEV4cHJlc3Npb24odmlzaXRvciwgY29udGV4dCk7XG4gICAgfVxuICB9XG5cbiAgb3ZlcnJpZGUgaXNFcXVpdmFsZW50KGU6IEV4cHJlc3Npb24pOiBib29sZWFuIHtcbiAgICByZXR1cm4gZSBpbnN0YW5jZW9mIENvbmRpdGlvbmFsQ2FzZUV4cHIgJiYgZS5leHByID09PSB0aGlzLmV4cHI7XG4gIH1cblxuICBvdmVycmlkZSBpc0NvbnN0YW50KCkge1xuICAgIHJldHVybiB0cnVlO1xuICB9XG5cbiAgb3ZlcnJpZGUgY2xvbmUoKTogQ29uZGl0aW9uYWxDYXNlRXhwciB7XG4gICAgcmV0dXJuIG5ldyBDb25kaXRpb25hbENhc2VFeHByKHRoaXMuZXhwciwgdGhpcy50YXJnZXQsIHRoaXMudGFyZ2V0U2xvdCk7XG4gIH1cblxuICBvdmVycmlkZSB0cmFuc2Zvcm1JbnRlcm5hbEV4cHJlc3Npb25zKFxuICAgIHRyYW5zZm9ybTogRXhwcmVzc2lvblRyYW5zZm9ybSxcbiAgICBmbGFnczogVmlzaXRvckNvbnRleHRGbGFnLFxuICApOiB2b2lkIHtcbiAgICBpZiAodGhpcy5leHByICE9PSBudWxsKSB7XG4gICAgICB0aGlzLmV4cHIgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbih0aGlzLmV4cHIsIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgIH1cbiAgfVxufVxuXG5leHBvcnQgY2xhc3MgQ29uc3RDb2xsZWN0ZWRFeHByIGV4dGVuZHMgRXhwcmVzc2lvbkJhc2Uge1xuICBvdmVycmlkZSByZWFkb25seSBraW5kID0gRXhwcmVzc2lvbktpbmQuQ29uc3RDb2xsZWN0ZWQ7XG5cbiAgY29uc3RydWN0b3IocHVibGljIGV4cHI6IG8uRXhwcmVzc2lvbikge1xuICAgIHN1cGVyKCk7XG4gIH1cblxuICBvdmVycmlkZSB0cmFuc2Zvcm1JbnRlcm5hbEV4cHJlc3Npb25zKFxuICAgIHRyYW5zZm9ybTogRXhwcmVzc2lvblRyYW5zZm9ybSxcbiAgICBmbGFnczogVmlzaXRvckNvbnRleHRGbGFnLFxuICApOiB2b2lkIHtcbiAgICB0aGlzLmV4cHIgPSB0cmFuc2Zvcm0odGhpcy5leHByLCBmbGFncyk7XG4gIH1cblxuICBvdmVycmlkZSB2aXNpdEV4cHJlc3Npb24odmlzaXRvcjogby5FeHByZXNzaW9uVmlzaXRvciwgY29udGV4dDogYW55KSB7XG4gICAgdGhpcy5leHByLnZpc2l0RXhwcmVzc2lvbih2aXNpdG9yLCBjb250ZXh0KTtcbiAgfVxuXG4gIG92ZXJyaWRlIGlzRXF1aXZhbGVudChlOiBvLkV4cHJlc3Npb24pOiBib29sZWFuIHtcbiAgICBpZiAoIShlIGluc3RhbmNlb2YgQ29uc3RDb2xsZWN0ZWRFeHByKSkge1xuICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cbiAgICByZXR1cm4gdGhpcy5leHByLmlzRXF1aXZhbGVudChlLmV4cHIpO1xuICB9XG5cbiAgb3ZlcnJpZGUgaXNDb25zdGFudCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gdGhpcy5leHByLmlzQ29uc3RhbnQoKTtcbiAgfVxuXG4gIG92ZXJyaWRlIGNsb25lKCk6IENvbnN0Q29sbGVjdGVkRXhwciB7XG4gICAgcmV0dXJuIG5ldyBDb25zdENvbGxlY3RlZEV4cHIodGhpcy5leHByKTtcbiAgfVxufVxuXG4vKipcbiAqIFZpc2l0cyBhbGwgYEV4cHJlc3Npb25gcyBpbiB0aGUgQVNUIG9mIGBvcGAgd2l0aCB0aGUgYHZpc2l0b3JgIGZ1bmN0aW9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gdmlzaXRFeHByZXNzaW9uc0luT3AoXG4gIG9wOiBDcmVhdGVPcCB8IFVwZGF0ZU9wLFxuICB2aXNpdG9yOiAoZXhwcjogby5FeHByZXNzaW9uLCBmbGFnczogVmlzaXRvckNvbnRleHRGbGFnKSA9PiB2b2lkLFxuKTogdm9pZCB7XG4gIHRyYW5zZm9ybUV4cHJlc3Npb25zSW5PcChcbiAgICBvcCxcbiAgICAoZXhwciwgZmxhZ3MpID0+IHtcbiAgICAgIHZpc2l0b3IoZXhwciwgZmxhZ3MpO1xuICAgICAgcmV0dXJuIGV4cHI7XG4gICAgfSxcbiAgICBWaXNpdG9yQ29udGV4dEZsYWcuTm9uZSxcbiAgKTtcbn1cblxuZXhwb3J0IGVudW0gVmlzaXRvckNvbnRleHRGbGFnIHtcbiAgTm9uZSA9IDBiMDAwMCxcbiAgSW5DaGlsZE9wZXJhdGlvbiA9IDBiMDAwMSxcbn1cblxuZnVuY3Rpb24gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkludGVycG9sYXRpb24oXG4gIGludGVycG9sYXRpb246IEludGVycG9sYXRpb24sXG4gIHRyYW5zZm9ybTogRXhwcmVzc2lvblRyYW5zZm9ybSxcbiAgZmxhZ3M6IFZpc2l0b3JDb250ZXh0RmxhZyxcbikge1xuICBmb3IgKGxldCBpID0gMDsgaSA8IGludGVycG9sYXRpb24uZXhwcmVzc2lvbnMubGVuZ3RoOyBpKyspIHtcbiAgICBpbnRlcnBvbGF0aW9uLmV4cHJlc3Npb25zW2ldID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24oXG4gICAgICBpbnRlcnBvbGF0aW9uLmV4cHJlc3Npb25zW2ldLFxuICAgICAgdHJhbnNmb3JtLFxuICAgICAgZmxhZ3MsXG4gICAgKTtcbiAgfVxufVxuXG4vKipcbiAqIFRyYW5zZm9ybSBhbGwgYEV4cHJlc3Npb25gcyBpbiB0aGUgQVNUIG9mIGBvcGAgd2l0aCB0aGUgYHRyYW5zZm9ybWAgZnVuY3Rpb24uXG4gKlxuICogQWxsIHN1Y2ggb3BlcmF0aW9ucyB3aWxsIGJlIHJlcGxhY2VkIHdpdGggdGhlIHJlc3VsdCBvZiBhcHBseWluZyBgdHJhbnNmb3JtYCwgd2hpY2ggbWF5IGJlIGFuXG4gKiBpZGVudGl0eSB0cmFuc2Zvcm1hdGlvbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIHRyYW5zZm9ybUV4cHJlc3Npb25zSW5PcChcbiAgb3A6IENyZWF0ZU9wIHwgVXBkYXRlT3AsXG4gIHRyYW5zZm9ybTogRXhwcmVzc2lvblRyYW5zZm9ybSxcbiAgZmxhZ3M6IFZpc2l0b3JDb250ZXh0RmxhZyxcbik6IHZvaWQge1xuICBzd2l0Y2ggKG9wLmtpbmQpIHtcbiAgICBjYXNlIE9wS2luZC5TdHlsZVByb3A6XG4gICAgY2FzZSBPcEtpbmQuU3R5bGVNYXA6XG4gICAgY2FzZSBPcEtpbmQuQ2xhc3NQcm9wOlxuICAgIGNhc2UgT3BLaW5kLkNsYXNzTWFwOlxuICAgIGNhc2UgT3BLaW5kLkJpbmRpbmc6XG4gICAgICBpZiAob3AuZXhwcmVzc2lvbiBpbnN0YW5jZW9mIEludGVycG9sYXRpb24pIHtcbiAgICAgICAgdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkludGVycG9sYXRpb24ob3AuZXhwcmVzc2lvbiwgdHJhbnNmb3JtLCBmbGFncyk7XG4gICAgICB9IGVsc2Uge1xuICAgICAgICBvcC5leHByZXNzaW9uID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24ob3AuZXhwcmVzc2lvbiwgdHJhbnNmb3JtLCBmbGFncyk7XG4gICAgICB9XG4gICAgICBicmVhaztcbiAgICBjYXNlIE9wS2luZC5Qcm9wZXJ0eTpcbiAgICBjYXNlIE9wS2luZC5Ib3N0UHJvcGVydHk6XG4gICAgY2FzZSBPcEtpbmQuQXR0cmlidXRlOlxuICAgICAgaWYgKG9wLmV4cHJlc3Npb24gaW5zdGFuY2VvZiBJbnRlcnBvbGF0aW9uKSB7XG4gICAgICAgIHRyYW5zZm9ybUV4cHJlc3Npb25zSW5JbnRlcnBvbGF0aW9uKG9wLmV4cHJlc3Npb24sIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgICAgfSBlbHNlIHtcbiAgICAgICAgb3AuZXhwcmVzc2lvbiA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKG9wLmV4cHJlc3Npb24sIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgICAgfVxuICAgICAgb3Auc2FuaXRpemVyID1cbiAgICAgICAgb3Auc2FuaXRpemVyICYmIHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKG9wLnNhbml0aXplciwgdHJhbnNmb3JtLCBmbGFncyk7XG4gICAgICBicmVhaztcbiAgICBjYXNlIE9wS2luZC5Ud29XYXlQcm9wZXJ0eTpcbiAgICAgIG9wLmV4cHJlc3Npb24gPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihvcC5leHByZXNzaW9uLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICAgIG9wLnNhbml0aXplciA9XG4gICAgICAgIG9wLnNhbml0aXplciAmJiB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihvcC5zYW5pdGl6ZXIsIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgICAgYnJlYWs7XG4gICAgY2FzZSBPcEtpbmQuSTE4bkV4cHJlc3Npb246XG4gICAgICBvcC5leHByZXNzaW9uID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24ob3AuZXhwcmVzc2lvbiwgdHJhbnNmb3JtLCBmbGFncyk7XG4gICAgICBicmVhaztcbiAgICBjYXNlIE9wS2luZC5JbnRlcnBvbGF0ZVRleHQ6XG4gICAgICB0cmFuc2Zvcm1FeHByZXNzaW9uc0luSW50ZXJwb2xhdGlvbihvcC5pbnRlcnBvbGF0aW9uLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICAgIGJyZWFrO1xuICAgIGNhc2UgT3BLaW5kLlN0YXRlbWVudDpcbiAgICAgIHRyYW5zZm9ybUV4cHJlc3Npb25zSW5TdGF0ZW1lbnQob3Auc3RhdGVtZW50LCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICAgIGJyZWFrO1xuICAgIGNhc2UgT3BLaW5kLlZhcmlhYmxlOlxuICAgICAgb3AuaW5pdGlhbGl6ZXIgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihvcC5pbml0aWFsaXplciwgdHJhbnNmb3JtLCBmbGFncyk7XG4gICAgICBicmVhaztcbiAgICBjYXNlIE9wS2luZC5Db25kaXRpb25hbDpcbiAgICAgIGZvciAoY29uc3QgY29uZGl0aW9uIG9mIG9wLmNvbmRpdGlvbnMpIHtcbiAgICAgICAgaWYgKGNvbmRpdGlvbi5leHByID09PSBudWxsKSB7XG4gICAgICAgICAgLy8gVGhpcyBpcyBhIGRlZmF1bHQgY2FzZS5cbiAgICAgICAgICBjb250aW51ZTtcbiAgICAgICAgfVxuICAgICAgICBjb25kaXRpb24uZXhwciA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKGNvbmRpdGlvbi5leHByLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICAgIH1cbiAgICAgIGlmIChvcC5wcm9jZXNzZWQgIT09IG51bGwpIHtcbiAgICAgICAgb3AucHJvY2Vzc2VkID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24ob3AucHJvY2Vzc2VkLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICAgIH1cbiAgICAgIGlmIChvcC5jb250ZXh0VmFsdWUgIT09IG51bGwpIHtcbiAgICAgICAgb3AuY29udGV4dFZhbHVlID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24ob3AuY29udGV4dFZhbHVlLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICAgIH1cbiAgICAgIGJyZWFrO1xuICAgIGNhc2UgT3BLaW5kLkxpc3RlbmVyOlxuICAgIGNhc2UgT3BLaW5kLlR3b1dheUxpc3RlbmVyOlxuICAgICAgZm9yIChjb25zdCBpbm5lck9wIG9mIG9wLmhhbmRsZXJPcHMpIHtcbiAgICAgICAgdHJhbnNmb3JtRXhwcmVzc2lvbnNJbk9wKGlubmVyT3AsIHRyYW5zZm9ybSwgZmxhZ3MgfCBWaXNpdG9yQ29udGV4dEZsYWcuSW5DaGlsZE9wZXJhdGlvbik7XG4gICAgICB9XG4gICAgICBicmVhaztcbiAgICBjYXNlIE9wS2luZC5FeHRyYWN0ZWRBdHRyaWJ1dGU6XG4gICAgICBvcC5leHByZXNzaW9uID1cbiAgICAgICAgb3AuZXhwcmVzc2lvbiAmJiB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihvcC5leHByZXNzaW9uLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICAgIG9wLnRydXN0ZWRWYWx1ZUZuID1cbiAgICAgICAgb3AudHJ1c3RlZFZhbHVlRm4gJiYgdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24ob3AudHJ1c3RlZFZhbHVlRm4sIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgICAgYnJlYWs7XG4gICAgY2FzZSBPcEtpbmQuUmVwZWF0ZXJDcmVhdGU6XG4gICAgICBvcC50cmFjayA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKG9wLnRyYWNrLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICAgIGlmIChvcC50cmFja0J5Rm4gIT09IG51bGwpIHtcbiAgICAgICAgb3AudHJhY2tCeUZuID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24ob3AudHJhY2tCeUZuLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICAgIH1cbiAgICAgIGJyZWFrO1xuICAgIGNhc2UgT3BLaW5kLlJlcGVhdGVyOlxuICAgICAgb3AuY29sbGVjdGlvbiA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKG9wLmNvbGxlY3Rpb24sIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgICAgYnJlYWs7XG4gICAgY2FzZSBPcEtpbmQuRGVmZXI6XG4gICAgICBpZiAob3AubG9hZGluZ0NvbmZpZyAhPT0gbnVsbCkge1xuICAgICAgICBvcC5sb2FkaW5nQ29uZmlnID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24ob3AubG9hZGluZ0NvbmZpZywgdHJhbnNmb3JtLCBmbGFncyk7XG4gICAgICB9XG4gICAgICBpZiAob3AucGxhY2Vob2xkZXJDb25maWcgIT09IG51bGwpIHtcbiAgICAgICAgb3AucGxhY2Vob2xkZXJDb25maWcgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihcbiAgICAgICAgICBvcC5wbGFjZWhvbGRlckNvbmZpZyxcbiAgICAgICAgICB0cmFuc2Zvcm0sXG4gICAgICAgICAgZmxhZ3MsXG4gICAgICAgICk7XG4gICAgICB9XG4gICAgICBpZiAob3AucmVzb2x2ZXJGbiAhPT0gbnVsbCkge1xuICAgICAgICBvcC5yZXNvbHZlckZuID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24ob3AucmVzb2x2ZXJGbiwgdHJhbnNmb3JtLCBmbGFncyk7XG4gICAgICB9XG4gICAgICBicmVhaztcbiAgICBjYXNlIE9wS2luZC5JMThuTWVzc2FnZTpcbiAgICAgIGZvciAoY29uc3QgW3BsYWNlaG9sZGVyLCBleHByXSBvZiBvcC5wYXJhbXMpIHtcbiAgICAgICAgb3AucGFyYW1zLnNldChwbGFjZWhvbGRlciwgdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24oZXhwciwgdHJhbnNmb3JtLCBmbGFncykpO1xuICAgICAgfVxuICAgICAgZm9yIChjb25zdCBbcGxhY2Vob2xkZXIsIGV4cHJdIG9mIG9wLnBvc3Rwcm9jZXNzaW5nUGFyYW1zKSB7XG4gICAgICAgIG9wLnBvc3Rwcm9jZXNzaW5nUGFyYW1zLnNldChcbiAgICAgICAgICBwbGFjZWhvbGRlcixcbiAgICAgICAgICB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihleHByLCB0cmFuc2Zvcm0sIGZsYWdzKSxcbiAgICAgICAgKTtcbiAgICAgIH1cbiAgICAgIGJyZWFrO1xuICAgIGNhc2UgT3BLaW5kLkRlZmVyV2hlbjpcbiAgICAgIG9wLmV4cHIgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihvcC5leHByLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICAgIGJyZWFrO1xuICAgIGNhc2UgT3BLaW5kLlN0b3JlTGV0OlxuICAgICAgb3AudmFsdWUgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihvcC52YWx1ZSwgdHJhbnNmb3JtLCBmbGFncyk7XG4gICAgICBicmVhaztcbiAgICBjYXNlIE9wS2luZC5BZHZhbmNlOlxuICAgIGNhc2UgT3BLaW5kLkNvbnRhaW5lcjpcbiAgICBjYXNlIE9wS2luZC5Db250YWluZXJFbmQ6XG4gICAgY2FzZSBPcEtpbmQuQ29udGFpbmVyU3RhcnQ6XG4gICAgY2FzZSBPcEtpbmQuRGVmZXJPbjpcbiAgICBjYXNlIE9wS2luZC5EaXNhYmxlQmluZGluZ3M6XG4gICAgY2FzZSBPcEtpbmQuRWxlbWVudDpcbiAgICBjYXNlIE9wS2luZC5FbGVtZW50RW5kOlxuICAgIGNhc2UgT3BLaW5kLkVsZW1lbnRTdGFydDpcbiAgICBjYXNlIE9wS2luZC5FbmFibGVCaW5kaW5nczpcbiAgICBjYXNlIE9wS2luZC5JMThuOlxuICAgIGNhc2UgT3BLaW5kLkkxOG5BcHBseTpcbiAgICBjYXNlIE9wS2luZC5JMThuQ29udGV4dDpcbiAgICBjYXNlIE9wS2luZC5JMThuRW5kOlxuICAgIGNhc2UgT3BLaW5kLkkxOG5TdGFydDpcbiAgICBjYXNlIE9wS2luZC5JY3VFbmQ6XG4gICAgY2FzZSBPcEtpbmQuSWN1U3RhcnQ6XG4gICAgY2FzZSBPcEtpbmQuTmFtZXNwYWNlOlxuICAgIGNhc2UgT3BLaW5kLlBpcGU6XG4gICAgY2FzZSBPcEtpbmQuUHJvamVjdGlvbjpcbiAgICBjYXNlIE9wS2luZC5Qcm9qZWN0aW9uRGVmOlxuICAgIGNhc2UgT3BLaW5kLlRlbXBsYXRlOlxuICAgIGNhc2UgT3BLaW5kLlRleHQ6XG4gICAgY2FzZSBPcEtpbmQuSTE4bkF0dHJpYnV0ZXM6XG4gICAgY2FzZSBPcEtpbmQuSWN1UGxhY2Vob2xkZXI6XG4gICAgY2FzZSBPcEtpbmQuRGVjbGFyZUxldDpcbiAgICAgIC8vIFRoZXNlIG9wZXJhdGlvbnMgY29udGFpbiBubyBleHByZXNzaW9ucy5cbiAgICAgIGJyZWFrO1xuICAgIGRlZmF1bHQ6XG4gICAgICB0aHJvdyBuZXcgRXJyb3IoYEFzc2VydGlvbkVycm9yOiB0cmFuc2Zvcm1FeHByZXNzaW9uc0luT3AgZG9lc24ndCBoYW5kbGUgJHtPcEtpbmRbb3Aua2luZF19YCk7XG4gIH1cbn1cblxuLyoqXG4gKiBUcmFuc2Zvcm0gYWxsIGBFeHByZXNzaW9uYHMgaW4gdGhlIEFTVCBvZiBgZXhwcmAgd2l0aCB0aGUgYHRyYW5zZm9ybWAgZnVuY3Rpb24uXG4gKlxuICogQWxsIHN1Y2ggb3BlcmF0aW9ucyB3aWxsIGJlIHJlcGxhY2VkIHdpdGggdGhlIHJlc3VsdCBvZiBhcHBseWluZyBgdHJhbnNmb3JtYCwgd2hpY2ggbWF5IGJlIGFuXG4gKiBpZGVudGl0eSB0cmFuc2Zvcm1hdGlvbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKFxuICBleHByOiBvLkV4cHJlc3Npb24sXG4gIHRyYW5zZm9ybTogRXhwcmVzc2lvblRyYW5zZm9ybSxcbiAgZmxhZ3M6IFZpc2l0b3JDb250ZXh0RmxhZyxcbik6IG8uRXhwcmVzc2lvbiB7XG4gIGlmIChleHByIGluc3RhbmNlb2YgRXhwcmVzc2lvbkJhc2UpIHtcbiAgICBleHByLnRyYW5zZm9ybUludGVybmFsRXhwcmVzc2lvbnModHJhbnNmb3JtLCBmbGFncyk7XG4gIH0gZWxzZSBpZiAoZXhwciBpbnN0YW5jZW9mIG8uQmluYXJ5T3BlcmF0b3JFeHByKSB7XG4gICAgZXhwci5saHMgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihleHByLmxocywgdHJhbnNmb3JtLCBmbGFncyk7XG4gICAgZXhwci5yaHMgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihleHByLnJocywgdHJhbnNmb3JtLCBmbGFncyk7XG4gIH0gZWxzZSBpZiAoZXhwciBpbnN0YW5jZW9mIG8uVW5hcnlPcGVyYXRvckV4cHIpIHtcbiAgICBleHByLmV4cHIgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihleHByLmV4cHIsIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICB9IGVsc2UgaWYgKGV4cHIgaW5zdGFuY2VvZiBvLlJlYWRQcm9wRXhwcikge1xuICAgIGV4cHIucmVjZWl2ZXIgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihleHByLnJlY2VpdmVyLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgfSBlbHNlIGlmIChleHByIGluc3RhbmNlb2Ygby5SZWFkS2V5RXhwcikge1xuICAgIGV4cHIucmVjZWl2ZXIgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihleHByLnJlY2VpdmVyLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICBleHByLmluZGV4ID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24oZXhwci5pbmRleCwgdHJhbnNmb3JtLCBmbGFncyk7XG4gIH0gZWxzZSBpZiAoZXhwciBpbnN0YW5jZW9mIG8uV3JpdGVQcm9wRXhwcikge1xuICAgIGV4cHIucmVjZWl2ZXIgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihleHByLnJlY2VpdmVyLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICBleHByLnZhbHVlID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24oZXhwci52YWx1ZSwgdHJhbnNmb3JtLCBmbGFncyk7XG4gIH0gZWxzZSBpZiAoZXhwciBpbnN0YW5jZW9mIG8uV3JpdGVLZXlFeHByKSB7XG4gICAgZXhwci5yZWNlaXZlciA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKGV4cHIucmVjZWl2ZXIsIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgIGV4cHIuaW5kZXggPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihleHByLmluZGV4LCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICBleHByLnZhbHVlID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24oZXhwci52YWx1ZSwgdHJhbnNmb3JtLCBmbGFncyk7XG4gIH0gZWxzZSBpZiAoZXhwciBpbnN0YW5jZW9mIG8uSW52b2tlRnVuY3Rpb25FeHByKSB7XG4gICAgZXhwci5mbiA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKGV4cHIuZm4sIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgIGZvciAobGV0IGkgPSAwOyBpIDwgZXhwci5hcmdzLmxlbmd0aDsgaSsrKSB7XG4gICAgICBleHByLmFyZ3NbaV0gPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihleHByLmFyZ3NbaV0sIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgIH1cbiAgfSBlbHNlIGlmIChleHByIGluc3RhbmNlb2Ygby5MaXRlcmFsQXJyYXlFeHByKSB7XG4gICAgZm9yIChsZXQgaSA9IDA7IGkgPCBleHByLmVudHJpZXMubGVuZ3RoOyBpKyspIHtcbiAgICAgIGV4cHIuZW50cmllc1tpXSA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKGV4cHIuZW50cmllc1tpXSwgdHJhbnNmb3JtLCBmbGFncyk7XG4gICAgfVxuICB9IGVsc2UgaWYgKGV4cHIgaW5zdGFuY2VvZiBvLkxpdGVyYWxNYXBFeHByKSB7XG4gICAgZm9yIChsZXQgaSA9IDA7IGkgPCBleHByLmVudHJpZXMubGVuZ3RoOyBpKyspIHtcbiAgICAgIGV4cHIuZW50cmllc1tpXS52YWx1ZSA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKFxuICAgICAgICBleHByLmVudHJpZXNbaV0udmFsdWUsXG4gICAgICAgIHRyYW5zZm9ybSxcbiAgICAgICAgZmxhZ3MsXG4gICAgICApO1xuICAgIH1cbiAgfSBlbHNlIGlmIChleHByIGluc3RhbmNlb2Ygby5Db25kaXRpb25hbEV4cHIpIHtcbiAgICBleHByLmNvbmRpdGlvbiA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKGV4cHIuY29uZGl0aW9uLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICBleHByLnRydWVDYXNlID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24oZXhwci50cnVlQ2FzZSwgdHJhbnNmb3JtLCBmbGFncyk7XG4gICAgaWYgKGV4cHIuZmFsc2VDYXNlICE9PSBudWxsKSB7XG4gICAgICBleHByLmZhbHNlQ2FzZSA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKGV4cHIuZmFsc2VDYXNlLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICB9XG4gIH0gZWxzZSBpZiAoZXhwciBpbnN0YW5jZW9mIG8uVHlwZW9mRXhwcikge1xuICAgIGV4cHIuZXhwciA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKGV4cHIuZXhwciwgdHJhbnNmb3JtLCBmbGFncyk7XG4gIH0gZWxzZSBpZiAoZXhwciBpbnN0YW5jZW9mIG8uV3JpdGVWYXJFeHByKSB7XG4gICAgZXhwci52YWx1ZSA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKGV4cHIudmFsdWUsIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICB9IGVsc2UgaWYgKGV4cHIgaW5zdGFuY2VvZiBvLkxvY2FsaXplZFN0cmluZykge1xuICAgIGZvciAobGV0IGkgPSAwOyBpIDwgZXhwci5leHByZXNzaW9ucy5sZW5ndGg7IGkrKykge1xuICAgICAgZXhwci5leHByZXNzaW9uc1tpXSA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKGV4cHIuZXhwcmVzc2lvbnNbaV0sIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgIH1cbiAgfSBlbHNlIGlmIChleHByIGluc3RhbmNlb2Ygby5Ob3RFeHByKSB7XG4gICAgZXhwci5jb25kaXRpb24gPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihleHByLmNvbmRpdGlvbiwgdHJhbnNmb3JtLCBmbGFncyk7XG4gIH0gZWxzZSBpZiAoZXhwciBpbnN0YW5jZW9mIG8uVGFnZ2VkVGVtcGxhdGVFeHByKSB7XG4gICAgZXhwci50YWcgPSB0cmFuc2Zvcm1FeHByZXNzaW9uc0luRXhwcmVzc2lvbihleHByLnRhZywgdHJhbnNmb3JtLCBmbGFncyk7XG4gICAgZXhwci50ZW1wbGF0ZS5leHByZXNzaW9ucyA9IGV4cHIudGVtcGxhdGUuZXhwcmVzc2lvbnMubWFwKChlKSA9PlxuICAgICAgdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24oZSwgdHJhbnNmb3JtLCBmbGFncyksXG4gICAgKTtcbiAgfSBlbHNlIGlmIChleHByIGluc3RhbmNlb2Ygby5BcnJvd0Z1bmN0aW9uRXhwcikge1xuICAgIGlmIChBcnJheS5pc0FycmF5KGV4cHIuYm9keSkpIHtcbiAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgZXhwci5ib2R5Lmxlbmd0aDsgaSsrKSB7XG4gICAgICAgIHRyYW5zZm9ybUV4cHJlc3Npb25zSW5TdGF0ZW1lbnQoZXhwci5ib2R5W2ldLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICAgIH1cbiAgICB9IGVsc2Uge1xuICAgICAgZXhwci5ib2R5ID0gdHJhbnNmb3JtRXhwcmVzc2lvbnNJbkV4cHJlc3Npb24oZXhwci5ib2R5LCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICB9XG4gIH0gZWxzZSBpZiAoZXhwciBpbnN0YW5jZW9mIG8uV3JhcHBlZE5vZGVFeHByKSB7XG4gICAgLy8gVE9ETzogRG8gd2UgbmVlZCB0byB0cmFuc2Zvcm0gYW55IFRTIG5vZGVzIG5lc3RlZCBpbnNpZGUgb2YgdGhpcyBleHByZXNzaW9uP1xuICB9IGVsc2UgaWYgKFxuICAgIGV4cHIgaW5zdGFuY2VvZiBvLlJlYWRWYXJFeHByIHx8XG4gICAgZXhwciBpbnN0YW5jZW9mIG8uRXh0ZXJuYWxFeHByIHx8XG4gICAgZXhwciBpbnN0YW5jZW9mIG8uTGl0ZXJhbEV4cHJcbiAgKSB7XG4gICAgLy8gTm8gYWN0aW9uIGZvciB0aGVzZSB0eXBlcy5cbiAgfSBlbHNlIHtcbiAgICB0aHJvdyBuZXcgRXJyb3IoYFVuaGFuZGxlZCBleHByZXNzaW9uIGtpbmQ6ICR7ZXhwci5jb25zdHJ1Y3Rvci5uYW1lfWApO1xuICB9XG4gIHJldHVybiB0cmFuc2Zvcm0oZXhwciwgZmxhZ3MpO1xufVxuXG4vKipcbiAqIFRyYW5zZm9ybSBhbGwgYEV4cHJlc3Npb25gcyBpbiB0aGUgQVNUIG9mIGBzdG10YCB3aXRoIHRoZSBgdHJhbnNmb3JtYCBmdW5jdGlvbi5cbiAqXG4gKiBBbGwgc3VjaCBvcGVyYXRpb25zIHdpbGwgYmUgcmVwbGFjZWQgd2l0aCB0aGUgcmVzdWx0IG9mIGFwcGx5aW5nIGB0cmFuc2Zvcm1gLCB3aGljaCBtYXkgYmUgYW5cbiAqIGlkZW50aXR5IHRyYW5zZm9ybWF0aW9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gdHJhbnNmb3JtRXhwcmVzc2lvbnNJblN0YXRlbWVudChcbiAgc3RtdDogby5TdGF0ZW1lbnQsXG4gIHRyYW5zZm9ybTogRXhwcmVzc2lvblRyYW5zZm9ybSxcbiAgZmxhZ3M6IFZpc2l0b3JDb250ZXh0RmxhZyxcbik6IHZvaWQge1xuICBpZiAoc3RtdCBpbnN0YW5jZW9mIG8uRXhwcmVzc2lvblN0YXRlbWVudCkge1xuICAgIHN0bXQuZXhwciA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKHN0bXQuZXhwciwgdHJhbnNmb3JtLCBmbGFncyk7XG4gIH0gZWxzZSBpZiAoc3RtdCBpbnN0YW5jZW9mIG8uUmV0dXJuU3RhdGVtZW50KSB7XG4gICAgc3RtdC52YWx1ZSA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKHN0bXQudmFsdWUsIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICB9IGVsc2UgaWYgKHN0bXQgaW5zdGFuY2VvZiBvLkRlY2xhcmVWYXJTdG10KSB7XG4gICAgaWYgKHN0bXQudmFsdWUgIT09IHVuZGVmaW5lZCkge1xuICAgICAgc3RtdC52YWx1ZSA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKHN0bXQudmFsdWUsIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgIH1cbiAgfSBlbHNlIGlmIChzdG10IGluc3RhbmNlb2Ygby5JZlN0bXQpIHtcbiAgICBzdG10LmNvbmRpdGlvbiA9IHRyYW5zZm9ybUV4cHJlc3Npb25zSW5FeHByZXNzaW9uKHN0bXQuY29uZGl0aW9uLCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICBmb3IgKGNvbnN0IGNhc2VTdGF0ZW1lbnQgb2Ygc3RtdC50cnVlQ2FzZSkge1xuICAgICAgdHJhbnNmb3JtRXhwcmVzc2lvbnNJblN0YXRlbWVudChjYXNlU3RhdGVtZW50LCB0cmFuc2Zvcm0sIGZsYWdzKTtcbiAgICB9XG4gICAgZm9yIChjb25zdCBjYXNlU3RhdGVtZW50IG9mIHN0bXQuZmFsc2VDYXNlKSB7XG4gICAgICB0cmFuc2Zvcm1FeHByZXNzaW9uc0luU3RhdGVtZW50KGNhc2VTdGF0ZW1lbnQsIHRyYW5zZm9ybSwgZmxhZ3MpO1xuICAgIH1cbiAgfSBlbHNlIHtcbiAgICB0aHJvdyBuZXcgRXJyb3IoYFVuaGFuZGxlZCBzdGF0ZW1lbnQga2luZDogJHtzdG10LmNvbnN0cnVjdG9yLm5hbWV9YCk7XG4gIH1cbn1cblxuLyoqXG4gKiBDaGVja3Mgd2hldGhlciB0aGUgZ2l2ZW4gZXhwcmVzc2lvbiBpcyBhIHN0cmluZyBsaXRlcmFsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gaXNTdHJpbmdMaXRlcmFsKGV4cHI6IG8uRXhwcmVzc2lvbik6IGV4cHIgaXMgby5MaXRlcmFsRXhwciAmIHt2YWx1ZTogc3RyaW5nfSB7XG4gIHJldHVybiBleHByIGluc3RhbmNlb2Ygby5MaXRlcmFsRXhwciAmJiB0eXBlb2YgZXhwci52YWx1ZSA9PT0gJ3N0cmluZyc7XG59XG4iXX0=