/**
 * @license
 * Copyright Google LLC All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
import { SecurityContext } from '../../../core';
import * as e from '../../../expression_parser/ast';
import * as i18n from '../../../i18n/i18n_ast';
import { splitNsName } from '../../../ml_parser/tags';
import * as o from '../../../output/output_ast';
import { ParseSourceSpan } from '../../../parse_util';
import * as t from '../../../render3/r3_ast';
import { icuFromI18nMessage } from '../../../render3/view/i18n/util';
import { DomElementSchemaRegistry } from '../../../schema/dom_element_schema_registry';
import * as ir from '../ir';
import { ComponentCompilationJob, HostBindingCompilationJob, } from './compilation';
import { BINARY_OPERATORS, namespaceForKey, prefixWithNamespace } from './conversion';
const compatibilityMode = ir.CompatibilityMode.TemplateDefinitionBuilder;
// Schema containing DOM elements and their properties.
const domSchema = new DomElementSchemaRegistry();
// Tag name of the `ng-template` element.
const NG_TEMPLATE_TAG_NAME = 'ng-template';
export function isI18nRootNode(meta) {
    return meta instanceof i18n.Message;
}
export function isSingleI18nIcu(meta) {
    return isI18nRootNode(meta) && meta.nodes.length === 1 && meta.nodes[0] instanceof i18n.Icu;
}
/**
 * Process a template AST and convert it into a `ComponentCompilation` in the intermediate
 * representation.
 * TODO: Refactor more of the ingestion code into phases.
 */
export function ingestComponent(componentName, template, constantPool, relativeContextFilePath, i18nUseExternalIds, deferMeta, allDeferrableDepsFn) {
    const job = new ComponentCompilationJob(componentName, constantPool, compatibilityMode, relativeContextFilePath, i18nUseExternalIds, deferMeta, allDeferrableDepsFn);
    ingestNodes(job.root, template);
    return job;
}
/**
 * Process a host binding AST and convert it into a `HostBindingCompilationJob` in the intermediate
 * representation.
 */
export function ingestHostBinding(input, bindingParser, constantPool) {
    const job = new HostBindingCompilationJob(input.componentName, constantPool, compatibilityMode);
    for (const property of input.properties ?? []) {
        let bindingKind = ir.BindingKind.Property;
        // TODO: this should really be handled in the parser.
        if (property.name.startsWith('attr.')) {
            property.name = property.name.substring('attr.'.length);
            bindingKind = ir.BindingKind.Attribute;
        }
        if (property.isAnimation) {
            bindingKind = ir.BindingKind.Animation;
        }
        const securityContexts = bindingParser
            .calcPossibleSecurityContexts(input.componentSelector, property.name, bindingKind === ir.BindingKind.Attribute)
            .filter((context) => context !== SecurityContext.NONE);
        ingestHostProperty(job, property, bindingKind, securityContexts);
    }
    for (const [name, expr] of Object.entries(input.attributes) ?? []) {
        const securityContexts = bindingParser
            .calcPossibleSecurityContexts(input.componentSelector, name, true)
            .filter((context) => context !== SecurityContext.NONE);
        ingestHostAttribute(job, name, expr, securityContexts);
    }
    for (const event of input.events ?? []) {
        ingestHostEvent(job, event);
    }
    return job;
}
// TODO: We should refactor the parser to use the same types and structures for host bindings as
// with ordinary components. This would allow us to share a lot more ingestion code.
export function ingestHostProperty(job, property, bindingKind, securityContexts) {
    let expression;
    const ast = property.expression.ast;
    if (ast instanceof e.Interpolation) {
        expression = new ir.Interpolation(ast.strings, ast.expressions.map((expr) => convertAst(expr, job, property.sourceSpan)), []);
    }
    else {
        expression = convertAst(ast, job, property.sourceSpan);
    }
    job.root.update.push(ir.createBindingOp(job.root.xref, bindingKind, property.name, expression, null, securityContexts, false, false, null, 
    /* TODO: How do Host bindings handle i18n attrs? */ null, property.sourceSpan));
}
export function ingestHostAttribute(job, name, value, securityContexts) {
    const attrBinding = ir.createBindingOp(job.root.xref, ir.BindingKind.Attribute, name, value, null, securityContexts, 
    /* Host attributes should always be extracted to const hostAttrs, even if they are not
     *strictly* text literals */
    true, false, null, 
    /* TODO */ null, 
    /** TODO: May be null? */ value.sourceSpan);
    job.root.update.push(attrBinding);
}
export function ingestHostEvent(job, event) {
    const [phase, target] = event.type !== e.ParsedEventType.Animation
        ? [null, event.targetOrPhase]
        : [event.targetOrPhase, null];
    const eventBinding = ir.createListenerOp(job.root.xref, new ir.SlotHandle(), event.name, null, makeListenerHandlerOps(job.root, event.handler, event.handlerSpan), phase, target, true, event.sourceSpan);
    job.root.create.push(eventBinding);
}
/**
 * Ingest the nodes of a template AST into the given `ViewCompilation`.
 */
function ingestNodes(unit, template) {
    for (const node of template) {
        if (node instanceof t.Element) {
            ingestElement(unit, node);
        }
        else if (node instanceof t.Template) {
            ingestTemplate(unit, node);
        }
        else if (node instanceof t.Content) {
            ingestContent(unit, node);
        }
        else if (node instanceof t.Text) {
            ingestText(unit, node, null);
        }
        else if (node instanceof t.BoundText) {
            ingestBoundText(unit, node, null);
        }
        else if (node instanceof t.IfBlock) {
            ingestIfBlock(unit, node);
        }
        else if (node instanceof t.SwitchBlock) {
            ingestSwitchBlock(unit, node);
        }
        else if (node instanceof t.DeferredBlock) {
            ingestDeferBlock(unit, node);
        }
        else if (node instanceof t.Icu) {
            ingestIcu(unit, node);
        }
        else if (node instanceof t.ForLoopBlock) {
            ingestForBlock(unit, node);
        }
        else if (node instanceof t.LetDeclaration) {
            ingestLetDeclaration(unit, node);
        }
        else {
            throw new Error(`Unsupported template node: ${node.constructor.name}`);
        }
    }
}
/**
 * Ingest an element AST from the template into the given `ViewCompilation`.
 */
function ingestElement(unit, element) {
    if (element.i18n !== undefined &&
        !(element.i18n instanceof i18n.Message || element.i18n instanceof i18n.TagPlaceholder)) {
        throw Error(`Unhandled i18n metadata type for element: ${element.i18n.constructor.name}`);
    }
    const id = unit.job.allocateXrefId();
    const [namespaceKey, elementName] = splitNsName(element.name);
    const startOp = ir.createElementStartOp(elementName, id, namespaceForKey(namespaceKey), element.i18n instanceof i18n.TagPlaceholder ? element.i18n : undefined, element.startSourceSpan, element.sourceSpan);
    unit.create.push(startOp);
    ingestElementBindings(unit, startOp, element);
    ingestReferences(startOp, element);
    // Start i18n, if needed, goes after the element create and bindings, but before the nodes
    let i18nBlockId = null;
    if (element.i18n instanceof i18n.Message) {
        i18nBlockId = unit.job.allocateXrefId();
        unit.create.push(ir.createI18nStartOp(i18nBlockId, element.i18n, undefined, element.startSourceSpan));
    }
    ingestNodes(unit, element.children);
    // The source span for the end op is typically the element closing tag. However, if no closing tag
    // exists, such as in `<input>`, we use the start source span instead. Usually the start and end
    // instructions will be collapsed into one `element` instruction, negating the purpose of this
    // fallback, but in cases when it is not collapsed (such as an input with a binding), we still
    // want to map the end instruction to the main element.
    const endOp = ir.createElementEndOp(id, element.endSourceSpan ?? element.startSourceSpan);
    unit.create.push(endOp);
    // If there is an i18n message associated with this element, insert i18n start and end ops.
    if (i18nBlockId !== null) {
        ir.OpList.insertBefore(ir.createI18nEndOp(i18nBlockId, element.endSourceSpan ?? element.startSourceSpan), endOp);
    }
}
/**
 * Ingest an `ng-template` node from the AST into the given `ViewCompilation`.
 */
function ingestTemplate(unit, tmpl) {
    if (tmpl.i18n !== undefined &&
        !(tmpl.i18n instanceof i18n.Message || tmpl.i18n instanceof i18n.TagPlaceholder)) {
        throw Error(`Unhandled i18n metadata type for template: ${tmpl.i18n.constructor.name}`);
    }
    const childView = unit.job.allocateView(unit.xref);
    let tagNameWithoutNamespace = tmpl.tagName;
    let namespacePrefix = '';
    if (tmpl.tagName) {
        [namespacePrefix, tagNameWithoutNamespace] = splitNsName(tmpl.tagName);
    }
    const i18nPlaceholder = tmpl.i18n instanceof i18n.TagPlaceholder ? tmpl.i18n : undefined;
    const namespace = namespaceForKey(namespacePrefix);
    const functionNameSuffix = tagNameWithoutNamespace === null ? '' : prefixWithNamespace(tagNameWithoutNamespace, namespace);
    const templateKind = isPlainTemplate(tmpl)
        ? ir.TemplateKind.NgTemplate
        : ir.TemplateKind.Structural;
    const templateOp = ir.createTemplateOp(childView.xref, templateKind, tagNameWithoutNamespace, functionNameSuffix, namespace, i18nPlaceholder, tmpl.startSourceSpan, tmpl.sourceSpan);
    unit.create.push(templateOp);
    ingestTemplateBindings(unit, templateOp, tmpl, templateKind);
    ingestReferences(templateOp, tmpl);
    ingestNodes(childView, tmpl.children);
    for (const { name, value } of tmpl.variables) {
        childView.contextVariables.set(name, value !== '' ? value : '$implicit');
    }
    // If this is a plain template and there is an i18n message associated with it, insert i18n start
    // and end ops. For structural directive templates, the i18n ops will be added when ingesting the
    // element/template the directive is placed on.
    if (templateKind === ir.TemplateKind.NgTemplate && tmpl.i18n instanceof i18n.Message) {
        const id = unit.job.allocateXrefId();
        ir.OpList.insertAfter(ir.createI18nStartOp(id, tmpl.i18n, undefined, tmpl.startSourceSpan), childView.create.head);
        ir.OpList.insertBefore(ir.createI18nEndOp(id, tmpl.endSourceSpan ?? tmpl.startSourceSpan), childView.create.tail);
    }
}
/**
 * Ingest a content node from the AST into the given `ViewCompilation`.
 */
function ingestContent(unit, content) {
    if (content.i18n !== undefined && !(content.i18n instanceof i18n.TagPlaceholder)) {
        throw Error(`Unhandled i18n metadata type for element: ${content.i18n.constructor.name}`);
    }
    let fallbackView = null;
    // Don't capture default content that's only made up of empty text nodes and comments.
    // Note that we process the default content before the projection in order to match the
    // insertion order at runtime.
    if (content.children.some((child) => !(child instanceof t.Comment) &&
        (!(child instanceof t.Text) || child.value.trim().length > 0))) {
        fallbackView = unit.job.allocateView(unit.xref);
        ingestNodes(fallbackView, content.children);
    }
    const id = unit.job.allocateXrefId();
    const op = ir.createProjectionOp(id, content.selector, content.i18n, fallbackView?.xref ?? null, content.sourceSpan);
    for (const attr of content.attributes) {
        const securityContext = domSchema.securityContext(content.name, attr.name, true);
        unit.update.push(ir.createBindingOp(op.xref, ir.BindingKind.Attribute, attr.name, o.literal(attr.value), null, securityContext, true, false, null, asMessage(attr.i18n), attr.sourceSpan));
    }
    unit.create.push(op);
}
/**
 * Ingest a literal text node from the AST into the given `ViewCompilation`.
 */
function ingestText(unit, text, icuPlaceholder) {
    unit.create.push(ir.createTextOp(unit.job.allocateXrefId(), text.value, icuPlaceholder, text.sourceSpan));
}
/**
 * Ingest an interpolated text node from the AST into the given `ViewCompilation`.
 */
function ingestBoundText(unit, text, icuPlaceholder) {
    let value = text.value;
    if (value instanceof e.ASTWithSource) {
        value = value.ast;
    }
    if (!(value instanceof e.Interpolation)) {
        throw new Error(`AssertionError: expected Interpolation for BoundText node, got ${value.constructor.name}`);
    }
    if (text.i18n !== undefined && !(text.i18n instanceof i18n.Container)) {
        throw Error(`Unhandled i18n metadata type for text interpolation: ${text.i18n?.constructor.name}`);
    }
    const i18nPlaceholders = text.i18n instanceof i18n.Container
        ? text.i18n.children
            .filter((node) => node instanceof i18n.Placeholder)
            .map((placeholder) => placeholder.name)
        : [];
    if (i18nPlaceholders.length > 0 && i18nPlaceholders.length !== value.expressions.length) {
        throw Error(`Unexpected number of i18n placeholders (${value.expressions.length}) for BoundText with ${value.expressions.length} expressions`);
    }
    const textXref = unit.job.allocateXrefId();
    unit.create.push(ir.createTextOp(textXref, '', icuPlaceholder, text.sourceSpan));
    // TemplateDefinitionBuilder does not generate source maps for sub-expressions inside an
    // interpolation. We copy that behavior in compatibility mode.
    // TODO: is it actually correct to generate these extra maps in modern mode?
    const baseSourceSpan = unit.job.compatibility ? null : text.sourceSpan;
    unit.update.push(ir.createInterpolateTextOp(textXref, new ir.Interpolation(value.strings, value.expressions.map((expr) => convertAst(expr, unit.job, baseSourceSpan)), i18nPlaceholders), text.sourceSpan));
}
/**
 * Ingest an `@if` block into the given `ViewCompilation`.
 */
function ingestIfBlock(unit, ifBlock) {
    let firstXref = null;
    let conditions = [];
    for (let i = 0; i < ifBlock.branches.length; i++) {
        const ifCase = ifBlock.branches[i];
        const cView = unit.job.allocateView(unit.xref);
        const tagName = ingestControlFlowInsertionPoint(unit, cView.xref, ifCase);
        if (ifCase.expressionAlias !== null) {
            cView.contextVariables.set(ifCase.expressionAlias.name, ir.CTX_REF);
        }
        let ifCaseI18nMeta = undefined;
        if (ifCase.i18n !== undefined) {
            if (!(ifCase.i18n instanceof i18n.BlockPlaceholder)) {
                throw Error(`Unhandled i18n metadata type for if block: ${ifCase.i18n?.constructor.name}`);
            }
            ifCaseI18nMeta = ifCase.i18n;
        }
        const templateOp = ir.createTemplateOp(cView.xref, ir.TemplateKind.Block, tagName, 'Conditional', ir.Namespace.HTML, ifCaseI18nMeta, ifCase.startSourceSpan, ifCase.sourceSpan);
        unit.create.push(templateOp);
        if (firstXref === null) {
            firstXref = cView.xref;
        }
        const caseExpr = ifCase.expression ? convertAst(ifCase.expression, unit.job, null) : null;
        const conditionalCaseExpr = new ir.ConditionalCaseExpr(caseExpr, templateOp.xref, templateOp.handle, ifCase.expressionAlias);
        conditions.push(conditionalCaseExpr);
        ingestNodes(cView, ifCase.children);
    }
    unit.update.push(ir.createConditionalOp(firstXref, null, conditions, ifBlock.sourceSpan));
}
/**
 * Ingest an `@switch` block into the given `ViewCompilation`.
 */
function ingestSwitchBlock(unit, switchBlock) {
    // Don't ingest empty switches since they won't render anything.
    if (switchBlock.cases.length === 0) {
        return;
    }
    let firstXref = null;
    let conditions = [];
    for (const switchCase of switchBlock.cases) {
        const cView = unit.job.allocateView(unit.xref);
        const tagName = ingestControlFlowInsertionPoint(unit, cView.xref, switchCase);
        let switchCaseI18nMeta = undefined;
        if (switchCase.i18n !== undefined) {
            if (!(switchCase.i18n instanceof i18n.BlockPlaceholder)) {
                throw Error(`Unhandled i18n metadata type for switch block: ${switchCase.i18n?.constructor.name}`);
            }
            switchCaseI18nMeta = switchCase.i18n;
        }
        const templateOp = ir.createTemplateOp(cView.xref, ir.TemplateKind.Block, tagName, 'Case', ir.Namespace.HTML, switchCaseI18nMeta, switchCase.startSourceSpan, switchCase.sourceSpan);
        unit.create.push(templateOp);
        if (firstXref === null) {
            firstXref = cView.xref;
        }
        const caseExpr = switchCase.expression
            ? convertAst(switchCase.expression, unit.job, switchBlock.startSourceSpan)
            : null;
        const conditionalCaseExpr = new ir.ConditionalCaseExpr(caseExpr, templateOp.xref, templateOp.handle);
        conditions.push(conditionalCaseExpr);
        ingestNodes(cView, switchCase.children);
    }
    unit.update.push(ir.createConditionalOp(firstXref, convertAst(switchBlock.expression, unit.job, null), conditions, switchBlock.sourceSpan));
}
function ingestDeferView(unit, suffix, i18nMeta, children, sourceSpan) {
    if (i18nMeta !== undefined && !(i18nMeta instanceof i18n.BlockPlaceholder)) {
        throw Error('Unhandled i18n metadata type for defer block');
    }
    if (children === undefined) {
        return null;
    }
    const secondaryView = unit.job.allocateView(unit.xref);
    ingestNodes(secondaryView, children);
    const templateOp = ir.createTemplateOp(secondaryView.xref, ir.TemplateKind.Block, null, `Defer${suffix}`, ir.Namespace.HTML, i18nMeta, sourceSpan, sourceSpan);
    unit.create.push(templateOp);
    return templateOp;
}
function ingestDeferBlock(unit, deferBlock) {
    let ownResolverFn = null;
    if (unit.job.deferMeta.mode === 0 /* DeferBlockDepsEmitMode.PerBlock */) {
        if (!unit.job.deferMeta.blocks.has(deferBlock)) {
            throw new Error(`AssertionError: unable to find a dependency function for this deferred block`);
        }
        ownResolverFn = unit.job.deferMeta.blocks.get(deferBlock) ?? null;
    }
    // Generate the defer main view and all secondary views.
    const main = ingestDeferView(unit, '', deferBlock.i18n, deferBlock.children, deferBlock.sourceSpan);
    const loading = ingestDeferView(unit, 'Loading', deferBlock.loading?.i18n, deferBlock.loading?.children, deferBlock.loading?.sourceSpan);
    const placeholder = ingestDeferView(unit, 'Placeholder', deferBlock.placeholder?.i18n, deferBlock.placeholder?.children, deferBlock.placeholder?.sourceSpan);
    const error = ingestDeferView(unit, 'Error', deferBlock.error?.i18n, deferBlock.error?.children, deferBlock.error?.sourceSpan);
    // Create the main defer op, and ops for all secondary views.
    const deferXref = unit.job.allocateXrefId();
    const deferOp = ir.createDeferOp(deferXref, main.xref, main.handle, ownResolverFn, unit.job.allDeferrableDepsFn, deferBlock.sourceSpan);
    deferOp.placeholderView = placeholder?.xref ?? null;
    deferOp.placeholderSlot = placeholder?.handle ?? null;
    deferOp.loadingSlot = loading?.handle ?? null;
    deferOp.errorSlot = error?.handle ?? null;
    deferOp.placeholderMinimumTime = deferBlock.placeholder?.minimumTime ?? null;
    deferOp.loadingMinimumTime = deferBlock.loading?.minimumTime ?? null;
    deferOp.loadingAfterTime = deferBlock.loading?.afterTime ?? null;
    unit.create.push(deferOp);
    // Configure all defer `on` conditions.
    // TODO: refactor prefetch triggers to use a separate op type, with a shared superclass. This will
    // make it easier to refactor prefetch behavior in the future.
    let prefetch = false;
    let deferOnOps = [];
    let deferWhenOps = [];
    for (const triggers of [deferBlock.triggers, deferBlock.prefetchTriggers]) {
        if (triggers.idle !== undefined) {
            const deferOnOp = ir.createDeferOnOp(deferXref, { kind: ir.DeferTriggerKind.Idle }, prefetch, triggers.idle.sourceSpan);
            deferOnOps.push(deferOnOp);
        }
        if (triggers.immediate !== undefined) {
            const deferOnOp = ir.createDeferOnOp(deferXref, { kind: ir.DeferTriggerKind.Immediate }, prefetch, triggers.immediate.sourceSpan);
            deferOnOps.push(deferOnOp);
        }
        if (triggers.timer !== undefined) {
            const deferOnOp = ir.createDeferOnOp(deferXref, { kind: ir.DeferTriggerKind.Timer, delay: triggers.timer.delay }, prefetch, triggers.timer.sourceSpan);
            deferOnOps.push(deferOnOp);
        }
        if (triggers.hover !== undefined) {
            const deferOnOp = ir.createDeferOnOp(deferXref, {
                kind: ir.DeferTriggerKind.Hover,
                targetName: triggers.hover.reference,
                targetXref: null,
                targetSlot: null,
                targetView: null,
                targetSlotViewSteps: null,
            }, prefetch, triggers.hover.sourceSpan);
            deferOnOps.push(deferOnOp);
        }
        if (triggers.interaction !== undefined) {
            const deferOnOp = ir.createDeferOnOp(deferXref, {
                kind: ir.DeferTriggerKind.Interaction,
                targetName: triggers.interaction.reference,
                targetXref: null,
                targetSlot: null,
                targetView: null,
                targetSlotViewSteps: null,
            }, prefetch, triggers.interaction.sourceSpan);
            deferOnOps.push(deferOnOp);
        }
        if (triggers.viewport !== undefined) {
            const deferOnOp = ir.createDeferOnOp(deferXref, {
                kind: ir.DeferTriggerKind.Viewport,
                targetName: triggers.viewport.reference,
                targetXref: null,
                targetSlot: null,
                targetView: null,
                targetSlotViewSteps: null,
            }, prefetch, triggers.viewport.sourceSpan);
            deferOnOps.push(deferOnOp);
        }
        if (triggers.when !== undefined) {
            if (triggers.when.value instanceof e.Interpolation) {
                // TemplateDefinitionBuilder supports this case, but it's very strange to me. What would it
                // even mean?
                throw new Error(`Unexpected interpolation in defer block when trigger`);
            }
            const deferOnOp = ir.createDeferWhenOp(deferXref, convertAst(triggers.when.value, unit.job, triggers.when.sourceSpan), prefetch, triggers.when.sourceSpan);
            deferWhenOps.push(deferOnOp);
        }
        // If no (non-prefetching) defer triggers were provided, default to `idle`.
        if (deferOnOps.length === 0 && deferWhenOps.length === 0) {
            deferOnOps.push(ir.createDeferOnOp(deferXref, { kind: ir.DeferTriggerKind.Idle }, false, null));
        }
        prefetch = true;
    }
    unit.create.push(deferOnOps);
    unit.update.push(deferWhenOps);
}
function ingestIcu(unit, icu) {
    if (icu.i18n instanceof i18n.Message && isSingleI18nIcu(icu.i18n)) {
        const xref = unit.job.allocateXrefId();
        unit.create.push(ir.createIcuStartOp(xref, icu.i18n, icuFromI18nMessage(icu.i18n).name, null));
        for (const [placeholder, text] of Object.entries({ ...icu.vars, ...icu.placeholders })) {
            if (text instanceof t.BoundText) {
                ingestBoundText(unit, text, placeholder);
            }
            else {
                ingestText(unit, text, placeholder);
            }
        }
        unit.create.push(ir.createIcuEndOp(xref));
    }
    else {
        throw Error(`Unhandled i18n metadata type for ICU: ${icu.i18n?.constructor.name}`);
    }
}
/**
 * Ingest an `@for` block into the given `ViewCompilation`.
 */
function ingestForBlock(unit, forBlock) {
    const repeaterView = unit.job.allocateView(unit.xref);
    // We copy TemplateDefinitionBuilder's scheme of creating names for `$count` and `$index`
    // that are suffixed with special information, to disambiguate which level of nested loop
    // the below aliases refer to.
    // TODO: We should refactor Template Pipeline's variable phases to gracefully handle
    // shadowing, and arbitrarily many levels of variables depending on each other.
    const indexName = `ɵ$index_${repeaterView.xref}`;
    const countName = `ɵ$count_${repeaterView.xref}`;
    const indexVarNames = new Set();
    // Set all the context variables and aliases available in the repeater.
    repeaterView.contextVariables.set(forBlock.item.name, forBlock.item.value);
    for (const variable of forBlock.contextVariables) {
        if (variable.value === '$index') {
            indexVarNames.add(variable.name);
        }
        if (variable.name === '$index') {
            repeaterView.contextVariables.set('$index', variable.value).set(indexName, variable.value);
        }
        else if (variable.name === '$count') {
            repeaterView.contextVariables.set('$count', variable.value).set(countName, variable.value);
        }
        else {
            repeaterView.aliases.add({
                kind: ir.SemanticVariableKind.Alias,
                name: null,
                identifier: variable.name,
                expression: getComputedForLoopVariableExpression(variable, indexName, countName),
            });
        }
    }
    const sourceSpan = convertSourceSpan(forBlock.trackBy.span, forBlock.sourceSpan);
    const track = convertAst(forBlock.trackBy, unit.job, sourceSpan);
    ingestNodes(repeaterView, forBlock.children);
    let emptyView = null;
    let emptyTagName = null;
    if (forBlock.empty !== null) {
        emptyView = unit.job.allocateView(unit.xref);
        ingestNodes(emptyView, forBlock.empty.children);
        emptyTagName = ingestControlFlowInsertionPoint(unit, emptyView.xref, forBlock.empty);
    }
    const varNames = {
        $index: indexVarNames,
        $implicit: forBlock.item.name,
    };
    if (forBlock.i18n !== undefined && !(forBlock.i18n instanceof i18n.BlockPlaceholder)) {
        throw Error('AssertionError: Unhandled i18n metadata type or @for');
    }
    if (forBlock.empty?.i18n !== undefined &&
        !(forBlock.empty.i18n instanceof i18n.BlockPlaceholder)) {
        throw Error('AssertionError: Unhandled i18n metadata type or @empty');
    }
    const i18nPlaceholder = forBlock.i18n;
    const emptyI18nPlaceholder = forBlock.empty?.i18n;
    const tagName = ingestControlFlowInsertionPoint(unit, repeaterView.xref, forBlock);
    const repeaterCreate = ir.createRepeaterCreateOp(repeaterView.xref, emptyView?.xref ?? null, tagName, track, varNames, emptyTagName, i18nPlaceholder, emptyI18nPlaceholder, forBlock.startSourceSpan, forBlock.sourceSpan);
    unit.create.push(repeaterCreate);
    const expression = convertAst(forBlock.expression, unit.job, convertSourceSpan(forBlock.expression.span, forBlock.sourceSpan));
    const repeater = ir.createRepeaterOp(repeaterCreate.xref, repeaterCreate.handle, expression, forBlock.sourceSpan);
    unit.update.push(repeater);
}
/**
 * Gets an expression that represents a variable in an `@for` loop.
 * @param variable AST representing the variable.
 * @param indexName Loop-specific name for `$index`.
 * @param countName Loop-specific name for `$count`.
 */
function getComputedForLoopVariableExpression(variable, indexName, countName) {
    switch (variable.value) {
        case '$index':
            return new ir.LexicalReadExpr(indexName);
        case '$count':
            return new ir.LexicalReadExpr(countName);
        case '$first':
            return new ir.LexicalReadExpr(indexName).identical(o.literal(0));
        case '$last':
            return new ir.LexicalReadExpr(indexName).identical(new ir.LexicalReadExpr(countName).minus(o.literal(1)));
        case '$even':
            return new ir.LexicalReadExpr(indexName).modulo(o.literal(2)).identical(o.literal(0));
        case '$odd':
            return new ir.LexicalReadExpr(indexName).modulo(o.literal(2)).notIdentical(o.literal(0));
        default:
            throw new Error(`AssertionError: unknown @for loop variable ${variable.value}`);
    }
}
function ingestLetDeclaration(unit, node) {
    const target = unit.job.allocateXrefId();
    unit.create.push(ir.createDeclareLetOp(target, node.name, node.sourceSpan));
    unit.update.push(ir.createStoreLetOp(target, node.name, convertAst(node.value, unit.job, node.valueSpan), node.sourceSpan));
}
/**
 * Convert a template AST expression into an output AST expression.
 */
