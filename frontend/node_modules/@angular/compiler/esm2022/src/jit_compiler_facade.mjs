/**
 * @license
 * Copyright Google LLC All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
import { ConstantPool } from './constant_pool';
import { ChangeDetectionStrategy, ViewEncapsulation, } from './core';
import { compileInjectable } from './injectable_compiler_2';
import { DEFAULT_INTERPOLATION_CONFIG, InterpolationConfig } from './ml_parser/defaults';
import { DeclareVarStmt, literal, LiteralExpr, StmtModifier, WrappedNodeExpr, } from './output/output_ast';
import { JitEvaluator } from './output/output_jit';
import { r3JitTypeSourceSpan } from './parse_util';
import { compileFactoryFunction, FactoryTarget } from './render3/r3_factory';
import { compileInjector } from './render3/r3_injector_compiler';
import { R3JitReflector } from './render3/r3_jit';
import { compileNgModule, compileNgModuleDeclarationExpression, R3NgModuleMetadataKind, R3SelectorScopeMode, } from './render3/r3_module_compiler';
import { compilePipeFromMetadata } from './render3/r3_pipe_compiler';
import { createMayBeForwardRefExpression, getSafePropertyAccessString, wrapReference, } from './render3/util';
import { R3TemplateDependencyKind, } from './render3/view/api';
import { compileComponentFromMetadata, compileDirectiveFromMetadata, parseHostBindings, verifyHostBindings, } from './render3/view/compiler';
import { R3TargetBinder } from './render3/view/t2_binder';
import { makeBindingParser, parseTemplate } from './render3/view/template';
import { ResourceLoader } from './resource_loader';
import { DomElementSchemaRegistry } from './schema/dom_element_schema_registry';
import { SelectorMatcher } from './selector';
export class CompilerFacadeImpl {
    constructor(jitEvaluator = new JitEvaluator()) {
        this.jitEvaluator = jitEvaluator;
        this.FactoryTarget = FactoryTarget;
        this.ResourceLoader = ResourceLoader;
        this.elementSchemaRegistry = new DomElementSchemaRegistry();
    }
    compilePipe(angularCoreEnv, sourceMapUrl, facade) {
        const metadata = {
            name: facade.name,
            type: wrapReference(facade.type),
            typeArgumentCount: 0,
            deps: null,
            pipeName: facade.pipeName,
            pure: facade.pure,
            isStandalone: facade.isStandalone,
        };
        const res = compilePipeFromMetadata(metadata);
        return this.jitExpression(res.expression, angularCoreEnv, sourceMapUrl, []);
    }
    compilePipeDeclaration(angularCoreEnv, sourceMapUrl, declaration) {
        const meta = convertDeclarePipeFacadeToMetadata(declaration);
        const res = compilePipeFromMetadata(meta);
        return this.jitExpression(res.expression, angularCoreEnv, sourceMapUrl, []);
    }
    compileInjectable(angularCoreEnv, sourceMapUrl, facade) {
        const { expression, statements } = compileInjectable({
            name: facade.name,
            type: wrapReference(facade.type),
            typeArgumentCount: facade.typeArgumentCount,
            providedIn: computeProvidedIn(facade.providedIn),
            useClass: convertToProviderExpression(facade, 'useClass'),
            useFactory: wrapExpression(facade, 'useFactory'),
            useValue: convertToProviderExpression(facade, 'useValue'),
            useExisting: convertToProviderExpression(facade, 'useExisting'),
            deps: facade.deps?.map(convertR3DependencyMetadata),
        }, 
        /* resolveForwardRefs */ true);
        return this.jitExpression(expression, angularCoreEnv, sourceMapUrl, statements);
    }
    compileInjectableDeclaration(angularCoreEnv, sourceMapUrl, facade) {
        const { expression, statements } = compileInjectable({
            name: facade.type.name,
            type: wrapReference(facade.type),
            typeArgumentCount: 0,
            providedIn: computeProvidedIn(facade.providedIn),
            useClass: convertToProviderExpression(facade, 'useClass'),
            useFactory: wrapExpression(facade, 'useFactory'),
            useValue: convertToProviderExpression(facade, 'useValue'),
            useExisting: convertToProviderExpression(facade, 'useExisting'),
            deps: facade.deps?.map(convertR3DeclareDependencyMetadata),
        }, 
        /* resolveForwardRefs */ true);
        return this.jitExpression(expression, angularCoreEnv, sourceMapUrl, statements);
    }
    compileInjector(angularCoreEnv, sourceMapUrl, facade) {
        const meta = {
            name: facade.name,
            type: wrapReference(facade.type),
            providers: facade.providers && facade.providers.length > 0
                ? new WrappedNodeExpr(facade.providers)
                : null,
            imports: facade.imports.map((i) => new WrappedNodeExpr(i)),
        };
        const res = compileInjector(meta);
        return this.jitExpression(res.expression, angularCoreEnv, sourceMapUrl, []);
    }
    compileInjectorDeclaration(angularCoreEnv, sourceMapUrl, declaration) {
        const meta = convertDeclareInjectorFacadeToMetadata(declaration);
        const res = compileInjector(meta);
        return this.jitExpression(res.expression, angularCoreEnv, sourceMapUrl, []);
    }
    compileNgModule(angularCoreEnv, sourceMapUrl, facade) {
        const meta = {
            kind: R3NgModuleMetadataKind.Global,
            type: wrapReference(facade.type),
            bootstrap: facade.bootstrap.map(wrapReference),
            declarations: facade.declarations.map(wrapReference),
            publicDeclarationTypes: null, // only needed for types in AOT
            imports: facade.imports.map(wrapReference),
            includeImportTypes: true,
            exports: facade.exports.map(wrapReference),
            selectorScopeMode: R3SelectorScopeMode.Inline,
            containsForwardDecls: false,
            schemas: facade.schemas ? facade.schemas.map(wrapReference) : null,
            id: facade.id ? new WrappedNodeExpr(facade.id) : null,
        };
        const res = compileNgModule(meta);
        return this.jitExpression(res.expression, angularCoreEnv, sourceMapUrl, []);
    }
    compileNgModuleDeclaration(angularCoreEnv, sourceMapUrl, declaration) {
        const expression = compileNgModuleDeclarationExpression(declaration);
        return this.jitExpression(expression, angularCoreEnv, sourceMapUrl, []);
    }
    compileDirective(angularCoreEnv, sourceMapUrl, facade) {
        const meta = convertDirectiveFacadeToMetadata(facade);
        return this.compileDirectiveFromMeta(angularCoreEnv, sourceMapUrl, meta);
    }
    compileDirectiveDeclaration(angularCoreEnv, sourceMapUrl, declaration) {
        const typeSourceSpan = this.createParseSourceSpan('Directive', declaration.type.name, sourceMapUrl);
        const meta = convertDeclareDirectiveFacadeToMetadata(declaration, typeSourceSpan);
        return this.compileDirectiveFromMeta(angularCoreEnv, sourceMapUrl, meta);
    }
    compileDirectiveFromMeta(angularCoreEnv, sourceMapUrl, meta) {
        const constantPool = new ConstantPool();
        const bindingParser = makeBindingParser();
        const res = compileDirectiveFromMetadata(meta, constantPool, bindingParser);
        return this.jitExpression(res.expression, angularCoreEnv, sourceMapUrl, constantPool.statements);
    }
    compileComponent(angularCoreEnv, sourceMapUrl, facade) {
        // Parse the template and check for errors.
        const { template, interpolation, defer } = parseJitTemplate(facade.template, facade.name, sourceMapUrl, facade.preserveWhitespaces, facade.interpolation, undefined);
        // Compile the component metadata, including template, into an expression.
        const meta = {
            ...facade,
            ...convertDirectiveFacadeToMetadata(facade),
            selector: facade.selector || this.elementSchemaRegistry.getDefaultComponentElementName(),
            template,
            declarations: facade.declarations.map(convertDeclarationFacadeToMetadata),
            declarationListEmitMode: 0 /* DeclarationListEmitMode.Direct */,
            defer,
            styles: [...facade.styles, ...template.styles],
            encapsulation: facade.encapsulation,
            interpolation,
            changeDetection: facade.changeDetection ?? null,
            animations: facade.animations != null ? new WrappedNodeExpr(facade.animations) : null,
            viewProviders: facade.viewProviders != null ? new WrappedNodeExpr(facade.viewProviders) : null,
            relativeContextFilePath: '',
            i18nUseExternalIds: true,
        };
        const jitExpressionSourceMap = `ng:///${facade.name}.js`;
        return this.compileComponentFromMeta(angularCoreEnv, jitExpressionSourceMap, meta);
    }
    compileComponentDeclaration(angularCoreEnv, sourceMapUrl, declaration) {
        const typeSourceSpan = this.createParseSourceSpan('Component', declaration.type.name, sourceMapUrl);
        const meta = convertDeclareComponentFacadeToMetadata(declaration, typeSourceSpan, sourceMapUrl);
        return this.compileComponentFromMeta(angularCoreEnv, sourceMapUrl, meta);
    }
    compileComponentFromMeta(angularCoreEnv, sourceMapUrl, meta) {
        const constantPool = new ConstantPool();
        const bindingParser = makeBindingParser(meta.interpolation);
        const res = compileComponentFromMetadata(meta, constantPool, bindingParser);
        return this.jitExpression(res.expression, angularCoreEnv, sourceMapUrl, constantPool.statements);
    }
    compileFactory(angularCoreEnv, sourceMapUrl, meta) {
        const factoryRes = compileFactoryFunction({
            name: meta.name,
            type: wrapReference(meta.type),
            typeArgumentCount: meta.typeArgumentCount,
            deps: convertR3DependencyMetadataArray(meta.deps),
            target: meta.target,
        });
        return this.jitExpression(factoryRes.expression, angularCoreEnv, sourceMapUrl, factoryRes.statements);
    }
    compileFactoryDeclaration(angularCoreEnv, sourceMapUrl, meta) {
        const factoryRes = compileFactoryFunction({
            name: meta.type.name,
            type: wrapReference(meta.type),
            typeArgumentCount: 0,
            deps: Array.isArray(meta.deps)
                ? meta.deps.map(convertR3DeclareDependencyMetadata)
                : meta.deps,
            target: meta.target,
        });
        return this.jitExpression(factoryRes.expression, angularCoreEnv, sourceMapUrl, factoryRes.statements);
    }
    createParseSourceSpan(kind, typeName, sourceUrl) {
        return r3JitTypeSourceSpan(kind, typeName, sourceUrl);
    }
    /**
     * JIT compiles an expression and returns the result of executing that expression.
     *
     * @param def the definition which will be compiled and executed to get the value to patch
     * @param context an object map of @angular/core symbol names to symbols which will be available
     * in the context of the compiled expression
     * @param sourceUrl a URL to use for the source map of the compiled expression
     * @param preStatements a collection of statements that should be evaluated before the expression.
     */
    jitExpression(def, context, sourceUrl, preStatements) {
        // The ConstantPool may contain Statements which declare variables used in the final expression.
        // Therefore, its statements need to precede the actual JIT operation. The final statement is a
        // declaration of $def which is set to the expression being compiled.
        const statements = [
            ...preStatements,
            new DeclareVarStmt('$def', def, undefined, StmtModifier.Exported),
        ];
        const res = this.jitEvaluator.evaluateStatements(sourceUrl, statements, new R3JitReflector(context), 
        /* enableSourceMaps */ true);
        return res['$def'];
    }
}
function convertToR3QueryMetadata(facade) {
    return {
        ...facade,
        isSignal: facade.isSignal,
        predicate: convertQueryPredicate(facade.predicate),
        read: facade.read ? new WrappedNodeExpr(facade.read) : null,
        static: facade.static,
        emitDistinctChangesOnly: facade.emitDistinctChangesOnly,
    };
}
function convertQueryDeclarationToMetadata(declaration) {
    return {
        propertyName: declaration.propertyName,
        first: declaration.first ?? false,
        predicate: convertQueryPredicate(declaration.predicate),
        descendants: declaration.descendants ?? false,
        read: declaration.read ? new WrappedNodeExpr(declaration.read) : null,
        static: declaration.static ?? false,
        emitDistinctChangesOnly: declaration.emitDistinctChangesOnly ?? true,
        isSignal: !!declaration.isSignal,
    };
}
function convertQueryPredicate(predicate) {
    return Array.isArray(predicate)
        ? // The predicate is an array of strings so pass it through.
            predicate
        : // The predicate is a type - assume that we will need to unwrap any `forwardRef()` calls.
            createMayBeForwardRefExpression(new WrappedNodeExpr(predicate), 1 /* ForwardRefHandling.Wrapped */);
}
function convertDirectiveFacadeToMetadata(facade) {
    const inputsFromMetadata = parseInputsArray(facade.inputs || []);
    const outputsFromMetadata = parseMappingStringArray(facade.outputs || []);
    const propMetadata = facade.propMetadata;
    const inputsFromType = {};
    const outputsFromType = {};
    for (const field in propMetadata) {
        if (propMetadata.hasOwnProperty(field)) {
            propMetadata[field].forEach((ann) => {
                if (isInput(ann)) {
                    inputsFromType[field] = {
                        bindingPropertyName: ann.alias || field,
                        classPropertyName: field,
                        required: ann.required || false,
                        // For JIT, decorators are used to declare signal inputs. That is because of
                        // a technical limitation where it's not possible to statically reflect class
                        // members of a directive/component at runtime before instantiating the class.
                        isSignal: !!ann.isSignal,
                        transformFunction: ann.transform != null ? new WrappedNodeExpr(ann.transform) : null,
                    };
                }
                else if (isOutput(ann)) {
                    outputsFromType[field] = ann.alias || field;
                }
            });
        }
    }
    return {
        ...facade,
        typeArgumentCount: 0,
        typeSourceSpan: facade.typeSourceSpan,
        type: wrapReference(facade.type),
        deps: null,
        host: {
            ...extractHostBindings(facade.propMetadata, facade.typeSourceSpan, facade.host),
        },
        inputs: { ...inputsFromMetadata, ...inputsFromType },
        outputs: { ...outputsFromMetadata, ...outputsFromType },
        queries: facade.queries.map(convertToR3QueryMetadata),
        providers: facade.providers != null ? new WrappedNodeExpr(facade.providers) : null,
        viewQueries: facade.viewQueries.map(convertToR3QueryMetadata),
        fullInheritance: false,
        hostDirectives: convertHostDirectivesToMetadata(facade),
    };
}
function convertDeclareDirectiveFacadeToMetadata(declaration, typeSourceSpan) {
    return {
        name: declaration.type.name,
        type: wrapReference(declaration.type),
        typeSourceSpan,
        selector: declaration.selector ?? null,
        inputs: declaration.inputs ? inputsPartialMetadataToInputMetadata(declaration.inputs) : {},
        outputs: declaration.outputs ?? {},
        host: convertHostDeclarationToMetadata(declaration.host),
        queries: (declaration.queries ?? []).map(convertQueryDeclarationToMetadata),
        viewQueries: (declaration.viewQueries ?? []).map(convertQueryDeclarationToMetadata),
        providers: declaration.providers !== undefined ? new WrappedNodeExpr(declaration.providers) : null,
        exportAs: declaration.exportAs ?? null,
        usesInheritance: declaration.usesInheritance ?? false,
        lifecycle: { usesOnChanges: declaration.usesOnChanges ?? false },
        deps: null,
        typeArgumentCount: 0,
        fullInheritance: false,
        isStandalone: declaration.isStandalone ?? false,
        isSignal: declaration.isSignal ?? false,
        hostDirectives: convertHostDirectivesToMetadata(declaration),
    };
}
function convertHostDeclarationToMetadata(host = {}) {
    return {
        attributes: convertOpaqueValuesToExpressions(host.attributes ?? {}),
        listeners: host.listeners ?? {},
        properties: host.properties ?? {},
        specialAttributes: {
            classAttr: host.classAttribute,
            styleAttr: host.styleAttribute,
        },
    };
}
function convertHostDirectivesToMetadata(metadata) {
    if (metadata.hostDirectives?.length) {
        return metadata.hostDirectives.map((hostDirective) => {
            return typeof hostDirective === 'function'
                ? {
                    directive: wrapReference(hostDirective),
                    inputs: null,
                    outputs: null,
                    isForwardReference: false,
                }
                : {
                    directive: wrapReference(hostDirective.directive),
                    isForwardReference: false,
                    inputs: hostDirective.inputs ? parseMappingStringArray(hostDirective.inputs) : null,
                    outputs: hostDirective.outputs ? parseMappingStringArray(hostDirective.outputs) : null,
                };
        });
    }
    return null;
}
function convertOpaqueValuesToExpressions(obj) {
    const result = {};
    for (const key of Object.keys(obj)) {
        result[key] = new WrappedNodeExpr(obj[key]);
    }
    return result;
}
function convertDeclareComponentFacadeToMetadata(decl, typeSourceSpan, sourceMapUrl) {
    const { template, interpolation, defer } = parseJitTemplate(decl.template, decl.type.name, sourceMapUrl, decl.preserveWhitespaces ?? false, decl.interpolation, decl.deferBlockDependencies);
    const declarations = [];
    if (decl.dependencies) {
        for (const innerDep of decl.dependencies) {
            switch (innerDep.kind) {
                case 'directive':
                case 'component':
                    declarations.push(convertDirectiveDeclarationToMetadata(innerDep));
                    break;
                case 'pipe':
                    declarations.push(convertPipeDeclarationToMetadata(innerDep));
                    break;
            }
        }
    }
    else if (decl.components || decl.directives || decl.pipes) {
        // Existing declarations on NPM may not be using the new `dependencies` merged field, and may
        // have separate fields for dependencies instead. Unify them for JIT compilation.
        decl.components &&
            declarations.push(...decl.components.map((dir) => convertDirectiveDeclarationToMetadata(dir, /* isComponent */ true)));
        decl.directives &&
            declarations.push(...decl.directives.map((dir) => convertDirectiveDeclarationToMetadata(dir)));
        decl.pipes && declarations.push(...convertPipeMapToMetadata(decl.pipes));
    }
    return {
        ...convertDeclareDirectiveFacadeToMetadata(decl, typeSourceSpan),
        template,
        styles: decl.styles ?? [],
        declarations,
        viewProviders: decl.viewProviders !== undefined ? new WrappedNodeExpr(decl.viewProviders) : null,
        animations: decl.animations !== undefined ? new WrappedNodeExpr(decl.animations) : null,
        defer,
        changeDetection: decl.changeDetection ?? ChangeDetectionStrategy.Default,
        encapsulation: decl.encapsulation ?? ViewEncapsulation.Emulated,
        interpolation,
        declarationListEmitMode: 2 /* DeclarationListEmitMode.ClosureResolved */,
        relativeContextFilePath: '',
        i18nUseExternalIds: true,
    };
}
function convertDeclarationFacadeToMetadata(declaration) {
    return {
        ...declaration,
        type: new WrappedNodeExpr(declaration.type),
    };
}
function convertDirectiveDeclarationToMetadata(declaration, isComponent = null) {
    return {
        kind: R3TemplateDependencyKind.Directive,
        isComponent: isComponent || declaration.kind === 'component',
        selector: declaration.selector,
        type: new WrappedNodeExpr(declaration.type),
        inputs: declaration.inputs ?? [],
        outputs: declaration.outputs ?? [],
        exportAs: declaration.exportAs ?? null,
    };
}
function convertPipeMapToMetadata(pipes) {
    if (!pipes) {
        return [];
    }
    return Object.keys(pipes).map((name) => {
        return {
            kind: R3TemplateDependencyKind.Pipe,
            name,
            type: new WrappedNodeExpr(pipes[name]),
        };
    });
}
function convertPipeDeclarationToMetadata(pipe) {
    return {
        kind: R3TemplateDependencyKind.Pipe,
        name: pipe.name,
        type: new WrappedNodeExpr(pipe.type),
    };
}
function parseJitTemplate(template, typeName, sourceMapUrl, preserveWhitespaces, interpolation, deferBlockDependencies) {
    const interpolationConfig = interpolation
        ? InterpolationConfig.fromArray(interpolation)
        : DEFAULT_INTERPOLATION_CONFIG;
    // Parse the template and check for errors.
    const parsed = parseTemplate(template, sourceMapUrl, {
        preserveWhitespaces,
        interpolationConfig,
    });
    if (parsed.errors !== null) {
        const errors = parsed.errors.map((err) => err.toString()).join(', ');
        throw new Error(`Errors during JIT compilation of template for ${typeName}: ${errors}`);
    }
    const binder = new R3TargetBinder(new SelectorMatcher());
    const boundTarget = binder.bind({ template: parsed.nodes });
    return {
        template: parsed,
        interpolation: interpolationConfig,
        defer: createR3ComponentDeferMetadata(boundTarget, deferBlockDependencies),
    };
}
/**
 * Convert the expression, if present to an `R3ProviderExpression`.
 *
 * In JIT mode we do not want the compiler to wrap the expression in a `forwardRef()` call because,
 * if it is referencing a type that has not yet been defined, it will have already been wrapped in
 * a `forwardRef()` - either by the application developer or during partial-compilation. Thus we can
 * use `ForwardRefHandling.None`.
 */
