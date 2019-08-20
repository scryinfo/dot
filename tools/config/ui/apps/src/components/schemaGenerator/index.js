'use strict';

import {Compiler} from './compiler';
import {AST} from './ast.js';


var jsonToSchema = function(json) {
  var compiler = new Compiler();
  var ast = new AST();
  ast.build(json);
  console.log(JSON.stringify(ast.tree, null, 4));
  compiler.compile(ast.tree);
  console.log(JSON.stringify(compiler.schema, null, 4));
  return compiler.schema;
};

export {jsonToSchema};