function convertAst(ast, job, baseSourceSpan) {
    if (ast instanceof e.ASTWithSource) {
        return convertAst(ast.ast, job, baseSourceSpan);
    }
    else if (ast instanceof e.PropertyRead) {
        const isThisReceiver = ast.receiver instanceof e.ThisReceiver;
        // Whether this is an implicit receiver, *excluding* explicit reads of `this`.
        const isImplicitReceiver = ast.receiver instanceof e.ImplicitReceiver && !(ast.receiver instanceof e.ThisReceiver);
        // Whether the  name of the read is a node that should be never retain its explicit this
        // receiver.
        const isSpecialNode = ast.name === '$any' || ast.name === '$event';
        // TODO: The most sensible condition here would be simply `isImplicitReceiver`, to convert only
        // actual implicit `this` reads, and not explicit ones. However, TemplateDefinitionBuilder (and
        // the Typecheck block!) both have the same bug, in which they also consider explicit `this`
        // reads to be implicit. This causes problems when the explicit `this` read is inside a
        // template with a context that also provides the variable name being read:
        // ```
        // <ng-template let-a>{{this.a}}</ng-template>
        // ```
        // The whole point of the explicit `this` was to access the class property, but TDB and the
        // current TCB treat the read as implicit, and give you the context property instead!
        //
        // For now, we emulate this old behavior by aggressively converting explicit reads to to
        // implicit reads, except for the special cases that TDB and the current TCB protect. However,
        // it would be an improvement to fix this.
        //
        // See also the corresponding comment for the TCB, in `type_check_block.ts`.
        if (isImplicitReceiver || (isThisReceiver && !isSpecialNode)) {
            return new ir.LexicalReadExpr(ast.name);
        }
        else {
            return new o.ReadPropExpr(convertAst(ast.receiver, job, baseSourceSpan), ast.name, null, convertSourceSpan(ast.span, baseSourceSpan));
        }
    }
    else if (ast instanceof e.PropertyWrite) {
        if (ast.receiver instanceof e.ImplicitReceiver) {
            return new o.WritePropExpr(
            // TODO: Is it correct to always use the root context in place of the implicit receiver?
            new ir.ContextExpr(job.root.xref), ast.name, convertAst(ast.value, job, baseSourceSpan), null, convertSourceSpan(ast.span, baseSourceSpan));
        }
        return new o.WritePropExpr(convertAst(ast.receiver, job, baseSourceSpan), ast.name, convertAst(ast.value, job, baseSourceSpan), undefined, convertSourceSpan(ast.span, baseSourceSpan));
    }
    else if (ast instanceof e.KeyedWrite) {
        return new o.WriteKeyExpr(convertAst(ast.receiver, job, baseSourceSpan), convertAst(ast.key, job, baseSourceSpan), convertAst(ast.value, job, baseSourceSpan), undefined, convertSourceSpan(ast.span, baseSourceSpan));
    }
    else if (ast instanceof e.Call) {
        if (ast.receiver instanceof e.ImplicitReceiver) {
            throw new Error(`Unexpected ImplicitReceiver`);
        }
        else {
            return new o.InvokeFunctionExpr(convertAst(ast.receiver, job, baseSourceSpan), ast.args.map((arg) => convertAst(arg, job, baseSourceSpan)), undefined, convertSourceSpan(ast.span, baseSourceSpan));
        }
    }
    else if (ast instanceof e.LiteralPrimitive) {
        return o.literal(ast.value, undefined, convertSourceSpan(ast.span, baseSourceSpan));
    }
    else if (ast instanceof e.Unary) {
        switch (ast.operator) {
            case '+':
                return new o.UnaryOperatorExpr(o.UnaryOperator.Plus, convertAst(ast.expr, job, baseSourceSpan), undefined, convertSourceSpan(ast.span, baseSourceSpan));
            case '-':
                return new o.UnaryOperatorExpr(o.UnaryOperator.Minus, convertAst(ast.expr, job, baseSourceSpan), undefined, convertSourceSpan(ast.span, baseSourceSpan));
            default:
                throw new Error(`AssertionError: unknown unary operator ${ast.operator}`);
        }
    }
    else if (ast instanceof e.Binary) {
        const operator = BINARY_OPERATORS.get(ast.operation);
        if (operator === undefined) {
            throw new Error(`AssertionError: unknown binary operator ${ast.operation}`);
        }
        return new o.BinaryOperatorExpr(operator, convertAst(ast.left, job, baseSourceSpan), convertAst(ast.right, job, baseSourceSpan), undefined, convertSourceSpan(ast.span, baseSourceSpan));
    }
    else if (ast instanceof e.ThisReceiver) {
        // TODO: should context expressions have source maps?
        return new ir.ContextExpr(job.root.xref);
    }
    else if (ast instanceof e.KeyedRead) {
        return new o.ReadKeyExpr(convertAst(ast.receiver, job, baseSourceSpan), convertAst(ast.key, job, baseSourceSpan), undefined, convertSourceSpan(ast.span, baseSourceSpan));
    }
    else if (ast instanceof e.Chain) {
        throw new Error(`AssertionError: Chain in unknown context`);
    }
    else if (ast instanceof e.LiteralMap) {
        const entries = ast.keys.map((key, idx) => {
            const value = ast.values[idx];
            // TODO: should literals have source maps, or do we just map the whole surrounding
            // expression?
            return new o.LiteralMapEntry(key.key, convertAst(value, job, baseSourceSpan), key.quoted);
        });
        return new o.LiteralMapExpr(entries, undefined, convertSourceSpan(ast.span, baseSourceSpan));
    }
    else if (ast instanceof e.LiteralArray) {
        // TODO: should literals have source maps, or do we just map the whole surrounding expression?
        return new o.LiteralArrayExpr(ast.expressions.map((expr) => convertAst(expr, job, baseSourceSpan)));
    }
    else if (ast instanceof e.Conditional) {
        return new o.ConditionalExpr(convertAst(ast.condition, job, baseSourceSpan), convertAst(ast.trueExp, job, baseSourceSpan), convertAst(ast.falseExp, job, baseSourceSpan), undefined, convertSourceSpan(ast.span, baseSourceSpan));
    }
    else if (ast instanceof e.NonNullAssert) {
        // A non-null assertion shouldn't impact generated instructions, so we can just drop it.
        return convertAst(ast.expression, job, baseSourceSpan);
    }
    else if (ast instanceof e.BindingPipe) {
        // TODO: pipes should probably have source maps; figure out details.
        return new ir.PipeBindingExpr(job.allocateXrefId(), new ir.SlotHandle(), ast.name, [
            convertAst(ast.exp, job, baseSourceSpan),
            ...ast.args.map((arg) => convertAst(arg, job, baseSourceSpan)),
        ]);
    }
    else if (ast instanceof e.SafeKeyedRead) {
        return new ir.SafeKeyedReadExpr(convertAst(ast.receiver, job, baseSourceSpan), convertAst(ast.key, job, baseSourceSpan), convertSourceSpan(ast.span, baseSourceSpan));
    }
    else if (ast instanceof e.SafePropertyRead) {
        // TODO: source span
        return new ir.SafePropertyReadExpr(convertAst(ast.receiver, job, baseSourceSpan), ast.name);
    }
    else if (ast instanceof e.SafeCall) {
        // TODO: source span
        return new ir.SafeInvokeFunctionExpr(convertAst(ast.receiver, job, baseSourceSpan), ast.args.map((a) => convertAst(a, job, baseSourceSpan)));
    }
    else if (ast instanceof e.EmptyExpr) {
        return new ir.EmptyExpr(convertSourceSpan(ast.span, baseSourceSpan));
    }
    else if (ast instanceof e.PrefixNot) {
        return o.not(convertAst(ast.expression, job, baseSourceSpan), convertSourceSpan(ast.span, baseSourceSpan));
    }
    else {
        throw new Error(`Unhandled expression type "${ast.constructor.name}" in file "${baseSourceSpan?.start.file.url}"`);
    }
}
function convertAstWithInterpolation(job, value, i18nMeta, sourceSpan) {
    let expression;
    if (value instanceof e.Interpolation) {
        expression = new ir.Interpolation(value.strings, value.expressions.map((e) => convertAst(e, job, sourceSpan ?? null)), Object.keys(asMessage(i18nMeta)?.placeholders ?? {}));
    }
    else if (value instanceof e.AST) {
        expression = convertAst(value, job, sourceSpan ?? null);
    }
    else {
        expression = o.literal(value);
    }
    return expression;
}
// TODO: Can we populate Template binding kinds in ingest?
const BINDING_KINDS = new Map([
    [e.BindingType.Property, ir.BindingKind.Property],
    [e.BindingType.TwoWay, ir.BindingKind.TwoWayProperty],
    [e.BindingType.Attribute, ir.BindingKind.Attribute],
    [e.BindingType.Class, ir.BindingKind.ClassName],
    [e.BindingType.Style, ir.BindingKind.StyleProperty],
    [e.BindingType.Animation, ir.BindingKind.Animation],
]);
/**
 * Checks whether the given template is a plain ng-template (as opposed to another kind of template
 * such as a structural directive template or control flow template). This is checked based on the
 * tagName. We can expect that only plain ng-templates will come through with a tagName of
 * 'ng-template'.
 *
 * Here are some of the cases we expect:
 *
 * | Angular HTML                       | Template tagName   |
 * | ---------------------------------- | ------------------ |
 * | `<ng-template>`                    | 'ng-template'      |
 * | `<div *ngIf="true">`               | 'div'              |
 * | `<svg><ng-template>`               | 'svg:ng-template'  |
 * | `@if (true) {`                     | 'Conditional'      |
 * | `<ng-template *ngIf>` (plain)      | 'ng-template'      |
 * | `<ng-template *ngIf>` (structural) | null               |
 */
function isPlainTemplate(tmpl) {
    return splitNsName(tmpl.tagName ?? '')[1] === NG_TEMPLATE_TAG_NAME;
}
/**
 * Ensures that the i18nMeta, if provided, is an i18n.Message.
 */
function asMessage(i18nMeta) {
    if (i18nMeta == null) {
        return null;
    }
    if (!(i18nMeta instanceof i18n.Message)) {
        throw Error(`Expected i18n meta to be a Message, but got: ${i18nMeta.constructor.name}`);
    }
    return i18nMeta;
}
/**
 * Process all of the bindings on an element in the template AST and convert them to their IR
 * representation.
 */
function ingestElementBindings(unit, op, element) {
    let bindings = new Array();
    let i18nAttributeBindingNames = new Set();
    for (const attr of element.attributes) {
        // Attribute literal bindings, such as `attr.foo="bar"`.
        const securityContext = domSchema.securityContext(element.name, attr.name, true);
        bindings.push(ir.createBindingOp(op.xref, ir.BindingKind.Attribute, attr.name, convertAstWithInterpolation(unit.job, attr.value, attr.i18n), null, securityContext, true, false, null, asMessage(attr.i18n), attr.sourceSpan));
        if (attr.i18n) {
            i18nAttributeBindingNames.add(attr.name);
        }
    }
    for (const input of element.inputs) {
        if (i18nAttributeBindingNames.has(input.name)) {
            console.error(`On component ${unit.job.componentName}, the binding ${input.name} is both an i18n attribute and a property. You may want to remove the property binding. This will become a compilation error in future versions of Angular.`);
        }
        // All dynamic bindings (both attribute and property bindings).
        bindings.push(ir.createBindingOp(op.xref, BINDING_KINDS.get(input.type), input.name, convertAstWithInterpolation(unit.job, astOf(input.value), input.i18n), input.unit, input.securityContext, false, false, null, asMessage(input.i18n) ?? null, input.sourceSpan));
    }
    unit.create.push(bindings.filter((b) => b?.kind === ir.OpKind.ExtractedAttribute));
    unit.update.push(bindings.filter((b) => b?.kind === ir.OpKind.Binding));
    for (const output of element.outputs) {
        if (output.type === e.ParsedEventType.Animation && output.phase === null) {
            throw Error('Animation listener should have a phase');
        }
        if (output.type === e.ParsedEventType.TwoWay) {
            unit.create.push(ir.createTwoWayListenerOp(op.xref, op.handle, output.name, op.tag, makeTwoWayListenerHandlerOps(unit, output.handler, output.handlerSpan), output.sourceSpan));
        }
        else {
            unit.create.push(ir.createListenerOp(op.xref, op.handle, output.name, op.tag, makeListenerHandlerOps(unit, output.handler, output.handlerSpan), output.phase, output.target, false, output.sourceSpan));
        }
    }
    // If any of the bindings on this element have an i18n message, then an i18n attrs configuration
    // op is also required.
    if (bindings.some((b) => b?.i18nMessage) !== null) {
        unit.create.push(ir.createI18nAttributesOp(unit.job.allocateXrefId(), new ir.SlotHandle(), op.xref));
    }
}
/**
 * Process all of the bindings on a template in the template AST and convert them to their IR
 * representation.
 */
function ingestTemplateBindings(unit, op, template, templateKind) {
    let bindings = new Array();
    for (const attr of template.templateAttrs) {
        if (attr instanceof t.TextAttribute) {
            const securityContext = domSchema.securityContext(NG_TEMPLATE_TAG_NAME, attr.name, true);
            bindings.push(createTemplateBinding(unit, op.xref, e.BindingType.Attribute, attr.name, attr.value, null, securityContext, true, templateKind, asMessage(attr.i18n), attr.sourceSpan));
        }
        else {
            bindings.push(createTemplateBinding(unit, op.xref, attr.type, attr.name, astOf(attr.value), attr.unit, attr.securityContext, true, templateKind, asMessage(attr.i18n), attr.sourceSpan));
        }
    }
    for (const attr of template.attributes) {
        // Attribute literal bindings, such as `attr.foo="bar"`.
        const securityContext = domSchema.securityContext(NG_TEMPLATE_TAG_NAME, attr.name, true);
        bindings.push(createTemplateBinding(unit, op.xref, e.BindingType.Attribute, attr.name, attr.value, null, securityContext, false, templateKind, asMessage(attr.i18n), attr.sourceSpan));
    }
    for (const input of template.inputs) {
        // Dynamic bindings (both attribute and property bindings).
        bindings.push(createTemplateBinding(unit, op.xref, input.type, input.name, astOf(input.value), input.unit, input.securityContext, false, templateKind, asMessage(input.i18n), input.sourceSpan));
    }
    unit.create.push(bindings.filter((b) => b?.kind === ir.OpKind.ExtractedAttribute));
    unit.update.push(bindings.filter((b) => b?.kind === ir.OpKind.Binding));
    for (const output of template.outputs) {
        if (output.type === e.ParsedEventType.Animation && output.phase === null) {
            throw Error('Animation listener should have a phase');
        }
        if (templateKind === ir.TemplateKind.NgTemplate) {
            if (output.type === e.ParsedEventType.TwoWay) {
                unit.create.push(ir.createTwoWayListenerOp(op.xref, op.handle, output.name, op.tag, makeTwoWayListenerHandlerOps(unit, output.handler, output.handlerSpan), output.sourceSpan));
            }
            else {
                unit.create.push(ir.createListenerOp(op.xref, op.handle, output.name, op.tag, makeListenerHandlerOps(unit, output.handler, output.handlerSpan), output.phase, output.target, false, output.sourceSpan));
            }
        }
        if (templateKind === ir.TemplateKind.Structural &&
            output.type !== e.ParsedEventType.Animation) {
            // Animation bindings are excluded from the structural template's const array.
            const securityContext = domSchema.securityContext(NG_TEMPLATE_TAG_NAME, output.name, false);
            unit.create.push(ir.createExtractedAttributeOp(op.xref, ir.BindingKind.Property, null, output.name, null, null, null, securityContext));
        }
    }
    // TODO: Perhaps we could do this in a phase? (It likely wouldn't change the slot indices.)
    if (bindings.some((b) => b?.i18nMessage) !== null) {
        unit.create.push(ir.createI18nAttributesOp(unit.job.allocateXrefId(), new ir.SlotHandle(), op.xref));
    }
}
/**
 * Helper to ingest an individual binding on a template, either an explicit `ng-template`, or an
 * implicit template created via structural directive.
 *
 * Bindings on templates are *extremely* tricky. I have tried to isolate all of the confusing edge
 * cases into this function, and to comment it well to document the behavior.
 *
 * Some of this behavior is intuitively incorrect, and we should consider changing it in the future.
 *
 * @param view The compilation unit for the view containing the template.
 * @param xref The xref of the template op.
 * @param type The binding type, according to the parser. This is fairly reasonable, e.g. both
 *     dynamic and static attributes have e.BindingType.Attribute.
 * @param name The binding's name.
 * @param value The bindings's value, which will either be an input AST expression, or a string
 *     literal. Note that the input AST expression may or may not be const -- it will only be a
 *     string literal if the parser considered it a text binding.
 * @param unit If the binding has a unit (e.g. `px` for style bindings), then this is the unit.
 * @param securityContext The security context of the binding.
 * @param isStructuralTemplateAttribute Whether this binding actually applies to the structural
 *     ng-template. For example, an `ngFor` would actually apply to the structural template. (Most
 *     bindings on structural elements target the inner element, not the template.)
 * @param templateKind Whether this is an explicit `ng-template` or an implicit template created by
 *     a structural directive. This should never be a block template.
 * @param i18nMessage The i18n metadata for the binding, if any.
 * @param sourceSpan The source span of the binding.
 * @returns An IR binding op, or null if the binding should be skipped.
 */
function createTemplateBinding(view, xref, type, name, value, unit, securityContext, isStructuralTemplateAttribute, templateKind, i18nMessage, sourceSpan) {
    const isTextBinding = typeof value === 'string';
    // If this is a structural template, then several kinds of bindings should not result in an
    // update instruction.
    if (templateKind === ir.TemplateKind.Structural) {
        if (!isStructuralTemplateAttribute) {
            switch (type) {
                case e.BindingType.Property:
                case e.BindingType.Class:
                case e.BindingType.Style:
                    // Because this binding doesn't really target the ng-template, it must be a binding on an
                    // inner node of a structural template. We can't skip it entirely, because we still need
                    // it on the ng-template's consts (e.g. for the purposes of directive matching). However,
                    // we should not generate an update instruction for it.
                    return ir.createExtractedAttributeOp(xref, ir.BindingKind.Property, null, name, null, null, i18nMessage, securityContext);
                case e.BindingType.TwoWay:
                    return ir.createExtractedAttributeOp(xref, ir.BindingKind.TwoWayProperty, null, name, null, null, i18nMessage, securityContext);
            }
        }
        if (!isTextBinding && (type === e.BindingType.Attribute || type === e.BindingType.Animation)) {
            // Again, this binding doesn't really target the ng-template; it actually targets the element
            // inside the structural template. In the case of non-text attribute or animation bindings,
            // the binding doesn't even show up on the ng-template const array, so we just skip it
            // entirely.
            return null;
        }
    }
    let bindingType = BINDING_KINDS.get(type);
    if (templateKind === ir.TemplateKind.NgTemplate) {
        // We know we are dealing with bindings directly on an explicit ng-template.
        // Static attribute bindings should be collected into the const array as k/v pairs. Property
        // bindings should result in a `property` instruction, and `AttributeMarker.Bindings` const
        // entries.
        //
        // The difficulty is with dynamic attribute, style, and class bindings. These don't really make
        // sense on an `ng-template` and should probably be parser errors. However,
        // TemplateDefinitionBuilder generates `property` instructions for them, and so we do that as
        // well.
        //
        // Note that we do have a slight behavior difference with TemplateDefinitionBuilder: although
        // TDB emits `property` instructions for dynamic attributes, styles, and classes, only styles
        // and classes also get const collected into the `AttributeMarker.Bindings` field. Dynamic
        // attribute bindings are missing from the consts entirely. We choose to emit them into the
        // consts field anyway, to avoid creating special cases for something so arcane and nonsensical.
        if (type === e.BindingType.Class ||
            type === e.BindingType.Style ||
            (type === e.BindingType.Attribute && !isTextBinding)) {
            // TODO: These cases should be parse errors.
            bindingType = ir.BindingKind.Property;
        }
    }
    return ir.createBindingOp(xref, bindingType, name, convertAstWithInterpolation(view.job, value, i18nMessage), unit, securityContext, isTextBinding, isStructuralTemplateAttribute, templateKind, i18nMessage, sourceSpan);
}
function makeListenerHandlerOps(unit, handler, handlerSpan) {
    handler = astOf(handler);
    const handlerOps = new Array();
    let handlerExprs = handler instanceof e.Chain ? handler.expressions : [handler];
    if (handlerExprs.length === 0) {
        throw new Error('Expected listener to have non-empty expression list.');
    }
    const expressions = handlerExprs.map((expr) => convertAst(expr, unit.job, handlerSpan));
    const returnExpr = expressions.pop();
    handlerOps.push(...expressions.map((e) => ir.createStatementOp(new o.ExpressionStatement(e, e.sourceSpan))));
    handlerOps.push(ir.createStatementOp(new o.ReturnStatement(returnExpr, returnExpr.sourceSpan)));
    return handlerOps;
}
function makeTwoWayListenerHandlerOps(unit, handler, handlerSpan) {
    handler = astOf(handler);
    const handlerOps = new Array();
    if (handler instanceof e.Chain) {
        if (handler.expressions.length === 1) {
            handler = handler.expressions[0];
        }
        else {
            // This is validated during parsing already, but we do it here just in case.
            throw new Error('Expected two-way listener to have a single expression.');
        }
    }
    const handlerExpr = convertAst(handler, unit.job, handlerSpan);
    const eventReference = new ir.LexicalReadExpr('$event');
    const twoWaySetExpr = new ir.TwoWayBindingSetExpr(handlerExpr, eventReference);
    handlerOps.push(ir.createStatementOp(new o.ExpressionStatement(twoWaySetExpr)));
    handlerOps.push(ir.createStatementOp(new o.ReturnStatement(eventReference)));
    return handlerOps;
}
function astOf(ast) {
    return ast instanceof e.ASTWithSource ? ast.ast : ast;
}
/**
 * Process all of the local references on an element-like structure in the template AST and
 * convert them to their IR representation.
 */
function ingestReferences(op, element) {
    assertIsArray(op.localRefs);
    for (const { name, value } of element.references) {
        op.localRefs.push({
            name,
            target: value,
        });
    }
}
/**
 * Assert that the given value is an array.
 */
function assertIsArray(value) {
    if (!Array.isArray(value)) {
        throw new Error(`AssertionError: expected an array`);
    }
}
/**
 * Creates an absolute `ParseSourceSpan` from the relative `ParseSpan`.
 *
 * `ParseSpan` objects are relative to the start of the expression.
 * This method converts these to full `ParseSourceSpan` objects that
 * show where the span is within the overall source file.
 *
 * @param span the relative span to convert.
 * @param baseSourceSpan a span corresponding to the base of the expression tree.
 * @returns a `ParseSourceSpan` for the given span or null if no `baseSourceSpan` was provided.
 */
function convertSourceSpan(span, baseSourceSpan) {
    if (baseSourceSpan === null) {
        return null;
    }
    const start = baseSourceSpan.start.moveBy(span.start);
    const end = baseSourceSpan.start.moveBy(span.end);
    const fullStart = baseSourceSpan.fullStart.moveBy(span.start);
    return new ParseSourceSpan(start, end, fullStart);
}
/**
 * With the directive-based control flow users were able to conditionally project content using
 * the `*` syntax. E.g. `<div *ngIf="expr" projectMe></div>` will be projected into
 * `<ng-content select="[projectMe]"/>`, because the attributes and tag name from the `div` are
 * copied to the template via the template creation instruction. With `@if` and `@for` that is
 * not the case, because the conditional is placed *around* elements, rather than *on* them.
 * The result is that content projection won't work in the same way if a user converts from
 * `*ngIf` to `@if`.
 *
 * This function aims to cover the most common case by doing the same copying when a control flow
 * node has *one and only one* root element or template node.
 *
 * This approach comes with some caveats:
 * 1. As soon as any other node is added to the root, the copying behavior won't work anymore.
 *    A diagnostic will be added to flag cases like this and to explain how to work around it.
 * 2. If `preserveWhitespaces` is enabled, it's very likely that indentation will break this
 *    workaround, because it'll include an additional text node as the first child. We can work
 *    around it here, but in a discussion it was decided not to, because the user explicitly opted
 *    into preserving the whitespace and we would have to drop it from the generated code.
 *    The diagnostic mentioned point #1 will flag such cases to users.
 *
 * @returns Tag name to be used for the control flow template.
 */
