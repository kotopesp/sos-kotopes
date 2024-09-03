/**
 * @license
 * Copyright Google LLC All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
/**
 * Distinguishes different kinds of IR operations.
 *
 * Includes both creation and update operations.
 */
export var OpKind;
(function (OpKind) {
    /**
     * A special operation type which is used to represent the beginning and end nodes of a linked
     * list of operations.
     */
    OpKind[OpKind["ListEnd"] = 0] = "ListEnd";
    /**
     * An operation which wraps an output AST statement.
     */
    OpKind[OpKind["Statement"] = 1] = "Statement";
    /**
     * An operation which declares and initializes a `SemanticVariable`.
     */
    OpKind[OpKind["Variable"] = 2] = "Variable";
    /**
     * An operation to begin rendering of an element.
     */
    OpKind[OpKind["ElementStart"] = 3] = "ElementStart";
    /**
     * An operation to render an element with no children.
     */
    OpKind[OpKind["Element"] = 4] = "Element";
    /**
     * An operation which declares an embedded view.
     */
    OpKind[OpKind["Template"] = 5] = "Template";
    /**
     * An operation to end rendering of an element previously started with `ElementStart`.
     */
    OpKind[OpKind["ElementEnd"] = 6] = "ElementEnd";
    /**
     * An operation to begin an `ng-container`.
     */
    OpKind[OpKind["ContainerStart"] = 7] = "ContainerStart";
    /**
     * An operation for an `ng-container` with no children.
     */
    OpKind[OpKind["Container"] = 8] = "Container";
    /**
     * An operation to end an `ng-container`.
     */
    OpKind[OpKind["ContainerEnd"] = 9] = "ContainerEnd";
    /**
     * An operation disable binding for subsequent elements, which are descendants of a non-bindable
     * node.
     */
    OpKind[OpKind["DisableBindings"] = 10] = "DisableBindings";
    /**
     * An op to conditionally render a template.
     */
    OpKind[OpKind["Conditional"] = 11] = "Conditional";
    /**
     * An operation to re-enable binding, after it was previously disabled.
     */
    OpKind[OpKind["EnableBindings"] = 12] = "EnableBindings";
    /**
     * An operation to render a text node.
     */
    OpKind[OpKind["Text"] = 13] = "Text";
    /**
     * An operation declaring an event listener for an element.
     */
    OpKind[OpKind["Listener"] = 14] = "Listener";
    /**
     * An operation to interpolate text into a text node.
     */
    OpKind[OpKind["InterpolateText"] = 15] = "InterpolateText";
    /**
     * An intermediate binding op, that has not yet been processed into an individual property,
     * attribute, style, etc.
     */
    OpKind[OpKind["Binding"] = 16] = "Binding";
    /**
     * An operation to bind an expression to a property of an element.
     */
    OpKind[OpKind["Property"] = 17] = "Property";
    /**
     * An operation to bind an expression to a style property of an element.
     */
    OpKind[OpKind["StyleProp"] = 18] = "StyleProp";
    /**
     * An operation to bind an expression to a class property of an element.
     */
    OpKind[OpKind["ClassProp"] = 19] = "ClassProp";
    /**
     * An operation to bind an expression to the styles of an element.
     */
    OpKind[OpKind["StyleMap"] = 20] = "StyleMap";
    /**
     * An operation to bind an expression to the classes of an element.
     */
    OpKind[OpKind["ClassMap"] = 21] = "ClassMap";
    /**
     * An operation to advance the runtime's implicit slot context during the update phase of a view.
     */
    OpKind[OpKind["Advance"] = 22] = "Advance";
    /**
     * An operation to instantiate a pipe.
     */
    OpKind[OpKind["Pipe"] = 23] = "Pipe";
    /**
     * An operation to associate an attribute with an element.
     */
    OpKind[OpKind["Attribute"] = 24] = "Attribute";
    /**
     * An attribute that has been extracted for inclusion in the consts array.
     */
    OpKind[OpKind["ExtractedAttribute"] = 25] = "ExtractedAttribute";
    /**
     * An operation that configures a `@defer` block.
     */
    OpKind[OpKind["Defer"] = 26] = "Defer";
    /**
     * An operation that controls when a `@defer` loads.
     */
    OpKind[OpKind["DeferOn"] = 27] = "DeferOn";
    /**
     * An operation that controls when a `@defer` loads, using a custom expression as the condition.
     */
    OpKind[OpKind["DeferWhen"] = 28] = "DeferWhen";
    /**
     * An i18n message that has been extracted for inclusion in the consts array.
     */
    OpKind[OpKind["I18nMessage"] = 29] = "I18nMessage";
    /**
     * A host binding property.
     */
    OpKind[OpKind["HostProperty"] = 30] = "HostProperty";
    /**
     * A namespace change, which causes the subsequent elements to be processed as either HTML or SVG.
     */
    OpKind[OpKind["Namespace"] = 31] = "Namespace";
    /**
     * Configure a content projeciton definition for the view.
     */
    OpKind[OpKind["ProjectionDef"] = 32] = "ProjectionDef";
    /**
     * Create a content projection slot.
     */
    OpKind[OpKind["Projection"] = 33] = "Projection";
    /**
     * Create a repeater creation instruction op.
     */
    OpKind[OpKind["RepeaterCreate"] = 34] = "RepeaterCreate";
    /**
     * An update up for a repeater.
     */
    OpKind[OpKind["Repeater"] = 35] = "Repeater";
    /**
     * An operation to bind an expression to the property side of a two-way binding.
     */
    OpKind[OpKind["TwoWayProperty"] = 36] = "TwoWayProperty";
    /**
     * An operation declaring the event side of a two-way binding.
     */
    OpKind[OpKind["TwoWayListener"] = 37] = "TwoWayListener";
    /**
     * A creation-time operation that initializes the slot for a `@let` declaration.
     */
    OpKind[OpKind["DeclareLet"] = 38] = "DeclareLet";
    /**
     * An update-time operation that stores the current value of a `@let` declaration.
     */
    OpKind[OpKind["StoreLet"] = 39] = "StoreLet";
    /**
     * The start of an i18n block.
     */
    OpKind[OpKind["I18nStart"] = 40] = "I18nStart";
    /**
     * A self-closing i18n on a single element.
     */
    OpKind[OpKind["I18n"] = 41] = "I18n";
    /**
     * The end of an i18n block.
     */
    OpKind[OpKind["I18nEnd"] = 42] = "I18nEnd";
    /**
     * An expression in an i18n message.
     */
    OpKind[OpKind["I18nExpression"] = 43] = "I18nExpression";
    /**
     * An instruction that applies a set of i18n expressions.
     */
    OpKind[OpKind["I18nApply"] = 44] = "I18nApply";
    /**
     * An instruction to create an ICU expression.
     */
    OpKind[OpKind["IcuStart"] = 45] = "IcuStart";
    /**
     * An instruction to update an ICU expression.
     */
    OpKind[OpKind["IcuEnd"] = 46] = "IcuEnd";
    /**
     * An instruction representing a placeholder in an ICU expression.
     */
    OpKind[OpKind["IcuPlaceholder"] = 47] = "IcuPlaceholder";
    /**
     * An i18n context containing information needed to generate an i18n message.
     */
    OpKind[OpKind["I18nContext"] = 48] = "I18nContext";
    /**
     * A creation op that corresponds to i18n attributes on an element.
     */
    OpKind[OpKind["I18nAttributes"] = 49] = "I18nAttributes";
})(OpKind || (OpKind = {}));
/**
 * Distinguishes different kinds of IR expressions.
 */
