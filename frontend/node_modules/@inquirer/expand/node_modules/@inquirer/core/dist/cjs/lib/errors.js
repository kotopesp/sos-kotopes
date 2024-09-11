"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ValidationError = exports.HookError = exports.ExitPromptError = exports.CancelPromptError = void 0;
class CancelPromptError extends Error {
    constructor() {
        super(...arguments);
        Object.defineProperty(this, "message", {
            enumerable: true,
            configurable: true,
            writable: true,
            value: 'Prompt was canceled'
        });
    }
}
exports.CancelPromptError = CancelPromptError;
class ExitPromptError extends Error {
}
exports.ExitPromptError = ExitPromptError;
class HookError extends Error {
}
exports.HookError = HookError;
class ValidationError extends Error {
}
exports.ValidationError = ValidationError;
