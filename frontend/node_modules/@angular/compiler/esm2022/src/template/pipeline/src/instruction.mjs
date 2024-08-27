/**
 * @license
 * Copyright Google LLC All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
import * as o from '../../../output/output_ast';
import { Identifiers } from '../../../render3/r3_identifiers';
import * as ir from '../ir';
// This file contains helpers for generating calls to Ivy instructions. In particular, each
// instruction type is represented as a function, which may select a specific instruction variant
// depending on the exact arguments.
export function element(slot, tag, constIndex, localRefIndex, sourceSpan) {
    return elementOrContainerBase(Identifiers.element, slot, tag, constIndex, localRefIndex, sourceSpan);
}
export function elementStart(slot, tag, constIndex, localRefIndex, sourceSpan) {
    return elementOrContainerBase(Identifiers.elementStart, slot, tag, constIndex, localRefIndex, sourceSpan);
}
function elementOrContainerBase(instruction, slot, tag, constIndex, localRefIndex, sourceSpan) {
    const args = [o.literal(slot)];
    if (tag !== null) {
        args.push(o.literal(tag));
    }
    if (localRefIndex !== null) {
        args.push(o.literal(constIndex), // might be null, but that's okay.
        o.literal(localRefIndex));
    }
    else if (constIndex !== null) {
        args.push(o.literal(constIndex));
    }
    return call(instruction, args, sourceSpan);
}
export function elementEnd(sourceSpan) {
    return call(Identifiers.elementEnd, [], sourceSpan);
}
export function elementContainerStart(slot, constIndex, localRefIndex, sourceSpan) {
    return elementOrContainerBase(Identifiers.elementContainerStart, slot, 
    /* tag */ null, constIndex, localRefIndex, sourceSpan);
}
export function elementContainer(slot, constIndex, localRefIndex, sourceSpan) {
    return elementOrContainerBase(Identifiers.elementContainer, slot, 
    /* tag */ null, constIndex, localRefIndex, sourceSpan);
}
export function elementContainerEnd() {
    return call(Identifiers.elementContainerEnd, [], null);
}
export function template(slot, templateFnRef, decls, vars, tag, constIndex, localRefs, sourceSpan) {
    const args = [
        o.literal(slot),
        templateFnRef,
        o.literal(decls),
        o.literal(vars),
        o.literal(tag),
        o.literal(constIndex),
    ];
    if (localRefs !== null) {
        args.push(o.literal(localRefs));
        args.push(o.importExpr(Identifiers.templateRefExtractor));
    }
    while (args[args.length - 1].isEquivalent(o.NULL_EXPR)) {
        args.pop();
    }
    return call(Identifiers.templateCreate, args, sourceSpan);
}
export function disableBindings() {
    return call(Identifiers.disableBindings, [], null);
}
export function enableBindings() {
    return call(Identifiers.enableBindings, [], null);
}
export function listener(name, handlerFn, eventTargetResolver, syntheticHost, sourceSpan) {
    const args = [o.literal(name), handlerFn];
    if (eventTargetResolver !== null) {
        args.push(o.literal(false)); // `useCapture` flag, defaults to `false`
        args.push(o.importExpr(eventTargetResolver));
    }
    return call(syntheticHost ? Identifiers.syntheticHostListener : Identifiers.listener, args, sourceSpan);
}
export function twoWayBindingSet(target, value) {
    return o.importExpr(Identifiers.twoWayBindingSet).callFn([target, value]);
}
export function twoWayListener(name, handlerFn, sourceSpan) {
    return call(Identifiers.twoWayListener, [o.literal(name), handlerFn], sourceSpan);
}
export function pipe(slot, name) {
    return call(Identifiers.pipe, [o.literal(slot), o.literal(name)], null);
}
export function namespaceHTML() {
    return call(Identifiers.namespaceHTML, [], null);
}
export function namespaceSVG() {
    return call(Identifiers.namespaceSVG, [], null);
}
export function namespaceMath() {
    return call(Identifiers.namespaceMathML, [], null);
}
export function advance(delta, sourceSpan) {
    return call(Identifiers.advance, delta > 1 ? [o.literal(delta)] : [], sourceSpan);
}
export function reference(slot) {
    return o.importExpr(Identifiers.reference).callFn([o.literal(slot)]);
}
export function nextContext(steps) {
    return o.importExpr(Identifiers.nextContext).callFn(steps === 1 ? [] : [o.literal(steps)]);
}
export function getCurrentView() {
    return o.importExpr(Identifiers.getCurrentView).callFn([]);
}
export function restoreView(savedView) {
    return o.importExpr(Identifiers.restoreView).callFn([savedView]);
}
export function resetView(returnValue) {
    return o.importExpr(Identifiers.resetView).callFn([returnValue]);
}
export function text(slot, initialValue, sourceSpan) {
    const args = [o.literal(slot, null)];
    if (initialValue !== '') {
        args.push(o.literal(initialValue));
    }
    return call(Identifiers.text, args, sourceSpan);
}
export function defer(selfSlot, primarySlot, dependencyResolverFn, loadingSlot, placeholderSlot, errorSlot, loadingConfig, placeholderConfig, enableTimerScheduling, sourceSpan) {
    const args = [
        o.literal(selfSlot),
        o.literal(primarySlot),
        dependencyResolverFn ?? o.literal(null),
        o.literal(loadingSlot),
        o.literal(placeholderSlot),
        o.literal(errorSlot),
        loadingConfig ?? o.literal(null),
        placeholderConfig ?? o.literal(null),
        enableTimerScheduling ? o.importExpr(Identifiers.deferEnableTimerScheduling) : o.literal(null),
    ];
    let expr;
    while ((expr = args[args.length - 1]) !== null &&
        expr instanceof o.LiteralExpr &&
        expr.value === null) {
        args.pop();
    }
    return call(Identifiers.defer, args, sourceSpan);
}
const deferTriggerToR3TriggerInstructionsMap = new Map([
    [ir.DeferTriggerKind.Idle, [Identifiers.deferOnIdle, Identifiers.deferPrefetchOnIdle]],
    [
        ir.DeferTriggerKind.Immediate,
        [Identifiers.deferOnImmediate, Identifiers.deferPrefetchOnImmediate],
    ],
    [ir.DeferTriggerKind.Timer, [Identifiers.deferOnTimer, Identifiers.deferPrefetchOnTimer]],
    [ir.DeferTriggerKind.Hover, [Identifiers.deferOnHover, Identifiers.deferPrefetchOnHover]],
    [
        ir.DeferTriggerKind.Interaction,
        [Identifiers.deferOnInteraction, Identifiers.deferPrefetchOnInteraction],
    ],
    [
        ir.DeferTriggerKind.Viewport,
        [Identifiers.deferOnViewport, Identifiers.deferPrefetchOnViewport],
    ],
]);
export function deferOn(trigger, args, prefetch, sourceSpan) {
    const instructions = deferTriggerToR3TriggerInstructionsMap.get(trigger);
    if (instructions === undefined) {
        throw new Error(`Unable to determine instruction for trigger ${trigger}`);
    }
    const instructionToCall = prefetch ? instructions[1] : instructions[0];
    return call(instructionToCall, args.map((a) => o.literal(a)), sourceSpan);
}
export function projectionDef(def) {
    return call(Identifiers.projectionDef, def ? [def] : [], null);
}
export function projection(slot, projectionSlotIndex, attributes, fallbackFnName, fallbackDecls, fallbackVars, sourceSpan) {
    const args = [o.literal(slot)];
    if (projectionSlotIndex !== 0 || attributes !== null || fallbackFnName !== null) {
        args.push(o.literal(projectionSlotIndex));
        if (attributes !== null) {
            args.push(attributes);
        }
        if (fallbackFnName !== null) {
            if (attributes === null) {
                args.push(o.literal(null));
            }
            args.push(o.variable(fallbackFnName), o.literal(fallbackDecls), o.literal(fallbackVars));
        }
    }
    return call(Identifiers.projection, args, sourceSpan);
}
export function i18nStart(slot, constIndex, subTemplateIndex, sourceSpan) {
    const args = [o.literal(slot), o.literal(constIndex)];
    if (subTemplateIndex !== null) {
        args.push(o.literal(subTemplateIndex));
    }
    return call(Identifiers.i18nStart, args, sourceSpan);
}
export function repeaterCreate(slot, viewFnName, decls, vars, tag, constIndex, trackByFn, trackByUsesComponentInstance, emptyViewFnName, emptyDecls, emptyVars, emptyTag, emptyConstIndex, sourceSpan) {
    const args = [
        o.literal(slot),
        o.variable(viewFnName),
        o.literal(decls),
        o.literal(vars),
        o.literal(tag),
        o.literal(constIndex),
        trackByFn,
    ];
    if (trackByUsesComponentInstance || emptyViewFnName !== null) {
        args.push(o.literal(trackByUsesComponentInstance));
        if (emptyViewFnName !== null) {
            args.push(o.variable(emptyViewFnName), o.literal(emptyDecls), o.literal(emptyVars));
            if (emptyTag !== null || emptyConstIndex !== null) {
                args.push(o.literal(emptyTag));
            }
            if (emptyConstIndex !== null) {
                args.push(o.literal(emptyConstIndex));
            }
        }
    }
    return call(Identifiers.repeaterCreate, args, sourceSpan);
}
export function repeater(collection, sourceSpan) {
    return call(Identifiers.repeater, [collection], sourceSpan);
}
export function deferWhen(prefetch, expr, sourceSpan) {
    return call(prefetch ? Identifiers.deferPrefetchWhen : Identifiers.deferWhen, [expr], sourceSpan);
}
export function declareLet(slot, sourceSpan) {
    return call(Identifiers.declareLet, [o.literal(slot)], sourceSpan);
}
export function storeLet(value, sourceSpan) {
    return o.importExpr(Identifiers.storeLet).callFn([value], sourceSpan);
}
export function readContextLet(slot) {
    return o.importExpr(Identifiers.readContextLet).callFn([o.literal(slot)]);
}
export function i18n(slot, constIndex, subTemplateIndex, sourceSpan) {
    const args = [o.literal(slot), o.literal(constIndex)];
    if (subTemplateIndex) {
        args.push(o.literal(subTemplateIndex));
    }
    return call(Identifiers.i18n, args, sourceSpan);
}
export function i18nEnd(endSourceSpan) {
    return call(Identifiers.i18nEnd, [], endSourceSpan);
}
export function i18nAttributes(slot, i18nAttributesConfig) {
    const args = [o.literal(slot), o.literal(i18nAttributesConfig)];
    return call(Identifiers.i18nAttributes, args, null);
}
export function property(name, expression, sanitizer, sourceSpan) {
    const args = [o.literal(name), expression];
    if (sanitizer !== null) {
        args.push(sanitizer);
    }
    return call(Identifiers.property, args, sourceSpan);
}
export function twoWayProperty(name, expression, sanitizer, sourceSpan) {
    const args = [o.literal(name), expression];
    if (sanitizer !== null) {
        args.push(sanitizer);
    }
    return call(Identifiers.twoWayProperty, args, sourceSpan);
}
export function attribute(name, expression, sanitizer, namespace) {
    const args = [o.literal(name), expression];
    if (sanitizer !== null || namespace !== null) {
        args.push(sanitizer ?? o.literal(null));
    }
    if (namespace !== null) {
        args.push(o.literal(namespace));
    }
    return call(Identifiers.attribute, args, null);
}
export function styleProp(name, expression, unit, sourceSpan) {
    const args = [o.literal(name), expression];
    if (unit !== null) {
        args.push(o.literal(unit));
    }
    return call(Identifiers.styleProp, args, sourceSpan);
}
export function classProp(name, expression, sourceSpan) {
    return call(Identifiers.classProp, [o.literal(name), expression], sourceSpan);
}
export function styleMap(expression, sourceSpan) {
    return call(Identifiers.styleMap, [expression], sourceSpan);
}
export function classMap(expression, sourceSpan) {
    return call(Identifiers.classMap, [expression], sourceSpan);
}
const PIPE_BINDINGS = [
    Identifiers.pipeBind1,
    Identifiers.pipeBind2,
    Identifiers.pipeBind3,
    Identifiers.pipeBind4,
];
export function pipeBind(slot, varOffset, args) {
    if (args.length < 1 || args.length > PIPE_BINDINGS.length) {
        throw new Error(`pipeBind() argument count out of bounds`);
    }
    const instruction = PIPE_BINDINGS[args.length - 1];
    return o.importExpr(instruction).callFn([o.literal(slot), o.literal(varOffset), ...args]);
}
export function pipeBindV(slot, varOffset, args) {
    return o.importExpr(Identifiers.pipeBindV).callFn([o.literal(slot), o.literal(varOffset), args]);
}
export function textInterpolate(strings, expressions, sourceSpan) {
    const interpolationArgs = collateInterpolationArgs(strings, expressions);
    return callVariadicInstruction(TEXT_INTERPOLATE_CONFIG, [], interpolationArgs, [], sourceSpan);
}
export function i18nExp(expr, sourceSpan) {
    return call(Identifiers.i18nExp, [expr], sourceSpan);
}
export function i18nApply(slot, sourceSpan) {
    return call(Identifiers.i18nApply, [o.literal(slot)], sourceSpan);
}
export function propertyInterpolate(name, strings, expressions, sanitizer, sourceSpan) {
    const interpolationArgs = collateInterpolationArgs(strings, expressions);
    const extraArgs = [];
    if (sanitizer !== null) {
        extraArgs.push(sanitizer);
    }
    return callVariadicInstruction(PROPERTY_INTERPOLATE_CONFIG, [o.literal(name)], interpolationArgs, extraArgs, sourceSpan);
}
export function attributeInterpolate(name, strings, expressions, sanitizer, sourceSpan) {
    const interpolationArgs = collateInterpolationArgs(strings, expressions);
    const extraArgs = [];
    if (sanitizer !== null) {
        extraArgs.push(sanitizer);
    }
    return callVariadicInstruction(ATTRIBUTE_INTERPOLATE_CONFIG, [o.literal(name)], interpolationArgs, extraArgs, sourceSpan);
}
export function stylePropInterpolate(name, strings, expressions, unit, sourceSpan) {
    const interpolationArgs = collateInterpolationArgs(strings, expressions);
    const extraArgs = [];
    if (unit !== null) {
        extraArgs.push(o.literal(unit));
    }
    return callVariadicInstruction(STYLE_PROP_INTERPOLATE_CONFIG, [o.literal(name)], interpolationArgs, extraArgs, sourceSpan);
}
export function styleMapInterpolate(strings, expressions, sourceSpan) {
    const interpolationArgs = collateInterpolationArgs(strings, expressions);
    return callVariadicInstruction(STYLE_MAP_INTERPOLATE_CONFIG, [], interpolationArgs, [], sourceSpan);
}
export function classMapInterpolate(strings, expressions, sourceSpan) {
    const interpolationArgs = collateInterpolationArgs(strings, expressions);
    return callVariadicInstruction(CLASS_MAP_INTERPOLATE_CONFIG, [], interpolationArgs, [], sourceSpan);
}
export function hostProperty(name, expression, sanitizer, sourceSpan) {
    const args = [o.literal(name), expression];
    if (sanitizer !== null) {
        args.push(sanitizer);
    }
    return call(Identifiers.hostProperty, args, sourceSpan);
}
export function syntheticHostProperty(name, expression, sourceSpan) {
    return call(Identifiers.syntheticHostProperty, [o.literal(name), expression], sourceSpan);
}
export function pureFunction(varOffset, fn, args) {
    return callVariadicInstructionExpr(PURE_FUNCTION_CONFIG, [o.literal(varOffset), fn], args, [], null);
}
/**
 * Collates the string an expression arguments for an interpolation instruction.
 */