export var ExpressionKind;
(function (ExpressionKind) {
    /**
     * Read of a variable in a lexical scope.
     */
    ExpressionKind[ExpressionKind["LexicalRead"] = 0] = "LexicalRead";
    /**
     * A reference to the current view context.
     */
    ExpressionKind[ExpressionKind["Context"] = 1] = "Context";
    /**
     * A reference to the view context, for use inside a track function.
     */
    ExpressionKind[ExpressionKind["TrackContext"] = 2] = "TrackContext";
    /**
     * Read of a variable declared in a `VariableOp`.
     */
    ExpressionKind[ExpressionKind["ReadVariable"] = 3] = "ReadVariable";
    /**
     * Runtime operation to navigate to the next view context in the view hierarchy.
     */
    ExpressionKind[ExpressionKind["NextContext"] = 4] = "NextContext";
    /**
     * Runtime operation to retrieve the value of a local reference.
     */
    ExpressionKind[ExpressionKind["Reference"] = 5] = "Reference";
    /**
     * A call storing the value of a `@let` declaration.
     */
    ExpressionKind[ExpressionKind["StoreLet"] = 6] = "StoreLet";
    /**
     * A reference to a `@let` declaration read from the context view.
     */
    ExpressionKind[ExpressionKind["ContextLetReference"] = 7] = "ContextLetReference";
    /**
     * Runtime operation to snapshot the current view context.
     */
    ExpressionKind[ExpressionKind["GetCurrentView"] = 8] = "GetCurrentView";
    /**
     * Runtime operation to restore a snapshotted view.
     */
    ExpressionKind[ExpressionKind["RestoreView"] = 9] = "RestoreView";
    /**
     * Runtime operation to reset the current view context after `RestoreView`.
     */
    ExpressionKind[ExpressionKind["ResetView"] = 10] = "ResetView";
    /**
     * Defines and calls a function with change-detected arguments.
     */
    ExpressionKind[ExpressionKind["PureFunctionExpr"] = 11] = "PureFunctionExpr";
    /**
     * Indicates a positional parameter to a pure function definition.
     */
    ExpressionKind[ExpressionKind["PureFunctionParameterExpr"] = 12] = "PureFunctionParameterExpr";
    /**
     * Binding to a pipe transformation.
     */
    ExpressionKind[ExpressionKind["PipeBinding"] = 13] = "PipeBinding";
    /**
     * Binding to a pipe transformation with a variable number of arguments.
     */
    ExpressionKind[ExpressionKind["PipeBindingVariadic"] = 14] = "PipeBindingVariadic";
    /*
     * A safe property read requiring expansion into a null check.
     */
    ExpressionKind[ExpressionKind["SafePropertyRead"] = 15] = "SafePropertyRead";
    /**
     * A safe keyed read requiring expansion into a null check.
     */
    ExpressionKind[ExpressionKind["SafeKeyedRead"] = 16] = "SafeKeyedRead";
    /**
     * A safe function call requiring expansion into a null check.
     */
    ExpressionKind[ExpressionKind["SafeInvokeFunction"] = 17] = "SafeInvokeFunction";
    /**
     * An intermediate expression that will be expanded from a safe read into an explicit ternary.
     */
    ExpressionKind[ExpressionKind["SafeTernaryExpr"] = 18] = "SafeTernaryExpr";
    /**
     * An empty expression that will be stipped before generating the final output.
     */
    ExpressionKind[ExpressionKind["EmptyExpr"] = 19] = "EmptyExpr";
    /*
     * An assignment to a temporary variable.
     */
    ExpressionKind[ExpressionKind["AssignTemporaryExpr"] = 20] = "AssignTemporaryExpr";
    /**
     * A reference to a temporary variable.
     */
    ExpressionKind[ExpressionKind["ReadTemporaryExpr"] = 21] = "ReadTemporaryExpr";
    /**
     * An expression that will cause a literal slot index to be emitted.
     */
    ExpressionKind[ExpressionKind["SlotLiteralExpr"] = 22] = "SlotLiteralExpr";
    /**
     * A test expression for a conditional op.
     */
    ExpressionKind[ExpressionKind["ConditionalCase"] = 23] = "ConditionalCase";
    /**
     * An expression that will be automatically extracted to the component const array.
     */
    ExpressionKind[ExpressionKind["ConstCollected"] = 24] = "ConstCollected";
    /**
     * Operation that sets the value of a two-way binding.
     */
    ExpressionKind[ExpressionKind["TwoWayBindingSet"] = 25] = "TwoWayBindingSet";
})(ExpressionKind || (ExpressionKind = {}));
export var VariableFlags;
(function (VariableFlags) {
    VariableFlags[VariableFlags["None"] = 0] = "None";
    /**
     * Always inline this variable, regardless of the number of times it's used.
     * An `AlwaysInline` variable may not depend on context, because doing so may cause side effects
     * that are illegal when multi-inlined. (The optimizer will enforce this constraint.)
     */
    VariableFlags[VariableFlags["AlwaysInline"] = 1] = "AlwaysInline";
})(VariableFlags || (VariableFlags = {}));
/**
 * Distinguishes between different kinds of `SemanticVariable`s.
 */
export var SemanticVariableKind;
(function (SemanticVariableKind) {
    /**
     * Represents the context of a particular view.
     */
    SemanticVariableKind[SemanticVariableKind["Context"] = 0] = "Context";
    /**
     * Represents an identifier declared in the lexical scope of a view.
     */
    SemanticVariableKind[SemanticVariableKind["Identifier"] = 1] = "Identifier";
    /**
     * Represents a saved state that can be used to restore a view in a listener handler function.
     */
    SemanticVariableKind[SemanticVariableKind["SavedView"] = 2] = "SavedView";
    /**
     * An alias generated by a special embedded view type (e.g. a `@for` block).
     */
    SemanticVariableKind[SemanticVariableKind["Alias"] = 3] = "Alias";
})(SemanticVariableKind || (SemanticVariableKind = {}));
/**
 * Whether to compile in compatibilty mode. In compatibility mode, the template pipeline will
 * attempt to match the output of `TemplateDefinitionBuilder` as exactly as possible, at the cost
 * of producing quirky or larger code in some cases.
 */
export var CompatibilityMode;
(function (CompatibilityMode) {
    CompatibilityMode[CompatibilityMode["Normal"] = 0] = "Normal";
    CompatibilityMode[CompatibilityMode["TemplateDefinitionBuilder"] = 1] = "TemplateDefinitionBuilder";
})(CompatibilityMode || (CompatibilityMode = {}));
/**
 * Enumeration of the types of attributes which can be applied to an element.
 */
export var BindingKind;
(function (BindingKind) {
    /**
     * Static attributes.
     */
    BindingKind[BindingKind["Attribute"] = 0] = "Attribute";
    /**
     * Class bindings.
     */
    BindingKind[BindingKind["ClassName"] = 1] = "ClassName";
    /**
     * Style bindings.
     */
    BindingKind[BindingKind["StyleProperty"] = 2] = "StyleProperty";
    /**
     * Dynamic property bindings.
     */
    BindingKind[BindingKind["Property"] = 3] = "Property";
    /**
     * Property or attribute bindings on a template.
     */
    BindingKind[BindingKind["Template"] = 4] = "Template";
    /**
     * Internationalized attributes.
     */
    BindingKind[BindingKind["I18n"] = 5] = "I18n";
    /**
     * Animation property bindings.
     */
    BindingKind[BindingKind["Animation"] = 6] = "Animation";
    /**
     * Property side of a two-way binding.
     */
    BindingKind[BindingKind["TwoWayProperty"] = 7] = "TwoWayProperty";
})(BindingKind || (BindingKind = {}));
/**
 * Enumeration of possible times i18n params can be resolved.
 */
