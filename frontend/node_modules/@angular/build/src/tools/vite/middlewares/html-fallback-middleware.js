"use strict";
/**
 * @license
 * Copyright Google LLC All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.dev/license
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.angularHtmlFallbackMiddleware = angularHtmlFallbackMiddleware;
const utils_1 = require("../utils");
function angularHtmlFallbackMiddleware(req, res, next) {
    // Similar to how it is handled in vite
    // https://github.com/vitejs/vite/blob/main/packages/vite/src/node/server/middlewares/htmlFallback.ts#L15C19-L15C45
    if ((req.method === 'GET' || req.method === 'HEAD') &&
        (!req.url || !(0, utils_1.lookupMimeTypeFromRequest)(req.url)) &&
        (!req.headers.accept ||
            req.headers.accept.includes('text/html') ||
            req.headers.accept.includes('text/*') ||
            req.headers.accept.includes('*/*'))) {
        req.url = '/index.html';
    }
    next();
}
