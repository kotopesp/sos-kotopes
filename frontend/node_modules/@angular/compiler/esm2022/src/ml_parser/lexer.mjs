/**
 * @license
 * Copyright Google LLC All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
import * as chars from '../chars';
import { ParseError, ParseLocation, ParseSourceFile, ParseSourceSpan } from '../parse_util';
import { DEFAULT_INTERPOLATION_CONFIG } from './defaults';
import { NAMED_ENTITIES } from './entities';
import { TagContentType } from './tags';
export class TokenError extends ParseError {
    constructor(errorMsg, tokenType, span) {
        super(span, errorMsg);
        this.tokenType = tokenType;
    }
}
export class TokenizeResult {
    constructor(tokens, errors, nonNormalizedIcuExpressions) {
        this.tokens = tokens;
        this.errors = errors;
        this.nonNormalizedIcuExpressions = nonNormalizedIcuExpressions;
    }
}
export function tokenize(source, url, getTagDefinition, options = {}) {
    const tokenizer = new _Tokenizer(new ParseSourceFile(source, url), getTagDefinition, options);
    tokenizer.tokenize();
    return new TokenizeResult(mergeTextTokens(tokenizer.tokens), tokenizer.errors, tokenizer.nonNormalizedIcuExpressions);
}
const _CR_OR_CRLF_REGEXP = /\r\n?/g;
function _unexpectedCharacterErrorMsg(charCode) {
    const char = charCode === chars.$EOF ? 'EOF' : String.fromCharCode(charCode);
    return `Unexpected character "${char}"`;
}
function _unknownEntityErrorMsg(entitySrc) {
    return `Unknown entity "${entitySrc}" - use the "&#<decimal>;" or  "&#x<hex>;" syntax`;
}
function _unparsableEntityErrorMsg(type, entityStr) {
    return `Unable to parse entity "${entityStr}" - ${type} character reference entities must end with ";"`;
}
var CharacterReferenceType;
(function (CharacterReferenceType) {
    CharacterReferenceType["HEX"] = "hexadecimal";
    CharacterReferenceType["DEC"] = "decimal";
})(CharacterReferenceType || (CharacterReferenceType = {}));
class _ControlFlowError {
    constructor(error) {
        this.error = error;
    }
}
// See https://www.w3.org/TR/html51/syntax.html#writing-html-documents
class _Tokenizer {
    /**
     * @param _file The html source file being tokenized.
     * @param _getTagDefinition A function that will retrieve a tag definition for a given tag name.
     * @param options Configuration of the tokenization.
     */
    constructor(_file, _getTagDefinition, options) {
        this._getTagDefinition = _getTagDefinition;
        this._currentTokenStart = null;
        this._currentTokenType = null;
        this._expansionCaseStack = [];
        this._inInterpolation = false;
        this.tokens = [];
        this.errors = [];
        this.nonNormalizedIcuExpressions = [];
        this._tokenizeIcu = options.tokenizeExpansionForms || false;
        this._interpolationConfig = options.interpolationConfig || DEFAULT_INTERPOLATION_CONFIG;
        this._leadingTriviaCodePoints =
            options.leadingTriviaChars && options.leadingTriviaChars.map((c) => c.codePointAt(0) || 0);
        const range = options.range || {
            endPos: _file.content.length,
            startPos: 0,
            startLine: 0,
            startCol: 0,
        };
        this._cursor = options.escapedString
            ? new EscapedCharacterCursor(_file, range)
            : new PlainCharacterCursor(_file, range);
        this._preserveLineEndings = options.preserveLineEndings || false;
        this._i18nNormalizeLineEndingsInICUs = options.i18nNormalizeLineEndingsInICUs || false;
        this._tokenizeBlocks = options.tokenizeBlocks ?? true;
        this._tokenizeLet = options.tokenizeLet ?? true;
        try {
            this._cursor.init();
        }
        catch (e) {
            this.handleError(e);
        }
    }
    _processCarriageReturns(content) {
        if (this._preserveLineEndings) {
            return content;
        }
        // https://www.w3.org/TR/html51/syntax.html#preprocessing-the-input-stream
        // In order to keep the original position in the source, we can not
        // pre-process it.
        // Instead CRs are processed right before instantiating the tokens.
        return content.replace(_CR_OR_CRLF_REGEXP, '\n');
    }
    tokenize() {
        while (this._cursor.peek() !== chars.$EOF) {
            const start = this._cursor.clone();
            try {
                if (this._attemptCharCode(chars.$LT)) {
                    if (this._attemptCharCode(chars.$BANG)) {
                        if (this._attemptCharCode(chars.$LBRACKET)) {
                            this._consumeCdata(start);
                        }
                        else if (this._attemptCharCode(chars.$MINUS)) {
                            this._consumeComment(start);
                        }
                        else {
                            this._consumeDocType(start);
                        }
                    }
                    else if (this._attemptCharCode(chars.$SLASH)) {
                        this._consumeTagClose(start);
                    }
                    else {
                        this._consumeTagOpen(start);
                    }
                }
                else if (this._tokenizeLet &&
                    // Use `peek` instead of `attempCharCode` since we
                    // don't want to advance in case it's not `@let`.
                    this._cursor.peek() === chars.$AT &&
                    !this._inInterpolation &&
                    this._attemptStr('@let')) {
                    this._consumeLetDeclaration(start);
                }
                else if (this._tokenizeBlocks && this._attemptCharCode(chars.$AT)) {
                    this._consumeBlockStart(start);
                }
                else if (this._tokenizeBlocks &&
                    !this._inInterpolation &&
                    !this._isInExpansionCase() &&
                    !this._isInExpansionForm() &&
                    this._attemptCharCode(chars.$RBRACE)) {
                    this._consumeBlockEnd(start);
                }
                else if (!(this._tokenizeIcu && this._tokenizeExpansionForm())) {
                    // In (possibly interpolated) text the end of the text is given by `isTextEnd()`, while
                    // the premature end of an interpolation is given by the start of a new HTML element.
                    this._consumeWithInterpolation(5 /* TokenType.TEXT */, 8 /* TokenType.INTERPOLATION */, () => this._isTextEnd(), () => this._isTagStart());
                }
            }
            catch (e) {
                this.handleError(e);
            }
        }
        this._beginToken(33 /* TokenType.EOF */);
        this._endToken([]);
    }
    _getBlockName() {
        // This allows us to capture up something like `@else if`, but not `@ if`.
        let spacesInNameAllowed = false;
        const nameCursor = this._cursor.clone();
        this._attemptCharCodeUntilFn((code) => {
            if (chars.isWhitespace(code)) {
                return !spacesInNameAllowed;
            }
            if (isBlockNameChar(code)) {
                spacesInNameAllowed = true;
                return false;
            }
            return true;
        });
        return this._cursor.getChars(nameCursor).trim();
    }
    _consumeBlockStart(start) {
        this._beginToken(24 /* TokenType.BLOCK_OPEN_START */, start);
        const startToken = this._endToken([this._getBlockName()]);
        if (this._cursor.peek() === chars.$LPAREN) {
            // Advance past the opening paren.
            this._cursor.advance();
            // Capture the parameters.
            this._consumeBlockParameters();
            // Allow spaces before the closing paren.
            this._attemptCharCodeUntilFn(isNotWhitespace);
            if (this._attemptCharCode(chars.$RPAREN)) {
                // Allow spaces after the paren.
                this._attemptCharCodeUntilFn(isNotWhitespace);
            }
            else {
                startToken.type = 28 /* TokenType.INCOMPLETE_BLOCK_OPEN */;
                return;
            }
        }
        if (this._attemptCharCode(chars.$LBRACE)) {
            this._beginToken(25 /* TokenType.BLOCK_OPEN_END */);
            this._endToken([]);
        }
        else {
            startToken.type = 28 /* TokenType.INCOMPLETE_BLOCK_OPEN */;
        }
    }
    _consumeBlockEnd(start) {
        this._beginToken(26 /* TokenType.BLOCK_CLOSE */, start);
        this._endToken([]);
    }
    _consumeBlockParameters() {
        // Trim the whitespace until the first parameter.
        this._attemptCharCodeUntilFn(isBlockParameterChar);
        while (this._cursor.peek() !== chars.$RPAREN && this._cursor.peek() !== chars.$EOF) {
            this._beginToken(27 /* TokenType.BLOCK_PARAMETER */);
            const start = this._cursor.clone();
            let inQuote = null;
            let openParens = 0;
            // Consume the parameter until the next semicolon or brace.
            // Note that we skip over semicolons/braces inside of strings.
            while ((this._cursor.peek() !== chars.$SEMICOLON && this._cursor.peek() !== chars.$EOF) ||
                inQuote !== null) {
                const char = this._cursor.peek();
                // Skip to the next character if it was escaped.
                if (char === chars.$BACKSLASH) {
                    this._cursor.advance();
                }
                else if (char === inQuote) {
                    inQuote = null;
                }
                else if (inQuote === null && chars.isQuote(char)) {
                    inQuote = char;
                }
                else if (char === chars.$LPAREN && inQuote === null) {
                    openParens++;
                }
                else if (char === chars.$RPAREN && inQuote === null) {
                    if (openParens === 0) {
                        break;
                    }
                    else if (openParens > 0) {
                        openParens--;
                    }
                }
                this._cursor.advance();
            }
            this._endToken([this._cursor.getChars(start)]);
            // Skip to the next parameter.
            this._attemptCharCodeUntilFn(isBlockParameterChar);
        }
    }
    _consumeLetDeclaration(start) {
        this._beginToken(29 /* TokenType.LET_START */, start);
        // Require at least one white space after the `@let`.
        if (chars.isWhitespace(this._cursor.peek())) {
            this._attemptCharCodeUntilFn(isNotWhitespace);
        }
        else {
            const token = this._endToken([this._cursor.getChars(start)]);
            token.type = 32 /* TokenType.INCOMPLETE_LET */;
            return;
        }
        const startToken = this._endToken([this._getLetDeclarationName()]);
        // Skip over white space before the equals character.
        this._attemptCharCodeUntilFn(isNotWhitespace);
        // Expect an equals sign.
        if (!this._attemptCharCode(chars.$EQ)) {
            startToken.type = 32 /* TokenType.INCOMPLETE_LET */;
            return;
        }
        // Skip spaces after the equals.
        this._attemptCharCodeUntilFn((code) => isNotWhitespace(code) && !chars.isNewLine(code));
        this._consumeLetDeclarationValue();
        // Terminate the `@let` with a semicolon.
        const endChar = this._cursor.peek();
        if (endChar === chars.$SEMICOLON) {
            this._beginToken(31 /* TokenType.LET_END */);
            this._endToken([]);
            this._cursor.advance();
        }
        else {
            startToken.type = 32 /* TokenType.INCOMPLETE_LET */;
            startToken.sourceSpan = this._cursor.getSpan(start);
        }
    }
    _getLetDeclarationName() {
        const nameCursor = this._cursor.clone();
        let allowDigit = false;
        this._attemptCharCodeUntilFn((code) => {
            if (chars.isAsciiLetter(code) ||
                code == chars.$$ ||
                code === chars.$_ ||
                // `@let` names can't start with a digit, but digits are valid anywhere else in the name.
                (allowDigit && chars.isDigit(code))) {
                allowDigit = true;
                return false;
            }
            return true;
        });
        return this._cursor.getChars(nameCursor).trim();
    }
    _consumeLetDeclarationValue() {
        const start = this._cursor.clone();
        this._beginToken(30 /* TokenType.LET_VALUE */, start);
        while (this._cursor.peek() !== chars.$EOF) {
            const char = this._cursor.peek();
            // `@let` declarations terminate with a semicolon.
            if (char === chars.$SEMICOLON) {
                break;
            }
            // If we hit a quote, skip over its content since we don't care what's inside.
            if (chars.isQuote(char)) {
                this._cursor.advance();
                this._attemptCharCodeUntilFn((inner) => {
                    if (inner === chars.$BACKSLASH) {
                        this._cursor.advance();
                        return false;
                    }
                    return inner === char;
                });
            }
            this._cursor.advance();
        }
        this._endToken([this._cursor.getChars(start)]);
    }
    /**
     * @returns whether an ICU token has been created
     * @internal
     */
    _tokenizeExpansionForm() {
        if (this.isExpansionFormStart()) {
            this._consumeExpansionFormStart();
            return true;
        }
        if (isExpansionCaseStart(this._cursor.peek()) && this._isInExpansionForm()) {
            this._consumeExpansionCaseStart();
            return true;
        }
        if (this._cursor.peek() === chars.$RBRACE) {
            if (this._isInExpansionCase()) {
                this._consumeExpansionCaseEnd();
                return true;
            }
            if (this._isInExpansionForm()) {
                this._consumeExpansionFormEnd();
                return true;
            }
        }
        return false;
    }
    _beginToken(type, start = this._cursor.clone()) {
        this._currentTokenStart = start;
        this._currentTokenType = type;
    }
    _endToken(parts, end) {
        if (this._currentTokenStart === null) {
            throw new TokenError('Programming error - attempted to end a token when there was no start to the token', this._currentTokenType, this._cursor.getSpan(end));
        }
        if (this._currentTokenType === null) {
            throw new TokenError('Programming error - attempted to end a token which has no token type', null, this._cursor.getSpan(this._currentTokenStart));
        }
        const token = {
            type: this._currentTokenType,
            parts,
            sourceSpan: (end ?? this._cursor).getSpan(this._currentTokenStart, this._leadingTriviaCodePoints),
        };
        this.tokens.push(token);
        this._currentTokenStart = null;
        this._currentTokenType = null;
        return token;
    }
    _createError(msg, span) {
        if (this._isInExpansionForm()) {
            msg += ` (Do you have an unescaped "{" in your template? Use "{{ '{' }}") to escape it.)`;
        }
        const error = new TokenError(msg, this._currentTokenType, span);
        this._currentTokenStart = null;
        this._currentTokenType = null;
        return new _ControlFlowError(error);
    }
    handleError(e) {
        if (e instanceof CursorError) {
            e = this._createError(e.msg, this._cursor.getSpan(e.cursor));
        }
        if (e instanceof _ControlFlowError) {
            this.errors.push(e.error);
        }
        else {
            throw e;
        }
    }
    _attemptCharCode(charCode) {
        if (this._cursor.peek() === charCode) {
            this._cursor.advance();
            return true;
        }
        return false;
    }
    _attemptCharCodeCaseInsensitive(charCode) {
        if (compareCharCodeCaseInsensitive(this._cursor.peek(), charCode)) {
            this._cursor.advance();
            return true;
        }
        return false;
    }
    _requireCharCode(charCode) {
        const location = this._cursor.clone();
        if (!this._attemptCharCode(charCode)) {
            throw this._createError(_unexpectedCharacterErrorMsg(this._cursor.peek()), this._cursor.getSpan(location));
        }
    }
    _attemptStr(chars) {
        const len = chars.length;
        if (this._cursor.charsLeft() < len) {
            return false;
        }
        const initialPosition = this._cursor.clone();
        for (let i = 0; i < len; i++) {
            if (!this._attemptCharCode(chars.charCodeAt(i))) {
                // If attempting to parse the string fails, we want to reset the parser
                // to where it was before the attempt
                this._cursor = initialPosition;
                return false;
            }
        }
        return true;
    }
    _attemptStrCaseInsensitive(chars) {
        for (let i = 0; i < chars.length; i++) {
            if (!this._attemptCharCodeCaseInsensitive(chars.charCodeAt(i))) {
                return false;
            }
        }
        return true;
    }
    _requireStr(chars) {
        const location = this._cursor.clone();
        if (!this._attemptStr(chars)) {
            throw this._createError(_unexpectedCharacterErrorMsg(this._cursor.peek()), this._cursor.getSpan(location));
        }
    }
    _attemptCharCodeUntilFn(predicate) {
        while (!predicate(this._cursor.peek())) {
            this._cursor.advance();
        }
    }
    _requireCharCodeUntilFn(predicate, len) {
        const start = this._cursor.clone();
        this._attemptCharCodeUntilFn(predicate);
        if (this._cursor.diff(start) < len) {
            throw this._createError(_unexpectedCharacterErrorMsg(this._cursor.peek()), this._cursor.getSpan(start));
        }
    }
    _attemptUntilChar(char) {
        while (this._cursor.peek() !== char) {
            this._cursor.advance();
        }
    }
    _readChar() {
        // Don't rely upon reading directly from `_input` as the actual char value
        // may have been generated from an escape sequence.
        const char = String.fromCodePoint(this._cursor.peek());
        this._cursor.advance();
        return char;
    }
    _consumeEntity(textTokenType) {
        this._beginToken(9 /* TokenType.ENCODED_ENTITY */);
        const start = this._cursor.clone();
        this._cursor.advance();
        if (this._attemptCharCode(chars.$HASH)) {
            const isHex = this._attemptCharCode(chars.$x) || this._attemptCharCode(chars.$X);
            const codeStart = this._cursor.clone();
            this._attemptCharCodeUntilFn(isDigitEntityEnd);
            if (this._cursor.peek() != chars.$SEMICOLON) {
                // Advance cursor to include the peeked character in the string provided to the error
                // message.
                this._cursor.advance();
                const entityType = isHex ? CharacterReferenceType.HEX : CharacterReferenceType.DEC;
                throw this._createError(_unparsableEntityErrorMsg(entityType, this._cursor.getChars(start)), this._cursor.getSpan());
            }
            const strNum = this._cursor.getChars(codeStart);
            this._cursor.advance();
            try {
                const charCode = parseInt(strNum, isHex ? 16 : 10);
                this._endToken([String.fromCharCode(charCode), this._cursor.getChars(start)]);
            }
            catch {
                throw this._createError(_unknownEntityErrorMsg(this._cursor.getChars(start)), this._cursor.getSpan());
            }
        }
        else {
            const nameStart = this._cursor.clone();
            this._attemptCharCodeUntilFn(isNamedEntityEnd);
            if (this._cursor.peek() != chars.$SEMICOLON) {
                // No semicolon was found so abort the encoded entity token that was in progress, and treat
                // this as a text token
                this._beginToken(textTokenType, start);
                this._cursor = nameStart;
                this._endToken(['&']);
            }
            else {
                const name = this._cursor.getChars(nameStart);
                this._cursor.advance();
                const char = NAMED_ENTITIES[name];
                if (!char) {
                    throw this._createError(_unknownEntityErrorMsg(name), this._cursor.getSpan(start));
                }
                this._endToken([char, `&${name};`]);
            }
        }
    }
    _consumeRawText(consumeEntities, endMarkerPredicate) {
        this._beginToken(consumeEntities ? 6 /* TokenType.ESCAPABLE_RAW_TEXT */ : 7 /* TokenType.RAW_TEXT */);
        const parts = [];
        while (true) {
            const tagCloseStart = this._cursor.clone();
            const foundEndMarker = endMarkerPredicate();
            this._cursor = tagCloseStart;
            if (foundEndMarker) {
                break;
            }
            if (consumeEntities && this._cursor.peek() === chars.$AMPERSAND) {
                this._endToken([this._processCarriageReturns(parts.join(''))]);
                parts.length = 0;
                this._consumeEntity(6 /* TokenType.ESCAPABLE_RAW_TEXT */);
                this._beginToken(6 /* TokenType.ESCAPABLE_RAW_TEXT */);
            }
            else {
                parts.push(this._readChar());
            }
        }
        this._endToken([this._processCarriageReturns(parts.join(''))]);
    }
    _consumeComment(start) {
        this._beginToken(10 /* TokenType.COMMENT_START */, start);
        this._requireCharCode(chars.$MINUS);
        this._endToken([]);
        this._consumeRawText(false, () => this._attemptStr('-->'));
        this._beginToken(11 /* TokenType.COMMENT_END */);
        this._requireStr('-->');
        this._endToken([]);
    }
    _consumeCdata(start) {
        this._beginToken(12 /* TokenType.CDATA_START */, start);
        this._requireStr('CDATA[');
        this._endToken([]);
        this._consumeRawText(false, () => this._attemptStr(']]>'));
        this._beginToken(13 /* TokenType.CDATA_END */);
        this._requireStr(']]>');
        this._endToken([]);
    }
    _consumeDocType(start) {
        this._beginToken(18 /* TokenType.DOC_TYPE */, start);
        const contentStart = this._cursor.clone();
        this._attemptUntilChar(chars.$GT);
        const content = this._cursor.getChars(contentStart);
        this._cursor.advance();
        this._endToken([content]);
    }
    _consumePrefixAndName() {
        const nameOrPrefixStart = this._cursor.clone();
        let prefix = '';
        while (this._cursor.peek() !== chars.$COLON && !isPrefixEnd(this._cursor.peek())) {
            this._cursor.advance();
        }
        let nameStart;
        if (this._cursor.peek() === chars.$COLON) {
            prefix = this._cursor.getChars(nameOrPrefixStart);
            this._cursor.advance();
            nameStart = this._cursor.clone();
        }
        else {
            nameStart = nameOrPrefixStart;
        }
        this._requireCharCodeUntilFn(isNameEnd, prefix === '' ? 0 : 1);
        const name = this._cursor.getChars(nameStart);
        return [prefix, name];
    }
    _consumeTagOpen(start) {
        let tagName;
        let prefix;
        let openTagToken;
        try {
            if (!chars.isAsciiLetter(this._cursor.peek())) {
                throw this._createError(_unexpectedCharacterErrorMsg(this._cursor.peek()), this._cursor.getSpan(start));
            }
            openTagToken = this._consumeTagOpenStart(start);
            prefix = openTagToken.parts[0];
            tagName = openTagToken.parts[1];
            this._attemptCharCodeUntilFn(isNotWhitespace);
            while (this._cursor.peek() !== chars.$SLASH &&
                this._cursor.peek() !== chars.$GT &&
                this._cursor.peek() !== chars.$LT &&
                this._cursor.peek() !== chars.$EOF) {
                this._consumeAttributeName();
                this._attemptCharCodeUntilFn(isNotWhitespace);
                if (this._attemptCharCode(chars.$EQ)) {
                    this._attemptCharCodeUntilFn(isNotWhitespace);
                    this._consumeAttributeValue();
                }
                this._attemptCharCodeUntilFn(isNotWhitespace);
            }
            this._consumeTagOpenEnd();
        }
        catch (e) {
            if (e instanceof _ControlFlowError) {
                if (openTagToken) {
                    // We errored before we could close the opening tag, so it is incomplete.
                    openTagToken.type = 4 /* TokenType.INCOMPLETE_TAG_OPEN */;
                }
                else {
                    // When the start tag is invalid, assume we want a "<" as text.
                    // Back to back text tokens are merged at the end.
                    this._beginToken(5 /* TokenType.TEXT */, start);
                    this._endToken(['<']);
                }
                return;
            }
            throw e;
        }
        const contentTokenType = this._getTagDefinition(tagName).getContentType(prefix);
        if (contentTokenType === TagContentType.RAW_TEXT) {
            this._consumeRawTextWithTagClose(prefix, tagName, false);
        }
        else if (contentTokenType === TagContentType.ESCAPABLE_RAW_TEXT) {
            this._consumeRawTextWithTagClose(prefix, tagName, true);
        }
    }
    _consumeRawTextWithTagClose(prefix, tagName, consumeEntities) {
        this._consumeRawText(consumeEntities, () => {
            if (!this._attemptCharCode(chars.$LT))
                return false;
            if (!this._attemptCharCode(chars.$SLASH))
                return false;
            this._attemptCharCodeUntilFn(isNotWhitespace);
            if (!this._attemptStrCaseInsensitive(tagName))
                return false;
            this._attemptCharCodeUntilFn(isNotWhitespace);
            return this._attemptCharCode(chars.$GT);
        });
        this._beginToken(3 /* TokenType.TAG_CLOSE */);
        this._requireCharCodeUntilFn((code) => code === chars.$GT, 3);
        this._cursor.advance(); // Consume the `>`
        this._endToken([prefix, tagName]);
    }
    _consumeTagOpenStart(start) {
        this._beginToken(0 /* TokenType.TAG_OPEN_START */, start);
        const parts = this._consumePrefixAndName();
        return this._endToken(parts);
    }
    _consumeAttributeName() {
        const attrNameStart = this._cursor.peek();
        if (attrNameStart === chars.$SQ || attrNameStart === chars.$DQ) {
            throw this._createError(_unexpectedCharacterErrorMsg(attrNameStart), this._cursor.getSpan());
        }
        this._beginToken(14 /* TokenType.ATTR_NAME */);
        const prefixAndName = this._consumePrefixAndName();
        this._endToken(prefixAndName);
    }
    _consumeAttributeValue() {
        if (this._cursor.peek() === chars.$SQ || this._cursor.peek() === chars.$DQ) {
            const quoteChar = this._cursor.peek();
            this._consumeQuote(quoteChar);
            // In an attribute then end of the attribute value and the premature end to an interpolation
            // are both triggered by the `quoteChar`.
            const endPredicate = () => this._cursor.peek() === quoteChar;
            this._consumeWithInterpolation(16 /* TokenType.ATTR_VALUE_TEXT */, 17 /* TokenType.ATTR_VALUE_INTERPOLATION */, endPredicate, endPredicate);
            this._consumeQuote(quoteChar);
        }
        else {
            const endPredicate = () => isNameEnd(this._cursor.peek());
            this._consumeWithInterpolation(16 /* TokenType.ATTR_VALUE_TEXT */, 17 /* TokenType.ATTR_VALUE_INTERPOLATION */, endPredicate, endPredicate);
        }
    }
    _consumeQuote(quoteChar) {
        this._beginToken(15 /* TokenType.ATTR_QUOTE */);
        this._requireCharCode(quoteChar);
        this._endToken([String.fromCodePoint(quoteChar)]);
    }
    _consumeTagOpenEnd() {
        const tokenType = this._attemptCharCode(chars.$SLASH)
            ? 2 /* TokenType.TAG_OPEN_END_VOID */
            : 1 /* TokenType.TAG_OPEN_END */;
        this._beginToken(tokenType);
        this._requireCharCode(chars.$GT);
        this._endToken([]);
    }
    _consumeTagClose(start) {
        this._beginToken(3 /* TokenType.TAG_CLOSE */, start);
        this._attemptCharCodeUntilFn(isNotWhitespace);
        const prefixAndName = this._consumePrefixAndName();
        this._attemptCharCodeUntilFn(isNotWhitespace);
        this._requireCharCode(chars.$GT);
        this._endToken(prefixAndName);
    }
    _consumeExpansionFormStart() {
        this._beginToken(19 /* TokenType.EXPANSION_FORM_START */);
        this._requireCharCode(chars.$LBRACE);
        this._endToken([]);
        this._expansionCaseStack.push(19 /* TokenType.EXPANSION_FORM_START */);
        this._beginToken(7 /* TokenType.RAW_TEXT */);
        const condition = this._readUntil(chars.$COMMA);
        const normalizedCondition = this._processCarriageReturns(condition);
        if (this._i18nNormalizeLineEndingsInICUs) {
            // We explicitly want to normalize line endings for this text.
            this._endToken([normalizedCondition]);
        }
        else {
            // We are not normalizing line endings.
            const conditionToken = this._endToken([condition]);
            if (normalizedCondition !== condition) {
                this.nonNormalizedIcuExpressions.push(conditionToken);
            }
        }
        this._requireCharCode(chars.$COMMA);
        this._attemptCharCodeUntilFn(isNotWhitespace);
        this._beginToken(7 /* TokenType.RAW_TEXT */);
        const type = this._readUntil(chars.$COMMA);
        this._endToken([type]);
        this._requireCharCode(chars.$COMMA);
        this._attemptCharCodeUntilFn(isNotWhitespace);
    }
    _consumeExpansionCaseStart() {
        this._beginToken(20 /* TokenType.EXPANSION_CASE_VALUE */);
        const value = this._readUntil(chars.$LBRACE).trim();
        this._endToken([value]);
        this._attemptCharCodeUntilFn(isNotWhitespace);
        this._beginToken(21 /* TokenType.EXPANSION_CASE_EXP_START */);
        this._requireCharCode(chars.$LBRACE);
        this._endToken([]);
        this._attemptCharCodeUntilFn(isNotWhitespace);
        this._expansionCaseStack.push(21 /* TokenType.EXPANSION_CASE_EXP_START */);
    }
    _consumeExpansionCaseEnd() {
        this._beginToken(22 /* TokenType.EXPANSION_CASE_EXP_END */);
        this._requireCharCode(chars.$RBRACE);
        this._endToken([]);
        this._attemptCharCodeUntilFn(isNotWhitespace);
        this._expansionCaseStack.pop();
    }
    _consumeExpansionFormEnd() {
        this._beginToken(23 /* TokenType.EXPANSION_FORM_END */);
        this._requireCharCode(chars.$RBRACE);
        this._endToken([]);
        this._expansionCaseStack.pop();
    }
    /**
     * Consume a string that may contain interpolation expressions.
     *
     * The first token consumed will be of `tokenType` and then there will be alternating
     * `interpolationTokenType` and `tokenType` tokens until the `endPredicate()` returns true.
     *
     * If an interpolation token ends prematurely it will have no end marker in its `parts` array.
     *
     * @param textTokenType the kind of tokens to interleave around interpolation tokens.
     * @param interpolationTokenType the kind of tokens that contain interpolation.
     * @param endPredicate a function that should return true when we should stop consuming.
     * @param endInterpolation a function that should return true if there is a premature end to an
     *     interpolation expression - i.e. before we get to the normal interpolation closing marker.
     */
    _consumeWithInterpolation(textTokenType, interpolationTokenType, endPredicate, endInterpolation) {
        this._beginToken(textTokenType);
        const parts = [];
        while (!endPredicate()) {
            const current = this._cursor.clone();
            if (this._interpolationConfig && this._attemptStr(this._interpolationConfig.start)) {
                this._endToken([this._processCarriageReturns(parts.join(''))], current);
                parts.length = 0;
                this._consumeInterpolation(interpolationTokenType, current, endInterpolation);
                this._beginToken(textTokenType);
            }
            else if (this._cursor.peek() === chars.$AMPERSAND) {
                this._endToken([this._processCarriageReturns(parts.join(''))]);
                parts.length = 0;
                this._consumeEntity(textTokenType);
                this._beginToken(textTokenType);
            }
            else {
                parts.push(this._readChar());
            }
        }
        // It is possible that an interpolation was started but not ended inside this text token.
        // Make sure that we reset the state of the lexer correctly.
        this._inInterpolation = false;
        this._endToken([this._processCarriageReturns(parts.join(''))]);
    }
    /**
     * Consume a block of text that has been interpreted as an Angular interpolation.
     *
     * @param interpolationTokenType the type of the interpolation token to generate.
     * @param interpolationStart a cursor that points to the start of this interpolation.
     * @param prematureEndPredicate a function that should return true if the next characters indicate
     *     an end to the interpolation before its normal closing marker.
     */
    _consumeInterpolation(interpolationTokenType, interpolationStart, prematureEndPredicate) {
        const parts = [];
        this._beginToken(interpolationTokenType, interpolationStart);
        parts.push(this._interpolationConfig.start);
        // Find the end of the interpolation, ignoring content inside quotes.
        const expressionStart = this._cursor.clone();
        let inQuote = null;
        let inComment = false;
        while (this._cursor.peek() !== chars.$EOF &&
            (prematureEndPredicate === null || !prematureEndPredicate())) {
            const current = this._cursor.clone();
            if (this._isTagStart()) {
                // We are starting what looks like an HTML element in the middle of this interpolation.
                // Reset the cursor to before the `<` character and end the interpolation token.
                // (This is actually wrong but here for backward compatibility).
                this._cursor = current;
                parts.push(this._getProcessedChars(expressionStart, current));
                this._endToken(parts);
                return;
            }
            if (inQuote === null) {
                if (this._attemptStr(this._interpolationConfig.end)) {
                    // We are not in a string, and we hit the end interpolation marker
                    parts.push(this._getProcessedChars(expressionStart, current));
                    parts.push(this._interpolationConfig.end);
                    this._endToken(parts);
                    return;
                }
                else if (this._attemptStr('//')) {
                    // Once we are in a comment we ignore any quotes
                    inComment = true;
                }
            }
            const char = this._cursor.peek();
            this._cursor.advance();
            if (char === chars.$BACKSLASH) {
                // Skip the next character because it was escaped.
                this._cursor.advance();
            }
            else if (char === inQuote) {
                // Exiting the current quoted string
                inQuote = null;
            }
            else if (!inComment && inQuote === null && chars.isQuote(char)) {
                // Entering a new quoted string
                inQuote = char;
            }
        }
        // We hit EOF without finding a closing interpolation marker
        parts.push(this._getProcessedChars(expressionStart, this._cursor));
        this._endToken(parts);
    }
    _getProcessedChars(start, end) {
        return this._processCarriageReturns(end.getChars(start));
    }
    _isTextEnd() {
        if (this._isTagStart() || this._cursor.peek() === chars.$EOF) {
            return true;
        }
        if (this._tokenizeIcu && !this._inInterpolation) {
            if (this.isExpansionFormStart()) {
                // start of an expansion form
                return true;
            }
            if (this._cursor.peek() === chars.$RBRACE && this._isInExpansionCase()) {
                // end of and expansion case
                return true;
            }
        }
        if (this._tokenizeBlocks &&
            !this._inInterpolation &&
            !this._isInExpansion() &&
            (this._cursor.peek() === chars.$AT || this._cursor.peek() === chars.$RBRACE)) {
            return true;
        }
        return false;
    }
    /**
     * Returns true if the current cursor is pointing to the start of a tag
     * (opening/closing/comments/cdata/etc).
     */
    _isTagStart() {
        if (this._cursor.peek() === chars.$LT) {
            // We assume that `<` followed by whitespace is not the start of an HTML element.
            const tmp = this._cursor.clone();
            tmp.advance();
            // If the next character is alphabetic, ! nor / then it is a tag start
            const code = tmp.peek();
            if ((chars.$a <= code && code <= chars.$z) ||
                (chars.$A <= code && code <= chars.$Z) ||
                code === chars.$SLASH ||
                code === chars.$BANG) {
                return true;
            }
        }
        return false;
    }
    _readUntil(char) {
        const start = this._cursor.clone();
        this._attemptUntilChar(char);
        return this._cursor.getChars(start);
    }
    _isInExpansion() {
        return this._isInExpansionCase() || this._isInExpansionForm();
    }
    _isInExpansionCase() {
        return (this._expansionCaseStack.length > 0 &&
            this._expansionCaseStack[this._expansionCaseStack.length - 1] ===
                21 /* TokenType.EXPANSION_CASE_EXP_START */);
    }
    _isInExpansionForm() {
        return (this._expansionCaseStack.length > 0 &&
            this._expansionCaseStack[this._expansionCaseStack.length - 1] ===
                19 /* TokenType.EXPANSION_FORM_START */);
    }
    isExpansionFormStart() {
        if (this._cursor.peek() !== chars.$LBRACE) {
            return false;
        }
        if (this._interpolationConfig) {
            const start = this._cursor.clone();
            const isInterpolation = this._attemptStr(this._interpolationConfig.start);
            this._cursor = start;
            return !isInterpolation;
        }
        return true;
    }
}
function isNotWhitespace(code) {
    return !chars.isWhitespace(code) || code === chars.$EOF;
}
function isNameEnd(code) {
    return (chars.isWhitespace(code) ||
        code === chars.$GT ||
        code === chars.$LT ||
        code === chars.$SLASH ||
        code === chars.$SQ ||
        code === chars.$DQ ||
        code === chars.$EQ ||
        code === chars.$EOF);
}
function isPrefixEnd(code) {
    return ((code < chars.$a || chars.$z < code) &&
        (code < chars.$A || chars.$Z < code) &&
        (code < chars.$0 || code > chars.$9));
}
function isDigitEntityEnd(code) {
    return code === chars.$SEMICOLON || code === chars.$EOF || !chars.isAsciiHexDigit(code);
}
function isNamedEntityEnd(code) {
    return code === chars.$SEMICOLON || code === chars.$EOF || !chars.isAsciiLetter(code);
}
function isExpansionCaseStart(peek) {
    return peek !== chars.$RBRACE;
}
function compareCharCodeCaseInsensitive(code1, code2) {
    return toUpperCaseCharCode(code1) === toUpperCaseCharCode(code2);
}
function toUpperCaseCharCode(code) {
    return code >= chars.$a && code <= chars.$z ? code - chars.$a + chars.$A : code;
}
function isBlockNameChar(code) {
    return chars.isAsciiLetter(code) || chars.isDigit(code) || code === chars.$_;
}
function isBlockParameterChar(code) {
    return code !== chars.$SEMICOLON && isNotWhitespace(code);
}
function mergeTextTokens(srcTokens) {
    const dstTokens = [];
    let lastDstToken = undefined;
    for (let i = 0; i < srcTokens.length; i++) {
        const token = srcTokens[i];
        if ((lastDstToken && lastDstToken.type === 5 /* TokenType.TEXT */ && token.type === 5 /* TokenType.TEXT */) ||
            (lastDstToken &&
                lastDstToken.type === 16 /* TokenType.ATTR_VALUE_TEXT */ &&
                token.type === 16 /* TokenType.ATTR_VALUE_TEXT */)) {
            lastDstToken.parts[0] += token.parts[0];
            lastDstToken.sourceSpan.end = token.sourceSpan.end;
        }
        else {
            lastDstToken = token;
            dstTokens.push(lastDstToken);
        }
    }
    return dstTokens;
}
class PlainCharacterCursor {
    constructor(fileOrCursor, range) {
        if (fileOrCursor instanceof PlainCharacterCursor) {
            this.file = fileOrCursor.file;
            this.input = fileOrCursor.input;
            this.end = fileOrCursor.end;
            const state = fileOrCursor.state;
            // Note: avoid using `{...fileOrCursor.state}` here as that has a severe performance penalty.
            // In ES5 bundles the object spread operator is translated into the `__assign` helper, which
            // is not optimized by VMs as efficiently as a raw object literal. Since this constructor is
            // called in tight loops, this difference matters.
            this.state = {
                peek: state.peek,
                offset: state.offset,
                line: state.line,
                column: state.column,
            };
        }
        else {
            if (!range) {
                throw new Error('Programming error: the range argument must be provided with a file argument.');
            }
            this.file = fileOrCursor;
            this.input = fileOrCursor.content;
            this.end = range.endPos;
            this.state = {
                peek: -1,
                offset: range.startPos,
                line: range.startLine,
                column: range.startCol,
            };
        }
    }
    clone() {
        return new PlainCharacterCursor(this);
    }
    peek() {
        return this.state.peek;
    }
    charsLeft() {
        return this.end - this.state.offset;
    }
    diff(other) {
        return this.state.offset - other.state.offset;
    }
    advance() {
        this.advanceState(this.state);
    }
    init() {
        this.updatePeek(this.state);
    }
    getSpan(start, leadingTriviaCodePoints) {
        start = start || this;
        let fullStart = start;
        if (leadingTriviaCodePoints) {
            while (this.diff(start) > 0 && leadingTriviaCodePoints.indexOf(start.peek()) !== -1) {
                if (fullStart === start) {
                    start = start.clone();
                }
                start.advance();
            }
        }
        const startLocation = this.locationFromCursor(start);
        const endLocation = this.locationFromCursor(this);
        const fullStartLocation = fullStart !== start ? this.locationFromCursor(fullStart) : startLocation;
        return new ParseSourceSpan(startLocation, endLocation, fullStartLocation);
    }
    getChars(start) {
        return this.input.substring(start.state.offset, this.state.offset);
    }
    charAt(pos) {
        return this.input.charCodeAt(pos);
    }
    advanceState(state) {
        if (state.offset >= this.end) {
            this.state = state;
            throw new CursorError('Unexpected character "EOF"', this);
        }
        const currentChar = this.charAt(state.offset);
        if (currentChar === chars.$LF) {
            state.line++;
            state.column = 0;
        }
        else if (!chars.isNewLine(currentChar)) {
            state.column++;
        }
        state.offset++;
        this.updatePeek(state);
    }
    updatePeek(state) {
        state.peek = state.offset >= this.end ? chars.$EOF : this.charAt(state.offset);
    }
    locationFromCursor(cursor) {
        return new ParseLocation(cursor.file, cursor.state.offset, cursor.state.line, cursor.state.column);
    }
}
class EscapedCharacterCursor extends PlainCharacterCursor {
    constructor(fileOrCursor, range) {
        if (fileOrCursor instanceof EscapedCharacterCursor) {
            super(fileOrCursor);
            this.internalState = { ...fileOrCursor.internalState };
        }
        else {
            super(fileOrCursor, range);
            this.internalState = this.state;
        }
    }
    advance() {
        this.state = this.internalState;
        super.advance();
        this.processEscapeSequence();
    }
    init() {
        super.init();
        this.processEscapeSequence();
    }
    clone() {
        return new EscapedCharacterCursor(this);
    }
    getChars(start) {
        const cursor = start.clone();
        let chars = '';
        while (cursor.internalState.offset < this.internalState.offset) {
            chars += String.fromCodePoint(cursor.peek());
            cursor.advance();
        }
        return chars;
    }
    /**
     * Process the escape sequence that starts at the current position in the text.
     *
     * This method is called to ensure that `peek` has the unescaped value of escape sequences.
     */
    processEscapeSequence() {
        const peek = () => this.internalState.peek;
        if (peek() === chars.$BACKSLASH) {
            // We have hit an escape sequence so we need the internal state to become independent
            // of the external state.
            this.internalState = { ...this.state };
            // Move past the backslash
            this.advanceState(this.internalState);
            // First check for standard control char sequences
            if (peek() === chars.$n) {
                this.state.peek = chars.$LF;
            }
            else if (peek() === chars.$r) {
                this.state.peek = chars.$CR;
            }
            else if (peek() === chars.$v) {
                this.state.peek = chars.$VTAB;
            }
            else if (peek() === chars.$t) {
                this.state.peek = chars.$TAB;
            }
            else if (peek() === chars.$b) {
                this.state.peek = chars.$BSPACE;
            }
            else if (peek() === chars.$f) {
                this.state.peek = chars.$FF;
            }
            // Now consider more complex sequences
            else if (peek() === chars.$u) {
                // Unicode code-point sequence
                this.advanceState(this.internalState); // advance past the `u` char
                if (peek() === chars.$LBRACE) {
                    // Variable length Unicode, e.g. `\x{123}`
                    this.advanceState(this.internalState); // advance past the `{` char
                    // Advance past the variable number of hex digits until we hit a `}` char
                    const digitStart = this.clone();
                    let length = 0;
                    while (peek() !== chars.$RBRACE) {
                        this.advanceState(this.internalState);
                        length++;
                    }
                    this.state.peek = this.decodeHexDigits(digitStart, length);
                }
                else {
                    // Fixed length Unicode, e.g. `\u1234`
                    const digitStart = this.clone();
                    this.advanceState(this.internalState);
                    this.advanceState(this.internalState);
                    this.advanceState(this.internalState);
                    this.state.peek = this.decodeHexDigits(digitStart, 4);
                }
            }
            else if (peek() === chars.$x) {
                // Hex char code, e.g. `\x2F`
                this.advanceState(this.internalState); // advance past the `x` char
                const digitStart = this.clone();
                this.advanceState(this.internalState);
                this.state.peek = this.decodeHexDigits(digitStart, 2);
            }
            else if (chars.isOctalDigit(peek())) {
                // Octal char code, e.g. `\012`,
                let octal = '';
                let length = 0;
                let previous = this.clone();
                while (chars.isOctalDigit(peek()) && length < 3) {
                    previous = this.clone();
                    octal += String.fromCodePoint(peek());
                    this.advanceState(this.internalState);
                    length++;
                }
                this.state.peek = parseInt(octal, 8);
                // Backup one char
                this.internalState = previous.internalState;
            }
            else if (chars.isNewLine(this.internalState.peek)) {
                // Line continuation `\` followed by a new line
                this.advanceState(this.internalState); // advance over the newline
                this.state = this.internalState;
            }
            else {
                // If none of the `if` blocks were executed then we just have an escaped normal character.
                // In that case we just, effectively, skip the backslash from the character.
                this.state.peek = this.internalState.peek;
            }
        }
    }
    decodeHexDigits(start, length) {
        const hex = this.input.slice(start.internalState.offset, start.internalState.offset + length);
        const charCode = parseInt(hex, 16);
        if (!isNaN(charCode)) {
            return charCode;
        }
        else {
            start.state = start.internalState;
            throw new CursorError('Invalid hexadecimal escape sequence', start);
        }
    }
}
export class CursorError {
    constructor(msg, cursor) {
        this.msg = msg;
        this.cursor = cursor;
    }
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibGV4ZXIuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi8uLi8uLi8uLi9wYWNrYWdlcy9jb21waWxlci9zcmMvbWxfcGFyc2VyL2xleGVyLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBOzs7Ozs7R0FNRztBQUVILE9BQU8sS0FBSyxLQUFLLE1BQU0sVUFBVSxDQUFDO0FBQ2xDLE9BQU8sRUFBQyxVQUFVLEVBQUUsYUFBYSxFQUFFLGVBQWUsRUFBRSxlQUFlLEVBQUMsTUFBTSxlQUFlLENBQUM7QUFFMUYsT0FBTyxFQUFDLDRCQUE0QixFQUFzQixNQUFNLFlBQVksQ0FBQztBQUM3RSxPQUFPLEVBQUMsY0FBYyxFQUFDLE1BQU0sWUFBWSxDQUFDO0FBQzFDLE9BQU8sRUFBQyxjQUFjLEVBQWdCLE1BQU0sUUFBUSxDQUFDO0FBR3JELE1BQU0sT0FBTyxVQUFXLFNBQVEsVUFBVTtJQUN4QyxZQUNFLFFBQWdCLEVBQ1QsU0FBMkIsRUFDbEMsSUFBcUI7UUFFckIsS0FBSyxDQUFDLElBQUksRUFBRSxRQUFRLENBQUMsQ0FBQztRQUhmLGNBQVMsR0FBVCxTQUFTLENBQWtCO0lBSXBDLENBQUM7Q0FDRjtBQUVELE1BQU0sT0FBTyxjQUFjO0lBQ3pCLFlBQ1MsTUFBZSxFQUNmLE1BQW9CLEVBQ3BCLDJCQUFvQztRQUZwQyxXQUFNLEdBQU4sTUFBTSxDQUFTO1FBQ2YsV0FBTSxHQUFOLE1BQU0sQ0FBYztRQUNwQixnQ0FBMkIsR0FBM0IsMkJBQTJCLENBQVM7SUFDMUMsQ0FBQztDQUNMO0FBK0VELE1BQU0sVUFBVSxRQUFRLENBQ3RCLE1BQWMsRUFDZCxHQUFXLEVBQ1gsZ0JBQW9ELEVBQ3BELFVBQTJCLEVBQUU7SUFFN0IsTUFBTSxTQUFTLEdBQUcsSUFBSSxVQUFVLENBQUMsSUFBSSxlQUFlLENBQUMsTUFBTSxFQUFFLEdBQUcsQ0FBQyxFQUFFLGdCQUFnQixFQUFFLE9BQU8sQ0FBQyxDQUFDO0lBQzlGLFNBQVMsQ0FBQyxRQUFR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QUFDLGFBQWEsQ0FBQztZQUNsQyxNQUFNLElBQUksV0FBVyxDQUFDLHFDQUFxQyxFQUFFLEtBQUssQ0FBQyxDQUFDO1FBQ3RFLENBQUM7SUFDSCxDQUFDO0NBQ0Y7QUFFRCxNQUFNLE9BQU8sV0FBVztJQUN0QixZQUNTLEdBQVcsRUFDWCxNQUF1QjtRQUR2QixRQUFHLEdBQUgsR0FBRyxDQUFRO1FBQ1gsV0FBTSxHQUFOLE1BQU0sQ0FBaUI7SUFDN0IsQ0FBQztDQUNMIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAbGljZW5zZVxuICogQ29weXJpZ2h0IEdvb2dsZSBMTEMgQWxsIFJpZ2h0cyBSZXNlcnZlZC5cbiAqXG4gKiBVc2Ugb2YgdGhpcyBzb3VyY2UgY29kZSBpcyBnb3Zlcm5lZCBieSBhbiBNSVQtc3R5bGUgbGljZW5zZSB0aGF0IGNhbiBiZVxuICogZm91bmQgaW4gdGhlIExJQ0VOU0UgZmlsZSBhdCBodHRwczovL2FuZ3VsYXIuaW8vbGljZW5zZVxuICovXG5cbmltcG9ydCAqIGFzIGNoYXJzIGZyb20gJy4uL2NoYXJzJztcbmltcG9ydCB7UGFyc2VFcnJvciwgUGFyc2VMb2NhdGlvbiwgUGFyc2VTb3VyY2VGaWxlLCBQYXJzZVNvdXJjZVNwYW59IGZyb20gJy4uL3BhcnNlX3V0aWwnO1xuXG5pbXBvcnQge0RFRkFVTFRfSU5URVJQT0xBVElPTl9DT05GSUcsIEludGVycG9sYXRpb25Db25maWd9IGZyb20gJy4vZGVmYXVsdHMnO1xuaW1wb3J0IHtOQU1FRF9FTlRJVElFU30gZnJvbSAnLi9lbnRpdGllcyc7XG5pbXBvcnQge1RhZ0NvbnRlbnRUeXBlLCBUYWdEZWZpbml0aW9ufSBmcm9tICcuL3RhZ3MnO1xuaW1wb3J0IHtJbmNvbXBsZXRlVGFnT3BlblRva2VuLCBUYWdPcGVuU3RhcnRUb2tlbiwgVG9rZW4sIFRva2VuVHlwZX0gZnJvbSAnLi90b2tlbnMnO1xuXG5leHBvcnQgY2xhc3MgVG9rZW5FcnJvciBleHRlbmRzIFBhcnNlRXJyb3Ige1xuICBjb25zdHJ1Y3RvcihcbiAgICBlcnJvck1zZzogc3RyaW5nLFxuICAgIHB1YmxpYyB0b2tlblR5cGU6IFRva2VuVHlwZSB8IG51bGwsXG4gICAgc3BhbjogUGFyc2VTb3VyY2VTcGFuLFxuICApIHtcbiAgICBzdXBlcihzcGFuLCBlcnJvck1zZyk7XG4gIH1cbn1cblxuZXhwb3J0IGNsYXNzIFRva2VuaXplUmVzdWx0IHtcbiAgY29uc3RydWN0b3IoXG4gICAgcHVibGljIHRva2VuczogVG9rZW5bXSxcbiAgICBwdWJsaWMgZXJyb3JzOiBUb2tlbkVycm9yW10sXG4gICAgcHVibGljIG5vbk5vcm1hbGl6ZWRJY3VFeHByZXNzaW9uczogVG9rZW5bXSxcbiAgKSB7fVxufVxuXG5leHBvcnQgaW50ZXJmYWNlIExleGVyUmFuZ2Uge1xuICBzdGFydFBvczogbnVtYmVyO1xuICBzdGFydExpbmU6IG51bWJlcjtcbiAgc3RhcnRDb2w6IG51bWJlcjtcbiAgZW5kUG9zOiBudW1iZXI7XG59XG5cbi8qKlxuICogT3B0aW9ucyB0aGF0IG1vZGlmeSBob3cgdGhlIHRleHQgaXMgdG9rZW5pemVkLlxuICovXG5leHBvcnQgaW50ZXJmYWNlIFRva2VuaXplT3B0aW9ucyB7XG4gIC8qKiBXaGV0aGVyIHRvIHRva2VuaXplIElDVSBtZXNzYWdlcyAoY29uc2lkZXJlZCBhcyB0ZXh0IG5vZGVzIHdoZW4gZmFsc2UpLiAqL1xuICB0b2tlbml6ZUV4cGFuc2lvbkZvcm1zPzogYm9vbGVhbjtcbiAgLyoqIEhvdyB0byB0b2tlbml6ZSBpbnRlcnBvbGF0aW9uIG1hcmtlcnMuICovXG4gIGludGVycG9sYXRpb25Db25maWc/OiBJbnRlcnBvbGF0aW9uQ29uZmlnO1xuICAvKipcbiAgICogVGhlIHN0YXJ0IGFuZCBlbmQgcG9pbnQgb2YgdGhlIHRleHQgdG8gcGFyc2Ugd2l0aGluIHRoZSBgc291cmNlYCBzdHJpbmcuXG4gICAqIFRoZSBlbnRpcmUgYHNvdXJjZWAgc3RyaW5nIGlzIHBhcnNlZCBpZiB0aGlzIGlzIG5vdCBwcm92aWRlZC5cbiAgICogKi9cbiAgcmFuZ2U/OiBMZXhlclJhbmdlO1xuICAvKipcbiAgICogSWYgdGhpcyB0ZXh0IGlzIHN0b3JlZCBpbiBhIEphdmFTY3JpcHQgc3RyaW5nLCB0aGVuIHdlIGhhdmUgdG8gZGVhbCB3aXRoIGVzY2FwZSBzZXF1ZW5jZXMuXG4gICAqXG4gICAqICoqRXhhbXBsZSAxOioqXG4gICAqXG4gICAqIGBgYFxuICAgKiBcImFiY1xcXCJkZWZcXG5naGlcIlxuICAgKiBgYGBcbiAgICpcbiAgICogLSBUaGUgYFxcXCJgIG11c3QgYmUgY29udmVydGVkIHRvIGBcImAuXG4gICAqIC0gVGhlIGBcXG5gIG11c3QgYmUgY29udmVydGVkIHRvIGEgbmV3IGxpbmUgY2hhcmFjdGVyIGluIGEgdG9rZW4sXG4gICAqICAgYnV0IGl0IHNob3VsZCBub3QgaW5jcmVtZW50IHRoZSBjdXJyZW50IGxpbmUgZm9yIHNvdXJjZSBtYXBwaW5nLlxuICAgKlxuICAgKiAqKkV4YW1wbGUgMjoqKlxuICAgKlxuICAgKiBgYGBcbiAgICogXCJhYmNcXFxuICAgKiAgZGVmXCJcbiAgICogYGBgXG4gICAqXG4gICAqIFRoZSBsaW5lIGNvbnRpbnVhdGlvbiAoYFxcYCBmb2xsb3dlZCBieSBhIG5ld2xpbmUpIHNob3VsZCBiZSByZW1vdmVkIGZyb20gYSB0b2tlblxuICAgKiBidXQgdGhlIG5ldyBsaW5lIHNob3VsZCBpbmNyZW1lbnQgdGhlIGN1cnJlbnQgbGluZSBmb3Igc291cmNlIG1hcHBpbmcuXG4gICAqL1xuICBlc2NhcGVkU3RyaW5nPzogYm9vbGVhbjtcbiAgLyoqXG4gICAqIElmIHRoaXMgdGV4dCBpcyBzdG9yZWQgaW4gYW4gZXh0ZXJuYWwgdGVtcGxhdGUgKGUuZy4gdmlhIGB0ZW1wbGF0ZVVybGApIHRoZW4gd2UgbmVlZCB0byBkZWNpZGVcbiAgICogd2hldGhlciBvciBub3QgdG8gbm9ybWFsaXplIHRoZSBsaW5lLWVuZGluZ3MgKGZyb20gYFxcclxcbmAgdG8gYFxcbmApIHdoZW4gcHJvY2Vzc2luZyBJQ1VcbiAgICogZXhwcmVzc2lvbnMuXG4gICAqXG4gICAqIElmIGB0cnVlYCB0aGVuIHdlIHdpbGwgbm9ybWFsaXplIElDVSBleHByZXNzaW9uIGxpbmUgZW5kaW5ncy5cbiAgICogVGhlIGRlZmF1bHQgaXMgYGZhbHNlYCwgYnV0IHRoaXMgd2lsbCBiZSBzd2l0Y2hlZCBpbiBhIGZ1dHVyZSBtYWpvciByZWxlYXNlLlxuICAgKi9cbiAgaTE4bk5vcm1hbGl6ZUxpbmVFbmRpbmdzSW5JQ1VzPzogYm9vbGVhbjtcbiAgLyoqXG4gICAqIEFuIGFycmF5IG9mIGNoYXJhY3RlcnMgdGhhdCBzaG91bGQgYmUgY29uc2lkZXJlZCBhcyBsZWFkaW5nIHRyaXZpYS5cbiAgICogTGVhZGluZyB0cml2aWEgYXJlIGNoYXJhY3RlcnMgdGhhdCBhcmUgbm90IGltcG9ydGFudCB0byB0aGUgZGV2ZWxvcGVyLCBhbmQgc28gc2hvdWxkIG5vdCBiZVxuICAgKiBpbmNsdWRlZCBpbiBzb3VyY2UtbWFwIHNlZ21lbnRzLiAgQSBjb21tb24gZXhhbXBsZSBpcyB3aGl0ZXNwYWNlLlxuICAgKi9cbiAgbGVhZGluZ1RyaXZpYUNoYXJzPzogc3RyaW5nW107XG4gIC8qKlxuICAgKiBJZiB0cnVlLCBkbyBub3QgY29udmVydCBDUkxGIHRvIExGLlxuICAgKi9cbiAgcHJlc2VydmVMaW5lRW5kaW5ncz86IGJvb2xlYW47XG5cbiAgLyoqXG4gICAqIFdoZXRoZXIgdG8gdG9rZW5pemUgQCBibG9jayBzeW50YXguIE90aGVyd2lzZSBjb25zaWRlcmVkIHRleHQsXG4gICAqIG9yIElDVSB0b2tlbnMgaWYgYHRva2VuaXplRXhwYW5zaW9uRm9ybXNgIGlzIGVuYWJsZWQuXG4gICAqL1xuICB0b2tlbml6ZUJsb2Nrcz86IGJvb2xlYW47XG5cbiAgLyoqXG4gICAqIFdoZXRoZXIgdG8gdG9rZW5pemUgdGhlIGBAbGV0YCBzeW50YXguIE90aGVyd2lzZSB3aWxsIGJlIGNvbnNpZGVyZWQgZWl0aGVyXG4gICAqIHRleHQgb3IgYW4gaW5jb21wbGV0ZSBibG9jaywgZGVwZW5kaW5nIG9uIHdoZXRoZXIgYHRva2VuaXplQmxvY2tzYCBpcyBlbmFibGVkLlxuICAgKi9cbiAgdG9rZW5pemVMZXQ/OiBib29sZWFuO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gdG9rZW5pemUoXG4gIHNvdXJjZTogc3RyaW5nLFxuICB1cmw6IHN0cmluZyxcbiAgZ2V0VGFnRGVmaW5pdGlvbjogKHRhZ05hbWU6IHN0cmluZykgPT4gVGFnRGVmaW5pdGlvbixcbiAgb3B0aW9uczogVG9rZW5pemVPcHRpb25zID0ge30sXG4pOiBUb2tlbml6ZVJlc3VsdCB7XG4gIGNvbnN0IHRva2VuaXplciA9IG5ldyBfVG9rZW5pemVyKG5ldyBQYXJzZVNvdXJjZUZpbGUoc291cmNlLCB1cmwpLCBnZXRUYWdEZWZpbml0aW9uLCBvcHRpb25zKTtcbiAgdG9rZW5pemVyLnRva2VuaXplKCk7XG4gIHJldHVybiBuZXcgVG9rZW5pemVSZXN1bHQoXG4gICAgbWVyZ2VUZXh0VG9rZW5zKHRva2VuaXplci50b2tlbnMpLFxuICAgIHRva2VuaXplci5lcnJvcnMsXG4gICAgdG9rZW5pemVyLm5vbk5vcm1hbGl6ZWRJY3VFeHByZXNzaW9ucyxcbiAgKTtcbn1cblxuY29uc3QgX0NSX09SX0NSTEZfUkVHRVhQID0gL1xcclxcbj8vZztcblxuZnVuY3Rpb24gX3VuZXhwZWN0ZWRDaGFyYWN0ZXJFcnJvck1zZyhjaGFyQ29kZTogbnVtYmVyKTogc3RyaW5nIHtcbiAgY29uc3QgY2hhciA9IGNoYXJDb2RlID09PSBjaGFycy4kRU9GID8gJ0VPRicgOiBTdHJpbmcuZnJvbUNoYXJDb2RlKGNoYXJDb2RlKTtcbiAgcmV0dXJuIGBVbmV4cGVjdGVkIGNoYXJhY3RlciBcIiR7Y2hhcn1cImA7XG59XG5cbmZ1bmN0aW9uIF91bmtub3duRW50aXR5RXJyb3JNc2coZW50aXR5U3JjOiBzdHJpbmcpOiBzdHJpbmcge1xuICByZXR1cm4gYFVua25vd24gZW50aXR5IFwiJHtlbnRpdHlTcmN9XCIgLSB1c2UgdGhlIFwiJiM8ZGVjaW1hbD47XCIgb3IgIFwiJiN4PGhleD47XCIgc3ludGF4YDtcbn1cblxuZnVuY3Rpb24gX3VucGFyc2FibGVFbnRpdHlFcnJvck1zZyh0eXBlOiBDaGFyYWN0ZXJSZWZlcmVuY2VUeXBlLCBlbnRpdHlTdHI6IHN0cmluZyk6IHN0cmluZyB7XG4gIHJldHVybiBgVW5hYmxlIHRvIHBhcnNlIGVudGl0eSBcIiR7ZW50aXR5U3RyfVwiIC0gJHt0eXBlfSBjaGFyYWN0ZXIgcmVmZXJlbmNlIGVudGl0aWVzIG11c3QgZW5kIHdpdGggXCI7XCJgO1xufVxuXG5lbnVtIENoYXJhY3RlclJlZmVyZW5jZVR5cGUge1xuICBIRVggPSAnaGV4YWRlY2ltYWwnLFxuICBERUMgPSAnZGVjaW1hbCcsXG59XG5cbmNsYXNzIF9Db250cm9sRmxvd0Vycm9yIHtcbiAgY29uc3RydWN0b3IocHVibGljIGVycm9yOiBUb2tlbkVycm9yKSB7fVxufVxuXG4vLyBTZWUgaHR0cHM6Ly93d3cudzMub3JnL1RSL2h0bWw1MS9zeW50YXguaHRtbCN3cml0aW5nLWh0bWwtZG9jdW1lbnRzXG5jbGFzcyBfVG9rZW5pemVyIHtcbiAgcHJpdmF0ZSBfY3Vyc29yOiBDaGFyYWN0ZXJDdXJzb3I7XG4gIHByaXZhdGUgX3Rva2VuaXplSWN1OiBib29sZWFuO1xuICBwcml2YXRlIF9pbnRlcnBvbGF0aW9uQ29uZmlnOiBJbnRlcnBvbGF0aW9uQ29uZmlnO1xuICBwcml2YXRlIF9sZWFkaW5nVHJpdmlhQ29kZVBvaW50czogbnVtYmVyW10gfCB1bmRlZmluZWQ7XG4gIHByaXZhdGUgX2N1cnJlbnRUb2tlblN0YXJ0OiBDaGFyYWN0ZXJDdXJzb3IgfCBudWxsID0gbnVsbDtcbiAgcHJpdmF0ZSBfY3VycmVudFRva2VuVHlwZTogVG9rZW5UeXBlIHwgbnVsbCA9IG51bGw7XG4gIHByaXZhdGUgX2V4cGFuc2lvbkNhc2VTdGFjazogVG9rZW5UeXBlW10gPSBbXTtcbiAgcHJpdmF0ZSBfaW5JbnRlcnBvbGF0aW9uOiBib29sZWFuID0gZmFsc2U7XG4gIHByaXZhdGUgcmVhZG9ubHkgX3ByZXNlcnZlTGluZUVuZGluZ3M6IGJvb2xlYW47XG4gIHByaXZhdGUgcmVhZG9ubHkgX2kxOG5Ob3JtYWxpemVMaW5lRW5kaW5nc0luSUNVczogYm9vbGVhbjtcbiAgcHJpdmF0ZSByZWFkb25seSBfdG9rZW5pemVCbG9ja3M6IGJvb2xlYW47XG4gIHByaXZhdGUgcmVhZG9ubHkgX3Rva2VuaXplTGV0OiBib29sZWFuO1xuICB0b2tlbnM6IFRva2VuW10gPSBbXTtcbiAgZXJyb3JzOiBUb2tlbkVycm9yW10gPSBbXTtcbiAgbm9uTm9ybWFsaXplZEljdUV4cHJlc3Npb25zOiBUb2tlbltdID0gW107XG5cbiAgLyoqXG4gICAqIEBwYXJhbSBfZmlsZSBUaGUgaHRtbCBzb3VyY2UgZmlsZSBiZWluZyB0b2tlbml6ZWQuXG4gICAqIEBwYXJhbSBfZ2V0VGFnRGVmaW5pdGlvbiBBIGZ1bmN0aW9uIHRoYXQgd2lsbCByZXRyaWV2ZSBhIHRhZyBkZWZpbml0aW9uIGZvciBhIGdpdmVuIHRhZyBuYW1lLlxuICAgKiBAcGFyYW0gb3B0aW9ucyBDb25maWd1cmF0aW9uIG9mIHRoZSB0b2tlbml6YXRpb24uXG4gICAqL1xuICBjb25zdHJ1Y3RvcihcbiAgICBfZmlsZTogUGFyc2VTb3VyY2VGaWxlLFxuICAgIHByaXZhdGUgX2dldFRhZ0RlZmluaXRpb246ICh0YWdOYW1lOiBzdHJpbmcpID0+IFRhZ0RlZmluaXRpb24sXG4gICAgb3B0aW9uczogVG9rZW5pemVPcHRpb25zLFxuICApIHtcbiAgICB0aGlzLl90b2tlbml6ZUljdSA9IG9wdGlvbnMudG9rZW5pemVFeHBhbnNpb25Gb3JtcyB8fCBmYWxzZTtcbiAgICB0aGlzLl9pbnRlcnBvbGF0aW9uQ29uZmlnID0gb3B0aW9ucy5pbnRlcnBvbGF0aW9uQ29uZmlnIHx8IERFRkFVTFRfSU5URVJQT0xBVElPTl9DT05GSUc7XG4gICAgdGhpcy5fbGVhZGluZ1RyaXZpYUNvZGVQb2ludHMgPVxuICAgICAgb3B0aW9ucy5sZWFkaW5nVHJpdmlhQ2hhcnMgJiYgb3B0aW9ucy5sZWFkaW5nVHJpdmlhQ2hhcnMubWFwKChjKSA9PiBjLmNvZGVQb2ludEF0KDApIHx8IDApO1xuICAgIGNvbnN0IHJhbmdlID0gb3B0aW9ucy5yYW5nZSB8fCB7XG4gICAgICBlbmRQb3M6IF9maWxlLmNvbnRlbnQubGVuZ3RoLFxuICAgICAgc3RhcnRQb3M6IDAsXG4gICAgICBzdGFydExpbmU6IDAsXG4gICAgICBzdGFydENvbDogMCxcbiAgICB9O1xuICAgIHRoaXMuX2N1cnNvciA9IG9wdGlvbnMuZXNjYXBlZFN0cmluZ1xuICAgICAgPyBuZXcgRXNjYXBlZENoYXJhY3RlckN1cnNvcihfZmlsZSwgcmFuZ2UpXG4gICAgICA6IG5ldyBQbGFpbkNoYXJhY3RlckN1cnNvcihfZmlsZSwgcmFuZ2UpO1xuICAgIHRoaXMuX3ByZXNlcnZlTGluZUVuZGluZ3MgPSBvcHRpb25zLnByZXNlcnZlTGluZUVuZGluZ3MgfHwgZmFsc2U7XG4gICAgdGhpcy5faTE4bk5vcm1hbGl6ZUxpbmVFbmRpbmdzSW5JQ1VzID0gb3B0aW9ucy5pMThuTm9ybWFsaXplTGluZUVuZGluZ3NJbklDVXMgfHwgZmFsc2U7XG4gICAgdGhpcy5fdG9rZW5pemVCbG9ja3MgPSBvcHRpb25zLnRva2VuaXplQmxvY2tzID8/IHRydWU7XG4gICAgdGhpcy5fdG9rZW5pemVMZXQgPSBvcHRpb25zLnRva2VuaXplTGV0ID8/IHRydWU7XG4gICAgdHJ5IHtcbiAgICAgIHRoaXMuX2N1cnNvci5pbml0KCk7XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgdGhpcy5oYW5kbGVFcnJvcihlKTtcbiAgICB9XG4gIH1cblxuICBwcml2YXRlIF9wcm9jZXNzQ2FycmlhZ2VSZXR1cm5zKGNvbnRlbnQ6IHN0cmluZyk6IHN0cmluZyB7XG4gICAgaWYgKHRoaXMuX3ByZXNlcnZlTGluZUVuZGluZ3MpIHtcbiAgICAgIHJldHVybiBjb250ZW50O1xuICAgIH1cbiAgICAvLyBodHRwczovL3d3dy53My5vcmcvVFIvaHRtbDUxL3N5bnRheC5odG1sI3ByZXByb2Nlc3NpbmctdGhlLWlucHV0LXN0cmVhbVxuICAgIC8vIEluIG9yZGVyIHRvIGtlZXAgdGhlIG9yaWdpbmFsIHBvc2l0aW9uIGluIHRoZSBzb3VyY2UsIHdlIGNhbiBub3RcbiAgICAvLyBwcmUtcHJvY2VzcyBpdC5cbiAgICAvLyBJbnN0ZWFkIENScyBhcmUgcHJvY2Vzc2VkIHJpZ2h0IGJlZm9yZSBpbnN0YW50aWF0aW5nIHRoZSB0b2tlbnMuXG4gICAgcmV0dXJuIGNvbnRlbnQucmVwbGFjZShfQ1JfT1JfQ1JMRl9SRUdFWFAsICdcXG4nKTtcbiAgfVxuXG4gIHRva2VuaXplKCk6IHZvaWQge1xuICAgIHdoaWxlICh0aGlzLl9jdXJzb3IucGVlaygpICE9PSBjaGFycy4kRU9GKSB7XG4gICAgICBjb25zdCBzdGFydCA9IHRoaXMuX2N1cnNvci5jbG9uZSgpO1xuICAgICAgdHJ5IHtcbiAgICAgICAgaWYgKHRoaXMuX2F0dGVtcHRDaGFyQ29kZShjaGFycy4kTFQpKSB7XG4gICAgICAgICAgaWYgKHRoaXMuX2F0dGVtcHRDaGFyQ29kZShjaGFycy4kQkFORykpIHtcbiAgICAgICAgICAgIGlmICh0aGlzLl9hdHRlbXB0Q2hhckNvZGUoY2hhcnMuJExCUkFDS0VUKSkge1xuICAgICAgICAgICAgICB0aGlzLl9jb25zdW1lQ2RhdGEoc3RhcnQpO1xuICAgICAgICAgICAgfSBlbHNlIGlmICh0aGlzLl9hdHRlbXB0Q2hhckNvZGUoY2hhcnMuJE1JTlVTKSkge1xuICAgICAgICAgICAgICB0aGlzLl9jb25zdW1lQ29tbWVudChzdGFydCk7XG4gICAgICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgICB0aGlzLl9jb25zdW1lRG9jVHlwZShzdGFydCk7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgfSBlbHNlIGlmICh0aGlzLl9hdHRlbXB0Q2hhckNvZGUoY2hhcnMuJFNMQVNIKSkge1xuICAgICAgICAgICAgdGhpcy5fY29uc3VtZVRhZ0Nsb3NlKHN0YXJ0KTtcbiAgICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgdGhpcy5fY29uc3VtZVRhZ09wZW4oc3RhcnQpO1xuICAgICAgICAgIH1cbiAgICAgICAgfSBlbHNlIGlmIChcbiAgICAgICAgICB0aGlzLl90b2tlbml6ZUxldCAmJlxuICAgICAgICAgIC8vIFVzZSBgcGVla2AgaW5zdGVhZCBvZiBgYXR0ZW1wQ2hhckNvZGVgIHNpbmNlIHdlXG4gICAgICAgICAgLy8gZG9uJ3Qgd2FudCB0byBhZHZhbmNlIGluIGNhc2UgaXQncyBub3QgYEBsZXRgLlxuICAgICAgICAgIHRoaXMuX2N1cnNvci5wZWVrKCkgPT09IGNoYXJzLiRBVCAmJlxuICAgICAgICAgICF0aGlzLl9pbkludGVycG9sYXRpb24gJiZcbiAgICAgICAgICB0aGlzLl9hdHRlbXB0U3RyKCdAbGV0JylcbiAgICAgICAgKSB7XG4gICAgICAgICAgdGhpcy5fY29uc3VtZUxldERlY2xhcmF0aW9uKHN0YXJ0KTtcbiAgICAgICAgfSBlbHNlIGlmICh0aGlzLl90b2tlbml6ZUJsb2NrcyAmJiB0aGlzLl9hdHRlbXB0Q2hhckNvZGUoY2hhcnMuJEFUKSkge1xuICAgICAgICAgIHRoaXMuX2NvbnN1bWVCbG9ja1N0YXJ0KHN0YXJ0KTtcbiAgICAgICAgfSBlbHNlIGlmIChcbiAgICAgICAgICB0aGlzLl90b2tlbml6ZUJsb2NrcyAmJlxuICAgICAgICAgICF0aGlzLl9pbkludGVycG9sYXRpb24gJiZcbiAgICAgICAgICAhdGhpcy5faXNJbkV4cGFuc2lvbkNhc2UoKSAmJlxuICAgICAgICAgICF0aGlzLl9pc0luRXhwYW5zaW9uRm9ybSgpICYmXG4gICAgICAgICAgdGhpcy5fYXR0ZW1wdENoYXJDb2RlKGNoYXJzLiRSQlJBQ0UpXG4gICAgICAgICkge1xuICAgICAgICAgIHRoaXMuX2NvbnN1bWVCbG9ja0VuZChzdGFydCk7XG4gICAgICAgIH0gZWxzZSBpZiAoISh0aGlzLl90b2tlbml6ZUljdSAmJiB0aGlzLl90b2tlbml6ZUV4cGFuc2lvbkZvcm0oKSkpIHtcbiAgICAgICAgICAvLyBJbiAocG9zc2libHkgaW50ZXJwb2xhdGVkKSB0ZXh0IHRoZSBlbmQgb2YgdGhlIHRleHQgaXMgZ2l2ZW4gYnkgYGlzVGV4dEVuZCgpYCwgd2hpbGVcbiAgICAgICAgICAvLyB0aGUgcHJlbWF0dXJlIGVuZCBvZiBhbiBpbnRlcnBvbGF0aW9uIGlzIGdpdmVuIGJ5IHRoZSBzdGFydCBvZiBhIG5ldyBIVE1MIGVsZW1lbnQuXG4gICAgICAgICAgdGhpcy5fY29uc3VtZVdpdGhJbnRlcnBvbGF0aW9uKFxuICAgICAgICAgICAgVG9rZW5UeXBlLlRFWFQsXG4gICAgICAgICAgICBUb2tlblR5cGUuSU5URVJQT0xBVElPTixcbiAgICAgICAgICAgICgpID0+IHRoaXMuX2lzVGV4dEVuZCgpLFxuICAgICAgICAgICAgKCkgPT4gdGhpcy5faXNUYWdTdGFydCgpLFxuICAgICAgICAgICk7XG4gICAgICAgIH1cbiAgICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgdGhpcy5oYW5kbGVFcnJvcihlKTtcbiAgICAgIH1cbiAgICB9XG4gICAgdGhpcy5fYmVnaW5Ub2tlbihUb2tlblR5cGUuRU9GKTtcbiAgICB0aGlzLl9lbmRUb2tlbihbXSk7XG4gIH1cblxuICBwcml2YXRlIF9nZXRCbG9ja05hbWUoKTogc3RyaW5nIHtcbiAgICAvLyBUaGlzIGFsbG93cyB1cyB0byBjYXB0dXJlIHVwIHNvbWV0aGluZyBsaWtlIGBAZWxzZSBpZmAsIGJ1dCBub3QgYEAgaWZgLlxuICAgIGxldCBzcGFjZXNJbk5hbWVBbGxvd2VkID0gZmFsc2U7XG4gICAgY29uc3QgbmFtZUN1cnNvciA9IHRoaXMuX2N1cnNvci5jbG9uZSgpO1xuXG4gICAgdGhpcy5fYXR0ZW1wdENoYXJDb2RlVW50aWxGbigoY29kZSkgPT4ge1xuICAgICAgaWYgKGNoYXJzLmlzV2hpdGVzcGFjZShjb2RlKSkge1xuICAgICAgICByZXR1cm4gIXNwYWNlc0luTmFtZUFsbG93ZWQ7XG4gICAgICB9XG4gICAgICBpZiAoaXNCbG9ja05hbWVDaGFyKGNvZGUpKSB7XG4gICAgICAgIHNwYWNlc0luTmFtZUFsbG93ZWQgPSB0cnVlO1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgICB9XG4gICAgICByZXR1cm4gdHJ1ZTtcbiAgICB9KTtcbiAgICByZXR1cm4gdGhpcy5fY3Vyc29yLmdldENoYXJzKG5hbWVDdXJzb3IpLnRyaW0oKTtcbiAgfVxuXG4gIHByaXZhdGUgX2NvbnN1bWVCbG9ja1N0YXJ0KHN0YXJ0OiBDaGFyYWN0ZXJDdXJzb3IpIHtcbiAgICB0aGlzLl9iZWdpblRva2VuKFRva2VuVHlwZS5CTE9DS19PUEVOX1NUQVJULCBzdGFydCk7XG4gICAgY29uc3Qgc3RhcnRUb2tlbiA9IHRoaXMuX2VuZFRva2VuKFt0aGlzLl9nZXRCbG9ja05hbWUoKV0pO1xuXG4gICAgaWYgKHRoaXMuX2N1cnNvci5wZWVrKCkgPT09IGNoYXJzLiRMUEFSRU4pIHtcbiAgICAgIC8vIEFkdmFuY2UgcGFzdCB0aGUgb3BlbmluZyBwYXJlbi5cbiAgICAgIHRoaXMuX2N1cnNvci5hZHZhbmNlKCk7XG4gICAgICAvLyBDYXB0dXJlIHRoZSBwYXJhbWV0ZXJzLlxuICAgICAgdGhpcy5fY29uc3VtZUJsb2NrUGFyYW1ldGVycygpO1xuICAgICAgLy8gQWxsb3cgc3BhY2VzIGJlZm9yZSB0aGUgY2xvc2luZyBwYXJlbi5cbiAgICAgIHRoaXMuX2F0dGVtcHRDaGFyQ29kZVVudGlsRm4oaXNOb3RXaGl0ZXNwYWNlKTtcblxuICAgICAgaWYgKHRoaXMuX2F0dGVtcHRDaGFyQ29kZShjaGFycy4kUlBBUkVOKSkge1xuICAgICAgICAvLyBBbGxvdyBzcGFjZXMgYWZ0ZXIgdGhlIHBhcmVuLlxuICAgICAgICB0aGlzLl9hdHRlbXB0Q2hhckNvZGVVbnRpbEZuKGlzTm90V2hpdGVzcGFjZSk7XG4gICAgICB9IGVsc2Uge1xuICAgICAgICBzdGFydFRva2VuLnR5cGUgPSBUb2tlblR5cGUuSU5DT01QTEVURV9CTE9DS19PUEVOO1xuICAgICAgICByZXR1cm47XG4gICAgICB9XG4gICAgfVxuXG4gICAgaWYgKHRoaXMuX2F0dGVtcHRDaGFyQ29kZShjaGFycy4kTEJSQUNFKSkge1xuICAgICAgdGhpcy5fYmVnaW5Ub2tlbihUb2tlblR5cGUuQkxPQ0tfT1BFTl9FTkQpO1xuICAgICAgdGhpcy5fZW5kVG9rZW4oW10pO1xuICAgIH0gZWxzZSB7XG4gICAgICBzdGFydFRva2VuLnR5cGUgPSBUb2tlblR5cGUuSU5DT01QTEVURV9CTE9DS19PUEVOO1xuICAgIH1cbiAgfVxuXG4gIHByaXZhdGUgX2NvbnN1bWVCbG9ja0VuZChzdGFydDogQ2hhcmFjdGVyQ3Vyc29yKSB7XG4gICAgdGhpcy5fYmVnaW5Ub2tlbihUb2tlblR5cGUuQkxPQ0tfQ0xPU0UsIHN0YXJ0KTtcbiAgICB0aGlzLl9lbmRUb2tlbihbXSk7XG4gIH1cblxuICBwcml2YXRlIF9jb25zdW1lQmxvY2tQYXJhbWV0ZXJzKCkge1xuICAgIC8vIFRyaW0gdGhlIHdoaXRlc3BhY2UgdW50aWwgdGhlIGZpcnN0IHBhcmFtZXRlci5cbiAgICB0aGlzLl9hdHRlbXB0Q2hhckNvZGVVbnRpbEZuKGlzQmxvY2tQYXJhbWV0ZXJDaGFyKTtcblxuICAgIHdoaWxlICh0aGlzLl9jdXJzb3IucGVlaygpICE9PSBjaGFycy4kUlBBUkVOICYmIHRoaXMuX2N1cnNvci5wZWVrKCkgIT09IGNoYXJzLiRFT0YpIHtcbiAgICAgIHRoaXMuX2JlZ2luVG9rZW4oVG9rZW5UeXBlLkJMT0NLX1BBUkFNRVRFUik7XG4gICAgICBjb25zdCBzdGFydCA9IHRoaXMuX2N1cnNvci5jbG9uZSgpO1xuICAgICAgbGV0IGluUXVvdGU6IG51bWJlciB8IG51bGwgPSBudWxsO1xuICAgICAgbGV0IG9wZW5QYXJlbnMgPSAwO1xuXG4gICAgICAvLyBDb25zdW1lIHRoZSBwYXJhbWV0ZXIgdW50aWwgdGhlIG5leHQgc2VtaWNvbG9uIG9yIGJyYWNlLlxuICAgICAgLy8gTm90ZSB0aGF0IHdlIHNraXAgb3ZlciBzZW1pY29sb25zL2JyYWNlcyBpbnNpZGUgb2Ygc3RyaW5ncy5cbiAgICAgIHdoaWxlIChcbiAgICAgICAgKHRoaXMuX2N1cnNvci5wZWVrKCkgIT09IGNoYXJzLiRTRU1JQ09MT04gJiYgdGhpcy5fY3Vyc29yLnBlZWsoKSAhPT0gY2hhcnMuJEVPRikgfHxcbiAgICAgICAgaW5RdW90ZSAhPT0gbnVsbFxuICAgICAgKSB7XG4gICAgICAgIGNvbnN0IGNoYXIgPSB0aGlzLl9jdXJzb3IucGVlaygpO1xuXG4gICAgICAgIC8vIFNraXAgdG8gdGhlIG5leHQgY2hhcmFjdGVyIGlmIGl0IHdhcyBlc2NhcGVkLlxuICAgICAgICBpZiAoY2hhciA9PT0gY2hhcnMuJEJBQ0tTTEFTSCkge1xuICAgICAgICAgIHRoaXMuX2N1cnNvci5hZHZhbmNlKCk7XG4gICAgICAgIH0gZWxzZSBpZiAoY2hhciA9PT0gaW5RdW90ZSkge1xuICAgICAgICAgIGluUXVvdGUgPSBudWxsO1xuICAgICAgICB9IGVsc2UgaWYgKGluUXVvdGUgPT09IG51bGwgJiYgY2hhcnMuaXNRdW90ZShjaGFyKSkge1xuICAgICAgICAgIGluUXVvdGUgPSBjaGFyO1xuICAgICAgICB9IGVsc2UgaWYgKGNoYXIgPT09IGNoYXJzLiRMUEFSRU4gJiYgaW5RdW90ZSA9PT0gbnVsbCkge1xuICAgICAgICAgIG9wZW5QYXJlbnMrKztcbiAgICAgICAgfSBlbHNlIGlmIChjaGFyID09PSBjaGFycy4kUlBBUkVOICYmIGluUXVvdGUgPT09IG51bGwpIHtcbiAgICAgICAgICBpZiAob3BlblBhcmVucyA9PT0gMCkge1xuICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgICAgfSBlbHNlIGlmIChvcGVuUGFyZW5zID4gMCkge1xuICAgICAgICAgICAgb3BlblBhcmVucy0tO1xuICAgICAgICAgIH1cbiAgICAgICAgfVxuXG4gICAgICAgIHRoaXMuX2N1cnNvci5hZHZhbmNlKCk7XG4gICAgICB9XG5cbiAgICAgIHRoaXMuX2VuZFRva2VuKFt0aGlzLl9jdXJzb3IuZ2V0Q2hhcnMoc3RhcnQpXSk7XG5cbiAgICAgIC8vIFNraXAgdG8gdGhlIG5leHQgcGFyYW1ldGVyLlxuICAgICAgdGhpcy5fYXR0ZW1wdENoYXJDb2RlVW50aWxGbihpc0Jsb2NrUGFyYW1ldGVyQ2hhcik7XG4gICAgfVxuICB9XG5cbiAgcHJpdmF0ZSBfY29uc3VtZUxldERlY2xhcmF0aW9uKHN0YXJ0OiBDaGFyYWN0ZXJDdXJzb3IpIHtcbiAgICB0aGlzLl9iZWdpblRva2VuKFRva2VuVHlwZS5MRVRfU1RBUlQsIHN0YXJ0KTtcblxuICAgIC8vIFJlcXVpcmUgYXQgbGVhc3Qgb25lIHdoaXRlIHNwYWNlIGFmdGVyIHRoZSBgQGxldGAuXG4gICAgaWYgKGNoYXJzLmlzV2hpdGVzcGFjZSh0aGlzLl9jdXJzb3IucGVlaygpKSkge1xuICAgICAgdGhpcy5fYXR0ZW1wdENoYXJDb2RlVW50aWxGbihpc05vdFdoaXRlc3BhY2UpO1xuICAgIH0gZWxzZSB7XG4gICAgICBjb25zdCB0b2tlbiA9IHRoaXMuX2VuZFRva2VuKFt0aGlzLl9jdXJzb3IuZ2V0Q2hhcnMoc3RhcnQpXSk7XG4gICAgICB0b2tlbi50eXBlID0gVG9rZW5UeXBlLklOQ09NUExFVEVfTEVUO1xuICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGNvbnN0IHN0YXJ0VG9rZW4gPSB0aGlzLl9lbmRUb2tlbihbdGhpcy5fZ2V0TGV0RGVjbGFyYXRpb25OYW1lKCldKTtcblxuICAgIC8vIFNraXAgb3ZlciB3aGl0ZSBzcGFjZSBiZWZvcmUgdGhlIGVxdWFscyBjaGFyYWN0ZXIuXG4gICAgdGhpcy5fYXR0ZW1wdENoYXJDb2RlVW50aWxGbihpc05vdFdoaXRlc3BhY2UpO1xuXG4gICAgLy8gRXhwZWN0IGFuIGVxdWFscyBzaWduLlxuICAgIGlmICghdGhpcy5fYXR0ZW1wdENoYXJDb2RlKGNoYXJzLiRFUSkpIHtcbiAgICAgIHN0YXJ0VG9rZW4udHlwZSA9IFRva2VuVHlwZS5JTkNPTVBMRVRFX0xFVDtcbiAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBTa2lwIHNwYWNlcyBhZnRlciB0aGUgZXF1YWxzLlxuICAgIHRoaXMuX2F0dGVtcHRDaGFyQ29kZVVudGlsRm4oKGNvZGUpID0+IGlzTm90V2hpdGVzcGFjZShjb2RlKSAmJiAhY2hhcnMuaXNOZXdMaW5lKGNvZGUpKTtcbiAgICB0aGlzLl9jb25zdW1lTGV0RGVjbGFyYXRpb25WYWx1ZSgpO1xuXG4gICAgLy8gVGVybWluYXRlIHRoZSBgQGxldGAgd2l0aCBhIHNlbWljb2xvbi5cbiAgICBjb25zdCBlbmRDaGFyID0gdGhpcy5fY3Vyc29yLnBlZWsoKTtcbiAgICBpZiAoZW5kQ2hhciA9PT0gY2hhcnMuJFNFTUlDT0xPTikge1xuICAgICAgdGhpcy5fYmVnaW5Ub2tlbihUb2tlblR5cGUuTEVUX0VORCk7XG4gICAgICB0aGlzLl9lbmRUb2tlbihbXSk7XG4gICAgICB0aGlzLl9jdXJzb3IuYWR2YW5jZSgpO1xuICAgIH0gZWxzZSB7XG4gICAgICBzdGFydFRva2VuLnR5cGUgPSBUb2tlblR5cGUuSU5DT01QTEVURV9MRVQ7XG4gICAgICBzdGFydFRva2VuLnNvdXJjZVNwYW4gPSB0aGlzLl9jdXJzb3IuZ2V0U3BhbihzdGFydCk7XG4gICAgfVxuICB9XG5cbiAgcHJpdmF0ZSBfZ2V0TGV0RGVjbGFyYXRpb25OYW1lKCk6IHN0cmluZyB7XG4gICAgY29uc3QgbmFtZUN1cnNvciA9IHRoaXMuX2N1cnNvci5jbG9uZSgpO1xuICAgIGxldCBhbGxvd0RpZ2l0ID0gZmFsc2U7XG5cbiAgICB0aGlzLl9hdHRlbXB0Q2hhckNvZGVVbnRpbEZuKChjb2RlKSA9PiB7XG4gICAgICBpZiAoXG4gICAgICAgIGNoYXJzLmlzQXNjaWlMZXR0ZXIoY29kZSkgfHxcbiAgICAgICAgY29kZSA9PSBjaGFycy4kJCB8fFxuICAgICAgICBjb2RlID09PSBjaGFycy4kXyB8fFxuICAgICAgICAvLyBgQGxldGAgbmFtZXMgY2FuJ3Qgc3RhcnQgd2l0aCBhIGRpZ2l0LCBidXQgZGlnaXRzIGFyZSB2YWxpZCBhbnl3aGVyZSBlbHNlIGluIHRoZSBuYW1lLlxuICAgICAgICAoYWxsb3dEaWdpdCAmJiBjaGFycy5pc0RpZ2l0KGNvZGUpKVxuICAgICAgKSB7XG4gICAgICAgIGFsbG93RGlnaXQgPSB0cnVlO1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgICB9XG4gICAgICByZXR1cm4gdHJ1ZTtcbiAgICB9KTtcblxuICAgIHJldHVybiB0aGlzLl9jdXJzb3IuZ2V0Q2hhcnMobmFtZUN1cnNvcikudHJpbSgpO1xuICB9XG5cbiAgcHJpdmF0ZSBfY29uc3VtZUxldERlY2xhcmF0aW9uVmFsdWUoKTogdm9pZCB7XG4gICAgY29uc3Qgc3RhcnQgPSB0aGlzLl9jdXJzb3IuY2xvbmUoKTtcbiAgICB0aGlzLl9iZWdpblRva2VuKFRva2VuVHlwZS5MRVRfVkFMVUUsIHN0YXJ0KTtcblxuICAgIHdoaWxlICh0aGlzLl9jdXJzb3IucGVlaygpICE9PSBjaGFycy4kRU9GKSB7XG4gICAgICBjb25zdCBjaGFyID0gdGhpcy5fY3Vyc29yLnBlZWsoKTtcblxuICAgICAgLy8gYEBsZXRgIGRlY2xhcmF0aW9ucyB0ZXJtaW5hdGUgd2l0aCBhIHNlbWljb2xvbi5cbiAgICAgIGlmIChjaGFyID09PSBjaGFycy4kU0VNSUNPTE9OKSB7XG4gICAgICAgIGJyZWFrO1xuICAgICAgfVxuXG4gICAgICAvLyBJZiB3ZSBoaXQgYSBxdW90ZSwgc2tpcCBvdmVyIGl0cyBjb250ZW50IHNpbmNlIHdlIGRvbid0IGNhcmUgd2hhdCdzIGluc2lkZS5cbiAgICAgIGlmIChjaGFycy5pc1F1b3RlKGNoYXIpKSB7XG4gICAgICAgIHRoaXMuX2N1cnNvci5hZHZhbmNlKCk7XG4gICAgICAgIHRoaXMuX2F0dGVtcHRDaGFyQ29kZVVudGlsRm4oKGlubmVyKSA9PiB7XG4gICAgICAgICAgaWYgKGlubmVyID09PSBjaGFycy4kQkFDS1NMQVNIKSB7XG4gICAgICAgICAgICB0aGlzLl9jdXJzb3IuYWR2YW5jZSgpO1xuICAgICAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgICAgICAgIH1cbiAgICAgICAgICByZXR1cm4gaW5uZXIgPT09IGNoYXI7XG4gICAgICAgIH0pO1xuICAgICAgfVxuXG4gICAgICB0aGlzLl9jdXJzb3IuYWR2YW5jZSgpO1xuICAgIH1cblxuICAgIHRoaXMuX2VuZFRva2VuKFt0aGlzLl9jdXJzb3IuZ2V0Q2hhcnMoc3RhcnQpXSk7XG4gIH1cblxuICAvKipcbiAgICogQHJldHVybnMgd2hldGhlciBhbiBJQ1UgdG9rZW4gaGFzIGJlZW4gY3JlYXRlZFxuICAgKiBAaW50ZXJuYWxcbiAgICovXG4gIHByaXZhdGUgX3Rva2VuaXplRXhwYW5zaW9uRm9ybSgpOiBib29sZWFuIHtcbiAgICBpZiAodGhpcy5pc0V4cGFuc2lvbkZvcm1TdGFydCgpKSB7XG4gICAgICB0aGlzLl9jb25zdW1lRXhwYW5zaW9uRm9ybVN0YXJ0KCk7XG4gICAgICByZXR1cm4gdHJ1ZTtcbiAgICB9XG5cbiAgICBpZiAoaXNFeHBhbnNpb25DYXNlU3RhcnQodGhpcy5fY3Vyc29yLnBlZWsoKSkgJiYgdGhpcy5faXNJbkV4cGFuc2lvbkZvcm0oKSkge1xuICAgICAgdGhpcy5fY29uc3VtZUV4cGFuc2lvbkNhc2VTdGFydCgpO1xuICAgICAgcmV0dXJuIHRydWU7XG4gICAgfVxuXG4gICAgaWYgKHRoaXMuX2N1cnNvci5wZWVrKCkgPT09IGNoYXJzLiRSQlJBQ0UpIHtcbiAgICAgIGlmICh0aGlzLl9pc0luRXhwYW5zaW9uQ2FzZSgpKSB7XG4gICAgICAgIHRoaXMuX2NvbnN1bWVFeHBhbnNpb25DYXNlRW5kKCk7XG4gICAgICAgIHJldHVybiB0cnVlO1xuICAgICAgfVxuXG4gICAgICBpZiAodGhpcy5faXNJbkV4cGFuc2lvbkZvcm0oKSkge1xuICAgICAgICB0aGlzLl9jb25zdW1lRXhwYW5zaW9uRm9ybUVuZCgpO1xuICAgICAgICByZXR1cm4gdHJ1ZTtcbiAgICAgIH1cbiAgICB9XG5cbiAgICByZXR1cm4gZmFsc2U7XG4gIH1cblxuICBwcml2YXRlIF9iZWdpblRva2VuKHR5cGU6IFRva2VuVHlwZSwgc3RhcnQgPSB0aGlzLl9jdXJzb3IuY2xvbmUoKSkge1xuICAgIHRoaXMuX2N1cnJlbnRUb2tlblN0YXJ0ID0gc3RhcnQ7XG4gICAgdGhpcy5fY3VycmVudFRva2VuVHlwZSA9IHR5cGU7XG4gIH1cblxuICBwcml2YXRlIF9lbmRUb2tlbihwYXJ0czogc3RyaW5nW10sIGVuZD86IENoYXJhY3RlckN1cnNvcik6IFRva2VuIHtcbiAgICBpZiAodGhpcy5fY3VycmVudFRva2VuU3RhcnQgPT09IG51bGwpIHtcbiAgICAgIHRocm93IG5ldyBUb2tlbkVycm9yKFxuICAgICAgICAnUHJvZ3JhbW1pbmcgZXJyb3IgLSBhdHRlbXB0ZWQgdG8gZW5kIGEgdG9rZW4gd2hlbiB0aGVyZSB3YXMgbm8gc3RhcnQgdG8gdGhlIHRva2VuJyxcbiAgICAgICAgdGhpcy5fY3VycmVudFRva2VuVHlwZSxcbiAgICAgICAgdGhpcy5fY3Vyc29yLmdldFNwYW4oZW5kKSxcbiAgICAgICk7XG4gICAgfVxuICAgIGlmICh0aGlzLl9jdXJyZW50VG9rZW5UeXBlID09PSBudWxsKSB7XG4gICAgICB0aHJvdyBuZXcgVG9rZW5FcnJvcihcbiAgICAgICAgJ1Byb2dyYW1taW5nIGVycm9yIC0gYXR0ZW1wdGVkIHRvIGVuZCBhIHRva2VuIHdoaWNoIGhhcyBubyB0b2tlbiB0eXBlJyxcbiAgICAgICAgbnVsbCxcbiAgICAgICAgdGhpcy5fY3Vyc29yLmdldFNwYW4odGhpcy5fY3VycmVudFRva2VuU3RhcnQpLFxuICAgICAgKTtcbiAgICB9XG4gICAgY29uc3QgdG9rZW4gPSB7XG4gICAgICB0eXBlOiB0aGlzLl9jdXJyZW50VG9rZW5UeXBlLFxuICAgICAgcGFydHMsXG4gICAgICBzb3VyY2VTcGFuOiAoZW5kID8/IHRoaXMuX2N1cnNvcikuZ2V0U3BhbihcbiAgICAgICAgdGhpcy5fY3VycmVudFRva2VuU3RhcnQsXG4gICAgICAgIHRoaXMuX2xlYWRpbmdUcml2aWFDb2RlUG9pbnRzLFxuICAgICAgKSxcbiAgICB9IGFzIFRva2VuO1xuICAgIHRoaXMudG9rZW5zLnB1c2godG9rZW4pO1xuICAgIHRoaXMuX2N1cnJlbnRUb2tlblN0YXJ0ID0gbnVsbDtcbiAgICB0aGlzLl9jdXJyZW50VG9rZW5UeXBlID0gbnVsbDtcbiAgICByZXR1cm4gdG9rZW47XG4gIH1cblxuICBwcml2YXRlIF9jcmVhdGVFcnJvcihtc2c6IHN0cmluZywgc3BhbjogUGFyc2VTb3VyY2VTcGFuKTogX0NvbnRyb2xGbG93RXJyb3Ige1xuICAgIGlmICh0aGlzLl9pc0luRXhwYW5zaW9uRm9ybSgpKSB7XG4gICAgICBtc2cgKz0gYCAoRG8geW91IGhhdmUgYW4gdW5lc2NhcGVkIFwie1wiIGluIHlvdXIgdGVtcGxhdGU/IFVzZSBcInt7ICd7JyB9fVwiKSB0byBlc2NhcGUgaXQuKWA7XG4gICAgfVxuICAgIGNvbnN0IGVycm9yID0gbmV3IFRva2VuRXJyb3IobXNnLCB0aGlzLl9jdXJyZW50VG9rZW5UeXBlLCBzcGFuKTtcbiAgICB0aGlzLl9jdXJyZW50VG9rZW5TdGFydCA9IG51bGw7XG4gICAgdGhpcy5fY3VycmVudFRva2VuVHlwZSA9IG51bGw7XG4gICAgcmV0dXJuIG5ldyBfQ29udHJvbEZsb3dFcnJvcihlcnJvcik7XG4gIH1cblxuICBwcml2YXRlIGhhbmRsZUVycm9yKGU6IGFueSkge1xuICAgIGlmIChlIGluc3RhbmNlb2YgQ3Vyc29yRXJyb3IpIHtcbiAgICAgIGUgPSB0aGlzLl9jcmVhdGVFcnJvcihlLm1zZywgdGhpcy5fY3Vyc29yLmdldFNwYW4oZS5jdXJzb3IpKTtcbiAgICB9XG4gICAgaWYgKGUgaW5zdGFuY2VvZiBfQ29udHJvbEZsb3dFcnJvcikge1xuICAgICAgdGhpcy5lcnJvcnMucHVzaChlLmVycm9yKTtcbiAgICB9IGVsc2Uge1xuICAgICAgdGhyb3cgZTtcbiAgICB9XG4gIH1cblxuICBwcml2YXRlIF9hdHRlbXB0Q2hhckNvZGUoY2hhckNvZGU6IG51bWJlcik6IGJvb2xlYW4ge1xuICAgIGlmICh0aGlzLl9jdXJzb3IucGVlaygpID09PSBjaGFyQ29kZSkge1xuICAgICAgdGhpcy5fY3Vyc29yLmFkdmFuY2UoKTtcbiAgICAgIHJldHVybiB0cnVlO1xuICAgIH1cbiAgICByZXR1cm4gZmFsc2U7XG4gIH1cblxuICBwcml2YXRlIF9hdHRlbXB0Q2hhckNvZGVDYXNlSW5zZW5zaXRpdmUoY2hhckNvZGU6IG51bWJlcik6IGJvb2xlYW4ge1xuICAgIGlmIChjb21wYXJlQ2hhckNvZGVDYXNlSW5zZW5zaXRpdmUodGhpcy5fY3Vyc29yLnBlZWsoKSwgY2hhckNvZGUpKSB7XG4gICAgICB0aGlzLl9jdXJzb3IuYWR2YW5jZSgpO1xuICAgICAgcmV0dXJuIHRydWU7XG4gICAgfVxuICAgIHJldHVybiBmYWxzZTtcbiAgfVxuXG4gIHByaXZhdGUgX3JlcXVpcmVDaGFyQ29kZShjaGFyQ29kZTogbnVtYmVyKSB7XG4gICAgY29uc3QgbG9jYXRpb24gPSB0aGlzLl9jdXJzb3IuY2xvbmUoKTtcbiAgICBpZiAoIXRoaXMuX2F0dGVtcHRDaGFyQ29kZShjaGFyQ29kZSkpIHtcbiAgICAgIHRocm93IHRoaXMuX2NyZWF0ZUVycm9yKFxuICAgICAgICBfdW5leHBlY3RlZENoYXJhY3RlckVycm9yTXNnKHRoaXMuX2N1cnNvci5wZWVrKCkpLFxuICAgICAgICB0aGlzLl9jdXJzb3IuZ2V0U3Bhbihsb2NhdGlvbiksXG4gICAgICApO1xuICAgIH1cbiAgfVxuXG4gIHByaXZhdGUgX2F0dGVtcHRTdHIoY2hhcnM6IHN0cmluZyk6IGJvb2xlYW4ge1xuICAgIGNvbnN0IGxlbiA9IGNoYXJzLmxlbmd0aDtcbiAgICBpZiAodGhpcy5fY3Vyc29yLmNoYXJzTGVmdCgpIDwgbGVuKSB7XG4gICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuICAgIGNvbnN0IGluaXRpYWxQb3NpdGlvbiA9IHRoaXMuX2N1cnNvci5jbG9uZSgpO1xuICAgIGZvciAobGV0IGkgPSAwOyBpIDwgbGVuOyBpKyspIHtcbiAgICAgIGlmICghdGhpcy5fYXR0ZW1wdENoYXJDb2RlKGNoYXJzLmNoYXJDb2RlQXQoaSkpKSB7XG4gICAgICAgIC8vIElmIGF0dGVtcHRpbmcgdG8gcGFyc2UgdGhlIHN0cmluZyBmYWlscywgd2Ugd2FudCB0byByZXNldCB0aGUgcGFyc2VyXG4gICAgICAgIC8vIHRvIHdoZXJlIGl0IHdhcyBiZWZvcmUgdGhlIGF0dGVtcHRcbiAgICAgICAgdGhpcy5fY3Vyc29yID0gaW5pdGlhbFBvc2l0aW9uO1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgICB9XG4gICAgfVxuICAgIHJldHVybiB0cnVlO1xuICB9XG5cbiAgcHJpdmF0ZSBfYXR0ZW1wdFN0ckNhc2VJbnNlbnNpdGl2ZShjaGFyczogc3RyaW5nKTogYm9vbGVhbiB7XG4gICAgZm9yIChsZXQgaSA9IDA7IGkgPCBjaGFycy5sZW5ndGg7IGkrKykge1xuICAgICAgaWYgKCF0aGlzLl9hdHRlbXB0Q2hhckNvZGVDYXNlSW5zZW5zaXRpdmUoY2hhcnMuY2hhckNvZGVBdChpKSkpIHtcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgICAgfVxuICAgIH1cbiAgICByZXR1cm4gdHJ1ZTtcbiAgfVxuXG4gIHByaXZhdGUgX3JlcXVpcmVTdHIoY2hhcnM6IHN0cmluZykge1xuICAgIGNvbnN0IGxvY2F0aW9uID0gdGhpcy5fY3Vyc29yLmNsb25lKCk7XG4gICAgaWYgKCF0aGlzLl9hdHRlbXB0U3RyKGNoYXJzKSkge1xuICAgICAgdGhyb3cgdGhpcy5fY3JlYXRlRXJyb3IoXG4gICAgICAgIF91bmV4cGVjdGVkQ2hhcmFjdGVyRXJyb3JNc2codGhpcy5fY3Vyc29yLnBlZWsoKSksXG4gICAgICAgIHRoaXMuX2N1cnNvci5nZXRTcGFuKGxvY2F0aW9uKSxcbiAgICAgICk7XG4gICAgfVxuICB9XG5cbiAgcHJpdmF0ZSBfYXR0ZW1wdENoYXJDb2RlVW50aWxGbihwcmVkaWNhdGU6IChjb2RlOiBudW1iZXIpID0+IGJvb2xlYW4pIHtcbiAgICB3aGlsZSAoIXByZWRpY2F0ZSh0aGlzLl9jdXJzb3IucGVlaygpKSkge1xuICAgICAgdGhpcy5fY3Vyc29yLmFkdmFuY2UoKTtcbiAgICB9XG4gIH1cblxuICBwcml2YXRlIF9yZXF1aXJlQ2hhckNvZGVVbnRpbEZuKHByZWRpY2F0ZTogKGNvZGU6IG51bWJlcikgPT4gYm9vbGVhbiwgbGVuOiBudW1iZXIpIHtcbiAgICBjb25zdCBzdGFydCA9IHRoaXMuX2N1cnNvci5jbG9uZSgpO1xuICAgIHRoaXMuX2F0dGVtcHRDaGFyQ29kZVVudGlsRm4ocHJlZGljYXRlKTtcbiAgICBpZiAodGhpcy5fY3Vyc29yLmRpZmYoc3RhcnQpIDwgbGVuKSB7XG4gICAgICB0aHJvdyB0aGlzLl9jcmVhdGVFcnJvcihcbiAgICAgICAgX3VuZXhwZWN0ZWRDaGFyYWN0ZXJFcnJvck1zZyh0aGlzLl9jdXJzb3IucGVlaygpKSxcbiAgICAgICAgdGhpcy5fY3Vyc29yLmdldFNwYW4oc3RhcnQpLFxuICAgICAgKTtcbiAgICB9XG4gIH1cblxuICBwcml2YXRlIF9hdHRlbXB0VW50aWxDaGFyKGNoYXI6IG51bWJlcikge1xuICAgIHdoaWxlICh0aGlzLl9jdXJzb3IucGVlaygpICE9PSBjaGFyKSB7XG4gICAgICB0aGlzLl9jdXJzb3IuYWR2YW5jZSgpO1xuICAgIH1cbiAgfVxuXG4gIHByaXZhdGUgX3JlYWRDaGFyKCk6IHN0cmluZyB7XG4gICAgLy8gRG9uJ3QgcmVseSB1cG9uIHJlYWRpbmcgZGlyZWN0bHkgZnJvbSBgX2lucHV0YCBhcyB0aGUgYWN0dWFsIGNoYXIgdmFsdWVcbiAgICAvLyBtYXkgaGF2ZSBiZWVuIGdlbmVyYXRlZCBmcm9tIGFuIGVzY2FwZSBzZXF1ZW5jZS5cbiAgICBjb25zdCBjaGFyID0gU3RyaW5nLmZyb21Db2RlUG9pbnQodGhpcy5fY3Vyc29yLnBlZWsoKSk7XG4gICAgdGhpcy5fY3Vyc29yLmFkdmFuY2UoKTtcbiAgICByZXR1cm4gY2hhcjtcbiAgfVxuXG4gIHByaXZhdGUgX2NvbnN1bWVFbnRpdHkodGV4dFRva2VuVHlwZTogVG9rZW5UeXBlKTogdm9pZCB7XG4gICAgdGhpcy5fYmVnaW5Ub2tlbihUb2tlblR5cGUuRU5DT0RFRF9FTlRJVFkpO1xuICAgIGNvbnN0IHN0YXJ0ID0gdGhpcy5fY3Vyc29yLmNsb25lKCk7XG4gICAgdGhpcy5fY3Vyc29yLmFkdmFuY2UoKTtcbiAgICBpZiAodGhpcy5fYXR0ZW1wdENoYXJDb2RlKGNoYXJzLiRIQVNIKSkge1xuICAgICAgY29uc3QgaXNIZXggPSB0aGlzLl9hdHRlbXB0Q2hhckNvZGUoY2hhcnMuJHgpIHx8IHRoaXMuX2F0dGVtcHRDaGFyQ29kZShjaGFycy4kWCk7XG4gICAgICBjb25zdCBjb2RlU3RhcnQgPSB0aGlzLl9jdXJzb3IuY2xvbmUoKTtcbiAgICAgIHRoaXMuX2F0dGVtcHRDaGFyQ29kZVVudGlsRm4oaXNEaWdpdEVudGl0eUVuZCk7XG4gICAgICBpZiAodGhpcy5fY3Vyc29yLnBlZWsoKSAhPSBjaGFycy4kU0VNSUNPTE9OKSB7XG4gICAgICAgIC8vIEFkdmFuY2UgY3Vyc29yIHRvIGluY2x1ZGUgdGhlIHBlZWtlZCBjaGFyYWN0ZXIgaW4gdGhlIHN0cmluZyBwcm92aWRlZCB0byB0aGUgZXJyb3JcbiAgICAgICAgLy8gbWVzc2FnZS5cbiAgICAgICAgdGhpcy5fY3Vyc29yLmFkdmFuY2UoKTtcbiAgICAgICAgY29uc3QgZW50aXR5VHlwZSA9IGlzSGV4ID8gQ2hhcmFjdGVyUmVmZXJlbmNlVHlwZS5IRVggOiBDaGFyYWN0ZXJSZWZlcmVuY2VUeXBlLkRFQztcbiAgICAgICAgdGhyb3cgdGhpcy5fY3JlYXRlRXJyb3IoXG4gICAgICAgICAgX3VucGFyc2FibGVFbnRpdHlFcnJvck1zZyhlbnRpdHlUeXBlLCB0aGlzLl9jdXJzb3IuZ2V0Q2hhcnMoc3RhcnQpKSxcbiAgICAgICAgICB0aGlzLl9jdXJzb3IuZ2V0U3BhbigpLFxuICAgICAgICApO1xuICAgICAgfVxuICAgICAgY29uc3Qgc3RyTnVtID0gdGhpcy5fY3Vyc29yLmdldENoYXJzKGNvZGVTdGFydCk7XG4gICAgICB0aGlzLl9jdXJzb3IuYWR2YW5jZSgpO1xuICAgICAgdHJ5IHtcbiAgICAgICAgY29uc3QgY2hhckNvZGUgPSBwYXJzZUludChzdHJOdW0sIGlzSGV4ID8gMTYgOiAxMCk7XG4gICAgICAgIHRoaXMuX2VuZFRva2VuKFtTdHJpbmcuZnJvbUNoYXJDb2RlKGNoYXJDb2RlKSwgdGhpcy5fY3Vyc29yLmdldENoYXJzKHN0YXJ0KV0pO1xuICAgICAgfSBjYXRjaCB7XG4gICAgICAgIHRocm93IHRoaXMuX2NyZWF0ZUVycm9yKFxuICAgICAgICAgIF91bmtub3duRW50aXR5RXJyb3JNc2codGhpcy5fY3Vyc29yLmdldENoYXJzKHN0YXJ0KSksXG4gICAgICAgICAgdGhpcy5fY3Vyc29yLmdldFNwYW4oKSxcbiAgICAgICAgKTtcbiAgICAgIH1cbiAgICB9IGVsc2Uge1xuICAgICAgY29uc3QgbmFtZVN0YXJ0ID0gdGhpcy5fY3Vyc29yLmNsb25lKCk7XG4gICAgICB0aGlzLl9hdHRlbXB0Q2hhckNvZGVVbnRpbEZuKGlzTmFtZWRFbnRpdHlFbmQpO1xuICAgICAgaWYgKHRoaXMuX2N1cnNvci5wZWVrKCkgIT0gY2hhcnMuJFNFTUlDT0xPTikge1xuICAgICAgICAvLyBObyBzZW1pY29sb24gd2FzIGZvdW5kIHNvIGFib3J0IHRoZSBlbmNvZGVkIGVudGl0eSB0b2tlbiB0aGF0IHdhcyBpbiBwcm9ncmVzcywgYW5kIHRyZWF0XG4gICAgICAgIC8vIHRoaXMgYXMgYSB0ZXh0IHRva2VuXG4gICAgICAgIHRoaXMuX2JlZ2luVG9rZW4odGV4dFRva2VuVHlwZSwgc3RhcnQpO1xuICAgICAgICB0aGlzLl9jdXJzb3IgPSBuYW1lU3RhcnQ7XG4gICAgICAgIHRoaXMuX2VuZFRva2VuKFsnJiddKTtcbiAgICAgIH0gZWxzZSB7XG4gICAgICAgIGNvbnN0IG5hbWUgPSB0aGlzLl9jdXJzb3IuZ2V0Q2hhcnMobmFtZVN0YXJ0KTtcbiAgICAgICAgdGhpcy5fY3Vyc29yLmFkdmFuY2UoKTtcbiAgICAgICAgY29uc3QgY2hhciA9IE5BTUVEX0VOVElUSUVTW25hbWVdO1xuICAgICAgICBpZiAoIWNoYXIpIHtcbiAgICAgICAgICB0aHJvdyB0aGlzLl9jcmVhdGVFcnJvcihfdW5rbm93bkVudGl0eUVycm9yTXNnKG5hbWUpLCB0aGlzLl9jdXJzb3IuZ2V0U3BhbihzdGFydCkpO1xuICAgICAgICB9XG4gICAgICAgIHRoaXMuX2VuZFRva2VuKFtjaGFyLCBgJiR7bmFtZX07YF0pO1xuICAgICAgfVxuICAgIH1cbiAgfVxuXG4gIHByaXZhdGUgX2NvbnN1bWVSYXdUZXh0KGNvbnN1bWVFbnRpdGllczogYm9vbGVhbiwgZW5kTWFya2VyUHJlZGljYXRlOiAoKSA9PiBib29sZWFuKTogdm9pZCB7XG4gICAgdGhpcy5fYmVnaW5Ub2tlbihjb25zdW1lRW50aXRpZXMgPyBUb2tlblR5cGUuRVNDQVBBQkxFX1JBV19URVhUIDogVG9rZW5UeXBlLlJBV19URVhUKTtcbiAgICBjb25zdCBwYXJ0czogc3RyaW5nW10gPSBbXTtcbiAgICB3aGlsZSAodHJ1ZSkge1xuICAgICAgY29uc3QgdGFnQ2xvc2VTdGFydCA9IHRoaXMuX2N1cnNvci5jbG9uZSgpO1xuICAgICAgY29uc3QgZm91bmRFbmRNYXJrZXIgPSBlbmRNYXJrZXJQcmVkaWNhdGUoKTtcbiAgICAgIHRoaXMuX2N1cnNvciA9IHRhZ0Nsb3NlU3RhcnQ7XG4gICAgICBpZiAoZm91bmRFbmRNYXJrZXIpIHtcbiAgICAgICAgYnJlYWs7XG4gICAgICB9XG4gICAgICBpZiAoY29uc3VtZUVudGl0aWVzICYmIHRoaXMuX2N1cnNvci5wZWVrKCkgPT09IGNoYXJzLiRBTVBFUlNBTkQpIHtcbiAgICAgICAgdGhpcy5fZW5kVG9rZW4oW3RoaXMuX3Byb2Nlc3NDYXJyaWFnZVJldHVybnMocGFydHMuam9pbignJykpXSk7XG4gICAgICAgIHBhcnRzLmxlbmd0aCA9IDA7XG4gICAgICAgIHRoaXMuX2NvbnN1bWVFbnRpdHkoVG9rZW5UeXBlLkVTQ0FQQUJMRV9SQVdfVEVYVCk7XG4gICAgICAgIHRoaXMuX2JlZ2luVG9rZW4oVG9rZW5UeXBlLkVTQ0FQQUJMRV9SQVdfVEVYVCk7XG4gICAgICB9IGVsc2Uge1xuICAgICAgICBwYXJ0cy5wdXNoKHRoaXMuX3JlYWRDaGFyKCkpO1xuICAgICAgfVxuICAgIH1cbiAgICB0aGlzLl9lbmRUb2tlbihbdGhpcy5fcHJvY2Vzc0NhcnJpYWdlUmV0dXJucyhwYXJ0cy5qb2luKCcnKSldKTtcbiAgfVxuXG4gIHByaXZhdGUgX2NvbnN1bWVDb21tZW50KHN0YXJ0OiBDaGFyYWN0ZXJDdXJzb3IpIHtcbiAgICB0aGlzLl9iZWdpblRva2VuKFRva2VuVHlwZS5DT01NRU5UX1NUQVJULCBzdGFydCk7XG4gICAgdGhpcy5fcmVxdWlyZUNoYXJDb2RlKGNoYXJzLiRNSU5VUyk7XG4gICAgdGhpcy5fZW5kVG9rZW4oW10pO1xuICAgIHRoaXMuX2NvbnN1bWVSYXdUZXh0KGZhbHNlLCAoKSA9PiB0aGlzLl9hdHRlbXB0U3RyKCctLT4nKSk7XG4gICAgdGhpcy5fYmVnaW5Ub2tlbihUb2tlblR5cGUuQ09NTUVOVF9FTkQpO1xuICAgIHRoaXMuX3JlcXVpcmVTdHIoJy0tPicpO1xuICAgIHRoaXMuX2VuZFRva2VuKFtdKTtcbiAgfVxuXG4gIHByaXZhdGUgX2NvbnN1bWVDZGF0YShzdGFydDogQ2hhcmFjdGVyQ3Vyc29yKSB7XG4gICAgdGhpcy5fYmVnaW5Ub2tlbihUb2tlblR5cGUuQ0RBVEFfU1RBUlQsIHN0YXJ0KTtcbiAgICB0aGlzLl9yZXF1aXJlU3RyKCdDREFUQVsnKTtcbiAgICB0aGlzLl9lbmRUb2tlbihbXSk7XG4gICAgdGhpcy5fY29uc3VtZVJhd1RleHQoZmFsc2UsICgpID0+IHRoaXMuX2F0dGVtcHRTdHIoJ11dPicpKTtcbiAgICB0aGlzLl9iZWdpblRva2VuKFRva2VuVHlwZS5DREFUQV9FTkQpO1xuICAgIHRoaXMuX3JlcXVpcmVTdHIoJ11dPicpO1xuICAgIHRoaXMuX2VuZFRva2VuKFtdKTtcbiAgfVxuXG4gIHByaXZhdGUgX2NvbnN1bWVEb2NUeXBlKHN0YXJ0OiBDaGFyYWN0ZXJDdXJzb3IpIHtcbiAgICB0aGlzLl9iZWdpblRva2VuKFRva2VuVHlwZS5ET0NfVFlQRSwgc3RhcnQpO1xuICAgIGNvbnN0IGNvbnRlbnRTdGFydCA9IHRoaXMuX2N1cnNvci5jbG9uZSgpO1xuICAgIHRoaXMuX2F0dGVtcHRVbnRpbENoYXIoY2hhcnMuJEdUKTtcbiAgICBjb25zdCBjb250ZW50ID0gdGhpcy5fY3Vyc29yLmdldENoYXJzKGNvbnRlbnRTdGFydCk7XG4gICAgdGhpcy5fY3Vyc29yLmFkdmFuY2UoKTtcbiAgICB0aGlzLl9lbmRUb2tlbihbY29udGVudF0pO1xuICB9XG5cbiAgcHJpdmF0ZSBfY29uc3VtZVByZWZpeEFuZE5hbWUoKTogc3RyaW5nW10ge1xuICAgIGNvbnN0IG5hbWVPclByZWZpeFN0YXJ0ID0gdGhpcy5fY3Vyc29yLmNsb25lKCk7XG4gICAgbGV0IHByZWZpeDogc3RyaW5nID0gJyc7XG4gICAgd2hpbGUgKHRoaXMuX2N1cnNvci5wZWVrKCkgIT09IGNoYXJzLiRDT0xPTiAmJiAhaXNQcmVmaXhFbmQodGhpcy5fY3Vyc29yLnBlZWsoKSkpIHtcbiAgICAgIHRoaXMuX2N1cnNvci5hZHZhbmNlKCk7XG4gICAgfVxuICAgIGxldCBuYW1lU3RhcnQ6IENoYXJhY3RlckN1cnNvcjtcbiAgICBpZiAodGhpcy5fY3Vyc29yLnBlZWsoKSA9PT0gY2hhcnMuJENPTE9OKSB7XG4gICAgICBwcmVmaXggPSB0aGlzLl9jdXJzb3IuZ2V0Q2hhcnMobmFtZU9yUHJlZml4U3RhcnQpO1xuICAgICAgdGhpcy5fY3Vyc29yLmFkdmFuY2UoKTtcbiAgICAgIG5hbWVTdGFydCA9IHRoaXMuX2N1cnNvci5jbG9uZSgpO1xuICAgIH0gZWxzZSB7XG4gICAgICBuYW1lU3RhcnQgPSBuYW1lT3JQcmVmaXhTdGFydDtcbiAgICB9XG4gICAgdGhpcy5fcmVxdWlyZUNoYXJDb2RlVW50aWxGbihpc05hbWVFbmQsIHByZWZpeCA9PT0gJycgPyAwIDogMSk7XG4gICAgY29uc3QgbmFtZSA9IHRoaXMuX2N1cnNvci5nZXRDaGFycyhuYW1lU3RhcnQpO1xuICAgIHJldHVybiBbcHJlZml4LCBuYW1lXTtcbiAgfVxuXG4gIHByaXZhdGUgX2NvbnN1bWVUYWdPcGVuKHN0YXJ0OiBDaGFyYWN0ZXJDdXJzb3IpIHtcbiAgICBsZXQgdGFnTmFtZTogc3RyaW5nO1xuICAgIGxldCBwcmVmaXg6IHN0cmluZztcbiAgICBsZXQgb3BlblRhZ1Rva2VuOiBUYWdPcGVuU3RhcnRUb2tlbiB8IEluY29tcGxldGVUYWdPcGVuVG9rZW4gfCB1bmRlZmluZWQ7XG4gICAgdHJ5IHtcbiAgICAgIGlmICghY2hhcnMuaXNBc2NpaUxldHRlcih0aGlzLl9jdXJzb3IucGVlaygpKSkge1xuICAgICAgICB0aHJvdyB0aGlzLl9jcmVhdGVFcnJvcihcbiAgICAgICAgICBfdW5leHBlY3RlZENoYXJhY3RlckVycm9yTXNnKHRoaXMuX2N1cnNvci5wZWVrKCkpLFxuICAgICAgICAgIHRoaXMuX2N1cnNvci5nZXRTcGFuKHN0YXJ0KSxcbiAgICAgICAgKTtcbiAgICAgIH1cblxuICAgICAgb3BlblRhZ1Rva2VuID0gdGhpcy5fY29uc3VtZVRhZ09wZW5TdGFydChzdGFydCk7XG4gICAgICBwcmVmaXggPSBvcGVuVGFnVG9rZW4ucGFydHNbMF07XG4gICAgICB0YWdOYW1lID0gb3BlblRhZ1Rva2VuLnBhcnRzWzFdO1xuICAgICAgdGhpcy5fYXR0ZW1wdENoYXJDb2RlVW50aWxGbihpc05vdFdoaXRlc3BhY2UpO1xuICAgICAgd2hpbGUgKFxuICAgICAgICB0aGlzLl9jdXJzb3IucGVlaygpICE9PSBjaGFycy4kU0xBU0ggJiZcbiAgICAgICAgdGhpcy5fY3Vyc29yLnBlZWsoKSAhPT0gY2hhcnMuJEdUICYmXG4gICAgICAgIHRoaXMuX2N1cnNvci5wZWVrKCkgIT09IGNoYXJzLiRMVCAmJlxuICAgICAgICB0aGlzLl9jdXJzb3IucGVlaygpICE9PSBjaGFycy4kRU9GXG4gICAgICApIHtcbiAgICAgICAgdGhpcy5fY29uc3VtZUF0dHJpYnV0ZU5hbWUoKTtcbiAgICAgICAgdGhpcy5fYXR0ZW1wdENoYXJDb2RlVW50aWxGbihpc05vdFdoaXRlc3BhY2UpO1xuICAgICAgICBpZiAodGhpcy5fYXR0ZW1wdENoYXJDb2RlKGNoYXJzLiRFUSkpIHtcbiAgICAgICAgICB0aGlzLl9hdHRlbXB0Q2hhckNvZGVVbnRpbEZuKGlzTm90V2hpdGVzcGFjZSk7XG4gICAgICAgICAgdGhpcy5fY29uc3VtZUF0dHJpYnV0ZVZhbHVlKCk7XG4gICAgICAgIH1cbiAgICAgICAgdGhpcy5fYXR0ZW1wdENoYXJDb2RlVW50aWxGbihpc05vdFdoaXRlc3BhY2UpO1xuICAgICAgfVxuICAgICAgdGhpcy5fY29uc3VtZVRhZ09wZW5FbmQoKTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICBpZiAoZSBpbnN0YW5jZW9mIF9Db250cm9sRmxvd0Vycm9yKSB7XG4gICAgICAgIGlmIChvcGVuVGFnVG9rZW4pIHtcbiAgICAgICAgICAvLyBXZSBlcnJvcmVkIGJlZm9yZSB3ZSBjb3VsZCBjbG9zZSB0aGUgb3BlbmluZyB0YWcsIHNvIGl0IGlzIGluY29tcGxldGUuXG4gICAgICAgICAgb3BlblRhZ1Rva2VuLnR5cGUgPSBUb2tlblR5cGUuSU5DT01QTEVURV9UQUdfT1BFTjtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAvLyBXaGVuIHRoZSBzdGFydCB0YWcgaXMgaW52YWxpZCwgYXNzdW1lIHdlIHdhbnQgYSBcIjxcIiBhcyB0ZXh0LlxuICAgICAgICAgIC8vIEJhY2sgdG8gYmFjayB0ZXh0IHRva2VucyBhcmUgbWVyZ2VkIGF0IHRoZSBlbmQuXG4gICAgICAgICAgdGhpcy5fYmVnaW5Ub2tlbihUb2tlblR5cGUuVEVYVCwgc3RhcnQpO1xuICAgICAgICAgIHRoaXMuX2VuZFRva2VuKFsnPCddKTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm47XG4gICAgICB9XG5cbiAgICAgIHRocm93IGU7XG4gICAgfVxuXG4gICAgY29uc3QgY29udGVudFRva2VuVHlwZSA9IHRoaXMuX2dldFRhZ0RlZmluaXRpb24odGFnTmFtZSkuZ2V0Q29udGVudFR5cGUocHJlZml4KTtcblxuICAgIGlmIChjb250ZW50VG9rZW5UeXBlID09PSBUYWdDb250ZW50VHlwZS5SQVdfVEVYVCkge1xuICAgICAgdGhpcy5fY29uc3VtZVJhd1RleHRXaXRoVGFnQ2xvc2UocHJlZml4LCB0YWdOYW1lLCBmYWxzZSk7XG4gICAgfSBlbHNlIGlmIChjb250ZW50VG9rZW5UeXBlID09PSBUYWdDb250ZW50VHlwZS5FU0NBUEFCTEVfUkFXX1RFWFQpIHtcbiAgICAgIHRoaXMuX2NvbnN1bWVSYXdUZXh0V2l0aFRhZ0Nsb3NlKHByZWZpeCwgdGFnTmFtZSwgdHJ1ZSk7XG4gICAgfVxuICB9XG5cbiAgcHJpdmF0ZSBfY29uc3VtZVJhd1RleHRXaXRoVGFnQ2xvc2UocHJlZml4OiBzdHJpbmcsIHRhZ05hbWU6IHN0cmluZywgY29uc3VtZUVudGl0aWVzOiBib29sZWFuKSB7XG4gICAgdGhpcy5fY29uc3VtZVJhd1RleHQoY29uc3VtZUVudGl0aWVzLCAoKSA9PiB7XG4gICAgICBpZiAoIXRoaXMuX2F0dGVtcHRDaGFyQ29kZShjaGFycy4kTFQpKSByZXR1cm4gZmFsc2U7XG4gICAgICBpZiAoIXRoaXMuX2F0dGVtcHRDaGFyQ29kZShjaGFycy4kU0xBU0gpKSByZXR1cm4gZmFsc2U7XG4gICAgICB0aGlzLl9hdHRlbXB0Q2hhckNvZGVVbnRpbEZuKGlzTm90V2hpdGVzcGFjZSk7XG4gICAgICBpZiAoIXRoaXMuX2F0dGVtcHRTdHJDYXNlSW5zZW5zaXRpdmUodGFnTmFtZSkpIHJldHVybiBmYWxzZTtcbiAgICAgIHRoaXMuX2F0dGVtcHRDaGFyQ29kZVVudGlsRm4oaXNOb3RXaGl0ZXNwYWNlKTtcbiAgICAgIHJldHVybiB0aGlzLl9hdHRlbXB0Q2hhckNvZGUoY2hhcnMuJEdUKTtcbiAgICB9KTtcbiAgICB0aGlzLl9iZWdpblRva2VuKFRva2VuVHlwZS5UQUdfQ0xPU0UpO1xuICAgIHRoaXMuX3JlcXVpcmVDaGFyQ29kZVVudGlsRm4oKGNvZGUpID0+IGNvZGUgPT09IGNoYXJzLiRHVCwgMyk7XG4gICAgdGhpcy5fY3Vyc29yLmFkdmFuY2UoKTsgLy8gQ29uc3VtZSB0aGUgYD5gXG4gICAgdGhpcy5fZW5kVG9rZW4oW3ByZWZpeCwgdGFnTmFtZV0pO1xuICB9XG5cbiAgcHJpdmF0ZSBfY29uc3VtZVRhZ09wZW5TdGFydChzdGFydDogQ2hhcmFjdGVyQ3Vyc29yKTogVGFnT3BlblN0YXJ0VG9rZW4ge1xuICAgIHRoaXMuX2JlZ2luVG9rZW4oVG9rZW5UeXBlLlRBR19PUEVOX1NUQVJULCBzdGFydCk7XG4gICAgY29uc3QgcGFydHMgPSB0aGlzLl9jb25zdW1lUHJlZml4QW5kTmFtZSgpO1xuICAgIHJldHVybiB0aGlzLl9lbmRUb2tlbihwYXJ0cykgYXMgVGFnT3BlblN0YXJ0VG9rZW47XG4gIH1cblxuICBwcml2YXRlIF9jb25zdW1lQXR0cmlidXRlTmFtZSgpIHtcbiAgICBjb25zdCBhdHRyTmFtZVN0YXJ0ID0gdGhpcy5fY3Vyc29yLnBlZWsoKTtcbiAgICBpZiAoYXR0ck5hbWVTdGFydCA9PT0gY2hhcnMuJFNRIHx8IGF0dHJOYW1lU3RhcnQgPT09IGNoYXJzLiREUSkge1xuICAgICAgdGhyb3cgdGhpcy5fY3JlYXRlRXJyb3IoX3VuZXhwZWN0ZWRDaGFyYWN0ZXJFcnJvck1zZyhhdHRyTmFtZVN0YXJ0KSwgdGhpcy5fY3Vyc29yLmdldFNwYW4oKSk7XG4gICAgfVxuICAgIHRoaXMuX2JlZ2luVG9rZW4oVG9rZW5UeXBlLkFUVFJfTkFNRSk7XG4gICAgY29uc3QgcHJlZml4QW5kTmFtZSA9IHRoaXMuX2NvbnN1bWVQcmVmaXhBbmROYW1lKCk7XG4gICAgdGhpcy5fZW5kVG9rZW4ocHJlZml4QW5kTmFtZSk7XG4gIH1cblxuICBwcml2YXRlIF9jb25zdW1lQXR0cmlidXRlVmFsdWUoKSB7XG4gICAgaWYgKHRoaXMuX2N1cnNvci5wZWVrKCkgPT09IGNoYXJzLiRTUSB8fCB0aGlzLl9jdXJzb3IucGVlaygpID09PSBjaGFycy4kRFEpIHtcbiAgICAgIGNvbnN0IHF1b3RlQ2hhciA9IHRoaXMuX2N1cnNvci5wZWVrKCk7XG4gICAgICB0aGlzLl9jb25zdW1lUXVvdGUocXVvdGVDaGFyKTtcbiAgICAgIC8vIEluIGFuIGF0dHJpYnV0ZSB0aGVuIGVuZCBvZiB0aGUgYXR0cmlidXRlIHZhbHVlIGFuZCB0aGUgcHJlbWF0dXJlIGVuZCB0byBhbiBpbnRlcnBvbGF0aW9uXG4gICAgICAvLyBhcmUgYm90aCB0cmlnZ2VyZWQgYnkgdGhlIGBxdW90ZUNoYXJgLlxuICAgICAgY29uc3QgZW5kUHJlZGljYXRlID0gKCkgPT4gdGhpcy5fY3Vyc29yLnBlZWsoKSA9PT0gcXVvdGVDaGFyO1xuICAgICAgdGhpcy5fY29uc3VtZVdpdGhJbnRlcnBvbGF0aW9uKFxuICAgICAgICBUb2tlblR5cGUuQVRUUl9WQUxVRV9URVhULFxuICAgICAgICBUb2tlblR5cGUuQVRUUl9WQUxVRV9JTlRFUlBPTEFUSU9OLFxuICAgICAgICBlbmRQcmVkaWNhdGUsXG4gICAgICAgIGVuZFByZWRpY2F0ZSxcbiAgICAgICk7XG4gICAgICB0aGlzLl9jb25zdW1lUXVvdGUocXVvdGVDaGFyKTtcbiAgICB9IGVsc2Uge1xuICAgICAgY29uc3QgZW5kUHJlZGljYXRlID0gKCkgPT4gaXNOYW1lRW5kKHRoaXMuX2N1cnNvci5wZWVrKCkpO1xuICAgICAgdGhpcy5fY29uc3VtZVdpdGhJbnRlcnBvbGF0aW9uKFxuICAgICAgICBUb2tlblR5cGUuQVRUUl9WQUxVRV9URVhULFxuICAgICAgICBUb2tlblR5cGUuQVRUUl9WQUxVRV9JTlRFUlBPTEFUSU9OLFxuICAgICAgICBlbmRQcmVkaWNhdGUsXG4gICAgICAgIGVuZFByZWRpY2F0ZSxcbiAgICAgICk7XG4gICAgfVxuICB9XG5cbiAgcHJpdmF0ZSBfY29uc3VtZVF1b3RlKHF1b3RlQ2hhcjogbnVtYmVyKSB7XG4gICAgdGhpcy5fYmVnaW5Ub2tlbihUb2tlblR5cGUuQVRUUl9RVU9URSk7XG4gICAgdGhpcy5fcmVxdWlyZUNoYXJDb2RlKHF1b3RlQ2hhcik7XG4gICAgdGhpcy5fZW5kVG9rZW4oW1N0cmluZy5mcm9tQ29kZVBvaW50KHF1b3RlQ2hhcildKTtcbiAgfVxuXG4gIHByaXZhdGUgX2NvbnN1bWVUYWdPcGVuRW5kKCkge1xuICAgIGNvbnN0IHRva2VuVHlwZSA9IHRoaXMuX2F0dGVtcHRDaGFyQ29kZShjaGFycy4kU0xBU0gpXG4gICAgICA/IFRva2VuVHlwZS5UQUdfT1BFTl9FTkRfVk9JRFxuICAgICAgOiBUb2tlblR5cGUuVEFHX09QRU5fRU5EO1xuICAgIHRoaXMuX2JlZ2luVG9rZW4odG9rZW5UeXBlKTtcbiAgICB0aGlzLl9yZXF1aXJlQ2hhckNvZGUoY2hhcnMuJEdUKTtcbiAgICB0aGlzLl9lbmRUb2tlbihbXSk7XG4gIH1cblxuICBwcml2YXRlIF9jb25zdW1lVGFnQ2xvc2Uoc3RhcnQ6IENoYXJhY3RlckN1cnNvcikge1xuICAgIHRoaXMuX2JlZ2luVG9rZW4oVG9rZW5UeXBlLlRBR19DTE9TRSwgc3RhcnQpO1xuICAgIHRoaXMuX2F0dGVtcHRDaGFyQ29kZVVudGlsRm4oaXNOb3RXaGl0ZXNwYWNlKTtcbiAgICBjb25zdCBwcmVmaXhBbmROYW1lID0gdGhpcy5fY29uc3VtZVByZWZpeEFuZE5hbWUoKTtcbiAgICB0aGlzLl9hdHRlbXB0Q2hhckNvZGVVbnRpbEZuKGlzTm90V2hpdGVzcGFjZSk7XG4gICAgdGhpcy5fcmVxdWlyZUNoYXJDb2RlKGNoYXJzLiRHVCk7XG4gICAgdGhpcy5fZW5kVG9rZW4ocHJlZml4QW5kTmFtZSk7XG4gIH1cblxuICBwcml2YXRlIF9jb25zdW1lRXhwYW5zaW9uRm9ybVN0YXJ0KCkge1xuICAgIHRoaXMuX2JlZ2luVG9rZW4oVG9rZW5UeXBlLkVYUEFOU0lPTl9GT1JNX1NUQVJUKTtcbiAgICB0aGlzLl9yZXF1aXJlQ2hhckNvZGUoY2hhcnMuJExCUkFDRSk7XG4gICAgdGhpcy5fZW5kVG9rZW4oW10pO1xuXG4gICAgdGhpcy5fZXhwYW5zaW9uQ2FzZVN0YWNrLnB1c2goVG9rZW5UeXBlLkVYUEFOU0lPTl9GT1JNX1NUQVJUKTtcblxuICAgIHRoaXMuX2JlZ2luVG9rZW4oVG9rZW5UeXBlLlJBV19URVhUKTtcbiAgICBjb25zdCBjb25kaXRpb24gPSB0aGlzLl9yZWFkVW50aWwoY2hhcnMuJENPTU1BKTtcbiAgICBjb25zdCBub3JtYWxpemVkQ29uZGl0aW9uID0gdGhpcy5fcHJvY2Vzc0NhcnJpYWdlUmV0dXJucyhjb25kaXRpb24pO1xuICAgIGlmICh0aGlzLl9pMThuTm9ybWFsaXplTGluZUVuZGluZ3NJbklDVXMpIHtcbiAgICAgIC8vIFdlIGV4cGxpY2l0bHkgd2FudCB0byBub3JtYWxpemUgbGluZSBlbmRpbmdzIGZvciB0aGlzIHRleHQuXG4gICAgICB0aGlzLl9lbmRUb2tlbihbbm9ybWFsaXplZENvbmRpdGlvbl0pO1xuICAgIH0gZWxzZSB7XG4gICAgICAvLyBXZSBhcmUgbm90IG5vcm1hbGl6aW5nIGxpbmUgZW5kaW5ncy5cbiAgICAgIGNvbnN0IGNvbmRpdGlvblRva2VuID0gdGhpcy5fZW5kVG9rZW4oW2NvbmRpdGlvbl0pO1xuICAgICAgaWYgKG5vcm1hbGl6ZWRDb25kaXRpb24gIT09IGNvbmRpdGlvbikge1xuICAgICAgICB0aGlzLm5vbk5vcm1hbGl6ZWRJY3VFeHByZXNzaW9ucy5wdXNoKGNvbmRpdGlvblRva2VuKTtcbiAgICAgIH1cbiAgICB9XG4gICAgdGhpcy5fcmVxdWlyZUNoYXJDb2RlKGNoYXJzLiRDT01NQSk7XG4gICAgdGhpcy5fYXR0ZW1wdENoYXJDb2RlVW50aWxGbihpc05vdFdoaXRlc3BhY2UpO1xuXG4gICAgdGhpcy5fYmVnaW5Ub2tlbihUb2tlblR5cGUuUkFXX1RFWFQpO1xuICAgIGNvbnN0IHR5cGUgPSB0aGlzLl9yZWFkVW50aWwoY2hhcnMuJENPTU1BKTtcbiAgICB0aGlzLl9lbmRUb2tlbihbdHlwZV0pO1xuICAgIHRoaXMuX3JlcXVpcmVDaGFyQ29kZShjaGFycy4kQ09NTUEpO1xuICAgIHRoaXMuX2F0dGVtcHRDaGFyQ29kZVVudGlsRm4oaXNOb3RXaGl0ZXNwYWNlKTtcbiAgfVxuXG4gIHByaXZhdGUgX2NvbnN1bWVFeHBhbnNpb25DYXNlU3RhcnQoKSB7XG4gICAgdGhpcy5fYmVnaW5Ub2tlbihUb2tlblR5cGUuRVhQQU5TSU9OX0NBU0VfVkFMVUUpO1xuICAgIGNvbnN0IHZhbHVlID0gdGhpcy5fcmVhZFVudGlsKGNoYXJzLiRMQlJBQ0UpLnRyaW0oKTtcbiAgICB0aGlzLl9lbmRUb2tlbihbdmFsdWVdKTtcbiAgICB0aGlzLl9hdHRlbXB0Q2hhckNvZGVVbnRpbEZuKGlzTm90V2hpdGVzcGFjZSk7XG5cbiAgICB0aGlzLl9iZWdpblRva2VuKFRva2VuVHlwZS5FWFBBTlNJT05fQ0FTRV9FWFBfU1RBUlQpO1xuICAgIHRoaXMuX3JlcXVpcmVDaGFyQ29kZShjaGFycy4kTEJSQUNFKTtcbiAgICB0aGlzLl9lbmRUb2tlbihbXSk7XG4gICAgdGhpcy5fYXR0ZW1wdENoYXJDb2RlVW50aWxGbihpc05vdFdoaXRlc3BhY2UpO1xuXG4gICAgdGhpcy5fZXhwYW5zaW9uQ2FzZVN0YWNrLnB1c2goVG9rZW5UeXBlLkVYUEFOU0lPTl9DQVNFX0VYUF9TVEFSVCk7XG4gIH1cblxuICBwcml2YXRlIF9jb25zdW1lRXhwYW5zaW9uQ2FzZUVuZCgpIHtcbiAgICB0aGlzLl9iZWdpblRva2VuKFRva2VuVHlwZS5FWFBBTlNJT05fQ0FTRV9FWFBfRU5EKTtcbiAgICB0aGlzLl9yZXF1aXJlQ2hhckNvZGUoY2hhcnMuJFJCUkFDRSk7XG4gICAgdGhpcy5fZW5kVG9rZW4oW10pO1xuICAgIHRoaXMuX2F0dGVtcHRDaGFyQ29kZVVudGlsRm4oaXNOb3RXaGl0ZXNwYWNlKTtcblxuICAgIHRoaXMuX2V4cGFuc2lvbkNhc2VTdGFjay5wb3AoKTtcbiAgfVxuXG4gIHByaXZhdGUgX2NvbnN1bWVFeHBhbnNpb25Gb3JtRW5kKCkge1xuICAgIHRoaXMuX2JlZ2luVG9rZW4oVG9rZW5UeXBlLkVYUEFOU0lPTl9GT1JNX0VORCk7XG4gICAgdGhpcy5fcmVxdWlyZUNoYXJDb2RlKGNoYXJzLiRSQlJBQ0UpO1xuICAgIHRoaXMuX2VuZFRva2VuKFtdKTtcblxuICAgIHRoaXMuX2V4cGFuc2lvbkNhc2VTdGFjay5wb3AoKTtcbiAgfVxuXG4gIC8qKlxuICAgKiBDb25zdW1lIGEgc3RyaW5nIHRoYXQgbWF5IGNvbnRhaW4gaW50ZXJwb2xhdGlvbiBleHByZXNzaW9ucy5cbiAgICpcbiAgICogVGhlIGZpcnN0IHRva2VuIGNvbnN1bWVkIHdpbGwgYmUgb2YgYHRva2VuVHlwZWAgYW5kIHRoZW4gdGhlcmUgd2lsbCBiZSBhbHRlcm5hdGluZ1xuICAgKiBgaW50ZXJwb2xhdGlvblRva2VuVHlwZWAgYW5kIGB0b2tlblR5cGVgIHRva2VucyB1bnRpbCB0aGUgYGVuZFByZWRpY2F0ZSgpYCByZXR1cm5zIHRydWUuXG4gICAqXG4gICAqIElmIGFuIGludGVycG9sYXRpb24gdG9rZW4gZW5kcyBwcmVtYXR1cmVseSBpdCB3aWxsIGhhdmUgbm8gZW5kIG1hcmtlciBpbiBpdHMgYHBhcnRzYCBhcnJheS5cbiAgICpcbiAgICogQHBhcmFtIHRleHRUb2tlblR5cGUgdGhlIGtpbmQgb2YgdG9rZW5zIHRvIGludGVybGVhdmUgYXJvdW5kIGludGVycG9sYXRpb24gdG9rZW5zLlxuICAgKiBAcGFyYW0gaW50ZXJwb2xhdGlvblRva2VuVHlwZSB0aGUga2luZCBvZiB0b2tlbnMgdGhhdCBjb250YWluIGludGVycG9sYXRpb24uXG4gICAqIEBwYXJhbSBlbmRQcmVkaWNhdGUgYSBmdW5jdGlvbiB0aGF0IHNob3VsZCByZXR1cm4gdHJ1ZSB3aGVuIHdlIHNob3VsZCBzdG9wIGNvbnN1bWluZy5cbiAgICogQHBhcmFtIGVuZEludGVycG9sYXRpb24gYSBmdW5jdGlvbiB0aGF0IHNob3VsZCByZXR1cm4gdHJ1ZSBpZiB0aGVyZSBpcyBhIHByZW1hdHVyZSBlbmQgdG8gYW5cbiAgICogICAgIGludGVycG9sYXRpb24gZXhwcmVzc2lvbiAtIGkuZS4gYmVmb3JlIHdlIGdldCB0byB0aGUgbm9ybWFsIGludGVycG9sYXRpb24gY2xvc2luZyBtYXJrZXIuXG4gICAqL1xuICBwcml2YXRlIF9jb25zdW1lV2l0aEludGVycG9sYXRpb24oXG4gICAgdGV4dFRva2VuVHlwZTogVG9rZW5UeXBlLFxuICAgIGludGVycG9sYXRpb25Ub2tlblR5cGU6IFRva2VuVHlwZSxcbiAgICBlbmRQcmVkaWNhdGU6ICgpID0+IGJvb2xlYW4sXG4gICAgZW5kSW50ZXJwb2xhdGlvbjogKCkgPT4gYm9vbGVhbixcbiAgKSB7XG4gICAgdGhpcy5fYmVnaW5Ub2tlbih0ZXh0VG9rZW5UeXBlKTtcbiAgICBjb25zdCBwYXJ0czogc3RyaW5nW10gPSBbXTtcblxuICAgIHdoaWxlICghZW5kUHJlZGljYXRlKCkpIHtcbiAgICAgIGNvbnN0IGN1cnJlbnQgPSB0aGlzLl9jdXJzb3IuY2xvbmUoKTtcbiAgICAgIGlmICh0aGlzLl9pbnRlcnBvbGF0aW9uQ29uZmlnICYmIHRoaXMuX2F0dGVtcHRTdHIodGhpcy5faW50ZXJwb2xhdGlvbkNvbmZpZy5zdGFydCkpIHtcbiAgICAgICAgdGhpcy5fZW5kVG9rZW4oW3RoaXMuX3Byb2Nlc3NDYXJyaWFnZVJldHVybnMocGFydHMuam9pbignJykpXSwgY3VycmVudCk7XG4gICAgICAgIHBhcnRzLmxlbmd0aCA9IDA7XG4gICAgICAgIHRoaXMuX2NvbnN1bWVJbnRlcnBvbGF0aW9uKGludGVycG9sYXRpb25Ub2tlblR5cGUsIGN1cnJlbnQsIGVuZEludGVycG9sYXRpb24pO1xuICAgICAgICB0aGlzLl9iZWdpblRva2VuKHRleHRUb2tlblR5cGUpO1xuICAgICAgfSBlbHNlIGlmICh0aGlzLl9jdXJzb3IucGVlaygpID09PSBjaGFycy4kQU1QRVJTQU5EKSB7XG4gICAgICAgIHRoaXMuX2VuZFRva2VuKFt0aGlzLl9wcm9jZXNzQ2FycmlhZ2VSZXR1cm5zKHBhcnRzLmpvaW4oJycpKV0pO1xuICAgICAgICBwYXJ0cy5sZW5ndGggPSAwO1xuICAgICAgICB0aGlzLl9jb25zdW1lRW50aXR5KHRleHRUb2tlblR5cGUpO1xuICAgICAgICB0aGlzLl9iZWdpblRva2VuKHRleHRUb2tlblR5cGUpO1xuICAgICAgfSBlbHNlIHtcbiAgICAgICAgcGFydHMucHVzaCh0aGlzLl9yZWFkQ2hhcigpKTtcbiAgICAgIH1cbiAgICB9XG5cbiAgICAvLyBJdCBpcyBwb3NzaWJsZSB0aGF0IGFuIGludGVycG9sYXRpb24gd2FzIHN0YXJ0ZWQgYnV0IG5vdCBlbmRlZCBpbnNpZGUgdGhpcyB0ZXh0IHRva2VuLlxuICAgIC8vIE1ha2Ugc3VyZSB0aGF0IHdlIHJlc2V0IHRoZSBzdGF0ZSBvZiB0aGUgbGV4ZXIgY29ycmVjdGx5LlxuICAgIHRoaXMuX2luSW50ZXJwb2xhdGlvbiA9IGZhbHNlO1xuXG4gICAgdGhpcy5fZW5kVG9rZW4oW3RoaXMuX3Byb2Nlc3NDYXJyaWFnZVJldHVybnMocGFydHMuam9pbignJykpXSk7XG4gIH1cblxuICAvKipcbiAgICogQ29uc3VtZSBhIGJsb2NrIG9mIHRleHQgdGhhdCBoYXMgYmVlbiBpbnRlcnByZXRlZCBhcyBhbiBBbmd1bGFyIGludGVycG9sYXRpb24uXG4gICAqXG4gICAqIEBwYXJhbSBpbnRlcnBvbGF0aW9uVG9rZW5UeXBlIHRoZSB0eXBlIG9mIHRoZSBpbnRlcnBvbGF0aW9uIHRva2VuIHRvIGdlbmVyYXRlLlxuICAgKiBAcGFyYW0gaW50ZXJwb2xhdGlvblN0YXJ0IGEgY3Vyc29yIHRoYXQgcG9pbnRzIHRvIHRoZSBzdGFydCBvZiB0aGlzIGludGVycG9sYXRpb24uXG4gICAqIEBwYXJhbSBwcmVtYXR1cmVFbmRQcmVkaWNhdGUgYSBmdW5jdGlvbiB0aGF0IHNob3VsZCByZXR1cm4gdHJ1ZSBpZiB0aGUgbmV4dCBjaGFyYWN0ZXJzIGluZGljYXRlXG4gICAqICAgICBhbiBlbmQgdG8gdGhlIGludGVycG9sYXRpb24gYmVmb3JlIGl0cyBub3JtYWwgY2xvc2luZyBtYXJrZXIuXG4gICAqL1xuICBwcml2YXRlIF9jb25zdW1lSW50ZXJwb2xhdGlvbihcbiAgICBpbnRlcnBvbGF0aW9uVG9rZW5UeXBlOiBUb2tlblR5cGUsXG4gICAgaW50ZXJwb2xhdGlvblN0YXJ0OiBDaGFyYWN0ZXJDdXJzb3IsXG4gICAgcHJlbWF0dXJlRW5kUHJlZGljYXRlOiAoKCkgPT4gYm9vbGVhbikgfCBudWxsLFxuICApOiB2b2lkIHtcbiAgICBjb25zdCBwYXJ0czogc3RyaW5nW10gPSBbXTtcbiAgICB0aGlzLl9iZWdpblRva2VuKGludGVycG9sYXRpb25Ub2tlblR5cGUsIGludGVycG9sYXRpb25TdGFydCk7XG4gICAgcGFydHMucHVzaCh0aGlzLl9pbnRlcnBvbGF0aW9uQ29uZmlnLnN0YXJ0KTtcblxuICAgIC8vIEZpbmQgdGhlIGVuZCBvZiB0aGUgaW50ZXJwb2xhdGlvbiwgaWdub3JpbmcgY29udGVudCBpbnNpZGUgcXVvdGVzLlxuICAgIGNvbnN0IGV4cHJlc3Npb25TdGFydCA9IHRoaXMuX2N1cnNvci5jbG9uZSgpO1xuICAgIGxldCBpblF1b3RlOiBudW1iZXIgfCBudWxsID0gbnVsbDtcbiAgICBsZXQgaW5Db21tZW50ID0gZmFsc2U7XG4gICAgd2hpbGUgKFxuICAgICAgdGhpcy5fY3Vyc29yLnBlZWsoKSAhPT0gY2hhcnMuJEVPRiAmJlxuICAgICAgKHByZW1hdHVyZUVuZFByZWRpY2F0ZSA9PT0gbnVsbCB8fCAhcHJlbWF0dXJlRW5kUHJlZGljYXRlKCkpXG4gICAgKSB7XG4gICAgICBjb25zdCBjdXJyZW50ID0gdGhpcy5fY3Vyc29yLmNsb25lKCk7XG5cbiAgICAgIGlmICh0aGlzLl9pc1RhZ1N0YXJ0KCkpIHtcbiAgICAgICAgLy8gV2UgYXJlIHN0YXJ0aW5nIHdoYXQgbG9va3MgbGlrZSBhbiBIVE1MIGVsZW1lbnQgaW4gdGhlIG1pZGRsZSBvZiB0aGlzIGludGVycG9sYXRpb24uXG4gICAgICAgIC8vIFJlc2V0IHRoZSBjdXJzb3IgdG8gYmVmb3JlIHRoZSBgPGAgY2hhcmFjdGVyIGFuZCBlbmQgdGhlIGludGVycG9sYXRpb24gdG9rZW4uXG4gICAgICAgIC8vIChUaGlzIGlzIGFjdHVhbGx5IHdyb25nIGJ1dCBoZXJlIGZvciBiYWNrd2FyZCBjb21wYXRpYmlsaXR5KS5cbiAgICAgICAgdGhpcy5fY3Vyc29yID0gY3VycmVudDtcbiAgICAgICAgcGFydHMucHVzaCh0aGlzLl9nZXRQcm9jZXNzZWRDaGFycyhleHByZXNzaW9uU3RhcnQsIGN1cnJlbnQpKTtcbiAgICAgICAgdGhpcy5fZW5kVG9rZW4ocGFydHMpO1xuICAgICAgICByZXR1cm47XG4gICAgICB9XG5cbiAgICAgIGlmIChpblF1b3RlID09PSBudWxsKSB7XG4gICAgICAgIGlmICh0aGlzLl9hdHRlbXB0U3RyKHRoaXMuX2ludGVycG9sYXRpb25Db25maWcuZW5kKSkge1xuICAgICAgICAgIC8vIFdlIGFyZSBub3QgaW4gYSBzdHJpbmcsIGFuZCB3ZSBoaXQgdGhlIGVuZCBpbnRlcnBvbGF0aW9uIG1hcmtlclxuICAgICAgICAgIHBhcnRzLnB1c2godGhpcy5fZ2V0UHJvY2Vzc2VkQ2hhcnMoZXhwcmVzc2lvblN0YXJ0LCBjdXJyZW50KSk7XG4gICAgICAgICAgcGFydHMucHVzaCh0aGlzLl9pbnRlcnBvbGF0aW9uQ29uZmlnLmVuZCk7XG4gICAgICAgICAgdGhpcy5fZW5kVG9rZW4ocGFydHMpO1xuICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfSBlbHNlIGlmICh0aGlzLl9hdHRlbXB0U3RyKCcvLycpKSB7XG4gICAgICAgICAgLy8gT25jZSB3ZSBhcmUgaW4gYSBjb21tZW50IHdlIGlnbm9yZSBhbnkgcXVvdGVzXG4gICAgICAgICAgaW5Db21tZW50ID0gdHJ1ZTtcbiAgICAgICAgfVxuICAgICAgfVxuXG4gICAgICBjb25zdCBjaGFyID0gdGhpcy5fY3Vyc29yLnBlZWsoKTtcbiAgICAgIHRoaXMuX2N1cnNvci5hZHZhbmNlKCk7XG4gICAgICBpZiAoY2hhciA9PT0gY2hhcnMuJEJBQ0tTTEFTSCkge1xuICAgICAgICAvLyBTa2lwIHRoZSBuZXh0IGNoYXJhY3RlciBiZWNhdXNlIGl0IHdhcyBlc2NhcGVkLlxuICAgICAgICB0aGlzLl9jdXJzb3IuYWR2YW5jZSgpO1xuICAgICAgfSBlbHNlIGlmIChjaGFyID09PSBpblF1b3RlKSB7XG4gICAgICAgIC8vIEV4aXRpbmcgdGhlIGN1cnJlbnQgcXVvdGVkIHN0cmluZ1xuICAgICAgICBpblF1b3RlID0gbnVsbDtcbiAgICAgIH0gZWxzZSBpZiAoIWluQ29tbWVudCAmJiBpblF1b3RlID09PSBudWxsICYmIGNoYXJzLmlzUXVvdGUoY2hhcikpIHtcbiAgICAgICAgLy8gRW50ZXJpbmcgYSBuZXcgcXVvdGVkIHN0cmluZ1xuICAgICAgICBpblF1b3RlID0gY2hhcjtcbiAgICAgIH1cbiAgICB9XG5cbiAgICAvLyBXZSBoaXQgRU9GIHdpdGhvdXQgZmluZGluZyBhIGNsb3NpbmcgaW50ZXJwb2xhdGlvbiBtYXJrZXJcbiAgICBwYXJ0cy5wdXNoKHRoaXMuX2dldFByb2Nlc3NlZENoYXJzKGV4cHJlc3Npb25TdGFydCwgdGhpcy5fY3Vyc29yKSk7XG4gICAgdGhpcy5fZW5kVG9rZW4ocGFydHMpO1xuICB9XG5cbiAgcHJpdmF0ZSBfZ2V0UHJvY2Vzc2VkQ2hhcnMoc3RhcnQ6IENoYXJhY3RlckN1cnNvciwgZW5kOiBDaGFyYWN0ZXJDdXJzb3IpOiBzdHJpbmcge1xuICAgIHJldHVybiB0aGlzLl9wcm9jZXNzQ2FycmlhZ2VSZXR1cm5zKGVuZC5nZXRDaGFycyhzdGFydCkpO1xuICB9XG5cbiAgcHJpdmF0ZSBfaXNUZXh0RW5kKCk6IGJvb2xlYW4ge1xuICAgIGlmICh0aGlzLl9pc1RhZ1N0YXJ0KCkgfHwgdGhpcy5fY3Vyc29yLnBlZWsoKSA9PT0gY2hhcnMuJEVPRikge1xuICAgICAgcmV0dXJuIHRydWU7XG4gICAgfVxuXG4gICAgaWYgKHRoaXMuX3Rva2VuaXplSWN1ICYmICF0aGlzLl9pbkludGVycG9sYXRpb24pIHtcbiAgICAgIGlmICh0aGlzLmlzRXhwYW5zaW9uRm9ybVN0YXJ0KCkpIHtcbiAgICAgICAgLy8gc3RhcnQgb2YgYW4gZXhwYW5zaW9uIGZvcm1cbiAgICAgICAgcmV0dXJuIHRydWU7XG4gICAgICB9XG5cbiAgICAgIGlmICh0aGlzLl9jdXJzb3IucGVlaygpID09PSBjaGFycy4kUkJSQUNFICYmIHRoaXMuX2lzSW5FeHBhbnNpb25DYXNlKCkpIHtcbiAgICAgICAgLy8gZW5kIG9mIGFuZCBleHBhbnNpb24gY2FzZVxuICAgICAgICByZXR1cm4gdHJ1ZTtcbiAgICAgIH1cbiAgICB9XG5cbiAgICBpZiAoXG4gICAgICB0aGlzLl90b2tlbml6ZUJsb2NrcyAmJlxuICAgICAgIXRoaXMuX2luSW50ZXJwb2xhdGlvbiAmJlxuICAgICAgIXRoaXMuX2lzSW5FeHBhbnNpb24oKSAmJlxuICAgICAgKHRoaXMuX2N1cnNvci5wZWVrKCkgPT09IGNoYXJzLiRBVCB8fCB0aGlzLl9jdXJzb3IucGVlaygpID09PSBjaGFycy4kUkJSQUNFKVxuICAgICkge1xuICAgICAgcmV0dXJuIHRydWU7XG4gICAgfVxuXG4gICAgcmV0dXJuIGZhbHNlO1xuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgY3VycmVudCBjdXJzb3IgaXMgcG9pbnRpbmcgdG8gdGhlIHN0YXJ0IG9mIGEgdGFnXG4gICAqIChvcGVuaW5nL2Nsb3NpbmcvY29tbWVudHMvY2RhdGEvZXRjKS5cbiAgICovXG4gIHByaXZhdGUgX2lzVGFnU3RhcnQoKTogYm9vbGVhbiB7XG4gICAgaWYgKHRoaXMuX2N1cnNvci5wZWVrKCkgPT09IGNoYXJzLiRMVCkge1xuICAgICAgLy8gV2UgYXNzdW1lIHRoYXQgYDxgIGZvbGxvd2VkIGJ5IHdoaXRlc3BhY2UgaXMgbm90IHRoZSBzdGFydCBvZiBhbiBIVE1MIGVsZW1lbnQuXG4gICAgICBjb25zdCB0bXAgPSB0aGlzLl9jdXJzb3IuY2xvbmUoKTtcbiAgICAgIHRtcC5hZHZhbmNlKCk7XG4gICAgICAvLyBJZiB0aGUgbmV4dCBjaGFyYWN0ZXIgaXMgYWxwaGFiZXRpYywgISBub3IgLyB0aGVuIGl0IGlzIGEgdGFnIHN0YXJ0XG4gICAgICBjb25zdCBjb2RlID0gdG1wLnBlZWsoKTtcbiAgICAgIGlmIChcbiAgICAgICAgKGNoYXJzLiRhIDw9IGNvZGUgJiYgY29kZSA8PSBjaGFycy4keikgfHxcbiAgICAgICAgKGNoYXJzLiRBIDw9IGNvZGUgJiYgY29kZSA8PSBjaGFycy4kWikgfHxcbiAgICAgICAgY29kZSA9PT0gY2hhcnMuJFNMQVNIIHx8XG4gICAgICAgIGNvZGUgPT09IGNoYXJzLiRCQU5HXG4gICAgICApIHtcbiAgICAgICAgcmV0dXJuIHRydWU7XG4gICAgICB9XG4gICAgfVxuICAgIHJldHVybiBmYWxzZTtcbiAgfVxuXG4gIHByaXZhdGUgX3JlYWRVbnRpbChjaGFyOiBudW1iZXIpOiBzdHJpbmcge1xuICAgIGNvbnN0IHN0YXJ0ID0gdGhpcy5fY3Vyc29yLmNsb25lKCk7XG4gICAgdGhpcy5fYXR0ZW1wdFVudGlsQ2hhcihjaGFyKTtcbiAgICByZXR1cm4gdGhpcy5fY3Vyc29yLmdldENoYXJzKHN0YXJ0KTtcbiAgfVxuXG4gIHByaXZhdGUgX2lzSW5FeHBhbnNpb24oKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIHRoaXMuX2lzSW5FeHBhbnNpb25DYXNlKCkgfHwgdGhpcy5faXNJbkV4cGFuc2lvbkZvcm0oKTtcbiAgfVxuXG4gIHByaXZhdGUgX2lzSW5FeHBhbnNpb25DYXNlKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiAoXG4gICAgICB0aGlzLl9leHBhbnNpb25DYXNlU3RhY2subGVuZ3RoID4gMCAmJlxuICAgICAgdGhpcy5fZXhwYW5zaW9uQ2FzZVN0YWNrW3RoaXMuX2V4cGFuc2lvbkNhc2VTdGFjay5sZW5ndGggLSAxXSA9PT1cbiAgICAgICAgVG9rZW5UeXBlLkVYUEFOU0lPTl9DQVNFX0VYUF9TVEFSVFxuICAgICk7XG4gIH1cblxuICBwcml2YXRlIF9pc0luRXhwYW5zaW9uRm9ybSgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gKFxuICAgICAgdGhpcy5fZXhwYW5zaW9uQ2FzZVN0YWNrLmxlbmd0aCA+IDAgJiZcbiAgICAgIHRoaXMuX2V4cGFuc2lvbkNhc2VTdGFja1t0aGlzLl9leHBhbnNpb25DYXNlU3RhY2subGVuZ3RoIC0gMV0gPT09XG4gICAgICAgIFRva2VuVHlwZS5FWFBBTlNJT05fRk9STV9TVEFSVFxuICAgICk7XG4gIH1cblxuICBwcml2YXRlIGlzRXhwYW5zaW9uRm9ybVN0YXJ0KCk6IGJvb2xlYW4ge1xuICAgIGlmICh0aGlzLl9jdXJzb3IucGVlaygpICE9PSBjaGFycy4kTEJSQUNFKSB7XG4gICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuICAgIGlmICh0aGlzLl9pbnRlcnBvbGF0aW9uQ29uZmlnKSB7XG4gICAgICBjb25zdCBzdGFydCA9IHRoaXMuX2N1cnNvci5jbG9uZSgpO1xuICAgICAgY29uc3QgaXNJbnRlcnBvbGF0aW9uID0gdGhpcy5fYXR0ZW1wdFN0cih0aGlzLl9pbnRlcnBvbGF0aW9uQ29uZmlnLnN0YXJ0KTtcbiAgICAgIHRoaXMuX2N1cnNvciA9IHN0YXJ0O1xuICAgICAgcmV0dXJuICFpc0ludGVycG9sYXRpb247XG4gICAgfVxuICAgIHJldHVybiB0cnVlO1xuICB9XG59XG5cbmZ1bmN0aW9uIGlzTm90V2hpdGVzcGFjZShjb2RlOiBudW1iZXIpOiBib29sZWFuIHtcbiAgcmV0dXJuICFjaGFycy5pc1doaXRlc3BhY2UoY29kZSkgfHwgY29kZSA9PT0gY2hhcnMuJEVPRjtcbn1cblxuZnVuY3Rpb24gaXNOYW1lRW5kKGNvZGU6IG51bWJlcik6IGJvb2xlYW4ge1xuICByZXR1cm4gKFxuICAgIGNoYXJzLmlzV2hpdGVzcGFjZShjb2RlKSB8fFxuICAgIGNvZGUgPT09IGNoYXJzLiRHVCB8fFxuICAgIGNvZGUgPT09IGNoYXJzLiRMVCB8fFxuICAgIGNvZGUgPT09IGNoYXJzLiRTTEFTSCB8fFxuICAgIGNvZGUgPT09IGNoYXJzLiRTUSB8fFxuICAgIGNvZGUgPT09IGNoYXJzLiREUSB8fFxuICAgIGNvZGUgPT09IGNoYXJzLiRFUSB8fFxuICAgIGNvZGUgPT09IGNoYXJzLiRFT0ZcbiAgKTtcbn1cblxuZnVuY3Rpb24gaXNQcmVmaXhFbmQoY29kZTogbnVtYmVyKTogYm9vbGVhbiB7XG4gIHJldHVybiAoXG4gICAgKGNvZGUgPCBjaGFycy4kYSB8fCBjaGFycy4keiA8IGNvZGUpICYmXG4gICAgKGNvZGUgPCBjaGFycy4kQSB8fCBjaGFycy4kWiA8IGNvZGUpICYmXG4gICAgKGNvZGUgPCBjaGFycy4kMCB8fCBjb2RlID4gY2hhcnMuJDkpXG4gICk7XG59XG5cbmZ1bmN0aW9uIGlzRGlnaXRFbnRpdHlFbmQoY29kZTogbnVtYmVyKTogYm9vbGVhbiB7XG4gIHJldHVybiBjb2RlID09PSBjaGFycy4kU0VNSUNPTE9OIHx8IGNvZGUgPT09IGNoYXJzLiRFT0YgfHwgIWNoYXJzLmlzQXNjaWlIZXhEaWdpdChjb2RlKTtcbn1cblxuZnVuY3Rpb24gaXNOYW1lZEVudGl0eUVuZChjb2RlOiBudW1iZXIpOiBib29sZWFuIHtcbiAgcmV0dXJuIGNvZGUgPT09IGNoYXJzLiRTRU1JQ09MT04gfHwgY29kZSA9PT0gY2hhcnMuJEVPRiB8fCAhY2hhcnMuaXNBc2NpaUxldHRlcihjb2RlKTtcbn1cblxuZnVuY3Rpb24gaXNFeHBhbnNpb25DYXNlU3RhcnQocGVlazogbnVtYmVyKTogYm9vbGVhbiB7XG4gIHJldHVybiBwZWVrICE9PSBjaGFycy4kUkJSQUNFO1xufVxuXG5mdW5jdGlvbiBjb21wYXJlQ2hhckNvZGVDYXNlSW5zZW5zaXRpdmUoY29kZTE6IG51bWJlciwgY29kZTI6IG51bWJlcik6IGJvb2xlYW4ge1xuICByZXR1cm4gdG9VcHBlckNhc2VDaGFyQ29kZShjb2RlMSkgPT09IHRvVXBwZXJDYXNlQ2hhckNvZGUoY29kZTIpO1xufVxuXG5mdW5jdGlvbiB0b1VwcGVyQ2FzZUNoYXJDb2RlKGNvZGU6IG51bWJlcik6IG51bWJlciB7XG4gIHJldHVybiBjb2RlID49IGNoYXJzLiRhICYmIGNvZGUgPD0gY2hhcnMuJHogPyBjb2RlIC0gY2hhcnMuJGEgKyBjaGFycy4kQSA6IGNvZGU7XG59XG5cbmZ1bmN0aW9uIGlzQmxvY2tOYW1lQ2hhcihjb2RlOiBudW1iZXIpOiBib29sZWFuIHtcbiAgcmV0dXJuIGNoYXJzLmlzQXNjaWlMZXR0ZXIoY29kZSkgfHwgY2hhcnMuaXNEaWdpdChjb2RlKSB8fCBjb2RlID09PSBjaGFycy4kXztcbn1cblxuZnVuY3Rpb24gaXNCbG9ja1BhcmFtZXRlckNoYXIoY29kZTogbnVtYmVyKTogYm9vbGVhbiB7XG4gIHJldHVybiBjb2RlICE9PSBjaGFycy4kU0VNSUNPTE9OICYmIGlzTm90V2hpdGVzcGFjZShjb2RlKTtcbn1cblxuZnVuY3Rpb24gbWVyZ2VUZXh0VG9rZW5zKHNyY1Rva2VuczogVG9rZW5bXSk6IFRva2VuW10ge1xuICBjb25zdCBkc3RUb2tlbnM6IFRva2VuW10gPSBbXTtcbiAgbGV0IGxhc3REc3RUb2tlbjogVG9rZW4gfCB1bmRlZmluZWQgPSB1bmRlZmluZWQ7XG4gIGZvciAobGV0IGkgPSAwOyBpIDwgc3JjVG9rZW5zLmxlbmd0aDsgaSsrKSB7XG4gICAgY29uc3QgdG9rZW4gPSBzcmNUb2tlbnNbaV07XG4gICAgaWYgKFxuICAgICAgKGxhc3REc3RUb2tlbiAmJiBsYXN0RHN0VG9rZW4udHlwZSA9PT0gVG9rZW5UeXBlLlRFWFQgJiYgdG9rZW4udHlwZSA9PT0gVG9rZW5UeXBlLlRFWFQpIHx8XG4gICAgICAobGFzdERzdFRva2VuICYmXG4gICAgICAgIGxhc3REc3RUb2tlbi50eXBlID09PSBUb2tlblR5cGUuQVRUUl9WQUxVRV9URVhUICYmXG4gICAgICAgIHRva2VuLnR5cGUgPT09IFRva2VuVHlwZS5BVFRSX1ZBTFVFX1RFWFQpXG4gICAgKSB7XG4gICAgICBsYXN0RHN0VG9rZW4ucGFydHNbMF0hICs9IHRva2VuLnBhcnRzWzBdO1xuICAgICAgbGFzdERzdFRva2VuLnNvdXJjZVNwYW4uZW5kID0gdG9rZW4uc291cmNlU3Bhbi5lbmQ7XG4gICAgfSBlbHNlIHtcbiAgICAgIGxhc3REc3RUb2tlbiA9IHRva2VuO1xuICAgICAgZHN0VG9rZW5zLnB1c2gobGFzdERzdFRva2VuKTtcbiAgICB9XG4gIH1cblxuICByZXR1cm4gZHN0VG9rZW5zO1xufVxuXG4vKipcbiAqIFRoZSBfVG9rZW5pemVyIHVzZXMgb2JqZWN0cyBvZiB0aGlzIHR5cGUgdG8gbW92ZSB0aHJvdWdoIHRoZSBpbnB1dCB0ZXh0LFxuICogZXh0cmFjdGluZyBcInBhcnNlZCBjaGFyYWN0ZXJzXCIuIFRoZXNlIGNvdWxkIGJlIG1vcmUgdGhhbiBvbmUgYWN0dWFsIGNoYXJhY3RlclxuICogaWYgdGhlIHRleHQgY29udGFpbnMgZXNjYXBlIHNlcXVlbmNlcy5cbiAqL1xuaW50ZXJmYWNlIENoYXJhY3RlckN1cnNvciB7XG4gIC8qKiBJbml0aWFsaXplIHRoZSBjdXJzb3IuICovXG4gIGluaXQoKTogdm9pZDtcbiAgLyoqIFRoZSBwYXJzZWQgY2hhcmFjdGVyIGF0IHRoZSBjdXJyZW50IGN1cnNvciBwb3NpdGlvbi4gKi9cbiAgcGVlaygpOiBudW1iZXI7XG4gIC8qKiBBZHZhbmNlIHRoZSBjdXJzb3IgYnkgb25lIHBhcnNlZCBjaGFyYWN0ZXIuICovXG4gIGFkdmFuY2UoKTogdm9pZDtcbiAgLyoqIEdldCBhIHNwYW4gZnJvbSB0aGUgbWFya2VkIHN0YXJ0IHBvaW50IHRvIHRoZSBjdXJyZW50IHBvaW50LiAqL1xuICBnZXRTcGFuKHN0YXJ0PzogdGhpcywgbGVhZGluZ1RyaXZpYUNvZGVQb2ludHM/OiBudW1iZXJbXSk6IFBhcnNlU291cmNlU3BhbjtcbiAgLyoqIEdldCB0aGUgcGFyc2VkIGNoYXJhY3RlcnMgZnJvbSB0aGUgbWFya2VkIHN0YXJ0IHBvaW50IHRvIHRoZSBjdXJyZW50IHBvaW50LiAqL1xuICBnZXRDaGFycyhzdGFydDogdGhpcyk6IHN0cmluZztcbiAgLyoqIFRoZSBudW1iZXIgb2YgY2hhcmFjdGVycyBsZWZ0IGJlZm9yZSB0aGUgZW5kIG9mIHRoZSBjdXJzb3IuICovXG4gIGNoYXJzTGVmdCgpOiBudW1iZXI7XG4gIC8qKiBUaGUgbnVtYmVyIG9mIGNoYXJhY3RlcnMgYmV0d2VlbiBgdGhpc2AgY3Vyc29yIGFuZCBgb3RoZXJgIGN1cnNvci4gKi9cbiAgZGlmZihvdGhlcjogdGhpcyk6IG51bWJlcjtcbiAgLyoqIE1ha2UgYSBjb3B5IG9mIHRoaXMgY3Vyc29yICovXG4gIGNsb25lKCk6IENoYXJhY3RlckN1cnNvcjtcbn1cblxuaW50ZXJmYWNlIEN1cnNvclN0YXRlIHtcbiAgcGVlazogbnVtYmVyO1xuICBvZmZzZXQ6IG51bWJlcjtcbiAgbGluZTogbnVtYmVyO1xuICBjb2x1bW46IG51bWJlcjtcbn1cblxuY2xhc3MgUGxhaW5DaGFyYWN0ZXJDdXJzb3IgaW1wbGVtZW50cyBDaGFyYWN0ZXJDdXJzb3Ige1xuICBwcm90ZWN0ZWQgc3RhdGU6IEN1cnNvclN0YXRlO1xuICBwcm90ZWN0ZWQgZmlsZTogUGFyc2VTb3VyY2VGaWxlO1xuICBwcm90ZWN0ZWQgaW5wdXQ6IHN0cmluZztcbiAgcHJvdGVjdGVkIGVuZDogbnVtYmVyO1xuXG4gIGNvbnN0cnVjdG9yKGZpbGVPckN1cnNvcjogUGxhaW5DaGFyYWN0ZXJDdXJzb3IpO1xuICBjb25zdHJ1Y3RvcihmaWxlT3JDdXJzb3I6IFBhcnNlU291cmNlRmlsZSwgcmFuZ2U6IExleGVyUmFuZ2UpO1xuICBjb25zdHJ1Y3RvcihmaWxlT3JDdXJzb3I6IFBhcnNlU291cmNlRmlsZSB8IFBsYWluQ2hhcmFjdGVyQ3Vyc29yLCByYW5nZT86IExleGVyUmFuZ2UpIHtcbiAgICBpZiAoZmlsZU9yQ3Vyc29yIGluc3RhbmNlb2YgUGxhaW5DaGFyYWN0ZXJDdXJzb3IpIHtcbiAgICAgIHRoaXMuZmlsZSA9IGZpbGVPckN1cnNvci5maWxlO1xuICAgICAgdGhpcy5pbnB1dCA9IGZpbGVPckN1cnNvci5pbnB1dDtcbiAgICAgIHRoaXMuZW5kID0gZmlsZU9yQ3Vyc29yLmVuZDtcblxuICAgICAgY29uc3Qgc3RhdGUgPSBmaWxlT3JDdXJzb3Iuc3RhdGU7XG4gICAgICAvLyBOb3RlOiBhdm9pZCB1c2luZyBgey4uLmZpbGVPckN1cnNvci5zdGF0ZX1gIGhlcmUgYXMgdGhhdCBoYXMgYSBzZXZlcmUgcGVyZm9ybWFuY2UgcGVuYWx0eS5cbiAgICAgIC8vIEluIEVTNSBidW5kbGVzIHRoZSBvYmplY3Qgc3ByZWFkIG9wZXJhdG9yIGlzIHRyYW5zbGF0ZWQgaW50byB0aGUgYF9fYXNzaWduYCBoZWxwZXIsIHdoaWNoXG4gICAgICAvLyBpcyBub3Qgb3B0aW1pemVkIGJ5IFZNcyBhcyBlZmZpY2llbnRseSBhcyBhIHJhdyBvYmplY3QgbGl0ZXJhbC4gU2luY2UgdGhpcyBjb25zdHJ1Y3RvciBpc1xuICAgICAgLy8gY2FsbGVkIGluIHRpZ2h0IGxvb3BzLCB0aGlzIGRpZmZlcmVuY2UgbWF0dGVycy5cbiAgICAgIHRoaXMuc3RhdGUgPSB7XG4gICAgICAgIHBlZWs6IHN0YXRlLnBlZWssXG4gICAgICAgIG9mZnNldDogc3RhdGUub2Zmc2V0LFxuICAgICAgICBsaW5lOiBzdGF0ZS5saW5lLFxuICAgICAgICBjb2x1bW46IHN0YXRlLmNvbHVtbixcbiAgICAgIH07XG4gICAgfSBlbHNlIHtcbiAgICAgIGlmICghcmFuZ2UpIHtcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKFxuICAgICAgICAgICdQcm9ncmFtbWluZyBlcnJvcjogdGhlIHJhbmdlIGFyZ3VtZW50IG11c3QgYmUgcHJvdmlkZWQgd2l0aCBhIGZpbGUgYXJndW1lbnQuJyxcbiAgICAgICAgKTtcbiAgICAgIH1cbiAgICAgIHRoaXMuZmlsZSA9IGZpbGVPckN1cnNvcjtcbiAgICAgIHRoaXMuaW5wdXQgPSBmaWxlT3JDdXJzb3IuY29udGVudDtcbiAgICAgIHRoaXMuZW5kID0gcmFuZ2UuZW5kUG9zO1xuICAgICAgdGhpcy5zdGF0ZSA9IHtcbiAgICAgICAgcGVlazogLTEsXG4gICAgICAgIG9mZnNldDogcmFuZ2Uuc3RhcnRQb3MsXG4gICAgICAgIGxpbmU6IHJhbmdlLnN0YXJ0TGluZSxcbiAgICAgICAgY29sdW1uOiByYW5nZS5zdGFydENvbCxcbiAgICAgIH07XG4gICAgfVxuICB9XG5cbiAgY2xvbmUoKTogUGxhaW5DaGFyYWN0ZXJDdXJzb3Ige1xuICAgIHJldHVybiBuZXcgUGxhaW5DaGFyYWN0ZXJDdXJzb3IodGhpcyk7XG4gIH1cblxuICBwZWVrKCkge1xuICAgIHJldHVybiB0aGlzLnN0YXRlLnBlZWs7XG4gIH1cbiAgY2hhcnNMZWZ0KCkge1xuICAgIHJldHVybiB0aGlzLmVuZCAtIHRoaXMuc3RhdGUub2Zmc2V0O1xuICB9XG4gIGRpZmYob3RoZXI6IHRoaXMpIHtcbiAgICByZXR1cm4gdGhpcy5zdGF0ZS5vZmZzZXQgLSBvdGhlci5zdGF0ZS5vZmZzZXQ7XG4gIH1cblxuICBhZHZhbmNlKCk6IHZvaWQge1xuICAgIHRoaXMuYWR2YW5jZVN0YXRlKHRoaXMuc3RhdGUpO1xuICB9XG5cbiAgaW5pdCgpOiB2b2lkIHtcbiAgICB0aGlzLnVwZGF0ZVBlZWsodGhpcy5zdGF0ZSk7XG4gIH1cblxuICBnZXRTcGFuKHN0YXJ0PzogdGhpcywgbGVhZGluZ1RyaXZpYUNvZGVQb2ludHM/OiBudW1iZXJbXSk6IFBhcnNlU291cmNlU3BhbiB7XG4gICAgc3RhcnQgPSBzdGFydCB8fCB0aGlzO1xuICAgIGxldCBmdWxsU3RhcnQgPSBzdGFydDtcbiAgICBpZiAobGVhZGluZ1RyaXZpYUNvZGVQb2ludHMpIHtcbiAgICAgIHdoaWxlICh0aGlzLmRpZmYoc3RhcnQpID4gMCAmJiBsZWFkaW5nVHJpdmlhQ29kZVBvaW50cy5pbmRleE9mKHN0YXJ0LnBlZWsoKSkgIT09IC0xKSB7XG4gICAgICAgIGlmIChmdWxsU3RhcnQgPT09IHN0YXJ0KSB7XG4gICAgICAgICAgc3RhcnQgPSBzdGFydC5jbG9uZSgpIGFzIHRoaXM7XG4gICAgICAgIH1cbiAgICAgICAgc3RhcnQuYWR2YW5jZSgpO1xuICAgICAgfVxuICAgIH1cbiAgICBjb25zdCBzdGFydExvY2F0aW9uID0gdGhpcy5sb2NhdGlvbkZyb21DdXJzb3Ioc3RhcnQpO1xuICAgIGNvbnN0IGVuZExvY2F0aW9uID0gdGhpcy5sb2NhdGlvbkZyb21DdXJzb3IodGhpcyk7XG4gICAgY29uc3QgZnVsbFN0YXJ0TG9jYXRpb24gPVxuICAgICAgZnVsbFN0YXJ0ICE9PSBzdGFydCA/IHRoaXMubG9jYXRpb25Gcm9tQ3Vyc29yKGZ1bGxTdGFydCkgOiBzdGFydExvY2F0aW9uO1xuICAgIHJldHVybiBuZXcgUGFyc2VTb3VyY2VTcGFuKHN0YXJ0TG9jYXRpb24sIGVuZExvY2F0aW9uLCBmdWxsU3RhcnRMb2NhdGlvbik7XG4gIH1cblxuICBnZXRDaGFycyhzdGFydDogdGhpcyk6IHN0cmluZyB7XG4gICAgcmV0dXJuIHRoaXMuaW5wdXQuc3Vic3RyaW5nKHN0YXJ0LnN0YXRlLm9mZnNldCwgdGhpcy5zdGF0ZS5vZmZzZXQpO1xuICB9XG5cbiAgY2hhckF0KHBvczogbnVtYmVyKTogbnVtYmVyIHtcbiAgICByZXR1cm4gdGhpcy5pbnB1dC5jaGFyQ29kZUF0KHBvcyk7XG4gIH1cblxuICBwcm90ZWN0ZWQgYWR2YW5jZVN0YXRlKHN0YXRlOiBDdXJzb3JTdGF0ZSkge1xuICAgIGlmIChzdGF0ZS5vZmZzZXQgPj0gdGhpcy5lbmQpIHtcbiAgICAgIHRoaXMuc3RhdGUgPSBzdGF0ZTtcbiAgICAgIHRocm93IG5ldyBDdXJzb3JFcnJvcignVW5leHBlY3RlZCBjaGFyYWN0ZXIgXCJFT0ZcIicsIHRoaXMpO1xuICAgIH1cbiAgICBjb25zdCBjdXJyZW50Q2hhciA9IHRoaXMuY2hhckF0KHN0YXRlLm9mZnNldCk7XG4gICAgaWYgKGN1cnJlbnRDaGFyID09PSBjaGFycy4kTEYpIHtcbiAgICAgIHN0YXRlLmxpbmUrKztcbiAgICAgIHN0YXRlLmNvbHVtbiA9IDA7XG4gICAgfSBlbHNlIGlmICghY2hhcnMuaXNOZXdMaW5lKGN1cnJlbnRDaGFyKSkge1xuICAgICAgc3RhdGUuY29sdW1uKys7XG4gICAgfVxuICAgIHN0YXRlLm9mZnNldCsrO1xuICAgIHRoaXMudXBkYXRlUGVlayhzdGF0ZSk7XG4gIH1cblxuICBwcm90ZWN0ZWQgdXBkYXRlUGVlayhzdGF0ZTogQ3Vyc29yU3RhdGUpOiB2b2lkIHtcbiAgICBzdGF0ZS5wZWVrID0gc3RhdGUub2Zmc2V0ID49IHRoaXMuZW5kID8gY2hhcnMuJEVPRiA6IHRoaXMuY2hhckF0KHN0YXRlLm9mZnNldCk7XG4gIH1cblxuICBwcml2YXRlIGxvY2F0aW9uRnJvbUN1cnNvcihjdXJzb3I6IHRoaXMpOiBQYXJzZUxvY2F0aW9uIHtcbiAgICByZXR1cm4gbmV3IFBhcnNlTG9jYXRpb24oXG4gICAgICBjdXJzb3IuZmlsZSxcbiAgICAgIGN1cnNvci5zdGF0ZS5vZmZzZXQsXG4gICAgICBjdXJzb3Iuc3RhdGUubGluZSxcbiAgICAgIGN1cnNvci5zdGF0ZS5jb2x1bW4sXG4gICAgKTtcbiAgfVxufVxuXG5jbGFzcyBFc2NhcGVkQ2hhcmFjdGVyQ3Vyc29yIGV4dGVuZHMgUGxhaW5DaGFyYWN0ZXJDdXJzb3Ige1xuICBwcm90ZWN0ZWQgaW50ZXJuYWxTdGF0ZTogQ3Vyc29yU3RhdGU7XG5cbiAgY29uc3RydWN0b3IoZmlsZU9yQ3Vyc29yOiBFc2NhcGVkQ2hhcmFjdGVyQ3Vyc29yKTtcbiAgY29uc3RydWN0b3IoZmlsZU9yQ3Vyc29yOiBQYXJzZVNvdXJjZUZpbGUsIHJhbmdlOiBMZXhlclJhbmdlKTtcbiAgY29uc3RydWN0b3IoZmlsZU9yQ3Vyc29yOiBQYXJzZVNvdXJjZUZpbGUgfCBFc2NhcGVkQ2hhcmFjdGVyQ3Vyc29yLCByYW5nZT86IExleGVyUmFuZ2UpIHtcbiAgICBpZiAoZmlsZU9yQ3Vyc29yIGluc3RhbmNlb2YgRXNjYXBlZENoYXJhY3RlckN1cnNvcikge1xuICAgICAgc3VwZXIoZmlsZU9yQ3Vyc29yKTtcbiAgICAgIHRoaXMuaW50ZXJuYWxTdGF0ZSA9IHsuLi5maWxlT3JDdXJzb3IuaW50ZXJuYWxTdGF0ZX07XG4gICAgfSBlbHNlIHtcbiAgICAgIHN1cGVyKGZpbGVPckN1cnNvciwgcmFuZ2UhKTtcbiAgICAgIHRoaXMuaW50ZXJuYWxTdGF0ZSA9IHRoaXMuc3RhdGU7XG4gICAgfVxuICB9XG5cbiAgb3ZlcnJpZGUgYWR2YW5jZSgpOiB2b2lkIHtcbiAgICB0aGlzLnN0YXRlID0gdGhpcy5pbnRlcm5hbFN0YXRlO1xuICAgIHN1cGVyLmFkdmFuY2UoKTtcbiAgICB0aGlzLnByb2Nlc3NFc2NhcGVTZXF1ZW5jZSgpO1xuICB9XG5cbiAgb3ZlcnJpZGUgaW5pdCgpOiB2b2lkIHtcbiAgICBzdXBlci5pbml0KCk7XG4gICAgdGhpcy5wcm9jZXNzRXNjYXBlU2VxdWVuY2UoKTtcbiAgfVxuXG4gIG92ZXJyaWRlIGNsb25lKCk6IEVzY2FwZWRDaGFyYWN0ZXJDdXJzb3Ige1xuICAgIHJldHVybiBuZXcgRXNjYXBlZENoYXJhY3RlckN1cnNvcih0aGlzKTtcbiAgfVxuXG4gIG92ZXJyaWRlIGdldENoYXJzKHN0YXJ0OiB0aGlzKTogc3RyaW5nIHtcbiAgICBjb25zdCBjdXJzb3IgPSBzdGFydC5jbG9uZSgpO1xuICAgIGxldCBjaGFycyA9ICcnO1xuICAgIHdoaWxlIChjdXJzb3IuaW50ZXJuYWxTdGF0ZS5vZmZzZXQgPCB0aGlzLmludGVybmFsU3RhdGUub2Zmc2V0KSB7XG4gICAgICBjaGFycyArPSBTdHJpbmcuZnJvbUNvZGVQb2ludChjdXJzb3IucGVlaygpKTtcbiAgICAgIGN1cnNvci5hZHZhbmNlKCk7XG4gICAgfVxuICAgIHJldHVybiBjaGFycztcbiAgfVxuXG4gIC8qKlxuICAgKiBQcm9jZXNzIHRoZSBlc2NhcGUgc2VxdWVuY2UgdGhhdCBzdGFydHMgYXQgdGhlIGN1cnJlbnQgcG9zaXRpb24gaW4gdGhlIHRleHQuXG4gICAqXG4gICAqIFRoaXMgbWV0aG9kIGlzIGNhbGxlZCB0byBlbnN1cmUgdGhhdCBgcGVla2AgaGFzIHRoZSB1bmVzY2FwZWQgdmFsdWUgb2YgZXNjYXBlIHNlcXVlbmNlcy5cbiAgICovXG4gIHByb3RlY3RlZCBwcm9jZXNzRXNjYXBlU2VxdWVuY2UoKTogdm9pZCB7XG4gICAgY29uc3QgcGVlayA9ICgpID0+IHRoaXMuaW50ZXJuYWxTdGF0ZS5wZWVrO1xuXG4gICAgaWYgKHBlZWsoKSA9PT0gY2hhcnMuJEJBQ0tTTEFTSCkge1xuICAgICAgLy8gV2UgaGF2ZSBoaXQgYW4gZXNjYXBlIHNlcXVlbmNlIHNvIHdlIG5lZWQgdGhlIGludGVybmFsIHN0YXRlIHRvIGJlY29tZSBpbmRlcGVuZGVudFxuICAgICAgLy8gb2YgdGhlIGV4dGVybmFsIHN0YXRlLlxuICAgICAgdGhpcy5pbnRlcm5hbFN0YXRlID0gey4uLnRoaXMuc3RhdGV9O1xuXG4gICAgICAvLyBNb3ZlIHBhc3QgdGhlIGJhY2tzbGFzaFxuICAgICAgdGhpcy5hZHZhbmNlU3RhdGUodGhpcy5pbnRlcm5hbFN0YXRlKTtcblxuICAgICAgLy8gRmlyc3QgY2hlY2sgZm9yIHN0YW5kYXJkIGNvbnRyb2wgY2hhciBzZXF1ZW5jZXNcbiAgICAgIGlmIChwZWVrKCkgPT09IGNoYXJzLiRuKSB7XG4gICAgICAgIHRoaXMuc3RhdGUucGVlayA9IGNoYXJzLiRMRjtcbiAgICAgIH0gZWxzZSBpZiAocGVlaygpID09PSBjaGFycy4kcikge1xuICAgICAgICB0aGlzLnN0YXRlLnBlZWsgPSBjaGFycy4kQ1I7XG4gICAgICB9IGVsc2UgaWYgKHBlZWsoKSA9PT0gY2hhcnMuJHYpIHtcbiAgICAgICAgdGhpcy5zdGF0ZS5wZWVrID0gY2hhcnMuJFZUQUI7XG4gICAgICB9IGVsc2UgaWYgKHBlZWsoKSA9PT0gY2hhcnMuJHQpIHtcbiAgICAgICAgdGhpcy5zdGF0ZS5wZWVrID0gY2hhcnMuJFRBQjtcbiAgICAgIH0gZWxzZSBpZiAocGVlaygpID09PSBjaGFycy4kYikge1xuICAgICAgICB0aGlzLnN0YXRlLnBlZWsgPSBjaGFycy4kQlNQQUNFO1xuICAgICAgfSBlbHNlIGlmIChwZWVrKCkgPT09IGNoYXJzLiRmKSB7XG4gICAgICAgIHRoaXMuc3RhdGUucGVlayA9IGNoYXJzLiRGRjtcbiAgICAgIH1cblxuICAgICAgLy8gTm93IGNvbnNpZGVyIG1vcmUgY29tcGxleCBzZXF1ZW5jZXNcbiAgICAgIGVsc2UgaWYgKHBlZWsoKSA9PT0gY2hhcnMuJHUpIHtcbiAgICAgICAgLy8gVW5pY29kZSBjb2RlLXBvaW50IHNlcXVlbmNlXG4gICAgICAgIHRoaXMuYWR2YW5jZVN0YXRlKHRoaXMuaW50ZXJuYWxTdGF0ZSk7IC8vIGFkdmFuY2UgcGFzdCB0aGUgYHVgIGNoYXJcbiAgICAgICAgaWYgKHBlZWsoKSA9PT0gY2hhcnMuJExCUkFDRSkge1xuICAgICAgICAgIC8vIFZhcmlhYmxlIGxlbmd0aCBVbmljb2RlLCBlLmcuIGBcXHh7MTIzfWBcbiAgICAgICAgICB0aGlzLmFkdmFuY2VTdGF0ZSh0aGlzLmludGVybmFsU3RhdGUpOyAvLyBhZHZhbmNlIHBhc3QgdGhlIGB7YCBjaGFyXG4gICAgICAgICAgLy8gQWR2YW5jZSBwYXN0IHRoZSB2YXJpYWJsZSBudW1iZXIgb2YgaGV4IGRpZ2l0cyB1bnRpbCB3ZSBoaXQgYSBgfWAgY2hhclxuICAgICAgICAgIGNvbnN0IGRpZ2l0U3RhcnQgPSB0aGlzLmNsb25lKCk7XG4gICAgICAgICAgbGV0IGxlbmd0aCA9IDA7XG4gICAgICAgICAgd2hpbGUgKHBlZWsoKSAhPT0gY2hhcnMuJFJCUkFDRSkge1xuICAgICAgICAgICAgdGhpcy5hZHZhbmNlU3RhdGUodGhpcy5pbnRlcm5hbFN0YXRlKTtcbiAgICAgICAgICAgIGxlbmd0aCsrO1xuICAgICAgICAgIH1cbiAgICAgICAgICB0aGlzLnN0YXRlLnBlZWsgPSB0aGlzLmRlY29kZUhleERpZ2l0cyhkaWdpdFN0YXJ0LCBsZW5ndGgpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgIC8vIEZpeGVkIGxlbmd0aCBVbmljb2RlLCBlLmcuIGBcXHUxMjM0YFxuICAgICAgICAgIGNvbnN0IGRpZ2l0U3RhcnQgPSB0aGlzLmNsb25lKCk7XG4gICAgICAgICAgdGhpcy5hZHZhbmNlU3RhdGUodGhpcy5pbnRlcm5hbFN0YXRlKTtcbiAgICAgICAgICB0aGlzLmFkdmFuY2VTdGF0ZSh0aGlzLmludGVybmFsU3RhdGUpO1xuICAgICAgICAgIHRoaXMuYWR2YW5jZVN0YXRlKHRoaXMuaW50ZXJuYWxTdGF0ZSk7XG4gICAgICAgICAgdGhpcy5zdGF0ZS5wZWVrID0gdGhpcy5kZWNvZGVIZXhEaWdpdHMoZGlnaXRTdGFydCwgNCk7XG4gICAgICAgIH1cbiAgICAgIH0gZWxzZSBpZiAocGVlaygpID09PSBjaGFycy4keCkge1xuICAgICAgICAvLyBIZXggY2hhciBjb2RlLCBlLmcuIGBcXHgyRmBcbiAgICAgICAgdGhpcy5hZHZhbmNlU3RhdGUodGhpcy5pbnRlcm5hbFN0YXRlKTsgLy8gYWR2YW5jZSBwYXN0IHRoZSBgeGAgY2hhclxuICAgICAgICBjb25zdCBkaWdpdFN0YXJ0ID0gdGhpcy5jbG9uZSgpO1xuICAgICAgICB0aGlzLmFkdmFuY2VTdGF0ZSh0aGlzLmludGVybmFsU3RhdGUpO1xuICAgICAgICB0aGlzLnN0YXRlLnBlZWsgPSB0aGlzLmRlY29kZUhleERpZ2l0cyhkaWdpdFN0YXJ0LCAyKTtcbiAgICAgIH0gZWxzZSBpZiAoY2hhcnMuaXNPY3RhbERpZ2l0KHBlZWsoKSkpIHtcbiAgICAgICAgLy8gT2N0YWwgY2hhciBjb2RlLCBlLmcuIGBcXDAxMmAsXG4gICAgICAgIGxldCBvY3RhbCA9ICcnO1xuICAgICAgICBsZXQgbGVuZ3RoID0gMDtcbiAgICAgICAgbGV0IHByZXZpb3VzID0gdGhpcy5jbG9uZSgpO1xuICAgICAgICB3aGlsZSAoY2hhcnMuaXNPY3RhbERpZ2l0KHBlZWsoKSkgJiYgbGVuZ3RoIDwgMykge1xuICAgICAgICAgIHByZXZpb3VzID0gdGhpcy5jbG9uZSgpO1xuICAgICAgICAgIG9jdGFsICs9IFN0cmluZy5mcm9tQ29kZVBvaW50KHBlZWsoKSk7XG4gICAgICAgICAgdGhpcy5hZHZhbmNlU3RhdGUodGhpcy5pbnRlcm5hbFN0YXRlKTtcbiAgICAgICAgICBsZW5ndGgrKztcbiAgICAgICAgfVxuICAgICAgICB0aGlzLnN0YXRlLnBlZWsgPSBwYXJzZUludChvY3RhbCwgOCk7XG4gICAgICAgIC8vIEJhY2t1cCBvbmUgY2hhclxuICAgICAgICB0aGlzLmludGVybmFsU3RhdGUgPSBwcmV2aW91cy5pbnRlcm5hbFN0YXRlO1xuICAgICAgfSBlbHNlIGlmIChjaGFycy5pc05ld0xpbmUodGhpcy5pbnRlcm5hbFN0YXRlLnBlZWspKSB7XG4gICAgICAgIC8vIExpbmUgY29udGludWF0aW9uIGBcXGAgZm9sbG93ZWQgYnkgYSBuZXcgbGluZVxuICAgICAgICB0aGlzLmFkdmFuY2VTdGF0ZSh0aGlzLmludGVybmFsU3RhdGUpOyAvLyBhZHZhbmNlIG92ZXIgdGhlIG5ld2xpbmVcbiAgICAgICAgdGhpcy5zdGF0ZSA9IHRoaXMuaW50ZXJuYWxTdGF0ZTtcbiAgICAgIH0gZWxzZSB7XG4gICAgICAgIC8vIElmIG5vbmUgb2YgdGhlIGBpZmAgYmxvY2tzIHdlcmUgZXhlY3V0ZWQgdGhlbiB3ZSBqdXN0IGhhdmUgYW4gZXNjYXBlZCBub3JtYWwgY2hhcmFjdGVyLlxuICAgICAgICAvLyBJbiB0aGF0IGNhc2Ugd2UganVzdCwgZWZmZWN0aXZlbHksIHNraXAgdGhlIGJhY2tzbGFzaCBmcm9tIHRoZSBjaGFyYWN0ZXIuXG4gICAgICAgIHRoaXMuc3RhdGUucGVlayA9IHRoaXMuaW50ZXJuYWxTdGF0ZS5wZWVrO1xuICAgICAgfVxuICAgIH1cbiAgfVxuXG4gIHByb3RlY3RlZCBkZWNvZGVIZXhEaWdpdHMoc3RhcnQ6IEVzY2FwZWRDaGFyYWN0ZXJDdXJzb3IsIGxlbmd0aDogbnVtYmVyKTogbnVtYmVyIHtcbiAgICBjb25zdCBoZXggPSB0aGlzLmlucHV0LnNsaWNlKHN0YXJ0LmludGVybmFsU3RhdGUub2Zmc2V0LCBzdGFydC5pbnRlcm5hbFN0YXRlLm9mZnNldCArIGxlbmd0aCk7XG4gICAgY29uc3QgY2hhckNvZGUgPSBwYXJzZUludChoZXgsIDE2KTtcbiAgICBpZiAoIWlzTmFOKGNoYXJDb2RlKSkge1xuICAgICAgcmV0dXJuIGNoYXJDb2RlO1xuICAgIH0gZWxzZSB7XG4gICAgICBzdGFydC5zdGF0ZSA9IHN0YXJ0LmludGVybmFsU3RhdGU7XG4gICAgICB0aHJvdyBuZXcgQ3Vyc29yRXJyb3IoJ0ludmFsaWQgaGV4YWRlY2ltYWwgZXNjYXBlIHNlcXVlbmNlJywgc3RhcnQpO1xuICAgIH1cbiAgfVxufVxuXG5leHBvcnQgY2xhc3MgQ3Vyc29yRXJyb3Ige1xuICBjb25zdHJ1Y3RvcihcbiAgICBwdWJsaWMgbXNnOiBzdHJpbmcsXG4gICAgcHVibGljIGN1cnNvcjogQ2hhcmFjdGVyQ3Vyc29yLFxuICApIHt9XG59XG4iXX0=