function convertToProviderExpression(obj, property) {
    if (obj.hasOwnProperty(property)) {
        return createMayBeForwardRefExpression(new WrappedNodeExpr(obj[property]), 0 /* ForwardRefHandling.None */);
    }
    else {
        return undefined;
    }
}
function wrapExpression(obj, property) {
    if (obj.hasOwnProperty(property)) {
        return new WrappedNodeExpr(obj[property]);
    }
    else {
        return undefined;
    }
}
function computeProvidedIn(providedIn) {
    const expression = typeof providedIn === 'function'
        ? new WrappedNodeExpr(providedIn)
        : new LiteralExpr(providedIn ?? null);
    // See `convertToProviderExpression()` for why this uses `ForwardRefHandling.None`.
    return createMayBeForwardRefExpression(expression, 0 /* ForwardRefHandling.None */);
}
function convertR3DependencyMetadataArray(facades) {
    return facades == null ? null : facades.map(convertR3DependencyMetadata);
}
function convertR3DependencyMetadata(facade) {
    const isAttributeDep = facade.attribute != null; // both `null` and `undefined`
    const rawToken = facade.token === null ? null : new WrappedNodeExpr(facade.token);
    // In JIT mode, if the dep is an `@Attribute()` then we use the attribute name given in
    // `attribute` rather than the `token`.
    const token = isAttributeDep ? new WrappedNodeExpr(facade.attribute) : rawToken;
    return createR3DependencyMetadata(token, isAttributeDep, facade.host, facade.optional, facade.self, facade.skipSelf);
}
function convertR3DeclareDependencyMetadata(facade) {
    const isAttributeDep = facade.attribute ?? false;
    const token = facade.token === null ? null : new WrappedNodeExpr(facade.token);
    return createR3DependencyMetadata(token, isAttributeDep, facade.host ?? false, facade.optional ?? false, facade.self ?? false, facade.skipSelf ?? false);
}
function createR3DependencyMetadata(token, isAttributeDep, host, optional, self, skipSelf) {
    // If the dep is an `@Attribute()` the `attributeNameType` ought to be the `unknown` type.
    // But types are not available at runtime so we just use a literal `"<unknown>"` string as a dummy
    // marker.
    const attributeNameType = isAttributeDep ? literal('unknown') : null;
    return { token, attributeNameType, host, optional, self, skipSelf };
}
function createR3ComponentDeferMetadata(boundTarget, deferBlockDependencies) {
    const deferredBlocks = boundTarget.getDeferBlocks();
    const blocks = new Map();
    for (let i = 0; i < deferredBlocks.length; i++) {
        const dependencyFn = deferBlockDependencies?.[i];
        blocks.set(deferredBlocks[i], dependencyFn ? new WrappedNodeExpr(dependencyFn) : null);
    }
    return { mode: 0 /* DeferBlockDepsEmitMode.PerBlock */, blocks };
}
function extractHostBindings(propMetadata, sourceSpan, host) {
    // First parse the declarations from the metadata.
    const bindings = parseHostBindings(host || {});
    // After that check host bindings for errors
    const errors = verifyHostBindings(bindings, sourceSpan);
    if (errors.length) {
        throw new Error(errors.map((error) => error.msg).join('\n'));
    }
    // Next, loop over the properties of the object, looking for @HostBinding and @HostListener.
    for (const field in propMetadata) {
        if (propMetadata.hasOwnProperty(field)) {
            propMetadata[field].forEach((ann) => {
                if (isHostBinding(ann)) {
                    // Since this is a decorator, we know that the value is a class member. Always access it
                    // through `this` so that further down the line it can't be confused for a literal value
                    // (e.g. if there's a property called `true`).
                    bindings.properties[ann.hostPropertyName || field] = getSafePropertyAccessString('this', field);
                }
                else if (isHostListener(ann)) {
                    bindings.listeners[ann.eventName || field] = `${field}(${(ann.args || []).join(',')})`;
                }
            });
        }
    }
    return bindings;
}
function isHostBinding(value) {
    return value.ngMetadataName === 'HostBinding';
}
function isHostListener(value) {
    return value.ngMetadataName === 'HostListener';
}
function isInput(value) {
    return value.ngMetadataName === 'Input';
}
function isOutput(value) {
    return value.ngMetadataName === 'Output';
}
function inputsPartialMetadataToInputMetadata(inputs) {
    return Object.keys(inputs).reduce((result, minifiedClassName) => {
        const value = inputs[minifiedClassName];
        // Handle legacy partial input output.
        if (typeof value === 'string' || Array.isArray(value)) {
            result[minifiedClassName] = parseLegacyInputPartialOutput(value);
        }
        else {
            result[minifiedClassName] = {
                bindingPropertyName: value.publicName,
                classPropertyName: minifiedClassName,
                transformFunction: value.transformFunction !== null ? new WrappedNodeExpr(value.transformFunction) : null,
                required: value.isRequired,
                isSignal: value.isSignal,
            };
        }
        return result;
    }, {});
}
/**
 * Parses the legacy input partial output. For more details see `partial/directive.ts`.
 * TODO(legacy-partial-output-inputs): Remove in v18.
 */
