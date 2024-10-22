"use strict";
/**
 * @license
 * Copyright Google LLC All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.dev/license
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.createAngularIndexHtmlMiddleware = createAngularIndexHtmlMiddleware;
const utils_1 = require("../utils");
function createAngularIndexHtmlMiddleware(server, outputFiles, indexHtmlTransformer) {
    return function (req, res, next) {
        if (!req.url) {
            next();
            return;
        }
        // Parse the incoming request.
        // The base of the URL is unused but required to parse the URL.
        const pathname = (0, utils_1.pathnameWithoutBasePath)(req.url, server.config.base);
        if (pathname !== '/' && pathname !== '/index.html') {
            next();
            return;
        }
        const rawHtml = outputFiles.get('/index.html')?.contents;
        if (!rawHtml) {
            next();
            return;
        }
        server
            .transformIndexHtml(req.url, Buffer.from(rawHtml).toString('utf-8'))
            .then(async (processedHtml) => {
            if (indexHtmlTransformer) {
                processedHtml = await indexHtmlTransformer(processedHtml);
            }
            res.setHeader('Content-Type', 'text/html');
            res.setHeader('Cache-Control', 'no-cache');
            (0, utils_1.appendServerConfiguredHeaders)(server, res);
            res.end(processedHtml);
        })
            .catch((error) => next(error));
    };
}