function collateInterpolationArgs(strings, expressions) {
    if (strings.length < 1 || expressions.length !== strings.length - 1) {
        throw new Error(`AssertionError: expected specific shape of args for strings/expressions in interpolation`);
    }
    const interpolationArgs = [];
    if (expressions.length === 1 && strings[0] === '' && strings[1] === '') {
        interpolationArgs.push(expressions[0]);
    }
    else {
        let idx;
        for (idx = 0; idx < expressions.length; idx++) {
            interpolationArgs.push(o.literal(strings[idx]), expressions[idx]);
        }
        // idx points at the last string.
        interpolationArgs.push(o.literal(strings[idx]));
    }
    return interpolationArgs;
}
function call(instruction, args, sourceSpan) {
    const expr = o.importExpr(instruction).callFn(args, sourceSpan);
    return ir.createStatementOp(new o.ExpressionStatement(expr, sourceSpan));
}
export function conditional(condition, contextValue, sourceSpan) {
    const args = [condition];
    if (contextValue !== null) {
        args.push(contextValue);
    }
    return call(Identifiers.conditional, args, sourceSpan);
}
/**
 * `InterpolationConfig` for the `textInterpolate` instruction.
 */
const TEXT_INTERPOLATE_CONFIG = {
    constant: [
        Identifiers.textInterpolate,
        Identifiers.textInterpolate1,
        Identifiers.textInterpolate2,
        Identifiers.textInterpolate3,
        Identifiers.textInterpolate4,
        Identifiers.textInterpolate5,
        Identifiers.textInterpolate6,
        Identifiers.textInterpolate7,
        Identifiers.textInterpolate8,
    ],
    variable: Identifiers.textInterpolateV,
    mapping: (n) => {
        if (n % 2 === 0) {
            throw new Error(`Expected odd number of arguments`);
        }
        return (n - 1) / 2;
    },
};
/**
 * `InterpolationConfig` for the `propertyInterpolate` instruction.
 */
const PROPERTY_INTERPOLATE_CONFIG = {
    constant: [
        Identifiers.propertyInterpolate,
        Identifiers.propertyInterpolate1,
        Identifiers.propertyInterpolate2,
        Identifiers.propertyInterpolate3,
        Identifiers.propertyInterpolate4,
        Identifiers.propertyInterpolate5,
        Identifiers.propertyInterpolate6,
        Identifiers.propertyInterpolate7,
        Identifiers.propertyInterpolate8,
    ],
    variable: Identifiers.propertyInterpolateV,
    mapping: (n) => {
        if (n % 2 === 0) {
            throw new Error(`Expected odd number of arguments`);
        }
        return (n - 1) / 2;
    },
};
/**
 * `InterpolationConfig` for the `stylePropInterpolate` instruction.
 */
const STYLE_PROP_INTERPOLATE_CONFIG = {
    constant: [
        Identifiers.styleProp,
        Identifiers.stylePropInterpolate1,
        Identifiers.stylePropInterpolate2,
        Identifiers.stylePropInterpolate3,
        Identifiers.stylePropInterpolate4,
        Identifiers.stylePropInterpolate5,
        Identifiers.stylePropInterpolate6,
        Identifiers.stylePropInterpolate7,
        Identifiers.stylePropInterpolate8,
    ],
    variable: Identifiers.stylePropInterpolateV,
    mapping: (n) => {
        if (n % 2 === 0) {
            throw new Error(`Expected odd number of arguments`);
        }
        return (n - 1) / 2;
    },
};
/**
 * `InterpolationConfig` for the `attributeInterpolate` instruction.
 */
const ATTRIBUTE_INTERPOLATE_CONFIG = {
    constant: [
        Identifiers.attribute,
        Identifiers.attributeInterpolate1,
        Identifiers.attributeInterpolate2,
        Identifiers.attributeInterpolate3,
        Identifiers.attributeInterpolate4,
        Identifiers.attributeInterpolate5,
        Identifiers.attributeInterpolate6,
        Identifiers.attributeInterpolate7,
        Identifiers.attributeInterpolate8,
    ],
    variable: Identifiers.attributeInterpolateV,
    mapping: (n) => {
        if (n % 2 === 0) {
            throw new Error(`Expected odd number of arguments`);
        }
        return (n - 1) / 2;
    },
};
/**
 * `InterpolationConfig` for the `styleMapInterpolate` instruction.
 */
const STYLE_MAP_INTERPOLATE_CONFIG = {
    constant: [
        Identifiers.styleMap,
        Identifiers.styleMapInterpolate1,
        Identifiers.styleMapInterpolate2,
        Identifiers.styleMapInterpolate3,
        Identifiers.styleMapInterpolate4,
        Identifiers.styleMapInterpolate5,
        Identifiers.styleMapInterpolate6,
        Identifiers.styleMapInterpolate7,
        Identifiers.styleMapInterpolate8,
    ],
    variable: Identifiers.styleMapInterpolateV,
    mapping: (n) => {
        if (n % 2 === 0) {
            throw new Error(`Expected odd number of arguments`);
        }
        return (n - 1) / 2;
    },
};
/**
 * `InterpolationConfig` for the `classMapInterpolate` instruction.
 */