function ingestControlFlowInsertionPoint(unit, xref, node) {
    let root = null;
    for (const child of node.children) {
        // Skip over comment nodes.
        if (child instanceof t.Comment) {
            continue;
        }
        // We can only infer the tag name/attributes if there's a single root node.
        if (root !== null) {
            return null;
        }
        // Root nodes can only elements or templates with a tag name (e.g. `<div *foo></div>`).
        if (child instanceof t.Element || (child instanceof t.Template && child.tagName !== null)) {
            root = child;
        }
    }
    // If we've found a single root node, its tag name and attributes can be
    // copied to the surrounding template to be used for content projection.
    if (root !== null) {
        // Collect the static attributes for content projection purposes.
        for (const attr of root.attributes) {
            const securityContext = domSchema.securityContext(NG_TEMPLATE_TAG_NAME, attr.name, true);
            unit.update.push(ir.createBindingOp(xref, ir.BindingKind.Attribute, attr.name, o.literal(attr.value), null, securityContext, true, false, null, asMessage(attr.i18n), attr.sourceSpan));
        }
        // Also collect the inputs since they participate in content projection as well.
        // Note that TDB used to collect the outputs as well, but it wasn't passing them into
        // the template instruction. Here we just don't collect them.
        for (const attr of root.inputs) {
            if (attr.type !== e.BindingType.Animation && attr.type !== e.BindingType.Attribute) {
                const securityContext = domSchema.securityContext(NG_TEMPLATE_TAG_NAME, attr.name, true);
                unit.create.push(ir.createExtractedAttributeOp(xref, ir.BindingKind.Property, null, attr.name, null, null, null, securityContext));
            }
        }
        const tagName = root instanceof t.Element ? root.name : root.tagName;
        // Don't pass along `ng-template` tag name since it enables directive matching.
        return tagName === NG_TEMPLATE_TAG_NAME ? null : tagName;
    }
    return null;
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaW5nZXN0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vcGFja2FnZXMvY29tcGlsZXIvc3JjL3RlbXBsYXRlL3BpcGVsaW5lL3NyYy9pbmdlc3QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IkFBQUE7Ozs7OztHQU1HO0FBR0gsT0FBTyxFQUFDLGVBQWUsRUFBQyxNQUFNLGVBQWUsQ0FBQztBQUM5QyxPQUFPLEtBQUssQ0FBQyxNQUFNLGdDQUFnQyxDQUFDO0FBQ3BELE9BQU8sS0FBSyxJQUFJLE1BQU0sd0JBQXdCLENBQUM7QUFDL0MsT0FBTyxFQUFDLFdBQVcsRUFBQyxNQUFNLHlCQUF5QixDQUFDO0FBQ3BELE9BQU8sS0FBSyxDQUFDLE1BQU0sNEJBQTRCLENBQUM7QUFDaEQsT0FBTyxFQUFDLGVBQWUsRUFBQyxNQUFNLHFCQUFxQixDQUFDO0FBQ3BELE9BQU8sS0FBSyxDQUFDLE1BQU0seUJBQXlCLENBQUM7QUFFN0MsT0FBTyxFQUFDLGtCQUFrQixFQUFDLE1BQU0saUNBQWlDLENBQUM7QUFDbkUsT0FBTyxFQUFDLHdCQUF3QixFQUFDLE1BQU0sNkNBQTZDLENBQUM7QUFFckYsT0FBTyxLQUFLLEVBQUUsTUFBTSxPQUFPLENBQUM7QUFFNUIsT0FBTyxFQUVMLHVCQUF1QixFQUN2Qix5QkFBeUIsR0FHMUIsTUFBTSxlQUFlLENBQUM7QUFDdkIsT0FBTyxFQUFDLGdCQUFnQixFQUFFLGVBQWUsRUFBRSxtQkFBbUIsRUFBQyxNQUFNLGNBQWMsQ0FBQztBQUVwRixNQUFNLGlCQUFpQixHQUFHLEVBQUUsQ0FBQyxpQkFBaUIsQ0FBQyx5QkFBeUIsQ0FBQztBQUV6RSx1REFBdUQ7QUFDdkQsTUFBTSxTQUFTLEdBQUcsSUFBSSx3QkFBd0IsRUFBRSxDQUFDO0FBRWpELHlDQUF5QztBQUN6QyxNQUFNLG9CQUFvQixHQUFHLGFBQWEsQ0FBQztBQUUzQyxNQUFNLFVBQVUsY0FBYyxDQUFDLElBQW9CO0lB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dkIsSUFBSSxFQUNKLElBQUksQ0FBQyxJQUFJLEVBQ1QsSUFBSSxFQUNKLElBQUksRUFDSixJQUFJLEVBQ0osZUFBZSxDQUNoQixDQUNGLENBQUM7WUFDSixDQUFDO1FBQ0gsQ0FBQztRQUVELE1BQU0sT0FBTyxHQUFHLElBQUksWUFBWSxDQUFDLENBQUMsT0FBTyxDQUFDLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQyxJQUFJLENBQUMsT0FBTyxDQUFDO1FBRXJFLCtFQUErRTtRQUMvRSxPQUFPLE9BQU8sS0FBSyxvQkFBb0IsQ0FBQyxDQUFDLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQyxPQUFPLENBQUM7SUFDM0QsQ0FBQztJQUVELE9BQU8sSUFBSSxDQUFDO0FBQ2QsQ0FBQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQGxpY2Vuc2VcbiAqIENvcHlyaWdodCBHb29nbGUgTExDIEFsbCBSaWdodHMgUmVzZXJ2ZWQuXG4gKlxuICogVXNlIG9mIHRoaXMgc291cmNlIGNvZGUgaXMgZ292ZXJuZWQgYnkgYW4gTUlULXN0eWxlIGxpY2Vuc2UgdGhhdCBjYW4gYmVcbiAqIGZvdW5kIGluIHRoZSBMSUNFTlNFIGZpbGUgYXQgaHR0cHM6Ly9hbmd1bGFyLmlvL2xpY2Vuc2VcbiAqL1xuXG5pbXBvcnQge0NvbnN0YW50UG9vbH0gZnJvbSAnLi4vLi4vLi4vY29uc3RhbnRfcG9vbCc7XG5pbXBvcnQge1NlY3VyaXR5Q29udGV4dH0gZnJvbSAnLi4vLi4vLi4vY29yZSc7XG5pbXBvcnQgKiBhcyBlIGZyb20gJy4uLy4uLy4uL2V4cHJlc3Npb25fcGFyc2VyL2FzdCc7XG5pbXBvcnQgKiBhcyBpMThuIGZyb20gJy4uLy4uLy4uL2kxOG4vaTE4bl9hc3QnO1xuaW1wb3J0IHtzcGxpdE5zTmFtZX0gZnJvbSAnLi4vLi4vLi4vbWxfcGFyc2VyL3RhZ3MnO1xuaW1wb3J0ICogYXMgbyBmcm9tICcuLi8uLi8uLi9vdXRwdXQvb3V0cHV0X2FzdCc7XG5pbXBvcnQge1BhcnNlU291cmNlU3Bhbn0gZnJvbSAnLi4vLi4vLi4vcGFyc2VfdXRpbCc7XG5pbXBvcnQgKiBhcyB0IGZyb20gJy4uLy4uLy4uL3JlbmRlcjMvcjNfYXN0JztcbmltcG9ydCB7RGVmZXJCbG9ja0RlcHNFbWl0TW9kZSwgUjNDb21wb25lbnREZWZlck1ldGFkYXRhfSBmcm9tICcuLi8uLi8uLi9yZW5kZXIzL3ZpZXcvYXBpJztcbmltcG9ydCB7aWN1RnJvbUkxOG5NZXNzYWdlfSBmcm9tICcuLi8uLi8uLi9yZW5kZXIzL3ZpZXcvaTE4bi91dGlsJztcbmltcG9ydCB7RG9tRWxlbWVudFNjaGVtYVJlZ2lzdHJ5fSBmcm9tICcuLi8uLi8uLi9zY2hlbWEvZG9tX2VsZW1lbnRfc2NoZW1hX3JlZ2lzdHJ5JztcbmltcG9ydCB7QmluZGluZ1BhcnNlcn0gZnJvbSAnLi4vLi4vLi4vdGVtcGxhdGVfcGFyc2VyL2JpbmRpbmdfcGFyc2VyJztcbmltcG9ydCAqIGFzIGlyIGZyb20gJy4uL2lyJztcblxuaW1wb3J0IHtcbiAgQ29tcGlsYXRpb25Vbml0LFxuICBDb21wb25lbnRDb21waWxhdGlvbkpvYixcbiAgSG9zdEJpbmRpbmdDb21waWxhdGlvbkpvYixcbiAgdHlwZSBDb21waWxhdGlvbkpvYixcbiAgdHlwZSBWaWV3Q29tcGlsYXRpb25Vbml0LFxufSBmcm9tICcuL2NvbXBpbGF0aW9uJztcbmltcG9ydCB7QklOQVJZX09QRVJBVE9SUywgbmFtZXNwYWNlRm9yS2V5LCBwcmVmaXhXaXRoTmFtZXNwYWNlfSBmcm9tICcuL2NvbnZlcnNpb24nO1xuXG5jb25zdCBjb21wYXRpYmlsaXR5TW9kZSA9IGlyLkNvbXBhdGliaWxpdHlNb2RlLlRlbXBsYXRlRGVmaW5pdGlvbkJ1aWxkZXI7XG5cbi8vIFNjaGVtYSBjb250YWluaW5nIERPTSBlbGVtZW50cyBhbmQgdGhlaXIgcHJvcGVydGllcy5cbmNvbnN0IGRvbVNjaGVtYSA9IG5ldyBEb21FbGVtZW50U2NoZW1hUmVnaXN0cnkoKTtcblxuLy8gVGFnIG5hbWUgb2YgdGhlIGBuZy10ZW1wbGF0ZWAgZWxlbWVudC5cbmNvbnN0IE5HX1RFTVBMQVRFX1RBR19OQU1FID0gJ25nLXRlbXBsYXRlJztcblxuZXhwb3J0IGZ1bmN0aW9uIGlzSTE4blJvb3ROb2RlKG1ldGE/OiBpMThuLkkxOG5NZXRhKTogbWV0YSBpcyBpMThuLk1lc3NhZ2Uge1xuICByZXR1cm4gbWV0YSBpbnN0YW5jZW9mIGkxOG4uTWVzc2FnZTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGlzU2luZ2xlSTE4bkljdShtZXRhPzogaTE4bi5JMThuTWV0YSk6IG1ldGEgaXMgaTE4bi5JMThuTWV0YSAmIHtub2RlczogW2kxOG4uSWN1XX0ge1xuICByZXR1cm4gaXNJMThuUm9vdE5vZGUobWV0YSkgJiYgbWV0YS5ub2Rlcy5sZW5ndGggPT09IDEgJiYgbWV0YS5ub2Rlc1swXSBpbnN0YW5jZW9mIGkxOG4uSWN1O1xufVxuXG4vKipcbiAqIFByb2Nlc3MgYSB0ZW1wbGF0ZSBBU1QgYW5kIGNvbnZlcnQgaXQgaW50byBhIGBDb21wb25lbnRDb21waWxhdGlvbmAgaW4gdGhlIGludGVybWVkaWF0ZVxuICogcmVwcmVzZW50YXRpb24uXG4gKiBUT0RPOiBSZWZhY3RvciBtb3JlIG9mIHRoZSBpbmdlc3Rpb24gY29kZSBpbnRvIHBoYXNlcy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGluZ2VzdENvbXBvbmVudChcbiAgY29tcG9uZW50TmFtZTogc3RyaW5nLFxuICB0ZW1wbGF0ZTogdC5Ob2RlW10sXG4gIGNvbnN0YW50UG9vbDogQ29uc3RhbnRQb29sLFxuICByZWxhdGl2ZUNvbnRleHRGaWxlUGF0aDogc3RyaW5nLFxuICBpMThuVXNlRXh0ZXJuYWxJZHM6IGJvb2xlYW4sXG4gIGRlZmVyTWV0YTogUjNDb21wb25lbnREZWZlck1ldGFkYXRhLFxuICBhbGxEZWZlcnJhYmxlRGVwc0ZuOiBvLlJlYWRWYXJFeHByIHwgbnVsbCxcbik6IENvbXBvbmVudENvbXBpbGF0aW9uSm9iIHtcbiAgY29uc3Qgam9iID0gbmV3IENvbXBvbmVudENvbXBpbGF0aW9uSm9iKFxuICAgIGNvbXBvbmVudE5hbWUsXG4gICAgY29uc3RhbnRQb29sLFxuICAgIGNvbXBhdGliaWxpdHlNb2RlLFxuICAgIHJlbGF0aXZlQ29udGV4dEZpbGVQYXRoLFxuICAgIGkxOG5Vc2VFeHRlcm5hbElkcyxcbiAgICBkZWZlck1ldGEsXG4gICAgYWxsRGVmZXJyYWJsZURlcHNGbixcbiAgKTtcbiAgaW5nZXN0Tm9kZXMoam9iLnJvb3QsIHRlbXBsYXRlKTtcbiAgcmV0dXJuIGpvYjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBIb3N0QmluZGluZ0lucHV0IHtcbiAgY29tcG9uZW50TmFtZTogc3RyaW5nO1xuICBjb21wb25lbnRTZWxlY3Rvcjogc3RyaW5nO1xuICBwcm9wZXJ0aWVzOiBlLlBhcnNlZFByb3BlcnR5W10gfCBudWxsO1xuICBhdHRyaWJ1dGVzOiB7W2tleTogc3RyaW5nXTogby5FeHByZXNzaW9ufTtcbiAgZXZlbnRzOiBlLlBhcnNlZEV2ZW50W10gfCBudWxsO1xufVxuXG4vKipcbiAqIFByb2Nlc3MgYSBob3N0IGJpbmRpbmcgQVNUIGFuZCBjb252ZXJ0IGl0IGludG8gYSBgSG9zdEJpbmRpbmdDb21waWxhdGlvbkpvYmAgaW4gdGhlIGludGVybWVkaWF0ZVxuICogcmVwcmVzZW50YXRpb24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBpbmdlc3RIb3N0QmluZGluZyhcbiAgaW5wdXQ6IEhvc3RCaW5kaW5nSW5wdXQsXG4gIGJpbmRpbmdQYXJzZXI6IEJpbmRpbmdQYXJzZXIsXG4gIGNvbnN0YW50UG9vbDogQ29uc3RhbnRQb29sLFxuKTogSG9zdEJpbmRpbmdDb21waWxhdGlvbkpvYiB7XG4gIGNvbnN0IGpvYiA9IG5ldyBIb3N0QmluZGluZ0NvbXBpbGF0aW9uSm9iKGlucHV0LmNvbXBvbmVudE5hbWUsIGNvbnN0YW50UG9vbCwgY29tcGF0aWJpbGl0eU1vZGUpO1xuICBmb3IgKGNvbnN0IHByb3BlcnR5IG9mIGlucHV0LnByb3BlcnRpZXMgPz8gW10pIHtcbiAgICBsZXQgYmluZGluZ0tpbmQgPSBpci5CaW5kaW5nS2luZC5Qcm9wZXJ0eTtcbiAgICAvLyBUT0RPOiB0aGlzIHNob3VsZCByZWFsbHkgYmUgaGFuZGxlZCBpbiB0aGUgcGFyc2VyLlxuICAgIGlmIChwcm9wZXJ0eS5uYW1lLnN0YXJ0c1dpdGgoJ2F0dHIuJykpIHtcbiAgICAgIHByb3BlcnR5Lm5hbWUgPSBwcm9wZXJ0eS5uYW1lLnN1YnN0cmluZygnYXR0ci4nLmxlbmd0aCk7XG4gICAgICBiaW5kaW5nS2luZCA9IGlyLkJpbmRpbmdLaW5kLkF0dHJpYnV0ZTtcbiAgICB9XG4gICAgaWYgKHByb3BlcnR5LmlzQW5pbWF0aW9uKSB7XG4gICAgICBiaW5kaW5nS2luZCA9IGlyLkJpbmRpbmdLaW5kLkFuaW1hdGlvbjtcbiAgICB9XG4gICAgY29uc3Qgc2VjdXJpdHlDb250ZXh0cyA9IGJpbmRpbmdQYXJzZXJcbiAgICAgIC5jYWxjUG9zc2libGVTZWN1cml0eUNvbnRleHRzKFxuICAgICAgICBpbnB1dC5jb21wb25lbnRTZWxlY3RvcixcbiAgICAgICAgcHJvcGVydHkubmFtZSxcbiAgICAgICAgYmluZGluZ0tpbmQgPT09IGlyLkJpbmRpbmdLaW5kLkF0dHJpYnV0ZSxcbiAgICAgIClcbiAgICAgIC5maWx0ZXIoKGNvbnRleHQpID0+IGNvbnRleHQgIT09IFNlY3VyaXR5Q29udGV4dC5OT05FKTtcbiAgICBpbmdlc3RIb3N0UHJvcGVydHkoam9iLCBwcm9wZXJ0eSwgYmluZGluZ0tpbmQsIHNlY3VyaXR5Q29udGV4dHMpO1xuICB9XG4gIGZvciAoY29uc3QgW25hbWUsIGV4cHJdIG9mIE9iamVjdC5lbnRyaWVzKGlucHV0LmF0dHJpYnV0ZXMpID8/IFtdKSB7XG4gICAgY29uc3Qgc2VjdXJpdHlDb250ZXh0cyA9IGJpbmRpbmdQYXJzZXJcbiAgICAgIC5jYWxjUG9zc2libGVTZWN1cml0eUNvbnRleHRzKGlucHV0LmNvbXBvbmVudFNlbGVjdG9yLCBuYW1lLCB0cnVlKVxuICAgICAgLmZpbHRlcigoY29udGV4dCkgPT4gY29udGV4dCAhPT0gU2VjdXJpdHlDb250ZXh0Lk5PTkUpO1xuICAgIGluZ2VzdEhvc3RBdHRyaWJ1dGUoam9iLCBuYW1lLCBleHByLCBzZWN1cml0eUNvbnRleHRzKTtcbiAgfVxuICBmb3IgKGNvbnN0IGV2ZW50IG9mIGlucHV0LmV2ZW50cyA/PyBbXSkge1xuICAgIGluZ2VzdEhvc3RFdmVudChqb2IsIGV2ZW50KTtcbiAgfVxuICByZXR1cm4gam9iO1xufVxuXG4vLyBUT0RPOiBXZSBzaG91bGQgcmVmYWN0b3IgdGhlIHBhcnNlciB0byB1c2UgdGhlIHNhbWUgdHlwZXMgYW5kIHN0cnVjdHVyZXMgZm9yIGhvc3QgYmluZGluZ3MgYXNcbi8vIHdpdGggb3JkaW5hcnkgY29tcG9uZW50cy4gVGhpcyB3b3VsZCBhbGxvdyB1cyB0byBzaGFyZSBhIGxvdCBtb3JlIGluZ2VzdGlvbiBjb2RlLlxuZXhwb3J0IGZ1bmN0aW9uIGluZ2VzdEhvc3RQcm9wZXJ0eShcbiAgam9iOiBIb3N0QmluZGluZ0NvbXBpbGF0aW9uSm9iLFxuICBwcm9wZXJ0eTogZS5QYXJzZWRQcm9wZXJ0eSxcbiAgYmluZGluZ0tpbmQ6IGlyLkJpbmRpbmdLaW5kLFxuICBzZWN1cml0eUNvbnRleHRzOiBTZWN1cml0eUNvbnRleHRbXSxcbik6IHZvaWQge1xuICBsZXQgZXhwcmVzc2lvbjogby5FeHByZXNzaW9uIHwgaXIuSW50ZXJwb2xhdGlvbjtcbiAgY29uc3QgYXN0ID0gcHJvcGVydHkuZXhwcmVzc2lvbi5hc3Q7XG4gIGlmIChhc3QgaW5zdGFuY2VvZiBlLkludGVycG9sYXRpb24pIHtcbiAgICBleHByZXNzaW9uID0gbmV3IGlyLkludGVycG9sYXRpb24oXG4gICAgICBhc3Quc3RyaW5ncyxcbiAgICAgIGFzdC5leHByZXNzaW9ucy5tYXAoKGV4cHIpID0+IGNvbnZlcnRBc3QoZXhwciwgam9iLCBwcm9wZXJ0eS5zb3VyY2VTcGFuKSksXG4gICAgICBbXSxcbiAgICApO1xuICB9IGVsc2Uge1xuICAgIGV4cHJlc3Npb24gPSBjb252ZXJ0QXN0KGFzdCwgam9iLCBwcm9wZXJ0eS5zb3VyY2VTcGFuKTtcbiAgfVxuICBqb2Iucm9vdC51cGRhdGUucHVzaChcbiAgICBpci5jcmVhdGVCaW5kaW5nT3AoXG4gICAgICBqb2Iucm9vdC54cmVmLFxuICAgICAgYmluZGluZ0tpbmQsXG4gICAgICBwcm9wZXJ0eS5uYW1lLFxuICAgICAgZXhwcmVzc2lvbixcbiAgICAgIG51bGwsXG4gICAgICBzZWN1cml0eUNvbnRleHRzLFxuICAgICAgZmFsc2UsXG4gICAgICBmYWxzZSxcbiAgICAgIG51bGwsXG4gICAgICAvKiBUT0RPOiBIb3cgZG8gSG9zdCBiaW5kaW5ncyBoYW5kbGUgaTE4biBhdHRycz8gKi8gbnVsbCxcbiAgICAgIHByb3BlcnR5LnNvdXJjZVNwYW4sXG4gICAgKSxcbiAgKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGluZ2VzdEhvc3RBdHRyaWJ1dGUoXG4gIGpvYjogSG9zdEJpbmRpbmdDb21waWxhdGlvbkpvYixcbiAgbmFtZTogc3RyaW5nLFxuICB2YWx1ZTogby5FeHByZXNzaW9uLFxuICBzZWN1cml0eUNvbnRleHRzOiBTZWN1cml0eUNvbnRleHRbXSxcbik6IHZvaWQge1xuICBjb25zdCBhdHRyQmluZGluZyA9IGlyLmNyZWF0ZUJpbmRpbmdPcChcbiAgICBqb2Iucm9vdC54cmVmLFxuICAgIGlyLkJpbmRpbmdLaW5kLkF0dHJpYnV0ZSxcbiAgICBuYW1lLFxuICAgIHZhbHVlLFxuICAgIG51bGwsXG4gICAgc2VjdXJpdHlDb250ZXh0cyxcbiAgICAvKiBIb3N0IGF0dHJpYnV0ZXMgc2hvdWxkIGFsd2F5cyBiZSBleHRyYWN0ZWQgdG8gY29uc3QgaG9zdEF0dHJzLCBldmVuIGlmIHRoZXkgYXJlIG5vdFxuICAgICAqc3RyaWN0bHkqIHRleHQgbGl0ZXJhbHMgKi9cbiAgICB0cnVlLFxuICAgIGZhbHNlLFxuICAgIG51bGwsXG4gICAgLyogVE9ETyAqLyBudWxsLFxuICAgIC8qKiBUT0RPOiBNYXkgYmUgbnVsbD8gKi8gdmFsdWUuc291cmNlU3BhbiEsXG4gICk7XG4gIGpvYi5yb290LnVwZGF0ZS5wdXNoKGF0dHJCaW5kaW5nKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGluZ2VzdEhvc3RFdmVudChqb2I6IEhvc3RCaW5kaW5nQ29tcGlsYXRpb25Kb2IsIGV2ZW50OiBlLlBhcnNlZEV2ZW50KSB7XG4gIGNvbnN0IFtwaGFzZSwgdGFyZ2V0XSA9XG4gICAgZXZlbnQudHlwZSAhPT0gZS5QYXJzZWRFdmVudFR5cGUuQW5pbWF0aW9uXG4gICAgICA/IFtudWxsLCBldmVudC50YXJnZXRPclBoYXNlXVxuICAgICAgOiBbZXZlbnQudGFyZ2V0T3JQaGFzZSwgbnVsbF07XG4gIGNvbnN0IGV2ZW50QmluZGluZyA9IGlyLmNyZWF0ZUxpc3RlbmVyT3AoXG4gICAgam9iLnJvb3QueHJlZixcbiAgICBuZXcgaXIuU2xvdEhhbmRsZSgpLFxuICAgIGV2ZW50Lm5hbWUsXG4gICAgbnVsbCxcbiAgICBtYWtlTGlzdGVuZXJIYW5kbGVyT3BzKGpvYi5yb290LCBldmVudC5oYW5kbGVyLCBldmVudC5oYW5kbGVyU3BhbiksXG4gICAgcGhhc2UsXG4gICAgdGFyZ2V0LFxuICAgIHRydWUsXG4gICAgZXZlbnQuc291cmNlU3BhbixcbiAgKTtcbiAgam9iLnJvb3QuY3JlYXRlLnB1c2goZXZlbnRCaW5kaW5nKTtcbn1cblxuLyoqXG4gKiBJbmdlc3QgdGhlIG5vZGVzIG9mIGEgdGVtcGxhdGUgQVNUIGludG8gdGhlIGdpdmVuIGBWaWV3Q29tcGlsYXRpb25gLlxuICovXG5mdW5jdGlvbiBpbmdlc3ROb2Rlcyh1bml0OiBWaWV3Q29tcGlsYXRpb25Vbml0LCB0ZW1wbGF0ZTogdC5Ob2RlW10pOiB2b2lkIHtcbiAgZm9yIChjb25zdCBub2RlIG9mIHRlbXBsYXRlKSB7XG4gICAgaWYgKG5vZGUgaW5zdGFuY2VvZiB0LkVsZW1lbnQpIHtcbiAgICAgIGluZ2VzdEVsZW1lbnQodW5pdCwgbm9kZSk7XG4gICAgfSBlbHNlIGlmIChub2RlIGluc3RhbmNlb2YgdC5UZW1wbGF0ZSkge1xuICAgICAgaW5nZXN0VGVtcGxhdGUodW5pdCwgbm9kZSk7XG4gICAgfSBlbHNlIGlmIChub2RlIGluc3RhbmNlb2YgdC5Db250ZW50KSB7XG4gICAgICBpbmdlc3RDb250ZW50KHVuaXQsIG5vZGUpO1xuICAgIH0gZWxzZSBpZiAobm9kZSBpbnN0YW5jZW9mIHQuVGV4dCkge1xuICAgICAgaW5nZXN0VGV4dCh1bml0LCBub2RlLCBudWxsKTtcbiAgICB9IGVsc2UgaWYgKG5vZGUgaW5zdGFuY2VvZiB0LkJvdW5kVGV4dCkge1xuICAgICAgaW5nZXN0Qm91bmRUZXh0KHVuaXQsIG5vZGUsIG51bGwpO1xuICAgIH0gZWxzZSBpZiAobm9kZSBpbnN0YW5jZW9mIHQuSWZCbG9jaykge1xuICAgICAgaW5nZXN0SWZCbG9jayh1bml0LCBub2RlKTtcbiAgICB9IGVsc2UgaWYgKG5vZGUgaW5zdGFuY2VvZiB0LlN3aXRjaEJsb2NrKSB7XG4gICAgICBpbmdlc3RTd2l0Y2hCbG9jayh1bml0LCBub2RlKTtcbiAgICB9IGVsc2UgaWYgKG5vZGUgaW5zdGFuY2VvZiB0LkRlZmVycmVkQmxvY2spIHtcbiAgICAgIGluZ2VzdERlZmVyQmxvY2sodW5pdCwgbm9kZSk7XG4gICAgfSBlbHNlIGlmIChub2RlIGluc3RhbmNlb2YgdC5JY3UpIHtcbiAgICAgIGluZ2VzdEljdSh1bml0LCBub2RlKTtcbiAgICB9IGVsc2UgaWYgKG5vZGUgaW5zdGFuY2VvZiB0LkZvckxvb3BCbG9jaykge1xuICAgICAgaW5nZXN0Rm9yQmxvY2sodW5pdCwgbm9kZSk7XG4gICAgfSBlbHNlIGlmIChub2RlIGluc3RhbmNlb2YgdC5MZXREZWNsYXJhdGlvbikge1xuICAgICAgaW5nZXN0TGV0RGVjbGFyYXRpb24odW5pdCwgbm9kZSk7XG4gICAgfSBlbHNlIHtcbiAgICAgIHRocm93IG5ldyBFcnJvcihgVW5zdXBwb3J0ZWQgdGVtcGxhdGUgbm9kZTogJHtub2RlLmNvbnN0cnVjdG9yLm5hbWV9YCk7XG4gICAgfVxuICB9XG59XG5cbi8qKlxuICogSW5nZXN0IGFuIGVsZW1lbnQgQVNUIGZyb20gdGhlIHRlbXBsYXRlIGludG8gdGhlIGdpdmVuIGBWaWV3Q29tcGlsYXRpb25gLlxuICovXG5mdW5jdGlvbiBpbmdlc3RFbGVtZW50KHVuaXQ6IFZpZXdDb21waWxhdGlvblVuaXQsIGVsZW1lbnQ6IHQuRWxlbWVudCk6IHZvaWQge1xuICBpZiAoXG4gICAgZWxlbWVudC5pMThuICE9PSB1bmRlZmluZWQgJiZcbiAgICAhKGVsZW1lbnQuaTE4biBpbnN0YW5jZW9mIGkxOG4uTWVzc2FnZSB8fCBlbGVtZW50LmkxOG4gaW5zdGFuY2VvZiBpMThuLlRhZ1BsYWNlaG9sZGVyKVxuICApIHtcbiAgICB0aHJvdyBFcnJvcihgVW5oYW5kbGVkIGkxOG4gbWV0YWRhdGEgdHlwZSBmb3IgZWxlbWVudDogJHtlbGVtZW50LmkxOG4uY29uc3RydWN0b3IubmFtZX1gKTtcbiAgfVxuXG4gIGNvbnN0IGlkID0gdW5pdC5qb2IuYWxsb2NhdGVYcmVmSWQoKTtcblxuICBjb25zdCBbbmFtZXNwYWNlS2V5LCBlbGVtZW50TmFtZV0gPSBzcGxpdE5zTmFtZShlbGVtZW50Lm5hbWUpO1xuXG4gIGNvbnN0IHN0YXJ0T3AgPSBpci5jcmVhdGVFbGVtZW50U3RhcnRPcChcbiAgICBlbGVtZW50TmFtZSxcbiAgICBpZCxcbiAgICBuYW1lc3BhY2VGb3JLZXkobmFtZXNwYWNlS2V5KSxcbiAgICBlbGVtZW50LmkxOG4gaW5zdGFuY2VvZiBpMThuLlRhZ1BsYWNlaG9sZGVyID8gZWxlbWVudC5pMThuIDogdW5kZWZpbmVkLFxuICAgIGVsZW1lbnQuc3RhcnRTb3VyY2VTcGFuLFxuICAgIGVsZW1lbnQuc291cmNlU3BhbixcbiAgKTtcbiAgdW5pdC5jcmVhdGUucHVzaChzdGFydE9wKTtcblxuICBpbmdlc3RFbGVtZW50QmluZGluZ3ModW5pdCwgc3RhcnRPcCwgZWxlbWVudCk7XG4gIGluZ2VzdFJlZmVyZW5jZXMoc3RhcnRPcCwgZWxlbWVudCk7XG5cbiAgLy8gU3RhcnQgaTE4biwgaWYgbmVlZGVkLCBnb2VzIGFmdGVyIHRoZSBlbGVtZW50IGNyZWF0ZSBhbmQgYmluZGluZ3MsIGJ1dCBiZWZvcmUgdGhlIG5vZGVzXG4gIGxldCBpMThuQmxvY2tJZDogaXIuWHJlZklkIHwgbnVsbCA9IG51bGw7XG4gIGlmIChlbGVtZW50LmkxOG4gaW5zdGFuY2VvZiBpMThuLk1lc3NhZ2UpIHtcbiAgICBpMThuQmxvY2tJZCA9IHVuaXQuam9iLmFsbG9jYXRlWHJlZklkKCk7XG4gICAgdW5pdC5jcmVhdGUucHVzaChcbiAgICAgIGlyLmNyZWF0ZUkxOG5TdGFydE9wKGkxOG5CbG9ja0lkLCBlbGVtZW50LmkxOG4sIHVuZGVmaW5lZCwgZWxlbWVudC5zdGFydFNvdXJjZVNwYW4pLFxuICAgICk7XG4gIH1cblxuICBpbmdlc3ROb2Rlcyh1bml0LCBlbGVtZW50LmNoaWxkcmVuKTtcblxuICAvLyBUaGUgc291cmNlIHNwYW4gZm9yIHRoZSBlbmQgb3AgaXMgdHlwaWNhbGx5IHRoZSBlbGVtZW50IGNsb3NpbmcgdGFnLiBIb3dldmVyLCBpZiBubyBjbG9zaW5nIHRhZ1xuICAvLyBleGlzdHMsIHN1Y2ggYXMgaW4gYDxpbnB1dD5gLCB3ZSB1c2UgdGhlIHN0YXJ0IHNvdXJjZSBzcGFuIGluc3RlYWQuIFVzdWFsbHkgdGhlIHN0YXJ0IGFuZCBlbmRcbiAgLy8gaW5zdHJ1Y3Rpb25zIHdpbGwgYmUgY29sbGFwc2VkIGludG8gb25lIGBlbGVtZW50YCBpbnN0cnVjdGlvbiwgbmVnYXRpbmcgdGhlIHB1cnBvc2Ugb2YgdGhpc1xuICAvLyBmYWxsYmFjaywgYnV0IGluIGNhc2VzIHdoZW4gaXQgaXMgbm90IGNvbGxhcHNlZCAoc3VjaCBhcyBhbiBpbnB1dCB3aXRoIGEgYmluZGluZyksIHdlIHN0aWxsXG4gIC8vIHdhbnQgdG8gbWFwIHRoZSBlbmQgaW5zdHJ1Y3Rpb24gdG8gdGhlIG1haW4gZWxlbWVudC5cbiAgY29uc3QgZW5kT3AgPSBpci5jcmVhdGVFbGVtZW50RW5kT3AoaWQsIGVsZW1lbnQuZW5kU291cmNlU3BhbiA/PyBlbGVtZW50LnN0YXJ0U291cmNlU3Bhbik7XG4gIHVuaXQuY3JlYXRlLnB1c2goZW5kT3ApO1xuXG4gIC8vIElmIHRoZXJlIGlzIGFuIGkxOG4gbWVzc2FnZSBhc3NvY2lhdGVkIHdpdGggdGhpcyBlbGVtZW50LCBpbnNlcnQgaTE4biBzdGFydCBhbmQgZW5kIG9wcy5cbiAgaWYgKGkxOG5CbG9ja0lkICE9PSBudWxsKSB7XG4gICAgaXIuT3BMaXN0Lmluc2VydEJlZm9yZTxpci5DcmVhdGVPcD4oXG4gICAgICBpci5jcmVhdGVJMThuRW5kT3AoaTE4bkJsb2NrSWQsIGVsZW1lbnQuZW5kU291cmNlU3BhbiA/PyBlbGVtZW50LnN0YXJ0U291cmNlU3BhbiksXG4gICAgICBlbmRPcCxcbiAgICApO1xuICB9XG59XG5cbi8qKlxuICogSW5nZXN0IGFuIGBuZy10ZW1wbGF0ZWAgbm9kZSBmcm9tIHRoZSBBU1QgaW50byB0aGUgZ2l2ZW4gYFZpZXdDb21waWxhdGlvbmAuXG4gKi9cbmZ1bmN0aW9uIGluZ2VzdFRlbXBsYXRlKHVuaXQ6IFZpZXdDb21waWxhdGlvblVuaXQsIHRtcGw6IHQuVGVtcGxhdGUpOiB2b2lkIHtcbiAgaWYgKFxuICAgIHRtcGwuaTE4biAhPT0gdW5kZWZpbmVkICYmXG4gICAgISh0bXBsLmkxOG4gaW5zdGFuY2VvZiBpMThuLk1lc3NhZ2UgfHwgdG1wbC5pMThuIGluc3RhbmNlb2YgaTE4bi5UYWdQbGFjZWhvbGRlcilcbiAgKSB7XG4gICAgdGhyb3cgRXJyb3IoYFVuaGFuZGxlZCBpMThuIG1ldGFkYXRhIHR5cGUgZm9yIHRlbXBsYXRlOiAke3RtcGwuaTE4bi5jb25zdHJ1Y3Rvci5uYW1lfWApO1xuICB9XG5cbiAgY29uc3QgY2hpbGRWaWV3ID0gdW5pdC5qb2IuYWxsb2NhdGVWaWV3KHVuaXQueHJlZik7XG5cbiAgbGV0IHRhZ05hbWVXaXRob3V0TmFtZXNwYWNlID0gdG1wbC50YWdOYW1lO1xuICBsZXQgbmFtZXNwYWNlUHJlZml4OiBzdHJpbmcgfCBudWxsID0gJyc7XG4gIGlmICh0bXBsLnRhZ05hbWUpIHtcbiAgICBbbmFtZXNwYWNlUHJlZml4LCB0YWdOYW1lV2l0aG91dE5hbWVzcGFjZV0gPSBzcGxpdE5zTmFtZSh0bXBsLnRhZ05hbWUpO1xuICB9XG5cbiAgY29uc3QgaTE4blBsYWNlaG9sZGVyID0gdG1wbC5pMThuIGluc3RhbmNlb2YgaTE4bi5UYWdQbGFjZWhvbGRlciA/IHRtcGwuaTE4biA6IHVuZGVmaW5lZDtcbiAgY29uc3QgbmFtZXNwYWNlID0gbmFtZXNwYWNlRm9yS2V5KG5hbWVzcGFjZVByZWZpeCk7XG4gIGNvbnN0IGZ1bmN0aW9uTmFtZVN1ZmZpeCA9XG4gICAgdGFnTmFtZVdpdGhvdXROYW1lc3BhY2UgPT09IG51bGwgPyAnJyA6IHByZWZpeFdpdGhOYW1lc3BhY2UodGFnTmFtZVdpdGhvdXROYW1lc3BhY2UsIG5hbWVzcGFjZSk7XG4gIGNvbnN0IHRlbXBsYXRlS2luZCA9IGlzUGxhaW5UZW1wbGF0ZSh0bXBsKVxuICAgID8gaXIuVGVtcGxhdGVLaW5kLk5nVGVtcGxhdGVcbiAgICA6IGlyLlRlbXBsYXRlS2luZC5TdHJ1Y3R1cmFsO1xuICBjb25zdCB0ZW1wbGF0ZU9wID0gaXIuY3JlYXRlVGVtcGxhdGVPcChcbiAgICBjaGlsZFZpZXcueHJlZixcbiAgICB0ZW1wbGF0ZUtpbmQsXG4gICAgdGFnTmFtZVdpdGhvdXROYW1lc3BhY2UsXG4gICAgZnVuY3Rpb25OYW1lU3VmZml4LFxuICAgIG5hbWVzcGFjZSxcbiAgICBpMThuUGxhY2Vob2xkZXIsXG4gICAgdG1wbC5zdGFydFNvdXJjZVNwYW4sXG4gICAgdG1wbC5zb3VyY2VTcGFuLFxuICApO1xuICB1bml0LmNyZWF0ZS5wdXNoKHRlbXBsYXRlT3ApO1xuXG4gIGluZ2VzdFRlbXBsYXRlQmluZGluZ3ModW5pdCwgdGVtcGxhdGVPcCwgdG1wbCwgdGVtcGxhdGVLaW5kKTtcbiAgaW5nZXN0UmVmZXJlbmNlcyh0ZW1wbGF0ZU9wLCB0bXBsKTtcbiAgaW5nZXN0Tm9kZXMoY2hpbGRWaWV3LCB0bXBsLmNoaWxkcmVuKTtcblxuICBmb3IgKGNvbnN0IHtuYW1lLCB2YWx1ZX0gb2YgdG1wbC52YXJpYWJsZXMpIHtcbiAgICBjaGlsZFZpZXcuY29udGV4dFZhcmlhYmxlcy5zZXQobmFtZSwgdmFsdWUgIT09ICcnID8gdmFsdWUgOiAnJGltcGxpY2l0Jyk7XG4gIH1cblxuICAvLyBJZiB0aGlzIGlzIGEgcGxhaW4gdGVtcGxhdGUgYW5kIHRoZXJlIGlzIGFuIGkxOG4gbWVzc2FnZSBhc3NvY2lhdGVkIHdpdGggaXQsIGluc2VydCBpMThuIHN0YXJ0XG4gIC8vIGFuZCBlbmQgb3BzLiBGb3Igc3RydWN0dXJhbCBkaXJlY3RpdmUgdGVtcGxhdGVzLCB0aGUgaTE4biBvcHMgd2lsbCBiZSBhZGRlZCB3aGVuIGluZ2VzdGluZyB0aGVcbiAgLy8gZWxlbWVudC90ZW1wbGF0ZSB0aGUgZGlyZWN0aXZlIGlzIHBsYWNlZCBvbi5cbiAgaWYgKHRlbXBsYXRlS2luZCA9PT0gaXIuVGVtcGxhdGVLaW5kLk5nVGVtcGxhdGUgJiYgdG1wbC5pMThuIGluc3RhbmNlb2YgaTE4bi5NZXNzYWdlKSB7XG4gICAgY29uc3QgaWQgPSB1bml0LmpvYi5hbGxvY2F0ZVhyZWZJZCgpO1xuICAgIGlyLk9wTGlzdC5pbnNlcnRBZnRlcihcbiAgICAgIGlyLmNyZWF0ZUkxOG5TdGFydE9wKGlkLCB0bXBsLmkxOG4sIHVuZGVmaW5lZCwgdG1wbC5zdGFydFNvdXJjZVNwYW4pLFxuICAgICAgY2hpbGRWaWV3LmNyZWF0ZS5oZWFkLFxuICAgICk7XG4gICAgaXIuT3BMaXN0Lmluc2VydEJlZm9yZShcbiAgICAgIGlyLmNyZWF0ZUkxOG5FbmRPcChpZCwgdG1wbC5lbmRTb3VyY2VTcGFuID8/IHRtcGwuc3RhcnRTb3VyY2VTcGFuKSxcbiAgICAgIGNoaWxkVmlldy5jcmVhdGUudGFpbCxcbiAgICApO1xuICB9XG59XG5cbi8qKlxuICogSW5nZXN0IGEgY29udGVudCBub2RlIGZyb20gdGhlIEFTVCBpbnRvIHRoZSBnaXZlbiBgVmlld0NvbXBpbGF0aW9uYC5cbiAqL1xuZnVuY3Rpb24gaW5nZXN0Q29udGVudCh1bml0OiBWaWV3Q29tcGlsYXRpb25Vbml0LCBjb250ZW50OiB0LkNvbnRlbnQpOiB2b2lkIHtcbiAgaWYgKGNvbnRlbnQuaTE4biAhPT0gdW5kZWZpbmVkICYmICEoY29udGVudC5pMThuIGluc3RhbmNlb2YgaTE4bi5UYWdQbGFjZWhvbGRlcikpIHtcbiAgICB0aHJvdyBFcnJvcihgVW5oYW5kbGVkIGkxOG4gbWV0YWRhdGEgdHlwZSBmb3IgZWxlbWVudDogJHtjb250ZW50LmkxOG4uY29uc3RydWN0b3IubmFtZX1gKTtcbiAgfVxuXG4gIGxldCBmYWxsYmFja1ZpZXc6IFZpZXdDb21waWxhdGlvblVuaXQgfCBudWxsID0gbnVsbDtcblxuICAvLyBEb24ndCBjYXB0dXJlIGRlZmF1bHQgY29udGVudCB0aGF0J3Mgb25seSBtYWRlIHVwIG9mIGVtcHR5IHRleHQgbm9kZXMgYW5kIGNvbW1lbnRzLlxuICAvLyBOb3RlIHRoYXQgd2UgcHJvY2VzcyB0aGUgZGVmYXVsdCBjb250ZW50IGJlZm9yZSB0aGUgcHJvamVjdGlvbiBpbiBvcmRlciB0byBtYXRjaCB0aGVcbiAgLy8gaW5zZXJ0aW9uIG9yZGVyIGF0IHJ1bnRpbWUuXG4gIGlmIChcbiAgICBjb250ZW50LmNoaWxkcmVuLnNvbWUoXG4gICAgICAoY2hpbGQpID0+XG4gICAgICAgICEoY2hpbGQgaW5zdGFuY2VvZiB0LkNvbW1lbnQpICYmXG4gICAgICAgICghKGNoaWxkIGluc3RhbmNlb2YgdC5UZXh0KSB8fCBjaGlsZC52YWx1ZS50cmltKCkubGVuZ3RoID4gMCksXG4gICAgKVxuICApIHtcbiAgICBmYWxsYmFja1ZpZXcgPSB1bml0LmpvYi5hbGxvY2F0ZVZpZXcodW5pdC54cmVmKTtcbiAgICBpbmdlc3ROb2RlcyhmYWxsYmFja1ZpZXcsIGNvbnRlbnQuY2hpbGRyZW4pO1xuICB9XG5cbiAgY29uc3QgaWQgPSB1bml0LmpvYi5hbGxvY2F0ZVhyZWZJZCgpO1xuICBjb25zdCBvcCA9IGlyLmNyZWF0ZVByb2plY3Rpb25PcChcbiAgICBpZCxcbiAgICBjb250ZW50LnNlbGVjdG9yLFxuICAgIGNvbnRlbnQuaTE4bixcbiAgICBmYWxsYmFja1ZpZXc/LnhyZWYgPz8gbnVsbCxcbiAgICBjb250ZW50LnNvdXJjZVNwYW4sXG4gICk7XG4gIGZvciAoY29uc3QgYXR0ciBvZiBjb250ZW50LmF0dHJpYnV0ZXMpIHtcbiAgICBjb25zdCBzZWN1cml0eUNvbnRleHQgPSBkb21TY2hlbWEuc2VjdXJpdHlDb250ZXh0KGNvbnRlbnQubmFtZSwgYXR0ci5uYW1lLCB0cnVlKTtcbiAgICB1bml0LnVwZGF0ZS5wdXNoKFxuICAgICAgaXIuY3JlYXRlQmluZGluZ09wKFxuICAgICAgICBvcC54cmVmLFxuICAgICAgICBpci5CaW5kaW5nS2luZC5BdHRyaWJ1dGUsXG4gICAgICAgIGF0dHIubmFtZSxcbiAgICAgICAgby5saXRlcmFsKGF0dHIudmFsdWUpLFxuICAgICAgICBudWxsLFxuICAgICAgICBzZWN1cml0eUNvbnRleHQsXG4gICAgICAgIHRydWUsXG4gICAgICAgIGZhbHNlLFxuICAgICAgICBudWxsLFxuICAgICAgICBhc01lc3NhZ2UoYXR0ci5pMThuKSxcbiAgICAgICAgYXR0ci5zb3VyY2VTcGFuLFxuICAgICAgKSxcbiAgICApO1xuICB9XG4gIHVuaXQuY3JlYXRlLnB1c2gob3ApO1xufVxuXG4vKipcbiAqIEluZ2VzdCBhIGxpdGVyYWwgdGV4dCBub2RlIGZyb20gdGhlIEFTVCBpbnRvIHRoZSBnaXZlbiBgVmlld0NvbXBpbGF0aW9uYC5cbiAqL1xuZnVuY3Rpb24gaW5nZXN0VGV4dCh1bml0OiBWaWV3Q29tcGlsYXRpb25Vbml0LCB0ZXh0OiB0LlRleHQsIGljdVBsYWNlaG9sZGVyOiBzdHJpbmcgfCBudWxsKTogdm9pZCB7XG4gIHVuaXQuY3JlYXRlLnB1c2goXG4gICAgaXIuY3JlYXRlVGV4dE9wKHVuaXQuam9iLmFsbG9jYXRlWHJlZklkKCksIHRleHQudmFsdWUsIGljdVBsYWNlaG9sZGVyLCB0ZXh0LnNvdXJjZVNwYW4pLFxuICApO1xufVxuXG4vKipcbiAqIEluZ2VzdCBhbiBpbnRlcnBvbGF0ZWQgdGV4dCBub2RlIGZyb20gdGhlIEFTVCBpbnRvIHRoZSBnaXZlbiBgVmlld0NvbXBpbGF0aW9uYC5cbiAqL1xuZnVuY3Rpb24gaW5nZXN0Qm91bmRUZXh0KFxuICB1bml0OiBWaWV3Q29tcGlsYXRpb25Vbml0LFxuICB0ZXh0OiB0LkJvdW5kVGV4dCxcbiAgaWN1UGxhY2Vob2xkZXI6IHN0cmluZyB8IG51bGwsXG4pOiB2b2lkIHtcbiAgbGV0IHZhbHVlID0gdGV4dC52YWx1ZTtcbiAgaWYgKHZhbHVlIGluc3RhbmNlb2YgZS5BU1RXaXRoU291cmNlKSB7XG4gICAgdmFsdWUgPSB2YWx1ZS5hc3Q7XG4gIH1cbiAgaWYgKCEodmFsdWUgaW5zdGFuY2VvZiBlLkludGVycG9sYXRpb24pKSB7XG4gICAgdGhyb3cgbmV3IEVycm9yKFxuICAgICAgYEFzc2VydGlvbkVycm9yOiBleHBlY3RlZCBJbnRlcnBvbGF0aW9uIGZvciBCb3VuZFRleHQgbm9kZSwgZ290ICR7dmFsdWUuY29uc3RydWN0b3IubmFtZX1gLFxuICAgICk7XG4gIH1cbiAgaWYgKHRleHQuaTE4biAhPT0gdW5kZWZpbmVkICYmICEodGV4dC5pMThuIGluc3RhbmNlb2YgaTE4bi5Db250YWluZXIpKSB7XG4gICAgdGhyb3cgRXJyb3IoXG4gICAgICBgVW5oYW5kbGVkIGkxOG4gbWV0YWRhdGEgdHlwZSBmb3IgdGV4dCBpbnRlcnBvbGF0aW9uOiAke3RleHQuaTE4bj8uY29uc3RydWN0b3IubmFtZX1gLFxuICAgICk7XG4gIH1cblxuICBjb25zdCBpMThuUGxhY2Vob2xkZXJzID1cbiAgICB0ZXh0LmkxOG4gaW5zdGFuY2VvZiBpMThuLkNvbnRhaW5lclxuICAgICAgPyB0ZXh0LmkxOG4uY2hpbGRyZW5cbiAgICAgICAgICAuZmlsdGVyKChub2RlKTogbm9kZSBpcyBpMThuLlBsYWNlaG9sZGVyID0+IG5vZGUgaW5zdGFuY2VvZiBpMThuLlBsYWNlaG9sZGVyKVxuICAgICAgICAgIC5tYXAoKHBsYWNlaG9sZGVyKSA9PiBwbGFjZWhvbGRlci5uYW1lKVxuICAgICAgOiBbXTtcbiAgaWYgKGkxOG5QbGFjZWhvbGRlcnMubGVuZ3RoID4gMCAmJiBpMThuUGxhY2Vob2xkZXJzLmxlbmd0aCAhPT0gdmFsdWUuZXhwcmVzc2lvbnMubGVuZ3RoKSB7XG4gICAgdGhyb3cgRXJyb3IoXG4gICAgICBgVW5leHBlY3RlZCBudW1iZXIgb2YgaTE4biBwbGFjZWhvbGRlcnMgKCR7dmFsdWUuZXhwcmVzc2lvbnMubGVuZ3RofSkgZm9yIEJvdW5kVGV4dCB3aXRoICR7dmFsdWUuZXhwcmVzc2lvbnMubGVuZ3RofSBleHByZXNzaW9uc2AsXG4gICAgKTtcbiAgfVxuXG4gIGNvbnN0IHRleHRYcmVmID0gdW5pdC5qb2IuYWxsb2NhdGVYcmVmSWQoKTtcbiAgdW5pdC5jcmVhdGUucHVzaChpci5jcmVhdGVUZXh0T3AodGV4dFhyZWYsICcnLCBpY3VQbGFjZWhvbGRlciwgdGV4dC5zb3VyY2VTcGFuKSk7XG4gIC8vIFRlbXBsYXRlRGVmaW5pdGlvbkJ1aWxkZXIgZG9lcyBub3QgZ2VuZXJhdGUgc291cmNlIG1hcHMgZm9yIHN1Yi1leHByZXNzaW9ucyBpbnNpZGUgYW5cbiAgLy8gaW50ZXJwb2xhdGlvbi4gV2UgY29weSB0aGF0IGJlaGF2aW9yIGluIGNvbXBhdGliaWxpdHkgbW9kZS5cbiAgLy8gVE9ETzogaXMgaXQgYWN0dWFsbHkgY29ycmVjdCB0byBnZW5lcmF0ZSB0aGVzZSBleHRyYSBtYXBzIGluIG1vZGVybiBtb2RlP1xuICBjb25zdCBiYXNlU291cmNlU3BhbiA9IHVuaXQuam9iLmNvbXBhdGliaWxpdHkgPyBudWxsIDogdGV4dC5zb3VyY2VTcGFuO1xuICB1bml0LnVwZGF0ZS5wdXNoKFxuICAgIGlyLmNyZWF0ZUludGVycG9sYXRlVGV4dE9wKFxuICAgICAgdGV4dFhyZWYsXG4gICAgICBuZXcgaXIuSW50ZXJwb2xhdGlvbihcbiAgICAgICAgdmFsdWUuc3RyaW5ncyxcbiAgICAgICAgdmFsdWUuZXhwcmVzc2lvbnMubWFwKChleHByKSA9PiBjb252ZXJ0QXN0KGV4cHIsIHVuaXQuam9iLCBiYXNlU291cmNlU3BhbikpLFxuICAgICAgICBpMThuUGxhY2Vob2xkZXJzLFxuICAgICAgKSxcbiAgICAgIHRleHQuc291cmNlU3BhbixcbiAgICApLFxuICApO1xufVxuXG4vKipcbiAqIEluZ2VzdCBhbiBgQGlmYCBibG9jayBpbnRvIHRoZSBnaXZlbiBgVmlld0NvbXBpbGF0aW9uYC5cbiAqL1xuZnVuY3Rpb24gaW5nZXN0SWZCbG9jayh1bml0OiBWaWV3Q29tcGlsYXRpb25Vbml0LCBpZkJsb2NrOiB0LklmQmxvY2spOiB2b2lkIHtcbiAgbGV0IGZpcnN0WHJlZjogaXIuWHJlZklkIHwgbnVsbCA9IG51bGw7XG4gIGxldCBjb25kaXRpb25zOiBBcnJheTxpci5Db25kaXRpb25hbENhc2VFeHByPiA9IFtdO1xuICBmb3IgKGxldCBpID0gMDsgaSA8IGlmQmxvY2suYnJhbmNoZXMubGVuZ3RoOyBpKyspIHtcbiAgICBjb25zdCBpZkNhc2UgPSBpZkJsb2NrLmJyYW5jaGVzW2ldO1xuICAgIGNvbnN0IGNWaWV3ID0gdW5pdC5qb2IuYWxsb2NhdGVWaWV3KHVuaXQueHJlZik7XG4gICAgY29uc3QgdGFnTmFtZSA9IGluZ2VzdENvbnRyb2xGbG93SW5zZXJ0aW9uUG9pbnQodW5pdCwgY1ZpZXcueHJlZiwgaWZDYXNlKTtcblxuICAgIGlmIChpZkNhc2UuZXhwcmVzc2lvbkFsaWFzICE9PSBudWxsKSB7XG4gICAgICBjVmlldy5jb250ZXh0VmFyaWFibGVzLnNldChpZkNhc2UuZXhwcmVzc2lvbkFsaWFzLm5hbWUsIGlyLkNUWF9SRUYpO1xuICAgIH1cblxuICAgIGxldCBpZkNhc2VJMThuTWV0YTogaTE4bi5CbG9ja1BsYWNlaG9sZGVyIHwgdW5kZWZpbmVkID0gdW5kZWZpbmVkO1xuICAgIGlmIChpZkNhc2UuaTE4biAhPT0gdW5kZWZpbmVkKSB7XG4gICAgICBpZiAoIShpZkNhc2UuaTE4biBpbnN0YW5jZW9mIGkxOG4uQmxvY2tQbGFjZWhvbGRlcikpIHtcbiAgICAgICAgdGhyb3cgRXJyb3IoYFVuaGFuZGxlZCBpMThuIG1ldGFkYXRhIHR5cGUgZm9yIGlmIGJsb2NrOiAke2lmQ2FzZS5pMThuPy5jb25zdHJ1Y3Rvci5uYW1lfWApO1xuICAgICAgfVxuICAgICAgaWZDYXNlSTE4bk1ldGEgPSBpZkNhc2UuaTE4bjtcbiAgICB9XG5cbiAgICBjb25zdCB0ZW1wbGF0ZU9wID0gaXIuY3JlYXRlVGVtcGxhdGVPcChcbiAgICAgIGNWaWV3LnhyZWYsXG4gICAgICBpci5UZW1wbGF0ZUtpbmQuQmxvY2ssXG4gICAgICB0YWdOYW1lLFxuICAgICAgJ0NvbmRpdGlvbmFsJyxcbiAgICAgIGlyLk5hbWVzcGFjZS5IVE1MLFxuICAgICAgaWZDYXNlSTE4bk1ldGEsXG4gICAgICBpZkNhc2Uuc3RhcnRTb3VyY2VTcGFuLFxuICAgICAgaWZDYXNlLnNvdXJjZVNwYW4sXG4gICAgKTtcbiAgICB1bml0LmNyZWF0ZS5wdXNoKHRlbXBsYXRlT3ApO1xuXG4gICAgaWYgKGZpcnN0WHJlZiA9PT0gbnVsbCkge1xuICAgICAgZmlyc3RYcmVmID0gY1ZpZXcueHJlZjtcbiAgICB9XG5cbiAgICBjb25zdCBjYXNlRXhwciA9IGlmQ2FzZS5leHByZXNzaW9uID8gY29udmVydEFzdChpZkNhc2UuZXhwcmVzc2lvbiwgdW5pdC5qb2IsIG51bGwpIDogbnVsbDtcbiAgICBjb25zdCBjb25kaXRpb25hbENhc2VFeHByID0gbmV3IGlyLkNvbmRpdGlvbmFsQ2FzZUV4cHIoXG4gICAgICBjYXNlRXhwcixcbiAgICAgIHRlbXBsYXRlT3AueHJlZixcbiAgICAgIHRlbXBsYXRlT3AuaGFuZGxlLFxuICAgICAgaWZDYXNlLmV4cHJlc3Npb25BbGlhcyxcbiAgICApO1xuICAgIGNvbmRpdGlvbnMucHVzaChjb25kaXRpb25hbENhc2VFeHByKTtcbiAgICBpbmdlc3ROb2RlcyhjVmlldywgaWZDYXNlLmNoaWxkcmVuKTtcbiAgfVxuICB1bml0LnVwZGF0ZS5wdXNoKGlyLmNyZWF0ZUNvbmRpdGlvbmFsT3AoZmlyc3RYcmVmISwgbnVsbCwgY29uZGl0aW9ucywgaWZCbG9jay5zb3VyY2VTcGFuKSk7XG59XG5cbi8qKlxuICogSW5nZXN0IGFuIGBAc3dpdGNoYCBibG9jayBpbnRvIHRoZSBnaXZlbiBgVmlld0NvbXBpbGF0aW9uYC5cbiAqL1xuZnVuY3Rpb24gaW5nZXN0U3dpdGNoQmxvY2sodW5pdDogVmlld0NvbXBpbGF0aW9uVW5pdCwgc3dpdGNoQmxvY2s6IHQuU3dpdGNoQmxvY2spOiB2b2lkIHtcbiAgLy8gRG9uJ3QgaW5nZXN0IGVtcHR5IHN3aXRjaGVzIHNpbmNlIHRoZXkgd29uJ3QgcmVuZGVyIGFueXRoaW5nLlxuICBpZiAoc3dpdGNoQmxvY2suY2FzZXMubGVuZ3RoID09PSAwKSB7XG4gICAgcmV0dXJuO1xuICB9XG5cbiAgbGV0IGZpcnN0WHJlZjogaXIuWHJlZklkIHwgbnVsbCA9IG51bGw7XG4gIGxldCBjb25kaXRpb25zOiBBcnJheTxpci5Db25kaXRpb25hbENhc2VFeHByPiA9IFtdO1xuICBmb3IgKGNvbnN0IHN3aXRjaENhc2Ugb2Ygc3dpdGNoQmxvY2suY2FzZXMpIHtcbiAgICBjb25zdCBjVmlldyA9IHVuaXQuam9iLmFsbG9jYXRlVmlldyh1bml0LnhyZWYpO1xuICAgIGNvbnN0IHRhZ05hbWUgPSBpbmdlc3RDb250cm9sRmxvd0luc2VydGlvblBvaW50KHVuaXQsIGNWaWV3LnhyZWYsIHN3aXRjaENhc2UpO1xuICAgIGxldCBzd2l0Y2hDYXNlSTE4bk1ldGE6IGkxOG4uQmxvY2tQbGFjZWhvbGRlciB8IHVuZGVmaW5lZCA9IHVuZGVmaW5lZDtcbiAgICBpZiAoc3dpdGNoQ2FzZS5pMThuICE9PSB1bmRlZmluZWQpIHtcbiAgICAgIGlmICghKHN3aXRjaENhc2UuaTE4biBpbnN0YW5jZW9mIGkxOG4uQmxvY2tQbGFjZWhvbGRlcikpIHtcbiAgICAgICAgdGhyb3cgRXJyb3IoXG4gICAgICAgICAgYFVuaGFuZGxlZCBpMThuIG1ldGFkYXRhIHR5cGUgZm9yIHN3aXRjaCBibG9jazogJHtzd2l0Y2hDYXNlLmkxOG4/LmNvbnN0cnVjdG9yLm5hbWV9YCxcbiAgICAgICAgKTtcbiAgICAgIH1cbiAgICAgIHN3aXRjaENhc2VJMThuTWV0YSA9IHN3aXRjaENhc2UuaTE4bjtcbiAgICB9XG4gICAgY29uc3QgdGVtcGxhdGVPcCA9IGlyLmNyZWF0ZVRlbXBsYXRlT3AoXG4gICAgICBjVmlldy54cmVmLFxuICAgICAgaXIuVGVtcGxhdGVLaW5kLkJsb2NrLFxuICAgICAgdGFnTmFtZSxcbiAgICAgICdDYXNlJyxcbiAgICAgIGlyLk5hbWVzcGFjZS5IVE1MLFxuICAgICAgc3dpdGNoQ2FzZUkxOG5NZXRhLFxuICAgICAgc3dpdGNoQ2FzZS5zdGFydFNvdXJjZVNwYW4sXG4gICAgICBzd2l0Y2hDYXNlLnNvdXJjZVNwYW4sXG4gICAgKTtcbiAgICB1bml0LmNyZWF0ZS5wdXNoKHRlbXBsYXRlT3ApO1xuICAgIGlmIChmaXJzdFhyZWYgPT09IG51bGwpIHtcbiAgICAgIGZpcnN0WHJlZiA9IGNWaWV3LnhyZWY7XG4gICAgfVxuICAgIGNvbnN0IGNhc2VFeHByID0gc3dpdGNoQ2FzZS5leHByZXNzaW9uXG4gICAgICA/IGNvbnZlcnRBc3Qoc3dpdGNoQ2FzZS5leHByZXNzaW9uLCB1bml0LmpvYiwgc3dpdGNoQmxvY2suc3RhcnRTb3VyY2VTcGFuKVxuICAgICAgOiBudWxsO1xuICAgIGNvbnN0IGNvbmRpdGlvbmFsQ2FzZUV4cHIgPSBuZXcgaXIuQ29uZGl0aW9uYWxDYXNlRXhwcihcbiAgICAgIGNhc2VFeHByLFxuICAgICAgdGVtcGxhdGVPcC54cmVmLFxuICAgICAgdGVtcGxhdGVPcC5oYW5kbGUsXG4gICAgKTtcbiAgICBjb25kaXRpb25zLnB1c2goY29uZGl0aW9uYWxDYXNlRXhwcik7XG4gICAgaW5nZXN0Tm9kZXMoY1ZpZXcsIHN3aXRjaENhc2UuY2hpbGRyZW4pO1xuICB9XG4gIHVuaXQudXBkYXRlLnB1c2goXG4gICAgaXIuY3JlYXRlQ29uZGl0aW9uYWxPcChcbiAgICAgIGZpcnN0WHJlZiEsXG4gICAgICBjb252ZXJ0QXN0KHN3aXRjaEJsb2NrLmV4cHJlc3Npb24sIHVuaXQuam9iLCBudWxsKSxcbiAgICAgIGNvbmRpdGlvbnMsXG4gICAgICBzd2l0Y2hCbG9jay5zb3VyY2VTcGFuLFxuICAgICksXG4gICk7XG59XG5cbmZ1bmN0aW9uIGluZ2VzdERlZmVyVmlldyhcbiAgdW5pdDogVmlld0NvbXBpbGF0aW9uVW5pdCxcbiAgc3VmZml4OiBzdHJpbmcsXG4gIGkxOG5NZXRhOiBpMThuLkkxOG5NZXRhIHwgdW5kZWZpbmVkLFxuICBjaGlsZHJlbj86IHQuTm9kZVtdLFxuICBzb3VyY2VTcGFuPzogUGFyc2VTb3VyY2VTcGFuLFxuKTogaXIuVGVtcGxhdGVPcCB8IG51bGwge1xuICBpZiAoaTE4bk1ldGEgIT09IHVuZGVmaW5lZCAmJiAhKGkxOG5NZXRhIGluc3RhbmNlb2YgaTE4bi5CbG9ja1BsYWNlaG9sZGVyKSkge1xuICAgIHRocm93IEVycm9yKCdVbmhhbmRsZWQgaTE4biBtZXRhZGF0YSB0eXBlIGZvciBkZWZlciBibG9jaycpO1xuICB9XG4gIGlmIChjaGlsZHJlbiA9PT0gdW5kZWZpbmVkKSB7XG4gICAgcmV0dXJuIG51bGw7XG4gIH1cbiAgY29uc3Qgc2Vjb25kYXJ5VmlldyA9IHVuaXQuam9iLmFsbG9jYXRlVmlldyh1bml0LnhyZWYpO1xuICBpbmdlc3ROb2RlcyhzZWNvbmRhcnlWaWV3LCBjaGlsZHJlbik7XG4gIGNvbnN0IHRlbXBsYXRlT3AgPSBpci5jcmVhdGVUZW1wbGF0ZU9wKFxuICAgIHNlY29uZGFyeVZpZXcueHJlZixcbiAgICBpci5UZW1wbGF0ZUtpbmQuQmxvY2ssXG4gICAgbnVsbCxcbiAgICBgRGVmZXIke3N1ZmZpeH1gLFxuICAgIGlyLk5hbWVzcGFjZS5IVE1MLFxuICAgIGkxOG5NZXRhLFxuICAgIHNvdXJjZVNwYW4hLFxuICAgIHNvdXJjZVNwYW4hLFxuICApO1xuICB1bml0LmNyZWF0ZS5wdXNoKHRlbXBsYXRlT3ApO1xuICByZXR1cm4gdGVtcGxhdGVPcDtcbn1cblxuZnVuY3Rpb24gaW5nZXN0RGVmZXJCbG9jayh1bml0OiBWaWV3Q29tcGlsYXRpb25Vbml0LCBkZWZlckJsb2NrOiB0LkRlZmVycmVkQmxvY2spOiB2b2lkIHtcbiAgbGV0IG93blJlc29sdmVyRm46IG8uRXhwcmVzc2lvbiB8IG51bGwgPSBudWxsO1xuXG4gIGlmICh1bml0LmpvYi5kZWZlck1ldGEubW9kZSA9PT0gRGVmZXJCbG9ja0RlcHNFbWl0TW9kZS5QZXJCbG9jaykge1xuICAgIGlmICghdW5pdC5qb2IuZGVmZXJNZXRhLmJsb2Nrcy5oYXMoZGVmZXJCbG9jaykpIHtcbiAgICAgIHRocm93IG5ldyBFcnJvcihcbiAgICAgICAgYEFzc2VydGlvbkVycm9yOiB1bmFibGUgdG8gZmluZCBhIGRlcGVuZGVuY3kgZnVuY3Rpb24gZm9yIHRoaXMgZGVmZXJyZWQgYmxvY2tgLFxuICAgICAgKTtcbiAgICB9XG4gICAgb3duUmVzb2x2ZXJGbiA9IHVuaXQuam9iLmRlZmVyTWV0YS5ibG9ja3MuZ2V0KGRlZmVyQmxvY2spID8/IG51bGw7XG4gIH1cblxuICAvLyBHZW5lcmF0ZSB0aGUgZGVmZXIgbWFpbiB2aWV3IGFuZCBhbGwgc2Vjb25kYXJ5IHZpZXdzLlxuICBjb25zdCBtYWluID0gaW5nZXN0RGVmZXJWaWV3KFxuICAgIHVuaXQsXG4gICAgJycsXG4gICAgZGVmZXJCbG9jay5pMThuLFxuICAgIGRlZmVyQmxvY2suY2hpbGRyZW4sXG4gICAgZGVmZXJCbG9jay5zb3VyY2VTcGFuLFxuICApITtcbiAgY29uc3QgbG9hZGluZyA9IGluZ2VzdERlZmVyVmlldyhcbiAgICB1bml0LFxuICAgICdMb2FkaW5nJyxcbiAgICBkZWZlckJsb2NrLmxvYWRpbmc/LmkxOG4sXG4gICAgZGVmZXJCbG9jay5sb2FkaW5nPy5jaGlsZHJlbixcbiAgICBkZWZlckJsb2NrLmxvYWRpbmc/LnNvdXJjZVNwYW4sXG4gICk7XG4gIGNvbnN0IHBsYWNlaG9sZGVyID0gaW5nZXN0RGVmZXJWaWV3KFxuICAgIHVuaXQsXG4gICAgJ1BsYWNlaG9sZGVyJyxcbiAgICBkZWZlckJsb2NrLnBsYWNlaG9sZGVyPy5pMThuLFxuICAgIGRlZmVyQmxvY2sucGxhY2Vob2xkZXI/LmNoaWxkcmVuLFxuICAgIGRlZmVyQmxvY2sucGxhY2Vob2xkZXI/LnNvdXJjZVNwYW4sXG4gICk7XG4gIGNvbnN0IGVycm9yID0gaW5nZXN0RGVmZXJWaWV3KFxuICAgIHVuaXQsXG4gICAgJ0Vycm9yJyxcbiAgICBkZWZlckJsb2NrLmVycm9yPy5pMThuLFxuICAgIGRlZmVyQmxvY2suZXJyb3I/LmNoaWxkcmVuLFxuICAgIGRlZmVyQmxvY2suZXJyb3I/LnNvdXJjZVNwYW4sXG4gICk7XG5cbiAgLy8gQ3JlYXRlIHRoZSBtYWluIGRlZmVyIG9wLCBhbmQgb3BzIGZvciBhbGwgc2Vjb25kYXJ5IHZpZXdzLlxuICBjb25zdCBkZWZlclhyZWYgPSB1bml0LmpvYi5hbGxvY2F0ZVhyZWZJZCgpO1xuICBjb25zdCBkZWZlck9wID0gaXIuY3JlYXRlRGVmZXJPcChcbiAgICBkZWZlclhyZWYsXG4gICAgbWFpbi54cmVmLFxuICAgIG1haW4uaGFuZGxlLFxuICAgIG93blJlc29sdmVyRm4sXG4gICAgdW5pdC5qb2IuYWxsRGVmZXJyYWJsZURlcHNGbixcbiAgICBkZWZlckJsb2NrLnNvdXJjZVNwYW4sXG4gICk7XG4gIGRlZmVyT3AucGxhY2Vob2xkZXJWaWV3ID0gcGxhY2Vob2xkZXI/LnhyZWYgPz8gbnVsbDtcbiAgZGVmZXJPcC5wbGFjZWhvbGRlclNsb3QgPSBwbGFjZWhvbGRlcj8uaGFuZGxlID8/IG51bGw7XG4gIGRlZmVyT3AubG9hZGluZ1Nsb3QgPSBsb2FkaW5nPy5oYW5kbGUgPz8gbnVsbDtcbiAgZGVmZXJPcC5lcnJvclNsb3QgPSBlcnJvcj8uaGFuZGxlID8/IG51bGw7XG4gIGRlZmVyT3AucGxhY2Vob2xkZXJNaW5pbXVtVGltZSA9IGRlZmVyQmxvY2sucGxhY2Vob2xkZXI/Lm1pbmltdW1UaW1lID8/IG51bGw7XG4gIGRlZmVyT3AubG9hZGluZ01pbmltdW1UaW1lID0gZGVmZXJCbG9jay5sb2FkaW5nPy5taW5pbXVtVGltZSA/PyBudWxsO1xuICBkZWZlck9wLmxvYWRpbmdBZnRlclRpbWUgPSBkZWZlckJsb2NrLmxvYWRpbmc/LmFmdGVyVGltZSA/PyBudWxsO1xuICB1bml0LmNyZWF0ZS5wdXNoKGRlZmVyT3ApO1xuXG4gIC8vIENvbmZpZ3VyZSBhbGwgZGVmZXIgYG9uYCBjb25kaXRpb25zLlxuICAvLyBUT0RPOiByZWZhY3RvciBwcmVmZXRjaCB0cmlnZ2VycyB0byB1c2UgYSBzZXBhcmF0ZSBvcCB0eXBlLCB3aXRoIGEgc2hhcmVkIHN1cGVyY2xhc3MuIFRoaXMgd2lsbFxuICAvLyBtYWtlIGl0IGVhc2llciB0byByZWZhY3RvciBwcmVmZXRjaCBiZWhhdmlvciBpbiB0aGUgZnV0dXJlLlxuICBsZXQgcHJlZmV0Y2ggPSBmYWxzZTtcbiAgbGV0IGRlZmVyT25PcHM6IGlyLkRlZmVyT25PcFtdID0gW107XG4gIGxldCBkZWZlcldoZW5PcHM6IGlyLkRlZmVyV2hlbk9wW10gPSBbXTtcbiAgZm9yIChjb25zdCB0cmlnZ2VycyBvZiBbZGVmZXJCbG9jay50cmlnZ2VycywgZGVmZXJCbG9jay5wcmVmZXRjaFRyaWdnZXJzXSkge1xuICAgIGlmICh0cmlnZ2Vycy5pZGxlICE9PSB1bmRlZmluZWQpIHtcbiAgICAgIGNvbnN0IGRlZmVyT25PcCA9IGlyLmNyZWF0ZURlZmVyT25PcChcbiAgICAgICAgZGVmZXJYcmVmLFxuICAgICAgICB7a2luZDogaXIuRGVmZXJUcmlnZ2VyS2luZC5JZGxlfSxcbiAgICAgICAgcHJlZmV0Y2gsXG4gICAgICAgIHRyaWdnZXJzLmlkbGUuc291cmNlU3BhbixcbiAgICAgICk7XG4gICAgICBkZWZlck9uT3BzLnB1c2goZGVmZXJPbk9wKTtcbiAgICB9XG4gICAgaWYgKHRyaWdnZXJzLmltbWVkaWF0ZSAhPT0gdW5kZWZpbmVkKSB7XG4gICAgICBjb25zdCBkZWZlck9uT3AgPSBpci5jcmVhdGVEZWZlck9uT3AoXG4gICAgICAgIGRlZmVyWHJlZixcbiAgICAgICAge2tpbmQ6IGlyLkRlZmVyVHJpZ2dlcktpbmQuSW1tZWRpYXRlfSxcbiAgICAgICAgcHJlZmV0Y2gsXG4gICAgICAgIHRyaWdnZXJzLmltbWVkaWF0ZS5zb3VyY2VTcGFuLFxuICAgICAgKTtcbiAgICAgIGRlZmVyT25PcHMucHVzaChkZWZlck9uT3ApO1xuICAgIH1cbiAgICBpZiAodHJpZ2dlcnMudGltZXIgIT09IHVuZGVmaW5lZCkge1xuICAgICAgY29uc3QgZGVmZXJPbk9wID0gaXIuY3JlYXRlRGVmZXJPbk9wKFxuICAgICAgICBkZWZlclhyZWYsXG4gICAgICAgIHtraW5kOiBpci5EZWZlclRyaWdnZXJLaW5kLlRpbWVyLCBkZWxheTogdHJpZ2dlcnMudGltZXIuZGVsYXl9LFxuICAgICAgICBwcmVmZXRjaCxcbiAgICAgICAgdHJpZ2dlcnMudGltZXIuc291cmNlU3BhbixcbiAgICAgICk7XG4gICAgICBkZWZlck9uT3BzLnB1c2goZGVmZXJPbk9wKTtcbiAgICB9XG4gICAgaWYgKHRyaWdnZXJzLmhvdmVyICE9PSB1bmRlZmluZWQpIHtcbiAgICAgIGNvbnN0IGRlZmVyT25PcCA9IGlyLmNyZWF0ZURlZmVyT25PcChcbiAgICAgICAgZGVmZXJYcmVmLFxuICAgICAgICB7XG4gICAgICAgICAga2luZDogaXIuRGVmZXJUcmlnZ2VyS2luZC5Ib3ZlcixcbiAgICAgICAgICB0YXJnZXROYW1lOiB0cmlnZ2Vycy5ob3Zlci5yZWZlcmVuY2UsXG4gICAgICAgICAgdGFyZ2V0WHJlZjogbnVsbCxcbiAgICAgICAgICB0YXJnZXRTbG90OiBudWxsLFxuICAgICAgICAgIHRhcmdldFZpZXc6IG51bGwsXG4gICAgICAgICAgdGFyZ2V0U2xvdFZpZXdTdGVwczogbnVsbCxcbiAgICAgICAgfSxcbiAgICAgICAgcHJlZmV0Y2gsXG4gICAgICAgIHRyaWdnZXJzLmhvdmVyLnNvdXJjZVNwYW4sXG4gICAgICApO1xuICAgICAgZGVmZXJPbk9wcy5wdXNoKGRlZmVyT25PcCk7XG4gICAgfVxuICAgIGlmICh0cmlnZ2Vycy5pbnRlcmFjdGlvbiAhPT0gdW5kZWZpbmVkKSB7XG4gICAgICBjb25zdCBkZWZlck9uT3AgPSBpci5jcmVhdGVEZWZlck9uT3AoXG4gICAgICAgIGRlZmVyWHJlZixcbiAgICAgICAge1xuICAgICAgICAgIGtpbmQ6IGlyLkRlZmVyVHJpZ2dlcktpbmQuSW50ZXJhY3Rpb24sXG4gICAgICAgICAgdGFyZ2V0TmFtZTogdHJpZ2dlcnMuaW50ZXJhY3Rpb24ucmVmZXJlbmNlLFxuICAgICAgICAgIHRhcmdldFhyZWY6IG51bGwsXG4gICAgICAgICAgdGFyZ2V0U2xvdDogbnVsbCxcbiAgICAgICAgICB0YXJnZXRWaWV3OiBudWxsLFxuICAgICAgICAgIHRhcmdldFNsb3RWaWV3U3RlcHM6IG51bGwsXG4gICAgICAgIH0sXG4gICAgICAgIHByZWZldGNoLFxuICAgICAgICB0cmlnZ2Vycy5pbnRlcmFjdGlvbi5zb3VyY2VTcGFuLFxuICAgICAgKTtcbiAgICAgIGRlZmVyT25PcHMucHVzaChkZWZlck9uT3ApO1xuICAgIH1cbiAgICBpZiAodHJpZ2dlcnMudmlld3BvcnQgIT09IHVuZGVmaW5lZCkge1xuICAgICAgY29uc3QgZGVmZXJPbk9wID0gaXIuY3JlYXRlRGVmZXJPbk9wKFxuICAgICAgICBkZWZlclhyZWYsXG4gICAgICAgIHtcbiAgICAgICAgICBraW5kOiBpci5EZWZlclRyaWdnZXJLaW5kLlZpZXdwb3J0LFxuICAgICAgICAgIHRhcmdldE5hbWU6IHRyaWdnZXJzLnZpZXdwb3J0LnJlZmVyZW5jZSxcbiAgICAgICAgICB0YXJnZXRYcmVmOiBudWxsLFxuICAgICAgICAgIHRhcmdldFNsb3Q6IG51bGwsXG4gICAgICAgICAgdGFyZ2V0VmlldzogbnVsbCxcbiAgICAgICAgICB0YXJnZXRTbG90Vmlld1N0ZXBzOiBudWxsLFxuICAgICAgICB9LFxuICAgICAgICBwcmVmZXRjaCxcbiAgICAgICAgdHJpZ2dlcnMudmlld3BvcnQuc291cmNlU3BhbixcbiAgICAgICk7XG4gICAgICBkZWZlck9uT3BzLnB1c2goZGVmZXJPbk9wKTtcbiAgICB9XG4gICAgaWYgKHRyaWdnZXJzLndoZW4gIT09IHVuZGVmaW5lZCkge1xuICAgICAgaWYgKHRyaWdnZXJzLndoZW4udmFsdWUgaW5zdGFuY2VvZiBlLkludGVycG9sYXRpb24pIHtcbiAgICAgICAgLy8gVGVtcGxhdGVEZWZpbml0aW9uQnVpbGRlciBzdXBwb3J0cyB0aGlzIGNhc2UsIGJ1dCBpdCdzIHZlcnkgc3RyYW5nZSB0byBtZS4gV2hhdCB3b3VsZCBpdFxuICAgICAgICAvLyBldmVuIG1lYW4/XG4gICAgICAgIHRocm93IG5ldyBFcnJvcihgVW5leHBlY3RlZCBpbnRlcnBvbGF0aW9uIGluIGRlZmVyIGJsb2NrIHdoZW4gdHJpZ2dlcmApO1xuICAgICAgfVxuICAgICAgY29uc3QgZGVmZXJPbk9wID0gaXIuY3JlYXRlRGVmZXJXaGVuT3AoXG4gICAgICAgIGRlZmVyWHJlZixcbiAgICAgICAgY29udmVydEFzdCh0cmlnZ2Vycy53aGVuLnZhbHVlLCB1bml0LmpvYiwgdHJpZ2dlcnMud2hlbi5zb3VyY2VTcGFuKSxcbiAgICAgICAgcHJlZmV0Y2gsXG4gICAgICAgIHRyaWdnZXJzLndoZW4uc291cmNlU3BhbixcbiAgICAgICk7XG4gICAgICBkZWZlcldoZW5PcHMucHVzaChkZWZlck9uT3ApO1xuICAgIH1cblxuICAgIC8vIElmIG5vIChub24tcHJlZmV0Y2hpbmcpIGRlZmVyIHRyaWdnZXJzIHdlcmUgcHJvdmlkZWQsIGRlZmF1bHQgdG8gYGlkbGVgLlxuICAgIGlmIChkZWZlck9uT3BzLmxlbmd0aCA9PT0gMCAmJiBkZWZlcldoZW5PcHMubGVuZ3RoID09PSAwKSB7XG4gICAgICBkZWZlck9uT3BzLnB1c2goXG4gICAgICAgIGlyLmNyZWF0ZURlZmVyT25PcChkZWZlclhyZWYsIHtraW5kOiBpci5EZWZlclRyaWdnZXJLaW5kLklkbGV9LCBmYWxzZSwgbnVsbCEpLFxuICAgICAgKTtcbiAgICB9XG4gICAgcHJlZmV0Y2ggPSB0cnVlO1xuICB9XG5cbiAgdW5pdC5jcmVhdGUucHVzaChkZWZlck9uT3BzKTtcbiAgdW5pdC51cGRhdGUucHVzaChkZWZlcldoZW5PcHMpO1xufVxuXG5mdW5jdGlvbiBpbmdlc3RJY3UodW5pdDogVmlld0NvbXBpbGF0aW9uVW5pdCwgaWN1OiB0LkljdSkge1xuICBpZiAoaWN1LmkxOG4gaW5zdGFuY2VvZiBpMThuLk1lc3NhZ2UgJiYgaXNTaW5nbGVJMThuSWN1KGljdS5pMThuKSkge1xuICAgIGNvbnN0IHhyZWYgPSB1bml0LmpvYi5hbGxvY2F0ZVhyZWZJZCgpO1xuICAgIHVuaXQuY3JlYXRlLnB1c2goaXIuY3JlYXRlSWN1U3RhcnRPcCh4cmVmLCBpY3UuaTE4biwgaWN1RnJvbUkxOG5NZXNzYWdlKGljdS5pMThuKS5uYW1lLCBudWxsISkpO1xuICAgIGZvciAoY29uc3QgW3BsYWNlaG9sZGVyLCB0ZXh0XSBvZiBPYmplY3QuZW50cmllcyh7Li4uaWN1LnZhcnMsIC4uLmljdS5wbGFjZWhvbGRlcnN9KSkge1xuICAgICAgaWYgKHRleHQgaW5zdGFuY2VvZiB0LkJvdW5kVGV4dCkge1xuICAgICAgICBpbmdlc3RCb3VuZFRleHQodW5pdCwgdGV4dCwgcGxhY2Vob2xkZXIpO1xuICAgICAgfSBlbHNlIHtcbiAgICAgICAgaW5nZXN0VGV4dCh1bml0LCB0ZXh0LCBwbGFjZWhvbGRlcik7XG4gICAgICB9XG4gICAgfVxuICAgIHVuaXQuY3JlYXRlLnB1c2goaXIuY3JlYXRlSWN1RW5kT3AoeHJlZikpO1xuICB9IGVsc2Uge1xuICAgIHRocm93IEVycm9yKGBVbmhhbmRsZWQgaTE4biBtZXRhZGF0YSB0eXBlIGZvciBJQ1U6ICR7aWN1LmkxOG4/LmNvbnN0cnVjdG9yLm5hbWV9YCk7XG4gIH1cbn1cblxuLyoqXG4gKiBJbmdlc3QgYW4gYEBmb3JgIGJsb2NrIGludG8gdGhlIGdpdmVuIGBWaWV3Q29tcGlsYXRpb25gLlxuICovXG5mdW5jdGlvbiBpbmdlc3RGb3JCbG9jayh1bml0OiBWaWV3Q29tcGlsYXRpb25Vbml0LCBmb3JCbG9jazogdC5Gb3JMb29wQmxvY2spOiB2b2lkIHtcbiAgY29uc3QgcmVwZWF0ZXJWaWV3ID0gdW5pdC5qb2IuYWxsb2NhdGVWaWV3KHVuaXQueHJlZik7XG5cbiAgLy8gV2UgY29weSBUZW1wbGF0ZURlZmluaXRpb25CdWlsZGVyJ3Mgc2NoZW1lIG9mIGNyZWF0aW5nIG5hbWVzIGZvciBgJGNvdW50YCBhbmQgYCRpbmRleGBcbiAgLy8gdGhhdCBhcmUgc3VmZml4ZWQgd2l0aCBzcGVjaWFsIGluZm9ybWF0aW9uLCB0byBkaXNhbWJpZ3VhdGUgd2hpY2ggbGV2ZWwgb2YgbmVzdGVkIGxvb3BcbiAgLy8gdGhlIGJlbG93IGFsaWFzZXMgcmVmZXIgdG8uXG4gIC8vIFRPRE86IFdlIHNob3VsZCByZWZhY3RvciBUZW1wbGF0ZSBQaXBlbGluZSdzIHZhcmlhYmxlIHBoYXNlcyB0byBncmFjZWZ1bGx5IGhhbmRsZVxuICAvLyBzaGFkb3dpbmcsIGFuZCBhcmJpdHJhcmlseSBtYW55IGxldmVscyBvZiB2YXJpYWJsZXMgZGVwZW5kaW5nIG9uIGVhY2ggb3RoZXIuXG4gIGNvbnN0IGluZGV4TmFtZSA9IGDJtSRpbmRleF8ke3JlcGVhdGVyVmlldy54cmVmfWA7XG4gIGNvbnN0IGNvdW50TmFtZSA9IGDJtSRjb3VudF8ke3JlcGVhdGVyVmlldy54cmVmfWA7XG4gIGNvbnN0IGluZGV4VmFyTmFtZXMgPSBuZXcgU2V0PHN0cmluZz4oKTtcblxuICAvLyBTZXQgYWxsIHRoZSBjb250ZXh0IHZhcmlhYmxlcyBhbmQgYWxpYXNlcyBhdmFpbGFibGUgaW4gdGhlIHJlcGVhdGVyLlxuICByZXBlYXRlclZpZXcuY29udGV4dFZhcmlhYmxlcy5zZXQoZm9yQmxvY2suaXRlbS5uYW1lLCBmb3JCbG9jay5pdGVtLnZhbHVlKTtcblxuICBmb3IgKGNvbnN0IHZhcmlhYmxlIG9mIGZvckJsb2NrLmNvbnRleHRWYXJpYWJsZXMpIHtcbiAgICBpZiAodmFyaWFibGUudmFsdWUgPT09ICckaW5kZXgnKSB7XG4gICAgICBpbmRleFZhck5hbWVzLmFkZCh2YXJpYWJsZS5uYW1lKTtcbiAgICB9XG4gICAgaWYgKHZhcmlhYmxlLm5hbWUgPT09ICckaW5kZXgnKSB7XG4gICAgICByZXBlYXRlclZpZXcuY29udGV4dFZhcmlhYmxlcy5zZXQoJyRpbmRleCcsIHZhcmlhYmxlLnZhbHVlKS5zZXQoaW5kZXhOYW1lLCB2YXJpYWJsZS52YWx1ZSk7XG4gICAgfSBlbHNlIGlmICh2YXJpYWJsZS5uYW1lID09PSAnJGNvdW50Jykge1xuICAgICAgcmVwZWF0ZXJWaWV3LmNvbnRleHRWYXJpYWJsZXMuc2V0KCckY291bnQnLCB2YXJpYWJsZS52YWx1ZSkuc2V0KGNvdW50TmFtZSwgdmFyaWFibGUudmFsdWUpO1xuICAgIH0gZWxzZSB7XG4gICAgICByZXBlYXRlclZpZXcuYWxpYXNlcy5hZGQoe1xuICAgICAgICBraW5kOiBpci5TZW1hbnRpY1ZhcmlhYmxlS2luZC5BbGlhcyxcbiAgICAgICAgbmFtZTogbnVsbCxcbiAgICAgICAgaWRlbnRpZmllcjogdmFyaWFibGUubmFtZSxcbiAgICAgICAgZXhwcmVzc2lvbjogZ2V0Q29tcHV0ZWRGb3JMb29wVmFyaWFibGVFeHByZXNzaW9uKHZhcmlhYmxlLCBpbmRleE5hbWUsIGNvdW50TmFtZSksXG4gICAgICB9KTtcbiAgICB9XG4gIH1cblxuICBjb25zdCBzb3VyY2VTcGFuID0gY29udmVydFNvdXJjZVNwYW4oZm9yQmxvY2sudHJhY2tCeS5zcGFuLCBmb3JCbG9jay5zb3VyY2VTcGFuKTtcbiAgY29uc3QgdHJhY2sgPSBjb252ZXJ0QXN0KGZvckJsb2NrLnRyYWNrQnksIHVuaXQuam9iLCBzb3VyY2VTcGFuKTtcblxuICBpbmdlc3ROb2RlcyhyZXBlYXRlclZpZXcsIGZvckJsb2NrLmNoaWxkcmVuKTtcblxuICBsZXQgZW1wdHlWaWV3OiBWaWV3Q29tcGlsYXRpb25Vbml0IHwgbnVsbCA9IG51bGw7XG4gIGxldCBlbXB0eVRhZ05hbWU6IHN0cmluZyB8IG51bGwgPSBudWxsO1xuICBpZiAoZm9yQmxvY2suZW1wdHkgIT09IG51bGwpIHtcbiAgICBlbXB0eVZpZXcgPSB1bml0LmpvYi5hbGxvY2F0ZVZpZXcodW5pdC54cmVmKTtcbiAgICBpbmdlc3ROb2RlcyhlbXB0eVZpZXcsIGZvckJsb2NrLmVtcHR5LmNoaWxkcmVuKTtcbiAgICBlbXB0eVRhZ05hbWUgPSBpbmdlc3RDb250cm9sRmxvd0luc2VydGlvblBvaW50KHVuaXQsIGVtcHR5Vmlldy54cmVmLCBmb3JCbG9jay5lbXB0eSk7XG4gIH1cblxuICBjb25zdCB2YXJOYW1lczogaXIuUmVwZWF0ZXJWYXJOYW1lcyA9IHtcbiAgICAkaW5kZXg6IGluZGV4VmFyTmFtZXMsXG4gICAgJGltcGxpY2l0OiBmb3JCbG9jay5pdGVtLm5hbWUsXG4gIH07XG5cbiAgaWYgKGZvckJsb2NrLmkxOG4gIT09IHVuZGVmaW5lZCAmJiAhKGZvckJsb2NrLmkxOG4gaW5zdGFuY2VvZiBpMThuLkJsb2NrUGxhY2Vob2xkZXIpKSB7XG4gICAgdGhyb3cgRXJyb3IoJ0Fzc2VydGlvbkVycm9yOiBVbmhhbmRsZWQgaTE4biBtZXRhZGF0YSB0eXBlIG9yIEBmb3InKTtcbiAgfVxuICBpZiAoXG4gICAgZm9yQmxvY2suZW1wdHk/LmkxOG4gIT09IHVuZGVmaW5lZCAmJlxuICAgICEoZm9yQmxvY2suZW1wdHkuaTE4biBpbnN0YW5jZW9mIGkxOG4uQmxvY2tQbGFjZWhvbGRlcilcbiAgKSB7XG4gICAgdGhyb3cgRXJyb3IoJ0Fzc2VydGlvbkVycm9yOiBVbmhhbmRsZWQgaTE4biBtZXRhZGF0YSB0eXBlIG9yIEBlbXB0eScpO1xuICB9XG4gIGNvbnN0IGkxOG5QbGFjZWhvbGRlciA9IGZvckJsb2NrLmkxOG47XG4gIGNvbnN0IGVtcHR5STE4blBsYWNlaG9sZGVyID0gZm9yQmxvY2suZW1wdHk/LmkxOG47XG5cbiAgY29uc3QgdGFnTmFtZSA9IGluZ2VzdENvbnRyb2xGbG93SW5zZXJ0aW9uUG9pbnQodW5pdCwgcmVwZWF0ZXJWaWV3LnhyZWYsIGZvckJsb2NrKTtcbiAgY29uc3QgcmVwZWF0ZXJDcmVhdGUgPSBpci5jcmVhdGVSZXBlYXRlckNyZWF0ZU9wKFxuICAgIHJlcGVhdGVyVmlldy54cmVmLFxuICAgIGVtcHR5Vmlldz8ueHJlZiA/PyBudWxsLFxuICAgIHRhZ05hbWUsXG4gICAgdHJhY2ssXG4gICAgdmFyTmFtZXMsXG4gICAgZW1wdHlUYWdOYW1lLFxuICAgIGkxOG5QbGFjZWhvbGRlcixcbiAgICBlbXB0eUkxOG5QbGFjZWhvbGRlcixcbiAgICBmb3JCbG9jay5zdGFydFNvdXJjZVNwYW4sXG4gICAgZm9yQmxvY2suc291cmNlU3BhbixcbiAgKTtcbiAgdW5pdC5jcmVhdGUucHVzaChyZXBlYXRlckNyZWF0ZSk7XG5cbiAgY29uc3QgZXhwcmVzc2lvbiA9IGNvbnZlcnRBc3QoXG4gICAgZm9yQmxvY2suZXhwcmVzc2lvbixcbiAgICB1bml0LmpvYixcbiAgICBjb252ZXJ0U291cmNlU3Bhbihmb3JCbG9jay5leHByZXNzaW9uLnNwYW4sIGZvckJsb2NrLnNvdXJjZVNwYW4pLFxuICApO1xuICBjb25zdCByZXBlYXRlciA9IGlyLmNyZWF0ZVJlcGVhdGVyT3AoXG4gICAgcmVwZWF0ZXJDcmVhdGUueHJlZixcbiAgICByZXBlYXRlckNyZWF0ZS5oYW5kbGUsXG4gICAgZXhwcmVzc2lvbixcbiAgICBmb3JCbG9jay5zb3VyY2VTcGFuLFxuICApO1xuICB1bml0LnVwZGF0ZS5wdXNoKHJlcGVhdGVyKTtcbn1cblxuLyoqXG4gKiBHZXRzIGFuIGV4cHJlc3Npb24gdGhhdCByZXByZXNlbnRzIGEgdmFyaWFibGUgaW4gYW4gYEBmb3JgIGxvb3AuXG4gKiBAcGFyYW0gdmFyaWFibGUgQVNUIHJlcHJlc2VudGluZyB0aGUgdmFyaWFibGUuXG4gKiBAcGFyYW0gaW5kZXhOYW1lIExvb3Atc3BlY2lmaWMgbmFtZSBmb3IgYCRpbmRleGAuXG4gKiBAcGFyYW0gY291bnROYW1lIExvb3Atc3BlY2lmaWMgbmFtZSBmb3IgYCRjb3VudGAuXG4gKi9cbmZ1bmN0aW9uIGdldENvbXB1dGVkRm9yTG9vcFZhcmlhYmxlRXhwcmVzc2lvbihcbiAgdmFyaWFibGU6IHQuVmFyaWFibGUsXG4gIGluZGV4TmFtZTogc3RyaW5nLFxuICBjb3VudE5hbWU6IHN0cmluZyxcbik6IG8uRXhwcmVzc2lvbiB7XG4gIHN3aXRjaCAodmFyaWFibGUudmFsdWUpIHtcbiAgICBjYXNlICckaW5kZXgnOlxuICAgICAgcmV0dXJuIG5ldyBpci5MZXhpY2FsUmVhZEV4cHIoaW5kZXhOYW1lKTtcblxuICAgIGNhc2UgJyRjb3VudCc6XG4gICAgICByZXR1cm4gbmV3IGlyLkxleGljYWxSZWFkRXhwcihjb3VudE5hbWUpO1xuXG4gICAgY2FzZSAnJGZpcnN0JzpcbiAgICAgIHJldHVybiBuZXcgaXIuTGV4aWNhbFJlYWRFeHByKGluZGV4TmFtZSkuaWRlbnRpY2FsKG8ubGl0ZXJhbCgwKSk7XG5cbiAgICBjYXNlICckbGFzdCc6XG4gICAgICByZXR1cm4gbmV3IGlyLkxleGljYWxSZWFkRXhwcihpbmRleE5hbWUpLmlkZW50aWNhbChcbiAgICAgICAgbmV3IGlyLkxleGljYWxSZWFkRXhwcihjb3VudE5hbWUpLm1pbnVzKG8ubGl0ZXJhbCgxKSksXG4gICAgICApO1xuXG4gICAgY2FzZSAnJGV2ZW4nOlxuICAgICAgcmV0dXJuIG5ldyBpci5MZXhpY2FsUmVhZEV4cHIoaW5kZXhOYW1lKS5tb2R1bG8oby5saXRlcmFsKDIpKS5pZGVudGljYWwoby5saXRlcmFsKDApKTtcblxuICAgIGNhc2UgJyRvZGQnOlxuICAgICAgcmV0dXJuIG5ldyBpci5MZXhpY2FsUmVhZEV4cHIoaW5kZXhOYW1lKS5tb2R1bG8oby5saXRlcmFsKDIpKS5ub3RJZGVudGljYWwoby5saXRlcmFsKDApKTtcblxuICAgIGRlZmF1bHQ6XG4gICAgICB0aHJvdyBuZXcgRXJyb3IoYEFzc2VydGlvbkVycm9yOiB1bmtub3duIEBmb3IgbG9vcCB2YXJpYWJsZSAke3ZhcmlhYmxlLnZhbHVlfWApO1xuICB9XG59XG5cbmZ1bmN0aW9uIGluZ2VzdExldERlY2xhcmF0aW9uKHVuaXQ6IFZpZXdDb21waWxhdGlvblVuaXQsIG5vZGU6IHQuTGV0RGVjbGFyYXRpb24pIHtcbiAgY29uc3QgdGFyZ2V0ID0gdW5pdC5qb2IuYWxsb2NhdGVYcmVmSWQoKTtcblxuICB1bml0LmNyZWF0ZS5wdXNoKGlyLmNyZWF0ZURlY2xhcmVMZXRPcCh0YXJnZXQsIG5vZGUubmFtZSwgbm9kZS5zb3VyY2VTcGFuKSk7XG4gIHVuaXQudXBkYXRlLnB1c2goXG4gICAgaXIuY3JlYXRlU3RvcmVMZXRPcChcbiAgICAgIHRhcmdldCxcbiAgICAgIG5vZGUubmFtZSxcbiAgICAgIGNvbnZlcnRBc3Qobm9kZS52YWx1ZSwgdW5pdC5qb2IsIG5vZGUudmFsdWVTcGFuKSxcbiAgICAgIG5vZGUuc291cmNlU3BhbixcbiAgICApLFxuICApO1xufVxuXG4vKipcbiAqIENvbnZlcnQgYSB0ZW1wbGF0ZSBBU1QgZXhwcmVzc2lvbiBpbnRvIGFuIG91dHB1dCBBU1QgZXhwcmVzc2lvbi5cbiAqL1xuZnVuY3Rpb24gY29udmVydEFzdChcbiAgYXN0OiBlLkFTVCxcbiAgam9iOiBDb21waWxhdGlvbkpvYixcbiAgYmFzZVNvdXJjZVNwYW46IFBhcnNlU291cmNlU3BhbiB8IG51bGwsXG4pOiBvLkV4cHJlc3Npb24ge1xuICBpZiAoYXN0IGluc3RhbmNlb2YgZS5BU1RXaXRoU291cmNlKSB7XG4gICAgcmV0dXJuIGNvbnZlcnRBc3QoYXN0LmFzdCwgam9iLCBiYXNlU291cmNlU3Bhbik7XG4gIH0gZWxzZSBpZiAoYXN0IGluc3RhbmNlb2YgZS5Qcm9wZXJ0eVJlYWQpIHtcbiAgICBjb25zdCBpc1RoaXNSZWNlaXZlciA9IGFzdC5yZWNlaXZlciBpbnN0YW5jZW9mIGUuVGhpc1JlY2VpdmVyO1xuICAgIC8vIFdoZXRoZXIgdGhpcyBpcyBhbiBpbXBsaWNpdCByZWNlaXZlciwgKmV4Y2x1ZGluZyogZXhwbGljaXQgcmVhZHMgb2YgYHRoaXNgLlxuICAgIGNvbnN0IGlzSW1wbGljaXRSZWNlaXZlciA9XG4gICAgICBhc3QucmVjZWl2ZXIgaW5zdGFuY2VvZiBlLkltcGxpY2l0UmVjZWl2ZXIgJiYgIShhc3QucmVjZWl2ZXIgaW5zdGFuY2VvZiBlLlRoaXNSZWNlaXZlcik7XG4gICAgLy8gV2hldGhlciB0aGUgIG5hbWUgb2YgdGhlIHJlYWQgaXMgYSBub2RlIHRoYXQgc2hvdWxkIGJlIG5ldmVyIHJldGFpbiBpdHMgZXhwbGljaXQgdGhpc1xuICAgIC8vIHJlY2VpdmVyLlxuICAgIGNvbnN0IGlzU3BlY2lhbE5vZGUgPSBhc3QubmFtZSA9PT0gJyRhbnknIHx8IGFzdC5uYW1lID09PSAnJGV2ZW50JztcbiAgICAvLyBUT0RPOiBUaGUgbW9zdCBzZW5zaWJsZSBjb25kaXRpb24gaGVyZSB3b3VsZCBiZSBzaW1wbHkgYGlzSW1wbGljaXRSZWNlaXZlcmAsIHRvIGNvbnZlcnQgb25seVxuICAgIC8vIGFjdHVhbCBpbXBsaWNpdCBgdGhpc2AgcmVhZHMsIGFuZCBub3QgZXhwbGljaXQgb25lcy4gSG93ZXZlciwgVGVtcGxhdGVEZWZpbml0aW9uQnVpbGRlciAoYW5kXG4gICAgLy8gdGhlIFR5cGVjaGVjayBibG9jayEpIGJvdGggaGF2ZSB0aGUgc2FtZSBidWcsIGluIHdoaWNoIHRoZXkgYWxzbyBjb25zaWRlciBleHBsaWNpdCBgdGhpc2BcbiAgICAvLyByZWFkcyB0byBiZSBpbXBsaWNpdC4gVGhpcyBjYXVzZXMgcHJvYmxlbXMgd2hlbiB0aGUgZXhwbGljaXQgYHRoaXNgIHJlYWQgaXMgaW5zaWRlIGFcbiAgICAvLyB0ZW1wbGF0ZSB3aXRoIGEgY29udGV4dCB0aGF0IGFsc28gcHJvdmlkZXMgdGhlIHZhcmlhYmxlIG5hbWUgYmVpbmcgcmVhZDpcbiAgICAvLyBgYGBcbiAgICAvLyA8bmctdGVtcGxhdGUgbGV0LWE+e3t0aGlzLmF9fTwvbmctdGVtcGxhdGU+XG4gICAgLy8gYGBgXG4gICAgLy8gVGhlIHdob2xlIHBvaW50IG9mIHRoZSBleHBsaWNpdCBgdGhpc2Agd2FzIHRvIGFjY2VzcyB0aGUgY2xhc3MgcHJvcGVydHksIGJ1dCBUREIgYW5kIHRoZVxuICAgIC8vIGN1cnJlbnQgVENCIHRyZWF0IHRoZSByZWFkIGFzIGltcGxpY2l0LCBhbmQgZ2l2ZSB5b3UgdGhlIGNvbnRleHQgcHJvcGVydHkgaW5zdGVhZCFcbiAgICAvL1xuICAgIC8vIEZvciBub3csIHdlIGVtdWxhdGUgdGhpcyBvbGQgYmVoYXZpb3IgYnkgYWdncmVzc2l2ZWx5IGNvbnZlcnRpbmcgZXhwbGljaXQgcmVhZHMgdG8gdG9cbiAgICAvLyBpbXBsaWNpdCByZWFkcywgZXhjZXB0IGZvciB0aGUgc3BlY2lhbCBjYXNlcyB0aGF0IFREQiBhbmQgdGhlIGN1cnJlbnQgVENCIHByb3RlY3QuIEhvd2V2ZXIsXG4gICAgLy8gaXQgd291bGQgYmUgYW4gaW1wcm92ZW1lbnQgdG8gZml4IHRoaXMuXG4gICAgLy9cbiAgICAvLyBTZWUgYWxzbyB0aGUgY29ycmVzcG9uZGluZyBjb21tZW50IGZvciB0aGUgVENCLCBpbiBgdHlwZV9jaGVja19ibG9jay50c2AuXG4gICAgaWYgKGlzSW1wbGljaXRSZWNlaXZlciB8fCAoaXNUaGlzUmVjZWl2ZXIgJiYgIWlzU3BlY2lhbE5vZGUpKSB7XG4gICAgICByZXR1cm4gbmV3IGlyLkxleGljYWxSZWFkRXhwcihhc3QubmFtZSk7XG4gICAgfSBlbHNlIHtcbiAgICAgIHJldHVybiBuZXcgby5SZWFkUHJvcEV4cHIoXG4gICAgICAgIGNvbnZlcnRBc3QoYXN0LnJlY2VpdmVyLCBqb2IsIGJhc2VTb3VyY2VTcGFuKSxcbiAgICAgICAgYXN0Lm5hbWUsXG4gICAgICAgIG51bGwsXG4gICAgICAgIGNvbnZlcnRTb3VyY2VTcGFuKGFzdC5zcGFuLCBiYXNlU291cmNlU3BhbiksXG4gICAgICApO1xuICAgIH1cbiAgfSBlbHNlIGlmIChhc3QgaW5zdGFuY2VvZiBlLlByb3BlcnR5V3JpdGUpIHtcbiAgICBpZiAoYXN0LnJlY2VpdmVyIGluc3RhbmNlb2YgZS5JbXBsaWNpdFJlY2VpdmVyKSB7XG4gICAgICByZXR1cm4gbmV3IG8uV3JpdGVQcm9wRXhwcihcbiAgICAgICAgLy8gVE9ETzogSXMgaXQgY29ycmVjdCB0byBhbHdheXMgdXNlIHRoZSByb290IGNvbnRleHQgaW4gcGxhY2Ugb2YgdGhlIGltcGxpY2l0IHJlY2VpdmVyP1xuICAgICAgICBuZXcgaXIuQ29udGV4dEV4cHIoam9iLnJvb3QueHJlZiksXG4gICAgICAgIGFzdC5uYW1lLFxuICAgICAgICBjb252ZXJ0QXN0KGFzdC52YWx1ZSwgam9iLCBiYXNlU291cmNlU3BhbiksXG4gICAgICAgIG51bGwsXG4gICAgICAgIGNvbnZlcnRTb3VyY2VTcGFuKGFzdC5zcGFuLCBiYXNlU291cmNlU3BhbiksXG4gICAgICApO1xuICAgIH1cbiAgICByZXR1cm4gbmV3IG8uV3JpdGVQcm9wRXhwcihcbiAgICAgIGNvbnZlcnRBc3QoYXN0LnJlY2VpdmVyLCBqb2IsIGJhc2VTb3VyY2VTcGFuKSxcbiAgICAgIGFzdC5uYW1lLFxuICAgICAgY29udmVydEFzdChhc3QudmFsdWUsIGpvYiwgYmFzZVNvdXJjZVNwYW4pLFxuICAgICAgdW5kZWZpbmVkLFxuICAgICAgY29udmVydFNvdXJjZVNwYW4oYXN0LnNwYW4sIGJhc2VTb3VyY2VTcGFuKSxcbiAgICApO1xuICB9IGVsc2UgaWYgKGFzdCBpbnN0YW5jZW9mIGUuS2V5ZWRXcml0ZSkge1xuICAgIHJldHVybiBuZXcgby5Xcml0ZUtleUV4cHIoXG4gICAgICBjb252ZXJ0QXN0KGFzdC5yZWNlaXZlciwgam9iLCBiYXNlU291cmNlU3BhbiksXG4gICAgICBjb252ZXJ0QXN0KGFzdC5rZXksIGpvYiwgYmFzZVNvdXJjZVNwYW4pLFxuICAgICAgY29udmVydEFzdChhc3QudmFsdWUsIGpvYiwgYmFzZVNvdXJjZVNwYW4pLFxuICAgICAgdW5kZWZpbmVkLFxuICAgICAgY29udmVydFNvdXJjZVNwYW4oYXN0LnNwYW4sIGJhc2VTb3VyY2VTcGFuKSxcbiAgICApO1xuICB9IGVsc2UgaWYgKGFzdCBpbnN0YW5jZW9mIGUuQ2FsbCkge1xuICAgIGlmIChhc3QucmVjZWl2ZXIgaW5zdGFuY2VvZiBlLkltcGxpY2l0UmVjZWl2ZXIpIHtcbiAgICAgIHRocm93IG5ldyBFcnJvcihgVW5leHBlY3RlZCBJbXBsaWNpdFJlY2VpdmVyYCk7XG4gICAgfSBlbHNlIHtcbiAgICAgIHJldHVybiBuZXcgby5JbnZva2VGdW5jdGlvbkV4cHIoXG4gICAgICAgIGNvbnZlcnRBc3QoYXN0LnJlY2VpdmVyLCBqb2IsIGJhc2VTb3VyY2VTcGFuKSxcbiAgICAgICAgYXN0LmFyZ3MubWFwKChhcmcpID0+IGNvbnZlcnRBc3QoYXJnLCBqb2IsIGJhc2VTb3VyY2VTcGFuKSksXG4gICAgICAgIHVuZGVmaW5lZCxcbiAgICAgICAgY29udmVydFNvdXJjZVNwYW4oYXN0LnNwYW4sIGJhc2VTb3VyY2VTcGFuKSxcbiAgICAgICk7XG4gICAgfVxuICB9IGVsc2UgaWYgKGFzdCBpbnN0YW5jZW9mIGUuTGl0ZXJhbFByaW1pdGl2ZSkge1xuICAgIHJldHVybiBvLmxpdGVyYWwoYXN0LnZhbHVlLCB1bmRlZmluZWQsIGNvbnZlcnRTb3VyY2VTcGFuKGFzdC5zcGFuLCBiYXNlU291cmNlU3BhbikpO1xuICB9IGVsc2UgaWYgKGFzdCBpbnN0YW5jZW9mIGUuVW5hcnkpIHtcbiAgICBzd2l0Y2ggKGFzdC5vcGVyYXRvcikge1xuICAgICAgY2FzZSAnKyc6XG4gICAgICAgIHJldHVybiBuZXcgby5VbmFyeU9wZXJhdG9yRXhwcihcbiAgICAgICAgICBvLlVuYXJ5T3BlcmF0b3IuUGx1cyxcbiAgICAgICAgICBjb252ZXJ0QXN0KGFzdC5leHByLCBqb2IsIGJhc2VTb3VyY2VTcGFuKSxcbiAgICAgICAgICB1bmRlZmluZWQsXG4gICAgICAgICAgY29udmVydFNvdXJjZVNwYW4oYXN0LnNwYW4sIGJhc2VTb3VyY2VTcGFuKSxcbiAgICAgICAgKTtcbiAgICAgIGNhc2UgJy0nOlxuICAgICAgICByZXR1cm4gbmV3IG8uVW5hcnlPcGVyYXRvckV4cHIoXG4gICAgICAgICAgby5VbmFyeU9wZXJhdG9yLk1pbnVzLFxuICAgICAgICAgIGNvbnZlcnRBc3QoYXN0LmV4cHIsIGpvYiwgYmFzZVNvdXJjZVNwYW4pLFxuICAgICAgICAgIHVuZGVmaW5lZCxcbiAgICAgICAgICBjb252ZXJ0U291cmNlU3Bhbihhc3Quc3BhbiwgYmFzZVNvdXJjZVNwYW4pLFxuICAgICAgICApO1xuICAgICAgZGVmYXVsdDpcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKGBBc3NlcnRpb25FcnJvcjogdW5rbm93biB1bmFyeSBvcGVyYXRvciAke2FzdC5vcGVyYXRvcn1gKTtcbiAgICB9XG4gIH0gZWxzZSBpZiAoYXN0IGluc3RhbmNlb2YgZS5CaW5hcnkpIHtcbiAgICBjb25zdCBvcGVyYXRvciA9IEJJTkFSWV9PUEVSQVRPUlMuZ2V0KGFzdC5vcGVyYXRpb24pO1xuICAgIGlmIChvcGVyYXRvciA9PT0gdW5kZWZpbmVkKSB7XG4gICAgICB0aHJvdyBuZXcgRXJyb3IoYEFzc2VydGlvbkVycm9yOiB1bmtub3duIGJpbmFyeSBvcGVyYXRvciAke2FzdC5vcGVyYXRpb259YCk7XG4gICAgfVxuICAgIHJldHVybiBuZXcgby5CaW5hcnlPcGVyYXRvckV4cHIoXG4gICAgICBvcGVyYXRvcixcbiAgICAgIGNvbnZlcnRBc3QoYXN0LmxlZnQsIGpvYiwgYmFzZVNvdXJjZVNwYW4pLFxuICAgICAgY29udmVydEFzdChhc3QucmlnaHQsIGpvYiwgYmFzZVNvdXJjZVNwYW4pLFxuICAgICAgdW5kZWZpbmVkLFxuICAgICAgY29udmVydFNvdXJjZVNwYW4oYXN0LnNwYW4sIGJhc2VTb3VyY2VTcGFuKSxcbiAgICApO1xuICB9IGVsc2UgaWYgKGFzdCBpbnN0YW5jZW9mIGUuVGhpc1JlY2VpdmVyKSB7XG4gICAgLy8gVE9ETzogc2hvdWxkIGNvbnRleHQgZXhwcmVzc2lvbnMgaGF2ZSBzb3VyY2UgbWFwcz9cbiAgICByZXR1cm4gbmV3IGlyLkNvbnRleHRFeHByKGpvYi5yb290LnhyZWYpO1xuICB9IGVsc2UgaWYgKGFzdCBpbnN0YW5jZW9mIGUuS2V5ZWRSZWFkKSB7XG4gICAgcmV0dXJuIG5ldyBvLlJlYWRLZXlFeHByKFxuICAgICAgY29udmVydEFzdChhc3QucmVjZWl2ZXIsIGpvYiwgYmFzZVNvdXJjZVNwYW4pLFxuICAgICAgY29udmVydEFzdChhc3Qua2V5LCBqb2IsIGJhc2VTb3VyY2VTcGFuKSxcbiAgICAgIHVuZGVmaW5lZCxcbiAgICAgIGNvbnZlcnRTb3VyY2VTcGFuKGFzdC5zcGFuLCBiYXNlU291cmNlU3BhbiksXG4gICAgKTtcbiAgfSBlbHNlIGlmIChhc3QgaW5zdGFuY2VvZiBlLkNoYWluKSB7XG4gICAgdGhyb3cgbmV3IEVycm9yKGBBc3NlcnRpb25FcnJvcjogQ2hhaW4gaW4gdW5rbm93biBjb250ZXh0YCk7XG4gIH0gZWxzZSBpZiAoYXN0IGluc3RhbmNlb2YgZS5MaXRlcmFsTWFwKSB7XG4gICAgY29uc3QgZW50cmllcyA9IGFzdC5rZXlzLm1hcCgoa2V5LCBpZHgpID0+IHtcbiAgICAgIGNvbnN0IHZhbHVlID0gYXN0LnZhbHVlc1tpZHhdO1xuICAgICAgLy8gVE9ETzogc2hvdWxkIGxpdGVyYWxzIGhhdmUgc291cmNlIG1hcHMsIG9yIGRvIHdlIGp1c3QgbWFwIHRoZSB3aG9sZSBzdXJyb3VuZGluZ1xuICAgICAgLy8gZXhwcmVzc2lvbj9cbiAgICAgIHJldHVybiBuZXcgby5MaXRlcmFsTWFwRW50cnkoa2V5LmtleSwgY29udmVydEFzdCh2YWx1ZSwgam9iLCBiYXNlU291cmNlU3BhbiksIGtleS5xdW90ZWQpO1xuICAgIH0pO1xuICAgIHJldHVybiBuZXcgby5MaXRlcmFsTWFwRXhwcihlbnRyaWVzLCB1bmRlZmluZWQsIGNvbnZlcnRTb3VyY2VTcGFuKGFzdC5zcGFuLCBiYXNlU291cmNlU3BhbikpO1xuICB9IGVsc2UgaWYgKGFzdCBpbnN0YW5jZW9mIGUuTGl0ZXJhbEFycmF5KSB7XG4gICAgLy8gVE9ETzogc2hvdWxkIGxpdGVyYWxzIGhhdmUgc291cmNlIG1hcHMsIG9yIGRvIHdlIGp1c3QgbWFwIHRoZSB3aG9sZSBzdXJyb3VuZGluZyBleHByZXNzaW9uP1xuICAgIHJldHVybiBuZXcgby5MaXRlcmFsQXJyYXlFeHByKFxuICAgICAgYXN0LmV4cHJlc3Npb25zLm1hcCgoZXhwcikgPT4gY29udmVydEFzdChleHByLCBqb2IsIGJhc2VTb3VyY2VTcGFuKSksXG4gICAgKTtcbiAgfSBlbHNlIGlmIChhc3QgaW5zdGFuY2VvZiBlLkNvbmRpdGlvbmFsKSB7XG4gICAgcmV0dXJuIG5ldyBvLkNvbmRpdGlvbmFsRXhwcihcbiAgICAgIGNvbnZlcnRBc3QoYXN0LmNvbmRpdGlvbiwgam9iLCBiYXNlU291cmNlU3BhbiksXG4gICAgICBjb252ZXJ0QXN0KGFzdC50cnVlRXhwLCBqb2IsIGJhc2VTb3VyY2VTcGFuKSxcbiAgICAgIGNvbnZlcnRBc3QoYXN0LmZhbHNlRXhwLCBqb2IsIGJhc2VTb3VyY2VTcGFuKSxcbiAgICAgIHVuZGVmaW5lZCxcbiAgICAgIGNvbnZlcnRTb3VyY2VTcGFuKGFzdC5zcGFuLCBiYXNlU291cmNlU3BhbiksXG4gICAgKTtcbiAgfSBlbHNlIGlmIChhc3QgaW5zdGFuY2VvZiBlLk5vbk51bGxBc3NlcnQpIHtcbiAgICAvLyBBIG5vbi1udWxsIGFzc2VydGlvbiBzaG91bGRuJ3QgaW1wYWN0IGdlbmVyYXRlZCBpbnN0cnVjdGlvbnMsIHNvIHdlIGNhbiBqdXN0IGRyb3AgaXQuXG4gICAgcmV0dXJuIGNvbnZlcnRBc3QoYXN0LmV4cHJlc3Npb24sIGpvYiwgYmFzZVNvdXJjZVNwYW4pO1xuICB9IGVsc2UgaWYgKGFzdCBpbnN0YW5jZW9mIGUuQmluZGluZ1BpcGUpIHtcbiAgICAvLyBUT0RPOiBwaXBlcyBzaG91bGQgcHJvYmFibHkgaGF2ZSBzb3VyY2UgbWFwczsgZmlndXJlIG91dCBkZXRhaWxzLlxuICAgIHJldHVybiBuZXcgaXIuUGlwZUJpbmRpbmdFeHByKGpvYi5hbGxvY2F0ZVhyZWZJZCgpLCBuZXcgaXIuU2xvdEhhbmRsZSgpLCBhc3QubmFtZSwgW1xuICAgICAgY29udmVydEFzdChhc3QuZXhwLCBqb2IsIGJhc2VTb3VyY2VTcGFuKSxcbiAgICAgIC4uLmFzdC5hcmdzLm1hcCgoYXJnKSA9PiBjb252ZXJ0QXN0KGFyZywgam9iLCBiYXNlU291cmNlU3BhbikpLFxuICAgIF0pO1xuICB9IGVsc2UgaWYgKGFzdCBpbnN0YW5jZW9mIGUuU2FmZUtleWVkUmVhZCkge1xuICAgIHJldHVybiBuZXcgaXIuU2FmZUtleWVkUmVhZEV4cHIoXG4gICAgICBjb252ZXJ0QXN0KGFzdC5yZWNlaXZlciwgam9iLCBiYXNlU291cmNlU3BhbiksXG4gICAgICBjb252ZXJ0QXN0KGFzdC5rZXksIGpvYiwgYmFzZVNvdXJjZVNwYW4pLFxuICAgICAgY29udmVydFNvdXJjZVNwYW4oYXN0LnNwYW4sIGJhc2VTb3VyY2VTcGFuKSxcbiAgICApO1xuICB9IGVsc2UgaWYgKGFzdCBpbnN0YW5jZW9mIGUuU2FmZVByb3BlcnR5UmVhZCkge1xuICAgIC8vIFRPRE86IHNvdXJjZSBzcGFuXG4gICAgcmV0dXJuIG5ldyBpci5TYWZlUHJvcGVydHlSZWFkRXhwcihjb252ZXJ0QXN0KGFzdC5yZWNlaXZlciwgam9iLCBiYXNlU291cmNlU3BhbiksIGFzdC5uYW1lKTtcbiAgfSBlbHNlIGlmIChhc3QgaW5zdGFuY2VvZiBlLlNhZmVDYWxsKSB7XG4gICAgLy8gVE9ETzogc291cmNlIHNwYW5cbiAgICByZXR1cm4gbmV3IGlyLlNhZmVJbnZva2VGdW5jdGlvbkV4cHIoXG4gICAgICBjb252ZXJ0QXN0KGFzdC5yZWNlaXZlciwgam9iLCBiYXNlU291cmNlU3BhbiksXG4gICAgICBhc3QuYXJncy5tYXAoKGEpID0+IGNvbnZlcnRBc3QoYSwgam9iLCBiYXNlU291cmNlU3BhbikpLFxuICAgICk7XG4gIH0gZWxzZSBpZiAoYXN0IGluc3RhbmNlb2YgZS5FbXB0eUV4cHIpIHtcbiAgICByZXR1cm4gbmV3IGlyLkVtcHR5RXhwcihjb252ZXJ0U291cmNlU3Bhbihhc3Quc3BhbiwgYmFzZVNvdXJjZVNwYW4pKTtcbiAgfSBlbHNlIGlmIChhc3QgaW5zdGFuY2VvZiBlLlByZWZpeE5vdCkge1xuICAgIHJldHVybiBvLm5vdChcbiAgICAgIGNvbnZlcnRBc3QoYXN0LmV4cHJlc3Npb24sIGpvYiwgYmFzZVNvdXJjZVNwYW4pLFxuICAgICAgY29udmVydFNvdXJjZVNwYW4oYXN0LnNwYW4sIGJhc2VTb3VyY2VTcGFuKSxcbiAgICApO1xuICB9IGVsc2Uge1xuICAgIHRocm93IG5ldyBFcnJvcihcbiAgICAgIGBVbmhhbmRsZWQgZXhwcmVzc2lvbiB0eXBlIFwiJHthc3QuY29uc3RydWN0b3IubmFtZX1cIiBpbiBmaWxlIFwiJHtiYXNlU291cmNlU3Bhbj8uc3RhcnQuZmlsZS51cmx9XCJgLFxuICAgICk7XG4gIH1cbn1cblxuZnVuY3Rpb24gY29udmVydEFzdFdpdGhJbnRlcnBvbGF0aW9uKFxuICBqb2I6IENvbXBpbGF0aW9uSm9iLFxuICB2YWx1ZTogZS5BU1QgfCBzdHJpbmcsXG4gIGkxOG5NZXRhOiBpMThuLkkxOG5NZXRhIHwgbnVsbCB8IHVuZGVmaW5lZCxcbiAgc291cmNlU3Bhbj86IFBhcnNlU291cmNlU3Bhbixcbik6IG8uRXhwcmVzc2lvbiB8IGlyLkludGVycG9sYXRpb24ge1xuICBsZXQgZXhwcmVzc2lvbjogby5FeHByZXNzaW9uIHwgaXIuSW50ZXJwb2xhdGlvbjtcbiAgaWYgKHZhbHVlIGluc3RhbmNlb2YgZS5JbnRlcnBvbGF0aW9uKSB7XG4gICAgZXhwcmVzc2lvbiA9IG5ldyBpci5JbnRlcnBvbGF0aW9uKFxuICAgICAgdmFsdWUuc3RyaW5ncyxcbiAgICAgIHZhbHVlLmV4cHJlc3Npb25zLm1hcCgoZSkgPT4gY29udmVydEFzdChlLCBqb2IsIHNvdXJjZVNwYW4gPz8gbnVsbCkpLFxuICAgICAgT2JqZWN0LmtleXMoYXNNZXNzYWdlKGkxOG5NZXRhKT8ucGxhY2Vob2xkZXJzID8/IHt9KSxcbiAgICApO1xuICB9IGVsc2UgaWYgKHZhbHVlIGluc3RhbmNlb2YgZS5BU1QpIHtcbiAgICBleHByZXNzaW9uID0gY29udmVydEFzdCh2YWx1ZSwgam9iLCBzb3VyY2VTcGFuID8/IG51bGwpO1xuICB9IGVsc2Uge1xuICAgIGV4cHJlc3Npb24gPSBvLmxpdGVyYWwodmFsdWUpO1xuICB9XG4gIHJldHVybiBleHByZXNzaW9uO1xufVxuXG4vLyBUT0RPOiBDYW4gd2UgcG9wdWxhdGUgVGVtcGxhdGUgYmluZGluZyBraW5kcyBpbiBpbmdlc3Q/XG5jb25zdCBCSU5ESU5HX0tJTkRTID0gbmV3IE1hcDxlLkJpbmRpbmdUeXBlLCBpci5CaW5kaW5nS2luZD4oW1xuICBbZS5CaW5kaW5nVHlwZS5Qcm9wZXJ0eSwgaXIuQmluZGluZ0tpbmQuUHJvcGVydHldLFxuICBbZS5CaW5kaW5nVHlwZS5Ud29XYXksIGlyLkJpbmRpbmdLaW5kLlR3b1dheVByb3BlcnR5XSxcbiAgW2UuQmluZGluZ1R5cGUuQXR0cmlidXRlLCBpci5CaW5kaW5nS2luZC5BdHRyaWJ1dGVdLFxuICBbZS5CaW5kaW5nVHlwZS5DbGFzcywgaXIuQmluZGluZ0tpbmQuQ2xhc3NOYW1lXSxcbiAgW2UuQmluZGluZ1R5cGUuU3R5bGUsIGlyLkJpbmRpbmdLaW5kLlN0eWxlUHJvcGVydHldLFxuICBbZS5CaW5kaW5nVHlwZS5BbmltYXRpb24sIGlyLkJpbmRpbmdLaW5kLkFuaW1hdGlvbl0sXG5dKTtcblxuLyoqXG4gKiBDaGVja3Mgd2hldGhlciB0aGUgZ2l2ZW4gdGVtcGxhdGUgaXMgYSBwbGFpbiBuZy10ZW1wbGF0ZSAoYXMgb3Bwb3NlZCB0byBhbm90aGVyIGtpbmQgb2YgdGVtcGxhdGVcbiAqIHN1Y2ggYXMgYSBzdHJ1Y3R1cmFsIGRpcmVjdGl2ZSB0ZW1wbGF0ZSBvciBjb250cm9sIGZsb3cgdGVtcGxhdGUpLiBUaGlzIGlzIGNoZWNrZWQgYmFzZWQgb24gdGhlXG4gKiB0YWdOYW1lLiBXZSBjYW4gZXhwZWN0IHRoYXQgb25seSBwbGFpbiBuZy10ZW1wbGF0ZXMgd2lsbCBjb21lIHRocm91Z2ggd2l0aCBhIHRhZ05hbWUgb2ZcbiAqICduZy10ZW1wbGF0ZScuXG4gKlxuICogSGVyZSBhcmUgc29tZSBvZiB0aGUgY2FzZXMgd2UgZXhwZWN0OlxuICpcbiAqIHwgQW5ndWxhciBIVE1MICAgICAgICAgICAgICAgICAgICAgICB8IFRlbXBsYXRlIHRhZ05hbWUgICB8XG4gKiB8IC0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0gfCAtLS0tLS0tLS0tLS0tLS0tLS0gfFxuICogfCBgPG5nLXRlbXBsYXRlPmAgICAgICAgICAgICAgICAgICAgIHwgJ25nLXRlbXBsYXRlJyAgICAgIHxcbiAqIHwgYDxkaXYgKm5nSWY9XCJ0cnVlXCI+YCAgICAgICAgICAgICAgIHwgJ2RpdicgICAgICAgICAgICAgIHxcbiAqIHwgYDxzdmc+PG5nLXRlbXBsYXRlPmAgICAgICAgICAgICAgICB8ICdzdmc6bmctdGVtcGxhdGUnICB8XG4gKiB8IGBAaWYgKHRydWUpIHtgICAgICAgICAgICAgICAgICAgICAgfCAnQ29uZGl0aW9uYWwnICAgICAgfFxuICogfCBgPG5nLXRlbXBsYXRlICpuZ0lmPmAgKHBsYWluKSAgICAgIHwgJ25nLXRlbXBsYXRlJyAgICAgIHxcbiAqIHwgYDxuZy10ZW1wbGF0ZSAqbmdJZj5gIChzdHJ1Y3R1cmFsKSB8IG51bGwgICAgICAgICAgICAgICB8XG4gKi9cbmZ1bmN0aW9uIGlzUGxhaW5UZW1wbGF0ZSh0bXBsOiB0LlRlbXBsYXRlKSB7XG4gIHJldHVybiBzcGxpdE5zTmFtZSh0bXBsLnRhZ05hbWUgPz8gJycpWzFdID09PSBOR19URU1QTEFURV9UQUdfTkFNRTtcbn1cblxuLyoqXG4gKiBFbnN1cmVzIHRoYXQgdGhlIGkxOG5NZXRhLCBpZiBwcm92aWRlZCwgaXMgYW4gaTE4bi5NZXNzYWdlLlxuICovXG5mdW5jdGlvbiBhc01lc3NhZ2UoaTE4bk1ldGE6IGkxOG4uSTE4bk1ldGEgfCBudWxsIHwgdW5kZWZpbmVkKTogaTE4bi5NZXNzYWdlIHwgbnVsbCB7XG4gIGlmIChpMThuTWV0YSA9PSBudWxsKSB7XG4gICAgcmV0dXJuIG51bGw7XG4gIH1cbiAgaWYgKCEoaTE4bk1ldGEgaW5zdGFuY2VvZiBpMThuLk1lc3NhZ2UpKSB7XG4gICAgdGhyb3cgRXJyb3IoYEV4cGVjdGVkIGkxOG4gbWV0YSB0byBiZSBhIE1lc3NhZ2UsIGJ1dCBnb3Q6ICR7aTE4bk1ldGEuY29uc3RydWN0b3IubmFtZX1gKTtcbiAgfVxuICByZXR1cm4gaTE4bk1ldGE7XG59XG5cbi8qKlxuICogUHJvY2VzcyBhbGwgb2YgdGhlIGJpbmRpbmdzIG9uIGFuIGVsZW1lbnQgaW4gdGhlIHRlbXBsYXRlIEFTVCBhbmQgY29udmVydCB0aGVtIHRvIHRoZWlyIElSXG4gKiByZXByZXNlbnRhdGlvbi5cbiAqL1xuZnVuY3Rpb24gaW5nZXN0RWxlbWVudEJpbmRpbmdzKFxuICB1bml0OiBWaWV3Q29tcGlsYXRpb25Vbml0LFxuICBvcDogaXIuRWxlbWVudE9wQmFzZSxcbiAgZWxlbWVudDogdC5FbGVtZW50LFxuKTogdm9pZCB7XG4gIGxldCBiaW5kaW5ncyA9IG5ldyBBcnJheTxpci5CaW5kaW5nT3AgfCBpci5FeHRyYWN0ZWRBdHRyaWJ1dGVPcCB8IG51bGw+KCk7XG5cbiAgbGV0IGkxOG5BdHRyaWJ1dGVCaW5kaW5nTmFtZXMgPSBuZXcgU2V0PHN0cmluZz4oKTtcblxuICBmb3IgKGNvbnN0IGF0dHIgb2YgZWxlbWVudC5hdHRyaWJ1dGVzKSB7XG4gICAgLy8gQXR0cmlidXRlIGxpdGVyYWwgYmluZGluZ3MsIHN1Y2ggYXMgYGF0dHIuZm9vPVwiYmFyXCJgLlxuICAgIGNvbnN0IHNlY3VyaXR5Q29udGV4dCA9IGRvbVNjaGVtYS5zZWN1cml0eUNvbnRleHQoZWxlbWVudC5uYW1lLCBhdHRyLm5hbWUsIHRydWUpO1xuICAgIGJpbmRpbmdzLnB1c2goXG4gICAgICBpci5jcmVhdGVCaW5kaW5nT3AoXG4gICAgICAgIG9wLnhyZWYsXG4gICAgICAgIGlyLkJpbmRpbmdLaW5kLkF0dHJpYnV0ZSxcbiAgICAgICAgYXR0ci5uYW1lLFxuICAgICAgICBjb252ZXJ0QXN0V2l0aEludGVycG9sYXRpb24odW5pdC5qb2IsIGF0dHIudmFsdWUsIGF0dHIuaTE4biksXG4gICAgICAgIG51bGwsXG4gICAgICAgIHNlY3VyaXR5Q29udGV4dCxcbiAgICAgICAgdHJ1ZSxcbiAgICAgICAgZmFsc2UsXG4gICAgICAgIG51bGwsXG4gICAgICAgIGFzTWVzc2FnZShhdHRyLmkxOG4pLFxuICAgICAgICBhdHRyLnNvdXJjZVNwYW4sXG4gICAgICApLFxuICAgICk7XG4gICAgaWYgKGF0dHIuaTE4bikge1xuICAgICAgaTE4bkF0dHJpYnV0ZUJpbmRpbmdOYW1lcy5hZGQoYXR0ci5uYW1lKTtcbiAgICB9XG4gIH1cblxuICBmb3IgKGNvbnN0IGlucHV0IG9mIGVsZW1lbnQuaW5wdXRzKSB7XG4gICAgaWYgKGkxOG5BdHRyaWJ1dGVCaW5kaW5nTmFtZXMuaGFzKGlucHV0Lm5hbWUpKSB7XG4gICAgICBjb25zb2xlLmVycm9yKFxuICAgICAgICBgT24gY29tcG9uZW50ICR7dW5pdC5qb2IuY29tcG9uZW50TmFtZX0sIHRoZSBiaW5kaW5nICR7aW5wdXQubmFtZX0gaXMgYm90aCBhbiBpMThuIGF0dHJpYnV0ZSBhbmQgYSBwcm9wZXJ0eS4gWW91IG1heSB3YW50IHRvIHJlbW92ZSB0aGUgcHJvcGVydHkgYmluZGluZy4gVGhpcyB3aWxsIGJlY29tZSBhIGNvbXBpbGF0aW9uIGVycm9yIGluIGZ1dHVyZSB2ZXJzaW9ucyBvZiBBbmd1bGFyLmAsXG4gICAgICApO1xuICAgIH1cbiAgICAvLyBBbGwgZHluYW1pYyBiaW5kaW5ncyAoYm90aCBhdHRyaWJ1dGUgYW5kIHByb3BlcnR5IGJpbmRpbmdzKS5cbiAgICBiaW5kaW5ncy5wdXNoKFxuICAgICAgaXIuY3JlYXRlQmluZGluZ09wKFxuICAgICAgICBvcC54cmVmLFxuICAgICAgICBCSU5ESU5HX0tJTkRTLmdldChpbnB1dC50eXBlKSEsXG4gICAgICAgIGlucHV0Lm5hbWUsXG4gICAgICAgIGNvbnZlcnRBc3RXaXRoSW50ZXJwb2xhdGlvbih1bml0LmpvYiwgYXN0T2YoaW5wdXQudmFsdWUpLCBpbnB1dC5pMThuKSxcbiAgICAgICAgaW5wdXQudW5pdCxcbiAgICAgICAgaW5wdXQuc2VjdXJpdHlDb250ZXh0LFxuICAgICAgICBmYWxzZSxcbiAgICAgICAgZmFsc2UsXG4gICAgICAgIG51bGwsXG4gICAgICAgIGFzTWVzc2FnZShpbnB1dC5pMThuKSA/PyBudWxsLFxuICAgICAgICBpbnB1dC5zb3VyY2VTcGFuLFxuICAgICAgKSxcbiAgICApO1xuICB9XG5cbiAgdW5pdC5jcmVhdGUucHVzaChcbiAgICBiaW5kaW5ncy5maWx0ZXIoKGIpOiBiIGlzIGlyLkV4dHJhY3RlZEF0dHJpYnV0ZU9wID0+IGI/LmtpbmQgPT09IGlyLk9wS2luZC5FeHRyYWN0ZWRBdHRyaWJ1dGUpLFxuICApO1xuICB1bml0LnVwZGF0ZS5wdXNoKGJpbmRpbmdzLmZpbHRlcigoYik6IGIgaXMgaXIuQmluZGluZ09wID0+IGI/LmtpbmQgPT09IGlyLk9wS2luZC5CaW5kaW5nKSk7XG5cbiAgZm9yIChjb25zdCBvdXRwdXQgb2YgZWxlbWVudC5vdXRwdXRzKSB7XG4gICAgaWYgKG91dHB1dC50eXBlID09PSBlLlBhcnNlZEV2ZW50VHlwZS5BbmltYXRpb24gJiYgb3V0cHV0LnBoYXNlID09PSBudWxsKSB7XG4gICAgICB0aHJvdyBFcnJvcignQW5pbWF0aW9uIGxpc3RlbmVyIHNob3VsZCBoYXZlIGEgcGhhc2UnKTtcbiAgICB9XG5cbiAgICBpZiAob3V0cHV0LnR5cGUgPT09IGUuUGFyc2VkRXZlbnRUeXBlLlR3b1dheSkge1xuICAgICAgdW5pdC5jcmVhdGUucHVzaChcbiAgICAgICAgaXIuY3JlYXRlVHdvV2F5TGlzdGVuZXJPcChcbiAgICAgICAgICBvcC54cmVmLFxuICAgICAgICAgIG9wLmhhbmRsZSxcbiAgICAgICAgICBvdXRwdXQubmFtZSxcbiAgICAgICAgICBvcC50YWcsXG4gICAgICAgICAgbWFrZVR3b1dheUxpc3RlbmVySGFuZGxlck9wcyh1bml0LCBvdXRwdXQuaGFuZGxlciwgb3V0cHV0LmhhbmRsZXJTcGFuKSxcbiAgICAgICAgICBvdXRwdXQuc291cmNlU3BhbixcbiAgICAgICAgKSxcbiAgICAgICk7XG4gICAgfSBlbHNlIHtcbiAgICAgIHVuaXQuY3JlYXRlLnB1c2goXG4gICAgICAgIGlyLmNyZWF0ZUxpc3RlbmVyT3AoXG4gICAgICAgICAgb3AueHJlZixcbiAgICAgICAgICBvcC5oYW5kbGUsXG4gICAgICAgICAgb3V0cHV0Lm5hbWUsXG4gICAgICAgICAgb3AudGFnLFxuICAgICAgICAgIG1ha2VMaXN0ZW5lckhhbmRsZXJPcHModW5pdCwgb3V0cHV0LmhhbmRsZXIsIG91dHB1dC5oYW5kbGVyU3BhbiksXG4gICAgICAgICAgb3V0cHV0LnBoYXNlLFxuICAgICAgICAgIG91dHB1dC50YXJnZXQsXG4gICAgICAgICAgZmFsc2UsXG4gICAgICAgICAgb3V0cHV0LnNvdXJjZVNwYW4sXG4gICAgICAgICksXG4gICAgICApO1xuICAgIH1cbiAgfVxuXG4gIC8vIElmIGFueSBvZiB0aGUgYmluZGluZ3Mgb24gdGhpcyBlbGVtZW50IGhhdmUgYW4gaTE4biBtZXNzYWdlLCB0aGVuIGFuIGkxOG4gYXR0cnMgY29uZmlndXJhdGlvblxuICAvLyBvcCBpcyBhbHNvIHJlcXVpcmVkLlxuICBpZiAoYmluZGluZ3Muc29tZSgoYikgPT4gYj8uaTE4bk1lc3NhZ2UpICE9PSBudWxsKSB7XG4gICAgdW5pdC5jcmVhdGUucHVzaChcbiAgICAgIGlyLmNyZWF0ZUkxOG5BdHRyaWJ1dGVzT3AodW5pdC5qb2IuYWxsb2NhdGVYcmVmSWQoKSwgbmV3IGlyLlNsb3RIYW5kbGUoKSwgb3AueHJlZiksXG4gICAgKTtcbiAgfVxufVxuXG4vKipcbiAqIFByb2Nlc3MgYWxsIG9mIHRoZSBiaW5kaW5ncyBvbiBhIHRlbXBsYXRlIGluIHRoZSB0ZW1wbGF0ZSBBU1QgYW5kIGNvbnZlcnQgdGhlbSB0byB0aGVpciBJUlxuICogcmVwcmVzZW50YXRpb24uXG4gKi9cbmZ1bmN0aW9uIGluZ2VzdFRlbXBsYXRlQmluZGluZ3MoXG4gIHVuaXQ6IFZpZXdDb21waWxhdGlvblVuaXQsXG4gIG9wOiBpci5FbGVtZW50T3BCYXNlLFxuICB0ZW1wbGF0ZTogdC5UZW1wbGF0ZSxcbiAgdGVtcGxhdGVLaW5kOiBpci5UZW1wbGF0ZUtpbmQgfCBudWxsLFxuKTogdm9pZCB7XG4gIGxldCBiaW5kaW5ncyA9IG5ldyBBcnJheTxpci5CaW5kaW5nT3AgfCBpci5FeHRyYWN0ZWRBdHRyaWJ1dGVPcCB8IG51bGw+KCk7XG5cbiAgZm9yIChjb25zdCBhdHRyIG9mIHRlbXBsYXRlLnRlbXBsYXRlQXR0cnMpIHtcbiAgICBpZiAoYXR0ciBpbnN0YW5jZW9mIHQuVGV4dEF0dHJpYnV0ZSkge1xuICAgICAgY29uc3Qgc2VjdXJpdHlDb250ZXh0ID0gZG9tU2NoZW1hLnNlY3VyaXR5Q29udGV4dChOR19URU1QTEFURV9UQUdfTkFNRSwgYXR0ci5uYW1lLCB0cnVlKTtcbiAgICAgIGJpbmRpbmdzLnB1c2goXG4gICAgICAgIGNyZWF0ZVRlbXBsYXRlQmluZGluZyhcbiAgICAgICAgICB1bml0LFxuICAgICAgICAgIG9wLnhyZWYsXG4gICAgICAgICAgZS5CaW5kaW5nVHlwZS5BdHRyaWJ1dGUsXG4gICAgICAgICAgYXR0ci5uYW1lLFxuICAgICAgICAgIGF0dHIudmFsdWUsXG4gICAgICAgICAgbnVsbCxcbiAgICAgICAgICBzZWN1cml0eUNvbnRleHQsXG4gICAgICAgICAgdHJ1ZSxcbiAgICAgICAgICB0ZW1wbGF0ZUtpbmQsXG4gICAgICAgICAgYXNNZXNzYWdlKGF0dHIuaTE4biksXG4gICAgICAgICAgYXR0ci5zb3VyY2VTcGFuLFxuICAgICAgICApLFxuICAgICAgKTtcbiAgICB9IGVsc2Uge1xuICAgICAgYmluZGluZ3MucHVzaChcbiAgICAgICAgY3JlYXRlVGVtcGxhdGVCaW5kaW5nKFxuICAgICAgICAgIHVuaXQsXG4gICAgICAgICAgb3AueHJlZixcbiAgICAgICAgICBhdHRyLnR5cGUsXG4gICAgICAgICAgYXR0ci5uYW1lLFxuICAgICAgICAgIGFzdE9mKGF0dHIudmFsdWUpLFxuICAgICAgICAgIGF0dHIudW5pdCxcbiAgICAgICAgICBhdHRyLnNlY3VyaXR5Q29udGV4dCxcbiAgICAgICAgICB0cnVlLFxuICAgICAgICAgIHRlbXBsYXRlS2luZCxcbiAgICAgICAgICBhc01lc3NhZ2UoYXR0ci5pMThuKSxcbiAgICAgICAgICBhdHRyLnNvdXJjZVNwYW4sXG4gICAgICAgICksXG4gICAgICApO1xuICAgIH1cbiAgfVxuXG4gIGZvciAoY29uc3QgYXR0ciBvZiB0ZW1wbGF0ZS5hdHRyaWJ1dGVzKSB7XG4gICAgLy8gQXR0cmlidXRlIGxpdGVyYWwgYmluZGluZ3MsIHN1Y2ggYXMgYGF0dHIuZm9vPVwiYmFyXCJgLlxuICAgIGNvbnN0IHNlY3VyaXR5Q29udGV4dCA9IGRvbVNjaGVtYS5zZWN1cml0eUNvbnRleHQoTkdfVEVNUExBVEVfVEFHX05BTUUsIGF0dHIubmFtZSwgdHJ1ZSk7XG4gICAgYmluZGluZ3MucHVzaChcbiAgICAgIGNyZWF0ZVRlbXBsYXRlQmluZGluZyhcbiAgICAgICAgdW5pdCxcbiAgICAgICAgb3AueHJlZixcbiAgICAgICAgZS5CaW5kaW5nVHlwZS5BdHRyaWJ1dGUsXG4gICAgICAgIGF0dHIubmFtZSxcbiAgICAgICAgYXR0ci52YWx1ZSxcbiAgICAgICAgbnVsbCxcbiAgICAgICAgc2VjdXJpdHlDb250ZXh0LFxuICAgICAgICBmYWxzZSxcbiAgICAgICAgdGVtcGxhdGVLaW5kLFxuICAgICAgICBhc01lc3NhZ2UoYXR0ci5pMThuKSxcbiAgICAgICAgYXR0ci5zb3VyY2VTcGFuLFxuICAgICAgKSxcbiAgICApO1xuICB9XG5cbiAgZm9yIChjb25zdCBpbnB1dCBvZiB0ZW1wbGF0ZS5pbnB1dHMpIHtcbiAgICAvLyBEeW5hbWljIGJpbmRpbmdzIChib3RoIGF0dHJpYnV0ZSBhbmQgcHJvcGVydHkgYmluZGluZ3MpLlxuICAgIGJpbmRpbmdzLnB1c2goXG4gICAgICBjcmVhdGVUZW1wbGF0ZUJpbmRpbmcoXG4gICAgICAgIHVuaXQsXG4gICAgICAgIG9wLnhyZWYsXG4gICAgICAgIGlucHV0LnR5cGUsXG4gICAgICAgIGlucHV0Lm5hbWUsXG4gICAgICAgIGFzdE9mKGlucHV0LnZhbHVlKSxcbiAgICAgICAgaW5wdXQudW5pdCxcbiAgICAgICAgaW5wdXQuc2VjdXJpdHlDb250ZXh0LFxuICAgICAgICBmYWxzZSxcbiAgICAgICAgdGVtcGxhdGVLaW5kLFxuICAgICAgICBhc01lc3NhZ2UoaW5wdXQuaTE4biksXG4gICAgICAgIGlucHV0LnNvdXJjZVNwYW4sXG4gICAgICApLFxuICAgICk7XG4gIH1cblxuICB1bml0LmNyZWF0ZS5wdXNoKFxuICAgIGJpbmRpbmdzLmZpbHRlcigoYik6IGIgaXMgaXIuRXh0cmFjdGVkQXR0cmlidXRlT3AgPT4gYj8ua2luZCA9PT0gaXIuT3BLaW5kLkV4dHJhY3RlZEF0dHJpYnV0ZSksXG4gICk7XG4gIHVuaXQudXBkYXRlLnB1c2goYmluZGluZ3MuZmlsdGVyKChiKTogYiBpcyBpci5CaW5kaW5nT3AgPT4gYj8ua2luZCA9PT0gaXIuT3BLaW5kLkJpbmRpbmcpKTtcblxuICBmb3IgKGNvbnN0IG91dHB1dCBvZiB0ZW1wbGF0ZS5vdXRwdXRzKSB7XG4gICAgaWYgKG91dHB1dC50eXBlID09PSBlLlBhcnNlZEV2ZW50VHlwZS5BbmltYXRpb24gJiYgb3V0cHV0LnBoYXNlID09PSBudWxsKSB7XG4gICAgICB0aHJvdyBFcnJvcignQW5pbWF0aW9uIGxpc3RlbmVyIHNob3VsZCBoYXZlIGEgcGhhc2UnKTtcbiAgICB9XG5cbiAgICBpZiAodGVtcGxhdGVLaW5kID09PSBpci5UZW1wbGF0ZUtpbmQuTmdUZW1wbGF0ZSkge1xuICAgICAgaWYgKG91dHB1dC50eXBlID09PSBlLlBhcnNlZEV2ZW50VHlwZS5Ud29XYXkpIHtcbiAgICAgICAgdW5pdC5jcmVhdGUucHVzaChcbiAgICAgICAgICBpci5jcmVhdGVUd29XYXlMaXN0ZW5lck9wKFxuICAgICAgICAgICAgb3AueHJlZixcbiAgICAgICAgICAgIG9wLmhhbmRsZSxcbiAgICAgICAgICAgIG91dHB1dC5uYW1lLFxuICAgICAgICAgICAgb3AudGFnLFxuICAgICAgICAgICAgbWFrZVR3b1dheUxpc3RlbmVySGFuZGxlck9wcyh1bml0LCBvdXRwdXQuaGFuZGxlciwgb3V0cHV0LmhhbmRsZXJTcGFuKSxcbiAgICAgICAgICAgIG91dHB1dC5zb3VyY2VTcGFuLFxuICAgICAgICAgICksXG4gICAgICAgICk7XG4gICAgICB9IGVsc2Uge1xuICAgICAgICB1bml0LmNyZWF0ZS5wdXNoKFxuICAgICAgICAgIGlyLmNyZWF0ZUxpc3RlbmVyT3AoXG4gICAgICAgICAgICBvcC54cmVmLFxuICAgICAgICAgICAgb3AuaGFuZGxlLFxuICAgICAgICAgICAgb3V0cHV0Lm5hbWUsXG4gICAgICAgICAgICBvcC50YWcsXG4gICAgICAgICAgICBtYWtlTGlzdGVuZXJIYW5kbGVyT3BzKHVuaXQsIG91dHB1dC5oYW5kbGVyLCBvdXRwdXQuaGFuZGxlclNwYW4pLFxuICAgICAgICAgICAgb3V0cHV0LnBoYXNlLFxuICAgICAgICAgICAgb3V0cHV0LnRhcmdldCxcbiAgICAgICAgICAgIGZhbHNlLFxuICAgICAgICAgICAgb3V0cHV0LnNvdXJjZVNwYW4sXG4gICAgICAgICAgKSxcbiAgICAgICAgKTtcbiAgICAgIH1cbiAgICB9XG4gICAgaWYgKFxuICAgICAgdGVtcGxhdGVLaW5kID09PSBpci5UZW1wbGF0ZUtpbmQuU3RydWN0dXJhbCAmJlxuICAgICAgb3V0cHV0LnR5cGUgIT09IGUuUGFyc2VkRXZlbnRUeXBlLkFuaW1hdGlvblxuICAgICkge1xuICAgICAgLy8gQW5pbWF0aW9uIGJpbmRpbmdzIGFyZSBleGNsdWRlZCBmcm9tIHRoZSBzdHJ1Y3R1cmFsIHRlbXBsYXRlJ3MgY29uc3QgYXJyYXkuXG4gICAgICBjb25zdCBzZWN1cml0eUNvbnRleHQgPSBkb21TY2hlbWEuc2VjdXJpdHlDb250ZXh0KE5HX1RFTVBMQVRFX1RBR19OQU1FLCBvdXRwdXQubmFtZSwgZmFsc2UpO1xuICAgICAgdW5pdC5jcmVhdGUucHVzaChcbiAgICAgICAgaXIuY3JlYXRlRXh0cmFjdGVkQXR0cmlidXRlT3AoXG4gICAgICAgICAgb3AueHJlZixcbiAgICAgICAgICBpci5CaW5kaW5nS2luZC5Qcm9wZXJ0eSxcbiAgICAgICAgICBudWxsLFxuICAgICAgICAgIG91dHB1dC5uYW1lLFxuICAgICAgICAgIG51bGwsXG4gICAgICAgICAgbnVsbCxcbiAgICAgICAgICBudWxsLFxuICAgICAgICAgIHNlY3VyaXR5Q29udGV4dCxcbiAgICAgICAgKSxcbiAgICAgICk7XG4gICAgfVxuICB9XG5cbiAgLy8gVE9ETzogUGVyaGFwcyB3ZSBjb3VsZCBkbyB0aGlzIGluIGEgcGhhc2U/IChJdCBsaWtlbHkgd291bGRuJ3QgY2hhbmdlIHRoZSBzbG90IGluZGljZXMuKVxuICBpZiAoYmluZGluZ3Muc29tZSgoYikgPT4gYj8uaTE4bk1lc3NhZ2UpICE9PSBudWxsKSB7XG4gICAgdW5pdC5jcmVhdGUucHVzaChcbiAgICAgIGlyLmNyZWF0ZUkxOG5BdHRyaWJ1dGVzT3AodW5pdC5qb2IuYWxsb2NhdGVYcmVmSWQoKSwgbmV3IGlyLlNsb3RIYW5kbGUoKSwgb3AueHJlZiksXG4gICAgKTtcbiAgfVxufVxuXG4vKipcbiAqIEhlbHBlciB0byBpbmdlc3QgYW4gaW5kaXZpZHVhbCBiaW5kaW5nIG9uIGEgdGVtcGxhdGUsIGVpdGhlciBhbiBleHBsaWNpdCBgbmctdGVtcGxhdGVgLCBvciBhblxuICogaW1wbGljaXQgdGVtcGxhdGUgY3JlYXRlZCB2aWEgc3RydWN0dXJhbCBkaXJlY3RpdmUuXG4gKlxuICogQmluZGluZ3Mgb24gdGVtcGxhdGVzIGFyZSAqZXh0cmVtZWx5KiB0cmlja3kuIEkgaGF2ZSB0cmllZCB0byBpc29sYXRlIGFsbCBvZiB0aGUgY29uZnVzaW5nIGVkZ2VcbiAqIGNhc2VzIGludG8gdGhpcyBmdW5jdGlvbiwgYW5kIHRvIGNvbW1lbnQgaXQgd2VsbCB0byBkb2N1bWVudCB0aGUgYmVoYXZpb3IuXG4gKlxuICogU29tZSBvZiB0aGlzIGJlaGF2aW9yIGlzIGludHVpdGl2ZWx5IGluY29ycmVjdCwgYW5kIHdlIHNob3VsZCBjb25zaWRlciBjaGFuZ2luZyBpdCBpbiB0aGUgZnV0dXJlLlxuICpcbiAqIEBwYXJhbSB2aWV3IFRoZSBjb21waWxhdGlvbiB1bml0IGZvciB0aGUgdmlldyBjb250YWluaW5nIHRoZSB0ZW1wbGF0ZS5cbiAqIEBwYXJhbSB4cmVmIFRoZSB4cmVmIG9mIHRoZSB0ZW1wbGF0ZSBvcC5cbiAqIEBwYXJhbSB0eXBlIFRoZSBiaW5kaW5nIHR5cGUsIGFjY29yZGluZyB0byB0aGUgcGFyc2VyLiBUaGlzIGlzIGZhaXJseSByZWFzb25hYmxlLCBlLmcuIGJvdGhcbiAqICAgICBkeW5hbWljIGFuZCBzdGF0aWMgYXR0cmlidXRlcyBoYXZlIGUuQmluZGluZ1R5cGUuQXR0cmlidXRlLlxuICogQHBhcmFtIG5hbWUgVGhlIGJpbmRpbmcncyBuYW1lLlxuICogQHBhcmFtIHZhbHVlIFRoZSBiaW5kaW5ncydzIHZhbHVlLCB3aGljaCB3aWxsIGVpdGhlciBiZSBhbiBpbnB1dCBBU1QgZXhwcmVzc2lvbiwgb3IgYSBzdHJpbmdcbiAqICAgICBsaXRlcmFsLiBOb3RlIHRoYXQgdGhlIGlucHV0IEFTVCBleHByZXNzaW9uIG1heSBvciBtYXkgbm90IGJlIGNvbnN0IC0tIGl0IHdpbGwgb25seSBiZSBhXG4gKiAgICAgc3RyaW5nIGxpdGVyYWwgaWYgdGhlIHBhcnNlciBjb25zaWRlcmVkIGl0IGEgdGV4dCBiaW5kaW5nLlxuICogQHBhcmFtIHVuaXQgSWYgdGhlIGJpbmRpbmcgaGFzIGEgdW5pdCAoZS5nLiBgcHhgIGZvciBzdHlsZSBiaW5kaW5ncyksIHRoZW4gdGhpcyBpcyB0aGUgdW5pdC5cbiAqIEBwYXJhbSBzZWN1cml0eUNvbnRleHQgVGhlIHNlY3VyaXR5IGNvbnRleHQgb2YgdGhlIGJpbmRpbmcuXG4gKiBAcGFyYW0gaXNTdHJ1Y3R1cmFsVGVtcGxhdGVBdHRyaWJ1dGUgV2hldGhlciB0aGlzIGJpbmRpbmcgYWN0dWFsbHkgYXBwbGllcyB0byB0aGUgc3RydWN0dXJhbFxuICogICAgIG5nLXRlbXBsYXRlLiBGb3IgZXhhbXBsZSwgYW4gYG5nRm9yYCB3b3VsZCBhY3R1YWxseSBhcHBseSB0byB0aGUgc3RydWN0dXJhbCB0ZW1wbGF0ZS4gKE1vc3RcbiAqICAgICBiaW5kaW5ncyBvbiBzdHJ1Y3R1cmFsIGVsZW1lbnRzIHRhcmdldCB0aGUgaW5uZXIgZWxlbWVudCwgbm90IHRoZSB0ZW1wbGF0ZS4pXG4gKiBAcGFyYW0gdGVtcGxhdGVLaW5kIFdoZXRoZXIgdGhpcyBpcyBhbiBleHBsaWNpdCBgbmctdGVtcGxhdGVgIG9yIGFuIGltcGxpY2l0IHRlbXBsYXRlIGNyZWF0ZWQgYnlcbiAqICAgICBhIHN0cnVjdHVyYWwgZGlyZWN0aXZlLiBUaGlzIHNob3VsZCBuZXZlciBiZSBhIGJsb2NrIHRlbXBsYXRlLlxuICogQHBhcmFtIGkxOG5NZXNzYWdlIFRoZSBpMThuIG1ldGFkYXRhIGZvciB0aGUgYmluZGluZywgaWYgYW55LlxuICogQHBhcmFtIHNvdXJjZVNwYW4gVGhlIHNvdXJjZSBzcGFuIG9mIHRoZSBiaW5kaW5nLlxuICogQHJldHVybnMgQW4gSVIgYmluZGluZyBvcCwgb3IgbnVsbCBpZiB0aGUgYmluZGluZyBzaG91bGQgYmUgc2tpcHBlZC5cbiAqL1xuZnVuY3Rpb24gY3JlYXRlVGVtcGxhdGVCaW5kaW5nKFxuICB2aWV3OiBWaWV3Q29tcGlsYXRpb25Vbml0LFxuICB4cmVmOiBpci5YcmVmSWQsXG4gIHR5cGU6IGUuQmluZGluZ1R5cGUsXG4gIG5hbWU6IHN0cmluZyxcbiAgdmFsdWU6IGUuQVNUIHwgc3RyaW5nLFxuICB1bml0OiBzdHJpbmcgfCBudWxsLFxuICBzZWN1cml0eUNvbnRleHQ6IFNlY3VyaXR5Q29udGV4dCxcbiAgaXNTdHJ1Y3R1cmFsVGVtcGxhdGVBdHRyaWJ1dGU6IGJvb2xlYW4sXG4gIHRlbXBsYXRlS2luZDogaXIuVGVtcGxhdGVLaW5kIHwgbnVsbCxcbiAgaTE4bk1lc3NhZ2U6IGkxOG4uTWVzc2FnZSB8IG51bGwsXG4gIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3Bhbixcbik6IGlyLkJpbmRpbmdPcCB8IGlyLkV4dHJhY3RlZEF0dHJpYnV0ZU9wIHwgbnVsbCB7XG4gIGNvbnN0IGlzVGV4dEJpbmRpbmcgPSB0eXBlb2YgdmFsdWUgPT09ICdzdHJpbmcnO1xuICAvLyBJZiB0aGlzIGlzIGEgc3RydWN0dXJhbCB0ZW1wbGF0ZSwgdGhlbiBzZXZlcmFsIGtpbmRzIG9mIGJpbmRpbmdzIHNob3VsZCBub3QgcmVzdWx0IGluIGFuXG4gIC8vIHVwZGF0ZSBpbnN0cnVjdGlvbi5cbiAgaWYgKHRlbXBsYXRlS2luZCA9PT0gaXIuVGVtcGxhdGVLaW5kLlN0cnVjdHVyYWwpIHtcbiAgICBpZiAoIWlzU3RydWN0dXJhbFRlbXBsYXRlQXR0cmlidXRlKSB7XG4gICAgICBzd2l0Y2ggKHR5cGUpIHtcbiAgICAgICAgY2FzZSBlLkJpbmRpbmdUeXBlLlByb3BlcnR5OlxuICAgICAgICBjYXNlIGUuQmluZGluZ1R5cGUuQ2xhc3M6XG4gICAgICAgIGNhc2UgZS5CaW5kaW5nVHlwZS5TdHlsZTpcbiAgICAgICAgICAvLyBCZWNhdXNlIHRoaXMgYmluZGluZyBkb2Vzbid0IHJlYWxseSB0YXJnZXQgdGhlIG5nLXRlbXBsYXRlLCBpdCBtdXN0IGJlIGEgYmluZGluZyBvbiBhblxuICAgICAgICAgIC8vIGlubmVyIG5vZGUgb2YgYSBzdHJ1Y3R1cmFsIHRlbXBsYXRlLiBXZSBjYW4ndCBza2lwIGl0IGVudGlyZWx5LCBiZWNhdXNlIHdlIHN0aWxsIG5lZWRcbiAgICAgICAgICAvLyBpdCBvbiB0aGUgbmctdGVtcGxhdGUncyBjb25zdHMgKGUuZy4gZm9yIHRoZSBwdXJwb3NlcyBvZiBkaXJlY3RpdmUgbWF0Y2hpbmcpLiBIb3dldmVyLFxuICAgICAgICAgIC8vIHdlIHNob3VsZCBub3QgZ2VuZXJhdGUgYW4gdXBkYXRlIGluc3RydWN0aW9uIGZvciBpdC5cbiAgICAgICAgICByZXR1cm4gaXIuY3JlYXRlRXh0cmFjdGVkQXR0cmlidXRlT3AoXG4gICAgICAgICAgICB4cmVmLFxuICAgICAgICAgICAgaXIuQmluZGluZ0tpbmQuUHJvcGVydHksXG4gICAgICAgICAgICBudWxsLFxuICAgICAgICAgICAgbmFtZSxcbiAgICAgICAgICAgIG51bGwsXG4gICAgICAgICAgICBudWxsLFxuICAgICAgICAgICAgaTE4bk1lc3NhZ2UsXG4gICAgICAgICAgICBzZWN1cml0eUNvbnRleHQsXG4gICAgICAgICAgKTtcbiAgICAgICAgY2FzZSBlLkJpbmRpbmdUeXBlLlR3b1dheTpcbiAgICAgICAgICByZXR1cm4gaXIuY3JlYXRlRXh0cmFjdGVkQXR0cmlidXRlT3AoXG4gICAgICAgICAgICB4cmVmLFxuICAgICAgICAgICAgaXIuQmluZGluZ0tpbmQuVHdvV2F5UHJvcGVydHksXG4gICAgICAgICAgICBudWxsLFxuICAgICAgICAgICAgbmFtZSxcbiAgICAgICAgICAgIG51bGwsXG4gICAgICAgICAgICBudWxsLFxuICAgICAgICAgICAgaTE4bk1lc3NhZ2UsXG4gICAgICAgICAgICBzZWN1cml0eUNvbnRleHQsXG4gICAgICAgICAgKTtcbiAgICAgIH1cbiAgICB9XG5cbiAgICBpZiAoIWlzVGV4dEJpbmRpbmcgJiYgKHR5cGUgPT09IGUuQmluZGluZ1R5cGUuQXR0cmlidXRlIHx8IHR5cGUgPT09IGUuQmluZGluZ1R5cGUuQW5pbWF0aW9uKSkge1xuICAgICAgLy8gQWdhaW4sIHRoaXMgYmluZGluZyBkb2Vzbid0IHJlYWxseSB0YXJnZXQgdGhlIG5nLXRlbXBsYXRlOyBpdCBhY3R1YWxseSB0YXJnZXRzIHRoZSBlbGVtZW50XG4gICAgICAvLyBpbnNpZGUgdGhlIHN0cnVjdHVyYWwgdGVtcGxhdGUuIEluIHRoZSBjYXNlIG9mIG5vbi10ZXh0IGF0dHJpYnV0ZSBvciBhbmltYXRpb24gYmluZGluZ3MsXG4gICAgICAvLyB0aGUgYmluZGluZyBkb2Vzbid0IGV2ZW4gc2hvdyB1cCBvbiB0aGUgbmctdGVtcGxhdGUgY29uc3QgYXJyYXksIHNvIHdlIGp1c3Qgc2tpcCBpdFxuICAgICAgLy8gZW50aXJlbHkuXG4gICAgICByZXR1cm4gbnVsbDtcbiAgICB9XG4gIH1cblxuICBsZXQgYmluZGluZ1R5cGUgPSBCSU5ESU5HX0tJTkRTLmdldCh0eXBlKSE7XG5cbiAgaWYgKHRlbXBsYXRlS2luZCA9PT0gaXIuVGVtcGxhdGVLaW5kLk5nVGVtcGxhdGUpIHtcbiAgICAvLyBXZSBrbm93IHdlIGFyZSBkZWFsaW5nIHdpdGggYmluZGluZ3MgZGlyZWN0bHkgb24gYW4gZXhwbGljaXQgbmctdGVtcGxhdGUuXG4gICAgLy8gU3RhdGljIGF0dHJpYnV0ZSBiaW5kaW5ncyBzaG91bGQgYmUgY29sbGVjdGVkIGludG8gdGhlIGNvbnN0IGFycmF5IGFzIGsvdiBwYWlycy4gUHJvcGVydHlcbiAgICAvLyBiaW5kaW5ncyBzaG91bGQgcmVzdWx0IGluIGEgYHByb3BlcnR5YCBpbnN0cnVjdGlvbiwgYW5kIGBBdHRyaWJ1dGVNYXJrZXIuQmluZGluZ3NgIGNvbnN0XG4gICAgLy8gZW50cmllcy5cbiAgICAvL1xuICAgIC8vIFRoZSBkaWZmaWN1bHR5IGlzIHdpdGggZHluYW1pYyBhdHRyaWJ1dGUsIHN0eWxlLCBhbmQgY2xhc3MgYmluZGluZ3MuIFRoZXNlIGRvbid0IHJlYWxseSBtYWtlXG4gICAgLy8gc2Vuc2Ugb24gYW4gYG5nLXRlbXBsYXRlYCBhbmQgc2hvdWxkIHByb2JhYmx5IGJlIHBhcnNlciBlcnJvcnMuIEhvd2V2ZXIsXG4gICAgLy8gVGVtcGxhdGVEZWZpbml0aW9uQnVpbGRlciBnZW5lcmF0ZXMgYHByb3BlcnR5YCBpbnN0cnVjdGlvbnMgZm9yIHRoZW0sIGFuZCBzbyB3ZSBkbyB0aGF0IGFzXG4gICAgLy8gd2VsbC5cbiAgICAvL1xuICAgIC8vIE5vdGUgdGhhdCB3ZSBkbyBoYXZlIGEgc2xpZ2h0IGJlaGF2aW9yIGRpZmZlcmVuY2Ugd2l0aCBUZW1wbGF0ZURlZmluaXRpb25CdWlsZGVyOiBhbHRob3VnaFxuICAgIC8vIFREQiBlbWl0cyBgcHJvcGVydHlgIGluc3RydWN0aW9ucyBmb3IgZHluYW1pYyBhdHRyaWJ1dGVzLCBzdHlsZXMsIGFuZCBjbGFzc2VzLCBvbmx5IHN0eWxlc1xuICAgIC8vIGFuZCBjbGFzc2VzIGFsc28gZ2V0IGNvbnN0IGNvbGxlY3RlZCBpbnRvIHRoZSBgQXR0cmlidXRlTWFya2VyLkJpbmRpbmdzYCBmaWVsZC4gRHluYW1pY1xuICAgIC8vIGF0dHJpYnV0ZSBiaW5kaW5ncyBhcmUgbWlzc2luZyBmcm9tIHRoZSBjb25zdHMgZW50aXJlbHkuIFdlIGNob29zZSB0byBlbWl0IHRoZW0gaW50byB0aGVcbiAgICAvLyBjb25zdHMgZmllbGQgYW55d2F5LCB0byBhdm9pZCBjcmVhdGluZyBzcGVjaWFsIGNhc2VzIGZvciBzb21ldGhpbmcgc28gYXJjYW5lIGFuZCBub25zZW5zaWNhbC5cbiAgICBpZiAoXG4gICAgICB0eXBlID09PSBlLkJpbmRpbmdUeXBlLkNsYXNzIHx8XG4gICAgICB0eXBlID09PSBlLkJpbmRpbmdUeXBlLlN0eWxlIHx8XG4gICAgICAodHlwZSA9PT0gZS5CaW5kaW5nVHlwZS5BdHRyaWJ1dGUgJiYgIWlzVGV4dEJpbmRpbmcpXG4gICAgKSB7XG4gICAgICAvLyBUT0RPOiBUaGVzZSBjYXNlcyBzaG91bGQgYmUgcGFyc2UgZXJyb3JzLlxuICAgICAgYmluZGluZ1R5cGUgPSBpci5CaW5kaW5nS2luZC5Qcm9wZXJ0eTtcbiAgICB9XG4gIH1cblxuICByZXR1cm4gaXIuY3JlYXRlQmluZGluZ09wKFxuICAgIHhyZWYsXG4gICAgYmluZGluZ1R5cGUsXG4gICAgbmFtZSxcbiAgICBjb252ZXJ0QXN0V2l0aEludGVycG9sYXRpb24odmlldy5qb2IsIHZhbHVlLCBpMThuTWVzc2FnZSksXG4gICAgdW5pdCxcbiAgICBzZWN1cml0eUNvbnRleHQsXG4gICAgaXNUZXh0QmluZGluZyxcbiAgICBpc1N0cnVjdHVyYWxUZW1wbGF0ZUF0dHJpYnV0ZSxcbiAgICB0ZW1wbGF0ZUtpbmQsXG4gICAgaTE4bk1lc3NhZ2UsXG4gICAgc291cmNlU3BhbixcbiAgKTtcbn1cblxuZnVuY3Rpb24gbWFrZUxpc3RlbmVySGFuZGxlck9wcyhcbiAgdW5pdDogQ29tcGlsYXRpb25Vbml0LFxuICBoYW5kbGVyOiBlLkFTVCxcbiAgaGFuZGxlclNwYW46IFBhcnNlU291cmNlU3Bhbixcbik6IGlyLlVwZGF0ZU9wW10ge1xuICBoYW5kbGVyID0gYXN0T2YoaGFuZGxlcik7XG4gIGNvbnN0IGhhbmRsZXJPcHMgPSBuZXcgQXJyYXk8aXIuVXBkYXRlT3A+KCk7XG4gIGxldCBoYW5kbGVyRXhwcnM6IGUuQVNUW10gPSBoYW5kbGVyIGluc3RhbmNlb2YgZS5DaGFpbiA/IGhhbmRsZXIuZXhwcmVzc2lvbnMgOiBbaGFuZGxlcl07XG4gIGlmIChoYW5kbGVyRXhwcnMubGVuZ3RoID09PSAwKSB7XG4gICAgdGhyb3cgbmV3IEVycm9yKCdFeHBlY3RlZCBsaXN0ZW5lciB0byBoYXZlIG5vbi1lbXB0eSBleHByZXNzaW9uIGxpc3QuJyk7XG4gIH1cbiAgY29uc3QgZXhwcmVzc2lvbnMgPSBoYW5kbGVyRXhwcnMubWFwKChleHByKSA9PiBjb252ZXJ0QXN0KGV4cHIsIHVuaXQuam9iLCBoYW5kbGVyU3BhbikpO1xuICBjb25zdCByZXR1cm5FeHByID0gZXhwcmVzc2lvbnMucG9wKCkhO1xuICBoYW5kbGVyT3BzLnB1c2goXG4gICAgLi4uZXhwcmVzc2lvbnMubWFwKChlKSA9PlxuICAgICAgaXIuY3JlYXRlU3RhdGVtZW50T3A8aXIuVXBkYXRlT3A+KG5ldyBvLkV4cHJlc3Npb25TdGF0ZW1lbnQoZSwgZS5zb3VyY2VTcGFuKSksXG4gICAgKSxcbiAgKTtcbiAgaGFuZGxlck9wcy5wdXNoKGlyLmNyZWF0ZVN0YXRlbWVudE9wKG5ldyBvLlJldHVyblN0YXRlbWVudChyZXR1cm5FeHByLCByZXR1cm5FeHByLnNvdXJjZVNwYW4pKSk7XG4gIHJldHVybiBoYW5kbGVyT3BzO1xufVxuXG5mdW5jdGlvbiBtYWtlVHdvV2F5TGlzdGVuZXJIYW5kbGVyT3BzKFxuICB1bml0OiBDb21waWxhdGlvblVuaXQsXG4gIGhhbmRsZXI6IGUuQVNULFxuICBoYW5kbGVyU3BhbjogUGFyc2VTb3VyY2VTcGFuLFxuKTogaXIuVXBkYXRlT3BbXSB7XG4gIGhhbmRsZXIgPSBhc3RPZihoYW5kbGVyKTtcbiAgY29uc3QgaGFuZGxlck9wcyA9IG5ldyBBcnJheTxpci5VcGRhdGVPcD4oKTtcblxuICBpZiAoaGFuZGxlciBpbnN0YW5jZW9mIGUuQ2hhaW4pIHtcbiAgICBpZiAoaGFuZGxlci5leHByZXNzaW9ucy5sZW5ndGggPT09IDEpIHtcbiAgICAgIGhhbmRsZXIgPSBoYW5kbGVyLmV4cHJlc3Npb25zWzBdO1xuICAgIH0gZWxzZSB7XG4gICAgICAvLyBUaGlzIGlzIHZhbGlkYXRlZCBkdXJpbmcgcGFyc2luZyBhbHJlYWR5LCBidXQgd2UgZG8gaXQgaGVyZSBqdXN0IGluIGNhc2UuXG4gICAgICB0aHJvdyBuZXcgRXJyb3IoJ0V4cGVjdGVkIHR3by13YXkgbGlzdGVuZXIgdG8gaGF2ZSBhIHNpbmdsZSBleHByZXNzaW9uLicpO1xuICAgIH1cbiAgfVxuXG4gIGNvbnN0IGhhbmRsZXJFeHByID0gY29udmVydEFzdChoYW5kbGVyLCB1bml0LmpvYiwgaGFuZGxlclNwYW4pO1xuICBjb25zdCBldmVudFJlZmVyZW5jZSA9IG5ldyBpci5MZXhpY2FsUmVhZEV4cHIoJyRldmVudCcpO1xuICBjb25zdCB0d29XYXlTZXRFeHByID0gbmV3IGlyLlR3b1dheUJpbmRpbmdTZXRFeHByKGhhbmRsZXJFeHByLCBldmVudFJlZmVyZW5jZSk7XG5cbiAgaGFuZGxlck9wcy5wdXNoKGlyLmNyZWF0ZVN0YXRlbWVudE9wPGlyLlVwZGF0ZU9wPihuZXcgby5FeHByZXNzaW9uU3RhdGVtZW50KHR3b1dheVNldEV4cHIpKSk7XG4gIGhhbmRsZXJPcHMucHVzaChpci5jcmVhdGVTdGF0ZW1lbnRPcChuZXcgby5SZXR1cm5TdGF0ZW1lbnQoZXZlbnRSZWZlcmVuY2UpKSk7XG4gIHJldHVybiBoYW5kbGVyT3BzO1xufVxuXG5mdW5jdGlvbiBhc3RPZihhc3Q6IGUuQVNUIHwgZS5BU1RXaXRoU291cmNlKTogZS5BU1Qge1xuICByZXR1cm4gYXN0IGluc3RhbmNlb2YgZS5BU1RXaXRoU291cmNlID8gYXN0LmFzdCA6IGFzdDtcbn1cblxuLyoqXG4gKiBQcm9jZXNzIGFsbCBvZiB0aGUgbG9jYWwgcmVmZXJlbmNlcyBvbiBhbiBlbGVtZW50LWxpa2Ugc3RydWN0dXJlIGluIHRoZSB0ZW1wbGF0ZSBBU1QgYW5kXG4gKiBjb252ZXJ0IHRoZW0gdG8gdGhlaXIgSVIgcmVwcmVzZW50YXRpb24uXG4gKi9cbmZ1bmN0aW9uIGluZ2VzdFJlZmVyZW5jZXMob3A6IGlyLkVsZW1lbnRPcEJhc2UsIGVsZW1lbnQ6IHQuRWxlbWVudCB8IHQuVGVtcGxhdGUpOiB2b2lkIHtcbiAgYXNzZXJ0SXNBcnJheTxpci5Mb2NhbFJlZj4ob3AubG9jYWxSZWZzKTtcbiAgZm9yIChjb25zdCB7bmFtZSwgdmFsdWV9IG9mIGVsZW1lbnQucmVmZXJlbmNlcykge1xuICAgIG9wLmxvY2FsUmVmcy5wdXNoKHtcbiAgICAgIG5hbWUsXG4gICAgICB0YXJnZXQ6IHZhbHVlLFxuICAgIH0pO1xuICB9XG59XG5cbi8qKlxuICogQXNzZXJ0IHRoYXQgdGhlIGdpdmVuIHZhbHVlIGlzIGFuIGFycmF5LlxuICovXG5mdW5jdGlvbiBhc3NlcnRJc0FycmF5PFQ+KHZhbHVlOiBhbnkpOiBhc3NlcnRzIHZhbHVlIGlzIEFycmF5PFQ+IHtcbiAgaWYgKCFBcnJheS5pc0FycmF5KHZhbHVlKSkge1xuICAgIHRocm93IG5ldyBFcnJvcihgQXNzZXJ0aW9uRXJyb3I6IGV4cGVjdGVkIGFuIGFycmF5YCk7XG4gIH1cbn1cblxuLyoqXG4gKiBDcmVhdGVzIGFuIGFic29sdXRlIGBQYXJzZVNvdXJjZVNwYW5gIGZyb20gdGhlIHJlbGF0aXZlIGBQYXJzZVNwYW5gLlxuICpcbiAqIGBQYXJzZVNwYW5gIG9iamVjdHMgYXJlIHJlbGF0aXZlIHRvIHRoZSBzdGFydCBvZiB0aGUgZXhwcmVzc2lvbi5cbiAqIFRoaXMgbWV0aG9kIGNvbnZlcnRzIHRoZXNlIHRvIGZ1bGwgYFBhcnNlU291cmNlU3BhbmAgb2JqZWN0cyB0aGF0XG4gKiBzaG93IHdoZXJlIHRoZSBzcGFuIGlzIHdpdGhpbiB0aGUgb3ZlcmFsbCBzb3VyY2UgZmlsZS5cbiAqXG4gKiBAcGFyYW0gc3BhbiB0aGUgcmVsYXRpdmUgc3BhbiB0byBjb252ZXJ0LlxuICogQHBhcmFtIGJhc2VTb3VyY2VTcGFuIGEgc3BhbiBjb3JyZXNwb25kaW5nIHRvIHRoZSBiYXNlIG9mIHRoZSBleHByZXNzaW9uIHRyZWUuXG4gKiBAcmV0dXJucyBhIGBQYXJzZVNvdXJjZVNwYW5gIGZvciB0aGUgZ2l2ZW4gc3BhbiBvciBudWxsIGlmIG5vIGBiYXNlU291cmNlU3BhbmAgd2FzIHByb3ZpZGVkLlxuICovXG5mdW5jdGlvbiBjb252ZXJ0U291cmNlU3BhbihcbiAgc3BhbjogZS5QYXJzZVNwYW4sXG4gIGJhc2VTb3VyY2VTcGFuOiBQYXJzZVNvdXJjZVNwYW4gfCBudWxsLFxuKTogUGFyc2VTb3VyY2VTcGFuIHwgbnVsbCB7XG4gIGlmIChiYXNlU291cmNlU3BhbiA9PT0gbnVsbCkge1xuICAgIHJldHVybiBudWxsO1xuICB9XG4gIGNvbnN0IHN0YXJ0ID0gYmFzZVNvdXJjZVNwYW4uc3RhcnQubW92ZUJ5KHNwYW4uc3RhcnQpO1xuICBjb25zdCBlbmQgPSBiYXNlU291cmNlU3Bhbi5zdGFydC5tb3ZlQnkoc3Bhbi5lbmQpO1xuICBjb25zdCBmdWxsU3RhcnQgPSBiYXNlU291cmNlU3Bhbi5mdWxsU3RhcnQubW92ZUJ5KHNwYW4uc3RhcnQpO1xuICByZXR1cm4gbmV3IFBhcnNlU291cmNlU3BhbihzdGFydCwgZW5kLCBmdWxsU3RhcnQpO1xufVxuXG4vKipcbiAqIFdpdGggdGhlIGRpcmVjdGl2ZS1iYXNlZCBjb250cm9sIGZsb3cgdXNlcnMgd2VyZSBhYmxlIHRvIGNvbmRpdGlvbmFsbHkgcHJvamVjdCBjb250ZW50IHVzaW5nXG4gKiB0aGUgYCpgIHN5bnRheC4gRS5nLiBgPGRpdiAqbmdJZj1cImV4cHJcIiBwcm9qZWN0TWU+PC9kaXY+YCB3aWxsIGJlIHByb2plY3RlZCBpbnRvXG4gKiBgPG5nLWNvbnRlbnQgc2VsZWN0PVwiW3Byb2plY3RNZV1cIi8+YCwgYmVjYXVzZSB0aGUgYXR0cmlidXRlcyBhbmQgdGFnIG5hbWUgZnJvbSB0aGUgYGRpdmAgYXJlXG4gKiBjb3BpZWQgdG8gdGhlIHRlbXBsYXRlIHZpYSB0aGUgdGVtcGxhdGUgY3JlYXRpb24gaW5zdHJ1Y3Rpb24uIFdpdGggYEBpZmAgYW5kIGBAZm9yYCB0aGF0IGlzXG4gKiBub3QgdGhlIGNhc2UsIGJlY2F1c2UgdGhlIGNvbmRpdGlvbmFsIGlzIHBsYWNlZCAqYXJvdW5kKiBlbGVtZW50cywgcmF0aGVyIHRoYW4gKm9uKiB0aGVtLlxuICogVGhlIHJlc3VsdCBpcyB0aGF0IGNvbnRlbnQgcHJvamVjdGlvbiB3b24ndCB3b3JrIGluIHRoZSBzYW1lIHdheSBpZiBhIHVzZXIgY29udmVydHMgZnJvbVxuICogYCpuZ0lmYCB0byBgQGlmYC5cbiAqXG4gKiBUaGlzIGZ1bmN0aW9uIGFpbXMgdG8gY292ZXIgdGhlIG1vc3QgY29tbW9uIGNhc2UgYnkgZG9pbmcgdGhlIHNhbWUgY29weWluZyB3aGVuIGEgY29udHJvbCBmbG93XG4gKiBub2RlIGhhcyAqb25lIGFuZCBvbmx5IG9uZSogcm9vdCBlbGVtZW50IG9yIHRlbXBsYXRlIG5vZGUuXG4gKlxuICogVGhpcyBhcHByb2FjaCBjb21lcyB3aXRoIHNvbWUgY2F2ZWF0czpcbiAqIDEuIEFzIHNvb24gYXMgYW55IG90aGVyIG5vZGUgaXMgYWRkZWQgdG8gdGhlIHJvb3QsIHRoZSBjb3B5aW5nIGJlaGF2aW9yIHdvbid0IHdvcmsgYW55bW9yZS5cbiAqICAgIEEgZGlhZ25vc3RpYyB3aWxsIGJlIGFkZGVkIHRvIGZsYWcgY2FzZXMgbGlrZSB0aGlzIGFuZCB0byBleHBsYWluIGhvdyB0byB3b3JrIGFyb3VuZCBpdC5cbiAqIDIuIElmIGBwcmVzZXJ2ZVdoaXRlc3BhY2VzYCBpcyBlbmFibGVkLCBpdCdzIHZlcnkgbGlrZWx5IHRoYXQgaW5kZW50YXRpb24gd2lsbCBicmVhayB0aGlzXG4gKiAgICB3b3JrYXJvdW5kLCBiZWNhdXNlIGl0J2xsIGluY2x1ZGUgYW4gYWRkaXRpb25hbCB0ZXh0IG5vZGUgYXMgdGhlIGZpcnN0IGNoaWxkLiBXZSBjYW4gd29ya1xuICogICAgYXJvdW5kIGl0IGhlcmUsIGJ1dCBpbiBhIGRpc2N1c3Npb24gaXQgd2FzIGRlY2lkZWQgbm90IHRvLCBiZWNhdXNlIHRoZSB1c2VyIGV4cGxpY2l0bHkgb3B0ZWRcbiAqICAgIGludG8gcHJlc2VydmluZyB0aGUgd2hpdGVzcGFjZSBhbmQgd2Ugd291bGQgaGF2ZSB0byBkcm9wIGl0IGZyb20gdGhlIGdlbmVyYXRlZCBjb2RlLlxuICogICAgVGhlIGRpYWdub3N0aWMgbWVudGlvbmVkIHBvaW50ICMxIHdpbGwgZmxhZyBzdWNoIGNhc2VzIHRvIHVzZXJzLlxuICpcbiAqIEByZXR1cm5zIFRhZyBuYW1lIHRvIGJlIHVzZWQgZm9yIHRoZSBjb250cm9sIGZsb3cgdGVtcGxhdGUuXG4gKi9cbmZ1bmN0aW9uIGluZ2VzdENvbnRyb2xGbG93SW5zZXJ0aW9uUG9pbnQoXG4gIHVuaXQ6IFZpZXdDb21waWxhdGlvblVuaXQsXG4gIHhyZWY6IGlyLlhyZWZJZCxcbiAgbm9kZTogdC5JZkJsb2NrQnJhbmNoIHwgdC5Td2l0Y2hCbG9ja0Nhc2UgfCB0LkZvckxvb3BCbG9jayB8IHQuRm9yTG9vcEJsb2NrRW1wdHksXG4pOiBzdHJpbmcgfCBudWxsIHtcbiAgbGV0IHJvb3Q6IHQuRWxlbWVudCB8IHQuVGVtcGxhdGUgfCBudWxsID0gbnVsbDtcblxuICBmb3IgKGNvbnN0IGNoaWxkIG9mIG5vZGUuY2hpbGRyZW4pIHtcbiAgICAvLyBTa2lwIG92ZXIgY29tbWVudCBub2Rlcy5cbiAgICBpZiAoY2hpbGQgaW5zdGFuY2VvZiB0LkNvbW1lbnQpIHtcbiAgICAgIGNvbnRpbnVlO1xuICAgIH1cblxuICAgIC8vIFdlIGNhbiBvbmx5IGluZmVyIHRoZSB0YWcgbmFtZS9hdHRyaWJ1dGVzIGlmIHRoZXJlJ3MgYSBzaW5nbGUgcm9vdCBub2RlLlxuICAgIGlmIChyb290ICE9PSBudWxsKSB7XG4gICAgICByZXR1cm4gbnVsbDtcbiAgICB9XG5cbiAgICAvLyBSb290IG5vZGVzIGNhbiBvbmx5IGVsZW1lbnRzIG9yIHRlbXBsYXRlcyB3aXRoIGEgdGFnIG5hbWUgKGUuZy4gYDxkaXYgKmZvbz48L2Rpdj5gKS5cbiAgICBpZiAoY2hpbGQgaW5zdGFuY2VvZiB0LkVsZW1lbnQgfHwgKGNoaWxkIGluc3RhbmNlb2YgdC5UZW1wbGF0ZSAmJiBjaGlsZC50YWdOYW1lICE9PSBudWxsKSkge1xuICAgICAgcm9vdCA9IGNoaWxkO1xuICAgIH1cbiAgfVxuXG4gIC8vIElmIHdlJ3ZlIGZvdW5kIGEgc2luZ2xlIHJvb3Qgbm9kZSwgaXRzIHRhZyBuYW1lIGFuZCBhdHRyaWJ1dGVzIGNhbiBiZVxuICAvLyBjb3BpZWQgdG8gdGhlIHN1cnJvdW5kaW5nIHRlbXBsYXRlIHRvIGJlIHVzZWQgZm9yIGNvbnRlbnQgcHJvamVjdGlvbi5cbiAgaWYgKHJvb3QgIT09IG51bGwpIHtcbiAgICAvLyBDb2xsZWN0IHRoZSBzdGF0aWMgYXR0cmlidXRlcyBmb3IgY29udGVudCBwcm9qZWN0aW9uIHB1cnBvc2VzLlxuICAgIGZvciAoY29uc3QgYXR0ciBvZiByb290LmF0dHJpYnV0ZXMpIHtcbiAgICAgIGNvbnN0IHNlY3VyaXR5Q29udGV4dCA9IGRvbVNjaGVtYS5zZWN1cml0eUNvbnRleHQoTkdfVEVNUExBVEVfVEFHX05BTUUsIGF0dHIubmFtZSwgdHJ1ZSk7XG4gICAgICB1bml0LnVwZGF0ZS5wdXNoKFxuICAgICAgICBpci5jcmVhdGVCaW5kaW5nT3AoXG4gICAgICAgICAgeHJlZixcbiAgICAgICAgICBpci5CaW5kaW5nS2luZC5BdHRyaWJ1dGUsXG4gICAgICAgICAgYXR0ci5uYW1lLFxuICAgICAgICAgIG8ubGl0ZXJhbChhdHRyLnZhbHVlKSxcbiAgICAgICAgICBudWxsLFxuICAgICAgICAgIHNlY3VyaXR5Q29udGV4dCxcbiAgICAgICAgICB0cnVlLFxuICAgICAgICAgIGZhbHNlLFxuICAgICAgICAgIG51bGwsXG4gICAgICAgICAgYXNNZXNzYWdlKGF0dHIuaTE4biksXG4gICAgICAgICAgYXR0ci5zb3VyY2VTcGFuLFxuICAgICAgICApLFxuICAgICAgKTtcbiAgICB9XG5cbiAgICAvLyBBbHNvIGNvbGxlY3QgdGhlIGlucHV0cyBzaW5jZSB0aGV5IHBhcnRpY2lwYXRlIGluIGNvbnRlbnQgcHJvamVjdGlvbiBhcyB3ZWxsLlxuICAgIC8vIE5vdGUgdGhhdCBUREIgdXNlZCB0byBjb2xsZWN0IHRoZSBvdXRwdXRzIGFzIHdlbGwsIGJ1dCBpdCB3YXNuJ3QgcGFzc2luZyB0aGVtIGludG9cbiAgICAvLyB0aGUgdGVtcGxhdGUgaW5zdHJ1Y3Rpb24uIEhlcmUgd2UganVzdCBkb24ndCBjb2xsZWN0IHRoZW0uXG4gICAgZm9yIChjb25zdCBhdHRyIG9mIHJvb3QuaW5wdXRzKSB7XG4gICAgICBpZiAoYXR0ci50eXBlICE9PSBlLkJpbmRpbmdUeXBlLkFuaW1hdGlvbiAmJiBhdHRyLnR5cGUgIT09IGUuQmluZGluZ1R5cGUuQXR0cmlidXRlKSB7XG4gICAgICAgIGNvbnN0IHNlY3VyaXR5Q29udGV4dCA9IGRvbVNjaGVtYS5zZWN1cml0eUNvbnRleHQoTkdfVEVNUExBVEVfVEFHX05BTUUsIGF0dHIubmFtZSwgdHJ1ZSk7XG4gICAgICAgIHVuaXQuY3JlYXRlLnB1c2goXG4gICAgICAgICAgaXIuY3JlYXRlRXh0cmFjdGVkQXR0cmlidXRlT3AoXG4gICAgICAgICAgICB4cmVmLFxuICAgICAgICAgICAgaXIuQmluZGluZ0tpbmQuUHJvcGVydHksXG4gICAgICAgICAgICBudWxsLFxuICAgICAgICAgICAgYXR0ci5uYW1lLFxuICAgICAgICAgICAgbnVsbCxcbiAgICAgICAgICAgIG51bGwsXG4gICAgICAgICAgICBudWxsLFxuICAgICAgICAgICAgc2VjdXJpdHlDb250ZXh0LFxuICAgICAgICAgICksXG4gICAgICAgICk7XG4gICAgICB9XG4gICAgfVxuXG4gICAgY29uc3QgdGFnTmFtZSA9IHJvb3QgaW5zdGFuY2VvZiB0LkVsZW1lbnQgPyByb290Lm5hbWUgOiByb290LnRhZ05hbWU7XG5cbiAgICAvLyBEb24ndCBwYXNzIGFsb25nIGBuZy10ZW1wbGF0ZWAgdGFnIG5hbWUgc2luY2UgaXQgZW5hYmxlcyBkaXJlY3RpdmUgbWF0Y2hpbmcuXG4gICAgcmV0dXJuIHRhZ05hbWUgPT09IE5HX1RFTVBMQVRFX1RBR19OQU1FID8gbnVsbCA6IHRhZ05hbWU7XG4gIH1cblxuICByZXR1cm4gbnVsbDtcbn1cbiJdfQ==