function parseLegacyInputPartialOutput(value) {
    if (typeof value === 'string') {
        return {
            bindingPropertyName: value,
            classPropertyName: value,
            transformFunction: null,
            required: false,
            // legacy partial output does not capture signal inputs.
            isSignal: false,
        };
    }
    return {
        bindingPropertyName: value[0],
        classPropertyName: value[1],
        transformFunction: value[2] ? new WrappedNodeExpr(value[2]) : null,
        required: false,
        // legacy partial output does not capture signal inputs.
        isSignal: false,
    };
}
function parseInputsArray(values) {
    return values.reduce((results, value) => {
        if (typeof value === 'string') {
            const [bindingPropertyName, classPropertyName] = parseMappingString(value);
            results[classPropertyName] = {
                bindingPropertyName,
                classPropertyName,
                required: false,
                // Signal inputs not supported for the inputs array.
                isSignal: false,
                transformFunction: null,
            };
        }
        else {
            results[value.name] = {
                bindingPropertyName: value.alias || value.name,
                classPropertyName: value.name,
                required: value.required || false,
                // Signal inputs not supported for the inputs array.
                isSignal: false,
                transformFunction: value.transform != null ? new WrappedNodeExpr(value.transform) : null,
            };
        }
        return results;
    }, {});
}
function parseMappingStringArray(values) {
    return values.reduce((results, value) => {
        const [alias, fieldName] = parseMappingString(value);
        results[fieldName] = alias;
        return results;
    }, {});
}
function parseMappingString(value) {
    // Either the value is 'field' or 'field: property'. In the first case, `property` will
    // be undefined, in which case the field name should also be used as the property name.
    const [fieldName, bindingPropertyName] = value.split(':', 2).map((str) => str.trim());
    return [bindingPropertyName ?? fieldName, fieldName];
}
function convertDeclarePipeFacadeToMetadata(declaration) {
    return {
        name: declaration.type.name,
        type: wrapReference(declaration.type),
        typeArgumentCount: 0,
        pipeName: declaration.name,
        deps: null,
        pure: declaration.pure ?? true,
        isStandalone: declaration.isStandalone ?? false,
    };
}
function convertDeclareInjectorFacadeToMetadata(declaration) {
    return {
        name: declaration.type.name,
        type: wrapReference(declaration.type),
        providers: declaration.providers !== undefined && declaration.providers.length > 0
            ? new WrappedNodeExpr(declaration.providers)
            : null,
        imports: declaration.imports !== undefined
            ? declaration.imports.map((i) => new WrappedNodeExpr(i))
            : [],
    };
}
export function publishFacade(global) {
    const ng = global.ng || (global.ng = {});
    ng.ɵcompilerFacade = new CompilerFacadeImpl();
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaml0X2NvbXBpbGVyX2ZhY2FkZS5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uLy4uLy4uL3BhY2thZ2VzL2NvbXBpbGVyL3NyYy9qaXRfY29tcGlsZXJfZmFjYWRlLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBOzs7Ozs7R0FNRztBQThCSCxPQUFPLEVBQUMsWUFBWSxFQUFDLE1BQU0saUJBQWlCLENBQUM7QUFDN0MsT0FBTyxFQUNMLHVCQUF1QixFQUt2QixpQkFBaUIsR0FDbEIsTUFBTSxRQUFRLENBQUM7QUFDaEIsT0FBTyxFQUFDLGlCQUFpQixFQUFDLE1BQU0seUJBQXlCLENBQUM7QUFDMUQsT0FBTyxFQUFDLDRCQUE0QixFQUFFLG1CQUFtQixFQUFDLE1BQU0sc0JBQXNCLENBQUM7QUFDdkYsT0FBTyxFQUNMLGNBQWMsRUFFZCxPQUFPLEVBQ1AsV0FBVyxFQUVYLFlBQVksRUFDWixlQUFlLEdBQ2hCLE1BQU0scUJBQXFCLENBQUM7QUFDN0IsT0FBTyxFQUFDLFlBQVksRUFBQyxNQUFNLHFCQUFxQixDQUFDO0FBQ2pELE9BQU8sRUFBOEIsbUJBQW1CLEVBQUMsTUFBTSxjQUFjLENBQUM7QUFFOUUsT0FBTyxFQUFDLHNCQUFzQixFQUFFLGFBQWEsRUFBdUIsTUFBTSxzQkFBc0IsQ0FBQztBQUNqRyxPQUFPLEVBQUMsZUFBZSxFQUFxQixNQUFNLGdDQUFnQyxDQUFDO0FBQ25GLE9BQU8sRUFBQyxjQUFjLEVBQUMsTUFBTSxrQkFBa0IsQ0FBQztBQUNoRCxPQUFPLEVBQ0wsZUFBZSxFQUNmLG9DQUFvQyxFQUVwQyxzQkFBc0IsRUFDdEIsbUJBQW1CLEdBQ3BCLE1BQU0sOEJBQThCLENBQUM7QUFDdEMsT0FBTyxFQUFDLHVCQUF1QixFQUFpQixNQUFNLDRCQUE0QixDQUFDO0FBQ25GLE9BQU8sRUFDTCwrQkFBK0IsRUFFL0IsMkJBQTJCLEVBRTNCLGFBQWEsR0FDZCxNQUFNLGdCQUFnQixDQUFDO0FBQ3hCLE9BQU8s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RCxDQUFDLENBQUMsRUFBRTtLQUNULENBQUM7QUFDSixDQUFDO0FBRUQsTUFBTSxVQUFVLGFBQWEsQ0FBQyxNQUFXO0lBQ3ZDLE1BQU0sRUFBRSxHQUEyQixNQUFNLENBQUMsRUFBRSxJQUFJLENBQUMsTUFBTSxDQUFDLEVBQUUsR0FBRyxFQUFFLENBQUMsQ0FBQztJQUNqRSxFQUFFLENBQUMsZUFBZSxHQUFHLElBQUksa0JBQWtCLEVBQUUsQ0FBQztBQUNoRCxDQUFDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAbGljZW5zZVxuICogQ29weXJpZ2h0IEdvb2dsZSBMTEMgQWxsIFJpZ2h0cyBSZXNlcnZlZC5cbiAqXG4gKiBVc2Ugb2YgdGhpcyBzb3VyY2UgY29kZSBpcyBnb3Zlcm5lZCBieSBhbiBNSVQtc3R5bGUgbGljZW5zZSB0aGF0IGNhbiBiZVxuICogZm91bmQgaW4gdGhlIExJQ0VOU0UgZmlsZSBhdCBodHRwczovL2FuZ3VsYXIuaW8vbGljZW5zZVxuICovXG5cbmltcG9ydCB7XG4gIENvbXBpbGVyRmFjYWRlLFxuICBDb3JlRW52aXJvbm1lbnQsXG4gIEV4cG9ydGVkQ29tcGlsZXJGYWNhZGUsXG4gIExlZ2FjeUlucHV0UGFydGlhbE1hcHBpbmcsXG4gIE9wYXF1ZVZhbHVlLFxuICBSM0NvbXBvbmVudE1ldGFkYXRhRmFjYWRlLFxuICBSM0RlY2xhcmVDb21wb25lbnRGYWNhZGUsXG4gIFIzRGVjbGFyZURlcGVuZGVuY3lNZXRhZGF0YUZhY2FkZSxcbiAgUjNEZWNsYXJlRGlyZWN0aXZlRGVwZW5kZW5jeUZhY2FkZSxcbiAgUjNEZWNsYXJlRGlyZWN0aXZlRmFjYWRlLFxuICBSM0RlY2xhcmVGYWN0b3J5RmFjYWRlLFxuICBSM0RlY2xhcmVJbmplY3RhYmxlRmFjYWRlLFxuICBSM0RlY2xhcmVJbmplY3RvckZhY2FkZSxcbiAgUjNEZWNsYXJlTmdNb2R1bGVGYWNhZGUsXG4gIFIzRGVjbGFyZVBpcGVEZXBlbmRlbmN5RmFjYWRlLFxuICBSM0RlY2xhcmVQaXBlRmFjYWRlLFxuICBSM0RlY2xhcmVRdWVyeU1ldGFkYXRhRmFjYWRlLFxuICBSM0RlcGVuZGVuY3lNZXRhZGF0YUZhY2FkZSxcbiAgUjNEaXJlY3RpdmVNZXRhZGF0YUZhY2FkZSxcbiAgUjNGYWN0b3J5RGVmTWV0YWRhdGFGYWNhZGUsXG4gIFIzSW5qZWN0YWJsZU1ldGFkYXRhRmFjYWRlLFxuICBSM0luamVjdG9yTWV0YWRhdGFGYWNhZGUsXG4gIFIzTmdNb2R1bGVNZXRhZGF0YUZhY2FkZSxcbiAgUjNQaXBlTWV0YWRhdGFGYWNhZGUsXG4gIFIzUXVlcnlNZXRhZGF0YUZhY2FkZSxcbiAgUjNUZW1wbGF0ZURlcGVuZGVuY3lGYWNhZGUsXG59IGZyb20gJy4vY29tcGlsZXJfZmFjYWRlX2ludGVyZmFjZSc7XG5pbXBvcnQge0NvbnN0YW50UG9vbH0gZnJvbSAnLi9jb25zdGFudF9wb29sJztcbmltcG9ydCB7XG4gIENoYW5nZURldGVjdGlvblN0cmF0ZWd5LFxuICBIb3N0QmluZGluZyxcbiAgSG9zdExpc3RlbmVyLFxuICBJbnB1dCxcbiAgT3V0cHV0LFxuICBWaWV3RW5jYXBzdWxhdGlvbixcbn0gZnJvbSAnLi9jb3JlJztcbmltcG9ydCB7Y29tcGlsZUluamVjdGFibGV9IGZyb20gJy4vaW5qZWN0YWJsZV9jb21waWxlcl8yJztcbmltcG9ydCB7REVGQVVMVF9JTlRFUlBPTEFUSU9OX0NPTkZJRywgSW50ZXJwb2xhdGlvbkNvbmZpZ30gZnJvbSAnLi9tbF9wYXJzZXIvZGVmYXVsdHMnO1xuaW1wb3J0IHtcbiAgRGVjbGFyZVZhclN0bXQsXG4gIEV4cHJlc3Npb24sXG4gIGxpdGVyYWwsXG4gIExpdGVyYWxFeHByLFxuICBTdGF0ZW1lbnQsXG4gIFN0bXRNb2RpZmllcixcbiAgV3JhcHBlZE5vZGVFeHByLFxufSBmcm9tICcuL291dHB1dC9vdXRwdXRfYXN0JztcbmltcG9ydCB7Sml0RXZhbHVhdG9yfSBmcm9tICcuL291dHB1dC9vdXRwdXRfaml0JztcbmltcG9ydCB7UGFyc2VFcnJvciwgUGFyc2VTb3VyY2VTcGFuLCByM0ppdFR5cGVTb3VyY2VTcGFufSBmcm9tICcuL3BhcnNlX3V0aWwnO1xuaW1wb3J0IHtEZWZlcnJlZEJsb2NrfSBmcm9tICcuL3JlbmRlcjMvcjNfYXN0JztcbmltcG9ydCB7Y29tcGlsZUZhY3RvcnlGdW5jdGlvbiwgRmFjdG9yeVRhcmdldCwgUjNEZXBlbmRlbmN5TWV0YWRhdGF9IGZyb20gJy4vcmVuZGVyMy9yM19mYWN0b3J5JztcbmltcG9ydCB7Y29tcGlsZUluamVjdG9yLCBSM0luamVjdG9yTWV0YWRhdGF9IGZyb20gJy4vcmVuZGVyMy9yM19pbmplY3Rvcl9jb21waWxlcic7XG5pbXBvcnQge1IzSml0UmVmbGVjdG9yfSBmcm9tICcuL3JlbmRlcjMvcjNfaml0JztcbmltcG9ydCB7XG4gIGNvbXBpbGVOZ01vZHVsZSxcbiAgY29tcGlsZU5nTW9kdWxlRGVjbGFyYXRpb25FeHByZXNzaW9uLFxuICBSM05nTW9kdWxlTWV0YWRhdGEsXG4gIFIzTmdNb2R1bGVNZXRhZGF0YUtpbmQsXG4gIFIzU2VsZWN0b3JTY29wZU1vZGUsXG59IGZyb20gJy4vcmVuZGVyMy9yM19tb2R1bGVfY29tcGlsZXInO1xuaW1wb3J0IHtjb21waWxlUGlwZUZyb21NZXRhZGF0YSwgUjNQaXBlTWV0YWRhdGF9IGZyb20gJy4vcmVuZGVyMy9yM19waXBlX2NvbXBpbGVyJztcbmltcG9ydCB7XG4gIGNyZWF0ZU1heUJlRm9yd2FyZFJlZkV4cHJlc3Npb24sXG4gIEZvcndhcmRSZWZIYW5kbGluZyxcbiAgZ2V0U2FmZVByb3BlcnR5QWNjZXNzU3RyaW5nLFxuICBNYXliZUZvcndhcmRSZWZFeHByZXNzaW9uLFxuICB3cmFwUmVmZXJlbmNlLFxufSBmcm9tICcuL3JlbmRlcjMvdXRpbCc7XG5pbXBvcnQge1xuICBEZWNsYXJhdGlvbkxpc3RFbWl0TW9kZSxcbiAgRGVmZXJCbG9ja0RlcHNFbWl0TW9kZSxcbiAgUjNDb21wb25lbnREZWZlck1ldGFkYXRhLFxuICBSM0NvbXBvbmVudE1ldGFkYXRhLFxuICBSM0RpcmVjdGl2ZURlcGVuZGVuY3lNZXRhZGF0YSxcbiAgUjNEaXJlY3RpdmVNZXRhZGF0YSxcbiAgUjNIb3N0RGlyZWN0aXZlTWV0YWRhdGEsXG4gIFIzSG9zdE1ldGFkYXRhLFxuICBSM0lucHV0TWV0YWRhdGEsXG4gIFIzUGlwZURlcGVuZGVuY3lNZXRhZGF0YSxcbiAgUjNRdWVyeU1ldGFkYXRhLFxuICBSM1RlbXBsYXRlRGVwZW5kZW5jeSxcbiAgUjNUZW1wbGF0ZURlcGVuZGVuY3lLaW5kLFxuICBSM1RlbXBsYXRlRGVwZW5kZW5jeU1ldGFkYXRhLFxufSBmcm9tICcuL3JlbmRlcjMvdmlldy9hcGknO1xuaW1wb3J0IHtcbiAgY29tcGlsZUNvbXBvbmVudEZyb21NZXRhZGF0YSxcbiAgY29tcGlsZURpcmVjdGl2ZUZyb21NZXRhZGF0YSxcbiAgUGFyc2VkSG9zdEJpbmRpbmdzLFxuICBwYXJzZUhvc3RCaW5kaW5ncyxcbiAgdmVyaWZ5SG9zdEJpbmRpbmdzLFxufSBmcm9tICcuL3JlbmRlcjMvdmlldy9jb21waWxlcic7XG5cbmltcG9ydCB0eXBlIHtCb3VuZFRhcmdldH0gZnJvbSAnLi9yZW5kZXIzL3ZpZXcvdDJfYXBpJztcbmltcG9ydCB7UjNUYXJnZXRCaW5kZXJ9IGZyb20gJy4vcmVuZGVyMy92aWV3L3QyX2JpbmRlcic7XG5pbXBvcnQge21ha2VCaW5kaW5nUGFyc2VyLCBwYXJzZVRlbXBsYXRlfSBmcm9tICcuL3JlbmRlcjMvdmlldy90ZW1wbGF0ZSc7XG5pbXBvcnQge1Jlc291cmNlTG9hZGVyfSBmcm9tICcuL3Jlc291cmNlX2xvYWRlcic7XG5pbXBvcnQge0RvbUVsZW1lbnRTY2hlbWFSZWdpc3RyeX0gZnJvbSAnLi9zY2hlbWEvZG9tX2VsZW1lbnRfc2NoZW1hX3JlZ2lzdHJ5JztcbmltcG9ydCB7U2VsZWN0b3JNYXRjaGVyfSBmcm9tICcuL3NlbGVjdG9yJztcblxuZXhwb3J0IGNsYXNzIENvbXBpbGVyRmFjYWRlSW1wbCBpbXBsZW1lbnRzIENvbXBpbGVyRmFjYWRlIHtcbiAgRmFjdG9yeVRhcmdldCA9IEZhY3RvcnlUYXJnZXQ7XG4gIFJlc291cmNlTG9hZGVyID0gUmVzb3VyY2VMb2FkZXI7XG4gIHByaXZhdGUgZWxlbWVudFNjaGVtYVJlZ2lzdHJ5ID0gbmV3IERvbUVsZW1lbnRTY2hlbWFSZWdpc3RyeSgpO1xuXG4gIGNvbnN0cnVjdG9yKHByaXZhdGUgaml0RXZhbHVhdG9yID0gbmV3IEppdEV2YWx1YXRvcigpKSB7fVxuXG4gIGNvbXBpbGVQaXBlKFxuICAgIGFuZ3VsYXJDb3JlRW52OiBDb3JlRW52aXJvbm1lbnQsXG4gICAgc291cmNlTWFwVXJsOiBzdHJpbmcsXG4gICAgZmFjYWRlOiBSM1BpcGVNZXRhZGF0YUZhY2FkZSxcbiAgKTogYW55IHtcbiAgICBjb25zdCBtZXRhZGF0YTogUjNQaXBlTWV0YWRhdGEgPSB7XG4gICAgICBuYW1lOiBmYWNhZGUubmFtZSxcbiAgICAgIHR5cGU6IHdyYXBSZWZlcmVuY2UoZmFjYWRlLnR5cGUpLFxuICAgICAgdHlwZUFyZ3VtZW50Q291bnQ6IDAsXG4gICAgICBkZXBzOiBudWxsLFxuICAgICAgcGlwZU5hbWU6IGZhY2FkZS5waXBlTmFtZSxcbiAgICAgIHB1cmU6IGZhY2FkZS5wdXJlLFxuICAgICAgaXNTdGFuZGFsb25lOiBmYWNhZGUuaXNTdGFuZGFsb25lLFxuICAgIH07XG4gICAgY29uc3QgcmVzID0gY29tcGlsZVBpcGVGcm9tTWV0YWRhdGEobWV0YWRhdGEpO1xuICAgIHJldHVybiB0aGlzLmppdEV4cHJlc3Npb24ocmVzLmV4cHJlc3Npb24sIGFuZ3VsYXJDb3JlRW52LCBzb3VyY2VNYXBVcmwsIFtdKTtcbiAgfVxuXG4gIGNvbXBpbGVQaXBlRGVjbGFyYXRpb24oXG4gICAgYW5ndWxhckNvcmVFbnY6IENvcmVFbnZpcm9ubWVudCxcbiAgICBzb3VyY2VNYXBVcmw6IHN0cmluZyxcbiAgICBkZWNsYXJhdGlvbjogUjNEZWNsYXJlUGlwZUZhY2FkZSxcbiAgKTogYW55IHtcbiAgICBjb25zdCBtZXRhID0gY29udmVydERlY2xhcmVQaXBlRmFjYWRlVG9NZXRhZGF0YShkZWNsYXJhdGlvbik7XG4gICAgY29uc3QgcmVzID0gY29tcGlsZVBpcGVGcm9tTWV0YWRhdGEobWV0YSk7XG4gICAgcmV0dXJuIHRoaXMuaml0RXhwcmVzc2lvbihyZXMuZXhwcmVzc2lvbiwgYW5ndWxhckNvcmVFbnYsIHNvdXJjZU1hcFVybCwgW10pO1xuICB9XG5cbiAgY29tcGlsZUluamVjdGFibGUoXG4gICAgYW5ndWxhckNvcmVFbnY6IENvcmVFbnZpcm9ubWVudCxcbiAgICBzb3VyY2VNYXBVcmw6IHN0cmluZyxcbiAgICBmYWNhZGU6IFIzSW5qZWN0YWJsZU1ldGFkYXRhRmFjYWRlLFxuICApOiBhbnkge1xuICAgIGNvbnN0IHtleHByZXNzaW9uLCBzdGF0ZW1lbnRzfSA9IGNvbXBpbGVJbmplY3RhYmxlKFxuICAgICAge1xuICAgICAgICBuYW1lOiBmYWNhZGUubmFtZSxcbiAgICAgICAgdHlwZTogd3JhcFJlZmVyZW5jZShmYWNhZGUudHlwZSksXG4gICAgICAgIHR5cGVBcmd1bWVudENvdW50OiBmYWNhZGUudHlwZUFyZ3VtZW50Q291bnQsXG4gICAgICAgIHByb3ZpZGVkSW46IGNvbXB1dGVQcm92aWRlZEluKGZhY2FkZS5wcm92aWRlZEluKSxcbiAgICAgICAgdXNlQ2xhc3M6IGNvbnZlcnRUb1Byb3ZpZGVyRXhwcmVzc2lvbihmYWNhZGUsICd1c2VDbGFzcycpLFxuICAgICAgICB1c2VGYWN0b3J5OiB3cmFwRXhwcmVzc2lvbihmYWNhZGUsICd1c2VGYWN0b3J5JyksXG4gICAgICAgIHVzZVZhbHVlOiBjb252ZXJ0VG9Qcm92aWRlckV4cHJlc3Npb24oZmFjYWRlLCAndXNlVmFsdWUnKSxcbiAgICAgICAgdXNlRXhpc3Rpbmc6IGNvbnZlcnRUb1Byb3ZpZGVyRXhwcmVzc2lvbihmYWNhZGUsICd1c2VFeGlzdGluZycpLFxuICAgICAgICBkZXBzOiBmYWNhZGUuZGVwcz8ubWFwKGNvbnZlcnRSM0RlcGVuZGVuY3lNZXRhZGF0YSksXG4gICAgICB9LFxuICAgICAgLyogcmVzb2x2ZUZvcndhcmRSZWZzICovIHRydWUsXG4gICAgKTtcblxuICAgIHJldHVybiB0aGlzLmppdEV4cHJlc3Npb24oZXhwcmVzc2lvbiwgYW5ndWxhckNvcmVFbnYsIHNvdXJjZU1hcFVybCwgc3RhdGVtZW50cyk7XG4gIH1cblxuICBjb21waWxlSW5qZWN0YWJsZURlY2xhcmF0aW9uKFxuICAgIGFuZ3VsYXJDb3JlRW52OiBDb3JlRW52aXJvbm1lbnQsXG4gICAgc291cmNlTWFwVXJsOiBzdHJpbmcsXG4gICAgZmFjYWRlOiBSM0RlY2xhcmVJbmplY3RhYmxlRmFjYWRlLFxuICApOiBhbnkge1xuICAgIGNvbnN0IHtleHByZXNzaW9uLCBzdGF0ZW1lbnRzfSA9IGNvbXBpbGVJbmplY3RhYmxlKFxuICAgICAge1xuICAgICAgICBuYW1lOiBmYWNhZGUudHlwZS5uYW1lLFxuICAgICAgICB0eXBlOiB3cmFwUmVmZXJlbmNlKGZhY2FkZS50eXBlKSxcbiAgICAgICAgdHlwZUFyZ3VtZW50Q291bnQ6IDAsXG4gICAgICAgIHByb3ZpZGVkSW46IGNvbXB1dGVQcm92aWRlZEluKGZhY2FkZS5wcm92aWRlZEluKSxcbiAgICAgICAgdXNlQ2xhc3M6IGNvbnZlcnRUb1Byb3ZpZGVyRXhwcmVzc2lvbihmYWNhZGUsICd1c2VDbGFzcycpLFxuICAgICAgICB1c2VGYWN0b3J5OiB3cmFwRXhwcmVzc2lvbihmYWNhZGUsICd1c2VGYWN0b3J5JyksXG4gICAgICAgIHVzZVZhbHVlOiBjb252ZXJ0VG9Qcm92aWRlckV4cHJlc3Npb24oZmFjYWRlLCAndXNlVmFsdWUnKSxcbiAgICAgICAgdXNlRXhpc3Rpbmc6IGNvbnZlcnRUb1Byb3ZpZGVyRXhwcmVzc2lvbihmYWNhZGUsICd1c2VFeGlzdGluZycpLFxuICAgICAgICBkZXBzOiBmYWNhZGUuZGVwcz8ubWFwKGNvbnZlcnRSM0RlY2xhcmVEZXBlbmRlbmN5TWV0YWRhdGEpLFxuICAgICAgfSxcbiAgICAgIC8qIHJlc29sdmVGb3J3YXJkUmVmcyAqLyB0cnVlLFxuICAgICk7XG5cbiAgICByZXR1cm4gdGhpcy5qaXRFeHByZXNzaW9uKGV4cHJlc3Npb24sIGFuZ3VsYXJDb3JlRW52LCBzb3VyY2VNYXBVcmwsIHN0YXRlbWVudHMpO1xuICB9XG5cbiAgY29tcGlsZUluamVjdG9yKFxuICAgIGFuZ3VsYXJDb3JlRW52OiBDb3JlRW52aXJvbm1lbnQsXG4gICAgc291cmNlTWFwVXJsOiBzdHJpbmcsXG4gICAgZmFjYWRlOiBSM0luamVjdG9yTWV0YWRhdGFGYWNhZGUsXG4gICk6IGFueSB7XG4gICAgY29uc3QgbWV0YTogUjNJbmplY3Rvck1ldGFkYXRhID0ge1xuICAgICAgbmFtZTogZmFjYWRlLm5hbWUsXG4gICAgICB0eXBlOiB3cmFwUmVmZXJlbmNlKGZhY2FkZS50eXBlKSxcbiAgICAgIHByb3ZpZGVyczpcbiAgICAgICAgZmFjYWRlLnByb3ZpZGVycyAmJiBmYWNhZGUucHJvdmlkZXJzLmxlbmd0aCA+IDBcbiAgICAgICAgICA/IG5ldyBXcmFwcGVkTm9kZUV4cHIoZmFjYWRlLnByb3ZpZGVycylcbiAgICAgICAgICA6IG51bGwsXG4gICAgICBpbXBvcnRzOiBmYWNhZGUuaW1wb3J0cy5tYXAoKGkpID0+IG5ldyBXcmFwcGVkTm9kZUV4cHIoaSkpLFxuICAgIH07XG4gICAgY29uc3QgcmVzID0gY29tcGlsZUluamVjdG9yKG1ldGEpO1xuICAgIHJldHVybiB0aGlzLmppdEV4cHJlc3Npb24ocmVzLmV4cHJlc3Npb24sIGFuZ3VsYXJDb3JlRW52LCBzb3VyY2VNYXBVcmwsIFtdKTtcbiAgfVxuXG4gIGNvbXBpbGVJbmplY3RvckRlY2xhcmF0aW9uKFxuICAgIGFuZ3VsYXJDb3JlRW52OiBDb3JlRW52aXJvbm1lbnQsXG4gICAgc291cmNlTWFwVXJsOiBzdHJpbmcsXG4gICAgZGVjbGFyYXRpb246IFIzRGVjbGFyZUluamVjdG9yRmFjYWRlLFxuICApOiBhbnkge1xuICAgIGNvbnN0IG1ldGEgPSBjb252ZXJ0RGVjbGFyZUluamVjdG9yRmFjYWRlVG9NZXRhZGF0YShkZWNsYXJhdGlvbik7XG4gICAgY29uc3QgcmVzID0gY29tcGlsZUluamVjdG9yKG1ldGEpO1xuICAgIHJldHVybiB0aGlzLmppdEV4cHJlc3Npb24ocmVzLmV4cHJlc3Npb24sIGFuZ3VsYXJDb3JlRW52LCBzb3VyY2VNYXBVcmwsIFtdKTtcbiAgfVxuXG4gIGNvbXBpbGVOZ01vZHVsZShcbiAgICBhbmd1bGFyQ29yZUVudjogQ29yZUVudmlyb25tZW50LFxuICAgIHNvdXJjZU1hcFVybDogc3RyaW5nLFxuICAgIGZhY2FkZTogUjNOZ01vZHVsZU1ldGFkYXRhRmFjYWRlLFxuICApOiBhbnkge1xuICAgIGNvbnN0IG1ldGE6IFIzTmdNb2R1bGVNZXRhZGF0YSA9IHtcbiAgICAgIGtpbmQ6IFIzTmdNb2R1bGVNZXRhZGF0YUtpbmQuR2xvYmFsLFxuICAgICAgdHlwZTogd3JhcFJlZmVyZW5jZShmYWNhZGUudHlwZSksXG4gICAgICBib290c3RyYXA6IGZhY2FkZS5ib290c3RyYXAubWFwKHdyYXBSZWZlcmVuY2UpLFxuICAgICAgZGVjbGFyYXRpb25zOiBmYWNhZGUuZGVjbGFyYXRpb25zLm1hcCh3cmFwUmVmZXJlbmNlKSxcbiAgICAgIHB1YmxpY0RlY2xhcmF0aW9uVHlwZXM6IG51bGwsIC8vIG9ubHkgbmVlZGVkIGZvciB0eXBlcyBpbiBBT1RcbiAgICAgIGltcG9ydHM6IGZhY2FkZS5pbXBvcnRzLm1hcCh3cmFwUmVmZXJlbmNlKSxcbiAgICAgIGluY2x1ZGVJbXBvcnRUeXBlczogdHJ1ZSxcbiAgICAgIGV4cG9ydHM6IGZhY2FkZS5leHBvcnRzLm1hcCh3cmFwUmVmZXJlbmNlKSxcbiAgICAgIHNlbGVjdG9yU2NvcGVNb2RlOiBSM1NlbGVjdG9yU2NvcGVNb2RlLklubGluZSxcbiAgICAgIGNvbnRhaW5zRm9yd2FyZERlY2xzOiBmYWxzZSxcbiAgICAgIHNjaGVtYXM6IGZhY2FkZS5zY2hlbWFzID8gZmFjYWRlLnNjaGVtYXMubWFwKHdyYXBSZWZlcmVuY2UpIDogbnVsbCxcbiAgICAgIGlkOiBmYWNhZGUuaWQgPyBuZXcgV3JhcHBlZE5vZGVFeHByKGZhY2FkZS5pZCkgOiBudWxsLFxuICAgIH07XG4gICAgY29uc3QgcmVzID0gY29tcGlsZU5nTW9kdWxlKG1ldGEpO1xuICAgIHJldHVybiB0aGlzLmppdEV4cHJlc3Npb24ocmVzLmV4cHJlc3Npb24sIGFuZ3VsYXJDb3JlRW52LCBzb3VyY2VNYXBVcmwsIFtdKTtcbiAgfVxuXG4gIGNvbXBpbGVOZ01vZHVsZURlY2xhcmF0aW9uKFxuICAgIGFuZ3VsYXJDb3JlRW52OiBDb3JlRW52aXJvbm1lbnQsXG4gICAgc291cmNlTWFwVXJsOiBzdHJpbmcsXG4gICAgZGVjbGFyYXRpb246IFIzRGVjbGFyZU5nTW9kdWxlRmFjYWRlLFxuICApOiBhbnkge1xuICAgIGNvbnN0IGV4cHJlc3Npb24gPSBjb21waWxlTmdNb2R1bGVEZWNsYXJhdGlvbkV4cHJlc3Npb24oZGVjbGFyYXRpb24pO1xuICAgIHJldHVybiB0aGlzLmppdEV4cHJlc3Npb24oZXhwcmVzc2lvbiwgYW5ndWxhckNvcmVFbnYsIHNvdXJjZU1hcFVybCwgW10pO1xuICB9XG5cbiAgY29tcGlsZURpcmVjdGl2ZShcbiAgICBhbmd1bGFyQ29yZUVudjogQ29yZUVudmlyb25tZW50LFxuICAgIHNvdXJjZU1hcFVybDogc3RyaW5nLFxuICAgIGZhY2FkZTogUjNEaXJlY3RpdmVNZXRhZGF0YUZhY2FkZSxcbiAgKTogYW55IHtcbiAgICBjb25zdCBtZXRhOiBSM0RpcmVjdGl2ZU1ldGFkYXRhID0gY29udmVydERpcmVjdGl2ZUZhY2FkZVRvTWV0YWRhdGEoZmFjYWRlKTtcbiAgICByZXR1cm4gdGhpcy5jb21waWxlRGlyZWN0aXZlRnJvbU1ldGEoYW5ndWxhckNvcmVFbnYsIHNvdXJjZU1hcFVybCwgbWV0YSk7XG4gIH1cblxuICBjb21waWxlRGlyZWN0aXZlRGVjbGFyYXRpb24oXG4gICAgYW5ndWxhckNvcmVFbnY6IENvcmVFbnZpcm9ubWVudCxcbiAgICBzb3VyY2VNYXBVcmw6IHN0cmluZyxcbiAgICBkZWNsYXJhdGlvbjogUjNEZWNsYXJlRGlyZWN0aXZlRmFjYWRlLFxuICApOiBhbnkge1xuICAgIGNvbnN0IHR5cGVTb3VyY2VTcGFuID0gdGhpcy5jcmVhdGVQYXJzZVNvdXJjZVNwYW4oXG4gICAgICAnRGlyZWN0aXZlJyxcbiAgICAgIGRlY2xhcmF0aW9uLnR5cGUubmFtZSxcbiAgICAgIHNvdXJjZU1hcFVybCxcbiAgICApO1xuICAgIGNvbnN0IG1ldGEgPSBjb252ZXJ0RGVjbGFyZURpcmVjdGl2ZUZhY2FkZVRvTWV0YWRhdGEoZGVjbGFyYXRpb24sIHR5cGVTb3VyY2VTcGFuKTtcbiAgICByZXR1cm4gdGhpcy5jb21waWxlRGlyZWN0aXZlRnJvbU1ldGEoYW5ndWxhckNvcmVFbnYsIHNvdXJjZU1hcFVybCwgbWV0YSk7XG4gIH1cblxuICBwcml2YXRlIGNvbXBpbGVEaXJlY3RpdmVGcm9tTWV0YShcbiAgICBhbmd1bGFyQ29yZUVudjogQ29yZUVudmlyb25tZW50LFxuICAgIHNvdXJjZU1hcFVybDogc3RyaW5nLFxuICAgIG1ldGE6IFIzRGlyZWN0aXZlTWV0YWRhdGEsXG4gICk6IGFueSB7XG4gICAgY29uc3QgY29uc3RhbnRQb29sID0gbmV3IENvbnN0YW50UG9vbCgpO1xuICAgIGNvbnN0IGJpbmRpbmdQYXJzZXIgPSBtYWtlQmluZGluZ1BhcnNlcigpO1xuICAgIGNvbnN0IHJlcyA9IGNvbXBpbGVEaXJlY3RpdmVGcm9tTWV0YWRhdGEobWV0YSwgY29uc3RhbnRQb29sLCBiaW5kaW5nUGFyc2VyKTtcbiAgICByZXR1cm4gdGhpcy5qaXRFeHByZXNzaW9uKFxuICAgICAgcmVzLmV4cHJlc3Npb24sXG4gICAgICBhbmd1bGFyQ29yZUVudixcbiAgICAgIHNvdXJjZU1hcFVybCxcbiAgICAgIGNvbnN0YW50UG9vbC5zdGF0ZW1lbnRzLFxuICAgICk7XG4gIH1cblxuICBjb21waWxlQ29tcG9uZW50KFxuICAgIGFuZ3VsYXJDb3JlRW52OiBDb3JlRW52aXJvbm1lbnQsXG4gICAgc291cmNlTWFwVXJsOiBzdHJpbmcsXG4gICAgZmFjYWRlOiBSM0NvbXBvbmVudE1ldGFkYXRhRmFjYWRlLFxuICApOiBhbnkge1xuICAgIC8vIFBhcnNlIHRoZSB0ZW1wbGF0ZSBhbmQgY2hlY2sgZm9yIGVycm9ycy5cbiAgICBjb25zdCB7dGVtcGxhdGUsIGludGVycG9sYXRpb24sIGRlZmVyfSA9IHBhcnNlSml0VGVtcGxhdGUoXG4gICAgICBmYWNhZGUudGVtcGxhdGUsXG4gICAgICBmYWNhZGUubmFtZSxcbiAgICAgIHNvdXJjZU1hcFVybCxcbiAgICAgIGZhY2FkZS5wcmVzZXJ2ZVdoaXRlc3BhY2VzLFxuICAgICAgZmFjYWRlLmludGVycG9sYXRpb24sXG4gICAgICB1bmRlZmluZWQsXG4gICAgKTtcblxuICAgIC8vIENvbXBpbGUgdGhlIGNvbXBvbmVudCBtZXRhZGF0YSwgaW5jbHVkaW5nIHRlbXBsYXRlLCBpbnRvIGFuIGV4cHJlc3Npb24uXG4gICAgY29uc3QgbWV0YTogUjNDb21wb25lbnRNZXRhZGF0YTxSM1RlbXBsYXRlRGVwZW5kZW5jeT4gPSB7XG4gICAgICAuLi5mYWNhZGUsXG4gICAgICAuLi5jb252ZXJ0RGlyZWN0aXZlRmFjYWRlVG9NZXRhZGF0YShmYWNhZGUpLFxuICAgICAgc2VsZWN0b3I6IGZhY2FkZS5zZWxlY3RvciB8fCB0aGlzLmVsZW1lbnRTY2hlbWFSZWdpc3RyeS5nZXREZWZhdWx0Q29tcG9uZW50RWxlbWVudE5hbWUoKSxcbiAgICAgIHRlbXBsYXRlLFxuICAgICAgZGVjbGFyYXRpb25zOiBmYWNhZGUuZGVjbGFyYXRpb25zLm1hcChjb252ZXJ0RGVjbGFyYXRpb25GYWNhZGVUb01ldGFkYXRhKSxcbiAgICAgIGRlY2xhcmF0aW9uTGlzdEVtaXRNb2RlOiBEZWNsYXJhdGlvbkxpc3RFbWl0TW9kZS5EaXJlY3QsXG4gICAgICBkZWZlcixcblxuICAgICAgc3R5bGVzOiBbLi4uZmFjYWRlLnN0eWxlcywgLi4udGVtcGxhdGUuc3R5bGVzXSxcbiAgICAgIGVuY2Fwc3VsYXRpb246IGZhY2FkZS5lbmNhcHN1bGF0aW9uLFxuICAgICAgaW50ZXJwb2xhdGlvbixcbiAgICAgIGNoYW5nZURldGVjdGlvbjogZmFjYWRlLmNoYW5nZURldGVjdGlvbiA/PyBudWxsLFxuICAgICAgYW5pbWF0aW9uczogZmFjYWRlLmFuaW1hdGlvbnMgIT0gbnVsbCA/IG5ldyBXcmFwcGVkTm9kZUV4cHIoZmFjYWRlLmFuaW1hdGlvbnMpIDogbnVsbCxcbiAgICAgIHZpZXdQcm92aWRlcnM6XG4gICAgICAgIGZhY2FkZS52aWV3UHJvdmlkZXJzICE9IG51bGwgPyBuZXcgV3JhcHBlZE5vZGVFeHByKGZhY2FkZS52aWV3UHJvdmlkZXJzKSA6IG51bGwsXG4gICAgICByZWxhdGl2ZUNvbnRleHRGaWxlUGF0aDogJycsXG4gICAgICBpMThuVXNlRXh0ZXJuYWxJZHM6IHRydWUsXG4gICAgfTtcbiAgICBjb25zdCBqaXRFeHByZXNzaW9uU291cmNlTWFwID0gYG5nOi8vLyR7ZmFjYWRlLm5hbWV9LmpzYDtcbiAgICByZXR1cm4gdGhpcy5jb21waWxlQ29tcG9uZW50RnJvbU1ldGEoYW5ndWxhckNvcmVFbnYsIGppdEV4cHJlc3Npb25Tb3VyY2VNYXAsIG1ldGEpO1xuICB9XG5cbiAgY29tcGlsZUNvbXBvbmVudERlY2xhcmF0aW9uKFxuICAgIGFuZ3VsYXJDb3JlRW52OiBDb3JlRW52aXJvbm1lbnQsXG4gICAgc291cmNlTWFwVXJsOiBzdHJpbmcsXG4gICAgZGVjbGFyYXRpb246IFIzRGVjbGFyZUNvbXBvbmVudEZhY2FkZSxcbiAgKTogYW55IHtcbiAgICBjb25zdCB0eXBlU291cmNlU3BhbiA9IHRoaXMuY3JlYXRlUGFyc2VTb3VyY2VTcGFuKFxuICAgICAgJ0NvbXBvbmVudCcsXG4gICAgICBkZWNsYXJhdGlvbi50eXBlLm5hbWUsXG4gICAgICBzb3VyY2VNYXBVcmwsXG4gICAgKTtcbiAgICBjb25zdCBtZXRhID0gY29udmVydERlY2xhcmVDb21wb25lbnRGYWNhZGVUb01ldGFkYXRhKGRlY2xhcmF0aW9uLCB0eXBlU291cmNlU3Bhbiwgc291cmNlTWFwVXJsKTtcbiAgICByZXR1cm4gdGhpcy5jb21waWxlQ29tcG9uZW50RnJvbU1ldGEoYW5ndWxhckNvcmVFbnYsIHNvdXJjZU1hcFVybCwgbWV0YSk7XG4gIH1cblxuICBwcml2YXRlIGNvbXBpbGVDb21wb25lbnRGcm9tTWV0YShcbiAgICBhbmd1bGFyQ29yZUVudjogQ29yZUVudmlyb25tZW50LFxuICAgIHNvdXJjZU1hcFVybDogc3RyaW5nLFxuICAgIG1ldGE6IFIzQ29tcG9uZW50TWV0YWRhdGE8UjNUZW1wbGF0ZURlcGVuZGVuY3k+LFxuICApOiBhbnkge1xuICAgIGNvbnN0IGNvbnN0YW50UG9vbCA9IG5ldyBDb25zdGFudFBvb2woKTtcbiAgICBjb25zdCBiaW5kaW5nUGFyc2VyID0gbWFrZUJpbmRpbmdQYXJzZXIobWV0YS5pbnRlcnBvbGF0aW9uKTtcbiAgICBjb25zdCByZXMgPSBjb21waWxlQ29tcG9uZW50RnJvbU1ldGFkYXRhKG1ldGEsIGNvbnN0YW50UG9vbCwgYmluZGluZ1BhcnNlcik7XG4gICAgcmV0dXJuIHRoaXMuaml0RXhwcmVzc2lvbihcbiAgICAgIHJlcy5leHByZXNzaW9uLFxuICAgICAgYW5ndWxhckNvcmVFbnYsXG4gICAgICBzb3VyY2VNYXBVcmwsXG4gICAgICBjb25zdGFudFBvb2wuc3RhdGVtZW50cyxcbiAgICApO1xuICB9XG5cbiAgY29tcGlsZUZhY3RvcnkoXG4gICAgYW5ndWxhckNvcmVFbnY6IENvcmVFbnZpcm9ubWVudCxcbiAgICBzb3VyY2VNYXBVcmw6IHN0cmluZyxcbiAgICBtZXRhOiBSM0ZhY3RvcnlEZWZNZXRhZGF0YUZhY2FkZSxcbiAgKSB7XG4gICAgY29uc3QgZmFjdG9yeVJlcyA9IGNvbXBpbGVGYWN0b3J5RnVuY3Rpb24oe1xuICAgICAgbmFtZTogbWV0YS5uYW1lLFxuICAgICAgdHlwZTogd3JhcFJlZmVyZW5jZShtZXRhLnR5cGUpLFxuICAgICAgdHlwZUFyZ3VtZW50Q291bnQ6IG1ldGEudHlwZUFyZ3VtZW50Q291bnQsXG4gICAgICBkZXBzOiBjb252ZXJ0UjNEZXBlbmRlbmN5TWV0YWRhdGFBcnJheShtZXRhLmRlcHMpLFxuICAgICAgdGFyZ2V0OiBtZXRhLnRhcmdldCxcbiAgICB9KTtcbiAgICByZXR1cm4gdGhpcy5qaXRFeHByZXNzaW9uKFxuICAgICAgZmFjdG9yeVJlcy5leHByZXNzaW9uLFxuICAgICAgYW5ndWxhckNvcmVFbnYsXG4gICAgICBzb3VyY2VNYXBVcmwsXG4gICAgICBmYWN0b3J5UmVzLnN0YXRlbWVudHMsXG4gICAgKTtcbiAgfVxuXG4gIGNvbXBpbGVGYWN0b3J5RGVjbGFyYXRpb24oXG4gICAgYW5ndWxhckNvcmVFbnY6IENvcmVFbnZpcm9ubWVudCxcbiAgICBzb3VyY2VNYXBVcmw6IHN0cmluZyxcbiAgICBtZXRhOiBSM0RlY2xhcmVGYWN0b3J5RmFjYWRlLFxuICApIHtcbiAgICBjb25zdCBmYWN0b3J5UmVzID0gY29tcGlsZUZhY3RvcnlGdW5jdGlvbih7XG4gICAgICBuYW1lOiBtZXRhLnR5cGUubmFtZSxcbiAgICAgIHR5cGU6IHdyYXBSZWZlcmVuY2UobWV0YS50eXBlKSxcbiAgICAgIHR5cGVBcmd1bWVudENvdW50OiAwLFxuICAgICAgZGVwczogQXJyYXkuaXNBcnJheShtZXRhLmRlcHMpXG4gICAgICAgID8gbWV0YS5kZXBzLm1hcChjb252ZXJ0UjNEZWNsYXJlRGVwZW5kZW5jeU1ldGFkYXRhKVxuICAgICAgICA6IG1ldGEuZGVwcyxcbiAgICAgIHRhcmdldDogbWV0YS50YXJnZXQsXG4gICAgfSk7XG4gICAgcmV0dXJuIHRoaXMuaml0RXhwcmVzc2lvbihcbiAgICAgIGZhY3RvcnlSZXMuZXhwcmVzc2lvbixcbiAgICAgIGFuZ3VsYXJDb3JlRW52LFxuICAgICAgc291cmNlTWFwVXJsLFxuICAgICAgZmFjdG9yeVJlcy5zdGF0ZW1lbnRzLFxuICAgICk7XG4gIH1cblxuICBjcmVhdGVQYXJzZVNvdXJjZVNwYW4oa2luZDogc3RyaW5nLCB0eXBlTmFtZTogc3RyaW5nLCBzb3VyY2VVcmw6IHN0cmluZyk6IFBhcnNlU291cmNlU3BhbiB7XG4gICAgcmV0dXJuIHIzSml0VHlwZVNvdXJjZVNwYW4oa2luZCwgdHlwZU5hbWUsIHNvdXJjZVVybCk7XG4gIH1cblxuICAvKipcbiAgICogSklUIGNvbXBpbGVzIGFuIGV4cHJlc3Npb24gYW5kIHJldHVybnMgdGhlIHJlc3VsdCBvZiBleGVjdXRpbmcgdGhhdCBleHByZXNzaW9uLlxuICAgKlxuICAgKiBAcGFyYW0gZGVmIHRoZSBkZWZpbml0aW9uIHdoaWNoIHdpbGwgYmUgY29tcGlsZWQgYW5kIGV4ZWN1dGVkIHRvIGdldCB0aGUgdmFsdWUgdG8gcGF0Y2hcbiAgICogQHBhcmFtIGNvbnRleHQgYW4gb2JqZWN0IG1hcCBvZiBAYW5ndWxhci9jb3JlIHN5bWJvbCBuYW1lcyB0byBzeW1ib2xzIHdoaWNoIHdpbGwgYmUgYXZhaWxhYmxlXG4gICAqIGluIHRoZSBjb250ZXh0IG9mIHRoZSBjb21waWxlZCBleHByZXNzaW9uXG4gICAqIEBwYXJhbSBzb3VyY2VVcmwgYSBVUkwgdG8gdXNlIGZvciB0aGUgc291cmNlIG1hcCBvZiB0aGUgY29tcGlsZWQgZXhwcmVzc2lvblxuICAgKiBAcGFyYW0gcHJlU3RhdGVtZW50cyBhIGNvbGxlY3Rpb24gb2Ygc3RhdGVtZW50cyB0aGF0IHNob3VsZCBiZSBldmFsdWF0ZWQgYmVmb3JlIHRoZSBleHByZXNzaW9uLlxuICAgKi9cbiAgcHJpdmF0ZSBqaXRFeHByZXNzaW9uKFxuICAgIGRlZjogRXhwcmVzc2lvbixcbiAgICBjb250ZXh0OiB7W2tleTogc3RyaW5nXTogYW55fSxcbiAgICBzb3VyY2VVcmw6IHN0cmluZyxcbiAgICBwcmVTdGF0ZW1lbnRzOiBTdGF0ZW1lbnRbXSxcbiAgKTogYW55IHtcbiAgICAvLyBUaGUgQ29uc3RhbnRQb29sIG1heSBjb250YWluIFN0YXRlbWVudHMgd2hpY2ggZGVjbGFyZSB2YXJpYWJsZXMgdXNlZCBpbiB0aGUgZmluYWwgZXhwcmVzc2lvbi5cbiAgICAvLyBUaGVyZWZvcmUsIGl0cyBzdGF0ZW1lbnRzIG5lZWQgdG8gcHJlY2VkZSB0aGUgYWN0dWFsIEpJVCBvcGVyYXRpb24uIFRoZSBmaW5hbCBzdGF0ZW1lbnQgaXMgYVxuICAgIC8vIGRlY2xhcmF0aW9uIG9mICRkZWYgd2hpY2ggaXMgc2V0IHRvIHRoZSBleHByZXNzaW9uIGJlaW5nIGNvbXBpbGVkLlxuICAgIGNvbnN0IHN0YXRlbWVudHM6IFN0YXRlbWVudFtdID0gW1xuICAgICAgLi4ucHJlU3RhdGVtZW50cyxcbiAgICAgIG5ldyBEZWNsYXJlVmFyU3RtdCgnJGRlZicsIGRlZiwgdW5kZWZpbmVkLCBTdG10TW9kaWZpZXIuRXhwb3J0ZWQpLFxuICAgIF07XG5cbiAgICBjb25zdCByZXMgPSB0aGlzLmppdEV2YWx1YXRvci5ldmFsdWF0ZVN0YXRlbWVudHMoXG4gICAgICBzb3VyY2VVcmwsXG4gICAgICBzdGF0ZW1lbnRzLFxuICAgICAgbmV3IFIzSml0UmVmbGVjdG9yKGNvbnRleHQpLFxuICAgICAgLyogZW5hYmxlU291cmNlTWFwcyAqLyB0cnVlLFxuICAgICk7XG4gICAgcmV0dXJuIHJlc1snJGRlZiddO1xuICB9XG59XG5cbmZ1bmN0aW9uIGNvbnZlcnRUb1IzUXVlcnlNZXRhZGF0YShmYWNhZGU6IFIzUXVlcnlNZXRhZGF0YUZhY2FkZSk6IFIzUXVlcnlNZXRhZGF0YSB7XG4gIHJldHVybiB7XG4gICAgLi4uZmFjYWRlLFxuICAgIGlzU2lnbmFsOiBmYWNhZGUuaXNTaWduYWwsXG4gICAgcHJlZGljYXRlOiBjb252ZXJ0UXVlcnlQcmVkaWNhdGUoZmFjYWRlLnByZWRpY2F0ZSksXG4gICAgcmVhZDogZmFjYWRlLnJlYWQgPyBuZXcgV3JhcHBlZE5vZGVFeHByKGZhY2FkZS5yZWFkKSA6IG51bGwsXG4gICAgc3RhdGljOiBmYWNhZGUuc3RhdGljLFxuICAgIGVtaXREaXN0aW5jdENoYW5nZXNPbmx5OiBmYWNhZGUuZW1pdERpc3RpbmN0Q2hhbmdlc09ubHksXG4gIH07XG59XG5cbmZ1bmN0aW9uIGNvbnZlcnRRdWVyeURlY2xhcmF0aW9uVG9NZXRhZGF0YShcbiAgZGVjbGFyYXRpb246IFIzRGVjbGFyZVF1ZXJ5TWV0YWRhdGFGYWNhZGUsXG4pOiBSM1F1ZXJ5TWV0YWRhdGEge1xuICByZXR1cm4ge1xuICAgIHByb3BlcnR5TmFtZTogZGVjbGFyYXRpb24ucHJvcGVydHlOYW1lLFxuICAgIGZpcnN0OiBkZWNsYXJhdGlvbi5maXJzdCA/PyBmYWxzZSxcbiAgICBwcmVkaWNhdGU6IGNvbnZlcnRRdWVyeVByZWRpY2F0ZShkZWNsYXJhdGlvbi5wcmVkaWNhdGUpLFxuICAgIGRlc2NlbmRhbnRzOiBkZWNsYXJhdGlvbi5kZXNjZW5kYW50cyA/PyBmYWxzZSxcbiAgICByZWFkOiBkZWNsYXJhdGlvbi5yZWFkID8gbmV3IFdyYXBwZWROb2RlRXhwcihkZWNsYXJhdGlvbi5yZWFkKSA6IG51bGwsXG4gICAgc3RhdGljOiBkZWNsYXJhdGlvbi5zdGF0aWMgPz8gZmFsc2UsXG4gICAgZW1pdERpc3RpbmN0Q2hhbmdlc09ubHk6IGRlY2xhcmF0aW9uLmVtaXREaXN0aW5jdENoYW5nZXNPbmx5ID8/IHRydWUsXG4gICAgaXNTaWduYWw6ICEhZGVjbGFyYXRpb24uaXNTaWduYWwsXG4gIH07XG59XG5cbmZ1bmN0aW9uIGNvbnZlcnRRdWVyeVByZWRpY2F0ZShcbiAgcHJlZGljYXRlOiBPcGFxdWVWYWx1ZSB8IHN0cmluZ1tdLFxuKTogTWF5YmVGb3J3YXJkUmVmRXhwcmVzc2lvbiB8IHN0cmluZ1tdIHtcbiAgcmV0dXJuIEFycmF5LmlzQXJyYXkocHJlZGljYXRlKVxuICAgID8gLy8gVGhlIHByZWRpY2F0ZSBpcyBhbiBhcnJheSBvZiBzdHJpbmdzIHNvIHBhc3MgaXQgdGhyb3VnaC5cbiAgICAgIHByZWRpY2F0ZVxuICAgIDogLy8gVGhlIHByZWRpY2F0ZSBpcyBhIHR5cGUgLSBhc3N1bWUgdGhhdCB3ZSB3aWxsIG5lZWQgdG8gdW53cmFwIGFueSBgZm9yd2FyZFJlZigpYCBjYWxscy5cbiAgICAgIGNyZWF0ZU1heUJlRm9yd2FyZFJlZkV4cHJlc3Npb24obmV3IFdyYXBwZWROb2RlRXhwcihwcmVkaWNhdGUpLCBGb3J3YXJkUmVmSGFuZGxpbmcuV3JhcHBlZCk7XG59XG5cbmZ1bmN0aW9uIGNvbnZlcnREaXJlY3RpdmVGYWNhZGVUb01ldGFkYXRhKGZhY2FkZTogUjNEaXJlY3RpdmVNZXRhZGF0YUZhY2FkZSk6IFIzRGlyZWN0aXZlTWV0YWRhdGEge1xuICBjb25zdCBpbnB1dHNGcm9tTWV0YWRhdGEgPSBwYXJzZUlucHV0c0FycmF5KGZhY2FkZS5pbnB1dHMgfHwgW10pO1xuICBjb25zdCBvdXRwdXRzRnJvbU1ldGFkYXRhID0gcGFyc2VNYXBwaW5nU3RyaW5nQXJyYXkoZmFjYWRlLm91dHB1dHMgfHwgW10pO1xuICBjb25zdCBwcm9wTWV0YWRhdGEgPSBmYWNhZGUucHJvcE1ldGFkYXRhO1xuICBjb25zdCBpbnB1dHNGcm9tVHlwZTogUmVjb3JkPHN0cmluZywgUjNJbnB1dE1ldGFkYXRhPiA9IHt9O1xuICBjb25zdCBvdXRwdXRzRnJvbVR5cGU6IFJlY29yZDxzdHJpbmcsIHN0cmluZz4gPSB7fTtcbiAgZm9yIChjb25zdCBmaWVsZCBpbiBwcm9wTWV0YWRhdGEpIHtcbiAgICBpZiAocHJvcE1ldGFkYXRhLmhhc093blByb3BlcnR5KGZpZWxkKSkge1xuICAgICAgcHJvcE1ldGFkYXRhW2ZpZWxkXS5mb3JFYWNoKChhbm4pID0+IHtcbiAgICAgICAgaWYgKGlzSW5wdXQoYW5uKSkge1xuICAgICAgICAgIGlucHV0c0Zyb21UeXBlW2ZpZWxkXSA9IHtcbiAgICAgICAgICAgIGJpbmRpbmdQcm9wZXJ0eU5hbWU6IGFubi5hbGlhcyB8fCBmaWVsZCxcbiAgICAgICAgICAgIGNsYXNzUHJvcGVydHlOYW1lOiBmaWVsZCxcbiAgICAgICAgICAgIHJlcXVpcmVkOiBhbm4ucmVxdWlyZWQgfHwgZmFsc2UsXG4gICAgICAgICAgICAvLyBGb3IgSklULCBkZWNvcmF0b3JzIGFyZSB1c2VkIHRvIGRlY2xhcmUgc2lnbmFsIGlucHV0cy4gVGhhdCBpcyBiZWNhdXNlIG9mXG4gICAgICAgICAgICAvLyBhIHRlY2huaWNhbCBsaW1pdGF0aW9uIHdoZXJlIGl0J3Mgbm90IHBvc3NpYmxlIHRvIHN0YXRpY2FsbHkgcmVmbGVjdCBjbGFzc1xuICAgICAgICAgICAgLy8gbWVtYmVycyBvZiBhIGRpcmVjdGl2ZS9jb21wb25lbnQgYXQgcnVudGltZSBiZWZvcmUgaW5zdGFudGlhdGluZyB0aGUgY2xhc3MuXG4gICAgICAgICAgICBpc1NpZ25hbDogISFhbm4uaXNTaWduYWwsXG4gICAgICAgICAgICB0cmFuc2Zvcm1GdW5jdGlvbjogYW5uLnRyYW5zZm9ybSAhPSBudWxsID8gbmV3IFdyYXBwZWROb2RlRXhwcihhbm4udHJhbnNmb3JtKSA6IG51bGwsXG4gICAgICAgICAgfTtcbiAgICAgICAgfSBlbHNlIGlmIChpc091dHB1dChhbm4pKSB7XG4gICAgICAgICAgb3V0cHV0c0Zyb21UeXBlW2ZpZWxkXSA9IGFubi5hbGlhcyB8fCBmaWVsZDtcbiAgICAgICAgfVxuICAgICAgfSk7XG4gICAgfVxuICB9XG5cbiAgcmV0dXJuIHtcbiAgICAuLi5mYWNhZGUsXG4gICAgdHlwZUFyZ3VtZW50Q291bnQ6IDAsXG4gICAgdHlwZVNvdXJjZVNwYW46IGZhY2FkZS50eXBlU291cmNlU3BhbixcbiAgICB0eXBlOiB3cmFwUmVmZXJlbmNlKGZhY2FkZS50eXBlKSxcbiAgICBkZXBzOiBudWxsLFxuICAgIGhvc3Q6IHtcbiAgICAgIC4uLmV4dHJhY3RIb3N0QmluZGluZ3MoZmFjYWRlLnByb3BNZXRhZGF0YSwgZmFjYWRlLnR5cGVTb3VyY2VTcGFuLCBmYWNhZGUuaG9zdCksXG4gICAgfSxcbiAgICBpbnB1dHM6IHsuLi5pbnB1dHNGcm9tTWV0YWRhdGEsIC4uLmlucHV0c0Zyb21UeXBlfSxcbiAgICBvdXRwdXRzOiB7Li4ub3V0cHV0c0Zyb21NZXRhZGF0YSwgLi4ub3V0cHV0c0Zyb21UeXBlfSxcbiAgICBxdWVyaWVzOiBmYWNhZGUucXVlcmllcy5tYXAoY29udmVydFRvUjNRdWVyeU1ldGFkYXRhKSxcbiAgICBwcm92aWRlcnM6IGZhY2FkZS5wcm92aWRlcnMgIT0gbnVsbCA/IG5ldyBXcmFwcGVkTm9kZUV4cHIoZmFjYWRlLnByb3ZpZGVycykgOiBudWxsLFxuICAgIHZpZXdRdWVyaWVzOiBmYWNhZGUudmlld1F1ZXJpZXMubWFwKGNvbnZlcnRUb1IzUXVlcnlNZXRhZGF0YSksXG4gICAgZnVsbEluaGVyaXRhbmNlOiBmYWxzZSxcbiAgICBob3N0RGlyZWN0aXZlczogY29udmVydEhvc3REaXJlY3RpdmVzVG9NZXRhZGF0YShmYWNhZGUpLFxuICB9O1xufVxuXG5mdW5jdGlvbiBjb252ZXJ0RGVjbGFyZURpcmVjdGl2ZUZhY2FkZVRvTWV0YWRhdGEoXG4gIGRlY2xhcmF0aW9uOiBSM0RlY2xhcmVEaXJlY3RpdmVGYWNhZGUsXG4gIHR5cGVTb3VyY2VTcGFuOiBQYXJzZVNvdXJjZVNwYW4sXG4pOiBSM0RpcmVjdGl2ZU1ldGFkYXRhIHtcbiAgcmV0dXJuIHtcbiAgICBuYW1lOiBkZWNsYXJhdGlvbi50eXBlLm5hbWUsXG4gICAgdHlwZTogd3JhcFJlZmVyZW5jZShkZWNsYXJhdGlvbi50eXBlKSxcbiAgICB0eXBlU291cmNlU3BhbixcbiAgICBzZWxlY3RvcjogZGVjbGFyYXRpb24uc2VsZWN0b3IgPz8gbnVsbCxcbiAgICBpbnB1dHM6IGRlY2xhcmF0aW9uLmlucHV0cyA/IGlucHV0c1BhcnRpYWxNZXRhZGF0YVRvSW5wdXRNZXRhZGF0YShkZWNsYXJhdGlvbi5pbnB1dHMpIDoge30sXG4gICAgb3V0cHV0czogZGVjbGFyYXRpb24ub3V0cHV0cyA/PyB7fSxcbiAgICBob3N0OiBjb252ZXJ0SG9zdERlY2xhcmF0aW9uVG9NZXRhZGF0YShkZWNsYXJhdGlvbi5ob3N0KSxcbiAgICBxdWVyaWVzOiAoZGVjbGFyYXRpb24ucXVlcmllcyA/PyBbXSkubWFwKGNvbnZlcnRRdWVyeURlY2xhcmF0aW9uVG9NZXRhZGF0YSksXG4gICAgdmlld1F1ZXJpZXM6IChkZWNsYXJhdGlvbi52aWV3UXVlcmllcyA/PyBbXSkubWFwKGNvbnZlcnRRdWVyeURlY2xhcmF0aW9uVG9NZXRhZGF0YSksXG4gICAgcHJvdmlkZXJzOlxuICAgICAgZGVjbGFyYXRpb24ucHJvdmlkZXJzICE9PSB1bmRlZmluZWQgPyBuZXcgV3JhcHBlZE5vZGVFeHByKGRlY2xhcmF0aW9uLnByb3ZpZGVycykgOiBudWxsLFxuICAgIGV4cG9ydEFzOiBkZWNsYXJhdGlvbi5leHBvcnRBcyA/PyBudWxsLFxuICAgIHVzZXNJbmhlcml0YW5jZTogZGVjbGFyYXRpb24udXNlc0luaGVyaXRhbmNlID8/IGZhbHNlLFxuICAgIGxpZmVjeWNsZToge3VzZXNPbkNoYW5nZXM6IGRlY2xhcmF0aW9uLnVzZXNPbkNoYW5nZXMgPz8gZmFsc2V9LFxuICAgIGRlcHM6IG51bGwsXG4gICAgdHlwZUFyZ3VtZW50Q291bnQ6IDAsXG4gICAgZnVsbEluaGVyaXRhbmNlOiBmYWxzZSxcbiAgICBpc1N0YW5kYWxvbmU6IGRlY2xhcmF0aW9uLmlzU3RhbmRhbG9uZSA/PyBmYWxzZSxcbiAgICBpc1NpZ25hbDogZGVjbGFyYXRpb24uaXNTaWduYWwgPz8gZmFsc2UsXG4gICAgaG9zdERpcmVjdGl2ZXM6IGNvbnZlcnRIb3N0RGlyZWN0aXZlc1RvTWV0YWRhdGEoZGVjbGFyYXRpb24pLFxuICB9O1xufVxuXG5mdW5jdGlvbiBjb252ZXJ0SG9zdERlY2xhcmF0aW9uVG9NZXRhZGF0YShcbiAgaG9zdDogUjNEZWNsYXJlRGlyZWN0aXZlRmFjYWRlWydob3N0J10gPSB7fSxcbik6IFIzSG9zdE1ldGFkYXRhIHtcbiAgcmV0dXJuIHtcbiAgICBhdHRyaWJ1dGVzOiBjb252ZXJ0T3BhcXVlVmFsdWVzVG9FeHByZXNzaW9ucyhob3N0LmF0dHJpYnV0ZXMgPz8ge30pLFxuICAgIGxpc3RlbmVyczogaG9zdC5saXN0ZW5lcnMgPz8ge30sXG4gICAgcHJvcGVydGllczogaG9zdC5wcm9wZXJ0aWVzID8/IHt9LFxuICAgIHNwZWNpYWxBdHRyaWJ1dGVzOiB7XG4gICAgICBjbGFzc0F0dHI6IGhvc3QuY2xhc3NBdHRyaWJ1dGUsXG4gICAgICBzdHlsZUF0dHI6IGhvc3Quc3R5bGVBdHRyaWJ1dGUsXG4gICAgfSxcbiAgfTtcbn1cblxuZnVuY3Rpb24gY29udmVydEhvc3REaXJlY3RpdmVzVG9NZXRhZGF0YShcbiAgbWV0YWRhdGE6IFIzRGVjbGFyZURpcmVjdGl2ZUZhY2FkZSB8IFIzRGlyZWN0aXZlTWV0YWRhdGFGYWNhZGUsXG4pOiBSM0hvc3REaXJlY3RpdmVNZXRhZGF0YVtdIHwgbnVsbCB7XG4gIGlmIChtZXRhZGF0YS5ob3N0RGlyZWN0aXZlcz8ubGVuZ3RoKSB7XG4gICAgcmV0dXJuIG1ldGFkYXRhLmhvc3REaXJlY3RpdmVzLm1hcCgoaG9zdERpcmVjdGl2ZSkgPT4ge1xuICAgICAgcmV0dXJuIHR5cGVvZiBob3N0RGlyZWN0aXZlID09PSAnZnVuY3Rpb24nXG4gICAgICAgID8ge1xuICAgICAgICAgICAgZGlyZWN0aXZlOiB3cmFwUmVmZXJlbmNlKGhvc3REaXJlY3RpdmUpLFxuICAgICAgICAgICAgaW5wdXRzOiBudWxsLFxuICAgICAgICAgICAgb3V0cHV0czogbnVsbCxcbiAgICAgICAgICAgIGlzRm9yd2FyZFJlZmVyZW5jZTogZmFsc2UsXG4gICAgICAgICAgfVxuICAgICAgICA6IHtcbiAgICAgICAgICAgIGRpcmVjdGl2ZTogd3JhcFJlZmVyZW5jZShob3N0RGlyZWN0aXZlLmRpcmVjdGl2ZSksXG4gICAgICAgICAgICBpc0ZvcndhcmRSZWZlcmVuY2U6IGZhbHNlLFxuICAgICAgICAgICAgaW5wdXRzOiBob3N0RGlyZWN0aXZlLmlucHV0cyA/IHBhcnNlTWFwcGluZ1N0cmluZ0FycmF5KGhvc3REaXJlY3RpdmUuaW5wdXRzKSA6IG51bGwsXG4gICAgICAgICAgICBvdXRwdXRzOiBob3N0RGlyZWN0aXZlLm91dHB1dHMgPyBwYXJzZU1hcHBpbmdTdHJpbmdBcnJheShob3N0RGlyZWN0aXZlLm91dHB1dHMpIDogbnVsbCxcbiAgICAgICAgICB9O1xuICAgIH0pO1xuICB9XG5cbiAgcmV0dXJuIG51bGw7XG59XG5cbmZ1bmN0aW9uIGNvbnZlcnRPcGFxdWVWYWx1ZXNUb0V4cHJlc3Npb25zKG9iajoge1trZXk6IHN0cmluZ106IE9wYXF1ZVZhbHVlfSk6IHtcbiAgW2tleTogc3RyaW5nXTogV3JhcHBlZE5vZGVFeHByPHVua25vd24+O1xufSB7XG4gIGNvbnN0IHJlc3VsdDoge1trZXk6IHN0cmluZ106IFdyYXBwZWROb2RlRXhwcjx1bmtub3duPn0gPSB7fTtcbiAgZm9yIChjb25zdCBrZXkgb2YgT2JqZWN0LmtleXMob2JqKSkge1xuICAgIHJlc3VsdFtrZXldID0gbmV3IFdyYXBwZWROb2RlRXhwcihvYmpba2V5XSk7XG4gIH1cbiAgcmV0dXJuIHJlc3VsdDtcbn1cblxuZnVuY3Rpb24gY29udmVydERlY2xhcmVDb21wb25lbnRGYWNhZGVUb01ldGFkYXRhKFxuICBkZWNsOiBSM0RlY2xhcmVDb21wb25lbnRGYWNhZGUsXG4gIHR5cGVTb3VyY2VTcGFuOiBQYXJzZVNvdXJjZVNwYW4sXG4gIHNvdXJjZU1hcFVybDogc3RyaW5nLFxuKTogUjNDb21wb25lbnRNZXRhZGF0YTxSM1RlbXBsYXRlRGVwZW5kZW5jeU1ldGFkYXRhPiB7XG4gIGNvbnN0IHt0ZW1wbGF0ZSwgaW50ZXJwb2xhdGlvbiwgZGVmZXJ9ID0gcGFyc2VKaXRUZW1wbGF0ZShcbiAgICBkZWNsLnRlbXBsYXRlLFxuICAgIGRlY2wudHlwZS5uYW1lLFxuICAgIHNvdXJjZU1hcFVybCxcbiAgICBkZWNsLnByZXNlcnZlV2hpdGVzcGFjZXMgPz8gZmFsc2UsXG4gICAgZGVjbC5pbnRlcnBvbGF0aW9uLFxuICAgIGRlY2wuZGVmZXJCbG9ja0RlcGVuZGVuY2llcyxcbiAgKTtcblxuICBjb25zdCBkZWNsYXJhdGlvbnM6IFIzVGVtcGxhdGVEZXBlbmRlbmN5TWV0YWRhdGFbXSA9IFtdO1xuICBpZiAoZGVjbC5kZXBlbmRlbmNpZXMpIHtcbiAgICBmb3IgKGNvbnN0IGlubmVyRGVwIG9mIGRlY2wuZGVwZW5kZW5jaWVzKSB7XG4gICAgICBzd2l0Y2ggKGlubmVyRGVwLmtpbmQpIHtcbiAgICAgICAgY2FzZSAnZGlyZWN0aXZlJzpcbiAgICAgICAgY2FzZSAnY29tcG9uZW50JzpcbiAgICAgICAgICBkZWNsYXJhdGlvbnMucHVzaChjb252ZXJ0RGlyZWN0aXZlRGVjbGFyYXRpb25Ub01ldGFkYXRhKGlubmVyRGVwKSk7XG4gICAgICAgICAgYnJlYWs7XG4gICAgICAgIGNhc2UgJ3BpcGUnOlxuICAgICAgICAgIGRlY2xhcmF0aW9ucy5wdXNoKGNvbnZlcnRQaXBlRGVjbGFyYXRpb25Ub01ldGFkYXRhKGlubmVyRGVwKSk7XG4gICAgICAgICAgYnJlYWs7XG4gICAgICB9XG4gICAgfVxuICB9IGVsc2UgaWYgKGRlY2wuY29tcG9uZW50cyB8fCBkZWNsLmRpcmVjdGl2ZXMgfHwgZGVjbC5waXBlcykge1xuICAgIC8vIEV4aXN0aW5nIGRlY2xhcmF0aW9ucyBvbiBOUE0gbWF5IG5vdCBiZSB1c2luZyB0aGUgbmV3IGBkZXBlbmRlbmNpZXNgIG1lcmdlZCBmaWVsZCwgYW5kIG1heVxuICAgIC8vIGhhdmUgc2VwYXJhdGUgZmllbGRzIGZvciBkZXBlbmRlbmNpZXMgaW5zdGVhZC4gVW5pZnkgdGhlbSBmb3IgSklUIGNvbXBpbGF0aW9uLlxuICAgIGRlY2wuY29tcG9uZW50cyAmJlxuICAgICAgZGVjbGFyYXRpb25zLnB1c2goXG4gICAgICAgIC4uLmRlY2wuY29tcG9uZW50cy5tYXAoKGRpcikgPT5cbiAgICAgICAgICBjb252ZXJ0RGlyZWN0aXZlRGVjbGFyYXRpb25Ub01ldGFkYXRhKGRpciwgLyogaXNDb21wb25lbnQgKi8gdHJ1ZSksXG4gICAgICAgICksXG4gICAgICApO1xuICAgIGRlY2wuZGlyZWN0aXZlcyAmJlxuICAgICAgZGVjbGFyYXRpb25zLnB1c2goXG4gICAgICAgIC4uLmRlY2wuZGlyZWN0aXZlcy5tYXAoKGRpcikgPT4gY29udmVydERpcmVjdGl2ZURlY2xhcmF0aW9uVG9NZXRhZGF0YShkaXIpKSxcbiAgICAgICk7XG4gICAgZGVjbC5waXBlcyAmJiBkZWNsYXJhdGlvbnMucHVzaCguLi5jb252ZXJ0UGlwZU1hcFRvTWV0YWRhdGEoZGVjbC5waXBlcykpO1xuICB9XG5cbiAgcmV0dXJuIHtcbiAgICAuLi5jb252ZXJ0RGVjbGFyZURpcmVjdGl2ZUZhY2FkZVRvTWV0YWRhdGEoZGVjbCwgdHlwZVNvdXJjZVNwYW4pLFxuICAgIHRlbXBsYXRlLFxuICAgIHN0eWxlczogZGVjbC5zdHlsZXMgPz8gW10sXG4gICAgZGVjbGFyYXRpb25zLFxuICAgIHZpZXdQcm92aWRlcnM6XG4gICAgICBkZWNsLnZpZXdQcm92aWRlcnMgIT09IHVuZGVmaW5lZCA/IG5ldyBXcmFwcGVkTm9kZUV4cHIoZGVjbC52aWV3UHJvdmlkZXJzKSA6IG51bGwsXG4gICAgYW5pbWF0aW9uczogZGVjbC5hbmltYXRpb25zICE9PSB1bmRlZmluZWQgPyBuZXcgV3JhcHBlZE5vZGVFeHByKGRlY2wuYW5pbWF0aW9ucykgOiBudWxsLFxuICAgIGRlZmVyLFxuXG4gICAgY2hhbmdlRGV0ZWN0aW9uOiBkZWNsLmNoYW5nZURldGVjdGlvbiA/PyBDaGFuZ2VEZXRlY3Rpb25TdHJhdGVneS5EZWZhdWx0LFxuICAgIGVuY2Fwc3VsYXRpb246IGRlY2wuZW5jYXBzdWxhdGlvbiA/PyBWaWV3RW5jYXBzdWxhdGlvbi5FbXVsYXRlZCxcbiAgICBpbnRlcnBvbGF0aW9uLFxuICAgIGRlY2xhcmF0aW9uTGlzdEVtaXRNb2RlOiBEZWNsYXJhdGlvbkxpc3RFbWl0TW9kZS5DbG9zdXJlUmVzb2x2ZWQsXG4gICAgcmVsYXRpdmVDb250ZXh0RmlsZVBhdGg6ICcnLFxuICAgIGkxOG5Vc2VFeHRlcm5hbElkczogdHJ1ZSxcbiAgfTtcbn1cblxuZnVuY3Rpb24gY29udmVydERlY2xhcmF0aW9uRmFjYWRlVG9NZXRhZGF0YShcbiAgZGVjbGFyYXRpb246IFIzVGVtcGxhdGVEZXBlbmRlbmN5RmFjYWRlLFxuKTogUjNUZW1wbGF0ZURlcGVuZGVuY3kge1xuICByZXR1cm4ge1xuICAgIC4uLmRlY2xhcmF0aW9uLFxuICAgIHR5cGU6IG5ldyBXcmFwcGVkTm9kZUV4cHIoZGVjbGFyYXRpb24udHlwZSksXG4gIH07XG59XG5cbmZ1bmN0aW9uIGNvbnZlcnREaXJlY3RpdmVEZWNsYXJhdGlvblRvTWV0YWRhdGEoXG4gIGRlY2xhcmF0aW9uOiBSM0RlY2xhcmVEaXJlY3RpdmVEZXBlbmRlbmN5RmFjYWRlLFxuICBpc0NvbXBvbmVudDogdHJ1ZSB8IG51bGwgPSBudWxsLFxuKTogUjNEaXJlY3RpdmVEZXBlbmRlbmN5TWV0YWRhdGEge1xuICByZXR1cm4ge1xuICAgIGtpbmQ6IFIzVGVtcGxhdGVEZXBlbmRlbmN5S2luZC5EaXJlY3RpdmUsXG4gICAgaXNDb21wb25lbnQ6IGlzQ29tcG9uZW50IHx8IGRlY2xhcmF0aW9uLmtpbmQgPT09ICdjb21wb25lbnQnLFxuICAgIHNlbGVjdG9yOiBkZWNsYXJhdGlvbi5zZWxlY3RvcixcbiAgICB0eXBlOiBuZXcgV3JhcHBlZE5vZGVFeHByKGRlY2xhcmF0aW9uLnR5cGUpLFxuICAgIGlucHV0czogZGVjbGFyYXRpb24uaW5wdXRzID8/IFtdLFxuICAgIG91dHB1dHM6IGRlY2xhcmF0aW9uLm91dHB1dHMgPz8gW10sXG4gICAgZXhwb3J0QXM6IGRlY2xhcmF0aW9uLmV4cG9ydEFzID8/IG51bGwsXG4gIH07XG59XG5cbmZ1bmN0aW9uIGNvbnZlcnRQaXBlTWFwVG9NZXRhZGF0YShcbiAgcGlwZXM6IFIzRGVjbGFyZUNvbXBvbmVudEZhY2FkZVsncGlwZXMnXSxcbik6IFIzUGlwZURlcGVuZGVuY3lNZXRhZGF0YVtdIHtcbiAgaWYgKCFwaXBlcykge1xuICAgIHJldHVybiBbXTtcbiAgfVxuXG4gIHJldHVybiBPYmplY3Qua2V5cyhwaXBlcykubWFwKChuYW1lKSA9PiB7XG4gICAgcmV0dXJuIHtcbiAgICAgIGtpbmQ6IFIzVGVtcGxhdGVEZXBlbmRlbmN5S2luZC5QaXBlLFxuICAgICAgbmFtZSxcbiAgICAgIHR5cGU6IG5ldyBXcmFwcGVkTm9kZUV4cHIocGlwZXNbbmFtZV0pLFxuICAgIH07XG4gIH0pO1xufVxuXG5mdW5jdGlvbiBjb252ZXJ0UGlwZURlY2xhcmF0aW9uVG9NZXRhZGF0YShcbiAgcGlwZTogUjNEZWNsYXJlUGlwZURlcGVuZGVuY3lGYWNhZGUsXG4pOiBSM1BpcGVEZXBlbmRlbmN5TWV0YWRhdGEge1xuICByZXR1cm4ge1xuICAgIGtpbmQ6IFIzVGVtcGxhdGVEZXBlbmRlbmN5S2luZC5QaXBlLFxuICAgIG5hbWU6IHBpcGUubmFtZSxcbiAgICB0eXBlOiBuZXcgV3JhcHBlZE5vZGVFeHByKHBpcGUudHlwZSksXG4gIH07XG59XG5cbmZ1bmN0aW9uIHBhcnNlSml0VGVtcGxhdGUoXG4gIHRlbXBsYXRlOiBzdHJpbmcsXG4gIHR5cGVOYW1lOiBzdHJpbmcsXG4gIHNvdXJjZU1hcFVybDogc3RyaW5nLFxuICBwcmVzZXJ2ZVdoaXRlc3BhY2VzOiBib29sZWFuLFxuICBpbnRlcnBvbGF0aW9uOiBbc3RyaW5nLCBzdHJpbmddIHwgdW5kZWZpbmVkLFxuICBkZWZlckJsb2NrRGVwZW5kZW5jaWVzOiAoKCkgPT4gUHJvbWlzZTx1bmtub3duPiB8IG51bGwpW10gfCB1bmRlZmluZWQsXG4pIHtcbiAgY29uc3QgaW50ZXJwb2xhdGlvbkNvbmZpZyA9IGludGVycG9sYXRpb25cbiAgICA/IEludGVycG9sYXRpb25Db25maWcuZnJvbUFycmF5KGludGVycG9sYXRpb24pXG4gICAgOiBERUZBVUxUX0lOVEVSUE9MQVRJT05fQ09ORklHO1xuICAvLyBQYXJzZSB0aGUgdGVtcGxhdGUgYW5kIGNoZWNrIGZvciBlcnJvcnMuXG4gIGNvbnN0IHBhcnNlZCA9IHBhcnNlVGVtcGxhdGUodGVtcGxhdGUsIHNvdXJjZU1hcFVybCwge1xuICAgIHByZXNlcnZlV2hpdGVzcGFjZXMsXG4gICAgaW50ZXJwb2xhdGlvbkNvbmZpZyxcbiAgfSk7XG4gIGlmIChwYXJzZWQuZXJyb3JzICE9PSBudWxsKSB7XG4gICAgY29uc3QgZXJyb3JzID0gcGFyc2VkLmVycm9ycy5tYXAoKGVycikgPT4gZXJyLnRvU3RyaW5nKCkpLmpvaW4oJywgJyk7XG4gICAgdGhyb3cgbmV3IEVycm9yKGBFcnJvcnMgZHVyaW5nIEpJVCBjb21waWxhdGlvbiBvZiB0ZW1wbGF0ZSBmb3IgJHt0eXBlTmFtZX06ICR7ZXJyb3JzfWApO1xuICB9XG4gIGNvbnN0IGJpbmRlciA9IG5ldyBSM1RhcmdldEJpbmRlcihuZXcgU2VsZWN0b3JNYXRjaGVyKCkpO1xuICBjb25zdCBib3VuZFRhcmdldCA9IGJpbmRlci5iaW5kKHt0ZW1wbGF0ZTogcGFyc2VkLm5vZGVzfSk7XG5cbiAgcmV0dXJuIHtcbiAgICB0ZW1wbGF0ZTogcGFyc2VkLFxuICAgIGludGVycG9sYXRpb246IGludGVycG9sYXRpb25Db25maWcsXG4gICAgZGVmZXI6IGNyZWF0ZVIzQ29tcG9uZW50RGVmZXJNZXRhZGF0YShib3VuZFRhcmdldCwgZGVmZXJCbG9ja0RlcGVuZGVuY2llcyksXG4gIH07XG59XG5cbi8qKlxuICogQ29udmVydCB0aGUgZXhwcmVzc2lvbiwgaWYgcHJlc2VudCB0byBhbiBgUjNQcm92aWRlckV4cHJlc3Npb25gLlxuICpcbiAqIEluIEpJVCBtb2RlIHdlIGRvIG5vdCB3YW50IHRoZSBjb21waWxlciB0byB3cmFwIHRoZSBleHByZXNzaW9uIGluIGEgYGZvcndhcmRSZWYoKWAgY2FsbCBiZWNhdXNlLFxuICogaWYgaXQgaXMgcmVmZXJlbmNpbmcgYSB0eXBlIHRoYXQgaGFzIG5vdCB5ZXQgYmVlbiBkZWZpbmVkLCBpdCB3aWxsIGhhdmUgYWxyZWFkeSBiZWVuIHdyYXBwZWQgaW5cbiAqIGEgYGZvcndhcmRSZWYoKWAgLSBlaXRoZXIgYnkgdGhlIGFwcGxpY2F0aW9uIGRldmVsb3BlciBvciBkdXJpbmcgcGFydGlhbC1jb21waWxhdGlvbi4gVGh1cyB3ZSBjYW5cbiAqIHVzZSBgRm9yd2FyZFJlZkhhbmRsaW5nLk5vbmVgLlxuICovXG5mdW5jdGlvbiBjb252ZXJ0VG9Qcm92aWRlckV4cHJlc3Npb24oXG4gIG9iajogYW55LFxuICBwcm9wZXJ0eTogc3RyaW5nLFxuKTogTWF5YmVGb3J3YXJkUmVmRXhwcmVzc2lvbiB8IHVuZGVmaW5lZCB7XG4gIGlmIChvYmouaGFzT3duUHJvcGVydHkocHJvcGVydHkpKSB7XG4gICAgcmV0dXJuIGNyZWF0ZU1heUJlRm9yd2FyZFJlZkV4cHJlc3Npb24oXG4gICAgICBuZXcgV3JhcHBlZE5vZGVFeHByKG9ialtwcm9wZXJ0eV0pLFxuICAgICAgRm9yd2FyZFJlZkhhbmRsaW5nLk5vbmUsXG4gICAgKTtcbiAgfSBlbHNlIHtcbiAgICByZXR1cm4gdW5kZWZpbmVkO1xuICB9XG59XG5cbmZ1bmN0aW9uIHdyYXBFeHByZXNzaW9uKG9iajogYW55LCBwcm9wZXJ0eTogc3RyaW5nKTogV3JhcHBlZE5vZGVFeHByPGFueT4gfCB1bmRlZmluZWQge1xuICBpZiAob2JqLmhhc093blByb3BlcnR5KHByb3BlcnR5KSkge1xuICAgIHJldHVybiBuZXcgV3JhcHBlZE5vZGVFeHByKG9ialtwcm9wZXJ0eV0pO1xuICB9IGVsc2Uge1xuICAgIHJldHVybiB1bmRlZmluZWQ7XG4gIH1cbn1cblxuZnVuY3Rpb24gY29tcHV0ZVByb3ZpZGVkSW4oXG4gIHByb3ZpZGVkSW46IEZ1bmN0aW9uIHwgc3RyaW5nIHwgbnVsbCB8IHVuZGVmaW5lZCxcbik6IE1heWJlRm9yd2FyZFJlZkV4cHJlc3Npb24ge1xuICBjb25zdCBleHByZXNzaW9uID1cbiAgICB0eXBlb2YgcHJvdmlkZWRJbiA9PT0gJ2Z1bmN0aW9uJ1xuICAgICAgPyBuZXcgV3JhcHBlZE5vZGVFeHByKHByb3ZpZGVkSW4pXG4gICAgICA6IG5ldyBMaXRlcmFsRXhwcihwcm92aWRlZEluID8/IG51bGwpO1xuICAvLyBTZWUgYGNvbnZlcnRUb1Byb3ZpZGVyRXhwcmVzc2lvbigpYCBmb3Igd2h5IHRoaXMgdXNlcyBgRm9yd2FyZFJlZkhhbmRsaW5nLk5vbmVgLlxuICByZXR1cm4gY3JlYXRlTWF5QmVGb3J3YXJkUmVmRXhwcmVzc2lvbihleHByZXNzaW9uLCBGb3J3YXJkUmVmSGFuZGxpbmcuTm9uZSk7XG59XG5cbmZ1bmN0aW9uIGNvbnZlcnRSM0RlcGVuZGVuY3lNZXRhZGF0YUFycmF5KFxuICBmYWNhZGVzOiBSM0RlcGVuZGVuY3lNZXRhZGF0YUZhY2FkZVtdIHwgbnVsbCB8IHVuZGVmaW5lZCxcbik6IFIzRGVwZW5kZW5jeU1ldGFkYXRhW10gfCBudWxsIHtcbiAgcmV0dXJuIGZhY2FkZXMgPT0gbnVsbCA/IG51bGwgOiBmYWNhZGVzLm1hcChjb252ZXJ0UjNEZXBlbmRlbmN5TWV0YWRhdGEpO1xufVxuXG5mdW5jdGlvbiBjb252ZXJ0UjNEZXBlbmRlbmN5TWV0YWRhdGEoZmFjYWRlOiBSM0RlcGVuZGVuY3lNZXRhZGF0YUZhY2FkZSk6IFIzRGVwZW5kZW5jeU1ldGFkYXRhIHtcbiAgY29uc3QgaXNBdHRyaWJ1dGVEZXAgPSBmYWNhZGUuYXR0cmlidXRlICE9IG51bGw7IC8vIGJvdGggYG51bGxgIGFuZCBgdW5kZWZpbmVkYFxuICBjb25zdCByYXdUb2tlbiA9IGZhY2FkZS50b2tlbiA9PT0gbnVsbCA/IG51bGwgOiBuZXcgV3JhcHBlZE5vZGVFeHByKGZhY2FkZS50b2tlbik7XG4gIC8vIEluIEpJVCBtb2RlLCBpZiB0aGUgZGVwIGlzIGFuIGBAQXR0cmlidXRlKClgIHRoZW4gd2UgdXNlIHRoZSBhdHRyaWJ1dGUgbmFtZSBnaXZlbiBpblxuICAvLyBgYXR0cmlidXRlYCByYXRoZXIgdGhhbiB0aGUgYHRva2VuYC5cbiAgY29uc3QgdG9rZW4gPSBpc0F0dHJpYnV0ZURlcCA/IG5ldyBXcmFwcGVkTm9kZUV4cHIoZmFjYWRlLmF0dHJpYnV0ZSkgOiByYXdUb2tlbjtcbiAgcmV0dXJuIGNyZWF0ZVIzRGVwZW5kZW5jeU1ldGFkYXRhKFxuICAgIHRva2VuLFxuICAgIGlzQXR0cmlidXRlRGVwLFxuICAgIGZhY2FkZS5ob3N0LFxuICAgIGZhY2FkZS5vcHRpb25hbCxcbiAgICBmYWNhZGUuc2VsZixcbiAgICBmYWNhZGUuc2tpcFNlbGYsXG4gICk7XG59XG5cbmZ1bmN0aW9uIGNvbnZlcnRSM0RlY2xhcmVEZXBlbmRlbmN5TWV0YWRhdGEoXG4gIGZhY2FkZTogUjNEZWNsYXJlRGVwZW5kZW5jeU1ldGFkYXRhRmFjYWRlLFxuKTogUjNEZXBlbmRlbmN5TWV0YWRhdGEge1xuICBjb25zdCBpc0F0dHJpYnV0ZURlcCA9IGZhY2FkZS5hdHRyaWJ1dGUgPz8gZmFsc2U7XG4gIGNvbnN0IHRva2VuID0gZmFjYWRlLnRva2VuID09PSBudWxsID8gbnVsbCA6IG5ldyBXcmFwcGVkTm9kZUV4cHIoZmFjYWRlLnRva2VuKTtcbiAgcmV0dXJuIGNyZWF0ZVIzRGVwZW5kZW5jeU1ldGFkYXRhKFxuICAgIHRva2VuLFxuICAgIGlzQXR0cmlidXRlRGVwLFxuICAgIGZhY2FkZS5ob3N0ID8/IGZhbHNlLFxuICAgIGZhY2FkZS5vcHRpb25hbCA/PyBmYWxzZSxcbiAgICBmYWNhZGUuc2VsZiA/PyBmYWxzZSxcbiAgICBmYWNhZGUuc2tpcFNlbGYgPz8gZmFsc2UsXG4gICk7XG59XG5cbmZ1bmN0aW9uIGNyZWF0ZVIzRGVwZW5kZW5jeU1ldGFkYXRhKFxuICB0b2tlbjogV3JhcHBlZE5vZGVFeHByPHVua25vd24+IHwgbnVsbCxcbiAgaXNBdHRyaWJ1dGVEZXA6IGJvb2xlYW4sXG4gIGhvc3Q6IGJvb2xlYW4sXG4gIG9wdGlvbmFsOiBib29sZWFuLFxuICBzZWxmOiBib29sZWFuLFxuICBza2lwU2VsZjogYm9vbGVhbixcbik6IFIzRGVwZW5kZW5jeU1ldGFkYXRhIHtcbiAgLy8gSWYgdGhlIGRlcCBpcyBhbiBgQEF0dHJpYnV0ZSgpYCB0aGUgYGF0dHJpYnV0ZU5hbWVUeXBlYCBvdWdodCB0byBiZSB0aGUgYHVua25vd25gIHR5cGUuXG4gIC8vIEJ1dCB0eXBlcyBhcmUgbm90IGF2YWlsYWJsZSBhdCBydW50aW1lIHNvIHdlIGp1c3QgdXNlIGEgbGl0ZXJhbCBgXCI8dW5rbm93bj5cImAgc3RyaW5nIGFzIGEgZHVtbXlcbiAgLy8gbWFya2VyLlxuICBjb25zdCBhdHRyaWJ1dGVOYW1lVHlwZSA9IGlzQXR0cmlidXRlRGVwID8gbGl0ZXJhbCgndW5rbm93bicpIDogbnVsbDtcbiAgcmV0dXJuIHt0b2tlbiwgYXR0cmlidXRlTmFtZVR5cGUsIGhvc3QsIG9wdGlvbmFsLCBzZWxmLCBza2lwU2VsZn07XG59XG5cbmZ1bmN0aW9uIGNyZWF0ZVIzQ29tcG9uZW50RGVmZXJNZXRhZGF0YShcbiAgYm91bmRUYXJnZXQ6IEJvdW5kVGFyZ2V0PGFueT4sXG4gIGRlZmVyQmxvY2tEZXBlbmRlbmNpZXM6ICgoKSA9PiBQcm9taXNlPHVua25vd24+IHwgbnVsbClbXSB8IHVuZGVmaW5lZCxcbik6IFIzQ29tcG9uZW50RGVmZXJNZXRhZGF0YSB7XG4gIGNvbnN0IGRlZmVycmVkQmxvY2tzID0gYm91bmRUYXJnZXQuZ2V0RGVmZXJCbG9ja3MoKTtcbiAgY29uc3QgYmxvY2tzID0gbmV3IE1hcDxEZWZlcnJlZEJsb2NrLCBFeHByZXNzaW9uIHwgbnVsbD4oKTtcblxuICBmb3IgKGxldCBpID0gMDsgaSA8IGRlZmVycmVkQmxvY2tzLmxlbmd0aDsgaSsrKSB7XG4gICAgY29uc3QgZGVwZW5kZW5jeUZuID0gZGVmZXJCbG9ja0RlcGVuZGVuY2llcz8uW2ldO1xuICAgIGJsb2Nrcy5zZXQoZGVmZXJyZWRCbG9ja3NbaV0sIGRlcGVuZGVuY3lGbiA/IG5ldyBXcmFwcGVkTm9kZUV4cHIoZGVwZW5kZW5jeUZuKSA6IG51bGwpO1xuICB9XG5cbiAgcmV0dXJuIHttb2RlOiBEZWZlckJsb2NrRGVwc0VtaXRNb2RlLlBlckJsb2NrLCBibG9ja3N9O1xufVxuXG5mdW5jdGlvbiBleHRyYWN0SG9zdEJpbmRpbmdzKFxuICBwcm9wTWV0YWRhdGE6IHtba2V5OiBzdHJpbmddOiBhbnlbXX0sXG4gIHNvdXJjZVNwYW46IFBhcnNlU291cmNlU3BhbixcbiAgaG9zdD86IHtba2V5OiBzdHJpbmddOiBzdHJpbmd9LFxuKTogUGFyc2VkSG9zdEJpbmRpbmdzIHtcbiAgLy8gRmlyc3QgcGFyc2UgdGhlIGRlY2xhcmF0aW9ucyBmcm9tIHRoZSBtZXRhZGF0YS5cbiAgY29uc3QgYmluZGluZ3MgPSBwYXJzZUhvc3RCaW5kaW5ncyhob3N0IHx8IHt9KTtcblxuICAvLyBBZnRlciB0aGF0IGNoZWNrIGhvc3QgYmluZGluZ3MgZm9yIGVycm9yc1xuICBjb25zdCBlcnJvcnMgPSB2ZXJpZnlIb3N0QmluZGluZ3MoYmluZGluZ3MsIHNvdXJjZVNwYW4pO1xuICBpZiAoZXJyb3JzLmxlbmd0aCkge1xuICAgIHRocm93IG5ldyBFcnJvcihlcnJvcnMubWFwKChlcnJvcjogUGFyc2VFcnJvcikgPT4gZXJyb3IubXNnKS5qb2luKCdcXG4nKSk7XG4gIH1cblxuICAvLyBOZXh0LCBsb29wIG92ZXIgdGhlIHByb3BlcnRpZXMgb2YgdGhlIG9iamVjdCwgbG9va2luZyBmb3IgQEhvc3RCaW5kaW5nIGFuZCBASG9zdExpc3RlbmVyLlxuICBmb3IgKGNvbnN0IGZpZWxkIGluIHByb3BNZXRhZGF0YSkge1xuICAgIGlmIChwcm9wTWV0YWRhdGEuaGFzT3duUHJvcGVydHkoZmllbGQpKSB7XG4gICAgICBwcm9wTWV0YWRhdGFbZmllbGRdLmZvckVhY2goKGFubikgPT4ge1xuICAgICAgICBpZiAoaXNIb3N0QmluZGluZyhhbm4pKSB7XG4gICAgICAgICAgLy8gU2luY2UgdGhpcyBpcyBhIGRlY29yYXRvciwgd2Uga25vdyB0aGF0IHRoZSB2YWx1ZSBpcyBhIGNsYXNzIG1lbWJlci4gQWx3YXlzIGFjY2VzcyBpdFxuICAgICAgICAgIC8vIHRocm91Z2ggYHRoaXNgIHNvIHRoYXQgZnVydGhlciBkb3duIHRoZSBsaW5lIGl0IGNhbid0IGJlIGNvbmZ1c2VkIGZvciBhIGxpdGVyYWwgdmFsdWVcbiAgICAgICAgICAvLyAoZS5nLiBpZiB0aGVyZSdzIGEgcHJvcGVydHkgY2FsbGVkIGB0cnVlYCkuXG4gICAgICAgICAgYmluZGluZ3MucHJvcGVydGllc1thbm4uaG9zdFByb3BlcnR5TmFtZSB8fCBmaWVsZF0gPSBnZXRTYWZlUHJvcGVydHlBY2Nlc3NTdHJpbmcoXG4gICAgICAgICAgICAndGhpcycsXG4gICAgICAgICAgICBmaWVsZCxcbiAgICAgICAgICApO1xuICAgICAgICB9IGVsc2UgaWYgKGlzSG9zdExpc3RlbmVyKGFubikpIHtcbiAgICAgICAgICBiaW5kaW5ncy5saXN0ZW5lcnNbYW5uLmV2ZW50TmFtZSB8fCBmaWVsZF0gPSBgJHtmaWVsZH0oJHsoYW5uLmFyZ3MgfHwgW10pLmpvaW4oJywnKX0pYDtcbiAgICAgICAgfVxuICAgICAgfSk7XG4gICAgfVxuICB9XG5cbiAgcmV0dXJuIGJpbmRpbmdzO1xufVxuXG5mdW5jdGlvbiBpc0hvc3RCaW5kaW5nKHZhbHVlOiBhbnkpOiB2YWx1ZSBpcyBIb3N0QmluZGluZyB7XG4gIHJldHVybiB2YWx1ZS5uZ01ldGFkYXRhTmFtZSA9PT0gJ0hvc3RCaW5kaW5nJztcbn1cblxuZnVuY3Rpb24gaXNIb3N0TGlzdGVuZXIodmFsdWU6IGFueSk6IHZhbHVlIGlzIEhvc3RMaXN0ZW5lciB7XG4gIHJldHVybiB2YWx1ZS5uZ01ldGFkYXRhTmFtZSA9PT0gJ0hvc3RMaXN0ZW5lcic7XG59XG5cbmZ1bmN0aW9uIGlzSW5wdXQodmFsdWU6IGFueSk6IHZhbHVlIGlzIElucHV0IHtcbiAgcmV0dXJuIHZhbHVlLm5nTWV0YWRhdGFOYW1lID09PSAnSW5wdXQnO1xufVxuXG5mdW5jdGlvbiBpc091dHB1dCh2YWx1ZTogYW55KTogdmFsdWUgaXMgT3V0cHV0IHtcbiAgcmV0dXJuIHZhbHVlLm5nTWV0YWRhdGFOYW1lID09PSAnT3V0cHV0Jztcbn1cblxuZnVuY3Rpb24gaW5wdXRzUGFydGlhbE1ldGFkYXRhVG9JbnB1dE1ldGFkYXRhKFxuICBpbnB1dHM6IE5vbk51bGxhYmxlPFIzRGVjbGFyZURpcmVjdGl2ZUZhY2FkZVsnaW5wdXRzJ10+LFxuKSB7XG4gIHJldHVybiBPYmplY3Qua2V5cyhpbnB1dHMpLnJlZHVjZTxSZWNvcmQ8c3RyaW5nLCBSM0lucHV0TWV0YWRhdGE+PihcbiAgICAocmVzdWx0LCBtaW5pZmllZENsYXNzTmFtZSkgPT4ge1xuICAgICAgY29uc3QgdmFsdWUgPSBpbnB1dHNbbWluaWZpZWRDbGFzc05hbWVdO1xuXG4gICAgICAvLyBIYW5kbGUgbGVnYWN5IHBhcnRpYWwgaW5wdXQgb3V0cHV0LlxuICAgICAgaWYgKHR5cGVvZiB2YWx1ZSA9PT0gJ3N0cmluZycgfHwgQXJyYXkuaXNBcnJheSh2YWx1ZSkpIHtcbiAgICAgICAgcmVzdWx0W21pbmlmaWVkQ2xhc3NOYW1lXSA9IHBhcnNlTGVnYWN5SW5wdXRQYXJ0aWFsT3V0cHV0KHZhbHVlKTtcbiAgICAgIH0gZWxzZSB7XG4gICAgICAgIHJlc3VsdFttaW5pZmllZENsYXNzTmFtZV0gPSB7XG4gICAgICAgICAgYmluZGluZ1Byb3BlcnR5TmFtZTogdmFsdWUucHVibGljTmFtZSxcbiAgICAgICAgICBjbGFzc1Byb3BlcnR5TmFtZTogbWluaWZpZWRDbGFzc05hbWUsXG4gICAgICAgICAgdHJhbnNmb3JtRnVuY3Rpb246XG4gICAgICAgICAgICB2YWx1ZS50cmFuc2Zvcm1GdW5jdGlvbiAhPT0gbnVsbCA/IG5ldyBXcmFwcGVkTm9kZUV4cHIodmFsdWUudHJhbnNmb3JtRnVuY3Rpb24pIDogbnVsbCxcbiAgICAgICAgICByZXF1aXJlZDogdmFsdWUuaXNSZXF1aXJlZCxcbiAgICAgICAgICBpc1NpZ25hbDogdmFsdWUuaXNTaWduYWwsXG4gICAgICAgIH07XG4gICAgICB9XG5cbiAgICAgIHJldHVybiByZXN1bHQ7XG4gICAgfSxcbiAgICB7fSxcbiAgKTtcbn1cblxuLyoqXG4gKiBQYXJzZXMgdGhlIGxlZ2FjeSBpbnB1dCBwYXJ0aWFsIG91dHB1dC4gRm9yIG1vcmUgZGV0YWlscyBzZWUgYHBhcnRpYWwvZGlyZWN0aXZlLnRzYC5cbiAqIFRPRE8obGVnYWN5LXBhcnRpYWwtb3V0cHV0LWlucHV0cyk6IFJlbW92ZSBpbiB2MTguXG4gKi9cbmZ1bmN0aW9uIHBhcnNlTGVnYWN5SW5wdXRQYXJ0aWFsT3V0cHV0KHZhbHVlOiBMZWdhY3lJbnB1dFBhcnRpYWxNYXBwaW5nKTogUjNJbnB1dE1ldGFkYXRhIHtcbiAgaWYgKHR5cGVvZiB2YWx1ZSA9PT0gJ3N0cmluZycpIHtcbiAgICByZXR1cm4ge1xuICAgICAgYmluZGluZ1Byb3BlcnR5TmFtZTogdmFsdWUsXG4gICAgICBjbGFzc1Byb3BlcnR5TmFtZTogdmFsdWUsXG4gICAgICB0cmFuc2Zvcm1GdW5jdGlvbjogbnVsbCxcbiAgICAgIHJlcXVpcmVkOiBmYWxzZSxcbiAgICAgIC8vIGxlZ2FjeSBwYXJ0aWFsIG91dHB1dCBkb2VzIG5vdCBjYXB0dXJlIHNpZ25hbCBpbnB1dHMuXG4gICAgICBpc1NpZ25hbDogZmFsc2UsXG4gICAgfTtcbiAgfVxuXG4gIHJldHVybiB7XG4gICAgYmluZGluZ1Byb3BlcnR5TmFtZTogdmFsdWVbMF0sXG4gICAgY2xhc3NQcm9wZXJ0eU5hbWU6IHZhbHVlWzFdLFxuICAgIHRyYW5zZm9ybUZ1bmN0aW9uOiB2YWx1ZVsyXSA/IG5ldyBXcmFwcGVkTm9kZUV4cHIodmFsdWVbMl0pIDogbnVsbCxcbiAgICByZXF1aXJlZDogZmFsc2UsXG4gICAgLy8gbGVnYWN5IHBhcnRpYWwgb3V0cHV0IGRvZXMgbm90IGNhcHR1cmUgc2lnbmFsIGlucHV0cy5cbiAgICBpc1NpZ25hbDogZmFsc2UsXG4gIH07XG59XG5cbmZ1bmN0aW9uIHBhcnNlSW5wdXRzQXJyYXkoXG4gIHZhbHVlczogKHN0cmluZyB8IHtuYW1lOiBzdHJpbmc7IGFsaWFzPzogc3RyaW5nOyByZXF1aXJlZD86IGJvb2xlYW47IHRyYW5zZm9ybT86IEZ1bmN0aW9ufSlbXSxcbikge1xuICByZXR1cm4gdmFsdWVzLnJlZHVjZTxSZWNvcmQ8c3RyaW5nLCBSM0lucHV0TWV0YWRhdGE+PigocmVzdWx0cywgdmFsdWUpID0+IHtcbiAgICBpZiAodHlwZW9mIHZhbHVlID09PSAnc3RyaW5nJykge1xuICAgICAgY29uc3QgW2JpbmRpbmdQcm9wZXJ0eU5hbWUsIGNsYXNzUHJvcGVydHlOYW1lXSA9IHBhcnNlTWFwcGluZ1N0cmluZyh2YWx1ZSk7XG4gICAgICByZXN1bHRzW2NsYXNzUHJvcGVydHlOYW1lXSA9IHtcbiAgICAgICAgYmluZGluZ1Byb3BlcnR5TmFtZSxcbiAgICAgICAgY2xhc3NQcm9wZXJ0eU5hbWUsXG4gICAgICAgIHJlcXVpcmVkOiBmYWxzZSxcbiAgICAgICAgLy8gU2lnbmFsIGlucHV0cyBub3Qgc3VwcG9ydGVkIGZvciB0aGUgaW5wdXRzIGFycmF5LlxuICAgICAgICBpc1NpZ25hbDogZmFsc2UsXG4gICAgICAgIHRyYW5zZm9ybUZ1bmN0aW9uOiBudWxsLFxuICAgICAgfTtcbiAgICB9IGVsc2Uge1xuICAgICAgcmVzdWx0c1t2YWx1ZS5uYW1lXSA9IHtcbiAgICAgICAgYmluZGluZ1Byb3BlcnR5TmFtZTogdmFsdWUuYWxpYXMgfHwgdmFsdWUubmFtZSxcbiAgICAgICAgY2xhc3NQcm9wZXJ0eU5hbWU6IHZhbHVlLm5hbWUsXG4gICAgICAgIHJlcXVpcmVkOiB2YWx1ZS5yZXF1aXJlZCB8fCBmYWxzZSxcbiAgICAgICAgLy8gU2lnbmFsIGlucHV0cyBub3Qgc3VwcG9ydGVkIGZvciB0aGUgaW5wdXRzIGFycmF5LlxuICAgICAgICBpc1NpZ25hbDogZmFsc2UsXG4gICAgICAgIHRyYW5zZm9ybUZ1bmN0aW9uOiB2YWx1ZS50cmFuc2Zvcm0gIT0gbnVsbCA/IG5ldyBXcmFwcGVkTm9kZUV4cHIodmFsdWUudHJhbnNmb3JtKSA6IG51bGwsXG4gICAgICB9O1xuICAgIH1cbiAgICByZXR1cm4gcmVzdWx0cztcbiAgfSwge30pO1xufVxuXG5mdW5jdGlvbiBwYXJzZU1hcHBpbmdTdHJpbmdBcnJheSh2YWx1ZXM6IHN0cmluZ1tdKTogUmVjb3JkPHN0cmluZywgc3RyaW5nPiB7XG4gIHJldHVybiB2YWx1ZXMucmVkdWNlPFJlY29yZDxzdHJpbmcsIHN0cmluZz4+KChyZXN1bHRzLCB2YWx1ZSkgPT4ge1xuICAgIGNvbnN0IFthbGlhcywgZmllbGROYW1lXSA9IHBhcnNlTWFwcGluZ1N0cmluZyh2YWx1ZSk7XG4gICAgcmVzdWx0c1tmaWVsZE5hbWVdID0gYWxpYXM7XG4gICAgcmV0dXJuIHJlc3VsdHM7XG4gIH0sIHt9KTtcbn1cblxuZnVuY3Rpb24gcGFyc2VNYXBwaW5nU3RyaW5nKHZhbHVlOiBzdHJpbmcpOiBbYWxpYXM6IHN0cmluZywgZmllbGROYW1lOiBzdHJpbmddIHtcbiAgLy8gRWl0aGVyIHRoZSB2YWx1ZSBpcyAnZmllbGQnIG9yICdmaWVsZDogcHJvcGVydHknLiBJbiB0aGUgZmlyc3QgY2FzZSwgYHByb3BlcnR5YCB3aWxsXG4gIC8vIGJlIHVuZGVmaW5lZCwgaW4gd2hpY2ggY2FzZSB0aGUgZmllbGQgbmFtZSBzaG91bGQgYWxzbyBiZSB1c2VkIGFzIHRoZSBwcm9wZXJ0eSBuYW1lLlxuICBjb25zdCBbZmllbGROYW1lLCBiaW5kaW5nUHJvcGVydHlOYW1lXSA9IHZhbHVlLnNwbGl0KCc6JywgMikubWFwKChzdHIpID0+IHN0ci50cmltKCkpO1xuICByZXR1cm4gW2JpbmRpbmdQcm9wZXJ0eU5hbWUgPz8gZmllbGROYW1lLCBmaWVsZE5hbWVdO1xufVxuXG5mdW5jdGlvbiBjb252ZXJ0RGVjbGFyZVBpcGVGYWNhZGVUb01ldGFkYXRhKGRlY2xhcmF0aW9uOiBSM0RlY2xhcmVQaXBlRmFjYWRlKTogUjNQaXBlTWV0YWRhdGEge1xuICByZXR1cm4ge1xuICAgIG5hbWU6IGRlY2xhcmF0aW9uLnR5cGUubmFtZSxcbiAgICB0eXBlOiB3cmFwUmVmZXJlbmNlKGRlY2xhcmF0aW9uLnR5cGUpLFxuICAgIHR5cGVBcmd1bWVudENvdW50OiAwLFxuICAgIHBpcGVOYW1lOiBkZWNsYXJhdGlvbi5uYW1lLFxuICAgIGRlcHM6IG51bGwsXG4gICAgcHVyZTogZGVjbGFyYXRpb24ucHVyZSA/PyB0cnVlLFxuICAgIGlzU3RhbmRhbG9uZTogZGVjbGFyYXRpb24uaXNTdGFuZGFsb25lID8/IGZhbHNlLFxuICB9O1xufVxuXG5mdW5jdGlvbiBjb252ZXJ0RGVjbGFyZUluamVjdG9yRmFjYWRlVG9NZXRhZGF0YShcbiAgZGVjbGFyYXRpb246IFIzRGVjbGFyZUluamVjdG9yRmFjYWRlLFxuKTogUjNJbmplY3Rvck1ldGFkYXRhIHtcbiAgcmV0dXJuIHtcbiAgICBuYW1lOiBkZWNsYXJhdGlvbi50eXBlLm5hbWUsXG4gICAgdHlwZTogd3JhcFJlZmVyZW5jZShkZWNsYXJhdGlvbi50eXBlKSxcbiAgICBwcm92aWRlcnM6XG4gICAgICBkZWNsYXJhdGlvbi5wcm92aWRlcnMgIT09IHVuZGVmaW5lZCAmJiBkZWNsYXJhdGlvbi5wcm92aWRlcnMubGVuZ3RoID4gMFxuICAgICAgICA/IG5ldyBXcmFwcGVkTm9kZUV4cHIoZGVjbGFyYXRpb24ucHJvdmlkZXJzKVxuICAgICAgICA6IG51bGwsXG4gICAgaW1wb3J0czpcbiAgICAgIGRlY2xhcmF0aW9uLmltcG9ydHMgIT09IHVuZGVmaW5lZFxuICAgICAgICA/IGRlY2xhcmF0aW9uLmltcG9ydHMubWFwKChpKSA9PiBuZXcgV3JhcHBlZE5vZGVFeHByKGkpKVxuICAgICAgICA6IFtdLFxuICB9O1xufVxuXG5leHBvcnQgZnVuY3Rpb24gcHVibGlzaEZhY2FkZShnbG9iYWw6IGFueSkge1xuICBjb25zdCBuZzogRXhwb3J0ZWRDb21waWxlckZhY2FkZSA9IGdsb2JhbC5uZyB8fCAoZ2xvYmFsLm5nID0ge30pO1xuICBuZy7JtWNvbXBpbGVyRmFjYWRlID0gbmV3IENvbXBpbGVyRmFjYWRlSW1wbCgpO1xufVxuIl19