const CLASS_MAP_INTERPOLATE_CONFIG = {
    constant: [
        Identifiers.classMap,
        Identifiers.classMapInterpolate1,
        Identifiers.classMapInterpolate2,
        Identifiers.classMapInterpolate3,
        Identifiers.classMapInterpolate4,
        Identifiers.classMapInterpolate5,
        Identifiers.classMapInterpolate6,
        Identifiers.classMapInterpolate7,
        Identifiers.classMapInterpolate8,
    ],
    variable: Identifiers.classMapInterpolateV,
    mapping: (n) => {
        if (n % 2 === 0) {
            throw new Error(`Expected odd number of arguments`);
        }
        return (n - 1) / 2;
    },
};
const PURE_FUNCTION_CONFIG = {
    constant: [
        Identifiers.pureFunction0,
        Identifiers.pureFunction1,
        Identifiers.pureFunction2,
        Identifiers.pureFunction3,
        Identifiers.pureFunction4,
        Identifiers.pureFunction5,
        Identifiers.pureFunction6,
        Identifiers.pureFunction7,
        Identifiers.pureFunction8,
    ],
    variable: Identifiers.pureFunctionV,
    mapping: (n) => n,
};
function callVariadicInstructionExpr(config, baseArgs, interpolationArgs, extraArgs, sourceSpan) {
    const n = config.mapping(interpolationArgs.length);
    if (n < config.constant.length) {
        // Constant calling pattern.
        return o
            .importExpr(config.constant[n])
            .callFn([...baseArgs, ...interpolationArgs, ...extraArgs], sourceSpan);
    }
    else if (config.variable !== null) {
        // Variable calling pattern.
        return o
            .importExpr(config.variable)
            .callFn([...baseArgs, o.literalArr(interpolationArgs), ...extraArgs], sourceSpan);
    }
    else {
        throw new Error(`AssertionError: unable to call variadic function`);
    }
}
function callVariadicInstruction(config, baseArgs, interpolationArgs, extraArgs, sourceSpan) {
    return ir.createStatementOp(callVariadicInstructionExpr(config, baseArgs, interpolationArgs, extraArgs, sourceSpan).toStmt());
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaW5zdHJ1Y3Rpb24uanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi9wYWNrYWdlcy9jb21waWxlci9zcmMvdGVtcGxhdGUvcGlwZWxpbmUvc3JjL2luc3RydWN0aW9uLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBOzs7Ozs7R0FNRztBQUVILE9BQU8sS0FBSyxDQUFDLE1BQU0sNEJBQTRCLENBQUM7QUFFaEQsT0FBTyxFQUFDLFdBQVcsRUFBQyxNQUFNLGlDQUFpQyxDQUFDO0FBQzVELE9BQU8sS0FBSyxFQUFFLE1BQU0sT0FBTyxDQUFDO0FBRTVCLDJGQUEyRjtBQUMzRixpR0FBaUc7QUFDakcsb0NBQW9DO0FBRXBDLE1BQU0sVUFBVSxPQUFPLENBQ3JCLElBQVksRUFDWixHQUFXLEVBQ1gsVUFBeUIsRUFDekIsYUFBNEIsRUFDNUIsVUFBMkI7SUFFM0IsT0FBTyxzQkFBc0IsQ0FDM0IsV0FBVyxDQUFDLE9BQU8sRUFDbkIsSUFBSSxFQUNKLEdBQUcsRUFDSCxVQUFVLEVBQ1YsYUFBYSxFQUNiLFVBQVUsQ0FDWCxDQUFDO0FBQ0osQ0FBQztBQUVELE1BQU0sVUFBVSxZQUFZLENBQzFCLElBQVksRUFDWixHQUFXLEVBQ1gsVUFBeUIsRUFDekIsYUFBNEIsRUFDNUIsVUFBMkI7SUFFM0IsT0FBTyxzQkFBc0IsQ0FDM0IsV0FBVyxDQUFDLFlBQVksRUFDeEIsSUFBSSxFQUNKLEdBQUcsRUFDSCxVQUFVLEVBQ1YsYUFBYSxFQUNiLFVBQVUsQ0FDWCxDQUFDO0FBQ0osQ0FBQztBQUVELFNBQVMsc0JBQXNCLENBQzdCLFdBQWdDLEVBQ2hDLElBQVksRUFDWixHQUFrQixFQUNsQixVQUF5QixFQUN6QixhQUE0QixFQUM1QixVQUEyQjtJQUUzQixNQUFNLElBQUksR0FBbUIsQ0FBQyxDQUFDLENBQUMsT0FBTyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUM7SUFDL0MsSUFBSSxHQUFHLEtBQUssSUFBSSxF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VSxDQUFDLENBQUM7SUFDM0UsQ0FBQztTQUFNLElBQUksTUFBTSxDQUFDLFFBQVEsS0FBSyxJQUFJLEVBQUUsQ0FBQztRQUNwQyw0QkFBNEI7UUFDNUIsT0FBTyxDQUFDO2FBQ0wsVUFBVSxDQUFDLE1BQU0sQ0FBQyxRQUFRLENBQUM7YUFDM0IsTUFBTSxDQUFDLENBQUMsR0FBRyxRQUFRLEVBQUUsQ0FBQyxDQUFDLFVBQVUsQ0FBQyxpQkFBaUIsQ0FBQyxFQUFFLEdBQUcsU0FBUyxDQUFDLEVBQUUsVUFBVSxDQUFDLENBQUM7SUFDdEYsQ0FBQztTQUFNLENBQUM7UUFDTixNQUFNLElBQUksS0FBSyxDQUFDLGtEQUFrRCxDQUFDLENBQUM7SUFDdEUsQ0FBQztBQUNILENBQUM7QUFFRCxTQUFTLHVCQUF1QixDQUM5QixNQUFpQyxFQUNqQyxRQUF3QixFQUN4QixpQkFBaUMsRUFDakMsU0FBeUIsRUFDekIsVUFBa0M7SUFFbEMsT0FBTyxFQUFFLENBQUMsaUJBQWlCLENBQ3pCLDJCQUEyQixDQUN6QixNQUFNLEVBQ04sUUFBUSxFQUNSLGlCQUFpQixFQUNqQixTQUFTLEVBQ1QsVUFBVSxDQUNYLENBQUMsTUFBTSxFQUFFLENBQ1gsQ0FBQztBQUNKLENBQUMiLCJzb3VyY2VzQ29udGVudCI6WyIvKipcbiAqIEBsaWNlbnNlXG4gKiBDb3B5cmlnaHQgR29vZ2xlIExMQyBBbGwgUmlnaHRzIFJlc2VydmVkLlxuICpcbiAqIFVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlIGlzIGdvdmVybmVkIGJ5IGFuIE1JVC1zdHlsZSBsaWNlbnNlIHRoYXQgY2FuIGJlXG4gKiBmb3VuZCBpbiB0aGUgTElDRU5TRSBmaWxlIGF0IGh0dHBzOi8vYW5ndWxhci5pby9saWNlbnNlXG4gKi9cblxuaW1wb3J0ICogYXMgbyBmcm9tICcuLi8uLi8uLi9vdXRwdXQvb3V0cHV0X2FzdCc7XG5pbXBvcnQge1BhcnNlU291cmNlU3Bhbn0gZnJvbSAnLi4vLi4vLi4vcGFyc2VfdXRpbCc7XG5pbXBvcnQge0lkZW50aWZpZXJzfSBmcm9tICcuLi8uLi8uLi9yZW5kZXIzL3IzX2lkZW50aWZpZXJzJztcbmltcG9ydCAqIGFzIGlyIGZyb20gJy4uL2lyJztcblxuLy8gVGhpcyBmaWxlIGNvbnRhaW5zIGhlbHBlcnMgZm9yIGdlbmVyYXRpbmcgY2FsbHMgdG8gSXZ5IGluc3RydWN0aW9ucy4gSW4gcGFydGljdWxhciwgZWFjaFxuLy8gaW5zdHJ1Y3Rpb24gdHlwZSBpcyByZXByZXNlbnRlZCBhcyBhIGZ1bmN0aW9uLCB3aGljaCBtYXkgc2VsZWN0IGEgc3BlY2lmaWMgaW5zdHJ1Y3Rpb24gdmFyaWFudFxuLy8gZGVwZW5kaW5nIG9uIHRoZSBleGFjdCBhcmd1bWVudHMuXG5cbmV4cG9ydCBmdW5jdGlvbiBlbGVtZW50KFxuICBzbG90OiBudW1iZXIsXG4gIHRhZzogc3RyaW5nLFxuICBjb25zdEluZGV4OiBudW1iZXIgfCBudWxsLFxuICBsb2NhbFJlZkluZGV4OiBudW1iZXIgfCBudWxsLFxuICBzb3VyY2VTcGFuOiBQYXJzZVNvdXJjZVNwYW4sXG4pOiBpci5DcmVhdGVPcCB7XG4gIHJldHVybiBlbGVtZW50T3JDb250YWluZXJCYXNlKFxuICAgIElkZW50aWZpZXJzLmVsZW1lbnQsXG4gICAgc2xvdCxcbiAgICB0YWcsXG4gICAgY29uc3RJbmRleCxcbiAgICBsb2NhbFJlZkluZGV4LFxuICAgIHNvdXJjZVNwYW4sXG4gICk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBlbGVtZW50U3RhcnQoXG4gIHNsb3Q6IG51bWJlcixcbiAgdGFnOiBzdHJpbmcsXG4gIGNvbnN0SW5kZXg6IG51bWJlciB8IG51bGwsXG4gIGxvY2FsUmVmSW5kZXg6IG51bWJlciB8IG51bGwsXG4gIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3Bhbixcbik6IGlyLkNyZWF0ZU9wIHtcbiAgcmV0dXJuIGVsZW1lbnRPckNvbnRhaW5lckJhc2UoXG4gICAgSWRlbnRpZmllcnMuZWxlbWVudFN0YXJ0LFxuICAgIHNsb3QsXG4gICAgdGFnLFxuICAgIGNvbnN0SW5kZXgsXG4gICAgbG9jYWxSZWZJbmRleCxcbiAgICBzb3VyY2VTcGFuLFxuICApO1xufVxuXG5mdW5jdGlvbiBlbGVtZW50T3JDb250YWluZXJCYXNlKFxuICBpbnN0cnVjdGlvbjogby5FeHRlcm5hbFJlZmVyZW5jZSxcbiAgc2xvdDogbnVtYmVyLFxuICB0YWc6IHN0cmluZyB8IG51bGwsXG4gIGNvbnN0SW5kZXg6IG51bWJlciB8IG51bGwsXG4gIGxvY2FsUmVmSW5kZXg6IG51bWJlciB8IG51bGwsXG4gIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3Bhbixcbik6IGlyLkNyZWF0ZU9wIHtcbiAgY29uc3QgYXJnczogby5FeHByZXNzaW9uW10gPSBbby5saXRlcmFsKHNsb3QpXTtcbiAgaWYgKHRhZyAhPT0gbnVsbCkge1xuICAgIGFyZ3MucHVzaChvLmxpdGVyYWwodGFnKSk7XG4gIH1cbiAgaWYgKGxvY2FsUmVmSW5kZXggIT09IG51bGwpIHtcbiAgICBhcmdzLnB1c2goXG4gICAgICBvLmxpdGVyYWwoY29uc3RJbmRleCksIC8vIG1pZ2h0IGJlIG51bGwsIGJ1dCB0aGF0J3Mgb2theS5cbiAgICAgIG8ubGl0ZXJhbChsb2NhbFJlZkluZGV4KSxcbiAgICApO1xuICB9IGVsc2UgaWYgKGNvbnN0SW5kZXggIT09IG51bGwpIHtcbiAgICBhcmdzLnB1c2goby5saXRlcmFsKGNvbnN0SW5kZXgpKTtcbiAgfVxuXG4gIHJldHVybiBjYWxsKGluc3RydWN0aW9uLCBhcmdzLCBzb3VyY2VTcGFuKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGVsZW1lbnRFbmQoc291cmNlU3BhbjogUGFyc2VTb3VyY2VTcGFuIHwgbnVsbCk6IGlyLkNyZWF0ZU9wIHtcbiAgcmV0dXJuIGNhbGwoSWRlbnRpZmllcnMuZWxlbWVudEVuZCwgW10sIHNvdXJjZVNwYW4pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gZWxlbWVudENvbnRhaW5lclN0YXJ0KFxuICBzbG90OiBudW1iZXIsXG4gIGNvbnN0SW5kZXg6IG51bWJlciB8IG51bGwsXG4gIGxvY2FsUmVmSW5kZXg6IG51bWJlciB8IG51bGwsXG4gIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3Bhbixcbik6IGlyLkNyZWF0ZU9wIHtcbiAgcmV0dXJuIGVsZW1lbnRPckNvbnRhaW5lckJhc2UoXG4gICAgSWRlbnRpZmllcnMuZWxlbWVudENvbnRhaW5lclN0YXJ0LFxuICAgIHNsb3QsXG4gICAgLyogdGFnICovIG51bGwsXG4gICAgY29uc3RJbmRleCxcbiAgICBsb2NhbFJlZkluZGV4LFxuICAgIHNvdXJjZVNwYW4sXG4gICk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBlbGVtZW50Q29udGFpbmVyKFxuICBzbG90OiBudW1iZXIsXG4gIGNvbnN0SW5kZXg6IG51bWJlciB8IG51bGwsXG4gIGxvY2FsUmVmSW5kZXg6IG51bWJlciB8IG51bGwsXG4gIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3Bhbixcbik6IGlyLkNyZWF0ZU9wIHtcbiAgcmV0dXJuIGVsZW1lbnRPckNvbnRhaW5lckJhc2UoXG4gICAgSWRlbnRpZmllcnMuZWxlbWVudENvbnRhaW5lcixcbiAgICBzbG90LFxuICAgIC8qIHRhZyAqLyBudWxsLFxuICAgIGNvbnN0SW5kZXgsXG4gICAgbG9jYWxSZWZJbmRleCxcbiAgICBzb3VyY2VTcGFuLFxuICApO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gZWxlbWVudENvbnRhaW5lckVuZCgpOiBpci5DcmVhdGVPcCB7XG4gIHJldHVybiBjYWxsKElkZW50aWZpZXJzLmVsZW1lbnRDb250YWluZXJFbmQsIFtdLCBudWxsKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIHRlbXBsYXRlKFxuICBzbG90OiBudW1iZXIsXG4gIHRlbXBsYXRlRm5SZWY6IG8uRXhwcmVzc2lvbixcbiAgZGVjbHM6IG51bWJlcixcbiAgdmFyczogbnVtYmVyLFxuICB0YWc6IHN0cmluZyB8IG51bGwsXG4gIGNvbnN0SW5kZXg6IG51bWJlciB8IG51bGwsXG4gIGxvY2FsUmVmczogbnVtYmVyIHwgbnVsbCxcbiAgc291cmNlU3BhbjogUGFyc2VTb3VyY2VTcGFuLFxuKTogaXIuQ3JlYXRlT3Age1xuICBjb25zdCBhcmdzID0gW1xuICAgIG8ubGl0ZXJhbChzbG90KSxcbiAgICB0ZW1wbGF0ZUZuUmVmLFxuICAgIG8ubGl0ZXJhbChkZWNscyksXG4gICAgby5saXRlcmFsKHZhcnMpLFxuICAgIG8ubGl0ZXJhbCh0YWcpLFxuICAgIG8ubGl0ZXJhbChjb25zdEluZGV4KSxcbiAgXTtcbiAgaWYgKGxvY2FsUmVmcyAhPT0gbnVsbCkge1xuICAgIGFyZ3MucHVzaChvLmxpdGVyYWwobG9jYWxSZWZzKSk7XG4gICAgYXJncy5wdXNoKG8uaW1wb3J0RXhwcihJZGVudGlmaWVycy50ZW1wbGF0ZVJlZkV4dHJhY3RvcikpO1xuICB9XG4gIHdoaWxlIChhcmdzW2FyZ3MubGVuZ3RoIC0gMV0uaXNFcXVpdmFsZW50KG8uTlVMTF9FWFBSKSkge1xuICAgIGFyZ3MucG9wKCk7XG4gIH1cbiAgcmV0dXJuIGNhbGwoSWRlbnRpZmllcnMudGVtcGxhdGVDcmVhdGUsIGFyZ3MsIHNvdXJjZVNwYW4pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gZGlzYWJsZUJpbmRpbmdzKCk6IGlyLkNyZWF0ZU9wIHtcbiAgcmV0dXJuIGNhbGwoSWRlbnRpZmllcnMuZGlzYWJsZUJpbmRpbmdzLCBbXSwgbnVsbCk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBlbmFibGVCaW5kaW5ncygpOiBpci5DcmVhdGVPcCB7XG4gIHJldHVybiBjYWxsKElkZW50aWZpZXJzLmVuYWJsZUJpbmRpbmdzLCBbXSwgbnVsbCk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBsaXN0ZW5lcihcbiAgbmFtZTogc3RyaW5nLFxuICBoYW5kbGVyRm46IG8uRXhwcmVzc2lvbixcbiAgZXZlbnRUYXJnZXRSZXNvbHZlcjogby5FeHRlcm5hbFJlZmVyZW5jZSB8IG51bGwsXG4gIHN5bnRoZXRpY0hvc3Q6IGJvb2xlYW4sXG4gIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3Bhbixcbik6IGlyLkNyZWF0ZU9wIHtcbiAgY29uc3QgYXJncyA9IFtvLmxpdGVyYWwobmFtZSksIGhhbmRsZXJGbl07XG4gIGlmIChldmVudFRhcmdldFJlc29sdmVyICE9PSBudWxsKSB7XG4gICAgYXJncy5wdXNoKG8ubGl0ZXJhbChmYWxzZSkpOyAvLyBgdXNlQ2FwdHVyZWAgZmxhZywgZGVmYXVsdHMgdG8gYGZhbHNlYFxuICAgIGFyZ3MucHVzaChvLmltcG9ydEV4cHIoZXZlbnRUYXJnZXRSZXNvbHZlcikpO1xuICB9XG4gIHJldHVybiBjYWxsKFxuICAgIHN5bnRoZXRpY0hvc3QgPyBJZGVudGlmaWVycy5zeW50aGV0aWNIb3N0TGlzdGVuZXIgOiBJZGVudGlmaWVycy5saXN0ZW5lcixcbiAgICBhcmdzLFxuICAgIHNvdXJjZVNwYW4sXG4gICk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiB0d29XYXlCaW5kaW5nU2V0KHRhcmdldDogby5FeHByZXNzaW9uLCB2YWx1ZTogby5FeHByZXNzaW9uKTogby5FeHByZXNzaW9uIHtcbiAgcmV0dXJuIG8uaW1wb3J0RXhwcihJZGVudGlmaWVycy50d29XYXlCaW5kaW5nU2V0KS5jYWxsRm4oW3RhcmdldCwgdmFsdWVdKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIHR3b1dheUxpc3RlbmVyKFxuICBuYW1lOiBzdHJpbmcsXG4gIGhhbmRsZXJGbjogby5FeHByZXNzaW9uLFxuICBzb3VyY2VTcGFuOiBQYXJzZVNvdXJjZVNwYW4sXG4pOiBpci5DcmVhdGVPcCB7XG4gIHJldHVybiBjYWxsKElkZW50aWZpZXJzLnR3b1dheUxpc3RlbmVyLCBbby5saXRlcmFsKG5hbWUpLCBoYW5kbGVyRm5dLCBzb3VyY2VTcGFuKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIHBpcGUoc2xvdDogbnVtYmVyLCBuYW1lOiBzdHJpbmcpOiBpci5DcmVhdGVPcCB7XG4gIHJldHVybiBjYWxsKElkZW50aWZpZXJzLnBpcGUsIFtvLmxpdGVyYWwoc2xvdCksIG8ubGl0ZXJhbChuYW1lKV0sIG51bGwpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gbmFtZXNwYWNlSFRNTCgpOiBpci5DcmVhdGVPcCB7XG4gIHJldHVybiBjYWxsKElkZW50aWZpZXJzLm5hbWVzcGFjZUhUTUwsIFtdLCBudWxsKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIG5hbWVzcGFjZVNWRygpOiBpci5DcmVhdGVPcCB7XG4gIHJldHVybiBjYWxsKElkZW50aWZpZXJzLm5hbWVzcGFjZVNWRywgW10sIG51bGwpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gbmFtZXNwYWNlTWF0aCgpOiBpci5DcmVhdGVPcCB7XG4gIHJldHVybiBjYWxsKElkZW50aWZpZXJzLm5hbWVzcGFjZU1hdGhNTCwgW10sIG51bGwpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gYWR2YW5jZShkZWx0YTogbnVtYmVyLCBzb3VyY2VTcGFuOiBQYXJzZVNvdXJjZVNwYW4pOiBpci5VcGRhdGVPcCB7XG4gIHJldHVybiBjYWxsKElkZW50aWZpZXJzLmFkdmFuY2UsIGRlbHRhID4gMSA/IFtvLmxpdGVyYWwoZGVsdGEpXSA6IFtdLCBzb3VyY2VTcGFuKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIHJlZmVyZW5jZShzbG90OiBudW1iZXIpOiBvLkV4cHJlc3Npb24ge1xuICByZXR1cm4gby5pbXBvcnRFeHByKElkZW50aWZpZXJzLnJlZmVyZW5jZSkuY2FsbEZuKFtvLmxpdGVyYWwoc2xvdCldKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIG5leHRDb250ZXh0KHN0ZXBzOiBudW1iZXIpOiBvLkV4cHJlc3Npb24ge1xuICByZXR1cm4gby5pbXBvcnRFeHByKElkZW50aWZpZXJzLm5leHRDb250ZXh0KS5jYWxsRm4oc3RlcHMgPT09IDEgPyBbXSA6IFtvLmxpdGVyYWwoc3RlcHMpXSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBnZXRDdXJyZW50VmlldygpOiBvLkV4cHJlc3Npb24ge1xuICByZXR1cm4gby5pbXBvcnRFeHByKElkZW50aWZpZXJzLmdldEN1cnJlbnRWaWV3KS5jYWxsRm4oW10pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gcmVzdG9yZVZpZXcoc2F2ZWRWaWV3OiBvLkV4cHJlc3Npb24pOiBvLkV4cHJlc3Npb24ge1xuICByZXR1cm4gby5pbXBvcnRFeHByKElkZW50aWZpZXJzLnJlc3RvcmVWaWV3KS5jYWxsRm4oW3NhdmVkVmlld10pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gcmVzZXRWaWV3KHJldHVyblZhbHVlOiBvLkV4cHJlc3Npb24pOiBvLkV4cHJlc3Npb24ge1xuICByZXR1cm4gby5pbXBvcnRFeHByKElkZW50aWZpZXJzLnJlc2V0VmlldykuY2FsbEZuKFtyZXR1cm5WYWx1ZV0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gdGV4dChcbiAgc2xvdDogbnVtYmVyLFxuICBpbml0aWFsVmFsdWU6IHN0cmluZyxcbiAgc291cmNlU3BhbjogUGFyc2VTb3VyY2VTcGFuIHwgbnVsbCxcbik6IGlyLkNyZWF0ZU9wIHtcbiAgY29uc3QgYXJnczogby5FeHByZXNzaW9uW10gPSBbby5saXRlcmFsKHNsb3QsIG51bGwpXTtcbiAgaWYgKGluaXRpYWxWYWx1ZSAhPT0gJycpIHtcbiAgICBhcmdzLnB1c2goby5saXRlcmFsKGluaXRpYWxWYWx1ZSkpO1xuICB9XG4gIHJldHVybiBjYWxsKElkZW50aWZpZXJzLnRleHQsIGFyZ3MsIHNvdXJjZVNwYW4pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gZGVmZXIoXG4gIHNlbGZTbG90OiBudW1iZXIsXG4gIHByaW1hcnlTbG90OiBudW1iZXIsXG4gIGRlcGVuZGVuY3lSZXNvbHZlckZuOiBvLkV4cHJlc3Npb24gfCBudWxsLFxuICBsb2FkaW5nU2xvdDogbnVtYmVyIHwgbnVsbCxcbiAgcGxhY2Vob2xkZXJTbG90OiBudW1iZXIgfCBudWxsLFxuICBlcnJvclNsb3Q6IG51bWJlciB8IG51bGwsXG4gIGxvYWRpbmdDb25maWc6IG8uRXhwcmVzc2lvbiB8IG51bGwsXG4gIHBsYWNlaG9sZGVyQ29uZmlnOiBvLkV4cHJlc3Npb24gfCBudWxsLFxuICBlbmFibGVUaW1lclNjaGVkdWxpbmc6IGJvb2xlYW4sXG4gIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3BhbiB8IG51bGwsXG4pOiBpci5DcmVhdGVPcCB7XG4gIGNvbnN0IGFyZ3M6IEFycmF5PG8uRXhwcmVzc2lvbj4gPSBbXG4gICAgby5saXRlcmFsKHNlbGZTbG90KSxcbiAgICBvLmxpdGVyYWwocHJpbWFyeVNsb3QpLFxuICAgIGRlcGVuZGVuY3lSZXNvbHZlckZuID8/IG8ubGl0ZXJhbChudWxsKSxcbiAgICBvLmxpdGVyYWwobG9hZGluZ1Nsb3QpLFxuICAgIG8ubGl0ZXJhbChwbGFjZWhvbGRlclNsb3QpLFxuICAgIG8ubGl0ZXJhbChlcnJvclNsb3QpLFxuICAgIGxvYWRpbmdDb25maWcgPz8gby5saXRlcmFsKG51bGwpLFxuICAgIHBsYWNlaG9sZGVyQ29uZmlnID8/IG8ubGl0ZXJhbChudWxsKSxcbiAgICBlbmFibGVUaW1lclNjaGVkdWxpbmcgPyBvLmltcG9ydEV4cHIoSWRlbnRpZmllcnMuZGVmZXJFbmFibGVUaW1lclNjaGVkdWxpbmcpIDogby5saXRlcmFsKG51bGwpLFxuICBdO1xuXG4gIGxldCBleHByOiBvLkV4cHJlc3Npb247XG4gIHdoaWxlIChcbiAgICAoZXhwciA9IGFyZ3NbYXJncy5sZW5ndGggLSAxXSkgIT09IG51bGwgJiZcbiAgICBleHByIGluc3RhbmNlb2Ygby5MaXRlcmFsRXhwciAmJlxuICAgIGV4cHIudmFsdWUgPT09IG51bGxcbiAgKSB7XG4gICAgYXJncy5wb3AoKTtcbiAgfVxuXG4gIHJldHVybiBjYWxsKElkZW50aWZpZXJzLmRlZmVyLCBhcmdzLCBzb3VyY2VTcGFuKTtcbn1cblxuY29uc3QgZGVmZXJUcmlnZ2VyVG9SM1RyaWdnZXJJbnN0cnVjdGlvbnNNYXAgPSBuZXcgTWFwKFtcbiAgW2lyLkRlZmVyVHJpZ2dlcktpbmQuSWRsZSwgW0lkZW50aWZpZXJzLmRlZmVyT25JZGxlLCBJZGVudGlmaWVycy5kZWZlclByZWZldGNoT25JZGxlXV0sXG4gIFtcbiAgICBpci5EZWZlclRyaWdnZXJLaW5kLkltbWVkaWF0ZSxcbiAgICBbSWRlbnRpZmllcnMuZGVmZXJPbkltbWVkaWF0ZSwgSWRlbnRpZmllcnMuZGVmZXJQcmVmZXRjaE9uSW1tZWRpYXRlXSxcbiAgXSxcbiAgW2lyLkRlZmVyVHJpZ2dlcktpbmQuVGltZXIsIFtJZGVudGlmaWVycy5kZWZlck9uVGltZXIsIElkZW50aWZpZXJzLmRlZmVyUHJlZmV0Y2hPblRpbWVyXV0sXG4gIFtpci5EZWZlclRyaWdnZXJLaW5kLkhvdmVyLCBbSWRlbnRpZmllcnMuZGVmZXJPbkhvdmVyLCBJZGVudGlmaWVycy5kZWZlclByZWZldGNoT25Ib3Zlcl1dLFxuICBbXG4gICAgaXIuRGVmZXJUcmlnZ2VyS2luZC5JbnRlcmFjdGlvbixcbiAgICBbSWRlbnRpZmllcnMuZGVmZXJPbkludGVyYWN0aW9uLCBJZGVudGlmaWVycy5kZWZlclByZWZldGNoT25JbnRlcmFjdGlvbl0sXG4gIF0sXG4gIFtcbiAgICBpci5EZWZlclRyaWdnZXJLaW5kLlZpZXdwb3J0LFxuICAgIFtJZGVudGlmaWVycy5kZWZlck9uVmlld3BvcnQsIElkZW50aWZpZXJzLmRlZmVyUHJlZmV0Y2hPblZpZXdwb3J0XSxcbiAgXSxcbl0pO1xuXG5leHBvcnQgZnVuY3Rpb24gZGVmZXJPbihcbiAgdHJpZ2dlcjogaXIuRGVmZXJUcmlnZ2VyS2luZCxcbiAgYXJnczogbnVtYmVyW10sXG4gIHByZWZldGNoOiBib29sZWFuLFxuICBzb3VyY2VTcGFuOiBQYXJzZVNvdXJjZVNwYW4gfCBudWxsLFxuKTogaXIuQ3JlYXRlT3Age1xuICBjb25zdCBpbnN0cnVjdGlvbnMgPSBkZWZlclRyaWdnZXJUb1IzVHJpZ2dlckluc3RydWN0aW9uc01hcC5nZXQodHJpZ2dlcik7XG4gIGlmIChpbnN0cnVjdGlvbnMgPT09IHVuZGVmaW5lZCkge1xuICAgIHRocm93IG5ldyBFcnJvcihgVW5hYmxlIHRvIGRldGVybWluZSBpbnN0cnVjdGlvbiBmb3IgdHJpZ2dlciAke3RyaWdnZXJ9YCk7XG4gIH1cbiAgY29uc3QgaW5zdHJ1Y3Rpb25Ub0NhbGwgPSBwcmVmZXRjaCA/IGluc3RydWN0aW9uc1sxXSA6IGluc3RydWN0aW9uc1swXTtcbiAgcmV0dXJuIGNhbGwoXG4gICAgaW5zdHJ1Y3Rpb25Ub0NhbGwsXG4gICAgYXJncy5tYXAoKGEpID0+IG8ubGl0ZXJhbChhKSksXG4gICAgc291cmNlU3BhbixcbiAgKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIHByb2plY3Rpb25EZWYoZGVmOiBvLkV4cHJlc3Npb24gfCBudWxsKTogaXIuQ3JlYXRlT3Age1xuICByZXR1cm4gY2FsbChJZGVudGlmaWVycy5wcm9qZWN0aW9uRGVmLCBkZWYgPyBbZGVmXSA6IFtdLCBudWxsKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIHByb2plY3Rpb24oXG4gIHNsb3Q6IG51bWJlcixcbiAgcHJvamVjdGlvblNsb3RJbmRleDogbnVtYmVyLFxuICBhdHRyaWJ1dGVzOiBvLkxpdGVyYWxBcnJheUV4cHIgfCBudWxsLFxuICBmYWxsYmFja0ZuTmFtZTogc3RyaW5nIHwgbnVsbCxcbiAgZmFsbGJhY2tEZWNsczogbnVtYmVyIHwgbnVsbCxcbiAgZmFsbGJhY2tWYXJzOiBudW1iZXIgfCBudWxsLFxuICBzb3VyY2VTcGFuOiBQYXJzZVNvdXJjZVNwYW4sXG4pOiBpci5DcmVhdGVPcCB7XG4gIGNvbnN0IGFyZ3M6IG8uRXhwcmVzc2lvbltdID0gW28ubGl0ZXJhbChzbG90KV07XG4gIGlmIChwcm9qZWN0aW9uU2xvdEluZGV4ICE9PSAwIHx8IGF0dHJpYnV0ZXMgIT09IG51bGwgfHwgZmFsbGJhY2tGbk5hbWUgIT09IG51bGwpIHtcbiAgICBhcmdzLnB1c2goby5saXRlcmFsKHByb2plY3Rpb25TbG90SW5kZXgpKTtcbiAgICBpZiAoYXR0cmlidXRlcyAhPT0gbnVsbCkge1xuICAgICAgYXJncy5wdXNoKGF0dHJpYnV0ZXMpO1xuICAgIH1cbiAgICBpZiAoZmFsbGJhY2tGbk5hbWUgIT09IG51bGwpIHtcbiAgICAgIGlmIChhdHRyaWJ1dGVzID09PSBudWxsKSB7XG4gICAgICAgIGFyZ3MucHVzaChvLmxpdGVyYWwobnVsbCkpO1xuICAgICAgfVxuICAgICAgYXJncy5wdXNoKG8udmFyaWFibGUoZmFsbGJhY2tGbk5hbWUpLCBvLmxpdGVyYWwoZmFsbGJhY2tEZWNscyksIG8ubGl0ZXJhbChmYWxsYmFja1ZhcnMpKTtcbiAgICB9XG4gIH1cbiAgcmV0dXJuIGNhbGwoSWRlbnRpZmllcnMucHJvamVjdGlvbiwgYXJncywgc291cmNlU3Bhbik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBpMThuU3RhcnQoXG4gIHNsb3Q6IG51bWJlcixcbiAgY29uc3RJbmRleDogbnVtYmVyLFxuICBzdWJUZW1wbGF0ZUluZGV4OiBudW1iZXIsXG4gIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3BhbiB8IG51bGwsXG4pOiBpci5DcmVhdGVPcCB7XG4gIGNvbnN0IGFyZ3MgPSBbby5saXRlcmFsKHNsb3QpLCBvLmxpdGVyYWwoY29uc3RJbmRleCldO1xuICBpZiAoc3ViVGVtcGxhdGVJbmRleCAhPT0gbnVsbCkge1xuICAgIGFyZ3MucHVzaChvLmxpdGVyYWwoc3ViVGVtcGxhdGVJbmRleCkpO1xuICB9XG4gIHJldHVybiBjYWxsKElkZW50aWZpZXJzLmkxOG5TdGFydCwgYXJncywgc291cmNlU3Bhbik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiByZXBlYXRlckNyZWF0ZShcbiAgc2xvdDogbnVtYmVyLFxuICB2aWV3Rm5OYW1lOiBzdHJpbmcsXG4gIGRlY2xzOiBudW1iZXIsXG4gIHZhcnM6IG51bWJlcixcbiAgdGFnOiBzdHJpbmcgfCBudWxsLFxuICBjb25zdEluZGV4OiBudW1iZXIgfCBudWxsLFxuICB0cmFja0J5Rm46IG8uRXhwcmVzc2lvbixcbiAgdHJhY2tCeVVzZXNDb21wb25lbnRJbnN0YW5jZTogYm9vbGVhbixcbiAgZW1wdHlWaWV3Rm5OYW1lOiBzdHJpbmcgfCBudWxsLFxuICBlbXB0eURlY2xzOiBudW1iZXIgfCBudWxsLFxuICBlbXB0eVZhcnM6IG51bWJlciB8IG51bGwsXG4gIGVtcHR5VGFnOiBzdHJpbmcgfCBudWxsLFxuICBlbXB0eUNvbnN0SW5kZXg6IG51bWJlciB8IG51bGwsXG4gIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3BhbiB8IG51bGwsXG4pOiBpci5DcmVhdGVPcCB7XG4gIGNvbnN0IGFyZ3MgPSBbXG4gICAgby5saXRlcmFsKHNsb3QpLFxuICAgIG8udmFyaWFibGUodmlld0ZuTmFtZSksXG4gICAgby5saXRlcmFsKGRlY2xzKSxcbiAgICBvLmxpdGVyYWwodmFycyksXG4gICAgby5saXRlcmFsKHRhZyksXG4gICAgby5saXRlcmFsKGNvbnN0SW5kZXgpLFxuICAgIHRyYWNrQnlGbixcbiAgXTtcbiAgaWYgKHRyYWNrQnlVc2VzQ29tcG9uZW50SW5zdGFuY2UgfHwgZW1wdHlWaWV3Rm5OYW1lICE9PSBudWxsKSB7XG4gICAgYXJncy5wdXNoKG8ubGl0ZXJhbCh0cmFja0J5VXNlc0NvbXBvbmVudEluc3RhbmNlKSk7XG4gICAgaWYgKGVtcHR5Vmlld0ZuTmFtZSAhPT0gbnVsbCkge1xuICAgICAgYXJncy5wdXNoKG8udmFyaWFibGUoZW1wdHlWaWV3Rm5OYW1lKSwgby5saXRlcmFsKGVtcHR5RGVjbHMpLCBvLmxpdGVyYWwoZW1wdHlWYXJzKSk7XG4gICAgICBpZiAoZW1wdHlUYWcgIT09IG51bGwgfHwgZW1wdHlDb25zdEluZGV4ICE9PSBudWxsKSB7XG4gICAgICAgIGFyZ3MucHVzaChvLmxpdGVyYWwoZW1wdHlUYWcpKTtcbiAgICAgIH1cbiAgICAgIGlmIChlbXB0eUNvbnN0SW5kZXggIT09IG51bGwpIHtcbiAgICAgICAgYXJncy5wdXNoKG8ubGl0ZXJhbChlbXB0eUNvbnN0SW5kZXgpKTtcbiAgICAgIH1cbiAgICB9XG4gIH1cbiAgcmV0dXJuIGNhbGwoSWRlbnRpZmllcnMucmVwZWF0ZXJDcmVhdGUsIGFyZ3MsIHNvdXJjZVNwYW4pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gcmVwZWF0ZXIoXG4gIGNvbGxlY3Rpb246IG8uRXhwcmVzc2lvbixcbiAgc291cmNlU3BhbjogUGFyc2VTb3VyY2VTcGFuIHwgbnVsbCxcbik6IGlyLlVwZGF0ZU9wIHtcbiAgcmV0dXJuIGNhbGwoSWRlbnRpZmllcnMucmVwZWF0ZXIsIFtjb2xsZWN0aW9uXSwgc291cmNlU3Bhbik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBkZWZlcldoZW4oXG4gIHByZWZldGNoOiBib29sZWFuLFxuICBleHByOiBvLkV4cHJlc3Npb24sXG4gIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3BhbiB8IG51bGwsXG4pOiBpci5VcGRhdGVPcCB7XG4gIHJldHVybiBjYWxsKHByZWZldGNoID8gSWRlbnRpZmllcnMuZGVmZXJQcmVmZXRjaFdoZW4gOiBJZGVudGlmaWVycy5kZWZlcldoZW4sIFtleHByXSwgc291cmNlU3Bhbik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBkZWNsYXJlTGV0KHNsb3Q6IG51bWJlciwgc291cmNlU3BhbjogUGFyc2VTb3VyY2VTcGFuKTogaXIuQ3JlYXRlT3Age1xuICByZXR1cm4gY2FsbChJZGVudGlmaWVycy5kZWNsYXJlTGV0LCBbby5saXRlcmFsKHNsb3QpXSwgc291cmNlU3Bhbik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBzdG9yZUxldCh2YWx1ZTogby5FeHByZXNzaW9uLCBzb3VyY2VTcGFuOiBQYXJzZVNvdXJjZVNwYW4pOiBvLkV4cHJlc3Npb24ge1xuICByZXR1cm4gby5pbXBvcnRFeHByKElkZW50aWZpZXJzLnN0b3JlTGV0KS5jYWxsRm4oW3ZhbHVlXSwgc291cmNlU3Bhbik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiByZWFkQ29udGV4dExldChzbG90OiBudW1iZXIpOiBvLkV4cHJlc3Npb24ge1xuICByZXR1cm4gby5pbXBvcnRFeHByKElkZW50aWZpZXJzLnJlYWRDb250ZXh0TGV0KS5jYWxsRm4oW28ubGl0ZXJhbChzbG90KV0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gaTE4bihcbiAgc2xvdDogbnVtYmVyLFxuICBjb25zdEluZGV4OiBudW1iZXIsXG4gIHN1YlRlbXBsYXRlSW5kZXg6IG51bWJlcixcbiAgc291cmNlU3BhbjogUGFyc2VTb3VyY2VTcGFuIHwgbnVsbCxcbik6IGlyLkNyZWF0ZU9wIHtcbiAgY29uc3QgYXJncyA9IFtvLmxpdGVyYWwoc2xvdCksIG8ubGl0ZXJhbChjb25zdEluZGV4KV07XG4gIGlmIChzdWJUZW1wbGF0ZUluZGV4KSB7XG4gICAgYXJncy5wdXNoKG8ubGl0ZXJhbChzdWJUZW1wbGF0ZUluZGV4KSk7XG4gIH1cbiAgcmV0dXJuIGNhbGwoSWRlbnRpZmllcnMuaTE4biwgYXJncywgc291cmNlU3Bhbik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBpMThuRW5kKGVuZFNvdXJjZVNwYW46IFBhcnNlU291cmNlU3BhbiB8IG51bGwpOiBpci5DcmVhdGVPcCB7XG4gIHJldHVybiBjYWxsKElkZW50aWZpZXJzLmkxOG5FbmQsIFtdLCBlbmRTb3VyY2VTcGFuKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGkxOG5BdHRyaWJ1dGVzKHNsb3Q6IG51bWJlciwgaTE4bkF0dHJpYnV0ZXNDb25maWc6IG51bWJlcik6IGlyLkNyZWF0ZU9wIHtcbiAgY29uc3QgYXJncyA9IFtvLmxpdGVyYWwoc2xvdCksIG8ubGl0ZXJhbChpMThuQXR0cmlidXRlc0NvbmZpZyldO1xuICByZXR1cm4gY2FsbChJZGVudGlmaWVycy5pMThuQXR0cmlidXRlcywgYXJncywgbnVsbCk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBwcm9wZXJ0eShcbiAgbmFtZTogc3RyaW5nLFxuICBleHByZXNzaW9uOiBvLkV4cHJlc3Npb24sXG4gIHNhbml0aXplcjogby5FeHByZXNzaW9uIHwgbnVsbCxcbiAgc291cmNlU3BhbjogUGFyc2VTb3VyY2VTcGFuLFxuKTogaXIuVXBkYXRlT3Age1xuICBjb25zdCBhcmdzID0gW28ubGl0ZXJhbChuYW1lKSwgZXhwcmVzc2lvbl07XG4gIGlmIChzYW5pdGl6ZXIgIT09IG51bGwpIHtcbiAgICBhcmdzLnB1c2goc2FuaXRpemVyKTtcbiAgfVxuICByZXR1cm4gY2FsbChJZGVudGlmaWVycy5wcm9wZXJ0eSwgYXJncywgc291cmNlU3Bhbik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiB0d29XYXlQcm9wZXJ0eShcbiAgbmFtZTogc3RyaW5nLFxuICBleHByZXNzaW9uOiBvLkV4cHJlc3Npb24sXG4gIHNhbml0aXplcjogby5FeHByZXNzaW9uIHwgbnVsbCxcbiAgc291cmNlU3BhbjogUGFyc2VTb3VyY2VTcGFuLFxuKTogaXIuVXBkYXRlT3Age1xuICBjb25zdCBhcmdzID0gW28ubGl0ZXJhbChuYW1lKSwgZXhwcmVzc2lvbl07XG4gIGlmIChzYW5pdGl6ZXIgIT09IG51bGwpIHtcbiAgICBhcmdzLnB1c2goc2FuaXRpemVyKTtcbiAgfVxuICByZXR1cm4gY2FsbChJZGVudGlmaWVycy50d29XYXlQcm9wZXJ0eSwgYXJncywgc291cmNlU3Bhbik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBhdHRyaWJ1dGUoXG4gIG5hbWU6IHN0cmluZyxcbiAgZXhwcmVzc2lvbjogby5FeHByZXNzaW9uLFxuICBzYW5pdGl6ZXI6IG8uRXhwcmVzc2lvbiB8IG51bGwsXG4gIG5hbWVzcGFjZTogc3RyaW5nIHwgbnVsbCxcbik6IGlyLlVwZGF0ZU9wIHtcbiAgY29uc3QgYXJncyA9IFtvLmxpdGVyYWwobmFtZSksIGV4cHJlc3Npb25dO1xuICBpZiAoc2FuaXRpemVyICE9PSBudWxsIHx8IG5hbWVzcGFjZSAhPT0gbnVsbCkge1xuICAgIGFyZ3MucHVzaChzYW5pdGl6ZXIgPz8gby5saXRlcmFsKG51bGwpKTtcbiAgfVxuICBpZiAobmFtZXNwYWNlICE9PSBudWxsKSB7XG4gICAgYXJncy5wdXNoKG8ubGl0ZXJhbChuYW1lc3BhY2UpKTtcbiAgfVxuICByZXR1cm4gY2FsbChJZGVudGlmaWVycy5hdHRyaWJ1dGUsIGFyZ3MsIG51bGwpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gc3R5bGVQcm9wKFxuICBuYW1lOiBzdHJpbmcsXG4gIGV4cHJlc3Npb246IG8uRXhwcmVzc2lvbixcbiAgdW5pdDogc3RyaW5nIHwgbnVsbCxcbiAgc291cmNlU3BhbjogUGFyc2VTb3VyY2VTcGFuLFxuKTogaXIuVXBkYXRlT3Age1xuICBjb25zdCBhcmdzID0gW28ubGl0ZXJhbChuYW1lKSwgZXhwcmVzc2lvbl07XG4gIGlmICh1bml0ICE9PSBudWxsKSB7XG4gICAgYXJncy5wdXNoKG8ubGl0ZXJhbCh1bml0KSk7XG4gIH1cbiAgcmV0dXJuIGNhbGwoSWRlbnRpZmllcnMuc3R5bGVQcm9wLCBhcmdzLCBzb3VyY2VTcGFuKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGNsYXNzUHJvcChcbiAgbmFtZTogc3RyaW5nLFxuICBleHByZXNzaW9uOiBvLkV4cHJlc3Npb24sXG4gIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3Bhbixcbik6IGlyLlVwZGF0ZU9wIHtcbiAgcmV0dXJuIGNhbGwoSWRlbnRpZmllcnMuY2xhc3NQcm9wLCBbby5saXRlcmFsKG5hbWUpLCBleHByZXNzaW9uXSwgc291cmNlU3Bhbik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBzdHlsZU1hcChleHByZXNzaW9uOiBvLkV4cHJlc3Npb24sIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3Bhbik6IGlyLlVwZGF0ZU9wIHtcbiAgcmV0dXJuIGNhbGwoSWRlbnRpZmllcnMuc3R5bGVNYXAsIFtleHByZXNzaW9uXSwgc291cmNlU3Bhbik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBjbGFzc01hcChleHByZXNzaW9uOiBvLkV4cHJlc3Npb24sIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3Bhbik6IGlyLlVwZGF0ZU9wIHtcbiAgcmV0dXJuIGNhbGwoSWRlbnRpZmllcnMuY2xhc3NNYXAsIFtleHByZXNzaW9uXSwgc291cmNlU3Bhbik7XG59XG5cbmNvbnN0IFBJUEVfQklORElOR1M6IG8uRXh0ZXJuYWxSZWZlcmVuY2VbXSA9IFtcbiAgSWRlbnRpZmllcnMucGlwZUJpbmQxLFxuICBJZGVudGlmaWVycy5waXBlQmluZDIsXG4gIElkZW50aWZpZXJzLnBpcGVCaW5kMyxcbiAgSWRlbnRpZmllcnMucGlwZUJpbmQ0LFxuXTtcblxuZXhwb3J0IGZ1bmN0aW9uIHBpcGVCaW5kKHNsb3Q6IG51bWJlciwgdmFyT2Zmc2V0OiBudW1iZXIsIGFyZ3M6IG8uRXhwcmVzc2lvbltdKTogby5FeHByZXNzaW9uIHtcbiAgaWYgKGFyZ3MubGVuZ3RoIDwgMSB8fCBhcmdzLmxlbmd0aCA+IFBJUEVfQklORElOR1MubGVuZ3RoKSB7XG4gICAgdGhyb3cgbmV3IEVycm9yKGBwaXBlQmluZCgpIGFyZ3VtZW50IGNvdW50IG91dCBvZiBib3VuZHNgKTtcbiAgfVxuXG4gIGNvbnN0IGluc3RydWN0aW9uID0gUElQRV9CSU5ESU5HU1thcmdzLmxlbmd0aCAtIDFdO1xuICByZXR1cm4gby5pbXBvcnRFeHByKGluc3RydWN0aW9uKS5jYWxsRm4oW28ubGl0ZXJhbChzbG90KSwgby5saXRlcmFsKHZhck9mZnNldCksIC4uLmFyZ3NdKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIHBpcGVCaW5kVihzbG90OiBudW1iZXIsIHZhck9mZnNldDogbnVtYmVyLCBhcmdzOiBvLkV4cHJlc3Npb24pOiBvLkV4cHJlc3Npb24ge1xuICByZXR1cm4gby5pbXBvcnRFeHByKElkZW50aWZpZXJzLnBpcGVCaW5kVikuY2FsbEZuKFtvLmxpdGVyYWwoc2xvdCksIG8ubGl0ZXJhbCh2YXJPZmZzZXQpLCBhcmdzXSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiB0ZXh0SW50ZXJwb2xhdGUoXG4gIHN0cmluZ3M6IHN0cmluZ1tdLFxuICBleHByZXNzaW9uczogby5FeHByZXNzaW9uW10sXG4gIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3Bhbixcbik6IGlyLlVwZGF0ZU9wIHtcbiAgY29uc3QgaW50ZXJwb2xhdGlvbkFyZ3MgPSBjb2xsYXRlSW50ZXJwb2xhdGlvbkFyZ3Moc3RyaW5ncywgZXhwcmVzc2lvbnMpO1xuICByZXR1cm4gY2FsbFZhcmlhZGljSW5zdHJ1Y3Rpb24oVEVYVF9JTlRFUlBPTEFURV9DT05GSUcsIFtdLCBpbnRlcnBvbGF0aW9uQXJncywgW10sIHNvdXJjZVNwYW4pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gaTE4bkV4cChleHByOiBvLkV4cHJlc3Npb24sIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3BhbiB8IG51bGwpOiBpci5VcGRhdGVPcCB7XG4gIHJldHVybiBjYWxsKElkZW50aWZpZXJzLmkxOG5FeHAsIFtleHByXSwgc291cmNlU3Bhbik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBpMThuQXBwbHkoc2xvdDogbnVtYmVyLCBzb3VyY2VTcGFuOiBQYXJzZVNvdXJjZVNwYW4gfCBudWxsKTogaXIuVXBkYXRlT3Age1xuICByZXR1cm4gY2FsbChJZGVudGlmaWVycy5pMThuQXBwbHksIFtvLmxpdGVyYWwoc2xvdCldLCBzb3VyY2VTcGFuKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIHByb3BlcnR5SW50ZXJwb2xhdGUoXG4gIG5hbWU6IHN0cmluZyxcbiAgc3RyaW5nczogc3RyaW5nW10sXG4gIGV4cHJlc3Npb25zOiBvLkV4cHJlc3Npb25bXSxcbiAgc2FuaXRpemVyOiBvLkV4cHJlc3Npb24gfCBudWxsLFxuICBzb3VyY2VTcGFuOiBQYXJzZVNvdXJjZVNwYW4sXG4pOiBpci5VcGRhdGVPcCB7XG4gIGNvbnN0IGludGVycG9sYXRpb25BcmdzID0gY29sbGF0ZUludGVycG9sYXRpb25BcmdzKHN0cmluZ3MsIGV4cHJlc3Npb25zKTtcbiAgY29uc3QgZXh0cmFBcmdzID0gW107XG4gIGlmIChzYW5pdGl6ZXIgIT09IG51bGwpIHtcbiAgICBleHRyYUFyZ3MucHVzaChzYW5pdGl6ZXIpO1xuICB9XG5cbiAgcmV0dXJuIGNhbGxWYXJpYWRpY0luc3RydWN0aW9uKFxuICAgIFBST1BFUlRZX0lOVEVSUE9MQVRFX0NPTkZJRyxcbiAgICBbby5saXRlcmFsKG5hbWUpXSxcbiAgICBpbnRlcnBvbGF0aW9uQXJncyxcbiAgICBleHRyYUFyZ3MsXG4gICAgc291cmNlU3BhbixcbiAgKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGF0dHJpYnV0ZUludGVycG9sYXRlKFxuICBuYW1lOiBzdHJpbmcsXG4gIHN0cmluZ3M6IHN0cmluZ1tdLFxuICBleHByZXNzaW9uczogby5FeHByZXNzaW9uW10sXG4gIHNhbml0aXplcjogby5FeHByZXNzaW9uIHwgbnVsbCxcbiAgc291cmNlU3BhbjogUGFyc2VTb3VyY2VTcGFuLFxuKTogaXIuVXBkYXRlT3Age1xuICBjb25zdCBpbnRlcnBvbGF0aW9uQXJncyA9IGNvbGxhdGVJbnRlcnBvbGF0aW9uQXJncyhzdHJpbmdzLCBleHByZXNzaW9ucyk7XG4gIGNvbnN0IGV4dHJhQXJncyA9IFtdO1xuICBpZiAoc2FuaXRpemVyICE9PSBudWxsKSB7XG4gICAgZXh0cmFBcmdzLnB1c2goc2FuaXRpemVyKTtcbiAgfVxuXG4gIHJldHVybiBjYWxsVmFyaWFkaWNJbnN0cnVjdGlvbihcbiAgICBBVFRSSUJVVEVfSU5URVJQT0xBVEVfQ09ORklHLFxuICAgIFtvLmxpdGVyYWwobmFtZSldLFxuICAgIGludGVycG9sYXRpb25BcmdzLFxuICAgIGV4dHJhQXJncyxcbiAgICBzb3VyY2VTcGFuLFxuICApO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gc3R5bGVQcm9wSW50ZXJwb2xhdGUoXG4gIG5hbWU6IHN0cmluZyxcbiAgc3RyaW5nczogc3RyaW5nW10sXG4gIGV4cHJlc3Npb25zOiBvLkV4cHJlc3Npb25bXSxcbiAgdW5pdDogc3RyaW5nIHwgbnVsbCxcbiAgc291cmNlU3BhbjogUGFyc2VTb3VyY2VTcGFuLFxuKTogaXIuVXBkYXRlT3Age1xuICBjb25zdCBpbnRlcnBvbGF0aW9uQXJncyA9IGNvbGxhdGVJbnRlcnBvbGF0aW9uQXJncyhzdHJpbmdzLCBleHByZXNzaW9ucyk7XG4gIGNvbnN0IGV4dHJhQXJnczogby5FeHByZXNzaW9uW10gPSBbXTtcbiAgaWYgKHVuaXQgIT09IG51bGwpIHtcbiAgICBleHRyYUFyZ3MucHVzaChvLmxpdGVyYWwodW5pdCkpO1xuICB9XG5cbiAgcmV0dXJuIGNhbGxWYXJpYWRpY0luc3RydWN0aW9uKFxuICAgIFNUWUxFX1BST1BfSU5URVJQT0xBVEVfQ09ORklHLFxuICAgIFtvLmxpdGVyYWwobmFtZSldLFxuICAgIGludGVycG9sYXRpb25BcmdzLFxuICAgIGV4dHJhQXJncyxcbiAgICBzb3VyY2VTcGFuLFxuICApO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gc3R5bGVNYXBJbnRlcnBvbGF0ZShcbiAgc3RyaW5nczogc3RyaW5nW10sXG4gIGV4cHJlc3Npb25zOiBvLkV4cHJlc3Npb25bXSxcbiAgc291cmNlU3BhbjogUGFyc2VTb3VyY2VTcGFuLFxuKTogaXIuVXBkYXRlT3Age1xuICBjb25zdCBpbnRlcnBvbGF0aW9uQXJncyA9IGNvbGxhdGVJbnRlcnBvbGF0aW9uQXJncyhzdHJpbmdzLCBleHByZXNzaW9ucyk7XG5cbiAgcmV0dXJuIGNhbGxWYXJpYWRpY0luc3RydWN0aW9uKFxuICAgIFNUWUxFX01BUF9JTlRFUlBPTEFURV9DT05GSUcsXG4gICAgW10sXG4gICAgaW50ZXJwb2xhdGlvbkFyZ3MsXG4gICAgW10sXG4gICAgc291cmNlU3BhbixcbiAgKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGNsYXNzTWFwSW50ZXJwb2xhdGUoXG4gIHN0cmluZ3M6IHN0cmluZ1tdLFxuICBleHByZXNzaW9uczogby5FeHByZXNzaW9uW10sXG4gIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3Bhbixcbik6IGlyLlVwZGF0ZU9wIHtcbiAgY29uc3QgaW50ZXJwb2xhdGlvbkFyZ3MgPSBjb2xsYXRlSW50ZXJwb2xhdGlvbkFyZ3Moc3RyaW5ncywgZXhwcmVzc2lvbnMpO1xuXG4gIHJldHVybiBjYWxsVmFyaWFkaWNJbnN0cnVjdGlvbihcbiAgICBDTEFTU19NQVBfSU5URVJQT0xBVEVfQ09ORklHLFxuICAgIFtdLFxuICAgIGludGVycG9sYXRpb25BcmdzLFxuICAgIFtdLFxuICAgIHNvdXJjZVNwYW4sXG4gICk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBob3N0UHJvcGVydHkoXG4gIG5hbWU6IHN0cmluZyxcbiAgZXhwcmVzc2lvbjogby5FeHByZXNzaW9uLFxuICBzYW5pdGl6ZXI6IG8uRXhwcmVzc2lvbiB8IG51bGwsXG4gIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3BhbiB8IG51bGwsXG4pOiBpci5VcGRhdGVPcCB7XG4gIGNvbnN0IGFyZ3MgPSBbby5saXRlcmFsKG5hbWUpLCBleHByZXNzaW9uXTtcbiAgaWYgKHNhbml0aXplciAhPT0gbnVsbCkge1xuICAgIGFyZ3MucHVzaChzYW5pdGl6ZXIpO1xuICB9XG4gIHJldHVybiBjYWxsKElkZW50aWZpZXJzLmhvc3RQcm9wZXJ0eSwgYXJncywgc291cmNlU3Bhbik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBzeW50aGV0aWNIb3N0UHJvcGVydHkoXG4gIG5hbWU6IHN0cmluZyxcbiAgZXhwcmVzc2lvbjogby5FeHByZXNzaW9uLFxuICBzb3VyY2VTcGFuOiBQYXJzZVNvdXJjZVNwYW4gfCBudWxsLFxuKTogaXIuVXBkYXRlT3Age1xuICByZXR1cm4gY2FsbChJZGVudGlmaWVycy5zeW50aGV0aWNIb3N0UHJvcGVydHksIFtvLmxpdGVyYWwobmFtZSksIGV4cHJlc3Npb25dLCBzb3VyY2VTcGFuKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIHB1cmVGdW5jdGlvbihcbiAgdmFyT2Zmc2V0OiBudW1iZXIsXG4gIGZuOiBvLkV4cHJlc3Npb24sXG4gIGFyZ3M6IG8uRXhwcmVzc2lvbltdLFxuKTogby5FeHByZXNzaW9uIHtcbiAgcmV0dXJuIGNhbGxWYXJpYWRpY0luc3RydWN0aW9uRXhwcihcbiAgICBQVVJFX0ZVTkNUSU9OX0NPTkZJRyxcbiAgICBbby5saXRlcmFsKHZhck9mZnNldCksIGZuXSxcbiAgICBhcmdzLFxuICAgIFtdLFxuICAgIG51bGwsXG4gICk7XG59XG5cbi8qKlxuICogQ29sbGF0ZXMgdGhlIHN0cmluZyBhbiBleHByZXNzaW9uIGFyZ3VtZW50cyBmb3IgYW4gaW50ZXJwb2xhdGlvbiBpbnN0cnVjdGlvbi5cbiAqL1xuZnVuY3Rpb24gY29sbGF0ZUludGVycG9sYXRpb25BcmdzKHN0cmluZ3M6IHN0cmluZ1tdLCBleHByZXNzaW9uczogby5FeHByZXNzaW9uW10pOiBvLkV4cHJlc3Npb25bXSB7XG4gIGlmIChzdHJpbmdzLmxlbmd0aCA8IDEgfHwgZXhwcmVzc2lvbnMubGVuZ3RoICE9PSBzdHJpbmdzLmxlbmd0aCAtIDEpIHtcbiAgICB0aHJvdyBuZXcgRXJyb3IoXG4gICAgICBgQXNzZXJ0aW9uRXJyb3I6IGV4cGVjdGVkIHNwZWNpZmljIHNoYXBlIG9mIGFyZ3MgZm9yIHN0cmluZ3MvZXhwcmVzc2lvbnMgaW4gaW50ZXJwb2xhdGlvbmAsXG4gICAgKTtcbiAgfVxuICBjb25zdCBpbnRlcnBvbGF0aW9uQXJnczogby5FeHByZXNzaW9uW10gPSBbXTtcblxuICBpZiAoZXhwcmVzc2lvbnMubGVuZ3RoID09PSAxICYmIHN0cmluZ3NbMF0gPT09ICcnICYmIHN0cmluZ3NbMV0gPT09ICcnKSB7XG4gICAgaW50ZXJwb2xhdGlvbkFyZ3MucHVzaChleHByZXNzaW9uc1swXSk7XG4gIH0gZWxzZSB7XG4gICAgbGV0IGlkeDogbnVtYmVyO1xuICAgIGZvciAoaWR4ID0gMDsgaWR4IDwgZXhwcmVzc2lvbnMubGVuZ3RoOyBpZHgrKykge1xuICAgICAgaW50ZXJwb2xhdGlvbkFyZ3MucHVzaChvLmxpdGVyYWwoc3RyaW5nc1tpZHhdKSwgZXhwcmVzc2lvbnNbaWR4XSk7XG4gICAgfVxuICAgIC8vIGlkeCBwb2ludHMgYXQgdGhlIGxhc3Qgc3RyaW5nLlxuICAgIGludGVycG9sYXRpb25BcmdzLnB1c2goby5saXRlcmFsKHN0cmluZ3NbaWR4XSkpO1xuICB9XG5cbiAgcmV0dXJuIGludGVycG9sYXRpb25BcmdzO1xufVxuXG5mdW5jdGlvbiBjYWxsPE9wVCBleHRlbmRzIGlyLkNyZWF0ZU9wIHwgaXIuVXBkYXRlT3A+KFxuICBpbnN0cnVjdGlvbjogby5FeHRlcm5hbFJlZmVyZW5jZSxcbiAgYXJnczogby5FeHByZXNzaW9uW10sXG4gIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3BhbiB8IG51bGwsXG4pOiBPcFQge1xuICBjb25zdCBleHByID0gby5pbXBvcnRFeHByKGluc3RydWN0aW9uKS5jYWxsRm4oYXJncywgc291cmNlU3Bhbik7XG4gIHJldHVybiBpci5jcmVhdGVTdGF0ZW1lbnRPcChuZXcgby5FeHByZXNzaW9uU3RhdGVtZW50KGV4cHIsIHNvdXJjZVNwYW4pKSBhcyBPcFQ7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBjb25kaXRpb25hbChcbiAgY29uZGl0aW9uOiBvLkV4cHJlc3Npb24sXG4gIGNvbnRleHRWYWx1ZTogby5FeHByZXNzaW9uIHwgbnVsbCxcbiAgc291cmNlU3BhbjogUGFyc2VTb3VyY2VTcGFuIHwgbnVsbCxcbik6IGlyLlVwZGF0ZU9wIHtcbiAgY29uc3QgYXJncyA9IFtjb25kaXRpb25dO1xuICBpZiAoY29udGV4dFZhbHVlICE9PSBudWxsKSB7XG4gICAgYXJncy5wdXNoKGNvbnRleHRWYWx1ZSk7XG4gIH1cbiAgcmV0dXJuIGNhbGwoSWRlbnRpZmllcnMuY29uZGl0aW9uYWwsIGFyZ3MsIHNvdXJjZVNwYW4pO1xufVxuXG4vKipcbiAqIERlc2NyaWJlcyBhIHNwZWNpZmljIGZsYXZvciBvZiBpbnN0cnVjdGlvbiB1c2VkIHRvIHJlcHJlc2VudCB2YXJpYWRpYyBpbnN0cnVjdGlvbnMsIHdoaWNoXG4gKiBoYXZlIHNvbWUgbnVtYmVyIG9mIHZhcmlhbnRzIGZvciBzcGVjaWZpYyBhcmd1bWVudCBjb3VudHMuXG4gKi9cbmludGVyZmFjZSBWYXJpYWRpY0luc3RydWN0aW9uQ29uZmlnIHtcbiAgY29uc3RhbnQ6IG8uRXh0ZXJuYWxSZWZlcmVuY2VbXTtcbiAgdmFyaWFibGU6IG8uRXh0ZXJuYWxSZWZlcmVuY2UgfCBudWxsO1xuICBtYXBwaW5nOiAoYXJnQ291bnQ6IG51bWJlcikgPT4gbnVtYmVyO1xufVxuXG4vKipcbiAqIGBJbnRlcnBvbGF0aW9uQ29uZmlnYCBmb3IgdGhlIGB0ZXh0SW50ZXJwb2xhdGVgIGluc3RydWN0aW9uLlxuICovXG5jb25zdCBURVhUX0lOVEVSUE9MQVRFX0NPTkZJRzogVmFyaWFkaWNJbnN0cnVjdGlvbkNvbmZpZyA9IHtcbiAgY29uc3RhbnQ6IFtcbiAgICBJZGVudGlmaWVycy50ZXh0SW50ZXJwb2xhdGUsXG4gICAgSWRlbnRpZmllcnMudGV4dEludGVycG9sYXRlMSxcbiAgICBJZGVudGlmaWVycy50ZXh0SW50ZXJwb2xhdGUyLFxuICAgIElkZW50aWZpZXJzLnRleHRJbnRlcnBvbGF0ZTMsXG4gICAgSWRlbnRpZmllcnMudGV4dEludGVycG9sYXRlNCxcbiAgICBJZGVudGlmaWVycy50ZXh0SW50ZXJwb2xhdGU1LFxuICAgIElkZW50aWZpZXJzLnRleHRJbnRlcnBvbGF0ZTYsXG4gICAgSWRlbnRpZmllcnMudGV4dEludGVycG9sYXRlNyxcbiAgICBJZGVudGlmaWVycy50ZXh0SW50ZXJwb2xhdGU4LFxuICBdLFxuICB2YXJpYWJsZTogSWRlbnRpZmllcnMudGV4dEludGVycG9sYXRlVixcbiAgbWFwcGluZzogKG4pID0+IHtcbiAgICBpZiAobiAlIDIgPT09IDApIHtcbiAgICAgIHRocm93IG5ldyBFcnJvcihgRXhwZWN0ZWQgb2RkIG51bWJlciBvZiBhcmd1bWVudHNgKTtcbiAgICB9XG4gICAgcmV0dXJuIChuIC0gMSkgLyAyO1xuICB9LFxufTtcblxuLyoqXG4gKiBgSW50ZXJwb2xhdGlvbkNvbmZpZ2AgZm9yIHRoZSBgcHJvcGVydHlJbnRlcnBvbGF0ZWAgaW5zdHJ1Y3Rpb24uXG4gKi9cbmNvbnN0IFBST1BFUlRZX0lOVEVSUE9MQVRFX0NPTkZJRzogVmFyaWFkaWNJbnN0cnVjdGlvbkNvbmZpZyA9IHtcbiAgY29uc3RhbnQ6IFtcbiAgICBJZGVudGlmaWVycy5wcm9wZXJ0eUludGVycG9sYXRlLFxuICAgIElkZW50aWZpZXJzLnByb3BlcnR5SW50ZXJwb2xhdGUxLFxuICAgIElkZW50aWZpZXJzLnByb3BlcnR5SW50ZXJwb2xhdGUyLFxuICAgIElkZW50aWZpZXJzLnByb3BlcnR5SW50ZXJwb2xhdGUzLFxuICAgIElkZW50aWZpZXJzLnByb3BlcnR5SW50ZXJwb2xhdGU0LFxuICAgIElkZW50aWZpZXJzLnByb3BlcnR5SW50ZXJwb2xhdGU1LFxuICAgIElkZW50aWZpZXJzLnByb3BlcnR5SW50ZXJwb2xhdGU2LFxuICAgIElkZW50aWZpZXJzLnByb3BlcnR5SW50ZXJwb2xhdGU3LFxuICAgIElkZW50aWZpZXJzLnByb3BlcnR5SW50ZXJwb2xhdGU4LFxuICBdLFxuICB2YXJpYWJsZTogSWRlbnRpZmllcnMucHJvcGVydHlJbnRlcnBvbGF0ZVYsXG4gIG1hcHBpbmc6IChuKSA9PiB7XG4gICAgaWYgKG4gJSAyID09PSAwKSB7XG4gICAgICB0aHJvdyBuZXcgRXJyb3IoYEV4cGVjdGVkIG9kZCBudW1iZXIgb2YgYXJndW1lbnRzYCk7XG4gICAgfVxuICAgIHJldHVybiAobiAtIDEpIC8gMjtcbiAgfSxcbn07XG5cbi8qKlxuICogYEludGVycG9sYXRpb25Db25maWdgIGZvciB0aGUgYHN0eWxlUHJvcEludGVycG9sYXRlYCBpbnN0cnVjdGlvbi5cbiAqL1xuY29uc3QgU1RZTEVfUFJPUF9JTlRFUlBPTEFURV9DT05GSUc6IFZhcmlhZGljSW5zdHJ1Y3Rpb25Db25maWcgPSB7XG4gIGNvbnN0YW50OiBbXG4gICAgSWRlbnRpZmllcnMuc3R5bGVQcm9wLFxuICAgIElkZW50aWZpZXJzLnN0eWxlUHJvcEludGVycG9sYXRlMSxcbiAgICBJZGVudGlmaWVycy5zdHlsZVByb3BJbnRlcnBvbGF0ZTIsXG4gICAgSWRlbnRpZmllcnMuc3R5bGVQcm9wSW50ZXJwb2xhdGUzLFxuICAgIElkZW50aWZpZXJzLnN0eWxlUHJvcEludGVycG9sYXRlNCxcbiAgICBJZGVudGlmaWVycy5zdHlsZVByb3BJbnRlcnBvbGF0ZTUsXG4gICAgSWRlbnRpZmllcnMuc3R5bGVQcm9wSW50ZXJwb2xhdGU2LFxuICAgIElkZW50aWZpZXJzLnN0eWxlUHJvcEludGVycG9sYXRlNyxcbiAgICBJZGVudGlmaWVycy5zdHlsZVByb3BJbnRlcnBvbGF0ZTgsXG4gIF0sXG4gIHZhcmlhYmxlOiBJZGVudGlmaWVycy5zdHlsZVByb3BJbnRlcnBvbGF0ZVYsXG4gIG1hcHBpbmc6IChuKSA9PiB7XG4gICAgaWYgKG4gJSAyID09PSAwKSB7XG4gICAgICB0aHJvdyBuZXcgRXJyb3IoYEV4cGVjdGVkIG9kZCBudW1iZXIgb2YgYXJndW1lbnRzYCk7XG4gICAgfVxuICAgIHJldHVybiAobiAtIDEpIC8gMjtcbiAgfSxcbn07XG5cbi8qKlxuICogYEludGVycG9sYXRpb25Db25maWdgIGZvciB0aGUgYGF0dHJpYnV0ZUludGVycG9sYXRlYCBpbnN0cnVjdGlvbi5cbiAqL1xuY29uc3QgQVRUUklCVVRFX0lOVEVSUE9MQVRFX0NPTkZJRzogVmFyaWFkaWNJbnN0cnVjdGlvbkNvbmZpZyA9IHtcbiAgY29uc3RhbnQ6IFtcbiAgICBJZGVudGlmaWVycy5hdHRyaWJ1dGUsXG4gICAgSWRlbnRpZmllcnMuYXR0cmlidXRlSW50ZXJwb2xhdGUxLFxuICAgIElkZW50aWZpZXJzLmF0dHJpYnV0ZUludGVycG9sYXRlMixcbiAgICBJZGVudGlmaWVycy5hdHRyaWJ1dGVJbnRlcnBvbGF0ZTMsXG4gICAgSWRlbnRpZmllcnMuYXR0cmlidXRlSW50ZXJwb2xhdGU0LFxuICAgIElkZW50aWZpZXJzLmF0dHJpYnV0ZUludGVycG9sYXRlNSxcbiAgICBJZGVudGlmaWVycy5hdHRyaWJ1dGVJbnRlcnBvbGF0ZTYsXG4gICAgSWRlbnRpZmllcnMuYXR0cmlidXRlSW50ZXJwb2xhdGU3LFxuICAgIElkZW50aWZpZXJzLmF0dHJpYnV0ZUludGVycG9sYXRlOCxcbiAgXSxcbiAgdmFyaWFibGU6IElkZW50aWZpZXJzLmF0dHJpYnV0ZUludGVycG9sYXRlVixcbiAgbWFwcGluZzogKG4pID0+IHtcbiAgICBpZiAobiAlIDIgPT09IDApIHtcbiAgICAgIHRocm93IG5ldyBFcnJvcihgRXhwZWN0ZWQgb2RkIG51bWJlciBvZiBhcmd1bWVudHNgKTtcbiAgICB9XG4gICAgcmV0dXJuIChuIC0gMSkgLyAyO1xuICB9LFxufTtcblxuLyoqXG4gKiBgSW50ZXJwb2xhdGlvbkNvbmZpZ2AgZm9yIHRoZSBgc3R5bGVNYXBJbnRlcnBvbGF0ZWAgaW5zdHJ1Y3Rpb24uXG4gKi9cbmNvbnN0IFNUWUxFX01BUF9JTlRFUlBPTEFURV9DT05GSUc6IFZhcmlhZGljSW5zdHJ1Y3Rpb25Db25maWcgPSB7XG4gIGNvbnN0YW50OiBbXG4gICAgSWRlbnRpZmllcnMuc3R5bGVNYXAsXG4gICAgSWRlbnRpZmllcnMuc3R5bGVNYXBJbnRlcnBvbGF0ZTEsXG4gICAgSWRlbnRpZmllcnMuc3R5bGVNYXBJbnRlcnBvbGF0ZTIsXG4gICAgSWRlbnRpZmllcnMuc3R5bGVNYXBJbnRlcnBvbGF0ZTMsXG4gICAgSWRlbnRpZmllcnMuc3R5bGVNYXBJbnRlcnBvbGF0ZTQsXG4gICAgSWRlbnRpZmllcnMuc3R5bGVNYXBJbnRlcnBvbGF0ZTUsXG4gICAgSWRlbnRpZmllcnMuc3R5bGVNYXBJbnRlcnBvbGF0ZTYsXG4gICAgSWRlbnRpZmllcnMuc3R5bGVNYXBJbnRlcnBvbGF0ZTcsXG4gICAgSWRlbnRpZmllcnMuc3R5bGVNYXBJbnRlcnBvbGF0ZTgsXG4gIF0sXG4gIHZhcmlhYmxlOiBJZGVudGlmaWVycy5zdHlsZU1hcEludGVycG9sYXRlVixcbiAgbWFwcGluZzogKG4pID0+IHtcbiAgICBpZiAobiAlIDIgPT09IDApIHtcbiAgICAgIHRocm93IG5ldyBFcnJvcihgRXhwZWN0ZWQgb2RkIG51bWJlciBvZiBhcmd1bWVudHNgKTtcbiAgICB9XG4gICAgcmV0dXJuIChuIC0gMSkgLyAyO1xuICB9LFxufTtcblxuLyoqXG4gKiBgSW50ZXJwb2xhdGlvbkNvbmZpZ2AgZm9yIHRoZSBgY2xhc3NNYXBJbnRlcnBvbGF0ZWAgaW5zdHJ1Y3Rpb24uXG4gKi9cbmNvbnN0IENMQVNTX01BUF9JTlRFUlBPTEFURV9DT05GSUc6IFZhcmlhZGljSW5zdHJ1Y3Rpb25Db25maWcgPSB7XG4gIGNvbnN0YW50OiBbXG4gICAgSWRlbnRpZmllcnMuY2xhc3NNYXAsXG4gICAgSWRlbnRpZmllcnMuY2xhc3NNYXBJbnRlcnBvbGF0ZTEsXG4gICAgSWRlbnRpZmllcnMuY2xhc3NNYXBJbnRlcnBvbGF0ZTIsXG4gICAgSWRlbnRpZmllcnMuY2xhc3NNYXBJbnRlcnBvbGF0ZTMsXG4gICAgSWRlbnRpZmllcnMuY2xhc3NNYXBJbnRlcnBvbGF0ZTQsXG4gICAgSWRlbnRpZmllcnMuY2xhc3NNYXBJbnRlcnBvbGF0ZTUsXG4gICAgSWRlbnRpZmllcnMuY2xhc3NNYXBJbnRlcnBvbGF0ZTYsXG4gICAgSWRlbnRpZmllcnMuY2xhc3NNYXBJbnRlcnBvbGF0ZTcsXG4gICAgSWRlbnRpZmllcnMuY2xhc3NNYXBJbnRlcnBvbGF0ZTgsXG4gIF0sXG4gIHZhcmlhYmxlOiBJZGVudGlmaWVycy5jbGFzc01hcEludGVycG9sYXRlVixcbiAgbWFwcGluZzogKG4pID0+IHtcbiAgICBpZiAobiAlIDIgPT09IDApIHtcbiAgICAgIHRocm93IG5ldyBFcnJvcihgRXhwZWN0ZWQgb2RkIG51bWJlciBvZiBhcmd1bWVudHNgKTtcbiAgICB9XG4gICAgcmV0dXJuIChuIC0gMSkgLyAyO1xuICB9LFxufTtcblxuY29uc3QgUFVSRV9GVU5DVElPTl9DT05GSUc6IFZhcmlhZGljSW5zdHJ1Y3Rpb25Db25maWcgPSB7XG4gIGNvbnN0YW50OiBbXG4gICAgSWRlbnRpZmllcnMucHVyZUZ1bmN0aW9uMCxcbiAgICBJZGVudGlmaWVycy5wdXJlRnVuY3Rpb24xLFxuICAgIElkZW50aWZpZXJzLnB1cmVGdW5jdGlvbjIsXG4gICAgSWRlbnRpZmllcnMucHVyZUZ1bmN0aW9uMyxcbiAgICBJZGVudGlmaWVycy5wdXJlRnVuY3Rpb240LFxuICAgIElkZW50aWZpZXJzLnB1cmVGdW5jdGlvbjUsXG4gICAgSWRlbnRpZmllcnMucHVyZUZ1bmN0aW9uNixcbiAgICBJZGVudGlmaWVycy5wdXJlRnVuY3Rpb243LFxuICAgIElkZW50aWZpZXJzLnB1cmVGdW5jdGlvbjgsXG4gIF0sXG4gIHZhcmlhYmxlOiBJZGVudGlmaWVycy5wdXJlRnVuY3Rpb25WLFxuICBtYXBwaW5nOiAobikgPT4gbixcbn07XG5cbmZ1bmN0aW9uIGNhbGxWYXJpYWRpY0luc3RydWN0aW9uRXhwcihcbiAgY29uZmlnOiBWYXJpYWRpY0luc3RydWN0aW9uQ29uZmlnLFxuICBiYXNlQXJnczogby5FeHByZXNzaW9uW10sXG4gIGludGVycG9sYXRpb25BcmdzOiBvLkV4cHJlc3Npb25bXSxcbiAgZXh0cmFBcmdzOiBvLkV4cHJlc3Npb25bXSxcbiAgc291cmNlU3BhbjogUGFyc2VTb3VyY2VTcGFuIHwgbnVsbCxcbik6IG8uRXhwcmVzc2lvbiB7XG4gIGNvbnN0IG4gPSBjb25maWcubWFwcGluZyhpbnRlcnBvbGF0aW9uQXJncy5sZW5ndGgpO1xuICBpZiAobiA8IGNvbmZpZy5jb25zdGFudC5sZW5ndGgpIHtcbiAgICAvLyBDb25zdGFudCBjYWxsaW5nIHBhdHRlcm4uXG4gICAgcmV0dXJuIG9cbiAgICAgIC5pbXBvcnRFeHByKGNvbmZpZy5jb25zdGFudFtuXSlcbiAgICAgIC5jYWxsRm4oWy4uLmJhc2VBcmdzLCAuLi5pbnRlcnBvbGF0aW9uQXJncywgLi4uZXh0cmFBcmdzXSwgc291cmNlU3Bhbik7XG4gIH0gZWxzZSBpZiAoY29uZmlnLnZhcmlhYmxlICE9PSBudWxsKSB7XG4gICAgLy8gVmFyaWFibGUgY2FsbGluZyBwYXR0ZXJuLlxuICAgIHJldHVybiBvXG4gICAgICAuaW1wb3J0RXhwcihjb25maWcudmFyaWFibGUpXG4gICAgICAuY2FsbEZuKFsuLi5iYXNlQXJncywgby5saXRlcmFsQXJyKGludGVycG9sYXRpb25BcmdzKSwgLi4uZXh0cmFBcmdzXSwgc291cmNlU3Bhbik7XG4gIH0gZWxzZSB7XG4gICAgdGhyb3cgbmV3IEVycm9yKGBBc3NlcnRpb25FcnJvcjogdW5hYmxlIHRvIGNhbGwgdmFyaWFkaWMgZnVuY3Rpb25gKTtcbiAgfVxufVxuXG5mdW5jdGlvbiBjYWxsVmFyaWFkaWNJbnN0cnVjdGlvbihcbiAgY29uZmlnOiBWYXJpYWRpY0luc3RydWN0aW9uQ29uZmlnLFxuICBiYXNlQXJnczogby5FeHByZXNzaW9uW10sXG4gIGludGVycG9sYXRpb25BcmdzOiBvLkV4cHJlc3Npb25bXSxcbiAgZXh0cmFBcmdzOiBvLkV4cHJlc3Npb25bXSxcbiAgc291cmNlU3BhbjogUGFyc2VTb3VyY2VTcGFuIHwgbnVsbCxcbik6IGlyLlVwZGF0ZU9wIHtcbiAgcmV0dXJuIGlyLmNyZWF0ZVN0YXRlbWVudE9wKFxuICAgIGNhbGxWYXJpYWRpY0luc3RydWN0aW9uRXhwcihcbiAgICAgIGNvbmZpZyxcbiAgICAgIGJhc2VBcmdzLFxuICAgICAgaW50ZXJwb2xhdGlvbkFyZ3MsXG4gICAgICBleHRyYUFyZ3MsXG4gICAgICBzb3VyY2VTcGFuLFxuICAgICkudG9TdG10KCksXG4gICk7XG59XG4iXX0=