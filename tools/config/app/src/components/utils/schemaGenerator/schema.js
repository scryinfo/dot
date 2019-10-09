'use strict';

import {Compiler} from './compiler';
import {AST} from './ast.js';


var jsonToSchema = function (json) {
    var compiler = new Compiler();
    var ast = new AST();
    ast.build(json);
    compiler.compile(ast.tree);
    return compiler.schema;
};

export {jsonToSchema};