export var I18nParamResolutionTime;
(function (I18nParamResolutionTime) {
    /**
     * Param is resolved at message creation time. Most params should be resolved at message creation
     * time. However, ICU params need to be handled in post-processing.
     */
    I18nParamResolutionTime[I18nParamResolutionTime["Creation"] = 0] = "Creation";
    /**
     * Param is resolved during post-processing. This should be used for params whose value comes from
     * an ICU.
     */
    I18nParamResolutionTime[I18nParamResolutionTime["Postproccessing"] = 1] = "Postproccessing";
})(I18nParamResolutionTime || (I18nParamResolutionTime = {}));
/**
 * The contexts in which an i18n expression can be used.
 */
export var I18nExpressionFor;
(function (I18nExpressionFor) {
    /**
     * This expression is used as a value (i.e. inside an i18n block).
     */
    I18nExpressionFor[I18nExpressionFor["I18nText"] = 0] = "I18nText";
    /**
     * This expression is used in a binding.
     */
    I18nExpressionFor[I18nExpressionFor["I18nAttribute"] = 1] = "I18nAttribute";
})(I18nExpressionFor || (I18nExpressionFor = {}));
/**
 * Flags that describe what an i18n param value. These determine how the value is serialized into
 * the final map.
 */
export var I18nParamValueFlags;
(function (I18nParamValueFlags) {
    I18nParamValueFlags[I18nParamValueFlags["None"] = 0] = "None";
    /**
     *  This value represents an element tag.
     */
    I18nParamValueFlags[I18nParamValueFlags["ElementTag"] = 1] = "ElementTag";
    /**
     * This value represents a template tag.
     */
    I18nParamValueFlags[I18nParamValueFlags["TemplateTag"] = 2] = "TemplateTag";
    /**
     * This value represents the opening of a tag.
     */
    I18nParamValueFlags[I18nParamValueFlags["OpenTag"] = 4] = "OpenTag";
    /**
     * This value represents the closing of a tag.
     */
    I18nParamValueFlags[I18nParamValueFlags["CloseTag"] = 8] = "CloseTag";
    /**
     * This value represents an i18n expression index.
     */
    I18nParamValueFlags[I18nParamValueFlags["ExpressionIndex"] = 16] = "ExpressionIndex";
})(I18nParamValueFlags || (I18nParamValueFlags = {}));
/**
 * Whether the active namespace is HTML, MathML, or SVG mode.
 */
export var Namespace;
(function (Namespace) {
    Namespace[Namespace["HTML"] = 0] = "HTML";
    Namespace[Namespace["SVG"] = 1] = "SVG";
    Namespace[Namespace["Math"] = 2] = "Math";
})(Namespace || (Namespace = {}));
/**
 * The type of a `@defer` trigger, for use in the ir.
 */
export var DeferTriggerKind;
(function (DeferTriggerKind) {
    DeferTriggerKind[DeferTriggerKind["Idle"] = 0] = "Idle";
    DeferTriggerKind[DeferTriggerKind["Immediate"] = 1] = "Immediate";
    DeferTriggerKind[DeferTriggerKind["Timer"] = 2] = "Timer";
    DeferTriggerKind[DeferTriggerKind["Hover"] = 3] = "Hover";
    DeferTriggerKind[DeferTriggerKind["Interaction"] = 4] = "Interaction";
    DeferTriggerKind[DeferTriggerKind["Viewport"] = 5] = "Viewport";
})(DeferTriggerKind || (DeferTriggerKind = {}));
/**
 * Kinds of i18n contexts. They can be created because of root i18n blocks, or ICUs.
 */
