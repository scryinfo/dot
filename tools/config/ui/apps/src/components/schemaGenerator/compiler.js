'use strict';

import {isArray} from './utils';

/**
Schema generator using a AST
tree.

@class Compiler
*/
var Compiler = function() {
  if (!this instanceof Compiler) {
    return new Compiler();
  }

  this.schema = {};
};

/**
Generates a JSON schema based on the provided AST tree.

@method generate
@param {Object} tree AST
@param {Object} schema The schema object
@param {Object} parent Schema node parent object
@return void
*/
Compiler.prototype.generate = function(tree, schema, parent) {
  if (tree.type === 'object'){
    parent.maxProperties = Object.keys(tree.children).length;
  }
  for (var i in tree.children) {
    var child = tree.children[i];
    if (child.type === 'object') {
      if (isArray(parent.required)) {
        parent.required.push(i);
      }
      if(tree.type === 'array'){
        schema = {
          type: 'object'
          ,properties: {}
          ,required: []
        };
        this.generate(child, schema.properties, schema);
        break;
      }else {
        schema[i] = {
          type: 'object'
          ,properties: {}
          ,required: []
        };
        this.generate(child, schema[i].properties, schema[i]);
      }
    } else if (child.type === 'array') {
      if (isArray(parent.required)) {
        parent.required.push(i);
      };
      schema[i] = {
        type: 'array'
        ,uniqueItems: child.uniqueItems
        ,minItems: 1
        ,items: {}
      }
      this.generate(child, schema[i].items, schema[i]);
    } else {
      if(tree.type === 'array'){
        if (child.type) {
          schema.type = child.type;
        };
        break;
      }else{
        schema[i] = {};
        if (child.type) {
          schema[i].type = child.type;
        }

        if (child.minLength) {
          schema[i].minLength = child.minLength;
        }

        if (child.required) {
          if (parent.items && isArray(parent.items.required)) {
            parent.items.required.push(i);
          } else {
            parent.required.push(i);
          }
        }
      }
    }
  }
};

/**
Initates compiling the given AST into a
JSON schema.

@method compile
@param {Object} tree AST object
@return void
*/
Compiler.prototype.compile = function(tree) {
  if (tree.type === 'object') {
   this.schema = {
    '$schema': 'http://json-schema.org/draft-04/schema#'
    ,description: ''
    ,type: 'object'
    ,properties: {}
    ,required: []
   };
   this.generate(tree, this.schema.properties, this.schema);
  } else {
    this.schema = {
      type: 'array'
      ,'$schema': 'http://json-schema.org/draft-04/schema#'
      ,'description': ''
      ,minItems: 1
      ,uniqueItems: false
      ,items: {}
    };
    this.generate(tree, this.schema.items.properties, this.schema.items);
  }
};

export {Compiler};
