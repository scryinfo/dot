'use strict';

/**
 Generates an Abstract Syntax Tree
 used for generating the schema.

 @see: https://en.wikipedia.org/wiki/Abstract_syntax_tree
 */
import {isObject, isArray, getType} from './utils';

/**
 Abstract Syntax Tree Class

 @class AST
 @return {Object}
 */
var AST = function () {
    if (!this instanceof AST) {
        return new AST();
    }

    this.tree = {};
};
/**
 Inspect primitatives and apply the correct type
 and mark as required if the element contains a value.

 @method buildPrimitive
 @param {Object} tree Schema which represents the node
 @param {Node} node The JSON node being inspected
 @return void
 */
AST.prototype.buildPrimitive = function (tree, node) {
    tree.type = getType(node);
    if (tree.type === 'string') {
        tree.minLength = (node.length > 0) ? 1 : 0;
    }

    if (node !== null && typeof node !== 'undefined') {
        tree.required = true;
    }
}

/**
 Inspect object properties and apply the correct
 type and mark as required if the element has set
 properties.

 @method buildObject
 @param {Object} tree Schema which represents the node
 @param {Node} node The JSON node being inspected
 */
AST.prototype.buildObjectTree = function (tree, node) {
    tree.type = tree.type || 'object';
    tree.children = tree.children || {};
    for (var i in node) {
        if (isObject(node[i])) {
            tree.children[i] = {};
            this.buildObjectTree(tree.children[i], node[i]);
            continue;
        } else if (isArray(node[i])) {
            tree.children[i] = {};
            this.buildArrayTree(tree.children[i], node[i]);
        } else {
            tree.children[i] = {};
            this.buildPrimitive(tree.children[i], node[i]);
        }
    }
};

/**
 Inspect array elements apply the correct
 type and mark as required if the element has
 set properties.

 @method buildObject
 @param {Object} tree Schema which represents the node
 @param {Node} node The JSON node being inspected
 */
AST.prototype.buildArrayTree = function (tree, node) {
    tree.type = 'array';
    tree.children = {};
    for (var i = 0; i < node.length; i++) {
        if (isObject(node[i])) {
            tree.children[i] = {};
            tree.children[i].type = 'object';
            var keys = Object.keys(node[i]);
            if (keys.length > 0) {
                tree.children[i].required = true;
            }
            this.buildObjectTree(tree.children[i], node[i]);
        } else if (isArray(node[i])) {
            tree.children[i] = {};
            tree.children[i].type = 'array';
            tree.children[i].uniqueItems = false;
            if (node[i].length > 0) {
                tree.children[i].required = true;
            }
            this.buildArrayTree(tree.children[i], node[i]);
        } else if (tree.type === 'object') {
            tree.children[i] = {};
            this.buildPrimitive(tree.children[i], node[i]);
        }
    }
};

/**
 Initiates generating the AST from the
 given JSON document.

 @param {Object} json JSON object
 @return void
 */
AST.prototype.build = function (json) {
    if (json instanceof Array) {
        this.buildArrayTree(this.tree, json);
    } else {
        this.buildObjectTree(this.tree, json);
    }
};

export {AST};