export var I18nContextKind;
(function (I18nContextKind) {
    I18nContextKind[I18nContextKind["RootI18n"] = 0] = "RootI18n";
    I18nContextKind[I18nContextKind["Icu"] = 1] = "Icu";
    I18nContextKind[I18nContextKind["Attr"] = 2] = "Attr";
})(I18nContextKind || (I18nContextKind = {}));
export var TemplateKind;
(function (TemplateKind) {
    TemplateKind[TemplateKind["NgTemplate"] = 0] = "NgTemplate";
    TemplateKind[TemplateKind["Structural"] = 1] = "Structural";
    TemplateKind[TemplateKind["Block"] = 2] = "Block";
})(TemplateKind || (TemplateKind = {}));
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZW51bXMuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi9wYWNrYWdlcy9jb21waWxlci9zcmMvdGVtcGxhdGUvcGlwZWxpbmUvaXIvc3JjL2VudW1zLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBOzs7Ozs7R0FNRztBQUVIOzs7O0dBSUc7QUFDSCxNQUFNLENBQU4sSUFBWSxNQTZQWDtBQTdQRCxXQUFZLE1BQU07SUFDaEI7OztPQUdHO0lBQ0gseUNBQU8sQ0FBQTtJQUVQOztPQUVHO0lBQ0gsNkNBQVMsQ0FBQTtJQUVUOztPQUVHO0lBQ0gsMkNBQVEsQ0FBQTtJQUVSOztPQUVHO0lBQ0gsbURBQVksQ0FBQTtJQUVaOztPQUVHO0lBQ0gseUNBQU8sQ0FBQTtJQUVQOztPQUVHO0lBQ0gsMkNBQVEsQ0FBQTtJQUVSOztPQUVHO0lBQ0gsK0NBQVUsQ0FBQTtJQUVWOztPQUVHO0lBQ0gsdURBQWMsQ0FBQTtJQUVkOztPQUVHO0lBQ0gsNkNBQVMsQ0FBQTtJQUVUOztPQUVHO0lBQ0gsbURBQVksQ0FBQTtJQUVaOzs7T0FHRztJQUNILDBEQUFlLENBQUE7SUFFZjs7T0FFRztJQUNILGtEQUFXLENBQUE7SUFFWDs7T0FFRztJQUNILHdEQUFjLENBQUE7SUFFZDs7T0FFRztJQUNILG9DQUFJLENBQUE7SUFFSjs7T0FFRztJQUNILDRDQUFRLENBQUE7SUFFUjs7T0FFRztJQUNILDBEQUFlLENBQUE7SUFFZjs7O09BR0c7SUFDSCwwQ0FBTyxDQUFBO0lBRVA7O09BRUc7SUFDSCw0Q0FBUSxDQUFBO0lBRVI7O09BRUc7SUFDSCw4Q0FBUyxDQUFBO0lBRVQ7O09BRUc7SUFDSCw4Q0FBUyxDQUFBO0lBRVQ7O09BRUc7SUFDSCw0Q0FBUSxDQUFBO0lBRVI7O09BRUc7SUFDSCw0Q0FBUSxDQUFBO0lBRVI7O09BRUc7SUFDSCwwQ0FBTyxDQUFBO0lBRVA7O09BRUc7SUFDSCxvQ0FBSSxDQUFBO0lBRUo7O09B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Yjs7T0FFRztJQUNILHlFQUFnQixDQUFBO0lBRWhCOztPQUVHO0lBQ0gsMkVBQWtCLENBQUE7SUFFbEI7O09BRUc7SUFDSCxtRUFBZ0IsQ0FBQTtJQUVoQjs7T0FFRztJQUNILHFFQUFpQixDQUFBO0lBRWpCOztPQUVHO0lBQ0gsb0ZBQXlCLENBQUE7QUFDM0IsQ0FBQyxFQTNCVyxtQkFBbUIsS0FBbkIsbUJBQW1CLFFBMkI5QjtBQUVEOztHQUVHO0FBQ0gsTUFBTSxDQUFOLElBQVksU0FJWDtBQUpELFdBQVksU0FBUztJQUNuQix5Q0FBSSxDQUFBO0lBQ0osdUNBQUcsQ0FBQTtJQUNILHlDQUFJLENBQUE7QUFDTixDQUFDLEVBSlcsU0FBUyxLQUFULFNBQVMsUUFJcEI7QUFFRDs7R0FFRztBQUNILE1BQU0sQ0FBTixJQUFZLGdCQU9YO0FBUEQsV0FBWSxnQkFBZ0I7SUFDMUIsdURBQUksQ0FBQTtJQUNKLGlFQUFTLENBQUE7SUFDVCx5REFBSyxDQUFBO0lBQ0wseURBQUssQ0FBQTtJQUNMLHFFQUFXLENBQUE7SUFDWCwrREFBUSxDQUFBO0FBQ1YsQ0FBQyxFQVBXLGdCQUFnQixLQUFoQixnQkFBZ0IsUUFPM0I7QUFFRDs7R0FFRztBQUNILE1BQU0sQ0FBTixJQUFZLGVBSVg7QUFKRCxXQUFZLGVBQWU7SUFDekIsNkRBQVEsQ0FBQTtJQUNSLG1EQUFHLENBQUE7SUFDSCxxREFBSSxDQUFBO0FBQ04sQ0FBQyxFQUpXLGVBQWUsS0FBZixlQUFlLFFBSTFCO0FBRUQsTUFBTSxDQUFOLElBQVksWUFJWDtBQUpELFdBQVksWUFBWTtJQUN0QiwyREFBVSxDQUFBO0lBQ1YsMkRBQVUsQ0FBQTtJQUNWLGlEQUFLLENBQUE7QUFDUCxDQUFDLEVBSlcsWUFBWSxLQUFaLFlBQVksUUFJdkIiLCJzb3VyY2VzQ29udGVudCI6WyIvKipcbiAqIEBsaWNlbnNlXG4gKiBDb3B5cmlnaHQgR29vZ2xlIExMQyBBbGwgUmlnaHRzIFJlc2VydmVkLlxuICpcbiAqIFVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlIGlzIGdvdmVybmVkIGJ5IGFuIE1JVC1zdHlsZSBsaWNlbnNlIHRoYXQgY2FuIGJlXG4gKiBmb3VuZCBpbiB0aGUgTElDRU5TRSBmaWxlIGF0IGh0dHBzOi8vYW5ndWxhci5pby9saWNlbnNlXG4gKi9cblxuLyoqXG4gKiBEaXN0aW5ndWlzaGVzIGRpZmZlcmVudCBraW5kcyBvZiBJUiBvcGVyYXRpb25zLlxuICpcbiAqIEluY2x1ZGVzIGJvdGggY3JlYXRpb24gYW5kIHVwZGF0ZSBvcGVyYXRpb25zLlxuICovXG5leHBvcnQgZW51bSBPcEtpbmQge1xuICAvKipcbiAgICogQSBzcGVjaWFsIG9wZXJhdGlvbiB0eXBlIHdoaWNoIGlzIHVzZWQgdG8gcmVwcmVzZW50IHRoZSBiZWdpbm5pbmcgYW5kIGVuZCBub2RlcyBvZiBhIGxpbmtlZFxuICAgKiBsaXN0IG9mIG9wZXJhdGlvbnMuXG4gICAqL1xuICBMaXN0RW5kLFxuXG4gIC8qKlxuICAgKiBBbiBvcGVyYXRpb24gd2hpY2ggd3JhcHMgYW4gb3V0cHV0IEFTVCBzdGF0ZW1lbnQuXG4gICAqL1xuICBTdGF0ZW1lbnQsXG5cbiAgLyoqXG4gICAqIEFuIG9wZXJhdGlvbiB3aGljaCBkZWNsYXJlcyBhbmQgaW5pdGlhbGl6ZXMgYSBgU2VtYW50aWNWYXJpYWJsZWAuXG4gICAqL1xuICBWYXJpYWJsZSxcblxuICAvKipcbiAgICogQW4gb3BlcmF0aW9uIHRvIGJlZ2luIHJlbmRlcmluZyBvZiBhbiBlbGVtZW50LlxuICAgKi9cbiAgRWxlbWVudFN0YXJ0LFxuXG4gIC8qKlxuICAgKiBBbiBvcGVyYXRpb24gdG8gcmVuZGVyIGFuIGVsZW1lbnQgd2l0aCBubyBjaGlsZHJlbi5cbiAgICovXG4gIEVsZW1lbnQsXG5cbiAgLyoqXG4gICAqIEFuIG9wZXJhdGlvbiB3aGljaCBkZWNsYXJlcyBhbiBlbWJlZGRlZCB2aWV3LlxuICAgKi9cbiAgVGVtcGxhdGUsXG5cbiAgLyoqXG4gICAqIEFuIG9wZXJhdGlvbiB0byBlbmQgcmVuZGVyaW5nIG9mIGFuIGVsZW1lbnQgcHJldmlvdXNseSBzdGFydGVkIHdpdGggYEVsZW1lbnRTdGFydGAuXG4gICAqL1xuICBFbGVtZW50RW5kLFxuXG4gIC8qKlxuICAgKiBBbiBvcGVyYXRpb24gdG8gYmVnaW4gYW4gYG5nLWNvbnRhaW5lcmAuXG4gICAqL1xuICBDb250YWluZXJTdGFydCxcblxuICAvKipcbiAgICogQW4gb3BlcmF0aW9uIGZvciBhbiBgbmctY29udGFpbmVyYCB3aXRoIG5vIGNoaWxkcmVuLlxuICAgKi9cbiAgQ29udGFpbmVyLFxuXG4gIC8qKlxuICAgKiBBbiBvcGVyYXRpb24gdG8gZW5kIGFuIGBuZy1jb250YWluZXJgLlxuICAgKi9cbiAgQ29udGFpbmVyRW5kLFxuXG4gIC8qKlxuICAgKiBBbiBvcGVyYXRpb24gZGlzYWJsZSBiaW5kaW5nIGZvciBzdWJzZXF1ZW50IGVsZW1lbnRzLCB3aGljaCBhcmUgZGVzY2VuZGFudHMgb2YgYSBub24tYmluZGFibGVcbiAgICogbm9kZS5cbiAgICovXG4gIERpc2FibGVCaW5kaW5ncyxcblxuICAvKipcbiAgICogQW4gb3AgdG8gY29uZGl0aW9uYWxseSByZW5kZXIgYSB0ZW1wbGF0ZS5cbiAgICovXG4gIENvbmRpdGlvbmFsLFxuXG4gIC8qKlxuICAgKiBBbiBvcGVyYXRpb24gdG8gcmUtZW5hYmxlIGJpbmRpbmcsIGFmdGVyIGl0IHdhcyBwcmV2aW91c2x5IGRpc2FibGVkLlxuICAgKi9cbiAgRW5hYmxlQmluZGluZ3MsXG5cbiAgLyoqXG4gICAqIEFuIG9wZXJhdGlvbiB0byByZW5kZXIgYSB0ZXh0IG5vZGUuXG4gICAqL1xuICBUZXh0LFxuXG4gIC8qKlxuICAgKiBBbiBvcGVyYXRpb24gZGVjbGFyaW5nIGFuIGV2ZW50IGxpc3RlbmVyIGZvciBhbiBlbGVtZW50LlxuICAgKi9cbiAgTGlzdGVuZXIsXG5cbiAgLyoqXG4gICAqIEFuIG9wZXJhdGlvbiB0byBpbnRlcnBvbGF0ZSB0ZXh0IGludG8gYSB0ZXh0IG5vZGUuXG4gICAqL1xuICBJbnRlcnBvbGF0ZVRleHQsXG5cbiAgLyoqXG4gICAqIEFuIGludGVybWVkaWF0ZSBiaW5kaW5nIG9wLCB0aGF0IGhhcyBub3QgeWV0IGJlZW4gcHJvY2Vzc2VkIGludG8gYW4gaW5kaXZpZHVhbCBwcm9wZXJ0eSxcbiAgICogYXR0cmlidXRlLCBzdHlsZSwgZXRjLlxuICAgKi9cbiAgQmluZGluZyxcblxuICAvKipcbiAgICogQW4gb3BlcmF0aW9uIHRvIGJpbmQgYW4gZXhwcmVzc2lvbiB0byBhIHByb3BlcnR5IG9mIGFuIGVsZW1lbnQuXG4gICAqL1xuICBQcm9wZXJ0eSxcblxuICAvKipcbiAgICogQW4gb3BlcmF0aW9uIHRvIGJpbmQgYW4gZXhwcmVzc2lvbiB0byBhIHN0eWxlIHByb3BlcnR5IG9mIGFuIGVsZW1lbnQuXG4gICAqL1xuICBTdHlsZVByb3AsXG5cbiAgLyoqXG4gICAqIEFuIG9wZXJhdGlvbiB0byBiaW5kIGFuIGV4cHJlc3Npb24gdG8gYSBjbGFzcyBwcm9wZXJ0eSBvZiBhbiBlbGVtZW50LlxuICAgKi9cbiAgQ2xhc3NQcm9wLFxuXG4gIC8qKlxuICAgKiBBbiBvcGVyYXRpb24gdG8gYmluZCBhbiBleHByZXNzaW9uIHRvIHRoZSBzdHlsZXMgb2YgYW4gZWxlbWVudC5cbiAgICovXG4gIFN0eWxlTWFwLFxuXG4gIC8qKlxuICAgKiBBbiBvcGVyYXRpb24gdG8gYmluZCBhbiBleHByZXNzaW9uIHRvIHRoZSBjbGFzc2VzIG9mIGFuIGVsZW1lbnQuXG4gICAqL1xuICBDbGFzc01hcCxcblxuICAvKipcbiAgICogQW4gb3BlcmF0aW9uIHRvIGFkdmFuY2UgdGhlIHJ1bnRpbWUncyBpbXBsaWNpdCBzbG90IGNvbnRleHQgZHVyaW5nIHRoZSB1cGRhdGUgcGhhc2Ugb2YgYSB2aWV3LlxuICAgKi9cbiAgQWR2YW5jZSxcblxuICAvKipcbiAgICogQW4gb3BlcmF0aW9uIHRvIGluc3RhbnRpYXRlIGEgcGlwZS5cbiAgICovXG4gIFBpcGUsXG5cbiAgLyoqXG4gICAqIEFuIG9wZXJhdGlvbiB0byBhc3NvY2lhdGUgYW4gYXR0cmlidXRlIHdpdGggYW4gZWxlbWVudC5cbiAgICovXG4gIEF0dHJpYnV0ZSxcblxuICAvKipcbiAgICogQW4gYXR0cmlidXRlIHRoYXQgaGFzIGJlZW4gZXh0cmFjdGVkIGZvciBpbmNsdXNpb24gaW4gdGhlIGNvbnN0cyBhcnJheS5cbiAgICovXG4gIEV4dHJhY3RlZEF0dHJpYnV0ZSxcblxuICAvKipcbiAgICogQW4gb3BlcmF0aW9uIHRoYXQgY29uZmlndXJlcyBhIGBAZGVmZXJgIGJsb2NrLlxuICAgKi9cbiAgRGVmZXIsXG5cbiAgLyoqXG4gICAqIEFuIG9wZXJhdGlvbiB0aGF0IGNvbnRyb2xzIHdoZW4gYSBgQGRlZmVyYCBsb2Fkcy5cbiAgICovXG4gIERlZmVyT24sXG5cbiAgLyoqXG4gICAqIEFuIG9wZXJhdGlvbiB0aGF0IGNvbnRyb2xzIHdoZW4gYSBgQGRlZmVyYCBsb2FkcywgdXNpbmcgYSBjdXN0b20gZXhwcmVzc2lvbiBhcyB0aGUgY29uZGl0aW9uLlxuICAgKi9cbiAgRGVmZXJXaGVuLFxuXG4gIC8qKlxuICAgKiBBbiBpMThuIG1lc3NhZ2UgdGhhdCBoYXMgYmVlbiBleHRyYWN0ZWQgZm9yIGluY2x1c2lvbiBpbiB0aGUgY29uc3RzIGFycmF5LlxuICAgKi9cbiAgSTE4bk1lc3NhZ2UsXG5cbiAgLyoqXG4gICAqIEEgaG9zdCBiaW5kaW5nIHByb3BlcnR5LlxuICAgKi9cbiAgSG9zdFByb3BlcnR5LFxuXG4gIC8qKlxuICAgKiBBIG5hbWVzcGFjZSBjaGFuZ2UsIHdoaWNoIGNhdXNlcyB0aGUgc3Vic2VxdWVudCBlbGVtZW50cyB0byBiZSBwcm9jZXNzZWQgYXMgZWl0aGVyIEhUTUwgb3IgU1ZHLlxuICAgKi9cbiAgTmFtZXNwYWNlLFxuXG4gIC8qKlxuICAgKiBDb25maWd1cmUgYSBjb250ZW50IHByb2plY2l0b24gZGVmaW5pdGlvbiBmb3IgdGhlIHZpZXcuXG4gICAqL1xuICBQcm9qZWN0aW9uRGVmLFxuXG4gIC8qKlxuICAgKiBDcmVhdGUgYSBjb250ZW50IHByb2plY3Rpb24gc2xvdC5cbiAgICovXG4gIFByb2plY3Rpb24sXG5cbiAgLyoqXG4gICAqIENyZWF0ZSBhIHJlcGVhdGVyIGNyZWF0aW9uIGluc3RydWN0aW9uIG9wLlxuICAgKi9cbiAgUmVwZWF0ZXJDcmVhdGUsXG5cbiAgLyoqXG4gICAqIEFuIHVwZGF0ZSB1cCBmb3IgYSByZXBlYXRlci5cbiAgICovXG4gIFJlcGVhdGVyLFxuXG4gIC8qKlxuICAgKiBBbiBvcGVyYXRpb24gdG8gYmluZCBhbiBleHByZXNzaW9uIHRvIHRoZSBwcm9wZXJ0eSBzaWRlIG9mIGEgdHdvLXdheSBiaW5kaW5nLlxuICAgKi9cbiAgVHdvV2F5UHJvcGVydHksXG5cbiAgLyoqXG4gICAqIEFuIG9wZXJhdGlvbiBkZWNsYXJpbmcgdGhlIGV2ZW50IHNpZGUgb2YgYSB0d28td2F5IGJpbmRpbmcuXG4gICAqL1xuICBUd29XYXlMaXN0ZW5lcixcblxuICAvKipcbiAgICogQSBjcmVhdGlvbi10aW1lIG9wZXJhdGlvbiB0aGF0IGluaXRpYWxpemVzIHRoZSBzbG90IGZvciBhIGBAbGV0YCBkZWNsYXJhdGlvbi5cbiAgICovXG4gIERlY2xhcmVMZXQsXG5cbiAgLyoqXG4gICAqIEFuIHVwZGF0ZS10aW1lIG9wZXJhdGlvbiB0aGF0IHN0b3JlcyB0aGUgY3VycmVudCB2YWx1ZSBvZiBhIGBAbGV0YCBkZWNsYXJhdGlvbi5cbiAgICovXG4gIFN0b3JlTGV0LFxuXG4gIC8qKlxuICAgKiBUaGUgc3RhcnQgb2YgYW4gaTE4biBibG9jay5cbiAgICovXG4gIEkxOG5TdGFydCxcblxuICAvKipcbiAgICogQSBzZWxmLWNsb3NpbmcgaTE4biBvbiBhIHNpbmdsZSBlbGVtZW50LlxuICAgKi9cbiAgSTE4bixcblxuICAvKipcbiAgICogVGhlIGVuZCBvZiBhbiBpMThuIGJsb2NrLlxuICAgKi9cbiAgSTE4bkVuZCxcblxuICAvKipcbiAgICogQW4gZXhwcmVzc2lvbiBpbiBhbiBpMThuIG1lc3NhZ2UuXG4gICAqL1xuICBJMThuRXhwcmVzc2lvbixcblxuICAvKipcbiAgICogQW4gaW5zdHJ1Y3Rpb24gdGhhdCBhcHBsaWVzIGEgc2V0IG9mIGkxOG4gZXhwcmVzc2lvbnMuXG4gICAqL1xuICBJMThuQXBwbHksXG5cbiAgLyoqXG4gICAqIEFuIGluc3RydWN0aW9uIHRvIGNyZWF0ZSBhbiBJQ1UgZXhwcmVzc2lvbi5cbiAgICovXG4gIEljdVN0YXJ0LFxuXG4gIC8qKlxuICAgKiBBbiBpbnN0cnVjdGlvbiB0byB1cGRhdGUgYW4gSUNVIGV4cHJlc3Npb24uXG4gICAqL1xuICBJY3VFbmQsXG5cbiAgLyoqXG4gICAqIEFuIGluc3RydWN0aW9uIHJlcHJlc2VudGluZyBhIHBsYWNlaG9sZGVyIGluIGFuIElDVSBleHByZXNzaW9uLlxuICAgKi9cbiAgSWN1UGxhY2Vob2xkZXIsXG5cbiAgLyoqXG4gICAqIEFuIGkxOG4gY29udGV4dCBjb250YWluaW5nIGluZm9ybWF0aW9uIG5lZWRlZCB0byBnZW5lcmF0ZSBhbiBpMThuIG1lc3NhZ2UuXG4gICAqL1xuICBJMThuQ29udGV4dCxcblxuICAvKipcbiAgICogQSBjcmVhdGlvbiBvcCB0aGF0IGNvcnJlc3BvbmRzIHRvIGkxOG4gYXR0cmlidXRlcyBvbiBhbiBlbGVtZW50LlxuICAgKi9cbiAgSTE4bkF0dHJpYnV0ZXMsXG59XG5cbi8qKlxuICogRGlzdGluZ3Vpc2hlcyBkaWZmZXJlbnQga2luZHMgb2YgSVIgZXhwcmVzc2lvbnMuXG4gKi9cbmV4cG9ydCBlbnVtIEV4cHJlc3Npb25LaW5kIHtcbiAgLyoqXG4gICAqIFJlYWQgb2YgYSB2YXJpYWJsZSBpbiBhIGxleGljYWwgc2NvcGUuXG4gICAqL1xuICBMZXhpY2FsUmVhZCxcblxuICAvKipcbiAgICogQSByZWZlcmVuY2UgdG8gdGhlIGN1cnJlbnQgdmlldyBjb250ZXh0LlxuICAgKi9cbiAgQ29udGV4dCxcblxuICAvKipcbiAgICogQSByZWZlcmVuY2UgdG8gdGhlIHZpZXcgY29udGV4dCwgZm9yIHVzZSBpbnNpZGUgYSB0cmFjayBmdW5jdGlvbi5cbiAgICovXG4gIFRyYWNrQ29udGV4dCxcblxuICAvKipcbiAgICogUmVhZCBvZiBhIHZhcmlhYmxlIGRlY2xhcmVkIGluIGEgYFZhcmlhYmxlT3BgLlxuICAgKi9cbiAgUmVhZFZhcmlhYmxlLFxuXG4gIC8qKlxuICAgKiBSdW50aW1lIG9wZXJhdGlvbiB0byBuYXZpZ2F0ZSB0byB0aGUgbmV4dCB2aWV3IGNvbnRleHQgaW4gdGhlIHZpZXcgaGllcmFyY2h5LlxuICAgKi9cbiAgTmV4dENvbnRleHQsXG5cbiAgLyoqXG4gICAqIFJ1bnRpbWUgb3BlcmF0aW9uIHRvIHJldHJpZXZlIHRoZSB2YWx1ZSBvZiBhIGxvY2FsIHJlZmVyZW5jZS5cbiAgICovXG4gIFJlZmVyZW5jZSxcblxuICAvKipcbiAgICogQSBjYWxsIHN0b3JpbmcgdGhlIHZhbHVlIG9mIGEgYEBsZXRgIGRlY2xhcmF0aW9uLlxuICAgKi9cbiAgU3RvcmVMZXQsXG5cbiAgLyoqXG4gICAqIEEgcmVmZXJlbmNlIHRvIGEgYEBsZXRgIGRlY2xhcmF0aW9uIHJlYWQgZnJvbSB0aGUgY29udGV4dCB2aWV3LlxuICAgKi9cbiAgQ29udGV4dExldFJlZmVyZW5jZSxcblxuICAvKipcbiAgICogUnVudGltZSBvcGVyYXRpb24gdG8gc25hcHNob3QgdGhlIGN1cnJlbnQgdmlldyBjb250ZXh0LlxuICAgKi9cbiAgR2V0Q3VycmVudFZpZXcsXG5cbiAgLyoqXG4gICAqIFJ1bnRpbWUgb3BlcmF0aW9uIHRvIHJlc3RvcmUgYSBzbmFwc2hvdHRlZCB2aWV3LlxuICAgKi9cbiAgUmVzdG9yZVZpZXcsXG5cbiAgLyoqXG4gICAqIFJ1bnRpbWUgb3BlcmF0aW9uIHRvIHJlc2V0IHRoZSBjdXJyZW50IHZpZXcgY29udGV4dCBhZnRlciBgUmVzdG9yZVZpZXdgLlxuICAgKi9cbiAgUmVzZXRWaWV3LFxuXG4gIC8qKlxuICAgKiBEZWZpbmVzIGFuZCBjYWxscyBhIGZ1bmN0aW9uIHdpdGggY2hhbmdlLWRldGVjdGVkIGFyZ3VtZW50cy5cbiAgICovXG4gIFB1cmVGdW5jdGlvbkV4cHIsXG5cbiAgLyoqXG4gICAqIEluZGljYXRlcyBhIHBvc2l0aW9uYWwgcGFyYW1ldGVyIHRvIGEgcHVyZSBmdW5jdGlvbiBkZWZpbml0aW9uLlxuICAgKi9cbiAgUHVyZUZ1bmN0aW9uUGFyYW1ldGVyRXhwcixcblxuICAvKipcbiAgICogQmluZGluZyB0byBhIHBpcGUgdHJhbnNmb3JtYXRpb24uXG4gICAqL1xuICBQaXBlQmluZGluZyxcblxuICAvKipcbiAgICogQmluZGluZyB0byBhIHBpcGUgdHJhbnNmb3JtYXRpb24gd2l0aCBhIHZhcmlhYmxlIG51bWJlciBvZiBhcmd1bWVudHMuXG4gICAqL1xuICBQaXBlQmluZGluZ1ZhcmlhZGljLFxuXG4gIC8qXG4gICAqIEEgc2FmZSBwcm9wZXJ0eSByZWFkIHJlcXVpcmluZyBleHBhbnNpb24gaW50byBhIG51bGwgY2hlY2suXG4gICAqL1xuICBTYWZlUHJvcGVydHlSZWFkLFxuXG4gIC8qKlxuICAgKiBBIHNhZmUga2V5ZWQgcmVhZCByZXF1aXJpbmcgZXhwYW5zaW9uIGludG8gYSBudWxsIGNoZWNrLlxuICAgKi9cbiAgU2FmZUtleWVkUmVhZCxcblxuICAvKipcbiAgICogQSBzYWZlIGZ1bmN0aW9uIGNhbGwgcmVxdWlyaW5nIGV4cGFuc2lvbiBpbnRvIGEgbnVsbCBjaGVjay5cbiAgICovXG4gIFNhZmVJbnZva2VGdW5jdGlvbixcblxuICAvKipcbiAgICogQW4gaW50ZXJtZWRpYXRlIGV4cHJlc3Npb24gdGhhdCB3aWxsIGJlIGV4cGFuZGVkIGZyb20gYSBzYWZlIHJlYWQgaW50byBhbiBleHBsaWNpdCB0ZXJuYXJ5LlxuICAgKi9cbiAgU2FmZVRlcm5hcnlFeHByLFxuXG4gIC8qKlxuICAgKiBBbiBlbXB0eSBleHByZXNzaW9uIHRoYXQgd2lsbCBiZSBzdGlwcGVkIGJlZm9yZSBnZW5lcmF0aW5nIHRoZSBmaW5hbCBvdXRwdXQuXG4gICAqL1xuICBFbXB0eUV4cHIsXG5cbiAgLypcbiAgICogQW4gYXNzaWdubWVudCB0byBhIHRlbXBvcmFyeSB2YXJpYWJsZS5cbiAgICovXG4gIEFzc2lnblRlbXBvcmFyeUV4cHIsXG5cbiAgLyoqXG4gICAqIEEgcmVmZXJlbmNlIHRvIGEgdGVtcG9yYXJ5IHZhcmlhYmxlLlxuICAgKi9cbiAgUmVhZFRlbXBvcmFyeUV4cHIsXG5cbiAgLyoqXG4gICAqIEFuIGV4cHJlc3Npb24gdGhhdCB3aWxsIGNhdXNlIGEgbGl0ZXJhbCBzbG90IGluZGV4IHRvIGJlIGVtaXR0ZWQuXG4gICAqL1xuICBTbG90TGl0ZXJhbEV4cHIsXG5cbiAgLyoqXG4gICAqIEEgdGVzdCBleHByZXNzaW9uIGZvciBhIGNvbmRpdGlvbmFsIG9wLlxuICAgKi9cbiAgQ29uZGl0aW9uYWxDYXNlLFxuXG4gIC8qKlxuICAgKiBBbiBleHByZXNzaW9uIHRoYXQgd2lsbCBiZSBhdXRvbWF0aWNhbGx5IGV4dHJhY3RlZCB0byB0aGUgY29tcG9uZW50IGNvbnN0IGFycmF5LlxuICAgKi9cbiAgQ29uc3RDb2xsZWN0ZWQsXG5cbiAgLyoqXG4gICAqIE9wZXJhdGlvbiB0aGF0IHNldHMgdGhlIHZhbHVlIG9mIGEgdHdvLXdheSBiaW5kaW5nLlxuICAgKi9cbiAgVHdvV2F5QmluZGluZ1NldCxcbn1cblxuZXhwb3J0IGVudW0gVmFyaWFibGVGbGFncyB7XG4gIE5vbmUgPSAwYjAwMDAsXG5cbiAgLyoqXG4gICAqIEFsd2F5cyBpbmxpbmUgdGhpcyB2YXJpYWJsZSwgcmVnYXJkbGVzcyBvZiB0aGUgbnVtYmVyIG9mIHRpbWVzIGl0J3MgdXNlZC5cbiAgICogQW4gYEFsd2F5c0lubGluZWAgdmFyaWFibGUgbWF5IG5vdCBkZXBlbmQgb24gY29udGV4dCwgYmVjYXVzZSBkb2luZyBzbyBtYXkgY2F1c2Ugc2lkZSBlZmZlY3RzXG4gICAqIHRoYXQgYXJlIGlsbGVnYWwgd2hlbiBtdWx0aS1pbmxpbmVkLiAoVGhlIG9wdGltaXplciB3aWxsIGVuZm9yY2UgdGhpcyBjb25zdHJhaW50LilcbiAgICovXG4gIEFsd2F5c0lubGluZSA9IDBiMDAwMSxcbn1cbi8qKlxuICogRGlzdGluZ3Vpc2hlcyBiZXR3ZWVuIGRpZmZlcmVudCBraW5kcyBvZiBgU2VtYW50aWNWYXJpYWJsZWBzLlxuICovXG5leHBvcnQgZW51bSBTZW1hbnRpY1ZhcmlhYmxlS2luZCB7XG4gIC8qKlxuICAgKiBSZXByZXNlbnRzIHRoZSBjb250ZXh0IG9mIGEgcGFydGljdWxhciB2aWV3LlxuICAgKi9cbiAgQ29udGV4dCxcblxuICAvKipcbiAgICogUmVwcmVzZW50cyBhbiBpZGVudGlmaWVyIGRlY2xhcmVkIGluIHRoZSBsZXhpY2FsIHNjb3BlIG9mIGEgdmlldy5cbiAgICovXG4gIElkZW50aWZpZXIsXG5cbiAgLyoqXG4gICAqIFJlcHJlc2VudHMgYSBzYXZlZCBzdGF0ZSB0aGF0IGNhbiBiZSB1c2VkIHRvIHJlc3RvcmUgYSB2aWV3IGluIGEgbGlzdGVuZXIgaGFuZGxlciBmdW5jdGlvbi5cbiAgICovXG4gIFNhdmVkVmlldyxcblxuICAvKipcbiAgICogQW4gYWxpYXMgZ2VuZXJhdGVkIGJ5IGEgc3BlY2lhbCBlbWJlZGRlZCB2aWV3IHR5cGUgKGUuZy4gYSBgQGZvcmAgYmxvY2spLlxuICAgKi9cbiAgQWxpYXMsXG59XG5cbi8qKlxuICogV2hldGhlciB0byBjb21waWxlIGluIGNvbXBhdGliaWx0eSBtb2RlLiBJbiBjb21wYXRpYmlsaXR5IG1vZGUsIHRoZSB0ZW1wbGF0ZSBwaXBlbGluZSB3aWxsXG4gKiBhdHRlbXB0IHRvIG1hdGNoIHRoZSBvdXRwdXQgb2YgYFRlbXBsYXRlRGVmaW5pdGlvbkJ1aWxkZXJgIGFzIGV4YWN0bHkgYXMgcG9zc2libGUsIGF0IHRoZSBjb3N0XG4gKiBvZiBwcm9kdWNpbmcgcXVpcmt5IG9yIGxhcmdlciBjb2RlIGluIHNvbWUgY2FzZXMuXG4gKi9cbmV4cG9ydCBlbnVtIENvbXBhdGliaWxpdHlNb2RlIHtcbiAgTm9ybWFsLFxuICBUZW1wbGF0ZURlZmluaXRpb25CdWlsZGVyLFxufVxuXG4vKipcbiAqIEVudW1lcmF0aW9uIG9mIHRoZSB0eXBlcyBvZiBhdHRyaWJ1dGVzIHdoaWNoIGNhbiBiZSBhcHBsaWVkIHRvIGFuIGVsZW1lbnQuXG4gKi9cbmV4cG9ydCBlbnVtIEJpbmRpbmdLaW5kIHtcbiAgLyoqXG4gICAqIFN0YXRpYyBhdHRyaWJ1dGVzLlxuICAgKi9cbiAgQXR0cmlidXRlLFxuXG4gIC8qKlxuICAgKiBDbGFzcyBiaW5kaW5ncy5cbiAgICovXG4gIENsYXNzTmFtZSxcblxuICAvKipcbiAgICogU3R5bGUgYmluZGluZ3MuXG4gICAqL1xuICBTdHlsZVByb3BlcnR5LFxuXG4gIC8qKlxuICAgKiBEeW5hbWljIHByb3BlcnR5IGJpbmRpbmdzLlxuICAgKi9cbiAgUHJvcGVydHksXG5cbiAgLyoqXG4gICAqIFByb3BlcnR5IG9yIGF0dHJpYnV0ZSBiaW5kaW5ncyBvbiBhIHRlbXBsYXRlLlxuICAgKi9cbiAgVGVtcGxhdGUsXG5cbiAgLyoqXG4gICAqIEludGVybmF0aW9uYWxpemVkIGF0dHJpYnV0ZXMuXG4gICAqL1xuICBJMThuLFxuXG4gIC8qKlxuICAgKiBBbmltYXRpb24gcHJvcGVydHkgYmluZGluZ3MuXG4gICAqL1xuICBBbmltYXRpb24sXG5cbiAgLyoqXG4gICAqIFByb3BlcnR5IHNpZGUgb2YgYSB0d28td2F5IGJpbmRpbmcuXG4gICAqL1xuICBUd29XYXlQcm9wZXJ0eSxcbn1cblxuLyoqXG4gKiBFbnVtZXJhdGlvbiBvZiBwb3NzaWJsZSB0aW1lcyBpMThuIHBhcmFtcyBjYW4gYmUgcmVzb2x2ZWQuXG4gKi9cbmV4cG9ydCBlbnVtIEkxOG5QYXJhbVJlc29sdXRpb25UaW1lIHtcbiAgLyoqXG4gICAqIFBhcmFtIGlzIHJlc29sdmVkIGF0IG1lc3NhZ2UgY3JlYXRpb24gdGltZS4gTW9zdCBwYXJhbXMgc2hvdWxkIGJlIHJlc29sdmVkIGF0IG1lc3NhZ2UgY3JlYXRpb25cbiAgICogdGltZS4gSG93ZXZlciwgSUNVIHBhcmFtcyBuZWVkIHRvIGJlIGhhbmRsZWQgaW4gcG9zdC1wcm9jZXNzaW5nLlxuICAgKi9cbiAgQ3JlYXRpb24sXG5cbiAgLyoqXG4gICAqIFBhcmFtIGlzIHJlc29sdmVkIGR1cmluZyBwb3N0LXByb2Nlc3NpbmcuIFRoaXMgc2hvdWxkIGJlIHVzZWQgZm9yIHBhcmFtcyB3aG9zZSB2YWx1ZSBjb21lcyBmcm9tXG4gICAqIGFuIElDVS5cbiAgICovXG4gIFBvc3Rwcm9jY2Vzc2luZyxcbn1cblxuLyoqXG4gKiBUaGUgY29udGV4dHMgaW4gd2hpY2ggYW4gaTE4biBleHByZXNzaW9uIGNhbiBiZSB1c2VkLlxuICovXG5leHBvcnQgZW51bSBJMThuRXhwcmVzc2lvbkZvciB7XG4gIC8qKlxuICAgKiBUaGlzIGV4cHJlc3Npb24gaXMgdXNlZCBhcyBhIHZhbHVlIChpLmUuIGluc2lkZSBhbiBpMThuIGJsb2NrKS5cbiAgICovXG4gIEkxOG5UZXh0LFxuXG4gIC8qKlxuICAgKiBUaGlzIGV4cHJlc3Npb24gaXMgdXNlZCBpbiBhIGJpbmRpbmcuXG4gICAqL1xuICBJMThuQXR0cmlidXRlLFxufVxuXG4vKipcbiAqIEZsYWdzIHRoYXQgZGVzY3JpYmUgd2hhdCBhbiBpMThuIHBhcmFtIHZhbHVlLiBUaGVzZSBkZXRlcm1pbmUgaG93IHRoZSB2YWx1ZSBpcyBzZXJpYWxpemVkIGludG9cbiAqIHRoZSBmaW5hbCBtYXAuXG4gKi9cbmV4cG9ydCBlbnVtIEkxOG5QYXJhbVZhbHVlRmxhZ3Mge1xuICBOb25lID0gMGIwMDAwLFxuXG4gIC8qKlxuICAgKiAgVGhpcyB2YWx1ZSByZXByZXNlbnRzIGFuIGVsZW1lbnQgdGFnLlxuICAgKi9cbiAgRWxlbWVudFRhZyA9IDBiMSxcblxuICAvKipcbiAgICogVGhpcyB2YWx1ZSByZXByZXNlbnRzIGEgdGVtcGxhdGUgdGFnLlxuICAgKi9cbiAgVGVtcGxhdGVUYWcgPSAwYjEwLFxuXG4gIC8qKlxuICAgKiBUaGlzIHZhbHVlIHJlcHJlc2VudHMgdGhlIG9wZW5pbmcgb2YgYSB0YWcuXG4gICAqL1xuICBPcGVuVGFnID0gMGIwMTAwLFxuXG4gIC8qKlxuICAgKiBUaGlzIHZhbHVlIHJlcHJlc2VudHMgdGhlIGNsb3Npbmcgb2YgYSB0YWcuXG4gICAqL1xuICBDbG9zZVRhZyA9IDBiMTAwMCxcblxuICAvKipcbiAgICogVGhpcyB2YWx1ZSByZXByZXNlbnRzIGFuIGkxOG4gZXhwcmVzc2lvbiBpbmRleC5cbiAgICovXG4gIEV4cHJlc3Npb25JbmRleCA9IDBiMTAwMDAsXG59XG5cbi8qKlxuICogV2hldGhlciB0aGUgYWN0aXZlIG5hbWVzcGFjZSBpcyBIVE1MLCBNYXRoTUwsIG9yIFNWRyBtb2RlLlxuICovXG5leHBvcnQgZW51bSBOYW1lc3BhY2Uge1xuICBIVE1MLFxuICBTVkcsXG4gIE1hdGgsXG59XG5cbi8qKlxuICogVGhlIHR5cGUgb2YgYSBgQGRlZmVyYCB0cmlnZ2VyLCBmb3IgdXNlIGluIHRoZSBpci5cbiAqL1xuZXhwb3J0IGVudW0gRGVmZXJUcmlnZ2VyS2luZCB7XG4gIElkbGUsXG4gIEltbWVkaWF0ZSxcbiAgVGltZXIsXG4gIEhvdmVyLFxuICBJbnRlcmFjdGlvbixcbiAgVmlld3BvcnQsXG59XG5cbi8qKlxuICogS2luZHMgb2YgaTE4biBjb250ZXh0cy4gVGhleSBjYW4gYmUgY3JlYXRlZCBiZWNhdXNlIG9mIHJvb3QgaTE4biBibG9ja3MsIG9yIElDVXMuXG4gKi9cbmV4cG9ydCBlbnVtIEkxOG5Db250ZXh0S2luZCB7XG4gIFJvb3RJMThuLFxuICBJY3UsXG4gIEF0dHIsXG59XG5cbmV4cG9ydCBlbnVtIFRlbXBsYXRlS2luZCB7XG4gIE5nVGVtcGxhdGUsXG4gIFN0cnVjdHVyYWwsXG4gIEJsb2NrLFxufVxuIl19